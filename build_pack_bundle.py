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
PROVIDER_REQUIREMENT_FIELDS = {
    "requirement_id", "provider_family", "gateway_class", "required_capabilities",
    "software_requirement_refs", "risk_class", "fallback_policy", "optional",
}
CAPABILITY_FIELDS = {"capability_id", "provider_requirement_ref", "description", "optional"}
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
    for key in ("flow_file", "role_slots_file", "capabilities_file", "knowledge_index",
                "knowledge_scopes_manifest"):
        rel = manifest.get(key)
        if rel and not os.path.isfile(os.path.join(pack_dir, rel)):
            raise ValueError("manifest 声明的 %s=%s 在 pack 内不存在" % (key, rel))
    software_requirements = manifest.get("software_requirements", [])
    provider_requirements = manifest.get("provider_requirements", [])
    _validate_software_requirements(software_requirements)
    _validate_provider_requirements(provider_requirements)
    capabilities_doc = _load_declared_json(pack_dir, manifest.get("capabilities_file"), "capabilities")
    flow_doc = _load_declared_json(pack_dir, manifest.get("flow_file"), "flow")
    _validate_provider_requirement_bindings(
        provider_requirements, software_requirements, capabilities_doc, flow_doc)
    return manifest


def _load_declared_json(pack_dir, relative_path, label):
    if not relative_path:
        if label == "capabilities":
            return {"capabilities": []}
        return {"nodes": []}
    with open(os.path.join(pack_dir, relative_path), encoding="utf-8") as handle:
        return json.load(handle)


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
        unknown = sorted(set(item) - PROVIDER_REQUIREMENT_FIELDS)
        if unknown:
            raise ValueError(
                "migration_warning: provider_requirement 含非法字段 %s；canonical pack 暂不可执行"
                % ", ".join(unknown))
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
        capabilities = item.get("required_capabilities")
        if not isinstance(capabilities, list) or not capabilities or any(
                not isinstance(value, str) or not value for value in capabilities):
            raise ValueError("provider_requirement required_capabilities 必须是非空字符串列表：%s" % rid)
        if len(capabilities) != len(set(capabilities)):
            raise ValueError("provider_requirement required_capabilities 重复：%s" % rid)
        software_refs = item.get("software_requirement_refs")
        if not isinstance(software_refs, list) or not software_refs or any(
                not isinstance(value, str) or not value for value in software_refs):
            raise ValueError("provider_requirement software_requirement_refs 必须是非空字符串列表：%s" % rid)
        if len(software_refs) != len(set(software_refs)):
            raise ValueError("provider_requirement software_requirement_refs 重复：%s" % rid)


def _validate_provider_requirement_bindings(provider_items, software_items, capabilities_doc, flow_doc):
    """校验 Pack 内 ProviderRequirement、软件需求、能力定义和流程引用的闭合性。"""
    providers = {item["requirement_id"]: item for item in provider_items}
    software = {item["requirement_id"]: item for item in software_items}
    referenced_software = set()
    for provider in provider_items:
        provider_id = provider["requirement_id"]
        for software_id in provider["software_requirement_refs"]:
            if "://" in software_id or software_id not in software:
                raise ValueError(
                    "provider_requirement %s 引用了未知或跨 Pack software_requirement：%s"
                    % (provider_id, software_id))
            referenced_software.add(software_id)
            software_family = software[software_id].get("provider_family")
            if software_family and software_family != provider["provider_family"]:
                raise ValueError(
                    "provider_requirement %s 与 software_requirement %s provider_family 不一致：%s != %s"
                    % (provider_id, software_id, provider["provider_family"], software_family))
            software_capabilities = set(software[software_id].get("required_capabilities", []))
            missing_capabilities = sorted(set(provider["required_capabilities"]) - software_capabilities)
            if missing_capabilities:
                raise ValueError(
                    "provider_requirement %s 的 capability 未被 software_requirement %s 声明：%s"
                    % (provider_id, software_id, ", ".join(missing_capabilities)))
    orphaned = sorted(set(software) - referenced_software)
    if orphaned:
        raise ValueError("software_requirement 孤儿声明（未被任何 ProviderRequirement 引用）：%s" % ", ".join(orphaned))

    if not isinstance(capabilities_doc, dict):
        raise ValueError("capabilities 文件必须是对象；canonical pack 暂不可执行")
    if "provider_requirements" in capabilities_doc:
        raise ValueError(
            "migration_warning: capabilities.json 不得重复 provider_requirements；canonical pack 暂不可执行")
    unknown_top = sorted(set(capabilities_doc) - {"capabilities"})
    if unknown_top:
        raise ValueError("capabilities.json 含非法顶层字段：%s" % ", ".join(unknown_top))
    capabilities = capabilities_doc.get("capabilities")
    if not isinstance(capabilities, list):
        raise ValueError("capabilities.json capabilities 必须是列表")
    capability_ids = set()
    capability_provider_refs = {}
    for capability in capabilities:
        if not isinstance(capability, dict):
            raise ValueError("capabilities.json 每项必须是对象")
        unknown = sorted(set(capability) - CAPABILITY_FIELDS)
        if unknown:
            raise ValueError("capabilities.json 含非法字段：%s" % ", ".join(unknown))
        capability_id = capability.get("capability_id")
        provider_ref = capability.get("provider_requirement_ref")
        if not capability_id or not provider_ref:
            raise ValueError("capabilities.json 缺 capability_id/provider_requirement_ref")
        if capability_id in capability_ids:
            raise ValueError("capability_id 重复：%s" % capability_id)
        capability_ids.add(capability_id)
        capability_provider_refs[capability_id] = provider_ref
        if "://" in provider_ref or provider_ref not in providers:
            raise ValueError("capability %s 引用了未知或跨 Pack ProviderRequirement：%s" %
                             (capability_id, provider_ref))

    expected_capabilities = set()
    for provider in provider_items:
        provider_id = provider["requirement_id"]
        expected_capabilities.update(provider["required_capabilities"])
        for capability_id in provider["required_capabilities"]:
            if capability_id not in capability_ids:
                raise ValueError("ProviderRequirement %s 缺 capability 定义：%s" % (provider_id, capability_id))
            if capability_provider_refs[capability_id] != provider_id:
                raise ValueError("capability %s 的 provider_requirement_ref 与 ProviderRequirement 不一致" %
                                 capability_id)
    orphaned_capabilities = sorted(capability_ids - expected_capabilities)
    if orphaned_capabilities:
        raise ValueError("capability 孤儿声明（未被 ProviderRequirement 使用）：%s" %
                         ", ".join(orphaned_capabilities))

    if not isinstance(flow_doc, dict):
        raise ValueError("flow 文件必须是对象")
    nodes = flow_doc.get("nodes", [])
    if not isinstance(nodes, list):
        raise ValueError("flow nodes 必须是列表")
    flow_provider_refs = set()
    node_ids = set()
    for node in nodes:
        if not isinstance(node, dict) or not node.get("id"):
            raise ValueError("flow node 必须含非空 id")
        if node["id"] in node_ids:
            raise ValueError("flow node id 重复：%s" % node["id"])
        node_ids.add(node["id"])
        provider_ref = node.get("provider_requirement_ref")
        if provider_ref is None:
            continue
        if not isinstance(provider_ref, str) or not provider_ref:
            raise ValueError("flow node provider_requirement_ref 必须是非空字符串")
        if "://" in provider_ref or provider_ref not in providers:
            raise ValueError("flow node 引用了未知或跨 Pack ProviderRequirement：%s" % provider_ref)
        flow_provider_refs.add(provider_ref)
    missing_flow_refs = sorted(set(providers) - flow_provider_refs)
    if missing_flow_refs:
        raise ValueError("ProviderRequirement 未接入 flow 节点：%s" % ", ".join(missing_flow_refs))


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
