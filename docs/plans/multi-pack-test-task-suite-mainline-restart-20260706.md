# 多 Pack 主线重启执行计划（2026-07-06，主仓更新后）

## 0. 目标与边界

- 在主仓同步更新后，按主线顺序执行三包场景回归：`housekeeping-ops-pack-v0`、`smart-home-owner-pack-v0`、`environmental-enforcement-pack-v0`。
- 本轮**只回归不改码**，不触碰 contracts/schema/DTO/Gate 接口；不做真实发送、真实执行、真实派工、真实 Frappe 写回。
- 外部动作按真实可用能力链路输出：候选、`provider_missing`、`not_ready`、`blocked`、`manual_handoff`。

## 1. 参与仓库与工作树

- `truzhen-packs`（工作树：`/Users/li/Documents/truzhenv3worktree/truzhen-packs-pack-suite-retest-20260706`，只读）
- `truzhenos`（工作树：`/Users/li/Documents/truzhenv3worktree/truzhenos-pack-suite-retest-20260706`，测试）
- `truzhen-client-web-desktop`（工作树：`/Users/li/Documents/truzhenv3worktree/truzhen-client-pack-suite-retest-20260706`，测试）
- `truzhen-contracts`（`/Users/li/Documents/truzhen-contracts`，只读）
- `truzhen-software`（`/Users/li/Documents/truzhen-software`，只读）
- `truzhen-cloud`（`/Users/li/Documents/truzhen-cloud`，只读）

## 2. 前置准备（P0）

- `git status --short --branch` 逐仓执行并记录。
- 记录端口占用（18119 / 5189）和旧 pids。
- 读取 `/Users/li/Documents/过程文档/pack-suite-mainline-retest-20260706/round-mainline/round-4-20260706/baseline-status-4-20260706.json`，并补写新基线文件 `.../round-mainline/round-4-20260706/baseline-status-mainline-restart-20260707.json`。
- 确认所有路径、分支、只读范围不越界。

## 3. 环境启动（P1）

- 使用新的隔离运行目录：`/Users/li/Documents/过程文档/pack-suite-mainline-retest-20260706/round-mainline/round-4-20260706`。
- 启动参数：
  - 后端 `TRUZHEN_DEVSERVER_PORT=18119`
  - 前端 `TRUZHEN_FRONTEND_PORT=5189`
  - `TRUZHEN_PACKS_ROOT=/Users/li/Documents/truzhenv3worktree/truzhen-packs-pack-suite-retest-20260706`
  - `TRUZHEN_CLIENT_REPO_DIR=/Users/li/Documents/truzhenv3worktree/truzhen-client-pack-suite-retest-20260706`
- 日志落地：
  - `round-mainline/round-4-20260706/logs/manual-backend.log`
  - `round-mainline/round-4-20260706/logs/manual-vite.log`
- 记录：后端 PID、前端 PID、真实端口监听结果。

## 4. 三包生命周期（P2）

按家政→智能家居→环保顺序执行：

1. `install.py`
2. `uninstall.py`
3. `install.py`（再次）

每个 pack 需留存：

- install/uninstall/reinstall 日志。
- 安装后 30 秒内 registry 快照。
- uninstall 后 registry 快照。
- 03 回执反查原文。

## 5. 场景相关 GUI 与读模型（P3）

- 关键 UI 路径截图：
  - 场景包管理启用/状态；
  - 事务/对象列表；
  - 03 回执反查。
- 三包各自至少一次关键行为验证：
  - 家政：无知识库误报、生命周期关键节点；
  - 智能家居：单角色、`provider` 投影；
  - 环保：知识挂载与回执可反查。

## 6. 门禁与专项验证（P4）

- `truzhen-packs`：JSON、`py_compile`、禁止项扫描、结构审计。
- `truzhenos`：`go test ./backend/tests/devserver ./backend/tests/sceneflow ./backend/tests/execution ./backend/tests/registry`（可执行范围按可用性）与场景包针对测试。
- `truzhen-client-web-desktop`：`npm run typecheck`、相关 vitest、`npm run smoke:scene-flow-product`。
- 若以上无红项，执行 `GOWORK=off TRUZHEN_PACKS_ROOT=... TRUZHEN_CLIENT_REPO_DIR=... ./scripts/verify.sh`。

## 7. 产出（P5）

- `report-mainline-restart-20260706-closure.md`
- `issue-register-mainline-restart-20260706.md`
- `screenshot-index-mainline-restart-20260706.md`
- `run-stack-mainline-restart-20260706.json`
- `repo-final-state-mainline-restart-20260706.md`

## 8. 通过条件

- 每项结论必须有日志、API 读模型、receipt 及截图证据。
- 未出现真实外部动作。
- 未产生 `git` 修改（只读仓）或仅在文档/证据目录新增。
