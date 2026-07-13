# 两个 Pack 续测 P1 修复设计（2026-07-11）

## 1. 派活卡

- **我要做的事**：修复环保执法 Pack 在真实 GUI 市场获取链无法出现可安装商品，以及智能家居项目页状态与 SceneFlow 当前运行态不同源的两个 P1；修复后提供两份新的三轮测试计划。
- **版本 / 优先级**：当前客户主线，P1。
- **真实场景证据**：环保 R1 在空 registry 下进入 GUI 市场但无可安装 Pack；智能家居 R1、R2 中运行态已为 `waiting / owner_gate`，项目页仍显示 `intake`。证据见本轮报告与问题台账。
- **真相源**：市场商品、授权与审批桥 readiness 归 `truzhen-cloud`；受控市场代理、Pack lifecycle 与 SceneFlow 运行态归 `truzhenos`；前端只消费投影；Pack 仓只保存 Pack 与测试计划。
- **仓库 / 层归属**：`truzhen-client-web-desktop` 修正 API 路由和状态投影；`truzhen-cloud` 修正本地审批桥构建、健康门和真实测试商品准备；`truzhenos` 先以现有实现作为权威链并补验证，只有测试证明实现缺陷才修改；`truzhen-packs` 保存测试计划和收尾材料。
- **风险颜色**：前端实现与本地测试编排为黄；受控代理、Gate、Receipt、权限为橙/红边界。本轮不修改 contracts、Gate、Receipt、权限语义或正式事实规则。
- **生命周期目标**：本轮把两条修复从“设计中”推进到“已实现、已接线”，以自动化验证证明；真实 GUI 三轮完成后才能进入“已验收”。

## 2. 根因与修复设计

### 2.1 环保市场获取链

根因由两个独立缺口叠加：

1. 前端使用同一个 `packMarketBaseUrl` 承载旧 `/pack-proxy/*` 和受控 `/v3/market/proxy/*`。测试设置 cloud 市场地址后，受控请求被直接送到 cloud，绕开本地 OS 会话门并得到 404。
2. cloud 本地审批桥在容器启动阶段临时安装编译依赖并执行 `npm ci`；ARM64/musl 下 `better-sqlite3` 需本地编译。启动脚本未在审批桥失败时终止，市场代理随后以 `pack_registry_unreachable` 暴露 `not_ready`。

设计：

- 前端按路径分流：`/v3/market/proxy/*`、`/v3/market/admin-proxy/*` 与 `/v3/pack-studio/*` 始终走本地 OS `backendBaseURL()`；仅旧 `/pack-proxy/*` 使用可配置市场基址。
- cloud 为 canonical 01 审批桥增加可缓存的依赖镜像；源码仍只读挂载自 canonical 01，不在 08 复制第二份实现。
- 为审批桥增加健康检查；`up.sh` 必须等待健康态，失败即停止，不得继续启动市场代理并伪装成可测环境。
- 复用并参数化现有本地 E2E 商品准备脚本，使其可接收真实 Pack 产物；它只准备 cloud 商品、授权和下载物，不直接写 OS registry、不绕过 GUI 安装。

### 2.2 智能家居状态投影

SQLite 证据表明 `current_step_refs` 已持久化为 `owner_gate`，OS 当前运行态真相正确。前端 AppShell 与工作台仍从“运行列表 + 自选 live run”推导项目摘要，和详情页使用的唯一当前运行端点形成两条选择链，出现 stale / 歧义时保留创建时的 `intake`。

设计：项目摘要和工作台统一调用 `/v3/scene-flow/runs/current?transaction_ref={transaction_ref}`，只消费 `current_run`；无唯一当前运行时诚实保留原状态或显示无当前运行，不自行挑选另一条 run。

## 3. 不做事项

- 不改 contracts、DTO、schema、Gateway、Base Gate、Receipt 或权限边界。
- 不让 client 直连 cloud 受控接口。
- 不在 OS 增加错误 cloud 路径别名或 fallback。
- 不用脚本/API 直接写 OS registry 冒充 GUI 安装成功。
- 不提交、不推送、不合并、不发布；不碰生产环境、真实客户数据或真实支付。

## 4. 验收设计

- 先写失败测试，分别证明受控市场请求会被错误市场基址劫持、项目摘要未消费唯一当前运行态、cloud 审批桥不具备 fail-fast readiness。
- 修复后运行每仓最小相关测试、类型检查/静态审计，再运行 client 与 OS/cloud 的本地受控链验证。
- 自动化验证只把实现推进到“已接线”；下一轮必须使用全新运行态，从 registry=0、真实 GUI 获取 Pack 开始跑三轮。
- 小问题可自修，大问题登记的边界写入两份新计划，防止测试者用任意修复扩大授权。
