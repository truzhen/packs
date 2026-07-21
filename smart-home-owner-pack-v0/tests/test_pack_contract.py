#!/usr/bin/env python3
"""智能家居 Pack 的离线商品化契约防线；不调用 OS 或任何 Provider。"""

import json
import pathlib
import unittest


PACK = pathlib.Path(__file__).resolve().parents[1]


def load(relative):
    return json.loads((PACK / relative).read_text(encoding="utf-8"))


class SmartHomeOwnerPackContractTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = load("manifest.json")
        cls.flow = load("flows/smart-home-owner-project-ops-flow.flow.json")
        cls.capabilities = load("capabilities/capabilities.json")

    def test_five_project_stages_are_candidates_before_gate(self):
        nodes = {node["id"]: node for node in self.flow["nodes"]}
        required = {
            "opportunity_candidate": "BusinessObjectCandidate",
            "project_initiation_candidate": "BusinessObjectCandidate",
            "progress_candidate": "TaskCandidate",
            "material_candidate": "BusinessObjectCandidate",
            "delivery_candidate": "TaskCandidate",
        }
        for node_id, candidate_type in required.items():
            self.assertEqual(nodes[node_id].get("candidate_type"), candidate_type)
        edges = {(edge["source"], edge["target"]): edge for edge in self.flow["edges"]}
        self.assertEqual(edges[("frappe_write_candidate", "owner_gate")].get("id"), "e11")
        self.assertEqual(edges[("owner_gate", "gateway_execution")].get("condition"), "approved")
        self.assertIn(("gateway_execution", "project_receipt"), edges)
        self.assertIn(("project_receipt", "history_query"), edges)

    def test_frappe_write_has_gate_gateway_and_receipt_requirements(self):
        capability = next(item for item in self.capabilities["provider_requirements"]
                          if item["requirement_id"] == "frappe_project_lifecycle_write_candidate")
        self.assertEqual(capability["gateway_class"], "execution")
        self.assertEqual(capability["fallback_policy"], "blocked")
        self.assertIn("Owner + Base Gate", capability["description"])
        self.assertTrue(self.manifest["security_profile"]["requires_base_gate"])
        self.assertTrue(self.manifest["receipt_policy"]["append_only"])

    def test_optional_home_assistant_cannot_unblock_project_mainline(self):
        nodes = {node["id"]: node for node in self.flow["nodes"]}
        device = nodes["optional_device_control_candidate"]
        self.assertTrue(device["optional"])
        self.assertEqual(device["candidate_type"], "ExecutionIntentCandidate")
        requirement = next(item for item in self.capabilities["provider_requirements"]
                           if item["requirement_id"] == "home_assistant_device_control_candidate")
        self.assertTrue(requirement["optional"])
        self.assertEqual(requirement["fallback_policy"], "not_ready")
        edge = next(item for item in self.flow["edges"]
                    if item["source"] == "optional_device_gate" and item["target"] == "optional_device_gateway")
        self.assertEqual(edge["condition"], "approved")

    def test_provider_and_software_version_declarations_are_explicit(self):
        software = {item["provider_family"]: item for item in self.manifest["software_requirements"]}
        self.assertEqual(software["frappe"]["version_range"], ">=16.0.0,<17.0.0")
        self.assertEqual(software["baserow"]["fallback_policy"], "provider_missing")
        self.assertTrue(software["home_assistant"]["optional"])
        self.assertEqual(software["home_assistant"]["fallback_policy"], "not_ready")


if __name__ == "__main__":
    unittest.main()
