#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""通过 Truzhen 既有 lifecycle 装入内容运营工作台 Pack 与声明式计划。

本脚本只提交 Pack/角色/槽位声明；不启动 Codex CLI，不读取凭据，不执行内容生产，
也不登录、上传或发布到任何平台。Provider readiness 不通过时必须停止。
"""
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

from pack_diagnostics import (  # noqa: E402
    INSTALL_CONNECTIVITY,
    INSTALL_GENERIC,
    INSTALL_LIFECYCLE_HTTP,
    INSTALL_READINESS,
    INSTALL_ROLE_BINDING,
    INSTALL_STATE_CONFLICT,
    emit_pack_error,
)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")


def load(rel):
    with open(os.path.join(PACK_DIR, rel), encoding="utf-8") as handle:
        return json.load(handle)


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    request = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        request.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(request, timeout=120) as response:
            raw = response.read().decode("utf-8")
            code = response.status
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8")
        code = exc.code
    except Exception as exc:  # transport boundary; the exact error is returned to diagnostics
        return 0, {"_transport_error": str(exc)}
    try:
        return code, json.loads(raw) if raw else {}
    except json.JSONDecodeError:
        return code, {"_raw": raw}


def die(message, error_code=INSTALL_GENERIC):
    emit_pack_error(
        pack_dir=PACK_DIR,
        base=BASE,
        action="install",
        error_code=error_code,
        message=message,
    )
    print("装入失败：" + message, file=sys.stderr)
    raise SystemExit(1)


def lifecycle_snapshot(pack_ref):
    query = urllib.parse.urlencode({"pack_ref": pack_ref})
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?" + query)
    if code == 0:
        die("连不上 devserver（%s）" % BASE, INSTALL_CONNECTIVITY)
    if code != 200:
        die("读取 lifecycle HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    for item in body.get("packs", []) or []:
        if item.get("pack_ref") == pack_ref:
            current = (item.get("enabled_pointer") or {}).get("current_version", "")
            records = item.get("records", []) or []
            if not records:
                return current, "", "", 0
            latest = max(records, key=lambda record: int(record.get("occ_version", 0)))
            return current, latest.get("version", ""), latest.get("state", ""), int(latest.get("occ_version", 0))
    return "", "", "", 0


def install_scene(manifest, flow, role_slots, capabilities):
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    pack_version_ref = pack_ref + "@" + version
    current, lifecycle_version, lifecycle_state, lifecycle_occ = lifecycle_snapshot(pack_ref)
    if current:
        if current != version:
            die(
                "已有不同版本 %s 处于 enabled；本脚本不自动升级或替换" % current,
                INSTALL_STATE_CONFLICT,
            )
        print("[1-5/7] 场景包已启用 @%s，幂等跳过。" % current)
        return pack_version_ref

    if lifecycle_state == "disabled" and lifecycle_version == version:
        code, body = call(
            "POST",
            "/v3/pack-studio/lifecycle/reactivate",
            {
                "pack_ref": pack_ref,
                "version": version,
                "owner_ref": OWNER,
                "idempotency_key": "content-ops-pack-reactivate:%s:%d" % (pack_version_ref, lifecycle_occ),
                "comment": manifest["name"] + "重新启用",
            },
        )
        if code == 200 and body.get("status") == "enabled":
            print("[1-5/7] 同版本 reactivate 成功。")
            return pack_version_ref
        die("disabled 版本 reactivate 失败 HTTP %d: %s" % (code, body), INSTALL_STATE_CONFLICT)
    if lifecycle_state and lifecycle_version == version:
        die("已有非 enabled/disabled lifecycle 记录，拒绝覆盖：%s" % lifecycle_state, INSTALL_STATE_CONFLICT)

    print("[1/7] 画布写穿 06 ...")
    canvas = {
        "flow_id": flow["flow_id"],
        "title": flow.get("title", ""),
        "occ_version": 0,
        "save_source": "pack_install",
        "flow_spec_draft": flow,
    }
    code, body = call("POST", "/v3/pack-studio/canvas", canvas)
    if code == 409 and isinstance(body.get("current_occ_version"), (int, float)):
        canvas["occ_version"] = int(body["current_occ_version"])
        code, body = call("POST", "/v3/pack-studio/canvas", canvas)
    if code != 200 or not ((body.get("engine_sync") or {}).get("synced")):
        die("canvas 未真实同步 HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

    print("[2/7] 提交 lifecycle draft ...")
    provider_requirements = []
    for requirement in capabilities.get("provider_requirements", []):
        provider_requirements.append(
            {
                key: requirement[key]
                for key in (
                    "requirement_id",
                    "capability",
                    "gateway_class",
                    "risk_class",
                    "fallback_policy",
                    "provider_family",
                )
                if key in requirement
            }
        )
    routes = manifest.get("notification_command_report_routes", {})
    draft = {
        "pack_ref": pack_ref,
        "version": version,
        "title": manifest["name"],
        "template_family": manifest["template_family"],
        "flow_id": flow["flow_id"],
        "role_slots": [
            {
                "slot_id": slot["slot_id"],
                "responsibility": slot["responsibility"],
                "required_role": slot["required_role"],
                "default_role_pack_ref": slot["default_role_pack_ref"],
                "slice_scope_policy": slot.get("slice_scope_policy", ""),
                "node_type": slot.get("node_type", ""),
            }
            for slot in role_slots["role_slots"]
        ],
        "provider_requirements": provider_requirements,
        "formalization_requirement": manifest["formalization_requirement"]["summary"],
        "notification_routes": routes.get("notification", []),
        "command_candidates": routes.get("command_candidate", []),
        "report_routes": routes.get("report", []),
        "moat_justification": manifest["moat_justification"],
        "knowledge_scopes": [],
        "risk_types": manifest.get("risk_types", []),
        "idempotency_key": "content-ops-pack-draft:" + pack_version_ref,
        "actor_ref": OWNER,
    }
    code, body = call("POST", "/v3/pack-studio/lifecycle/draft", draft)
    if code != 200:
        die("draft HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

    print("[3/7] 检查 readiness ...")
    code, body = call(
        "POST",
        "/v3/pack-studio/lifecycle/readiness",
        {"pack_ref": pack_ref, "version": version, "actor_ref": OWNER},
    )
    if code != 200:
        die("readiness HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    readiness = (body.get("record") or {}).get("readiness_report") or {}
    if not readiness.get("ready"):
        die("readiness 未通过：%s" % body, INSTALL_READINESS)

    print("[4/7] promote ...")
    code, body = call(
        "POST",
        "/v3/pack-studio/lifecycle/promote",
        {"pack_ref": pack_ref, "version": version, "actor_ref": OWNER},
    )
    if code != 200:
        die("promote HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

    print("[5/7] confirm（由 01 Base 签发并落 03 Receipt）...")
    code, body = call(
        "POST",
        "/v3/pack-studio/lifecycle/confirm",
        {
            "pack_ref": pack_ref,
            "version": version,
            "idempotency_key": "content-ops-pack-confirm:" + pack_version_ref,
            "owner_ref": OWNER,
            "approve": True,
            "comment": manifest["name"] + "启用",
        },
    )
    if code != 200 or body.get("status") != "enabled":
        die("confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    return pack_version_ref


def install_role_pack(role_pack):
    role_id = role_pack["role_pack_id"]
    safe_id = re.sub(r"[/:@\s]+", "-", role_id).strip("-")
    draft_id = "role-pack-draft-content-ops-" + safe_id
    idem = "content-ops-role-" + safe_id
    style = role_pack.get("communication_style", {})
    code, body = call(
        "POST",
        "/v3/agent-orchestration/role-packs/drafts",
        {
            "draft_id": draft_id,
            "role_pack_id": role_id,
            "version": role_pack["version"],
            "display_name": role_pack["display_name"],
            "description": role_pack["description"],
            "target_use": role_pack.get("target_use", []),
            "style_summary": role_pack.get("style_summary", ""),
            "decision_style_summary": role_pack.get("decision_style_summary", ""),
            "model_policy_ref": role_pack.get("model_policy_ref", ""),
            "communication_style": {
                "structure": style.get("structure", ""),
                "tone": style.get("tone", ""),
                "forbidden_phrases": style.get("forbidden_phrases", []),
            },
            "forbidden_behavior_policy_ref": role_pack.get("forbidden_policy_ref", ""),
            "risk_level": role_pack.get("risk_level", "medium"),
            "owner_ref": OWNER,
            "evidence_refs": ["evidence://pack-asset/" + safe_id],
            "idempotency_key": idem + "-draft",
        },
    )
    if code != 200:
        die("role draft HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    draft_id = body.get("draft_id") or (body.get("draft") or {}).get("draft_id") or draft_id
    code, body = call(
        "POST",
        "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/readiness-check",
        {},
    )
    if code != 200:
        die("role readiness HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call(
        "POST",
        "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/promote-candidate",
        {
            "owner_ref": OWNER,
            "target_agent_ref": role_id,
            "idempotency_key": idem + "-promote",
            "evidence_refs": ["evidence://pack-asset/role-promote"],
        },
    )
    if code != 200 or not body.get("ok"):
        die("role promote HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call(
        "POST",
        "/v3/agent-orchestration/role-packs/enable-candidate",
        {
            "role_pack_id": role_id,
            "version": role_pack["version"],
            "target_agent_ref": role_id,
            "owner_ref": OWNER,
            "idempotency_key": idem + "-enable",
            "evidence_refs": ["evidence://pack-asset/role-enable"],
        },
    )
    if code != 200 or not body.get("ok"):
        die("role enable candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call(
        "POST",
        "/v3/agent-orchestration/role-packs/enable-confirm",
        {
            "role_pack_id": role_id,
            "version": role_pack["version"],
            "target_agent_ref": role_id,
            "owner_ref": OWNER,
            "idempotency_key": idem + "-confirm",
            "approve": True,
            "comment": role_pack["display_name"] + "启用",
            "evidence_refs": ["evidence://pack-asset/role-enable-confirm"],
        },
    )
    if code != 200 or body.get("status") != "enabled":
        die("role confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def install_roles_and_bindings(role_slots, pack_version_ref):
    code, body = call("GET", "/v3/agent-orchestration/role-packs/readmodel")
    if code != 200:
        die("role readmodel HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    enabled = {
        item.get("role_pack_id", "").split("@")[0]
        for item in body.get("enabled_versions", []) or []
    }
    for filename in sorted(os.listdir(os.path.join(PACK_DIR, "role-packs"))):
        if not filename.endswith(".json"):
            continue
        role_pack = load(os.path.join("role-packs", filename))
        if role_pack["role_pack_id"] not in enabled:
            install_role_pack(role_pack)

    code, body = call("GET", "/v3/agent-orchestration/agent-slots/readmodel")
    if code != 200:
        die("slot readmodel HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    enabled_bindings = {
        (item.get("slot_ref", ""), item.get("scope_ref", ""))
        for item in body.get("agent_slot_bindings", []) or []
        if item.get("enabled_state") == "enabled"
    }
    for binding in role_slots.get("bindings", []):
        binding_key = (binding["slot_id"], pack_version_ref)
        if binding_key in enabled_bindings:
            print("    槽位 %s 已绑定，跳过" % binding["slot_id"])
            continue
        code, body = call(
            "POST",
            "/v3/agent-orchestration/agent-slots/bind-candidate",
            {
                "slot_ref": binding["slot_id"],
                "scope_ref": pack_version_ref,
                "source_pack_ref": pack_version_ref,
                "required_role": binding["required_role"],
                "requested_agent_ref": binding["agent_ref"],
                "requested_role_pack_id": binding["role_pack_id"],
                "ttl": "8760h",
                "evidence_refs": ["evidence://pack-asset/slot/" + binding["slot_id"]],
            },
        )
        if code != 200 or not body.get("ok"):
            die("slot bind candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
        code, body = call(
            "POST",
            "/v3/agent-orchestration/agent-slots/confirm",
            {
                "binding_ref": body.get("binding_ref"),
                "idempotency_key": "content-ops-slot-confirm:" + binding["slot_id"],
                "approve": True,
                "evidence_refs": ["evidence://pack-asset/slot-confirm/" + binding["slot_id"]],
            },
        )
        if code != 200 or body.get("status") != "enabled":
            die("slot confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def schedule_transaction_ref(pack_ref, schedule_key):
    slug = pack_ref.removeprefix("scene_pack://")
    return "transaction://" + slug + "/schedule/" + schedule_key


def install_schedules(manifest, schedule_doc, pack_version_ref):
    code, body = call("GET", "/v3/task-governance/schedules")
    if code != 200:
        die("schedule readmodel HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    existing = {
        item.get("transaction_ref", ""): item
        for item in body.get("schedules", []) or []
        if item.get("transaction_ref")
    }
    for declaration in schedule_doc.get("schedules", []):
        key = declaration["schedule_key"]
        transaction_ref = schedule_transaction_ref(manifest["pack_ref"], key)
        schedule = existing.get(transaction_ref)
        if schedule:
            status = schedule.get("status")
            if status == "active":
                print("    计划 %s 已 active，跳过" % key)
                continue
            if status == "cancelled":
                die("计划 %s 已 cancelled；不可自动复活，需发布新 Pack 版本" % key, INSTALL_STATE_CONFLICT)
            if status == "paused":
                code, response = call(
                    "POST",
                    "/v3/task-governance/schedules/resume",
                    {
                        "schedule_id": schedule["id"],
                        "idempotency_key": "content-ops-schedule-resume:%s:%s" % (
                            schedule["id"],
                            schedule.get("occ_version", 0),
                        ),
                        "reason": "重新启用" + manifest["name"],
                    },
                )
                if code != 200:
                    die("schedule resume HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
                print("    计划 %s 已恢复" % key)
                continue
            die("计划 %s 状态未知：%s" % (key, status), INSTALL_STATE_CONFLICT)

        source_event_id = "pack_schedule_install:%s:%s" % (pack_version_ref, key)
        event = {
            "event_id": source_event_id,
            "intent_ref": pack_version_ref,
            "transaction_ref": transaction_ref,
            "source_event_id": source_event_id,
            "source": "intent_event",
            "task_type": "scheduled",
            "spec": {
                "name": declaration["display_name"],
                "description": declaration["description"],
                "action": declaration["skill_id"],
                "params": {
                    "cron": declaration["cron"],
                    "template_risk_level": declaration["template_risk_level"],
                    "risk_ceiling": declaration["risk_ceiling"],
                    "misfire_policy": declaration["misfire_policy"],
                    "standing_approval_scope": "%s:%s" % (pack_version_ref, declaration["skill_id"]),
                    "trigger_task_type": declaration["trigger_task_type"],
                },
            },
            "risk_level": "low",
            "risk_reason": "Pack-owned schedule declaration; each Hands run remains separately gated",
            "evidence_refs": [
                "evidence://pack-asset/content-operations/schedule/" + key,
                "evidence://pack-version/" + pack_version_ref,
            ],
            "idempotency_key": "content-ops-schedule-candidate:" + pack_version_ref + ":" + key,
            "reason": declaration["local_time_note"],
        }
        code, response = call("POST", "/v3/task-governance/candidates/intake", event)
        if code not in (200, 201):
            die("schedule intake HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
        candidate = response.get("candidate") or {}
        candidate_id = candidate.get("id")
        if not candidate_id:
            die("schedule intake 未返回 candidate id: %s" % response, INSTALL_LIFECYCLE_HTTP)
        state = candidate.get("state")
        if state == "candidate":
            code, response = call(
                "POST",
                "/v3/task-governance/candidates/submit-review",
                {"candidate_id": candidate_id},
            )
            if code != 200:
                die("schedule submit-review HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
        elif state != "reviewing":
            die("schedule candidate 状态不可批准：%s" % state, INSTALL_STATE_CONFLICT)
        code, response = call(
            "POST",
            "/v3/task-governance/schedules/approve",
            {
                "candidate_id": candidate_id,
                "idempotency_key": "content-ops-schedule-approve:" + candidate_id,
                "reason": "启用" + declaration["display_name"],
            },
        )
        if code != 200 or (response.get("schedule") or {}).get("status") != "active":
            die("schedule approve 未达 active HTTP %d: %s" % (code, response), INSTALL_LIFECYCLE_HTTP)
        print("    计划 %s 已 active" % key)


def main():
    manifest = load("manifest.json")
    flow = load(manifest["flow_file"])
    role_slots = load(manifest["role_slots_file"])
    capabilities = load(manifest["capabilities_file"])
    schedules = load(manifest["schedules_file"])
    print("== 装入 %s @ %s 到 %s ==" % (manifest["pack_ref"], manifest["version"], BASE))
    pack_version_ref = install_scene(manifest, flow, role_slots, capabilities)
    print("[6/7] 装入角色包并绑定 Pack 槽位 ...")
    install_roles_and_bindings(role_slots, pack_version_ref)
    print("[7/7] 经 07 + 01 + 03 装入 Pack 计划 ...")
    install_schedules(manifest, schedules, pack_version_ref)
    print("装入成功：%s 已 enabled；外部发布能力仍不存在。" % pack_version_ref)


if __name__ == "__main__":
    main()
