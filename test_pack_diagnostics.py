#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Y4 行为测试：证明 install 脚本在不同失败阶段发出对应细分错误码（非仅结构存在）。

无需 devserver：monkeypatch install 模块的 call()，把主流程逼进各失败分支，
捕获 stderr 的 TRUZHEN_PACK_ERROR JSON，断言 error_code 与阶段一致。

运行：python3 -m unittest test_pack_diagnostics   （在仓根，Go test 之外的行为级验证）
"""
import contextlib
import importlib.util
import io
import json
import os
import shutil
import tempfile
import unittest

REPO = os.path.dirname(os.path.abspath(__file__))
PACK = os.path.join(REPO, "housekeeping-ops-pack-v0")
ENV_PACK = os.path.join(REPO, "environmental-enforcement-pack-v0")


def load_install_module():
    """按路径加载带连字符目录里的 install.py（不能用普通 import）。每次全新加载。"""
    spec = importlib.util.spec_from_file_location(
        "hk_install_under_test", os.path.join(PACK, "install.py"))
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod


def run_and_capture_error(mod):
    """跑 main()，期望 die→SystemExit(1)，返回解析出的 TRUZHEN_PACK_ERROR payload。"""
    err = io.StringIO()
    with contextlib.redirect_stderr(err):
        with contextlib.suppress(SystemExit):
            mod.main()
    for line in err.getvalue().splitlines():
        if line.startswith("TRUZHEN_PACK_ERROR "):
            return json.loads(line[len("TRUZHEN_PACK_ERROR "):])
    raise AssertionError("未捕获到 TRUZHEN_PACK_ERROR，stderr=\n" + err.getvalue())


class TestPackDiagnosticsRegistry(unittest.TestCase):
    def test_emit_carries_given_code_and_action(self):
        import pack_diagnostics as pd
        err = io.StringIO()
        with contextlib.redirect_stderr(err):
            pd.emit_pack_error(pack_dir=PACK, base="http://127.0.0.1:18080",
                               action="install", error_code=pd.INSTALL_KNOWLEDGE,
                               message="boom")
        payload = json.loads(err.getvalue().split("TRUZHEN_PACK_ERROR ", 1)[1])
        self.assertEqual(payload["error_code"], "TZ-PACK-INSTALL-007")
        self.assertEqual(payload["action"], "install")

    def test_registry_membership_and_shape(self):
        import pack_diagnostics as pd
        self.assertTrue(pd.is_registered_code(pd.INSTALL_CONNECTIVITY))
        self.assertTrue(pd.is_registered_code(pd.UNINSTALL_LIFECYCLE_HTTP))
        self.assertFalse(pd.is_registered_code("TZ-PACK-INSTALL-999"))  # 未登记
        self.assertFalse(pd.is_registered_code("bogus"))               # 形状非法
        install = [c for c in pd.PACK_ERROR_CODES if "-INSTALL-" in c]
        uninstall = [c for c in pd.PACK_ERROR_CODES if "-UNINSTALL-" in c]
        self.assertGreaterEqual(len(install), 8)
        self.assertGreaterEqual(len(uninstall), 3)


class TestInstallStageErrorCodes(unittest.TestCase):
    def setUp(self):
        self._journal_tmp = tempfile.mkdtemp(prefix="journal-diag-")
        os.environ["TRUZHEN_PACK_INSTALL_STATE_DIR"] = self._journal_tmp

    def tearDown(self):
        shutil.rmtree(self._journal_tmp, ignore_errors=True)
        os.environ.pop("TRUZHEN_PACK_INSTALL_STATE_DIR", None)

    def test_connectivity_failure_maps_to_002(self):
        mod = load_install_module()
        # call() 恒返回传输层错误（code 0）→ 健康检查即判连不上。
        mod.call = lambda method, path, body=None: (0, {"_transport_error": "refused"})
        payload = run_and_capture_error(mod)
        self.assertEqual(payload["error_code"], mod.INSTALL_CONNECTIVITY)
        self.assertEqual(payload["error_code"], "TZ-PACK-INSTALL-002")

    def test_readiness_not_ready_maps_to_004(self):
        mod = load_install_module()

        def fake_call(method, path, body=None):
            if path.startswith("/v3/pack-studio/lifecycle/packs"):
                return 200, {"packs": []}                     # 未装 → enabled_version=None
            if path.endswith("/reactivate"):
                return 409, {}                                # reactivate 失败 → 走 do_seed_scene
            if path.endswith("/canvas"):
                return 200, {"engine_sync": {"synced": True}} # canvas 同步成功
            if path.endswith("/draft"):
                return 200, {}                                # draft 成功
            if path.endswith("/readiness"):
                return 200, {"record": {"readiness_report": {"ready": False}}}  # 六件事未就绪
            return 200, {}
        mod.call = fake_call
        payload = run_and_capture_error(mod)
        self.assertEqual(payload["error_code"], mod.INSTALL_READINESS)
        self.assertEqual(payload["error_code"], "TZ-PACK-INSTALL-004")


class TestEnvironmentalHighRiskLifecycle(unittest.TestCase):
    def test_manifest_high_risk_is_forwarded_to_lifecycle_draft(self):
        with open(os.path.join(ENV_PACK, "manifest.json"), encoding="utf-8") as f:
            manifest = json.load(f)
        with open(os.path.join(ENV_PACK, "install.py"), encoding="utf-8") as f:
            install_source = f.read()

        self.assertEqual(manifest.get("risk_level"), "high")
        self.assertIn('"risk_level": manifest.get("risk_level", "medium")', install_source)
        self.assertIn('"verify_authority": False', install_source)
        self.assertNotIn('"verify_authority": True', install_source)

    def test_enforcement_elite_grounding_is_declared_and_forwarded(self):
        with open(os.path.join(ENV_PACK, "role-packs", "enforcement-elite.rolepack.json"), encoding="utf-8") as f:
            role_pack = json.load(f)
        with open(os.path.join(ENV_PACK, "install.py"), encoding="utf-8") as f:
            install_source = f.read()

        self.assertIn("涉嫌超标排放", role_pack.get("scenario", ""))
        self.assertIn("最终决定需您裁定", role_pack.get("opening_line_candidate", ""))
        self.assertEqual(2, len(role_pack.get("example_dialogues", [])))
        self.assertIn('"scenario": rp.get("scenario", "")', install_source)
        self.assertIn('"opening_line_candidate": rp.get("opening_line_candidate", "")', install_source)
        self.assertIn('"example_dialogues": rp.get("example_dialogues", None)', install_source)


if __name__ == "__main__":
    unittest.main()
