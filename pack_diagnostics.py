#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Structured lifecycle diagnostics for pack install / uninstall scripts.

错误码按失败阶段细分（Y4），契约形状对齐 truzhen-contracts monitoring error_code
pattern `TZ-<仓>-<域>-<NNN>`。追加式：永不复用、永不改义；废弃只标注不删除。
新增码必须登记进本模块的 PACK_ERROR_CODES，否则 pack_error_code_taxonomy_test 报红。
"""

import json
import os
import re
import sys
import time
import urllib.parse
import webbrowser

# --- install 阶段码 -------------------------------------------------------
INSTALL_GENERIC = "TZ-PACK-INSTALL-001"        # 未归类失败（兜底，避免误用）
INSTALL_CONNECTIVITY = "TZ-PACK-INSTALL-002"   # 连不上 devserver / 传输层错误
INSTALL_LIFECYCLE_HTTP = "TZ-PACK-INSTALL-003" # 场景包生命周期 HTTP 非预期（canvas/draft/promote/confirm）
INSTALL_READINESS = "TZ-PACK-INSTALL-004"      # readiness 未通过（六件事声明不满足）
INSTALL_STATE_CONFLICT = "TZ-PACK-INSTALL-005" # 版本状态冲突（draft_frozen / 遗留态占用耗尽）
INSTALL_ROLE_BINDING = "TZ-PACK-INSTALL-006"   # 角色包生命周期 / 绑槽失败
INSTALL_KNOWLEDGE = "TZ-PACK-INSTALL-007"      # 知识库入库 / 审批失败
INSTALL_BASE_GATE = "TZ-PACK-INSTALL-008"      # Base gated-action 签发失败
INSTALL_KNOWLEDGE_CHECKSUM = "TZ-PACK-INSTALL-009"  # 知识内容与 index checksum 漂移（装入前 fail-fast 拒绝）

# --- uninstall 阶段码 -----------------------------------------------------
UNINSTALL_GENERIC = "TZ-PACK-UNINSTALL-001"        # 未归类失败（兜底）
UNINSTALL_CONNECTIVITY = "TZ-PACK-UNINSTALL-002"   # 连不上 devserver / 传输层错误
UNINSTALL_LIFECYCLE_HTTP = "TZ-PACK-UNINSTALL-003" # 停用 / 卸载 HTTP 非预期

# 登记簿：所有合法码。guard 测试据此拒绝裸造未登记码。
PACK_ERROR_CODES = (
    INSTALL_GENERIC, INSTALL_CONNECTIVITY, INSTALL_LIFECYCLE_HTTP, INSTALL_READINESS,
    INSTALL_STATE_CONFLICT, INSTALL_ROLE_BINDING, INSTALL_KNOWLEDGE, INSTALL_BASE_GATE,
    INSTALL_KNOWLEDGE_CHECKSUM,
    UNINSTALL_GENERIC, UNINSTALL_CONNECTIVITY, UNINSTALL_LIFECYCLE_HTTP,
)

_CODE_SHAPE = re.compile(r"^TZ-PACK-[A-Z0-9]{2,10}-\d{3}$")


def is_registered_code(error_code):
    """码是否已登记且形状合法。"""
    return bool(_CODE_SHAPE.match(str(error_code))) and error_code in PACK_ERROR_CODES


def emit_pack_error(*, pack_dir, base, action, error_code, message):
    manifest = _load_manifest(pack_dir)
    payload = {
        "event_type": "pack_lifecycle_error",
        "source_kind": "pack_lifecycle",
        "action": action,
        "error_code": error_code,
        "pack_ref": manifest.get("pack_ref", ""),
        "pack_version": manifest.get("version", ""),
        "devserver_base": _sanitize_base(base),
        "message": _truncate(str(message), 2048),
    }
    print("TRUZHEN_PACK_ERROR " + json.dumps(payload, ensure_ascii=False, sort_keys=True), file=sys.stderr)


def pack_enabled_from_readmodel(body, pack_ref):
    """只读解析 os-14 lifecycle 真相；缺字段时保守视为仍启用。"""
    if not isinstance(body, dict) or not isinstance(body.get("packs"), list):
        return None
    for entry in body["packs"]:
        if isinstance(entry, dict) and entry.get("pack_ref") == pack_ref:
            pointer = entry.get("enabled_pointer")
            if not isinstance(pointer, dict):
                return None
            return bool(pointer.get("current_version"))
    return False


def wait_for_owner_disabled(call, pack_ref, timeout_seconds, poll_seconds=1.0, sleep=time.sleep):
    """等待 GUI Owner 动作落入 os-14；本函数永远只调用 GET。"""
    deadline = time.monotonic() + max(0.0, float(timeout_seconds))
    path = "/v3/pack-studio/lifecycle/packs?pack_ref=" + urllib.parse.quote(pack_ref, safe="")
    while True:
        code, body = call("GET", path)
        if code == 0:
            return False, "connectivity"
        if code != 200:
            return False, "readmodel_http_%s" % code
        enabled = pack_enabled_from_readmodel(body, pack_ref)
        if enabled is False:
            return True, "disabled"
        if enabled is None:
            return False, "readmodel_invalid"
        if time.monotonic() >= deadline:
            return False, "owner_presence_required"
        sleep(max(0.05, float(poll_seconds)))


def present_owner_disable_handoff(client_url, pack_ref, open_gui=False):
    """展示前台交接；只有显式 --open-gui 才打开浏览器，不注入身份。"""
    url = str(client_url or "").strip()
    print("需要 Owner 在可信前台完成停用：%s" % pack_ref)
    print("请打开 Truzhen → 场景平台 → 场景包管理，选择该 Pack 并确认停用。")
    if url:
        print("前台地址：%s" % url)
        if open_gui:
            webbrowser.open(url, new=2, autoraise=True)


def _load_manifest(pack_dir):
    try:
        with open(os.path.join(pack_dir, "manifest.json"), encoding="utf-8") as f:
            data = json.load(f)
            if isinstance(data, dict):
                return data
    except Exception:
        return {}
    return {}


def _sanitize_base(value):
    parsed = urllib.parse.urlsplit(str(value or ""))
    if not parsed.scheme or not parsed.netloc:
        return _truncate(str(value or ""), 256)
    host = parsed.hostname or ""
    if parsed.port:
        host = "%s:%d" % (host, parsed.port)
    return urllib.parse.urlunsplit((parsed.scheme, host, parsed.path, "", ""))


def _truncate(value, limit):
    if len(value) <= limit:
        return value
    return value[:limit]
