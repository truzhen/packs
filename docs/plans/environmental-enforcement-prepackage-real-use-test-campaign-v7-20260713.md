# 环保执法 Pack 测试计划 v7（真实云登录与真实市场三轮终验版）

日期：2026-07-13
状态：`契约已定 / 修复已实现 / 待新鲜运行态执行`
版本定位：当前客户主线，不是 V4 建议或 backlog。

## 1. 派活卡与最小交付

| 维度 | v7 裁定 |
|---|---|
| 我要做的事 | 在全新隔离运行态以真实 GUI 登录本地 cloud，通过真实云市场商品、授权、制品下载和 OS 安装链完成环保执法 Pack 三轮测试。 |
| 客户/场景证据 | v6 Chrome 已可访问 GUI，但 OS 固定调用 `/cloud-auth/password-login`，本地网关误代理到不存在的 `/auth/password-login` 并返回 HTML 404，R1 因此冻结。 |
| 最小可交付 | 先证明真实登录 JSON、Session/Receipt、market registry 0→1、sandbox Entitlement、真实制品 digest 和安装 0→1，再跑 C01 PDF 四路、双角色 compare、C02 隔离和急停零副作用。 |
| 真相源 | 身份/Session、Listing、Entitlement、下载审计归 cloud；安装、06 Run、03 Receipt、Base Gate/急停归 OS；Pack 只声明；client 只投影。 |
| 仓库归属 | cloud 修复密码登录代理；OS 消费既有稳定入口；client 走 GUI；packs 持计划与行业声明。contracts/software 只读。 |
| 风险颜色 | 登录/市场为橙；Receipt、Gate、急停、真实动作攻击为红，只允许受控验证和 fail-closed。 |
| 契约影响 | 不改 contracts/schema/DTO；只修既有实现接线与兼容响应。 |
| 生命周期 | 修复为`已实现`且代码验证通过；真实 GUI 新鲜运行态通过后才升为`已验收`。 |

本轮明确砍掉：生产账号、真实资金、生产市场、客户生产数据、真实执法执行、真实外发、合同变更、发布上架。它们不构成本轮 PASS 条件。

## 2. Owner 全权限授权口径

Owner 授权执行者在本计划列明的四个隔离 worktree 与本机隔离测试基础设施内，直接读取、修改、构建、测试、启动、停止、重启服务和容器；可创建/销毁本轮专用 HOME、SQLite、Docker volume、browser profile、测试账号、市场项目和临时制品；可使用本地 cloud seed 测试凭据，经 GUI 真实登录；可使用本地 cloud 市场完成 sandbox 授权、真实制品上传/下载、digest 核验和 OS 安装。执行中无需为每个中小修复再次请示。

“所有权限”仅指上述隔离测试闭环，不授权：修改 contracts；访问生产或客户系统/数据；真实付款；真实执法、发送或外部执行；读取或提交真实 secret；覆盖他人 WIP；提交、合并、推送、发布或清理非本轮资源。若测试工具需要凭据，只能使用本轮 seed/沙箱凭据，并在证据中脱敏。

允许修改/测试的四仓：

1. `truzhen-cloud`：本地身份、市场、Entitlement、下载链；不得接生产。
2. `truzhenos`：AuthIntent、Pack lifecycle、Run/Gate/Receipt；不得放宽主权边界。
3. `truzhen-client-web-desktop`：真实 GUI 与诚实错误呈现；不得用 API 结果冒充 GUI。
4. `truzhen-packs`：计划、Pack 声明和静态审计；不得实现 Provider。

## 3. 新鲜基线与真实登录硬门

过程目录使用 `/Users/li/Documents/过程文档/env-prepack-v7-20260713/`。R1-R3 每轮使用新 HOME、DB、cloud project、browser profile、origin、C01/C02、Run、Candidate、Receipt，并记录四仓 SHA、分支、WIP、diff hash、端口和容器。

真实登录必须同时满足：

1. Chrome 中从产品 GUI 输入本地 seed 测试账号和密码；禁止预置 session、直接写 localStorage、API 登录或复用旧 Cookie。
2. 网络证据显示 `POST /cloud-auth/password-login` 返回 JSON，`logged_in=true` 与非空 `session_id`；不得接受 HTML、mock 或静态 fixture。
3. OS Session ReadModel、GUI 左下角身份态和 cloud 服务端 session 同源；退出后均失效。
4. 错密码、无 cloud、非 JSON、过期 session 均 fail-closed，且不产生 Entitlement、下载或安装副作用。

## 4. 真实本地云市场硬门

1. 从 GUI 打开市场并读取 cloud Listing；空注册表起步，禁止 seed OS registry 冒充安装。
2. 通过本地 sandbox 商品授权链生成可反查 Entitlement。sandbox 不发生真实扣款，但必须走真实 market/commerce 状态机，不得手填 entitlement ref。
3. GUI 点击下载；cloud 返回真实 Pack 制品，记录制品 digest、下载审计和 entitlement 绑定。
4. GUI 点击安装；OS 校验制品并使 registry `0→1`。核对 enabled version、2 Role Slots、15 KnowledgeMount、752 FormalKnowledge。
5. 错 entitlement、错 digest、重放和跨账号下载均 fail-closed，不能污染 registry。

## 5. 三轮矩阵

### R1：核心闭环与隔离

- C01 PDF 经 GUI 提交，保存 GUI、network、06 ReadModel、SQLite 四路同源证据。
- 双角色真实本地模型 Candidate、Receipt、compare Gate/waiting；冷启动超时后手工候选只可回灌同一 run/step/candidate。
- 新建 C02 时不得显示 C01 正文、Candidate 或事务引用；再回 C01 历史仍可反查。
- 启用急停后 PDF、模型、candidate-input、Receipt、usage 增量均为 0。

### R2：生命周期与恢复

- disable/re-enable/reinstall、OS/cloud 重启、session 过期重登；停用不删历史 Receipt/Knowledge，停用期不允许新调用。
- 重复 C01 幂等复用同一合法回执并明确标重放；不得重复产生正式副作用。

### R3：对抗与多视口

- 伪/错 run/错 step/跨事务/过期/重放 receipt 与 decision ref 全部 fail-closed。
- 双 Tab、390×844、1024×768、1440×900；current run、候选正文、任务、Receipt 不串案。
- Provider missing/timeout/recovery、坏 PDF/加密/伪装/超限；只产受控 Candidate，不越过 Gate。

## 6. 问题处理与放行

中小问题由执行者直接完成“复现→根因→最小修复→定向及邻接回归→回到原 GUI 点→继续”，不得因可修小问题停止整轮。大问题必须登记严重性、污染范围、可继续 lane 与恢复门：安全独立 lane 能继续就继续；只有身份、权限、Receipt、Gate、跨事务隔离或真实外部副作用不可信，且不存在安全独立 lane 时才暂停。

R1-R3、真实登录、真实本地云市场、PDF 四路、跨事务隔离、急停零副作用和 Owner UAT 全有新证据，且开放 P0/阻塞 P1 为 0，才可判`已验收（打包前）`；否则保持`未验收 / 禁止打包放行`。
