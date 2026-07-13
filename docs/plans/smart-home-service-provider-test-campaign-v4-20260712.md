# 智能家居服务商 Pack 测试任务 v4（三轮自主修复与条件续跑版）

日期：2026-07-12
状态：`契约已定 / 待执行`
版本定位：当前客户主线、打包前放行门
生命周期上限：当前为 `已接线`；R1-R3 与 Owner UAT 全绿后才可标记 `已验收（打包前）`。

## 1. 真实证据与本版目标

权威输入：

- v3 测试计划：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v3-20260711.md`
- v3 中止报告：`/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/final-report.md`
- v3 大问题台账：`/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/issues/big-issue-register.md`
- v2 完整用例库：`/Users/li/Documents/truzhenv3worktree/smart-home-service-provider-test-plan-20260711/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v2-20260711.md`

真实问题：OS 06 已稳定为 `waiting / owner_gate`，项目头却持续显示 `进行中 · intake`。根因已修复为持续读取唯一 current-run 端点、焦点/可见/停留刷新、请求序号与 `occ_version` 防乱序；独立 GUI 已看到前台自动收敛，但旧 R1 已冻结，不能追认总 PASS。

本轮最小可交付：R1 用全新环境证明项目/场景/BOW/详情与 06 current run 同源，并完成售后主链、任务、Provider 诚实态和跨项目候选隔离；R2 验证恢复、生命周期、双 Tab；R3 完成急停对抗、多视口、Dendrite readiness 与 Owner UAT。生产 Frappe 写回、真实外发/支付、客户数据、contracts 变更与治理控制台均砍入本轮之外。

## 2. 派活卡与权威边界

| 维度 | v4 裁定 |
|---|---|
| 我要做的事 | 连续跑 R1-R3；中小问题由执行者在隔离分支修复并推进；大问题登记后判断安全续跑或暂停。 |
| 场景证据 | “王先生灯组离线售后”真实 GUI 显示 intake，但 06 current run 已 waiting/owner_gate。 |
| 真相源 | 06 持 Run/current step；05 持事务生命周期；07 持任务；03 持 Receipt；Base 持 Gate/急停；Provider 持外部能力；client 只投影。 |
| 仓库归属 | client 消费 current run 与事务级投影；OS 持 05/06/07/03；packs 声明角色和 ProviderRequirement；cloud 仅市场/登录。 |
| 风险颜色 | 绿：文档/只读；黄：UI/普通流程；橙：Gateway/Provider 接线；红：权限、Gate、Receipt、急停、真实写回/发送、跨事务泄漏。 |
| 契约影响 | 不改 contracts/schema/DTO，不改变 Gate、Receipt、Candidate/Formal 或 ProviderRequirement 语义。 |
| 禁止边界 | 不触生产、不用真实客户、不真实 Frappe 写回/发送/支付，不写 main，不提交/合并/推送，不用 API/DB 冒充 GUI PASS。 |
| 用户验收 | Owner 从 GUI 复走安装→售后事务→owner_gate→任务→Provider→Receipt→急停→生命周期/历史。 |

## 3. 修复基线与新过程目录

测试直接使用以下含未提交修复的隔离 worktree，并逐仓记录 SHA、branch、WIP、diff hash：

- client：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-client-web-desktop`
- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhenos`
- packs：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs`
- cloud：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-cloud`

contracts、software 只读。新过程目录：

`/Users/li/Documents/过程文档/smart-home-service-provider-v4-20260712/`

不得续写 v3。每轮使用全新 HOME、SQLite、browser profile、origin、Docker project、账号隔离态、transaction、run、Candidate 与 Receipt。与环保测试不得共享 DB、HOME、profile、Vite origin、测试账号写环境或焦点控制。

## 4. 问题分级、自主修复与续跑规则

### 4.1 小问题：自行修复并继续

P2/P3 或局部展示、刷新、选择、文案、测试脚本问题，且可逆、归属明确、不改契约/权限/Gate/Receipt/认证、不触真实外部副作用时：

`SF 记录 → 最短复现 → 红灯测试 → 根因 → 最小修改 → 定向测试 → 邻接回归 → 原 GUI 检查点重跑 → PASS → 继续`

无需等待 Owner 再授权；写入 `R<N>/small-fixes/SF-NN.md`。

### 4.2 中等问题：既有授权内自行修复并推进

P1 或 client/OS/packs/cloud 的普通业务接线问题，只要根因明确、修改限于上述四仓隔离分支、不改 contracts/schema/DTO、Gateway 信任边界、Gate/Receipt/权限/认证、不产生客户数据或不可逆状态，并能补自动红灯与 GUI 复验，即可自行修复。写入 `R<N>/medium-fixes/MF-NN.md`，列影响、兼容和回滚点。

定向测试、邻接回归与失败 GUI lane PASS 后继续本轮；不因曾出现 P1 自动终止。连续两次最小修复失败、根因扩散或需要改契约时，升级为大问题。

### 4.3 大问题：登记后判断继续或暂停

P0、跨事务/跨用户暴露、权限/认证/急停/Gateway/Gate/Receipt/Candidate-Formal 穿透、真实 Frappe/发送/支付副作用、生产/客户数据、真相源迁移、contracts/schema/DTO、数据库迁移删除、根因不明或超授权均为大问题。

立即登记 `issues/big-issue-register.md` 并保存现场，然后做可继续性判定：

- **继续跑**：只污染单一 lane；共享 DB、身份、Receipt 与证据归属可信；危险入口可完全冻结；后续 lane 与故障无依赖。大问题保持开放，该轮不得总 PASS，但独立安全 lane 继续取证。
- **暂停当前轮**：存在信息暴露、主权/权限/急停穿透、正式或外部副作用、共享状态/证据污染；后续依赖故障链；无法物理冻结危险入口；继续会破坏现场。
- **记录要求**：影响面、已发生/未发生副作用、污染范围、可继续性、允许继续 lane 白名单、暂停 lane、恢复门槛。

不得刷新/清缓存掩盖问题，不得 API/DB 绕过，不得 fixture 假绿，不得用后续 PASS 稀释开放 P0/P1。

## 5. R1：新鲜安装、current-run 同源与售后主链

### 5.1 安装与资产硬门

1. 全新 registry=0，经真实 GUI 市场安装 `smart-home-owner-pack-v0`，禁止预装或脚本写 registry。
2. 对账 enabled version、1 个 Role Pack/SlotBinding、0 KnowledgeMount、3 项 Frappe ProviderRequirement 与 lifecycle Receipt。
3. Provider 必须诚实：未接 Frappe 时为 `provider_missing/not_ready/blocked`，不得以静态卡片显示 ready。

### 5.2 current-run 状态同源首门

GUI 新建“王先生灯组离线售后”并启动 flow。每次状态变化保留：创建响应、`/v3/scene-flow/runs/current`、run ReadModel、SQLite、项目侧栏、项目头、BOW 卡片、项目详情、场景页。

当 06 返回 `selection_state=selected`、`current_run.status=waiting`、`current_step_refs[0]=...:owner_gate` 后，所有前台表面必须在 5 秒内收敛为“等待中 · owner_gate”；停留轮询、窗口失焦再聚焦、隐藏再显示、刷新、切页、浏览器重开后均一致。不得回潮 `intake/待启动流程`，不得从 run 列表自行挑选。

必须追加乱序测试：先发旧 `occ_version` 响应、后发新响应并逆序返回，UI 不得倒退；`ambiguous/unavailable` 必须诚实显示歧义/不可用，只有无 current run 时才回落 05 lifecycle。

### 5.3 跨项目候选隔离

1. 在项目 A 生成角色 advice Candidate，记录正文、transaction/run/Candidate/Receipt refs。
2. 同一 profile 新建/切换项目 B；B 未生成前，DOM、技术详情和无障碍树不得包含 A 的正文或任何 A ref。
3. A 请求仍 pending 时切到 B，A 迟到响应必须被丢弃。
4. B 独立生成候选，只能引用 B；切回 A 后 A 仍正确。双 Tab A/B 再跑一次。
5. 任一串案按 P0 登记；证据归属污染时默认暂停当前轮。

### 5.4 售后闭环硬门

- 07 队列显示本 transaction 的真实 TaskCandidate，示例卡不得覆盖真实任务。
- 项目/客户快照与 Frappe 写回保持诚实 blocked；报价、采购、工期、沟通均只生成 Candidate。
- 单角色始终是 advice/proposer；正式任务、记忆、发送、写回必须经 Owner + Base Gate + Gateway + Receipt。
- 完成 SH-3R-01~08、LX 与 Receipt 时间线；R1 所有硬门 PASS 后才进入 R2。

## 6. R2：多项目、生命周期、恢复与竞态

使用新事务创建“李女士新装咨询”“施工增加灯带/门磁”。覆盖 disable/re-enable/reinstall、OS/cloud 重启、断网、重复点击、刷新、浏览器重开、双 Tab、OCC、Receipt replay、版本固定和历史回放。每次状态变化重复 current-run 六表面对账和跨项目旧正文/ref 抽检。

本地 alpha 账号完成作者工作台真实登录、刷新保持和退出；凭据只存在隔离运行态、证据脱敏。覆盖 SH-3R-01~09、11、13。中小问题按 4.1/4.2 修复后推进；大问题按 4.3 条件续跑或暂停。

## 7. R3：急停对抗、三视口、Dendrite 与 Owner UAT

1. 新建“验收后零星售后复盘”，验证任务历史、Receipt、MemoryCandidate，FormalMemory=0。
2. 伪 ref、非白名单发送、无效急停 body；急停下模型/PDF/任务正式化/沟通/Frappe 写回均 fail-closed，模型 usage/Candidate/Receipt/Formal/发送/执行增量为 0。
3. 解除急停只恢复候选能力；项目状态仍来自唯一 current run，正式动作仍需 Owner + Base Gate。
4. 双 Tab与 390×844、1024×768、1440×900 核项目状态、任务、Provider、Gate 和跨项目隔离。
5. 运行 Dendrite readiness/外部 E2E；未就绪保留原始日志并标 `blocked/not_ready`，不得扩大为联邦全量验收。
6. Owner 复走：市场安装→售后 transaction→owner_gate→任务→Provider 诚实态→Receipt→急停→停用/历史，并抽查三条 `transaction_ref → current_run → receipt_ref → ReadModel`。

## 8. 证据、清理与放行

每个关键 PASS 至少有 GUI、network、03 Receipt、ReadModel/SQLite 四路证据；安全断言再证明 Formal、外发、Frappe 写入为 0。每轮保存 baseline、参数、用例矩阵、small/medium fix、big issue 与可继续性结论、monitor doctor、端口/容器释放记录。旧 transaction、截图、fixture、API-only 或单元测试不能替代 GUI PASS。

放行标准：R1-R3 使用新 transaction/run/Receipt；SH-3R-01~13 按轮覆盖且 LX PASS；所有前台与 06 current 同源，无 intake 回潮或跨项目泄漏；P0=0、阻塞性 P1=0；Provider 诚实；Candidate/Formal、急停、Gate、Receipt 全程成立；中小修复均有红绿和 GUI 重跑；开放大问题为 0 或有 Owner 书面豁免、替代控制和到期日；Owner UAT 通过。否则保持 `未验收 / 禁止打包`。
