---
title: 中性事实性催办话术模板集
category: 应收催办
source: Truzhen Pack Studio 合成撰写（无真实客户语料）
---

## 中性事实性催办话术模板集

本文件是应收催办的**中性事实性话术模板集**的知识陈述，供人工与角色查阅参照。**正文真相在 truzhenos `backend/internal/communicationgateway/dunningpolicy/policy.go`（`defaultTemplates`），Pack 只引用 `template_id`；Pack 层改写正文不生效**——os 侧沟通网关渲染候选时只按 `template_id` 从 Go 侧硬编码白名单取正文，本文件的正文抄自 Go 侧且必须逐字一致（含占位符名），若与 Go 侧不一致以 Go 侧为准并回报修正本文件。全部内容为合成撰写，**不含任何真实客户、金额、人名、单位名或单据号**。

### 占位符词表（`AllowedPlaceholders`，模板中只允许出现这七类，单花括号）

| 占位符 | 含义 | 取值来源 |
|---|---|---|
| `{invoice_ref}` | 单据号 | Sales Invoice 只读快照 `name` |
| `{amount}` | 未收金额 | 只读快照 `outstanding_amount` |
| `{currency}` | 币种 | 只读快照 `currency` |
| `{due_date}` | 到期日 | 只读快照 `due_date`；缺失时按过账日回退并显式标注 |
| `{payment_method}` | 付款方式 / 收款账户说明 | Owner 维护的收款方式声明 |
| `{contact_name}` | 联系人称呼 | 只读快照客户联系人字段；取不到时 Go 侧回退 `CustomerName` |
| `{company_name}` | 我方公司名 | Owner 维护的公司名声明 |

模板正文不得出现上述七类之外的变量；不得插入催办次数、逾期天数以外的推断性表述，也不得内嵌链接、附件或二维码占位符（通道能力与风险另归通道层）。占位符取不到事实时，Go 侧渲染以中性「待确认」填充，不编造付款方式或联系人。

### 模板 receivable_due_reminder_v1 · 到期前提醒（未逾期）

> {contact_name}您好，这里是{company_name}。单据 {invoice_ref} 的应付金额为 {amount} {currency}，到期日 {due_date}。如已安排付款请忽略本条。

适用：到期日前的一次性事实提醒；本波默认模板（`DefaultTemplateID`）。

### 模板 receivable_overdue_notice_v1 · 逾期事实提醒

> {contact_name}您好，单据 {invoice_ref}（{amount} {currency}）的到期日为 {due_date}，系统目前显示未结清。如与贵方记录不一致，请回复告知。

适用：到期后事实核对提醒，不施加压力。

### 模板 receivable_payment_method_v1 · 付款方式提示

> {contact_name}您好，关于单据 {invoice_ref}，应付金额 {amount} {currency}，到期日 {due_date}，可用付款方式为 {payment_method}。如需对账单请回复。

适用：附带付款方式说明的事实提醒。

### 模板 receivable_statement_confirm_v1 · 对账信息确认

> {contact_name}您好，{company_name}已就单据 {invoice_ref} 生成对账信息：金额 {amount} {currency}，到期日 {due_date}。请确认金额与到期日是否一致。

适用：请求核对对账信息，不承诺减免、延期或法律主张。

### 使用纪律

1. 模板只是**候选草稿素材的知识陈述**：真正的组装、占位符渲染、策略裁剪与白名单校验在 os 侧沟通网关（`dunningpolicy.Evaluate` / `dunningpolicy.TemplateByID`）完成，Pack 不发送、不执行、不改写、不是第二真相源。
2. 模板文本不得由 Pack 单方新增或改写：任何新增模板必须先落地为 Go 侧 `defaultTemplates` 新条目（中性事实性表述，重新走白名单），再回填本文件同步 `template_id` 与正文；未进 Go 侧白名单的自由文本一律 `blocked_template_not_allowlisted`（`CustomText` 非空同样一律阻断，本波无高级模式）。
3. 模板命中禁用词表（见 `forbidden-phrases.md`）时一票否决（`vetoed_forbidden_phrase`），只产 blocked/vetoed 候选并附命中词，不得产出可发送草稿；判定对象是**渲染后**文本，占位符注入内容命中同样否决。
4. 模板不承诺减免、延期、付款计划或任何法律主张；这些属高风险动作，须回 Owner + Base Gate。
