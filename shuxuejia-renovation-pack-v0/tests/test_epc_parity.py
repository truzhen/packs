#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""墅学家 EPC 可执行投影的静态 parity 门。

本测试只审计 Pack 声明：不连接 Provider，不产生项目、合同、付款、通知或文件副作用。
运行：python3 -m unittest discover -s shuxuejia-renovation-pack-v0/tests -v
"""

import json
import pathlib
import re
import unittest
from collections import Counter, defaultdict, deque


PACK = pathlib.Path(__file__).resolve().parents[1]
FLOW_PATH = PACK / "flows" / "shuxuejia-epc-executable-projection.flow.json"
MANIFEST_PATH = PACK / "manifest.json"

EXPECTED_STAGE_COUNTS = {
    "S1": 61,
    "S2": 78,
    "S3": 61,
    "S4": 86,
    "S5": 66,
    "S6": 39,
    "S7": 52,
    "S8": 14,
}
EXPECTED_GATE = "shuxuejia_owner_base_gate"
STAGE_ID = re.compile(r"^(S[1-8])_\d{3}$")


def load_json(path):
    with path.open(encoding="utf-8") as stream:
        return json.load(stream)


class EpcParityTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.flow = load_json(FLOW_PATH)
        cls.manifest = load_json(MANIFEST_PATH)
        cls.nodes = cls.flow["nodes"]
        cls.edges = cls.flow["edges"]
        cls.node_ids = {node["id"] for node in cls.nodes}

    def test_version_and_source_are_traceable(self):
        self.assertEqual("shuxuejia-epc-executable-projection", self.flow["flow_id"])
        self.assertEqual(self.manifest["version"], self.flow["version"])
        self.assertEqual("candidate_only_governed_runtime", self.flow["execution_mode"])
        self.assertEqual("reference_flow_with_executable_projection", self.flow["graph_kind"])
        self.assertEqual(
            "3c161aff4e2cc9b6ee896d0ef1ec4b37aaf4b062",
            self.flow["source_commit"],
        )
        self.assertEqual("docs/epc-reference-topology.md", self.flow["source_topology"])

    def test_457_unique_nodes_and_eight_stage_coverage(self):
        ids = [node["id"] for node in self.nodes]
        self.assertEqual(457, len(ids))
        self.assertEqual(457, len(set(ids)), "节点 ID 必须唯一")
        self.assertTrue(all(STAGE_ID.fullmatch(node_id) for node_id in ids))
        stages = Counter(STAGE_ID.fullmatch(node_id).group(1) for node_id in ids)
        self.assertEqual(EXPECTED_STAGE_COUNTS, dict(stages))

    def test_543_edges_are_unique_and_reference_existing_nodes(self):
        edge_ids = [edge["id"] for edge in self.edges]
        pairs = [(edge["source"], edge["target"]) for edge in self.edges]
        self.assertEqual(543, len(edge_ids))
        self.assertEqual(543, len(set(edge_ids)), "边 ID 必须唯一")
        self.assertEqual(543, len(set(pairs)), "不得重复声明同一跳转")
        for edge in self.edges:
            self.assertIn(edge["source"], self.node_ids, edge)
            self.assertIn(edge["target"], self.node_ids, edge)
            self.assertNotEqual(edge["source"], edge["target"], edge)

    def test_every_node_is_topologically_reachable_from_a_start_set(self):
        """投影允许多个 EPC 分支入口，但不允许环或没有入口的孤立子图。"""
        adjacency = defaultdict(list)
        incoming = Counter()
        for edge in self.edges:
            adjacency[edge["source"]].append(edge["target"])
            incoming[edge["target"]] += 1
        starts = {node_id for node_id in self.node_ids if not incoming[node_id]}
        self.assertIn("S1_001", starts)
        self.assertTrue(starts, "流程必须有入口")
        seen = set(starts)
        queue = deque(starts)
        while queue:
            current = queue.popleft()
            for nxt in adjacency[current]:
                if nxt not in seen:
                    seen.add(nxt)
                    queue.append(nxt)
        self.assertEqual(self.node_ids, seen, "不得存在不可从任一 EPC 入口到达的节点")

        indegree = {node_id: incoming[node_id] for node_id in self.node_ids}
        queue = deque(starts)
        visited = 0
        while queue:
            current = queue.popleft()
            visited += 1
            for nxt in adjacency[current]:
                indegree[nxt] -= 1
                if indegree[nxt] == 0:
                    queue.append(nxt)
        self.assertEqual(457, visited, "可执行投影不得包含有向环")

    def test_candidate_gate_and_evidence_requirements_are_preserved(self):
        allowed = {
            "TaskCandidate",
            "AdviceCandidate",
            "MaterialCoordinationCandidate",
            "AcceptanceEvidenceCandidate",
            "SovereignDecisionCandidate",
            "AfterSalesTicketCandidate",
        }
        candidate_types = [node.get("candidate_type") for node in self.nodes]
        self.assertTrue(set(candidate_types) <= allowed)
        self.assertEqual(77, candidate_types.count("AcceptanceEvidenceCandidate"))

        gate_nodes = [node for node in self.nodes if node.get("gate_policy")]
        self.assertEqual(96, len(gate_nodes))
        for node in gate_nodes:
            policy = node["gate_policy"]
            self.assertEqual(EXPECTED_GATE, policy.get("required_gate"), node)
            self.assertTrue(policy.get("pending_owner_confirmation"), node)
            self.assertIn(node["candidate_type"], allowed, node)

        self.assertEqual(EXPECTED_GATE, self.manifest["gates"]["project_gate"])
        self.assertTrue(self.manifest["security_profile"]["candidate_only"])
        self.assertTrue(self.manifest["security_profile"]["requires_base_gate"])
        self.assertTrue(self.manifest["receipt_policy"]["append_only"])

    def test_all_declared_role_slots_are_used_by_collaboration_nodes(self):
        declared = {slot["slot_id"] for slot in self.manifest["role_slots"]}
        referenced = {node["slot_ref"] for node in self.nodes if node.get("slot_ref")}
        self.assertTrue(referenced <= declared, referenced - declared)
        self.assertEqual(declared, referenced)


if __name__ == "__main__":
    unittest.main()
