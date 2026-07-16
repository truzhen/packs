#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""F10 RC-B：Pack glue 只消费 Base confirm 返回的完整 issued binding。"""

import contextlib
import importlib.util
import io
import os
import unittest


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

                def fake_call(_method, path, body=None):
                    if "/lifecycle/packs?" in path:
                        return 200, {"packs": [{"pack_ref": pack_ref, "enabled_pointer": {"current_version": "0.1.0"}}]}
                    if path.endswith("/prepare"):
                        return 200, {"issue": {"issue_ref": "issue-1"}}
                    if path.endswith("/confirm"):
                        return 200, complete_confirm()
                    if path.endswith("/lifecycle/disable"):
                        disabled.append(body)
                        return 200, {}
                    self.fail("unexpected path: " + path)

                mod.call = fake_call
                with contextlib.redirect_stdout(io.StringIO()):
                    mod.main()
                self.assertEqual(len(disabled), 1)
                self.assertEqual(disabled[0]["owner_action_evidence_ref"], "owner_action_evidence://base/confirmed-1")


if __name__ == "__main__":
    unittest.main()
