# 环保执法 Pack 测试计划 v6（受控 GUI 恢复与 PDF 四路终验版）

日期：2026-07-13
状态：`契约已定 / 等待受控 GUI 会话`
定位：当前客户主线；v5 自动修复已接线，但 GUI 终验受 Codex 浏览器 loopback 安全策略阻断。本计划不把环境阻断冒充产品故障或产品通过。

## 1. 派活卡

| 维度 | v6 裁定 |
|---|---|
| 真实证据 | v5 C01 的 06/SQLite 已有 `candidate_input_applied_*` 与 PDF 状态；前端因仅读 snake_case 误报阻断，现已兼容运行端 PascalCase。Codex 浏览器仍拒绝宿主 `127.0.0.1`。 |
| 最小可交付 | 在产品策略允许的受控 GUI 会话中，用新事务完成 PDF GUI/network/06/SQLite 四路复验，再继续 C01/C02 隔离、角色回灌和 R2/R3。 |
| 真相源 | 是否回灌归 OS 06 timeline/state；client 只投影；浏览器策略归测试基础设施，不归 Pack 或业务代码。 |
| 仓库归属 | client 已修字段兼容；OS 保持 06 真相；packs 持计划；cloud 仅市场。不得为了浏览器可达性修改产品网络边界。 |
| 风险/契约 | 前端兼容为黄；不改 contracts/schema/DTO、Gate、Receipt、Provider、Candidate/Formal。 |
| 禁止边界 | 禁止绕过 URL 安全策略、改用未授权浏览器/CDP、绑定 `0.0.0.0`、公网暴露、隧道或 API/DB 冒充 GUI PASS。 |
| 生命周期 | 代码为 `已接线`；四路 GUI 证据前保持 `未验收`。 |

新过程目录：`/Users/li/Documents/过程文档/env-prepack-v6-20260713/`。全新 HOME、SQLite、cloud project、browser profile、origin、C01/C02、Run、Candidate、Receipt；使用既有四个修复 worktree并记录 SHA/WIP/diff hash。

## 2. GUI 恢复硬门

允许的入口只有两类：Codex 产品层明确允许本任务访问隔离 loopback；或 Owner 明确提供并亲自打开符合隔离策略的受控 GUI 会话。执行者只能认领已允许标签页，不能自行建立转发或扩大监听面。

若仍被 URL 策略拒绝：记录原始错误、监听与清理证据，GUI lane 标 `environment_blocked`；可继续自动测试、静态审计和只读 DB 对账，但不得关闭 P1、不得进入依赖 GUI PASS 的后续轮次。

## 3. PDF 四路首门

1. GUI 安装环保 Pack 并创建新 C01；通过本地受控路径生成 PDF Candidate。
2. 12 秒内正常返回则显示已回灌；超时则必须回读 06。运行端可能返回 `EventType/EvidenceRefs`，历史投影可能为 `event_type/evidence_refs`，两者均须识别。
3. 只有 `candidate_input_applied_*` 的 evidence 含本次 source candidate ref 时，GUI 才显示“已回灌、后台处理中”并退出“解析中”；不匹配时保持受控阻断。
4. 保存四路证据：GUI 读屏/截图、candidate-input 超时与 readmodel 网络原文、06 run ReadModel、SQLite run_event/state。四路 transaction/run/candidate/receipt 必须一致。
5. 反向用例：提供另一 Candidate 的 timeline、空 timeline、ReadModel 不可用，均不得显示已回灌。

首门 PASS 后继续 v5/v4 的 C01→C02→C01 跨事务隔离、fault proxy 冷启动手工双角色回灌、compare Gate、EmergencyStop 零副作用，再进入 R2/R3。

## 4. 问题处理与放行

中小问题在四仓授权内按“红灯→根因→最小修复→定向/邻接回归→原 GUI 点重跑→继续”处理。大问题登记污染、允许继续 lane、暂停 lane和恢复门槛；安全独立 lane 可继续，证据归属、主权、权限、急停或正式副作用不可信时暂停。

放行仍要求 R1-R3、LX、PDF 四路、跨事务隔离、Candidate/Formal、Gate、Receipt、急停与 Owner UAT 全绿；开放 P0/阻塞 P1 为 0。否则 `未验收 / 禁止打包`。
