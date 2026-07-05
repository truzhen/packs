# 能力 Pack 制作台商用化下一步影响清单

> 日期：2026-07-04
> 归属：`truzhen-packs` 文档计划。本文件只做影响评估和下一切片定义；P11 已获授权并完成接线，P12 仍未获跨仓施工授权。
> 关联计划：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/gui-capability-pack-workbench-github-oss-test-plan-20260704.md`

## 1. 当前结论

短视频运营压测已推进到 P3-P11：制作台可以生成 Code Assistant 调用候选、11 run 请求候选、candidate bundle、candidate bundle dry-run 阻断、lifecycle preflight 阻断、PatchCandidate 承接和复核候选。

这仍不是商用可交付。商用缺口不是“给候选目录补一个独立 Capability Pack folder loader”这么简单。基座治理已明确：Domain Work Pack 是可外迁 folder pack；Capability Pack 是基座常驻执行 provider + 声明。`truzhen-packs` 至多保存能力 descriptor / 候选声明，真实执行 provider、readiness、Gate、Receipt 和启用指针归 `truzhenos`。

P11 已证明下一步不应让 P10 dry-run 通过安装，而是先把 `candidate bundle -> lifecycle draft preflight`。它把候选资产映射到现有 04 lifecycle 的前置检查，明确哪些条件会阻断，哪些条件满足后才能进入 draft / readiness / promote / confirm。

P11 详细施工规格见：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-lifecycle-preflight-execution-spec-20260704.md`

下一步建议进入 P12：只用基座内置安全 provider / fixture 验证 lifecycle draft / readiness / promote / confirm 最小闭环，不直接使用 MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload，不运行真实 Codex CLI，不执行第三方 OSS。

P12 详细施工规格见：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`

## 2. 目标仓与允许动作

| 仓库 | 本切片职责 | 默认允许动作 | 需要 Owner 重新授权的动作 |
|---|---|---|---|
| `truzhen-packs` | 记录影响清单、候选资产台账、短视频能力 Pack descriptor 候选 | 文档和候选声明 | 把候选标成已启用、写 provider 实现、发布 |
| `truzhenos` | 04 Capability Studio / lifecycle / readiness；01 Gate；03 Receipt；11 Execution Gateway Code Assistant provider | 当前只读分析 | 修改后端 lifecycle、启动 devserver、接 Gate / Receipt、运行 Codex CLI |
| `truzhen-client-web-desktop` | 制作台 GUI 显示 lifecycle preflight、blocked reason、readiness issue | 当前只读分析 | 修改页面、跑 GUI 自动化、录屏 |
| `truzhen-contracts` | 若现有 DTO 不能表达 lifecycle preflight，再作为真相源先改 | 只读 | 新增或修改跨仓 schema / DTO |
| `truzhen-software` | 后续 adapter / provider 候选落点 | 不碰 | 生成或修改 provider / adapter 候选 |
| `truzhen-cloud` | 后续市场上架 / License / Entitlement | 不碰 | 上传、上架、购买、下载、安装 |

## 3. 风险颜色

| 风险 | 内容 | 处理 |
|---|---|---|
| 绿 | 本文件、候选台账、readiness 缺口记录 | 可直接补文档并跑文档级验证 |
| 黄 | GUI 显示 lifecycle preflight、blocked reason、readiness issue | 跨仓授权后可施工，必须有前端测试 |
| 橙 | lifecycle draft / readiness / promote 语义、DTO、跨仓接口 | 先出影响清单和兼容策略，再施工 |
| 红 | 01 Base Gate、03 Receipt、真实启用、真实 Codex CLI run、第三方 repo 执行、社媒登录 / 上传 | 只能在明确授权、隔离环境和回执链齐全后执行 |

## 4. 推荐切片

### P11：Lifecycle preflight（已接线）

目标：让制作台对 candidate bundle 做 lifecycle 前置检查，而不是假装安装。

验收：

- 后端新增或复用 preflight 能力，输入 candidate bundle / session delivery ref。
- 没有 studio delivery 时返回 `delivery_artifact_required`。
- candidate bundle 仍是 candidate-only 时返回 `lifecycle_preflight_blocked`，原因含 `candidate_bundle_not_delivery`。
- evaluation readiness 不是 `ready` 时返回 `provider_missing / not_ready` issue。
- 返回可读的下一步：需要 delivery、golden cases、evaluation ready、provider dependency、Owner/Base Gate。
- 不写 enabled version，不写 Formal Receipt，不运行 provider。

### P12：安全内置能力样本 lifecycle

目标：只选择一个已在基座内有安全 provider / fixture 的能力样本，验证 draft / readiness / promote / confirm 的最小闭环。不得用 MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload 直接跑真实代码。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

验收：

- draft 由服务端按 studio delivery 组装，不接收客户端伪 spec。
- readiness 必须聚合 delivery、golden cases、evaluation ready 和 provider dependency。
- confirm 必须经过 01 Gate 两段，并绑定 03 Receipt。
- 前端显示 enabled pointer 和 receipt ref。

### P13：制作台 GUI lifecycle 面板

目标：让 GUI 在短视频候选包页显示 lifecycle 状态，而不只显示导出文件。

验收：

- 显示 candidate bundle、delivery、readiness、promote、confirm 的分层状态。
- blocked reason 用人能理解的文案展示。
- 用户不能点击“启用”绕过 readiness / Gate / Receipt。
- 未接 provider 时按钮禁用并显示 `provider_missing / not_ready`。

### P14：商用化缺口审计

目标：用短视频运营三类能力 Pack 候选形成商用缺口台账。

验收：

- 每个候选 Pack 有 lifecycle preflight 结果。
- 每个阻断项有仓库归属、风险色、修复建议、验收命令。
- 明确哪些进入当前主线，哪些进入 backlog。
- 不把 candidate-only、readmodel-only、demo-only 写成已发布。

## 5. 禁止路径

1. 禁止把 Capability Pack 当 Domain Work Pack 一样补 folder loader 后直接 enabled。
2. 禁止把 `candidate_bundle_dry_run` 的静态校验通过写成 lifecycle 成功。
3. 禁止让 `truzhen-packs` 保存 provider 实现、真实执行脚本、raw token 或社媒凭据。
4. 禁止 04 直接运行 Codex CLI；真实执行必须走 11 Execution Gateway、Owner/Base Gate 和 Receipt。
5. 禁止直接安装、运行或测试 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
6. 禁止真实登录、上传、发布短视频或读取浏览器 cookie。

## 6. 需要 Owner 裁定

1. 下一施工切片是否按 P12 `安全内置能力样本 lifecycle` 开始，而不是直接做第三方 OSS provider 或独立 Capability Pack loader？
2. 是否授权修改 `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` 的 04 lifecycle / readiness / safe sample 相关代码和测试？
3. 是否授权修改 `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` 的制作台 lifecycle GUI？
4. 本轮是否继续不改 `truzhen-contracts`，除非 P12 证明现有 DTO 表达不了安全样本 lifecycle？
5. 是否继续禁止真实 Codex CLI run、第三方 OSS 执行、社媒登录 / 上传，只做安全 fixture 与阻断？

## 7. 建议验收命令

文档 / 本仓：

```sh
git diff --check
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

后续若获准改 `truzhenos`：

```sh
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

后续若获准改 `truzhen-client-web-desktop`：

```sh
npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts
npm run typecheck
npm run smoke:frontend-shell
git diff --check
```
