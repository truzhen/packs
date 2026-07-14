# 环保执法 Pack 单项目完整 lifecycle 测试计划 v15

日期：2026-07-14  
状态：`当前客户主线；ENV-V14 两项 P1 已实现 -> 已接线，待真实 GUI 验收；未发布`

前序证据：

- `/Users/li/Documents/过程文档/env-prepack-v14-20260714/R1/r1-recovery-gate-and-c01-gui-20260714.md`
- `/Users/li/Documents/过程文档/env-prepack-v14-20260714/R1/issue-register.md`

被替代计划：`/Users/li/Documents/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v14-20260714.md`

## 1. 派活卡与 Owner 裁定

| 维度 | v15 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，单独验收 `environmental-enforcement-pack-v0`；智能家居结果不冻结本计划。 |
| 我要做的事 | 从最新三仓 `main` 建全新隔离分支，以一个 C01 项目贯穿直接安装、启用、19 节点运行、停用、历史查询、重新启用、卸载、卸载后历史查询和同版本重装。 |
| 真实客户 / 场景证据 | v14 真实 GUI 证明 PDF 前基座能零副作用阻断，但页面丢失具体原因；PDF 后真实模型把未提供的采样、影像、证件、监测和检测写成本案事实。两项均有过程报告与问题台账。 |
| 最小可交付 | 同一 Pack v1.0.0、同一 digest、同一安装实例、同一 `transaction_ref` 完成可回放 lifecycle；PDF 前显示可解释阻断；PDF 后被接受事实必须逐字来自 06 案件切片，法条只来自 09 实际检索 ref，未知项全部标 `pending_human_review`。 |
| 本轮砍掉项 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部跳过并记 `not_tested_owner_directed`；不启动 cloud、不建商品/订单/授权，不把这四项写成 PASS。 |
| 真相源 | Pack/digest 归本仓；lifecycle 归 14；项目归 05；Run 与案件事实切片归 06；知识检索归 09；模型调用归 08；候选与任务归 06/07；回执归 03；client 只展示。 |
| 仓库 / 层归属 | `truzhenos` 负责事实锚定验收门；client 负责受控拒绝投影；本仓负责 Pack 资产和计划。contracts、cloud、software 不改。 |
| 风险 / AI 角色 | 绿：只读与展示；黄：直接安装、启停、候选；橙：Provider；红：正式法律裁定、Gate、Receipt、发送、执行。AI 仅作施工者和 Proposer。 |
| 契约影响 | 零 contracts/schema/DTO/Gate/Receipt/Candidate/ReadModel/ProviderRequirement 变更；只改 OS/client 内部实现与测试。 |
| 上下文 | 只读本计划、Pack 文件、v14 报告、合成 C01、隔离 HOME/SQLite、GUI 与对应诊断；不得读取真实客户资料或真实凭据。 |
| 禁止边界 | 不改库伪造状态，不正式化，不真实送达，不执行执法动作，不用 prompt 声明、旧截图、mock 或 API-only 冒充 GUI PASS。 |
| 变更影响 | 06 模型输出接受条件、client 错误展示、单 Pack lifecycle 与覆盖证据；历史项目和既有 Receipt 不改写。 |
| 生命周期档位 | 修复为`已实现 -> 已接线`；本计划与独立复核通过后才升`已验收（打包前）`；发布另行裁定。 |

## 2. 执行权限、隔离和自治

测试执行者已获本计划范围内全部本地权限：建立隔离 worktree、专用 HOME/SQLite/浏览器 profile，启停 OS/client，本地真实模型调用，使用合成 PDF，直接安装/停用/启用/卸载/重装，截图、只读 API/SQLite/Receipt 对账，以及修复不改契约/主权链的中小问题。不得触发真实市场、付款、发送、执法执行、生产 Frappe 或客户数据。

每仓必须从固定主仓最新 `origin/main` 拉分支，独立记录 branch、HEAD、`git status --short --branch`、既有 WIP、端口和验证命令。中小问题必须定位根因、补测试、从最近可靠点继续；P0/P1 登记并冻结受影响路径，但其它独立安全 lane 继续。只有所有剩余步骤都依赖已污染状态或继续会扩大风险时才暂停。测试任务不得自行 merge/push，除非 Owner 另行授权。

## 3. 唯一项目和直接安装口径

唯一主项目：`C01-v15-固体废物露天堆放线索核查`。同一 `transaction_ref`、Pack v1.0.0、digest、安装实例、浏览器 profile 和证据目录贯穿 R1/R2/R3。不得新建项目来清空缺陷；只有 05/06 权威状态被不可逆污染时允许建立一次 `C01-RECOVERY-01`，并登记旧/新 refs、污染范围和原因。

1. 从 `/Users/li/Documents/truzhen-packs/environmental-enforcement-pack-v0` 构建净制品，登记清单、SHA-256、知识条目/scope 对账和 forbidden artifacts。
2. 跳过市场四阶段，直接以 `install.py` 调隔离 OS 真实 lifecycle；若已有本地文件安装 GUI，优先用 GUI 触发同一端点。
3. `uninstall.py` 只撤销当前启用与运行访问，不删除 05 项目、FormalKnowledge、历史 Candidate 或 03 Receipt。
4. 每次 lifecycle 动作必须由 GUI、14 ReadModel 和 03 Receipt 三方对账；脚本退出码不能单独算 PASS。

## 4. 开跑恢复硬门

1. Pack JSON/Python/结构审计、45 条知识与 15 scope 对账、路径穿越/符号链接/禁品扫描、digest、`git diff --check` 全绿。
2. OS 定向测试证明越界模型输出在角色候选接受前返回 `role_advice_grounding_violation`，不保存 `Reasoning` 或角色引用 Receipt；真实模型调用已经发生时，08 usage Receipt 必须诚实保留，不能写成零调用。
3. client 定向测试证明 423 响应中的 `role_review_evidence_not_ready`、`missing_fields=[pdf_parse_status]` 和模型/候选/回执三项 false 不丢失、不被泛化错误覆盖。
4. client 全量测试、typecheck/build，OS 定向/race/vet 与全量 EGR，本仓文档/结构验证全绿后才启动 GUI。

## 5. R1：安装、启用和 C01 全节点运行

过程目录：`/Users/li/Documents/过程文档/env-prepack-v15-20260714/`。

1. 全新 registry=0、KnowledgeMount=0；直接安装后对账 registry 0→1、enabled v1.0.0、2 Role Pack、2 SlotBinding、15/15 mount、45/45 知识和安装 Receipt。
2. GUI 创建唯一 C01 与 06 Run，保存 03/05/06/07/08/09/13/14 基线。PDF 前分别点 advice/challenge：页面必须明确显示 `role_review_evidence_not_ready`、缺 `pdf_parse_status`、模型/候选/回执均 0；不得产生未处理前端异常或泛化失败。
3. GUI 上传合成 PDF，核对 PDF Candidate/Receipt、`pdf_parse_status=parsed` 回灌同一 Run。生成执法精英与挑剔律师候选：每条确认事实必须能在 06 `knowledge_query/knowledge_uncertainty` 中逐字定位；法律候选仅显示 09 实际返回的 `knowledge_ref`；未知问题均带 `pending_human_review`；不得出现未提供的主体、采样、影像、证件、监测、检测、污染介质、数值、违法认定或处置决定。
4. 主动投喂一次会诱导编造的合成提示。若模型不遵守严格格式，系统必须 `role_advice_grounding_violation` fail-closed；不能为了“跑通”降级 validator、拼模板或接受自然语言。
5. 按 flow 节点表逐项覆盖 19 节点：`intake`、`screening`、`inspect_prep`、`credentials`、`onsite`、`scan`、`evidence_chain`、`evidence_review`、`fact_advice`、`legal_advice`、`lawyer_challenge`、`compare`、`disposition_shadow`、`disposition`、`correction_task`、`doc_draft`、`owner_gate`、`archive`、`done`。每项标 `PASS / 受控阻断 / not_ready / not_tested`，不得因后段高风险动作不执行就漏记节点。
6. `compare` 及后段只验证 Candidate、诚实 Provider 状态和 Gate 阻断，不作 Owner 正式批准；保持 0 FormalTask、0 文书外发、0 执法执行、0 生产写回。

## 6. R2：停用、历史查询、重新启用和重启

1. 在同一 C01 停用 Pack；14=disabled、15 mount inactive。攻击角色候选、PDF 重放、旧 request、刷新和重复点击，03/06/07/08/09 业务增量均为 0。
2. 停用态从项目列表、详情、时间线、候选和 Receipt 入口查询同一 C01；标题、Run、PDF、候选、Receipt、Pack v1.0.0 均只读可回放，不删除、不迁移、不串项目。
3. 重新启用同版本并重启 OS/client；历史 C01 仍固定 v1.0.0，仅全新请求可恢复一次候选，停用期请求不得延迟补写。
4. 覆盖双 Tab、断网/恢复、重复点击、旧页面刷新和 390×844 / 1024×768 / 1440×900；逐项记录全部可见按钮。

## 7. R3：卸载、历史查询、重装和作者/Provider 覆盖

1. 卸载后 registry/enable 指针、Role Slot/Binding、KnowledgeMount 均不可运行；重复卸载幂等。卸载态继续查询同一 C01 历史，所有旧 refs 与时间线保持可见；新请求零副作用阻断。
2. 同 digest 重装后 registry 只保留一个当前版本，不重复 mount、不复制历史项目/候选/Receipt；C01 仍引用原 v1.0.0。
3. 对每个 ProviderRequirement 验证可达的 `ready/not_ready/provider_missing/blocked/degraded` 诚实投影；无法安全诱发的状态写 `not_tested`，不能伪造。
4. 进入 Pack 作者工作台，只读核对 manifest、flow、角色槽、知识 scope、ProviderRequirement 和版本；覆盖全部可见按钮，禁止提交审核或发布。
5. 急停下角色、PDF、重复请求和双 Tab 均零业务副作用；解除后只恢复全新请求。

## 8. 证据、验收和收尾

每个关键动作同时保存 GUI 截图、网络状态、03/05/06/07/08/09/13/14 前后计数、refs 和时间戳；API/SQLite 只作对账。交付轮报告、事实逐句审计、19 节点表、按钮覆盖表、Provider 状态矩阵、lifecycle/历史回放时间线、问题台账、截图/Receipt 索引和端口清理证据。

独立验收主体确认：开放 P0=0、阻塞 P1=0；PDF 前可解释零副作用；PDF 后无越界事实；不合格模型输出 fail-closed；单项目安装→停用→启用→卸载→重装和历史查询全程通过；全部安全可继续 lane 已跑完。达到后才标`已验收（打包前）`；四个市场阶段始终为 `not_tested_owner_directed`。
