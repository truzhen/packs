# 环保执法 Pack 直接 lifecycle 单项目全程测试计划 v14

日期：2026-07-14

状态：`当前客户主线；ENV-V13-P1-001 已实现 -> 已接线，待真实 GUI 复验；未验收、未发布`

前序证据：`/Users/li/Documents/过程文档/env-prepack-v13-20260714/R1/r1-direct-lifecycle-and-c01-run-20260714.md`、`/Users/li/Documents/过程文档/env-prepack-v13-20260714/R1/issue-register.md`

被替代计划：`/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v13-20260714.md`

## 1. 派活卡与本轮裁定

| 维度 | v14 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，只验 `environmental-enforcement-pack-v0`；智能家居 Pack 独立验收、互不冻结。 |
| 我要做的事 | 跳过市场交易链，使用本仓净制品和真实 OS lifecycle 直接安装；以一个 C01 项目贯穿安装、启用、项目运行、停用、历史查询、重新启用、卸载、卸载后历史查询和重装恢复。 |
| 真实客户 / 场景证据 | v13 已证明直接安装、15 个 KnowledgeMount、C01 创建、停用、历史查询及重新启用可工作；同时发现 PDF 前角色模型越出事实锚定，生成未提供事实、法条适用和处置建议。 |
| 最小可交付 | 同一 Pack v1.0.0、同一安装实例和同一 C01 完成一条可回放 lifecycle；PDF 前角色入口确定性零副作用阻断，PDF 后候选只引用已提供事实；停用/卸载不删历史，重装不迁移旧项目版本。 |
| 本轮明确砍掉 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部记为 `not_tested_owner_directed`，不启动 cloud、不建立商品/订单/授权，不以缺少市场证据阻断本轮。生产上架、真实付款、真实执法、正式处罚、外发和客户生产材料仍不在范围。 |
| 真相源 | Pack 文件与 digest 归本仓；安装/启停/卸载状态归 14；知识挂载归 09；项目归 05；Run 归 06；角色与模型归 13/08；候选/回执归 07/03。历史项目不因 Pack 停用或卸载被重写。 |
| 仓库 / 层归属 | `truzhenos` 负责 06 EvidenceGap 前置门禁与 lifecycle；client 负责真实状态投影；本仓负责 Pack 资产、净制品与计划。cloud、software、contracts 本轮不改。 |
| 风险 / AI 角色 | 黄：安装、候选、启停；橙：Provider；红：Gate、Receipt、正式法律裁定、发送和执行。AI 只作施工者和 Proposer，不取得正式裁定权。 |
| 契约影响 | 只改 OS 内部实现与测试，不改 contracts、schema/DTO、Candidate、Gate、Receipt、ReadModel、ProviderRequirement 或主权链。 |
| 上下文与禁止边界 | 只读本计划、该 Pack、隔离 HOME/SQLite、合成 C01 材料和相关页面；不读真实凭据/客户材料，不改数据库，不伪造运行态，不覆盖既有 WIP，不正式化、不外发、不执行执法动作。 |
| 变更影响 | 仅角色候选运行前门禁、client 投影复验、Pack lifecycle 测试与证据文档；不改变已安装 Pack 内容和历史 Receipt。 |
| 生命周期档位 | 修复当前为`已实现 -> 已接线`；直接 lifecycle 全程和独立复核通过后才升为`已验收（打包前）`；发布另行裁定。 |

## 2. 执行授权、隔离与直接安装口径

执行者已获隔离 worktree、专用 HOME/SQLite、受控浏览器、真实本地模型、服务启停、端口、截图、只读对账以及不改契约/主权链的中小问题自治修复权限。每仓仍须独立登记 branch、HEAD、`git status --short --branch`、既有 WIP、运行参数和验证命令；不得提交、合并、推送或覆盖他人改动。

1. 从 `/Users/li/Documents/truzhenv3worktree/two-packs-v12-repair-20260713/truzhen-packs/environmental-enforcement-pack-v0` 构建净本地制品，登记 SHA-256、文件清单和 forbidden artifact 扫描。
2. 直接通过 Pack 自带 `install.py` 对隔离 OS 的真实 lifecycle 端点安装；这不是市场下载的替代 PASS，而是 Owner 指定的本轮安装入口。若产品已有本地 Pack 安装 GUI，优先从 GUI 触发同一 lifecycle；没有 GUI 时脚本执行、14 真相和前台状态三方对账即可。
3. `uninstall.py` 的权威语义是撤销当前启用/挂载，不物理删除 05 项目、03 Receipt 或历史版本引用；不得把“数据库删除”写成卸载成功条件。
4. cloud 市场、登录、提交审核、模拟支付、Entitlement、下载不得启动、补测或伪造；报告固定写 `not_tested_owner_directed`。

## 3. 单一 C01 与尽量跑完全程纪律

唯一主项目为 `C01-v14-固废露天堆放线索核查`。同一 `transaction_ref`、Pack v1.0.0、安装实例、浏览器 profile 和证据目录贯穿 R1/R2/R3；不得以新项目清空历史状态。只有真相被缺陷污染且无法恢复时才可建立一个 `C01-RECOVERY-01`，并登记旧/新 refs、污染范围和迁移原因。跨项目泄漏只允许最小只读探针，不得建立第二套验收主链。

遇到 P2/P3 或影响清晰、可逆、不改契约/主权链的中等问题，执行者必须自行定位根因、修复、补测试并从最近可靠点继续。P0/P1 必须登记和冻结受影响路径，但安装审计、历史查询、响应式、只读对账等独立安全 lane 继续；只有所有剩余 lane 都依赖失真状态时才暂停。不得降低断言、吞错、改库、复用旧 PASS 或用重建项目掩盖问题。

## 4. 开跑恢复门

1. Pack JSON/Python/结构校验、45 条知识源与 15 scope 一致性、禁品扫描、ZIP 路径穿越/符号链接扫描、digest 和 `git diff --check` 全绿。
2. OS 定向测试必须证明：06 Run 存在 blocked `EvidenceGapCandidate(pdf_parse_status)` 时，角色入口返回 `role_review_evidence_not_ready`，且 `model_call_performed=false`、`candidate_created=false`、`receipt_created=false`；不得仅靠 prompt。
3. PDF 回灌并清除该 EvidenceGap 后，角色候选才可调用真实模型；停用、版本不匹配或急停仍在 07/03/08 前 fail-closed。
4. OS/client 定向与邻接回归、typecheck/build、Pack 文档级验证全绿后再起 GUI。

## 5. R1：直接安装、启用与同一 C01 核心运行

过程目录：`/Users/li/Documents/过程文档/env-prepack-v14-20260714/`。

1. 新隔离态确认 registry=0、KnowledgeMount=0；执行直接安装，验证 registry 0→1、Pack enabled v1.0.0、2 个 Role Slot/Role Pack、15/15 KnowledgeMount、45/45 本 Pack 知识源及安装 Receipt 可反查。
2. GUI 创建唯一 C01，记录 05 事务、06 Run、Pack version 和初始 03/07/08/09 基线。PDF 前主动点击 advice/challenge，必须被 `role_review_evidence_not_ready` 阻断且业务计数全 0。
3. GUI 上传合成 PDF，验证解析 Candidate/Receipt、`pdf_parse_status=parsed` 回灌同一 06 Run；随后生成执法精英与挑剔律师候选，逐句审计只引用 C01 已给事实，未知项显式待核验，不新增主体、污染介质、监测数值、违法认定、法条适用或处置结论。
4. 对照门展示双方候选并安全停在 compare / Owner Gate。逐个核对 flow 的 intake、screening、inspect_prep、credentials、onsite、scan、evidence_chain、evidence_review、fact/legal advice、challenge、compare 及后续受控门状态；后续正式化、文书发送和执法执行只验证受控阻断，保持 0 FormalTask、0 外发、0 执法执行。

## 6. R2：停用、历史查询、重新启用与重启

1. 在同一 C01 上停用 Pack，验证 14 disabled、15 个 mount inactive、角色/模型/PDF 新请求在写入前阻断；03/06/07/08/09 业务增量为 0。
2. 停用态从 GUI 项目列表、项目详情和历史时间线查询 C01：标题、事务、Run、Candidate、Receipt、PDF 证据和 `pack_version_ref=v1.0.0` 仍可只读回放，不自动迁移、不丢失、不串项目。
3. 执行停用态主动攻击：角色候选、PDF 重放、旧 request、重复点击和刷新恢复；必须零 Provider、零模型、零 Candidate、零业务 Receipt、Run 不前移。
4. 重新启用同一版本并重启 OS，历史 C01 仍固定 v1.0.0；仅一个全新请求可恢复一次候选能力，停用期请求不得延迟补写。

## 7. R3：卸载、卸载后历史与重装幂等

1. 执行 `uninstall.py`，对账当前 registry/enable 指针、Role Slot/Binding 和 KnowledgeMount 全部不可运行；不得物理删除历史项目、FormalKnowledge 或 Receipt。
2. 卸载态继续从 GUI 查询同一 C01 历史，所有旧 refs、version、证据和时间线保持只读可见；任何新角色/PDF/流程推进请求均零副作用受控阻断。
3. 重复卸载一次验证幂等；随后直接重装同一 digest，registry 只保留一个当前版本，不重复挂载知识、不复制历史项目。C01 仍引用原 v1.0.0，新请求恢复但旧请求不补写。
4. 补双 Tab、断网/恢复、重复点击竞态和 390×844 / 1024×768 / 1440×900；急停 lane 验证顶部与 01 同源、所有业务增量为 0，解除后只恢复新请求。

## 8. 验收、证据与收尾

关键动作同时保存 GUI、03、05/06、07、08、09、13、14 的前后计数、refs、截图与时间戳；只读 API/SQLite 只作对账，不能替代 GUI。交付轮报告、lifecycle 时间线、按钮覆盖表、事实逐句审计、问题台账、截图/Receipt 索引和清理证据。独立验收主体复核：开放 P0=0、阻塞 P1=0、ENV-V13-P1-001 真实 GUI 关闭、直接安装至卸载/重装全程通过、停用/卸载历史可查且零副作用，方可标记`已验收（打包前）`；四个市场阶段始终是 `not_tested_owner_directed`，不得写成 PASS。
