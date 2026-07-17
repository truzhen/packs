# 环保执法 Pack 单项目完整生命周期复验计划 v17

日期：2026-07-16  
状态：`当前客户主线；待材料完整性、lifecycle 与 cloud-auth 修复进入最新 main 后，以全新隔离环境复验；未发布`

本计划替代 `/Users/li/Documents/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v16-20260715.md`。

前序证据：

- `/Users/li/Documents/过程文档/env-prepack-v15-20260714/rerun-20260715/v15-rerun-closeout-20260715.md`
- `/Users/li/Documents/过程文档/env-prepack-v15-20260714/audits/environmental-enforcement-v1-v16-git-audit-20260716.md`
- `/Users/li/Documents/过程文档/env-prepack-v15-20260714/rerun-after-p1-fix-20260715/issues/authenticated-author-readonly-coverage-20260716.md`

## 1. 派活卡与 Owner 裁定

| 维度 | v17 裁定 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线；只验收 `environmental-enforcement-pack-v0`，智能家居结果不冻结本计划。 |
| 真实客户 / 场景证据 | 既有真实 GUI 已证明安装、15 scope/45 源文档、C01 PDF、角色候选、停用/急停零副作用等主链；Git 审计发现材料完整性硬门、可逆 lifecycle、材料矩阵 UI、云认证会话复验仍是未提交 WIP。 |
| 我要做的事 | 从四仓最新 `origin/main` 建全新隔离分支，以唯一 C01 跑直接安装、启用、完整 flow、材料补齐、停用、历史查询、重新启用、断网、急停、卸载和重装；证明修复已脱离旧 WIP。 |
| 最小可交付 | 关键材料缺失/冲突/待人工核验时，在任何角色模型和 compare 前生成 `EvidenceGapCandidate` 并稳定阻断；材料经人工核验后才进入双角色候选；全生命周期和历史查询可回放。 |
| 明确跳过 | `提交审核`、`模拟支付`、`Entitlement`、`下载`全部记 `not_tested_owner_directed`；不启动市场购买链，不把跳过项写成 PASS。作者只读登录仅在本地 cloud-auth 已接线时作为独立非主链 lane。 |
| 真相源 | Pack/digest 归本仓；lifecycle 归 14；项目归 05；Run/材料候选矩阵归 06；知识归 09；模型归 08；central Receipt 归 03；云会话归 cloud-01；client 只投影。 |
| 仓库 / 层归属 | OS 负责通用材料规则、Gate/Receipt/runtime 和云认证消费；client 展示矩阵/阻断；本仓声明 flow/证据需求/lifecycle；cloud 维护既有会话轨。contracts/software 不改。 |
| 风险与 AI | L0 知识/manifest；L1 角色推理；L2 OCR/模型/网络/本地软件。AI 仅为测试施工者和 Proposer；法律裁定、正式签发及红色动作属于具备资格的人 + Owner/Base。 |
| 契约影响 | 零 contracts、FormalEvidence、Gate/Receipt schema 或跨仓 DTO 变更；材料需求是 Pack flow 元数据，矩阵是 06 候选态投影。 |
| 禁止边界 | 不读真实案件、客户凭据或生产数据；不正式处罚、送达、删除、对外提交或执法执行；不手改数据库、降 validator、拼模板、接受越界事实或用 mock/旧截图冒充 PASS。 |
| 生命周期 | 修复为`已实现 -> 已接线`；v17 与独立复核通过后升`已验收（打包前）`；云生产部署、发布和上架另行裁定。 |

## 2. 权限、自治、隔离与停止规则

测试执行者获本计划范围内全部本地权限：隔离 worktree/HOME/SQLite/浏览器 profile，启停 OS/client/local cloud，真实本地模型，合成 C01 PDF，直接安装/启停/卸载/重装，截图及只读 API/SQLite/Receipt 对账，并修复不改契约与主权链的中小问题。

小问题和中等问题必须定位根因、补测试并从最近可信检查点继续。P0/P1 或其它大问题登记后，只冻结受影响危险 lane；其它独立安全 lane 继续跑完。仅当继续会扩大法律、数据、真实副作用或主权风险，或所有剩余步骤都依赖被污染的权威状态时暂停。不得为了跑通静默降级或伪造材料/Receipt。

L0/L1/L2 级别写入每条证据。OCR、网络、模型和本地软件只能经 L2 隔离 Provider/Gateway；容器、端口、登录页或 `health=200` 不能替代 readiness、Gate、Receipt 和业务结果。

## 3. 共通硬门：Pack 不得腐化基座

1. manifest、schema、flow、role、capabilities、knowledge scope/index 通过权威校验；Pack 不复制、放宽或私造 contracts/行业 schema。
2. Candidate/Formal 在存储、接口和 UI 三层隔离；Role/Model/Pack 只能产 Candidate。正式事实、文书或动作必须由具备资格的人确认并经 Owner + Base Gate。
3. 权限取 Owner 授权、角色白名单、Pack、ContextSlice、KnowledgeMount 与 Provider 能力最小交集；Pack 不持凭据、不直连网络/文件系统/外部软件。
4. 缺配置、缺 Provider、缺证据、未人工核验、Provider 未 ready、法源/Pack/软件版本冲突时，必须明确 `blocked/not_ready/provider_missing` 与稳定 reason/missing fields；不得静默换 Provider、补假数据、空成功或用 `degraded` 掩盖硬缺口。
5. Gate 裁定和真实动作只能由 Base/Gateway 完成；Provider ready、容器存活、前端提示或 Pack 声明不能代替 Gate/Receipt。
6. 接受、阻断、失败、重放和执行尝试都有可反查 Receipt，关联 candidate、decision、evidence、法源/Pack 版本、幂等键、结果和副作用计数；历史回放不从当前状态重算。
7. 安装、启停、PDF 回灌、角色请求、重试、双击、断网、旧请求、卸载和重装使用稳定幂等键，不重复 Candidate、Formal、Receipt 或外部动作。
8. 停用/卸载后 RoleBinding、KnowledgeMount、Provider session 与临时资源不可运行；历史 C01、FormalKnowledge 与 Receipt 只读保留，不与其它 Pack 串包。
9. 静态扫描 `truzhenos` 不得出现 environmental Pack ref、案件类型或 flow 节点专用分支、专用 Gate/Gateway/Receipt、行业 schema 副本或 seed；client 不持真相；Pack 不含 Provider/运行产物。

## 4. 唯一 C01、直接安装与开跑硬门

唯一项目：`C01-v17-固体废物露天堆放线索核查`。同一 `transaction_ref`、Pack v1.0.0、digest、安装实例、浏览器 profile 和证据目录贯穿全轮；不得新建项目清空缺陷。仅 05/06 权威状态不可逆污染时允许一次 `C01-RECOVERY-01`，必须登记旧/新 refs 和原因。

直接从净制品安装，跳过市场四阶段。GUI 前必须完成：四仓最新主线 baseline；Pack JSON/Python/结构、45 源文档/15 scope、符号链接/路径/禁品/digest；OS 06/14/16 定向/race 和全量 EGR；client 材料矩阵/阻断定向、全量测试/typecheck/build/smoke；cloud-01 会话伪造凭据 fail-closed。任何失败先修根因。

## 5. R1：安装、材料硬门与完整 flow

1. 空 registry/KnowledgeMount 安装并启用 v1.0.0；核对 registry 0→1、2 Role Pack、2 SlotBinding、15/15 mount、45/45 源文档、digest 和 03 Receipt。
2. GUI 创建唯一 C01 和 06 Run。PDF 前材料矩阵逐项显示 missing/pending；点击 advice/challenge 必须返回 `evidence_requirements_not_ready` / `EvidenceGapCandidate`，列真实 missing fields；角色模型、角色 Candidate、Task、外发和执行增量为 0。
3. 上传合成 PDF 后，只把线索原件改为待人工核验；“已上传”不得自动等于 provided。依次补齐主体身份、管辖地点、物料性质/样品、现场影像/记录、监测检测、许可台账的合成证据和人工核验状态；每次只更新对应矩阵项并保留 Receipt/时间线。
4. 任一关键材料 missing/conflicting/pending 时不能进入角色/compare。全部 required 项明确 verified 后，才允许执法精英和挑剔律师各生成一次候选；每个事实逐字定位到 06 evidence/ContextSlice，每个法律引用定位到 09 `knowledge_ref`，未知项标 `pending_human_review`。
5. 法源逐项核验名称、发布机关、版本/修订、生效/失效时间、辖区、条款和引用证据；过期、辖区不符、版本冲突或无法证明现行有效时 `blocked/not_ready`。
6. 敏感信息在 GUI、截图、日志、PDF 与 Receipt 中按最小展示脱敏。文书草稿与正式签发分区、分 ref、分状态；处罚、送达、删除、对外提交均为红色动作，只验证 Candidate/模拟和具备资格者确认门，不作正式动作。
7. 覆盖 flow 全节点和所有可见按钮，逐项标 `PASS/受控阻断/not_ready/provider_missing/not_tested`；保持 0 正式处罚、0 送达、0 删除、0 对外提交、0 执法执行。

## 6. R2：停用、历史、恢复与故障

1. 同一 C01 停用 Pack；攻击角色、PDF 重放、旧请求、重复点击和双 Tab，03/06/07/08/09 业务增量为 0。
2. 停用态从项目列表、详情、材料矩阵、时间线、候选和 Receipt 查询同一 C01；标题、Run、PDF、矩阵、候选、Receipt 与 v1.0.0 只读可回放。
3. 重新启用并重启 OS/client；历史固定原版本，只有新请求恢复，停用期请求不得延迟补写。
4. 覆盖 Provider missing/not_ready、版本冲突、断网前已提交/未提交边界、恢复、幂等重试、390×844 / 1024×768 / 1440×900；错误必须具体、稳定、可诊断。

## 7. R3：急停、卸载/重装与认证只读 lane

1. 启用急停后角色、PDF、旧 request、双 Tab 和伪 receipt 全部零业务副作用；急停状态、01、central 03 Receipt 与 GUI 同源，重放不增量；解除后只恢复新请求。
2. 卸载后运行访问、RoleBinding、KnowledgeMount 和 Provider session 撤销；新请求受控阻断、重复卸载幂等，同一 C01 历史仍可查。同 digest 重装后 registry/mount 不重复、历史不复制、C01 固定原版本。
3. 本地 cloud-auth 已接线时，以真实本地测试账号验证 password-login 后 `X-Truzhen-Session-Id` 的 Kratos cookie/API token 双形态 `/cloud-auth/session` 复验；伪造凭据未登录且不泄漏 raw token。该 lane 只验证作者只读会话，不做提交审核、支付、权益、下载、发布或生产部署。
4. cloud 未部署或不可达时记 `not_ready`，不阻断已经完整跑通的直接 Pack lifecycle 主链，也不把本地验证写成公网已发布。

## 8. 证据与放行

过程目录：`/Users/li/Documents/过程文档/env-prepack-v17-20260716/`。交付基线、材料矩阵逐步变化表、事实/法源逐句审计、完整 flow/按钮表、Provider 状态矩阵、三视口/双 Tab/断网记录、生命周期和历史回放时间线、问题台账、截图/Receipt 索引、完整门禁日志及资源清理证据。

放行要求：开放 P0=0、阻塞 P1=0；所有安全可继续 lane 已跑完；材料硬门在任何角色模型前 fail-closed，核验后双角色候选事实/法源/辖区/证据完全锚定；脱敏、草稿/正式隔离、停用、急停、卸载、历史、幂等、断网、重装和基座防腐通过。满足后仅标 `已验收（打包前）`；市场四阶段仍为 `not_tested_owner_directed`，生产 cloud 未部署不得写成已发布。
