# truzhen-packs 包清单

每个包是一个自包含文件夹（场景荚 = Domain Work Pack），声明领域治理策略需求，经基座真实 lifecycle 装入。

| 包 | pack 标识 | template_family | 成熟度 | 职责（一句话） |
|---|---|---|---|---|
| `environmental-enforcement-pack-v0/` | `scene_pack://environmental-enforcement-flow` | 合规审查执法证据链型 | 完整（install.py + 15 知识域 + 角色包） | 生态环境执法全领域证据链：线索→立案→取证→证据三性→执法精英/挑剔律师多角色对照→处置→文书→Owner/Base 裁定→Receipt。知识库含生态环境法典（带 as_of 择法）。 |
| `smart-home-owner-pack-v0/` | `scene_pack://smart-home-owner-project-ops` | 长周期项目交付型 | 完整（install.py + 项目经理角色包 + Frappe Provider） | 智能家居老板项目经营全周期：机会→Frappe 项目/客户快照→项目经理建议→里程碑/采购/施工任务候选→客户沟通→Frappe 写回→Owner/Base 经营确认门→项目回执。 |
| `housekeeping-ops-pack-v0/` | `scene_pack://housekeeping-ops`（保留兼容 `pack_id: pack_housekeeping_ops_v0`） | 客户服务全生命周期型 | 可装入文件夹包（manifest + flow + capabilities + 2 角色包 + role-slots + install.py；knowledge / uninstall.py 待补） | 家政公司客户服务全生命周期：受理咨询→顾问出方案→质检质询→排期报价→派工任务主权确认→上门执行→归档。 |
| `templates/` | — | — | 脚手架 | 文件夹包模板（如 `scene-pack-software-template`），不参与分发。 |

> **待统一项（记入交接）**：`housekeeping-ops-pack-v0` 已完成 W4 结构升级并可经 `install.py` 装入；仍保留旧 `pack_id` 兼容字段，且尚无 `knowledge/` 与 `uninstall.py`。后续补齐时不得改变其既有 lifecycle / role slot / provider requirement 行为。

## 文件夹包标准结构

| 路径 | 内容 |
|---|---|
| `manifest.json` | 场景荚规格：`pack_ref` / `version` / `template_family` + 六件事（`person_strategy` / `formalization_requirement` / 门控 `gates` / `provider_requirements` / `notification_command_report_routes` / `multi_role_comparison`）+ `moat_justification` + `knowledge_scopes`。 |
| `flows/*.flow.json` | GateFlowSpec 门控流程图（propose / interrogate / compare / owner_gate / base_gate / controlled_execute / receipt / review）。 |
| `role-slots/role-slots.json` | Role Slot 声明（如执法精英 advice + 挑剔律师 challenge）。 |
| `role-packs/*.json` | 绑定的角色包（人格 / 口吻 / 决策习惯 / 模型策略）。 |
| `capabilities/capabilities.json` | 能力需求声明。 |
| `knowledge/` | `knowledge-scopes.json`（知识域）+ `knowledge-index.json`（条目索引）+ 各 `.md`（权威资料结构化，verification_status=pending_human_review）。 |
| `_source-materials/` | Owner 投放的原始资料区（不进 Git，`.gitignore` 屏蔽 `*`，只留 `.gitignore` + `README`）。 |
| `install.py` / `uninstall.py` | 经基座真实 lifecycle 端点装入 / 卸载（canvas→draft→readiness→promote→confirm→role-pack→knowledge）。当前 env / smart-home 具备 install/uninstall；housekeeping 具备 install，uninstall 待补。 |

## 五落点协作边界

| 边界 | 本仓规则 |
|---|---|
| 面向 contracts | Pack manifest / flow / role / capability 引用必须面向 `truzhen-contracts` schema；schema 变更先改 contracts。 |
| 面向 truzhenos | 本仓只放数据和安装脚本；基座负责 Base Gate、Receipt、Gateway、runtime、loader。 |
| 面向 truzhen-software | ProviderRequirement 只声明需要什么 provider；Baserow / Frappe / OCR / IM / sidecar 的真实安装、端口、runtime profile 归 `truzhen-software`。 |
| 面向 client repo | Pack UI 只声明 Surface / visual unit 意图；具体 Web/Desktop/移动渲染归 client repo。 |
