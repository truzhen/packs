# Code Assistant 调用候选台账

本文件记录短视频运营能力包制作台压测中的 P3 调用候选信封。它不是 Codex CLI 真实运行日志，也不是 PatchCandidate、正式 Receipt 或 provider 接通证明。

## 候选信封

| 字段 | 值 |
| --- | --- |
| endpoint | `/v3/capability/studio/code_assistant_invocation` |
| artifact status | `code_assistant_invocation_candidate` |
| source session | `capability-studio session` |
| source package | `program_package` artifact |
| code assistant pack | `capability-pack://software-code-assistant` |
| allowed skills | `code.change_candidate`, `test.run` |
| model quota source | `owner_local_codex_cli_login` |
| credential policy | `credential_ref_only_no_raw_token` |
| candidate mode | `PatchCandidateOnly=true`, `CandidateOnly=true`, `NonFormal=true` |

## 隔离目录

| 目录 | 口径 |
| --- | --- |
| work dir | `process-dir://capability-pack-studio/<session_id>/glue-code-candidates` |
| output root | `artifact-root://capability-pack-studio/<session_id>/glue-code-candidates` |

## 允许文件

- `adapter/**`
- `tests/**`
- `manifest/**`
- `docs/**`
- `README.md`

## 禁止动作

- `no_raw_secret`
- `no_platform_login`
- `no_social_publish`
- `no_third_party_repo_run`
- `no_direct_provider_execution`

## 证据来源

- `MoneyPrinterTurbo` / `Pixelle-Video` / `social-auto-upload` 的 GitHub OSS evidence refs 由 04 `open_source_radar` artifact 向后传递。
- `IntegrationPlanRefs` 由 04 `program_package` artifact 向后传递。
- 调用候选信封由 04 写入 03 回执事件 `capability_studio_code_assistant_invocation_recorded`。

## 未发生的事

- 未执行 Codex CLI。
- 未消耗 Owner 的 GPT / Codex 会员模型额度。
- 未读取、保存或打印 raw token。
- 未安装、下载或运行第三方 GitHub repo。
- 未生成真实 PatchCandidate。
- 未生成正式执行 Receipt。

## 下一步门槛

真实 Codex CLI run 必须另走 11 Execution Gateway，并在 Owner + Base Gate 通过后产生可反查 Receipt。运行产物只能落在候选目录，且必须继续保留 `no_third_party_repo_run`、`credential_ref_only_no_raw_token` 和候选态边界。
