# P16 受控 Code Assistant 最小 Run 跨仓授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：通过 11 Execution Gateway 执行一次受控 Code Assistant 最小 run，产出 PatchCandidate 文件和 03 Receipt，不运行第三方 OSS，不应用补丁。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p16-cross-repo-authorization-scope.json`

scope_ref：`p16-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P16 受控 run 前置门；仍为 `pending_owner_authorization` 时，不得执行真实 Code Assistant run。

## 1. 本次要做的事

P16 允许一次真实 Code Assistant run，但范围极小：

- 从既有 `code_assistant_run_request_candidate` 进入 11 Gateway。
- 必须经 Owner/Base Gate 签发 `DecisionRef`、`RunID`、`Nonce`。
- 真实 run 只读取候选资产和 README。
- 输出只能是 PatchCandidate / README / adapter scaffold 候选。
- 结果必须绑定 03 Receipt。
- GUI 展示 run receipt、PatchCandidate refs、review required。
- 本仓写入 `code-assistant-controlled-run-ledger.md`。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 11 Gateway / 04 Studio 受控 run、Gate / Receipt、测试 | 修改、测试、可执行一次受控 run |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | GUI 展示 run receipt / PatchCandidate refs / review required | 修改、测试 |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | P16 evidence ledger 和候选集状态 | 文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改；不应用 adapter | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改 | 本批不授权修改 |

## 3. 允许动作

- 可运行一次受控 Code Assistant run。
- 可消耗 Owner 本机 GPT/Codex 会员额度。
- 可写入 11 Gateway 分配的隔离过程目录。
- 可生成 PatchCandidate 文件、summary 和 receipt。

## 4. 禁止边界

- 不安装、不运行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
- 不读取 raw token、cookie、账号密码、浏览器登录态、ssh key。
- 不访问社媒、不登录、不上传、不发布。
- 不写 `truzhen-packs` 正式 pack 目录。
- 不应用补丁、不提交、不 push。
- 不把 PatchCandidate 写成 provider ready 或 enabled。
- 不改 `truzhen-contracts`、`truzhen-software`、`truzhen-cloud`。

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

## 6. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P16 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop，并允许通过 11 Gateway 执行一次受控 Code Assistant 最小 run；本批不改 contracts、software、cloud，不执行第三方 OSS，不登录或上传社媒，不应用补丁。
```
