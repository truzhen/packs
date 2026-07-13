# 环保执法 Pack 打包前真实使用测试战役计划 v3（三轮持续修复版）

日期：2026-07-11
状态：`契约已定 / 待执行`
归属：当前客户主线，打包前放行门
目标：使用修复后的 Truzhen 与真实 `environmental-enforcement-pack-v0`，从空 registry 开始连续跑 R1-R3；小问题在隔离分支自行修复并复验直至通过，大问题登记并冻结危险路径。

## 1. 权威输入与本版关系

- v2 用例库：`/Users/li/Documents/truzhenv3worktree/truzhen-packs-env-prepackage-v2-20260711/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v2-20260711.md`
- R1 阻断报告：`/Users/li/Documents/过程文档/env-prepack-v2-20260711/three-rounds-20260711/R1/R1-round-report.md`
- P1 lane：`/Users/li/Documents/过程文档/env-prepack-v2-20260711/three-rounds-20260711/R1/lanes/R1-E1-GUI-market-blocked.md`
- 旧问题台账：`/Users/li/Documents/过程文档/env-prepack-v2-20260711/three-rounds-20260711/issues/environmental-enforcement-v2-register.md`
- 本次修复设计：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs/docs/plans/two-packs-rerun-p1-repair-design-20260711.md`
- 本次实施计划：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs/docs/plans/two-packs-rerun-p1-repair-implementation-plan-20260711.md`

v2 的 E0-E12、LX、C01-C08、752 知识、PDF、法律时点、急停、生命周期和三视口仍是完整用例库。本文件是新三轮的执行权威；与 v2 的“发现问题只登记、不边测边修”冲突时，以本文件第 5 节为准。

## 2. 派活卡

| 维度 | v3 裁定 |
|---|---|
| 我要做的事 | 修复后重新验证 GUI 市场获取、安装、环保办案、知识、PDF、角色、Gate、Receipt、急停、生命周期和回放。 |
| 真实场景证据 | 上轮从 registry=0 进入 GUI 市场后没有可安装环保 Pack；cloud 审批桥未 ready，client 受控市场路径被 cloud 基址劫持。 |
| 最小可交付 | R1 必须完成真实商品→entitlement→GUI 安装→registry 0→1→C01；R2 完成恢复/幂等；R3 完成对抗/UAT。 |
| 真相源 | cloud 持商品、artifact、entitlement；OS 持受控代理、registry、lifecycle、Gate、Receipt；client 只投影；packs 持 Pack 声明。 |
| 仓库归属 | client 修市场路由；cloud 修审批桥 readiness 与测试商品准备；OS 维持受控代理/lifecycle；packs 提供真实 artifact 和计划。 |
| 风险颜色 | 普通 UI、脚本、测试为黄；contracts/Gateway 为橙；Gate、Receipt、权限、生产动作、数据删除为红。 |
| 契约影响 | 不改契约；Pack manifest 只消费既有 `kind/min_truzhen_version` 市场字段。 |
| 禁止边界 | 不触生产、不用真实客户、不真实支付/处罚/移送/外发，不改 contracts，不绕过 GUI 写 OS registry，不提交/推送/合并。 |
| 生命周期上限 | 自动化修复为 `已接线`；三轮 GUI 与 Owner UAT 全绿后才为 `已验收（打包前）`。 |

## 3. 新过程目录与基线

本轮只写新目录：

`/Users/li/Documents/过程文档/env-prepack-v3-20260711/`

固定结构：

```text
baseline/
R1/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
R2/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
R3/{parameters,lanes,shots,network,receipts,readmodels,monitor,small-fixes}/
issues/big-issue-register.md
final-report.md
```

每轮必须记录六仓完整 SHA、分支、WIP、未提交 diff hash、二进制来源、HOME、DB、端口、Docker project、浏览器 profile、账号档位和时间。当前修复未提交时，禁止从旧测试 worktree 开跑；必须使用以下修复 worktree 的实际文件并登记 diff：

- client：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-client-web-desktop`
- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhenos`
- cloud：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-cloud`
- packs：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs`

contracts、software 只读；每仓独立 `git status --short --branch`。测试修复只能进入上述隔离分支，不能改 main。

## 4. R1 开跑硬门：真实市场获取链

### 4.1 cloud 启动

使用独立 Docker project/runtime/端口启动 full-stack；必须显式开启本地支付模拟，且所有回执保留 `local_payment_simulation=true` 与 `must_replace_before_upload=true`。审批桥必须先达到 `healthy`，不健康时启动脚本必须失败，禁止继续市场 lane。

必须证明：

1. node `/health` 返回 200、`database=ok`。
2. market health 200，`pack_registry`、commerce、license ready。
3. 使用真实环保 Pack 压缩产物，根目录含 `manifest.json/install.py/uninstall.py`；不得使用合成 manifest 或空包。
4. 通过 `scripts/ensure-pack-e2e-market-product.sh` 准备 paid sandbox 商品和明示模拟 entitlement；只写 cloud/Forgejo，不写 OS registry。
5. catalog 中 `pack_environmental_enforcement_v0@1.0.0` 同时为 `listed / has_forgejo / has_commerce`，artifact 下载 200 且 digest 与上传物一致。

### 4.2 client 与 OS 路由

1. client 的 `/v3/market/proxy/*`、`/v3/market/admin-proxy/*` 必须走 OS base；`packMarketBase` 只可影响旧 `/pack-proxy/*`。
2. 禁止把 cloud `18087` 直接配置为受控 V3 市场代理基址。
3. OS 无 session 请求 `/v3/market/proxy/license/products` 必须 401 `session_ref_required`；GUI 登录后同一路径才可 200。
4. cloud 直连接口只用于环境 readiness 取证，不可代替 GUI 获取/安装 PASS。

### 4.3 registry=0 与副作用门

开 GUI 前证明 OS registry=0、Pack/Role/KnowledgeMount/FormalKnowledge 均无本轮残留。cloud 商品与 entitlement 已存在不等于本地已安装；只有用户在 GUI 点击获取/安装并取得 lifecycle Receipt 后才允许 registry 0→1。

任一硬门不满足：先按第 5 节判断小/大问题；不得用 `install.py`、curl 或 DB 写入代替 GUI 安装。

## 5. 测试中持续修复机制

### 5.1 小问题：自行修复直至该检查点通过

满足全部条件才是小问题：

- 影响局部、可逆，真相源和仓库归属明确；
- 不改 schema/DTO/contracts，不改 Gate、Receipt、权限、认证、主权链；
- 不触生产、真实外部副作用、数据迁移或删除；
- 可在 client/cloud/OS/packs 现有实现或测试内最小修改完成；
- 有可自动复现的失败测试。

执行循环：

`登记 small-fix → 最短复现 → 失败测试 → 根因 → 最小修改 → 定向测试 → 邻接回归 → 从失败 GUI 检查点重跑 → PASS 后继续本轮`

小问题不得口头略过。每条写入本轮 `small-fixes/SF-NN.md`，包含 diff、命令、红/绿证据、重跑截图和 neuter/反证。允许同一轮多次修复，直至用例通过。

### 5.2 大问题：登记，不擅改

任一条件即为大问题：

- P0，或涉及 contracts/schema/DTO、真相源迁移、Gateway 边界；
- Gate、Receipt、权限、认证、生产支付/发送/执行、客户数据；
- 跨仓架构改向、数据库迁移/删除、不可逆外部状态；
- 根因不明确，或同一根因连续两次最小修复仍失败；
- 修复需要扩大当前四仓授权或触碰 main/生产。

写入 `/Users/li/Documents/过程文档/env-prepack-v3-20260711/issues/big-issue-register.md`，状态使用：

`已登记 → 已复现 → 根因已确认 → 待 Owner 裁定 → 修复获授权 → 已实现 → 已接线 → 独立回归 PASS → 终轮 PASS`

冻结受影响危险 lane，保留现场；没有依赖关系的安全 lane 可继续取证，但该轮不得判总 PASS。禁止 fallback、mock 成功或降低断言。

### 5.3 已知大问题观察项

本次修复验证发现：免费 Forgejo 商品经 node 下载时仍可能被许可 Gate 判为 `payment_required`。它涉及 Gate 语义，本计划不授权修复。三轮主路径固定使用 paid sandbox 商品 + 明示模拟 entitlement；另在 R3 负向 lane 复现并登记，不把免费路径结果混入 paid entitlement PASS。

## 6. 三轮执行

### R1：全新市场获取与 C01 核心闭环

1. 使用全新 OS DB、全新 cloud market DB、全新浏览器 profile；registry=0。
2. 完成第 4 节全部硬门，用户从真实 GUI 市场找到环保 Pack，查看商品/授权状态并点击安装。
3. 对账 registry 0→1、Pack enabled version、2 Role Pack、slot binding、15 KnowledgeMount、752 FormalKnowledge、lifecycle Receipt；不得预装。
4. GUI 新建 C01，真实 run 走至 compare/owner gate；执行 v2 `ENV-3R-01~09`、E1、E3-E6、E9 与 LX。
5. 重复 PDF 验证同 Candidate/Receipt 与 `replayed`；急停下 PDF/沟通/任务/正式化零副作用。
6. 只有 GUI + network + Receipt + ReadModel/DB 四路一致才 PASS。

R1 阻断门：市场商品不可见、OS 受控路径未走通、GUI 安装不能 0→1、752 对账失败、C01 主链走不完、急停/Gate 穿透，均不得进入 R2。

### R2：恢复、幂等、生命周期与竞态

1. 新 transaction，不复用 R1 run/Receipt。
2. 双击上传、断网重试、刷新、浏览器重开、OS 重启、cloud/node 重启、disable/re-enable/reinstall、双 Tab 同时推进。
3. 重启后审批桥必须由 health gate 控制；不得复现 `pack_registry_unreachable` 后继续测试。
4. 执行 v2 E1、E3-E5、E9-E12 与 LX；重点验证 Receipt replay、OCC、版本固定、历史回放、停用时知识不可新用。
5. 所有小问题按第 5.1 节修到相关检查点 PASS；出现大问题按第 5.2 节冻结。

### R3：对抗、法律时点、三视口与 Owner UAT

1. 新 transaction；覆盖坏/加密/伪装/超限 PDF、无命中、冲突时点、伪 ref、非白名单发送、停用后访问。
2. 急停请求使用无效 body，证明门控先于解析、候选、Receipt 和外部副作用。
3. 390×844、1024×768、1440×900 与双 Tab 各跑关键旅程；关键入口必须可达。
4. 复现免费商品 Gate 观察项并只登记；paid entitlement 主路径必须仍可下载/安装。
5. Owner 亲手复走：`市场获取 → 安装 → C01 → PDF → 三角色引用 → compare Gate → Receipt → 急停 → 停用后历史回放`。

## 7. 每轮证据与停止条件

每轮必须新增 baseline、lane、截图/读屏树、网络原文、Receipt、ReadModel/SQLite、monitor doctor、small-fix 日志、big issue 台账引用、端口/容器释放记录和轮报告。旧 transaction、旧 Receipt、静态测试、fixture 或 API-only 不能替代 GUI PASS。

轮末必须停止本轮服务，证明所有登记端口释放、容器归零；不清理其他任务资源。若有未关闭大问题，结论写 `未验收 / 禁止打包`，不得用后续轮稀释。

## 8. 放行标准

- R1-R3 必测项全部有新证据；LX PASS。
- P0=0、阻塞性 P1=0；big issue 均关闭或有 Owner 书面豁免、替代控制和到期日。
- Candidate/Formal、急停、Gate、Receipt、跨案隔离全程成立。
- paid sandbox 商品、模拟 entitlement、artifact digest、GUI lifecycle 四层一致。
- 小问题均有红/绿测试和 GUI 重跑证据；不存在“修了但没重跑”。
- Owner UAT 通过；否则生命周期保持 `已接线` 或 `未验收`，禁止打包。

## 9. 启动前 500 字简报模板

只汇报：本轮编号、四仓修复分支与 diff hash、六仓 SHA、registry/DB 是否全新、cloud node/market readiness、真实商品与 entitlement 是否 ready、OS 受控路由是否 401→登录后 200、端口/profile、已知大问题、是否允许进入 GUI。
