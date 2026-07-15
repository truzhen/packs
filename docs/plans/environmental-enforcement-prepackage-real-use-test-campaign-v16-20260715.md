# 环保执法 Pack 单项目完整生命周期测试计划 v16

日期：2026-07-15  
状态：`当前客户主线；修复已实现 -> 已接线；待全新真实 GUI 独立验收；未发布`

前序证据：

- `/Users/li/Documents/过程文档/env-prepack-v15-20260714/rerun-20260715/v15-rerun-closeout-20260715.md`
- `/Users/li/Documents/过程文档/env-prepack-v15-20260714/rerun-20260715/issues/P1-emergency-stop-enable-receipt-missing-20260715.md`

本计划替代 `/Users/li/Documents/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v15-20260714.md`。

## 1. 派活卡与 Owner 裁定

| 维度 | v16 裁定 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线；只验收 `environmental-enforcement-pack-v0`，智能家居结果不冻结本计划。 |
| 真实场景证据 | v15 真实 GUI 已通过空 registry 安装、角色/知识挂载、新 C01、PDF 候选隔离、缺证据及急停沟通零副作用；唯一 P1 是急停启用缺 central 03 Receipt。 |
| 最小可交付 | 用唯一 C01 从直接安装、启用、完整 flow、停用、历史查询、重新启用、卸载、卸载后历史查询到同版本重装；急停启用、重放和解除均有唯一 central 03 Receipt，可反查 decision/transaction/evidence。 |
| 跳过项 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部跳过，标 `not_tested_owner_directed`；不启动 cloud、不建商品、订单或授权，不把跳过项写成 PASS。 |
| 真相源 | Pack/digest 归本仓；lifecycle 归 truzhenos 14；项目归 05；Run 归 06；知识归 09；模型调用归 08；候选/任务归 06/07；central Receipt 归 03；client 只展示。 |
| 仓库 / 层归属 | truzhenos 修 central Receipt 与运行安全；client 只展示 ReadModel/阻断/Receipt；本仓持 Pack 声明和计划。contracts、cloud、software 不改。 |
| 风险 | L0 知识/manifest 为绿；L1 Agent sandbox 与候选为黄；L2 OCR/模型/本地软件为橙；处罚、送达、删除、对外提交、正式签发和真实执法为红。 |
| 契约影响 | 不改 contracts、Receipt schema、Gate、Candidate、ReadModel 或 ProviderRequirement；只修既有 Base 急停 transition 的 central 03 记账实现。 |
| 禁止边界 | 不读真实客户资料或凭据；不正式化、不代签、不送达、不删除、不对外提交、不执行执法；不以 mock、旧证据、API-only 或容器存活冒充 GUI/Gate/Receipt。 |
| 生命周期 | 代码修复为`已实现 -> 已接线`；本计划跑通并由独立验收主体复核后才升`已验收（打包前）`；发布另行裁定。 |

## 2. 执行授权、自治与停止规则

测试执行者获本计划内全部本地权限：建立隔离 worktree/HOME/SQLite/浏览器 profile，启停 OS/client，调用真实本地模型，使用合成 PDF，直接安装、启停、卸载、重装，截图及只读 API/SQLite/Receipt 对账，并修复不改契约与主权链的中小问题。

遇到小问题或中等问题，必须定位根因、补回归并从最近可靠点继续，不得因可局部修复问题停掉全程。遇到 P0/P1 或其它大问题，登记问题卡并冻结受影响危险路径；若其它 lane 可安全独立运行，必须继续跑完。只有继续会扩大安全、法律、数据或主权风险，或所有剩余步骤都依赖污染状态时才暂停。测试任务不得自行合并、推送或发布。

每仓从固定主仓最新 `origin/main` 建分支，记录 HEAD、WIP、端口、容器和验证命令。L0/L1/L2 隔离级别写入每条证据；Provider/容器启动成功不能替代 readiness、Gate、Gateway 或 Receipt。

## 3. 共通放行硬门：Pack 不得腐化基座

1. manifest、schema、flow、role、capabilities、knowledge scope/index 全部通过权威校验；Pack 不复制、放宽或私造 contracts schema。
2. Candidate 与 Formal 在存储、接口和 UI 三层隔离；Role/Model/Pack 只能生成 Candidate，Formal 必须由具备资格的人确认并经过 Owner + Base Gate。
3. 权限只能是 Owner 授权、Role 白名单、Pack、ContextSlice、KnowledgeMount 和 Provider 能力的最小交集；Pack 不持凭据、不直连网络/文件系统/外部软件。
4. 缺配置、缺 Provider、缺证据、Provider 未 ready 或版本冲突，必须明确 `blocked/not_ready/provider_missing`；版本冲突使用稳定 reason，不得静默降级、换 Provider、补假数据或空成功。
5. 真实动作仅由对应 Gateway 执行；Gate 不能由 Pack、前端、Provider readiness、容器或模型代替。
6. 接受、阻断和执行尝试都必须有可回放 Receipt，并关联 candidate、decision、evidence、输入摘要、版本和副作用计数。
7. 安装、启停、重放、重试、双击、断网恢复和卸载使用稳定幂等键，不重复 Candidate、Formal、Receipt 或外部动作。
8. 卸载后 RoleBinding、KnowledgeMount、Provider session 和临时资源不可运行；历史项目与 Receipt 只读保留，不与其它 Pack 串包。
9. 静态审计 truzhenos 不含环保专用分支、专用 Gate/Gateway/Receipt、行业 schema 副本或 Pack seed；client 不铸真相；本仓不含 Provider 实现或运行产物。

## 4. 开跑恢复硬门

1. JSON/Python/结构审计、45 条知识与 15 scope 对账、路径/符号链接/禁品扫描、digest、`git diff --check` 全绿。
2. OS 定向测试证明急停首次启用产生恰好一条 `base_emergency_stop_enabled` central 03 Receipt；相同幂等重放不增量；解除产生恰好一条 `base_emergency_stop_disabled` Receipt；Receipt append 失败不得声称动作成功。
3. Receipt 必须含稳定 transaction、decision、candidate/evidence 关联和真实 reason 摘要；急停状态、Gate 结论、03 Ledger 与 GUI 同源可对账。
4. ENV-V14 事实锚定回归仍通过：PDF 前显示 `role_review_evidence_not_ready` 与缺失字段；PDF 后越界模型输出 `role_advice_grounding_violation` fail-closed。
5. OS EGR、client 全量测试/typecheck/build、本仓结构验证全部通过后才启动 GUI。

## 5. 唯一 C01 与直接生命周期

唯一项目为 `C01-v16-固体废物露天堆放线索核查`。同一 `transaction_ref`、Pack v1.0.0、digest、安装实例、浏览器 profile 和证据目录贯穿三轮；不得通过新建项目清空缺陷。只有权威状态被不可逆污染时，才允许一次 `C01-RECOVERY-01`，并登记旧/新 refs 与污染原因。

从净制品直接执行真实安装链，优先用可用 GUI 触发同一 lifecycle 端点；市场四阶段全部跳过。每次安装、启用、停用、卸载、重装必须用 GUI、14 ReadModel 和 03 Receipt 三方对账，脚本退出码不能单独算 PASS。

## 6. R1：安装、急停与 C01 完整 flow

1. 全新 registry=0、KnowledgeMount=0；安装后核对 registry 0→1、enabled v1.0.0、2 Role Pack、2 SlotBinding、15/15 mount、45/45 知识、安装 Receipt。
2. 在无业务请求时 GUI 启用急停：顶部与安全核心均显示生效；03 恰好新增一条 central Receipt。重复点击/重放不新增；急停下角色、PDF、模型、Candidate、Task、Provider、外发和执行增量全部为 0。
3. GUI 解除急停：必须新增唯一解除 Receipt；只恢复之后的新请求，不补跑急停期请求。
4. 创建唯一 C01，PDF 前 advice/challenge 必须显示具体 evidence not ready、缺 `pdf_parse_status` 与三项零副作用。
5. 上传合成 PDF，核对解析 Candidate/Receipt 与同一 06 Run 回灌；角色候选的每个事实、法条、辖区、材料引用都能定位到 ContextSlice/knowledge/evidence，未知项标 `pending_human_review`。
6. 完整覆盖 flow 全节点；处罚、送达、删除、对外提交为红色动作，只验证 Candidate、具备资格者确认门与零副作用阻断，不作正式执行。
7. 核对材料完整性、敏感信息脱敏、文书草稿与正式签发分区/分 ref/分状态；保持 0 正式处罚、0 送达、0 删除、0 对外提交、0 执法执行。

## 7. R2：停用、历史、恢复与故障

1. 同一 C01 停用 Pack，攻击角色、PDF 重放、旧请求、双击和双 Tab；03/06/07/08/09 业务增量为 0。
2. 停用态从项目列表、详情、时间线、候选和 Receipt 入口查询同一 C01；标题、Run、PDF、候选、Receipt、Pack 版本只读可回放。
3. 重新启用并重启 OS/client；历史固定原版本，只有新请求恢复，停用期请求不延迟补写。
4. 覆盖断网前已提交/未提交边界、恢复、幂等重试、390×844 / 1024×768 / 1440×900 以及全部可见按钮。

## 8. R3：卸载、重装、Provider 与作者视图

1. 卸载后 registry/enable 指针、RoleBinding、KnowledgeMount 和 Provider session 不可运行；重复卸载幂等；同一 C01 历史与 Receipt 仍可查。
2. 同 digest 重装后 registry 只保留一个当前版本，不重复 mount、不复制项目/Candidate/Receipt；历史 C01 仍引用原版本。
3. 逐一验证 ProviderRequirement 的可达 `ready/not_ready/provider_missing/blocked`；无法安全诱发的状态写 `not_tested`，不得伪造。
4. 作者工作台只读核对 manifest、flow、角色、知识、ProviderRequirement 和版本；不提交审核、不发布。

## 9. 证据与放行

过程目录：`/Users/li/Documents/过程文档/env-prepack-v16-20260715/`。每个关键动作保存 GUI、网络状态、03/05/06/07/08/09/13/14 前后计数、refs 和时间戳；API/SQLite 只作对账。交付三轮报告、flow 节点表、法律事实逐句审计、按钮表、Provider 矩阵、生命周期时间线、问题台账、截图/Receipt 索引和资源清理证据。

放行要求：开放 P0=0、阻塞 P1=0；全部可安全继续 lane 已跑完；central 急停启用/重放/解除 Receipt 通过；共通硬门、环保专项门与完整 lifecycle 均有真实证据。满足后仅可标 `已验收（打包前）`；市场四阶段保持 `not_tested_owner_directed`。
