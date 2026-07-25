#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
墅学家大宅装修设计指导 Pack —— 从正在运行的 Truzhen devserver 卸载（可卸载）。

只消费 Owner 前台 / Base 已签发并外部注入的卸载证明，调用正式 lifecycle/uninstall
端点。脚本不会自行 Prepare/Confirm、也不会把环境变量伪装成主权
裁定。正式卸载会级联停用知识域；历史对象和 Receipt 仍可反查（卸载≠删历史）。

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
from pack_diagnostics import (
    emit_pack_error, UNINSTALL_CONNECTIVITY, UNINSTALL_GENERIC,
    UNINSTALL_LIFECYCLE_HTTP)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")
PROOF_ENV = "TRUZHEN_PACK_UNINSTALL_PROOF_JSON"
# 这是可信前台应签发的 canonical action 描述，不是 Pack 发起的请求体。
OWNER_UNINSTALL_HANDOFF = {"action_type": "14.pack-studio.lifecycle.uninstall"}
UNINSTALL_ACTION = OWNER_UNINSTALL_HANDOFF["action_type"]


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


def load_uninstall_proof(pack_ref, version):
    """只接受已签发的外部证明；本地只做形状和不可变绑定核验。"""
    raw = os.environ.get(PROOF_ENV, "").strip()
    if not raw:
        die("卸载需要可信 Owner 前台签发的 %s；脚本不会自行 Prepare/Confirm。" % PROOF_ENV,
            UNINSTALL_LIFECYCLE_HTTP)
    try:
        proof = json.loads(raw)
    except json.JSONDecodeError:
        die("Owner 卸载证明不是合法 JSON。", UNINSTALL_LIFECYCLE_HTTP)
    if not isinstance(proof, dict):
        die("Owner 卸载证明必须是对象。", UNINSTALL_LIFECYCLE_HTTP)

    expected_transaction = "transaction://pack-uninstall:%s@%s" % (pack_ref, version)
    required = ("decision_ref", "run_id", "nonce", "owner_action_evidence_ref")
    if (proof.get("action_type") != UNINSTALL_ACTION or
            proof.get("target_ref") != pack_ref or
            proof.get("transaction_ref") != expected_transaction or
            any(not isinstance(proof.get(key), str) or not proof[key].strip() for key in required)):
        die("Owner 卸载证明与 Pack、版本、action、transaction 或签发绑定不匹配。",
            UNINSTALL_LIFECYCLE_HTTP)
    return proof


def main():
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as f:
        manifest = json.load(f)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    print("== 卸载 %s @ %s（%s）==" % (pack_ref, version, BASE))
    proof = load_uninstall_proof(pack_ref, version)
    code, body = call("POST", "/v3/pack-studio/lifecycle/uninstall", {
        "pack_ref": pack_ref, "owner_ref": OWNER, "reason": "Owner 卸载墅学家大宅装修设计指导 pack",
        "decision_ref": proof["decision_ref"], "run_id": proof["run_id"], "nonce": proof["nonce"],
        "owner_action_evidence_ref": proof["owner_action_evidence_ref"]})
    if code != 200:
        error_code = UNINSTALL_CONNECTIVITY if code == 0 and "_transport_error" in body else UNINSTALL_LIFECYCLE_HTTP
        die("uninstall HTTP %d: %s" % (code, body), error_code)
    print("\n✅ 卸载成功：%s 已进入 uninstalled，知识域已级联停用。" % pack_ref)
    print("   （已产生的案件对象与 03 回执仍可反查——卸载不删历史。）")


if __name__ == "__main__":
    main()
