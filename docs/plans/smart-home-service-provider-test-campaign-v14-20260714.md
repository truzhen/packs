# 智能家居服务商 Pack 直接 lifecycle 单项目全程测试计划 v14

日期：2026-07-14

状态：`当前客户主线；SHV13-P2-002 已实现 -> 已接线，待真实 GUI 复验；未验收、未发布`

前序证据：`/Users/li/Documents/过程文档/smart-home-service-provider-v13-20260714/测试报告.md`、`/Users/li/Documents/过程文档/smart-home-service-provider-v13-20260714/问题台账.md`

被替代计划：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v13-20260714.md`

## 1. 派活卡与最小闭环

| 维度 | v14 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，只验 `smart-home-owner-pack-v0`；环保结果不冻结本计划。 |
| 我要做的事 | 跳过市场交易链，从本仓净制品直接安装，以一个智能家居主项目贯穿启用、完整候选流程、停用、历史查询、重新启用、重启、卸载、卸载后历史查询和同版本重装。 |
| 真实客户 / 场景证据 | v13 已通过真实 GUI 市场安装、项目候选、停用/急停零副作用和恢复；未覆盖双 Tab、断网、重复点击、三视口，且已打开的顶部急停浮层未随安全核心真相刷新。 |
| 最小可交付 | 同一 Pack digest、安装实例和 `SH-C01` 连续完成三阶段；每个 enabled 请求最多一次真实本地模型和一个待 Owner Candidate；disabled/uninstalled/emergency 全链零副作用；历史项目始终可查询且版本固定。 |
| 本轮明确砍掉 | `提交审核`、`模拟支付`、`Entitlement`、`下载`均记 `not_tested_owner_directed`；不启动 cloud、不登录市场、不建商品/订单/授权。生产扣款、正式发送、真实执行、Frappe 写回、生产数据、发布和完整 Dendrite 仍为 backlog。 |
| 真相源 | 本仓持 Pack/digest；14 持 lifecycle；05/06 持项目与 Run；13 持会话；01 持急停；07 持候选任务；08 持模型；03 持回执；11 持执行/Frappe。client 只重读投影。 |
| 仓库 / 层归属 | OS 负责 lifecycle/急停硬门；client 负责流式终态和安全投影刷新；本仓负责 Pack 与计划。cloud、contracts、software 本轮不改。 |
| 风险 / AI 角色 | 黄：安装、候选、启停；橙：Provider；红：急停、Gate、Receipt、发送、执行、Frappe。AI 只产 Candidate；Owner + Base 保留裁定。 |
| 契约影响 | 只改 client 内部刷新实现与测试，不改 contracts、DTO/schema、Gate、Receipt、Candidate、ReadModel、Surface 或主权边界。 |
| 上下文 / 禁止边界 | 只读本计划、Pack、隔离运行库和相关 GUI；不读真实 secret/客户数据，不改数据库，不真实发送/执行/Frappe，不正式化，不用 mock/旧证据冒充通过。 |
| 变更影响 | 顶部安全面板收到变更信号后重读 Base ReadModel；不缓存、不携带、不铸造急停真相。lifecycle 和历史数据语义不变。 |
| 生命周期档位 | 修复`已实现 -> 已接线`；本计划全程与独立复核通过后才为`已验收（打包前）`；发布另行裁定。 |

## 2. 权限、隔离与直接 lifecycle

执行者获隔离 worktree、专用 HOME/SQLite/浏览器 profile、真实本地模型、服务启停、端口、截图、只读对账和不改契约/主权链的中小问题自治修复权限。每仓独立记录 branch/HEAD/status/WIP；不得提交、合并、推送或覆盖既有改动。

1. 从 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/smart-home-owner-pack-v0` 构建净制品并登记 SHA-256、清单和 forbidden artifact 扫描。
2. 直接执行 Pack `install.py` 连接隔离 OS 真实 lifecycle；若存在本地 Pack 安装 GUI，优先用 GUI 触发同一链。安装结果必须由 14 真相、GUI 和 Receipt 三方确认。
3. 卸载使用 `uninstall.py` 的非破坏语义：撤销当前启用与运行访问，不删除 05 项目、06 历史、03 Receipt；历史项目必须继续只读查询。
4. cloud、提交审核、模拟支付、Entitlement、下载不启动、不补测、不伪造，统一写 `not_tested_owner_directed`。

## 3. 单一项目与自治推进

唯一主项目为 `SH-C01-v14-王先生灯组离线售后`。同一 `transaction_ref`、Pack v1.0.0、安装实例和浏览器 profile 贯穿 R1/R2/R3；刷新、重启、断网、停启、卸载、重装和急停都在该项目验证。只有项目真相被污染且不可恢复时可建一个 `SH-C01-RECOVERY-01`，并登记旧/新 refs 与原因；跨项目隔离仅允许只读探针。

P2/P3 及影响清晰、可逆、不改契约/主权链的中等问题必须现场修根因、补定向/邻接回归并继续。P0/P1 登记污染和冻结路径，但独立安全 lane 继续；仅当所有剩余 lane 都依赖失真状态时暂停。不得因一个问题笼统冻结后续轮次，不降低断言、不改库、不复用旧 PASS。

## 4. 开跑恢复门

1. Pack JSON/Python/结构/禁品/digest 与 `git diff --check` 全绿；registry、项目、队列、模型、Receipt 基线为空或已明确快照。
2. client 定向测试证明：安全核心成功启用/解除后发出无状态变更信号；已打开的顶部急停/系统健康面板立即重读 Base ReadModel，不能消费动作响应或保持旧文案。
3. OS 定向回归证明 disabled/uninstalled/emergency 均在 dialogue/stream/07/03/08/turn 前 fail-closed；流式 `done` 后立即收口，不残留“思考中”。
4. client 全量测试、typecheck、build、OS 相关 test/race/vet 和 Pack 审计通过后再起 GUI。

## 5. R1：直接安装、启用与项目运行全环节

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v14-20260714/`。

1. registry=0 起步，直接安装后验证 registry 0→1、Pack enabled v1.0.0、Role Pack/SlotBinding 生效、ProviderRequirement 呈现真实 `ready/not_ready/blocked/degraded`，不把 missing 写成 ready。
2. GUI 创建唯一 SH-C01，启动同一 06 Run。逐项覆盖 flow：intake、project_object、Frappe snapshot 诚实态、manager advice、milestone candidate、customer draft candidate、Frappe write candidate、Owner Gate、Receipt archive/done 的可达或受控阻断状态。
3. enabled 下分别执行一条流式和一条非流式新请求；每条最多一次真实本地模型、一个待 Owner Candidate，13/08/03/07 与 SH-C01/Pack 同源，完成态无残留 loading/streaming 卡。
4. 不真实发送、不执行、不写 Frappe、不正式化。后半流程通过候选态和 Gate 受控阻断完成覆盖，保持 0 FormalTask、0 Send、0 Execution、0 Frappe、0 正式 Memory。

## 6. R2：停用、历史查询、重新启用与重启

1. GUI 或真实 lifecycle 停用 Pack，确认 14 disabled；从项目上下文对话页攻击流式、非流式、tool confirm、旧 request、重复点击，模型、turn、Candidate、Receipt、Task、Run、Memory、Send、Execution/Frappe 业务增量全部为 0。
2. 停用态从项目清单、详情、对话历史和时间线查询同一 SH-C01：事务、Run、Candidate、Receipt 与 `pack_version_ref=v1.0.0` 仍可读，不清空、不串项目、不自动迁移。
3. 重新启用并重启 OS/client，刷新后仍恢复同一当前项目；只有全新请求恢复一次模型/候选，停用期请求不延迟补写。
4. 全局沟通中心若仍不保留项目上下文，按 SHV13-P2-003 明示“全局无项目上下文”，项目 Pack 验收固定走项目上下文对话页；不得把无上下文阻断当作停用 PASS。

## 7. R3：卸载、历史查询、重装、急停与体验

1. 执行卸载，确认当前启用指针/SlotBinding/运行访问撤销；卸载态查询 SH-C01 历史仍完整，所有新 Pack 请求零副作用阻断。重复卸载必须幂等。
2. 同 digest 重装：registry 仅一个当前版本、不复制项目/候选/Receipt；旧 SH-C01 保持 v1.0.0，新请求恢复，卸载期请求不补写。
3. 顶部急停浮层保持打开，在安全核心启用急停并返回项目页：浮层必须自动显示 `emergency_stop_enabled`；单/双 Tab、流式/非流式、旧 request、伪 receipt、双击并发全部零副作用。解除后浮层同源刷新，只有新请求恢复。
4. 补断网/恢复、浏览器刷新、390×844 / 1024×768 / 1440×900、Provider 各诚实态、作者工作台只读入口和历史回放。每个可见按钮/状态标为 PASS、受控阻断、缺口或不适用。

## 8. 证据、验收与收尾

每个关键动作同时记录 GUI、01、03、05/06、07、08、09、10、11、13、14 的前后计数、refs、截图和时间戳；只读 API/SQLite 仅用于对账。交付 lifecycle 时间线、按钮覆盖表、三视口/双 Tab/断网记录、问题台账、截图/Receipt 索引、服务与端口清理证据。独立验收主体确认开放 P0=0、阻塞 P1=0、SHV13-P2-002 GUI 关闭、单项目从直接安装到卸载/重装全程通过、停用/卸载/急停零副作用且历史可查，方可标记`已验收（打包前）`；市场四阶段只能是 `not_tested_owner_directed`，不得写 PASS。
