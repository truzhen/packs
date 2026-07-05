# P12 安全内置 Lifecycle 样本证据台账

> 状态：`pending_authorization`
> 本文件只记录 P12 安全样本 lifecycle 证据。当前未获 P12 授权，未修改或测试 `truzhenos` / `truzhen-client-web-desktop`，未创建 draft、未 promote、未 confirm、未生成 enabled pointer 或 03 Receipt。

## 样本边界

| 字段 | 当前值 |
|---|---|
| sample_scope | `safe_builtin_or_fixture_only` |
| target_repositories | `truzhenos`, `truzhen-client-web-desktop` |
| short_video_oss_usage | `evidence_only_not_executed` |
| codex_cli_run | false |
| third_party_oss_execution | false |
| social_login_or_upload | false |
| contracts_change | false |
| current_status | `pending_authorization` |
| evidence_contract | `docs/p12-safe-lifecycle-evidence-contract.json` |

## 授权前本仓预检记录

| 字段 | 当前值 |
|---|---|
| preflight_recorded_at | `2026-07-05` |
| authorization_status | `missing_owner_authorization` |
| cross_repo_work_started | false |
| truzhenos_modified_for_p12 | false |
| client_modified_for_p12 | false |
| packs_json_check | `JSON 合法` |
| packs_script_syntax_check | `脚本语法 OK` |
| packs_go_test | `go test ./... -count=1` 通过 |
| packs_diff_check | `git diff --check` 通过 |
| tracked_forbidden_artifacts | `tracked forbidden artifacts OK` |

说明：本记录只证明 P12 授权前的本仓候选资产和验证脚本处于可继续状态；它不是 P12 执行证据，不代表 lifecycle draft / readiness / promote / confirm 已接线，也不能解除 `docs/p12-pre-run-gate.json` 中的授权阻断。

## 必填证据

P12 获授权并执行后，必须先按 `docs/p12-safe-lifecycle-evidence-contract.json` 补齐机器可验条目，再补本台账人工摘要：

### 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `server_derived_draft` | `pending_authorization` | draft response / backend test output |
| `readiness_issues_or_ready_state` | `pending_authorization` | readiness response / backend test output |
| `promote_gate_preserved` | `pending_authorization` | promote response / lifecycle history |
| `confirm_requires_owner_base_gate` | `pending_authorization` | confirm blocked case / backend test output |
| `receipt_ref_bound_to_confirm` | `pending_authorization` | confirm response / 03 Receipt ref |
| `enabled_pointer_visible_after_confirm` | `pending_authorization` | history query / packs query |
| `gui_safe_sample_lifecycle_controls` | `pending_authorization` | frontend test / screenshot |
| `confirm_disabled_until_gate_receipt` | `pending_authorization` | frontend assertion / screenshot |
| `enabled_pointer_receipt_after_confirm` | `pending_authorization` | frontend smoke / receipt rendering |
| `third_party_oss_not_runnable` | `pending_authorization` | GUI smoke / blocked sample evidence |
| `codex_cli_not_run` | `pending_authorization` | execution log / no-run evidence |
| `third_party_oss_not_executed` | `pending_authorization` | artifact scan / execution log |
| `social_login_upload_not_happened` | `pending_authorization` | GUI/network no-login evidence |
| `raw_secret_not_saved` | `pending_authorization` | secret scan / redaction review |
| `contracts_not_changed` | `pending_authorization` | contracts git status / diff evidence |

- `truzhenos_worktree`
- `truzhenos_branch`
- `backend_test_command`
- `backend_test_result`
- `draft_response_summary`
- `readiness_response_summary`
- `promote_response_summary`
- `confirm_response_summary`
- `history_query_summary`
- `packs_query_enabled_pointer`
- `receipt_ref`
- `blocked_cases`
- `frontend_worktree`
- `frontend_test_command`
- `frontend_test_result`
- `frontend_smoke_result`
- `forbidden_action_checks`

## 通过条件

P12 通过必须同时证明：

- draft 由服务端安全样本 / fixture 组装，客户端伪 spec 不能覆盖。
- readiness 能解释 `ready / provider_missing / not_ready`。
- promote 不写 enabled pointer。
- confirm 必须带 Owner/Base Gate 结果和 03 Receipt ref。
- lifecycle packs 查询只在 confirm 成功后出现 enabled pointer。
- MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 未被执行。
- Codex CLI 未运行，未消耗 token。
- 未登录或上传社媒，未保存 raw secret。

任一证据缺失时，P12 不得标为完成。
