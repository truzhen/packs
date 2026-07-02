# truzhen-packs 模块与包清单

本文件记录 `/Users/li/Documents/truzhen-packs` 当前包模块、职责边界、成熟度和协作关系。完整纪律见 `AGENTS.md`。

## 1. 仓库总览

`truzhen-packs` 是 Truzhen 包层资产仓。每个正式 pack 是一个自包含文件夹，通过 `manifest.json`、flow、role-slots、role-packs、capabilities、knowledge 和 install / uninstall 脚本声明领域治理资产，并经基座真实 lifecycle 装入。

核心链路：

```text
Pack 声明 -> 基座 lifecycle enabled -> Scene Runtime 解释 -> Candidate -> Owner + Base Gate -> Gateway -> Receipt
```

Pack 不直接生成 Formal Record，不直接执行 provider，不直接写基座数据库。

## 2. 当前包清单

| 包 | pack 标识 | template_family | 成熟度 | 职责 |
|---|---|---|---|---|
| `environmental-enforcement-pack-v0/` | `scene_pack://environmental-enforcement-flow` | 合规审查执法证据链型 | 完整文件夹包：install / uninstall、flow、2 角色包、capabilities、knowledge | 生态环境执法证据链：线索、立案、取证、证据三性、执法精英 / 挑剔律师对照、处置、文书、Owner/Base 裁定、Receipt。 |
| `smart-home-owner-pack-v0/` | `scene_pack://smart-home-owner-project-ops` | 长周期项目交付型 | 完整文件夹包：install / uninstall、flow、1 角色包、Frappe ProviderRequirement；无知识库 | 智能家居老板项目经营：机会、Frappe 项目 / 客户快照、项目经理建议、里程碑 / 采购 / 施工候选、客户沟通、Frappe 写回候选、项目回执。 |
| `housekeeping-ops-pack-v0/` | `scene_pack://housekeeping-ops`，兼容 `pack_housekeeping_ops_v0` | 客户服务全生命周期型 | 可装入文件夹包：manifest、flow、capabilities、2 角色包、role-slots、install；`knowledge/` 与 `uninstall.py` 待补 | 家政客户服务全生命周期：受理咨询、顾问出方案、质检质询、对照确认门、排期报价、派工确认、通知客户、上门服务、服务回执、归档。 |
| `templates/` | — | — | 作者端脚手架 | Pack 作者工程模板，不参与 enabled pack 分发。 |

## 3. 标准文件夹包结构

| 路径 | 职责 | 真相源归属 |
|---|---|---|
| `README.md` | 包说明、主权红线、加载 / 卸载和验收入口 | 本仓文档 |
| `manifest.json` | pack 标识、版本、template_family、治理六件事、ProviderRequirement、护城河、知识域 | Pack 声明；schema 归 contracts / 基座 |
| `flows/*.flow.json` | GateFlowSpec / Scene Flow 图纸，声明节点、边、等待点和门控 | Pack 声明；运行解释归基座 06 |
| `role-slots/role-slots.json` | 场景角色槽、责任、默认角色包、绑定期望 | Pack 声明；绑定执行归基座 13 |
| `role-packs/*.json` | 角色人格、口吻、决策习惯、模型策略和权限边界 | Pack 声明；启用归基座 13 |
| `capabilities/capabilities.json` | 能力需求、ProviderRequirement、gateway class、risk class、fallback policy | Pack 声明；provider 真相归基座 / 外部 provider |
| `knowledge/knowledge-scopes.json` | 知识域、挂载策略、知识 kind | Pack 知识索引；正式知识真相归基座 09 |
| `knowledge/knowledge-index.json` | 知识条目、source_ref、scope、kind、verification status | Pack 知识索引；正式知识真相归基座 09 |
| `knowledge/**/*.md` | 结构化知识内容 | 本仓可分发资料；正式适用需人工核验 |
| `_source-materials/` | Owner 本地投放原始资料入口 | 原始资料不进 Git |
| `install.py` | 调基座 lifecycle 装入 pack、角色、槽位、知识 | 胶水脚本；正式状态归基座 |
| `uninstall.py` | 调基座 Base Gate + disable 卸载 pack | 胶水脚本；历史 Receipt 不删除 |
| `docs/` | 测试报告、验收证据、设计说明 | 本仓文档 |

## 4. Pack 治理六件事

正式 Domain Work Pack 必须声明：

1. `person_strategy`：谁只是 Proposer，高风险动作如何回 Owner + Base。
2. `formalization_requirement`：哪些候选可进入正式态，裁定权归谁。
3. `gates`：Owner Gate、Base Gate、角色对照门、执行门、发送门等。
4. `provider_requirements`：需要哪些能力和 provider readiness，不实现 provider。
5. `notification_command_report_routes`：通知、命令候选、回报和回执路由。
6. `multi_role_comparison`：多角色建议 / 质询 / 对照确认门，禁止隐藏 agent 回路。

此外必须声明 `moat_justification`，回答“用户能否一句话让最强模型直接做出同样结果”。如果能，则应降级为 prompt / 模板 / 能力说明，不应冒充场景荚。

## 5. 跨仓协作边界

| 目标 | 本仓规则 |
|---|---|
| `truzhen-contracts` | Pack manifest / flow / role / capability 引用必须面向契约；schema 变更先改 contracts，并列影响清单。 |
| `truzhenos`（旧 `truzhenv3` 已冻结） | 本仓只放声明数据和装载脚本；Base Gate、Receipt、Gateway、runtime、loader、ReadModel、Pack Studio 归基座。 |
| `truzhen-software` / provider 仓 | Provider、sidecar、外部软件安装、runtime profile、端口、健康检查归外部 provider / software 线。 |
| client repo / frontend | Pack UI 只声明 Surface / visual unit 意图；具体 Web / Desktop / Mobile 渲染归 client repo 或基座前端。 |
| `truzhen-cloud` | 本仓可准备可分发资产和市场展示候选元数据；正式 `PackListing`、审核发布、价格、订单、支付、License / Entitlement、Release、下载分发状态、作者后台和云端分发服务归 `truzhen-cloud`；发布、上架、授权动作必须 Owner 明确授权。 |

### 5.1 云边界红线

- 不保存 `payment_webhook`、`license_server`、`entitlement_db`、`cloud_secret`、`admin_secret`。
- 不把订单、支付状态、License、Entitlement 写成 Pack 真相。
- 不承载官方云端 server、官方云端网页或云端部署脚本。
- 需要云市场发布时，只写 Pack 元数据和发布引用，由 `truzhen-cloud` 持有发布状态和回执。

## 6. 当前包细节

### 6.1 `environmental-enforcement-pack-v0/`

**定位**：生态环境执法全领域证据链 Pack。

**核心资产**：

- `manifest.json`
- `flows/environmental-enforcement-flow.flow.json`
- 角色包：`role_pack://enforcement-elite`、`role_pack://critical-lawyer`
- 能力需求：OCR、文书生成、文书送达、执法执行、在线监测读取
- 知识库：`knowledge-scopes.json`、`knowledge-index.json`、法规 / SOP / 案例 / 索引 Markdown
- `install.py` / `uninstall.py`

**纪律**：

- 法律和执法知识默认 `pending_human_review`。
- 文书、处罚、送达、查封扣押、移送公安等只能是候选，正式动作需 Owner + Base Gate。
- Provider 未接通时必须 `blocked / provider_missing / not_ready`。
- `manifest.json` 的知识域声明必须与 `knowledge/knowledge-scopes.json` 和实际目录一致。

### 6.2 `smart-home-owner-pack-v0/`

**定位**：智能家居老板项目经营 Pack。

**核心资产**：

- `manifest.json`
- `flows/smart-home-owner-project-ops-flow.flow.json`
- 角色包：`role_pack://smart-home-project-manager`
- 能力需求：Frappe 项目快照、客户快照、项目写回候选
- `install.py` / `uninstall.py`

**纪律**：

- Frappe 只是 ProviderRequirement，不是真相源。
- 里程碑、采购、施工、对外承诺、Frappe 写回必须经 Owner + Base Gate。
- 本 pack 当前无知识库；不得声称有 FormalKnowledge 装入。

### 6.3 `housekeeping-ops-pack-v0/`

**定位**：家政 / 保洁客户服务全生命周期 Pack。

**核心资产**：

- `manifest.json`
- `flows/customer-service-lifecycle.flow.json`
- 角色包：`role_pack://housekeeping-consultant`、`role_pack://quality-auditor`
- 能力需求：排期、报价草稿、上门执行意图
- `install.py`

**纪律**：

- 当前 `knowledge/` 缺失，不能标成完整知识包。
- 当前 `uninstall.py` 待补，不能标成完整可卸载闭环。
- 报价发送、派工、上门执行和归档正式化必须经 Owner + Base Gate。
- legacy seed 脚本如继续保留，必须标明不替代正式 install lifecycle。

### 6.4 `templates/scene-pack-software-template/`

**定位**：作者端软件工程模板。

**负责**：

- 给 Pack 作者展示可登记、可测试、可说明的软件工程材料结构。
- 提供 objects、flows、views、policies、capabilities、adapters、software、tests、docs 等目录范式。

**不负责**：

- 不作为 enabled pack。
- 不接真实 provider。
- 不保存构建产物、真实密钥或运行数据库。

## 7. 验收映射

| 改动类型 | 必跑验证 |
|---|---|
| Markdown 文档 | `git diff --check` |
| JSON / manifest / flow / role / capability | JSON 合法性 + 结构审计 |
| install / uninstall 脚本 | `python3 -m py_compile` |
| knowledge | JSON 合法性 + scope / index / 文件一致性 + 来源与 `pending_human_review` 检查 |
| forbidden artifacts | `git ls-files` 禁入产物扫描 |
| lifecycle 行为 | 隔离基座 devserver E2E 装入 / 卸载，提供 registry、enabled version、Receipt 证据 |
| ProviderRequirement 语义 | 影响清单 + readiness / blocked 行为说明；必要时跨仓验收 |

## 8. 待统一项

- `environmental-enforcement-pack-v0/` 无固体废物（solid-waste）专属知识域：15 个知识域中无该域目录，固废内容仅散见 penalty / eia-permit / code 等域文件内。2026-07-02 已从 manifest 删除该 scope 声明对齐事实；如 Owner 权威资料含固废专章，后续按知识导入流程真实建域后再恢复声明，不得空声明。
- `housekeeping-ops-pack-v0/` 仍缺 `knowledge/` 与 `uninstall.py`，后续补齐时不得改变既有 lifecycle / role slot / provider requirement 行为。
- `environmental-enforcement-pack-v0/` 的知识域声明需要持续与实际 `knowledge-scopes.json` 同步。
- 模板 manifest 与实际 pack manifest 目前不是同一 schema，后续如要作为新包脚手架，应补一个更贴近当前 Domain Work Pack 文件夹结构的模板。

## 9. 完成口径

Pack 商品化或成熟度升级至少要能证明：

- 可安装、可启用、可停用或明确说明卸载缺口。
- Pack Readiness / 护城河理由明确。
- Role Pack 和 SlotBinding 能闭环或明确 `not_ready`。
- Provider readiness 诚实显示 `ready / degraded / provider_missing / blocked`。
- 正式动作不绕过 Owner + Base Gate。
- 关键动作有 Evidence / Receipt 可反查。
- 缺口写入 README、MODULES 或测试报告，不能静默包装成完成。

