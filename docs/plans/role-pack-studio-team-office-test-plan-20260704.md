# 角色制作台完善与测试计划

> 计划日期：2026-07-04
> 当前档位：`设计中`
> 计划类型：完善计划 + 测试计划。本文不代表已经实现、不代表已改契约、不代表团队设置或音色 / VRM 已接线。
> 参考文件：
> - `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/gui-capability-pack-workbench-github-oss-test-plan-20260704.md`
> - `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/FEATURE_LEDGER.md`

## 0. 明确目标

本计划的目标是：让“角色制作台”能生产一类可复用的 `Role Pack Candidate`，同一个角色既能被能力 Pack 作为角色能力 / 角色约束引用，也能被团队办公室用作“秘书长”和“五顾问”的可替换角色；测试必须由用户视角智能体只通过 GUI 操作完成，我作为组织者、协调者、记录者参与，最终收口到商品阶段可交付的前端和后端，并能进入云端商品链路：上传到云服务器、支付购买、下载安装。

本轮只做计划，不执行实现。Owner 批准前不改代码、不改契约、不启动跨仓测试、不生成正式角色、不替换真实团队设置。

目标模式可用的一句话目标：

> 制作并验收角色制作台计划：由用户视角智能体通过 GUI 创建 candidate-only 角色包，验证能力 Pack 引用、团队办公室秘书长 / 五顾问槽位替换，并为秘书长提供音色与 VRM 形象选择；我负责组织、协调、记录和问题回填；最终交付商品阶段可用的前端与后端闭环，并支持角色包上传云服务器、支付购买、下载安装；全程保持 Role Pack 只是 Proposer、团队设置真相归基座、音色 / VRM 只存 asset ref，订单 / 支付 / License / Entitlement 真相归 `truzhen-cloud`。

## 1. 派活卡

| 维度 | 结论 |
|---|---|
| 我要做的事 | 参考能力 Pack 制作台测试计划与账本，制定角色制作台完善与测试计划。测试方式必须是用户视角智能体只走 GUI 操作；我作为组织者、协调者、记录者，负责排阶段、协调前后端、登记证据和回填问题。重点验证角色从“制作台候选”到“能力 Pack 引用”再到“团队设置替换”的闭环，并扩展到云端商品链路：上传、审核 / 上架候选、支付购买、License / Entitlement、下载和安装。 |
| 真实客户 / 场景证据 | Owner 2026-07-04 当前指令：要求角色可被能力 Pack 引用，也可用作秘书长、团队办公室五顾问，并能在团队设置页面 tab 下替换；秘书长可选音色、VRM。真实客户原话、团队使用记录、五顾问岗位定义仍缺，标记为“缺证据”。 |
| 版本 / 优先级 | Owner 尚未裁定这是当前客户主线、V4、未来功能建议还是 backlog。建议按“当前设计中测试计划”处理，执行前再裁定优先级。 |
| 最小可交付 | 不做完整角色市场。最小闭环是：用户视角智能体通过 GUI 完成 1 个秘书长角色候选 + 5 个顾问角色候选 + 1 个能力 Pack 引用样例 + 1 个团队设置替换样例 + 1 组音色 / VRM asset ref 选择样例；再完成 1 个云端商品化 dry-run：上传候选包、生成商品草稿、模拟支付购买、获得 entitlement、下载并安装到隔离团队环境。 |
| 商品阶段目标 | 最终不是只给候选 JSON 或 demo，而是给出可商品化的角色制作台前端与后端：前端具备完整 GUI 流程、团队设置 tab、音色 / VRM 选择、错误态和回执展示、上传 / 购买 / 下载 / 安装入口；后端具备 candidate、binding、Owner Gate、Receipt、asset ref、兼容性校验、审计证据、云端商品草稿、License / Entitlement 校验和安装回执。 |
| 真相源 | Role Pack 声明资产归 `truzhen-packs`；角色候选、正式化、启用和团队绑定事实归 `truzhenos`；跨仓 schema 归 `truzhen-contracts`；团队设置 UI 归 client；音色 / VRM 真实资源和渲染 provider 归基座资产库或 `truzhen-software` / 外部 provider；能力 Pack 只引用角色需求，不拥有角色真相；云上传、商品草稿、订单、支付、License、Entitlement、下载分发和云端安装包真相归 `truzhen-cloud`。 |
| 仓库 / 层归属 | 本轮计划文件归 `truzhen-packs` worktree。未来执行会涉及 `truzhen-client-web-desktop`、`truzhenos`、`truzhen-contracts`、`truzhen-packs`、`truzhen-cloud`，必要时涉及 provider / asset 仓。跨仓执行前必须另获 Owner 授权。 |
| 风险颜色 | 绿：计划文档、只读审计。黄：角色制作台 GUI、候选角色包、团队设置替换交互、商品草稿。橙：Role Pack / TeamRoleBinding / CapabilityRoleRequirement schema、秘书长表现层 asset ref、云端上传和下载安装链路。红：真实支付扣款、真实发布上架、真实语音克隆、真人肖像 / VRM 未授权资产、绕过 Owner Gate 替用户做正式决策。 |
| 是否改契约 | 本计划不改。执行时若现有契约无法表达 `RolePackCandidate`、`CapabilityRoleRequirement`、`TeamRoleSlotBinding`、`SecretaryAppearancePreference`，必须先出 `truzhen-contracts` 影响清单。 |
| 不允许碰的边界 | 不替换真实团队设置；不保存 raw 音频、raw VRM、真人肖像素材、密钥或 token；不把秘书长 / 顾问写成正式裁定者；不让能力 Pack 直接执行角色；不把 candidate-only 角色冒充已发布 Role Pack；未获 Owner 单独授权前不做真实扣款、不真实上架、不推送生产云、不生成真实 License。 |
| 用户如何验收 | 看计划是否覆盖目标、真相源、仓库归属、风险、契约、测试矩阵和禁区。执行后验收 GUI 截图、候选 JSON、团队设置绑定回执、音色 / VRM asset ref、负例阻断、结构审计和无 secret 扫描。 |
| 先输出什么 | 本轮先输出本计划。若 Owner 批准执行，先输出跨仓影响清单，再按 P0-P9 执行。 |

## 2. 设计原则

1. 角色包永远是 Proposer。秘书长和五顾问只能提出建议、草稿、质询和候选，不能批准、发送、执行或写正式事实。
2. 角色定义与团队绑定分离。角色制作台生产 `RolePackCandidate`；团队设置只保存“某团队槽位当前绑定哪个角色”的事实。
3. 能力 Pack 只引用角色需求。能力 Pack 可以声明“推荐 / 要求某类角色参与”，不能拥有角色本体，也不能直接变更团队设置。
4. 音色 / VRM 是表现层偏好。秘书长角色可选择音色和 VRM，但保存的是 `asset_ref` / `provider_ref`，不是原始文件、凭据或授权事实。
5. 所有替换可回滚。团队设置替换应产生候选、Owner 确认、基座登记和可反查回执；历史对话、历史 Receipt 不被改写。
6. 测试必须以用户视角为准。用户视角智能体只能通过 GUI 点击、输入、选择、确认和观察，不得直接改 JSON、直接打接口或绕过前端制造成功。
7. 商品化链路必须分层。`truzhen-packs` 只提供可商品化角色包资产；`truzhen-cloud` 持有上传、审核、商品、支付、License、Entitlement、下载分发真相；`truzhenos` 负责安装、启用、绑定和运行回执。

## 2.1 测试组织法

本轮执行时采用三方结构：

| 参与方 | 职责 | 禁止边界 |
|---|---|---|
| 用户视角智能体 | 扮演真实用户，只通过 GUI 使用角色制作台、能力 Pack 引用入口、团队设置 tab、音色 / VRM 选择入口。 | 不直接改文件、不直接调用后端、不读取数据库、不用脚本替代 GUI。 |
| 组织者 / 协调者 / 记录者（我） | 编排 P0-P10 阶段，协调 client / `truzhenos` / contracts / packs 的问题定位，记录截图、日志、候选 ID、回执、失败点和改进卡。 | 不替用户智能体操作 GUI 成功路径，不用后端假数据掩盖前端缺口。 |
| 验收智能体 / 监控程序 | 独立复核 GUI 行为、候选 / formal 隔离、Owner Gate、Receipt、secret 禁入、provider readiness。 | 不把测试日志当正式事实，不越过 Owner 授权启动红色动作。 |

用户视角智能体的每一步都要留下可复核证据：页面截图、操作时间、页面状态、网络响应摘要、候选 ID、后端 receipt 或明确 `blocked / provider_missing / not_ready`。我负责把这些证据整理成阶段报告和商品阶段前后端缺口清单。

## 2.2 商品阶段前后端目标

前端商品阶段必须具备：

- 角色制作台完整 GUI：创建、编辑、预览、校验、导出候选、错误提示。
- 团队设置页面 tab：展示秘书长和五顾问槽位、当前绑定、可替换角色、兼容性、风险色、确认 / 回滚入口。
- 秘书长表现层设置：音色、VRM 形象、默认值、清空、回退、provider readiness 显示。
- 能力 Pack 引用入口：能在能力 Pack 制作 / 编辑流中选择角色需求，并显示角色兼容性。
- 证据展示：候选 ID、Owner Gate 状态、Receipt ref、blocked / not_ready 原因、审计记录。

后端商品阶段必须具备：

- `RolePackCandidate` 创建、校验、导出和正式化前状态管理。
- `CapabilityRoleRequirement` 兼容性校验，不允许能力 Pack 复制或篡改角色本体。
- `TeamRoleSlotBindingCandidate`、Owner Gate、绑定回执和回滚策略。
- `SecretaryAppearancePreference` 只保存 `voice_asset_ref` / `vrm_asset_ref` / `provider_ref`。
- candidate / formal 隔离、权限校验、审计日志、Receipt 反查和 secret 禁入扫描。
- 云端商品化 API：上传候选包、生成商品草稿、提交审核候选、下载包签名、Entitlement 校验和安装回执。
- 支付 / 购买链路：测试环境可模拟支付成功 / 失败 / 退款 / entitlement 失效；真实扣款必须单独授权。

商品阶段收口标准：用户视角智能体能从 GUI 完成完整路径；我能用后端证据证明每个 GUI 成功都对应真实候选、Gate、Receipt 或诚实失败态；前后端缺口全部登记为可执行改进卡。

## 2.3 云端商品链路

角色制作台产物要能走完整商品链路，但每层真相必须分清：

| 环节 | 用户视角 GUI | 后端真相源 | 验收证据 |
|---|---|---|---|
| 打包 | 在角色制作台点击“导出 / 准备上架” | `truzhenos` 生成 candidate bundle；`truzhen-packs` 保存候选资产 | bundle 文件清单、hash、candidate ID |
| 上传云服务器 | 在 GUI 点击“上传到云端草稿” | `truzhen-cloud` 保存上传对象、草稿 listing、artifact hash | cloud upload receipt、listing draft ID |
| 审核 / 上架候选 | 在云端商品页提交审核 | `truzhen-cloud` 保存审核状态；正式上架需 Owner 授权 | review candidate、状态流转 |
| 支付购买 | 买家从商品页点击购买 | `truzhen-cloud` 保存订单、支付状态、License / Entitlement | sandbox payment receipt、entitlement ID |
| 下载 | 买家在已购页面下载 | `truzhen-cloud` 校验 entitlement 并签发下载链接 | download receipt、artifact hash |
| 安装 | 买家在本地 / 基座点击安装 | `truzhenos` 校验包、启用角色、写安装回执 | install receipt、enabled role pack version |
| 替换团队角色 | 买家在团队设置 tab 选择已安装角色 | `truzhenos` 写 TeamRoleSlotBinding 回执 | binding receipt、团队设置刷新证据 |

本轮如果没有 `truzhen-cloud` 授权，只能写商品链路计划和 mock / sandbox 验收设计，不能声称云上传、支付购买、下载或安装已经通过。

## 3. 三种实现路线

### 推荐路线：候选先行，契约后升格

先在 `truzhenos` 内部用 candidate-only 对象打通角色制作台、能力 Pack 引用样例和团队设置替换样例；所有产物标 `candidate_only / non_formal`。跑通后再决定是否升格到 `truzhen-contracts`。

优点：最小闭环快，能暴露 UI / 绑定 / 验收缺口。缺点：正式跨仓引用前还要做契约升格。

### 契约先行

先在 `truzhen-contracts` 定义 `RolePackCandidate`、`CapabilityRoleRequirement`、`TeamRoleSlotBinding`、`SecretaryAppearancePreference`，再改基座和前端。

优点：边界最稳。缺点：橙色风险高，若五顾问定义尚未确定，容易把临时产品判断固化成契约。

### 仅前端原型

只在 client 做角色制作与团队设置原型，不接候选对象和基座回执。

优点：最快看交互。缺点：容易把 UI 状态误读成真相源，不建议作为本轮目标。

本计划建议采用“候选先行，契约后升格”。

## 4. 目标角色与槽位假设

执行前需以产品定义或 contracts / 基座现状为准。本计划先采用以下测试假设，避免空泛：

| 槽位 | slot_ref | 角色用途 | 可替换 | 表现层 |
|---|---|---|---|---|
| 秘书长 | `team_office.secretary_general` | 组织对话、任务拆解、召集顾问、汇总候选 | 是 | 可选音色、VRM |
| 战略顾问 | `team_office.advisor.strategy` | 方向、优先级、取舍建议 | 是 | 默认无表现层 |
| 产品顾问 | `team_office.advisor.product` | 用户需求、产品路径、体验建议 | 是 | 默认无表现层 |
| 运营顾问 | `team_office.advisor.operations` | 执行节奏、流程、资源与交付风险 | 是 | 默认无表现层 |
| 财务顾问 | `team_office.advisor.finance` | 成本、现金流、ROI、预算风险 | 是 | 默认无表现层 |
| 法务风控顾问 | `team_office.advisor.legal_risk` | 合规、合同、授权、红线质询 | 是 | 默认无表现层 |

若 Owner 已有固定“五顾问”命名，执行前替换本表；不得在未裁定时把本表写成正式产品事实。

## 5. 候选数据形态

### 5.1 Role Pack Candidate

角色制作台至少要产出：

```json
{
  "role_pack_id": "role_pack://team-office-secretary-general",
  "candidate_only": true,
  "non_formal": true,
  "display_name": "团队秘书长",
  "target_use": ["team_office.secretary_general"],
  "compatible_slots": ["team_office.secretary_general"],
  "capability_reference_tags": ["coordination", "summarization", "task_candidate"],
  "style_summary": "先澄清目标，再拆任务、召集顾问、汇总候选。",
  "decision_style_summary": "只输出任务候选、沟通草稿、风险提示，不替 Owner 裁定。",
  "communication_style": {
    "structure": "目标复述 -> 事实缺口 -> 顾问分工 -> 候选汇总 -> 待 Owner 裁定项",
    "tone": "克制、清晰、组织者式",
    "forbidden_phrases": ["我已批准", "我已执行", "已经正式生效"]
  },
  "appearance_preferences": {
    "voice_asset_ref": "voice_asset://default-secretary-neutral",
    "vrm_asset_ref": "vrm_asset://default-secretary-vrm",
    "asset_ref_only": true
  },
  "risk_level": "medium"
}
```

### 5.2 Capability Role Requirement

能力 Pack 引用角色时，引用的是需求与兼容性，不是复制角色本体：

```json
{
  "capability_pack_ref": "capability-pack://sample-team-research",
  "role_requirements": [
    {
      "requirement_ref": "role-requirement://team-office/strategy-advisor",
      "accepted_role_pack_refs": ["role_pack://team-office-strategy-advisor"],
      "accepted_slot_refs": ["team_office.advisor.strategy"],
      "required_behavior": "proposal_and_challenge_only",
      "formal_authority": "none"
    }
  ]
}
```

### 5.3 Team Role Slot Binding

团队设置保存绑定候选，正式化归基座：

```json
{
  "team_ref": "team://owner-default",
  "slot_ref": "team_office.secretary_general",
  "role_pack_ref": "role_pack://team-office-secretary-general",
  "binding_status": "candidate_pending_owner",
  "replaces_role_pack_ref": "role_pack://previous-secretary-general",
  "requires_owner_gate": true,
  "history_policy": "do_not_rewrite_existing_receipts"
}
```

## 6. 预期候选资产结构

若 Owner 批准执行，本仓候选资产建议落在：

```text
/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/role-pack-candidates/team-office-v0/
  README.md
  candidate-set.json
  role-packs/
    team-office-secretary-general.rolepack.json
    team-office-strategy-advisor.rolepack.json
    team-office-product-advisor.rolepack.json
    team-office-operations-advisor.rolepack.json
    team-office-finance-advisor.rolepack.json
    team-office-legal-risk-advisor.rolepack.json
  role-slots/
    team-office-role-slots.json
  capability-role-requirements/
    sample-team-research.role-requirements.json
  appearance/
    secretary-appearance-preferences.json
  docs/
    role-studio-issue-ledger.md
    team-settings-binding-ledger.md
    role-pack-export-ledger.md
```

这些文件只能是候选资产，不参与 enabled pack 分发，不代表角色已发布。

## 7. 分阶段测试计划

### P0：开工预检

目标：确认范围、优先级、跨仓授权和当前 WIP。

验收：

- 每个目标仓独立运行 `git status --short --branch`。
- Owner 裁定这是当前客户主线、V4、未来功能建议还是 backlog。
- Owner 裁定是否允许读取 / 修改 / 测试 client、`truzhenos`、contracts。
- 指定用户视角智能体，只允许其通过 GUI 操作；我只做组织、协调、记录和问题定位。
- 明确截图、操作日志、后端候选 ID、Receipt ref、issue ledger 的证据目录。
- 若只批准计划，不进入 P1。

### P1：现状审计

目标：只读确认现有角色制作台、团队设置页面、能力 Pack 引用字段、音色 / VRM 资产入口。

验收：

- 用户视角智能体打开 GUI，按真实用户路径寻找角色制作台、团队设置 tab 和能力 Pack 引用入口。
- 列出已有 `RolePackManagePage` / 角色制作台入口。
- 列出团队设置页面 tab 结构和是否存在角色替换入口。
- 列出现有 Role Pack 字段与缺口。
- 我记录 GUI 页面证据、路由、按钮、空状态、错误状态和缺口归属。
- 若发现必须改 contracts，先停下输出影响清单。

### P2：角色制作台字段闭环

目标：制作台能创建秘书长和五顾问的 `RolePackCandidate`。

验收：

- 用户视角智能体必须通过 GUI 表单创建角色，不允许直接写 JSON。
- 可输入 display name、目标槽位、口吻、决策习惯、禁用短语、模型策略、风险等级。
- 每个角色都自动标 `candidate_only / non_formal`。
- 禁用短语覆盖“已批准 / 已执行 / 已发送 / 已生效”。
- 保存后生成候选 ID，不写正式角色。

### P3：候选导出

目标：把 1 个秘书长 + 5 个顾问导出为候选 bundle。

验收：

- 导出 `candidate-set.json`、6 个 `rolepack.json`、`team-office-role-slots.json`。
- JSON 合法。
- `target_use` 与 `compatible_slots` 对齐。
- 不含 raw secret、raw 音频、raw VRM 文件、运行态日志。

### P4：能力 Pack 引用角色

目标：能力 Pack 可以引用角色需求。

验收：

- 用户视角智能体必须通过能力 Pack 制作 / 编辑 GUI 选择角色需求。
- 样例能力 Pack 能引用 `role-requirement://...`。
- 兼容角色可被选择，不兼容角色被阻断。
- 能力 Pack 不复制角色全文，不改团队设置。
- 执行能力时仍只产 `CapabilityInvocationCandidate`，不因角色存在而变成正式执行。

### P5：团队设置替换

目标：团队设置页面 tab 下能替换秘书长和五顾问角色。

验收：

- 用户视角智能体必须从团队设置页面 tab 完成替换路径。
- 团队设置存在清晰的“角色 / 办公室”tab。
- 每个槽位显示当前角色、可替换角色、兼容性说明、风险提示。
- 替换动作先生成 `TeamRoleSlotBindingCandidate`。
- Owner 确认后才更新绑定事实并产生回执。
- 历史对话、历史候选、历史 Receipt 不被重写。

### P6：秘书长音色与 VRM

目标：秘书长角色可选择音色和 VRM 形象。

验收：

- 用户视角智能体必须从 GUI 中选择音色和 VRM，不允许直接写 asset ref。
- 设置只保存 `voice_asset_ref`、`vrm_asset_ref`、`provider_ref`。
- 默认音色 / VRM 可选、可清空、可回退默认。
- 未授权 asset、raw 文件路径、外链不可信资源被阻断。
- 不做真实语音克隆，不上传真人肖像，不保存 raw VRM。
- 没有表现层 provider 时显示 `provider_missing / not_ready`，不影响文本角色候选。

### P7：秘书长 + 五顾问运行烟测

目标：团队办公室实际对话时使用替换后的角色。

验收：

- 用户视角智能体从 GUI 启动团队办公室会话。
- 启动一个隔离团队会话。
- 秘书长按新角色口吻组织问题。
- 五顾问各自按绑定角色产建议 / 质询候选。
- 全部输出带 `candidate_only / non_formal`。
- Owner Gate 前不生成正式任务、正式记忆、正式发送或真实执行。

### P8：负例与安全测试

目标：证明系统不会把角色制作台变成主权绕过口。

负例：

- 把法务顾问替换到财务顾问槽位，若不兼容必须阻断或要求 Owner 高风险确认。
- 角色文案包含“我已批准”，保存必须提示 forbidden phrase。
- 能力 Pack 要求角色直接发送消息，必须 blocked。
- 秘书长选择 raw 本地音频路径，必须拒绝。
- VRM 资源缺授权信息，必须 `pending_human_review` 或 blocked。
- 未接通 provider 时，音色 / VRM 不能显示为 ready。

### P9：问题回填与验收报告

目标：所有缺口都转成制作台改进卡。

验收：

- 每个缺口有 issue ID、复现步骤、目标仓、风险色、契约影响、截图 / 日志证据。
- 不允许测试者绕过 GUI 手工改 JSON 后声称制作台可用。
- 收尾报告明确生命周期档位，不写“完成”冒充已发布。

### P10：商品阶段前后端收口

目标：把角色制作台从候选演示收口成商品阶段可交付的前端和后端。

验收：

- 前端：用户视角智能体能通过 GUI 完成创建角色、导出候选、能力 Pack 引用、团队设置替换、秘书长音色 / VRM 选择、回滚和错误恢复。
- 后端：每个前端成功动作都能反查对应 candidate、Gate、Receipt、asset ref 或诚实失败态。
- 前后端字段一致，不出现 UI 有字段、后端丢字段，或后端有状态、前端不展示。
- 负例完整：未授权 asset、provider 缺失、不兼容角色、主权越权、raw secret 输入都被阻断并有可读错误。
- 输出商品阶段收口报告：已达标前端能力、已达标后端能力、剩余 P0/P1 缺口、阻断项、下一轮执行卡。

### P11：云上传、支付购买、下载和安装

目标：验证角色制作台产物能成为可购买、可下载、可安装的云端商品。

验收：

- 用户视角智能体从 GUI 触发“准备上架 / 上传云端草稿”，不得手工把文件复制到云端。
- 云端草稿 listing 展示角色包名称、版本、适用槽位、风险声明、音色 / VRM asset ref 边界、价格候选和 license 候选。
- sandbox 购买能生成订单、支付回执、License / Entitlement；真实扣款必须另获 Owner 授权。
- 未购买用户不能下载；已购用户可下载对应 artifact，下载 hash 与上传 hash 一致。
- 安装动作经 `truzhenos` 校验包结构、签名 / hash、entitlement、Role Pack schema 和 forbidden artifacts。
- 安装后角色可出现在团队设置 tab 的可替换角色列表，替换仍需 Owner Gate 并产生 binding receipt。
- 退款、entitlement 失效、版本下架、包 hash 不一致时，下载 / 安装必须 blocked，并显示可读原因。

## 8. 测试任务矩阵

| ID | 测试任务 | 压测能力 | 预期证据 |
|---|---|---|---|
| TC-ROLE-01 | 创建秘书长角色候选 | 角色制作台字段与 Proposer 约束 | `RolePackCandidate` + 截图 |
| TC-ROLE-02 | 创建五顾问角色候选 | 多角色批量创建与兼容槽位 | 5 个 rolepack JSON |
| TC-ROLE-03 | forbidden phrase 检查 | 主权语言阻断 | 保存失败证据 |
| TC-GUI-USER-01 | 用户视角智能体只走 GUI | 真实用户路径，不绕过前端 | 操作日志 + 截图 + 页面状态 |
| TC-EXPORT-01 | 导出候选 bundle | 本仓候选资产结构 | 文件清单 + JSON 合法 |
| TC-CAP-ROLE-01 | 能力 Pack 引用角色需求 | `CapabilityRoleRequirement` | 引用样例 + 兼容校验 |
| TC-CAP-ROLE-02 | 不兼容角色引用 | 负例阻断 | blocked 证据 |
| TC-TEAM-01 | 替换秘书长 | 团队设置 slot binding | binding candidate + 回执 |
| TC-TEAM-02 | 替换五顾问 | 团队设置批量替换 | 5 个 binding candidate |
| TC-TEAM-03 | 回滚角色绑定 | 团队设置历史策略 | 新绑定回执，不改旧 Receipt |
| TC-VOICE-01 | 秘书长选择音色 | asset ref 保存 | `voice_asset_ref` |
| TC-VRM-01 | 秘书长选择 VRM | asset ref 保存 | `vrm_asset_ref` |
| TC-NEG-01 | raw 音频路径 | 凭据 / 资产边界 | 拒绝证据 |
| TC-NEG-02 | 角色直接批准 | 主权链 | blocked 证据 |
| TC-NEG-03 | provider 缺失 | readiness 诚实性 | `provider_missing / not_ready` |
| TC-PRODUCT-FE-01 | 商品阶段前端收口 | GUI 完整性、错误态、证据展示 | 前端验收截图 + smoke |
| TC-PRODUCT-BE-01 | 商品阶段后端收口 | candidate / Gate / Receipt / asset ref | 后端测试 + 反查证据 |
| TC-CLOUD-01 | 上传角色包到云端草稿 | 云端 artifact 与 listing draft | upload receipt + artifact hash |
| TC-CLOUD-02 | 提交审核 / 上架候选 | 商品状态机 | review candidate + 状态证据 |
| TC-PAY-01 | sandbox 支付购买 | 订单 / 支付 / License / Entitlement | sandbox payment receipt |
| TC-DOWNLOAD-01 | 已购下载 | entitlement 校验 + 下载签名 | download receipt + hash 校验 |
| TC-INSTALL-01 | 下载后安装角色包 | `truzhenos` 安装启用 | install receipt + enabled version |
| TC-INSTALL-NEG-01 | 未购买安装 | 授权阻断 | blocked: entitlement_missing |
| TC-INSTALL-NEG-02 | 包 hash 不一致 | 安全阻断 | blocked: artifact_hash_mismatch |

## 9. 验收命令建议

计划文档本身：

```sh
git diff --check -- docs/plans/role-pack-studio-team-office-test-plan-20260704.md
```

候选资产生成后：

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
# 运行本仓约定的敏感信息扫描，覆盖密钥、凭据、私钥、会话信息和云访问令牌模式。
```

未来跨仓执行时的验证建议：

- client：角色制作台测试、团队设置 tab 测试、Playwright smoke、typecheck。
- `truzhenos`：Role Studio candidate、TeamRoleSlotBinding、Owner Gate、Receipt 反查、provider readiness 测试。
- `truzhen-contracts`：schema / DTO 兼容性、反向依赖检查。
- `truzhen-packs`：Role Pack candidate JSON 合法、forbidden artifacts、候选资产口径审计。
- `truzhen-cloud`：upload draft、listing status、sandbox payment、License / Entitlement、download artifact hash、refund / entitlement invalidation。

无隔离 devserver、无回执、无截图或无日志时，不得声称 E2E 通过。

## 10. 待 Owner 裁定项

1. 本计划属于当前客户主线、V4、未来功能建议还是 backlog？
2. “团队办公室五顾问”的固定槽位名称是否就是战略、产品、运营、财务、法务风控？若不是，请给正式槽位表。
3. 是否允许下一步读取 / 修改 / 测试 `truzhen-client-web-desktop`、`truzhenos`、`truzhen-contracts`？
4. 角色制作台先走“候选先行”，还是必须先改 `truzhen-contracts`？
5. 秘书长音色 / VRM 的资产真相源归哪里：基座资产库、provider 仓、用户本地目录，还是云端资产服务？
6. 本轮是否只做秘书长有音色 / VRM，五顾问暂不做表现层？
7. 是否允许启动隔离 devserver 做团队设置替换回执验证？
8. 是否需要把候选资产同步登记到 `FEATURE_LEDGER.md`？若需要，应在实现或测试开始时单独更新账本，避免污染现有 WIP。
9. 是否批准采用“用户视角智能体只走 GUI + 我作为组织者 / 协调者 / 记录者”的测试组织法？
10. “商品阶段的前端和后端”是否要求本轮直接实现并验收，还是本轮只输出收口清单和执行卡，待下一轮授权施工？
11. 是否批准把 `truzhen-cloud` 纳入后续执行范围？允许动作是只读、sandbox 测试、修改、提交还是推送？
12. 云端商品化本轮是否只做 sandbox：上传草稿、模拟支付、模拟 License / Entitlement、下载和安装？真实扣款 / 正式上架默认红色动作，需要单独授权。
13. 角色包下载后的安装入口放在 `truzhenos` 本地 Pack 管理、团队设置页，还是云端已购页面触发本地安装？

## 11. 当前结论

角色制作台不应只做“编辑一段人设文本”。它要产出可审计、可引用、可替换、可回滚的 Role Pack 候选资产。

建议第一轮只做最小闭环：由用户视角智能体通过 GUI 完成秘书长 + 五顾问候选、能力 Pack 引用样例、团队设置替换样例、秘书长音色 / VRM asset ref。所有产物保持 candidate-only，所有正式绑定必须经 Owner + Base Gate，所有真实语音 / VRM provider 未接通时必须诚实显示 `provider_missing / not_ready / blocked`。

最终收口不以“生成候选文件”为完成口径，而以商品阶段前后端可用为目标：前端能让真实用户完整操作，后端能为每个操作提供候选、门控、回执、审计和诚实失败态；云端能完成上传、购买、下载、安装的 sandbox 闭环，真实支付和正式上架必须另经 Owner 授权。
