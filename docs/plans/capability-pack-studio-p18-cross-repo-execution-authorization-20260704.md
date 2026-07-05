# P18 云市场 Sandbox 跨仓授权卡

> 日期：2026-07-04
> 状态：`待 Owner 授权`
> 目标：在 `truzhen-cloud` sandbox 验证短视频能力 Pack 候选的商品草稿、License / Entitlement、下载分发和安装预检链路，不做生产发布或真实支付。

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p18-cross-repo-authorization-scope.json`

scope_ref：`p18-cross-repo-authorization-scope://short-video-ops-v0`

该 JSON 是 P18 云市场 sandbox 前置门；仍为 `pending_owner_authorization` 时，不得修改 `truzhen-cloud`。

## 1. 本次要做的事

P18 只做 sandbox：

- 创建 sandbox PackListingDraft。
- 创建 sandbox review / License / Entitlement fixture。
- 验证 missing entitlement 时下载 / 安装 blocked。
- 验证 sandbox entitlement 可进入 download / install preflight。
- client 展示 sandbox listing 和 entitlement 状态。
- `truzhenos` 在 lifecycle 前检查 entitlement。
- 本仓记录 cloud sandbox ledger。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | sandbox listing、review、License / Entitlement、download fixture、测试 | 修改、测试 |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | sandbox market / download / install preflight 展示 | 修改、测试 |
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | entitlement check before lifecycle | 修改、测试 |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | cloud sandbox ledger 和 candidate refs | 文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改 | 本批不授权修改 |

## 3. 允许动作

- 使用 sandbox 数据库 / fixture。
- 生成 sandbox listing、sandbox entitlement、sandbox download ref。
- 跑云市场 sandbox 测试。
- 前端展示 sandbox badge。
- `truzhenos` 做安装预检，不做生产 enabled。

## 4. 禁止边界

- 不真实支付、不扣款。
- 不生产发布、不公开上架。
- 不生成生产 License / Entitlement。
- 不上传 raw secret。
- 不把本仓 manifest 字段当云真相。
- 不运行第三方 OSS。
- 不登录或上传社媒。
- 不改 `truzhen-contracts` / `truzhen-software`。

## 5. 验收命令

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

## 6. Owner 可回复的授权短语

推荐授权语：

```text
授权按 P18 授权卡修改和测试 truzhen-cloud、truzhen-client-web-desktop 与 truzhenos，并在 truzhen-packs 写入云市场 sandbox 证据；本批不改 contracts、software，不真实支付、不生产发布、不执行第三方 OSS、不登录或上传社媒。
```
