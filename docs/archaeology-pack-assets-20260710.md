# KWeaver/Palantir 考古 → Pack 体系触发式资产登记（2026-07-10）

> 用途：登记两轮考古中与 Pack 体系相关、**等触发条件到了才施工**的资产，防止后续会话重复考古或提前施工。节奏权威 = `truzhen-contracts/docs/governance/development-pacing-20260710.md`（阶段 1 明确不做 backlog 项）；裁定权威 = `/Users/li/Documents/systong/truzhen-notes/kweaver-archaeology/unified-decision-table-20260710.md`。
> 红线：KWeaver（爱数，破产重整、IP 争议）**代码一行不进本仓**；只读设计自己重写；无 LICENSE 仓（kweaver-eval/sandbox/workflow 等 24 个）连思路引用都注出处。存证：`/Users/li/Documents/systong/oss-sources/LICENSES.md`。

## 触发式资产（Pack 侧）

| 资产 | 触发条件 | 一句话内容 | 笔记出处 |
| --- | --- | --- | --- |
| `.skill` 三段式规格 | 阶段 2 装修 Pack 作者轮 | 技能包显式分区：config + 领域方法论 references/*.md + 交付模板 assets/——「高手经验」的可交付载体形状 | K-dip-ecosystem-notes |
| SKILL.md「不做」负面清单 | 同上（并轮） | 每个技能包声明严格负面边界（唯一 Apache-2.0 可放心直引出处：kweaver-engineering） | K-light-sweep |
| install 原子回滚协议 | 智能家居客户交付前评估提级 | 安装失败自动回滚前一版本；本仓 install.py 现无回滚档 | K2-sdk-skill-notes |
| type:id 粒度 CHECKSUM + 确定性导入五能力 | 内容防漂移需求出现时 | validate/diff/dry_run/apply/export round-trip；对照本仓现有 forbidden artifacts 扫描的空档=内容级校验 | K1-semantic-notes |
| RiskType 五件套判事策略骨架 | 作者轮验证字段需求后（橙，涉 contracts） | 管控范围/管控策略/前置检查/回滚方案/审计要求——「risk_level 说多危险、RiskType 说如何管控」 | K1-semantic-notes |
| 制作台作者旅程 9 步表 | 14 制作台施工时直接取用 | 顺/坎/启示逐步登记；填空脚手架+@ 联想+1 必填渐进 | K3-studio-governance-notes |
| 行业本体参考样例 | 智能家居/装修包业务对象设计参考 | 害虫防治 7 对象+47 数据行（Palantir 侧）+ 供应商/物料/BOM 种子 SQL | palantir 03 + K-dip-ecosystem-notes |

## 反面教材（Pack 形态红线，长期有效）

- **行业方案 ≠ 装进平台的行业 App**：dip-for-supply-chain = 82% 行业 UI 代码 + KN ID 硬编码前端 + 种子 SQL 手工导入——场景荚「声明治理策略+知识+流程图，不持主权」路线的反面。本仓 Pack 永远是声明性资产，不携带自己的 UI 应用。
- **治理声明必须有运行时下半场**：KWeaver 全生态「requires_approval 只在规格、风险闸只有脚手架」——本仓场景荚声明的判人/判事/门控策略，正式裁定永远在 Owner+Base Gate（六件事纪律不变），声明与 enforcement 的接线验收归基座轮。

## 现在补（唯一进行中项）

场景荚 manifest 补 `lifecycle_status` 字段（八档生命周期自洽修补）——计划：`truzhen-contracts/docs/plans/kweaver-two-items-plan-20260710.md`，**待 Owner「开工」，本仓侧任务=4+1 个 manifest 增字段（跟随 contracts schema 先行）**。
