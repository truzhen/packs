# 智能家居服务商 Pack 单包跑通测试计划 v12

日期：2026-07-13
状态：`当前客户主线；SHV11-P0-001 修复已实现 -> 已接线；待全新隔离三轮真实 GUI 验收，未发布`
前序证据：`/Users/li/Documents/过程文档/smart-home-service-provider-v11-20260713/最终测试报告.md`、`/Users/li/Documents/过程文档/smart-home-service-provider-v11-20260713/问题台账.md`

## 1. 派活卡与 Owner 已裁定边界

| 维度 | v12 裁定与 Agent 建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线。只验 `smart-home-owner-pack-v0`；环保 Pack 的结果不得冻结本计划。 |
| 我要做的事 | 从真实本地 cloud 市场安装单 Pack，连续跑完 R1/R2/R3；重点关闭 v11 急停后 07 队列仍新增 TaskCandidate 的 P0。 |
| 真实客户 / 场景证据 | v11 真实 GUI：顶部急停已生效，页面提示未生成候选，但 07 队列从 6 增至 7；危险正式化/发送/Frappe 路径随后冻结。 |
| 最小可交付 | 单 Pack 净制品 + GUI 市场安装 + 双项目候选隔离 + 启停恢复 + 急停下非流式/流式零 TaskCandidate、零模型、零回执、零 turn、零外部副作用。 |
| 明确砍掉 / backlog | 生产账号/扣款、客户生产数据、真实外发、Frappe 真实写回、生产部署、上架发布、完整 Dendrite 兼容。 |
| 真相源 | cloud：商品/订单/Entitlement/digest；14：Pack lifecycle；05/06：项目/Run；13：秘书对话编排；01：急停；07：任务候选；08：模型；03：回执；11：执行/Frappe。client 仅投影。 |
| 仓库 / 层归属 | truzhenos：01 急停真相注入 13，dialogue/07/03/08 前置硬门并在 07 路由防御复核；client 已消费后端诚实 blocked；packs 只持声明与本计划。 |
| 风险颜色 | 绿：文档/只读；黄：候选/lifecycle；橙：市场/Provider；红：急停、Gate、Receipt、真实发送/执行/Frappe。AI 只能产 Candidate。 |
| 契约影响 | 只改内部实现和测试；不改 contracts、DTO/schema、Gate、Receipt、Candidate、ReadModel 或 Surface 契约。 |
| 上下文范围 | 本计划、v11 报告/台账、智能家居 Pack、cloud 市场 lane、truzhenos 01/03/05/06/07/08/11/13/14 与相关 client 页面。 |
| 禁止边界 | 不读真实 secret；不生产支付；不 API 注入 Entitlement；不真实发送/执行/Frappe 写回；不把浏览器存储当项目真相；不复用旧 DB/项目/回执/截图。 |
| 变更影响 | 秘书非流式/流式入口、07 intake 顺序、急停响应与零增量诊断；不改变业务主权链。 |
| 生命周期 | 当前为`已实现 -> 已接线`；三轮真实 GUI 全通过后才升为`已验收（打包前）`，发布另行裁定。 |

## 2. 权限、真实登录与隔离运行

执行者已获隔离 worktree、专用 HOME/SQLite/Docker/Chrome profile、真实本地测试账号登录、本地 cloud 市场、`MOCK_NOT_PAYABLE` sandbox 模拟支付、代码中小修复、构建、服务启停、端口、截图与只读数据库对账权限。模拟支付不得产生真实扣款。

启动基线：OS 使用 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhenos` 的 `codex/two-packs-v12-runtime-repair-20260713`；Pack 资产/计划使用 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs` 的 `codex/two-packs-v12-plans-20260713`；client 使用干净当前主线 `/Users/li/Documents/truzhen-client-web-desktop`。开跑前逐仓记录实际 HEAD 和 dirty 状态；不得从旧 v11 worktree 起栈。

每轮使用全新商品 version、profile、订单、Entitlement、OS 数据库、项目与 Run；单包 lane 禁止安装环保 Pack。登记所有分支/commit、端口、制品 SHA-256 和运行参数。

## 3. 开跑前自动恢复门

1. 非流式 `/v3/secretary/runtime/chat-candidate` 与流式 `/v3/secretary/runtime/chat-candidate-stream` 在解析合法请求后、任何 dialogue 分类/stream header/07 intake/03/08/turn/发送候选之前读取 01 Base 急停真相。
2. `TestSecretaryChatCandidateBlocksEmergencyBeforeDialogueTaskAndModel` 与流式同类测试：返回 423、`base_emergency_stop_active`，并明确 `model_call_performed / candidate_created / task_candidate_created / receipt_created / turn_created / external_send_candidate_created=false`。
3. `TestSecretaryEmergencyStopKeepsTaskQueueAtZeroDelta`：经真实 Base enable-candidate 启用急停，使用 v11 原攻击语句分别打非流式/流式，07 全队列 before=after；不得只统计某个 zone。
4. 07 路由缝必须再读同一急停 guard，防未来内部调用绕开 HTTP 入口；急停真相读取失败也 fail-closed。
5. 既有 Pack 停用硬门、项目刷新恢复、双项目隔离回归必须保持；OS 定向 test/race/vet、秘书体量棘轮、模块体量门、client Vitest/typecheck/build、Pack 静态审计及 `git diff --check` 全绿后才启动 GUI。

## 4. 测试执行自治与停线规则

- P2/P3 小问题和范围明确、可逆、不改契约/主权链的中等问题：执行者自行修根因、补定向与邻接回归，从最近可靠 GUI 点继续，不得只登记后停线。
- P0/P1 大问题：登记等级、真相源、污染范围、副作用、冻结路径和恢复门。市场、身份、只读对账或其它独立 lane 可信时继续跑；只有 Gate/Receipt/权限/跨项目隔离不可信且无安全 lane 时暂停。
- 代码修复、服务重启或 lifecycle 变化后，用全新项目/lane 重取受影响证据；旧失败只追加、不删除、不修库。
- 禁止 mock 成功、脚本/API 注入代替 GUI、吞错、降低断言、复用旧 PASS、fallback 冒充完成。

## 5. 急停零副作用证据矩阵

每次急停攻击前后同时记录 GUI、01 ReadModel、13 turn/session、07 全队列、08 usage/ModelRun、03 Receipt、05/06 Candidate/Run、09 Memory、10 Send、11 Execution/Frappe。必须逐项计算精确 delta：

| 指标 | 急停期允许增量 |
| --- | ---: |
| Provider/stream/model calls/tokens/ModelRun | 0 |
| 13 turn / 对话候选 | 0 |
| 07 任意 zone TaskCandidate / FormalTask / dispatch | 0 |
| 03 candidate/model/task/send/execute Receipt | 0 |
| 05/06 Candidate / Run 推进 | 0 |
| 09 FormalMemory/MemoryCandidate | 0 |
| 10 外发尝试 / 11 执行 / Frappe 写回 | 0 |

仅允许 01 自身急停启用的治理事件/回执；必须单独列出，不能混入业务副作用分母。页面“未生成候选”必须与该矩阵完全一致。

## 6. R1：真实市场与业务核心链

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v12-20260713/`。

1. canonical builder 生成净 ZIP/SHA-256；cloud 上传前检查 manifest、路径穿越、符号链接和禁品。新商品/version、新 Chrome profile；GUI 先证明未授权阻断，再完成本地模拟支付、Entitlement、下载/安装；digest 一致、registry 0→1。
2. 对账 Pack enabled、Role Pack/SlotBinding 与 ProviderRequirement 诚实态。创建 A/B 两个智能家居项目和各自 06 Run，正文/ref/状态不得互串。
3. enabled 下每项目秘书请求恰好一次真实本地模型，只新增一个待 Owner Candidate；13/08/03/07 与 transaction/Pack 同源；0 FormalTask、0 Send、0 Frappe。
4. 刷新与双 Tab 后仍由 05 当前授权项目真相恢复各自 transaction；失权/删除/不存在则清除选择提示，不擅自选第一项。

## 7. R2：停用零副作用与生命周期恢复

GUI 停用并确认 14 disabled 后，攻击无项目领域短语、已选项目、非流式、流式、conversation turn、tool confirm、旧 request 重放。全部须在模型前 blocked/not_ready；模型、turn、Candidate、Receipt、Task、Send、Execution、Frappe 增量为 0。

随后只对智能家居 Pack 执行 `reactivate -> disable -> reactivate`，覆盖 OS/cloud 重启、登录过期重登、Owner/actor 差异、幂等、Provider timeout/recovery。重新启用后仅全新请求恢复，停用请求不得延迟补写。

## 8. R3：急停、并发与体验对抗

1. GUI 启用急停，01 ReadModel 与顶部状态一致后记录第 5 节基线。
2. 使用 v11 原攻击语句依次执行非流式、流式、双击并发、双 Tab、旧 request/伪 receipt 重放；每次均核对完整零增量矩阵，07 全队列不得从 N 变 N+1。
3. 急停中不得点击正式化、发送、执行或 Frappe；验证入口阻断已足够，不以真实危险动作探测。
4. 由 GUI 显式解除急停；只有解除后的全新请求可以恢复一次模型/候选链，急停期请求不得补写。
5. 补测断网/恢复、390×844 / 1024×768 / 1440×900、Provider `ready/not_ready/blocked/degraded`、历史回放与 Owner UAT。跨项目正文/ref 泄漏立即按 P0 冻结污染路径。

## 9. 独立验收与收尾

由未参与实现的专门验收主体复核自动门、真实 GUI 证据、市场/Entitlement/digest、急停完整零增量矩阵、停用竞态、项目隔离和端口/容器清理。开放 P0=0、阻塞 P1=0 且 R1/R2/R3 全 PASS，才标记`已验收（打包前）`；否则按第 4 节继续安全 lane 或诚实暂停，不得发布。
