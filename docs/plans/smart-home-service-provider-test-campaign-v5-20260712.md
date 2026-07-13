# 智能家居服务商 Pack 测试计划 v5（EmergencyStop 单一真相续跑版）

日期：2026-07-12
状态：`契约已定 / 待执行`
定位：当前客户主线；v4 R1 与 R2 的既有新鲜证据保持有效，本文件覆盖 R3 重启入口和 `SH-V4-P1-03` 专项，不追认原冻结 R3。

## 1. 派活卡

| 维度 | v5 裁定 |
|---|---|
| 真实证据 | v4 R3 顶部显示急停已生效，安全核心 `SecurityCoreReadModel` 同时显示 `standby_preview`。 |
| 最小可交付 | 用全新 R3 运行态证明顶部、Base readmodel、安全核心 readmodel/GUI 全部消费同一 EmergencyStop 真值，再完成剩余 R3 对抗与 UAT。 |
| 真相源 | EmergencyStop 唯一真相归 Base store；顶部和 securitycoredev 均只能投影，client 不得合成或保留 preview。 |
| 仓库归属 | OS securitycoredev/组合根负责真值接线；client 只显示；packs 持计划。 |
| 风险/契约 | ReadModel 实现为黄，急停动作本身为红；不改 contracts、Gate/Receipt、启停授权或权限语义。 |
| 禁止边界 | 未完成三方对账前不正式化、不发送、不执行、不写 Provider、不解除急停；不触生产/客户数据，不提交/推送。 |
| 生命周期 | 自动回归只到 `已接线`；新鲜 GUI 终轮后才关闭 P1。 |

新过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v5-20260712/`。使用现有四个修复 worktree，记录全新 baseline；R3 使用新 HOME、DB、profile、origin、Pack 安装、事务、Receipt。v4 R1/R2 报告只读，不复制运行态。

## 2. R3 首门：EmergencyStop 三方一致

1. 急停未启用时，同时抓取 `/v3/base/readmodel/security-core`、`/v3/security-core/readmodel`、顶部面板和安全核心页；四处必须一致为 `emergency_stop_disabled`，不得出现 `standby_preview`。
2. 仅从 GUI 启用急停，保存 Base decision/candidate 证据；再次抓取四处，必须一致为 `emergency_stop_enabled`。刷新、切页、窗口失焦/回焦、浏览器重开、OS 重启后仍一致。
3. 注入只影响本轮隔离 OS 的 Base store 读取失败，两个 readmodel 必须诚实 `backend_unavailable`，GUI 不得显示 disabled/standby；恢复后重新对账。
4. 急停 enabled 下验证模型、PDF、任务正式化、沟通、发送、执行、Frappe 写回的 usage/Candidate/Receipt/Formal/外部副作用增量均为 0。
5. 三方一致和零副作用通过前，解除急停入口保持冻结；不得以顶部单一路径 PASS 替代安全核心。

## 3. 剩余 R3 与问题处理

首门通过后继续 v4 R3：伪 ref、非白名单发送、多视口、双 Tab/跨项目隔离、Dendrite readiness、停用历史回放与 Owner UAT。解除急停仍必须走 Owner + Base Gate；解除后只恢复候选能力，不自动产生正式任务、发送或 Provider 写回。

中小问题按“红灯→最小修复→定向/邻接回归→GUI 复跑→继续”自行处理。大问题登记影响与污染；独立安全 lane 可继续，若状态真相、身份、Receipt 或危险动作证据不可信则暂停。Dendrite 未就绪只收窄结论，不得伪称联邦全量 PASS。

放行要求：v4 R1/R2 证据有效、v5 R3 全门通过；current-run 同源、跨项目隔离、Provider 诚实、Candidate/Formal、EmergencyStop、Gate、Receipt、Owner UAT 全绿；开放 P0/阻塞 P1 为 0。否则保持 `未验收 / 禁止打包`。
