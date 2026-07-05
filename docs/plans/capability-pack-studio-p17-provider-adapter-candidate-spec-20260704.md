# P17 Provider / Adapter Candidate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为短视频运营能力 Pack 生成 provider / adapter candidate 的归属、目录、readiness check 和验收边界，确保真实执行实现不进入 `truzhen-packs`。

**Architecture:** P17 把 P16 的 PatchCandidate 转成 provider / adapter candidate，但默认只生成 scaffold 和 readiness check，不安装、不运行第三方 OSS。真实 provider 归 `truzhen-software` 或外部 provider 仓；`truzhen-packs` 只保留 ProviderRequirement 和证据索引。

**Tech Stack:** `truzhen-software` 或外部 provider 仓、`truzhenos` readiness / Gateway、`truzhen-client-web-desktop` GUI readiness 展示、`truzhen-packs` capability descriptor。

---

## 1. 派活卡

| 维度 | P17 裁定 |
|---|---|
| 版本 / 优先级 | P17 是 provider 归属切片；未获 Owner 单独授权前只在本仓记录规格。 |
| 真实客户 / 场景证据 | 仍缺真实短视频运营客户证据。P17 只证明 adapter 候选归属清楚，不证明 OSS 可生产使用。 |
| 最小可交付 | 每个短视频能力候选有一个 adapter candidate 归属、readiness check、禁入凭据规则和 `provider_missing / blocked` 行为。 |
| 真相源 | provider 实现归 `truzhen-software` 或外部 provider；readiness 归 `truzhenos`；Pack 声明归本仓。 |
| 仓库 / 层归属 | 授权后可能改 `truzhen-software`、`truzhenos`、client；本仓只记录 descriptor 和 evidence ledger。 |
| 风险颜色 | 橙：provider / adapter scaffold；红：安装运行第三方 OSS、真实社媒登录、上传、生产凭据。 |
| 契约影响 | 默认不改 contracts。若 ProviderRequirement schema 不够表达 adapter candidate，先出 contracts 影响清单。 |
| 禁止边界 | 不把 adapter 放进 `truzhen-packs`；不执行 OSS；不保存 raw secret；不登录社媒；不上传。 |

## 2. 三类 adapter candidate

| 候选 | 参考 OSS | adapter candidate | 默认 readiness |
|---|---|---|---|
| 草稿生成 | MoneyPrinterTurbo | `provider-candidate://short-video-draft-generation-adapter` | `provider_missing`，只可 scaffold。 |
| 合成编排 | Pixelle-Video | `provider-candidate://short-video-composition-adapter` | `provider_missing`，模型 / 媒体执行 blocked。 |
| 发布草稿 | social-auto-upload | `provider-candidate://short-video-social-publish-adapter` | `blocked_by_external_send_risk`，不登录、不上传。 |

## 3. 目标文件

### `truzhen-software` 或外部 provider 仓

授权后才允许创建：

```text
providers/short-video-draft-generation-adapter/
providers/short-video-composition-adapter/
providers/short-video-social-publish-adapter/
```

每个目录至少包含：

```text
README.md
provider.manifest.json
readiness_check.py 或 readiness_check.ts
tests/
docs/risk-boundary.md
```

### `truzhenos`

| 文件 | 预期变更 |
|---|---|
| `modules/04-capability-management/acceptance.md` | 登记 provider candidate readiness 口径。 |
| `backend/tests/capability/...` | 如果 readiness 汇入 04，则补 provider_missing / blocked 测试。 |

### `truzhen-packs`

| 文件 | 预期变更 |
|---|---|
| `capability-pack-candidates/short-video-ops-v0/docs/provider-adapter-candidate-ledger.md` | P17 provider / adapter candidate 台账。 |
| `capability-pack-candidates/short-video-ops-v0/candidate-set.json` | 登记 P17 状态和候选 refs。 |

## 4. Provider manifest 最小字段

每个 adapter candidate 的 `provider.manifest.json` 至少包含：

```json
{
  "provider_candidate_ref": "provider-candidate://short-video-draft-generation-adapter",
  "capability_refs": ["capability-pack://short-video-draft-generation"],
  "source_evidence_refs": ["github://..."],
  "risk_class": "orange",
  "default_readiness": "provider_missing",
  "credential_policy": "secret_ref_only_no_raw_secret",
  "execution_policy": {
    "no_third_party_oss_execution_by_default": true,
    "owner_base_gate_required": true,
    "receipt_required": true
  }
}
```

`social-auto-upload` adapter 必须额外声明：

```json
{
  "external_send_risk": true,
  "social_login_blocked_by_default": true,
  "upload_blocked_by_default": true
}
```

## 5. 执行任务

### Task 1: 写 provider candidate 台账

**Files:**
- Create: `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/provider-adapter-candidate-ledger.md`

- [ ] **Step 1: Add three candidate rows**

每行必须包含：

```text
provider_candidate_ref
target_repository
capability_ref
source_oss
readiness_default
risk_class
credential_policy
forbidden_actions
```

### Task 2: 授权后创建 adapter scaffold

**Files:**
- Create under authorized `truzhen-software` worktree only.

- [ ] **Step 1: Create scaffold for draft generation**

Create:

```text
providers/short-video-draft-generation-adapter/README.md
providers/short-video-draft-generation-adapter/provider.manifest.json
providers/short-video-draft-generation-adapter/docs/risk-boundary.md
```

The README must state:

```text
This adapter candidate does not vendor or run MoneyPrinterTurbo.
Readiness stays provider_missing until Owner authorizes installation and tests.
```

- [ ] **Step 2: Create scaffold for composition**

Same structure, with Pixelle-Video and media/model execution blocked.

- [ ] **Step 3: Create scaffold for social publish**

Same structure, with social login and upload blocked by default.

### Task 3: readiness check

**Files:**
- Create readiness check file in each provider candidate directory.

- [ ] **Step 1: Return safe default**

Readiness check must return:

```json
{
  "status": "provider_missing",
  "reason": "third_party_runtime_not_installed_or_authorized",
  "safe_to_execute": false
}
```

For social publish:

```json
{
  "status": "blocked",
  "reason": "external_send_requires_separate_owner_authorization",
  "safe_to_execute": false
}
```

## 6. 验收命令

本仓：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

授权后 `truzhen-software`：

```sh
cd /Users/li/Documents/truzhen-software
git status --short --branch
git diff --check
```

如果仓库已有测试脚本，必须追加运行项目规定的 provider readiness 测试。

## 7. 完成口径

P17 完成只代表 provider / adapter candidate 归属清楚并能诚实返回 readiness。它不代表：

- 第三方 OSS 已安装或运行。
- 社媒已登录或上传。
- provider 已生产 ready。
- Capability Pack 已 enabled。
- 云市场已发布。
