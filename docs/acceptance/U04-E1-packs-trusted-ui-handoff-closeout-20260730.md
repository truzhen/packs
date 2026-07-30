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

## U04-T1 FAIL 与修复回边 1/2

独立验收 `U04-T1` 于 2026-07-30 封存 FAIL：同一目标 `pack_ref` 同时返回 `1.1.0` 与冲突 `9.9.9` 时，旧 helper 首条记录短路并误报成功。报告 `/Users/li/.codex/truzhenv3-process/closeouts/truzhen-v4-unified-goal-U04-T1-packs-independent-acceptance-20260730.md`，其 SHA256 为 `bb47fcc799387b42081b4455f74942baec11664319cfd8ddc5d04b705620dfef`。

修复回边 `1/2` 仅调整既有 `pack_enabled_version_from_readmodel`：目标记录必须恰好一条，`current_version` 必须是未归一化的 canonical 版本字符串；重复目标记录（包括同值）、enabled/disabled 冲突、前后空白、空白-only 与非字符串均返回非法状态。三条消费路径均在初始 lifecycle ReadModel 处 fail closed。

| 命令 | exit | 日志 / SHA256 |
| --- | --- | --- |
| 新反例红灯：重复冲突与三消费路径 | 1 | `/tmp/U04-E1-R1-red-lifecycle-readmodel-20260730.log` / `26c3595d02f58fb16e6e4ddfec7f5b48fb0dc5fefced2fce1f38ead9ab6744c9` |
| 修复后对抗测试 | 0 | `/tmp/U04-E1-R1-focused-lifecycle-readmodel-20260730.log` / `9cfb6ac2dcad4d4640e4a2d12e3753cd2a0cd710238a511ab402c0da8d9ac3d4` |
| 回边后完整 E1 Packs 门 | 0 | `/tmp/U04-E1-R1-final-packs-gates-20260730.log` / `716e484b042f40fc1607db53a3d3d00529234ceb4ad9450991cbbd8a68e448f8` |

这轮只修复 T1 的解析根因，不启动 T1-R1、T2、OS 兼容验收或任何外部动作。候选仍为`已实现 -> 待独立验收`，应由全新独立验收上下文重新判断。

## U04-T1-R1 FAIL 与修复回边 2/2（最后一次）

独立复验 `U04-T1-R1` 于 2026-07-30 封存 FAIL：os-07 的既有 dict comprehension 会以最后一条记录覆盖同一 `transaction_ref` 的前序状态，使 `paused -> active`（装入）或 `active -> paused`（停用）冲突被误判确认。报告 `/Users/li/.codex/truzhenv3-process/closeouts/truzhen-v4-unified-goal-U04-T1-R1-packs-independent-reacceptance-20260730.md`，SHA256 `9f4a2e13962cc1f35a9e834e4b1521896c4f396a959f5e37ce003e35e5d3a7af`。

修复回边 `2/2` 仅调整既有 `wait_for_owner_schedule_states`：声明 `transaction_ref` 必须唯一且合法；响应内每条 schedule 必须是合法对象，`transaction_ref` 与 `status` 必须是未归一化的 canonical 字符串；声明引用在 ReadModel 中必须恰好一条。重复同值、冲突、非对象 sibling、缺字段、空白/非字符串状态和非 canonical 引用一律 fail closed。未声明但合法的计划仍按原语义忽略，`allow_missing=True` 仍仅表示声明引用确实不存在。

| 命令 | exit | 日志 / SHA256 |
| --- | --- | --- |
| 新反例红灯：schedule 重复/畸形与 content 双路径 | 1 | `/tmp/U04-E1-R2-red-schedule-readmodel-20260730.log` / `9b1e1dfb77a84ad2737fa02a27d1cd8c94c560904e78827596875ade79deb3ca` |
| 修复后 schedule 对抗矩阵 | 0 | `/tmp/U04-E1-R2-focused-schedule-readmodel-20260730.log` / `5d0e586d193c5184e047f0b6468ee812eb8537779b30df0bc7b018d901ac602a` |
| 回边后完整 E1 Packs 门 | 0 | `/tmp/U04-E1-R2-final-packs-gates-20260730.log` / `4d954e4075a7d4484ef680f7f5c2696cb333f32b73930e995220b7ac01ed4a81` |

这是 sealed acceptance 的最后一次修复预算。不得在本节点启动 T1-R2、T2 或第三次修复；若下一次独立验收仍 FAIL，必须停止并报告。候选保持`已实现 -> 待独立验收`。

## U04-R2-P-E1：首次安装空态精确兼容

`U04-T2` 的独立验收报告 `/Users/li/.codex/truzhenv3-process/closeouts/truzhen-v4-unified-goal-U04-T2-packs-first-install-openapi-acceptance-20260730.md` 确认：真实 os-14 首次安装查询会返回目标 `pack_ref`、`records: []`，并且因为尚未启用而**省略** `enabled_pointer`。旧共享解析器把该合法空态当畸形，三个可信 GUI 交接脚本都不能进入只读 handoff。

本回边从 Packs `cb17c7f1643cc3161e05cf732685003dadcd83b4` 开始，只改既有 `pack_enabled_version_from_readmodel` 和直接测试：只有目标 entry 的 `records` 是精确空数组且 `enabled_pointer` 字段缺失时，才返回标准空态 `""`。目标 entry 缺失、重复 entry、`records` 缺失/非数组、`records` 非空却缺 pointer、pointer 非对象、非 canonical 版本和冲突形状全部继续 `None` / fail closed；已有 disabled/reactivate 与已安装指针行为不放宽。

固定 OS handler `/Users/li/Documents/truzhenos/backend/internal/devserver/packcapabilityhttp/pack_studio_stage4.go:923-937` 只读核对显示，每个目标 entry 都带 `records`；存在 `enabled_pointer` 时它对应 lifecycle record。故既有 issued-binding 测试 fixture 改为含 `pack_ref`、`version`、`state=enabled|disabled` 的非空 record，而非用空数组机械放绿。`smart-home-owner-pack-v0/uninstall.py` 未改。

| 命令 | exit | 日志 / SHA256 |
| --- | --- | --- |
| 新反例红灯：合法首装空态 | 1 | `/tmp/U04-R2-P-E1-red-first-install-20260730.log` / `2fe21ed7e29205b15d9d82f955cc5e7e28ba3de9d02cd90d14c9648b57880519` |
| 首装、已安装、disabled/reactivate、uninstall、重复/畸形与 issued-binding 定向矩阵 | 0 | `/tmp/U04-R2-P-E1-focused-first-install-final-20260730.log` / `5cf54de757d9b0b1df20403a28759108c3dab995333b04534a829e88b9392144` |
| 代码冻结后全量 Python、Go、JSON、语法、Pack 结构、禁品、敏感项、智能家居卸载不变与 `git diff --check` | 0 | `/tmp/U04-R2-P-E1-final-packs-gates-rerun-20260730.log` / `792c5afb8dfa3f0a130db601246f137b3c153194dd2214190b3a839741fd6149` |

本回边生命周期为`已实现 -> 待独立验收`：未合并、未推送、未启动 OS devserver / EGR / Provider / 登录或外部动作。下一节点只能由新的独立 `P-T1` 在隔离 OS ReadModel 上验收这条首装形状；本节点不启动 P-T1 或 O-E2。
