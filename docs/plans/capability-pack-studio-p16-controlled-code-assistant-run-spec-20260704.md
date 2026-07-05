# P16 受控 Code Assistant 最小 Run Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 通过 11 Execution Gateway 做一次受控 Code Assistant 最小 run，只在隔离过程目录生成 PatchCandidate 文件，证明制作台能从 run request candidate 进入真实回执链。

**Architecture:** P16 不由 04 Capability Studio 直接运行 Codex CLI，而是从 P8 `code_assistant_run_request_candidate` 进入 11 Gateway、Owner/Base Gate、Execution Receipt。输出只能是 PatchCandidate / README / adapter scaffold 候选，不应用补丁、不安装 provider、不运行第三方 OSS。

**Tech Stack:** `truzhenos` 11 Execution Gateway、01 Gate、03 Receipt、04 Capability Studio、`truzhen-client-web-desktop` GUI 状态、`truzhen-packs` evidence ledger。

---

## 1. 派活卡

| 维度 | P16 裁定 |
|---|---|
| 版本 / 优先级 | P16 是 P12/P13/P15 后的真实执行链切片；未授权前保持 `specified_pending_authorization`。 |
| 真实客户 / 场景证据 | 仍缺真实短视频运营客户记录。P16 只验证 Code Assistant 受控执行链，不证明短视频 provider 商用。 |
| 最小可交付 | 一个最小 run 由 11 Gateway 执行，产出 PatchCandidate 文件、run receipt、stdout/stderr 摘要、禁止动作证明。 |
| 真相源 | run 决策归 01 Gate；执行归 11 Gateway；回执归 03 Receipt；候选资产归 `truzhen-packs`。 |
| 仓库 / 层归属 | 授权后目标为 `truzhenos` + `truzhen-client-web-desktop` + 本仓 evidence ledger；不改 contracts / software / cloud。 |
| 风险颜色 | 橙：真实 Code Assistant run 可能消耗 Owner 本机模型额度。红：若访问外部网络、读取 secret、执行第三方 repo，必须禁止。 |
| 契约影响 | 默认不改 `truzhen-contracts`。若 run receipt 或 PatchCandidate shape 需跨仓稳定，先停工出 contracts 影响清单。 |
| 禁止边界 | 不安装、不运行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload；不读取 raw token；不登录社媒；不上传；不应用补丁。 |
| 验收设计 | 先写 11 Gateway / 04 Studio 集成测试，再跑最小真实 run；失败时必须记录 `provider_missing / not_ready / blocked`。 |

## 2. 目标文件

### `truzhenos`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/capability_studio_code_assistant_run_test.go` | 新增 P16 受控 run 测试。 |
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/service.go` | 仅在需要时把 11 run receipt 回填为 PatchCandidate intake。 |
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/api_contract.md` | 登记 P16 run receipt 回填口径。 |
| `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/acceptance.md` | 登记 P16 验收条款。 |

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx` | 显示 11 run receipt / blocked 状态。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 P16 GUI 测试。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs` | 增加 run receipt 或 blocked 断言。 |

### `truzhen-packs`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/code-assistant-controlled-run-ledger.md` | P16 受控 run 证据台账。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/candidate-set.json` | 登记 P16 状态和证据路径。 |

## 3. Run 输入限制

P16 最小 run 的 prompt 必须限制为：

```text
读取 capability-pack-candidates/short-video-ops-v0 下的候选声明和 README。
生成一个 adapter scaffold 说明或 README patch candidate。
禁止安装依赖。
禁止运行第三方 repo。
禁止访问社媒。
禁止读取 raw token、cookie、ssh key、浏览器登录态。
禁止写入本仓正式 pack 目录以外的位置。
输出 PatchCandidate 文件和 summary。
```

允许写入目录只能是 11 Gateway 分配的隔离过程目录，例如：

```text
process-artifacts/code-assistant-runs/<run_id>/
```

## 4. 执行任务

### Task 1: 写受控 run 失败测试

**Files:**
- Create: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/tests/capability/capability_studio_code_assistant_run_test.go`

- [ ] **Step 1: Add test `TestCapabilityStudioControlledCodeAssistantRunRequiresGateReceipt`**

断言：

```text
run request candidate exists
DecisionRef is required before execution
RunID and Nonce are issued by Base / 11 Gateway
ReceiptRef is required after execution
PatchCandidateRefs are captured
ReadyForInstall=false
FormalWrite=false
```

- [ ] **Step 2: Run failing test**

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityStudioControlledCodeAssistantRun' -count=1
```

Expected before implementation: FAIL because P16 controlled run fixture is absent.

### Task 2: 接入 run receipt 回填

**Files:**
- Modify: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/service.go`
- Modify if needed: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/backend/internal/capability/studio/http.go`

- [ ] **Step 1: Accept only 11 Gateway receipt**

规则：

```text
reject client-supplied fake receipt_ref
reject run without DecisionRef
reject run without Nonce
reject run output outside isolated artifact directory
accept only PatchCandidate / summary artifact kinds
```

- [ ] **Step 2: Keep output candidate-only**

规则：

```text
no_auto_apply=true
ready_for_install=false
enable_supported=false
formal_write=false
review_required=true
```

### Task 3: 前端展示 run 状态

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/CapabilityStudioPage.tsx`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/src/pages/__tests__/capabilityStudioWizard.test.tsx`

- [ ] **Step 1: Show controlled run state**

GUI 必须展示：

```text
Gate pending
DecisionRef
RunID
ReceiptRef
PatchCandidateRefs
No auto apply
Review required
```

- [ ] **Step 2: Disable run button unless gate receipt exists**

如果 01 Gate / 03 Receipt / 11 readiness 缺失，按钮必须显示 `blocked / not_ready`。

## 5. 验收命令

后端：

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

前端：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

本仓：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

## 6. 完成口径

P16 完成只代表一次受控 Code Assistant run 可产生 PatchCandidate 和 Receipt。它不代表：

- 补丁已应用。
- provider 已接通。
- 第三方 OSS 已运行。
- 短视频能力 Pack 已 enabled。
- 社媒发布已接通。
- 云市场已发布。
