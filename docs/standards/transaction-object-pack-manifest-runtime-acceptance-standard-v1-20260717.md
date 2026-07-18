# Transaction Object、Pack Manifest 与 Pack Runtime 验收标准 v1

日期：2026-07-17  
状态：`当前客户主线；契约不变；v18 暴露安装 Saga 与 Owner-presence 卸载 P1；待 v19 双 Pack 独立验收`

## 1. 目的与证据

本标准把智能家居 v17 的单项目 11 次受控闭环，以及环保执法 v17 的 C01 全流程、法律检索噪声、历史复盘和“Flow 完成≠案件完成”问题，转成两个 Pack 共用的打包前验收口径。最小交付不是“页面能开”或“容器健康”，而是一个真实项目在固定 Pack 版本上经历完整生命周期，并能从权威状态和 Receipt 重放。

事实归属不变：Pack 文件与 digest 归 `truzhen-packs`；安装/启停/卸载归 14；Transaction Object 归 05；SceneFlowRun 归 06；Gate 归 01；Receipt 归 03；知识归 09；Provider 执行归 11 与外部系统；client 只投影。本文不新增 contracts，不改变 Candidate、Formal、Gate 或 Receipt 主权链。

## 2. Transaction Object 标准

1. 一个客户项目只有一个稳定 `transaction_ref`。商机、立项、任务、材料、沟通、候选、执行、回执和历史查询都挂在该事务轴上；不得新建项目逃避失败。
2. 创建时固定 `pack_ref`、`pack_version_ref`、digest 与 flow spec hash。停用、重启、卸载、重装后历史仍固定原版本；升级必须走显式 migration candidate。
3. 事务状态由 05 权威维护。06 的节点或 Flow end、前端页面、Pack 文案和 Provider 状态都不能自行宣布 Transaction 完成。
4. `SceneFlow completion` 只表示本次编排已走到终点。`Transaction completion` 必须满足该业务对象声明的正式证明矩阵：所需 Formal 对象、GateDecision、Gateway 执行结果和 central Receipt 全部存在且相互绑定；缺任一项保持 `blocked/not_ready/in_progress`。
5. Timeline 是追加只读聚合：至少能回放候选产生、Owner 裁定、Gate、Gateway、Provider 结果、失败、重试、启停和版本；不得从当前状态重算历史冒充当时事实。
6. Candidate 与 Formal 分 ref、分状态、分存储语义。模型、角色和 Pack 只产 Candidate；Formal 必须由 Owner + Base Gate 产生。
7. 所有写入使用稳定业务幂等键和 OCC。重复点击、双 Tab、旧响应、断网恢复和服务重启不得重复 Candidate、Formal、Provider 动作或 Receipt。
8. 停用/卸载撤销运行访问，不删除 Transaction、FormalKnowledge 和 Receipt；跨项目、跨版本、跨 Pack 的 candidate/receipt/provider session 不得串用。

## 3. Pack Manifest 标准

Manifest 至少声明并通过权威 schema/结构校验：

- 身份：`pack_id`、`pack_ref`、`version`、`manifest_version`、`template_family`、flow/role/capability/knowledge 文件坐标与 digest。
- 主权：`person_strategy`、`formalization_requirement`、`gates`、Candidate/Formal 边界、高风险动作回 Owner + Base。
- 运行依赖：`provider_requirements`、`software_requirements`、gateway/risk/fallback；缺失必须对应稳定 `provider_missing/not_ready/blocked`，不得声明 mock 成功。
- 路由与协作：通知/命令候选/回报路由、显式多角色对照、ContextSlice 与权限最小交集。
- Receipt 与生命周期：append-only、安装/启停/卸载/重装语义、历史保留、Flow completion 与 Transaction completion 的区别及正式证明矩阵。
- 安全：隔离 namespace、forbidden artifacts、无 secret/数据库/日志/缓存/Provider 实现、无直接发送/执行。
- 产品诚实态：统一生命周期档位 `想法 -> 设计中 -> 契约已定 -> 已实现 -> 已接线 -> 已验收 -> 已发布 -> 已弃用`。静态声明、mock、单测或局部 E2E 不能写成已验收/已发布。

Manifest 只声明，不能持有运行真相。新增字段不得复制 contracts 或变成 Pack 私有 Gate/Receipt schema；消费方不认识的声明必须保持向后兼容。

## 4. Pack Runtime 标准

1. 装入前校验 manifest、引用文件、schema、digest、禁品和路径；错误 fail-closed，registry 不产生半安装记录。
2. 安装、启用、停用、卸载、重装均幂等并有 03 Receipt。跨 01/03/09/13/14 的安装必须消费同一 os-01 reservation，以 os-14 持久 Saga 日志逐步落账；任一步失败或崩溃均逆序补偿或在冷启动后恢复，不能半启用、半停用或为子资源另开 Gate。
3. 每个运行入口至少在“接收候选、持久化前、Gate 前、Gateway 执行前”重新核验 transaction-bound Pack version、急停与 Provider readiness；前端预检不能替代后端真闸。
4. KnowledgeMount 只授予可见性，不覆盖检索相关性。显式 mount refs 仍必须服从 query/as-of/辖区/权限；零相关结果应诚实 `not_ready`，不得把整库噪声塞给模型。
5. 外部候选只能回灌到同一 transaction、run、pack version 和可恢复步骤。已完成、取消、失败或停在无关 Owner Gate 的陈旧信号返回 `stale/blocked`，零状态与 Receipt 增量。
6. 缺配置、缺证据、缺 Provider、版本冲突、网络失败均明确返回稳定 reason；禁止静默 fallback、空成功、示例正文、补假事实或切换未授权 Provider。
7. 真实动作必须走 Candidate → Owner → Base Gate → Gateway → L2 Provider → central Receipt。Provider health、容器启动、前端提示不能替代业务结果与 Receipt。
8. Receipt 可反查 candidate、transaction、pack/provider version、decision、evidence、idempotency key、结果和副作用；Owner 的 approve/reject/needs_revision 刷新后可回放。
9. 服务重启、断网、超时和重放后从权威 store 恢复；client 不依赖纯 `useState` 保留历史，也不在本地自铸“已确认”。
10. 停用/卸载停止 RoleBinding、KnowledgeMount、Provider session、订阅、轮询和临时资源；重装不复制 registry、历史或外部对象。
11. 基座防腐：OS 不出现具体 Pack ref、行业对象或节点专用分支；client 不持真相；Pack 不含 Provider；software 不保存 runtime/secret；contracts 不因单 Pack 缺口被现场扩张。
12. Owner-presence 动作只能在可信 GUI Origin 内完成。CLI `uninstall.py` 只能只读发现 lifecycle、打开受控 GUI 交接并等待 os-14 真相，不得伪造 Origin、Owner presence、decision 或在脚本内串行 prepare/confirm/disable。

## 5. 三层隔离与故障断言

- L0：manifest、flow、知识与只读声明。
- L1：Agent/模型推理 sandbox，只产 Candidate。
- L2：设备、网络、本地软件、Frappe、Home Assistant、OCR 等隔离 Provider，经 Gateway/Gate/Receipt。

任何缺配置、缺 Provider、缺证据、版本冲突必须返回 `blocked/not_ready/provider_missing` 与稳定原因；受影响 lane 副作用增量为零。中小问题由测试执行者修根因、补回归并继续；P0/P1 登记并只冻结受影响危险 lane，其余安全 lane继续。仅当继续会扩大主权、法律、设备或数据风险，或所有剩余步骤依赖已污染真相时暂停。

## 6. 打包前证据矩阵

每个 Pack 至少交付：净基线与制品 digest；manifest/schema/禁品报告；唯一项目完整时间线；全节点/按钮覆盖；Candidate/Formal 对账；Provider 状态矩阵；Gate/Gateway/Receipt 索引；10 次以上受控闭环（适用 Pack）；启停/历史/卸载重装；双 Tab/幂等/断网/重启；三视口 GUI；资源清理；问题台账。开放 P0=0、阻塞 P1=0 且独立复核后，才能从`已接线`升为`已验收（打包前）`。
