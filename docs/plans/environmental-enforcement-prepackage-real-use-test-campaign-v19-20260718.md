# 环保执法 Pack 单项目全生命周期计划 v19

日期：2026-07-18
状态：`当前客户主线；修复分支门禁通过并进入最新主线后可执行；待独立验收；未发布`

依据：`docs/standards/transaction-object-pack-manifest-runtime-acceptance-standard-v1-20260717.md`、v18 跨链问题 `/Users/li/Documents/过程文档/env-prepack-v18-20260717/issues/P1-06-local-pack-child-lifecycle-saga-20260718.md` 与设计实施记录 `/Users/li/Documents/过程文档/env-prepack-v18-20260717/p1-05-gate-receipt-cross-chain-atomicity-recovery-design-20260718.md`。本计划替代 v18。

## 1. 派活卡与授权

| 维度 | v19 裁定 |
| --- | --- |
| 版本/证据 | 当前客户主线。v18 正确冻结旧加载器：Role、Slot、Knowledge 各自开 Gate，不能构成同一可恢复原子安装链。 |
| 中心目标 | 用唯一 C01 完整验证 Transaction Object、Pack Manifest、Pack Runtime；补齐共享 reservation、持久 Saga、补偿与恢复，Pack 不腐化 Base。 |
| 最小闭环 | 直接安装→启用→C01 全候选生命周期→停用→历史查询→重新启用→故障/急停→Owner GUI 卸载→同 digest 重装→历史回查；红色正式动作始终为 0。 |
| 跳过 | `提交审核、模拟支付、Entitlement、下载`全部记 `not_tested_owner_directed`，不进入市场链。 |
| 真相源 | Pack/digest 在 packs；14 lifecycle/Saga；05 Transaction；06 Run/材料矩阵；09 法源；01 Gate；03 Receipt；client 只投影。 |
| 契约/归属 | 不改 contracts。OS 只实现通用 Saga/Runtime/query；Pack 声明行业 flow、证据与完成证明；不得新增环保专用 Base 分支。 |
| 风险 | L0 manifest/知识，L1 角色推理，L2 OCR/模型/网络。AI 仅 Proposer；处罚、送达、删除、对外提交、执法执行均为红色动作。 |
| 禁止 | 不读真实案件/凭据，不补假事实，不手改 DB，不绕 Gate，不用当前法替代案发时法，不让 Pack 或 Client 持 Base 真相。 |

执行者获隔离 worktree、HOME、SQLite、受控 GUI、本地模型、合成 PDF、直接安装/启停/卸载/重装与故障注入权限。中小问题必须现场修根因、补回归并继续；P0/P1 登记后只冻结受影响危险 lane，其余安全 lane 必须跑完。只有继续会扩大法律、数据或主权风险，或剩余断言全部依赖污染真相时才暂停。

## 2. R0：v18 P1 解冻与三标准硬门

1. 从最新六仓 main 创建全新隔离 worktree，登记 SHA、clean、Pack digest、端口、数据库和浏览器 profile。OS 起栈时把 `TRUZHEN_LOCAL_PACK_SOURCES_ROOT` 显式设为该轮最新 `truzhen-packs` worktree 根目录；不得设为用户主目录或由浏览器传路径，未配置时应诚实 `not_ready`。
2. 从空 registry 执行 GUI 本地安装：来源路径只能由 OS 配置，Client 只提交服务器 artifact digest 与 Base proof。01/03/09/13/14 必须共用一个 reservation；role enable、slot binding、15 KnowledgeMount、FormalKnowledge、scene enable 与 completion receipt 逐步写入 os-14 持久 Saga。
3. 在 1/8/15 scope、role、slot、knowledge formalization、scene enable、completion receipt、decision commit 各注入一次失败/崩溃；重启后必须前向恢复或逆序补偿，重复请求幂等。任何时点 registry、RoleBinding、KnowledgeMount、FormalKnowledge、Receipt 不得半装、串 reservation 或跨 Pack。
4. 验证 Owner-presence lifecycle：CLI `uninstall.py` 只读并交接可信 GUI；伪 Origin、过期 presence、伪 decision 均 fail-closed；Owner GUI 当次确认后才产生停用/卸载 Receipt。
5. 跑 manifest/schema、JSON/Python、45 源文档/15 scope 全量分页、禁品、namespace；Manifest 明确 `SceneFlow completion != Transaction completion`。
6. 跑 OS Saga/Runtime/query 定向与 race/EGR、client 定向/全量/typecheck/build/smoke、Pack 静态门；失败先修根因。

## 3. R1：唯一 C01 完整候选生命周期

唯一项目 `C01-v19-固体废物露天堆放线索核查`，固定 transaction、Pack version/digest、浏览器 profile 与证据目录；不得新建项目逃避缺陷。

1. 安装完成后核对 2 Role Pack、2 SlotBinding、15 mounts、45 源文档、version/digest、单一 reservation 与 central Receipt；未完成 Saga 不得对 Runtime 可见。
2. GUI 创建 C01。PDF 前材料缺失/待核验时返回 `role_review_evidence_not_ready/EvidenceGapCandidate`，模型/Candidate/Task/Receipt/外发/执行增量 0。
3. 上传合成 PDF，补齐主体、辖区、物料性质、现场记录、监测检测、许可台账；上传不自动等于 verified，冲突继续阻断。
4. 材料齐备后各生成一次执法精英与挑剔律师候选。每个事实必须引用 evidence；法律检索绑定案发 as-of、辖区、法源版本与有效性。注入无关辐射/水/案例知识时上下文和引用为 0；零相关结果返回 `not_ready`，不得凭模型常识补法条。
5. 对材料/文书候选执行 Owner `needs_revision`；刷新和重启后从 06/07/03 回放 PDF spine、候选、裁定、引用和 Receipt。
6. 形成处置、补证和文书草稿候选并走最终 Owner/Base Gate，但不正式处罚、送达、删除、对外提交或执法执行；敏感信息在 UI、模型 ContextSlice、日志和 Receipt 公共元数据中脱敏。
7. SceneFlow end 只表示候选链收束；因执行、送达、整改复查正式 Receipt 为 0，05 Transaction 必须保持未完成并列缺口。

## 4. R2：Runtime 故障、启停与历史

1. 缺 Provider、模型超时、断网、版本冲突、重复 PDF、旧请求、错误 transaction 与双 Tab 分别返回稳定 `blocked/not_ready/provider_missing`，零重复候选/模型/Receipt。
2. 停用后攻击 PDF、角色、文书和旧请求，03/06/07/08/09 业务增量 0；同一 C01 的材料矩阵、候选、needs_revision、Gate 与 Receipt 可只读回放。
3. 重新启用并重启 OS/client，只恢复新请求；停用期请求不延迟补写，C01 固定原版本。
4. 390×844、1024×768、1440×900 覆盖全按钮、console/network、断网恢复；原始 refs 进入技术详情，人话状态准确。

## 5. R3：急停、Owner GUI 卸载重装与完成证明攻击

1. 急停下 PDF、角色、旧请求、伪 receipt、双 Tab 全部零副作用；01/03/GUI 同源，解除后只恢复新请求。
2. Owner 在可信 GUI 当次确认卸载；运行访问、RoleBinding、KnowledgeMount、session/临时资源撤销；C01、FormalKnowledge、Receipt 与历史保留。CLI 只能交接/等待，不能代签。
3. 同 digest GUI 重装并启用；registry、mount、知识、Receipt 和历史不重复，C01 仍固定原版本且可查询。
4. 用 Flow end、前端状态、候选 Receipt、伪 execution/delivery ref 或 Provider health 宣告 Transaction 完成，均须 fail-closed；只有 05 可依据可反查正式证明矩阵改变完成态。

## 6. 放行证据

证据目录 `/Users/li/Documents/过程文档/env-prepack-v19-20260718/`。交付三标准矩阵、Saga 故障注入与恢复表、Owner-presence 卸载证据、query 相关性/零结果、事实法源逐句审计、材料与裁定回放、全节点/按钮、启停/卸载重装、双 Tab/断网/三视口、Receipt 索引与清理。开放 P0=0、阻塞 P1=0、无关知识 0、事实/法源完全锚定、Flow/Transaction 完成语义正确，才可标 `已验收（打包前）`。
