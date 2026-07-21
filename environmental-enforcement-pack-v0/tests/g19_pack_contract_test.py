#!/usr/bin/env python3
"""G19 Pack 内静态防线：验证声明与去敏验收输入，不调用 Provider 或真实案件。"""

import hashlib
import json
import re
import sys
from datetime import date
from pathlib import Path


PACK = Path(__file__).resolve().parents[1]


def load(relative):
    return json.loads((PACK / relative).read_text(encoding="utf-8"))


def require(condition, message):
    if not condition:
        raise AssertionError(message)


def is_effective(source, as_of):
    start = date.fromisoformat(source["effective_from"])
    end = source.get("effective_to")
    return start <= as_of and (not end or as_of <= date.fromisoformat(end))


def test_legal_time_fence(manifest, flow, fence):
    require(manifest["legal_time_fence_file"] == "knowledge/legal-time-fence.json", "manifest 未引用法律时点防线")
    policy = manifest["legal_time_fence"]
    require(policy["policy_ref"] == fence["policy_ref"], "法律时点 policy_ref 漂移")
    require(policy["missing_or_unverified_outcome"] == "blocked_insufficient_legal_context", "缺法源必须 fail-closed")
    require(policy["future_law_outcome"] == "future_law_notice_only", "未来法不得直接适用")
    source = fence["authoritative_sources"][0]
    require(source["official_url"].startswith("https://www.npc.gov.cn/"), "法典必须保留中国人大网权威出处")
    require(source["effective_from"] == "2026-08-15", "法典生效日漂移")
    require(source["verification_status"] == "pending_human_review", "法律资料不得伪称人工核验完成")
    legal = next(node for node in flow["nodes"] if node["id"] == "legal_advice")
    require(legal["params"]["legal_time_fence_ref"] == fence["policy_ref"], "flow 未声明法律时点防线")
    require(set(legal["params"]["required_case_inputs"]) == set(fence["required_case_inputs"]), "案情法律输入不完整")

    before = date(2026, 8, 14)
    after = date(2026, 8, 15)
    require(not is_effective(source, before), "2026-08-15 前不得将法典当作直接依据")
    require(is_effective(source, after), "2026-08-15 起法典应可作为候选来源")


def test_evidence_and_fail_closed_controls(manifest, flow, fixture):
    require(fixture["fixture_status"] == "deidentified_synthetic_structurally_realistic", "G19 只允许去敏合成案情")
    require("@" not in json.dumps(fixture, ensure_ascii=False), "fixture 不得含邮箱或个人标识")
    for evidence in fixture["evidence"]:
        require(evidence["source_ref"].startswith("source_candidate://"), "证据只能引用候选来源")
        require(re.fullmatch(r"[0-9a-f]{64}", evidence["sha256"]), "证据必须带 SHA-256")
        require(evidence["verification_status"] == "verified", "进入对照前的 fixture 证据必须核验")
    controls = fixture["expected_controls"]
    require(controls["ocr_without_provider"] == "manual_handoff", "OCR 未接通必须人工交接")
    require(controls["model_without_owner_authorized_profile"] == "blocked", "模型未获授权必须阻断")
    require(controls["gate_deny"] == "candidate_only_no_formal_action", "Gate deny 不得形成正式动作")
    require(controls["receipt_failure"] == "transaction_not_complete", "Receipt failure 不得完成案件")
    require(controls["scene_flow_end"] == "transaction_not_complete", "SceneFlow end 不得等同交易完成")

    doc = next(node for node in flow["nodes"] if node["id"] == "doc_draft")
    require(doc["candidate_type"] == "DocumentDraft", "文书必须路由到 Model Gateway 的 DocumentDraft")
    require(doc["params"]["draft_generation"] == "model_gateway_non_template", "文书不得退化为模板拼接")
    require(doc["params"]["template_only_forbidden"] is True, "必须禁止模板冒充模型正文")
    cap = next(x for x in load("capabilities/capabilities.json")["provider_requirements"] if x["requirement_id"] == "legal_doc_draft")
    require(cap["fallback_policy"] == "blocked", "模型文书未接通不得假成功")
    ocr = next(x for x in load("capabilities/capabilities.json")["provider_requirements"] if x["requirement_id"] == "doc_ocr")
    require(ocr["fallback_policy"] == "manual_handoff", "OCR 降级必须人工交接")

    receipts = manifest["receipt_policy"]["transaction_completion_proofs"]
    require(receipts == fixture["receipt_classes"], "三类完成 Receipt 口径漂移")
    require(len(set(receipts)) == 3, "执行、送达、整改复查 Receipt 不可互代")
    done = next(node for node in flow["nodes"] if node["id"] == "done")
    require("不等于案件 Transaction 完成" in done["stage_guide"], "flow end 必须与案件完成隔离")


def main():
    manifest = load("manifest.json")
    flow = load(manifest["flow_file"])
    fence = load(manifest["legal_time_fence_file"])
    fixture = load("tests/fixtures/g19-deidentified-case.json")
    test_legal_time_fence(manifest, flow, fence)
    test_evidence_and_fail_closed_controls(manifest, flow, fixture)
    print("G19 Pack 法律时点与候选/正式隔离静态验收通过")


if __name__ == "__main__":
    try:
        main()
    except (AssertionError, KeyError, ValueError, json.JSONDecodeError) as exc:
        print("G19 Pack 静态验收失败: %s" % exc, file=sys.stderr)
        raise SystemExit(1)
