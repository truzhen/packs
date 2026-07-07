# 能力包相关全量重测报告

> 日期：2026-07-06
> 分支：`codex/capability-pack-test-replan-20260706`
> 计划：`/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706/docs/plans/capability-pack-related-full-retest-plan-20260706.md`
> 当前结论：Pack 仓能力包候选资产重测通过；client 能力制作台相关类型检查、目标测试、全量 Vitest、静态红线、前端 shell smoke 和 build 通过；`scripts/smoke-pack-01.cjs` 失败，根因为旧 smoke 仍断言 `/v2/pack-studio/*` 和 `src/api/client.ts`，与当前 `/v3/pack-studio/*` 和域拆分实现不一致。

## 范围

本轮只覆盖：

- `truzhen-packs` 的短视频运营 Capability Pack 候选资产、计划、证据门禁和禁入扫描。
- `truzhen-client-web-desktop` 的能力制作台 / Pack Studio / Code Assistant Pack Panel 只读测试。

本轮不覆盖：

- `truzhenos` lifecycle / Gate / Receipt live 验证。
- 真实 Codex CLI run。
- MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 的源码拉取、安装或执行。
- 社媒登录、上传、发布。
- `truzhen-contracts`、`truzhen-software`、`truzhen-cloud` 修改或测试。

## 已跑命令与结果

### truzhen-packs

工作目录：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/capability-pack-test-replan-20260706`

| 命令 | 结果 |
| --- | --- |
| `go test ./... -count=1` | 通过，`ok github.com/truzhen/packs` |
| `python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"` | 通过，`JSON 合法` |
| Pack 结构审计脚本 | 通过，`Pack 结构审计 OK` |
| `find . -name install.py -o -name uninstall.py \| xargs -r python3 -m py_compile` | 通过，`脚本语法 OK` |
| `git ls-files \| rg '(^\|/)(__pycache__\|node_modules\|dist\|build\|\.vite)(/\|$)\|\.(db\|sqlite\|log\|jsonl\|pyc)$'` | 通过，已跟踪文件无禁入产物 |

### truzhen-client-web-desktop

工作目录：

`/Users/li/Documents/truzhen-client-web-desktop`

当前分支：

`feat/o9-market-pack-install-client-20260705`

| 命令 | 结果 |
| --- | --- |
| `npm run typecheck` | 通过 |
| `npm run test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/pages/__tests__/capabilityStudioGateHint.test.tsx src/pages/__tests__/capabilityStudioReceipt.test.tsx src/pages/__tests__/capabilityStudioCopyGuard.test.tsx src/pages/__tests__/capabilityPackNearCardStatus.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx` | 通过，6 个文件、39 个测试通过 |
| `npm run test` | 通过，159 个文件、1184 个测试通过；期间本地 `127.0.0.1:18080` 未启动产生 `ECONNREFUSED` 日志，但退出码为 0 |
| `npm run check:static` | 通过，335 个文件断言通过 |
| `npm run smoke:frontend-shell` | 通过 |
| `npm run build` | 通过；生成的 `dist/` 已清理 |
| `node scripts/smoke-pack-01.cjs` | 失败，失败断言为 `Pack Studio API 方法匹配后端路径` |

## 失败项

### CPK-RT-001：Pack Studio 旧 smoke 断言过期

失败命令：

```bash
node scripts/smoke-pack-01.cjs
```

失败输出：

```text
FAIL Pack Studio API 方法匹配后端路径
```

根因证据：

- `scripts/smoke-pack-01.cjs` 第 33 行仍要求 `src/api/client.ts` 同时包含 `/v2/pack-studio/status` 和 `/v2/pack-studio/drafts`。
- 当前前端 Pack Studio 相关实现已经迁移到 `/v3/pack-studio/*`：
  - `src/api/client.ts` 包含 `/v3/pack-studio/drafts/candidate`、`/v3/pack-studio/dry-run`、`/v3/pack-studio/export`、`/v3/pack-studio/lifecycle/*`。
  - `src/api/domains/authMarket.ts` 承载 `safePackStudioFallback()`、`createPackDraft()`、`updatePackDraft()`、`riskScanPackDraft()` 等域拆分后的 Pack Studio fallback 方法。
- `git blame` 显示该 smoke 断言来自 2026-06-26 初始脚本，未随 `/v3` 和 API 域拆分更新。

当前处理：

- 本轮未获修改 client 授权，因此不改 `scripts/smoke-pack-01.cjs`。
- 将该失败登记为 client smoke 维护缺口，不把它包装成通过。

建议后续修复方向：

- 将 `scripts/smoke-pack-01.cjs` 的读取范围从仅 `src/api/client.ts` 扩展到 `src/api/client.ts` + `src/api/domains/authMarket.ts`。
- 将 `/v2/pack-studio/status`、`/v2/pack-studio/drafts` 旧断言替换为当前 `/v3/pack-studio/*` 或明确的 fallback preview 口径。
- 修复后运行 `node scripts/smoke-pack-01.cjs`、`npm run check:static`、`npm run test`、`npm run build`。

## 阻断项

| 阻断项 | 原因 | 后续切片 |
| --- | --- | --- |
| `truzhenos` lifecycle live 验证未跑 | 当前 workspace 未开放 `truzhenos`，本轮也未重新获得修改 / 测试授权 | P12 / P13 / P15 |
| 真实 Codex CLI run 未跑 | 本轮禁止真实 Codex CLI，不使用 Owner GPT 会员模型 token | P16 |
| provider / adapter candidate 未跑 | 本轮不碰 `truzhen-software`，不执行第三方 OSS | P17 |
| 云市场 sandbox 未跑 | 本轮不碰 `truzhen-cloud`，不处理支付、License、Entitlement、下载分发 | P18 |
| 社媒发布未跑 | 本轮禁止登录、上传、发布 | 红色单独授权 |

## 不计入完成的证据

- `go test` 通过不等于短视频能力包商用可用。
- `npm run test` 通过不等于 lifecycle enabled 或 Receipt live 通过。
- `node scripts/smoke-pack-01.cjs` 失败不能忽略，也不能用其他绿色测试覆盖。
- candidate bundle、run request candidate、PatchCandidate review candidate 仍不能替代正式 Gate / Receipt。
- `social-auto-upload` 仍只能作为发布草稿 / 上传意图候选样本，不代表真实社媒发布能力。

## 下一授权切片

| 切片 | 最小目标 | 需要 Owner 授权 |
| --- | --- | --- |
| CPK-FIX-001 | 修复 client 旧 smoke 对 `/v2/pack-studio/*` 的过期断言 | 修改和测试 `truzhen-client-web-desktop` |
| P12 | 安全内置能力 lifecycle 样本 draft / readiness / promote / confirm | 修改和测试 `truzhenos` 与 `truzhen-client-web-desktop`；不改 contracts、software、cloud |
| P13 | GUI lifecycle 面板 | 修改和测试 `truzhen-client-web-desktop`；必要时 `truzhenos` 只读聚合字段 |
| P15 | 三候选 GUI 实操 | 测试 `truzhen-client-web-desktop` 与 `truzhenos`，并在 `truzhen-packs` 写入 GUI 证据 |
| P16 | 受控 Code Assistant 最小 run | 允许通过 11 Gateway 执行一次受控 run；不执行第三方 OSS，不应用补丁 |
| P17 | provider / adapter candidate | 修改和测试 `truzhen-software` 或外部 provider 仓 |
| P18 | 云市场 sandbox | 修改和测试 `truzhen-cloud`、`truzhen-client-web-desktop`、`truzhenos`；不真实支付、不生产发布 |
