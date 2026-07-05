# P11 lifecycle preflight 跨仓施工授权卡

> 日期：2026-07-04
> 状态：`已授权并已执行`
> 目标：把短视频运营 Capability Pack 候选从 P10 dry-run 阻断推进到 P11 lifecycle preflight，让制作台能在后端和前端诚实展示“候选 bundle 离正式 lifecycle enabled 还差什么”。

> 备注：Owner 已在 2026-07-04 用本卡授权修改和测试 `truzhenos` 与 `truzhen-client-web-desktop`。本卡是 P11 历史记录，不可复用为 P12 或后续真实启用授权。

## 1. 本次要做的事

实现 P11 `lifecycle_preflight`，但不直接启用 Pack：

- 后端新增 `/v3/capability/studio/lifecycle_preflight`。
- 前端能力制作台增加“检查 lifecycle preflight”按钮和结果卡。
- live smoke 证明 GUI 能看到 `candidate_bundle_not_delivery`、`delivery_artifact_required`、`provider_missing / not_ready` 等商用阻断。
- 短视频候选 Pack 台账记录 preflight 结果，不把候选资产写成已商用。

## 2. 目标仓库与授权动作

| 仓库 | 路径 | 本次职责 | 需要授权的动作 |
|---|---|---|---|
| `truzhenos` | `/Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss` | 04 Capability Studio 后端 endpoint、候选 builder、store stage、HTTP 回归、API 契约 / acceptance | 修改、测试 |
| `truzhen-client-web-desktop` | `/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/capability-pack-studio-short-video-gui` | 制作台 GUI、api client、组件卡片、单测、smoke | 修改、测试、可启动前端 smoke |
| `truzhen-packs` | `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan` | 本仓候选台账、计划、验证登记 | 已允许文档 / 候选资产维护 |
| `truzhen-contracts` | `/Users/li/Documents/truzhen-contracts` | 默认不改；只有现有 DTO 表达不了 P11 时再另开影响清单 | 本批不授权修改 |
| `truzhen-software` | `/Users/li/Documents/truzhen-software` | 默认不改；后续 provider / adapter 候选才需要 | 本批不授权修改 |
| `truzhen-cloud` | `/Users/li/Documents/truzhen-cloud` | 默认不改；市场发布 / License / Entitlement 后续再做 | 本批不授权修改 |

## 3. 预计改动文件

### `truzhenos`

| 文件 | 预期变更 |
|---|---|
| `backend/internal/capability/capability_studio_lifecycle_preflight.go` | 新增 P11 候选类型和 builder |
| `backend/internal/capability/studio/service.go` | 新增 `LifecyclePreflight` 服务方法 |
| `backend/internal/capability/studio/http.go` | 新增 route 和 handler |
| `backend/internal/capability/studio/store.go` | stage 白名单增加 `lifecycle_preflight` |
| `backend/tests/capability/capability_studio_short_video_oss_test.go` | 新增域层 blocked / readiness issue 测试 |
| `backend/tests/capability/studio_pipeline_http_test.go` | 新增 HTTP roundtrip |
| `modules/04-capability-management/api_contract.md` | 补 P11 API 契约 |
| `modules/04-capability-management/acceptance.md` | 补 P11 验收条款 |
| `modules/04-capability-management/capability-studio/short-video-oss-pack-studio-asset-acceptance-20260704.md` | 补 P11 阶段证据 |

### `truzhen-client-web-desktop`

| 文件 | 预期变更 |
|---|---|
| `src/api/client.ts` | 新增 `runCapabilityStudioLifecyclePreflight` |
| `src/components/capability-studio/CapabilityStudioSurfaces.tsx` | 新增 `LifecyclePreflightCard` |
| `src/pages/CapabilityStudioPage.tsx` | P10 后增加 P11 表单 / 按钮 / state |
| `src/pages/__tests__/capabilityStudioWizard.test.tsx` | 新增 P11 GUI 单测 |
| `src/api/__tests__/executionLiveSmokeScript.test.ts` | 增加 smoke 断言 |
| `scripts/smoke-frontend-behavior.cjs` | 增加 live smoke 操作 |
| `docs/frontend-adjustment-impact-ledger-20260704.md` | 登记 P11 影响 |
| `FEATURE_LEDGER.md` | 更新前端功能账本 |

### `truzhen-packs`

| 文件 | 预期变更 |
|---|---|
| `capability-pack-candidates/short-video-ops-v0/docs/lifecycle-preflight-commercial-gap-ledger.md` | P11 实施后补实际 endpoint / evidence |
| `capability-pack-candidates/short-video-ops-v0/candidate-set.json` | P11 状态从 `specified_pending_authorization` 更新为后端 / 前端已接线状态 |
| `docs/plans/capability-pack-studio-short-video-oss-first-slice-report-20260704.md` | 追加 P11 证据 |
| `FEATURE_LEDGER.md`、`README.md`、`MODULES.md` | 更新状态口径 |

## 4. 后端验收命令

```sh
cd /Users/li/Documents/truzhenv3worktree/truzhenos-mod-04-capability-pack-studio-short-video-oss
GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1
git diff --check
```

预期证据：

- `TestCapabilityStudioShortVideoLifecyclePreflightBlocksCandidateBundleBeforeDelivery` 通过。
- `TestCapabilityStudioShortVideoLifecyclePreflightReportsProviderMissingReadiness` 通过。
- HTTP roundtrip 覆盖 `/v3/capability/studio/lifecycle_preflight`。
- 返回 `lifecycle_preflight_blocked / candidate_bundle_not_delivery` 时仍是 200 + artifact，不冒充 lifecycle draft。

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

- GUI 显示“候选 bundle 不是正式 delivery”。
- GUI 显示 delivery、黄金用例、evaluation ready、provider dependency、Owner/Base Gate、Receipt 缺口。
- GUI 不出现可绕过的“启用成功”。
- live smoke 仍证明 Codex CLI 未真实运行、第三方 OSS 未执行、社媒未上传。

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
- 不把 candidate-only、dry-run、preflight blocked 写成已启用或已商用。

## 8. Owner 可回复的授权短语

历史授权语：

```text
授权按 P11 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
```

后续 P12 需使用新的授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`
