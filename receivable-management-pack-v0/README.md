# 应收管理 Pack（v0）

应收账款**账龄巡检与催办候选**场景荚（Domain Work Pack / Scene Pack）。G4/G5 W-C 窗口交付的**声明骨架**：Pack 只声明流程、角色、能力引用与知识，不发送、不执行、不持主权。

> 成熟度诚实标注：`lifecycle_status = 设计中`。本目录的声明资产已齐备，但**尚未**做隔离 devserver 的装入 / 卸载 E2E，os 侧账龄读模型与催办策略模块亦在同批次另行实现。未跑 devserver 铁证前，不得声称本荚"已装入""已验收"或"催办已接通"。

## 生命周期主线

```
账龄巡检(task / task_governance，只读账龄快照 → 逾期跟进任务候选)
  → 应收顾问出催办建议(advice / 角色包 receivable_advisor，经 08 模型网关)
  → 生成催办草稿(communication / communication_gateway，CommunicationDraftCandidate)
       ├─ 命中时段 / 频次 / 非白名单模板 / 禁用词 → blocked 候选（附 reason 与命中词，不产可发送草稿）
       └─ 通过策略裁剪 → 草稿候选
  → Owner 裁定发送(human_approval，五要素确认 + Base Gate)
       ├─ approved → 归档（发送与回执由 Owner 确认后的受控链路产出）
       └─ rejected → 退回归档(close_rejected，零下游真实动作)
```

## 主权与红线

- AI 全程 Proposer：账龄跟进任务、催办建议、催办话术都只是候选；正式裁定权在 Owner + Base。
- 高风险动作一律回 Owner + Base Gate：`dunning_message_send_confirm`、`payment_plan_change_confirm`、`invoice_write_candidate_confirm`、`formal_memory_confirm`；`delegation_allowed = false`。
- 催办草稿被策略裁剪为 `blocked` 时**仍是候选**，不得包装成"已发送/已催办"。
- 应收事实来源 = Frappe / ERPNext Sales Invoice **只读快照**（`docstatus=1`），Pack 不写回；快照缺失或 Provider 未接通时诚实 `blocked`，不得用 fixture / 模板 / 估算冒充账龄事实。
- 沟通通道未接通时 `not_ready`；本荚不实现任何通道、不直连 IM / ERP。

## 硬约束的归属（重要）

催办的**时段 09:00–18:00 允许、22:00–08:00 硬禁、同一笔应收 7 日滚动窗 ≤ 2 次、模板白名单、禁用词一票否决、默认态无高级模式**，是 **os 侧沟通网关催办策略的硬编码产品默认值**。本荚 `knowledge/dunning/` 只做**知识陈述**，不是执行实现，也不是这些约束的第二真相源；若与 os 侧实现不一致，以 os 侧为准并回报修正知识文件。

## 角色包 / 能力引用

- 角色槽 `receivable_advisor`（单角色）：经 13 SlotBinding 绑定 `role_pack://receivable-advisor`，advice 节点经 08 模型网关真出候选；未绑定 / 模型未接通时诚实 `provider_missing`，不用模板字符串伪装模型输出。
- 能力**只引用 04 能力目录既有条目**，不新造能力 id：
  - `accounting.get_receivable_summary`（只读应收汇总）
  - `accounting.list_invoices`（只读单据列表）
  - `truzhen.communication.draft`（催办草稿，不直接发送）

## 知识域

| 知识域 | 内容 | kind |
|---|---|---|
| `knowledge_scope://receivable/dunning` | `neutral-dunning-templates.md`（中性事实性话术模板集）、`forbidden-phrases.md`（威胁性 / 暗示法律后果 / 涉第三人三类禁用词表）、`dunning-sop.md`（时段 / 频次 / 默认态 / 账龄口径 SOP） | `sop`、`checklist` |

知识全部为 Pack 作者**合成撰写**，不含任何真实客户、金额、人名、单位名或单据号；每条以 `verification_status=pending_human_review` 入库，正式对外适用前须经 Owner / 法务逐条核验。

## 覆盖范围与缺口（诚实清单）

已表达：事务流程（flow）、角色引用（role slot + role pack）、能力引用（04 既有条目）、知识域、判人 / 判事 / 门控 / Provider / 软件需求 / 通知路由 / 多角色对照 / 护城河声明、UI Surface 意图声明（只声明，不实现前端组件）。

尚未表达（backlog，不得当作已完成）：
- 业务对象 schema 与领域语义 glossary（账龄桶、到期日口径、`due_basis` 回退等术语目前只在 SOP 内陈述，未落 glossary 条目）。
- 工作模式集（多工作区 / 模式）。
- 付款计划、争议单据、坏账核销等分支流程。
- 隔离 devserver 的装入 / 卸载 E2E 证据。

## 加载 / 卸载

```sh
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 python3 receivable-management-pack-v0/install.py
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 python3 receivable-management-pack-v0/uninstall.py
```

`uninstall.py` 是**停用型**卸载：经 Base gated-action prepare→confirm 取真签发的 `decision_ref/run_id/nonce`（禁自铸）后调 `14.pack-studio.lifecycle.disable`；停用只收回当前版本与知识域挂载，历史事务对象与 03 回执保留可反查。

## 结构校验

```sh
python3 knowledge_checksums.py --verify receivable-management-pack-v0
python3 -m py_compile receivable-management-pack-v0/install.py receivable-management-pack-v0/uninstall.py
GOWORK=off go test ./...
```
