# 能力 Pack 制作台短视频 OSS 首批切片验收报告（2026-07-04）

## 结论

本轮已把“短视频运营能力 Pack 制作任务”从计划推进到首批可验证切片，并补上 P3-P11 候选链：04 后端候选链、12 前端制作台展示、11 Code Assistant readiness / run gate 负例、`code_assistant_invocation_candidate` 候选信封、`code_assistant_run_request_candidate` 运行请求候选、candidate-only Pack bundle、candidate bundle dry-run 阻断、lifecycle preflight 阻断、PatchCandidate 结果承接、PatchCandidate 复核候选，以及 `truzhen-packs` 候选能力包资产均已接线并有验证证据。

当前生命周期档位：`契约已定 -> 已接线`。仍不能标为 `已验收 / 已发布`，因为独立 Capability Pack loader / lifecycle enabled、真实受控 Codex CLI run 产物、第三方 OSS 执行、社媒发布、正式 provider、contracts schema、云市场发布均未完成。

## 涉及仓库

| 仓库 | 本轮动作 | 边界 |
| --- | --- | --- |
| `truzhenos` | 04 Capability Studio 补 GitHub OSS 证据、`IntegrationPlan.CodeAssistantPolicy`、`MinimalGluePlan.AssemblyManifest`，并新增 Code Assistant 调用候选信封、11 run 请求候选、candidate bundle、candidate bundle dry-run、lifecycle preflight、PatchCandidate 结果承接、PatchCandidate 复核端点和 HTTP/session 回归 | 不运行 Codex CLI，不执行第三方 repo，不改 contracts |
| `truzhen-client-web-desktop` | 能力制作台 GUI 展示 OSS 集成计划、代码助手边界、装配清单、Code Assistant readiness，并新增“生成 Code Assistant 调用候选”“准备 Code Assistant 运行请求”“导出候选 Pack bundle”“校验候选 bundle dry-run”“检查 lifecycle preflight”“登记 PatchCandidate 结果”“提交 PatchCandidate 复核”按钮 / 卡片；live smoke 补 readiness / run gate 负例 / 候选信封、运行请求候选、bundle、dry-run、preflight、PatchCandidate intake 和 review 正例 | 前端只读投影，不保存 token，不生成真实胶水代码、不保存 draft、不应用补丁 |
| `truzhen-packs` | 新增短视频运营 Capability Pack 候选集、3 个能力包候选声明、Code Assistant 调用候选台账、运行请求候选台账、candidate bundle 导出台账、candidate bundle dry-run 台账、lifecycle preflight 缺口台账、PatchCandidate 承接台账和 PatchCandidate 复核台账 | 候选资产，不安装、不启用、不发布 |

## 产物

| 类型 | 路径 |
| --- | --- |
| 测试计划 + 执行登记 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/gui-capability-pack-workbench-github-oss-test-plan-20260704.md` |
| 短视频能力包候选集 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/` |
| Code Assistant 调用候选台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/code-assistant-invocation-ledger.md` |
| Code Assistant 运行请求候选台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/code-assistant-run-request-ledger.md` |
| Candidate bundle 导出台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/candidate-bundle-ledger.md` |
| Candidate bundle dry-run 台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/candidate-bundle-dry-run-ledger.md` |
| Lifecycle preflight 商用缺口台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/lifecycle-preflight-commercial-gap-ledger.md` |
| PatchCandidate 承接台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/patch-candidate-intake-ledger.md` |
| PatchCandidate 复核台账 | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/patch-candidate-review-ledger.md` |
| 04 后端接收说明 | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/capability-studio/short-video-oss-pack-studio-asset-acceptance-20260704.md` |
| 前端影响账本 | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/docs/frontend-adjustment-impact-ledger-20260704.md` |

## 候选能力包

| 候选 | pack ref | 样本 OSS | 当前状态 |
| --- | --- | --- | --- |
| 短视频草稿生成能力包 | `capability-pack://short-video-draft-generation` | MoneyPrinterTurbo | candidate-only；允许胶水补丁候选；未接 provider 时 `provider_missing` |
| 短视频合成编排能力包 | `capability-pack://short-video-composition-orchestration` | Pixelle-Video | candidate-only；允许胶水补丁候选；真实合成归 11 |
| 短视频发布草稿能力包 | `capability-pack://short-video-social-publish-draft` | social-auto-upload | blocked-by-default；只生成发布草稿 / 上传意图候选 |

## 验证证据

| 仓库 | 命令 / 结果 |
| --- | --- |
| `truzhenos` | `GOWORK=off go test ./backend/tests/capability -run 'Test(CapabilityStudioShortVideo|StudioShortVideoOSS|StudioTenStage|StudioSession|StudioSkip|StudioHTTP)' -count=1` 通过 |
| `truzhenos` | `TestCapabilityStudioShortVideoCodeAssistantInvocationCandidateCarriesP3LedgerBoundary` 覆盖 P3 调用候选信封边界；`TestStudioShortVideoOSSIntegrationPlanRoundTripsOverHTTP` 覆盖 `/v3/capability/studio/code_assistant_invocation` |
| `truzhenos` | `TestCapabilityStudioShortVideoCodeAssistantRunRequestCandidateRequiresBaseIssuedEnvelope` 覆盖 P8 运行请求候选、Base 待签发字段、不可执行和 Owner/Base Gate pending；HTTP 回归覆盖 `/v3/capability/studio/code_assistant_run_request` |
| `truzhenos` | `TestCapabilityStudioShortVideoCandidateBundleExportsPackFilesWithoutProviderExecution` 覆盖 P4-P6 候选 bundle 文件清单、证据贯通、不可安装 / 不含 provider 边界；HTTP 回归覆盖 `/v3/capability/studio/candidate_bundle` |
| `truzhenos` | `TestCapabilityStudioShortVideoCandidateBundleDryRunBlocksUnsupportedCapabilityLoader` 覆盖 P10 dry-run：静态校验通过但独立 Capability Pack loader 未接线时返回 `candidate_bundle_dry_run_blocked / capability_pack_loader_missing`；HTTP 回归覆盖 `/v3/capability/studio/candidate_bundle_dry_run` |
| `truzhenos` | `TestCapabilityStudioShortVideoLifecyclePreflightBlocksCandidateBundleBeforeDelivery`、`TestCapabilityStudioShortVideoLifecyclePreflightReportsProviderMissingReadiness` 覆盖 P11 preflight：candidate bundle 未形成正式 delivery 或 provider readiness 缺失时返回 `lifecycle_preflight_blocked`，并保持 `NoDraft=true`、`NoInstall=true`、`NoEnable=true`、`FormalWrite=false`；HTTP 回归覆盖 `/v3/capability/studio/lifecycle_preflight` |
| `truzhenos` | `TestCapabilityStudioShortVideoCodeAssistantRunResultIntakesPatchCandidateWithoutApply` 覆盖 P7 PatchCandidate 结果承接、11 回执、不可自动应用和需复核边界；HTTP 回归覆盖 `/v3/capability/studio/code_assistant_result` |
| `truzhenos` | `TestCapabilityStudioShortVideoPatchCandidateReviewKeepsApplyBehindGate` 覆盖 P9 PatchCandidate 复核、apply Gate pending、不可自动应用、不可安装 / 启用边界；HTTP 回归覆盖 `/v3/capability/studio/patch_candidate_review` |
| `truzhenos` | `GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1` 通过 |
| `truzhenos` | `git diff --check` 通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts`，34 tests passed |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "Code Assistant invocation"` 红绿通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "gated Code Assistant run request"` 红绿通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "candidate Pack bundle"` 红绿通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "dry-runs a candidate Pack bundle"` 红绿通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t lifecycle` 红绿通过，覆盖“检查 lifecycle preflight”按钮、请求体和 `LifecyclePreflightCard` |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "PatchCandidate run result"` 红绿通过 |
| `truzhen-client-web-desktop` | `npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "reviews a PatchCandidate"` 红绿通过 |
| `truzhen-client-web-desktop` | `src/api/__tests__/executionLiveSmokeScript.test.ts` 已纳入三文件回归，覆盖 `/v3/capability/studio/candidate_bundle_dry_run`、`/v3/capability/studio/lifecycle_preflight`、`candidate_bundle_dry_run_blocked`、`lifecycle_preflight_blocked`、`candidate_bundle_not_delivery`、`capability_pack_loader_missing`、`/v3/capability/studio/patch_candidate_review`、`patch_candidate_review_ready` 和 `ApplyGateRequired` |
| `truzhen-client-web-desktop` | `npm run typecheck`、`npm run smoke:frontend-shell`、`git diff --check` 通过 |
| `truzhen-client-web-desktop + truzhenos` | `TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior` 通过；Code Assistant readiness 为 `ready`，假决策 run 被 `blocked`，能力制作台生成 `code_assistant_invocation_candidate`、`code_assistant_run_request_candidate`、candidate-only bundle、`candidate_bundle_dry_run_blocked`、`lifecycle_preflight_blocked`、PatchCandidate intake 和 PatchCandidate review，evaluation 为 `provider_missing`，delivery 为 `evaluation_not_ready` |
| `truzhenos` | 2026-07-05 复验：`GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1`、`git diff --check`、`GOWORK=off bash scripts/verify.sh` 均通过；全量门禁包含 race scoped 测试，通信真实 E2E 因 `TRUZHEN_REAL_E2E_COMMUNICATION` 未开启按脚本策略跳过 |
| `truzhen-client-web-desktop` | 2026-07-05 复验：`npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts` 通过，34 tests passed；`npm run typecheck`、`npm run smoke:frontend-shell`、`npm run build`、`git diff --check` 均通过 |
| `truzhen-client-web-desktop + truzhenos` | 2026-07-05 复验：`TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior` 通过，确认 `lifecycle_preflight_blocked / candidate_bundle_not_delivery`，且不保存 draft、不安装、不启用 |
| `truzhen-packs` | JSON 合法、install / uninstall 脚本语法、禁入产物扫描、Pack 结构审计、`git diff --check` 均通过 |

## 未完成项

1. 独立 Capability Pack loader / lifecycle enabled 尚未接线；本轮 dry-run 只证明静态文件清单可校验，P11 preflight 只证明正式 delivery 缺口可展示，并诚实暴露 `candidate_bundle_not_delivery`。
2. 真实受控 Codex CLI run 尚未产出真实 PatchCandidate 文件 / artifact / receipt；本轮只生成了调用候选信封、11 run 请求候选、candidate-only bundle 文件清单、dry-run 阻断、PatchCandidate 结果承接和复核候选能力。
3. 能力制作台还未通过真实 GUI 操作录像或截图完整走完 3 个短视频候选包导出。
4. 本轮没有安装、运行或测试 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
5. 社媒上传仍是 blocked-by-default，未登录账号、未上传、未发送。
6. `truzhen-software` provider / adapter 候选未生成，contracts 未变更，云市场未发布。

## 下一步建议

1. 先进入 P12 安全内置能力 lifecycle 样本：只用基座安全 fixture 验证 draft / readiness / promote / confirm，不运行真实 Codex CLI、不执行第三方 OSS。
2. 再用 GUI 实操制作台跑 3 个短视频能力包候选导出，补截图 / 行为日志 / issue ledger。
3. 之后才考虑用 11 Code Assistant Gateway 做一次真实受控最小 run：只允许读取候选目录、生成 README 或 adapter scaffold 补丁候选，必须有 Owner + Base Gate + Receipt。
4. 若要进入商用路径，再拆 `truzhen-software` provider 候选和 `truzhen-cloud` 上架审核链；社媒发布保持红色动作，另行授权。
