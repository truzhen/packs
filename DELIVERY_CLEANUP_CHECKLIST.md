# truzhen-packs 交付清理清单

更新时间：2026-07-23

本文件登记本仓交付级遗留。它不是功能完成声明，也不替代 Pack 安装、启用、停用、回执或商品化验收。

> **账本制度（2026-07-11）**：本文件是 packs 仓**分账**；六仓交付债务**总账** = truzhenos 仓 `DELIVERY_CLEANUP_CHECKLIST.md`。集中处理计划：truzhenos `docs/plans/six-repo-delivery-cleanup-debt-consolidation-plan-20260711.md`。

## 2026-07-23 第五批 Packs 收口

- [x] G18/G19/G20 候选已组成 K04R；独立 Python discovery 的 `sys.argv` 污染与吞异常已根因修复，28/28 discovery、Go、JSON、结构、语法和禁品门全绿。
- [x] K05R 已以独立端口和 test-store 重验三条 lane；原 G19 知识挂载阻断由 OS K01R 显式 Owner 授权消费闭环，原 G20 `1.0.0` 硬编码由 OS K02 改为权威 ReadModel 版本消费。
- [ ] 真实处罚/送达/外发、Home Assistant/设备动作、卸载 Owner-presence、生产备份仓与在线数据恢复均未执行；这些是发布运行门，不得由本仓 Pack 静态或隔离验收替代。

## 2026-07-07 分支吸收 E 项

| 项 | 状态 | 处理边界 | 后续裁定 |
| --- | --- | --- | --- |
| `pack-market-e2e` dirty worktree 中 `smart-home` `manifest.json` 的 `"kind": "scene_pack"` 1 行差异 | 已登记，未收编 | 本次不碰 dirty worktree、不覆盖责任会话 WIP、不声明 packs 主线有缺口 | 等责任会话 closeout 或 Owner 单独授权后再判断是否补入主线 |

## 2026-07-10 预存测试断链（发现即登记，非本轮改动引入）——✅ 2026-07-11 已按方案 b 修复

- 原登记：`capability_pack_candidates_test.go` 多个 TestShortVideo* 用例断言 `process_worktree_ref://truzhen-packs/gui-capability-pack-test-plan/...` 存在，但该 superpowers worktree 已不在本机——干净 main 同样 FAIL，`go test ./...` 全量常红。
- **范围补充（2026-07-11 核验）**：断链比原登记更广——除 `capability_pack_candidates_test.go` 外，`role_pack_candidates_test.go` 另有 3 处引用同一外部 worktree 前缀（经共享 helper 走同一解析路径）。
- **修复（2026-07-11 集中处理轮，Owner「按计划开始施工」= 裁方案 b）**：`process_worktree_ref://` 引用降级为「引用登记存在性」断言（前缀须在 `processWorktreeRefRegistry` 登记，未登记引用必 FAIL——反证探针已验证）；三处内容级读取点（requireExistingPath / 授权卡镜像 / readiness 审计标记）对外部引用统一跳过内容断言并注释防回潮；**仓内路径仍必须真实存在**，其余覆盖不变。修复后 `GOWORK=off go test ./...` 全量绿。
- 残留裁定项：P12-P18 执行规格/授权卡文档是否按方案 a 资产入仓（引用改仓内路径、恢复内容级断言），留 Owner 后续裁定；未裁前引用登记口径为准。
