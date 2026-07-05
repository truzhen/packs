# 候选 Pack Bundle 导出台账

本文件记录短视频运营能力包制作台压测中的 P4-P6 候选 bundle 导出。它不是正式 Pack 发布包，不代表 provider 接通，也不能用于 install / enable。

## 候选导出

| 字段 | 值 |
| --- | --- |
| endpoint | `/v3/capability/studio/candidate_bundle` |
| artifact status | `candidate_bundle_ready` |
| source artifact | `program_package` |
| readiness expectation | `provider_missing` |
| installable | `false` |
| enable supported | `false` |
| provider implementation included | `false` |
| candidate only | `true` |

## 文件清单

- `README.md`
- `manifest/capability-pack.json`
- `capabilities/capabilities.json`
- `docs/oss-evidence.md`
- `docs/code-assistant-invocation-ledger.md`

## 覆盖能力

- `capability-pack://short-video-draft-generation`
- `capability-pack://short-video-composition-orchestration`
- `capability-pack://short-video-social-publish-draft`

## 禁止边界

- 不写入真实 provider 实现。
- 不安装、不启用、不注册正式 lifecycle。
- 不运行 MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload。
- 不登录短视频平台，不读取 cookie，不上传或发布视频。
- 不把 candidate bundle 当成正式 delivery。

## 下一步门槛

P7 已补制作台 PatchCandidate 结果承接能力，但进入商用前仍需 11 受控 Code Assistant run 真实产出 PatchCandidate 文件、`truzhen-software` provider / adapter 候选、黄金用例回归、正式 delivery、Base Gate、Receipt 和云市场审核链。

P8 已补制作台 Code Assistant 运行请求候选能力，candidate bundle 文件清单需包含 `docs/code-assistant-run-request-ledger.md`，用于记录 11 run envelope 候选、Owner/Base Gate 待签发字段和“不可执行”状态。

P9 已补制作台 PatchCandidate 复核候选能力，candidate bundle 文件清单需包含 `docs/patch-candidate-review-ledger.md`，用于记录复核结论、review evidence、findings 和 apply Gate pending 边界。
