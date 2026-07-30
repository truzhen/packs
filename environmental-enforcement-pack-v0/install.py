#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
环保执法 Pack —— 装入正在运行的 Truzhen devserver（可加载）。

脚本只负责 lifecycle candidate staging 与只读组合 readiness。场景包正式 confirm、
Base prepare/confirm 和知识候选 approve 必须由可信 GUI 中的 Owner 操作完成；脚本不代办。
只有 exact enabled pointer、每个 required scope 的 active mount 与对应 FormalReceipt
全部可反查时，才续接角色、绑槽与 09 知识候选暂存。

前置：先在隔离 worktree 起 devserver（且已从 server.go 摘除环保自动 seed），并显式指定
受控地址；脚本不会猜测或写入默认端口：
  TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 \
    python3 packs/environmental-enforcement-pack-v0/install.py

幂等：首次 staging 以稳定 key 建 candidate；Owner 操作后重跑同一命令只读 OS 权威状态。
已启用角色包/绑定会跳过，知识批量以 versioned source_ref + content_hash 去重。
"""
import hashlib
import argparse
import json
import os
import re
import sys
import urllib.error
import urllib.parse
import urllib.request

PACK_DIR = os.path.dirname(os.path.abspath(__file__))
REPO_DIR = os.path.dirname(PACK_DIR)
if REPO_DIR not in sys.path:
    sys.path.insert(0, REPO_DIR)
from pack_diagnostics import (
    emit_pack_error, INSTALL_GENERIC, INSTALL_CONNECTIVITY, INSTALL_LIFECYCLE_HTTP,
    INSTALL_READINESS, INSTALL_ROLE_BINDING, INSTALL_KNOWLEDGE,
    INSTALL_KNOWLEDGE_CHECKSUM)
from knowledge_checksums import verify_entries
from pack_install_journal import InstallJournal

# 由 main() 在参数解析后赋值。lifecycle 是写操作，绝不以隐式 localhost 默认值
# 猜测目标实例，避免测试或运维命令越出已登记的隔离端口。
BASE = ""
# 用本地规范 Owner（前端记忆中心默认查询 owner_id='owner://local/default'，后端运行时
# 也用此 owner）：知识/挂载/角色/绑定都落在这个 owner 下，记忆中心与运行时 advice 才看得到。
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")

JOURNAL = None  # 装入事务日志（#8）：main() 内初始化；die() 失败时落半装状态

LIFECYCLE_RECORD_STATES = {
    "draft",
    "readiness_checked",
    "pack_spec_candidate",
    "owner_confirmed",
    "gate_approved",
    "pack_enabled_version",
    "rolled_back",
    "disabled",
    "uninstalled",
    "expired",
}


def load(rel):
    with open(os.path.join(PACK_DIR, rel), encoding="utf-8") as f:
        return json.load(f)


def load_opt(rel, default):
    """可选清单：字段缺失或文件不存在时返回 default（让本 loader 同时适用无知识库的 pack）。"""
    if not rel:
        return default
    path = os.path.join(PACK_DIR, rel)
    if not os.path.exists(path):
        return default
    with open(path, encoding="utf-8") as f:
        return json.load(f)


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            raw = resp.read().decode("utf-8")
            code = resp.status
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8")
        code = e.code
    except Exception as e:
        return 0, {"_transport_error": str(e)}
    try:
        return code, json.loads(raw) if raw else {}
    except json.JSONDecodeError:
        return code, {"_raw": raw}


def die(msg, error_code=INSTALL_GENERIC):
    if JOURNAL is not None:
        JOURNAL.fail(error_code=error_code, message=msg)
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="install", error_code=error_code, message=msg)
    print("装入失败：" + msg, file=sys.stderr)
    sys.exit(1)


def strip_frontmatter(text):
    m = re.match(r"^---\s*\n.*?\n---\s*\n", text, re.S)
    return text[m.end():] if m else text


def emit_owner_handoff(status, *, pack_ref, version, reason, audit_refs=None,
                       readiness=None, candidate_refs=None):
    """输出机器可读 Owner 交接；调用方重跑同一命令只做幂等续接。"""
    handoff = {
        "status": status,
        "pack_ref": pack_ref,
        "target_version": version,
        "pack_version_ref": pack_ref + "@" + version,
        "reason": reason,
        "owner_action": "请在可信 Truzhen GUI 中完成所列 Owner 确认；本脚本不代办 Gate。",
        "resume": {
            "command": "TRUZHEN_DEVSERVER_BASE=%s python3 %s" % (
                BASE, os.path.join(PACK_DIR, "install.py")),
            "rule": "Owner 操作后重跑同一命令；每次都重新读取 OS 权威状态。",
        },
        "audit_refs": sorted(set(audit_refs or [])),
        "candidate_refs": sorted(set(candidate_refs or [])),
    }
    if readiness is not None:
        handoff["readiness"] = readiness
    print("TRUZHEN_PACK_HANDOFF=" + json.dumps(handoff, ensure_ascii=False, sort_keys=True))


def lifecycle_entry(body, pack_ref):
    matches = [
        entry for entry in (body.get("packs") or [])
        if entry.get("pack_ref") == pack_ref
    ]
    if len(matches) > 1:
        die("lifecycle ReadModel 返回多个 exact pack_ref，拒绝猜测：%s" % pack_ref,
            INSTALL_LIFECYCLE_HTTP)
    return matches[0] if matches else {}


def target_record_states(entry, version):
    records = entry.get("records", [])
    if not isinstance(records, list):
        raise ValueError("lifecycle records 不是数组")
    states = []
    for record in records:
        if not isinstance(record, dict):
            raise ValueError("lifecycle record 不是对象")
        if record.get("version") != version:
            continue
        state = record.get("state")
        if state not in LIFECYCLE_RECORD_STATES:
            raise ValueError("目标版本 lifecycle state 缺失或未知")
        states.append(state)
    return states


def mount_audit_snapshot(mounts):
    refs = []
    rows = []
    for mount in mounts:
        for key in ("enabled_receipt_ref", "disabled_receipt_ref", "last_receipt_ref"):
            if mount.get(key):
                refs.append(mount[key])
        rows.append({
            "mount_ref": mount.get("mount_ref", ""),
            "knowledge_scope_ref": mount.get("knowledge_scope_ref", ""),
            "status": mount.get("status", ""),
            "blocked_reason": mount.get("blocked_reason", ""),
            "enabled_receipt_ref": mount.get("enabled_receipt_ref", ""),
            "disabled_receipt_ref": mount.get("disabled_receipt_ref", ""),
            "last_receipt_ref": mount.get("last_receipt_ref", ""),
        })
    return refs, rows


def verify_combined_readiness(pack_ref, version, scopes_doc):
    """只消费 O-T1 第 8 节：exact pointer + required active mounts + FormalReceipt。"""
    pack_version_ref = pack_ref + "@" + version
    required_scopes = sorted(
        scope["scope_ref"] for scope in (scopes_doc.get("scopes") or [])
        if scope.get("required", True)
    )
    base_query = {
        "owner_id": OWNER,
        "pack_ref": pack_ref,
        "pack_version_ref": pack_version_ref,
    }
    code, body = call(
        "GET",
        "/v3/memory/knowledge/mounts?" + urllib.parse.urlencode(base_query),
    )
    if code != 200:
        return {
            "ready": False,
            "status": "not_ready",
            "reason_codes": ["knowledge_mount_truth_unavailable"],
            "audit_refs": [],
            "mounts": [],
            "http_status": code,
        }
    all_exact_mounts = [
        mount for mount in (body.get("mounts") or [])
        if mount.get("owner_id") == OWNER
        and mount.get("pack_ref") == pack_ref
        and mount.get("pack_version_ref") == pack_version_ref
    ]
    audit_refs, audit_rows = mount_audit_snapshot(all_exact_mounts)
    reason_codes = []
    ready_mounts = []

    for scope_ref in required_scopes:
        exact_scope_mounts = [
            mount for mount in all_exact_mounts
            if mount.get("knowledge_scope_ref") == scope_ref
        ]
        if any(mount.get("status") != "active" for mount in exact_scope_mounts):
            reason_codes.append("required_scope_has_non_active_mount:" + scope_ref)
            continue
        query = dict(base_query)
        query.update({"knowledge_scope_ref": scope_ref, "status": "active"})
        code, active_body = call(
            "GET",
            "/v3/memory/knowledge/mounts?" + urllib.parse.urlencode(query),
        )
        if code != 200:
            reason_codes.append("active_mount_query_failed:" + scope_ref)
            continue
        active_mounts = [
            mount for mount in (active_body.get("mounts") or [])
            if mount.get("owner_id") == OWNER
            and mount.get("pack_ref") == pack_ref
            and mount.get("pack_version_ref") == pack_version_ref
            and mount.get("knowledge_scope_ref") == scope_ref
            and mount.get("status") == "active"
        ]
        if len(active_mounts) != 1:
            reason_codes.append("required_scope_not_exactly_active:" + scope_ref)
            continue
        mount = active_mounts[0]
        receipt_ref = mount.get("enabled_receipt_ref", "")
        if not receipt_ref:
            reason_codes.append("active_mount_missing_enabled_receipt:" + scope_ref)
            continue
        audit_refs.append(receipt_ref)
        code, receipt = call(
            "GET",
            "/v3/receipts/" + urllib.parse.quote(receipt_ref, safe=""),
        )
        # O-T1 明确 200 响应即为 FormalReceipt；不在 Packs 自造额外 status 枚举。
        if code != 200 or receipt.get("receipt_ref") != receipt_ref:
            reason_codes.append("formal_receipt_lookup_failed:" + scope_ref)
            continue
        ready_mounts.append({
            "mount_ref": mount.get("mount_ref", ""),
            "knowledge_scope_ref": scope_ref,
            "enabled_receipt_ref": receipt_ref,
        })

    if reason_codes:
        has_recovery_signal = any(
            row["status"] in ("blocked", "disabled")
            or bool(row["blocked_reason"])
            or bool(row["disabled_receipt_ref"])
            for row in audit_rows
        )
        return {
            "ready": False,
            "status": "recovery" if has_recovery_signal else "not_ready",
            "reason_codes": reason_codes,
            "audit_refs": sorted(set(audit_refs)),
            "mounts": audit_rows,
            "required_scope_count": len(required_scopes),
            "active_scope_count": len(ready_mounts),
        }
    return {
        "ready": True,
        "status": "ready",
        "reason_codes": [],
        "audit_refs": sorted(set(audit_refs)),
        "mounts": ready_mounts,
        "required_scope_count": len(required_scopes),
        "active_scope_count": len(ready_mounts),
    }


def main():
    parser = argparse.ArgumentParser(
        description="将环保执法 Pack 装入显式指定的受控 Truzhen devserver。"
    )
    parser.add_argument(
        "--devserver-base",
        default=os.environ.get("TRUZHEN_DEVSERVER_BASE", "").strip(),
        help="受控 devserver 根地址；也可通过 TRUZHEN_DEVSERVER_BASE 指定（必填）。",
    )
    args = parser.parse_args()
    global BASE
    BASE = args.devserver_base.rstrip("/")
    if not BASE:
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base；拒绝猜测默认端口", INSTALL_CONNECTIVITY)
    manifest = load("manifest.json")
    flow = load(manifest["flow_file"])
    role_slots_doc = load(manifest["role_slots_file"])
    caps = load(manifest["capabilities_file"])
    scopes_doc = load_opt(manifest.get("knowledge_scopes_manifest"), {"scopes": []})
    kindex = load_opt(manifest.get("knowledge_index"), {"entries": []})
    # 装入前先证 pack 自身完整性：知识内容与 index checksum 漂移即拒绝装入（防漂移 #10）。
    checksum_problems = verify_entries(PACK_DIR, kindex.get("entries", []))
    if checksum_problems:
        die("知识内容与 index checksum 漂移，拒绝装入：" + "; ".join(checksum_problems),
            INSTALL_KNOWLEDGE_CHECKSUM)

    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    flow_id = flow["flow_id"]
    pack_version_ref = pack_ref + "@" + version

    print("== 装入 %s @ %s 到 %s ==" % (pack_ref, version, BASE))
    global JOURNAL
    JOURNAL = InstallJournal.open(pack_ref, BASE)
    JOURNAL.step("scene")

    # lifecycle ReadModel 是组合 readiness 的第一项，但绝不是单一真相。
    lifecycle_path = (
        "/v3/pack-studio/lifecycle/packs?"
        + urllib.parse.urlencode({"pack_ref": pack_ref})
    )
    code, body = call("GET", lifecycle_path)
    if code == 0:
        die("连不上 devserver（%s）。请先 go run ./backend/cmd/devserver" % BASE, INSTALL_CONNECTIVITY)
    if code != 200:
        die("lifecycle ReadModel HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    entry = lifecycle_entry(body, pack_ref)
    pointer = entry.get("enabled_pointer")
    try:
        if pointer is None:
            enabled_version = ""
        elif not isinstance(pointer, dict):
            raise ValueError("enabled_pointer 不是对象")
        else:
            enabled_version = pointer.get("current_version", "")
            if not isinstance(enabled_version, str):
                raise ValueError("enabled pointer current_version 不是字符串")
        record_states = target_record_states(entry, version)
    except ValueError as exc:
        emit_owner_handoff(
            "not_ready",
            pack_ref=pack_ref,
            version=version,
            reason="malformed_lifecycle_readmodel_fail_closed",
            readiness={"readmodel_error": str(exc)},
        )
        return

    def sync_canvas():
        # 只在尚无目标 lifecycle 候选时写穿规格；resume readiness 阶段保持只读。
        print("[1/6] 画布写穿 06 ...")
        code, body = call("POST", "/v3/pack-studio/canvas", {
            "flow_id": flow_id, "title": flow.get("title", ""),
            "occ_version": 0, "save_source": "pack_install", "flow_spec_draft": flow})
        if code == 409 and isinstance(body.get("current_occ_version"), (int, float)):
            code, body = call("POST", "/v3/pack-studio/canvas", {
                "flow_id": flow_id, "title": flow.get("title", ""),
                "occ_version": int(body["current_occ_version"]),
                "save_source": "pack_install", "flow_spec_draft": flow})
        if code != 200:
            die("canvas HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
        if not ((body.get("engine_sync") or {}).get("synced")):
            die("canvas 未同步进 06：%s" % body.get("engine_sync"), INSTALL_LIFECYCLE_HTTP)

    def stage_scene_candidate(start_state):
        pvr = pack_ref + "@" + version
        provider_reqs = []
        for p in caps.get("provider_requirements", []):
            provider_reqs.append({k: p[k] for k in ("requirement_id", "capability", "gateway_class",
                                                    "risk_class", "fallback_policy", "provider_family") if k in p}
                                 | ({"optional": p["optional"]} if p.get("optional") else {}))
        ks = []
        for s in scopes_doc.get("scopes", []):
            ks.append({"scope_ref": s["scope_ref"], "display_name": s["display_name"],
                       "scene_ref": s.get("scene_ref", ""), "mount_policy": s.get("mount_policy", "mount_on_pack_enable"),
                       "knowledge_kinds": s.get("knowledge_kinds", []), "required": s.get("required", True)})
        routes = manifest.get("notification_command_report_routes", {})
        draft = {
            "pack_ref": pack_ref, "version": version, "title": manifest["name"],
            "template_family": manifest.get("template_family", "合规审查执法证据链型"),
            "risk_level": manifest.get("risk_level", "medium"),
            "flow_id": flow_id,
            "role_slots": [{"slot_id": s["slot_id"], "responsibility": s["responsibility"],
                            "required_role": s["required_role"], "default_role_pack_ref": s["default_role_pack_ref"],
                            "slice_scope_policy": s.get("slice_scope_policy", ""),
                            "node_type": s.get("node_type", "")} for s in role_slots_doc["role_slots"]],
            "provider_requirements": provider_reqs,
            "formalization_requirement": manifest["formalization_requirement"]["summary"],
            "notification_routes": routes.get("notification", []),
            "command_candidates": routes.get("command_candidate", []),
            "report_routes": routes.get("report", []),
            "moat_justification": manifest["moat_justification"],
            "knowledge_scopes": ks,
            # 判事策略结构化骨架（#11）：声明不授权，Base RiskTypeGate 裁定；缺省为空。
            "risk_types": manifest.get("risk_types", []),
            "idempotency_key": "pack-install-draft:" + pvr,
            "actor_ref": OWNER,
        }
        if not start_state:
            print("[2/6] lifecycle draft（六件事 + %d 知识域）..." % len(scopes_doc.get("scopes", [])))
            code, response = call("POST", "/v3/pack-studio/lifecycle/draft", draft)
            if code != 200:
                die("draft HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
        if start_state in ("", "draft"):
            print("[3/6] readiness ...")
            code, response = call("POST", "/v3/pack-studio/lifecycle/readiness", {
                "pack_ref": pack_ref, "version": version, "actor_ref": OWNER})
            if code != 200:
                die("readiness HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
            rr = (response.get("record") or {}).get("readiness_report") or {}
            if not rr.get("ready"):
                die("readiness 未通过：%s" % response, INSTALL_READINESS)
        if start_state in ("", "draft", "readiness_checked"):
            print("[4/6] promote candidate ...")
            code, response = call("POST", "/v3/pack-studio/lifecycle/promote", {
                "pack_ref": pack_ref, "version": version, "actor_ref": OWNER})
            if code != 200:
                die("promote HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)

    if enabled_version != version:
        JOURNAL.set_version(version)
        if enabled_version:
            emit_owner_handoff(
                "not_ready",
                pack_ref=pack_ref,
                version=version,
                reason="enabled_version_mismatch",
                readiness={
                    "enabled_pointer_version": enabled_version,
                    "target_record_states": record_states,
                },
            )
            return
        if any(state in ("owner_confirmed", "gate_approved", "pack_enabled_version")
               for state in record_states):
            emit_owner_handoff(
                "recovery",
                pack_ref=pack_ref,
                version=version,
                reason="lifecycle_success_projection_incomplete",
                readiness={
                    "enabled_pointer_version": "",
                    "target_record_states": record_states,
                },
            )
            return
        if any(state in ("pack_spec_candidate", "rolled_back", "disabled",
                         "uninstalled", "expired") for state in record_states):
            emit_owner_handoff(
                "awaiting_owner_confirmation",
                pack_ref=pack_ref,
                version=version,
                reason="lifecycle_candidate_requires_trusted_gui",
                readiness={
                    "enabled_pointer_version": "",
                    "target_record_states": record_states,
                },
            )
            return
        start_state = record_states[-1] if record_states else ""
        if start_state not in ("", "draft", "readiness_checked"):
            emit_owner_handoff(
                "not_ready",
                pack_ref=pack_ref,
                version=version,
                reason="unknown_lifecycle_state_fail_closed",
                readiness={
                    "enabled_pointer_version": "",
                    "target_record_states": record_states,
                },
            )
            return
        sync_canvas()
        stage_scene_candidate(start_state)
        JOURNAL.mark("scene_candidate_staged")
        emit_owner_handoff(
            "awaiting_owner_confirmation",
            pack_ref=pack_ref,
            version=version,
            reason="lifecycle_candidate_staged",
            readiness={
                "enabled_pointer_version": "",
                "target_record_states": record_states,
            },
        )
        return

    JOURNAL.set_version(version)
    readiness = verify_combined_readiness(pack_ref, version, scopes_doc)
    if not readiness["ready"]:
        emit_owner_handoff(
            readiness["status"],
            pack_ref=pack_ref,
            version=version,
            reason="combined_readiness_incomplete",
            audit_refs=readiness["audit_refs"],
            readiness=readiness,
        )
        return
    JOURNAL.mark("scene_combined_ready")

    # 6a. 角色包
    print("[6/6] 组合 readiness 通过；续接角色、绑槽与知识候选 ...")
    JOURNAL.step("role_packs")
    code, rpbody = call("GET", "/v3/agent-orchestration/role-packs/readmodel")
    if code != 200:
        die("role pack ReadModel HTTP %d: %s" % (code, rpbody), INSTALL_ROLE_BINDING)
    enabled_rp = set()
    for ev in (rpbody.get("enabled_versions") or []):
        rid = ev.get("role_pack_id", "")
        enabled_rp.add(rid.split("@")[0])
    for fn in sorted(os.listdir(os.path.join(PACK_DIR, "role-packs"))):
        if not fn.endswith(".json"):
            continue
        rp = load(os.path.join("role-packs", fn))
        rid = rp["role_pack_id"]
        if rid in enabled_rp:
            print("    角色包 %s 已启用，跳过" % rid)
            continue
        install_role_pack(rp)
        print("    角色包 %s 已启用" % rid)
        JOURNAL.mark_item("role_packs", rid)

    # 6b. 绑槽
    JOURNAL.step("bindings")
    scope_ref = pack_version_ref
    for b in role_slots_doc.get("bindings", []):
        bind_slot(b, scope_ref)
        JOURNAL.mark_item("bindings", b["slot_id"])
    print("    绑槽完成")

    # 6c. 知识库入库（按知识域分组）
    JOURNAL.step("knowledge")
    entries = kindex.get("entries", [])
    pending_candidate_refs = []
    if not entries:
        print("    本 pack 无知识库，跳过知识入库")
    else:
        groups = {}
        for e in entries:
            groups.setdefault(e["knowledge_scope_ref"], []).append(e)
        for scope_ref_k, items in groups.items():
            pending_candidate_refs.extend(
                ingest_knowledge_scope(pack_ref, pack_version_ref, scope_ref_k, items)
            )
        print("    知识候选已幂等暂存：%d 条；本脚本不代办 Base 确认。" % len(pending_candidate_refs))

    JOURNAL.mark("downstream_candidates_staged")
    emit_owner_handoff(
        "awaiting_owner_confirmation",
        pack_ref=pack_ref,
        version=version,
        reason="knowledge_candidates_require_trusted_owner_review",
        audit_refs=readiness["audit_refs"],
        readiness=readiness,
        candidate_refs=pending_candidate_refs,
    )


def install_role_pack(rp):
    sani = re.sub(r'[/:@\s]+', "-", rp["role_pack_id"]).strip("-")
    draft_id = "role_pack_draft-pack-install-" + sani
    idem = "pack-install-" + sani
    cs = rp.get("communication_style", {})
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts", {
        "draft_id": draft_id, "role_pack_id": rp["role_pack_id"], "version": rp["version"],
        "display_name": rp["display_name"], "description": rp["description"],
        "target_use": rp.get("target_use", []), "style_summary": rp.get("style_summary", ""),
        "decision_style_summary": rp.get("decision_style_summary", ""),
        "scenario": rp.get("scenario", ""),
        "opening_line_candidate": rp.get("opening_line_candidate", ""),
        "example_dialogues": rp.get("example_dialogues", None),
        "model_policy_ref": rp.get("model_policy_ref", ""),
        "communication_style": {"structure": cs.get("structure", ""), "tone": cs.get("tone", ""),
                                "forbidden_phrases": cs.get("forbidden_phrases", [])},
        "forbidden_behavior_policy_ref": rp.get("forbidden_policy_ref", ""),
        "risk_level": rp.get("risk_level", "medium"), "owner_ref": OWNER,
        "evidence_refs": ["evidence://pack-install/role-pack/" + sani + "/draft"],
        "idempotency_key": idem + "-draft"})
    if code != 200:
        die("role draft HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    draft_id = body.get("draft_id") or (body.get("draft") or {}).get("draft_id") or draft_id
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/readiness-check", {})
    if code != 200:
        die("role readiness HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/promote-candidate", {
        "owner_ref": OWNER, "target_agent_ref": rp["role_pack_id"], "idempotency_key": idem + "-promote",
        "evidence_refs": ["evidence://pack-install/role-pack/promote"]})
    if code != 200 or not body.get("ok"):
        die("role promote HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/enable-candidate", {
        "role_pack_id": rp["role_pack_id"], "version": rp["version"], "target_agent_ref": rp["role_pack_id"],
        "owner_ref": OWNER, "idempotency_key": idem + "-enable",
        "evidence_refs": ["evidence://pack-install/role-pack/enable"]})
    if code != 200 or not body.get("ok"):
        die("role enable-candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/enable-confirm", {
        "role_pack_id": rp["role_pack_id"], "version": rp["version"], "target_agent_ref": rp["role_pack_id"],
        "owner_ref": OWNER, "idempotency_key": idem + "-confirm", "approve": True,
        "comment": rp["display_name"] + " 启用", "evidence_refs": ["evidence://pack-install/role-pack/enable-confirm"]})
    if code != 200 or body.get("status") != "enabled":
        die("role enable-confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def bind_slot(b, scope_ref):
    code, rm = call("GET", "/v3/agent-orchestration/agent-slots/readmodel")
    if code == 200:
        for bd in (rm.get("agent_slot_bindings") or []):
            if bd.get("slot_ref") == b["slot_id"] and bd.get("enabled_state") == "enabled" \
               and bd.get("scope_ref") in (scope_ref, ""):
                return
    code, body = call("POST", "/v3/agent-orchestration/agent-slots/bind-candidate", {
        "slot_ref": b["slot_id"], "scope_ref": scope_ref, "source_pack_ref": scope_ref,
        "required_role": b["required_role"], "requested_agent_ref": b["agent_ref"],
        "requested_role_pack_id": b["role_pack_id"], "ttl": "8760h",
        "evidence_refs": ["evidence://pack-install/agent-slot/" + b["slot_id"] + "/bind"]})
    if code != 200 or not body.get("ok"):
        die("bind-candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    binding_ref = body.get("binding_ref")
    code, body = call("POST", "/v3/agent-orchestration/agent-slots/confirm", {
        "binding_ref": binding_ref, "idempotency_key": "pack-install-slot-confirm-" + b["slot_id"],
        "approve": True, "evidence_refs": ["evidence://pack-install/agent-slot/" + b["slot_id"] + "/confirm"]})
    if code != 200 or body.get("status") != "enabled":
        die("slot confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def ingest_knowledge_scope(pack_ref, pack_version_ref, scope_ref, items):
    source_files = []
    scene_ref = items[0].get("scene_ref", "")
    ver = pack_version_ref.split("@")[-1]
    for e in items:
        with open(os.path.join(PACK_DIR, e["file"]), encoding="utf-8") as f:
            body_text = strip_frontmatter(f.read())
        content = "# " + e["title"] + "\n" + body_text
        # source_ref 带 pack 版本：把知识归属到本 pack 版本；Owner 操作后重跑时
        # 仍使用相同 source_ref/content_hash，由 09 按稳定身份幂等返回候选。
        versioned_source_ref = e["source_ref"] + "@" + ver
        law_meta = {"verification_status": "pending_human_review", "source_authority": "reference_only"}
        if e.get("authority"):
            law_meta["authority"] = e["authority"]
        if e.get("effective_from"):
            law_meta["effective_from"] = e["effective_from"]
        if e.get("effective_to"):
            law_meta["effective_to"] = e["effective_to"]
        source_files.append({
            "source_ref": versioned_source_ref, "file_name": os.path.basename(e["file"]),
            "content": content, "content_hash": "sha1-" + hashlib.sha1((content + "|" + pack_version_ref).encode("utf-8")).hexdigest(),
            "kind": e["kind"], "evidence_refs": ["evidence://09/pack-knowledge/" + pack_ref + "/" + e["source_ref"]],
            "law_meta": law_meta})
    batch = {
        "owner_id": OWNER, "scope": "Formal",
        "transaction_ref": "transaction://pack-knowledge:" + pack_version_ref,
        "scene_ref": scene_ref, "tags": ["pack 知识库", "环保执法"],
        "policy_snapshot_ref": "policy_snapshot://pack-knowledge/import",
        "source_files": source_files, "pack_ref": pack_ref, "pack_version_ref": pack_version_ref,
        "knowledge_scope_ref": scope_ref}
    code, body = call("POST", "/v3/memory/knowledge/batches", batch)
    if code not in (200, 201):
        die("knowledge batches HTTP %d (scope=%s): %s" % (code, scope_ref, body), INSTALL_KNOWLEDGE)
    JOURNAL.mark_item("knowledge_batches", scope_ref)
    candidates = body.get("candidates", []) or []
    pending_refs = []
    for c in candidates:
        cref = c.get("candidate_ref")
        if cref and c.get("status") == "pending":
            pending_refs.append(cref)
    return pending_refs


if __name__ == "__main__":
    main()
