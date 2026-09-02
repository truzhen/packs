#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""市场制品 canonical manifest 投影守卫。

作者侧 manifest 保留工作台、Gate、binding、知识等富声明；云市场 ZIP 根只允许
Contracts v0.19 PackManifest 字段。设置 TRUZHEN_CONTRACTS_DIR 时，本测试会直接
读取并执行该仓的 pack/provider schema，不复制 schema 真相源。
"""

import hashlib
import json
import os
import pathlib
import shutil
import tempfile
import zipfile

from build_pack_bundle import build_market_artifact


HERE = pathlib.Path(__file__).resolve().parent
FORMAL_PACK_DIRS = (
    "backup-administrator-workbench-v0",
    "content-operations-workbench-v0",
    "environmental-enforcement-pack-v0",
    "housekeeping-ops-pack-v0",
    "project-watch-pack-v0",
    "shuxuejia-renovation-pack-v0",
    "smart-home-owner-pack-v0",
)
PACK_SCHEMA_SHA256 = "20a983b5f657eecbe0d725304f41e1f5a4610bb148b51da1c88a7ee77b0f4534"
PROVIDER_SCHEMA_SHA256 = "48c6237608a5d374368142fcd00a6bc4721c90c349a9a2ccc8a9db9dcb03bbc1"


def _sha256_bytes(payload):
    return hashlib.sha256(payload).hexdigest()


def _contracts_dir():
    configured = os.environ.get("TRUZHEN_CONTRACTS_DIR", "").strip()
    if configured:
        return pathlib.Path(configured).resolve()
    fixed = pathlib.Path("/Users/li/Documents/truzhen-contracts")
    if fixed.is_dir():
        return fixed
    raise AssertionError(
        "缺 Contracts v0.19 schema；请设置 TRUZHEN_CONTRACTS_DIR 指向只读 truzhen-contracts"
    )


def _load_contract_schemas():
    contracts_dir = _contracts_dir()
    pack_path = contracts_dir / "pack-manifest.schema.json"
    provider_path = contracts_dir / "provider-requirement.schema.json"
    pack_bytes = pack_path.read_bytes()
    provider_bytes = provider_path.read_bytes()
    assert _sha256_bytes(pack_bytes) == PACK_SCHEMA_SHA256, "pack schema 不是冻结的 Contracts v0.19"
    assert _sha256_bytes(provider_bytes) == PROVIDER_SCHEMA_SHA256, \
        "provider schema 不是冻结的 Contracts v0.19"
    return pack_path, json.loads(pack_bytes), {
        provider_path.resolve(): json.loads(provider_bytes),
    }


def _json_pointer(document, pointer):
    value = document
    for token in pointer.lstrip("/").split("/"):
        token = token.replace("~1", "/").replace("~0", "~")
        value = value[token]
    return value


def _validate_schema(instance, schema, schema_path, root_schema, external_schemas, at="$"):
    """执行本任务使用到的 JSON Schema 2020-12 关键字；未知关键字只作注释。"""
    ref = schema.get("$ref")
    if ref:
        if ref.startswith("#"):
            target = _json_pointer(root_schema, ref[1:])
            return _validate_schema(
                instance, target, schema_path, root_schema, external_schemas, at
            )
        ref_path = (schema_path.parent / ref).resolve()
        target = external_schemas.get(ref_path)
        if target is None:
            target = json.loads(ref_path.read_text(encoding="utf-8"))
            external_schemas[ref_path] = target
        return _validate_schema(instance, target, ref_path, target, external_schemas, at)

    expected_type = schema.get("type")
    if expected_type == "object":
        assert isinstance(instance, dict), f"{at}: expected object"
    elif expected_type == "array":
        assert isinstance(instance, list), f"{at}: expected array"
    elif expected_type == "string":
        assert isinstance(instance, str), f"{at}: expected string"
    elif expected_type == "boolean":
        assert isinstance(instance, bool), f"{at}: expected boolean"

    if "enum" in schema:
        assert instance in schema["enum"], f"{at}: {instance!r} not in enum"
    if isinstance(instance, str) and "minLength" in schema:
        assert len(instance) >= schema["minLength"], f"{at}: shorter than minLength"
    if isinstance(instance, list):
        if "minItems" in schema:
            assert len(instance) >= schema["minItems"], f"{at}: fewer than minItems"
        if schema.get("uniqueItems"):
            encoded = [json.dumps(item, sort_keys=True, ensure_ascii=False) for item in instance]
            assert len(encoded) == len(set(encoded)), f"{at}: duplicate array items"
        item_schema = schema.get("items")
        if item_schema:
            for index, item in enumerate(instance):
                _validate_schema(
                    item, item_schema, schema_path, root_schema, external_schemas,
                    f"{at}[{index}]",
                )
    if isinstance(instance, dict):
        required = schema.get("required", ())
        missing = [key for key in required if key not in instance]
        assert not missing, f"{at}: missing required {missing}"
        properties = schema.get("properties", {})
        if schema.get("additionalProperties") is False:
            extras = sorted(set(instance) - set(properties))
            assert not extras, f"{at}: additional properties {extras}"
        for key, value in instance.items():
            if key in properties:
                _validate_schema(
                    value, properties[key], schema_path, root_schema, external_schemas,
                    f"{at}.{key}",
                )


def _read_market_manifest(artifact):
    with zipfile.ZipFile(artifact) as archive:
        return json.loads(archive.read("manifest.json"))


def test_all_real_market_roots_match_contracts_v019():
    schema_path, schema, external = _load_contract_schemas()
    failures = []
    temp_root = pathlib.Path(tempfile.mkdtemp(prefix="packs-market-projection-"))
    try:
        for pack_name in FORMAL_PACK_DIRS:
            pack_dir = HERE / pack_name
            try:
                artifact = build_market_artifact(str(pack_dir), str(temp_root / pack_name))
                projected = _read_market_manifest(artifact)
                _validate_schema(projected, schema, schema_path, schema, external)
            except Exception as error:  # 汇总 6/6 真实 Pack，便于一次看清旧根漂移。
                failures.append(f"{pack_name}: {error}")
        assert not failures, "Contracts v0.19 market root failures:\n" + "\n".join(failures)
        print("PASS test_all_real_market_roots_match_contracts_v019")
    finally:
        shutil.rmtree(temp_root, ignore_errors=True)


def test_projection_preserves_author_source_and_strips_internal_fields():
    temp_root = pathlib.Path(tempfile.mkdtemp(prefix="packs-market-author-source-"))
    try:
        for pack_name in FORMAL_PACK_DIRS:
            pack_dir = HERE / pack_name
            source_path = pack_dir / "manifest.json"
            source_before = source_path.read_bytes()
            source = json.loads(source_before)
            artifact = build_market_artifact(str(pack_dir), str(temp_root / pack_name))
            projected = _read_market_manifest(artifact)

            assert source_path.read_bytes() == source_before, f"{pack_name}: 打包不得改写作者 manifest"
            assert "author" in source, f"{pack_name}: 作者富 manifest 缺 author"
            assert "author" not in projected, f"{pack_name}: author 不得偷渡 canonical 根"
            assert "gates" not in projected, f"{pack_name}: gates 不得偷渡 canonical 根"

            source_providers = source.get("provider_requirements", [])
            projected_providers = projected.get("provider_requirements", [])
            assert len(projected_providers) == len(source_providers)
            for rich, canonical in zip(source_providers, projected_providers):
                expected = rich.get("required_capabilities") or [rich.get("capability")]
                assert canonical.get("required_capabilities") == expected
                assert "binding" not in canonical
                assert "capability" not in canonical
                assert "description" not in canonical

            if pack_name == "content-operations-workbench-v0":
                rich_openmontage = source["software_requirements"][1]
                canonical_openmontage = projected["software_requirements"][1]
                assert rich_openmontage["license_policy"] == "review_required"
                assert rich_openmontage["license_evidence"]["declared_license"] == "AGPL-3.0"
                assert canonical_openmontage["license_policy"] == "review_required"
                assert "license_evidence" not in canonical_openmontage
        print("PASS test_projection_preserves_author_source_and_strips_internal_fields")
    finally:
        shutil.rmtree(temp_root, ignore_errors=True)


def test_market_sidecar_hashes_match_canonical_zip_bytes():
    temp_root = pathlib.Path(tempfile.mkdtemp(prefix="packs-market-checksum-"))
    try:
        for pack_name in FORMAL_PACK_DIRS:
            artifact = pathlib.Path(build_market_artifact(
                str(HERE / pack_name), str(temp_root / pack_name)
            ))
            sidecar = json.loads(
                artifact.with_name(artifact.name.removesuffix(".zip") + ".manifest.json")
                .read_text(encoding="utf-8")
            )
            declared = {item["path"]: item["sha256"] for item in sidecar["files"]}
            assert list(declared).count("manifest.json") == 1
            with zipfile.ZipFile(artifact) as archive:
                names = archive.namelist()
                assert names.count("manifest.json") == 1
                for name in names:
                    assert _sha256_bytes(archive.read(name)) == declared[name], \
                        f"{pack_name}: {name} sidecar checksum 漂移"
            assert _sha256_bytes(artifact.read_bytes()) == sidecar["artifact_sha256"]
        print("PASS test_market_sidecar_hashes_match_canonical_zip_bytes")
    finally:
        shutil.rmtree(temp_root, ignore_errors=True)


if __name__ == "__main__":
    test_all_real_market_roots_match_contracts_v019()
    test_projection_preserves_author_source_and_strips_internal_fields()
    test_market_sidecar_hashes_match_canonical_zip_bytes()
    print("ALL PASS")
