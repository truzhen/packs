#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""F10 RC-B：Pack glue 只消费 Base confirm 返回的完整 issued binding。"""

import contextlib
import importlib.util
import io
import os
import sys
import unittest
from unittest import mock


REPO = os.path.dirname(os.path.abspath(__file__))
PACKS = [
    "backup-administrator-workbench-v0",
    "environmental-enforcement-pack-v0",
    "housekeeping-ops-pack-v0",
    "shuxuejia-renovation-pack-v0",
    "smart-home-owner-pack-v0",
]


def load_script(pack, name):
    path = os.path.join(REPO, pack, name + ".py")
    spec = importlib.util.spec_from_file_location(
        "%s_%s_under_test" % (pack.replace("-", "_"), name), path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def incomplete_confirm():
    return {"issue": {"decision_ref": "decision-1", "run_id": "run-1", "nonce": "nonce-1"}}


def complete_confirm():
    return {"issue": {
        "decision_ref": "decision-1", "run_id": "run-1", "nonce": "nonce-1",
        "owner_action_evidence_ref": "owner_action_evidence://base/confirmed-1",
    }}


class TestPackIssuedBinding(unittest.TestCase):
    def test_install_mint_rejects_confirm_without_evidence(self):
        for pack in PACKS:
            with self.subTest(pack=pack):
                mod = load_script(pack, "install")

                def fake_call(_method, path, _body=None):
                    if path.endswith("/prepare"):
                        return 200, {"issue": {"issue_ref": "issue-1"}}
                    if path.endswith("/confirm"):
                        return 200, incomplete_confirm()
                    self.fail("unexpected path: " + path)

                mod.call = fake_call
                with contextlib.redirect_stderr(io.StringIO()):
                    with self.assertRaises(SystemExit):
                        mod.mint_decision("pack-version", "knowledge-candidate")

    def test_install_mint_returns_exact_confirm_evidence(self):
        for pack in PACKS:
            with self.subTest(pack=pack):
                mod = load_script(pack, "install")

                def fake_call(_method, path, _body=None):
                    if path.endswith("/prepare"):
                        return 200, {"issue": {"issue_ref": "issue-1"}}
                    if path.endswith("/confirm"):
                        return 200, complete_confirm()
                    self.fail("unexpected path: " + path)

                mod.call = fake_call
                self.assertEqual(
                    mod.mint_decision("pack-version", "knowledge-candidate"),
                    ("decision-1", "run-1", "nonce-1", "owner_action_evidence://base/confirmed-1"),
                )

    def test_uninstall_rejects_confirm_without_evidence(self):
        for pack in PACKS:
            with self.subTest(pack=pack):
                mod = load_script(pack, "uninstall")

                def fake_call(_method, path, _body=None):
                    if "/lifecycle/packs?" in path:
                        return 200, {"packs": [{
                            "pack_ref": mod.json.load(open(os.path.join(mod.PACK_DIR, "manifest.json"), encoding="utf-8"))["pack_ref"],
                            "enabled_pointer": {"current_version": "0.1.0"},
                        }]}
                    if path.endswith("/prepare"):
                        return 200, {"issue": {"issue_ref": "issue-1"}}
                    if path.endswith("/confirm"):
                        return 200, incomplete_confirm()
                    self.fail("unexpected path: " + path)

                mod.call = fake_call
                with contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
                    with mock.patch.object(sys, "argv", [mod.__file__]):
                        with self.assertRaises(SystemExit):
                            mod.main()

    def test_uninstall_forwards_exact_confirm_evidence(self):
        for pack in PACKS:
            with self.subTest(pack=pack):
                mod = load_script(pack, "uninstall")
                manifest_path = os.path.join(mod.PACK_DIR, "manifest.json")
                with open(manifest_path, encoding="utf-8") as stream:
                    pack_ref = mod.json.load(stream)["pack_ref"]
                disabled = []
                lifecycle_calls = [0]

                # Read-only handoff is an explicit production declaration.
                # Do not probe production call() and do not swallow an unknown
                # exception: an undeclared mode must follow the legacy contract.
                is_read_only = hasattr(mod, "OWNER_DISABLE_HANDOFF")

                def fake_call(_method, path, body=None):
                    if "/lifecycle/packs?" in path:
                        lifecycle_calls[0] += 1
                        # First call: pack is enabled. Second call (from
                        # wait_for_owner_disabled): pack is disabled.
                        if lifecycle_calls[0] == 1:
                            return 200, {"packs": [{"pack_ref": pack_ref, "enabled_pointer": {"current_version": "0.1.0"}}]}
                        return 200, {"packs": [{"pack_ref": pack_ref, "enabled_pointer": {"current_version": ""}}]}
                    if path.endswith("/prepare"):
                        return 200, {"issue": {"issue_ref": "issue-1"}}
                    if path.endswith("/confirm"):
                        return 200, complete_confirm()
                    if path.endswith("/lifecycle/disable"):
                        disabled.append(body)
                        return 200, {}
                    self.fail("unexpected path: " + path)

                mod.call = fake_call
                # The production uninstall.py requires an explicit devserver base;
                # tests must set it before invoking main(). Also zero the owner
                # handoff wait so the test does not sleep for the default 300s.
                old_base = os.environ.get("TRUZHEN_DEVSERVER_BASE")
                old_wait = os.environ.get("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS")
                os.environ["TRUZHEN_DEVSERVER_BASE"] = "http://127.0.0.1:18080"
                os.environ["TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS"] = "0"
                try:
                    with contextlib.redirect_stdout(io.StringIO()):
                        with mock.patch.object(sys, "argv", [mod.__file__]):
                            mod.main()
                finally:
                    if old_base is None:
                        os.environ.pop("TRUZHEN_DEVSERVER_BASE", None)
                    else:
                        os.environ["TRUZHEN_DEVSERVER_BASE"] = old_base
                    if old_wait is None:
                        os.environ.pop("TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS", None)
                    else:
                        os.environ["TRUZHEN_OWNER_HANDOFF_WAIT_SECONDS"] = old_wait
                # Read-only uninstall.py (environmental/smart-home) never POSTs
                # /lifecycle/disable; it waits for the Owner to disable via GUI.
                # Legacy uninstall.py (backup/housekeeping/shuxuejia) POSTs and
                # must forward the exact confirm evidence.
                if is_read_only:
                    self.assertEqual(len(disabled), 0, "read-only uninstall must not POST /lifecycle/disable")
                else:
                    self.assertEqual(len(disabled), 1)
                    self.assertEqual(disabled[0]["owner_action_evidence_ref"], "owner_action_evidence://base/confirmed-1")


if __name__ == "__main__":
    unittest.main()
