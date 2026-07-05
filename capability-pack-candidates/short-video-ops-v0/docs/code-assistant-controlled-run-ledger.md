# P16 受控 Code Assistant Run 证据台账

> 状态：`pending_authorization`
> 本文件只记录 P16 受控 run 证据。当前未获 P16 授权，未运行 Codex CLI，未消耗 token，未生成 PatchCandidate 文件，未产生 03 Receipt。

## 运行边界

| 字段 | 当前值 |
|---|---|
| run_status | `pending_authorization` |
| gateway | 11 Execution Gateway |
| gate_required | Owner + Base Gate |
| receipt_required | 03 Receipt |
| allowed_output | PatchCandidate / summary / adapter scaffold candidate |
| no_auto_apply | true |
| no_third_party_oss_execution | true |
| no_social_login_or_upload | true |
| no_raw_secret | true |

## 非充分证据

- 静态胶水代码、README patch 或 adapter scaffold 草稿不能替代 11 Gateway 受控 run receipt。
- Owner 本机 Codex CLI 登录态或 readiness 显示不能替代本次受控 run 证据。
- 手工放入的 PatchCandidate 文件不能解锁 P17，除非绑定 `DecisionRef`、`RunID`、`Nonce`、`ReceiptRef`、隔离输出目录和 no-auto-apply 证据。
- P8 `code_assistant_run_request_candidate` 只是运行请求候选；没有 11 Gateway receipt 时不能计入 P16 完成。

## 待填证据

P16 获授权并执行后，必须补齐：

evidence_contract：`docs/p16-controlled-code-assistant-run-evidence-contract.json`

### 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `owner_base_gate_approved_11_run_request` | `pending_authorization` | 11 Gateway decision / receipt evidence |
| `patch_candidate_created_no_auto_apply` | `pending_authorization` | PatchCandidate refs / no-auto-apply evidence |
| `gui_displays_run_receipt_and_candidate_refs` | `pending_authorization` | GUI test / smoke evidence |
| `third_party_oss_not_executed` | `pending_authorization` | run ledger / artifact scan |
| `no_social_login_upload` | `pending_authorization` | run and GUI no-login evidence |
| `no_auto_apply_no_formal_enable` | `pending_authorization` | PatchCandidate intake/review evidence |

- `DecisionRef`
- `RunID`
- `Nonce`
- `ReceiptRef`
- `PatchCandidateRefs`
- 隔离输出目录
- stdout / stderr 摘要
- 禁止动作检查结果

任一证据缺失时，P16 不得标为完成。
