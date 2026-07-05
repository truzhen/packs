# P15 GUI 实操证据台账

> 状态：`pending_gui_walkthrough`
> 本文件只记录 GUI 实操证据。当前尚未获 P15 授权，未启动前端 / 后端，未截图，未执行第三方 OSS，未登录或上传社媒。

evidence_contract：`docs/p15-three-candidate-gui-walkthrough-evidence-contract.json`

## 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `three_candidates_have_gui_screenshots` | `pending_gui_walkthrough` | screenshots |
| `clicked_steps_and_network_summaries_recorded` | `pending_gui_walkthrough` | clicked steps / network summaries |
| `blocked_reasons_match_backend_lifecycle` | `pending_authorization` | backend lifecycle evidence |
| `no_codex_cli_run` | `pending_gui_walkthrough` | walkthrough log / no-run evidence |
| `no_third_party_oss_execution` | `pending_gui_walkthrough` | artifact scan / execution log |
| `no_social_login_upload_or_formal_enable` | `pending_gui_walkthrough` | GUI/network evidence |

## P15-A draft generation

| 字段 | 证据 |
|---|---|
| candidate_pack_ref | `capability-pack://short-video-draft-generation` |
| github_repo_evidence_ref | MoneyPrinterTurbo GitHub evidence |
| frontend_url | `pending_gui_walkthrough` |
| backend_base_url | `pending_gui_walkthrough` |
| clicked_steps | `pending_gui_walkthrough` |
| network_responses | `pending_gui_walkthrough` |
| artifact_refs | `pending_gui_walkthrough` |
| screenshots | `pending_gui_walkthrough` |
| final_status | `pending_gui_walkthrough` |
| blocked_reasons | `provider_missing / candidate_only expected until verified` |
| forbidden_actions | no Codex CLI run; no third-party OSS execution; no social login/upload; no raw secret; no formal enable |

## P15-B composition orchestration

| 字段 | 证据 |
|---|---|
| candidate_pack_ref | `capability-pack://short-video-composition-orchestration` |
| github_repo_evidence_ref | Pixelle-Video GitHub evidence |
| frontend_url | `pending_gui_walkthrough` |
| backend_base_url | `pending_gui_walkthrough` |
| clicked_steps | `pending_gui_walkthrough` |
| network_responses | `pending_gui_walkthrough` |
| artifact_refs | `pending_gui_walkthrough` |
| screenshots | `pending_gui_walkthrough` |
| final_status | `pending_gui_walkthrough` |
| blocked_reasons | `provider_missing / candidate_only expected until verified` |
| forbidden_actions | no Codex CLI run; no third-party OSS execution; no social login/upload; no raw secret; no formal enable |

## P15-C social publish draft

| 字段 | 证据 |
|---|---|
| candidate_pack_ref | `capability-pack://short-video-social-publish-draft` |
| github_repo_evidence_ref | social-auto-upload GitHub evidence |
| frontend_url | `pending_gui_walkthrough` |
| backend_base_url | `pending_gui_walkthrough` |
| clicked_steps | `pending_gui_walkthrough` |
| network_responses | `pending_gui_walkthrough` |
| artifact_refs | `pending_gui_walkthrough` |
| screenshots | `pending_gui_walkthrough` |
| final_status | `pending_gui_walkthrough` |
| blocked_reasons | `external_send_risk / social_login_blocked / candidate_only expected until verified` |
| forbidden_actions | no Codex CLI run; no third-party OSS execution; no social login/upload; no raw secret; no formal enable |

## 完成判定

P15 通过前，每个剧本必须补齐：

- 至少一张制作台截图路径。
- 至少一条网络响应摘要。
- invocation / run request / candidate bundle / dry-run / preflight artifact refs。
- 明确 blocked reason。
- 明确未发生 Codex CLI run、第三方 OSS 执行、社媒登录、社媒上传、raw secret 泄漏和 formal enable。

任一证据缺失时，P15 状态保持 `walkthrough_incomplete` 或 `pending_gui_walkthrough`。
