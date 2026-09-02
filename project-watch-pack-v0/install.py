#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""项目关注 Pack 可信 GUI 装入交接：脚本全程只读。

沿用 smart-home 的只读交接骨架——脚本不伪造浏览器 Origin、Cookie、Owner presence
或 Base 决议，也不调用任何 prepare / confirm / draft / promote / 知识 approve 写端点。
它只做三件事：
  1. 装入前的本地自检（声明 readiness + 知识 checksum 防漂移），不满足就早停；
  2. 展示可信前台交接，并只读等待 os-14 lifecycle ReadModel 证明精确版本已启用；
  3. 启用后只读复核角色槽绑定与知识入库两个下游阶段，缺一即 fail closed。
任何「未观察到」都按阶段化错误码诚实阻断，不猜测、不宣称装入成功。
"""

import argparse
import json
import os
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

PACK_DIR = os.path.dirname(os.path.abspath(__file__))
REPO_DIR = os.path.dirname(PACK_DIR)
if REPO_DIR not in sys.path:
    sys.path.insert(0, REPO_DIR)

from knowledge_checksums import verify_entries  # noqa: E402
from pack_diagnostics import (  # noqa: E402
    INSTALL_BASE_GATE,
    INSTALL_CONNECTIVITY,
    INSTALL_GENERIC,
    INSTALL_KNOWLEDGE,
    INSTALL_KNOWLEDGE_CHECKSUM,
    INSTALL_LIFECYCLE_HTTP,
    INSTALL_READINESS,
    INSTALL_ROLE_BINDING,
    INSTALL_STATE_CONFLICT,
    emit_pack_error,
    pack_enabled_version_from_readmodel,
    present_owner_install_handoff,
    wait_for_owner_enabled,
)

BASE = ""

# 六件事 + 护城河的声明字段（AGENTS.md §3.1）；缺一即 readiness 不通过。
REQUIRED_GOVERNANCE_FIELDS = (
    "person_strategy",
    "formalization_requirement",
    "gates",
    "provider_requirements",
    "notification_command_report_routes",
    "multi_role_comparison",
    "moat_justification",
)
# manifest 与 capabilities 必须逐条同源的字段（description / binding 各自私有）。
ALIGNED_REQUIREMENT_FIELDS = (
    "capability",
    "gateway_class",
    "risk_class",
    "fallback_policy",
    "provider_family",
    "execution_level",
    "runtime_requirement",
)


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


def load(rel):
    with open(os.path.join(PACK_DIR, rel), encoding="utf-8") as stream:
        return json.load(stream)


def declaration_readiness_problems(manifest):
    """装入前的本地声明自检；只读本 pack 自身文件，不联网、不猜测。"""
    problems = []
    for field in REQUIRED_GOVERNANCE_FIELDS:
        if field not in manifest:
            problems.append("缺治理声明字段 %s" % field)
    for key in ("flow_file", "role_slots_file", "capabilities_file"):
        rel = manifest.get(key)
        if not rel or not os.path.exists(os.path.join(PACK_DIR, rel)):
            problems.append("%s 指向的文件不存在：%s" % (key, rel))
    if problems:
        return problems

    scopes_manifest = manifest.get("knowledge_scopes_manifest")
    if manifest.get("knowledge_scopes") and scopes_manifest:
        declared = set(manifest["knowledge_scopes"])
        actual = {s.get("scope_ref") for s in load(scopes_manifest).get("scopes", [])}
        if declared != actual:
            problems.append("knowledge_scopes 与 %s 不一致：%s" % (
                scopes_manifest, sorted(declared.symmetric_difference(actual))))

    caps = load(manifest["capabilities_file"]).get("provider_requirements", [])
    caps_by_id = {item.get("requirement_id"): item for item in caps}
    manifest_by_id = {item.get("requirement_id"): item for item in manifest.get("provider_requirements", [])}
    if set(caps_by_id) != set(manifest_by_id):
        problems.append("manifest 与 capabilities 的 ProviderRequirement 集合不一致：%s" % sorted(
            set(caps_by_id).symmetric_difference(set(manifest_by_id))))
    for requirement_id in sorted(set(caps_by_id) & set(manifest_by_id)):
        for field in ALIGNED_REQUIREMENT_FIELDS:
            if caps_by_id[requirement_id].get(field) != manifest_by_id[requirement_id].get(field):
                problems.append("ProviderRequirement %s 的 %s 在 manifest 与 capabilities 不一致" % (
                    requirement_id, field))
    return problems


def observe_role_bindings(bindings, scope_ref):
    """只读核对 os-13 角色槽绑定；返回 (ok, reason)。"""
    code, body = call("GET", "/v3/agent-orchestration/agent-slots/readmodel")
    if code == 0:
        return False, "connectivity"
    if code != 200:
        return False, "role_binding_readmodel_http_%s" % code
    if not isinstance(body, dict) or not isinstance(body.get("agent_slot_bindings"), list):
        return False, "role_binding_readmodel_invalid"
    enabled = set()
    for record in body["agent_slot_bindings"]:
        if not isinstance(record, dict):
            return False, "role_binding_readmodel_invalid"
        if record.get("enabled_state") == "enabled" and record.get("scope_ref") in (scope_ref, ""):
            enabled.add(record.get("slot_ref"))
    missing = sorted({b["slot_id"] for b in bindings} - enabled)
    if missing:
        return False, "slot_binding_missing:%s" % ",".join(missing)
    return True, "role_bindings_observed"


def observe_knowledge(entries, owner_ref, pack_ref, version):
    """只读核对 os-09 FormalKnowledge 入库与其 Base 门回执；返回 (ok, reason)。"""
    for entry in entries:
        source_ref = entry["source_ref"] + "@" + version
        path = ("/v3/memory/knowledge/formal?owner_id=" + urllib.parse.quote(owner_ref, safe="")
                + "&scope=Formal&pack_ref=" + urllib.parse.quote(pack_ref, safe="")
                + "&scene_ref=" + urllib.parse.quote(entry.get("scene_ref", ""), safe="")
                + "&source_ref=" + urllib.parse.quote(source_ref, safe="")
                + "&include_unmounted=true&limit=50")
        code, body = call("GET", path)
        if code == 0:
            return False, "connectivity"
        if code != 200:
            return False, "knowledge_readmodel_http_%s" % code
        if not isinstance(body, dict) or not isinstance(body.get("items"), list):
            return False, "knowledge_readmodel_invalid"
        item = next((i for i in body["items"] if isinstance(i, dict) and i.get("source_ref") == source_ref), None)
        if item is None:
            return False, "knowledge_not_ingested:%s" % source_ref
        if not str(item.get("receipt_ref") or "").strip():
            # 入库了却无可反查回执 = Base 门这一步没有留下签发证据，属 Base gate 阶段失败。
            return False, "knowledge_receipt_missing:%s" % source_ref
    return True, "knowledge_observed"


def wait_observation(observer, timeout_seconds, poll_seconds, sleep=time.sleep):
    """只读轮询下游阶段观测；超时返回最后一次原因，绝不改判为成功。"""
    deadline = time.monotonic() + max(0.0, float(timeout_seconds))
    while True:
        ok, reason = observer()
        if ok:
            return True, reason
        if reason == "connectivity" or time.monotonic() >= deadline:
            return False, reason
        sleep(max(0.05, float(poll_seconds)))


def main():
    parser = argparse.ArgumentParser(description="把项目关注 Pack 的装入交给可信 Truzhen GUI。")
    parser.add_argument("--devserver-base", default=os.environ.get("TRUZHEN_DEVSERVER_BASE", "").strip())
    parser.add_argument("--client-url", default=os.environ.get("TRUZHEN_CLIENT_URL", ""))
    parser.add_argument("--owner-ref", default=os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default"),
                        help="只用于只读 ReadModel 查询键；脚本不注入身份，也不代 Owner 决策。")
    parser.add_argument("--open-gui", action="store_true", help="仅显式打开前台；不注入登录态或 Owner presence。")
    parser.add_argument("--wait-seconds", type=float,
                        default=float(os.environ.get("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS", "300")))
    parser.add_argument("--poll-seconds", type=float, default=1.0)
    args = parser.parse_args()

    global BASE
    BASE = args.devserver_base.rstrip("/")
    if not BASE:
        die("必须显式指定 TRUZHEN_DEVSERVER_BASE 或 --devserver-base", INSTALL_CONNECTIVITY)

    manifest = load("manifest.json")
    pack_ref, version = manifest["pack_ref"], manifest["version"]

    problems = declaration_readiness_problems(manifest)
    if problems:
        die("装入前声明自检未通过：" + "；".join(problems), INSTALL_READINESS)

    index_doc = load(manifest["knowledge_index"]) if manifest.get("knowledge_index") else {"entries": []}
    entries = index_doc.get("entries", [])
    checksum_problems = verify_entries(PACK_DIR, entries)
    if checksum_problems:
        die("知识内容与 index checksum 漂移，拒绝装入：" + "; ".join(checksum_problems),
            INSTALL_KNOWLEDGE_CHECKSUM)

    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + urllib.parse.quote(pack_ref, safe=""))
    if code == 0:
        die("连不上 devserver（%s）" % BASE, INSTALL_CONNECTIVITY)
    if code != 200:
        die("lifecycle ReadModel HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
    enabled_version = pack_enabled_version_from_readmodel(body, pack_ref)
    if enabled_version is None:
        die("lifecycle ReadModel 形状不完整，拒绝猜测状态", INSTALL_LIFECYCLE_HTTP)
    if enabled_version and enabled_version != version:
        die("已启用版本 %s 与声明版本 %s 不一致，拒绝覆盖" % (enabled_version, version), INSTALL_STATE_CONFLICT)
    if enabled_version != version:
        present_owner_install_handoff(args.client_url, pack_ref, version, open_gui=args.open_gui)
        ok, reason = wait_for_owner_enabled(call, pack_ref, version, args.wait_seconds, args.poll_seconds)
        if not ok:
            die("%s；未观察到 os-14 的精确版本状态" % reason,
                INSTALL_CONNECTIVITY if reason == "connectivity" else INSTALL_LIFECYCLE_HTTP)
    print("os-14 已证明精确 Pack 版本启用。")

    pack_version_ref = pack_ref + "@" + version
    role_slots_doc = load(manifest["role_slots_file"])
    bindings = role_slots_doc.get("bindings", [])
    if bindings:
        ok, reason = wait_observation(
            lambda: observe_role_bindings(bindings, pack_version_ref), args.wait_seconds, args.poll_seconds)
        if not ok:
            die("%s；未观察到本 Pack 声明的角色槽绑定生效" % reason,
                INSTALL_CONNECTIVITY if reason == "connectivity" else INSTALL_ROLE_BINDING)
        print("os-13 已证明 %d 个角色槽绑定生效。" % len(bindings))

    if entries:
        ok, reason = wait_observation(
            lambda: observe_knowledge(entries, args.owner_ref, pack_ref, version),
            args.wait_seconds, args.poll_seconds)
        if not ok:
            if reason == "connectivity":
                code = INSTALL_CONNECTIVITY
            elif reason.startswith("knowledge_receipt_missing"):
                code = INSTALL_BASE_GATE
            else:
                code = INSTALL_KNOWLEDGE
            die("%s；未观察到本 Pack 声明的知识已正式入库并留下可反查回执" % reason, code)
        print("os-09 已证明 %d 条声明知识入库且带可反查回执（仍为 pending_human_review）。" % len(entries))

    print("✅ 只读交接完成：场景包、角色槽绑定与知识入库三个阶段均由权威 ReadModel 证明。")


if __name__ == "__main__":
    main()
