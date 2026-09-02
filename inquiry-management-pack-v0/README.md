# 询盘管理 Pack（v0）

零售客户获取链段 1-4 的独立文件夹场景荚（Domain Work Pack）：**线索 → 线索池 → 商机分诊 → 报价单**。当前生命周期 `设计中`——本目录是声明骨架，没有任何真实 Provider、模型或前端接线证据，不得据此宣称已实现或已接通。

## 定位与真相源

**真相源 = 外部 ERPNext（`provider_requirements[].provider_family = frappe`）**：`Lead` / `Quotation` / `Customer` 三类对象的事实归外部 ERP。基座 05 只持受控快照，本 Pack 只持声明与候选编排：

- 不自建客户主数据，不做通用 CRM，不在 Pack 层自建对象真相；
- 不照抄其它场景荚的字段当 05 对象 schema；
- 快照按白名单投影，联系方式 / 证件 / 账户类 PII 不进 Pack 侧存储与看板。

manifest 内部一律用厂商中立表述（`external_truth_source.authority = external_erp`）；厂商词只出现在 `provider_requirements[].provider_family` 与 `software_requirements[].provider_family / software_family`，不进 `capability_operation` 与 flow 节点标题。

## Owner 裁定引用（2026-09-02）

| # | 裁定 | 本 Pack 体现位置 |
|---|---|---|
| 真相源 | 外部 ERPNext 持 Lead / Quotation / Customer；Pack 只持受控快照 | `manifest.external_truth_source`、所有读节点 `read_only: true` |
| 零售客户口径 | 渠道（设计师 / 工长 / 异业）与业主都是 ERPNext `Customer`；渠道 = 线索与商机来源，业主 = 服务主体，业主记录链接来源渠道 | `manifest.inquiry_object_declaration.retail_party_semantics`、flow `channel_customer_read` |
| 询盘主锚 | 询盘主锚 = ERPNext `Lead`；`inquiry_source ∈ {channel, direct}` | `inquiry_object_declaration.primary_anchor` / `.inquiry_source_enum`、flow 边 `inquiry_source_channel` / `inquiry_source_direct` |
| 报价抬头策略 | `quotation_party_policy` 默认 `owner_party`，可配置 `channel`——**待 Owner 第 15 条裁定，未裁定前 channel 不得默认启用** | `inquiry_object_declaration.quotation_party_policy`、flow `quotation_candidate.quotation_party_policy` |
| Pack 边界 | 零售 CRM 七段链的段 1-4；「生成合同」是向项目关注 Pack 的交接点；回款 / 开票归应收 Pack | `manifest.chain_scope`、flow `handoff_to_project` |
| 商机对象 | 第一版不依赖独立 Opportunity doctype；分诊结果由 Lead status + 询盘 `payload.triage` 承载 | `inquiry_object_declaration.opportunity_doctype_dependency = false`、flow `triage_candidate.triage_carrier` |
| 能力命名（A-6） | `capabilities.json` 用声明层别名；真实读写在 flow 节点显式写 04 契约字段 `capability_domain` + `capability_operation` | `capabilities/capabilities.json` 的 `capability_bindings`、flow 六个能力节点 |
| 候选类型 | 只用 AGENTS §4.2 白名单八类，全部 `candidate_only` / `non_formal` | `manifest.security_profile`、flow 各 `candidate_type` |
| 不属本 Pack | 月 4 条企微群发上限与催办话术归私域运营 / 沟通 Pack | `chain_scope.handoff_out`、flow `gateway_send.send_note` |

## 能力落法（声明层别名 ↔ 04 契约字段）

| 声明层 `capability_ref` | flow 节点 | `capability_domain` | `capability_operation` |
|---|---|---|---|
| `truzhen.external_snapshot.read` | `inquiry_pool_scan` | `customer-relationship` | `customer.list_leads` |
| `truzhen.external_snapshot.read` | `lead_snapshot_read` | `customer-relationship` | `customer.get_lead` |
| `truzhen.external_snapshot.read` | `channel_customer_read`（可选） | `customer-relationship` | `customer.get_customer` |
| `truzhen.external_snapshot.read` | `quotation_history_read`（可选） | `selling` | `selling.list_quotations` |
| `truzhen.external_snapshot.read` | `quotation_snapshot_read` | `selling` | `selling.get_quotation` |
| `truzhen.capability.quote` | `quotation_candidate` / `gateway_execution` | `selling` | `selling.create_quotation_candidate` |
| `truzhen.capability.inquiry_triage` | `triage_candidate` | —（分诊在基座侧，无外部 operation） | — |
| `truzhen.communication.draft` | `followup_draft` / `gateway_send` | —（草稿经模型网关、发送经沟通网关） | — |

`truzhen.capability.quote` 与 `truzhen.communication.draft` 沿用与家政荚同源的别名；`truzhen.capability.inquiry_triage` 与 `truzhen.external_snapshot.read` 为本 Pack 新增别名。

## 候选与主权链

| 环节 | 候选类型（八类白名单） | 门 |
|---|---|---|
| 询盘对象化 | `BusinessObjectCandidate`（`object_type = inquiry`） | `inquiry_object_confirm` |
| 价值分诊 | `TaskCandidate`（GM 汇入 `department=inquiry` / `topic=inquiry`） | `triage_confirm` |
| 跟进候选 | `CommunicationDraftCandidate`（节点另标 `model_gateway_candidate_type = DocumentDraft`，草稿经 08 模型网关真实生成） | `followup_send_confirm` |
| 报价候选 | `ExecutionIntentCandidate`（`topic=quotation_followup`） | `quotation_confirm` |

四道门在 flow 中都是显式节点：`候选 → owner_gate_* → gateway_* → receipt`。`DocumentDraft` 是 08 模型网关侧的草稿类型名，与 §4.2 的候选八类不是同一命名空间，因此以独立字段承载，`candidate_type` 仍守白名单。

## 完整垂直职业工作台六部分：覆盖 / 缺口（诚实表）

| # | 部分 | 状态 | 说明 |
|---|---|---|---|
| 1 | 工作模式集 | **缺口（backlog）** | 未声明询盘经理的工作区 / 模式划分。 |
| 2 | 事务流程 | **已声明** | `flows/inquiry-management-flow.flow.json`：20 条边、四道显式门、两个回执节点、一个交接终点。只是声明，未经 06 真实装载或运行。 |
| 3 | 业务对象与领域语义 | **部分** | 询盘声明字段、`inquiry_source` 枚举、渠道 / 业主口径与快照 PII 边界已在 manifest 声明；**缺**：术语表、别名、判断口径与来源核验状态（无 `knowledge/`），故不得宣称「语义完整」。 |
| 4 | 能力引用 | **已声明** | 4 个能力别名 + 4 条 ProviderRequirement + 6 个 04 契约操作，manifest / capabilities / flow 三处名字闭合（有测试断言）。**未接通**任何 Provider。 |
| 5 | 角色引用 | **已声明** | 单槽 `inquiry_manager` + `role_pack://inquiry-manager`，全员 Proposer，`delegation_allowed = false`。 |
| 6 | 工作台 UI 声明 | **部分** | 只声明七类主权视觉单元 `ui_surface_slots`；不实现任何前端组件，也未声明 tab / 工作区。 |

其它诚实标注：

- **无知识库**：本 Pack 没有 `knowledge/`，`install.py` 不含知识入库步骤，不得声称有 FormalKnowledge 装入。
- **Provider 未接通**：读能力 fallback `provider_missing`，写能力 fallback `not_ready`；不存在任何 mock / fixture 成功路径。
- **无运行证据**：未做隔离 devserver 装入 / 卸载 E2E，未做真实模型草稿，未做前端用户视角检查。

## 待 Owner 裁定

1. **第 15 条：`quotation_party_policy`**——报价单抬头默认业主（`owner_party`）还是允许按渠道（`channel`）出具。当前实现：默认 `owner_party`，`channel` 只登记为可配置值，未裁定前不得默认启用。
2. `template_family` 取值 `线索到报价客户获取型`——本 Pack 只覆盖链段 1-4，与 smart-home 的「长周期项目交付型」不同族；若 Owner 要求与项目关注 Pack 归同族，改此一字段即可。

## 不做清单

- 不自建客户 / 线索主数据，不做通用 CRM，不在 Pack 层建立第二份对象真相。
- 不依赖独立 Opportunity doctype（第一版）。
- 不做生成合同、项目交付、回款、开票（分属项目关注 Pack 与应收 Pack）。
- 不做群发额度、催办话术、自动回复、联系人抓取。
- 不直连 ERP / 模型 / IM，不实现 Provider，不绕 Gateway，不自铸 `decision_ref` / `run_id` / `nonce` / `receipt_ref`。
- 不在本仓写入任何真实客户名、手机号或金额。

## 加载 / 卸载

```sh
# 装入：脚本只读 os-14 lifecycle ReadModel，Owner 在可信前台完成启用。
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 \
TRUZHEN_CLIENT_URL=http://127.0.0.1:5197 \
  python3 inquiry-management-pack-v0/install.py --open-gui

# 正式卸载：只消费可信前台为同一 Pack/action/transaction 取得的服务端签发证明。
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 \
TRUZHEN_PACK_UNINSTALL_PROOF_JSON='<可信前台签发的证明 JSON>' \
  python3 inquiry-management-pack-v0/uninstall.py
```

两个脚本都不伪造 Owner presence、Cookie 或 Base 决议。`uninstall.py` 走 os-14 正式卸载正门 `/v3/pack-studio/lifecycle/uninstall`；历史询盘事务与 03 Receipt 不因卸载被删除。以上命令**尚未在任何基座实例上跑过**，不构成 E2E 证据。

## 离线契约验证

```sh
python3 -m unittest discover -s inquiry-management-pack-v0/tests -v
```

该测试只验证声明闭合：manifest ↔ capabilities 的能力名与 ProviderRequirement 闭合、flow 节点 04 契约字段与厂商词边界、候选八类白名单、四道门节点化、Owner 裁定口径字段与 PII 静态扫描。不连接 OS、模型或任何 Provider。
