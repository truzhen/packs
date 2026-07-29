# Truzhen v4 第六批 Packs 集成收口

日期：2026-07-29

## 范围与结论

- 基线：`truzhen-packs` `7289f83684d697756be7f51034138fd0afa1d568`。
- G21 / G23：接入正式卸载语义、共享静态防线、EPC 拓扑测试及各自 closeout，结论为候选包已验收、未发布。
- G22：接入完整审计演进与 Owner RC 裁定。状态保持 `blocked`；Owner 接受 `sandbox_not_ready` 作为本次 RC 的受控边界，不代表 Hands、模型或 `pack.candidate.generate` 已验收。
- 真相源边界：Pack 仓只保存声明、脚本、测试与去敏 closeout；正式 Gate、Receipt、Provider、runtime 和发布事实不由本仓持有。

## 验证

| 门 | 结果 | 日志 / SHA256 |
| --- | --- | --- |
| 递归 JSON | PASS | `/tmp/truzhen-batch6-packs-json-20260729.log` / `48bfe0757f501c22d336cebf221fbf1a6a6dc5eecaf44d0499c9c7787bc176ca` |
| install / uninstall Python compile | PASS | `/tmp/truzhen-batch6-packs-pycompile-20260729.log` / `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Python 全量与 Pack 专项 | PASS；根级 28/28、家政 1/1、墅学家 6/6 + 7/7、智能家居 4/4 | `/tmp/truzhen-batch6-packs-python-tests-20260729.log` / `7ba6e6868d46d5ff4a21e0093e2ec725460824453bea07e235c677f8fcfcfbb8` |
| `GOWORK=off go test ./...` | PASS | `/tmp/truzhen-batch6-packs-go-test-20260729.log` / `c47793550d0a003ea2316516f5dbf8631b9c6f3c3a55ab14adb9ce294eb43451` |
| Pack 结构与 tracked forbidden artifacts | PASS | `/tmp/truzhen-batch6-packs-structure-forbidden-20260729.log` / `b5b119ebf948e2ed9de2d63b2319fa6ecd25a84e1f24f5a843748ed3912b5a4f` |
| 34 个变更文件敏感边界扫描 | PASS；无私钥、长 token、Kratos 会话或正式 Receipt 原文 | 本次集成终检 |
| `git diff --check` | PASS | 本次集成终检 |

测试输出中的“半装状态已落盘”来自断点恢复负向用例，测试进程整体成功，不是发布链失败。

## 禁止边界

- 未执行 Pack 发布、市场上架、平台登录、上传、发送或发布。
- 未调用真实 Provider / 模型，未启动容器，未写生产运行态。
- 未把 G22 的受控阻断包装成 pass。
