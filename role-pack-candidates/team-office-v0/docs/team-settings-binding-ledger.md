# 团队设置绑定台账

当前没有真实团队设置替换回执。本文件只登记候选绑定预期。

绑定候选文件：`bindings/team-office-role-binding-candidate.json`。
安装后可替换角色目录候选文件：`bindings/team-settings-installed-role-catalog-candidate.json`。

| slot_ref | role_pack_ref | binding_status | owner_gate | receipt |
|---|---|---|---|---|
| `team_office.secretary_general` | `role_pack://team-office-secretary-general` | `candidate_pending_owner` | required | pending |
| `team_office.advisor.strategy` | `role_pack://team-office-strategy-advisor` | `candidate_pending_owner` | required | pending |
| `team_office.advisor.product` | `role_pack://team-office-product-advisor` | `candidate_pending_owner` | required | pending |
| `team_office.advisor.operations` | `role_pack://team-office-operations-advisor` | `candidate_pending_owner` | required | pending |
| `team_office.advisor.finance` | `role_pack://team-office-finance-advisor` | `candidate_pending_owner` | required | pending |
| `team_office.advisor.legal_risk` | `role_pack://team-office-legal-risk-advisor` | `candidate_pending_owner` | required | pending |

目录刷新要求：真实 `truzhenos` 安装回执、entitlement 校验、enabled role pack version 和团队槽位兼容检查齐全后，才能在团队设置 tab 显示为可替换角色。

P11 sandbox 执行请求：`tests/p11-sandbox-run-request-candidate.json` 已把“安装后团队设置 tab 替换秘书长 / 五顾问”列为授权后 GUI-only 阶段；真实替换仍必须由用户视角智能体在 GUI 中完成，并产生 Owner Gate 与 binding receipt。
