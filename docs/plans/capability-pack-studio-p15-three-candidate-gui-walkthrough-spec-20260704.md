# P15 三候选 GUI 实操验收 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 用真实 GUI 路径跑完 3 个短视频能力 Pack 候选的制作台操作，并产出截图、行为日志、网络响应摘要和 issue ledger，证明制作台不是只靠 API / 文档推进。

**Architecture:** P15 是验收切片，不接真实 provider，不运行第三方 OSS，不运行真实 Codex CLI。它使用现有前端、隔离 `truzhenos` devserver 和浏览器自动化记录用户视角证据；所有输出作为候选证据写回 `truzhen-packs` 过程目录。

**Tech Stack:** `truzhen-client-web-desktop` 前端、`truzhenos` devserver、Playwright / in-app browser smoke、Markdown evidence ledger、现有 P3-P13 API。

---

## 1. 派活卡

| 维度 | P15 裁定 |
|---|---|
| 版本 / 优先级 | P15 是 P13 后的用户视角验收切片；没有 P13 面板时可先跑候选导出和阻断证据，但不能声称商用完成。 |
| 真实客户 / 场景证据 | 真实短视频运营客户记录仍缺。P15 只证明制作台操作链路，不证明业务 ROI 或 provider 可用。 |
| 最小可交付 | 3 个候选包各有 GUI 操作路径、截图、网络响应摘要、候选 artifact ref、阻断原因和 issue ledger。 |
| 真相源 | GUI 行为归 client；后端候选与 receipt candidate 归 `truzhenos`；证据台账归 `truzhen-packs`。 |
| 仓库 / 层归属 | 授权后测试 `truzhen-client-web-desktop` 和 `truzhenos`；本仓只新增 evidence ledger，不改 contracts / software / cloud。 |
| 风险颜色 | 黄：GUI 自动化和候选导出。橙：如启动隔离 devserver。红：真实 Codex run、第三方 OSS 执行、社媒登录 / 上传，均不做。 |
| 禁止边界 | 不登录任何社媒；不上传文件；不读取浏览器 cookie；不执行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload；不保存 raw token。 |
| 验收设计 | 每个候选包必须有截图、操作日志、网络响应摘要和最终 blocked / candidate-only 判断。 |

## 2. 三个 GUI 剧本

| 剧本 | 候选包 | GitHub 样本 | 必须证明 |
|---|---|---|---|
| P15-A | `capability-pack://short-video-draft-generation` | MoneyPrinterTurbo | 可查 OSS 证据，可生成草稿生成能力候选，可导出 bundle，provider 未接通时 blocked。 |
| P15-B | `capability-pack://short-video-composition-orchestration` | Pixelle-Video | 可查 license / 模型风险，可生成合成编排能力候选，可导出 bundle，媒体执行仍 blocked。 |
| P15-C | `capability-pack://short-video-social-publish-draft` | social-auto-upload | 可识别发布 / 上传红色动作，只生成发布草稿候选，不登录、不上传、不发送。 |

## 3. 目标文件

### `truzhen-packs`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-evidence-ledger.md` | 新增 P15 GUI 实操证据台账。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/candidate-set.json` | 登记 P15 状态和证据路径。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/README.md` | 登记 P15 证据资产。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/MODULES.md` | 登记 P15 成熟度。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/FEATURE_LEDGER.md` | 登记 P15 状态。 |

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs` | 可选：增加 P15 GUI walkthrough 自动化脚本。 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/docs/frontend-adjustment-impact-ledger-20260704.md` | 若脚本变更，登记 P15 影响。 |

## 4. 证据格式

P15 每个剧本在 evidence ledger 中必须记录：

```text
scenario_id
candidate_pack_ref
operator_ref
timestamp
frontend_url
backend_base_url
github_repo_evidence_ref
clicked_steps[]
network_responses[]
artifact_refs[]
screenshots[]
final_status
blocked_reasons[]
forbidden_actions_observed_false[]
issue_refs[]
```

`forbidden_actions_observed_false[]` 必须至少包含：

```text
no_codex_cli_run
no_third_party_oss_execution
no_social_login
no_social_upload
no_raw_secret_seen
no_formal_enable
```

## 5. 执行任务

### Task 1: 准备 evidence ledger

**Files:**
- Create: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-evidence-ledger.md`

- [ ] **Step 1: Create ledger skeleton**

每个剧本使用以下章节：

```markdown
## P15-A draft generation

| 字段 | 证据 |
|---|---|
| candidate_pack_ref | `capability-pack://short-video-draft-generation` |
| github_repo_evidence_ref | MoneyPrinterTurbo GitHub evidence |
| final_status | `pending_gui_walkthrough` |
| forbidden_actions | no Codex CLI run; no third-party OSS execution; no social login/upload |
```

重复 P15-B 和 P15-C，不合并证据。

### Task 2: 跑 GUI 剧本

**Files:**
- Read: `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs`
- Read: `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss/modules/04-capability-management/api_contract.md`

- [ ] **Step 1: Start isolated backend and frontend**

Run commands must be recorded in the ledger. Use existing project scripts; if server startup fails, record failure and stop.

- [ ] **Step 2: Execute P15-A**

Required GUI steps:

```text
open capability studio
select short-video draft generation candidate
view MoneyPrinterTurbo evidence
generate Code Assistant invocation candidate
generate run request candidate
export candidate bundle
run candidate bundle dry-run
run lifecycle preflight
record final blocked / candidate-only state
```

- [ ] **Step 3: Execute P15-B**

Required GUI steps:

```text
open capability studio
select short-video composition orchestration candidate
view Pixelle-Video evidence
confirm media/model provider is provider_missing or blocked
export candidate bundle
run dry-run and preflight
record blocked state
```

- [ ] **Step 4: Execute P15-C**

Required GUI steps:

```text
open capability studio
select short-video social publish draft candidate
view social-auto-upload evidence
verify upload / login action is blocked-by-default
export publish draft candidate
run dry-run and preflight
record no social login and no upload
```

### Task 3: 更新证据台账

**Files:**
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-evidence-ledger.md`

- [ ] **Step 1: Fill evidence**

For each scenario, ledger must include:

```text
at least one screenshot path
at least one network response summary
artifact refs for invocation/run request/bundle/dry-run/preflight
blocked reason
issue ledger refs
forbidden action checks
```

- [ ] **Step 2: Refuse completion if evidence is weak**

If screenshot, network response or artifact ref is missing, final status remains:

```text
walkthrough_incomplete
```

Do not mark P15 as passed.

## 6. 验收命令

Frontend / behavior:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

Backend, if devserver or HTTP behavior is changed:

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

Packs:

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

## 7. 完成口径

P15 完成只代表 GUI 实操路径有证据。它不代表：

- provider 已接通。
- 第三方 OSS 已运行。
- Codex CLI 已真实执行。
- 社媒账号已登录或上传。
- Pack 已正式 enabled。
- 云市场已发布。
