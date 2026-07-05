# OSS 证据矩阵

本文件只记录本轮能力包制作台压测使用的 GitHub 样本和风险口径，不复制第三方仓库内容，不形成法律结论。正式采用前必须重新拉取 GitHub 证据、检查 license、依赖、模型 / 素材来源、平台条款和安全风险，并经 Owner + Base Gate。

| 样本 | URL | 本轮用途 | 观察到的 license | 风险口径 | Pack 处理 |
| --- | --- | --- | --- | --- | --- |
| MoneyPrinterTurbo | https://github.com/harry0703/MoneyPrinterTurbo | 短视频草稿生成候选 | MIT | 可能涉及模型、素材、媒体工具和本地写入 | 可进入胶水补丁候选；真实 provider 未接通时 `provider_missing` |
| Pixelle-Video | https://github.com/AIDC-AI/Pixelle-Video | 短视频合成编排候选 | Apache-2.0 | 可能涉及模型工作流、ComfyUI、本地媒体处理 | 可进入胶水补丁候选；真实合成归 11 Execution Gateway |
| social-auto-upload | https://github.com/dreammis/social-auto-upload | 社媒发布草稿风险样本 | MIT | 涉及账号登录、网络上传、外部平台发布 | 默认 blocked；只生成发布草稿 / provider requirement，不生成直接执行胶水 |

## 证据要求

- `license`、`readme`、`maintenance`、`release`、`security`、`package_manager`、`docs` 至少要有可追溯 evidence ref。
- 只给 repo URL 不足以进入 `MinimalGluePlan`。
- 出现 external send、平台登录、账号凭据、cookie、真实发布、网络上传时，必须 blocked 或拆到 10/11 Gateway provider。
- Code Assistant 只能生成 PatchCandidate，不能保存 GPT/Codex 凭据原文，不能从 Pack 直接运行 Codex CLI。
