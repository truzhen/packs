# 智能家居服务商 Pack 单包跑通测试计划 v10

日期：2026-07-13
状态：`契约已定 -> 急停与两个 P2 修复已实现 -> 待补齐单 Pack 三轮真实 GUI 验收`
版本 / 优先级：当前客户主线。真实场景证据来自 v9：安装、启停、服务重启和急停恢复已通过真实前台；P1 是急停虽未触达 Provider/stream，仍把预检阻断写成 `calls=1`、Receipt Candidate 和 ModelRun。

## 0. v10 R1 补测修订

- 购买弹窗批准周期固定为周/月/年，对应 `7d/30d/365d`；即使旧 Market ReadModel 仍携带 `price_lifetime_minor`，GUI 也不得展示“永久”或提交 `lifetime`。不得为消除 UI 问题修改 Market DTO 或支付契约。
- 本地 E2E 专用 Cloud `up.sh` 必须显式向 market-proxy 传入 `TRUZHEN_LOCAL_E2E_PAYMENT_SIMULATION=1`，启动日志与容器 config 都要留证；只允许本地 sandbox 模拟支付，响应必须含不可发布标记。生产入口不得继承该默认值。
- v10 已完成的安装、秘书、停用、重启与急停链只算阶段证据，不等于总验收。修复后从购买弹窗做一次真实 GUI 回归，再继续剩余双 Tab、断网、三视口和 R2/R3 对抗；P2 若复现，执行者自行修复并继续，不得再次冻结整轮。

## 1. 派活卡与最小可交付

| 维度 | v10 裁定与建议 |
| --- | --- |
| 我要做的事 | 只验收 `smart-home-owner-pack-v0`，独立跑通真实 cloud 市场、sandbox Entitlement、安装、项目/秘书候选、启停、重启、急停与安全对抗。 |
| 最小可交付 | R1 单包市场与业务闭环；R2 单包生命周期/恢复；R3 急停零模型副作用与 Owner UAT。环保 Pack 的状态不得冻结本计划。 |
| 真相源 | cloud 持商品/订单/Entitlement/制品；14 持 lifecycle；05/06 持项目/Run；13 持秘书会话；08 持模型使用；03 持回执；07 持任务；11 持执行/Frappe。client 只投影。 |
| 仓库 / 层 | OS 08 实现急停最前置硬门；cloud 提供真实本地市场；packs 提供 Pack 声明和净制品；client 展示真实 blocked/not_ready。 |
| 风险颜色 | 普通 UI/文档为绿；生命周期为黄；Provider/市场为橙；急停、Gate、Receipt、正式任务、发送、Frappe 为红。AI 只能产生候选。 |
| 契约影响 | 只改实现：急停不进入 08 持久化；保留其它 Provider/配额/隐私受控失败的既有审计。不改 contracts/schema/DTO/Gate/Receipt 契约。 |
| 生命周期 | 修复为`已实现`；全新三轮 GUI 与数据库增量证明后才为`已验收（打包前）`。 |

明确砍掉 / backlog：生产账号和真实付款、客户生产数据、真实发送/Frappe 写回、生产部署、发布、提交/合并/推送、完整 Dendrite 兼容。

## 2. 权限、真实登录和禁止边界

执行者拥有四个隔离 worktree、本轮专用 HOME/DB/Docker/browser profile、修改/构建/测试/启停权限；必须使用真实 Chrome GUI 登录本地 cloud、浏览真实云市场、点击本地 sandbox 模拟支付取得 Entitlement、下载并安装真实制品。可自行完成中小修复并继续。禁止 API 注入 Entitlement、预装 registry、复用 v9 DB/订单/Receipt/截图、读取真实 secret、真实扣款、真实发送、Frappe 写回或覆盖他人 WIP。

上下文只读/修改范围限定于本计划、v9 报告、smart-home Pack、cloud 市场 lane、OS 05/06/08/13/14/03/07/11 和相关 client 页面。单包 lane 禁止安装环保 Pack。

## 3. 问题连续推进规则

- 小/中问题必须由执行者直接修根因、跑定向及邻接回归，然后从原 GUI 点继续；不得只登记后停止。
- 大问题登记等级、真相源、污染范围、副作用、冻结路径、恢复门。若不受影响的身份、市场、生命周期或只读 lane 安全，必须继续；只有 Gate/Receipt/权限/跨项目隔离不可信且无安全 lane 时暂停。
- 禁止 mock 成功、脚本代替 GUI、吞错误、复用证据或降低断言。

## 4. R1：单包真实市场与业务闭环

过程目录新建 `/Users/li/Documents/过程文档/smart-home-service-provider-v10-20260713/`。

1. 用 canonical market builder 生成净化 ZIP/sha256；cloud 上传前再做禁品、路径和 manifest 身份审计。全新商品/version、Chrome GUI 登录，证明未授权阻断；起栈证据确认本地模拟支付开关为 1，购买弹窗只显示周/月/年；点击本地 sandbox 模拟支付（无真实扣款），取得新 Entitlement，GUI 下载/安装，digest 对齐，registry 0→1。
2. 核对 Pack enabled、Role Pack/SlotBinding、ProviderRequirement 诚实态。新建真实 Smart Home 项目和 06 Run；项目页、通信中心和 ReadModel 必须同源，不得跨项目复用正文/ref。
3. enabled 下秘书请求恰好触发一次真实本地模型调用；只新增一个待 Owner Candidate，13 turn、08 usage、03 Receipt 同 transaction 可反查；0 FormalTask、0 发送、0 Frappe 写回。
4. GUI 停用必须先完成目标 mounts，再提交 14 disabled；停用后 `chat-candidate`、stream、conversation turn、tool confirm 四入口的模型、turn、Candidate、Receipt、Task、发送、执行、Frappe 增量均为 0。

## 5. R2：生命周期、重启与并发

只对智能家居 Pack 执行 disable→reactivate→disable→reactivate；覆盖 cloud/OS 重启、session 过期重登、Owner/actor 差异、幂等重放、Provider timeout/recovery、宿主双 lane 不冲突。双项目/双 Tab 切换后，05 项目、06 current-run、13 会话和候选正文/ref 必须严格同源；历史 Receipt 可反查，停用不删历史。

## 6. R3：急停零模型副作用硬门

1. 在全新 DB 记录急停前 08 usage、model receipt index、03 ledger、ModelRun、13 turn、Candidate、07 Task、09 Memory、发送/Frappe 基线。
2. GUI 启用全域急停，顶部与安全核心 ReadModel 同源；攻击非流式、流式、conversation turn、tool confirm、伪造/重放 receipt_ref。每次只允许返回瞬时 `blocked_reason=base_emergency_stop_active`。
3. 硬断言增量均为 0：Provider/stream、08 calls/tokens/failure、模型 Receipt Candidate、03 Formal Receipt、ModelRun、13 turn、业务 Candidate、FormalTask、FormalMemory、发送、执行、Frappe。响应不得携带 invocation/receipt/model-run/resolver ref。
4. 解除急停后仅发起一个全新请求，证明恰好一次模型调用和一套同源候选/回执恢复；急停期间请求不得延迟补写。
5. 补测 390×844/1024×768/1440×900、双 Tab、刷新、历史回放、Provider 四态与 Owner UAT。

## 7. 专门验收与完成门

验收主体独立复核代码定向测试、真实 Chrome 证据、市场/Entitlement/digest、05/06/08/13/14/03/07/11 增量表、急停竞态和服务清理。开放 P0=0、阻塞 P1=0，R1/R2/R3 均 PASS，才判智能家居单 Pack `已验收（打包前）`；否则登记大问题并按安全 lane 规则继续，不得因环保 Pack 状态否决本计划。
