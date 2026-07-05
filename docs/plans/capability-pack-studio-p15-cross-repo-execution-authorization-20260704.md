# P15 三候选 GUI 实操验收跨仓授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：用真实 GUI 路径跑 3 个短视频 Capability Pack 候选的制作台操作，收集截图、行为日志、网络响应摘要和 issue ledger，证明制作台用户路径可验收。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p15-cross-repo-authorization-scope.json`

scope_ref：`p15-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P15 跨仓测试前置门；仍为 `pending_owner_authorization` 时，本卡只作为待授权计划。

## 1. 本次要做的事

P15 做 GUI 实操验收，不做真实 provider 接通：

- 跑 `short-video-draft-generation`、`short-video-composition-orchestration`、`short-video-social-publish-draft` 三条 GUI 剧本。
- 每条剧本记录截图、点击步骤、网络响应摘要、artifact refs、blocked reasons。
- 证明 Code Assistant 仍只到候选或受控 Gateway，不直接运行 Codex CLI。
- 证明 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 只作为 GitHub evidence，不安装、不运行。
- 证明 social-auto-upload 路径不登录、不上传、不发布。
- 本仓写入 `gui-walkthrough-evidence-ledger.md`。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | 启动前端、跑 GUI smoke / 浏览器操作、必要时补 walkthrough 脚本 | 测试；必要时修改测试脚本 |
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 启动隔离 devserver，响应 P3-P13 端点 | 测试；必要时修复 walkthrough 阻断 bug |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | 写入 GUI evidence ledger 和候选集状态 | 文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改；不生成 adapter | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改；不做上架 / 支付 / License | 本批不授权修改 |

## 3. 预计新增证据文件

| 文件 | 用途 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-evidence-ledger.md` | P15 三候选 GUI 实操证据台账 |
| `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-shots/` | 截图目录；只保存制作台截图，不保存账号、cookie、token |

如前端需要补自动化脚本，预计改：

| 文件 | 用途 |
|---|---|
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/scripts/smoke-frontend-behavior.cjs` | 增加三候选 walkthrough 操作 |
| `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui/docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P15 影响 |

## 4. 验收命令

前端 / 行为：

```sh
cd /Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui
npm run smoke:frontend-shell
TRUZHENOS_BACKEND_ROOT=/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss GOWORK=off npm run smoke:frontend-behavior
git diff --check
```

后端，如改动或启动 devserver 发现行为缺口：

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

## 5. 禁止边界

- 不改 `truzhen-contracts`。
- 不改 `truzhen-software`，不生成真实 provider / adapter。
- 不改 `truzhen-cloud`。
- 不运行 Codex CLI，不消耗 GPT 会员 token。
- 不安装、不运行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
- 不读取或保存 raw token、cookie、账号密码、浏览器登录态。
- 不登录、不上传、不发布任何社媒内容。
- 不把截图、GUI 成功点击、blocked 结果写成正式 enabled 或商用完成。

## 6. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P15 授权卡测试 truzhen-client-web-desktop 与 truzhenos，并在 truzhen-packs 写入 GUI 实操证据；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
```
