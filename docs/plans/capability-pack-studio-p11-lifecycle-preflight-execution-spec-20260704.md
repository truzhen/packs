# P11 能力 Pack lifecycle preflight 执行规格

> 日期：2026-07-04
> 状态：`已接线`
> 归属：`truzhen-packs` 计划规格。本文记录 P11 已执行施工切片，不代表已进入 lifecycle draft、enabled 或商用发布。
> 前置影响清单：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-loader-lifecycle-commercialization-impact-20260704.md`
> 跨仓授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-cross-repo-execution-authorization-20260704.md`

## 1. 目的

P11 的目标不是让 `candidate_bundle_dry_run` 变成安装成功，而是让制作台明确回答：一个短视频运营 Capability Pack 候选离正式 lifecycle draft / readiness / promote / confirm 还差什么。

当前 P10 证明候选 bundle 静态文件可校验，但独立 Capability Pack loader 未接线时必须 blocked。P11 在这个基础上增加 lifecycle 前置检查，让 GUI 直接展示候选 bundle、正式 delivery、evaluation readiness、provider dependency、Owner/Base Gate 和 Receipt 的缺口。

## 2. 已有事实源

| 事实 | 当前位置 |
|---|---|
| 候选 bundle / dry-run 服务 | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/service.go` |
| 能力 lifecycle readiness 规则 | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/lifecycle_spec.go` |
| lifecycle enable chain 测试样板 | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/capability_lifecycle_chain_test.go` |
| P10 API 契约 | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/api_contract.md` |
| P10 前端卡片 | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/components/capability-studio/CapabilityStudioSurfaces.tsx` |
| P10 前端页面接线 | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx` |

## 3. 后端接口草案

建议新增：

```text
POST /v3/capability/studio/lifecycle_preflight
```

请求：

```json
{
  "session_id": "capability-studio-session://...",
  "actor_ref": "owner",
  "candidate_bundle_ref": "candidate-bundle://...",
  "preflight_evidence_ref": "preflight-evidence://...",
  "intended_action": "draft",
  "notes": []
}
```

响应仍走 `StageOutcome`：

```json
{
  "ok": true,
  "session_id": "...",
  "stage": "lifecycle_preflight",
  "receipt_ref": "receipt_candidate://...",
  "artifact": {
    "Preflight": {
      "Status": "lifecycle_preflight_blocked",
      "BlockedReason": "candidate_bundle_not_delivery",
      "CandidateBundleRef": "candidate-bundle://...",
      "DeliveryPresent": false,
      "EvaluationStatus": "provider_missing",
      "ReadinessIssues": [
        {"Key": "delivery_artifact_required", "Severity": "error"},
        {"Key": "provider_missing", "Severity": "error"}
      ],
      "RequiredNextSteps": [
        "生成正式 studio delivery",
        "补黄金用例并让 evaluation readiness=ready",
        "补 provider dependency readiness",
        "进入 lifecycle draft/readiness/promote/confirm"
      ],
      "CandidateOnly": true,
      "FormalWrite": false,
      "NoDraft": true,
      "NoEnable": true,
      "NoRealExecution": true
    }
  }
}
```

## 4. 后端规则

1. 缺 session：沿用现有 `requireSession` 行为。
2. 缺 `candidate_bundle`：返回 `409 stage_prerequisite_missing`，`missing_stage=candidate_bundle`。
3. 只有 candidate bundle、没有 delivery：返回 200 + `lifecycle_preflight_blocked / candidate_bundle_not_delivery`，让前端可展示。
4. 有 delivery 但 evaluation 缺失：返回 `lifecycle_preflight_blocked / evaluation_missing`。
5. evaluation 是 `provider_missing` 或 `not_ready`：返回 blocked，并生成 readiness issue。
6. delivery 存在且 evaluation `ready`：返回 `lifecycle_preflight_ready_for_draft`，但仍不创建 draft。
7. P11 不调用 `lifecycle.SaveDraft`、不写 enabled pointer、不写 Formal Receipt、不运行 provider、不执行 Codex CLI。
8. P11 只写制作台 stage artifact / receipt candidate，供 GUI 和 issue ledger 消费。

## 5. 后端施工点

| 文件 | 变更 |
|---|---|
| `backend/internal/capability/capability_studio_lifecycle_preflight.go` | 新增 `CapabilityLifecyclePreflightDraft` / `CapabilityLifecyclePreflightCandidate` 和 builder |
| `backend/internal/capability/studio/service.go` | 新增 `LifecyclePreflight`，加载 `candidate_bundle`、`delivery`、`evaluation` 后生成候选 |
| `backend/internal/capability/studio/http.go` | 新增 handler 和 route `/v3/capability/studio/lifecycle_preflight` |
| `backend/internal/capability/studio/store.go` | stage 白名单加入 `lifecycle_preflight` |
| `modules/04-capability-management/api_contract.md` | 记录 P11 API 和候选边界 |
| `modules/04-capability-management/acceptance.md` | 记录 P11 blocked / ready 判定 |

建议测试：

- `TestCapabilityStudioShortVideoLifecyclePreflightBlocksCandidateBundleBeforeDelivery`
- `TestCapabilityStudioShortVideoLifecyclePreflightReportsProviderMissingReadiness`
- `TestStudioShortVideoOSSIntegrationPlanRoundTripsOverHTTP` 增加 `/lifecycle_preflight` HTTP roundtrip

## 6. 前端施工点

| 文件 | 变更 |
|---|---|
| `src/api/client.ts` | 新增 `runCapabilityStudioLifecyclePreflight` |
| `src/components/capability-studio/CapabilityStudioSurfaces.tsx` | 新增 `LifecyclePreflightCard`，展示 delivery / evaluation / readiness / Gate / Receipt 缺口 |
| `src/pages/CapabilityStudioPage.tsx` | 在 P10 dry-run 后显示 P11 表单和按钮“检查 lifecycle preflight” |
| `src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 blocked 用例：candidate bundle 不是 delivery |
| `src/api/__tests__/executionLiveSmokeScript.test.ts` / `scripts/smoke-frontend-behavior.cjs` | live smoke 增加 endpoint 断言 |
| `docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P11 影响 |

前端文案必须明确：

- “候选 bundle 不是正式 delivery”。
- “此检查不安装、不启用、不运行 provider”。
- “下一步需要 delivery、黄金用例、evaluation ready、provider dependency、Owner/Base Gate、Receipt”。

## 7. 验收口径

后端：

```sh
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

前端：

```sh
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

本仓：

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
git diff --check
```

## 8. 授权边界

P11 已按 Owner 2026-07-04 明确授权完成以下工作树的修改和验证：

- `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss`
- `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui`

本切片未改 `truzhen-contracts`、`truzhen-software` 或 `truzhen-cloud`；未运行真实 Codex CLI；未执行第三方 OSS；未登录或上传社媒。P11 授权不自动延伸到 P12。

历史授权卡见：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-cross-repo-execution-authorization-20260704.md`
