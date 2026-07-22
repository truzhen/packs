# K04R 第五批 Packs 候选补修复验

日期：2026-07-22  
分支：`codex/v4-release-K04R-batch5-packs-candidate`  
基线：`5bb1f9bcaf001a50823d681a41276ceb4820e611`

## 结论

K04R `pass`。K04 原候选的 G18、G19、G20 集成内容未重写；本轮只修复仓根独立 Python discovery 暴露的测试隔离缺口，并补齐候选级完整静态验证。

## 根因与修复

`test_pack_issued_binding.py` 直接调用多个 `uninstall.py` 的 `main()`，但没有隔离 `unittest` 自身传入的 `sys.argv`，导致 `argparse` 将 discovery 参数识别为未知参数并产生 2 项 ERROR。测试现使用 `unittest.mock.patch.object(sys, "argv", [mod.__file__])` 隔离并自动恢复参数。

原测试还会调用生产 `call("POST", "/dummy")` 探测只读模式，并以 `except Exception: pass` 吞掉未知异常。现改为只信任 Pack 明确声明的 `OWNER_DISABLE_HANDOFF`；未声明者按 legacy 写链严格验收，未知异常不再被吞掉。

## 验证证据

- `python3 -m unittest discover -s . -p 'test*.py' -v`：28/28 OK，退出码 0。
- 全仓 JSON 递归解析：`JSON 合法`，退出码 0。
- 全部 `install.py` / `uninstall.py`：`python3 -m py_compile` 退出码 0。
- 正式 Pack 一级目录结构审计：`Pack 结构审计 OK`，退出码 0。
- `GOWORK=off go test ./...`：PASS，退出码 0。
- tracked forbidden artifacts 扫描：无命中，退出码 0。
- `git diff --check`：PASS，退出码 0。

K04R 不执行 Pack 安装 E2E；G18/G19/G20 三条隔离 lane 属于 K05R，必须使用 K03R OS 候选、独立端口与独立 test-store 重新验收。
