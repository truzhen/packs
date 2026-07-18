# 智能家居服务商 Pack 单项目全生命周期计划 v19

日期：2026-07-18
状态：`当前客户主线；修复分支门禁通过并进入最新主线后可执行；待独立验收；未发布`

依据：`docs/standards/transaction-object-pack-manifest-runtime-acceptance-standard-v1-20260717.md`、v18 最终报告 `/Users/li/Documents/过程文档/smart-home-service-provider-v18-rerun-20260718/最终测试报告.md` 与问题台账 `/Users/li/Documents/过程文档/smart-home-service-provider-v18-rerun-20260718/问题台账.md`。本计划替代 v18。

## 1. 派活卡与授权

| 维度 | v19 裁定 |
| --- | --- |
| 版本/证据 | 当前客户主线。v18 已完成 13 次 Frappe 受控写、Provider/急停/启停/历史，但 CLI 卸载缺可信 Origin 与 Owner presence，不能验收完整生命周期。 |
| 中心目标 | 用唯一项目完整验证 Transaction Object、Pack Manifest、Pack Runtime；Pack 不得腐化 Base。 |
| 最小闭环 | 直接安装→启用→项目全周期→停用→历史查询→重新启用→故障/急停→Owner GUI 卸载→同 digest 重装→历史回查；至少 10 次 Candidate→确认→执行→Receipt。 |
| 跳过 | `提交审核、模拟支付、Entitlement、下载`全部记 `not_tested_owner_directed`，不进入市场链。 |
| 真相源 | Pack/digest 在 packs；05 Transaction；06 Run；14 lifecycle/Saga；01 Gate；03 Receipt；11/Frappe Provider；client 只投影。 |
| 仓库/层 | Pack 声明；OS 通用 Runtime/Gate/Saga；client GUI；Frappe/Home Assistant 为 L2 Provider。硬件控制不另造程序；只可复用 `/Users/li/Documents/trae/home-assistant-core`。 |
| 契约/风险 | 不改 contracts。AI 仅 Proposer；真实写入必须 Owner + Base Gate + Gateway + Receipt。 |
| 禁止 | 不用生产客户/家庭设备/付款/外发；不手改 DB，不伪造 Origin/presence/Receipt，不用 API 替代标明 GUI 的步骤，不让 Pack 持 Base 主权。 |

执行者已获隔离 worktree、受控 GUI、本地测试账户、模型、Frappe、启停/卸载/重装与故障注入权限。中小问题必须现场修根因、补回归并继续；P0/P1 登记后只冻结受影响危险 lane，能安全继续的其余 lane 必须跑完。只有继续会扩大设备、数据或主权风险，或剩余断言全部依赖已污染真相时才暂停。

## 2. R0：v18 P1 解冻与三标准硬门

1. 从最新六仓 main 创建全新隔离 worktree，登记 SHA、clean、Pack digest、端口、数据库、Frappe site 与浏览器 profile。OS 起栈时把 `TRUZHEN_LOCAL_PACK_SOURCES_ROOT` 显式设为该轮最新 `truzhen-packs` worktree 根目录；不得设为用户主目录或由浏览器传路径，未配置时应诚实 `not_ready`。
2. 先验证 Owner-presence 卸载修复：`uninstall.py` 网络写请求必须为 0；脚本只展示/打开可信 GUI 交接并等待 os-14 `disabled`。从非可信 Origin、过期 presence、伪 decision 调用均 fail-closed；GUI Owner 当次确认后恰好一条 lifecycle Receipt。
3. 直接 GUI 本地安装必须使用 OS 受信来源摘要；浏览器不提交路径或安装计划。01/03/09/13/14 共用一个 reservation，安装完成才发布 enabled；中断在 1/8/15 scope、role、slot、scene、completion receipt、decision commit 时，分别验证逆序补偿、冷启动恢复和幂等重放，任何时点不得残留半装 Pack。
4. Manifest/schema、flow/role/capability 引用、namespace、forbidden artifacts 与 lifecycle 语义全绿；静态声明不得写“已验收”。
5. Frappe readiness 必须包含 Lead、Project、Task、Item 与售后/缺陷对象能力；容器启动或 tools 数量不能替代业务 ready。缺项返回 `not_ready`。
6. 跑 Pack 静态门、OS Saga/Runtime 定向与 race/EGR、client 定向/全量/typecheck/build/smoke；失败先修根因。

## 3. R1：唯一项目完整生命周期

唯一项目 `SH-C01-v19-全屋智能交付项目`，固定一个 `transaction_ref`、Pack version/digest、Frappe site 和浏览器 profile。不得新建 recovery 项目逃避缺陷。

1. 从空 registry 直接 GUI 安装并启用，核对 manifest、全 flow、RoleBinding、ProviderRequirement、version/digest、同一 reservation 与 central Receipt。
2. GUI 依次完成商机 Lead→Project 立项→里程碑/Task→物料 Item→采购/到货候选→进度更新→交付→缺陷/售后 Ticket→项目复盘；不得用 Task 冒充 Project。
3. 连续至少 10 次不同业务意图闭环。每次确认前展示项目、对象、动作、影响、Provider、幂等键与停止/回滚说明；确认前外部增量 0，确认后恰好一次 Frappe 动作和一条可反查 Receipt；重复点击/双 Tab 重放零新增。
4. 所有对象、Candidate、Owner 裁定、Gate、Gateway、Provider 与 Receipt 都挂同一 transaction timeline；Flow end 不能在交付/售后证明缺失时把 05 Transaction 写完成。
5. 覆盖全部 flow 节点和可见按钮；秘书单意图最多一次模型调用和一个 Owner Candidate；保持 0 未确认 Formal、0 外发、0 正式 Memory。

## 4. R2：Runtime 故障、启停与历史

1. Provider 停止、缺主数据、超时、断网、版本冲突、错误对象绑定分别返回稳定 `not_ready/provider_missing/blocked`，零 Candidate/Gate/Provider/Receipt 业务增量。
2. 双 Tab、重复确认、旧响应、断网恢复和 OS/Frappe 重启核验 transaction/Pack version/急停二次门控与 OCC；不得跨项目或重复外部对象。
3. 停用 Pack 后攻击秘书与全部写入口，模型、Candidate、Task、Receipt、Frappe 增量均 0；同一 SH-C01 的历史、版本和 Receipt 仍可只读查询。
4. 重新启用并重启 OS/client/Frappe，从权威 ReadModel 恢复当前项目和历史；停用期请求不延迟补写。
5. 用 390×844、1024×768、1440×900 覆盖全按钮、console/network、断连恢复和上下文保持。

## 5. R3：急停、Owner GUI 卸载重装与设备可选 lane

1. 急停下单/双 Tab、旧请求、伪 receipt、秘书与写回全部零副作用；01/03/GUI 同源，解除后只恢复新请求。
2. 在 GUI 由 Owner 当次确认卸载；RoleBinding、Provider session/轮询、运行访问撤销，新请求阻断；SH-C01、Frappe 正式对象、Receipt 与版本历史保留。CLI 仅负责交接/等待，不得代签。
3. 同 digest GUI 重装并重新启用；registry、Receipt、历史和 Frappe 对象不重复，旧项目仍固定原版本，历史查询可用。
4. 仅 Home Assistant L2 Provider 已正式接线且测试实体授权时，测试设备离线、状态漂移、重复指令、错误实体绑定、超时与人工急停；真实执行前展示设备/动作/影响/新鲜度/回滚，经 Gate/Gateway/Receipt。否则诚实记 `provider_missing/not_ready`，不阻断项目经营主链。

## 6. 放行证据

证据目录 `/Users/li/Documents/过程文档/smart-home-service-provider-v19-20260718/`。交付三标准矩阵、Saga 故障注入矩阵、Owner-presence 卸载证据、10+闭环、全节点/按钮、三视口、双 Tab/断网、Provider 状态、生命周期/历史、Receipt 索引与清理。开放 P0=0、阻塞 P1=0、GUI 与后端证据齐全，才可标 `已验收（打包前）`。
