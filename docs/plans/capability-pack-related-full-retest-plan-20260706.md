# 能力包相关全量重测 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to execute this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在主仓更新后，重新组织并启动能力包相关资产、制作台前端表面、候选证据门禁和禁入边界的全量重测。

**Architecture:** 本计划以 `truzhen-packs` 为计划和候选资产真相源，以 `truzhen-client-web-desktop` 为 GUI / ReadModel / Candidate / Receipt 消费方。第一轮只跑本地可复核测试和只读前端测试；任何需要 `truzhenos`、真实 Codex CLI、第三方 OSS、社媒账号或云市场的动作都保持 blocked，并写入后续授权切片。

**Tech Stack:** Go test、Python JSON / py_compile 审计、Vitest、TypeScript、Vite smoke、Pack candidate JSON、Markdown 计划和证据台账。

---

## 1. Owner 裁定记录

| 项 | 本轮记录 |
| --- | --- |
| 当前指令 | Owner 要求“主仓已经更新，按照前面的计划文件重新写计划，再次启动对能力包相关的所有内容进行测试”。 |
| 版本 / 优先级 | 按当前客户主线的“能力包制作台完善测试”执行；商用化、云市场、真实 provider、真实 Codex CLI 仍不进入本轮。 |
| 真实客户 / 场景证据 | 仍缺真实短视频运营客户原话、账号矩阵、真实日更压力和失败案例。现有证据只来自 Owner 纠偏、短视频运营过程样本和三个 OSS 种子：MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。 |
| 最小可交付 | 一份新计划文件、一轮 Pack 仓能力包候选资产验证、一轮 client 能力制作台相关只读测试、一个重测结果报告。 |
| 生命周期档位 | `设计中 -> 已接线` 之间的重测；本轮测试通过不等于 `已验收` 或 `已发布`。 |

## 2. 派活卡

| 维度 | 结论 |
| --- | --- |
| 我要做的事 | 用已经沉淀的短视频运营能力包候选集，重新测试能力包制作台相关资产、前端表面、候选 / formal 隔离、证据门禁和禁入边界。 |
| 真相源 | Pack 声明和候选证据归 `truzhen-packs`；GUI 投影和用户操作表面归 `truzhen-client-web-desktop`；Gate / Receipt / lifecycle runtime 归 `truzhenos`；schema 归 `truzhen-contracts`；provider 归 `truzhen-software` 或外部 provider；商品、License、Entitlement 归 `truzhen-cloud`。 |
| 仓库 / 层归属 | 本轮允许修改 `truzhen-packs` 的计划和测试报告；允许只读测试 `truzhen-client-web-desktop`。本轮不修改 client，除非 Owner 后续明确授权。 |
| 风险颜色 | 绿：计划、报告、JSON / Markdown / Go 测试；黄：client 能力制作台只读测试；橙：需要 `truzhenos` lifecycle / Receipt 的 live 验证；红：真实 Codex CLI、第三方 OSS 执行、社媒登录 / 上传、真实支付或生产发布。 |
| 契约影响 | 本轮不改 contracts，不新增 DTO，不改变 Candidate / Receipt / ProviderRequirement schema。若测试发现契约缺口，只写 issue 和影响清单。 |
| 禁止边界 | 不运行真实 Codex CLI；不执行 MoneyPrinterTurbo / Pixelle-Video / social-auto-upload；不登录或上传社媒；不保存 token / cookie / secret；不改 contracts、software、cloud；不把 candidate-only 写成 enabled。 |
| 用户如何验收 | 查看新计划、测试输出摘要、重测报告、失败项 issue、阻断项是否明确归属到后续授权切片。 |
| 先输出什么 | 先输出本计划文件，然后启动 T0-T2 测试；T3 以后需要 `truzhenos` 和更高风险授权时只登记 blocked。 |

## 3. 本轮涉及仓库

| 仓库 | 本轮职责 | 允许动作 | 目标文件 / 模块 | 验证命令 | 禁止边界 |
| --- | --- | --- | --- | --- | --- |
| `/Users/li/Documents/truzhen-packs` | 计划、候选资产、测试报告、能力包 candidate 静态门禁 | 修改计划和报告；运行本仓测试 | `docs/plans/`、`capability-pack-candidates/short-video-ops-v0/`、`capability_pack_candidates_test.go` | `go test ./... -count=1`、JSON 审计、py_compile、Pack 结构审计、禁入产物扫描、`git diff --check` | 不放 provider 实现、不存源码依赖、不存 secret、不声称 enabled |
| `/Users/li/Documents/truzhen-client-web-desktop` | 能力制作台、能力包管理、Code Assistant 面板、Pack lifecycle 前端表面 | 只读测试 | `src/pages/CapabilityStudioPage.tsx`、`src/components/capability-studio/`、`src/components/pack-lifecycle/`、相关 tests、`scripts/smoke-pack-01.cjs` | `npm run typecheck`、能力制作台相关 Vitest、`node scripts/smoke-pack-01.cjs` | 本轮不改 UI、不 mock 真实接线、不直连 provider |
| `truzhenos` | lifecycle / Gate / Receipt / Execution Gateway | 本轮不读写不测试，除非 Owner 后续开放 workspace 和授权 | 04 capability、11 execution、03 receipt、01 gate | 后续授权后再列 | 不绕 Base Gate、不启动真实执行 |
| `truzhen-contracts` | DTO / schema | 本轮不碰 | Candidate / Receipt / ProviderRequirement schema | 无 | 不由消费方先造 schema |
| `truzhen-software` | provider / adapter | 本轮不碰 | provider readiness / adapter candidate | 无 | 不把 provider 塞进 packs |
| `truzhen-cloud` | 商品 / License / Entitlement / 下载分发 | 本轮不碰 | cloud market sandbox | 无 | 不上架、不真实支付、不发 License |

## 4. 测试范围

| 范围 | 覆盖内容 | 本轮判定规则 |
| --- | --- | --- |
| A. Pack 候选资产 | 3 个短视频 Capability Pack 候选、candidate-set、P12-P18 evidence contract、readiness / go-no-go / goal map | JSON 合法、candidate-only、non-formal、禁入动作覆盖、后验收门引用完整 |
| B. 计划和证据台账 | 前一轮 P11-P18 计划、授权路线图、商用审计、新计划 | 路径不依赖旧 worktree；阻断项不冒充授权；缺证据标明 |
| C. client 能力制作台表面 | 能力制作台路由、copy guard、receipt / gate hint、Code Assistant pack 面板、Pack lifecycle smoke | 测试只证明 UI 表面和只读断言；不能证明真实 lifecycle enabled |
| D. 禁入扫描 | raw token、第三方源码、构建产物、日志、数据库、pyc、candidate enabled 冒充 | 已跟踪文件和本次变更均不得包含禁入产物 |
| E. 后续 blocked 切片 | P12/P13/P15/P16/P17/P18 | 缺 `truzhenos` / cloud / provider / Codex CLI 授权时必须 blocked，不跳过依赖 |

## 5. 执行步骤

### Task 1: Pack 仓隔离基线

**Files:**
- Read: `/Users/li/Documents/truzhen-packs/AGENTS.md`
- Read: `/Users/li/Documents/truzhen-packs/CLAUDE.md`
- Read: `/Users/li/Documents/truzhen-packs/README.md`
- Read: `/Users/li/Documents/truzhen-packs/MODULES.md`
- Create: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/docs/plans/capability-pack-related-full-retest-plan-20260706.md`

- [x] **Step 1: 确认主仓状态**

Run:

```bash
git -C /Users/li/Documents/truzhen-packs fetch origin
git -C /Users/li/Documents/truzhen-packs status --short --branch
```

Expected: `main...origin/main` 且无脏文件。

- [x] **Step 2: 创建隔离 worktree**

Run:

```bash
git -C /Users/li/Documents/truzhen-packs worktree add "$HOME/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706" -b codex/capability-pack-test-replan-20260706
```

Expected: 新分支基于更新后的 `main`。

- [x] **Step 3: 跑 Pack 基线**

Run:

```bash
go test ./... -count=1
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
```

Expected: Go tests pass；JSON 合法。

### Task 2: Pack 候选资产重测

**Files:**
- Read: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/capability-pack-candidates/short-video-ops-v0/README.md`
- Read: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/capability_pack_candidates_test.go`
- Report: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/docs/plans/capability-pack-related-retest-report-20260706.md`

- [x] **Step 1: 跑全量 Go 语义测试**

Run:

```bash
go test ./... -count=1
```

Expected: `ok github.com/truzhen/packs`。失败时进入系统化定位，不继续声明通过。

- [x] **Step 2: 跑 JSON 合法性**

Run:

```bash
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
```

Expected: `JSON 合法`。

- [x] **Step 3: 跑 install / uninstall 脚本语法并清理缓存**

Run:

```bash
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
find . -path ./.git -prune -o -type f -name '*.pyc' -delete
find . -path ./.git -prune -o -type d -name '__pycache__' -empty -delete
```

Expected: `脚本语法 OK`，并且没有 `*.pyc` 留下。

- [x] **Step 4: 跑 Pack 结构审计**

Run:

```bash
python3 - <<'PY'
import json, pathlib, sys
bad = []
for manifest_path in sorted(pathlib.Path('.').glob('*/manifest.json')):
    if str(manifest_path).startswith('templates/'):
        continue
    m = json.load(open(manifest_path, encoding='utf-8'))
    base = manifest_path.parent
    for key in ['flow_file', 'role_slots_file', 'capabilities_file']:
        rel = m.get(key)
        if not rel or not (base / rel).exists():
            bad.append(f'{manifest_path}: {key} missing or not found: {rel}')
    required = [
        'person_strategy', 'formalization_requirement', 'gates',
        'provider_requirements', 'notification_command_report_routes',
        'multi_role_comparison', 'moat_justification'
    ]
    miss = [k for k in required if k not in m]
    if miss:
        bad.append(f'{manifest_path}: missing governance fields {miss}')
    if m.get('knowledge_scopes') and m.get('knowledge_scopes_manifest'):
        declared = set(m.get('knowledge_scopes') or [])
        scopes = json.load(open(base / m['knowledge_scopes_manifest'], encoding='utf-8')).get('scopes', [])
        actual = {s.get('scope_ref') for s in scopes}
        if declared != actual:
            bad.append(f'{manifest_path}: knowledge scopes mismatch declared-only={sorted(declared-actual)} actual-only={sorted(actual-declared)}')
if bad:
    print('\n'.join(bad))
    sys.exit(1)
print('Pack 结构审计 OK')
PY
```

Expected: `Pack 结构审计 OK`。

- [x] **Step 5: 跑禁入产物扫描**

Run:

```bash
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' && exit 1 || echo "已跟踪禁入产物扫描 OK"
```

Expected: `已跟踪禁入产物扫描 OK`。

### Task 3: client 能力制作台只读测试

**Files:**
- Read: `/Users/li/Documents/truzhen-client-web-desktop/AGENTS.md`
- Read: `/Users/li/Documents/truzhen-client-web-desktop/CLIENT_LAYER.md`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/pages/__tests__/capabilityStudioWizard.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/pages/__tests__/capabilityStudioGateHint.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/pages/__tests__/capabilityStudioReceipt.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/pages/__tests__/capabilityStudioCopyGuard.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/pages/__tests__/capabilityPackNearCardStatus.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx`
- Test: `/Users/li/Documents/truzhen-client-web-desktop/scripts/smoke-pack-01.cjs`

- [x] **Step 1: 确认 client 仓状态**

Run:

```bash
git -C /Users/li/Documents/truzhen-client-web-desktop status --short --branch
```

Expected: 当前分支和 WIP 清楚；若存在脏文件，报告并不覆盖。

- [x] **Step 2: 跑 TypeScript 类型检查**

Run:

```bash
npm run typecheck
```

Expected: TypeScript build check exit 0。

- [x] **Step 3: 跑能力制作台目标 Vitest**

Run:

```bash
npm run test -- \
  src/pages/__tests__/capabilityStudioWizard.test.tsx \
  src/pages/__tests__/capabilityStudioGateHint.test.tsx \
  src/pages/__tests__/capabilityStudioReceipt.test.tsx \
  src/pages/__tests__/capabilityStudioCopyGuard.test.tsx \
  src/pages/__tests__/capabilityPackNearCardStatus.test.tsx \
  src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx
```

Expected: 目标测试全部 pass。该结果只证明前端只读表面和候选 UI，不证明真实 lifecycle enabled。

- [x] **Step 4: 跑 Pack Studio 静态 smoke**

Result: FAIL，`scripts/smoke-pack-01.cjs` 仍断言旧 `/v2/pack-studio/*` 和 `src/api/client.ts`，与当前 `/v3/pack-studio/*` 及 `src/api/domains/authMarket.ts` 域拆分不一致。详见 `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/docs/plans/capability-pack-related-retest-report-20260706.md#失败项`。

Run:

```bash
node scripts/smoke-pack-01.cjs
```

Expected: Pack 制作台入口、Capability Pack 草案入口、安全 fallback、禁止直连模型 API 等断言通过。

### Task 4: 写入重测报告

**Files:**
- Create: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/docs/plans/capability-pack-related-retest-report-20260706.md`
- Modify: `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/FEATURE_LEDGER.md` only if the report adds durable status that should enter the ledger.

- [x] **Step 1: 写报告**

Report must include:

```markdown
# 能力包相关全量重测报告

## 范围

## 已跑命令与结果

## 失败项

## 阻断项

## 不计入完成的证据

## 下一授权切片
```

- [x] **Step 2: 对报告跑文档检查**

Run:

```bash
git diff --check
```

Expected: exit 0。

### Task 5: 收尾验证

**Files:**
- Verify all changed files under `/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/`

- [x] **Step 1: 跑 Pack 仓最终验证**

Run:

```bash
go test ./... -count=1
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
find . -path ./.git -prune -o -type f -name '*.pyc' -delete
find . -path ./.git -prune -o -type d -name '__pycache__' -empty -delete
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' && exit 1 || echo "已跟踪禁入产物扫描 OK"
git diff --check
```

Expected: 全部 exit 0。

- [x] **Step 2: 状态检查**

Run:

```bash
git status --short --branch
```

Expected: 只显示本计划和本轮报告相关变更。

## 6. 本轮不做

- 不把 P12-P18 标成完成。
- 不运行真实 Codex CLI。
- 不执行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
- 不登录或上传抖音、小红书、Bilibili、TikTok、YouTube 等社媒。
- 不修改 `truzhen-contracts`、`truzhen-software`、`truzhen-cloud`。
- 不把 `truzhen-packs` 中的 Capability Pack candidate 装成正式 enabled Pack。

## 7. 后续授权入口

| 切片 | 下一步 | 授权要求 |
| --- | --- | --- |
| P12 | 安全内置能力 lifecycle 样本 | 需要 `truzhenos` 和 `truzhen-client-web-desktop` 修改 / 测试授权；不运行真实 Codex CLI，不执行第三方 OSS。 |
| P13 | GUI lifecycle 面板 | 需要 client 修改 / 测试授权；必要时 `truzhenos` 只读聚合字段授权。 |
| P15 | 三候选 GUI 实操 | 需要 client + `truzhenos` 测试授权，并允许写 Pack 仓 GUI 证据报告。 |
| P16 | 受控 Code Assistant 最小 run | 需要 11 Gateway、Owner/Base Gate、Receipt 授权；不应用补丁。 |
| P17 | provider / adapter candidate | 需要 `truzhen-software` 或外部 provider 仓授权。 |
| P18 | 云市场 sandbox | 需要 `truzhen-cloud`、client、`truzhenos` 授权；不真实支付、不生产发布。 |

## 8. 自审

| 检查项 | 结果 |
| --- | --- |
| Spec coverage | 已覆盖 Owner 新指令、旧计划 P3-P18、Pack 仓候选资产、client 能力制作台只读测试和 blocked 边界。 |
| Placeholder scan | 无 `TBD`、`TODO`、无“稍后补”；未授权事项均写为 blocked 或后续授权切片。 |
| Type / path consistency | 路径使用当前新 worktree 和实际 client 仓路径；不再引用旧 worktree 作为本轮计划的执行根。 |
