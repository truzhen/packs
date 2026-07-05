# P12 安全内置能力 lifecycle 样本跨仓施工授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：在不运行真实 Codex CLI、不执行第三方 OSS、不登录或上传社媒的前提下，用基座内置安全能力样本或 fixture 验证 Capability Pack lifecycle draft / readiness / promote / confirm 最小闭环。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p12-cross-repo-authorization-scope.json`

scope_ref：`p12-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P12 跨仓施工前置门：只有记录 Owner 明确授权后，才允许按本卡修改 / 测试 `truzhenos` 与 `truzhen-client-web-desktop`。

## 1. 本次要做的事

P12 只做安全 lifecycle 样本，不做短视频 provider 接通：

- 后端用服务端安全样本 / fixture 生成 lifecycle draft。
- readiness 聚合 delivery、golden cases、evaluation ready、provider dependency。
- promote 不绕过 01 Gate。
- confirm 经 01 Gate 两段并绑定 03 Receipt 测试链。
- 前端制作台展示 draft、readiness issue、promote / confirm 阻断、enabled pointer、receipt ref。
- 本仓候选集记录 P12 状态，仍标记短视频第三方 OSS 为未执行 / blocked。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 04 lifecycle 安全样本、readiness、promote / confirm 测试链、API / acceptance 文档 | 修改、测试 |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | 制作台 lifecycle 面板、确认卡、api client、单测、smoke | 修改、测试、可启动前端 smoke |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | P12 计划、授权卡、候选资产台账 | 已允许文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改；只有现有 DTO 表达不了 P12 时再另开影响清单 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改；后续真实 provider / adapter 才需要 | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改；市场发布 / License / Entitlement 后续再做 | 本批不授权修改 |

## 2.1 机器可验授权门

P12 授权必须同时满足：

- `p12-cross-repo-authorization-scope.json` 中 `current_authorization_evidence.status` 被更新为 Owner 明确授权证据。
- 授权语包含 `required_authorization_phrase_contains` 的全部边界片段。
- 只允许 `authorized_repositories` 中列出的仓库与动作。
- `disallowed_repositories` 和 `forbidden_actions` 保持未触碰、未发生。

在该授权门仍为 `missing` 时，本卡只能作为计划和待授权卡，不得作为跨仓施工授权。

## 3. 预计改动文件

### `truzhenos`

| 文件 | 预期变更 |
|---|---|
| `backend/tests/capability/capability_studio_safe_lifecycle_sample_test.go` | 新增 P12 安全样本 lifecycle 单测 |
| `backend/internal/capability/studio/lifecycle_safe_sample.go` | 新增安全样本 fixture 或服务端组装器 |
| `backend/internal/capability/studio/service.go` | 仅在必要时接入安全样本 lifecycle 入口 |
| `backend/internal/capability/studio/lifecycle_spec.go` | 仅在必要时补 readiness 聚合字段，不改跨仓 contracts |
| `backend/tests/capability/studio_pipeline_http_test.go` | 增加 lifecycle draft / readiness / promote / confirm HTTP 回归 |
| `modules/04-capability-management/api_contract.md` | 登记 P12 安全样本 lifecycle 用法 |
| `modules/04-capability-management/acceptance.md` | 登记 P12 验收条款 |

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `src/api/client.ts` | 复用或补齐 lifecycle draft / readiness / promote / confirm client 方法 |
| `src/components/capability-studio/CapabilityStudioSurfaces.tsx` | 展示安全样本 lifecycle 状态、issue、enabled pointer、receipt ref |
| `src/pages/CapabilityStudioPage.tsx` | P11 preflight 后增加 P12 安全样本 lifecycle 操作区 |
| `src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 P12 GUI 单测 |
| `src/api/__tests__/executionLiveSmokeScript.test.ts` | 增加 P12 smoke 断言 |
| `scripts/smoke-frontend-behavior.cjs` | 增加安全样本 lifecycle 行为验证 |
| `docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P12 前端影响 |
| `FEATURE_LEDGER.md` | 更新前端功能账本 |

### `truzhen-packs`

| 文件 | 预期变更 |
|---|---|
| `capability-pack-candidates/short-video-ops-v0/candidate-set.json` | P12 从 `specified_pending_authorization` 更新为实际接线状态，只有验收通过后才改 |
| `capability-pack-candidates/short-video-ops-v0/README.md` | 更新 P12 完成口径 |
| `README.md`、`MODULES.md`、`FEATURE_LEDGER.md` | 更新 P12 状态和边界 |

## 4. 后端验收命令

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityPackLifecycleSafeFixture' -count=1
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

预期证据：

- draft 由服务端安全样本 / fixture 组装，不接收客户端伪 spec。
- readiness 能返回 `ready` 或明确 `provider_missing / not_ready` issue。
- promote 不写 enabled pointer。
- confirm 必须带 Owner/Base Gate 结果和 03 Receipt ref，之后 lifecycle packs 查询才出现 enabled pointer。

## 5. 前端验收命令

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

预期证据：

- GUI 展示安全样本 draft / readiness / promote / confirm 的分层状态。
- GUI 不能在 readiness / Gate / Receipt 缺失时点击出“启用成功”。
- enabled pointer 和 receipt ref 只在 confirm 成功后出现。
- GUI 仍显示短视频第三方 OSS 只是证据 / blocked，不可运行。

## 6. 本仓验收命令

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
git diff --check
```

## 7. 禁止边界

- 不改 `truzhen-contracts`，除非先输出单独 contracts 影响清单并获授权。
- 不改 `truzhen-software`，不生成真实 provider / adapter。
- 不改 `truzhen-cloud`，不上架、不发布、不支付、不生成 License。
- 不运行 Codex CLI，不消耗 GPT 会员 token。
- 不安装、不运行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
- 不读取或保存 raw token、cookie、账号密码、浏览器登录态。
- 不登录、不上传、不发布任何社媒内容。
- 不把安全 fixture 通过写成短视频 provider 已接通。
- 不把 candidate-only、preflight blocked、fixture-only 写成已商用。

## 8. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P12 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒；只使用基座内置安全样本或 fixture 验证 lifecycle draft/readiness/promote/confirm。
```
