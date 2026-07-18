#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""智能家居项目经营 Pack：Owner 在场停用交接。

脚本不伪造浏览器 Origin、Cookie、Owner presence 或 Base 决议，也不调用任何
prepare/confirm/disable 写端点。它只展示可信前台交接，并只读等待 os-14
lifecycle ReadModel 证明 Pack 已停用；历史项目与 Receipt 不删除。
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
    UNINSTALL_CONNECTIVITY,
    UNINSTALL_GENERIC,
    UNINSTALL_LIFECYCLE_HTTP,
    emit_pack_error,
    pack_enabled_from_readmodel,
    present_owner_disable_handoff,
    wait_for_owner_disabled,
)

BASE = ""


def call(method, path, body=None):
    if method != "GET" or body is not None:
        raise RuntimeError("uninstall handoff is read-only; Owner writes must come from trusted GUI")
    req = urllib.request.Request(BASE + path, method="GET")
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return resp.status, json.loads(resp.read().decode("utf-8") or "{}")
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
    print("卸载未完成：" + message, file=sys.stderr)
    raise SystemExit(1)


def main():
    parser = argparse.ArgumentParser(description="交给可信 Truzhen 前台由 Owner 显式停用智能家居 Pack。")
    parser.add_argument("--devserver-base", default=os.environ.get("TRUZHEN_DEVSERVER_BASE", "").strip(), help="必填；拒绝猜测隔离实例。")
    parser.add_argument("--client-url", default=os.environ.get("TRUZHEN_CLIENT_URL", ""))
    parser.add_argument("--open-gui", action="store_true", help="显式打开前台；不注入登录态或 Owner presence。")
    parser.add_argument("--wait-seconds", type=float, default=float(os.environ.get("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS", "300")))
    parser.add_argument("--poll-seconds", type=float, default=1.0)
    args = parser.parse_args()

    global BASE
    BASE = args.devserver_base.rstrip("/")
    if not BASE:
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base", UNINSTALL_CONNECTIVITY)

    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
        manifest = json.load(stream)
    pack_ref, version = manifest["pack_ref"], manifest["version"]
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + urllib.parse.quote(pack_ref, safe=""))
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)
    if code != 200:
        die("lifecycle ReadModel HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    enabled = pack_enabled_from_readmodel(body, pack_ref)
    if enabled is False:
        print("场景包已停用或不存在，无需重复操作。")
        return
    if enabled is None:
        die("lifecycle ReadModel 形状不完整，拒绝猜测状态", UNINSTALL_LIFECYCLE_HTTP)

    print("== Owner 在场停用交接 %s @ %s ==" % (pack_ref, version))
    present_owner_disable_handoff(args.client_url, pack_ref, args.open_gui)
    ok, reason = wait_for_owner_disabled(call, pack_ref, args.wait_seconds, args.poll_seconds)
    if not ok:
        code = UNINSTALL_CONNECTIVITY if reason == "connectivity" else UNINSTALL_LIFECYCLE_HTTP
        die("%s；未收到 os-14 已停用真相，不宣称卸载成功" % reason, code)
    print("✅ os-14 已证明 Pack 停用。历史项目、FormalKnowledge 与 03 Receipt 保留可反查。")


if __name__ == "__main__":
    main()
