#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
环保执法 Pack —— 装入正在运行的 Truzhen devserver（可加载）。

脚本只负责 lifecycle candidate staging 与只读组合 readiness。场景包正式 confirm、
Base prepare/confirm 和知识候选 approve 必须由可信 GUI 中的 Owner 操作完成；脚本不代办。
只有 exact enabled pointer、每个 required scope 的 active mount 与对应 FormalReceipt
全部可反查时，才输出角色、槽位与知识的可信 GUI Owner 操作清单；不发下游写请求。

前置：先在隔离 worktree 起 devserver（且已从 server.go 摘除环保自动 seed），并显式指定
受控地址；脚本不会猜测或写入默认端口：
  TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 \
    python3 packs/environmental-enforcement-pack-v0/install.py

幂等：首次 staging 以稳定 key 建 lifecycle candidate；Owner 操作后重跑同一命令只读
OS 权威状态并输出稳定 handoff。下游正式对象只能由可信 GUI 与 backend-issued evidence 驱动。
"""
import argparse
import json
import os
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
    INSTALL_READINESS, INSTALL_KNOWLEDGE_CHECKSUM)
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


def emit_owner_handoff(status, *, pack_ref, version, reason, audit_refs=None,
                       readiness=None, candidate_refs=None, owner_steps=None):
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
        "owner_steps": owner_steps or [],
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


def downstream_owner_steps(role_slots_doc, kindex, version):
    """只描述可信 GUI 待办；目标身份来自 Pack 资产，不铸 evidence 或正式事实。"""
    role_targets = []
    role_pack_dir = os.path.join(PACK_DIR, "role-packs")
    for file_name in sorted(os.listdir(role_pack_dir)):
        if not file_name.endswith(".json"):
            continue
        role_pack = load(os.path.join("role-packs", file_name))
        role_targets.append({
            "role_pack_id": role_pack["role_pack_id"],
            "version": role_pack["version"],
        })
    slot_targets = [
        {
            "slot_ref": binding["slot_id"],
            "role_pack_id": binding["role_pack_id"],
            "agent_ref": binding["agent_ref"],
        }
        for binding in sorted(
            role_slots_doc.get("bindings", []),
            key=lambda item: item["slot_id"],
        )
    ]
    knowledge_targets = [
        {
            "source_ref": entry["source_ref"] + "@" + version,
            "knowledge_scope_ref": entry["knowledge_scope_ref"],
        }
        for entry in sorted(
            kindex.get("entries", []),
            key=lambda item: (item["knowledge_scope_ref"], item["source_ref"]),
        )
    ]
    return [
        {
            "step_id": "role_pack_candidate_and_enable",
            "target_count": len(role_targets),
            "targets": role_targets,
            "required_action": "在可信 GUI 审核角色候选；仅用 backend-issued evidence 完成启用。",
        },
        {
            "step_id": "agent_slot_candidate_and_confirm",
            "target_count": len(slot_targets),
            "targets": slot_targets,
            "required_action": "在可信 GUI 审核槽位绑定候选；由 Owner 完成确认。",
        },
        {
            "step_id": "knowledge_candidate_review_and_formalize",
            "target_count": len(knowledge_targets),
            "targets": knowledge_targets,
            "required_action": "在可信 GUI 暂存并审阅知识候选；保持 pending_human_review，不自动正式化。",
        },
    ]


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
    owner_steps = downstream_owner_steps(role_slots_doc, kindex, version)
    print("[6/6] 组合 readiness 通过；下游正式对象交可信 GUI Owner 处理。")
    emit_owner_handoff(
        "awaiting_owner_confirmation",
        pack_ref=pack_ref,
        version=version,
        reason="downstream_owner_confirmation_required",
        audit_refs=readiness["audit_refs"],
        readiness=readiness,
        owner_steps=owner_steps,
    )


if __name__ == "__main__":
    main()
