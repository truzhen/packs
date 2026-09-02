# 项目关注 Pack（v0）

面向长周期项目交付的「项目关注」场景荚（Domain Work Pack）。主链：**按需刷新项目外部快照 → 读 05 项目异常清单只读投影 → 项目关注候选 / 异常事项候选 → 异常质询对照 → 通报草稿与处置意图 → Owner + Base Gate → 受控 Gateway → Receipt**。

两个铁口径先说清楚：

- **阈值判定不在本 Pack。** 延期 / 停滞 / 付款节点逾期三条规则的计算、分级、缺字段处理与快照新鲜度全部在 truzhenos 05 的只读投影（`GET /v3/business-object/project-anomalies`）里。Pack 只以 `policy_ref: project_watch_policy://default` 引用策略对象，flow 内不写任何条件表达式、阈值数值或分支判定（packs `AGENTS.md` §4.2 硬禁）。
- **真相源是 ERPNext（外部系统）。** 05 只持快照与投影：读模型 ≠ 真相源，候选 ≠ 正式。零快照时读模型诚实返回 `not_ready`，快照过期带 `stale` 标记，缺字段的规则输出 `skipped`——都不得被表述为「项目一切正常」。

生命周期：`manifest.lifecycle_status = "设计中"`。本 Pack 目前只有声明与离线契约测试，**没有**跑过隔离 devserver 的端到端装入 / 启用 / SceneFlowRun 验收，不得据此宣称已接线或已验收。

## 一、标识

| 项 | 值 |
| --- | --- |
| `pack_id` | `project-watch-pack-v0` |
| `pack_ref` | `scene_pack://project-watch` |
| `version` | `0.1.0` |
| 版本化引用 | `scene_pack://project-watch@0.1.0`（= `pack_ref` + `@` + `version`，设计稿 D5 的写法；仓内 lifecycle ReadModel 按 `pack_ref` 查询、按 `version` 比对，故 manifest 内两者分列） |
| `template_family` | 长周期项目交付型 |

## 二、垂直职业工作台六部分（AGENTS.md §3.2）：已实现 / backlog

| # | 部分 | 本 Pack 状态 | 说明 |
| --- | --- | --- | --- |
| 1 | 工作模式集 | **部分实现** | 只有「项目关注」一个工作模式（一条 flow）；无多工作区切换。 |
| 2 | 事务流程 | **已实现（声明层）** | `flows/project-watch.flow.json` 14 节点 15 边，候选 → 质询 → Owner/Base 门 → Gateway → Receipt 全链声明齐备；运行解释归基座。 |
| 3 | 业务对象与领域语义 | **部分实现** | 异常事项候选以 `BusinessObjectCandidate` 基类声明，幂等键取读模型的 `anomaly_key`，字段口径 100% 回推自 05 读模型响应，**不在 Pack 层自建对象真相、不照抄 smart-home 的 `project_*` 字段**。规则语义与默认阈值落在 `knowledge/`，`pending_human_review`。术语表（glossary）为 **backlog**——缺真实客户样本，不宣称语义完整。 |
| 4 | 能力引用 | **已实现（声明层）** | 5 条 ProviderRequirement，`manifest.provider_requirements` 与 `capabilities/capabilities.json` 逐条同源，全部标注 L1–L6 执行级别与 `runtime_requirement`；均未接通，诚实 `provider_missing / not_ready / blocked`。 |
| 5 | 角色引用 | **已实现（声明层）** | 两槽两包：`project_watcher`(advice) + `exception_challenger`(challenge)，全员 Proposer，`proposer_only` 写进角色包决策风格。 |
| 6 | 工作台 UI 声明 | **只声明** | `ui_surface_slots` 为七类主权视觉单元；**不实现任何前端组件**，UI 归 client 仓。 |

其它明确的 backlog：

- `material_watch` 节点：v1 只承接读模型输出的 Task 类异常，**不做任何物料判定**（`backlog: true` + `backlog_reason` 写在节点上，设计稿 D4）。真正的物料口径待真实客户证据后另立。
- 跟进任务写回 ERPNext：本卡只出 `ExecutionIntentCandidate`，写回通路不接，`fallback_policy: blocked`。
- GM 队列汇入：只声明契约（见第五节），不实现、不 stub。
- 阈值策略的写路径（setting 端点 + 03 回执）归 05 侧另一批，本 Pack 只读默认策略。

## 三、领域候选名 → 8 个候选基类映射

Pack 只能声明 `AGENTS.md` §4.2 的 8 个候选基类。下表左列是 flow 内的领域特化命名（`domain_candidate_type`），**只是命名，不是新候选种类**；右列是实际生效的 `candidate_type`。

| flow 节点 | 领域候选名 | 8 基类（`candidate_type`） |
| --- | --- | --- |
| `snapshot_refresh` | —— | `CapabilityInvocationCandidate` |
| `anomaly_scan` | —— | `CapabilityInvocationCandidate` |
| `progress_watch` | `ProjectWatchTaskCandidate`（项目关注候选） | `TaskCandidate` |
| `material_watch` | `ProjectMaterialWatchCandidate`（backlog） | `BusinessObjectCandidate` |
| `payment_watch` | `ProjectPaymentWatchCandidate` | `TaskCandidate` |
| `anomaly_item_candidate` | `ProjectAnomalyCandidate`（异常事项候选） | `BusinessObjectCandidate` |
| `challenger_review` | `AcceptanceEvidenceCandidate`（质询验收证据） | `BusinessObjectCandidate` |
| `owner_notice_draft` | —— | `CommunicationDraftCandidate` |
| `followup_intent` | —— | `ExecutionIntentCandidate` |
| `gateway_execution` | —— | `CapabilityInvocationCandidate` |

全部候选 `candidate_only: true`；引用一律由 os 服务端签发，Pack 不自铸 `decision` / `run_id` / `nonce` / 回执引用。

## 四、能力与 Provider（全部未接通）

| requirement_id | capability | gateway | risk | fallback | provider_family | 执行级别 | runtime |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `frappe_project_snapshot` | `project_external_snapshot` | execution | low | `provider_missing` | frappe | L1 | local_preferred |
| `req_project_anomaly_readmodel` | `project_anomaly_readmodel` | memory | low | `not_ready` | readmodel | L1 | local_required |
| `frappe_project_followup_task_write_candidate` | `project_followup_task_write_candidate` | execution | medium | `blocked` | frappe | L1 | local_preferred |
| `project_watch_notice_draft` | `truzhen.communication.draft` | communication | high | `not_ready` | communication-gateway | L6 | local_required |
| `project_watch_followup_execution_intent` | `truzhen.execution.intent` | execution | high | `not_ready` | execution-gateway | L6 | local_required |

L6 = 本 Pack 不声明任何自动发送 / 执行手，送达与处置由 Owner 手动兜底。`software_requirements` 显式为空数组（附 `software_requirements_note` 说明理由），本 Pack 不引入任何新的底层软件本体。

## 五、GM 汇入约定（只声明）

`gm_handoff` 节点按窗口 D 定稿 v1 声明候选汇入 GM 队列的口径，**本 Pack 不实现汇入**：

- 队列定位：`department: "project"`、`topic: "project_watch"`。
- 端口 / 端点：`businessruntime.CandidateIntake.Submit` / `POST /v3/business-runtime/candidates/intake`。
- 去重键：`project:<kind>:<project_external_id>:<as_of>`，与读模型的 `anomaly_key` 同源。
- `priority_hint: high`；`object_refs` 取 Project 外部链接引用，`evidence_refs` 取读模型 evidence。
- GM 端口未接通时诚实 `not_ready`，不 stub 冒充 queued。任何带正式回执 / decision / run_id / nonce 的字段会被 GM 侧拒收。

## 六、知识域

一个知识域 `knowledge_scope://project-watch/anomaly-rules`（kinds: `sop` + `checklist`），当前 1 条：三条规则语义、两级分级说明、**默认阈值是可改的默认值**、诚实边界。整体 `verification_status: pending_human_review`——缺真实客户样本核验，不得据此对外承诺。`knowledge/knowledge-index.json` 的 checksum 由仓根 `knowledge_checksums.py` 生成与校验。

## 七、装入 / 卸载

```sh
# 装入：脚本全程只读，写入只能来自可信前台。
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 \
TRUZHEN_CLIENT_URL=http://127.0.0.1:5197 \
  python3 packs/project-watch-pack-v0/install.py --open-gui

# 正式卸载：只消费可信前台 / Base 已签发并外部注入的卸载证明。
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 \
TRUZHEN_PACK_UNINSTALL_PROOF_JSON='{...}' \
  python3 packs/project-watch-pack-v0/uninstall.py
```

`install.py` 不伪造 Owner presence，也不发任何写请求：先做本地声明自检（治理六件事 + 三路径 + `knowledge_scopes` 一致 + manifest↔capabilities 逐条对齐）与知识 checksum 防漂移，再交接可信前台，然后只读等待 os-14 证明**精确版本**启用，并只读复核 os-13 角色槽绑定与 os-09 知识入库（含可反查回执）。任一阶段未观察到即按阶段化错误码 fail closed，不宣称装入成功。

`uninstall.py` 是**正式卸载**语义（`14.pack-studio.lifecycle.uninstall`），级联停用知识域；历史候选、业务对象与 03 回执保留可反查——卸载不删历史。

## 八、离线契约验证

```sh
python3 -m unittest discover -s project-watch-pack-v0/tests -v
```

该测试只验证声明一致性：候选类型落在 8 基类、领域名映射、manifest↔capabilities 逐条对齐、role_slots 镜像与 `node_type` 一致、门控词汇三处同名、`knowledge_scopes` 三处一致、flow 内无阈值字段、无禁品字面串。不连接 OS，不调用任何 Provider，不读取任何真实业务数据；本 Pack 内所有示例值均为合成脱敏值。
