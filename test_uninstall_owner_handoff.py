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


def lifecycle_entry(pack_ref, version="", records=None, include_pointer=True):
    if records is None:
        records = [] if not include_pointer else [{
            "pack_ref": pack_ref,
            "version": version,
            "state": "enabled" if version else "disabled",
        }]
    entry = {"pack_ref": pack_ref, "records": records}
    if include_pointer:
        entry["enabled_pointer"] = {"current_version": version}
    return entry


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
            return 200, {"packs": [lifecycle_entry(pack_ref, "1.0.0" if enabled else "")]}

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
            return 200, {"packs": [lifecycle_entry(pack_ref, enabled)]}

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
            return 200, {"packs": [lifecycle_entry(pack_ref, version)]}

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
                return 200, {"packs": [lifecycle_entry(pack_ref, current_version)]}

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
            body = {"packs": [lifecycle_entry(pack_ref, version) for version in versions]}
            self.assertIsNone(
                diagnostics.pack_enabled_version_from_readmodel(body, pack_ref),
                versions,
            )

    def test_lifecycle_helper_accepts_only_canonical_first_install_empty_state(self):
        diagnostics = load_diagnostics()
        pack_ref = "scene_pack://smart-home-owner-project-ops"
        self.assertEqual(
            diagnostics.pack_enabled_version_from_readmodel(
                {"packs": [lifecycle_entry(pack_ref, include_pointer=False)]}, pack_ref
            ),
            "",
        )
        self.assertEqual(
            diagnostics.pack_enabled_version_from_readmodel(
                {"packs": [lifecycle_entry(pack_ref, "", records=[{"state": "disabled"}])]}, pack_ref
            ),
            "",
        )
        self.assertEqual(
            diagnostics.pack_enabled_version_from_readmodel(
                {"packs": [lifecycle_entry(pack_ref, "1.1.0", records=[{"state": "enabled"}])]}, pack_ref
            ),
            "1.1.0",
        )
        self.assertIsNone(
            diagnostics.pack_enabled_version_from_readmodel({"packs": []}, pack_ref),
        )
        self.assertIsNone(
            diagnostics.pack_enabled_version_from_readmodel(
                {"packs": [lifecycle_entry("scene_pack://another-pack", "1.1.0")]},
                pack_ref,
            ),
        )
        malformed = (
            {"pack_ref": pack_ref},
            {"pack_ref": pack_ref, "records": {}},
            {"pack_ref": pack_ref, "records": ["record"]},
            lifecycle_entry(pack_ref, "1.1.0", records=[]),
            lifecycle_entry(pack_ref, "", records=[]),
            {"pack_ref": pack_ref, "records": [], "enabled_pointer": None},
            {"pack_ref": pack_ref, "records": [], "enabled_pointer": {}},
            {"pack_ref": pack_ref, "records": [{"state": "draft"}]},
        )
        for entry in malformed:
            self.assertIsNone(
                diagnostics.pack_enabled_version_from_readmodel({"packs": [entry]}, pack_ref),
                entry,
            )

    def test_all_consumers_accept_canonical_first_install_and_reactivate_states(self):
        for pack_dir, version in (
            ("content-operations-workbench-v0", "0.2.0"),
            ("smart-home-owner-pack-v0", "1.1.0"),
        ):
            module = load_script("first_install_" + pack_dir.replace("-", "_"), pack_dir, "install")
            with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
                manifest = module.json.load(stream)
            pack_ref = manifest["pack_ref"]
            for initial in (
                lifecycle_entry(pack_ref, include_pointer=False),
                lifecycle_entry(pack_ref),
            ):
                lifecycle_states = [initial, lifecycle_entry(pack_ref, version)]

                def fake_call(method, path, body=None):
                    if path == "/v3/task-governance/schedules":
                        schedule_doc = module.load_json(manifest["schedules_file"])
                        return 200, {"schedules": [{
                            "transaction_ref": module.schedule_transaction_ref(pack_ref, item["schedule_key"]),
                            "status": "active",
                        } for item in schedule_doc["schedules"]]}
                    return 200, {"packs": [lifecycle_states.pop(0) if lifecycle_states else lifecycle_entry(pack_ref, version)]}

                module.call = fake_call
                with mock.patch.object(sys, "argv", [
                    "install.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
                ]):
                    module.main()

        module = load_script("first_install_content_uninstall", "content-operations-workbench-v0", "uninstall")
        manifest = module.load_json("manifest.json")
        pack_ref = manifest["pack_ref"]

        def uninstall_call(method, path, body=None):
            if path == "/v3/task-governance/schedules":
                return 200, {"schedules": []}
            return 200, {"packs": [lifecycle_entry(pack_ref, include_pointer=False)]}

        module.call = uninstall_call
        with mock.patch.object(sys, "argv", [
            "uninstall.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
        ]):
            module.main()

    def test_all_consumers_reject_orphan_enabled_pointers_with_empty_records(self):
        targets = (
            ("content-operations-workbench-v0", "install", "0.2.0"),
            ("content-operations-workbench-v0", "uninstall", ""),
            ("smart-home-owner-pack-v0", "install", "1.1.0"),
        )
        for pack_dir, action, orphan_version in targets:
            module = load_script("orphan_pointer_" + pack_dir.replace("-", "_") + "_" + action, pack_dir, action)
            with open(os.path.join(module.PACK_DIR, "manifest.json"), encoding="utf-8") as stream:
                manifest = module.json.load(stream)
            pack_ref = manifest["pack_ref"]
            calls = []

            def fake_call(method, path, body=None):
                calls.append((method, path, body))
                if path == "/v3/task-governance/schedules":
                    return 200, {"schedules": []}
                return 200, {"packs": [lifecycle_entry(pack_ref, orphan_version, records=[])]}

            module.call = fake_call
            with mock.patch.object(sys, "argv", [
                action + ".py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
            ]):
                with self.assertRaises(SystemExit):
                    module.main()
            self.assertTrue(all(method == "GET" and body is None for method, _, body in calls), calls)

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
                    lifecycle_entry(pack_ref, canonical),
                    lifecycle_entry(pack_ref, conflict),
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
            return 200, {"packs": [lifecycle_entry(pack_ref)]}

        module.call = fake_call
        with mock.patch.object(sys, "argv", [
            "uninstall.py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
        ]):
            with self.assertRaises(SystemExit):
                module.main()

    def test_schedule_helper_rejects_duplicate_and_malformed_records(self):
        diagnostics = load_diagnostics()
        ref = diagnostics.schedule_transaction_ref(
            "scene_pack://content-operations-workbench", "weekday_direction_radar"
        )
        malformed = (
            [{"transaction_ref": ref, "status": "paused"}, {"transaction_ref": ref, "status": "active"}],
            [{"transaction_ref": ref, "status": "active"}, {"transaction_ref": ref, "status": "active"}],
            [{"transaction_ref": ref, "status": "active"}, "not-an-object"],
            [{"transaction_ref": ref}],
            [{"transaction_ref": ref, "status": "active "}],
            [{"transaction_ref": ref, "status": 1}],
            [{"transaction_ref": ref + " ", "status": "active"}],
        )
        for schedules in malformed:
            def fake_call(method, path, body=None):
                self.assertEqual((method, path, body), ("GET", "/v3/task-governance/schedules", None))
                return 200, {"schedules": schedules}

            ok, reason = diagnostics.wait_for_owner_schedule_states(
                fake_call, [ref], {"active"}, timeout_seconds=0
            )
            self.assertFalse(ok, schedules)
            self.assertEqual(reason, "schedule_readmodel_invalid", schedules)

        ok, reason = diagnostics.wait_for_owner_schedule_states(
            lambda *_args: self.fail("duplicate declared refs must fail before GET"),
            [ref, ref], {"active"}, timeout_seconds=0,
        )
        self.assertFalse(ok)
        self.assertEqual(reason, "schedule_readmodel_invalid")

    def test_content_consumers_reject_full_schedule_malformed_matrix(self):
        cases = (
            "duplicate_conflict",
            "duplicate_same",
            "non_object_sibling",
            "missing_status",
            "noncanonical_status",
            "non_string_status",
            "noncanonical_ref",
        )
        for action in ("install", "uninstall"):
            module = load_script("content_schedule_" + action, "content-operations-workbench-v0", action)
            manifest = module.load_json("manifest.json")
            pack_ref = manifest["pack_ref"]
            schedule_doc = module.load_json(manifest["schedules_file"])
            refs = [module.schedule_transaction_ref(pack_ref, item["schedule_key"]) for item in schedule_doc["schedules"]]
            target_status = "active" if action == "install" else "paused"
            conflicting_status = "paused" if action == "install" else "active"
            lifecycle_version = manifest["version"] if action == "install" else ""
            for case in cases:
                if case == "duplicate_conflict":
                    schedules = [{"transaction_ref": ref, "status": target_status} for ref in refs] + [{"transaction_ref": refs[0], "status": conflicting_status}]
                elif case == "duplicate_same":
                    schedules = [{"transaction_ref": ref, "status": target_status} for ref in refs] + [{"transaction_ref": refs[0], "status": target_status}]
                elif case == "non_object_sibling":
                    schedules = [{"transaction_ref": ref, "status": target_status} for ref in refs] + ["not-an-object"]
                elif case == "missing_status":
                    schedules = [{"transaction_ref": refs[0]}] + [{"transaction_ref": ref, "status": target_status} for ref in refs[1:]]
                elif case == "noncanonical_status":
                    schedules = [{"transaction_ref": refs[0], "status": target_status + " "}] + [{"transaction_ref": ref, "status": target_status} for ref in refs[1:]]
                elif case == "non_string_status":
                    schedules = [{"transaction_ref": refs[0], "status": 1}] + [{"transaction_ref": ref, "status": target_status} for ref in refs[1:]]
                else:
                    schedules = [{"transaction_ref": refs[0] + " ", "status": target_status}] + [{"transaction_ref": ref, "status": target_status} for ref in refs[1:]]

                def fake_call(method, path, body=None):
                    if path == "/v3/task-governance/schedules":
                        return 200, {"schedules": schedules}
                    return 200, {"packs": [lifecycle_entry(pack_ref, lifecycle_version)]}

                def die_as_runtime_error(message, *_args):
                    raise RuntimeError(message)

                module.call = fake_call
                module.die = die_as_runtime_error
                with mock.patch.object(sys, "argv", [
                    action + ".py", "--devserver-base", "http://127.0.0.1:18080", "--wait-seconds", "0",
                ]):
                    with self.assertRaisesRegex(RuntimeError, "schedule_readmodel_invalid"):
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
