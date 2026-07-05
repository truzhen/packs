# Lifecycle preflight 商用缺口台账

> 日期：2026-07-04
> 状态：`已接线`
> 本文件记录短视频运营三类 Capability Pack 候选进入正式 lifecycle 前的缺口。它不是 enabled 证据，不代表 provider 已接通。

## 1. 总结

P10 candidate bundle dry-run 已能证明候选文件清单静态可校验，并诚实返回 `capability_pack_loader_missing`。P11 已在 `truzhenos` 与 `truzhen-client-web-desktop` 接线为 `/v3/capability/studio/lifecycle_preflight` 和“检查 lifecycle preflight”按钮，能把候选 bundle 不是正式 delivery 的事实展示为 `lifecycle_preflight_blocked / candidate_bundle_not_delivery`。它仍不是 lifecycle draft、enabled 或商用证据。

## 2. 候选缺口矩阵

| 候选 | 当前资产 | P11 预期状态 | 商用缺口 | 归属 |
|---|---|---|---|---|
| `capability-pack://short-video-draft-generation` | candidate bundle、Code Assistant 候选、PatchCandidate 复核候选 | `lifecycle_preflight_blocked / candidate_bundle_not_delivery` | 需要正式 studio delivery；黄金用例需 ready；MoneyPrinterTurbo 只能作 OSS 证据和胶水候选，不能当 provider ready | `truzhenos` + client；provider 候选未来归 `truzhen-software` |
| `capability-pack://short-video-composition-orchestration` | candidate bundle、Pixelle-Video OSS evidence、合成编排能力声明 | `lifecycle_preflight_blocked / candidate_bundle_not_delivery` | 需要合成 provider dependency；媒体生成 / 剪辑 / 模型能力不能由 pack 仓实现；真实执行必须 11 Gateway | `truzhenos` + `truzhen-software` |
| `capability-pack://short-video-social-publish-draft` | blocked-by-default 发布草稿能力声明 | `lifecycle_preflight_blocked / external_send_risk` | social-auto-upload 只能作为证据；真实登录、上传、发布是红色动作，需账号 secretref、Owner/Base Gate、Communication / Execution Gateway 和 Receipt | `truzhenos` + 10/11 Gateway + 16 Auth，未来可涉及 `truzhen-cloud` |

## 3. P11 已接线证据

- 后端返回 `lifecycle_preflight_blocked` 而不是 install success。
- 返回 `candidate_bundle_not_delivery` 时，UI 解释“候选 bundle 不是正式 delivery”。
- 返回 readiness issue：`delivery_artifact_required`。
- `FormalWrite=false`、`NoDraft=true`、`NoInstall=true`、`NoEnable=true`、`NoRealExecution=true`。
- 阻断项给出下一步：`delivery`、`evaluation_ready`、`owner_base_gate`。

## 4. 本批验证

- 后端定向：`GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityStudioShortVideoLifecyclePreflight|TestStudioShortVideoOSSIntegrationPlanRoundTripsOverHTTP' -count=1`。
- 前端定向：`npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t lifecycle`。
- 本批未运行真实 Codex CLI，未执行第三方 OSS，未登录或上传社媒。

## 4.1 2026-07-05 复验记录

- 后端工作树：`/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss`。
- 前端工作树：`/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui`。
- 后端定向验证：`GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1` 通过；`git diff --check` 通过。
- 前端验证：`npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts` 通过，34 tests passed；`npm run typecheck` 通过；`npm run smoke:frontend-shell` 通过；`npm run build` 通过；`git diff --check` 通过。
- 跨仓 live smoke：`TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior` 通过，确认 `lifecycle_preflight_blocked / candidate_bundle_not_delivery`，并确认 candidate bundle 不保存 draft、不安装、不启用。
- 后端全量门禁：`GOWORK=off bash scripts/verify.sh` 通过；其中 `TRUZHEN_REAL_E2E_COMMUNICATION` 未开启，通信真实 E2E 按脚本策略跳过。
- 复验边界：未运行真实 Codex CLI，未执行 MoneyPrinterTurbo / Pixelle-Video / social-auto-upload，未登录或上传社媒，未修改 `truzhen-contracts`、`truzhen-software`、`truzhen-cloud`。

## 5. 禁止冒充

- 不得把 P10 静态校验通过当成 lifecycle draft。
- 不得用 Domain Work Pack loader 代替 Capability Pack lifecycle。
- 不得在本仓保存 provider 实现、运行脚本、token、cookie、账号或上传回执。
- 不得把社媒发布草稿能力写成已能上传。

## 6. 下一步

若 Owner 继续授权，下一步建议进入 P12 安全内置能力 lifecycle 样本：只使用基座安全 fixture 验证 draft / readiness / promote / confirm 最小闭环。短视频第三方 OSS 仍保持 evidence-only / blocked；本批到 P11 preflight 为止，不进入正式 enabled。

P12 执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`

P12 授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`
