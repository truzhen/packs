#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
环保执法 Pack —— 从正在运行的 Truzhen devserver 受控卸载（可逆停用）。

走真实 lifecycle/disable 端点：先经 Base gated-action prepare→confirm 取得真签发的
decision_ref/run_id/nonce（停用与启用同纪律，禁自铸），再 disable。pack 停用会经
driveDisableMounts 级联卸载知识域。该 lifecycle 语义是“卸载运行访问权”：
保留本地已安装版本元数据，以便后续 install.py 幂等重新启用；已产生的案件对象与
回执在 03 仍可反查（卸载≠物理删除历史）。

用法：
  TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 \
    python3 packs/environmental-enforcement-pack-v0/uninstall.py
"""
import argparse
import json
import os
import sys
import urllib.error
import urllib.request

PACK_DIR = os.path.dirname(os.path.abspath(__file__))
REPO_DIR = os.path.dirname(PACK_DIR)
if REPO_DIR not in sys.path:
    sys.path.insert(0, REPO_DIR)
from pack_diagnostics import (
    emit_pack_error, UNINSTALL_GENERIC, UNINSTALL_CONNECTIVITY, UNINSTALL_LIFECYCLE_HTTP)

# 由 main() 在参数解析后赋值。lifecycle 是写操作，绝不以隐式 localhost 默认值
# 猜测目标实例，避免测试或运维命令越出已登记的隔离端口。
BASE = ""
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            return resp.status, json.loads(resp.read().decode("utf-8") or "{}")
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8")
        try:
            return e.code, json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            return e.code, {"_raw": raw}
    except Exception as e:
        return 0, {"_transport_error": str(e)}


def die(msg, error_code=UNINSTALL_GENERIC):
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="uninstall", error_code=error_code, message=msg)
    print("卸载失败：" + msg, file=sys.stderr)
    sys.exit(1)


def main():
    parser = argparse.ArgumentParser(
        description="通过受控 devserver 停用环保执法 Pack；不会删除历史项目或 Receipt。"
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
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base；拒绝猜测默认端口", UNINSTALL_CONNECTIVITY)
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as f:
        manifest = json.load(f)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    print("== 卸载 %s @ %s（%s）==" % (pack_ref, version, BASE))

    # 已是停用/不存在则幂等返回
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + pack_ref)
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)
    enabled = False
    for entry in (body.get("packs", []) or []):
        if entry.get("pack_ref") == pack_ref:
            ptr = entry.get("enabled_pointer") or {}
            if ptr.get("current_version"):
                enabled = True
    if not enabled:
        print("场景包未处于启用态，无需卸载。")
        return

    # 1. Base 签发 disable decision
    code, body = call("POST", "/v3/base/gated-actions/prepare", {
        "action_type": "14.pack-studio.lifecycle.disable", "target_ref": pack_ref,
        "content_summary": "停用环保执法 pack：" + pack_ref,
        "impact_summary": "停用场景包并级联卸载其知识域；历史回执保留可反查",
        "transaction_ref": "transaction://pack-disable:" + pack_ref + "@" + version})
    if code != 200:
        die("gated prepare HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    issue_ref = (body.get("issue") or {}).get("issue_ref")
    if not issue_ref:
        die("gated prepare 无 issue_ref: %s" % body, UNINSTALL_LIFECYCLE_HTTP)
    code, body = call("POST", "/v3/base/gated-actions/confirm", {"issue_ref": issue_ref})
    if code != 200:
        die("gated confirm HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    iss = body.get("issue") or {}
    if not all(iss.get(key) for key in ("decision_ref", "run_id", "nonce", "owner_action_evidence_ref")):
        die("gated confirm 返回的 issued binding 不完整: %s" % body, UNINSTALL_LIFECYCLE_HTTP)

    # 2. disable
    code, body = call("POST", "/v3/pack-studio/lifecycle/disable", {
        "pack_ref": pack_ref, "owner_ref": OWNER, "reason": "Owner 卸载环保执法 pack",
        "decision_ref": iss["decision_ref"], "run_id": iss.get("run_id", ""), "nonce": iss.get("nonce", ""),
        "owner_action_evidence_ref": iss["owner_action_evidence_ref"]})
    if code != 200:
        die("disable HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    print("\n✅ 受控卸载成功：%s 已停用，知识域已级联卸载。" % pack_ref)
    print("   场景包管理会保留“本地已安装 / 未启用”的版本元数据，供 install.py 幂等重装；")
    print("   已产生的案件对象与 03 回执仍可反查——卸载不物理删除历史。")


if __name__ == "__main__":
    main()
