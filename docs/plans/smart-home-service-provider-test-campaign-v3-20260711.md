# 智能家居服务商 Pack 测试任务 v3（三轮持续修复版）

日期：2026-07-11
状态：`契约已定 / 待执行`
归属：当前客户主线，打包前放行门
目标：使用修复后的唯一 SceneFlow 当前运行态投影重新跑 R1-R3；小问题在隔离分支自行修复并重跑直至通过，大问题登记并冻结危险路径。

## 1. 权威输入

- v2 用例库：`/Users/li/Documents/truzhenv3worktree/smart-home-service-provider-test-plan-20260711/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v2-20260711.md`
- 上轮最终报告：`/Users/li/Documents/过程文档/smart-home-service-provider-v2-rerun-20260711/final-report.md`
- 上轮问题台账：`/Users/li/Documents/过程文档/smart-home-service-provider-v2-rerun-20260711/issue-register.md`
- 四仓冻结基线：`/Users/li/Documents/过程文档/smart-home-service-provider-v2-rerun-20260711/baseline/four-repo-baseline.txt`
- 本次修复设计：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs/docs/plans/two-packs-rerun-p1-repair-design-20260711.md`

v2 的 L1-L11、LX、四业务剧本、SH-3R-01~13 仍是完整用例库。本文件覆盖新的状态同源断言、持续小修机制、证据目录和三轮停走门。

## 2. 派活卡

| 维度 | v3 裁定 |
|---|---|
| 我要做的事 | 复验项目页、工作台、详情、任务、Provider、生命周期、Candidate/Formal、急停、登录、Dendrite readiness 与三视口。 |
| 真实场景证据 | R1、R2 中 06/SQLite 已为 `waiting / owner_gate`，项目/场景前台仍显示创建态 `intake/待启动流程`。 |
| 最小可交付 | R1 证明唯一 current run 投影；R2 证明重启/双 Tab/多项目不回潮；R3 对抗与 Owner UAT。 |
| 真相源 | SceneFlow 当前运行态归 OS 06；项目生命周期归 OS 05；前端不得从 run 列表自行挑选或保留 stale 创建态。 |
| 仓库归属 | client 消费 `/runs/current`；OS 持 current run 与 `current_step_refs`；packs 声明单角色/ProviderRequirement；cloud 只在登录/市场 lane 使用。 |
| 风险颜色 | UI 投影、测试为黄；contracts/Gateway 为橙；Gate、Receipt、权限、真实 Frappe/发送为红。 |
| 契约影响 | 不改契约；只统一消费现有 current run 端点。 |
| 禁止边界 | 不真实 Frappe 写回、不外发、不生产支付、不用真实客户、不改 Gate/Receipt/权限、不提交/推送/合并。 |

## 3. 新过程目录与执行基线

本轮使用：

`/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/`

结构：

```text
baseline/
R1/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
R2/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
R3/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
issues/big-issue-register.md
final-report.md
```

修复未提交时必须直接使用并登记以下 worktree 的完整 SHA、WIP 与 diff hash，不得从旧 rerun worktree 启动：

- client：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-client-web-desktop`
- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhenos`
- packs：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs`
- cloud：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-cloud`

每轮使用全新 HOME/DB/browser profile/origin 和新 transaction/run/Receipt；建议 OS `18161`、Vite `5191`，但以参数单为准。与环保 lane 不得共享 DB、HOME、profile、origin、测试账号写环境或焦点控制。

## 4. 状态同源硬门

### 4.1 网络与 ReadModel

对每个新 transaction 必须保留：

1. 创建返回的 lifecycle 状态；
2. `/v3/scene-flow/runs/current?transaction_ref={transaction_ref}` 原始响应；
3. `current_run.status`、`current_step_refs[0]`、run_ref；
4. SQLite/ReadModel 对应行；
5. 项目侧栏、项目头、BOW 卡片、项目详情、场景页截图。

当前端点返回 `selection_state=selected`、`status=waiting`、当前步为 `owner_gate` 后，所有前台表面必须在 5 秒内收敛为“等待中 · owner_gate”，刷新、切页、浏览器重开后仍一致。不得显示 `intake`、`待启动流程` 或从 run 列表自行选另一条 run。

若 `selection_state=ambiguous`，前端必须诚实显示状态不可用/歧义，不得猜测；若无 current run，才允许回落 05 lifecycle。

### 4.2 OS 证明

OS 必须证明 current run 选择唯一、`current_step_refs` 已持久化并由 current 端点返回。若 OS 定向测试全绿而 UI 不一致，优先归 client 投影；不得为了修 UI 改 OS 真相。

## 5. 测试中持续修复机制

### 5.1 小问题自动修复

小问题必须同时满足：局部、可逆、仓库归属明确；不改 contracts/schema/DTO、Gate、Receipt、权限、认证、主权链；不触生产/真实外部副作用/迁移删除；能用失败测试复现。

循环：

`small-fix 日志 → 失败测试 → 根因 → 最小修改 → 定向测试 → 邻接回归 → 从 GUI 失败检查点重跑 → PASS 后继续`

允许在 client/OS/cloud/packs 隔离修复分支内修改现有实现与测试；contracts/software 只读。每条写入本轮 `small-fixes/SF-NN.md`，保存 diff、红/绿证据、GUI 重跑、Receipt/ReadModel 对账。小问题必须修到该用例通过，不能仅登记后继续冒充 PASS。

### 5.2 大问题登记门

以下只登记，不擅改：P0；contracts/schema/DTO；真相源或跨仓边界变化；Gateway、Gate、Receipt、权限、认证；生产支付/发送/执行、真实 Frappe、客户数据；数据库迁移/删除；同一根因连续两次最小修复仍失败；需要扩大当前授权。

登记到 `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/issues/big-issue-register.md`，冻结受影响危险 lane并保留现场。独立安全 lane 可继续取证，但该轮不得总 PASS。禁止 fallback、fixture 或静态文案假绿。

问题状态：`已登记 → 已复现 → 根因已确认 → 待 Owner 裁定 → 修复获授权 → 已实现 → 已接线 → 独立回归 PASS → 终轮 PASS`。

## 6. 三轮执行

### R1：全新安装、售后主链与状态同源

1. 全新 DB，registry=0；真实安装智能家居 Pack，核 enabled version、单角色、slot binding、无 KnowledgeMount、3 项 Frappe requirement 与 lifecycle Receipt。
2. GUI 创建“王先生灯组离线售后”新 transaction并启动 flow。
3. 在 `owner_gate` 前同时抓 current endpoint、SQLite、项目侧栏、项目头、BOW 卡片、详情、场景页；全部满足第 4 节。
4. 07 任务系统显示本 transaction 的真实 TaskCandidate；示例卡不得覆盖真实队列。
5. 项目/客户快照显示 `provider_missing`，Frappe 写回显示 `blocked`；单角色保持 advice/proposer。
6. Candidate/Formal、沟通、Memory、Receipt、时间线执行 v2 SH-3R-01~08 与 LX。

R1 放行门：任一页面仍显示 intake/待启动、真实任务不可见、Provider 假 ready、Candidate 写 Formal、Receipt 断链，均先按第 5 节处理；未 PASS 不得进入 R2。

### R2：多项目、生命周期、重启与登录

1. 新建“李女士新装咨询”“施工增加灯带/门磁”两个 transaction，各有独立 current run。
2. disable/re-enable/reinstall、OS 重启、刷新、浏览器重开、断网重试、重复点击、双 Tab 同时推进。
3. 每次状态变化后重复第 4 节六表面对账；不得跨项目串 run，不得由 waiting 倒退 intake。
4. 任务、Receipt、历史、Provider 状态不重复、不丢失；报价、工期、采购、发送均只产候选。
5. 用本地 alpha 测试账号完成作者工作台真实登录、刷新保持、退出；凭据只在隔离运行态，证据脱敏。
6. 执行 v2 SH-3R-01~09、11、13；小问题修到检查点 PASS。

### R3：复盘、急停对抗、三视口与 Owner UAT

1. 新建“验收后零星售后复盘”transaction；验证任务历史、Receipt、MemoryCandidate，FormalMemory=0。
2. 伪 ref、非白名单发送、急停下模型/PDF/任务正式化/沟通/Frappe 写回、停用后访问均 fail-closed。
3. 解除急停只恢复候选能力，项目状态仍来自唯一 current run，正式动作仍需 Owner + Base Gate。
4. 双 Tab与 390×844、1024×768、1440×900 分别核项目状态、07 队列、Provider、急停和 Gate 可达。
5. 运行 Dendrite readiness/外部 E2E；未就绪保留原始日志并标 `blocked/not_ready`，不得宣称联邦全量验收。
6. Owner 亲手复走售后最短旅程，并抽查三条 `transaction_ref → current_run → receipt_ref → ReadModel`。

## 7. 证据与独立复核

每个关键 PASS 至少有 GUI、network、03 Receipt、ReadModel/SQLite 四路证据；安全断言再证明 Formal 表、外发和真实执行为 0。旧 transaction、旧截图、单元测试、fixture、API-only 或执行者自报不能单独构成 PASS。

每轮报告必须列：参数单、用例矩阵、small-fix 日志、big issue、Receipt/ReadModel 索引、SKIP/外部依赖、端口/容器释放、独立 Validator 结论。执行者不能自审关闭自己修复的问题。

## 8. 放行标准

- R1-R3 全部使用新 transaction/run/Receipt，SH-3R-01~13 按轮覆盖，LX PASS。
- 项目工作台、项目详情、场景页、BOW 与 current endpoint/SQLite 同源，无 intake 回潮。
- P0=0、阻塞性 P1=0；big issue 已关闭或有 Owner 书面豁免、替代控制、到期日。
- 小问题都有失败测试、修复、邻接回归和 GUI 重跑；不存在“只修自动化、不重跑前台”。
- Provider 状态诚实；Candidate/Formal、急停、Gate、Receipt 全程成立。
- Dendrite 未就绪只能收窄结论；Owner UAT 未通过时保持 `未验收 / 禁止打包`。

## 9. 启动前 500 字简报模板

只汇报：轮号、四仓修复分支/diff hash、六仓 SHA、HOME/DB/profile/端口、registry 是否为 0、OS current/readiness、client 构建来源、Provider 档位、已知大问题、是否允许进入 GUI。不要复述整份计划。
