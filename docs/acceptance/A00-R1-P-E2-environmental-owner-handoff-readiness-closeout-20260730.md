# A00-R1-P-E2 环保 installer Owner handoff 与组合 readiness closeout

## 结论

P-E2 已在 P-D1 clean commit `cc0ef86cfd7b9cc924a33772aea2bf76e00ba1ed` 上完成初次实现：

- 环保 `install.py` 不再调用场景 lifecycle `confirm/reactivate`、Base `prepare/confirm` 或知识候选 `approve`，也不自铸 Owner proof。
- 首次安装只执行 canvas 同步与 lifecycle `draft/readiness/promote` candidate staging，随后输出唯一的 `TRUZHEN_PACK_HANDOFF`，状态为 `awaiting_owner_confirmation`，包含可信 GUI Owner 操作说明与同命令 resume 信息。
- Owner 操作后重跑时，只有 exact enabled pointer、manifest 声明的 15 个 required scope 各自唯一 active mount，以及每个 `enabled_receipt_ref` 均可经 os-03 反查为同 ref 的 FormalReceipt，才允许继续。
- pointer-only、14/15 partial、active+blocked 冲突、blocked/recovery、FormalReceipt 缺失或 lifecycle 投影不完整均 fail closed；角色、槽位与知识下游不会在组合 readiness 前运行，并保留 mount ReadModel 可获得的 enabled/disabled/last Receipt refs 与 blocked reason。
- readiness 完整后，角色与槽位沿既有幂等链续接，45 份知识只暂存稳定 candidate，不自动 formalize；同一状态连续重跑时角色、槽位、source identity 与 candidate identity 不增加。
- 15 scopes、45 documents、30 cases、checksum、`pending_human_review` 与 `reference_only` 均未退化。

生命周期为 `已实现 -> 待 A00-R1-P-T2 独立验收`，未发布。

## 权威交接与边界

| 字段 | 本节点事实 |
| --- | --- |
| OS 独立验收 | O-T1 修复回边 1 候选 `f00d7efaafc853aa5bd4140c66c92f198721ed7b` PASS |
| O-T1 报告 | `/Users/li/.codex/truzhenv3-process/closeouts/truzhen-v4-unified-goal-A00-R1-O-T1-independent-validation-rerun-1-20260730.md` |
| O-T1 报告 SHA-256 | `914fd658569bc23d2bf462f091714afda49cb94bc4853f8b5f8076def4274476` |
| 实际消费 surface | `GET /v3/pack-studio/lifecycle/packs?pack_ref=...`、`GET /v3/memory/knowledge/mounts?...`、`GET /v3/receipts/{receipt_ref}` |
| 未消费 | Memory Center readiness、Saga journal/payload、failure/recovery lifecycle payload、error history |
| 契约影响 | 未新增或修改 route、DTO、schema、状态枚举、Client consumer；未自造公开 readiness 真相源 |
| 仓边界 | 只写 `truzhen-packs` 当前 worktree；未读写 Contracts、Client、Cloud、Software 产品状态 |
| 禁止动作 | 未启动 OS devserver/EGR，未登录，未调用 Provider/模型，未执行处罚、送达、安装、发布、网络或生产动作 |
| 外部动作 | `external_actions=0` |
| Packs 修复回边 | 初次实现，不消费独立验收回边；保持 `0/2` |

## 组合 readiness 判定

| 输入状态 | installer 输出 | 下游 |
| --- | --- | --- |
| 无目标 lifecycle record | 暂存 candidate，`awaiting_owner_confirmation` | 停止 |
| 已有 `pack_spec_candidate` | 只读等待，`awaiting_owner_confirmation` | 不重复 staging |
| 目标版本 record 缺失/未知 `state`，或 pointer 形状畸形 | `not_ready` | 停止且不重复 staging |
| 目标版本 pointer 缺失，但出现成功态 record | `recovery` | 停止 |
| pointer 指向其它版本 | `not_ready` | 停止 |
| exact pointer，但任一 required scope missing/pending/disabled/blocked | `not_ready` 或 `recovery` | 停止 |
| 同 scope 同时 active 与 blocked | `recovery`，保留失败 Receipt ref | 停止 |
| active mount 缺 enabled Receipt 或 Receipt 404/ref 不符 | `not_ready` | 停止 |
| exact pointer + 15/15 active mounts + 15/15 FormalReceipt | `ready` 内部判定 | 续接角色、槽位、45 个知识候选，最终 handoff 等待知识治理 |

Packs 没有读取不可公开的 Saga payload；当组合真相不完整时只投影已有三态，不用“没有看到 recovery Receipt”作为成功条件。

## TDD 证据

父提交快照 `cc0ef86cfd7b9cc924a33772aea2bf76e00ba1ed` 运行同一 9 项行为验收：9/9 RED。旧实现会自动 `reactivate/confirm`、调用 Base prepare/confirm 并 approve 知识；pointer-only 也会继续下游。

| 阶段 | 结果 | 日志 / SHA-256 |
| --- | --- | --- |
| RED：P-D1 父提交 | 9/9 FAIL，exit 1 | `/tmp/A00-R1-P-E2-red-20260730.log` / `4107033df8d5cfef8d9fc4bb26658f1af72293fe4daf31c7c970f9b7726f8d39` |
| GREEN：P-E2 当前实现 | 9/9 PASS，exit 0 | `/tmp/A00-R1-P-E2-green-20260730.log` / `71fd6d87ed46bd4b806f485c9310eb0381665caed28fc608fc46908705d560c3` |

行为矩阵覆盖首次 staging、已有 candidate 重试、目标 record 缺失 state、pointer-only 14/15、blocked recovery audit ref、active+blocked 冲突、FormalReceipt 缺失、15/15 ready 续接和两次重跑 45 candidate identity 稳定。

初次实现自测还发现并根修三项：

1. 首次 staging 残留旧局部变量 `install_version`，触发 `NameError`；改为唯一 manifest `version`。
2. 同一 scope 同时返回 active 与 blocked 时，旧组合判断只看 active filter 会误放行；现先检查 exact scope 的全部 mounts，任一非 active 即 fail closed。
3. 目标版本 record 存在但缺失/携带未知 `state` 时，旧判断会误当首装空态并重复 staging；现严格消费 O-T1 枚举，畸形 record/pointer 受控投影为 `not_ready`。

仓级门发现两条旧测试契约与新边界冲突，均更新为更严格的当前语义：

- Python 高风险知识测试不再要求存在自动 approve 的 `verify_authority=False`，改为要求 `pending_human_review` 且禁止任何 `verify_authority`/知识 approve。
- Go 错误码 guard 不再要求环保 installer 接线已删除的版本 bump/Base 阶段，改为明确要求 connectivity、lifecycle、readiness、role、knowledge 与 checksum 六类仍可发生的阶段码。

## 全仓门禁

| 门 | 结果 | 日志 / SHA-256 |
| --- | --- | --- |
| 代码与文档冻结后的聚合全门 | 40/40 + 9/9 Python、Go、192 JSON、语法、结构、checksum、禁品、主权边界与 diff 全部 PASS | `/tmp/A00-R1-P-E2-final-packs-gates-20260730.log` / `160dbc5ec6749eac09e867b5f4e3f2da168d005bb92530e2d42215a784f12db0` |

本报告、实现、测试、账本和治理同步必须进入同一个新的 clean local commit。不得 amend P-D1，不 merge、不 push、不启动 P-T2。
