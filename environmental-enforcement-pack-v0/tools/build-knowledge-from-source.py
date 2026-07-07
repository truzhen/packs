#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
环保执法 Pack 知识库构建器。

从 Owner 显式指定的权威资料知识库把真实的
法规/执法指南/案例/索引导入到本 pack 的 knowledge/ 目录（自包含、可分发），
并生成两个权威清单（loader 据此装入）：
  - knowledge/knowledge-scopes.json   各知识域声明（mount_on_pack_enable）
  - knowledge/knowledge-index.json    每条知识 → scope/kind/source_ref/title/scene/生效日期

用法：
  TRUZHEN_ENV_PACK_SOURCE_DIR=/path/to/source python3 tools/build-knowledge-from-source.py
  python3 tools/build-knowledge-from-source.py [源知识库路径]

知识内容来自 Owner 提供的权威资料。每条以 verification_status=pending_human_review
入库，正式适用前须经法务/业务核验，以现行有效官方法规标准原文为准。
"""
import json
import os
import re
import shutil
import sys

PACK_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DEFAULT_SRC = os.environ.get("TRUZHEN_ENV_PACK_SOURCE_DIR", "")
SCENE_REF = "scene://environmental-enforcement"

# 路径片段（按出现顺序匹配，先具体后笼统）→ (scope_short, display_name, kind)
SCOPE_RULES = [
    ("01_法律法规",                 ("code", "生态环境法典与法律法规域", "law_article")),
    ("02_水污染源执法",             ("water", "水污染源执法知识域", "sop")),
    ("03_大气污染源执法",           ("air", "大气污染源执法知识域", "sop")),
    ("04_固体废物执法",             ("solid-waste", "固体废物执法知识域", "sop")),
    ("05_核安全与放射性执法",       ("radiation", "核安全与放射性执法知识域", "sop")),
    ("06_噪声污染源执法",           ("noise", "噪声污染源执法知识域", "sop")),
    ("01_污染源执法概述",           ("pollution-source-overview", "污染源执法概述域", "sop")),
    ("01_执法概述与改革",           ("guide-overview", "执法概述与改革域", "sop")),
    ("02_执法依据与法规体系",       ("legal-basis", "执法依据与法规体系域", "law_article")),
    ("04_建设项目与排污许可",       ("eia-permit", "建设项目与排污许可域", "sop")),
    ("05_自然生态保护",             ("ecology", "自然生态保护执法域", "sop")),
    ("06_行政处罚",                 ("penalty", "行政处罚程序与裁量域", "law_article")),
    ("07_行刑衔接",                 ("criminal", "行刑衔接（两法衔接）域", "law_article")),
    ("08_执法风险防范",             ("risk", "执法风险防范域", "sop")),
    ("03_案例库",                   ("cases", "环保执法案例库域", "case")),
    ("05_索引",                     ("index", "法条速查与违法行为分类索引域", "index")),
]

def classify(relpath):
    for needle, meta in SCOPE_RULES:
        if needle in relpath:
            return meta
    # 03_污染源环境执法 概述兜底
    if "03_污染源环境执法" in relpath:
        return ("pollution-source-overview", "污染源执法概述域", "sop")
    return ("guide-overview", "执法概述与改革域", "sop")

def parse_frontmatter(text):
    title, eff = "", ""
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, re.S)
    if m:
        fm = m.group(1)
        t = re.search(r"^title:\s*(.+)$", fm, re.M)
        if t:
            title = t.group(1).strip().lstrip("#").strip()
        d = re.search(r"^effective_date:\s*(.+)$", fm, re.M)
        if d:
            eff = d.group(1).strip()
    return title, eff

def slug(name):
    s = re.sub(r"\.md$", "", name)
    s = re.sub(r'[\s"\'()（）:：,，、]+', "-", s)
    s = re.sub(r"-+", "-", s).strip("-")
    return s[:80]

def main():
    src = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_SRC
    if not os.path.isdir(src):
        print("源知识库不存在: %s" % src, file=sys.stderr)
        sys.exit(2)
    kdir = os.path.join(PACK_DIR, "knowledge")
    # 清掉旧的导入内容（保留 scopes/index 由本脚本重写）
    for sub in os.listdir(kdir) if os.path.isdir(kdir) else []:
        p = os.path.join(kdir, sub)
        if os.path.isdir(p):
            shutil.rmtree(p)
    os.makedirs(kdir, exist_ok=True)

    index = []
    scopes = {}
    md_files = []
    for root, _, files in os.walk(src):
        for fn in files:
            if fn.endswith(".md"):
                md_files.append(os.path.join(root, fn))
    md_files.sort()

    for full in md_files:
        rel = os.path.relpath(full, src)
        scope_short, display_name, kind = classify(rel)
        with open(full, encoding="utf-8") as f:
            text = f.read()
        title, eff = parse_frontmatter(text)
        if not title:
            title = re.sub(r"\.md$", "", os.path.basename(full))
        scope_ref = "knowledge_scope://environmental/%s" % scope_short
        dest_dir = os.path.join(kdir, scope_short)
        os.makedirs(dest_dir, exist_ok=True)
        dest = os.path.join(dest_dir, os.path.basename(full))
        shutil.copyfile(full, dest)
        entry = {
            "file": os.path.relpath(dest, PACK_DIR),
            "source_ref": "source://environmental/%s/%s" % (scope_short, slug(os.path.basename(full))),
            "title": title,
            "kind": kind,
            "knowledge_scope_ref": scope_ref,
            "scene_ref": SCENE_REF,
            "authority": "reference_only",
            "verification_status": "pending_human_review",
        }
        if eff:
            entry["effective_from"] = eff
        index.append(entry)
        sc = scopes.setdefault(scope_ref, {
            "scope_ref": scope_ref,
            "display_name": display_name,
            "scene_ref": SCENE_REF,
            "mount_policy": "mount_on_pack_enable",
            "required": True,
            "knowledge_kinds": set(),
            "source_dir": "knowledge/%s" % scope_short,
        })
        sc["knowledge_kinds"].add(kind)

    scope_list = []
    for sc in scopes.values():
        sc["knowledge_kinds"] = sorted(sc["knowledge_kinds"])
        scope_list.append(sc)
    scope_list.sort(key=lambda x: x["scope_ref"])

    with open(os.path.join(kdir, "knowledge-scopes.json"), "w", encoding="utf-8") as f:
        json.dump({"scene_ref": SCENE_REF, "default_mount_policy": "mount_on_pack_enable",
                   "scopes": scope_list}, f, ensure_ascii=False, indent=2)
    with open(os.path.join(kdir, "knowledge-index.json"), "w", encoding="utf-8") as f:
        json.dump({"scene_ref": SCENE_REF, "count": len(index), "entries": index},
                  f, ensure_ascii=False, indent=2)

    print("导入完成：%d 条知识，%d 个知识域" % (len(index), len(scope_list)))
    for sc in scope_list:
        n = sum(1 for e in index if e["knowledge_scope_ref"] == sc["scope_ref"])
        print("  %-46s %2d 条  (%s)" % (sc["scope_ref"], n, ",".join(sc["knowledge_kinds"])))

if __name__ == "__main__":
    main()
