# PatchCandidate 复核候选台账

本文件记录短视频运营能力包制作台压测中的 P9：对 11 Code Assistant run 返回并由制作台承接的 PatchCandidate 做复核。它不是补丁应用记录，不代表文件已修改，也不代表能力包已安装、启用或发布。

## 复核候选

| 字段 | 值 |
| --- | --- |
| endpoint | `/v3/capability/studio/patch_candidate_review` |
| artifact status | `patch_candidate_review_ready` |
| source artifact | `code_assistant_result` / `PatchCandidateIntake` |
| decision values | `approved_for_apply_candidate`, `changes_requested`, `rejected` |
| default commercial path | `approved_for_apply_candidate -> apply Gate pending` |
| readiness effect | `patch_candidate_reviewed_pending_apply_gate` |

## 必须保留的证据

- `ReviewerRef`
- `ReviewEvidenceRef`
- `Summary`
- `Findings[]`
- `PatchCandidateRefs`
- `ArtifactRefs`
- `ExecutionReceiptRef`
- `SourceEvidenceRefs`
- `IntegrationPlanRefs`

## 主权边界

- `ApplyGateRequired=true`
- `ApplySupported=false`
- `ReadyForInstall=false`
- `Installable=false`
- `EnableSupported=false`
- `NoAutoApply=true`
- `FormalWrite=false`
- `CandidateOnly=true`
- `NonFormal=true`

复核通过只代表该 PatchCandidate 可以进入后续 apply 候选 Gate。实际 apply、install、enable、publish 必须另走 Owner + Base Gate + Gateway + Receipt。

## 禁止边界

- 不自动应用 diff。
- 不改工作区文件。
- 不把 review 通过写成 provider ready。
- 不把 review 通过写成正式 delivery。
- 不跳过黄金用例、Base Gate 或 Receipt。
- 不把 `changes_requested` / `rejected` 的补丁继续推进到 apply 候选。

## 下一步门槛

进入商用前需要补真实 PatchCandidate 文件审阅、apply 候选 Gate、应用后的构建 / 测试 / 黄金用例回归、正式 delivery、lifecycle install / enable、云市场审核链，以及 provider / adapter 的真实归属登记。
