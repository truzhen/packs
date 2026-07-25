#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""墅学家 uninstall proof 的本地 fail-closed 约束测试（不连接 devserver）。"""

import importlib.util
import json
import os
import pathlib
import unittest


PACK = pathlib.Path(__file__).resolve().parent
SCRIPT = PACK / "uninstall.py"
PACK_REF = "scene_pack://shuxuejia-large-home-renovation"
VERSION = "1.0.0"


def load_module():
    spec = importlib.util.spec_from_file_location("shuxuejia_uninstall_proof_test", SCRIPT)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class ProofRejected(Exception):
    pass


class UninstallProofTest(unittest.TestCase):
    def setUp(self):
        self.module = load_module()
        self.original_die = self.module.die
        self.original_proof = os.environ.get(self.module.PROOF_ENV)
        self.module.die = lambda message, _code: (_ for _ in ()).throw(ProofRejected(message))

    def tearDown(self):
        self.module.die = self.original_die
        if self.original_proof is None:
            os.environ.pop(self.module.PROOF_ENV, None)
        else:
            os.environ[self.module.PROOF_ENV] = self.original_proof

    def proof(self, **override):
        value = {
            "action_type": self.module.UNINSTALL_ACTION,
            "target_ref": PACK_REF,
            "transaction_ref": "transaction://pack-uninstall:%s@%s" % (PACK_REF, VERSION),
            "decision_ref": "decision://test/issued",
            "run_id": "run-test",
            "nonce": "nonce-test",
            "owner_action_evidence_ref": "evidence://owner/test",
        }
        value.update(override)
        return json.dumps(value)

    def load(self, raw):
        if raw is None:
            os.environ.pop(self.module.PROOF_ENV, None)
        else:
            os.environ[self.module.PROOF_ENV] = raw
        return self.module.load_uninstall_proof(PACK_REF, VERSION)

    def test_valid_external_proof_is_accepted_without_network(self):
        self.assertEqual("nonce-test", self.load(self.proof())["nonce"])

    def test_missing_proof_is_rejected(self):
        with self.assertRaises(ProofRejected):
            self.load(None)

    def test_malformed_proof_is_rejected(self):
        with self.assertRaises(ProofRejected):
            self.load("{")

    def test_wrong_action_is_rejected(self):
        with self.assertRaises(ProofRejected):
            self.load(self.proof(action_type="14.pack-studio.lifecycle.disable"))

    def test_script_only_declares_formal_uninstall_without_legacy_paths(self):
        source = SCRIPT.read_text(encoding="utf-8")
        self.assertIn('"action_type": "14.pack-studio.lifecycle.uninstall"', source)
        self.assertNotIn('"action_type": "14.pack-studio.lifecycle.disable"', source)
        self.assertIn('"/v3/pack-studio/lifecycle/uninstall"', source)
        for forbidden in (
            '"/v3/pack-studio/lifecycle/disable"',
            '"/v3/base/gated-actions/prepare"',
            '"/v3/base/gated-actions/confirm"',
            "owner_action_evidence://",
        ):
            self.assertNotIn(forbidden, source)

    def test_wrong_target_is_rejected(self):
        with self.assertRaises(ProofRejected):
            self.load(self.proof(target_ref="scene_pack://other"))

    def test_wrong_transaction_is_rejected(self):
        with self.assertRaises(ProofRejected):
            self.load(self.proof(transaction_ref="transaction://pack-uninstall:wrong@1.0.0"))


if __name__ == "__main__":
    unittest.main()
