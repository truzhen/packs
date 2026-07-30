# A00-R1 / P-E2 修复回边 1/2：环保 Pack installer 不再铸造 Owner evidence

日期：2026-07-30

仓库：`truzhen-packs`

分支：`codex/v4-unified-A00-R1-packs-owner-handoff-readiness-20260730`

基线提交：`d4d4fb0843856b171890a755accb791b7eeafeff`

生命周期：`修复回边 1/2 已实现 -> 待 P-T2 独立复验`

## 1. 触发原因

P-T2 独立验收报告：

- 路径：`/Users/li/.codex/truzhenv3-process/closeouts/truzhen-v4-unified-goal-A00-R1-P-T2-independent-validation-20260730.md`
- SHA-256：`1a6fd53cdc4e04c18e642984aefe88c12c270b43f4e0485050c1e46f0eabe5b9`
- 结论：`FAIL`

F-01 指出：旧 P-E2 在 os-14 已证明精确版本 `ready` 后，由 installer 自行生成 27 个 `evidence://` 引用，并对 Role enable-confirm、Agent Slot confirm 使用 `approve=true`，把调用方自造 proof 带入下游正式化动作。P-T2 的普通候选、部分态、恢复态和幂等边界通过，但 `no_caller_minted_evidence` 失败，因此旧候选不得进入 X-A3。

## 2. 修复范围与真相源

本次只修 `truzhen-packs` 环保 Pack installer 及其直接测试、治理说明和账本：

1. 保留 installer 对场景 lifecycle candidate staging、ReadModel readiness 和 Receipt 反查的既有职责。
2. 精确版本已启用、挂载与 FormalReceipt 均可反查后，installer 不再调用 Role Pack draft/enable-confirm、Agent Slot binding/confirm、Knowledge batch candidate/formalize。
3. installer 不生成 `evidence://`，不携带 `approve`、`confirm`、`formalize`、`decision_ref`、`owner_action_evidence_ref`、`run_id`、`nonce`、Origin 或 Cookie 权威材料。
4. 输出稳定、可机读的 `awaiting_owner_confirmation` handoff/resume：2 个 Role Pack、2 个 Agent Slot、45 个 Knowledge source 目标由可信 GUI 逐项完成；后端必须签发并核验 Owner evidence。
5. 该 handoff 只证明“下一步需要谁做什么”，不证明下游正式权威已经存在。

未新增或修改 HTTP route、DTO、schema、状态枚举、Client consumer、Saga payload、MemoryCenter 真相源或跨仓契约；未启动 devserver、GUI、登录、Provider、网络请求、真实装入、真实正式化、完整 EGR、push 或 merge。

## 3. 关键行为

### 3.1 精确 readiness

`scene_combined_ready` 仍要求：

- 场景 Pack 的精确 `pack_ref + version` 已启用；
- mount 指针与 lifecycle record 一致；
- readiness 中存在可反查 FormalReceipt；
- 15 knowledge scopes、45 documents、30 cases 与 audit refs 保持稳定。

满足上述条件后返回：

- `status=awaiting_owner_confirmation`
- `reason=downstream_owner_confirmation_required`
- `candidate_refs=[]`
- `owner_steps` 固定为 2 / 2 / 45 三组可信 GUI 任务

重复运行只执行同一组 GET，并返回完全相同的 handoff；不会重复创建候选，更不会创建正式事实。

### 3.2 部分态与恢复态

pointer-only、unknown、blocked recovery、active+blocked duplicate、缺 FormalReceipt 等边界继续 fail closed，并在首个不满足的事实处停止；不会进入下游写入。

## 4. 红绿证据

| 证据 | 结果 | 日志 / SHA-256 |
|---|---|---|
| 修复前新增对抗测试 | 9 项中 2 项失败：ready 行为仍写下游、源码禁止项命中 | `/tmp/A00-R1-P-E2-R1-red-20260730.log` / `c11d9608c48d1491e74f6dc3d7bd7530bb23a53efcaea105c3f9e18f74eae05c` |
| 修复后专项测试 | 9/9 PASS | `/tmp/A00-R1-P-E2-R1-green-20260730.log` / `f2a601175420b4894ca8a502f6fe0eeca0e0ab83a32c60404ca627adb4d13378` |
| P-T2 对抗矩阵复跑 | 动态 9/9；无 caller-minted evidence、无 installer formal confirm、无重命名权威绕过；部分/恢复态、重试和 15/45/30 均通过 | `/tmp/A00-R1-P-E2-R1-p-t2-adversarial-20260730.log` / `bb69580a5c0271d5b7f9dc4acbb78256069e86013ace59f6f49543f8b82e8122` |
| 实现与治理正文冻结后的 Packs 聚合门禁 | Python 40/40、专项 9/9、Go、192 JSON、脚本语法、Pack 结构、knowledge checksum、禁品、环保主权源码扫描、P-T2 矩阵与 `git diff --check` 全部 PASS | `/tmp/A00-R1-P-E2-R1-final-packs-gates-20260730.log` / `617f943f47cb0079587eca55997473c7a6b75cc5d95e2af38a1dc6e7be8cc163` |

## 5. 防回潮断言

专项测试同时锁定：

- installer 源码不得出现 `evidence://`；
- 不得出现 `approve=true`；
- 不得调用 Role enable-confirm、Agent Slot confirm 或 Knowledge batch 写入；
- exact-ready 重试只能读，不得产生候选或正式事实；
- 部分态、恢复态、重复态和缺 Receipt 必须提前停止；
- Role / Slot / Knowledge handoff 目标稳定为 2 / 2 / 45；
- 15 scopes、45 documents、30 cases、checksum、`pending_human_review`、`reference_only` 与 audit refs 不漂移。

## 6. 独立验收边界

本文件只记录执行节点完成了修复回边 1/2 和本地自测，不构成独立 PASS。下一步必须由 P-T2 使用裁剪后的新上下文，基于本次新提交亲自复跑对抗矩阵；只有 P-T2 的独立结论可以决定是否恢复进入 X-A3。
