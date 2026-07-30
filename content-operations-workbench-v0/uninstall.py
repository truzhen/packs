#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""内容运营 Pack 可信 GUI 停用交接：脚本只读 os-14 与 os-07。"""

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

from pack_diagnostics import (  # noqa: E402
    UNINSTALL_CONNECTIVITY,
    UNINSTALL_GENERIC,
    UNINSTALL_LIFECYCLE_HTTP,
    emit_pack_error,
    pack_enabled_from_readmodel,
    present_owner_disable_handoff,
    schedule_transaction_ref,
    wait_for_owner_disabled,
    wait_for_owner_schedule_states,
)

BASE = ""
OWNER_DISABLE_HANDOFF = {"action_type": "14.pack-studio.lifecycle.disable"}


def call(method, path, body=None):
    """uninstall handoff is read-only: 任何写入只能来自可信 GUI。"""
    if method != "GET" or body is not None:
        raise RuntimeError("uninstall handoff is read-only; Owner writes must come from trusted GUI")
    request = urllib.request.Request(BASE + path, method="GET")
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            return response.status, json.loads(response.read().decode("utf-8") or "{}")
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8")
        try:
            return exc.code, json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            return exc.code, {"_raw": raw}
    except Exception as exc:
        return 0, {"_transport_error": str(exc)}


def die(message, error_code=UNINSTALL_GENERIC):
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="uninstall", error_code=error_code, message=message)
    print("停用交接未完成：" + message, file=sys.stderr)
    raise SystemExit(1)


def load_json(relative_path):
    with open(os.path.join(PACK_DIR, relative_path), encoding="utf-8") as stream:
        return json.load(stream)


def main():
    parser = argparse.ArgumentParser(description="把内容运营 Pack 停用交给可信 Truzhen GUI。")
    parser.add_argument("--devserver-base", default=os.environ.get("TRUZHEN_DEVSERVER_BASE", "").strip())
    parser.add_argument("--client-url", default=os.environ.get("TRUZHEN_CLIENT_URL", ""))
    parser.add_argument("--open-gui", action="store_true", help="仅显式打开前台；不注入登录态或 Owner presence。")
    parser.add_argument("--wait-seconds", type=float, default=float(os.environ.get("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS", "300")))
    parser.add_argument("--poll-seconds", type=float, default=1.0)
    args = parser.parse_args()

    global BASE
    BASE = args.devserver_base.rstrip("/")
    if not BASE:
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base", UNINSTALL_CONNECTIVITY)

    manifest = load_json("manifest.json")
    schedule_doc = load_json(manifest["schedules_file"])
    pack_ref = manifest["pack_ref"]
    schedule_refs = [
        schedule_transaction_ref(pack_ref, item["schedule_key"])
        for item in schedule_doc.get("schedules", [])
    ]
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + urllib.parse.quote(pack_ref, safe=""))
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)
    if code != 200:
        die("lifecycle ReadModel HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    enabled = pack_enabled_from_readmodel(body, pack_ref)
    if enabled is None:
        die("lifecycle ReadModel 形状不完整，拒绝猜测状态", UNINSTALL_LIFECYCLE_HTTP)

    if enabled:
        present_owner_disable_handoff(args.client_url, pack_ref, args.open_gui, len(schedule_refs))
        ok, reason = wait_for_owner_disabled(call, pack_ref, args.wait_seconds, args.poll_seconds)
        if not ok:
            die("%s；未观察到 os-14 的停用状态" % reason, UNINSTALL_CONNECTIVITY if reason == "connectivity" else UNINSTALL_LIFECYCLE_HTTP)

    ok, reason = wait_for_owner_schedule_states(call, schedule_refs, {"paused", "cancelled"}, args.wait_seconds, args.poll_seconds, allow_missing=True)
    if not ok:
        die("%s；仍未观察到 os-07 声明计划停止" % reason, UNINSTALL_CONNECTIVITY if reason == "connectivity" else UNINSTALL_LIFECYCLE_HTTP)
    print("os-14 已证明 Pack 停用；os-07 已证明声明计划非 active。")


if __name__ == "__main__":
    main()
