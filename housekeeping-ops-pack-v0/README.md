# 家政公司日常运营 Pack（v0）

家政/保洁公司**客户服务全生命周期**场景荚（Domain Work Pack / Scene Pack）。以用户视角制作，用于在 Truzhen V3 真实场景流程引擎上端到端运行验证。

## 生命周期主线

```
受理咨询(05 业务对象)
  → 家政顾问出方案(advice / 角色包 housekeeping_consultant)
  → 质检质询(challenge / 角色包 quality_auditor)
  → 对照确认门(compare_gate，wait①：Owner 对照顾问建议 vs 质检质询裁定)
  → 排期报价(capability / 能力包)
  → 整理派工报价单(task，高风险，停 pending_decision，需 Owner 主权确认)
  → Owner 派工审批(human_approval，wait②)
       ├─ approved → 发送服务确认(communication，受控候选)
       │             → 上门执行保洁(execution，受控候选)
       │             → 服务回执(receipt) → 归档进历史项目
       └─ rejected → 退回归档(close_rejected，零下游真实动作)
```

> 多智能体协作（顾问 + 质检）经 compare_gate 对照确认门显式编排，无 runtime 内隐藏 agent 回路；对照门未放行前不产派工候选。两个主权等待点（对照门 + 派工审批）均经真实 Owner resume 裁定，绝不自动放行。

## 主权与红线

- AI 全程 Proposer：顾问建议、质检质询、报价、派工单都只是候选；正式裁定权在 Owner + Base。
- 高风险 `dispatch_task` 停 `pending_decision`，必须 Owner 经真 01 Gate 裁定（`decision_ref/run_id/nonce` 由 Base 签发，禁前端自铸）。
- `communication`/`execution` 节点在 provider 未接通时诚实产受控 candidate（candidate_only / non_formal），不假发送、不假执行。
- 全链路 Receipt 落中央 03 账本，可回查 candidate→decision→formal→completion。

## 角色包 / 能力包

- 角色槽 `housekeeping_consultant`、`quality_auditor`：经 13 SlotBinding 绑定角色包，advice/challenge 节点经 08 模型网关真出候选；未绑定/模型未接通时诚实 `provider_missing`。
- 能力 `truzhen.capability.scheduling` / `truzhen.capability.quote`：经 04 能力中心 invoke 产 `CapabilityInvocationCandidate`。

## 运行验收

```sh
go test -race -count=1 -run TestHousekeepingOpsPackEndToEnd ./backend/tests/devserver/
```

测试读取本目录 `flows/customer-service-lifecycle.flow.json` 作为真实运行的流程图纸，经 `devserver.NewHandler`（与正式 devserver 二进制同一组合根）驱动全生命周期。测试报告见 `docs/测试报告.md`。
