# 智能家居服务商 Pack 单主项目全程测试计划 v13

日期：2026-07-14

状态：`当前客户主线；v13 计划设计中；急停零副作用已复验，流式完成态 P2 已实现 -> 已接线、待真实 GUI 复验；未验收、未发布`

前序证据：`/Users/li/Documents/过程文档/smart-home-service-provider-v12-20260713/测试报告.md`、`/Users/li/Documents/过程文档/smart-home-service-provider-v12-20260713/问题台账.md`

被替代计划：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v12-20260713.md`

## 1. 派活卡与最小闭环

| 维度 | v13 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，只验 `smart-home-owner-pack-v0`；环保 Pack 结果不得冻结本计划。 |
| 我要做的事 | 真实登录、真实本地 cloud 市场安装，以**一个智能家居主项目**贯穿候选、停启、重启、断网、重放、急停、响应式、作者工作台与独立复核，补齐 v12 未覆盖项。 |
| 真实场景证据 | v12 市场安装、授权、双项目隔离、停启和急停单/双标签已通过；P2 是模型正文已完成后流式卡和“思考中”不消失。非流式、重启/断网/重放、多视口、作者工作台仍未覆盖。 |
| 最小可交付 | 同一制品、安装实例和主项目完整跑完 R1/R2/R3；证明一次真实本地模型只生成一个待 Owner Candidate，停用和急停均零副作用，恢复不补写旧请求。 |
| 砍掉 / backlog | 生产账号/扣款、客户生产数据、真实发送、正式化、真实执行、Frappe 写回、生产部署、发布、完整 Dendrite。 |
| 真相源 | cloud 持商品/订单/Entitlement/digest；14 持 lifecycle；05/06 持项目/Run；13 持会话；01 持急停；07 持任务；08 持模型；03 持回执；11 持执行/Frappe。client 只投影。 |
| 仓库 / 层归属 | OS 负责急停/lifecycle 硬门；client 负责流式终态和 ReadModel 展示；本仓负责 Pack 声明/计划；cloud 负责真实市场。无 contracts 变更。 |
| 风险 / AI 角色 | 黄：候选/lifecycle；橙：市场/Provider；红：急停、Gate、Receipt、发送、执行、Frappe。AI 只产 Candidate，Owner + Base 保留裁定权。 |
| 契约影响 | 只改内部实现和测试，不改 contracts、DTO/schema、Gate、Receipt、Candidate、ReadModel、Surface 或主权边界。 |
| 上下文 | 只读本计划、v12 证据、该 Pack、隔离运行数据库和相关 cloud/OS/client 页面；不读取其它客户资产。 |
| 禁止边界 | 不读真实 secret、不真实扣款、不 API 注入 Entitlement、不真实发送/执行/Frappe、不正式化、不把浏览器状态当项目真相、不用 mock 冒充通过。 |
| 生命周期 | 当前`已实现 -> 已接线`；本计划完成并独立复核后才升为`已验收（打包前）`，发布另行裁定。 |

## 2. 执行授权与启动基线

本文件只创建测试任务，尚未启动服务、浏览器或跨仓操作。正式执行前，Owner 必须分别明确授权下表所列仓库的测试与必要的小问题修复，并在每仓记录固定主仓 `main`、隔离 worktree 的 branch/HEAD/status/WIP 和完整 diff hash；不得以本计划文件或旧轮次授权替代该确认。模拟支付仅可使用本地 `MOCK_NOT_PAYABLE` sandbox，绝不产生真实扣款。

| 仓库 | 真相与职责 | v13 允许动作 | 禁止边界 |
| --- | --- | --- | --- |
| `truzhen-packs` | Pack 制品、声明、静态审计和计划 | 静态校验、制品构建、隔离安装脚本、测试证据登记 | 不保存运行态、订单、Entitlement 或 Provider 实现；不发布。 |
| `truzhenos` | lifecycle、项目/Run、急停、Gateway、Receipt 和队列真相 | 隔离 devserver、只读对账、P2/P3 小修及定向回归 | 不改 contracts、Gate/Receipt/权限语义；发现 P0/P1 先冻结危险路径。 |
| `truzhen-client-web-desktop` | 用户可见的流式终态、项目与安全投影 | 真实 GUI、Computer Use、P2/P3 小修、定向 UI 回归 | 不以接口或数据库查询替代 GUI 通过；不改 Surface/ReadModel 契约。 |
| `truzhen-cloud` | 本地市场、商品、模拟订单、Entitlement、下载 digest | 仅隔离 full-stack lane 的 GUI 验收和只读对账 | 不接生产支付、不写生产订单/授权、不发布/上架。 |
| `truzhen-contracts` | schema / DTO 权威契约 | 只读核对 | 本轮不得修改。 |

- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhenos`。
- Pack：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs`。
- Client：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-client-web-desktop`，包含 `done` 后不等待 HTTP EOF 的流式完成态修复。
- Cloud：干净授权分支的隔离 full-stack lane。

四仓分别登记 branch/HEAD/status/WIP；不得覆盖既有 WIP。登记端口、容器、运行参数、商品 version 与制品 SHA-256。每轮使用全新的 HOME、SQLite、浏览器 profile 和端口；测试报告、截图、网络异常记录及运行态快照统一写入 `/Users/li/Documents/过程文档/smart-home-service-provider-v13-20260714/`，不进入 Git。

## 3. 单主项目与尽量跑完全程纪律

1. 创建唯一主项目 `SH-C01-v13-王先生灯组离线售后`，同一 `transaction_ref`、Run、安装实例和浏览器 profile 贯穿 R1/R2/R3；不再用每轮新项目替代连续性验收。
2. lifecycle、OS/cloud 重启、刷新、断网、旧请求重放和急停都在该项目上验证。恢复动作必须使用新请求，但不换项目。
3. 仅当主项目真相被缺陷污染、删除或不可恢复时可创建 `SH-C01-RECOVERY-01`；必须登记旧/新 refs、原因和污染边界，并用恢复项目跑完余下全程。
4. 如必须复核跨项目隔离，只能建立一个最小只读 `SH-C02-ISOLATION-PROBE`；它只验证正文/ref 不串，不生成第二套 Pack 全流程证据，不计为第二验收项目。
5. 失败后先继续所有可信安全 lane，目标是把计划尽量跑完，而不是遇到一个问题就结束轮次。

## 4. 自治修复与停线规则

- P2/P3 及影响清晰、可逆、不改契约/主权链的中等问题：执行者自行修根因、补定向和邻接回归，从主项目最近可靠点继续。本轮流式完成态 P2 必须按此复验关闭。
- P0/P1：登记真相源、污染范围、副作用、冻结路径与恢复门。市场、登录、静态审计、响应式、只读回放或其它独立 lane 可信时继续。
- 只有 Gate/Receipt/权限/跨项目隔离失真，且所有剩余 lane 都依赖该状态时暂停；必须写明重启门，不得笼统“发现 P1 所以停止 R2/R3”。
- 不降低断言、不吞错、不改库、不复用旧 PASS、不以重建项目掩盖缺陷。

## 5. 开跑前恢复门

1. `secretaryCommunicationPipeline` 测试证明收到协议 `done` 后即使 HTTP 不 EOF 也立即完成并取消 reader；页面测试证明“秘书正在思考中”和“候选摘要等待主人评定”消失。
2. 非流式与流式入口在 dialogue/stream header/07/03/08/turn/发送候选前读取同一 01 急停真相；急停读取失败也 fail-closed。
3. Pack 停用硬门、项目刷新恢复和 transaction 同源回归保持；OS 定向 test/race/vet、client 定向测试/typecheck/build、Pack 审计、cloud 静态门和 `git diff --check` 全绿。

## 6. 统一证据矩阵

每个关键动作同时记录 GUI、01、03、05/06、07 全队列、08 usage/ModelRun、09 Memory、10 Send、11 Execution/Frappe、13 turn、14 lifecycle。GUI 证据必须来自真实前端可见控件、菜单、状态和读屏；只读 API/SQLite 只用于三方对账，不能替代 GUI 通过。停用或急停期间以下业务增量必须全部为 0：Provider/stream/model/tokens/ModelRun、turn、Candidate、TaskCandidate/FormalTask、业务 Receipt、Run 推进、Memory、Send、Execution、Frappe。01 自身急停治理事件单列，不混入业务副作用。

## 7. R1：真实市场与单项目核心链

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v13-20260714/`。

1. canonical builder 生成净 ZIP/SHA-256；真实 GUI 先证明未授权阻断，再经页面本地模拟支付获得 Entitlement 并安装；上传、下载、安装 digest 一致，registry 0→1。
2. 对账 Pack enabled、Role Pack/SlotBinding、ProviderRequirement 诚实态。创建唯一主项目，启动 06 Run。
3. enabled 下流式和独立 GUI 非流式入口各发一条明确的新请求；每次恰好一次真实本地模型、一个待 Owner Candidate，13/08/03/07 与主项目/Pack 同源；完成态不残留“思考中/流式生成”。
4. 0 FormalTask、0 Send、0 Execution、0 Frappe、0 正式记忆。刷新后仍从 05 恢复同一当前项目，失权/不存在时清除并提示，不擅自选择第一项。

## 8. R2：同一项目停启、重启与断网

GUI 停用并确认 14 disabled 后，攻击无项目领域短语、已选主项目、非流式、流式、conversation turn、tool confirm、旧 request 重放；完整证据矩阵增量为 0。随后在同一项目执行 `reactivate -> disable -> reactivate`，覆盖 OS/cloud 重启、登录过期重登、断网/恢复、Owner/actor 差异、幂等与 Provider timeout/recovery。重新启用后仅新请求恢复，停用期请求不得延迟补写。

## 9. R3：同一项目急停、重放与体验

1. GUI 启用急停，01 ReadModel 与顶部一致后执行非流式、流式、双击并发、双 Tab、旧 request 和伪 receipt 重放；07 全队列保持 N→N，完整矩阵全 0。
2. 急停中不点击正式化、发送、执行或 Frappe；入口阻断足够，不以危险真实动作探测。
3. GUI 解除急停；只有全新请求恢复一次模型/候选，急停请求不补写。
4. 同一主项目完成 390×844 / 1024×768 / 1440×900、Provider `ready/not_ready/blocked/degraded`、历史回放、作者工作台只读链和 Owner UAT。每个在本轮范围内的可见按钮、菜单和状态须逐项标为通过、受控阻断、缺口或不适用；前端没有独立 GUI 入口的项目只能登记覆盖缺口，不能用 API 操作补记为 GUI PASS。可选隔离探针若发现正文/ref 泄漏，按 P0 冻结污染路径，但继续不依赖该路径的安全项。

## 10. 独立验收与收尾

未参与实现的验收主体复核自动门、真实 GUI、市场/Entitlement/digest、主项目三阶段连续时间线、流式终态、停用/急停完整零增量矩阵、可选隔离探针、响应式与端口/容器清理。交付物至少包括轮次参数单、按钮/状态覆盖表、截图索引、Receipt 索引、运行态快照、问题台账和最终报告。开放 P0=0、阻塞 P1=0、所有 P2 有关闭证据或 Owner/期限/复测入口、R1/R2/R3 全 PASS 才标记`已验收（打包前）`；否则按第 4 节继续安全 lane 或诚实暂停，不得发布。
