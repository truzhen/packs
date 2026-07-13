# 两个 Pack 续测 P1 修复实施计划

> 执行日期：2026-07-11
> 状态：已获 Owner 选择方案 A；不改契约，不提交、不推送。

## 目标

修复两条 P1 的根因并形成可复核自动化证据，为环保执法与智能家居各自全新三轮 GUI 复测提供稳定起点。

## 仓库与允许动作

| 仓库 | 职责 | 允许动作 | 禁止边界 |
|---|---|---|---|
| `truzhen-client-web-desktop` | 市场 API 分流、SceneFlow 当前态投影 | 修改、相关测试、类型检查 | 不直连受控 cloud、不铸真相 |
| `truzhen-cloud` | 本地审批桥 readiness、真实测试商品准备 | 修改本地 E2E 编排与脚本、隔离验证 | 不触生产、不改商业契约 |
| `truzhenos` | 受控代理与 SceneFlow 真相链 | 只读核对、相关测试；仅失败证据触发实现修改 | 不改 Gate/Receipt/权限/契约 |
| `truzhen-packs` | 设计、实施、续测计划与报告 | 文档修改、Markdown/路径审计 | 不改 Pack 主权声明和产品契约 |

## 实施步骤

### 任务 1：前端受控市场路由回归

1. 在 API client 测试中设置外部 `packMarketBase`，断言 `/v3/market/proxy/license/products` 仍请求本地 OS。
2. 运行单测确认先失败。
3. 在统一请求入口按受控路径与旧路径拆分 base URL；同步 XHR 上传/下载调用，避免旁路。
4. 运行市场域测试与类型检查。

### 任务 2：前端唯一当前运行态回归

1. 为 AppShell / 项目工作台补异步投影测试：创建态为 `intake`，当前端点返回 `waiting` 且 `current_step_refs` 为 `owner_gate`，最终页面必须显示 `waiting / owner_gate`。
2. 运行测试确认先失败。
3. 将摘要投影切换到 `getCurrentSceneFlowRunByTransaction`，不再自行从列表挑选 live run。
4. 运行相关页面测试、前端测试和类型检查。

### 任务 3：cloud 审批桥 fail-fast readiness

1. 扩充 full-stack isolation smoke，要求审批桥使用可构建镜像、配置健康检查，启动脚本在未健康时失败。
2. 运行 smoke 确认先失败。
3. 在 canonical 01 审批桥增加多阶段依赖镜像；08 compose 只引用构建上下文与 canonical 源。
4. 修改 `up.sh`：显式构建、等待 healthy、失败输出日志并退出，再启动 market。
5. 运行静态 smoke；在隔离端口执行 compose 构建和健康验证。

### 任务 4：真实环保 Pack 市场准备

1. 为现有 `ensure-pack-e2e-market-product.sh` 增加可选真实 Pack 产物输入的失败测试。
2. 参数化脚本：上传真实压缩包、创建/更新测试商品与 entitlement；不得直接写 OS registry。
3. 在隔离 cloud 环境准备 `environmental-enforcement-pack-v0`，验证 catalog、entitlement、download gate 均 ready。

### 任务 5：OS 真相链证明

1. 运行受控市场代理相关 Go 测试，证明无 session 被拒绝、有 session 才可转发。
2. 运行 SceneFlow current/readmodel 相关测试，证明 `waiting / owner_gate` 与 `current_step_refs` 可持久化并由 current 端点返回。
3. 若测试通过则不改 OS 产品代码；若失败，只修已证明根因且保持契约不变。

### 任务 6：两份三轮续测计划

分别创建 v3 环保与智能家居测试计划。两份计划必须：

- 使用全新运行态、全新证据目录、固定四/六仓 SHA 与端口；禁止复用已污染数据库。
- 每轮从真实 GUI 开始，不用脚本/API 伪造最终用户动作。
- 定义小问题自修循环：定位、最小修复、相关回归、从失败检查点重跑，直至该轮通过。
- 定义大问题登记门：P0/P1、契约/真相源/Gate/Receipt/权限、生产或真实外部副作用、数据迁移/删除、跨仓架构变化、同一根因两次修复仍失败。登记后冻结危险 lane；无依赖的安全 lane 可继续。
- 保留每轮报告、lane 证据、问题台账、端口/容器释放与仓库洁净证明。

### 任务 7：收尾验证

逐仓运行相关测试、静态检查与 `git diff --check`，记录实际输出。生成修复收尾报告，明确自动化已通过的范围与仍需 GUI 三轮验证的范围。
