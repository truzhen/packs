# 销售 / CRM / 私域第一现金流 Pack 最小闭环设计

> 日期：2026-07-10  
> 生命周期：`设计中`  
> 版本定位：当前客户主线的最小验证候选；没有真实客户证据前，不进入契约、实现、接线或发布。  
> 本文只做设计与影响矩阵，不新增 Pack manifest、生产 API、数据库对象或运行时代码。

## 1. 派活卡

| 维度 | 裁定 |
|---|---|
| 我要做的事 | 把“第一现金流垂直场景荚”砍成唯一主轴：一个销售 / CRM / 私域客户，对一个真实客户对象做一次跟进，形成一个 Owner 可核对并确认的草稿，确认后经 Communication Gateway 产生一条可在 03 反查的受控回执。 |
| 真实客户 / 场景证据 | **缺证据**。当前只有 Owner 战略与跨版本文档，不具备客户原话、真实跟进记录或明确业务压力。本轮先准备最小验证脚本，不以推演推动工程。 |
| 最小可交付 | 一页访谈 / 走查脚本、一条闭环状态图、一份真相源与跨仓影响矩阵、一套 go / no-go 条件。 |
| 非目标 | 不做通用 CRM、不自建客户主数据、不做群发 / 自动回复 / 联系人抓取、不做直播打赏、不建五套独立平台、不实现生产接口、不修改 contracts。 |
| 真相源 | 外部 CRM 继续持客户、联系人和商机事实；`truzhenos` 05 只保存受控事务对象与必要快照，10 持发送候选与受控发送过程，03 持可反查回执；Pack 只声明工作模式、字段需求、流程和 Surface 意图。 |
| 仓库 / 层归属 | 本仓只拥有 Pack 设计；未来若验证通过，contracts 先判断现有形状能否复用，os 05/10/03 承接受控对象、沟通和回执，client 只渲染 ReadModel / Candidate / GatedAction，外部 CRM provider 保留外部权威。 |
| 风险颜色 | 本轮文档为绿；客户字段映射与 Pack 声明为黄；CRM provider / Communication Gateway 接入为橙；真实发送、批量触达、凭据和正式化为红。 |
| 契约影响 | 本轮无。验证通过后优先复用现有 TransactionObject、SendCandidate、OwnerActionEvidence、Receipt；只有不能表达最小闭环时才另出 contracts 影响清单。 |
| AI 角色 | AI 是草稿建议者，不是客户事实源、审批者或发送者。草稿不得因置信度、历史偏好或授权档位自动发送。 |
| 允许上下文 | 经 Owner 选择的单个客户最小字段、单次跟进目标、已有沟通摘要、必要产品资料；默认不读取整个通讯录、全部 CRM、历史私聊或无关客户资料。 |
| 禁止边界 | 不读取或保存 raw credential；不把 Pack / client 变成真相源；不让调用方自铸 decision / run / nonce / receipt；不绕过 Owner、Base Gate、03 Receipt；不向真实客户发送验证消息。 |
| 用户验收 | Owner 用一条脱敏但真实的跟进案例完成 20 分钟桌面走查，能够指出输入事实、草稿修改、确认动作、发送目标和回执反查位置；任何一步说不清即 no-go。 |
| 独立验收 | 未来施工后由独立验收主体证明：未确认不发送、目标错配阻断、重复幂等不重发、provider 缺失 fail-closed、回执能在 03 反查。 |

## 2. 唯一业务主轴

```text
外部 CRM 的一个客户 / 商机事实
        ↓ 只读映射与最小快照
一个跟进事务（05）
        ↓ AI 只生成 Candidate
一个可编辑草稿 + 目标 / 影响 / 依据
        ↓ Owner 明确确认，Base issuer 签发
一次 Communication Gateway 受控发送（10）
        ↓
一条可在 03 反查的 Receipt
```

最小闭环只接受单目标、单消息、单次确认。批量任务、自动节奏、线索评分、全量客户同步、营销自动化、佣金或支付全部砍入 backlog。

## 3. 最小字段与状态

### 3.1 输入快照

- 外部权威引用：`provider_ref`、`external_customer_ref`、可选 `external_opportunity_ref`。
- 人类可核对字段：客户显示名、当前关系阶段、上次联系时间、一次跟进目标。
- 依据：Owner 选中的沟通摘要或资料引用，不复制整套 CRM 数据。
- 运行引用：`transaction_ref`、`source_event_id`、Pack / version ref。

### 3.2 草稿候选

- 明确显示目标、正文、依据、风险提示和数据来源。
- AI 产物始终是 Candidate；Owner 修改后仍需本次明确确认。
- public request 只携业务字段与 `idempotency_key`，不得接受调用方提供 authority refs。

### 3.3 生命周期

`想法 -> 设计中` 是本轮终点。未来真实验证通过后才允许进入：

`契约已定 -> 已实现 -> 已接线 -> 已验收 -> 已发布`

文档、原型、fixture、mock、ReadModel 或 Candidate 单独存在都不能越级。

## 4. Top5 能力的最小切片

| 能力 | 本闭环只需要的薄片 | 明确不做 |
|---|---|---|
| Sovereignty Radar / 主权雷达 | 草稿确认面展示“谁拥有事实、谁将收到、动作风险、是否可撤销、回执落点”五项检查；本质是确认前摘要。 | 不做全局雷达、跨系统扫描或新治理控制台。 |
| Future Shadow | 仅展示“发送 / 不发送”两个短期结果提示，内容来源和不确定性可见。 | 不做长期预测引擎、自动决策或后端事实改写。 |
| 限时授权 | 仅用于同一客户、同一事务、短时有效的后续草稿候选；每次真实发送仍过 Base 硬地板，high / critical 永不委托。 | 不做永久授权、跨客户授权或自动发送授权。 |
| 关系状态 | 仅保存本事务需要的关系阶段快照与外部引用；外部 CRM 仍是权威。 | 不在 Pack / OS 内复制完整社交图或自建 CRM 真相。 |
| Surface Composer | 只组合现有 `object`、`candidate`、`receipt` 三类标准视觉单元。 | 不创造第八类卡片、不引入任意脚本式 UI 或独立 Surface 平台。 |

## 5. 跨仓影响矩阵

| 仓库 / 系统 | 未来职责 | 本轮动作 | 真相与边界 |
|---|---|---|---|
| `truzhen-packs` | 声明工作模式、字段需求、流程、风险和 Surface 意图 | 新增本文并登记账本 | 不持客户、关系、发送或回执事实 |
| 外部 CRM / provider | 客户、联系人、商机的外部权威 | 无 | 只经受控 provider 读取 / 写回；凭据不入 Pack |
| `truzhenos` 05 | 单次跟进事务与必要来源快照 | 无 | 快照不反客为主；保留 external refs |
| `truzhenos` 10 | 草稿候选、Owner 确认后的受控发送 | 无 | provider_missing / not_ready 必须阻断，不得假发送 |
| `truzhenos` 03 | 追加式回执与原文反查 | 无 | `receipt_ref` 必须真实 AppendReceipt，不接受拼接 ref |
| client | 只读展示、编辑 Candidate、提交 GatedAction | 无 | 不直连 CRM / provider，不自铸 authority refs |
| `truzhen-contracts` | 必要时定义跨仓 wire shape | 无 | 优先复用；确认缺口后才提兼容影响清单 |
| cloud | 可选的云会话 / 市场能力，不持本地客户事务 | 无 | 本最小闭环不依赖新增 cloud API |

## 6. 最小客户验证

在任何代码任务前，Owner 需提供一条脱敏真实案例，并完成以下 20 分钟走查：

1. 这次跟进由什么真实业务压力触发？
2. 外部 CRM 中哪些字段必须看，哪些字段绝不能复制？
3. Owner 在确认前最担心目标错、语气错、事实错还是时机错？
4. 哪些修改必须由 Owner 完成，哪些可以由 AI 建议？
5. 发送后什么回执足以证明动作发生且可追责？

go 条件：一个真实用户愿意用自己的单次跟进完成走查，且明确愿意为减少的时间 / 风险付费或继续试用。no-go 条件：只能描述“销售普遍需要”、无法提供单次案例，或价值依赖批量自动发送。

## 7. 未来施工的验收矩阵

| 验收项 | 必须证明 |
|---|---|
| 未确认 | 不产生 provider 调用和正式回执 |
| 目标错配 | session / customer / target 不一致时 fail-closed |
| 幂等 | 同一 `idempotency_key` 不重复发送或重复记账 |
| 授权 | authority refs 全由服务端签发；撤销 / 过期后不可使用 |
| provider | 缺失、不可用或超时返回明确 blocked 状态，不生成假 ready |
| 回执 | 发送结果可用 `receipt_ref` 在 03 查到原文与关联 transaction |
| 数据 | token、credential、完整联系人库不进入 JSON、日志、DB、ReadModel 或回执正文 |
| UI | Candidate / Formal 视觉隔离，确认面显示目标、内容、影响、run 和 evidence |
| 监控 | 复用 `truzhen-monitor` 记录候选、阻断、确认、provider 结果和回执链，不另建体系 |

## 8. 本轮裁定

- 状态保持 `设计中`。
- 不创建 Pack 成品、不改生产 API、不改 contracts、不声明客户验证完成。
- 下一步只有一个：取得一条脱敏真实销售 / CRM / 私域跟进记录，完成 §6 的桌面走查并留下 Owner 裁定。

