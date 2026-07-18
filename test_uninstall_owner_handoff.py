#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Owner 在场卸载交接行为测试：脚本只读、等待 os-14，不自铸主权。"""

import contextlib
import importlib.util
import io
import os
import sys
import unittest
from unittest import mock

REPO = os.path.dirname(os.path.abspath(__file__))


def load_uninstall(name, pack_dir):
    spec = importlib.util.spec_from_file_location(name, os.path.join(REPO, pack_dir, "uninstall.py"))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class OwnerHandoffUninstallTest(unittest.TestCase):
    def _run_success(self, pack_dir, args):
        module = load_uninstall("uninstall_" + pack_dir.replace("-", "_"), pack_dir)
        calls = []
        states = [True, False]
        with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
            pack_ref = module.json.load(stream)["pack_ref"]

        def fake_call(method, path, body=None):
            calls.append((method, path, body))
            enabled = states.pop(0) if states else False
            return 200, {"packs": [{"pack_ref": pack_ref, "enabled_pointer": {"current_version": "1.0.0" if enabled else ""}}]}

        module.call = fake_call
        out = io.StringIO()
        with mock.patch.object(sys, "argv", ["uninstall.py", *args, "--wait-seconds", "0"]):
            with contextlib.redirect_stdout(out):
                module.main()
        self.assertTrue(calls)
        self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)
        self.assertIn("os-14 已证明 Pack 停用", out.getvalue())

    def test_smart_home_handoff_is_read_only(self):
        self._run_success("smart-home-owner-pack-v0", ["--devserver-base", "http://127.0.0.1:18080"])

    def test_environmental_handoff_is_read_only(self):
        self._run_success("environmental-enforcement-pack-v0", ["--devserver-base", "http://127.0.0.1:18099"])

    def test_scripts_have_no_background_gate_or_disable_write(self):
        for pack_dir in ("smart-home-owner-pack-v0", "environmental-enforcement-pack-v0"):
            with open(os.path.join(REPO, pack_dir, "uninstall.py"), encoding="utf-8") as stream:
                source = stream.read()
            self.assertNotIn("/v3/base/gated-actions/confirm", source)
            self.assertNotIn("/v3/base/gated-actions/prepare", source)
            self.assertNotIn("/v3/pack-studio/lifecycle/disable", source)
            self.assertIn("wait_for_owner_disabled", source)


if __name__ == "__main__":
    unittest.main()
