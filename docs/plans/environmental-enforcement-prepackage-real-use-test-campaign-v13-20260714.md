# 环保执法 Pack 单主项目全程测试计划 v13

日期：2026-07-14

状态：`当前客户主线；修复已实现 -> 已接线；待同一 C01 主项目贯穿 R1/R2/R3 的真实 GUI 验收，未发布`

前序证据：`/Users/li/Documents/过程文档/env-prepack-v12-20260713/R1/r1-baseline-and-market-upload-20260713.md`、`/Users/li/Documents/过程文档/env-prepack-v12-20260713/R1/issue-register.md`

被替代计划：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v12-20260713.md`

## 1. 派活卡与最小闭环

| 维度 | v13 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，只验 `environmental-enforcement-pack-v0`；其它 Pack 不参与、不互相冻结。 |
| 我要做的事 | 使用真实本地账号和真实本地 cloud 市场取得净制品，以**一个 C01 主项目**连续跑完安装、知识、PDF、双角色、compare、停启、重启、断网、重放、急停和多视口。 |
| 真实场景证据 | v11 角色候选曾把知识常见情形冒充本案事实；v12 已完成真实登录与 C01 合成 PDF，但 Chrome 原生文件选择器需要人手选择 ZIP，测试未进入市场安装。 |
| 最小可交付 | 同一制品 digest、同一安装实例、同一 C01 事务贯穿三阶段；最终证明 45/45 源文档、15/15 scope、PDF/角色/compare 主链和停用/急停零副作用。 |
| 砍掉 / backlog | 生产上架、真实扣款、真实执法/处罚/送达、客户生产材料、生产外发、正式法律结论、完整 Dendrite。 |
| 真相源 | cloud 持商品/订单/Entitlement/digest；14 持 lifecycle；09 持 KnowledgeMount/FormalKnowledge；05 持事务事实；06 持 Run；13/08 持角色/模型；07/03 持 Candidate/Receipt。client 只投影，Pack 只声明。 |
| 仓库 / 层归属 | `truzhenos` 负责运行硬门和真相桥；client 负责真实状态展示；本仓负责 Pack 资产与计划；cloud 负责真实市场链。无 contracts 变更。 |
| 风险 / AI 角色 | 黄：候选与 lifecycle；橙：市场/Provider；红：法律正式化、Gate、Receipt、急停、真实执法。AI 是测试施工者；行业角色始终只是 Proposer。 |
| 契约影响 | 仅内部实现和测试，不改 schema/DTO、Gate、Receipt、Candidate、ReadModel、ProviderRequirement、Surface 或主权边界。 |
| 上下文 | 只读本计划、前序证据、该 Pack、受控运行数据库和对应 cloud/OS/client 页面；不读取无关客户数据。 |
| 禁止边界 | 不读真实 secret，不真实扣款，不伪造 Entitlement/安装，不改数据库，不正式化、不外发、不执行执法动作，不用 mock/fixture 冒充通过。 |
| 生命周期 | 当前`已实现 -> 已接线`；只有本计划跑完且独立复核通过，才升为`已验收（打包前）`。发布另行裁定。 |

## 2. 全权限与启动基线

执行者已获隔离 worktree、专用 HOME/SQLite/Docker/浏览器 profile、真实本地测试账号登录、真实本地 cloud 市场、本地 `MOCK_NOT_PAYABLE` sandbox 模拟支付、中小问题代码修复、构建、服务启停、端口、截图及只读对账权限；模拟支付不得产生真实扣款。

- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhenos`，实际分支/HEAD/WIP 开跑时登记。
- Pack：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs`。
- Client：优先使用包含本轮修复的 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-client-web-desktop`；开跑前登记实际分支/HEAD。
- Cloud：从干净授权分支起隔离 full-stack lane，登记端口、容器、商品 version 和制品 SHA-256。

## 3. 单主项目与尽量跑完全程纪律

1. R1 创建唯一主项目 `C01-v13-固废露天堆放线索核查`，同一 `transaction_ref`、Run 和浏览器 profile 贯穿 R1/R2/R3；不得每轮另建“新鲜 C01”逃避历史状态。
2. lifecycle 切换、OS/cloud 重启、浏览器刷新、断网恢复都在该 C01 上验证历史连续性。恢复后只发**新请求**，但不换项目。
3. 只有 C01 真相已被产品缺陷污染、删除或不可恢复时，才可建立一个 `C01-RECOVERY-01`；必须登记原因、旧/新 refs、污染边界，且最终仍以一个恢复项目跑完剩余全程。
4. 跨事务泄漏若必须对抗，可建立一个最小只读探针 `C02-ISOLATION-PROBE`；它只检查 C01 正文/ref 不泄漏，不执行 Pack 主链，也不计为第二个验收项目。
5. 测试目标是跑到所有安全 lane 结束。某一步失败不得自动终止整轮；先判断哪些后续项目内步骤仍可信并继续。

## 4. 原生文件选择器协作规则

Chrome 原生文件选择器不接受自动 `setFiles` 属受控 GUI 工具边界，不直接定级为产品 P0/P1。执行者应保留登录态和当前页面，只请求 Owner 完成一次明确动作：

- 市场 ZIP：`/Users/li/Documents/过程文档/env-prepack-v12-20260713/R1/evidence/environmental-enforcement-pack-v0.market.zip`
- C01 PDF：`/Users/li/Documents/过程文档/env-prepack-v12-20260713/R1/evidence/C01-synthetic-solid-waste-lead-20260713.pdf`

Owner 选择后从同一页面继续，不重启整轮、不用 API/脚本代传、不把协作等待写成产品失败。若选择后产品上传报错，才按真实产品问题定级和修复。

## 5. 自治修复与停线规则

- P2/P3 及影响范围清晰、可逆、不改契约/主权链的中等问题：执行者必须自行定位根因、修复、补定向与邻接回归，从同一 C01 最近可靠点继续。
- P0/P1：登记事实、真相源、污染范围、副作用、冻结路径与恢复门。身份、市场、静态审计、只读对账、响应式或其它独立安全 lane 仍可信时继续跑。
- 仅当 Gate/Receipt/权限/跨事务隔离失去可信度，且所有剩余 lane 都依赖该失真状态时暂停；报告必须列清“为什么不能继续”和重启门。
- 不降低断言、不吞错、不改库、不复用旧 PASS，不以重建项目掩盖缺陷。

## 6. 开跑前恢复门

1. Pack JSON、Python 语法、结构、forbidden artifact、bundle digest 全绿；ZIP 内 `__pycache__`、`*.pyc`、DB、日志、构建目录、路径穿越、符号链接均为 0。
2. 角色候选只可把同一 06 Run 事实切片作为本案事实；知识、persona、示例和行业常见情形只能作待核验假设。缺 Run 或 ref 不一致须在模型/Candidate/Receipt 前 `role_review_grounding_unavailable`。
3. 停用/版本漂移必须在 07/03/08/06 回灌前阻断，模型、Candidate、Receipt 增量为 0。
4. OS 定向 test/race/vet、Pack 审计、client 定向测试/typecheck/build 和 `git diff --check` 通过后再起 GUI。

## 7. R1：真实市场与 C01 核心链

过程目录：`/Users/li/Documents/过程文档/env-prepack-v13-20260714/`。

1. 真实 GUI 登录；未授权安装明确失败；经页面本地模拟支付获得 Entitlement；上传、下载、安装 digest 一致，registry 0→1。
2. 对账 45/45 源文档、15/15 scope/mount，并按 Pack/version/source 全量分页覆盖；历史全局 752 不作为本 Pack 固定契约。
3. 创建唯一 C01：PDF 前角色候选只能列证据缺口；经 GUI 上传合成 PDF 后生成全新 `enforcement_elite` 候选。逐句核对未新增废水/废气、监测结论、数值、主体或现场事实。
4. `critical_lawyer` ready 时生成质询候选；不 ready 时 GUI/06/08/03 一致诚实。流程停在 compare/Owner Gate，0 FormalTask、0 外发、0 执法执行、0 正式法律结论。

## 8. R2：同一 C01 生命周期与恢复

在同一 C01 执行 `disable -> reactivate -> disable -> reactivate`，覆盖 1/8/15 mount 中途失败补偿、OS/cloud 重启、session 过期重登、Owner/actor 差异、幂等重放和断网恢复。停用期攻击历史 advice/challenge、PDF 回灌、秘书与 Owner Confirm：Provider、模型、Candidate、Receipt、FormalTask、发送、执行增量均为 0，Run 不前移；重新启用后只允许新请求恢复，历史证据仍可回放。

## 9. R3：同一 C01 急停与体验对抗

覆盖坏/加密/伪装 PDF、重复 PDF 幂等、旧 request/伪 receipt 重放、双 Tab、刷新、390×844 / 1024×768 / 1440×900、Provider `ready/not_ready/blocked/degraded`、作者工作台只读链。急停期间 PDF/秘书/角色路径必须 0 Provider、0 stream、0 08 usage、0 Receipt、0 ModelRun、0 Task/Memory/Send/Execution；解除后仅新请求恢复，急停期请求不得延迟补写。

## 10. 独立验收与收尾

未参与实现的验收主体复核制品 hash/禁品、C01 全程时间线、逐句事实审计、05/06/07/08/09/13/14/03 增量矩阵、可选隔离探针和端口/容器清理。开放 P0=0、阻塞 P1=0、R1/R2/R3 全 PASS 才标记`已验收（打包前）`；否则按第 5 节继续安全 lane 或诚实暂停，不得发布。
