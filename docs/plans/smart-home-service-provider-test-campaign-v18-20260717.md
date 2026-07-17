# 智能家居服务商 Pack 单项目标准化全生命周期计划 v18

日期：2026-07-17  
状态：`当前客户主线；修复合并后可立即执行；待独立验收；未发布`

依据：`docs/standards/transaction-object-pack-manifest-runtime-acceptance-standard-v1-20260717.md`、v17 后端/Provider 报告 `/Users/li/Documents/truzhenv3worktree/smart-home-v17-20260717/evidence/v17-backend-provider-closeout.md` 与台账 `/Users/li/Documents/truzhenv3worktree/smart-home-v17-20260717/evidence/issues.md`。本计划替代 v17。

## 1. 派活卡

| 维度 | v18 裁定 |
| --- | --- |
| 版本/证据 | 当前客户主线。v17 已证明单项目 11 次后端受控闭环、启停、重启、Provider `not_ready`；GUI 因会话 URL 策略无证据，不得沿用为 PASS。 |
| 目标 | 以唯一项目证明 Transaction Object、Manifest 和 Runtime 三项标准；补齐真实受控 GUI，不重做市场购买链。 |
| 最小交付 | 10 次以上可回放 Candidate→Owner→Gate→Gateway→Frappe→Receipt；项目从商机到交付/售后，经历启停、历史、故障、卸载重装。 |
| 跳过 | `提交审核、模拟支付、Entitlement、下载`均记 `not_tested_owner_directed`；直接测试安装。 |
| 真相源 | Pack/digest 本仓；05 Transaction；06 Run；01 Gate；03 Receipt；11/Frappe Provider；client 只投影。 |
| 层归属 | Pack 声明；OS 通用 Runtime/Gate；client GUI；software/Frappe L2。Home Assistant 仅复用 `/Users/li/Documents/trae/home-assistant-core`，不另造硬件程序。 |
| 契约/风险 | 不改 contracts。L0 声明、L1 Agent、L2 Provider；AI 是施工者与 Proposer。 |
| 禁止 | 不用生产客户、生产家庭设备、真实付款或外发；不手改 DB、不绕 Gate、不用 API 代替要求的 GUI、不用 Task 冒充 Project。 |

执行者获本地隔离环境、测试账户、GUI、模型、Frappe 与安全测试实体所需权限。中小问题现场修复并继续；大问题登记并仅冻结相关 lane，安全 lane 跑完；只在继续会扩大真实风险或真相已污染时暂停。

## 2. R0：标准与 GUI 开跑硬门

1. 从最新主线建全新 worktree；记录 OS/client/packs/software SHA、clean、制品 digest。
2. Pack Manifest 全字段、schema、flow/role/capability 引用、禁品、namespace 和 lifecycle 语义通过；静态声明不得写已验收。
3. 启动前先在受控 Browser/Chrome 实际打开 GUI URL并截图；若策略阻断，换允许的受控会话或登记环境 blocker，但继续可独立后端 lane，不绕过安全策略。
4. Frappe readiness 必须包含 Company、Warehouse Type 等 Project 主数据和所需对象能力；仅容器/42 tools 不算业务 ready。缺 Project 能力返回 `not_ready`，禁止 Task/HD Ticket 冒充 Project。
5. 跑 Pack 审计、OS 定向/race/EGR、client 定向/全量/typecheck/build/smoke；失败先修根因。

## 3. R1：唯一项目完整生命周期与 10+ 闭环

唯一项目 `SH-C01-v18-全屋智能交付项目`，固定一个 `transaction_ref`、Pack version/digest、Frappe site 和浏览器 profile。除权威状态不可逆污染并登记外，不得另建 recovery 项目。

1. 直接空 registry 安装、启用；核对 manifest、flow 全节点、角色、ProviderRequirement、version/digest 与安装/启用 Receipt。
2. GUI 创建 SH-C01，依次完成真实业务对象：Lead→Project 立项→里程碑/Task→物料 Item→进度更新→交付→缺陷/售后 Ticket；用不同业务意图连续完成至少 10 次闭环。
3. 每次 GUI 均展示项目、对象、动作、影响、Provider、幂等键和回滚/停止说明；确认前外部增量 0，确认后恰好一次 Frappe 动作与一条可反查 Receipt。重复点击/重放零新增。
4. Transaction timeline 汇总所有 Candidate、Owner、Gate、Gateway、Frappe 与 Receipt；Flow end 不能自动把未交付项目写完成。
5. 覆盖所有 flow 节点和可见按钮；秘书每意图最多一次模型调用和一个 Owner Candidate；保持 0 未确认 Formal、0 外发、0 正式 Memory。

## 4. R2：Runtime 故障、启停与历史

1. Provider 停止、缺主数据、超时、断网、版本冲突分别返回稳定 `not_ready/provider_missing/blocked`，零 Candidate/Gate/Provider/Receipt 业务增量；恢复只处理新请求。
2. 双 Tab、重复确认、旧请求、断网前后重放核验 transaction/pack version/急停二次门控与 OCC；不得跨项目或重复 Frappe 对象。
3. 停用 Pack 后攻击秘书与全部写入口，模型、Candidate、Task、Receipt、Frappe 增量均 0；同一 SH-C01 历史、版本和 Receipt 仍可查。
4. 重新启用并重启 OS/client/Frappe；从权威 ReadModel 恢复当前项目和历史，不依赖前端内存。
5. 三视口 390×844、1024×768、1440×900，console/network、断连恢复、全部按钮逐项留证。

## 5. R3：急停、卸载重装与设备可选 lane

1. 急停下单/双 Tab、旧请求、伪 receipt、秘书和写回全部零副作用；01/03/GUI 同源，解除后只恢复新请求。
2. 卸载后运行资源、RoleBinding、Provider session/轮询撤销；新请求阻断，历史可查。相同 digest 重装后 registry/历史/Frappe 对象不重复，旧项目仍固定原版本。
3. 仅在 Home Assistant 已作为 L2 Provider 正式接线且测试实体授权时，覆盖设备离线、状态漂移、重复指令、错误实体绑定、超时和人工急停；执行前展示实体/动作/影响/新鲜度/回滚，经 Gate/Gateway/Receipt。否则记 `provider_missing/not_ready`，不阻断项目经营主链。

## 6. 放行

证据目录 `/Users/li/Documents/过程文档/smart-home-service-provider-v18-20260717/`。交付三标准矩阵、10+闭环表、全节点/按钮表、三视口、双 Tab/断网、Provider 状态、生命周期/历史、Receipt 索引、门禁与清理。开放 P0=0、阻塞 P1=0、GUI 与后端证据均完成后，只标`已验收（打包前）`。
