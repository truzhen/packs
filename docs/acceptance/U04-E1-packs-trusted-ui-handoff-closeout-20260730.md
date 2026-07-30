# U04-E1 Packs 可信 GUI 交接 closeout（2026-07-30）

## 范围与状态

- 候选分支：`codex/v4-unified-U04-packs-trusted-ui-handoff-20260730`
- 基线：Packs `de02cac54c0b3623a53fac8faa0079ae4f31c4e8`
- 本次只改内容运营的 `install.py` / `uninstall.py`、智能家居的 `install.py`、共享 `pack_diagnostics.py`、直接测试与治理文档。
- `smart-home-owner-pack-v0/uninstall.py` 未改动。
- 生命周期：`已实现 -> 待独立验收`。未合并、未推送、未启动 OS devserver、未执行 Provider、登录或外部动作。

## 实现事实

三个目标脚本都拒绝非 `GET` 方法与非空请求体；不构造或提交 Base prepare/confirm、decision、evidence、Receipt、lifecycle 或计划写入。

- 内容装入：只读要求 os-14 返回 manifest 的精确 enabled 版本，并要求 os-07 全部声明计划 `active`。
- 内容停用：只读要求 os-14 停用状态，并要求 os-07 声明计划 `paused`、`cancelled` 或缺失。
- 智能家居装入：只读要求 os-14 返回 manifest 的精确 enabled 版本。
- ReadModel 缺字段、版本不符、计划仍 `active`、连接异常或超时均 fail closed；前台提示只表达目标、当前状态与下一步。

## TDD 与验证证据

先运行新增 handoff 测试，旧脚本红灯：内容/智能家居 install 仍试图写 lifecycle、角色与计划；内容 uninstall 缺只读 handoff。命令 `python3 -B -m unittest -v test_uninstall_owner_handoff.py`，exit `1`，日志 `/tmp/U04-E1-red-handoff-tests-20260730.log`，SHA256 `ce55f58005539761c823e2bd7439b493c27a1ba33158de4901254d3c898b365e`。

转绿后的证据：

| 命令 | exit | 日志 / SHA256 |
| --- | --- | --- |
| `python3 -B -m unittest -v test_uninstall_owner_handoff.py test_pack_issued_binding.py` | 0 | `/tmp/U04-E1-focused-python-rerun-20260730.log` / `4749ed1775db57983ba9004d0baa1191f6ab283f80e07a3d34e462dfdd9e22e3` |
| `python3 -B -m unittest discover -v` | 0 | `/tmp/U04-E1-python-discovery-20260730.log` / `7cd925dda00dcf922bc7874911ab682dbaf0dcd5a72cd40c7c01eb492ebef4e0` |
| `GOWORK=off go test ./...` | 0 | `/tmp/U04-E1-go-test-20260730.log` / `4214de23505807e67d2e025d96adcd381b9d41ed895672e908d7279301b42256` |
| JSON、install/uninstall `py_compile`、Pack 结构、禁品、敏感新增行与 `git diff --check` | 0 | `/tmp/U04-E1-static-recheck-20260730.log` / `6647ab3e776fd8c420faa7a0744d1e1c5c8f19182758dc205708c1640669a207` |
| 代码冻结后的全量 Python、Go、JSON、语法、结构、禁品、目标代码敏感项、范围与 `git diff --check` | 0 | `/tmp/U04-E1-final-packs-gates-20260730.log` / `aa52b7b0aa2afaa1bc59600c64e002a52495b1aaa30a6f8cd07c68b882edcf24` |
| 精确变更白名单、目标代码敏感项、智能家居 `uninstall.py` 不变与最终 `git diff --check` | 0 | `/tmp/U04-E1-final-scope-20260730.log` / `a9462e0eba150dc85a4ed1a28efbd3edfb5077affc28e067f5f2756cd5ae247d` |

所有命令都在候选 worktree 运行；其中 Python discovery 约 0.24 秒，Go 测试约 2.62 秒。初次聚合静态扫描误把既有治理文档的历史绝对路径当作本次敏感内容，已缩小为仅扫描本次新增的运行/测试代码行；该修正日志不作为通过证据。

## 独立验收下一步

本候选不替代 OS 兼容验收。独立验收应在新的隔离 OS 实例仅作 ReadModel 验证：确认 `GET /v3/pack-studio/lifecycle/packs` 与 `GET /v3/task-governance/schedules` 形状、精确版本、计划状态和缺字段 fail-closed 行为；不得启动 OS EGR、不得写入、不得使用真实身份或 Provider。通过后才由协调线程按任务图决定后续集成。
