# 智能家居服务商 Pack 测试计划 v7（急停回填主权链与云市场三轮终验版）

日期：2026-07-13
状态：`契约已定 / P0 修复已实现 / 待新鲜运行态执行`
版本定位：当前客户主线，不是 V4 建议或 backlog。

## 1. 派活卡与最小交付

| 维度 | v7 裁定 |
|---|---|
| 我要做的事 | 以全新隔离运行态重跑三轮，证明急停下能力结果零写入、Receipt 与 run/step/candidate 强绑定、沟通中心保留 05 事务上下文，并完成真实云登录/市场安装及剩余 R3。 |
| 客户/场景证据 | v6 在急停已启用时接受伪造 `receipt_ref` 并写入 SceneFlowRun（P0）；项目返回沟通中心丢失 `transaction_ref`（P2）。 |
| 最小可交付 | 先过伪 Receipt + 急停双攻击硬门，再过真实 GUI 登录/市场安装、项目/Run/六类 Candidate/Provider/Receipt 主链；最后收口多 Tab、多视口、生命周期、Dendrite 和 Owner UAT。 |
| 真相源 | Base 持急停/Gate；03 持正式 Receipt；06 持 Run；05 持事务对象；07 持任务；cloud 持身份/Listing/Entitlement；Provider 持外部能力；client 只投影。 |
| 仓库归属 | OS 修 capability-result 前置校验；client 修显式沟通上下文；cloud 提供真实登录/市场；packs 持计划/声明。contracts/software 只读。 |
| 风险/契约 | P0 为红：只强化既有 Gate/Receipt 实现，不改契约、不放宽权限。P2 为黄：只按 transaction_ref 重读 05。 |
| 生命周期 | 修复为`已实现`且自动测试通过；新鲜 GUI/四路证据通过后才为`已验收`。 |

本轮不包含生产账号、真实资金、客户生产数据、真实外发、真实 Frappe 写入、发布上架或 contracts 修改；缺少受控真实 Provider 时必须诚实收窄，不得 fixture 冒充。

## 2. Owner 全权限授权口径

Owner 授权执行者在四个隔离 worktree 与本轮本机隔离基础设施内，直接读取、修改、构建、测试、启动、停止、重启服务/容器；可创建/销毁本轮专用 HOME、SQLite、Docker volume、browser profile、测试账号、market project、Frappe/Dendrite 受控测试实例与临时制品；可用本地 cloud seed 测试凭据经真实 GUI 登录；可经真实本地云市场完成 sandbox 授权、制品上传/下载/digest 核验和 OS 安装。中小问题无需逐项请示，修后继续跑。

“所有权限”不扩张到生产/客户系统、真实付款、真实外发/执行、真实 secret、contracts 变更、覆盖他人 WIP、提交/合并/推送/发布或清理非本轮资源。白名单发送只验证受控 Gateway 链，除非 Owner 另行明确授权真实发送。

四仓职责：`truzhenos` 负责 Base/03/05/06/07/Gateway；`truzhen-client-web-desktop` 负责 GUI 投影；`truzhen-cloud` 负责身份/市场/Entitlement/制品；`truzhen-packs` 负责声明、计划和静态审计。依赖方向不变，contracts/software 只读。

## 3. 新鲜基线与 P0 首门

过程目录使用 `/Users/li/Documents/过程文档/smart-home-service-provider-v7-20260713/`。R1-R3 每轮使用新 HOME、DB、cloud project、browser profile、origin、Pack 安装、项目、Run、Candidate、Receipt，并记录四仓 SHA/WIP/diff hash及服务清理证据。

进入任何 Provider 写回、正式化或 resume 前必须完成：

1. 建立带真实 capability candidate 的 run/step，并由受控 Gateway 产生 03 `recorded` Receipt；合法回填必须成功且只更新目标 run。
2. 伪造、404、错 candidate、错 step、跨 run/事务、非 recorded、重复使用 Receipt 均必须在写入前拒绝；06 state/timeline/evidence、03、usage、Formal 计数无增量。
3. GUI 显式启用急停后，合法或伪造 Receipt 均返回 `emergency_stop_active`；运行时 merge 函数不得被调用，Run/Receipt/外部系统无副作用。
4. Base 状态不可读或 Receipt verifier 不可用时返回服务不可用并 fail-closed，不得降级接受非空字符串。
5. 前端必须把拒绝显示为失败，不得显示“已回写”。

P0 首门任一失败，冻结能力回填、Frappe、正式化、resume、发送和真实执行；可继续只读及不依赖该边界的安全 lane。

## 4. 真实登录、市场和沟通上下文

1. Chrome 从产品 GUI 输入本地 seed 测试账号/密码，`POST /cloud-auth/password-login` 返回 JSON `logged_in=true`；禁止 API 登录、旧 Cookie、localStorage 注入或预造 Session。
2. GUI 从真实本地 cloud Listing 完成 sandbox Entitlement、真实制品下载、digest 校验和 registry `0→1` 安装；错 entitlement/digest/账号必须拒绝。
3. 安装后核对 smart-home Pack enabled version、Role Slot、ProviderRequirement、flow 和 manifest。
4. 从项目页点击“返回沟通中心”必须携带 transaction_ref，并由 AppShell 从 05 ReadModel 重新选中；草稿 Candidate 与右侧上下文显示同一真实事务。普通侧栏进入沟通中心仍为全局新对话，不得继承旧项目。
5. 事务不存在或 05 未投影时诚实显示缺失，禁止前端临时对象冒充真相。

## 5. 三轮矩阵

### R1：核心与主权链

- 真实登录/市场安装→新项目→current run→六类 Candidate→任务→Provider→Receipt→Owner/Base Gate。
- 运行 P0 首门全部攻击；核对 transaction/run/candidate/receipt 四路引用。
- 回到沟通中心验证 P2；再切另一项目验证上下文不串案。

### R2：生命周期与故障恢复

- disable/re-enable/reinstall、OS/cloud/Provider 重启、session 过期重登。
- Provider missing/timeout/disconnect/recovery，Frappe query/write 未接通时诚实 blocked；PDF 正常/坏文件/加密/伪装/超限/重复。
- 停用不删历史 Receipt，停用期不允许新能力调用。

### R3：并发、终端与历史

- 双 Tab 同项目 OCC/幂等，双项目 current run/正文/task/receipt 隔离。
- 390×844、1024×768、1440×900 覆盖项目状态、任务、Provider、急停和 Gate。
- Dendrite readiness/外部 E2E；未就绪保留原始日志并收窄联邦结论。
- 历史回放与增量统计：每个攻击 lane 记录模型 usage、Candidate、Receipt、FormalTask/Memory、发送、执行、Frappe 写入前后计数。

## 6. 问题处理与放行

中小问题必须由执行者自行“复现→根因→最小修复→定向/邻接回归→原 GUI 点复跑→继续”，不能因可修问题提前交卷。大问题登记严重性、污染范围、可继续 lane、暂停 lane和恢复门；若安全独立 lane 能继续就继续，只有身份、Receipt、Gate、权限、跨事务隔离或真实副作用不可信且无安全替代 lane 时暂停。

R1-R3、P0 攻击门、P2 上下文、真实登录/市场、生命周期、多 Tab/多视口、历史、增量统计和 Owner UAT 全有新证据，且开放 P0/阻塞 P1 为 0，才可判`已验收（打包前）`；否则保持`未验收 / 禁止打包放行`。
