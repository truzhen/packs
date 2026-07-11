#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""knowledge-index 内容 checksum 防漂移（统一决策表 #10）。

真相源永远是知识文件本体；checksum 只是 index 对文件字节的完整性声明，
不构成第二真相源。三向漂移检查：
  1) checksum_mismatch —— 文件内容被改而 index 不知情；
  2) missing_file / missing_checksum —— 条目断链或未声明完整性；
  3) unindexed_file —— knowledge/ 下存在未登记的 *.md（游离知识不入库也不该潜伏）。

用法：
  python3 knowledge_checksums.py --update [pack_dir ...]   # 重算并写回 index
  python3 knowledge_checksums.py --verify [pack_dir ...]   # 任一漂移 exit 1（CI 用）
不带 pack_dir 时自动发现仓内所有含 knowledge/knowledge-index.json 的 pack。
"""
import argparse
import glob
import hashlib
import json
import os
import sys

REPO_DIR = os.path.dirname(os.path.abspath(__file__))
CHECKSUM_PREFIX = "sha256:"


def file_checksum(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return CHECKSUM_PREFIX + h.hexdigest()


def find_pack_dirs(root):
    out = []
    for idx in sorted(glob.glob(os.path.join(root, "*", "knowledge", "knowledge-index.json"))):
        out.append(os.path.dirname(os.path.dirname(idx)))
    return out


def _index_path(pack_dir):
    return os.path.join(pack_dir, "knowledge", "knowledge-index.json")


def _load_index(pack_dir):
    with open(_index_path(pack_dir), encoding="utf-8") as f:
        return json.load(f)


def update_pack(pack_dir):
    """按 index 条目重算 checksum 写回；返回实际变更条数。"""
    doc = _load_index(pack_dir)
    changed = 0
    for e in doc.get("entries", []):
        fpath = os.path.join(pack_dir, e["file"])
        want = file_checksum(fpath)
        if e.get("checksum") != want:
            e["checksum"] = want
            changed += 1
    if changed:
        with open(_index_path(pack_dir), "w", encoding="utf-8") as f:
            json.dump(doc, f, ensure_ascii=False, indent=2)
            f.write("\n")
    return changed


def verify_entries(pack_dir, entries):
    """条目级两向检查（内容漂移 / 断链）。install.py 装入前复用本函数 fail-fast。"""
    problems = []
    for e in entries:
        rel = e.get("file", "")
        fpath = os.path.join(pack_dir, rel)
        if not os.path.exists(fpath):
            problems.append("missing_file %s" % rel)
            continue
        want = e.get("checksum", "")
        if not want:
            problems.append("missing_checksum %s" % rel)
            continue
        got = file_checksum(fpath)
        if got != want:
            problems.append("checksum_mismatch %s（index=%s 实际=%s）" % (rel, want, got))
    return problems


def verify_pack(pack_dir):
    """pack 级三向检查：条目两向 + knowledge/ 下未登记 *.md。"""
    doc = _load_index(pack_dir)
    entries = doc.get("entries", [])
    problems = verify_entries(pack_dir, entries)
    indexed = {e.get("file", "") for e in entries}
    for md in sorted(glob.glob(os.path.join(pack_dir, "knowledge", "**", "*.md"), recursive=True)):
        rel = os.path.relpath(md, pack_dir)
        if rel not in indexed:
            problems.append("unindexed_file %s" % rel)
    return problems


def main():
    ap = argparse.ArgumentParser(description="knowledge-index 内容 checksum 防漂移")
    mode = ap.add_mutually_exclusive_group(required=True)
    mode.add_argument("--update", action="store_true")
    mode.add_argument("--verify", action="store_true")
    ap.add_argument("pack_dirs", nargs="*")
    args = ap.parse_args()

    packs = [os.path.abspath(p) for p in args.pack_dirs] or find_pack_dirs(REPO_DIR)
    if not packs:
        print("未发现任何含 knowledge-index 的 pack")
        return 0

    bad = 0
    for p in packs:
        name = os.path.basename(p.rstrip("/"))
        if args.update:
            n = update_pack(p)
            print("%s: 更新 %d 条 checksum" % (name, n))
        else:
            problems = verify_pack(p)
            if problems:
                bad += 1
                for msg in problems:
                    print("::error::%s %s" % (name, msg))
            else:
                print("%s: checksum 全部一致" % name)
    if bad:
        print("知识内容与 index 漂移，先跑 --update 并复核改动来源", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
