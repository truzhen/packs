# P13 GUI Lifecycle 面板证据台账

> 状态：`pending_authorization`
> 本文件只记录 P13 GUI lifecycle 面板证据。当前未获 P13 授权，未修改或测试 `truzhen-client-web-desktop` / `truzhenos`，未截图，未运行前端 smoke。

## 面板边界

| 字段 | 当前值 |
|---|---|
| target_repositories | `truzhen-client-web-desktop`, `truzhenos` |
| depends_on | `p12_safe_lifecycle_sample` |
| gui_status | `pending_authorization` |
| backend_status | `pending_authorization` |
| codex_cli_run | false |
| third_party_oss_execution | false |
| social_login_or_upload | false |
| formal_enable | false |

## 必填证据

P13 获授权并执行后，必须补齐：

evidence_contract：`docs/p13-gui-lifecycle-panel-evidence-contract.json`

### 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `candidate_delivery_readiness_panel_rendered` | `pending_authorization` | frontend test / screenshot |
| `enable_controls_gate_blocked` | `pending_authorization` | GUI assertion / disabled control evidence |
| `blocked_reasons_and_repo_ownership_visible` | `pending_authorization` | shell smoke / screen evidence |
| `lifecycle_status_summary_readable_by_gui` | `pending_authorization` | backend API / handler test output |
| `third_party_oss_run_from_gui_blocked` | `pending_authorization` | GUI test / blocked provider evidence |
| `codex_cli_run_without_11_gateway_blocked` | `pending_authorization` | GUI test / no direct run evidence |

- `client_worktree`
- `client_branch`
- `test_command`
- `test_result`
- `typecheck_result`
- `frontend_shell_smoke_result`
- `frontend_behavior_smoke_result`
- `screenshots`
- `network_response_summaries`
- `candidate_bundle_panel_state`
- `delivery_panel_state`
- `readiness_panel_state`
- `promote_panel_state`
- `confirm_panel_state`
- `enabled_pointer_rendering`
- `receipt_ref_rendering`
- `blocked_reason_rendering`
- `issue_ledger_refs`

## 通过条件

P13 通过必须同时证明：

- GUI 能区分 candidate bundle、delivery、readiness、promote、confirm、enabled pointer 和 receipt ref。
- 用户不能绕过 readiness / Owner+Base Gate / Receipt 点击出“启用成功”。
- `provider_missing / not_ready / blocked` 有可读说明和目标仓归属。
- MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 仍只是证据或 blocked sample。
- GUI 不直接运行 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
- 制作台缺口写入 issue ledger，而不是靠手工绕过。

任一证据缺失时，P13 不得标为完成。
