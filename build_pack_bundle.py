#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""pack 交付 bundle 打包器（商用就绪 C1 pack 侧）。

把一个场景包目录打成「自包含可交付」bundle：买家拿到 bundle 解压后，
`python3 <pack>/install.py` 能直接对已运行的生产基座（TRUZHEN_DEVSERVER_BASE）装入，
不需要 packs 源码树、不需要手动 go run devserver。

关键：install.py 通过 `REPO_DIR = dirname(PACK_DIR); from pack_diagnostics import ...`
从 pack 目录的父层导入共享诊断模块。因此 bundle 布局必须是：

    <pack>.bundle.zip
      ├── pack_diagnostics.py        # 置于 pack 父层，满足 install.py 父目录导入
      └── <pack-name>/
          ├── install.py  uninstall.py  manifest.json
          ├── flows/ role-slots/ role-packs/ capabilities/ ...

产物：
  - 交付 bundle：dist/<pack-name>.bundle.zip（自包含父目录布局）；
  - 市场制品：dist/<pack-name>.market.zip（manifest.json 位于 ZIP 根）。
两类制品均生成每文件 sha256 + 制品 sha256 清单，并统一排除运行态禁品。

用法：
  python3 build_pack_bundle.py smart-home-owner-pack-v0            # 交付 bundle
  python3 build_pack_bundle.py --market smart-home-owner-pack-v0   # 云市场制品
"""
import hashlib
import json
import os
import re
import sys
import tempfile
import zipfile

REPO_ROOT = os.path.dirname(os.path.abspath(__file__))
DIAGNOSTICS = os.path.join(REPO_ROOT, "pack_diagnostics.py")
REQUIRED_FILES = ("manifest.json", "install.py", "uninstall.py")
REQUIRED_MANIFEST_FIELDS = (
    "pack_id", "name", "version", "kind", "min_truzhen_version", "lifecycle_status",
)
PACK_KINDS = {"scene_pack", "capability_pack", "role_pack", "skill_bundle"}
LIFECYCLE_STATUSES = {
    "想法", "设计中", "契约已定", "已实现", "已接线", "已验收", "已发布", "已弃用",
}
SEMVER_PATTERN = re.compile(r"^\d+\.\d+\.\d+$")
VERSION_RANGE_PATTERN = re.compile(r"^(>=|<=|>|<|=)?\d+\.\d+\.\d+(,(>=|<=|>|<|=)?\d+\.\d+\.\d+)*$")
ISOLATION_POLICIES = {"reuse_preferred", "isolated_install", "coexist_multi_version", "blocked"}
FALLBACK_POLICIES = {"blocked", "provider_missing", "manual_handoff", "not_ready"}
GATEWAY_CLASSES = {"execution", "communication", "model", "memory"}
RISK_CLASSES = {"low", "medium", "high", "critical"}
_EXCLUDE_DIRS = {"__pycache__", ".git", "node_modules", "dist", "build", ".vite"}
_EXCLUDE_SUFFIXES = (".pyc", ".db", ".sqlite", ".sqlite3", ".log", ".jsonl")


def _sha256(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def _validate(pack_dir):
    if not os.path.isdir(pack_dir):
        raise ValueError("pack 目录不存在: %s" % pack_dir)
    for req in REQUIRED_FILES:
        if not os.path.isfile(os.path.join(pack_dir, req)):
            raise ValueError("残缺 pack：缺 %s（不产半成品 bundle）" % req)
    if not os.path.isfile(DIAGNOSTICS):
        raise ValueError("缺共享 pack_diagnostics.py（%s）；install.py 父目录导入会失败" % DIAGNOSTICS)
    # manifest 声明的文件必须真实存在（交付前置门）
    with open(os.path.join(pack_dir, "manifest.json"), encoding="utf-8") as handle:
        manifest = json.load(handle)
    missing = [field for field in REQUIRED_MANIFEST_FIELDS if not manifest.get(field)]
    if missing:
        raise ValueError("manifest 不符合 canonical PackManifest：缺 %s" % ", ".join(missing))
    if not SEMVER_PATTERN.fullmatch(str(manifest["version"])):
        raise ValueError("manifest version 必须是 Cloud 接受的 SemVer（x.y.z）")
    if manifest["kind"] not in PACK_KINDS:
        raise ValueError("manifest kind 非 canonical 值：%s" % manifest["kind"])
    if manifest["lifecycle_status"] not in LIFECYCLE_STATUSES:
        raise ValueError("manifest lifecycle_status 必须是八档单值：%s" % manifest["lifecycle_status"])
    if manifest["kind"] == "scene_pack":
        for script in ("install.py", "uninstall.py"):
            if not os.path.isfile(os.path.join(pack_dir, script)):
                raise ValueError("场景荚市场制品缺 %s" % script)
    _validate_software_requirements(manifest.get("software_requirements", []))
    _validate_provider_requirements(manifest.get("provider_requirements", []))
    for key in ("flow_file", "role_slots_file", "capabilities_file", "knowledge_index",
                "knowledge_scopes_manifest"):
        rel = manifest.get(key)
        if rel and not os.path.isfile(os.path.join(pack_dir, rel)):
            raise ValueError("manifest 声明的 %s=%s 在 pack 内不存在" % (key, rel))
    return manifest


def _validate_software_requirements(items):
    seen = set()
    required = (
        "requirement_id", "software_family", "version_range", "isolation_policy",
        "fallback_policy", "gateway_class", "risk_class",
    )
    for item in items:
        missing = [field for field in required if not item.get(field)]
        if missing:
            raise ValueError("software_requirement 缺 %s" % ", ".join(missing))
        rid = item["requirement_id"]
        if rid in seen:
            raise ValueError("software_requirement requirement_id 重复：%s" % rid)
        seen.add(rid)
        if not VERSION_RANGE_PATTERN.fullmatch(str(item["version_range"])):
            raise ValueError("software_requirement version_range 非 canonical：%s" % item["version_range"])
        if item["isolation_policy"] not in ISOLATION_POLICIES:
            raise ValueError("software_requirement isolation_policy 非 canonical：%s" % item["isolation_policy"])
        if item["fallback_policy"] not in FALLBACK_POLICIES:
            raise ValueError("software_requirement fallback_policy 非 canonical：%s" % item["fallback_policy"])
        if item["gateway_class"] not in GATEWAY_CLASSES or item["risk_class"] not in RISK_CLASSES:
            raise ValueError("software_requirement gateway_class/risk_class 非 canonical")


def _validate_provider_requirements(items):
    seen = set()
    required = ("requirement_id", "provider_family", "fallback_policy", "gateway_class", "risk_class")
    for item in items:
        missing = [field for field in required if not item.get(field)]
        if missing:
            raise ValueError("provider_requirement 缺 %s" % ", ".join(missing))
        rid = item["requirement_id"]
        if rid in seen:
            raise ValueError("provider_requirement requirement_id 重复：%s" % rid)
        seen.add(rid)
        if item["fallback_policy"] not in FALLBACK_POLICIES:
            raise ValueError("provider_requirement fallback_policy 非 canonical：%s" % item["fallback_policy"])
        if item["gateway_class"] not in GATEWAY_CLASSES or item["risk_class"] not in RISK_CLASSES:
            raise ValueError("provider_requirement gateway_class/risk_class 非 canonical")


def _iter_pack_files(pack_dir):
    for root, dirs, files in os.walk(pack_dir):
        dirs[:] = [d for d in dirs if d not in _EXCLUDE_DIRS]
        for fn in files:
            if fn.lower().endswith(_EXCLUDE_SUFFIXES):
                continue
            yield os.path.join(root, fn)


def _write_artifact(path, members, pack_name, artifact_kind, install_hint):
    with zipfile.ZipFile(path, "w", zipfile.ZIP_DEFLATED) as z:
        for arc, disk in members:
            z.write(disk, arc)

    files_manifest = [{"path": arc, "sha256": _sha256(disk)} for arc, disk in members]
    manifest = {
        "artifact": os.path.basename(path),
        "artifact_kind": artifact_kind,
        "pack_name": pack_name,
        "artifact_sha256": _sha256(path),
        "file_count": len(files_manifest),
        "files": files_manifest,
        "install_hint": install_hint,
    }
    # 保留既有 bundle sidecar 字段；新增通用 artifact 字段仅为加法兼容。
    if artifact_kind == "self_contained_delivery_bundle":
        manifest["bundle"] = os.path.basename(path)
        manifest["bundle_sha256"] = manifest["artifact_sha256"]
    man_path = path.rsplit(".zip", 1)[0] + ".manifest.json"
    with open(man_path, "w", encoding="utf-8") as handle:
        json.dump(manifest, handle, ensure_ascii=False, indent=2)
        handle.write("\n")
    return path


def build_pack_bundle(pack_dir, out_dir=None):
    """打一个 pack 目录成自包含 bundle.zip，返回 bundle 路径。"""
    pack_dir = os.path.abspath(pack_dir.rstrip("/"))
    _validate(pack_dir)
    name = os.path.basename(pack_dir)
    out_dir = os.path.abspath(out_dir or os.path.join(REPO_ROOT, "dist"))
    os.makedirs(out_dir, exist_ok=True)
    bundle_path = os.path.join(out_dir, name + ".bundle.zip")

    # 收集 (arcname, 磁盘路径)：pack_diagnostics.py 置根，pack 内容置 <name>/ 下
    members = [("pack_diagnostics.py", DIAGNOSTICS)]
    for disk in sorted(_iter_pack_files(pack_dir)):
        arc = os.path.join(name, os.path.relpath(disk, pack_dir))
        members.append((arc, disk))

    return _write_artifact(
        bundle_path,
        members,
        name,
        "self_contained_delivery_bundle",
        "解压后：TRUZHEN_DEVSERVER_BASE=http://<基座地址> python3 %s/install.py" % name,
    )


def build_market_artifact(pack_dir, out_dir=None):
    """生成 manifest.json 位于 ZIP 根的云市场制品。"""
    pack_dir = os.path.abspath(pack_dir.rstrip("/"))
    manifest = _validate(pack_dir)
    name = os.path.basename(pack_dir)
    out_dir = os.path.abspath(out_dir or os.path.join(REPO_ROOT, "dist"))
    os.makedirs(out_dir, exist_ok=True)
    artifact_path = os.path.join(out_dir, name + ".market.zip")
    members = [(os.path.relpath(disk, pack_dir), disk) for disk in sorted(_iter_pack_files(pack_dir))]
    if not any(arc == "manifest.json" for arc, _ in members):
        raise ValueError("市场制品根目录缺少 manifest.json")
    return _write_artifact(
        artifact_path,
        members,
        name,
        "cloud_market_pack",
        "经 truzhen-cloud 市场下载并由 truzhenos lifecycle 受控安装：%s@%s"
        % (manifest.get("pack_id") or manifest.get("pack_ref"), manifest.get("version")),
    )


if __name__ == "__main__":
    args = sys.argv[1:]
    market = bool(args and args[0] == "--market")
    if market:
        args = args[1:]
    if not args:
        print("用法: python3 build_pack_bundle.py [--market] <pack-dir> [out-dir]", file=sys.stderr)
        sys.exit(2)
    out = args[1] if len(args) > 1 else None
    path = build_market_artifact(args[0], out) if market else build_pack_bundle(args[0], out)
    print("制品已产出:", path)
    print("manifest:", path.rsplit(".zip", 1)[0] + ".manifest.json")
