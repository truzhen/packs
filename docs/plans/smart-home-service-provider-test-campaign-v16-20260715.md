# 智能家居服务商 Pack 单项目全周期测试计划 v16

日期：2026-07-15  
状态：`当前客户主线；Pack v1.1.0 与跨仓修复已实现 -> 已接线；待真实 GUI 独立验收；未发布`

前序证据：

- `/Users/li/Documents/过程文档/smart-home-service-provider-v15.1-20260715/测试报告.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v15.1-20260715/issues/big-issue-register.md`

本计划替代 `/Users/li/Documents/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v15-20260714.md`。

## 1. 派活卡与 Owner 裁定

| 维度 | v16 裁定 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线；单独验收 `smart-home-owner-pack-v0`，环保结果不冻结本计划。 |
| 真实场景证据 | v15.1 已真实跑通隔离 Frappe 的项目任务创建、更新与交付缺陷，并验证 Provider 停止、Pack 停用、卸载重装；缺口是商机、立项、物料未形成一个项目的完整周期，顶部急停缺安全核心直达入口。Owner 明确硬件控制不另造程序，可尝试复用 Home Assistant。 |
| 最小可交付 | 用唯一 SH-C01，从商机、立项、进度、物料到交付完整跑通；每阶段至少一条 Candidate → Owner/Base Gate → Execution Gateway → Frappe → Receipt → 项目时间线闭环；再完成启停、历史查询、卸载与重装。 |
| 跳过项 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部跳过并标 `not_tested_owner_directed`；不启动 cloud，不建商品、订单、授权或下载。 |
| 真相源 | 项目/Run/Gate/Receipt/Provider readiness 归 truzhenos；Frappe 持外部项目业务对象；Home Assistant 持设备实体状态；client 只展示并发候选；本仓只声明 Pack；software 只登记 Provider/软件资源。 |
| 仓库 / 层归属 | truzhenos 负责 Gate/Gateway/Receipt 与 live tool readiness；client 负责项目五阶段候选和安全入口；software 登记 Frappe/Home Assistant；本仓持 v1.1.0 flow/ProviderRequirement。contracts/cloud 不改。 |
| 风险与 AI | L0 声明只读；L1 Agent 推理在 sandbox；Frappe/Home Assistant 为 L2 隔离 Provider。AI 仅是 Proposer/施工者，不拥有确认、执行、发送或项目真相。 |
| Home Assistant 边界 | 不在 Pack、client 或 truzhenos 新造硬件控制程序。Home Assistant 仅为可选 ProviderRequirement；未实现 adapter 或未 ready 时明确 `not_ready/provider_missing`，不影响五阶段项目主链放行。 |
| 契约影响 | 不改 contracts、Gate、Receipt schema、Candidate 或跨仓 DTO；优化既有实现、UI、Pack 声明和 software Registry。 |
| 生命周期 | v1.1.0 为`已实现 -> 已接线`；真实三轮与独立复核通过后才升`已验收（打包前）`；发布另行裁定。 |

## 2. 执行授权、自治与停止规则

测试执行者获本计划内全部本地权限：建立隔离 worktree/HOME/SQLite/浏览器 profile，启停 OS/client/Frappe/Home Assistant 测试实例，直接安装、启停、卸载、重装，使用本地模拟数据执行受控 Provider 写入，截图及只读 API/SQLite/Receipt 对账，并修复不改契约/主权链的中小问题。禁止真实付款、外发、生产客户数据、生产 Frappe、真实家庭设备和未授权凭据。

小问题和中等问题必须定位根因、补测试并继续跑；P0/P1 或其它大问题登记后冻结受影响危险路径，但安全、独立 lane 必须继续。只有继续会扩大安全/数据/设备风险，或剩余步骤全部依赖污染状态时暂停。不得因一个按钮、提示或局部 Provider 问题把整个三轮提前结束，也不得为了跑通绕过 Gate、伪造 Receipt 或静默降级。

每仓从固定主仓最新 `origin/main` 建独立分支，记录 HEAD、WIP、端口、容器与命令。测试任务不自行合并、推送或发布。

## 3. 共通放行硬门：Pack 不得腐化基座

1. manifest/schema/flow/role/capabilities 通过权威校验；Pack 不复制、放宽或私造 contracts schema。
2. Candidate/Formal 在存储、接口和 UI 三层隔离；模型、角色和 Pack 只能生成 Candidate，Formal 必须由 Owner/具备资格的人确认并经 Base Gate。
3. 权限是 Owner 授权、Role 白名单、Pack、ContextSlice 和 Provider 能力的最小交集；Pack 不持凭据、不直连 Provider。
4. 缺配置、缺 Provider、缺证据、Provider 未 ready 或版本冲突必须明确 `blocked/not_ready/provider_missing`，不得静默降级、补假结果、换 Provider 或返回空成功。
5. 真实 Frappe 写入或设备动作只能经 Execution Gateway；容器启动、tools/list 或 health 不能替代 Gate 与 Receipt。
6. 每次接受、阻断和执行尝试都有可反查 Receipt，关联 candidate、decision、evidence、项目阶段、Provider 结果和副作用计数；回放不能用当前状态重算。
7. 双击、重试、重放、断网恢复、启停和卸载使用稳定幂等键；不得重复 Candidate、Formal、Receipt、Frappe 对象或设备动作。
8. 卸载后 RoleBinding、Provider session 和临时资源不可运行；历史项目与 Receipt 只读保留，数据/缓存/端口/Provider 不与其它 Pack 串包。
9. truzhenos 不得出现 smart-home/Pack ref/flow node 专用业务分支、行业 schema 副本、专用 Gate/Gateway/Receipt 或 seed；client 不成为真相源；Pack 不含 Provider/设备驱动；software 不保存 raw endpoint、secret、数据库或运行态。

## 4. 唯一项目、直接安装与证据轴

唯一项目为 `SH-C01-v16-全屋智能交付项目`。同一 `transaction_ref`、Pack v1.1.0、digest、Frappe site、浏览器 profile 和证据目录贯穿三轮；不得通过新建项目逃避缺陷。仅权威状态不可逆污染时允许一次 `SH-C01-RECOVERY-01`，必须登记旧/新 refs 和原因。

直接从净 Pack 制品安装；市场四阶段全部跳过。每个阶段使用同一证据轴：GUI 动作 → Candidate → Owner 确认 → Base Gate decision → Execution Gateway → Provider 对象 → 03 Receipt → 项目时间线/历史查询。API、SQLite、Frappe DB 只作只读对账，不能替代 GUI。

## 5. 开跑恢复硬门

1. Pack v1.1.0 JSON/Python/结构/禁品/digest 校验全绿；software Registry 只含引用和 readiness，不含 Provider 源码、凭据或运行数据。
2. OS 定向测试证明 live tools catalog 缺工具时 `not_ready` 且零 Candidate/Gate/Receipt/Provider 写入；幂等重放不要求 Provider 再在线，也不重复执行。
3. client 定向测试覆盖商机 Lead、立项 Project、进度 Task、物料 Item、交付 Ticket；候选生成不执行，Owner/Base 确认前零外部写入。
4. 顶部急停浮层能直达安全核心，但不自行启停、不自铸 evidence/decision/Receipt。
5. OS EGR、client 全量测试/typecheck/build、本仓与 software 结构验证全部通过后再启动 GUI。

## 6. R1：安装与单项目五阶段主链

1. 全新 registry=0，直接安装并启用 Pack v1.1.0；核对 registry 0→1、角色槽/绑定、flow、安装/启用 Receipt 和 digest。
2. GUI 创建 SH-C01 商机：生成 `business_create_crm_lead` 候选，确认前 Frappe Lead 增量 0；确认后恰好写入 1 个 Lead 并有 Receipt/时间线。
3. 同一 SH-C01 立项：生成 Project 候选并受控写入；项目 ref 回挂同一事务，不新建平行项目真相。
4. 进度：创建 Task，再对同一 Task 更新一次里程碑；重复点击/重放只复用幂等结果，不新增第二对象或执行。
5. 物料：创建测试 Item/物料候选，确认卡显示编码、名称、影响和隔离 Provider；确认后恰好一条写入与 Receipt。
6. 交付：创建 Helpdesk Ticket/交付缺陷候选，完成确认、写入、Receipt 和项目时间线；客户发送保持草稿，不真实外发。
7. 在项目页、沟通中心、任务、时间线与 Receipt 入口查询同一 SH-C01，五阶段关联一致；保持 0 未确认 Formal、0 外发、0 生产设备动作。

## 7. R2：Provider 故障、停用、历史与恢复

1. 停止 Frappe Provider 后分别尝试五阶段新候选；每项明确 `not_ready/provider_missing`，且零新 Candidate/Gate/Receipt/Frappe 写入。恢复后只有新请求执行，旧请求不补写。
2. 停用 Pack 后攻击秘书、五阶段按钮、旧页面、双击和双 Tab；模型、Candidate、Task、Receipt 与 Provider 调用增量为 0。
3. 停用态查询同一 SH-C01 历史；五阶段对象、版本、Receipt 与时间线只读保留。
4. 重新启用并重启 OS/client/Frappe；历史固定 v1.1.0，新候选恢复；覆盖断网前已提交/未提交边界、恢复和幂等重试。
5. 覆盖 390×844、1024×768、1440×900 与全部可见按钮；提示必须区分 `blocked/not_ready/provider_missing`。

## 8. R3：急停、卸载/重装与可选 Home Assistant

1. 由顶部入口进入安全核心启用全域急停；五阶段候选、模型、Frappe 写入和可选设备动作全部零副作用。双 Tab 与旧请求重放不得绕过；解除后仅恢复新请求。
2. 卸载后新请求受控阻断，同一 SH-C01 历史可查；重复卸载幂等。同 digest 重装后 registry 不重复、历史不复制、项目仍引用原版本。
3. 作者工作台只读核对 v1.1.0 manifest、flow、role、Frappe 与可选 Home Assistant ProviderRequirement；不提交审核、不发布。
4. Home Assistant lane 是可选增强，不是五阶段主链放行前提。若 Registry/adapter 未就绪，记录 `not_ready/provider_missing` 并继续其它测试；禁止现场在 Pack/client/OS 内写临时驱动冒充接通。
5. 若 Home Assistant L2 Provider 已真实就绪，仅使用测试实体或 Owner 明确授权的受控设备：先读状态，再展示 entity、service/action、影响和回滚/不可逆说明，经 Owner + Base Gate 后调用 Gateway 并产 Receipt。
6. 可选设备对抗至少覆盖：设备离线、状态漂移、重复指令、错误实体绑定、超时和人工急停；任一故障必须 fail-closed，不得切换设备、补假状态或无 Receipt 执行。禁止测试真实门锁、燃气、安防、生命安全或生产家庭设备。
7. 只有 Home Assistant L2 Provider 已真实就绪且测试实体获 Owner 明确授权时，才连续完成至少 10 次“设备动作 Candidate → Owner 确认 → Base Gate → Gateway 执行 → Receipt”闭环；每次执行前必须展示设备、动作与影响，10 次均须核对实体状态、幂等键和 Receipt。Provider 未就绪时诚实记录 `not_ready/provider_missing`，不得为凑数自造驱动、mock 成功或把该可选 lane 伪报为已验收，也不得据此否定已完整跑通的项目五阶段主链。

## 9. 证据、独立验收与放行

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v16-20260715/`。输出基线、参数单、五阶段证据表、三轮报告、按钮覆盖表、Provider 状态矩阵、生命周期/历史时间线、问题台账、截图/Receipt 索引和资源清理证据。

独立验收主体逐项确认：开放 P0=0、阻塞 P1=0；全部安全可继续 lane 已跑完；单一 SH-C01 五阶段均完成 Candidate→Gate→Gateway→Provider→Receipt；停用、急停、Provider 缺失和卸载均零副作用；历史查询、幂等、断网恢复、重装隔离和基座防腐通过。Home Assistant 未接通只可记可选 lane `not_ready`，不得伪报通过，也不得阻断已完整验证的项目主链。满足后仅可标 `已验收（打包前）`；市场四阶段保持 `not_tested_owner_directed`。
