# 智能家居服务商 Pack 测试计划 v8（停用零副作用与秘书上下文终验版）

日期：2026-07-13
状态：`契约已定 / 修复已实现 / 待全新三轮 GUI 验收`
版本定位：当前客户主线；承接 v7，不把 v7 失败证据复用为 PASS。

## 1. 派活卡

| 维度 | v8 裁定 |
|---|---|
| 真实证据 | v7 证明 Smart Home Pack 停用后，同一 GUI 请求仍产生真实本地模型调用与 `model_usage_recorded`；07 因幂等未重复，但零副作用已失败。 |
| 最小可交付 | 新鲜项目中证明 enabled 可用；disabled 后秘书流式、非流式、动作交接、Owner Confirm 四入口均在写入前 blocked，模型/Candidate/Receipt/任务/外部副作用增量为 0。 |
| 真相源 | 05 持事务→Pack 绑定，14 持 enabled pointer/版本状态，09 持 KnowledgeMount，08 持模型 usage，03 持 Receipt，13 持秘书会话；client 只投影。 |
| 仓库归属 | OS 负责 05→14 前置门与 14→09 停用顺序；client 负责携带结构化 Pack 上下文与诚实阻断；packs 只持计划/声明；cloud 仅真实登录市场。 |
| 风险/契约 | 红：模型、Receipt、Owner Confirm、真实动作；黄：GUI 上下文。只强化实现，不改 contracts/schema/主权边界。 |
| 生命周期 | 修复为`已实现`；完成本计划新鲜 GUI 与四路对账后才升`已验收`。 |

明确非目标：生产账号、真实资金、客户数据、真实发送/Frappe 写回、contracts 修改、提交/合并/推送/发布。

## 2. 权限与问题处理

延续 Owner v7 授权：执行者可在四个隔离 worktree 和本轮专用本机基础设施内修改、构建、测试、重启/停止服务与容器，使用真实本地 seed 登录、真实本地云市场 sandbox Entitlement、真实制品下载/digest/OS 安装。中小问题必须自行修复、定向与邻接回归后继续；大问题登记污染范围和恢复门，安全独立 lane 能继续就继续，只有身份、权限、Receipt、Gate、跨事务隔离或真实副作用不可信且无安全 lane 时暂停。

授权不包含生产、真实付款/外发/执行、真实 secret、contracts、覆盖他人 WIP、提交、合并、推送或发布。

## 3. R1 首门：启停原子性与零副作用

使用 `/Users/li/Documents/过程文档/smart-home-service-provider-v8-20260713/`，全新 HOME、DB、browser profile、cloud project、Pack 安装、项目、Run、Candidate、Receipt。

1. 真实 GUI 登录与云市场安装，创建绑定 Smart Home Pack 的真实 05 项目；禁止 `tx-secretary-runtime-*` 合成事务冒充项目。
2. enabled 状态从项目页请求秘书候选：请求中的 transaction/pack/version 必须与 05、14 一致；允许 1 次受治理模型调用并形成同源 13 turn、Candidate、08 usage、03 Receipt。
3. GUI 停用：14 提交 disabled 前必须先把该版本全部 09 mount 置 disabled；任一 mount 失败时 14 record/pointer 保持 enabled，GUI 显示失败且不得出现分叉。
4. 停用成功时 GUI、14 record/pointer、09 mounts、03 lifecycle/mount receipts 四方一致；历史记录保留。
5. disabled 后分别攻击：`chat-candidate`、`chat-candidate-stream`、`conversation/turn`、`tool/confirm`。四者必须在 dialogue、Issuer.Confirm、model、candidate、receipt 前返回 `blocked`。
6. 每次攻击记录前后：13 turn、08 usage、03 Receipt、07 Candidate/FormalTask、发送、执行、Frappe 写入；全部增量必须为 0。
7. Base/14/05 任一真相不可读时返回 `not_ready`，不得降级调用模型。

## 4. 无项目与跨项目上下文

- 全局沟通中心检测到 Pack 意图但没有真实 05 项目时，不调用模型，只显示“先创建/选择项目”的生命周期引导；Pack 已停用则直接显示停用阻断。
- 从项目返回沟通中心必须由 05 ReadModel 重选 transaction；普通侧栏进入仍为无项目全局对话。
- 双项目轮换、双 Tab 与刷新后，客户端 pack_ref/version 不得覆盖 05 真相；错 Pack、错版本、伪 transaction 必须 blocked。

## 5. R2/R3 与放行

R2：disable→reactivate→disable、不同安装 owner/GUI actor、OS/cloud 重启、session 过期、Provider missing/timeout/recovery；每轮重新核对 enabled pointer 和 mount 状态。
R3：双 Tab OCC、多视口 390×844/1024×768/1440×900、历史 Receipt 回放、伪/跨项目/重放 refs、Dendrite 收窄结论及 Owner UAT。

开放 P0/阻塞 P1 为 0，R1-R3、零副作用增量表、真实登录/市场和 Owner UAT 全有新证据，才可判`已验收（打包前）`；否则`未验收 / 禁止打包放行`。
