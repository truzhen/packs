# 环保执法 Pack 测试计划 v9（15 scope 补偿事务终验版）

日期：2026-07-13
状态：`契约已定 / 修复已实现 / 待全新三轮真实 GUI 验收`
版本定位：当前客户主线；真实证据来自 v8 C01 PDF/秘书只读对账及第 N 个 KnowledgeMount 停用失败会形成半停用的 P1 分析。

## 1. 派活卡与最小闭环

| 维度 | v9 裁定 |
|---|---|
| 最小可交付 | 在全新运行态证明 15 个 mount 正常停用全收口；第 1、8、15 项注入失败均补偿为 15/15 active，14 继续 enabled；随后证明 disabled 后零副作用。 |
| 真相源 | 14 持 Pack lifecycle；09 持 KnowledgeMount；05 持案件绑定；13/08/03/07 分别持秘书会话、模型用量、回执与任务。client 仅投影。 |
| 仓库/层 | OS 实现 14→09 可审计补偿事务；cloud 提供真实本地市场；client 只消费真相；packs 只声明计划与 Pack 资产。 |
| 风险/契约 | Gate、Receipt、法律、模型为红色；AI 仅施工和提出 Candidate。不改 contracts/schema/DTO/主权边界。 |
| 生命周期 | 代码为`已实现`；完成 v9 三轮四路证据前保持`未验收`。 |

非目标/backlog：生产云、真实付款、真实执法/外发、客户材料、Dendrite 全量兼容、提交/合并/推送/发布。

## 2. 权限、上下文与连续推进

执行者拥有四个隔离 worktree、本轮 HOME/DB/Docker/browser profile 的修改、构建、测试和启停权限；必须真实 Chrome GUI 登录本地 cloud、真实 sandbox Entitlement、真实市场制品下载/digest/OS 安装。只读取本计划、v8 报告及相关 05/09/13/14、market、client 模块，不扫描或修改无关 WIP、真实凭据和客户数据。

中小问题自行定位根因、修复、跑定向及邻接回归，然后从原 GUI 证据点继续；大问题登记污染范围、已发生副作用、暂停 lane 与恢复门。安全独立 lane 能继续必须继续，只有身份、权限、Gate、Receipt、跨事务隔离或真实副作用不可信且无安全 lane 时暂停。

## 3. R1：补偿原子性、真实市场与停用零副作用

过程目录：`/Users/li/Documents/过程文档/env-prepack-v9-20260713/`；不得复用 v8 DB、profile、订单、Receipt 或截图。

1. cloud `up.sh` 必须先通过四个宿主 HTTP 硬门；真实 GUI 登录、授权、安装，核对 2 Role Slots、15 scopes、752 knowledge 和制品 digest。
2. enabled 下核对 14 pointer、15/15 active mounts；C01 PDF 与一次秘书请求保持 0 FormalTask、0 外发。
3. 在隔离复制运行态分别注入第 1、8、15 项 disable 失败：接口失败；已停用项逆序补偿；09 最终 15/15 active；14 record/pointer 仍 enabled；没有 `pack_version_disabled`。补偿失败必须显式双错误并冻结，不得声称回滚成功。
4. 正常 GUI 停用：15/15 disabled 后 14 才提交 disabled；GUI、09、14、03 四方同源。
5. disabled 后攻击流式、非流式、动作交接、Owner Confirm：13 turn、08 usage、Candidate、Receipt、FormalTask、发送、执行增量全部为 0；真相不可读则 `not_ready`。

## 4. R2/R3

R2：disable/reactivate/reinstall、安装 owner≠停用 actor、cloud/OS 重启、session 过期、补偿幂等重放、历史 Receipt 可反查。
R3：C01→C08、C01→C02 材料路径隔离、坏/加密/伪装 PDF、双角色 compare Gate、急停、双 Tab、多视口、Provider 四态与 Owner UAT。

开放 P0/阻塞 P1 为 0，R1-R3、补偿矩阵、零副作用增量、真实登录/市场及 Owner UAT 均有新证据，才可判`已验收（打包前）`；否则`未验收 / 禁止打包放行`。
