#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""询盘管理 Pack 的离线契约防线；不调用 OS、模型或任何 Provider。

运行：python3 -m unittest discover -s inquiry-management-pack-v0/tests -v
"""

import json
import pathlib
import re
import unittest


PACK = pathlib.Path(__file__).resolve().parents[1]

# AGENTS.md §4.2 候选类型白名单，八类，不得发明第九类。
CANDIDATE_WHITELIST = {
    "TaskCandidate",
    "MemoryRequestCandidate",
    "CommunicationDraftCandidate",
    "ExecutionIntentCandidate",
    "BusinessObjectCandidate",
    "SceneFlowRunCandidate",
    "CapabilityInvocationCandidate",
    "PackCandidate",
}

# 唯一例外：AdviceCandidate 是 contracts 实有类型（truzhen-contracts
# candidates/advice.go，带 cited_knowledge_refs），属基座 06 collaboration.advice
# 节点的产出语义，不是 Pack 自行「生成」的候选，故只允许出现在 advice 节点上。
ADVICE_NODE_CANDIDATE = "AdviceCandidate"
ADVICE_NODE_TYPES = {"collaboration.advice", "collaboration.challenge", "collaboration.compare_gate"}

# 14 制作台词表（truzhenos backend/internal/packstudio/nodeinfo/nodeinfo.go 固定注册）。
# 未注册的 type 在 convertCanvasDraftToSceneFlowSpec 会失败；扩词表属跨仓治理动作。
REGISTERED_NODE_TYPES = {
    "capability.invoke", "collaboration.advice", "collaboration.challenge",
    "collaboration.compare_gate", "draft.spec_document", "flow.end",
    "flow.stage_task_candidate", "gateway.communication_draft", "gateway.execution_intent",
    "input.user_need", "judgment.ai", "judgment.rule", "object.business_schema",
    "policy.gate_config", "receipt.link", "simulation.futureshadow", "wait.contact_confirm",
}

# truzhenos backend/internal/packstudio/softwareproject/template_family.go 的 12 族目录，
# 由 template_family_matrix_test.go 钉死；Pack 不得自造族名。
REGISTERED_TEMPLATE_FAMILIES = {
    "长周期项目交付型", "平台撮合运营型", "工单服务履约型", "关系经营推进型",
    "合规审查执法证据链型", "内容生产流水线型", "资产库存供应链型", "长期照护长期陪跑型",
    "软件数字产品交付型", "文书审批办文办会型", "交易合同账款履约型", "经营财务管理会计型",
}

# 厂商词只允许出现在 manifest 的 provider_family / software_family 值里。
VENDOR_WORDS = ("frappe", "erpnext", "odoo", "salesforce", "kingdee", "baserow")
VENDOR_ALLOWED_KEYS = {"provider_family", "software_family"}

# 询盘声明字段与快照都不得携带联系方式 / 身份 / 账户类 PII 键。
PII_KEYS = ("mobile_no", "email_id", "phone", "id_number", "bank_account")


def load(relative):
    return json.loads((PACK / relative).read_text(encoding="utf-8"))


def walk(node, path=()):
    """深度遍历 JSON，产出 (路径, 键, 值) 三元组。"""
    if isinstance(node, dict):
        for key, value in node.items():
            yield path, key, value
            yield from walk(value, path + (key,))
    elif isinstance(node, list):
        for index, value in enumerate(node):
            yield from walk(value, path + (str(index),))


class InquiryManagementPackContractTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = load("manifest.json")
        cls.flow = load("flows/inquiry-management-flow.flow.json")
        cls.capabilities = load("capabilities/capabilities.json")
        cls.role_slots = load("role-slots/role-slots.json")
        cls.role_pack = load("role-packs/inquiry-manager.rolepack.json")
        cls.nodes = {node["id"]: node for node in cls.flow["nodes"]}
        cls.edges = {(edge["source"], edge["target"]): edge for edge in cls.flow["edges"]}

    # ---- 1. manifest ↔ capabilities 名字闭合（smart-home 自身的已知不闭合缺陷，本 Pack 必须闭合） ----

    def test_manifest_and_capabilities_capability_names_are_closed(self):
        declared = set(self.manifest["capabilities_required"])
        actual = {item["capability_ref"] for item in self.capabilities["capabilities_required"]}
        self.assertEqual(declared, actual, "manifest.capabilities_required 与 capabilities.json 的 capability_ref 必须一一闭合")
        self.assertEqual(set(self.manifest["required_capabilities"]), actual)

    def test_manifest_and_capabilities_provider_requirement_ids_are_closed(self):
        manifest_ids = {item["requirement_id"] for item in self.manifest["provider_requirements"]}
        capability_ids = {item["requirement_id"] for item in self.capabilities["provider_requirements"]}
        self.assertEqual(manifest_ids, capability_ids)
        for manifest_item in self.manifest["provider_requirements"]:
            twin = next(item for item in self.capabilities["provider_requirements"]
                        if item["requirement_id"] == manifest_item["requirement_id"])
            for field in ("capability", "gateway_class", "risk_class", "fallback_policy"):
                self.assertEqual(manifest_item[field], twin[field],
                                 "%s 的 %s 在 manifest 与 capabilities 之间漂移" % (manifest_item["requirement_id"], field))

    def test_every_capability_ref_used_in_flow_is_declared(self):
        declared = {item["capability_ref"] for item in self.capabilities["capabilities_required"]}
        used = {node["capability_ref"] for node in self.flow["nodes"] if "capability_ref" in node}
        self.assertTrue(used, "flow 必须至少显式引用一个能力")
        self.assertTrue(used <= declared, "flow 引用了未声明的能力：%s" % sorted(used - declared))
        self.assertEqual(used, declared, "声明了却从不被 flow 引用的能力必须删除或落到节点上")

    def test_every_provider_requirement_is_referenced_by_flow(self):
        declared = {item["requirement_id"] for item in self.manifest["provider_requirements"]}
        used = {node["provider_requirement_ref"] for node in self.flow["nodes"]
                if "provider_requirement_ref" in node}
        self.assertEqual(used, declared, "ProviderRequirement 必须与 flow 节点引用闭合")

    # ---- 2. flow 的 04 契约字段：显式 capability_domain / capability_operation，且无厂商词 ----

    def test_capability_nodes_declare_domain_and_operation(self):
        declared_operations = set()
        for capability in self.capabilities["capabilities_required"]:
            for binding in capability.get("capability_bindings", []):
                declared_operations.add((binding["capability_domain"], binding["capability_operation"]))

        seen = set()
        for node in self.flow["nodes"]:
            if node["type"] not in ("capability.invoke", "gateway.execution_intent"):
                continue
            self.assertIn("capability_domain", node, "%s 缺 04 契约字段 capability_domain" % node["id"])
            self.assertIn("capability_operation", node, "%s 缺 04 契约字段 capability_operation" % node["id"])
            pair = (node["capability_domain"], node["capability_operation"])
            self.assertIn(pair, declared_operations,
                          "%s 的 %s 未在 capabilities.json 的 capability_bindings 中声明" % (node["id"], pair))
            seen.add(pair)
        self.assertEqual(seen, declared_operations, "声明的能力操作必须全部落到 flow 节点上，不留悬空声明")

    def test_expected_six_capability_operations_are_present(self):
        expected = {
            ("customer-relationship", "customer.list_leads"),
            ("customer-relationship", "customer.get_lead"),
            ("customer-relationship", "customer.get_customer"),
            ("selling", "selling.list_quotations"),
            ("selling", "selling.get_quotation"),
            ("selling", "selling.create_quotation_candidate"),
        }
        actual = {(node["capability_domain"], node["capability_operation"])
                  for node in self.flow["nodes"] if "capability_operation" in node}
        self.assertEqual(actual, expected)

    def test_flow_carries_no_vendor_words(self):
        raw = (PACK / "flows" / "inquiry-management-flow.flow.json").read_text(encoding="utf-8").lower()
        for vendor in VENDOR_WORDS:
            self.assertNotIn(vendor, raw, "flow 不得出现厂商词 %s（含节点标题与能力操作名）" % vendor)

    def test_manifest_vendor_words_only_in_provider_family_fields(self):
        raw_manifest = json.loads((PACK / "manifest.json").read_text(encoding="utf-8"))
        offenders = []
        for _path, key, value in walk(raw_manifest):
            if not isinstance(value, str):
                continue
            lowered = value.lower()
            if any(vendor in lowered for vendor in VENDOR_WORDS) and key not in VENDOR_ALLOWED_KEYS:
                offenders.append((key, value))
        self.assertEqual(offenders, [], "厂商词只允许出现在 provider_family / software_family：%s" % offenders)

    # ---- 3. 候选类型八类白名单 + 主权链节点化 ----

    def test_candidate_types_are_within_whitelist(self):
        for node in self.flow["nodes"]:
            if "candidate_type" not in node:
                continue
            candidate_type = node["candidate_type"]
            if candidate_type == ADVICE_NODE_CANDIDATE:
                self.assertIn(node["type"], ADVICE_NODE_TYPES,
                              "AdviceCandidate 只允许出现在 06 collaboration.* 节点：%s" % node["id"])
                continue
            self.assertIn(candidate_type, CANDIDATE_WHITELIST,
                          "%s 使用了白名单外的候选类型 %s" % (node["id"], candidate_type))

    def test_advice_candidate_is_the_only_exception_and_only_on_advice_nodes(self):
        extras = {node["candidate_type"] for node in self.flow["nodes"]
                  if node.get("candidate_type") and node["candidate_type"] not in CANDIDATE_WHITELIST}
        self.assertEqual(extras, {ADVICE_NODE_CANDIDATE},
                         "八类白名单之外只允许 AdviceCandidate 一个例外（contracts candidates/advice.go）")
        carriers = {node["id"] for node in self.flow["nodes"]
                    if node.get("candidate_type") == ADVICE_NODE_CANDIDATE}
        self.assertEqual(carriers, {"inquiry_manager_advice"})

    def test_node_types_are_registered_in_pack_studio_vocabulary(self):
        used = {node["type"] for node in self.flow["nodes"]}
        self.assertTrue(used <= REGISTERED_NODE_TYPES,
                        "使用了 14 制作台未注册的节点词：%s" % sorted(used - REGISTERED_NODE_TYPES))

    def test_template_family_is_from_the_registered_catalog(self):
        self.assertIn(self.manifest["template_family"], REGISTERED_TEMPLATE_FAMILIES,
                      "template_family 必须取自 12 族固定目录，不得自造")

    def test_four_business_candidates_use_the_ruled_types(self):
        self.assertEqual(self.nodes["inquiry_object_candidate"]["candidate_type"], "BusinessObjectCandidate")
        self.assertEqual(self.nodes["inquiry_object_candidate"]["object_type"], "inquiry")
        self.assertEqual(self.nodes["triage_candidate"]["candidate_type"], "TaskCandidate")
        self.assertEqual(self.nodes["triage_candidate"]["report_route"],
                         {"department": "inquiry", "topic": "inquiry"})
        self.assertEqual(self.nodes["followup_draft"]["candidate_type"], "CommunicationDraftCandidate")
        self.assertEqual(self.nodes["quotation_candidate"]["candidate_type"], "ExecutionIntentCandidate")
        self.assertEqual(self.nodes["quotation_candidate"]["report_route"], {"topic": "quotation_followup"})

    def test_real_send_is_delegated_to_communication_gateway_not_a_flow_node(self):
        types = {node["type"] for node in self.flow["nodes"]}
        self.assertNotIn("gateway.communication_send", types,
                         "制作台词表未注册发送类节点，不得自造")
        params = self.nodes["owner_gate_send"]["params"]
        self.assertIn("ControlledExternalSend", params["send_execution_note"])
        self.assertIn("03", params["send_execution_note"])
        self.assertIn("群发", params["send_scope"])

    def test_followup_draft_goes_through_model_gateway_as_document_draft(self):
        draft = self.nodes["followup_draft"]
        self.assertEqual(draft["model_gateway_candidate_type"], "DocumentDraft")
        self.assertEqual(draft["draft_generation_gateway"], "model_gateway")
        self.assertEqual(draft["send_gateway"], "communication_gateway")

    def test_advice_node_holds_no_formal_authority(self):
        advice = self.nodes["inquiry_manager_advice"]
        self.assertEqual(advice["type"], "collaboration.advice")
        self.assertEqual(advice["slot_ref"], "inquiry_manager")
        self.assertEqual(advice["formal_authority"], "none")
        self.assertEqual(advice["candidate_type"], ADVICE_NODE_CANDIDATE)

    def test_owner_gate_then_gateway_then_receipt_is_explicit(self):
        for gate_id in ("owner_gate_inquiry_object", "owner_gate_triage", "owner_gate_send", "owner_gate_quotation"):
            gate = self.nodes[gate_id]
            self.assertEqual(gate["type"], "policy.gate_config")
            self.assertTrue(gate["gate_policy"]["pending_owner_confirmation"])
        for source, target in (("owner_gate_inquiry_object", "inquiry_pool_scan"),
                               ("owner_gate_triage", "followup_draft"),
                               ("owner_gate_send", "followup_receipt"),
                               ("owner_gate_quotation", "gateway_execution")):
            self.assertEqual(self.edges[(source, target)].get("condition"), "approved",
                             "%s → %s 必须以 approved 为条件" % (source, target))
        self.assertIn(("gateway_execution", "quotation_snapshot_read"), self.edges)
        self.assertIn(("quotation_snapshot_read", "receipt"), self.edges)
        self.assertIn(("receipt", "handoff_to_project"), self.edges)

    def test_gate_flags_match_manifest_gates(self):
        flags = set(self.manifest["gate_flags"])
        self.assertEqual(flags, set(self.manifest["gates"].keys()))
        self.assertEqual(flags, set(self.manifest["person_strategy"]["high_risk_actions_return_to_owner"]))
        node_flags = {node["gate_policy"]["gate_flag"] for node in self.flow["nodes"]
                      if node["type"] == "policy.gate_config"}
        self.assertEqual(node_flags, flags, "四道门必须都在 flow 中显式节点化")

    def test_edges_reference_existing_nodes_only(self):
        for edge in self.flow["edges"]:
            self.assertIn(edge["source"], self.nodes)
            self.assertIn(edge["target"], self.nodes)

    # ---- 4. Owner 2026-09-02 裁定口径 ----

    def test_truth_source_is_external_and_pack_holds_snapshot_only(self):
        truth = self.manifest["external_truth_source"]
        self.assertEqual(truth["authority"], "external_erp")
        self.assertTrue(truth["snapshot_only"])
        self.assertEqual(sorted(truth["authoritative_object_kinds"]), ["customer", "lead", "quotation"])
        for node in self.flow["nodes"]:
            if node.get("capability_operation", "").startswith(("customer.", "selling.get", "selling.list")):
                self.assertTrue(node.get("read_only"), "%s 触碰外部真相源必须只读" % node["id"])

    def test_inquiry_source_and_quotation_to_follow_the_ruling(self):
        declaration = self.manifest["inquiry_object_declaration"]
        self.assertEqual(declaration["object_type"], "inquiry")
        self.assertEqual(declaration["inquiry_source_enum"], ["channel", "direct"])
        self.assertFalse(declaration["opportunity_doctype_dependency"],
                         "第一版不依赖独立商机对象类型")
        self.assertNotIn("quotation_party_policy", declaration["declared_fields"],
                         "Owner 第 15 条已裁：报价抬头不再是可配置的询盘声明字段")

    def test_quotation_to_is_fixed_to_the_service_party(self):
        """Owner 第 15 条已裁 2026-09-02：报价单发给业主、业主付费，无可配置策略。"""
        for rule in (self.manifest["inquiry_object_declaration"]["quotation_to_rule"],
                     self.nodes["quotation_candidate"]["quotation_to_rule"]):
            self.assertEqual(rule["quotation_to"], "service_party")
            self.assertTrue(rule["channel_party_forbidden"],
                            "渠道客户不得作 quotation_to / party_name")
            self.assertIn("第 15 条已裁", rule["owner_ruling"])
            self.assertNotIn("configurable", rule, "第 15 条已裁，不得保留可配置策略")
            resolution = {item["owner_party_kind"]: item["quotation_to"] for item in rule["resolution"]}
            self.assertEqual(resolution, {"lead": "lead", "customer": "customer"},
                             "业主为线索时 quotation_to=线索；业主已成客户时 quotation_to=业主客户")

    def test_no_configurable_quotation_party_policy_remains(self):
        for relative in ("manifest.json", "flows/inquiry-management-flow.flow.json",
                         "capabilities/capabilities.json", "README.md", "docs/派活卡.md"):
            raw = (PACK / relative).read_text(encoding="utf-8")
            self.assertNotIn("quotation_party_policy", raw,
                             "%s 仍残留已作废的可配置报价抬头策略" % relative)
            self.assertNotIn("owner_ruling_pending", raw)

    def test_pack_boundary_stops_at_quotation_and_declares_handoff(self):
        self.assertEqual(self.manifest["chain_scope"]["in_scope_segments"],
                         ["1-线索", "2-线索池", "3-商机分诊", "4-报价单"])
        handoff = self.nodes["handoff_to_project"]["handoff"]
        self.assertEqual(handoff["trigger"], "生成合同")
        joined = " ".join(handoff["out_of_scope"])
        self.assertIn("回款", joined)
        self.assertIn("开票", joined)
        self.assertIn("群发", joined)
        raw = (PACK / "manifest.json").read_text(encoding="utf-8")
        for out_of_scope in ("合同审批", "催办话术"):
            self.assertNotIn(out_of_scope + "候选", raw)

    def test_person_strategy_forbids_delegation(self):
        self.assertFalse(self.manifest["person_strategy"]["delegation_allowed"])
        self.assertFalse(self.role_pack["person_strategy"]["delegation_allowed"])
        self.assertEqual(self.role_pack["person_strategy"]["role_authority"], "proposer_only")
        self.assertEqual(self.manifest["person_strategy"]["proposer_roles"], ["inquiry_manager"])

    def test_role_slot_binding_uses_role_pack_scheme(self):
        slot = self.role_slots["role_slots"][0]
        binding = self.role_slots["bindings"][0]
        self.assertEqual(slot["slot_id"], "inquiry_manager")
        self.assertEqual(slot["node_type"], "advice")
        self.assertTrue(slot["default_role_pack_ref"].startswith("role_pack://"))
        self.assertTrue(binding["agent_ref"].startswith("role_pack://"))
        self.assertEqual(binding["role_pack_id"], self.role_pack["role_pack_id"])
        self.assertEqual(self.manifest["multi_role_comparison"]["explicit_nodes"],
                         ["inquiry_manager_advice(inquiry_manager/advice)"])
        self.assertTrue(self.manifest["multi_role_comparison"]["hidden_agent_loops_forbidden"])

    # ---- 5. 主权与诚实态硬性字段 ----

    def test_security_profile_is_candidate_only_and_non_formal(self):
        profile = self.manifest["security_profile"]
        for field in ("candidate_only", "non_formal", "requires_base_gate", "requires_receipt_candidate"):
            self.assertTrue(profile[field], "%s 必须为真" % field)
        for field in ("no_real_send", "no_real_execute"):
            self.assertTrue(profile[field])
        self.assertTrue(self.manifest["receipt_policy"]["append_only"])

    def test_provider_fallbacks_are_honest(self):
        by_id = {item["requirement_id"]: item for item in self.manifest["provider_requirements"]}
        for read_id in ("inquiry_lead_snapshot_read", "inquiry_customer_snapshot_read",
                        "inquiry_quotation_snapshot_read"):
            self.assertEqual(by_id[read_id]["risk_class"], "low")
            self.assertEqual(by_id[read_id]["fallback_policy"], "provider_missing")
        write = by_id["inquiry_quotation_write_candidate"]
        self.assertEqual(write["risk_class"], "medium")
        self.assertEqual(write["fallback_policy"], "not_ready")
        self.assertEqual({item["provider_family"] for item in self.manifest["provider_requirements"]}, {"frappe"})

    def test_no_knowledge_base_is_claimed(self):
        self.assertNotIn("knowledge_scopes", self.manifest)
        self.assertNotIn("knowledge_index", self.manifest)
        self.assertFalse((PACK / "knowledge").exists(), "本 Pack 无知识库；出现 knowledge/ 必须同时补齐索引与 checksum")
        readme = (PACK / "README.md").read_text(encoding="utf-8")
        self.assertIn("无知识库", readme)

    def test_no_pii_or_business_data_in_pack_assets(self):
        phone = re.compile(r":\s*\"1[3-9]\d{9}\"")
        for path in sorted(PACK.rglob("*.json")):
            raw = path.read_text(encoding="utf-8")
            self.assertIsNone(phone.search(raw), "%s 疑似携带真实手机号" % path)
            document = json.loads(raw)
            for _path, key, _value in walk(document):
                for pii in PII_KEYS:
                    self.assertNotEqual(key, pii, "%s 不得把 %s 作为字段键携带真实值" % (path, pii))

    def test_lifecycle_status_is_honest(self):
        self.assertEqual(self.manifest["lifecycle_status"], "设计中")
        self.assertEqual(self.manifest["version"], "0.1.0")
        self.assertEqual(self.manifest["kind"], "scene_pack")
        self.assertEqual(self.manifest["pack_type"], "domain_work_pack")
        self.assertEqual(self.manifest["pack_ref"], "scene_pack://inquiry-management")
        self.assertEqual(self.manifest["pack_id"], "inquiry-management-pack")

    def test_risk_types_carry_canonical_fields(self):
        required = {"risk_type_id", "definition", "trigger_action_types",
                    "evidence_requirement", "escalation_path", "fallback"}
        self.assertTrue(self.manifest["risk_types"])
        for item in self.manifest["risk_types"]:
            self.assertTrue(required <= set(item), "risk_type 缺 canonical 字段：%s" % item.get("risk_type_id"))
            self.assertIn(item["fallback"], {"blocked", "provider_missing", "not_ready", "manual_handoff"})


if __name__ == "__main__":
    unittest.main()
