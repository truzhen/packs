#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""#8 行为级闭环：install.py × 事务日志 的断点续装真流程（无需 devserver）。

用 shuxuejia pack（带 5 条知识）进程内加载 install.py，monkeypatch call()：
  轮 1：知识 approve 第 3 次注入 500 → die(007) → journal=failed，已批 2 条入账；
  轮 2：全通 fake 重跑 → 断点续装：已批 2 条不再重复过 Base gate，剩余 3 条补齐，
        journal=completed。
主权断言：journal 从不出现「自动反做」；失败提示指向 uninstall/Owner 禁用。
"""
import contextlib
import importlib.util
import io
import json
import os
import shutil
import tempfile
import unittest
import urllib.parse

from pack_install_journal import InstallJournal

REPO = os.path.dirname(os.path.abspath(__file__))
PACK = os.path.join(REPO, "shuxuejia-renovation-pack-v0")
PACK_REF = json.load(open(os.path.join(PACK, "manifest.json"), encoding="utf-8"))["pack_ref"]


def load_install_module():
    spec = importlib.util.spec_from_file_location(
        "sx_install_under_test", os.path.join(PACK, "install.py"))
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod


def make_fake_call(approve_log, fail_at_approve_n=None):
    def fake_call(method, path, body=None):
        if path.startswith("/v3/pack-studio/lifecycle/packs"):
            return 200, {"packs": []}
        if path.endswith("/lifecycle/reactivate"):
            return 409, {}
        if path.endswith("/pack-studio/canvas"):
            return 200, {"engine_sync": {"synced": True}}
        if path.endswith("/lifecycle/draft"):
            return 200, {}
        if path.endswith("/lifecycle/readiness"):
            return 200, {"record": {"readiness_report": {"ready": True}}}
        if path.endswith("/lifecycle/promote"):
            return 200, {}
        if path.endswith("/lifecycle/confirm"):
            return 200, {"status": "enabled"}
        if path.endswith("/role-packs/readmodel"):
            return 200, {"enabled_versions": []}
        if path.endswith("/role-packs/drafts"):
            return 200, {"draft_id": (body or {}).get("draft_id", "d")}
        if path.endswith("/readiness-check"):
            return 200, {}
        if path.endswith("/promote-candidate"):
            return 200, {"ok": True}
        if path.endswith("/enable-candidate"):
            return 200, {"ok": True}
        if path.endswith("/enable-confirm"):
            return 200, {"status": "enabled"}
        if path.endswith("/agent-slots/readmodel"):
            return 200, {"agent_slot_bindings": []}
        if path.endswith("/bind-candidate"):
            return 200, {"ok": True, "binding_ref": "b1"}
        if path.endswith("/agent-slots/confirm"):
            return 200, {"status": "enabled"}
        if path.endswith("/memory/knowledge/batches"):
            cands = [{"candidate_ref": "candidate://" + sf["source_ref"], "status": "pending"}
                     for sf in (body or {}).get("source_files", [])]
            return 200, {"candidates": cands}
        if path.endswith("/gated-actions/prepare"):
            return 200, {"issue": {"issue_ref": "i1"}}
        if path.endswith("/gated-actions/confirm"):
            return 200, {"issue": {
                "decision_ref": "dr", "run_id": "r", "nonce": "n",
                "owner_action_evidence_ref": "owner_action_evidence://base/test-confirm",
            }}
        if "/memory/knowledge/candidates/" in path and path.endswith("/approve"):
            cref = urllib.parse.unquote(
                path.split("/memory/knowledge/candidates/")[1].rsplit("/approve", 1)[0])
            if fail_at_approve_n is not None and len(approve_log) + 1 == fail_at_approve_n:
                return 500, {"error": "injected approve failure"}
            approve_log.append(cref)
            return 200, {}
        raise AssertionError("fake_call 未覆盖的路径：" + path)
    return fake_call


def run_main(mod):
    out, err = io.StringIO(), io.StringIO()
    with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
        with contextlib.suppress(SystemExit):
            mod.main()
    return out.getvalue(), err.getvalue()


class InstallJournalWiringTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp(prefix="journal-wiring-")
        os.environ["TRUZHEN_PACK_INSTALL_STATE_DIR"] = self.tmp

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)
        os.environ.pop("TRUZHEN_PACK_INSTALL_STATE_DIR", None)

    def _journal_raw(self):
        with open(InstallJournal.path_for(PACK_REF), encoding="utf-8") as f:
            return json.load(f)

    def test_partial_failure_then_resume_skips_approved(self):
        # 轮 1：第 3 次 approve 注入失败 → 半装
        mod = load_install_module()
        log1 = []
        mod.call = make_fake_call(log1, fail_at_approve_n=3)
        out1, err1 = run_main(mod)
        self.assertIn("TZ-PACK-INSTALL-007", err1)
        raw = self._journal_raw()
        self.assertEqual("failed", raw["status"])
        self.assertEqual("knowledge", raw["last_error"]["step"])
        approved_round1 = [r for refs in raw["approved"].values() for r in refs]
        self.assertEqual(sorted(log1), sorted(approved_round1),
                         "journal 入账必须与真实 approve 成功集一致")
        self.assertEqual(2, len(approved_round1))
        # 半装提示必须指向受控撤销，不得声称自动反做
        self.assertIn("不自动反做", out1)

        # 轮 2：全通重跑 → 断点续装
        mod2 = load_install_module()
        log2 = []
        mod2.call = make_fake_call(log2)
        out2, _ = run_main(mod2)
        self.assertIn("断点续装报告", out2)
        for done in approved_round1:
            self.assertNotIn(done, log2, "已批候选不得重复走 Base gate")
        raw2 = self._journal_raw()
        self.assertEqual("completed", raw2["status"])
        approved_total = [r for refs in raw2["approved"].values() for r in refs]
        self.assertEqual(5, len(approved_total), "5 条知识最终全部入账")
        self.assertEqual(3, len(log2), "续装只补剩余 3 条")

    def test_clean_run_completes_journal(self):
        mod = load_install_module()
        log = []
        mod.call = make_fake_call(log)
        out, err = run_main(mod)
        self.assertNotIn("TRUZHEN_PACK_ERROR", err)
        raw = self._journal_raw()
        self.assertEqual("completed", raw["status"])
        self.assertIn("scene_enabled", raw["marks"])
        self.assertEqual(5, len(log))


if __name__ == "__main__":
    unittest.main()
