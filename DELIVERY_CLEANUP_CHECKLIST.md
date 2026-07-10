# truzhen-packs 交付清理清单

更新时间：2026-07-07

本文件登记六仓分支吸收计划中的低风险 E 项。它不是功能完成声明，也不替代 Pack 安装、启用、停用、回执或商品化验收。

## 2026-07-07 分支吸收 E 项

| 项 | 状态 | 处理边界 | 后续裁定 |
| --- | --- | --- | --- |
| `pack-market-e2e` dirty worktree 中 `smart-home` `manifest.json` 的 `"kind": "scene_pack"` 1 行差异 | 已登记，未收编 | 本次不碰 dirty worktree、不覆盖责任会话 WIP、不声明 packs 主线有缺口 | 等责任会话 closeout 或 Owner 单独授权后再判断是否补入主线 |

## 2026-07-10 预存测试断链（发现即登记，非本轮改动引入）

- `capability_pack_candidates_test.go` 多个 TestShortVideo* 用例断言 `process_worktree_ref://truzhen-packs/gui-capability-pack-test-plan/...` 存在，但该 superpowers worktree（`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/`）已不在本机——干净 main 同样 FAIL（2026-07-10 复现实证）。`go test ./...` 全量因此常红。
- 处置建议（待 Owner 裁）：a) 把 P12-P18 执行规格/授权卡文档收进本仓 `docs/plans/`（资产入仓，引用改仓内路径）；b) 或将这些断言降级为「引用登记存在性」不 stat 外部路径。未裁前，相关会话验证用범위内套件替代全量。
