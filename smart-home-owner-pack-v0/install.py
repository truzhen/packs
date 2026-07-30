#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""智能家居 Pack 可信 GUI 装入交接：脚本只读 os-14。"""

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
    INSTALL_CONNECTIVITY,
    INSTALL_GENERIC,
    INSTALL_LIFECYCLE_HTTP,
    emit_pack_error,
    pack_enabled_version_from_readmodel,
    present_owner_install_handoff,
    wait_for_owner_enabled,
)

BASE = ""


def call(method, path, body=None):
    """install handoff is read-only: 任何写入只能来自可信 GUI。"""
    if method != "GET" or body is not None:
        raise RuntimeError("install handoff is read-only; Owner writes must come from trusted GUI")
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


def die(message, error_code=INSTALL_GENERIC):
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="install", error_code=error_code, message=message)
    print("装入交接未完成：" + message, file=sys.stderr)
    raise SystemExit(1)


def main():
    parser = argparse.ArgumentParser(description="把智能家居 Pack 装入交给可信 Truzhen GUI。")
    parser.add_argument("--devserver-base", default=os.environ.get("TRUZHEN_DEVSERVER_BASE", "").strip())
    parser.add_argument("--client-url", default=os.environ.get("TRUZHEN_CLIENT_URL", ""))
    parser.add_argument("--open-gui", action="store_true", help="仅显式打开前台；不注入登录态或 Owner presence。")
    parser.add_argument("--wait-seconds", type=float, default=float(os.environ.get("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS", "300")))
    parser.add_argument("--poll-seconds", type=float, default=1.0)
    args = parser.parse_args()

    global BASE
    BASE = args.devserver_base.rstrip("/")
    if not BASE:
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base", INSTALL_CONNECTIVITY)

    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
        manifest = json.load(stream)
    pack_ref, version = manifest["pack_ref"], manifest["version"]
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + urllib.parse.quote(pack_ref, safe=""))
    if code == 0:
        die("连不上 devserver（%s）" % BASE, INSTALL_CONNECTIVITY)
    if code != 200:
        die("lifecycle ReadModel HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    enabled_version = pack_enabled_version_from_readmodel(body, pack_ref)
    if enabled_version is None:
        die("lifecycle ReadModel 形状不完整，拒绝猜测状态", INSTALL_LIFECYCLE_HTTP)
    if enabled_version and enabled_version != version:
        die("已启用版本 %s 与声明版本 %s 不一致，拒绝覆盖" % (enabled_version, version), INSTALL_LIFECYCLE_HTTP)
    if enabled_version != version:
        present_owner_install_handoff(args.client_url, pack_ref, version, open_gui=args.open_gui)
        ok, reason = wait_for_owner_enabled(call, pack_ref, version, args.wait_seconds, args.poll_seconds)
        if not ok:
            die("%s；未观察到 os-14 的精确版本状态" % reason, INSTALL_CONNECTIVITY if reason == "connectivity" else INSTALL_LIFECYCLE_HTTP)
    print("os-14 已证明精确 Pack 版本启用。")


if __name__ == "__main__":
    main()
