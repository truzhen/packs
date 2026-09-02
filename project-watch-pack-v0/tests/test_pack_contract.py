#!/usr/bin/env python3
"""项目关注 Pack 的离线契约防线；不调用 OS、不连 Provider、不读任何真实业务数据。"""

import json
import pathlib
import re
import unittest


PACK = pathlib.Path(__file__).resolve().parents[1]

# AGENTS.md §4.2 的 8 个候选基类；Pack 只能声明这 8 种，领域名只能作 domain_candidate_type。
BASE_CANDIDATE_TYPES = {
    "TaskCandidate",
    "MemoryRequestCandidate",
    "CommunicationDraftCandidate",
    "ExecutionIntentCandidate",
    "BusinessObjectCandidate",
    "SceneFlowRunCandidate",
    "CapabilityInvocationCandidate",
    "PackCandidate",
}

# manifest 与 capabilities 必须逐条同源的字段（binding / description 各自私有）。
ALIGNED_REQUIREMENT_FIELDS = (
    "capability",
    "gateway_class",
    "risk_class",
    "fallback_policy",
    "provider_family",
    "execution_level",
    "runtime_requirement",
)

# 阈值判定只许在 truzhenos 05 读模型；flow 内出现任一阈值字段即视为 Pack 层写判定。
THRESHOLD_FIELD_NAMES = (
    "delay_days",
    "stall_days",
    "payment_overdue_days",
    "payment_critical_days",
    "max_snapshot_age_hours",
)


def load(relative):
    return json.loads((PACK / relative).read_text(encoding="utf-8"))


class ProjectWatchPackContractTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = load("manifest.json")
        cls.flow = load("flows/project-watch.flow.json")
        cls.capabilities = load("capabilities/capabilities.json")
        cls.role_slots = load("role-slots/role-slots.json")
        cls.scopes = load("knowledge/knowledge-scopes.json")
        cls.index = load("knowledge/knowledge-index.json")
        cls.nodes = {node["id"]: node for node in cls.flow["nodes"]}
        cls.edges = {(edge["source"], edge["target"]): edge for edge in cls.flow["edges"]}

    # ---- 候选类型 ----------------------------------------------------------
    def test_candidate_types_stay_inside_eight_base_classes(self):
        for node_id, node in self.nodes.items():
            candidate_type = node.get("candidate_type")
            if candidate_type is None:
                continue
            self.assertIn(candidate_type, BASE_CANDIDATE_TYPES,
                          "节点 %s 的 candidate_type 越出 8 基类" % node_id)

    def test_domain_candidate_names_map_to_base_classes(self):
        expected = {
            "snapshot_refresh": ("CapabilityInvocationCandidate", None),
            "anomaly_scan": ("CapabilityInvocationCandidate", None),
            "progress_watch": ("TaskCandidate", "ProjectWatchTaskCandidate"),
            "material_watch": ("BusinessObjectCandidate", "ProjectMaterialWatchCandidate"),
            "payment_watch": ("TaskCandidate", "ProjectPaymentWatchCandidate"),
            "anomaly_item_candidate": ("BusinessObjectCandidate", "ProjectAnomalyCandidate"),
            "challenger_review": ("BusinessObjectCandidate", "AcceptanceEvidenceCandidate"),
            "owner_notice_draft": ("CommunicationDraftCandidate", None),
            "followup_intent": ("ExecutionIntentCandidate", None),
            "gateway_execution": ("CapabilityInvocationCandidate", None),
        }
        for node_id, (base, domain) in expected.items():
            node = self.nodes[node_id]
            self.assertEqual(node.get("candidate_type"), base, node_id)
            self.assertEqual(node.get("domain_candidate_type"), domain, node_id)

    def test_anomaly_item_is_idempotent_on_anomaly_key(self):
        self.assertEqual(self.nodes["anomaly_item_candidate"].get("idempotency_key_source"), "anomaly_key")

    def test_material_watch_is_honest_backlog(self):
        node = self.nodes["material_watch"]
        self.assertTrue(node.get("backlog"))
        self.assertTrue(str(node.get("backlog_reason", "")).strip())

    # ---- 阈值只引用不判定 --------------------------------------------------
    def test_flow_only_references_threshold_policy(self):
        self.assertEqual(self.nodes["anomaly_scan"].get("policy_ref"), "project_watch_policy://default")
        self.assertTrue(self.flow["threshold_policy"]["pack_declares_no_conditions"])
        raw = (PACK / "flows/project-watch.flow.json").read_text(encoding="utf-8")
        for field in THRESHOLD_FIELD_NAMES:
            self.assertNotIn('"%s"' % field, raw, "flow 内不得出现阈值字段 %s" % field)

    # ---- 主链形状 ----------------------------------------------------------
    def test_main_chain_goes_candidate_gate_gateway_receipt(self):
        self.assertEqual(self.edges[("owner_gate", "gateway_execution")].get("condition"), "approved")
        self.assertIn(("gateway_execution", "project_watch_receipt"), self.edges)
        self.assertIn(("followup_intent", "owner_gate"), self.edges)
        gate_policy = self.nodes["owner_gate"]["gate_policy"]
        self.assertEqual(gate_policy["required_gate"], "project_watch_gate")
        self.assertTrue(gate_policy["pending_owner_confirmation"])
        self.assertEqual(self.nodes["challenger_review"].get("slot_ref"), "exception_challenger")

    def test_gm_handoff_is_declaration_only(self):
        node = self.nodes["gm_handoff"]
        self.assertTrue(node["declaration_only"])
        self.assertEqual(node["department"], "project")
        self.assertEqual(node["topic"], "project_watch")
        self.assertEqual(node["fallback_policy"], "not_ready")

    # ---- manifest ↔ capabilities 逐条对齐 -----------------------------------
    def test_manifest_and_capabilities_requirements_align_one_by_one(self):
        manifest_by_id = {item["requirement_id"]: item for item in self.manifest["provider_requirements"]}
        caps_by_id = {item["requirement_id"]: item for item in self.capabilities["provider_requirements"]}
        self.assertEqual(set(manifest_by_id), set(caps_by_id))
        self.assertEqual(len(manifest_by_id), 5)
        for requirement_id, requirement in manifest_by_id.items():
            for field in ALIGNED_REQUIREMENT_FIELDS:
                self.assertEqual(requirement.get(field), caps_by_id[requirement_id].get(field),
                                 "%s.%s 不一致" % (requirement_id, field))
            self.assertIn(requirement["execution_level"], {"L1", "L2", "L3", "L4", "L5", "L6"})
            self.assertIn(requirement["runtime_requirement"], {"cloud_ok", "local_preferred", "local_required"})
            self.assertTrue(str(caps_by_id[requirement_id].get("description", "")).strip())
        self.assertEqual(
            {item["capability"] for item in manifest_by_id.values()},
            set(self.manifest["capabilities_required"]))
        self.assertEqual(
            {item["capability_ref"] for item in self.capabilities["capabilities_required"]},
            set(self.manifest["capabilities_required"]))

    def test_readmodel_capability_is_not_ready_by_default_and_names_the_endpoint(self):
        requirement = next(item for item in self.capabilities["provider_requirements"]
                           if item["requirement_id"] == "req_project_anomaly_readmodel")
        self.assertEqual(requirement["fallback_policy"], "not_ready")
        self.assertEqual(requirement["gateway_class"], "memory")
        self.assertEqual(requirement["provider_family"], "readmodel")
        self.assertIn("GET /v3/business-object/project-anomalies", requirement["description"])

    # ---- role slots 镜像 ---------------------------------------------------
    def test_manifest_role_slots_mirror_role_slots_file(self):
        canonical = {slot["slot_id"]: slot for slot in self.role_slots["role_slots"]}
        mirrored = {slot["slot_id"]: slot for slot in self.manifest["role_slots"]}
        self.assertEqual(set(canonical), {"project_watcher", "exception_challenger"})
        self.assertEqual(set(canonical), set(mirrored))
        for slot_id, slot in mirrored.items():
            for field in ("responsibility", "node_type", "default_role_pack_ref"):
                self.assertEqual(slot[field], canonical[slot_id][field], "%s.%s" % (slot_id, field))
        bindings = {b["slot_id"]: b for b in self.role_slots["bindings"]}
        self.assertEqual(set(bindings), set(canonical))
        for slot_id, binding in bindings.items():
            self.assertEqual(binding["role_pack_id"], canonical[slot_id]["default_role_pack_ref"])
            self.assertEqual(binding["required_role"], canonical[slot_id]["required_role"])
            rolepack = load("role-packs/%s.rolepack.json" % slot_id.replace("_", "-"))
            self.assertEqual(rolepack["role_pack_id"], binding["role_pack_id"])
            self.assertIn(slot_id, rolepack["target_use"])

    def test_multi_role_comparison_nodes_match_slot_node_types(self):
        canonical = {slot["slot_id"]: slot["node_type"] for slot in self.role_slots["role_slots"]}
        declared = set(self.manifest["multi_role_comparison"]["explicit_nodes"])
        self.assertEqual(declared, {"%s/%s" % (slot_id, node_type) for slot_id, node_type in canonical.items()})
        self.assertTrue(self.manifest["multi_role_comparison"]["hidden_agent_loops_forbidden"])

    def test_gate_flags_high_risk_actions_and_command_routes_share_one_vocabulary(self):
        gate_flags = set(self.manifest["gate_flags"])
        high_risk = set(self.manifest["person_strategy"]["high_risk_actions_return_to_owner"])
        commands = set(self.manifest["notification_command_report_routes"]["command_candidate"])
        self.assertEqual(gate_flags, high_risk)
        self.assertEqual(gate_flags, commands)
        self.assertEqual(gate_flags, {
            "anomaly_escalate_confirm",
            "followup_task_write_confirm",
            "owner_notice_send_confirm",
            "formal_memory_confirm",
        })
        self.assertFalse(self.manifest["person_strategy"]["delegation_allowed"])
        self.assertEqual(set(self.manifest["person_strategy"]["proposer_roles"]),
                         {"project_watcher", "exception_challenger"})

    # ---- knowledge 一致性 ---------------------------------------------------
    def test_knowledge_scopes_are_consistent_across_manifest_scopes_and_index(self):
        declared = set(self.manifest["knowledge_scopes"])
        actual = {scope["scope_ref"] for scope in self.scopes["scopes"]}
        self.assertEqual(declared, actual)
        self.assertEqual(declared, {"knowledge_scope://project-watch/anomaly-rules"})
        entry_scopes = {entry["knowledge_scope_ref"] for entry in self.index["entries"]}
        self.assertTrue(entry_scopes <= declared)
        self.assertEqual(self.index["count"], len(self.index["entries"]))
        scene_refs = {self.scopes["scene_ref"], self.index["scene_ref"]}
        self.assertEqual(scene_refs, {"scene://project-watch"})
        for entry in self.index["entries"]:
            self.assertEqual(entry["verification_status"], "pending_human_review")
            self.assertTrue((PACK / entry["file"]).exists())
            self.assertTrue(entry["checksum"].startswith("sha256:"))
            scope = next(s for s in self.scopes["scopes"] if s["scope_ref"] == entry["knowledge_scope_ref"])
            self.assertIn(entry["kind"], scope["knowledge_kinds"])

    # ---- 诚实标注与禁品 -----------------------------------------------------
    def test_lifecycle_status_is_honestly_design_stage(self):
        self.assertEqual(self.manifest["lifecycle_status"], "设计中")
        self.assertTrue(self.manifest["security_profile"]["candidate_only"])
        self.assertTrue(self.manifest["security_profile"]["requires_base_gate"])
        self.assertTrue(self.manifest["receipt_policy"]["append_only"])
        self.assertEqual(self.manifest["software_requirements"], [])
        self.assertTrue(str(self.manifest["software_requirements_note"]).strip())

    def test_pack_assets_carry_no_forbidden_literals(self):
        forbidden_substrings = (
            "owner_action_evidence" + "://",
            "receipt" + "://",
            "decision" + "://",
            "pack_version" + "://",
        )
        forbidden_patterns = (
            re.compile(r':\s*"1[3-9]\d{9}"'),
            re.compile(r':\s*"\d{17}[\dXx]"'),
        )
        checked = 0
        for path in sorted(PACK.rglob("*")):
            if not path.is_file() or path.suffix not in {".json", ".py", ".md"}:
                continue
            text = path.read_text(encoding="utf-8")
            checked += 1
            for token in forbidden_substrings:
                self.assertNotIn(token, text, "%s 出现禁用字面串 %s" % (path, token))
            for pattern in forbidden_patterns:
                self.assertIsNone(pattern.search(text), "%s 出现疑似真实 PII" % path)
        self.assertGreaterEqual(checked, 10)


if __name__ == "__main__":
    unittest.main()
