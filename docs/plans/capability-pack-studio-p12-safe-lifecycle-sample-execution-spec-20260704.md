# P12 安全内置能力 lifecycle 样本 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 用一个基座内置安全能力样本或 fixture 验证 Capability Pack 从 studio delivery 到 lifecycle draft / readiness / promote / confirm 的最小闭环。

**Architecture:** P12 不接短视频第三方 OSS，也不运行真实 Codex CLI。后端只使用 `truzhenos` 已有 04 lifecycle、01 Gate 测试链和 03 Receipt 测试链；前端只把已有 lifecycle 状态、阻断原因、确认卡和 receipt ref 展示给制作台用户。

**Tech Stack:** Go 后端测试与 HTTP handler、React/TypeScript 前端、Vitest/Testing Library、现有 npm smoke、`truzhen-packs` 候选资产台账。

---

## 1. 派活卡

| 维度 | P12 裁定 |
|---|---|
| 版本 / 优先级 | 当前是 P11 后的商用化下一切片，生命周期档位为 `设计中`，待 Owner 授权后进入 `契约已定 -> 已接线`。 |
| 真实客户 / 场景证据 | Owner 要求“完善能力 Pack 制作台，直到能够商用”。真实短视频运营客户记录仍缺，P12 只验证制作台 lifecycle 底座，不声称短视频 provider 已商用。 |
| 最小可交付 | 一个安全样本能力从服务端 studio delivery 生成 draft，readiness 可解释，promote 不绕过门禁，confirm 经 01 Gate + 03 Receipt 测试链生成 enabled pointer。 |
| 要做的事 | 补后端测试和必要实现，补前端 lifecycle 面板 / 确认卡测试，让制作台能证明“安全样本可走通 lifecycle，短视频 OSS 仍 blocked”。 |
| 真相源 | lifecycle 状态、enabled pointer、Gate、Receipt 归 `truzhenos`；GUI 展示归 `truzhen-client-web-desktop`；P12 计划和候选索引归本仓。 |
| 仓库 / 层归属 | 授权后只改 `truzhenos` 与 `truzhen-client-web-desktop`；本仓只更新文档和候选台账。 |
| 风险颜色 | 黄：GUI 展示与普通 lifecycle 流程。橙：promote / confirm 语义和 01/03 测试链。红色真实动作全部禁止。 |
| 契约影响 | 默认不改 `truzhen-contracts`。若现有 DTO 不能表达 P12，必须先停工输出 contracts 影响清单。 |
| 上下文维度 | 仅允许读取 04 capability lifecycle、01 Gate 测试 fixture、03 Receipt 测试 fixture、前端制作台和 pack lifecycle 组件。 |
| 禁止边界 | 不运行 Codex CLI；不执行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload；不登录或上传社媒；不改 software / cloud；不保存 raw token。 |
| 验收设计 | 由后端单测证明 lifecycle 语义，由前端单测和 smoke 证明 GUI 不能绕过 readiness / Gate / Receipt，由本仓 JSON / Markdown 验证证明计划资产一致。 |
| 验收维度 | 改了 lifecycle 就跑 lifecycle 后端测试；改了 GUI 就跑前端定向测试、typecheck 和 smoke；改了本仓候选索引就跑 JSON 合法性和结构审计。 |
| 变更影响 | 影响 04 lifecycle 最小样本、制作台 lifecycle 面板、确认卡和候选资产台账；不影响 contracts、provider、云市场、真实执行。 |

## 2. 目标仓库与文件

| 仓库 | 路径 | P12 职责 |
|---|---|---|
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | lifecycle draft / readiness / promote / confirm 的安全样本后端测试与必要实现。 |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | 制作台展示安全样本 lifecycle 状态、阻断项、确认卡、enabled pointer 和 receipt ref。 |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | 本规格、授权卡、候选集 P12 状态和验证登记。 |

默认不改：

- `/Users/li/Documents/truzhen-contracts`
- `/Users/li/Documents/truzhen-software`
- `/Users/li/Documents/truzhen-cloud`

## 3. 复用接口

P12 优先复用现有 lifecycle 端点，不新增跨仓 schema：

```text
POST /v3/capability/lifecycle/draft
POST /v3/capability/lifecycle/readiness
POST /v3/capability/lifecycle/promote
POST /v3/capability/lifecycle/confirm
GET  /v3/capability/lifecycle/history?pack_ref=
GET  /v3/capability/lifecycle/packs?pack_ref=
```

P11 的 `/v3/capability/studio/lifecycle_preflight` 仍保留为短视频 candidate-only 阻断入口，不升级为正式 delivery。

## 4. 后端执行任务

### Task 1: 写安全样本 lifecycle 失败测试

**Files:**
- Create: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/capability_studio_safe_lifecycle_sample_test.go`
- Read: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/capability_lifecycle_chain_test.go`
- Read: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/lifecycle_spec.go`

- [ ] **Step 1: Add failing tests**

Test names and required assertions:

```text
TestCapabilityPackLifecycleSafeFixtureDraftIsServerDerived
- arrange a safe built-in or fixture studio delivery
- call lifecycle draft through the existing service or HTTP route
- assert draft pack_ref and version come from server delivery, not client-supplied fake spec
- assert no third_party_oss_ref and no code_assistant_run_ref are present
- assert no enabled pointer is written at draft stage

TestCapabilityPackLifecycleSafeFixtureReadinessExplainsProviderMissing
- arrange the same fixture with provider readiness forced to provider_missing
- call readiness
- assert status is blocked or not_ready
- assert issue key includes provider_missing
- assert promote is not allowed

TestCapabilityPackLifecycleSafeFixtureConfirmRequiresGateAndReceipt
- arrange a ready safe fixture using existing test Gate / Receipt helpers
- call draft, readiness, promote, confirm in order
- assert confirm requires Owner/Base decision input
- assert result contains receipt_ref from 03 test receipt chain
- assert lifecycle packs query returns enabled pointer only after confirm
```

- [ ] **Step 2: Verify failure**

Run:

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityPackLifecycleSafeFixture' -count=1
```

Expected before implementation: tests fail because the safe lifecycle fixture or confirm wiring is absent.

### Task 2: Implement minimal safe lifecycle fixture

**Files:**
- Create or modify: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/lifecycle_safe_sample.go`
- Modify only if needed: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/service.go`
- Modify only if needed: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/lifecycle_spec.go`

- [ ] **Step 1: Add a named safe fixture**

Required fixture properties:

```text
pack_ref: capability-pack://safe-builtin-lifecycle-sample
version: 0.1.0
title: 安全内置 lifecycle 样本
risk: low
source: server_side_fixture
third_party_oss_ref: empty
code_assistant_run_ref: empty
social_platform_ref: empty
golden_case_count: at least 1
provider_dependency: safe_builtin_fixture
```

- [ ] **Step 2: Ensure service refuses client-forged spec**

Acceptance:

```text
client may send session_id and idempotency_key
server derives PackRefValue, Version, Title, Risk and GoldenCaseCount
client-supplied pack_ref or version does not override server fixture
```

- [ ] **Step 3: Re-run backend targeted tests**

Run:

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityPackLifecycleSafeFixture' -count=1
```

Expected: targeted tests pass.

### Task 3: Add HTTP / devserver regression

**Files:**
- Modify: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/studio_pipeline_http_test.go`
- Modify: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/api_contract.md`
- Modify: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/acceptance.md`

- [ ] **Step 1: Add HTTP roundtrip**

Required route coverage:

```text
POST /v3/capability/lifecycle/draft
POST /v3/capability/lifecycle/readiness
POST /v3/capability/lifecycle/promote
POST /v3/capability/lifecycle/confirm
GET  /v3/capability/lifecycle/history?pack_ref=capability-pack://safe-builtin-lifecycle-sample
GET  /v3/capability/lifecycle/packs?pack_ref=capability-pack://safe-builtin-lifecycle-sample
```

- [ ] **Step 2: Run full backend capability verification**

Run:

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

Expected: all tests pass; no whitespace errors.

## 5. 前端执行任务

### Task 4: Add GUI lifecycle sample test

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/__tests__/capabilityStudioWizard.test.tsx`
- Modify if needed: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/api/client.ts`

- [ ] **Step 1: Add failing GUI test**

Required assertions:

```text
renders safe lifecycle sample action after P11 preflight area
draft button calls /v3/capability/lifecycle/draft
readiness result shows provider_missing or ready with issue list
confirm action is disabled until promote / Gate / Receipt prerequisites are present
enabled pointer and receipt ref render only after confirm response
MoneyPrinterTurbo, Pixelle-Video and social-auto-upload are not presented as runnable providers
```

- [ ] **Step 2: Run test and confirm failure**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t lifecycle
```

Expected before implementation: test fails because P12 GUI action is absent.

### Task 5: Implement GUI wiring

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/components/capability-studio/CapabilityStudioSurfaces.tsx`
- Modify if needed: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/api/client.ts`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/docs/frontend-adjustment-impact-ledger-20260704.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/FEATURE_LEDGER.md`

- [ ] **Step 1: Add safe sample lifecycle controls**

Required controls:

```text
button: 生成安全样本 draft
button: 检查安全样本 readiness
button: 提交 promote 候选
confirmation card: confirm requires Owner/Base Gate and Receipt
result fields: pack_ref, version, readiness status, issue keys, enabled pointer, receipt_ref
```

- [ ] **Step 2: Preserve blocked short-video OSS messaging**

Required UI boundary:

```text
MoneyPrinterTurbo / Pixelle-Video / social-auto-upload remain evidence or blocked samples
no UI copy says third-party provider ready
no UI control runs Codex CLI
no UI control logs in or uploads social media
```

- [ ] **Step 3: Run frontend verification**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

Expected: tests pass and smoke shows safe sample lifecycle status without any real Codex CLI, third-party OSS, or social upload.

## 6. 本仓收尾任务

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/candidate-set.json`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/README.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/README.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/MODULES.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/FEATURE_LEDGER.md`

- [ ] **Step 1: Update P12 status after implementation**

After P12 is actually implemented and verified, change `p12_safe_lifecycle_sample.status` from `specified_pending_authorization` to `safe_lifecycle_sample_connected`. Do not change it before cross-repo tests pass.

- [ ] **Step 2: Run packs validation**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

Expected: JSON 合法、脚本语法 OK、无禁入产物、diff 无空白错误。

## 7. P12 完成口径

P12 完成只代表安全内置样本证明 lifecycle 底座可用。它不代表：

- 短视频三个 Capability Pack 已启用。
- MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload 已运行。
- Codex CLI 已真实执行或消耗 token。
- 社媒账号已登录、上传或发布。
- provider / adapter 已进入 `truzhen-software`。
- 云市场、License 或 Entitlement 已接线。
