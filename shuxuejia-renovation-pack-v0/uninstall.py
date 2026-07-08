#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
墅学家大宅装修设计指导 Pack —— 从正在运行的 Truzhen devserver 卸载（可卸载）。

走真实 lifecycle/disable 端点：先经 Base gated-action prepare→confirm 取得真签发的
decision_ref/run_id/nonce（停用与启用同纪律，禁自铸），再 disable。pack 停用会经
driveDisableMounts 级联卸载知识域。已产生的案件对象与回执在 03 仍可反查（卸载≠删历史）。

用法：
  python3 packs/shuxuejia-renovation-pack-v0/uninstall.py
"""
import json
import os
import sys
import urllib.error
import urllib.request

PACK_DIR = os.path.dirname(os.path.abspath(__file__))
REPO_DIR = os.path.dirname(PACK_DIR)
if REPO_DIR not in sys.path:
    sys.path.insert(0, REPO_DIR)
from pack_diagnostics import emit_pack_error

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
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


def die(msg):
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="uninstall", error_code="TZ-PACK-UNINSTALL-001", message=msg)
    print("卸载失败：" + msg, file=sys.stderr)
    sys.exit(1)


def main():
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as f:
        manifest = json.load(f)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    print("== 卸载 %s @ %s（%s）==" % (pack_ref, version, BASE))

    # 已是停用/不存在则幂等返回
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + pack_ref)
    if code == 0:
        die("连不上 devserver（%s）" % BASE)
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
        "action_type": "pack_disable", "target_ref": pack_ref,
        "content_summary": "停用墅学家大宅装修设计指导 pack：" + pack_ref,
        "impact_summary": "停用墅学家场景包并级联卸载其知识域；历史回执保留可反查",
        "transaction_ref": "transaction://pack-disable:" + pack_ref + "@" + version})
    if code != 200:
        die("gated prepare HTTP %d: %s" % (code, body))
    issue_ref = (body.get("issue") or {}).get("issue_ref")
    if not issue_ref:
        die("gated prepare 无 issue_ref: %s" % body)
    code, body = call("POST", "/v3/base/gated-actions/confirm", {"issue_ref": issue_ref})
    if code != 200:
        die("gated confirm HTTP %d: %s" % (code, body))
    iss = body.get("issue") or {}
    if not iss.get("decision_ref"):
        die("gated confirm 无 decision_ref: %s" % body)

    # 2. disable
    code, body = call("POST", "/v3/pack-studio/lifecycle/disable", {
        "pack_ref": pack_ref, "owner_ref": OWNER, "reason": "Owner 卸载墅学家大宅装修设计指导 pack",
        "decision_ref": iss["decision_ref"], "run_id": iss.get("run_id", ""), "nonce": iss.get("nonce", ""),
        "owner_action_evidence_ref": "owner_action_evidence://pack-disable/" + pack_ref})
    if code != 200:
        die("disable HTTP %d: %s" % (code, body))
    print("\n✅ 卸载成功：%s 已停用，知识域已级联卸载。前端「场景包管理」刷新即消失。" % pack_ref)
    print("   （已产生的案件对象与 03 回执仍可反查——卸载不删历史。）")


if __name__ == "__main__":
    main()
