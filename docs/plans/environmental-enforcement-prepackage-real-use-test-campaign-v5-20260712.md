# 环保执法 Pack 打包前真实使用测试计划 v5（PDF 回灌可观察性修复续跑版）

日期：2026-07-12
状态：`契约已定 / 待执行`
定位：当前客户主线；v4 的完整 R1-R3 用例继续有效，本文件覆盖 v4 的续跑入口与 `ENV-V4-P1-001` 专项。

## 1. 派活卡

| 维度 | v5 裁定 |
|---|---|
| 真实证据 | v4 retry1 GUI 生成 PDF Candidate 后长期“解析中”；最终 SQLite 证明 `pdf_parse_status=parsed` 且 Run 已到 compare，属于前台可观察性/同步等待故障，不是 06 数据丢失。 |
| 最小可交付 | 全新 R1 先证明 PDF 回灌在有界时间内可判定，再完成 C01/C02 隔离与冷启动手工回灌；通过后续跑 R2/R3。 |
| 真相源 | PDF Candidate 归执行链；是否已回灌归 06 run timeline/state；client 只根据 06 证据显示已回灌或阻断。 |
| 仓库归属 | client 修有界等待与 06 回读；OS 保持 candidate-input 先持久化再推进；packs 持测试计划；cloud 只负责市场链。 |
| 风险/契约 | 黄级实现修复；不改 contracts/schema/DTO、Gate、Receipt、Candidate/Formal 或模型 Provider 契约。 |
| 禁止边界 | 不用 API/SQLite 补写状态，不干扰共享模型，不伪造 timeline/ref，不提交/合并/推送，不触生产/真实执法/外发。 |
| 生命周期 | 自动回归仅为 `已接线`；新鲜 GUI 四路证据通过后才为 `已验收`。 |

新过程目录：`/Users/li/Documents/过程文档/env-prepack-v5-20260712/`。使用 v4 计划列出的四个修复 worktree，但必须记录新的 SHA/WIP/diff hash，并使用全新 HOME、DB、profile、origin、Docker project、C01/C02、Run、Candidate、Receipt。

## 2. R1 首门：PDF Candidate 回灌可观察性

1. 完成 GUI 登录、paid sandbox entitlement、安装与 registry 0→1，对账 2 角色、15 mount、752 knowledge。
2. GUI 创建新 C01，通过本地路径生成新 PDF Candidate。记录 candidate ref、receipt ref、transaction/run ref 与点击时间。
3. candidate-input 正常快速返回时，GUI 显示“已回灌到当前流程”；若后续模型使请求超过 12 秒，GUI 必须退出“解析中”，回读 06 后仅在 timeline 存在 `candidate_input_applied_*` 且 evidence 含本 PDF candidate ref 时显示“已回灌、后台处理中”。
4. 四路对账：GUI、candidate-input network、06 run ReadModel/timeline、SQLite state 必须一致包含本 PDF candidate 与 `pdf_parse_status=parsed`。不得只看按钮文案。
5. 若回读无该 candidate 证据，必须显示受控阻断，不得继续 C01；登记新问题。不得把 HTTP timeout 本身等同回灌失败。
6. 等待/轮询至 compare；fault proxy 初始 not_ready、后切真实模型的 P1 专项继续按 v4 执行。任何迟到响应不得使 GUI 倒退为解析中。

## 3. 后续三轮与问题处理

R1 继续执行 v4 的 P0-004 C01→C02→C01 隔离、双 Tab、EmergencyStop 零副作用和 P1 手工回灌。R1 全门通过方可进入 R2；R2/R3 完整继承 v4。

中小问题在四仓既有授权内按“红灯→最小修复→定向/邻接回归→原 GUI 检查点复跑→继续”处理。大问题必须登记影响、污染、允许继续 lane、暂停 lane 与恢复门槛；独立安全 lane 可继续，证据归属污染、主权/权限/急停穿透或后续依赖故障链时暂停。

放行仍要求 R1-R3、LX、跨事务隔离、PDF、Candidate/Formal、Gate、Receipt、急停、Owner UAT 全绿；开放 P0/阻塞 P1 为 0。否则 `未验收 / 禁止打包`。
