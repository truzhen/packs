#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
家政运营 Pack（客户服务全生命周期）—— 从正在运行的 Truzhen devserver 卸载。

走真实 lifecycle/disable 端点：先经 Base gated-action prepare→confirm 取得真签发的
decision_ref/run_id/nonce（禁自铸），再 disable。卸载只停用 Pack 当前版本；
已产生的事务对象、候选和回执在 03 仍可反查，卸载不删除历史。

用法：
  TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 python3 housekeeping-ops-pack-v0/uninstall.py
"""
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
    emit_pack_error, UNINSTALL_GENERIC, UNINSTALL_CONNECTIVITY, UNINSTALL_LIFECYCLE_HTTP)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            raw = resp.read().decode("utf-8")
            return resp.status, json.loads(raw) if raw else {}
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
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as f:
        manifest = json.load(f)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    pack_name = manifest.get("name", "家政运营 Pack")
    print("== 卸载 %s @ %s（%s）==" % (pack_ref, version, BASE))

    query = urllib.parse.urlencode({"pack_ref": pack_ref})
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?" + query)
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)

    enabled = False
    for entry in body.get("packs", []) or []:
        if entry.get("pack_ref") == pack_ref:
            ptr = entry.get("enabled_pointer") or {}
            if ptr.get("current_version"):
                enabled = True
                break
    if not enabled:
        print("场景包未处于启用态，无需卸载。")
        return

    code, body = call("POST", "/v3/base/gated-actions/prepare", {
        "action_type": "14.pack-studio.lifecycle.disable",
        "target_ref": pack_ref,
        "content_summary": "停用家政运营 Pack：" + pack_ref,
        "impact_summary": "停用客户服务全生命周期场景包；历史事务对象、候选和回执保留可反查",
        "transaction_ref": "transaction://pack-disable:" + pack_ref + "@" + version,
    })
    if code != 200:
        die("gated prepare HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    issue_ref = (body.get("issue") or {}).get("issue_ref")
    if not issue_ref:
        die("gated prepare 无 issue_ref: %s" % body, UNINSTALL_LIFECYCLE_HTTP)

    code, body = call("POST", "/v3/base/gated-actions/confirm", {"issue_ref": issue_ref})
    if code != 200:
        die("gated confirm HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    issue = body.get("issue") or {}
    if not all(issue.get(key) for key in ("decision_ref", "run_id", "nonce", "owner_action_evidence_ref")):
        die("gated confirm 返回的 issued binding 不完整: %s" % body, UNINSTALL_LIFECYCLE_HTTP)

    code, body = call("POST", "/v3/pack-studio/lifecycle/disable", {
        "pack_ref": pack_ref,
        "owner_ref": OWNER,
        "reason": "Owner 卸载" + pack_name,
        "decision_ref": issue["decision_ref"],
        "run_id": issue.get("run_id", ""),
        "nonce": issue.get("nonce", ""),
        "owner_action_evidence_ref": issue["owner_action_evidence_ref"],
    })
    if code != 200:
        die("disable HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)

    print("\n卸载成功：%s 已停用。" % pack_ref)
    print("已产生的事务对象、候选和 03 回执仍可反查；卸载不删除历史。")


if __name__ == "__main__":
    main()
