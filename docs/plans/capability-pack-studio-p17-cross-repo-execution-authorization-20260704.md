# P17 Provider / Adapter Candidate 跨仓授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：把短视频能力 Pack 的真实 provider / adapter candidate 放到 `truzhen-software` 或外部 provider 仓，而不是放进 `truzhen-packs`。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p17-cross-repo-authorization-scope.json`

scope_ref：`p17-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P17 provider / adapter candidate 前置门；仍为 `pending_owner_authorization` 时，不得修改 `truzhen-software`。

## 1. 本次要做的事

P17 只做 adapter candidate 和 readiness，不执行第三方 OSS：

- 为 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 三条能力分别建立 provider candidate。
- 每个 candidate 有 manifest、README、risk boundary、readiness check。
- readiness 默认 `provider_missing / blocked`。
- `truzhen-packs` 只登记 ProviderRequirement 和 evidence。
- `truzhenos` / client 只消费 readiness，不直接执行 provider。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | provider / adapter candidate scaffold、readiness check、风险文档 | 修改、测试 |
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 如需汇入 readiness，则补只读 readiness 消费和测试 | 修改、测试 |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | 如需展示 provider candidate readiness，则补 GUI | 修改、测试 |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | provider candidate ledger 和 capability candidate refs | 文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改 | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改 | 本批不授权修改 |

## 3. 允许动作

- 创建 provider / adapter candidate scaffold。
- 创建 readiness check，默认返回 `provider_missing / blocked`。
- 创建测试证明未接通时不假成功。
- 更新 GUI readiness 展示。

## 4. 禁止边界

- 不把 provider 代码放进 `truzhen-packs`。
- 不 vendor MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 源码。
- 不安装、不运行第三方依赖。
- 不登录、不上传、不发布社媒内容。
- 不保存 raw token、cookie、账号密码。
- 不把 readiness check 存在写成 provider ready。
- 不改 `truzhen-contracts` / `truzhen-cloud`。

## 5. 验收命令

`truzhen-software`：

```sh
cd /Users/li/Documents/truzhen-software
git status --short --branch
git diff --check
```

如仓库已有测试脚本，运行 provider readiness 相关测试。

`truzhenos`，如改动：

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

`truzhen-client-web-desktop`，如改动：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
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

## 6. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P17 授权卡修改和测试 truzhen-software；必要时只为 provider readiness 展示修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、cloud，不执行第三方 OSS，不登录或上传社媒，不保存 raw secret。
```
