#!/usr/bin/env python3
"""智能家居 Pack 的离线商品化契约防线；不调用 OS 或任何 Provider。"""

import json
import pathlib
import shutil
import sys
import tempfile
import unittest


PACK = pathlib.Path(__file__).resolve().parents[1]
REPO = PACK.parents[0]
if str(REPO) not in sys.path:
    sys.path.insert(0, str(REPO))


def load(relative):
    return json.loads((PACK / relative).read_text(encoding="utf-8"))


class SmartHomeOwnerPackContractTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = load("manifest.json")
        cls.flow = load("flows/smart-home-owner-project-ops-flow.flow.json")
        cls.capabilities = load("capabilities/capabilities.json")

    def provider(self, requirement_id):
        return next(item for item in self.manifest["provider_requirements"]
                    if item["requirement_id"] == requirement_id)

    def software(self, requirement_id):
        return next(item for item in self.manifest["software_requirements"]
                    if item["requirement_id"] == requirement_id)

    def test_pack_lifecycle_is_single_current_state(self):
        self.assertEqual(self.manifest["version"], "1.2.0")
        self.assertEqual(self.manifest["lifecycle_status"], "已实现")

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
        requirement = self.provider("frappe_project_lifecycle_write_candidate")
        self.assertEqual(requirement["required_capabilities"], ["project_lifecycle_write_candidate"])
        self.assertEqual(requirement["software_requirement_refs"],
                         ["frappe-suite-runtime", "frappe-mcp-runtime"])
        self.assertEqual(requirement["gateway_class"], "execution")
        self.assertEqual(requirement["fallback_policy"], "blocked")
        self.assertTrue(self.manifest["security_profile"]["requires_base_gate"])
        self.assertTrue(self.manifest["receipt_policy"]["append_only"])

    def test_optional_home_assistant_cannot_unblock_project_mainline(self):
        nodes = {node["id"]: node for node in self.flow["nodes"]}
        device = nodes["optional_device_control_candidate"]
        self.assertTrue(device["optional"])
        self.assertEqual(device["candidate_type"], "ExecutionIntentCandidate")
        requirement = self.provider("home_assistant_device_control_candidate")
        self.assertTrue(requirement["optional"])
        self.assertEqual(requirement["required_capabilities"], ["smart_home_device_control_candidate"])
        self.assertEqual(requirement["fallback_policy"], "not_ready")
        edge = next(item for item in self.flow["edges"]
                    if item["source"] == "optional_device_gate" and item["target"] == "optional_device_gateway")
        self.assertEqual(edge["condition"], "approved")

    def test_provider_and_software_version_declarations_are_explicit(self):
        self.assertEqual(self.software("frappe-suite-runtime")["version_range"], ">=16.0.0,<17.0.0")
        self.assertEqual(self.software("frappe-mcp-runtime")["version_range"], ">=0.2.1,<0.3.0")
        self.assertTrue(self.software("home-assistant-runtime")["optional"])
        self.assertEqual(self.software("home-assistant-runtime")["fallback_policy"], "not_ready")
        self.assertNotIn("baserow-runtime-a", {item["requirement_id"]
                                                for item in self.manifest["software_requirements"]})
        self.assertNotIn("shared-document-ocr", {item["requirement_id"]
                                                  for item in self.manifest["software_requirements"]})

    def test_capabilities_are_definitions_with_refs_only(self):
        self.assertEqual(set(self.capabilities), {"capabilities"})
        provider_ids = {item["requirement_id"] for item in self.manifest["provider_requirements"]}
        capability_ids = set()
        for capability in self.capabilities["capabilities"]:
            self.assertTrue(set(capability) <= {"capability_id", "provider_requirement_ref",
                                                "description", "optional"})
            self.assertNotIn("binding", capability)
            self.assertNotIn("software_requirement_refs", capability)
            self.assertNotIn("gateway_class", capability)
            self.assertNotIn("risk_class", capability)
            capability_ids.add(capability["capability_id"])
            self.assertIn(capability["provider_requirement_ref"], provider_ids)
        required = {capability for provider in self.manifest["provider_requirements"]
                    for capability in provider["required_capabilities"]}
        self.assertEqual(capability_ids, required)

    def test_provider_requirements_are_canonical_and_closed(self):
        allowed = {"requirement_id", "provider_family", "gateway_class", "required_capabilities",
                   "software_requirement_refs", "risk_class", "fallback_policy", "optional"}
        providers = self.manifest["provider_requirements"]
        self.assertEqual(len({item["requirement_id"] for item in providers}), len(providers))
        software_ids = {item["requirement_id"] for item in self.manifest["software_requirements"]}
        for provider in providers:
            self.assertTrue(set(provider) <= allowed)
            self.assertTrue(provider["required_capabilities"])
            self.assertTrue(provider["software_requirement_refs"])
            self.assertTrue(set(provider["software_requirement_refs"]) <= software_ids)
            for software_id in provider["software_requirement_refs"]:
                self.assertEqual(provider["provider_family"], self.software(software_id)["provider_family"])

    def test_flow_provider_refs_are_explicit_and_closed(self):
        provider_ids = {item["requirement_id"] for item in self.manifest["provider_requirements"]}
        flow_refs = {node["provider_requirement_ref"] for node in self.flow["nodes"]
                     if "provider_requirement_ref" in node}
        self.assertEqual(flow_refs, provider_ids)
        self.assertEqual(self.flow["version"], "1.2.0")
        self.assertEqual(self.manifest["flow_spec_ref"].rsplit("@", 1)[-1], "1.2.0")

    def assert_builder_rejects(self, mutate, expected):
        from build_pack_bundle import _validate

        tmp = pathlib.Path(tempfile.mkdtemp(prefix="smart-home-pack-contract-"))
        try:
            copied = tmp / PACK.name
            shutil.copytree(PACK, copied)
            manifest_path = copied / "manifest.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            mutate(manifest, copied)
            manifest_path.write_text(json.dumps(manifest, ensure_ascii=False), encoding="utf-8")
            with self.assertRaisesRegex(ValueError, expected):
                _validate(str(copied))
        finally:
            shutil.rmtree(tmp, ignore_errors=True)

    def test_builder_rejects_illegal_provider_fields(self):
        self.assert_builder_rejects(
            lambda manifest, _: manifest["provider_requirements"][0].update({"binding": "legacy"}),
            "非法字段")

    def test_builder_rejects_missing_software_ref_and_family_mismatch(self):
        self.assert_builder_rejects(
            lambda manifest, _: manifest["provider_requirements"][0]["software_requirement_refs"].append(
                "pack://other/software"),
            "未知或跨 Pack")
        self.assert_builder_rejects(
            lambda manifest, _: manifest["provider_requirements"][0].update({"provider_family": "baserow"}),
            "provider_family 不一致")

    def test_builder_rejects_duplicate_provider_and_orphan_software_ids(self):
        self.assert_builder_rejects(
            lambda manifest, _: manifest["provider_requirements"].append(
                dict(manifest["provider_requirements"][0])),
            "requirement_id 重复")

        def add_orphan(manifest, _):
            orphan = dict(manifest["software_requirements"][0])
            orphan["requirement_id"] = "orphan-runtime"
            manifest["software_requirements"].append(orphan)

        self.assert_builder_rejects(add_orphan, "孤儿声明")

    def test_builder_rejects_dangling_flow_ref(self):
        def mutate(manifest, copied):
            flow_path = copied / manifest["flow_file"]
            flow = json.loads(flow_path.read_text(encoding="utf-8"))
            next(node for node in flow["nodes"] if node["id"] == "frappe_snapshot")["provider_requirement_ref"] = \
                "pack://other/provider"
            flow_path.write_text(json.dumps(flow, ensure_ascii=False), encoding="utf-8")

        self.assert_builder_rejects(mutate, "未知或跨 Pack")

    def test_builder_rejects_duplicate_flow_ids(self):
        def mutate(manifest, copied):
            flow_path = copied / manifest["flow_file"]
            flow = json.loads(flow_path.read_text(encoding="utf-8"))
            duplicate = dict(flow["nodes"][0])
            flow["nodes"].append(duplicate)
            flow_path.write_text(json.dumps(flow, ensure_ascii=False), encoding="utf-8")

        self.assert_builder_rejects(mutate, "flow node id 重复")

    def test_old_capability_requirement_shape_is_warning_and_not_executable(self):
        def mutate(manifest, copied):
            caps_path = copied / manifest["capabilities_file"]
            caps_path.write_text(json.dumps({"provider_requirements": []}), encoding="utf-8")

        self.assert_builder_rejects(mutate, "migration_warning")


if __name__ == "__main__":
    unittest.main()
