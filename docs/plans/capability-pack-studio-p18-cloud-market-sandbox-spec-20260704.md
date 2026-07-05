# P18 云市场 Sandbox 链路 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 `truzhen-cloud` sandbox 中验证能力 Pack 候选从商品草稿、审核、License / Entitlement、下载分发到安装入口的最小链路，不触碰生产支付或真实发布。

**Architecture:** `truzhen-packs` 只提供候选资产和展示元数据；正式商品、订单、支付、License / Entitlement、Release 和下载分发真相归 `truzhen-cloud`。P18 使用 sandbox 商品和测试 entitlement，不处理真实钱款、不发布生产市场。

**Tech Stack:** `truzhen-cloud` sandbox、`truzhen-client-web-desktop` 市场 / 安装入口、`truzhenos` pack lifecycle / entitlement check、`truzhen-packs` candidate artifact manifest。

---

## 1. 派活卡

| 维度 | P18 裁定 |
|---|---|
| 版本 / 优先级 | P18 是商用链路切片；P12-P17 未形成证据前，只能做 sandbox，不做生产发布。 |
| 真实客户 / 场景证据 | 仍缺真实购买客户记录。P18 只验证云市场技术链路，不证明市场需求。 |
| 最小可交付 | 一个 sandbox listing 能引用短视频能力候选包，生成测试 License / Entitlement，客户端下载并进入安装预检。 |
| 真相源 | 商品、订单、License / Entitlement、Release、下载分发归 `truzhen-cloud`；安装和 enabled 归 `truzhenos`；Pack 声明归本仓。 |
| 仓库 / 层归属 | 授权后可能改 `truzhen-cloud`、client、`truzhenos`；本仓只登记 artifact manifest 和云边界。 |
| 风险颜色 | 橙：sandbox 商品 / entitlement / download。红：真实支付、生产上架、真实 License、生产云发布。 |
| 契约影响 | 默认不改 contracts。若 listing / entitlement DTO 需跨仓稳定，先出 contracts 影响清单。 |
| 禁止边界 | 不生产发布、不真实扣款、不生成生产 License、不上传 raw secret、不把本仓字段当云真相。 |

## 2. Sandbox 链路

P18 最小链路：

```text
truzhen-packs candidate artifact manifest
-> truzhen-cloud sandbox PackListingDraft
-> sandbox review approved
-> sandbox License / Entitlement issued
-> client download/install preflight
-> truzhenos lifecycle / entitlement check
-> result receipt or blocked reason
```

## 3. 目标文件

### `truzhen-cloud`

授权后才允许修改：

```text
apps/market/
apps/entitlement/
apps/download/
tests/sandbox/
```

预期新增 / 修改：

```text
PackListingDraft sandbox fixture
License / Entitlement sandbox fixture
download release fixture
market review sandbox test
entitlement check test
```

### `truzhen-client-web-desktop`

预期新增 / 修改：

```text
market sandbox listing display
download/install preflight UI
entitlement blocked / ready display
```

### `truzhenos`

预期新增 / 修改：

```text
entitlement check before lifecycle install / enable
sandbox receipt evidence
blocked reason if entitlement missing
```

### `truzhen-packs`

预期新增 / 修改：

```text
capability-pack-candidates/short-video-ops-v0/docs/cloud-market-sandbox-ledger.md
candidate-set.json p18 status
```

## 4. Sandbox Listing 字段

`truzhen-cloud` 的 sandbox listing 至少需要：

```json
{
  "listing_ref": "sandbox-pack-listing://short-video-ops-capability-candidates",
  "artifact_ref": "capability-pack-candidate-set://short-video-ops-v0",
  "seller_ref": "seller://sandbox-owner",
  "review_status": "sandbox_approved",
  "license_policy": "sandbox_entitlement_required",
  "download_policy": "sandbox_signed_url",
  "production_release": false
}
```

本仓不得把这些字段写成正式云真相，只能作为候选或 evidence ref。

## 5. 执行任务

### Task 1: 写 cloud sandbox 失败测试

**Files:**
- Create under authorized `truzhen-cloud` sandbox tests.

- [ ] **Step 1: Add listing draft test**

断言：

```text
candidate artifact ref is accepted as sandbox listing draft
review status starts pending or sandbox_approved fixture
production_release=false
no payment capture
no production license issued
```

- [ ] **Step 2: Add entitlement test**

断言：

```text
download requires sandbox entitlement
missing entitlement returns blocked
issued sandbox entitlement returns ready_for_download
```

### Task 2: 客户端安装预检

**Files:**
- Modify authorized client market / install UI files.

- [ ] **Step 1: Display sandbox listing**

GUI 必须显示：

```text
sandbox badge
artifact ref
entitlement status
download preflight
install preflight
```

- [ ] **Step 2: Block missing entitlement**

没有 entitlement 时不能下载 / 安装。

### Task 3: truzhenos entitlement check

**Files:**
- Modify authorized lifecycle / entitlement boundary files.

- [ ] **Step 1: Require entitlement before lifecycle**

规则：

```text
missing entitlement -> blocked
sandbox entitlement -> install preflight only
production entitlement -> not covered by P18
```

## 6. 验收命令

`truzhen-cloud`：

```sh
cd /Users/li/Documents/truzhen-cloud
git status --short --branch
git diff --check
```

再运行该仓已有 sandbox / market / entitlement 测试命令。

`truzhen-client-web-desktop`，如改动：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
git diff --check
```

`truzhenos`，如改动：

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
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

## 7. 完成口径

P18 完成只代表 sandbox 云市场链路可验证。它不代表：

- 生产市场已发布。
- 真实支付已接通。
- 生产 License / Entitlement 已签发。
- 第三方 OSS provider 已 ready。
- 短视频能力 Pack 已生产 enabled。
