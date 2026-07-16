#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""通过 Base gated action 和 Pack lifecycle 停用内容运营工作台 Pack。"""
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

from pack_diagnostics import (  # noqa: E402
    UNINSTALL_CONNECTIVITY,
    UNINSTALL_GENERIC,
    UNINSTALL_LIFECYCLE_HTTP,
    emit_pack_error,
)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    request = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        request.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(request, timeout=60) as response:
            raw = response.read().decode("utf-8")
            code = response.status
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8")
        code = exc.code
    except Exception as exc:  # transport boundary; preserve the exact diagnostic
        return 0, {"_transport_error": str(exc)}
    try:
        return code, json.loads(raw) if raw else {}
    except json.JSONDecodeError:
        return code, {"_raw": raw}


def die(message, error_code=UNINSTALL_GENERIC):
    emit_pack_error(
        pack_dir=PACK_DIR,
        base=BASE,
        action="uninstall",
        error_code=error_code,
        message=message,
    )
    print("卸载失败：" + message, file=sys.stderr)
    raise SystemExit(1)


def pause_pack_schedules(pack_ref, pack_name):
    prefix = "transaction://" + pack_ref.removeprefix("scene_pack://") + "/schedule/"
    code, body = call("GET", "/v3/task-governance/schedules")
    if code != 200:
        die("schedule readmodel HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    for schedule in body.get("schedules", []) or []:
        if not str(schedule.get("transaction_ref", "")).startswith(prefix):
            continue
        if schedule.get("status") in ("paused", "cancelled"):
            continue
        if schedule.get("status") != "active":
            die("Pack 计划状态未知：%s" % schedule, UNINSTALL_LIFECYCLE_HTTP)
        code, response = call(
            "POST",
            "/v3/task-governance/schedules/pause",
            {
                "schedule_id": schedule["id"],
                "idempotency_key": "content-ops-schedule-pause:%s:%s" % (
                    schedule["id"],
                    schedule.get("occ_version", 0),
                ),
                "reason": "停用" + pack_name,
            },
        )
        if code != 200:
            die("schedule pause HTTP %d: %s" % (code, response), UNINSTALL_LIFECYCLE_HTTP)


def main():
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as handle:
        manifest = json.load(handle)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    query = urllib.parse.urlencode({"pack_ref": pack_ref})
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?" + query)
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)
    if code != 200:
        die("读取 lifecycle HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    current = ""
    for item in body.get("packs", []) or []:
        if item.get("pack_ref") == pack_ref:
            current = (item.get("enabled_pointer") or {}).get("current_version", "")
    pause_pack_schedules(pack_ref, manifest["name"])
    if not current:
        print("场景包未启用；Pack 计划已确认非 active。")
        return

    code, body = call(
        "POST",
        "/v3/base/gated-actions/prepare",
        {
            "action_type": "pack_disable",
            "target_ref": pack_ref,
            "content_summary": "停用内容运营工作台 Pack",
            "impact_summary": "停止新候选 run；历史候选、Evidence 与 Receipt 保留",
            "transaction_ref": "transaction://pack-disable:" + pack_ref + "@" + current,
        },
    )
    if code != 200:
        die("gated prepare HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    issue_ref = (body.get("issue") or {}).get("issue_ref")
    if not issue_ref:
        die("gated prepare 无 issue_ref: %s" % body, UNINSTALL_LIFECYCLE_HTTP)
    code, body = call("POST", "/v3/base/gated-actions/confirm", {"issue_ref": issue_ref})
    if code != 200:
        die("gated confirm HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    issue = body.get("issue") or {}
    required = ("decision_ref", "run_id", "nonce", "owner_action_evidence_ref")
    missing = [field for field in required if not issue.get(field)]
    if missing:
        die("Base confirm 缺少签发证明 %s: %s" % (missing, body), UNINSTALL_LIFECYCLE_HTTP)

    code, body = call(
        "POST",
        "/v3/pack-studio/lifecycle/disable",
        {
            "pack_ref": pack_ref,
            "owner_ref": OWNER,
            "reason": "Owner 停用" + manifest["name"],
            "decision_ref": issue["decision_ref"],
            "run_id": issue["run_id"],
            "nonce": issue["nonce"],
            "owner_action_evidence_ref": issue["owner_action_evidence_ref"],
        },
    )
    if code != 200:
        die("disable HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    print("停用成功：%s；历史候选、Evidence 与 Receipt 未删除。" % pack_ref)


if __name__ == "__main__":
    main()
