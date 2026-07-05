# 能力 Pack 制作台商用完成审计

> 日期：2026-07-04
> 归属：`truzhen-packs` 文档审计。本文件只记录当前商用缺口和验收证据要求，不代表 Pack 制作台已商用。
> 关联计划：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/gui-capability-pack-workbench-github-oss-test-plan-20260704.md`

## 1. 审计结论

当前能力 Pack 制作台已完成 P3-P11 候选链：GitHub OSS 证据、Code Assistant 调用候选、11 run 请求候选、candidate bundle、candidate bundle dry-run、PatchCandidate 承接 / 复核和 lifecycle preflight 阻断。

这还不能称为商用。商用最小闭环至少还缺：

1. P12 安全内置能力 lifecycle 样本：draft / readiness / promote / confirm 有后端和前端实证。
2. P13 GUI lifecycle 面板：用户不能绕过 readiness / Gate / Receipt，能看到阻断原因和下一步。
3. 3 个短视频能力候选的 GUI 实操证据：截图 / 行为日志 / 网络响应摘要 / issue ledger。
4. 真实 Code Assistant run 的受控回执：必须经 11 Gateway、Owner/Base Gate 和 Receipt；当前只到 run request candidate。
5. provider / adapter 候选归属：真实 adapter 不能进 `truzhen-packs`，后续应归 `truzhen-software` 或外部 provider。
6. 社媒发布红线：social-auto-upload 只能保持 blocked，真实登录 / 上传 / 发布需单独红色授权。
7. 云市场链路：Pack 商品、License / Entitlement、下载分发归 `truzhen-cloud`，本轮未接线。
8. 商用 go/no-go 写回门禁：P12-P18 的证据写回必须在机器门禁中归零，不能只靠人工台账或最终证据包说明。

## 2. 商用完成门槛

| 门槛 | 当前证据 | 当前判定 | 还缺什么 |
|---|---|---|---|
| 后端能生成能力 Pack 候选 | P3-P11 后端端点和回归已接线 | 部分满足 | P12 lifecycle 安全样本闭环；真实 Code Assistant run receipt；apply 候选 Gate。 |
| 前端能让用户完成制作流程 | GUI 已展示 OSS 证据、Code Assistant 候选、bundle、dry-run、preflight、PatchCandidate 复核 | 部分满足 | P12/P13 lifecycle GUI；3 个候选包 GUI 实操截图 / 行为日志。 |
| 候选与正式事实隔离 | P10/P11 已返回 `candidate_bundle_dry_run_blocked`、`lifecycle_preflight_blocked` | 满足当前切片 | 后续 enabled 时必须证明 Owner/Base Gate + Receipt。 |
| lifecycle enabled 可控 | P11 只证明 candidate bundle 不是 delivery | 未满足 | P12 安全 fixture 证明 draft / readiness / promote / confirm。 |
| Code Assistant 可控执行 | 已有 invocation candidate 和 11 run request candidate | 未满足 | 真实最小 run 需单独授权，不读取 raw token，不执行第三方 OSS。 |
| OSS 证据可追溯 | MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 已登记为样本 | 部分满足 | 需要 GUI 查询证据截图、license 人工复核、repo 变更时间记录。 |
| provider / adapter 归属清楚 | 文档已声明不进 packs | 部分满足 | 需要 `truzhen-software` 授权后生成 adapter candidate 和 readiness check。 |
| 社媒外发受控 | 发布能力 blocked-by-default | 满足当前切片 | 真实发布需账号 secretref、Owner/Base Gate、Gateway、Receipt 和单独红色授权。 |
| 云市场可售卖 | 本仓只有候选资产和商品化口径 | 未满足 | `truzhen-cloud` 上架、License / Entitlement、下载分发、支付沙箱均未接线。 |
| 商用证据写回门禁 | `commercial-go-no-go-gate.json#evidence_writeback_gate` 已阻断完成声明 | 未满足 | 需要 `docs/commercial-cross-repo-execution-queue.json#execution_entries[].evidence_writeback_plan_summary` 与 `docs/commercial-readiness-verifier.json#current_blockers.evidence_writeback_summary` 全部完成，并满足 `total_pending_entry_count=0`。 |

机器可检位置：

```text
docs/commercial-cross-repo-execution-queue.json
docs/commercial-cross-repo-execution-queue.json#execution_entries[].evidence_writeback_plan_summary
docs/commercial-readiness-verifier.json
docs/commercial-readiness-verifier.json#current_blockers.evidence_writeback_summary
docs/commercial-go-no-go-gate.json
docs/commercial-go-no-go-gate.json#evidence_writeback_gate
docs/commercial-go-no-go-gate.json#completion_rule.evidence_writeback_gate_must_pass=true
all_evidence_writebacks_completed_and_verified
```

## 3. 后续切片

| 切片 | 目标 | 目标仓 | 状态 | 进入条件 |
|---|---|---|---|---|
| P12 | 安全内置能力 lifecycle 样本 | `truzhenos` + `truzhen-client-web-desktop` | 待 Owner 授权 | 回复 P12 授权卡；本批不改 contracts / software / cloud。 |
| P13 | GUI lifecycle 面板 | `truzhen-client-web-desktop`，必要时 `truzhenos` | 待 Owner 授权；规格和授权卡已起草 | P12 后端和前端最小闭环存在。 |
| P14 | 商用缺口审计 | `truzhen-packs` | 本文件已起草 | P12/P13 后更新为实证审计。 |
| P15 | GUI 实操 3 个短视频候选包 | `truzhen-client-web-desktop` + `truzhenos` + `truzhen-packs` | 待授权；规格、授权卡和证据台账骨架已起草 | 允许录屏 / 截图 / 行为日志；仍不执行第三方 OSS。 |
| P16 | 真实受控 Code Assistant 最小 run | `truzhenos` 11 Gateway + client | 待单独授权；规格、授权卡和证据台账骨架已起草 | Owner/Base Gate + Receipt；不执行第三方 repo；只生成 PatchCandidate 文件。 |
| P17 | provider / adapter candidate | `truzhen-software` 或外部 provider 仓 | 待单独授权；规格、授权卡和证据台账骨架已起草 | 明确目标 provider、隔离环境、测试命令和禁入凭据。 |
| P18 | 云市场 sandbox 链路 | `truzhen-cloud` + client | 待单独授权；规格、授权卡和证据台账骨架已起草 | 商品草稿、License / Entitlement、下载分发、支付沙箱边界明确。 |

## 4. P12 前置授权

P12 是当前最小可前进切片，因为它不用第三方 OSS、不运行 Codex CLI，也不碰社媒账号。推荐授权语：

```text
授权按 P12 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒；只使用基座内置安全样本或 fixture 验证 lifecycle draft/readiness/promote/confirm。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

## 4.1 P13 / P15 后续授权入口

P13 GUI lifecycle 面板规格：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-gui-lifecycle-panel-execution-spec-20260704.md`

P13 授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-cross-repo-execution-authorization-20260704.md`

P15 三候选 GUI 实操验收规格：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-three-candidate-gui-walkthrough-spec-20260704.md`

P15 授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-cross-repo-execution-authorization-20260704.md`

## 4.2 P16 / P17 / P18 后续授权入口

P16 受控 Code Assistant 最小 run 规格：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-controlled-code-assistant-run-spec-20260704.md`

P16 授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-cross-repo-execution-authorization-20260704.md`

P17 provider / adapter candidate 规格：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-provider-adapter-candidate-spec-20260704.md`

P17 授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-cross-repo-execution-authorization-20260704.md`

P18 云市场 sandbox 规格：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cloud-market-sandbox-spec-20260704.md`

P18 授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cross-repo-execution-authorization-20260704.md`

## 5. 禁止把以下内容当成完成

- P10 静态文件校验通过。
- P11 `lifecycle_preflight_blocked`。
- 只存在 candidate bundle。
- 只存在 run request candidate。
- 只存在 PatchCandidate review candidate。
- 只存在 README / 计划 / mock / fixture 说明。
- 社媒发布能力只写了 blocked。
- 云商品字段只写在本仓候选资产中。

## 6. 完成证据包要求

商用前必须形成一份最终证据包，至少包含：

1. 后端测试输出：04 capability、01 Gate、03 Receipt、11 Gateway 相关测试。
2. 前端测试输出：制作台 GUI 单测、typecheck、shell smoke、behavior smoke。
3. GUI 证据：3 个短视频能力候选包的截图 / 行为日志 / 网络响应摘要。
4. lifecycle 证据：draft、readiness、promote、confirm、enabled pointer、receipt ref。
5. Code Assistant 证据：真实 run receipt 或明确 `provider_missing / blocked`。
6. provider 证据：每个 ProviderRequirement 的 `ready / provider_missing / blocked` 归属。
7. 禁入证据：无 raw token、无第三方源码、无构建产物、无社媒登录态。
8. 云边界证据：若涉及商品化，必须来自 `truzhen-cloud`，不能由本仓冒充。
9. 写回门禁证据：`evidence_writeback_gate_must_pass=true`，`total_pending_entry_count=0`，且 `all_evidence_writebacks_completed_and_verified` 已由独立验收主体复核。

最终证据包模板：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-commercial-go-live-evidence-package-20260704.md`

授权路线图：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-authorization-roadmap-20260704.md`
