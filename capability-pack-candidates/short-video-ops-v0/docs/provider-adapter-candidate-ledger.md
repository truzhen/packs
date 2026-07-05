# P17 Provider / Adapter Candidate 台账

> 状态：`pending_authorization`
> 本文件只记录 provider / adapter candidate 归属。当前未获 P17 授权，未修改 `truzhen-software`，未安装或运行第三方 OSS。

evidence_contract：`docs/p17-provider-adapter-candidate-evidence-contract.json`

## 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `adapter_candidate_lands_outside_packs` | `pending_authorization` | truzhen-software diff / artifact refs |
| `readiness_reports_provider_missing_or_blocked` | `pending_authorization` | readiness API output |
| `gui_shows_provider_ownership_and_blocked_state` | `pending_authorization` | GUI evidence |
| `provider_implementation_not_in_packs` | `pending_authorization` | packs artifact scan |
| `no_third_party_oss_execution` | `pending_authorization` | provider readiness / no-run evidence |
| `no_raw_secret` | `pending_authorization` | secret scan |

| provider_candidate_ref | capability_ref | source_oss | target_repository | readiness_default | risk_class | credential_policy | forbidden_actions |
|---|---|---|---|---|---|---|---|
| `provider-candidate://short-video-draft-generation-adapter` | `capability-pack://short-video-draft-generation` | MoneyPrinterTurbo | `truzhen-software` 或外部 provider 仓 | `provider_missing` | orange | `secret_ref_only_no_raw_secret` | no OSS execution; no raw secret; no formal enable |
| `provider-candidate://short-video-composition-adapter` | `capability-pack://short-video-composition-orchestration` | Pixelle-Video | `truzhen-software` 或外部 provider 仓 | `provider_missing` | orange | `secret_ref_only_no_raw_secret` | no model/media execution; no raw secret; no formal enable |
| `provider-candidate://short-video-social-publish-adapter` | `capability-pack://short-video-social-publish-draft` | social-auto-upload | `truzhen-software` 或外部 provider 仓 | `blocked` | red | `secret_ref_only_no_raw_secret` | no social login; no upload; no external send |

## 完成条件

P17 通过前必须有：

- adapter scaffold 所在仓库路径。
- readiness check 输出。
- 风险边界文档。
- 禁入凭据扫描结果。
- `truzhen-packs` 内无 provider 实现。

## 非充分证据

- `truzhen-packs` 内的 scaffold、README、manifest 或示例目录不能替代 `truzhen-software` / 外部 provider 仓证据。
- 只有 provider manifest、没有 readiness endpoint / readiness receipt 时，不能解锁 P18。
- 一句 `provider_ready` 声明不能替代 Owner + Base Gate + Gateway + Receipt 绑定的 readiness 证据。
- vendor 第三方 OSS 源码、构建产物或运行时进本仓不仅不能算完成，还必须阻断 P17。

任一证据缺失时，P17 保持 `pending_authorization` 或 `adapter_candidate_incomplete`。
