# 智能家居服务商 Pack 单项目全周期复验计划 v17

日期：2026-07-16  
状态：`当前客户主线；v16 已验收；待本轮修复进入最新 main 后，以全新隔离环境复验；未发布`

本计划替代 `/Users/li/Documents/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v16-20260715.md`。

前序证据：

- `/Users/li/Documents/过程文档/smart-home-service-provider-v16-20260715/rerun-after-p0-p1-fix-20260715/最终测试报告.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v16-20260715/rerun-after-p0-p1-fix-20260715/问题台账.md`
- `/Users/li/Documents/过程文档/smart-home-v1-v16-repair-merge-audit-20260716.md`

## 1. 派活卡与 Owner 裁定

| 维度 | v17 裁定 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线；只验收 `smart-home-owner-pack-v0`，环保执法结果不冻结本计划。 |
| 真实客户 / 场景证据 | v16 真实 GUI 已通过六笔隔离 Frappe 写回、Provider/停用/双 Tab/秘书阻断、急停、卸载重装、三视口；审计发现 v16 的 OS 二次门控、前端回放/移动端及 Pack flow 节点类型尚未进入主线。 |
| 我要做的事 | 从四仓最新 `origin/main` 建全新隔离分支，用唯一项目贯穿商机、立项、进度、物料、交付、启停、故障、历史查询、急停、卸载和重装；证明合并后的主线没有依赖旧 WIP。 |
| 最小可交付 | 单一项目连续完成至少 10 次可回放的业务闭环；停用、急停、Provider 缺失、版本漂移和卸载均零新增业务副作用；历史项目始终可查且固定原 Pack 版本。 |
| 明确跳过 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部记 `not_tested_owner_directed`；不启动市场购买链，不把跳过项写成 PASS。 |
| 真相源 | Pack/digest 归本仓；项目/Run 归 05/06；Gate 归 01；Receipt 归 03；Frappe readiness/写回归 11；Frappe 持外部业务对象；Home Assistant 持设备实体状态；client 只投影。 |
| 仓库 / 层归属 | `truzhenos` 负责通用 lifecycle/Gate/Gateway/Receipt；client 负责安全状态和移动投影；本仓负责 flow/角色/能力声明；software 只登记 Provider。contracts/cloud 不改。 |
| 风险与 AI 角色 | L0 声明只读；L1 Agent sandbox 只产 Candidate；Frappe/Home Assistant 为 L2 隔离 Provider。AI 是测试施工者和 Proposer，不拥有正式化、执行、发送或项目真相。 |
| 契约影响 | 零 contracts、Gate/Receipt schema、Candidate schema 或跨仓 DTO 变更；发现契约缺口登记为大问题，不在测试任务现场扩契约。 |
| 禁止边界 | 不用真实客户/生产 Frappe/生产家庭设备，不保存 raw secret，不手改数据库，不绕 Gate/Gateway，不以 mock、容器存活或旧截图冒充 PASS。 |
| 生命周期 | v16 为`已验收（打包前）`；合并修复为`已实现 -> 已接线`；v17 复验通过后恢复`已验收（打包前）`；发布另行裁定。 |

## 2. 权限、自治、隔离与停止规则

测试执行者获本计划范围内全部本地权限：建立隔离 worktree、HOME、SQLite、浏览器 profile，启停 OS/client/Frappe/Home Assistant 测试实例，直接安装、启用、停用、卸载和重装 Pack，使用合成项目数据和真实本地模型，执行经 Owner + Base Gate 的隔离 Provider 写入，截图及只读 API/SQLite/Receipt/Frappe 对账。

小问题和中等问题必须现场定位根因、补测试并从最近可信检查点继续；不得只改表象或清空状态。P0/P1 或其它大问题必须登记问题卡并冻结受影响危险 lane；其余独立、安全 lane 继续跑完。仅当继续会扩大真实副作用、数据、设备或主权风险，或所有剩余步骤都依赖已经污染的权威状态时暂停。测试任务不得自行发布、部署或上架。

隔离级别必须写入证据：manifest/知识声明为 L0；模型和角色推理为 L1 sandbox；Frappe、网络、本地软件和 Home Assistant 为 L2 隔离 Provider。`health=200`、容器已启动或 `tools/list` 成功不等于 readiness、Gate、Gateway、业务结果或 Receipt 通过。

## 3. 共通硬门：Pack 不得腐化基座

1. `manifest.json`、flow、role、capabilities、software/provider requirements 通过权威 schema、结构和禁品校验；Pack 不复制、放宽或私造 contracts schema。
2. Candidate/Formal 在存储、接口和 UI 三层隔离。Role/Model/Pack 只能生成 Candidate；正式对象或执行意图必须由 Owner 明确确认并通过 Base Gate。
3. 权限取 Owner 授权、角色白名单、Pack、项目 ContextSlice 与 Provider 能力的最小交集；Pack 不持凭据、不直连 Provider、不申请通配设备/网络/文件权限。
4. 缺配置、缺 Provider、缺证据、Provider 未 ready、Pack/Provider/设备版本冲突时，必须明确返回 `blocked/not_ready/provider_missing` 和稳定 reason；不得静默降级、切换设备、补假数据、返回空成功或用 `degraded` 掩盖硬缺口。
5. 每个真实 Frappe/设备动作必须走 Candidate → Owner 确认 → Base Gate → Execution Gateway → L2 Provider → Receipt。Provider ready、前端提示或容器存活不能代替 Gate/Receipt。
6. 接受、阻断、失败、重放和执行尝试都有可反查 Receipt，关联项目、Pack/Provider 版本、candidate、decision、evidence、幂等键、结果与副作用计数；回放不得从当前状态重算冒充历史。
7. 重复点击、确认、旧请求、双 Tab、断网恢复、启停、卸载和重装使用稳定幂等键，不重复 Candidate、Formal、Receipt、Frappe 对象或设备动作。
8. 停用/卸载后 RoleBinding、Provider session、订阅、轮询、缓存和临时资源不可运行；历史项目和 Receipt 只读保留，不与其它 Pack/项目串用设备、会话、端口或缓存。
9. 静态扫描 `truzhenos` 不得出现 smart-home Pack ref、具体设备型号或 flow 节点专用分支、专用 Gate/Gateway/Receipt、行业 schema 副本或 seed；client 不持真相；Pack 不含 Provider/驱动；software 不保存 runtime/secret。

## 4. 唯一项目、直接安装与开跑硬门

唯一项目：`SH-C01-v17-全屋智能交付项目`。同一 `transaction_ref`、Pack v1.1.0、digest、Frappe site、浏览器 profile 和证据目录贯穿全轮；不得新建项目逃避缺陷。只有 05/06 权威状态不可逆污染时允许一次 `SH-C01-RECOVERY-01`，并登记旧/新 refs、污染原因与影响。

直接从净 Pack 制品安装，跳过市场四阶段。启动 GUI 前必须完成：四仓 HEAD/clean baseline；Pack JSON/Python/结构/禁品/digest；OS lifecycle 二次门控与重放定向/race；client 项目页、移动端、回放和无 `AbortError` 全量测试/typecheck/build/smoke；跨仓 EGR。任一失败先修根因，不带病开跑。

## 5. R1：安装与单项目十次业务闭环

1. 空 registry 安装并启用 v1.1.0；核对 registry 0→1、flow 全节点、Role Pack/SlotBinding、ProviderRequirement、digest 和 03 安装/启用 Receipt。
2. GUI 创建唯一 SH-C01；以同一证据轴完成商机 Lead、立项 Project、进度 Task 创建、Task 更新、物料 Item、交付 Ticket，另以不同稳定业务意图补足至少 10 次闭环。每次均为 GUI Candidate → Owner 确认 → Base Gate → Execution Gateway → 隔离 Frappe → 03 Receipt → 项目时间线。
3. 10 次闭环每次确认前 Frappe 增量必须为 0；确认页展示对象、动作、字段摘要、影响、Provider、幂等键与停止/回滚说明；确认后恰好一个外部对象或一次受控更新。重复确认/重放复用原结果，不新增 Provider 调用或 Formal 对象。
4. 项目页、沟通中心、任务、时间线和 Receipt 入口对同一 SH-C01 的阶段、refs、Pack 版本一致；秘书每次最多一次真实模型调用、一个待 Owner Candidate，完成或阻断后无流式占位残留。
5. 覆盖 flow 每个节点及全部可见按钮；不可执行项写明缺失前置或风险，不静默漏测。保持 0 对外发送、0 生产写入、0 未确认 Formal、0 正式 Memory。

## 6. R2：故障、停用、历史与恢复

1. 停止 Frappe 后分别攻击五阶段入口：必须 `not_ready/provider_missing`，零新 Candidate/Gate/业务 Receipt/Provider 写入；恢复后只处理新请求，旧请求不补写。
2. 制造版本漂移、重复点击、旧请求、双 Tab 和断网边界；候选持久化前与确认执行前都必须重新核验 lifecycle/事务绑定/急停，陈旧确认不能穿透。
3. 停用 Pack 后攻击秘书和全部写入入口；模型、turn、Candidate、Task、Receipt、Frappe 调用业务增量为 0。停用态仍能从各入口只读查询同一 SH-C01 历史。
4. 重新启用并重启 OS/client/Frappe；当前项目上下文、v1.1.0、时间线和 Receipt 保持，只有新请求恢复一次。
5. 复验 390×844、1024×768、1440×900；图标导航有 `aria-label/title`，正文高对比可读，横向 tabs 不挤压主内容；全程记录浏览器 console/network。

## 7. R3：急停、卸载/重装与 Home Assistant 可选 lane

1. 从顶部入口进入安全核心启用急停；单/双 Tab、旧 request、伪 receipt、重复点击、秘书和五阶段写入全部零业务副作用。急停状态、01、03 与 GUI 同源；解除后只恢复新请求。
2. 卸载后运行访问、SlotBinding、Provider session 和轮询撤销；新请求受控阻断，重复卸载幂等；同一 SH-C01 历史可查。同 digest 重装后 registry 不重复、历史不复制、项目仍固定原版本。
3. Home Assistant 不作为项目经营主链放行前提，不在 Pack/client/OS 自造硬件控制程序。未登记或 adapter 未 ready 时诚实 `provider_missing/not_ready` 并继续其它 lane。
4. 若 `/Users/li/Documents/trae/home-assistant-core` 已由 software/OS 以 L2 Provider 正式接线且测试实体获 Owner 授权，才测试读取实体状态和受控动作；执行前展示 entity、动作、影响、状态新鲜度、风险与回滚，经 Owner + Base Gate/Gateway 后执行并产 Receipt。
5. 可选设备对抗覆盖离线、状态漂移、重复指令、错误实体绑定、超时和人工急停；任一故障 fail-closed。禁止门锁、燃气、安防、生命安全、生产家庭设备以及为凑数 mock 成功。

## 8. 证据与放行

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v17-20260716/`。交付基线、十次闭环表、flow/按钮覆盖表、Provider 状态矩阵、三视口/双 Tab/断网记录、生命周期与历史回放时间线、问题台账、截图/Receipt 索引、完整门禁日志和资源清理证据。

放行要求：开放 P0=0、阻塞 P1=0；所有安全可继续 lane 已跑完；10/10 项目业务闭环可逐条回放；停用、急停、Provider 缺失、版本漂移和卸载零业务副作用；历史、幂等、断网、重装和基座防腐通过；完整前端回归无 happy-dom `AbortError` 噪声。满足后仅标 `已验收（打包前）`；市场四阶段仍为 `not_tested_owner_directed`。
