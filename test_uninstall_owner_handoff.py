#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Owner 在场 Pack 脚本测试：只读等待 os-14/os-07，不自铸主权。"""

import contextlib
import importlib.util
import io
import os
import sys
import unittest
from unittest import mock

REPO = os.path.dirname(os.path.abspath(__file__))


def load_script(name, pack_dir, action):
    spec = importlib.util.spec_from_file_location(
        name, os.path.join(REPO, pack_dir, action + ".py")
    )
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def load_diagnostics():
    spec = importlib.util.spec_from_file_location(
        "u04_pack_diagnostics", os.path.join(REPO, "pack_diagnostics.py")
    )
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class OwnerHandoffPackScriptTest(unittest.TestCase):
    def _run_uninstall_success(self, pack_dir, args):
        module = load_script("uninstall_" + pack_dir.replace("-", "_"), pack_dir, "uninstall")
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
        self._run_uninstall_success("smart-home-owner-pack-v0", ["--devserver-base", "http://127.0.0.1:18080"])

    def test_environmental_handoff_is_read_only(self):
        self._run_uninstall_success("environmental-enforcement-pack-v0", ["--devserver-base", "http://127.0.0.1:18099"])

    def _run_install_success(self, pack_dir, version):
        module = load_script("install_" + pack_dir.replace("-", "_"), pack_dir, "install")
        calls = []
        lifecycle_states = ["", version]
        with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
            pack_ref = module.json.load(stream)["pack_ref"]

        def fake_call(method, path, body=None):
            calls.append((method, path, body))
            if path == "/v3/task-governance/schedules":
                manifest = module.load_json("manifest.json")
                schedule_doc = module.load_json(manifest["schedules_file"])
                return 200, {"schedules": [
                    {
                        "transaction_ref": module.schedule_transaction_ref(pack_ref, item["schedule_key"]),
                        "status": "active",
                    }
                    for item in schedule_doc["schedules"]
                ]}
            enabled = lifecycle_states.pop(0) if lifecycle_states else version
            return 200, {"packs": [{
                "pack_ref": pack_ref,
                "enabled_pointer": {"current_version": enabled},
            }]}

        module.call = fake_call
        out = io.StringIO()
        with mock.patch.object(sys, "argv", [
            "install.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
        ]):
            with contextlib.redirect_stdout(out):
                module.main()
        self.assertTrue(calls)
        self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)
        self.assertIn("os-14 已证明精确 Pack 版本启用", out.getvalue())

    def test_content_install_handoff_waits_for_exact_version_and_schedules(self):
        self._run_install_success("content-operations-workbench-v0", "0.2.0")

    def test_smart_home_install_handoff_waits_for_exact_version(self):
        self._run_install_success("smart-home-owner-pack-v0", "1.1.0")

    def test_content_uninstall_handoff_waits_for_pack_and_schedules(self):
        module = load_script("content_uninstall", "content-operations-workbench-v0", "uninstall")
        calls = []
        lifecycle_states = ["0.2.0", ""]
        schedule_states = ["paused"]
        manifest = module.load_json("manifest.json")
        pack_ref = manifest["pack_ref"]
        schedule_doc = module.load_json(manifest["schedules_file"])
        refs = {module.schedule_transaction_ref(pack_ref, item["schedule_key"]) for item in schedule_doc["schedules"]}

        def fake_call(method, path, body=None):
            calls.append((method, path, body))
            if path == "/v3/task-governance/schedules":
                state = schedule_states.pop(0) if schedule_states else "paused"
                return 200, {"schedules": [
                    {"transaction_ref": ref, "status": state} for ref in refs
                ]}
            version = lifecycle_states.pop(0) if lifecycle_states else ""
            return 200, {"packs": [{
                "pack_ref": pack_ref,
                "enabled_pointer": {"current_version": version},
            }]}

        module.call = fake_call
        out = io.StringIO()
        with mock.patch.object(sys, "argv", [
            "uninstall.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
        ]):
            with contextlib.redirect_stdout(out):
                module.main()
        self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)
        self.assertIn("os-07 已证明声明计划非 active", out.getvalue())

    def test_install_rejects_wrong_enabled_version_and_malformed_readmodel(self):
        module = load_script("smart_install_negative", "smart-home-owner-pack-v0", "install")
        with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
            manifest = module.json.load(stream)
        pack_ref = manifest["pack_ref"]
        for current_version in ("9.9.9", 9):
            calls = []

            def fake_call(method, path, body=None):
                calls.append((method, path, body))
                return 200, {"packs": [{
                    "pack_ref": pack_ref,
                    "enabled_pointer": {"current_version": current_version},
                }]}

            module.call = fake_call
            with mock.patch.object(sys, "argv", [
                "install.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
            ]):
                with self.assertRaises(SystemExit):
                    module.main()
            self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)

    def test_lifecycle_helper_rejects_duplicate_conflicts_and_noncanonical_versions(self):
        diagnostics = load_diagnostics()
        pack_ref = "scene_pack://smart-home-owner-project-ops"
        malformed = (
            ["1.1.0", "9.9.9"],
            ["1.1.0", ""],
            ["1.1.0", "1.1.0"],
            [" 1.1.0 "],
            ["   "],
            [9],
        )
        for versions in malformed:
            body = {"packs": [{
                "pack_ref": pack_ref,
                "enabled_pointer": {"current_version": version},
            } for version in versions]}
            self.assertIsNone(
                diagnostics.pack_enabled_version_from_readmodel(body, pack_ref),
                versions,
            )

    def test_all_three_target_consumers_reject_duplicate_lifecycle_records(self):
        targets = (
            ("content-operations-workbench-v0", "install", "0.2.0", "9.9.9"),
            ("content-operations-workbench-v0", "uninstall", "0.2.0", ""),
            ("smart-home-owner-pack-v0", "install", "1.1.0", "9.9.9"),
        )
        for pack_dir, action, canonical, conflict in targets:
            module = load_script("duplicate_" + pack_dir.replace("-", "_") + "_" + action, pack_dir, action)
            with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
                pack_ref = module.json.load(stream)["pack_ref"]
            calls = []

            def fake_call(method, path, body=None):
                calls.append((method, path, body))
                return 200, {"packs": [
                    {"pack_ref": pack_ref, "enabled_pointer": {"current_version": canonical}},
                    {"pack_ref": pack_ref, "enabled_pointer": {"current_version": conflict}},
                ]}

            def die_as_runtime_error(message, *_args):
                raise RuntimeError(message)

            module.call = fake_call
            module.die = die_as_runtime_error
            with mock.patch.object(sys, "argv", [
                action + ".py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
            ]):
                with self.assertRaisesRegex(RuntimeError, "lifecycle ReadModel 形状不完整"):
                    module.main()
            self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)

    def test_content_uninstall_rejects_active_schedule(self):
        module = load_script("content_uninstall_active_schedule", "content-operations-workbench-v0", "uninstall")
        manifest = module.load_json("manifest.json")
        pack_ref = manifest["pack_ref"]
        schedule_doc = module.load_json(manifest["schedules_file"])
        refs = {module.schedule_transaction_ref(pack_ref, item["schedule_key"]) for item in schedule_doc["schedules"]}

        def fake_call(method, path, body=None):
            if path == "/v3/task-governance/schedules":
                return 200, {"schedules": [{"transaction_ref": ref, "status": "active"} for ref in refs]}
            return 200, {"packs": [{"pack_ref": pack_ref, "enabled_pointer": {"current_version": ""}}]}

        module.call = fake_call
        with mock.patch.object(sys, "argv", [
            "uninstall.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
        ]):
            with self.assertRaises(SystemExit):
                module.main()

    def test_service_collaboration_scripts_have_no_owner_write_endpoints(self):
        scripts = (
            ("content-operations-workbench-v0", "install"),
            ("content-operations-workbench-v0", "uninstall"),
            ("smart-home-owner-pack-v0", "install"),
        )
        forbidden = (
            "/v3/base/gated-actions/confirm",
            "/v3/base/gated-actions/prepare",
            "/v3/pack-studio/lifecycle/confirm",
            "/v3/pack-studio/lifecycle/disable",
            "/v3/pack-studio/lifecycle/reactivate",
            "/v3/task-governance/schedules/approve",
            "/v3/task-governance/schedules/pause",
            "/v3/task-governance/schedules/resume",
        )
        for pack_dir, action in scripts:
            with open(os.path.join(REPO, pack_dir, action + ".py"), encoding="utf-8") as stream:
                source = stream.read()
            for endpoint in forbidden:
                self.assertNotIn(endpoint, source)
            self.assertNotRegex(source, r"call\(\s*['\"](?:POST|PUT|PATCH|DELETE)['\"]")
            self.assertIn('method != "GET"', source)
            self.assertIn("body is not None", source)
            if action == "install":
                self.assertIn("wait_for_owner_enabled", source)
            else:
                self.assertIn("wait_for_owner_disabled", source)


if __name__ == "__main__":
    unittest.main()
