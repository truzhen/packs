# Truzhen 两 Pack 产品修复收尾实施计划

**执行状态：已验收（2026-07-11）；未提交、未推送、未发布。核心 P0/P1 已归零，P2 仅保留作者真实账户登录与外部 Dendrite 运行覆盖。**

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 关闭智能家居与环保执法两份台账的全部 P0/P1/P2，使 Truzhen 共性链路和两个 Pack 达到可独立验收状态。

**Architecture:** 事实继续归 truzhenos 的 Base、SceneFlow、Task、Memory 与 Receipt；client 只消费权威 ReadModel；cloud 只提供隔离本地全栈；packs 只声明、携带知识与记录生命周期。所有上下文增量复用既有 `payload`，不修改 `truzhen-contracts`。

**Tech Stack:** Go、SQLite/GORM、React/TypeScript、Vitest、Vite、Docker Compose、Python unittest、Markdown。

## Global Constraints

- 当前客户主线；两个 issue register 的 P0/P1/P2 全部进入本轮。
- 只在四个现有隔离 worktree 施工；`truzhen-contracts`、`truzhen-software` 零改动。
- 不提交、不推送、不 merge、不清理其他任务容器或 worktree。
- 不真实外发、不付款、不写 Frappe、不读取生产凭据、不绕过 Owner/Base Gate。
- 所有行为变更遵循 TDD：先看见失败，再最小实现，再专项与相关回归。
- Dendrite 外部 E2E 必须诚实记录；若本机外部依赖未就绪，不得把该门禁包装成通过。

---

### Task 1: 让高风险知识上下文真实进入 SceneFlowRun

**Files:**
- Modify: `truzhenos/backend/internal/businessobject/service/productization_store.go:1482`
- Modify: `truzhenos/backend/internal/businessobject/service/capsule_scene_fields_test.go`
- Modify: `truzhen-client-web-desktop/src/api/client.ts:5130`
- Modify: `truzhen-client-web-desktop/src/api/__tests__/startSceneFlowRun.test.ts`
- Modify: `truzhen-client-web-desktop/src/components/business-objects/CreateTransactionObjectDialog.tsx`
- Test: `truzhenos/backend/internal/devserver/sceneflowdev/scene_flow_knowledge_runtime_test.go`

**Interfaces:**
- Consumes: `TransactionObjectCreateWithSceneRequest.Payload map[string]interface{}`、Pack spec 的 `knowledge_scopes`。
- Produces: `SceneFlowLoadCandidate.StatePayload` 中的 `owner_id`、`pack_ref`、`pack_version_ref`、`knowledge_scope_refs`、`knowledge_as_of_date`、`knowledge_query`、`knowledge_uncertainty`。

- [x] **Step 1: 写 05 失败测试**

在 `capsule_scene_fields_test.go` 新增真实 `CreateTransactionWithScene` 测试，输入：

```go
OwnerRef: "owner://local/default",
Payload: map[string]interface{}{
    "knowledge_scope_refs": []interface{}{
        "knowledge_scope://environmental/legal-basis",
        "knowledge_scope://environmental/water",
    },
    "knowledge_as_of_date": "2026-07-08",
    "knowledge_query": "核验排污许可、管辖与证据缺口",
    "knowledge_uncertainty": "仅有举报照片",
    "raw_secret": "must-not-propagate",
},
```

断言上述允许字段进入 `created.LoadCandidate.StatePayload`，`raw_secret` 不进入，且 `pack_ref/pack_version_ref/owner_id` 来自请求结构而非前端自造正式事实。

- [x] **Step 2: 运行测试确认 RED**

Run: `go test ./internal/businessobject/service -run TestCreateWithScenePropagatesAllowlistedKnowledgeContext -count=1`（workdir=`truzhenos/backend`）  
Expected: FAIL，缺少 `knowledge_as_of_date` 或 `knowledge_scope_refs`。

- [x] **Step 3: 最小实现 05 allowlist**

在 `buildLoadCandidateLocked` 中只复制下列 payload key：

```go
for _, key := range []string{"knowledge_scope_refs", "knowledge_scope_ref", "knowledge_as_of_date", "knowledge_query", "knowledge_uncertainty"} {
    if value, ok := request.Payload[key]; ok {
        statePayload[key] = value
    }
}
statePayload["owner_id"] = request.OwnerRef
statePayload["pack_ref"] = request.ScenePack.ScenePackRef
statePayload["pack_version_ref"] = request.ScenePack.PackVersionRef
```

不得复制任意 payload、secret、token、password 或 authorization 字段。

- [x] **Step 4: 运行 05 测试确认 GREEN**

Run: `go test ./internal/businessobject/service -run TestCreateWithScenePropagatesAllowlistedKnowledgeContext -count=1`  
Expected: PASS。

- [x] **Step 5: 写 client 失败测试**

扩展 `startSceneFlowRun.test.ts`，调用输入增加：

```ts
ownerRef: 'owner://local/default',
knowledgeScopeRefs: [
  'knowledge_scope://environmental/legal-basis',
  'knowledge_scope://environmental/water',
],
knowledgeAsOfDate: '2026-07-08',
knowledgeQuery: '核验排污许可、管辖与证据缺口',
knowledgeUncertainty: '仅有举报照片',
```

断言 create body 的 `owner_ref` 与 `payload` 完整包含这些值。

- [x] **Step 6: 运行 client 测试确认 RED**

Run: `npm test -- --run src/api/__tests__/startSceneFlowRun.test.ts`  
Expected: FAIL，`knowledge_scope_refs` 缺失或默认 owner 不规范。

- [x] **Step 7: 最小实现 client 透传**

给 `startSceneFlowRun` 输入增加 `knowledgeScopeRefs?: string[]`；默认 owner 改为 `owner://local/default`；payload 只在数组非空时加入去空、去重后的 `knowledge_scope_refs`。`CreateTransactionObjectDialog` 从选中 Pack 的 `spec.knowledge_scopes` 提取 scope ref 后传入。

- [x] **Step 8: 运行 client 与 SceneFlow 检索测试**

Run: `npm test -- --run src/api/__tests__/startSceneFlowRun.test.ts src/components/business-objects/__tests__/createTransactionObjectDialogEsc.test.tsx`  
Run: `go test ./internal/devserver/sceneflowdev ./internal/memory/contextcompiler -count=1`  
Expected: 全部 PASS；跨 active mount citation 测试通过。

### Task 2: 修复顶部急停状态与三入口前置阻断证明

**Files:**
- Modify: `truzhen-client-web-desktop/src/components/layout/AppShell.tsx:652`
- Modify: `truzhen-client-web-desktop/src/i18n/statusLabels.ts`
- Modify: `truzhen-client-web-desktop/src/components/layout/__tests__/appShellCopyGuard.test.tsx`
- Modify: `truzhenos/backend/internal/execution/httpapi/pdf_light_parse_candidate_test.go`

**Interfaces:**
- Consumes: Base 返回 `emergency_stop_enabled`、Gate `allow`、Receipt/evidence refs。
- Produces: 单一有效状态“急停已生效”；Gate 行只解释“已允许开启急停”，不解释为真实动作放行。

- [x] **Step 1: 写 AppShell 失败测试**

静态/行为断言要求源码：

```ts
expect(src).toContain("body.emergency_stop === 'emergency_stop_enabled'");
expect(src).toContain("{ label: '急停请求裁定', value: decision.status === 'allow' ? '已允许开启急停' : systemStatusLabel(decision.status) }");
expect(src).not.toContain("{ label: '系统裁定', value: systemStatusLabel(decision.status) }");
```

- [x] **Step 2: 运行测试确认 RED**

Run: `npm test -- --run src/components/layout/__tests__/appShellCopyGuard.test.tsx`  
Expected: FAIL，旧“未知状态 / 放行”投影仍在。

- [x] **Step 3: 最小实现状态归一化**

`AppShell` 将 `enabled` 与 `emergency_stop_enabled` 都识别为急停已生效；摘要写“急停已生效，所有真实动作将被阻断”；行分别显示急停事实、急停请求裁定、候选、原因、解除规则。`statusLabels` 增加 `emergency_stop_enabled: '急停已生效'`。

- [x] **Step 4: 运行 AppShell 测试确认 GREEN**

Run: `npm test -- --run src/components/layout/__tests__/appShellCopyGuard.test.tsx src/components/__tests__/frontendStrict.test.tsx`  
Expected: PASS。

- [x] **Step 5: 补 PDF light/local RED 测试**

新增两个测试：急停时 light multipart body 零读取；local JSON body 零读取，均返回 403 `blocked_by_emergency_stop`，且不写 Receipt。

- [x] **Step 6: 运行三入口测试并确认 GREEN**

Run: `go test ./internal/execution/httpapi -run 'TestPDF(Parse|Light|Local).*EmergencyStop' -count=1`  
Expected: upload/light/local 全部 PASS。

### Task 3: 在能力包管理展示场景包 ProviderRequirement

**Files:**
- Modify: `truzhen-client-web-desktop/src/pages/ScenePlatformPage.tsx`
- Modify: `truzhen-client-web-desktop/src/pages/__tests__/scenePlatformPackLifecycle.test.tsx`
- Reuse: `truzhen-client-web-desktop/src/pages/ScenePackManagePage.tsx`

**Interfaces:**
- Consumes: `getScenePackLifecyclePacks()` 返回的已启用 Pack spec/provider requirements。
- Produces: 能力包管理中的只读“场景包所需外部能力”区域；状态仅允许 `blocked/provider_missing/manual_handoff/not_ready`。

- [x] **Step 1: 写能力页失败测试**

在现有 lifecycle 测试中 mock 一个启用的智能家居 Pack，含三项 Frappe requirement；打开 `能力包管理` 后断言显示场景包名、三项 capability 与 `provider_missing/blocked`，并断言不出现 `ready`。

- [x] **Step 2: 运行测试确认 RED**

Run: `npm test -- --run src/pages/__tests__/scenePlatformPackLifecycle.test.tsx`  
Expected: FAIL，能力页没有场景包 requirements 区域。

- [x] **Step 3: 最小实现只读投影**

进入能力包管理时读取 lifecycle packs，筛选 enabled 版本，扁平化其 `provider_requirements`；按 `pack_ref + requirement_id` 去重。只展示 requirement/fallback/gateway/risk，不创建 Capability Pack、不触发 Provider。

- [x] **Step 4: 运行测试确认 GREEN**

Run: `npm test -- --run src/pages/__tests__/scenePlatformPackLifecycle.test.tsx src/pages/__tests__/capabilityPackNearCardStatus.test.tsx`  
Expected: PASS。

### Task 4: 补作者登录边界与浏览器隔离验收

**Files:**
- Modify: `truzhen-client-web-desktop/src/pages/__tests__/AuthorWorkbenchRouting.test.tsx`
- Create: `truzhen-cloud/scripts/two-pack-browser-lane-preflight.sh`
- Modify: `truzhen-cloud/scripts/verify.sh`

**Interfaces:**
- Consumes: `PageContent` 的 `authState`、每 lane 的 HOME/端口/profile path。
- Produces: 未登录必须落 AuthPage，已登录才能落 AuthorWorkbench；两个浏览器 lane 使用不同 user-data-dir 和 origin。

- [x] **Step 1: 写作者边界失败/回归测试**

加入未登录用例：

```tsx
render(<PageContent pageId="authorWorkbench" authState="unauthenticated" ... />);
expect(await screen.findByText('登录途真账户')).toBeTruthy();
expect(screen.queryByTestId('page-author-workbench')).toBeNull();
```

保留已有 authenticated 用例，证明没有绕过 auth。

- [x] **Step 2: 运行作者测试**

Run: `npm test -- --run src/pages/__tests__/AuthorWorkbenchRouting.test.tsx src/pages/__tests__/AuthorWorkbenchPage.test.tsx`  
Expected: PASS；若失败，只修测试装配或真实边界，不增加免登录入口。

- [x] **Step 3: 写浏览器 lane preflight**

脚本一次接受两组 `LANE_NAME CLIENT_ORIGIN PROFILE_DIR`，拒绝空值以及重复的 lane、origin 或 profile，输出 JSON：

```json
{"lanes":[{"lane":"smart-home","origin":"http://127.0.0.1:5191","profile":"/tmp/truzhen-chrome-smart-home"},{"lane":"environmental","origin":"http://127.0.0.1:5193","profile":"/tmp/truzhen-chrome-environmental"}]}
```

脚本只验证目录隔离，不启动或控制用户现有 Chrome。

- [x] **Step 4: 接入 cloud verify 并运行**

Run: `bash scripts/two-pack-browser-lane-preflight.sh smart-home http://127.0.0.1:5191 /tmp/truzhen-chrome-smart-home environmental http://127.0.0.1:5193 /tmp/truzhen-chrome-environmental`  
Expected: 两份 JSON 的 origin/profile 均不同。

### Task 5: 证明 cloud 两套真实 full-stack 可并行

**Files:**
- Modify: `truzhen-cloud/scripts/full-stack-isolation-smoke.sh`
- Modify: `truzhen-cloud/modules/08-local-e2e-and-release/full-stack/up.sh`
- Modify: `truzhen-cloud/modules/08-local-e2e-and-release/full-stack/README.md`

**Interfaces:**
- Consumes: `TRUZHEN_FS_PROJECT_NAME/RUNTIME_DIR/*_PORT`。
- Produces: alpha/beta 独立 Compose project、network、volume、container 与入口。

- [x] **Step 1: 先运行真实 preflight**

Run: `TRUZHEN_FS_PROJECT_NAME=tz-pack-alpha TRUZHEN_FS_RUNTIME_DIR=/tmp/tz-pack-alpha TRUZHEN_FS_FORGEJO_PORT=23000 TRUZHEN_FS_NODE_PORT=23333 TRUZHEN_FS_MARKET_PORT=28087 TRUZHEN_FS_GATEWAY_PORT=28090 bash modules/08-local-e2e-and-release/full-stack/up.sh --preflight`  
Run: `TRUZHEN_FS_PROJECT_NAME=tz-pack-beta TRUZHEN_FS_RUNTIME_DIR=/tmp/tz-pack-beta TRUZHEN_FS_FORGEJO_PORT=24000 TRUZHEN_FS_NODE_PORT=24333 TRUZHEN_FS_MARKET_PORT=29087 TRUZHEN_FS_GATEWAY_PORT=29090 bash modules/08-local-e2e-and-release/full-stack/up.sh --preflight`  
Expected: 两者 config 成功且端口、runtime 不同。

- [x] **Step 2: 启动 alpha/beta 最小真实栈**

为两 lane 指定不重叠的 Forgejo、Node、Market、Gateway 端口和 runtime 目录，分别执行 `up.sh`；检查 `docker compose --project-name ... ps` 和两个 gateway health/静态入口。

- [x] **Step 3: 验证隔离后按本 lane down**

只对 `tz-pack-alpha`、`tz-pack-beta` 执行对应 Compose down；不得使用全局 prune、不得停止其它容器。保存 `ps/config/health` 输出到修复证据目录。

- [x] **Step 4: 若环境依赖失败，记录根因**

如果 Docker、镜像或 Dendrite 外部服务不就绪，记录命令、退出码、容器状态与最短恢复动作；不得把静态 smoke 写成真实并行通过。

### Task 6: 两个 Pack 独立 GUI 三轮回归与台账关闭

**Files:**
- Modify: `/Users/li/Documents/过程文档/env-prepack-v2-20260711/issues/register.md`
- Modify: `/Users/li/Documents/过程文档/env-prepack-v2-20260711/R1-repair-closeout-20260711.md`
- Modify: `truzhen-packs/docs/plans/smart-home-service-provider-test-v2-20260711-issue-register.md`（如该测试 worktree 文档不在当前分支，则新建独立关闭报告，不覆盖原报告）
- Create: `/Users/li/Documents/过程文档/two-packs-product-repair-20260711/final-acceptance.md`

**Interfaces:**
- Consumes: 冻结后的四仓 diff、两个 Pack install 脚本、独立 HOME/端口/profile。
- Produces: 每个问题 ID 的 UI/API/Receipt 或静态证据，以及最终生命周期结论。

- [x] **Step 1: 环保三轮**

R1：install、15 mounts、752/752、ProviderRequirement；R2：C01 日期/查询/引用、同 PDF replay；R3：急停 upload/light/local、沟通 403、零发送/零落库。每轮前台操作，API/SQLite 只作对账。

- [x] **Step 2: 智能家居三轮**

R1：项目最新 run 状态；R2：真实 TaskCandidate 首屏与能力页三项 Frappe requirement；R3：顶部急停单一状态、阻断结果/Receipt、作者未登录/受控登录边界。

- [x] **Step 3: 台账逐项关闭**

每个 ID 写：修复文件、失败测试、通过测试、GUI 截图、API/Receipt/SQLite 证据、剩余限制。没有完整证据的项保持 open/blocked。

- [x] **Step 4: 全量相关门禁**

Run (os): `go test ./internal/businessobject/service ./internal/devserver/sceneflowdev ./internal/memory/contextcompiler ./internal/execution/httpapi ./internal/devserver/communicationfabricmessagedraftdev ./tests/execution -count=1`  
Run (client): `npm test -- --run src/api/__tests__/clientGuards.test.ts src/api/__tests__/startSceneFlowRun.test.ts src/components/layout/__tests__/appShellCopyGuard.test.tsx src/components/__tests__/frontendStrict.test.tsx src/components/__tests__/projectDetailTabStatusLabel.test.tsx src/components/__tests__/taskDualStateRealExecution.test.tsx src/components/business-objects/__tests__/createTransactionObjectDialogEsc.test.tsx src/pages/__tests__/scenePlatformPackLifecycle.test.tsx src/pages/__tests__/capabilityPackNearCardStatus.test.tsx src/pages/__tests__/AuthorWorkbenchRouting.test.tsx src/pages/__tests__/AuthorWorkbenchPage.test.tsx && npm run typecheck && npm run build`  
Run (cloud): `bash scripts/full-stack-isolation-smoke.sh && bash scripts/pack-admin-full-stack-static-smoke.sh`  
Run (packs): `python3 -m unittest test_pack_diagnostics test_knowledge_checksums test_pack_bundle test_pack_install_journal test_install_journal_wiring`，随后 JSON/脚本/结构/forbidden-artifact 审计。  
Expected: 相关门禁 0 failure；已知外部 Dendrite 失败单独列出，不混入通过口径。

- [x] **Step 5: 生命周期裁定**

只有 13 项全部有关闭证据、两个 Pack 三轮 GUI 通过、cloud 真实并行通过时，才把状态从 `已接线 / 未验收` 升为 `已验收`。否则保持当前状态并明确剩余 blocker。

## 自审

- 覆盖设计中的 W0-W6 与全部 13 个 issue ID。
- 没有修改 contracts/software、没有新增 Provider 或 Pack 类型。
- 每个产品行为变更都有 RED/GREEN；环境和 GUI 项有明确证据门，不用静态 smoke 替代真实运行。
- 所有步骤均给出具体文件、命令与期望结果；执行中的未知外部状态通过实际命令采集，不预判成功。
- Owner 已要求当前会话直接收尾，因此省略提交步骤；所有分支保持未提交、未推送。
