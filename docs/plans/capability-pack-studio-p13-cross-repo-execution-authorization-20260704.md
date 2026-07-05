# P13 GUI lifecycle 面板跨仓施工授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：让能力 Pack 制作台用 GUI 面板展示 candidate bundle、delivery、readiness、promote、confirm、enabled pointer、receipt ref 和阻断原因，并禁止用户绕过 readiness / Gate / Receipt。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p13-cross-repo-authorization-scope.json`

scope_ref：`p13-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P13 跨仓施工前置门；仍为 `pending_owner_authorization` 时，本卡只作为待授权计划。

## 1. 本次要做的事

P13 做 GUI lifecycle 面板，不做真实 provider 接通：

- 前端新增 lifecycle 分层面板。
- 面板展示 candidate bundle 与正式 delivery 的区别。
- 面板展示 `provider_missing / not_ready / blocked` issue 和目标仓归属。
- 启用按钮在 readiness、Owner/Base Gate、Receipt 任一缺失时禁用。
- smoke 证明不会从 GUI 触发第三方 OSS 执行、真实 Codex CLI run 或社媒上传。
- 本仓候选集记录 P13 状态，仍保持短视频第三方 OSS evidence-only / blocked。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | GUI lifecycle 面板、按钮禁用、测试、smoke | 修改、测试、可启动前端 smoke |
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 仅在前端无法从现有响应得到字段时补只读聚合和 HTTP 回归 | 修改、测试，限 P13 所需字段 |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | P13 计划、授权卡、候选资产台账 | 已允许文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改；若需要稳定 schema 另开影响清单 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改；不生成真实 provider / adapter | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改；不做上架、支付、License | 本批不授权修改 |

## 3. 预计改动文件

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `src/components/capability-studio/CapabilityStudioSurfaces.tsx` | 新增或扩展 lifecycle 面板组件 |
| `src/pages/CapabilityStudioPage.tsx` | 将 P11/P12 状态接入面板，禁用不安全启用动作 |
| `src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 P13 GUI 单测 |
| `src/api/__tests__/executionLiveSmokeScript.test.ts` | 增加面板 smoke 断言 |
| `scripts/smoke-frontend-behavior.cjs` | 增加 live GUI 行为检查 |
| `docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P13 影响 |
| `FEATURE_LEDGER.md` | 更新商用缺口口径 |

### `truzhenos`

| 文件 | 预期变更 |
|---|---|
| `backend/internal/capability/studio/service.go` | 仅在必要时补只读 panel 字段聚合 |
| `backend/tests/capability/studio_pipeline_http_test.go` | 仅在后端改动时补 HTTP 回归 |
| `modules/04-capability-management/acceptance.md` | 记录 P13 展示和禁用动作验收 |

## 4. 前端验收命令

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

预期证据：

- GUI 展示 lifecycle 面板。
- 面板展示 candidate bundle、delivery、readiness、promote、confirm、enabled pointer、receipt ref。
- 阻断状态显示 target repository 和 next action。
- readiness / Gate / Receipt 缺失时启用按钮禁用。
- GUI 不触发第三方 OSS 执行、不触发真实 Codex CLI run、不触发社媒上传。

## 5. 后端验收命令

仅在 P13 改了 `truzhenos` 时运行：

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

预期证据：

- panel 所需字段来自只读聚合，不创建 draft、不 enabled、不写 Formal Receipt。
- `lifecycle_preflight_blocked` 仍保持 candidate-only。

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
- 不把 GUI 面板显示、blocked、fixture-only 写成已商用。

## 8. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P13 授权卡修改和测试 truzhen-client-web-desktop；必要时只为 P13 只读聚合字段修改和测试 truzhenos；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
```
