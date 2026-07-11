#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""TDD 守卫：pack 交付 bundle 打包器（商用就绪 C1 pack 侧）。

验证 build_pack_bundle 产出的 bundle 是「自包含可交付」的：
  1. 含 pack 目录（install.py/uninstall.py/manifest.json + 声明文件）；
  2. 含 pack_diagnostics.py 且置于 pack 目录的父层——满足 install.py
     `REPO_DIR = dirname(PACK_DIR); from pack_diagnostics import ...` 的父目录导入；
  3. 缺 install.py 的残缺 pack 被拒（不产半成品 bundle）；
  4. bundle manifest 含每文件 sha256（交付核验）。

跑法：python3 test_pack_bundle.py
"""
import json
import os
import shutil
import sys
import tempfile
import zipfile

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)


def _make_fake_pack(root, name, with_install=True):
    pd = os.path.join(root, name)
    for sub in ("flows", "role-slots", "capabilities"):
        os.makedirs(os.path.join(pd, sub))
    json.dump(
        {"manifest_version": "v3", "pack_ref": "scene_pack://fake", "version": "1.0.0",
         "flow_file": "flows/f.flow.json", "role_slots_file": "role-slots/role-slots.json",
         "capabilities_file": "capabilities/capabilities.json"},
        open(os.path.join(pd, "manifest.json"), "w", encoding="utf-8"))
    open(os.path.join(pd, "uninstall.py"), "w").write("# uninstall\n")
    open(os.path.join(pd, "flows", "f.flow.json"), "w").write("{}")
    open(os.path.join(pd, "role-slots", "role-slots.json"), "w").write('{"role_slots":[]}')
    open(os.path.join(pd, "capabilities", "capabilities.json"), "w").write('{"provider_requirements":[]}')
    if with_install:
        open(os.path.join(pd, "install.py"), "w").write("from pack_diagnostics import emit_pack_error\n")
    return pd


def test_bundle_self_contained():
    from build_pack_bundle import build_pack_bundle
    tmp = tempfile.mkdtemp()
    try:
        pd = _make_fake_pack(tmp, "fake-pack-v0")
        out = os.path.join(tmp, "dist")
        bundle = build_pack_bundle(pd, out)
        assert os.path.exists(bundle), "bundle 未产出"
        with zipfile.ZipFile(bundle) as z:
            names = z.namelist()
        # pack_diagnostics.py 必须在 pack 目录父层（bundle 根），使 install.py 父目录导入可解析
        assert "pack_diagnostics.py" in names, "bundle 根缺 pack_diagnostics.py（install.py 会 import 失败）"
        assert "fake-pack-v0/install.py" in names, "bundle 缺 install.py"
        assert "fake-pack-v0/uninstall.py" in names, "bundle 缺 uninstall.py"
        assert "fake-pack-v0/manifest.json" in names, "bundle 缺 manifest.json"
        assert "fake-pack-v0/flows/f.flow.json" in names, "bundle 缺声明的 flow_file"
        # manifest 含 sha256
        man = os.path.join(out, "fake-pack-v0.bundle.manifest.json")
        assert os.path.exists(man), "缺 bundle manifest"
        md = json.load(open(man, encoding="utf-8"))
        assert md.get("bundle_sha256"), "manifest 缺 bundle_sha256"
        assert any(f.get("path") == "pack_diagnostics.py" and f.get("sha256") for f in md.get("files", [])), \
            "manifest 未登记 pack_diagnostics.py 的 sha256"
        print("PASS test_bundle_self_contained")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


def test_rejects_pack_missing_install():
    from build_pack_bundle import build_pack_bundle
    tmp = tempfile.mkdtemp()
    try:
        pd = _make_fake_pack(tmp, "broken-pack-v0", with_install=False)
        out = os.path.join(tmp, "dist")
        raised = False
        try:
            build_pack_bundle(pd, out)
        except ValueError as e:
            raised = True
            assert "install.py" in str(e), "拒绝原因应点名缺 install.py"
        assert raised, "缺 install.py 的残缺 pack 必须被拒，不得产半成品 bundle"
        print("PASS test_rejects_pack_missing_install")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


if __name__ == "__main__":
    test_bundle_self_contained()
    test_rejects_pack_missing_install()
    print("ALL PASS")
