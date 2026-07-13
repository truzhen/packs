# 两个 Pack v3 阻断修复计划（2026-07-12）

状态：`已实现 / 已接线 / 独立回归 PASS；待全新 R1`
版本归属：当前客户主线
优先级：环保登录 P0、智能家居状态投影 P1
契约结论：不改 `truzhen-contracts`、schema、DTO、Base Gate、Receipt、权限语义或 Pack 主权声明。

权威证据输入：

- `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/final-report.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/issues/big-issue-register.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/R1/small-fixes/SF-01-current-run-projection.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/R1/readmodels/r1-retry4-current-run-final.json`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v3-20260711/R1/shots/r1-retry4-ui-stale-intake-final.png`
- `/Users/li/Documents/过程文档/env-prepack-v3-20260711/issues/big-issue-register.md`
- `/Users/li/Documents/过程文档/env-prepack-v3-20260711/R1/lanes/R1-E1-GUI-market-login-P0.md`
- `/Users/li/Documents/过程文档/env-prepack-v3-20260711/baseline/baseline-20260712.md`
- `/Users/li/Documents/过程文档/env-prepack-v3-20260711/final-report.md`

## 1. 派活卡与 Owner 已给边界

| 维度 | 本次裁定 |
|---|---|
| 我要做的事 | 修复 GUI 本地市场登录错误转向生产认证，以及项目头未持续投影 06 当前运行态；完成独立自动回归后，分别从全新运行态恢复环保与智能家居 R1。 |
| 真实客户 / 场景证据 | 当前客户主线。环保在 registry=0 的真实 GUI 市场登录触发生产认证默认值；智能家居真实 Chrome 项目头显示 `进行中 · intake`，同一事务的 06 ReadModel 已为 `waiting / owner_gate`。 |
| 最小可交付 | 环保：本地 GUI 登录全程不访问生产域，OS 建立 opaque session，未登录 401、登录后受控市场代理 200。智能家居：真实 GUI 在 06 变化后收敛到 `等待中 · owner_gate`，刷新、切页、重新打开仍一致。 |
| 明确砍掉 / backlog | 不顺带改免费商品 Gate 语义，不建设新的事件总线，不改 Pack 内容，不做生产登录、真实支付、真实 Frappe、真实发送或正式执法。免费 Forgejo 商品 `payment_required` 继续作为独立 Gate 大问题。 |
| AI 角色 | 施工者：仅在既有隔离修复分支做最小实现与测试；建议者：涉及认证目的地址信任边界时先按本计划防护；旁观者：Gate、Receipt、权限及生产系统只取证不施工。 |
| 生命周期目标 | 修复完成最高推进到 `已接线 / 独立回归 PASS`；只有新的三轮 GUI 与 Owner UAT 全绿后才是 `已验收`，本计划本身不构成验收。 |

## 2. 检查结论

### 2.1 SHV3-SH3R01-P1-01：06 正确，前端投影会过期

证据链：

- GUI 截图显示 `R1-Retry4-Wang-Light-Offline-AfterSales` 为 `进行中 · intake`。
- 同一 `transaction_ref=tx_35665efe94850d7f` 的 `/v3/scene-flow/runs/current` 返回 `selection_state=selected`、`status=waiting`、`current_step_refs[0]=...:owner_gate`。
- OS 修复 worktree 无产品代码 diff，既有 current 端点和持久化证据未显示真相错误。
- client 当前未提交实现把 `run_status/run_node_label` 合并进 `sessions`，并只在新建后执行最多 20 次、每次 500ms 的观察。观察窗口结束后没有持续失效、重新取数或版本防倒退；因此初始 `running/intake` 可以长期滞留。

根因结论：这是 client 的 ReadModel 投影生命周期缺陷，不是 05/06 真相迁移问题。继续增加轮询次数只能延后复现，不能解决已有 run 后续变化、切页、双 Tab 和慢机器下的陈旧状态。

### 2.2 ENV-V3-P0-001：浏览器默认值抢走 OS 的可信路由选择

证据链：

- `AuthPage.readCloudAuthBase()` 在没有运行时注入时默认 `https://www.truzhen.com`，密码登录把它作为 `cloud_auth_base_url` 发送给 `/v3/auth-gateway/login-intent`。
- OS `authintentdev` 对非空 `cloud_auth_base_url` 走 `/cloud-auth/password-login`；因此本地测试没有进入已存在的 controlled real E2E 路径。
- OS 在 `realE2E` 开启且 `TRUZHEN_MARKET_PROXY_BASE_URL` 已配置时，既有 `controlledDevAuthClient` 会通过 market proxy `/auth/login` 获取真实 cloud session，并由 OS 导入为本地 opaque session；这条链无需新契约。
- 本地 cloud 对同一受控买家 `/auth/login` 已返回 200，说明 cloud 身份和商品授权不是本次登录失败根因。

根因结论：密码类认证的上游目的地址不应由浏览器默认值或 query 参数决定；可信目的地址应由 OS 运行配置选择。生产浏览器跳转 URL 与本地 Auth Gateway 密码代理目的地必须拆开。

## 3. 真相源、仓库归属与影响

| 仓库 / 层 | 真相与职责 | 计划动作 | 风险 |
|---|---|---|---|
| `truzhen-client-web-desktop` | 只消费 06 current run；只向本地 Auth Gateway 提交认证意图 | 修正当前 run 投影刷新；密码/注册/找回/SMS 不再携带浏览器选定的上游 base；保留扫码跳转所需展示 URL | 黄；认证调用边界按红色审查 |
| `truzhenos` | 06 持 current run；Auth Gateway 持可信路由选择与 opaque session 导入 | 先补失败测试；确认 client 不传 base 时走既有 realE2E market login；增加 fail-closed 防护，拒绝或忽略普通客户端覆盖认证上游 | 橙/红，只改路由防护，不改 Gate/Receipt/权限 |
| `truzhen-cloud` | 身份、session、商品、entitlement 真相 | 原则上不改产品实现；补/复用本地登录链集成验证和 full-stack readiness 证据 | 黄，禁止触生产 |
| `truzhen-packs` | Pack 声明与测试/修复计划 | 仅保存本计划；两个 manifest 和 Pack 主权链不因本问题变化 | 绿 |
| `truzhen-contracts` / `truzhen-software` | 跨仓形状 / Provider | 只读，不改 | 禁止施工 |

依赖方向：`cloud 身份真相 → OS 可信认证代理与 opaque session → client 登录表面`；`OS 06 current run → client 项目投影`。不得让 client 或 Pack 成为真相源。

## 4. 上下文包与禁止边界

施工只允许读取/修改以下现有隔离 worktree：

- client：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-client-web-desktop`
- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhenos`
- cloud：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-cloud`
- packs：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs`

开工基线：client HEAD `0908661c`、OS `df584740`、cloud `06a08a7a`、packs `4e41ae21`；保留四仓现有未提交修复，不覆盖、不回滚。每仓施工前重新记录完整 HEAD、status 和 diff hash。

禁止：改 contracts/schema/DTO；改 Base Gate、Receipt、权限、session 数据形状或主权语义；允许浏览器提供任意认证上游；向生产域发送测试手机号/密码；用 API/curl 写 OS session 或 registry 冒充 GUI；真实支付、外发、Frappe 写入、数据迁移/删除；提交、push、merge、发布；清理其它任务资源。

## 5. 实施步骤

### 任务 A：先建立可观测红灯，不猜竞态

1. 在既有 `frontendRuntimeMonitor` 增加脱敏的 current-run 投影事件：`transaction_ref` 短 hash、请求序号、触发原因、`selection_state/status/occ_version/current_step`、是否因旧响应被丢弃；不得记录凭据或业务原文。
2. 用真实响应序列复现：`none → running/intake → 延迟超过旧窗口 → waiting/owner_gate`，确认旧实现最终仍显示 intake。
3. 增加并发反例：较早发出的 `running/intake` 晚于较新的 `waiting/owner_gate` 返回时，UI 不得被旧响应倒退。

验收：先红后绿；报告能说明每次 UI 更新来自哪次 current-run 响应，而不是只看最终截图。

### 任务 B：把 current run 改成可失效的单一投影

1. 复用现有 `getCurrentSceneFlowRunByTransaction`，不新增端点、不从 run 列表自行挑选。
2. 将 active transaction 的 current-run 读取从“一次性合并进 sessions + 10 秒观察”收敛为一个可取消、带请求序号的刷新器；事务切换、创建完成、页面重新可见、窗口获得焦点、相关受控动作完成时刷新。
3. active transaction 在 `running/waiting/blocked` 时采用有界退避刷新；页面隐藏或事务切换立即取消。超过同步期限时显示“状态同步中/暂不可用”，不得继续把旧 intake 当作当前真相。
4. 只有同一 transaction 且更新的 `occ_version`/请求序号才能覆盖现有投影；旧响应、前一个事务响应必须丢弃。
5. 侧栏、项目头、BOW 卡片和详情继续消费同一个投影映射；`selection_state=ambiguous` 显示歧义，不猜；确无 current run 才回落 05 lifecycle。

定向测试至少覆盖：慢推进超过 10 秒、响应乱序、切换事务、刷新/重新挂载、双 Tab、`ambiguous`、无 current run、接口失败不伪造状态。随后运行 client 相关测试、全量测试、typecheck、static、build 和真实 OS live smoke。

### 任务 C：切断浏览器对密码认证目的地址的选择权

1. client 新增红灯测试：未注入任何认证配置时，密码登录请求体不得出现 `cloud_auth_base_url=https://www.truzhen.com`；带 `?cloudAuthBase=` 也不得改变 Auth Gateway 上游。
2. `AuthPage` 将“扫码/浏览器跳转 URL”与“本地 Auth Gateway 意图”拆开。扫码仍可使用受信任运行时/产品默认 URL；密码、注册、找回和 SMS 只发业务字段到本地 `/v3/auth-gateway/*`，不发送浏览器选择的上游 base。
3. 保留现有 API JSON 形状的向后兼容，不改 contracts；删除/收窄 client 内部可选参数只属于实现清理。
4. 对 return-ticket 单独审计：若仍需 cloud 来源信息，OS 必须用自身可信配置或已签发 ticket 元数据解析，不能恢复 query 参数选择权；本 P0 未证明需要改 ticket 契约时不扩做。

### 任务 D：OS 认证路由 fail-closed 加固

1. 增加 `authintentdev`/组合根测试：client 不传 base、`realE2E=true`、market endpoint 指向本地受控服务时，只调用 `/auth/login`，导入 cloud session，外部仅返回 opaque local session。
2. 增加网络反证：测试内设置生产域陷阱 transport；本地模式登录不得产生任何发往 `www.truzhen.com` 的请求。
3. 普通登录意图中的 `cloud_auth_base_url` 不得覆盖 OS 可信配置。兼容策略优先为“忽略并记录受控诊断”；如现有生产调用依赖该字段，则改为显式拒绝 `client_auth_upstream_override_forbidden`，由独立回归决定，不能静默接受。
4. 保持 prod build 的服务器侧生产默认；保持无真实 endpoint 时 `not_ready`，不得退回假 cloud session；保持 raw password 不落盘、不进日志，session 导入仍需 Base 签发 authority。

定向测试：`authintentdev`、production no-backdoor、market proxy session、cloud session import、allowed endpoint、防生产外发。若仅 client 修复已使真实链通过，OS 仍必须保留上述防回潮测试；是否需要最小产品代码改动由红灯结果决定。

### 任务 E：cloud 与四仓集成证明

1. 不改 cloud 身份/商业契约，启动既有隔离 full-stack，确认 node、market、registry、entitlement healthy。
2. 使用本地受控买家证明 node `/auth/login` 200、market `/auth/login` 200；记录 session 只作脱敏 hash。
3. 启动 OS 时显式配置本地 `TRUZHEN_MARKET_PROXY_BASE_URL`，开启 controlled real E2E；不要设置任何指向生产的 auth base。
4. 从真实 GUI 登录：抓取 client→OS、OS→local market 两段网络证据；生产域陷阱计数必须为 0；OS current session 为 opaque ref。
5. 未登录 `/v3/market/proxy/license/products` 为 401 `session_ref_required`；登录后同一路径 200；退出后再次 401。不得直接把 cloud 200 当 GUI PASS。

## 6. 独立验收与恢复测试门

修复执行者完成后，由独立验收主体按“改了什么证明什么”复核：

- client：相关红灯、全量 test/typecheck/static/build、真实浏览器 current-run 收敛与登录网络证据。
- OS：`authintentdev`、market proxy/session import、scene-flow current/readmodel、production no-backdoor 与生产域零请求反证。
- cloud：full-stack health、真实本地买家登录、entitlement 和 catalog；不以 mock/fixture 代替。
- packs：JSON/结构/forbidden artifacts 保持通过；确认本轮没有 Pack 产品变更。

恢复环保 R1 的硬门：全新 registry=0、全新 OS HOME/SQLite、全新 cloud runtime/DB、全新浏览器 profile；本地 GUI 登录生产域请求=0；opaque session 建立；401→200；随后才允许 GUI 市场安装。不得从已中止 R1 断点续跑。

恢复智能家居 R1 的硬门：全新事务；真实 06 先捕获 `running/intake` 再达到 `waiting/owner_gate`；项目头、侧栏、BOW、详情在 5 秒内一致，等待 30 秒、切页、刷新、重开后不回潮；通过后才恢复 R1 余项，R2/R3 仍按 v3 计划执行。

## 7. 小问题自修与大问题登记

施工/复测中的小问题必须同时满足：局部可逆、归属明确、不改契约/认证信任边界/Gate/Receipt/权限/主权、不触生产和数据迁移，有自动失败测试。循环为：

`small-fix 记录 → 红灯 → 根因 → 最小修复 → 定向回归 → 邻接回归 → GUI 检查点重跑 → PASS 后继续`

任一 P0、认证信任边界新变化、契约/schema/DTO、Gate/Receipt/权限、生产或客户数据、迁移删除、跨仓架构改向、同一根因两次最小修复仍失败，必须登记大问题并冻结危险 lane。禁止通过延长轮询、硬编码环境地址、预写 session/registry、fallback 或降低断言取得假绿。

## 8. 完成物与状态更新

完成物：四仓基线、红/绿测试证据、脱敏网络/monitor 日志、独立回归报告、更新后的两组大问题台账、全新 R1 启动参数单。代码仍保持未提交、未推送，直到 Owner 另行授权。

状态推进：

`设计中（本计划） → 修复获授权 → 已实现 → 已接线 → 独立回归 PASS → 全新 R1/R2/R3 → 已验收`。

在“独立回归 PASS”之前不得启动新的三轮；在 R1-R3 与 Owner UAT 之前不得打包放行。
