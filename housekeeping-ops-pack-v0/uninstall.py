#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Formally uninstall the housekeeping Scene Pack from a running devserver.

The script deliberately consumes a server-issued governed proof rather than
creating, guessing, or replaying Owner presence.  A trusted Owner surface (or
the OS controlled test handshake) must obtain that proof for the exact Pack,
action, and transaction first.
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
    emit_pack_error, UNINSTALL_CONNECTIVITY, UNINSTALL_GENERIC, UNINSTALL_LIFECYCLE_HTTP,
)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080").rstrip("/")
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")
PROOF_RAW = os.environ.get("TRUZHEN_PACK_UNINSTALL_PROOF_JSON", "").strip()


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    request = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        request.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(request, timeout=120) as response:
            raw = response.read().decode("utf-8")
            return response.status, json.loads(raw) if raw else {}
    except urllib.error.HTTPError as error:
        raw = error.read().decode("utf-8")
        try:
            return error.code, json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            return error.code, {"_raw": raw}
    except Exception as error:
        return 0, {"_transport_error": str(error)}


def die(message, error_code=UNINSTALL_GENERIC):
    emit_pack_error(
        pack_dir=PACK_DIR, base=BASE, action="uninstall", error_code=error_code, message=message
    )
    print("卸载失败：" + message, file=sys.stderr)
    sys.exit(1)


def load_proof(pack_ref, version):
    if not PROOF_RAW:
        die("缺少可信 Owner 前台签发的 TRUZHEN_PACK_UNINSTALL_PROOF_JSON")
    try:
        proof = json.loads(PROOF_RAW)
    except json.JSONDecodeError as error:
        die("Owner 卸载证明不是合法 JSON：" + str(error))
    required = ("decision_ref", "run_id", "nonce", "owner_action_evidence_ref")
    expected_transaction = "transaction://pack-uninstall:" + pack_ref + "@" + version
    if (
        proof.get("action_type") != "14.pack-studio.lifecycle.uninstall"
        or proof.get("target_ref") != pack_ref
        or proof.get("transaction_ref") != expected_transaction
        or any(not proof.get(key) for key in required)
    ):
        die("Owner 卸载证明与 Pack/action/transaction 不匹配")
    return proof


def main():
    with open(os.path.join(PACK_DIR, "manifest.json"), encoding="utf-8") as file:
        manifest = json.load(file)
    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    pack_name = manifest.get("name", "家政运营 Pack")
    print("== 正式卸载 %s @ %s（%s）==" % (pack_ref, version, BASE))

    code, body = call(
        "GET", "/v3/pack-studio/lifecycle/packs?" + urllib.parse.urlencode({"pack_ref": pack_ref})
    )
    if code == 0:
        die("连不上 devserver（%s）" % BASE, UNINSTALL_CONNECTIVITY)
    enabled = any(
        entry.get("pack_ref") == pack_ref
        and (entry.get("enabled_pointer") or {}).get("current_version")
        for entry in body.get("packs", []) or []
    )
    if not enabled:
        die("场景包未处于启用态，拒绝把非正式状态冒充为已卸载")

    proof = load_proof(pack_ref, version)
    code, body = call(
        "POST",
        "/v3/pack-studio/lifecycle/uninstall",
        {
            "pack_ref": pack_ref,
            "owner_ref": OWNER,
            "reason": "Owner 正式卸载 " + pack_name,
            "decision_ref": proof["decision_ref"],
            "run_id": proof["run_id"],
            "nonce": proof["nonce"],
            "owner_action_evidence_ref": proof["owner_action_evidence_ref"],
        },
    )
    if code != 200:
        die("formal uninstall HTTP %d: %s" % (code, body), UNINSTALL_LIFECYCLE_HTTP)
    print("正式卸载成功：%s 已进入 uninstalled；历史事务与 Receipt 保留可反查。" % pack_ref)


if __name__ == "__main__":
    main()
