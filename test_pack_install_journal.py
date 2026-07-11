#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""pack_install_journal 行为级测试（统一决策表 #8：install 事务日志+断点续装）。

主权边界：journal 只是本机恢复辅助账目，不是真相源；不提供任何「自动反做
正式域对象」能力——撤销走 uninstall / Owner 侧禁用状态机。
"""
import json
import os
import shutil
import tempfile
import unittest

from pack_install_journal import InstallJournal


class InstallJournalTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp(prefix="journal-test-")

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)

    def _open(self):
        return InstallJournal.open("pack://demo/pack", "http://127.0.0.1:18080",
                                   state_dir=self.tmp)

    def _raw(self):
        path = InstallJournal.path_for("pack://demo/pack", state_dir=self.tmp)
        with open(path, encoding="utf-8") as f:
            return json.load(f)

    def test_fresh_journal_persists_every_mutation(self):
        j = self._open()
        self.assertEqual("in_progress", self._raw()["status"])
        j.step("scene")
        j.set_version("0.1.0")
        j.mark("scene_enabled")
        raw = self._raw()
        self.assertEqual("scene", raw["current_step"])
        self.assertEqual("0.1.0", raw["resolved_version"])
        self.assertIn("scene_enabled", raw["marks"])

    def test_partial_failure_records_step_error_and_approved_refs(self):
        j = self._open()
        j.step("knowledge")
        j.mark_item("knowledge_batches", "knowledge_scope://demo/sop")
        j.add_approved("knowledge_scope://demo/sop", "candidate://k/1")
        j.fail(error_code="TZ-PACK-INSTALL-007", message="approve HTTP 500")
        raw = self._raw()
        self.assertEqual("failed", raw["status"])
        self.assertEqual("knowledge", raw["last_error"]["step"])
        self.assertEqual("TZ-PACK-INSTALL-007", raw["last_error"]["error_code"])
        self.assertEqual(["candidate://k/1"],
                         raw["approved"]["knowledge_scope://demo/sop"])

    def test_resume_reports_prior_state_and_skips_approved(self):
        j = self._open()
        j.step("role_packs")
        j.mark_item("role_packs", "role://demo/advisor")
        j.step("knowledge")
        j.add_approved("scope-a", "candidate://k/1")
        j.fail(error_code="TZ-PACK-INSTALL-007", message="boom")

        j2 = self._open()  # 续装：同 pack 重开
        report = j2.prior_report()
        self.assertTrue(any("failed" in line or "上次" in line for line in report), report)
        self.assertTrue(any("role://demo/advisor" in line for line in report), report)
        self.assertIn("candidate://k/1", j2.approved_refs("scope-a"))
        self.assertEqual("in_progress", self._raw()["status"])
        # 半装报告不得声称可自动反做正式域
        self.assertFalse(any("自动回滚" in line or "自动撤销" in line for line in report), report)

    def test_complete_marks_journal_completed(self):
        j = self._open()
        j.step("scene")
        j.complete()
        self.assertEqual("completed", self._raw()["status"])
        # 完成后重开 = 新一次安装，之前的 completed 不算半装
        j2 = self._open()
        self.assertEqual([], j2.prior_report())

    def test_fail_without_step_uses_unknown(self):
        j = self._open()
        j.fail(error_code="TZ-PACK-INSTALL-002", message="连不上")
        self.assertEqual("unknown", self._raw()["last_error"]["step"])

    def test_state_dir_env_override(self):
        os.environ["TRUZHEN_PACK_INSTALL_STATE_DIR"] = self.tmp
        try:
            j = InstallJournal.open("pack://demo/env", "http://x")
            self.assertTrue(InstallJournal.path_for("pack://demo/env").startswith(self.tmp))
            j.complete()
        finally:
            del os.environ["TRUZHEN_PACK_INSTALL_STATE_DIR"]


if __name__ == "__main__":
    unittest.main()
