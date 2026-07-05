# PatchCandidate 结果承接台账

本文件记录短视频运营能力包制作台压测中的 P7：承接 11 Code Assistant run 返回的 PatchCandidate 结果摘要。它不是补丁应用记录，不代表真实 Codex CLI 已在本轮运行，也不代表 provider 接通。

## 候选承接

| 字段 | 值 |
| --- | --- |
| endpoint | `/v3/capability/studio/code_assistant_result` |
| artifact status | `patch_candidate_intake_ready` |
| source artifact | `program_package` + 11 Code Assistant run result summary |
| output kind | `PatchCandidate` |
| no auto apply | `true` |
| review required | `true` |
| formal write | `false` |
| installable | `false` |
| enable supported | `false` |

## 必需证据

- `RunID`
- `ProviderID`：`codex-cli-host` 或 `codex-cli-docker`
- `SkillID`：`code.change_candidate` 或 `bugfix.candidate`
- `PatchCandidateRefs`
- `ArtifactRefs`
- `ExecutionReceiptRef`
- OSS `SourceEvidenceRefs`
- `IntegrationPlanRefs`

## 禁止边界

- 不自动应用补丁。
- 不把 PatchCandidate 当作正式 provider 实现。
- 不安装、不启用、不发布。
- 不把 provider stub 当作真实 Codex 结果。
- 不读取或保存 GPT/Codex token。
- 不运行 MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload。

## 下一步门槛

进入商用前仍需真实受控 Code Assistant run 产生可复查 PatchCandidate 文件、`truzhen-software` provider / adapter 候选、黄金用例回归、正式 delivery、Base Gate、Receipt 和云市场审核链。
