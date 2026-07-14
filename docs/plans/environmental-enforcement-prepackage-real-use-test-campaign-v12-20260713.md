# 环保执法 Pack 单包跑通测试计划 v12

日期：2026-07-13
状态：`当前客户主线；R1-ENV-012 修复已实现 -> 已接线；待全新隔离三轮真实 GUI 验收，未发布`
前序证据：`/Users/li/Documents/过程文档/env-prepack-v11-20260713/R1/r1-c01-gui-execution-20260713.md`、`/Users/li/Documents/过程文档/env-prepack-v11-20260713/R1/issue-register.md`

## 1. 派活卡与 Owner 已裁定边界

| 维度 | v12 裁定与 Agent 建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线。只验 `environmental-enforcement-pack-v0`，智能家居或其它 Pack 的状态不得冻结本计划。 |
| 我要做的事 | 从真实本地 cloud 市场取得净制品，连续跑完 R1/R2/R3；重点关闭 v11 C01 角色候选把知识/人格常见情形冒充本案事实的 P1。 |
| 真实客户 / 场景证据 | v11 真实 Chrome GUI：C01 仅提供“合成固废违法线索核查”，`enforcement_elite` 真实模型候选却新增未提供的“废水/废气超标排放、现场监测数据”。PDF、Candidate、Receipt、06 compare 均有真实记录。 |
| 最小可交付 | 单 Pack 净制品 + GUI 市场安装 + C01 PDF + 两角色诚实态 + compare/Owner Gate 安全停点 + 启停/急停零副作用；三轮使用同一制品 digest。 |
| 明确砍掉 / backlog | 生产上架、真实扣款、真实执法/处罚/送达、客户生产材料、生产外发、Dendrite 全量兼容、发布。法律知识继续 `pending_human_review`。 |
| 真相源 | cloud：商品、订单、Entitlement、分发 digest；14：Pack lifecycle；09：KnowledgeMount/FormalKnowledge；05：事务事实；06：Run `StatePayload`/Evidence；13/08：角色/模型；07/03：Candidate/Receipt。client 只投影，Pack 只声明。 |
| 仓库 / 层归属 | truzhenos：手工角色入口只读加载权威 06 Run、14 lifecycle 硬门及模型事实边界；client：展示真实状态；packs：单包资产与本计划。无 contracts 变更。 |
| 风险颜色 | 绿：文档/只读；黄：候选/lifecycle；橙：市场/Provider；红：法律正式化、Gate、Receipt、急停、真实执法。AI 是施工者和测试者，但业务角色始终是 Proposer。 |
| 契约影响 | 只改内部实现与测试；不改 DTO/schema、Gate、Receipt、Candidate、ReadModel、ProviderRequirement 或 Surface 契约。 |
| 上下文范围 | 本计划、v11 证据、该 Pack 资产、cloud 市场 lane、truzhenos 03/05/06/07/08/09/13/14、相关 client 页面。不得无边界读取其它客户数据。 |
| 禁止边界 | 不读真实 secret；不改库；不复用旧事务/订单/回执/截图；不以 API/脚本伪造 GUI 市场成功；不正式化、不外发、不执行执法动作。 |
| 变更影响 | 角色建议 prompt、05→06 只读上下文桥、14 lifecycle 预检、07/03/08 副作用顺序、相关回归与诊断证据；不改变主权链。 |
| 生命周期 | 当前为`已实现 -> 已接线`；只有本计划三轮通过才升为`已验收（打包前）`，之后仍须另行发布裁定。 |

## 2. 权限、真实登录与隔离运行

执行者已获本测试所需的隔离 worktree、专用 HOME/SQLite/Docker/Chrome profile、真实本地测试账号登录、本地 cloud 市场、`MOCK_NOT_PAYABLE` sandbox 模拟支付、代码中小修复、构建、服务启停、端口、截图及只读数据库对账权限。模拟支付不得产生真实扣款。

启动基线：OS 使用 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhenos` 的 `codex/two-packs-v12-runtime-repair-20260713`；Pack 资产/计划使用 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs` 的 `codex/two-packs-v12-plans-20260713`；client 使用干净当前主线 `/Users/li/Documents/truzhen-client-web-desktop`。开跑前逐仓记录实际 HEAD 和 dirty 状态；不得从旧 v11 worktree 起栈。

每轮必须使用全新 cloud 商品 version、浏览器 profile、订单、Entitlement、OS 数据库、事务和 Run；不得安装智能家居 Pack。所有命令、端口、commit、制品 SHA-256 与截图时间写入该轮基线。

## 3. 开跑前自动恢复门

1. `TestRoleReviewCandidateUsesAuthoritativeRunGrounding`：手工角色候选必须把同一 06 Run 的 `title / knowledge_query / knowledge_uncertainty / EvidenceRefs` 送入 13/08；请求 refs 与 Run refs 去重合并。
2. 事实边界 prompt 回归必须包含：只有 06 Run 事实切片可作为本案事实；知识切片、角色人格、示例话术、行业常见情形都不是本案事实；不得新增污染介质、监测结论、数值、行为、主体或现场情况。未提供内容只能列为待核验假设/证据缺口。
3. 06 Run 缺失或 `transaction_ref / pack_version_ref` 不一致，须在 07/03/08 前返回 `role_review_grounding_unavailable`，模型、Candidate、Receipt 均为 0。
4. `TestRoleReviewCreateDisabledPackWritesNoCandidateReceiptOrModelUsage`：14 停用/版本漂移在 07/03/08/06 回灌前阻断，明确三项 false。
5. OS 定向 test/race/vet、秘书体量棘轮、模块体量门、Pack JSON/脚本/结构/禁品、client 相关 typecheck/build 与 `git diff --check` 全绿后才启动 GUI。

## 4. 测试执行自治与停线规则

- P2/P3 小问题和影响清楚、可逆、不改契约/主权链的中等问题：执行者必须自行定位根因、修复、补定向及邻接回归，从最近可靠 GUI 点继续，不得只登记后结束。
- P0/P1 大问题：立即登记事实、真相源、污染范围、已发生副作用、冻结路径、恢复门。若身份、市场、静态审计、只读对账或其它独立 lane 仍可信，继续跑安全 lane；只有 Gate/Receipt/权限/跨事务隔离不可信且没有安全 lane 时暂停。
- 修复、服务重启或 lifecycle 状态变化后，受影响结论必须用全新事务/lane 重取证；旧失败证据只追加、不删除、不改库。
- 禁止 mock/fixture 成功替代真实链、降低断言、吞错、fallback 占位、复用旧 PASS 或脚本/API 注入 Entitlement/安装状态。

## 5. 四路证据与事实污染判定

每个角色候选同时保存：GUI 原文、06 Run JSON/SQLite、08 模型请求摘要与 usage、03/07 Candidate/Receipt。对 C01 输入建立“允许事实表”和“明确未提供表”，由独立验收主体逐句审计候选：

- 允许直接陈述：标题、`knowledge_query`、`knowledge_uncertainty`、PDF 已解析且有 Evidence ref 的原文事实。
- 只允许作为待核验：未有 Evidence 的物料性质、位置、主体、时间、污染介质、监测结论与数值。
- 禁止冒充本案事实：知识库案例、角色 persona 示例、行业常见情形、模型常识。
- 发现“废水/废气超标、现场监测数据”等无输入断言即 R1-ENV-012 复现；冻结该候选正式化，但继续安全只读 lane。

## 6. R1：净制品、真实市场与 C01 核心链

过程目录：`/Users/li/Documents/过程文档/env-prepack-v12-20260713/`。

1. canonical builder 以净模式生成 ZIP 和 SHA-256；根含 manifest，路径穿越、符号链接、`__pycache__`、`*.pyc`、DB、日志、构建目录为 0。上传、下载、安装 digest 一致。
2. GUI 真实登录本地测试账号；先证明未授权安装失败，再点击本地 sandbox 模拟支付，取得 Entitlement 并从真实市场点击安装，registry 0→1。
3. 对账 45/45 源文档、15/15 scope/mount、09 按 Pack/version/source 全量分页覆盖；不把历史全局 752 当单 Pack 固定契约。
4. 新建 C01，输入仅含固废露天堆放线索及明确不确定性。PDF 前先生成一次角色候选：只能输出证据缺口，不能产生具体污染事实。随后经 GUI 上传/受控文件选择器解析合成 PDF，再生成全新角色候选。
5. `enforcement_elite` 候选逐句通过第 5 节事实审计；`critical_lawyer` 若真实 Provider ready 则生成质询候选，若 `model_not_ready` 则 GUI/06/08/03 必须一致诚实，不得伪造双角色 PASS。
6. 流程安全停在 compare/Owner Gate；0 FormalTask、0 外发、0 执法执行、0 正式法律结论。

R1 通过门：市场/制品/知识/PDF/事实锚定/角色诚实态/compare 与零正式副作用均有四路证据。

## 7. R2：生命周期、停用零副作用与恢复

仅对环保 Pack 执行 `disable -> reactivate -> disable -> reactivate`，覆盖第 1/8/15 mount 中途失败补偿、OS/cloud 重启、session 过期重登、Owner/actor 差异、幂等重放。

停用后从 GUI 攻击历史 C01 advice/challenge、PDF 回灌、秘书及 Owner Confirm；每条路径模型、Candidate、Receipt、FormalTask、发送、执行增量均为 0，Run 不前移。重新启用后仅全新请求恢复，历史回执保持可查。

## 8. R3：急停、隔离与对抗

覆盖 C01→C02/C08 跨事务隔离、坏/加密/伪装 PDF、重复 PDF 幂等、双 Tab、刷新、390×844 / 1024×768 / 1440×900、Provider `ready/not_ready/blocked/degraded`。急停期间 PDF/秘书/角色路径必须 0 Provider、0 stream、0 08 usage、0 Receipt、0 ModelRun、0 Task/Memory/Send/Execution；解除后仅新请求恢复。

## 9. 独立验收与收尾

由未参与实现的专门验收主体复核制品 hash/禁品、四路逐句事实审计、05/06/07/08/09/13/14/03 增量矩阵、跨事务攻击和端口/容器清理。开放 P0=0、阻塞 P1=0 且 R1/R2/R3 全 PASS，才标记`已验收（打包前）`；否则按第 4 节继续安全 lane 或诚实暂停，不得发布。
