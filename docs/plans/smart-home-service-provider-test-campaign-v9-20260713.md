# 智能家居服务商 Pack 测试计划 v9（cloud 宿主转发与启停终验版）

日期：2026-07-13
状态：`契约已定 / 环境与产品修复已实现 / 待全新三轮真实 GUI 验收`
版本定位：当前客户主线；真实证据为 v8 两个 lane 容器内健康但宿主端口超时，以及 v7 停用后秘书仍调用模型。

## 1. 派活卡与最小闭环

| 维度 | v9 裁定 |
|---|---|
| 最小可交付 | 新 cloud lane 从安全未占用网段启动并通过四端口宿主硬门；真实 GUI 完成登录、授权、安装；Pack 停用后四个秘书入口零副作用。 |
| 真相源 | cloud 持市场/订单/Entitlement；14 持 lifecycle；05 持项目绑定；09/13/08/03/07 持挂载、会话、用量、回执与任务。 |
| 仓库/层 | cloud 08 负责隔离网段及宿主 readiness；OS 负责启停与秘书前置门；client 负责诚实投影；packs 只声明。 |
| 风险/契约 | 身份、Entitlement、Gate、Receipt、模型为红色。只修实现和本地 E2E 编排，不改 contracts、生产网络或支付契约。 |
| 生命周期 | cloud/OS 修复为`已实现`；新鲜 GUI 三轮完成前仍`未验收`。 |

非目标/backlog：生产账号/付款、客户数据、真实发送/Frappe 写回、生产部署、contracts、提交/合并/推送/发布。

## 2. 权限、上下文与问题推进

执行者拥有四个隔离 worktree、本轮专用 HOME/DB/Docker/browser profile 的修改、构建、测试、启停权限；必须真实 Chrome GUI 登录本地 cloud、真实 sandbox Entitlement、真实市场制品下载/digest/OS 安装。只接触 v8 报告及 cloud 08、OS 05/09/13/14、相关 client 页面和 Pack 声明；不得覆盖无关 WIP、读取真实 secret 或使用旧订单伪造成功。

中小问题必须自行修复并回原 GUI 点继续；大问题登记污染范围和恢复门。存在安全独立 lane 时继续跑完，只有身份、权限、Gate、Receipt、跨事务隔离或真实副作用不可信且没有安全 lane 时暂停。

## 3. R1：cloud 宿主硬门与启停零副作用

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v9-20260713/`；必须新建 cloud project/runtime/端口/HOME/profile/订单/Pack/项目/Run。

1. `up.sh` 自动选择与现有 Docker networks 不重叠的 `10.200-249.x.0/24`，禁止默认 `192.168.64.0/20`。返回成功前，宿主必须分别访问 Forgejo、node approval、market proxy、gateway；仅容器内 200 不算通过。
2. 真实 GUI 登录、市场展示、未授权阻断、sandbox 授权、制品下载/digest 和 OS 安装；禁止 API 注入 Entitlement、预装 registry 或复用旧订单。
3. 创建真实 05 Smart Home 项目；enabled 秘书只允许 1 次受治理模型调用，并形成同源 13 turn、Candidate、08 usage、03 Receipt。
4. GUI 停用时，所有目标 mounts 先 disabled，14 后提交；任一点失败必须补偿回 active 且 14 仍 enabled。
5. disabled 后攻击 `chat-candidate`、stream、`conversation/turn`、`tool/confirm`：模型、turn、Candidate、Receipt、FormalTask、发送、执行、Frappe 写入增量全部为 0；05/14 不可读则 `not_ready`。

## 4. R2/R3

R2：disable→reactivate→disable、安装 owner≠GUI actor、cloud/OS 重启、session 过期、Provider timeout/recovery、宿主网段并行 lane。
R3：双项目/双 Tab/OCC、390×844/1024×768/1440×900、多视口、伪/跨项目/重放 refs、历史 Receipt、急停、Dendrite 收窄与 Owner UAT。

开放 P0/阻塞 P1 为 0，R1-R3、四端口宿主证据、真实登录/市场、启停零副作用增量表和 Owner UAT 均为新证据，才可判`已验收（打包前）`；否则`未验收 / 禁止打包放行`。
