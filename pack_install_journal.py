#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Pack 装入事务日志（统一决策表 #8：install 原子边界 · 断点续装 · 半装显式化）。

主权边界（不可越）：
  - journal 是**本机恢复辅助账目**，不是真相源；安装进度真相源永远是 os 侧
    各状态机（lifecycle / 角色包 / 绑槽 / 09 candidates），冲突时以端点实况为准。
  - journal 不提供、也永远不得提供「自动反做正式域对象」能力：已 approve 的
    FormalKnowledge、已 enabled 的场景包/角色包属正式域，撤销只能走 uninstall /
    Owner 侧禁用状态机（真 Base Gate + receipt）。
  - 失败时的职责只有三件：把半装状态写成机器可读事实、向用户显式报告
    「哪些已成正式 / 哪些悬空」、给出续装指引（重跑同命令，幂等续装）。

落盘：默认 ~/.truzhen/pack-install-journals/<pack>.json，可用
TRUZHEN_PACK_INSTALL_STATE_DIR 覆盖；原子写（tmp+rename），每次变更即刷。
"""
import json
import os
import re
import tempfile

_DEFAULT_DIR = os.path.join(os.path.expanduser("~"), ".truzhen", "pack-install-journals")


def _sanitize(ref):
    return re.sub(r"[/:@\s]+", "-", str(ref)).strip("-")


class InstallJournal(object):
    def __init__(self, path, doc):
        self._path = path
        self._doc = doc
        self._prior_report = []

    # ---------- 构造 ----------

    @staticmethod
    def state_dir(state_dir=None):
        return state_dir or os.environ.get("TRUZHEN_PACK_INSTALL_STATE_DIR") or _DEFAULT_DIR

    @classmethod
    def path_for(cls, pack_ref, state_dir=None):
        return os.path.join(cls.state_dir(state_dir), _sanitize(pack_ref) + ".json")

    @classmethod
    def open(cls, pack_ref, base, state_dir=None):
        """打开（或续接）本 pack 的安装日志。上次未完成时生成半装报告。"""
        path = cls.path_for(pack_ref, state_dir)
        prior = None
        if os.path.exists(path):
            try:
                with open(path, encoding="utf-8") as f:
                    prior = json.load(f)
            except Exception:
                prior = None  # 日志损坏不阻塞安装：以端点实况为准重建
        doc = {
            "pack_ref": pack_ref,
            "base": base,
            "status": "in_progress",
            "current_step": "",
            "resolved_version": "",
            "marks": [],
            "items": {},      # 组 -> [已完成项]，如 role_packs / bindings / knowledge_batches
            "approved": {},   # knowledge scope_ref -> [已 approve 的 candidate_ref]
            "last_error": None,
        }
        report = []
        if prior and prior.get("status") in ("in_progress", "failed"):
            # 断点续装：保留已完成账目（幂等续装的依据），重置状态与错误。
            doc["resolved_version"] = prior.get("resolved_version", "")
            doc["marks"] = list(prior.get("marks", []))
            doc["items"] = {k: list(v) for k, v in (prior.get("items") or {}).items()}
            doc["approved"] = {k: list(v) for k, v in (prior.get("approved") or {}).items()}
            report.append("上次装入未完成（status=%s），本次为断点续装。" % prior.get("status"))
            err = prior.get("last_error") or {}
            if err:
                report.append("上次失败点：step=%s code=%s %s" % (
                    err.get("step", "unknown"), err.get("error_code", ""),
                    (err.get("message") or "")[:200]))
            for mark in doc["marks"]:
                report.append("已完成：%s" % mark)
            for group, items in doc["items"].items():
                for it in items:
                    report.append("已完成：%s %s" % (group, it))
            for scope, refs in doc["approved"].items():
                for r in refs:
                    report.append("已成正式（不自动反做，撤销走 uninstall/Owner 禁用）：%s %s" % (scope, r))
        j = cls(path, doc)
        j._prior_report = report
        j._flush()
        if report:
            print("↻ 断点续装报告：")
            for line in report:
                print("    " + line)
        return j

    # ---------- 查询 ----------

    def prior_report(self):
        return list(self._prior_report)

    def approved_refs(self, scope_ref):
        return set(self._doc["approved"].get(scope_ref, []))

    # ---------- 记账 ----------

    def step(self, name):
        self._doc["current_step"] = name
        self._flush()

    def set_version(self, version):
        self._doc["resolved_version"] = version
        self._flush()

    def mark(self, key):
        if key not in self._doc["marks"]:
            self._doc["marks"].append(key)
        self._flush()

    def mark_item(self, group, item_key):
        items = self._doc["items"].setdefault(group, [])
        if item_key not in items:
            items.append(item_key)
        self._flush()

    def add_approved(self, scope_ref, candidate_ref):
        refs = self._doc["approved"].setdefault(scope_ref, [])
        if candidate_ref not in refs:
            refs.append(candidate_ref)
        self._flush()

    def fail(self, *, error_code, message):
        self._doc["status"] = "failed"
        self._doc["last_error"] = {
            "step": self._doc.get("current_step") or "unknown",
            "error_code": error_code,
            "message": str(message)[:2048],
        }
        self._flush()
        print("✗ 半装状态已落盘：%s" % self._path)
        print("    恢复：重跑同一条 install 命令即断点续装（幂等）。")
        print("    撤销已成正式的对象：走 uninstall / Owner 侧禁用状态机，本日志不自动反做。")

    def complete(self):
        self._doc["status"] = "completed"
        self._doc["current_step"] = ""
        self._flush()

    # ---------- 落盘 ----------

    def _flush(self):
        d = os.path.dirname(self._path)
        os.makedirs(d, exist_ok=True)
        fd, tmp = tempfile.mkstemp(prefix=".journal-", dir=d)
        try:
            with os.fdopen(fd, "w", encoding="utf-8") as f:
                json.dump(self._doc, f, ensure_ascii=False, indent=2)
                f.write("\n")
            os.replace(tmp, self._path)
        finally:
            if os.path.exists(tmp):
                os.unlink(tmp)
