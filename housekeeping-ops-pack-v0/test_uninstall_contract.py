"""Guard the Pack-owned formal uninstall protocol.

The real behavior is exercised against the fixed OS devserver in G21.  This
small Pack-local test prevents the script from silently regressing to the
weaker disable/reactivate lifecycle.
"""
from pathlib import Path
import unittest


class FormalUninstallContractTest(unittest.TestCase):
    def test_uninstall_requires_governed_uninstall_proof(self):
        source = Path(__file__).with_name("uninstall.py").read_text(encoding="utf-8")
        self.assertIn("TRUZHEN_PACK_UNINSTALL_PROOF_JSON", source)
        self.assertIn("14.pack-studio.lifecycle.uninstall", source)
        self.assertIn("/v3/pack-studio/lifecycle/uninstall", source)
        self.assertNotIn("14.pack-studio.lifecycle.disable", source)
        self.assertNotIn("/v3/pack-studio/lifecycle/disable", source)
        self.assertNotIn("/v3/base/gated-actions/prepare", source)
        self.assertNotIn("/v3/base/gated-actions/confirm", source)


if __name__ == "__main__":
    unittest.main()
