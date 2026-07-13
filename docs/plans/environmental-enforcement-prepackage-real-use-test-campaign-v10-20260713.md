# 环保执法 Pack 单包跑通测试计划 v10

日期：2026-07-13
状态：`契约已定 -> R1-ENV-003 修复已实现 -> 待全新 C01 与单 Pack 三轮真实 GUI 复验`
版本 / 优先级：当前客户主线。真实场景证据来自 v9：单包市场安装、C01 PDF、双角色候选、秘书候选、停用和停用后零副作用主链已走通，但旧计划把 45 份 Pack 源文档与前台全局 752 条知识混为一谈，且临时 ZIP 混入 `__pycache__/*.pyc`。

## 0. v10 R1 阻断修复后的续跑硬门

- `R1-ENV-003` 已按根因修复：06 `candidate-input` 完成请求校验与 run 定位后，使用脱离 HTTP 取消、默认 5 分钟封顶的编排 Context 同步推进后续模型节点并刷新 Receipt Candidate；不改 contracts、Gate、Receipt 或正式化边界。
- 旧 transaction `tx_f8166ea03c9fce59` 与旧 run `run-environmental-enforcement-flow-1783926016888-2` 只保留为失败证据，不得修库、resume 或复用为 PASS。重启修复后 OS，以全新 C01 transaction/run 从 PDF 回灌点重跑；市场、授权、安装证据只有在仍属本 v10 同一隔离运行态且可反查时才可沿用。
- 回灌后必须在同一时间窗四路对账：GUI 不得继续无期限“进行中”；06 ReadModel 与 SQLite 必须同为 `waiting/compare`，或同为带真实 `blocked_reason` 的 `blocked/not_ready`。不得出现 `running + model_not_ready`、`context canceled` 或前后端状态不同源。
- 自动回归还必须覆盖“HTTP request Context 已取消”攻击；请求断开不影响已授权的 Candidate 编排完成，但超出有界 timeout 时必须诚实阻断，不得后台无界运行。

## 1. 派活卡与最小可交付

| 维度 | v10 裁定与建议 |
| --- | --- |
| 我要做的事 | 只验收 `environmental-enforcement-pack-v0`，从真实 cloud 市场制品、sandbox Entitlement、下载、digest、安装、启停，到 C01 核心业务和急停安全链全部跑通。 |
| 最小可交付 | R1 单包安装与 C01 核心链通过；R2 单包生命周期/重启通过；R3 单包对抗与 Owner UAT 通过。不得以另一个 Pack 是否通过冻结本计划。 |
| 真相源 | cloud 持商品、订单、Entitlement、分发制品；14 持 Pack lifecycle；09 持 KnowledgeMount/FormalKnowledge；05/06 持案件与 Run；08/03/07 持模型用量、回执、任务。Pack 仓只持 45 份源文档与 15 scope 声明，client 只投影。 |
| 仓库 / 层 | packs 生成净化市场制品；cloud 上传前做禁品和路径审计；OS 消费制品并提供真实 lifecycle/Gateway；client 只做真实 GUI。 |
| 风险颜色 | 文档/制品审计为绿；Pack 生命周期为黄；市场/Provider 为橙；Gate、Receipt、急停、法律、真实动作均为红。AI 是施工者与测试者，但仍只能产 Candidate。 |
| 契约影响 | 只改实现、构建审计和测试口径；不改 contracts/schema/DTO、Receipt/Gate 主权边界。 |
| 生命周期 | 修复为`已实现`；三轮新鲜证据齐全才升为`已验收（打包前）`，不宣称已发布。 |

明确砍掉 / backlog：生产云、真实扣款、真实执法、真实外发、客户生产资料、发布/上架、提交/合并/推送、Dendrite 全量兼容。法律知识一律保持 `pending_human_review`。

## 2. 权限、上下文和禁止边界

执行者已获本任务隔离 worktree、专用 HOME/DB/Docker/browser profile、服务启停、代码小中修复和定向测试权限；必须真实 Chrome GUI 登录本地 cloud、使用本地 sandbox 模拟支付授予 Entitlement、从云市场点击下载和安装。禁止复用 v9 DB、登录态、订单、商品版本、截图、Receipt、案件或 Run；禁止读取真实 secret、修改生产配置、真实扣款、正式发送/执法、覆盖他人 WIP。

上下文包限于本计划、v9 报告、Pack 资产及 cloud 市场制品审计、OS 05/06/08/09/14、相关 client 页面。单包 lane 中禁止安装智能家居 Pack；系统基础数据必须单独标注，不得混入本 Pack 计数。

## 3. 连续推进纪律

- P2/P3 及范围明确、可逆的中等问题：执行者自行定位根因、修复、跑定向与邻接回归，从原 GUI 证据点继续，不得因小问题把任务交回 Owner。
- P0/P1：登记污染范围、已发生副作用、冻结路径和恢复门。只要身份、权限、Gate、Receipt、跨事务隔离仍可信，安全独立 lane 必须继续跑；只有没有安全 lane 时才暂停。
- 禁止 mock/API 注入代替 GUI 成功，禁止跳过错误、改断言、用旧证据或 fallback 冒充 PASS。

## 4. R1：净制品、真实市场、知识归属与 C01

过程目录新建 `/Users/li/Documents/过程文档/env-prepack-v10-20260713/`。

1. 用 `python3 -B build_pack_bundle.py --market environmental-enforcement-pack-v0 <新输出目录>` 生成新制品和 sha256 清单；ZIP 根有 `manifest.json`，且 `__pycache__`、`*.pyc`、数据库、日志、构建目录、符号链接、路径穿越均为 0。cloud 上传前必须再次审计，旧临时 ZIP 不得复用。
2. 全新 cloud 商品/version、真实 GUI 登录；先证明未授权下载/安装被阻断，再在 GUI 点击本地 sandbox 模拟支付（明确无真实扣款），得到新订单/Entitlement，点击下载与安装；下载哈希=上传制品哈希，registry 0→1。
3. 资产对账使用三个互不混淆的口径：Pack ZIP 中 45/45 `knowledge-index` 源条目与 Markdown 一一对应；15/15 scope 声明与 mount 一一对应；09 按 `pack_ref + pack_version_ref` 全量分页，45 个 versioned `source_ref` 全覆盖、无外 Pack 混入。运行态分片数只记录观测值，不把历史全局 752 写成固定契约。
4. 核对 2 Role Slots、Role Pack 与绑定；新建 C01，真实 GUI 上传合成 PDF，解析状态回灌同一 06 Run；按第 0 节核对 GUI/06/SQLite 一致后，本地真实模型生成双角色候选，compare Gate 保持 waiting/Owner gate，0 FormalTask、0 外发、0 执法执行。
5. enabled 下秘书只允许一次受治理调用；13 turn、Candidate、08 usage、03 Receipt 与同一 transaction/Pack 归属可反查。

R1 通过门：制品净化、真实登录/市场、45 源覆盖、15 mounts、C01 PDF、双角色、秘书和零正式副作用全部有 GUI + 05/06/08/09/03 四路证据。

## 5. R2：单包生命周期与恢复

只对环保 Pack 执行 disable→reactivate→disable→reactivate；第 1/8/15 mount 注入失败分别证明逆序补偿后 15/15 active、14 仍 enabled，补偿失败必须双错误。正常停用必须先 15/15 disabled 再提交 14 disabled。覆盖 OS/cloud 重启、session 过期重登、安装 Owner 与操作者差异、幂等重放和历史 Receipt 反查。停用后 PDF、秘书、角色候选、Owner Confirm 攻击均保持模型、Candidate、Receipt、FormalTask、发送、执行增量 0。

## 6. R3：单包安全与用户验收

覆盖 C01→C02/C08 跨事务隔离、坏/加密/伪装 PDF、重复 PDF 幂等、双 Tab、390×844/1024×768/1440×900、Provider `ready/not_ready/blocked/degraded`、急停与解除。急停期间 PDF/秘书/角色路径必须 0 Provider、0 stream、0 08 usage、0 模型 Receipt Candidate、0 Formal Receipt、0 ModelRun、0 任务/记忆/发送/执行；解除后仅新请求恢复。最后由 Owner 逐项 UAT。

## 7. 验收与收尾

专门验收主体须复核：制品 sha256/禁品、45 个源引用覆盖、15 mount、GUI 录屏/截图、05/06/08/09/03 增量表、跨事务攻击和端口/容器清理。开放 P0=0、阻塞 P1=0，R1/R2/R3 均 PASS，才判该单 Pack `已验收（打包前）`；否则诚实标记未验收，但不得用另一个 Pack 的问题否决已完成证据。
