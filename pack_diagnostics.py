#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Structured lifecycle diagnostics for pack install / uninstall scripts."""

import json
import os
import sys
import urllib.parse


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
