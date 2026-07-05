# P13 GUI Lifecycle 面板 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让能力 Pack 制作台在短视频候选包页展示 candidate bundle、delivery、readiness、promote、confirm、enabled pointer 和 receipt ref 的分层状态，并阻止用户绕过 readiness / Gate / Receipt。

**Architecture:** P13 默认复用 P11 `lifecycle_preflight` 和 P12 安全样本 lifecycle 输出，不新增 contracts。前端负责汇总展示状态和禁用危险动作；后端只在现有响应缺少展示字段时补只读聚合，不创建真实 provider、不运行 Codex CLI。

**Tech Stack:** React/TypeScript、现有 `src/pages/CapabilityStudioPage.tsx`、`src/components/capability-studio/CapabilityStudioSurfaces.tsx`、Vitest/Testing Library、现有 smoke 脚本、Go capability HTTP 回归。

---

## 1. 派活卡

| 维度 | P13 裁定 |
|---|---|
| 版本 / 优先级 | P13 是 P12 后的 GUI 商用可用性切片；P12 未完成前，P13 只能先接 blocked / fixture 状态展示。 |
| 真实客户 / 场景证据 | Owner 目标是完善能力 Pack 制作台直到可商用；真实短视频运营客户记录仍缺，P13 只改善制作台可验收性。 |
| 最小可交付 | 用户能在一个面板看到 lifecycle 全链路状态、阻断原因、下一步、目标仓归属，并且不能点击出假启用。 |
| 真相源 | lifecycle 状态归 `truzhenos`；面板展示归 `truzhen-client-web-desktop`；候选资产和验收台账归 `truzhen-packs`。 |
| 仓库 / 层归属 | 授权后默认改 `truzhen-client-web-desktop`，必要时只读或小改 `truzhenos` 聚合响应；不改 contracts / software / cloud。 |
| 风险颜色 | 黄：GUI 状态面板和禁用按钮；橙：若需改后端 lifecycle 聚合字段；红色真实执行全部禁止。 |
| 契约影响 | 默认不改 `truzhen-contracts`。若现有响应无法稳定表达面板字段，先停工输出 contracts 影响清单。 |
| 禁止边界 | 不运行 Codex CLI；不执行第三方 OSS；不登录或上传社媒；不把 disabled / fixture / candidate-only 写成商用成功。 |
| 验收设计 | 前端单测覆盖按钮禁用和字段展示；smoke 覆盖 live 后端返回 blocked / ready 状态；本仓记录 P13 状态。 |

## 2. 目标文件

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/components/capability-studio/CapabilityStudioSurfaces.tsx` | 新增或扩展 `LifecyclePanel` / `LifecycleTimeline`，展示分层状态、issue、下一步和目标仓。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx` | 将 P11/P12 输出接入面板，启用按钮只在 readiness、Gate、Receipt 条件满足时可用。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 P13 GUI 单测。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/api/__tests__/executionLiveSmokeScript.test.ts` | 增加 lifecycle panel smoke 断言。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs` | 增加面板行为检查。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P13 前端影响。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/FEATURE_LEDGER.md` | 更新制作台商用缺口口径。 |

### `truzhenos`

仅当现有 P11/P12 响应不能支撑面板时才改：

| 文件 | 预期变更 |
|---|---|
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/service.go` | 只读聚合 lifecycle panel 所需字段。 |
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/studio_pipeline_http_test.go` | 增加 panel 字段 HTTP 回归。 |
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/acceptance.md` | 记录 P13 展示字段和禁用动作。 |

## 3. 面板字段

P13 面板至少展示：

```text
candidate_bundle.status
candidate_bundle.blocked_reason
delivery.present
delivery.status
readiness.status
readiness.issues[]
promote.status
confirm.status
enabled_pointer.ref
receipt_ref
next_required_action
target_repository
forbidden_action_notice
```

字段含义：

- `candidate_bundle.status`：候选包是否存在，候选仍不等于 delivery。
- `delivery.present`：是否已生成服务端正式 delivery。
- `readiness.status`：`ready / provider_missing / not_ready / blocked`。
- `promote.status`：是否可进入候选 promote；不能代表 enabled。
- `confirm.status`：是否已有 Owner/Base Gate 和 Receipt。
- `enabled_pointer.ref`：仅 confirm 成功后显示。
- `forbidden_action_notice`：提醒不运行第三方 OSS、不运行 Codex CLI、不登录社媒。

## 4. 执行任务

### Task 1: 写 GUI 面板失败测试

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/__tests__/capabilityStudioWizard.test.tsx`

- [ ] **Step 1: Add lifecycle panel blocked test**

测试名称：

```text
renders lifecycle panel and blocks enable before readiness gate receipt
```

断言：

```text
shows candidate bundle status
shows candidate_bundle_not_delivery
shows provider_missing or not_ready issue
shows target repository for each issue
enable button is disabled
no text says enabled successfully
```

- [ ] **Step 2: Run failing test**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx -t "lifecycle panel"
```

Expected before implementation: FAIL because the lifecycle panel assertions are absent.

### Task 2: 实现面板组件

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/components/capability-studio/CapabilityStudioSurfaces.tsx`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx`

- [ ] **Step 1: Add `LifecyclePanel` props shape**

使用本地类型即可，不改 contracts：

```ts
type LifecyclePanelIssue = {
  key: string
  status: 'ready' | 'provider_missing' | 'not_ready' | 'blocked'
  targetRepository: string
  nextAction: string
}

type LifecyclePanelState = {
  candidateBundleStatus: string
  deliveryStatus: string
  readinessStatus: string
  promoteStatus: string
  confirmStatus: string
  enabledPointerRef?: string
  receiptRef?: string
  blockedReason?: string
  issues: LifecyclePanelIssue[]
  canEnable: boolean
}
```

- [ ] **Step 2: Render status sections**

必须渲染这些稳定文本或等价 `aria-label`：

```text
候选包
正式 delivery
Readiness
Promote
Confirm
Enabled pointer
Receipt
下一步
```

- [ ] **Step 3: Disable enable action**

规则：

```text
canEnable=false when readinessStatus != ready
canEnable=false when confirmStatus does not include gate_and_receipt_confirmed
canEnable=false when receiptRef is empty
canEnable=false when enabledPointerRef is empty
```

### Task 3: 增加 smoke 覆盖

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/api/__tests__/executionLiveSmokeScript.test.ts`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs`

- [ ] **Step 1: Assert blocked panel state**

Smoke 必须断言：

```text
lifecycle panel visible
candidate_bundle_not_delivery visible
enable action disabled
provider_missing or not_ready visible
no request to third-party OSS execution endpoint
no Codex CLI run request beyond candidate request unless 11 Gateway receipt exists
```

- [ ] **Step 2: Run frontend verification**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

Expected: tests pass; smoke proves blocked state is visible and enable action cannot bypass Gate / Receipt.

## 5. 本仓收尾

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/candidate-set.json`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/README.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/MODULES.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/FEATURE_LEDGER.md`

- [ ] **Step 1: Update status only after verification**

After P13 implementation passes frontend and any needed backend tests, set:

```json
"p13_gui_lifecycle_panel": {
  "status": "gui_lifecycle_panel_connected"
}
```

Do not change this status before verification.

- [ ] **Step 2: Run packs verification**

Run:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

Expected: JSON 合法、脚本语法 OK、无禁入产物、diff 无空白错误。

## 6. 完成口径

P13 完成只代表 GUI 能诚实展示 lifecycle 状态并阻止假启用。它不代表：

- P12 安全样本已完成。
- 短视频能力 Pack 已正式 enabled。
- 第三方 OSS 已运行。
- Codex CLI 已真实执行。
- 社媒上传已接通。
- 云市场已发布。
