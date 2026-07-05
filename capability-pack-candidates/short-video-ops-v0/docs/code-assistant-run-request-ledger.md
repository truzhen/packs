# Code Assistant 运行请求候选台账

本文件记录短视频运营能力包制作台压测中的 P8：把 P3 调用候选信封整理成 11 Execution Gateway 可消费的运行请求候选。它不是 Codex CLI 真实运行日志，不代表模型额度已消耗，也不代表 PatchCandidate 已生成。

## 候选信封

| 字段 | 值 |
| --- | --- |
| endpoint | `/v3/capability/studio/code_assistant_run_request` |
| artifact status | `code_assistant_run_request_candidate` |
| source artifact | `code_assistant_invocation` |
| execution endpoint | `/v3/execution/code-assistant/runs` |
| gate action type | `code_assistant_run` |
| gate status | `owner_base_gate_pending` |
| provider id | `codex-cli-host` / `codex-cli-docker` |
| skill id | `code.change_candidate` / `bugfix.candidate` |
| network policy | `none` |
| ready for execution | `false` |

## Base 待签发字段

- `DecisionRef`
- `RunID`
- `Nonce`
- `ReceiptRef`

04 制作台必须保持这些字段为空，不得自行铸造。真实 Codex run 只能在 Owner + Base Gate 签发完整 envelope 后，由 11 Execution Gateway 执行并写入正式回执。

## 证据与范围

- `TaskPrompt` 来自 P3 调用候选的 `PromptSummary`。
- `WorktreeRef` 来自 P3 调用候选的隔离工作目录。
- `AllowedScope` 只允许 `adapter/**`、`tests/**`、`manifest/**`、`docs/**`、`README.md` 等候选包相关文件。
- `SourceEvidenceRefs` 和 `IntegrationPlanRefs` 继续贯通 GitHub OSS evidence 与装配计划。
- `IdempotencyKey` 由 04 针对候选请求稳定生成，但不等于正式执行授权。

## 禁止边界

- 不执行 Codex CLI。
- 不消耗 GPT / Codex 会员模型额度。
- 不读取、保存或打印 raw token。
- 不下载、安装或运行第三方 GitHub repo。
- 不生成 PatchCandidate 文件。
- 不把 `owner_base_gate_pending` 写成 `ready`。
- 不把本候选当成正式 11 `RunRequest`。

## 下一步门槛

进入商用前，需要 Owner 在 Base Gate 中确认本候选，Base 签发 `DecisionRef`、`RunID`、`Nonce`、`ReceiptRef`，然后由 11 Execution Gateway 调用真实 Code Assistant provider。真实 run 成功后，制作台只能通过 `/v3/capability/studio/code_assistant_result` 承接 PatchCandidate 证据，仍不得自动应用补丁。
