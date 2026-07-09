#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
环保执法 Pack —— 装入正在运行的 Truzhen devserver（可加载）。

走的是前端/产品种子完全相同的真实 lifecycle 端点（canvas → 06 同步 → lifecycle
draft/readiness/promote/confirm → 角色包 lifecycle → 绑槽 → 09 知识库批量入库 +
Base 签发 approve decision + approve），全程产 03 receipt，不手写裸塞、不绕主权链。

前置：先在本 worktree 起 devserver（且已从 server.go 摘除环保自动 seed）：
  go run ./backend/cmd/devserver        # 默认 127.0.0.1:18080
然后：
  python3 packs/environmental-enforcement-pack-v0/install.py

幂等：已启用的场景包/角色包/绑定会跳过；知识批量入库按内容去重。
"""
import hashlib
import json
import os
import re
import sys
import uuid
import urllib.error
import urllib.parse
import urllib.request

PACK_DIR = os.path.dirname(os.path.abspath(__file__))
REPO_DIR = os.path.dirname(PACK_DIR)
if REPO_DIR not in sys.path:
    sys.path.insert(0, REPO_DIR)
from pack_diagnostics import (
    emit_pack_error, INSTALL_GENERIC, INSTALL_CONNECTIVITY, INSTALL_LIFECYCLE_HTTP,
    INSTALL_READINESS, INSTALL_STATE_CONFLICT, INSTALL_ROLE_BINDING, INSTALL_KNOWLEDGE,
    INSTALL_BASE_GATE)

BASE = os.environ.get("TRUZHEN_DEVSERVER_BASE", "http://127.0.0.1:18080")
# 用本地规范 Owner（前端记忆中心默认查询 owner_id='owner://local/default'，后端运行时
# 也用此 owner）：知识/挂载/角色/绑定都落在这个 owner 下，记忆中心与运行时 advice 才看得到。
OWNER = os.environ.get("TRUZHEN_PACK_OWNER", "owner://local/default")


def load(rel):
    with open(os.path.join(PACK_DIR, rel), encoding="utf-8") as f:
        return json.load(f)


def load_opt(rel, default):
    """可选清单：字段缺失或文件不存在时返回 default（让本 loader 同时适用无知识库的 pack）。"""
    if not rel:
        return default
    path = os.path.join(PACK_DIR, rel)
    if not os.path.exists(path):
        return default
    with open(path, encoding="utf-8") as f:
        return json.load(f)


class PackVersionConflict(Exception):
    """目标版本被历史 disabled 记录占用（unload 后 reload 场景）。"""


def bump_patch(v):
    parts = v.split(".")
    if parts and parts[-1].isdigit():
        parts[-1] = str(int(parts[-1]) + 1)
        return ".".join(parts)
    return v + ".1"


def is_state_conflict(body):
    """判断响应是否为「该版本已被 disabled/冻结、无法再次 draft/check」的状态冲突。"""
    s = json.dumps(body, ensure_ascii=False)
    return any(k in s for k in ("is not checkable", "state=disabled", "DraftFrozen",
                                "draft is frozen", "not_confirmable", "ErrDraftFrozen"))


def call(method, path, body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            raw = resp.read().decode("utf-8")
            code = resp.status
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8")
        code = e.code
    except Exception as e:
        return 0, {"_transport_error": str(e)}
    try:
        return code, json.loads(raw) if raw else {}
    except json.JSONDecodeError:
        return code, {"_raw": raw}


def die(msg, error_code=INSTALL_GENERIC):
    emit_pack_error(pack_dir=PACK_DIR, base=BASE, action="install", error_code=error_code, message=msg)
    print("装入失败：" + msg, file=sys.stderr)
    sys.exit(1)


def strip_frontmatter(text):
    m = re.match(r"^---\s*\n.*?\n---\s*\n", text, re.S)
    return text[m.end():] if m else text


def main():
    manifest = load("manifest.json")
    flow = load(manifest["flow_file"])
    role_slots_doc = load(manifest["role_slots_file"])
    caps = load(manifest["capabilities_file"])
    scopes_doc = load_opt(manifest.get("knowledge_scopes_manifest"), {"scopes": []})
    kindex = load_opt(manifest.get("knowledge_index"), {"entries": []})

    pack_ref = manifest["pack_ref"]
    version = manifest["version"]
    flow_id = flow["flow_id"]
    pack_version_ref = pack_ref + "@" + version

    print("== 装入 %s @ %s 到 %s ==" % (pack_ref, version, BASE))

    # 健康检查
    code, _ = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + pack_ref)
    if code == 0:
        die("连不上 devserver（%s）。请先 go run ./backend/cmd/devserver" % BASE, INSTALL_CONNECTIVITY)

    # 解析装入版本：
    #   - 该 pack 已有 enabled 版本 → 幂等跳过场景包步骤（已装）。
    #   - 否则用 manifest 版本装入；若该版本被历史 disabled 记录占用（unload 后 reload，
    #     共享 lifecycle 不支持同版本 disabled→enabled 经 HTTP 重启），自动 bump patch，
    #     走 lifecycle 支持的「新版本 draft→…→enabled」路径，实现 unload 后可 reload。
    enabled_version = None
    code, body = call("GET", "/v3/pack-studio/lifecycle/packs?pack_ref=" + pack_ref)
    if code == 200:
        for entry in body.get("packs", []) or []:
            if entry.get("pack_ref") == pack_ref:
                cur = (entry.get("enabled_pointer") or {}).get("current_version")
                if cur:
                    enabled_version = cur

    def do_seed_scene(install_version):
        pvr = pack_ref + "@" + install_version
        # 1. canvas → 06（flow_id 与版本无关，幂等 upsert）
        print("[1/6] 画布写穿 06 ...")
        code, body = call("POST", "/v3/pack-studio/canvas", {
            "flow_id": flow_id, "title": flow.get("title", ""),
            "occ_version": 0, "save_source": "pack_install", "flow_spec_draft": flow})
        if code == 409 and isinstance(body.get("current_occ_version"), (int, float)):
            code, body = call("POST", "/v3/pack-studio/canvas", {
                "flow_id": flow_id, "title": flow.get("title", ""),
                "occ_version": int(body["current_occ_version"]),
                "save_source": "pack_install", "flow_spec_draft": flow})
        if code != 200:
            die("canvas HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
        if not ((body.get("engine_sync") or {}).get("synced")):
            die("canvas 未同步进 06：%s" % body.get("engine_sync"), INSTALL_LIFECYCLE_HTTP)

        # 2. lifecycle draft（六件事 + 知识域）
        print("[2/6] lifecycle draft（六件事 + %d 知识域）..." % len(scopes_doc.get("scopes", [])))
        provider_reqs = []
        for p in caps.get("provider_requirements", []):
            provider_reqs.append({k: p[k] for k in ("requirement_id", "capability", "gateway_class",
                                                    "risk_class", "fallback_policy", "provider_family") if k in p}
                                 | ({"optional": p["optional"]} if p.get("optional") else {}))
        ks = []
        for s in scopes_doc.get("scopes", []):
            ks.append({"scope_ref": s["scope_ref"], "display_name": s["display_name"],
                       "scene_ref": s.get("scene_ref", ""), "mount_policy": s.get("mount_policy", "mount_on_pack_enable"),
                       "knowledge_kinds": s.get("knowledge_kinds", []), "required": s.get("required", True)})
        routes = manifest.get("notification_command_report_routes", {})
        draft = {
            "pack_ref": pack_ref, "version": install_version, "title": manifest["name"],
            "template_family": manifest.get("template_family", "合规审查执法证据链型"),
            "flow_id": flow_id,
            "role_slots": [{"slot_id": s["slot_id"], "responsibility": s["responsibility"],
                            "required_role": s["required_role"], "default_role_pack_ref": s["default_role_pack_ref"],
                            "slice_scope_policy": s.get("slice_scope_policy", ""),
                            "node_type": s.get("node_type", "")} for s in role_slots_doc["role_slots"]],
            "provider_requirements": provider_reqs,
            "formalization_requirement": manifest["formalization_requirement"]["summary"],
            "notification_routes": routes.get("notification", []),
            "command_candidates": routes.get("command_candidate", []),
            "report_routes": routes.get("report", []),
            "moat_justification": manifest["moat_justification"],
            "knowledge_scopes": ks,
            "idempotency_key": "pack-install-draft:" + pvr,
            "actor_ref": OWNER,
        }
        code, body = call("POST", "/v3/pack-studio/lifecycle/draft", draft)
        if code != 200:
            if is_state_conflict(body):
                raise PackVersionConflict(body)
            die("draft HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

        # 3. readiness
        print("[3/6] readiness ...")
        code, body = call("POST", "/v3/pack-studio/lifecycle/readiness",
                          {"pack_ref": pack_ref, "version": install_version, "actor_ref": OWNER})
        if code != 200:
            if is_state_conflict(body):
                raise PackVersionConflict(body)
            die("readiness HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)
        rr = (body.get("record") or {}).get("readiness_report") or {}
        if not rr.get("ready"):
            die("readiness 未通过：%s" % body, INSTALL_READINESS)

        # 4. promote
        print("[4/6] promote ...")
        code, body = call("POST", "/v3/pack-studio/lifecycle/promote",
                          {"pack_ref": pack_ref, "version": install_version, "actor_ref": OWNER})
        if code != 200:
            die("promote HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

        # 5. confirm → enabled
        print("[5/6] confirm（真 01 Base Gate + 03 receipt）...")
        code, body = call("POST", "/v3/pack-studio/lifecycle/confirm", {
            "pack_ref": pack_ref, "version": install_version,
            "idempotency_key": "pack-install-confirm:" + pvr,
            "owner_ref": OWNER, "approve": True, "comment": manifest["name"] + " 启用"})
        if code != 200 or body.get("status") != "enabled":
            die("confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_LIFECYCLE_HTTP)

    reactivated = False
    if enabled_version:
        install_version = enabled_version
        print("[1-5/6] 场景包已启用 @%s，跳过场景包步骤。" % install_version)
    else:
        # 优先同版本 reactivate（unload→reload 回路；disabled→enabled，过真 Base gate +
        # receipt，知识域随版本重挂，无需重灌、不 bump 版本）。
        rk = "pack-reactivate:" + pack_ref + "@" + version + ":" + uuid.uuid4().hex
        code, body = call("POST", "/v3/pack-studio/lifecycle/reactivate", {
            "pack_ref": pack_ref, "version": version, "owner_ref": OWNER,
            "idempotency_key": rk, "comment": manifest["name"] + " 重新启用"})
        if code == 200 and body.get("status") == "enabled":
            install_version = version
            reactivated = True
            print("[1-5/6] 同版本 reactivate 成功（unload→reload）@%s，知识域已重挂。" % version)
        else:
            # 记录从未存在（全新装）→ 完整 draft→…→confirm；极端遗留态用 bump patch 兜底。
            install_version = version
            for _ in range(20):
                try:
                    do_seed_scene(install_version)
                    break
                except PackVersionConflict:
                    nv = bump_patch(install_version)
                    print("    版本 %s 遗留态占用，改用 %s 重试" % (install_version, nv))
                    install_version = nv
            else:
                die("连续 20 个版本都被占用，无法装入", INSTALL_STATE_CONFLICT)
    pack_version_ref = pack_ref + "@" + install_version

    # 6a. 角色包
    print("[6/6] 角色包 + 绑槽 + 知识库 ...")
    code, rpbody = call("GET", "/v3/agent-orchestration/role-packs/readmodel")
    enabled_rp = set()
    for ev in (rpbody.get("enabled_versions") or []):
        rid = ev.get("role_pack_id", "")
        enabled_rp.add(rid.split("@")[0])
    for fn in sorted(os.listdir(os.path.join(PACK_DIR, "role-packs"))):
        if not fn.endswith(".json"):
            continue
        rp = load(os.path.join("role-packs", fn))
        rid = rp["role_pack_id"]
        if rid in enabled_rp:
            print("    角色包 %s 已启用，跳过" % rid)
            continue
        install_role_pack(rp)
        print("    角色包 %s 已启用" % rid)

    # 6b. 绑槽
    scope_ref = pack_version_ref
    for b in role_slots_doc.get("bindings", []):
        bind_slot(b, scope_ref)
    print("    绑槽完成")

    # 6c. 知识库入库（按知识域分组）
    entries = kindex.get("entries", [])
    if reactivated:
        print("    reactivate 路径：知识域已随版本重挂，跳过重灌")
    elif not entries:
        print("    本 pack 无知识库，跳过知识入库")
    else:
        groups = {}
        for e in entries:
            groups.setdefault(e["knowledge_scope_ref"], []).append(e)
        total = 0
        for scope_ref_k, items in groups.items():
            n = ingest_knowledge_scope(pack_ref, pack_version_ref, scope_ref_k, items)
            total += n
        print("    知识库入库完成：%d 条已确认为 FormalKnowledge（pending_human_review）" % total)

    print("\n✅ 装入成功：%s @ %s 已 enabled。前端「场景包管理」刷新可见。" % (pack_ref, install_version))


def install_role_pack(rp):
    sani = re.sub(r'[/:@\s]+', "-", rp["role_pack_id"]).strip("-")
    draft_id = "role_pack_draft-pack-install-" + sani
    idem = "pack-install-" + sani
    cs = rp.get("communication_style", {})
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts", {
        "draft_id": draft_id, "role_pack_id": rp["role_pack_id"], "version": rp["version"],
        "display_name": rp["display_name"], "description": rp["description"],
        "target_use": rp.get("target_use", []), "style_summary": rp.get("style_summary", ""),
        "decision_style_summary": rp.get("decision_style_summary", ""),
        "model_policy_ref": rp.get("model_policy_ref", ""),
        "communication_style": {"structure": cs.get("structure", ""), "tone": cs.get("tone", ""),
                                "forbidden_phrases": cs.get("forbidden_phrases", [])},
        "forbidden_behavior_policy_ref": rp.get("forbidden_policy_ref", ""),
        "risk_level": rp.get("risk_level", "medium"), "owner_ref": OWNER,
        "evidence_refs": ["evidence://pack-install/role-pack/" + sani + "/draft"],
        "idempotency_key": idem + "-draft"})
    if code != 200:
        die("role draft HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    draft_id = body.get("draft_id") or (body.get("draft") or {}).get("draft_id") or draft_id
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/readiness-check", {})
    if code != 200:
        die("role readiness HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/drafts/" + draft_id + "/promote-candidate", {
        "owner_ref": OWNER, "target_agent_ref": rp["role_pack_id"], "idempotency_key": idem + "-promote",
        "evidence_refs": ["evidence://pack-install/role-pack/promote"]})
    if code != 200 or not body.get("ok"):
        die("role promote HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/enable-candidate", {
        "role_pack_id": rp["role_pack_id"], "version": rp["version"], "target_agent_ref": rp["role_pack_id"],
        "owner_ref": OWNER, "idempotency_key": idem + "-enable",
        "evidence_refs": ["evidence://pack-install/role-pack/enable"]})
    if code != 200 or not body.get("ok"):
        die("role enable-candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    code, body = call("POST", "/v3/agent-orchestration/role-packs/enable-confirm", {
        "role_pack_id": rp["role_pack_id"], "version": rp["version"], "target_agent_ref": rp["role_pack_id"],
        "owner_ref": OWNER, "idempotency_key": idem + "-confirm", "approve": True,
        "comment": rp["display_name"] + " 启用", "evidence_refs": ["evidence://pack-install/role-pack/enable-confirm"]})
    if code != 200 or body.get("status") != "enabled":
        die("role enable-confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def bind_slot(b, scope_ref):
    code, rm = call("GET", "/v3/agent-orchestration/agent-slots/readmodel")
    if code == 200:
        for bd in (rm.get("agent_slot_bindings") or []):
            if bd.get("slot_ref") == b["slot_id"] and bd.get("enabled_state") == "enabled" \
               and bd.get("scope_ref") in (scope_ref, ""):
                return
    code, body = call("POST", "/v3/agent-orchestration/agent-slots/bind-candidate", {
        "slot_ref": b["slot_id"], "scope_ref": scope_ref, "source_pack_ref": scope_ref,
        "required_role": b["required_role"], "requested_agent_ref": b["agent_ref"],
        "requested_role_pack_id": b["role_pack_id"], "ttl": "8760h",
        "evidence_refs": ["evidence://pack-install/agent-slot/" + b["slot_id"] + "/bind"]})
    if code != 200 or not body.get("ok"):
        die("bind-candidate HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)
    binding_ref = body.get("binding_ref")
    code, body = call("POST", "/v3/agent-orchestration/agent-slots/confirm", {
        "binding_ref": binding_ref, "idempotency_key": "pack-install-slot-confirm-" + b["slot_id"],
        "approve": True, "evidence_refs": ["evidence://pack-install/agent-slot/" + b["slot_id"] + "/confirm"]})
    if code != 200 or body.get("status") != "enabled":
        die("slot confirm 未达 enabled HTTP %d: %s" % (code, body), INSTALL_ROLE_BINDING)


def ingest_knowledge_scope(pack_ref, pack_version_ref, scope_ref, items):
    source_files = []
    scene_ref = items[0].get("scene_ref", "")
    ver = pack_version_ref.split("@")[-1]
    for e in items:
        with open(os.path.join(PACK_DIR, e["file"]), encoding="utf-8") as f:
            body_text = strip_frontmatter(f.read())
        content = "# " + e["title"] + "\n" + body_text
        # source_ref 带 pack 版本：把知识归属到本 pack 版本；同版本重装按 source_ref 幂等
        # 去重。unload→reload 走 reactivate（同版本 disabled→enabled）会重挂这批知识，无需重灌。
        versioned_source_ref = e["source_ref"] + "@" + ver
        law_meta = {"verification_status": "pending_human_review", "source_authority": "reference_only"}
        if e.get("authority"):
            law_meta["authority"] = e["authority"]
        if e.get("effective_from"):
            law_meta["effective_from"] = e["effective_from"]
        if e.get("effective_to"):
            law_meta["effective_to"] = e["effective_to"]
        source_files.append({
            "source_ref": versioned_source_ref, "file_name": os.path.basename(e["file"]),
            "content": content, "content_hash": "sha1-" + hashlib.sha1((content + "|" + pack_version_ref).encode("utf-8")).hexdigest(),
            "kind": e["kind"], "evidence_refs": ["evidence://09/pack-knowledge/" + pack_ref + "/" + e["source_ref"]],
            "law_meta": law_meta})
    batch = {
        "owner_id": OWNER, "scope": "Formal",
        "transaction_ref": "transaction://pack-knowledge:" + pack_version_ref,
        "scene_ref": scene_ref, "tags": ["pack 知识库", "环保执法"],
        "owner_action_evidence_ref": "owner_action_evidence://pack-knowledge/" + pack_ref + "/import",
        "policy_snapshot_ref": "policy_snapshot://pack-knowledge/import",
        "source_files": source_files, "pack_ref": pack_ref, "pack_version_ref": pack_version_ref,
        "knowledge_scope_ref": scope_ref}
    code, body = call("POST", "/v3/memory/knowledge/batches", batch)
    if code not in (200, 201):
        die("knowledge batches HTTP %d (scope=%s): %s" % (code, scope_ref, body), INSTALL_KNOWLEDGE)
    candidates = body.get("candidates", []) or []
    approved = 0
    for c in candidates:
        cref = c.get("candidate_ref")
        if not cref or c.get("status") != "pending":
            continue
        dref, run_id, nonce = mint_decision(pack_version_ref, cref)
        code, abody = call("POST", "/v3/memory/knowledge/candidates/" + urllib.parse.quote(cref, safe="") + "/approve", {
            "actor_ref": OWNER,
            "owner_action_evidence_ref": "owner_action_evidence://pack-knowledge/" + pack_ref + "/approve",
            "decision_ref": dref, "run_id": run_id, "nonce": nonce,
            "policy_snapshot_ref": "policy_snapshot://pack-knowledge/approve",
            # 装入 pack = Owner 选择信任 → 知识直接 owner_verified（不堆「待确认」逐条核验）。
            "verify_authority": True,
            "reason": "随场景包装入并经 Owner 信任确认（owner_verified）"})
        if code != 200:
            die("knowledge approve HTTP %d for %s: %s" % (code, cref, abody), INSTALL_KNOWLEDGE)
        approved += 1
    return approved


def mint_decision(scope_ref, candidate_ref):
    code, body = call("POST", "/v3/base/gated-actions/prepare", {
        "action_type": "memory_knowledge_approve", "target_ref": candidate_ref,
        "content_summary": "pack 知识库入库审批：" + candidate_ref,
        "impact_summary": "把 pack 声明的领域知识候选确认为 FormalKnowledge（pending_human_review）",
        "transaction_ref": "transaction://pack-knowledge:" + scope_ref})
    if code != 200:
        die("gated prepare HTTP %d: %s" % (code, body), INSTALL_BASE_GATE)
    issue_ref = (body.get("issue") or {}).get("issue_ref")
    if not issue_ref:
        die("gated prepare 无 issue_ref: %s" % body, INSTALL_BASE_GATE)
    code, body = call("POST", "/v3/base/gated-actions/confirm", {"issue_ref": issue_ref})
    if code != 200:
        die("gated confirm HTTP %d: %s" % (code, body), INSTALL_BASE_GATE)
    iss = body.get("issue") or {}
    if not iss.get("decision_ref"):
        die("gated confirm 无 decision_ref: %s" % body, INSTALL_BASE_GATE)
    return iss["decision_ref"], iss.get("run_id", ""), iss.get("nonce", "")


if __name__ == "__main__":
    main()
