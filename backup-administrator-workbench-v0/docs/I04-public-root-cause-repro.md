# I04 公共根因最小复现（G18）

日期：2026-07-21。此文件只记录公共测试根因，不修改 Pack 以外的资产。

## 复现

在本仓根执行：

```sh
python3 test_pack_issued_binding.py
GOWORK=off go test ./...
```

## 结果

1. `test_pack_issued_binding.py` 的 `test_uninstall_forwards_exact_confirm_evidence` 对 `environmental-enforcement-pack-v0` 与 `smart-home-owner-pack-v0` 失败：两个外部 Pack 的 `uninstall.py` 会解析测试进程的命令行参数，并在测试替身尚未调用前要求显式 `TRUZHEN_DEVSERVER_BASE`。本 G18 Pack 子用例未失败。
2. `GOWORK=off go test ./...` 的 `TestPackGlueDoesNotMintOwnerActionEvidence` 失败：`environmental-enforcement-pack-v0/uninstall.py` 缺少测试要求的 os-14 canonical `action_type`。

## 边界与建议

这两处均存在于从 `main` 派生的非 G18 Pack；本目标禁止修改它们及公共根工具。I04 应在拥有对应 Pack / 公共测试授权的任务中修复后，重新运行上述两条命令。G18 自身的隔离 lifecycle、restic backup / restore 和 Pack 静态审计不依赖这两处失败。
