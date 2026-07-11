#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""knowledge_checksums 行为级测试（统一决策表 #10）。

覆盖：update→verify 闭环、内容漂移、缺 checksum、条目文件缺失、
knowledge/ 下未登记 *.md、真实 pack 全量 verify 干净。
"""
import json
import os
import shutil
import tempfile
import unittest

import knowledge_checksums as kc

REPO_DIR = os.path.dirname(os.path.abspath(__file__))


def _write(path, text):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(text)


class KnowledgeChecksumTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp(prefix="kc-test-")
        self.pack = os.path.join(self.tmp, "fake-pack-v0")
        _write(os.path.join(self.pack, "knowledge", "sop", "a.md"), "# A\n正文A\n")
        _write(os.path.join(self.pack, "knowledge", "sop", "b.md"), "# B\n正文B\n")
        index = {
            "count": 2,
            "entries": [
                {"file": "knowledge/sop/a.md", "source_ref": "source://fake/a",
                 "title": "A", "kind": "sop",
                 "knowledge_scope_ref": "knowledge_scope://fake/sop"},
                {"file": "knowledge/sop/b.md", "source_ref": "source://fake/b",
                 "title": "B", "kind": "sop",
                 "knowledge_scope_ref": "knowledge_scope://fake/sop"},
            ],
        }
        _write(os.path.join(self.pack, "knowledge", "knowledge-index.json"),
               json.dumps(index, ensure_ascii=False, indent=2) + "\n")

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)

    def _entries(self):
        with open(os.path.join(self.pack, "knowledge", "knowledge-index.json"), encoding="utf-8") as f:
            return json.load(f)["entries"]

    def test_verify_before_update_reports_missing_checksum(self):
        problems = kc.verify_pack(self.pack)
        self.assertEqual(2, len([p for p in problems if p.startswith("missing_checksum")]),
                         "未生成 checksum 前 verify 必须逐条报 missing_checksum：%s" % problems)

    def test_update_then_verify_green(self):
        updated = kc.update_pack(self.pack)
        self.assertEqual(2, updated)
        entries = self._entries()
        for e in entries:
            self.assertTrue(e["checksum"].startswith("sha256:"), e)
        self.assertEqual([], kc.verify_pack(self.pack))
        # update 幂等：内容没变第二次不再改
        self.assertEqual(0, kc.update_pack(self.pack))

    def test_content_drift_detected(self):
        kc.update_pack(self.pack)
        _write(os.path.join(self.pack, "knowledge", "sop", "a.md"), "# A\n被篡改的正文\n")
        problems = kc.verify_pack(self.pack)
        self.assertTrue(any(p.startswith("checksum_mismatch knowledge/sop/a.md") for p in problems),
                        "内容漂移必须被抓：%s" % problems)

    def test_missing_entry_file_detected(self):
        kc.update_pack(self.pack)
        os.remove(os.path.join(self.pack, "knowledge", "sop", "b.md"))
        problems = kc.verify_pack(self.pack)
        self.assertTrue(any(p.startswith("missing_file knowledge/sop/b.md") for p in problems), problems)

    def test_unindexed_md_detected(self):
        kc.update_pack(self.pack)
        _write(os.path.join(self.pack, "knowledge", "sop", "stray.md"), "# 野文件\n")
        problems = kc.verify_pack(self.pack)
        self.assertTrue(any(p.startswith("unindexed_file knowledge/sop/stray.md") for p in problems), problems)

    def test_verify_entries_reusable_for_install(self):
        """install.py 复用的 verify_entries 只做内容/缺失两向（不做 unindexed，装入按 index 为准）。"""
        kc.update_pack(self.pack)
        _write(os.path.join(self.pack, "knowledge", "sop", "stray.md"), "# 野文件\n")
        self.assertEqual([], kc.verify_entries(self.pack, self._entries()))

    def test_real_packs_all_verify_clean(self):
        """仓内所有带 knowledge-index 的真实 pack 必须全量 verify 干净（CI 同口径）。"""
        packs = kc.find_pack_dirs(REPO_DIR)
        self.assertGreaterEqual(len(packs), 2, "至少应发现环保/shuxuejia 两个知识型 pack")
        for p in packs:
            self.assertEqual([], kc.verify_pack(p), "真实 pack 漂移：%s" % p)


if __name__ == "__main__":
    unittest.main()
