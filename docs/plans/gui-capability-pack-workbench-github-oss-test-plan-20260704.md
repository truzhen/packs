# 能力 Pack 制作台短视频运营压测与产出计划

> 计划日期：2026-07-04
> 当前档位：`设计中`（分片进度：P3-P11 候选导出、dry-run 阻断、运行请求候选、PatchCandidate 承接与复核、lifecycle preflight 已接线；P12 安全内置能力 lifecycle 样本待 Owner 授权）
> 计划类型：测试计划 + 首批执行登记。本文件不代表真实 Pack 已发布、不代表 provider 已接通、不写 provider 实现。
> 参考报告：`/Users/li/Documents/truzhenv3/docs/reports/cloud-account-market-session-audit-and-improvement-20260703.md`

## 0. 纠偏后的目标

上一版把核心路径写成“最终用户一句话生成工作台”的演示题，这不准确。

本轮真实目标是：

1. 用“短视频运营过程中需要的一组能力 Pack 制作任务”来压测能力 Pack 制作台。
2. 在制作过程中使用代码助手能力，重点是受控调用 `codex cli`，并使用 Owner 的 GPT 会员模型额度 / token 环境。
3. 通过这一组制作任务，反向发现并完善能力 Pack 制作台。
4. 最终得到两类产物：完善后的能力 Pack 制作台问题清单 / 改进卡，以及一组短视频运营能力 Pack 候选。

短视频运营不是最终 GUI 里一句话生成工作台的演示题，而是制作台的测试素材和验收压力源。

## 0.1 首批执行进展

截至 2026-07-04，本计划已完成首批切片，并推进到 P3-P11 可验证切片：Code Assistant 调用候选信封、11 运行请求候选、候选 Pack bundle 导出、candidate bundle dry-run 阻断、PatchCandidate 结果承接与复核候选、lifecycle preflight 阻断。

1. `truzhenos` 04 能力制作台后端已补短视频 OSS 候选字段：GitHub 证据、`IntegrationPlan.CodeAssistantPolicy`、`MinimalGluePlan.AssemblyManifest` 和 HTTP/session 回归。
2. `truzhenos` 04 已新增 `/v3/capability/studio/code_assistant_invocation`，能生成 `code_assistant_invocation_candidate`，保留 GitHub evidence、integration plan、隔离输出目录、禁止动作、`owner_local_codex_cli_login` 和 `credential_ref_only_no_raw_token`。
3. `truzhenos` 04 已新增 `/v3/capability/studio/code_assistant_run_request`，能生成 `code_assistant_run_request_candidate`，映射 11 run envelope 字段，同时保持 `DecisionRef` / `RunID` / `Nonce` / `ReceiptRef` 待 Base 签发，`ReadyForExecution=false`。
4. `truzhen-client-web-desktop` 能力制作台 GUI 已展示 OSS 集成计划、代码助手能力包、补丁候选边界、GitHub evidence ref、Codex CLI readiness 只读检查，并新增“生成 Code Assistant 调用候选”“准备 Code Assistant 运行请求”按钮 / 卡片；live smoke 已验证 Code Assistant readiness 为 `ready`，假决策 run 被 `blocked`，候选信封和运行请求候选可生成。
5. `truzhenos` 04 已新增 `/v3/capability/studio/candidate_bundle`，能从 `program_package` 导出 candidate-only bundle 文件清单，明确 `provider_missing`、不可安装、不可启用、不含 provider 实现。
6. `truzhenos` 04 已新增 `/v3/capability/studio/code_assistant_result`，能承接 11 Code Assistant run 返回的 PatchCandidate 结果摘要，明确不可自动应用、需复核、不能正式写入。
7. `truzhenos` 04 已新增 `/v3/capability/studio/patch_candidate_review`，能从 PatchCandidate intake 生成复核候选；即使 approved，也只能进入 apply 候选 Gate，不自动应用、不安装、不启用。
8. `truzhenos` 04 已新增 `/v3/capability/studio/candidate_bundle_dry_run`，能对已导出的候选 bundle 做静态校验；当前独立 Capability Pack loader 未接线，预期返回 `candidate_bundle_dry_run_blocked / capability_pack_loader_missing`，不安装、不启用、不正式写入。
9. `truzhenos` 04 已新增 `/v3/capability/studio/lifecycle_preflight`，能把 candidate bundle 与正式 studio delivery 隔离：candidate-only bundle 返回 `lifecycle_preflight_blocked / candidate_bundle_not_delivery`，不创建 draft、不启用、不写正式回执。
10. `truzhen-packs` 已沉淀候选资产：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/`，含 3 个 Capability Pack 候选声明、Code Assistant 调用候选台账、运行请求候选台账、candidate bundle 导出台账、candidate bundle dry-run 台账、lifecycle preflight 商用缺口台账、PatchCandidate 承接台账和 PatchCandidate 复核台账。
11. 仍未完成：正式 lifecycle draft / readiness / promote / confirm 闭环、真实受控 Codex CLI run 产物文件、真实 apply 候选 Gate、第三方仓库执行、社媒登录 / 上传、真实 provider 接通、contracts schema 变更、云市场发布。

## 0.2 商用化下一步纠偏

P10 暴露的 `capability_pack_loader_missing` 不能简单理解为“补一个像 Domain Work Pack 一样的文件夹 loader”。基座治理已明确：Domain Work Pack 是可外迁 folder pack；Capability Pack 是基座常驻执行 provider + 声明。`truzhen-packs` 只能沉淀能力 descriptor / 候选声明，不拥有真实执行 provider、lifecycle enabled 指针、Gate 或 Receipt。

P11 `lifecycle preflight` 已完成接线：candidate bundle 现在会被明确识别为 candidate-only，而不是被误装成正式 Capability Pack。详细影响清单见：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-loader-lifecycle-commercialization-impact-20260704.md`

P11 执行规格和历史授权卡见：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-lifecycle-preflight-execution-spec-20260704.md`

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-cross-repo-execution-authorization-20260704.md`

商用化下一切片建议改为 P12 `安全内置能力 lifecycle 样本`：只选择基座内置安全 provider / fixture，验证 draft / readiness / promote / confirm 最小闭环。P12 不使用 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。

P12 执行规格和授权卡见：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

## 1. 派活卡

| 维度 | 结论 |
|---|---|
| 我要做的事 | 设计一组能力 Pack 制作台测试任务：让测试操作者像 Pack 作者一样，在 Truzhen 制作台中逐个制作短视频运营能力 Pack，过程中触发 GitHub 开源软件查询、Codex CLI 代码助手、胶水代码候选、Pack 候选导出、治理检查和问题回填。 |
| 真实客户 / 场景证据 | Owner 本轮纠偏明确：目标是完善能力 Pack 制作台，并产出一组短视频运营能力 Pack。示例软件：`MoneyPrinterTurbo`、`social-auto-upload`、`Pixelle-Video`。真实客户原话、运营团队真实流程、账号矩阵、日更压力、失败案例仍缺，标记为“缺证据”。 |
| 版本 / 优先级 | Owner 已裁定这是“制作台完善测试 + 能力 Pack 产出”的当前任务目标；是否属于当前客户主线、V4 还是 backlog 仍未明示。本计划按“当前设计中测试计划”处理，实际跨仓施工前需 Owner 再裁定。 |
| 最小可交付 | 不追求一次做完整短视频运营宇宙。最小闭环是制作 3 个能力 Pack 候选：短视频草稿生成、短视频合成编排、短视频发布草稿；每个都能暴露制作台缺口、生成候选文件、生成代码助手边界，并明确 readiness。 |
| 要做的事是什么 | 不是做最终短视频工作台，也不是直接接通三个开源软件；而是用这些能力 Pack 制作任务验证并完善“能力 Pack 制作台”的实际生产能力。 |
| 真相源 | GitHub repo 元数据归 GitHub API / repo 原文；代码助手调用事实归本地 Codex CLI provider / Execution Gateway；GPT/Codex 会员凭据和模型额度归 Owner 用户环境；Pack 声明归 `truzhen-packs`；制作台、Gate、Receipt、候选正式化归 `truzhenos`；真实 provider / adapter 归 `truzhen-software` 或外部 provider 仓；发布 / 商品 / License 归 `truzhen-cloud`。 |
| 仓库 / 层归属 | 本轮文档归 `truzhen-packs`。未来执行将涉及：client GUI、`truzhenos` Pack Studio / Gateway / Receipt、`truzhen-contracts` 候选对象契约、`truzhen-software` provider / adapter、`truzhen-packs` 能力声明资产。跨仓执行前必须重新授权。 |
| 风险颜色 | 绿：文档计划、只读 GitHub 查询。黄：制作台 GUI 生成候选 Pack。橙：Codex CLI 代码助手、胶水代码候选、ProviderRequirement 语义、候选 DTO。红：运行第三方代码、真实社交平台登录 / 上传、真实凭据、生产发布。 |
| 是否改契约 | 本计划不改。实际执行若新增 `CapabilityPackBuildCandidate`、`CodeAssistantInvocationCandidate`、`OpenSourceSoftwareCandidate`、`GlueCodeCandidate`、`ProviderAdapterCandidate`，必须先从 `truzhen-contracts` 出影响清单。 |
| AI 风险维度 | AI / Codex CLI 是代码助手和 Proposer，只能产候选、diff、测试建议、风险说明；不能批准、不能执行真实上传、不能写正式事实、不能绕过 Owner + Base Gate。 |
| 不允许碰的边界 | 不执行第三方 repo 代码；不运行安装脚本；不保存 GPT / OpenAI / 平台 token；不读取浏览器 cookie；不登录抖音、小红书、Bilibili、TikTok、YouTube；不真实上传视频；不把 provider 实现塞进 `truzhen-packs`；不把候选冒充已接通。 |
| 用户如何验收 | 看制作台操作录像 / 截图、每个能力 Pack 候选文件 diff、Codex CLI 调用回执、GitHub 查询证据、胶水代码候选 diff、制作台问题台账、Gate / Receipt 或明确 `not_ready` 证据。 |
| 先输出什么 | 本轮先输出优化后的测试计划。实际开跑前先输出跨仓影响清单和最小执行范围。 |

## 2. 参考报告抽象成测试纪律

从参考报告复用五条纪律：

1. 真相源先行：制作台显示、ReadModel、候选文件都不是正式事实。
2. 每个测试项写清仓库归属、风险色、契约影响、禁止边界、验收证据。
3. GUI 成功不等于完成，必须追到候选对象、Gate、Receipt 或 `not_ready / blocked`。
4. 未接通 provider、缺凭据、外部执行不可用时，必须诚实显示 `provider_missing`、`not_ready` 或 `blocked`。
5. 红 / 橙项只先做影响清单，不默认执行真实外部动作。

## 3. 短视频运营能力地图

本轮用短视频运营过程拆出能力 Pack 制作任务，而不是先做一个大而全工作台。

| 阶段 | 需要的能力 | 能力 Pack 候选 |
|---|---|---|
| 选题 / 热点 | 选题生成、关键词扩展、竞品拆解、脚本目标设定 | `short-video-topic-script-capability-pack-v0` |
| 脚本 / 分镜 | 文案、分镜、镜头提示词、口播稿 | `short-video-topic-script-capability-pack-v0` |
| 图片 / 视频生成 | AI 图片、AI 视频片段、素材检索、风格一致性 | `short-video-ai-generation-capability-pack-v0` |
| 配音 / 字幕 / 合成 | TTS、字幕、BGM、片段剪辑、9:16 导出 | `short-video-composition-capability-pack-v0` |
| 封面 / 标题 | 标题候选、封面候选、平台适配 | `short-video-cover-title-capability-pack-v0` |
| 排期 / 发布 | 账号矩阵、平台草稿、定时发布候选 | `short-video-social-publishing-capability-pack-v0` |
| 账号 readiness | 登录状态、平台风控、发布权限、凭据 secretref | `short-video-account-readiness-capability-pack-v0` |
| 回执 / 数据 | 发布回执、链接、播放 / 点赞 / 评论数据回收 | `short-video-publish-receipt-analytics-capability-pack-v0` |

本轮最小执行只覆盖 3 个候选：

1. `capability-pack://short-video-draft-generation`：参考 `MoneyPrinterTurbo`。
2. `capability-pack://short-video-composition-orchestration`：参考 `Pixelle-Video`。
3. `capability-pack://short-video-social-publish-draft`：参考 `social-auto-upload`。

其余能力作为 backlog 或第二轮。

## 4. 制作台需要被压测的能力

### M1：能力 Pack 类型入口

制作台必须能让作者选择“能力 Pack / ProviderRequirement / adapter 候选”路径，而不是只能做 Domain Work Pack。若当前契约或 loader 还不支持独立 Capability Pack 文件夹分发，制作台必须诚实显示“能力 Pack 候选”，不得假装已可装入。

### M2：需求拆解到能力边界

制作台要把短视频运营过程拆成多个能力边界，而不是生成一个巨大 Pack。每个能力必须写清：

- 输入 / 输出。
- 适用阶段。
- 风险色。
- 是否需要外部软件。
- 是否需要账号 / token / secretref。
- 是否会触发真实执行。
- readiness 期望。

### M3：GitHub 开源软件发现

制作台应支持只读 GitHub 查询，记录 repo URL、license、语言、更新时间、stars / forks、README 摘要、依赖、潜在风险和不确定项。

种子 repo：

| 种子 | 公开 repo | 测试用途 |
|---|---|---|
| `MoneyPrinterTurbo` | `https://github.com/harry0703/MoneyPrinterTurbo` | 一句话短视频生成、脚本 / 配音 / 合成流水线候选 |
| `Pixelle-Video` | `https://github.com/AIDC-AI/Pixelle-Video` | AI 自动短视频引擎候选 |
| `social-auto-upload` | `https://github.com/dreammis/social-auto-upload` | 多平台上传 / 排期发布候选 |

这些 repo 只作为测试种子；实际测试以当次 GitHub API / repo 原文为准。

### M4：Codex CLI 代码助手

制作台要能调用代码助手，但调用边界必须受控：

- 代码助手是外部能力 provider，不是 Pack 真相源。
- 调用入口应走基座受控 Execution Gateway / Code Assistant Gateway。
- Codex CLI 使用 Owner 本地 GPT/Codex 会员模型额度和本机登录态；Truzhen 不复制、不保存、不打印凭据原文。
- 每次调用要有 `CodeAssistantInvocationCandidate`：prompt、工作目录、允许文件、禁止动作、模型 / 额度来源、输出摘要、退出码、日志引用。
- 生成的胶水代码只能是候选，不能自动接线到真实 provider。
- 代码助手不得读取浏览器 cookie、平台账号、真实 token 或生产密钥。
- 超出模型额度、未登录、CLI 不存在、网络失败时，必须返回 `provider_missing / not_ready / blocked`。

### M5：胶水代码候选

制作台要能让 Codex CLI 生成 adapter scaffold、readiness check、错误映射、测试骨架和 README，但必须落在隔离过程目录或 future `truzhen-software` 候选目录，不写进正式 Pack 目录。

### M6：候选 Pack 导出

制作台要导出候选 Pack 文件或 candidate bundle。候选必须包含：

- `manifest.json` 或等价候选声明。
- `capabilities/capabilities.json`。
- README。
- 风险说明。
- GitHub repo 证据索引。
- Codex CLI 调用证据索引。
- readiness：`provider_missing`、`not_ready`、`blocked` 或 `degraded`。

### M7：问题回填

每个制作失败、字段缺失、UI 误导、代码助手不可用、GitHub 查询失败，都必须回填为制作台 issue，而不是让测试者绕过 GUI 手工改文件。

## 5. 上下文包

### 本轮已读

- `/Users/li/Documents/truzhen-packs/AGENTS.md`
- `/Users/li/Documents/truzhen-packs/CLAUDE.md`
- `/Users/li/Documents/truzhen-packs/README.md`
- `/Users/li/Documents/truzhen-packs/MODULES.md`
- `/Users/li/Documents/truzhen-packs/FEATURE_LEDGER.md`
- `/Users/li/Documents/truzhenv3/docs/reports/cloud-account-market-session-audit-and-improvement-20260703.md`（只读参考）

### 实际执行前必须读并重新授权的仓

| 仓库 | 职责 | 允许动作默认值 | 禁止边界 |
|---|---|---|---|
| `truzhen-client-web-desktop` | 制作台 GUI、用户操作流、截图 / Playwright | 只读 + 测试；修改需 Owner 授权 | 不伪造 UI 成功 |
| `truzhenos` | Pack Studio、Gateway、Receipt、Codex CLI provider、candidate / formal | 只读 + 隔离测试；修改需授权 | 不绕 Base Gate，不直接执行外部软件 |
| `truzhen-contracts` | 候选对象、capability schema、CodeAssistant DTO | 只读；改契约需另开卡 | 不消费方先行 |
| `truzhen-software` | provider / adapter 候选落点 | 只读；生成候选需授权 | 不提交真实密钥或第三方构建产物 |
| `truzhen-packs` | Pack 声明资产、测试报告、账本 | 文档 / 候选资产 | 不把 provider 实现写进 Pack |
| `truzhen-cloud` | 可选：作者上传、市场发布、License | 默认不碰 | 不发布、不上架、不处理真实支付 |

## 6. 测试环境设计

### 角色

| 角色 | 职责 |
|---|---|
| Pack 作者模拟用户 | 使用制作台逐个制作短视频运营能力 Pack。 |
| 制作台组织者 | 记录每步 GUI、GitHub 查询、Codex CLI 调用、候选文件、失败点。 |
| Codex CLI 代码助手 | 只在受控目录生成候选代码 / 测试 / README，不执行真实外部动作。 |
| 验收智能体 | 独立检查候选 / formal 隔离、provider 禁入、凭据禁入、readiness 诚实性。 |
| truzhen-monitor | 收集操作日志、诊断包、候选转态、异常。 |

### 数据隔离

- 使用隔离 devserver、隔离 DB、隔离过程目录。
- GitHub 查询可联网，但需记录 query、时间、provider、repo URL。
- Codex CLI 工作目录必须是隔离目录。
- 禁止使用真实社交平台账号。
- 禁止读取本机浏览器 cookie。
- 禁止从环境变量读取真实平台 token。
- GPT / Codex CLI token 只允许由用户本地 CLI 自己使用，Truzhen 不接触明文。

### 证据目录建议

```text
/Users/li/Documents/过程文档/capability-pack-studio-shortvideo-20260704/
  logs/
    main-task-register.md
    user-log.md
    behavior-matrix.md
    github-query-evidence.json
    codex-cli-invocation-ledger.md
    capability-pack-candidate-ledger.md
    studio-issue-ledger.md
  shots/
  artifacts/
    short-video-ai-generation-capability-pack-v0/
    short-video-composition-capability-pack-v0/
    short-video-social-publishing-capability-pack-v0/
    glue-code-candidates/
```

## 7. 分阶段测试

### P0：开工预检

目标：确认范围、跨仓授权、工具 readiness。

验收：

- 每个目标仓独立 `git status --short --branch`。
- 过程目录创建完成。
- 明确是否允许启动隔离 devserver。
- 明确 Codex CLI 是否已安装、是否可用、是否走 Owner 本地 GPT 会员模型额度。
- 如果 Codex CLI 不可用，制作台必须显示 `provider_missing`，不能继续假生成。

### P1：制作台能力 Pack 入口

步骤：

1. 打开制作台。
2. 选择“能力 Pack 制作”。
3. 创建 `short-video-ai-generation-capability-pack-v0`。
4. 输入能力目标：短视频脚本 / AI 视频生成 / 合成流水线能力。

验收：

- GUI 有明确能力 Pack / ProviderRequirement 路径。
- 若没有，记录 `studio_issue: capability_pack_entry_missing`。
- 系统生成候选 ID，而不是直接写正式 Pack。
- 页面展示风险色和真相源归属。

### P2：GitHub 查询与候选软件选择

步骤：

1. 制作台根据能力目标查询 GitHub。
2. 覆盖 `MoneyPrinterTurbo`、`Pixelle-Video`。
3. 用户选择可参考的软件。

验收：

- 每个 repo 有来源证据、license、更新时间、语言、主要能力、风险。
- license / 依赖 / 外部模型不确定时标 `pending_human_review`。
- 不因 repo 存在就标 provider ready。

### P3：Codex CLI 代码助手生成胶水候选

步骤：

1. 制作台创建受控 Codex CLI 调用候选。
2. prompt 要求 Codex CLI 阅读 repo README / API 文档，生成 adapter scaffold 方案。
3. 输出胶水代码候选到隔离目录。

验收：

- 先有 `code-assistant-invocation-ledger.md`，记录候选信封来源、endpoint、目录、允许动作、禁止动作、模型额度来源、凭据策略和 evidence refs。
- 真实运行后再追加 `codex-cli-invocation-ledger.md`，记录 prompt、目录、允许动作、禁止动作、退出码、输出摘要和正式 run receipt。
- 不记录 GPT token。
- 不运行第三方 repo。
- 生成文件只在候选目录。
- 失败时状态为 `not_ready`，并回填制作台 issue。

### P4：导出短视频 AI 生成能力 Pack 候选

目标候选：`short-video-ai-generation-capability-pack-v0`。

验收：

- 声明能力：脚本到视频候选、AI 图像 / 视频生成候选、素材生成 readiness。
- ProviderRequirement 引用 `MoneyPrinterTurbo`、`Pixelle-Video` 作为候选 provider，不写死等号关系。
- readiness 默认 `not_ready` 或 `provider_missing`。
- README 明确未接通真实 provider。
- 制作台可通过 `/v3/capability/studio/candidate_bundle` 导出 candidate-only bundle 文件清单，但不得标为 installable。

### P5：导出短视频合成能力 Pack 候选

目标候选：`short-video-composition-capability-pack-v0`。

验收：

- 声明能力：TTS、字幕、BGM、剪辑、9:16 导出。
- Codex CLI 产出 adapter 候选，但不运行 FFmpeg / MoviePy / 第三方依赖。
- 高风险或依赖缺失用 `blocked / not_ready`。
- 制作台能把上一步软件证据复用到新能力，而不是重复无证据生成。
- candidate bundle 必须复用 `SourceEvidenceRefs` 和 `IntegrationPlanRefs`。

### P6：导出多平台发布候选能力 Pack

目标候选：`short-video-social-publishing-capability-pack-v0`。

参考软件：`social-auto-upload`。

验收：

- 声明能力：平台草稿、定时发布候选、发布前人工确认、发布回执采集。
- 抖音、小红书、Bilibili、快手、TikTok、YouTube 登录 / 上传均为红色真实动作。
- 未提供账号 secretref 和 Owner Gate 时，全部 `blocked`。
- 不保存平台账号、不读取浏览器 cookie、不真实上传。
- candidate bundle 仍必须 `provider_missing / blocked-by-default`，不得解锁真实上传。

### P7：制作台缺口收敛

本轮先补一个制作台关键缺口：Code Assistant 真正跑完后，制作台必须能承接 11 返回的 PatchCandidate 结果摘要，而不是停留在“调用候选已生成”。

PatchCandidate 结果承接验收：

- 制作台通过 `/v3/capability/studio/code_assistant_result` 承接 run 摘要。
- 必须记录 `RunID`、`ProviderID`、`SkillID`、`PatchCandidateRefs`、`ArtifactRefs`、`ExecutionReceiptRef`。
- 必须保留 OSS `SourceEvidenceRefs` 和 `IntegrationPlanRefs`。
- 必须 `no_auto_apply=true`、`review_required=true`、`formal_write=false`。
- provider stub、非成功 run、缺 PatchCandidate 或缺 11 回执必须 blocked。

每个测试失败点都要转成制作台改进卡：

| 缺口类型 | 改进卡示例 |
|---|---|
| 没有能力 Pack 类型入口 | `studio-capability-entry` |
| GitHub 查询没有证据字段 | `studio-github-evidence-schema` |
| Codex CLI 不可配置 | `studio-code-assistant-provider-readiness` |
| 无法记录 token / 额度来源边界 | `studio-code-assistant-secret-boundary` |
| 胶水代码输出路径混乱 | `studio-glue-candidate-output-root` |
| 不能导出候选 Pack | `studio-capability-candidate-export` |
| readiness 只能写 ready | `studio-readiness-state-honesty` |
| 缺少失败回填 | `studio-issue-ledger-writeback` |

验收：

- 每个缺口有 issue ID、截图、复现步骤、目标仓归属、风险色。
- 不允许测试者绕过制作台手工修文件后声称制作台可用。

### P8：Code Assistant 运行请求候选

本轮补上“调用候选 -> 11 run 请求候选”的中间层，避免制作台停在 P3，也避免 04 直接触发真实 Codex CLI。

验收：

- 制作台通过 `/v3/capability/studio/code_assistant_run_request` 生成 `code_assistant_run_request_candidate`。
- 必须映射 11 run envelope 所需字段：`ExecutionEndpointRef`、`GateActionType`、`ProviderID`、`SkillID`、`TaskPrompt`、`WorktreeRef`、`TransactionRef`、`IdempotencyKey`、`AllowedScope`、`NetworkPolicy=none`。
- `DecisionRef`、`RunID`、`Nonce`、`ReceiptRef` 必须为空，等待 Base 签发。
- 必须显示 `GateStatus=owner_base_gate_pending`、`GateRequired=true`、`ReadyForExecution=false`。
- 不执行 Codex CLI、不消耗 token、不生成 PatchCandidate 文件。

### P9：PatchCandidate 复核候选

本轮补上 PatchCandidate intake 后的人工 / Owner 复核环节。复核通过不等于应用补丁，只代表可以进入后续 apply 候选 Gate。

验收：

- 制作台通过 `/v3/capability/studio/patch_candidate_review` 生成 `patch_candidate_review_ready`。
- 必须记录 `ReviewerRef`、`ReviewEvidenceRef`、`Findings[]`、`PatchCandidateRefs`、`ArtifactRefs`、`ExecutionReceiptRef`。
- 必须保留 OSS `SourceEvidenceRefs` 和 `IntegrationPlanRefs`。
- `Decision` 只能是 `approved_for_apply_candidate`、`changes_requested` 或 `rejected`。
- 即使 approved，也必须 `apply_gate_required=true`、`apply_supported=false`、`ready_for_install=false`、`enable_supported=false`、`no_auto_apply=true`、`formal_write=false`。

### P10：候选 bundle dry-run

验收：

- 制作台通过 `/v3/capability/studio/candidate_bundle_dry_run` 对已保存的候选 bundle 做 dry-run 校验。
- 当前独立 Capability Pack loader 未接线时，必须显示 `candidate_bundle_dry_run_blocked` 和 `capability_pack_loader_missing`。
- 静态文件清单校验可为 `passed`，但生命周期能力必须仍是 `NoInstall=true`、`NoEnable=true`、`InstallSupported=false`、`EnableSupported=false`。
- 不得用 Domain Work Pack 装入结果冒充 Capability Pack 支持。
- 无隔离 devserver 铁证时不得声称 E2E 通过。

### P11：Lifecycle preflight（已接线）

目标：让制作台从 candidate bundle 进入 lifecycle 前置检查，而不是假装能安装 / 启用独立 Capability Pack 文件夹。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-lifecycle-preflight-execution-spec-20260704.md`

历史授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-cross-repo-execution-authorization-20260704.md`

验收：

- 制作台能显示 candidate bundle 与正式 studio delivery 的区别。
- 无 delivery 时返回 `delivery_artifact_required`。
- candidate-only bundle 返回 `lifecycle_preflight_blocked / candidate_bundle_not_delivery`。
- evaluation readiness 未 ready 时返回 `provider_missing / not_ready` issue。
- UI 展示下一步所需条件：delivery、golden cases、evaluation ready、provider dependency、Owner/Base Gate、Receipt。
- 不写 enabled version、不写 Formal Receipt、不运行 Codex CLI、不执行第三方 repo。

### P12：安全内置能力 lifecycle 样本（待 Owner 授权）

目标：只用已在基座内有安全 provider / fixture 的能力样本验证 draft / readiness / promote / confirm 最小闭环，不直接跑短视频第三方 OSS。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

验收：

- draft 由服务端按 studio delivery 组装，不接收客户端伪 spec。
- readiness 聚合 delivery、golden cases、evaluation ready 和 provider dependency。
- confirm 经 01 Gate 两段并绑定 03 Receipt。
- 前端显示 enabled pointer、readiness issue 和 receipt ref。

### P13：GUI lifecycle 面板（待 Owner 授权）

目标：制作台在短视频候选包页显示 candidate bundle、delivery、readiness、promote、confirm 的分层状态和阻断原因。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-gui-lifecycle-panel-execution-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-cross-repo-execution-authorization-20260704.md`

验收：

- 用户不能点击“启用”绕过 readiness / Gate / Receipt。
- `provider_missing / not_ready / blocked` 有可读说明和目标仓归属。
- 生成制作台 issue ledger，而不是让测试者手工绕过。

### P14：商用化缺口审计（待 Owner 授权）

目标：对三个短视频能力 Pack 候选形成商用缺口台账。

当前审计草案：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-commercial-readiness-audit-20260704.md`

验收：

- 每个候选 Pack 有 lifecycle preflight 结果。
- 每个阻断项有仓库归属、风险色、修复建议和验收命令。
- 明确哪些进入当前主线，哪些进入 backlog。
- candidate-only、readmodel-only、demo-only 不能写成已发布。

### P15：三候选 GUI 实操验收（待 Owner 授权）

目标：用真实 GUI 路径跑完 3 个短视频能力 Pack 候选制作流程，并沉淀截图、行为日志、网络响应摘要和 issue ledger。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-three-candidate-gui-walkthrough-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-cross-repo-execution-authorization-20260704.md`

证据台账骨架：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/gui-walkthrough-evidence-ledger.md`

验收：

- 3 个候选包各有至少一张制作台截图。
- 每个候选包有关键点击步骤、网络响应摘要和 artifact refs。
- 每个候选包记录 blocked reason 和 issue refs。
- 明确没有 Codex CLI run、第三方 OSS 执行、社媒登录 / 上传、raw secret 泄漏和 formal enable。
- 任一证据缺失时保持 `pending_gui_walkthrough / walkthrough_incomplete`。

### P16：受控 Code Assistant 最小 run（待 Owner 授权）

目标：通过 11 Execution Gateway 做一次最小 Code Assistant run，只生成 PatchCandidate 文件和 03 Receipt，不应用补丁、不运行第三方 OSS。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-controlled-code-assistant-run-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-cross-repo-execution-authorization-20260704.md`

证据台账：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/code-assistant-controlled-run-ledger.md`

### P17：Provider / Adapter candidate（待 Owner 授权）

目标：把短视频能力所需真实 provider / adapter candidate 放到 `truzhen-software` 或外部 provider 仓，本仓只保留 ProviderRequirement 和证据索引。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-provider-adapter-candidate-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-cross-repo-execution-authorization-20260704.md`

证据台账：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/provider-adapter-candidate-ledger.md`

### P18：云市场 sandbox 链路（待 Owner 授权）

目标：在 `truzhen-cloud` sandbox 中验证商品草稿、License / Entitlement、下载分发和安装预检，不做真实支付或生产发布。

执行规格：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cloud-market-sandbox-spec-20260704.md`

授权卡：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cross-repo-execution-authorization-20260704.md`

证据台账：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/cloud-market-sandbox-ledger.md`

### 商用最终证据包（模板已建立）

目标：在 P12-P18 全部实证后，用同一份证据包判定是否能从 `not_commercial_ready` 进入 `commercial_ready_candidate`。

证据包模板：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-commercial-go-live-evidence-package-20260704.md`

授权路线图：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-authorization-roadmap-20260704.md`

## 8. 测试任务矩阵

| ID | 制作任务 | 压测制作台能力 | 预期产物 |
|---|---|---|---|
| TC-STUDIO-01 | 创建能力 Pack 项目 | 能力 Pack 入口 / 类型选择 | 候选 ID + GUI 证据 |
| TC-GH-01 | 查询 `MoneyPrinterTurbo` | GitHub 查询 / 证据结构 | repo 证据 + 能力映射 |
| TC-GH-02 | 查询 `Pixelle-Video` | GitHub 查询 / license 风险 | repo 证据 + AI 生成能力映射 |
| TC-GH-03 | 查询 `social-auto-upload` | GitHub 查询 / 红色上传动作识别 | repo 证据 + blocked 上传声明 |
| TC-CODEX-01 | Codex CLI readiness | 代码助手 provider readiness | `ready / provider_missing / blocked` 证据 |
| TC-CODEX-02 | 生成 adapter scaffold | 代码助手候选输出 | `GlueCodeCandidate` |
| TC-CODEX-03 | token 边界检查 | GPT 会员 token 不落地 | 日志无 token、只记录来源类型 |
| TC-PACK-01 | 导出 AI 生成能力 Pack | capability manifest / readiness | `short-video-ai-generation...` 候选 |
| TC-PACK-02 | 导出合成能力 Pack | 复用 repo 证据 / 多能力声明 | `short-video-composition...` 候选 |
| TC-PACK-03 | 导出发布候选能力 Pack | 红色动作 Gate / blocked | `short-video-social-publishing...` 候选 |
| TC-PACK-04 | 候选 bundle dry-run | 独立 Capability Pack loader 缺口诚实暴露 | `candidate_bundle_dry_run_blocked / capability_pack_loader_missing` |
| TC-LIFE-01 | Lifecycle preflight | candidate bundle 与正式 delivery 隔离 | `lifecycle_preflight_blocked / candidate_bundle_not_delivery` |
| TC-LIFE-02 | Readiness issue | provider / evaluation 未 ready 时诚实阻断 | `provider_missing / not_ready` issue |
| TC-LIFE-03 | 安全内置样本启用 | draft / readiness / promote / confirm 最小闭环 | 安全样本 enabled pointer + 03 receipt |
| TC-NEG-01 | Codex CLI 未登录 | 失败诚实性 | `provider_missing` |
| TC-NEG-02 | GitHub API 失败 | 失败诚实性 | 不编造 repo |
| TC-NEG-03 | 要求直接上传 | 主权链 | Owner Gate 缺失时 blocked |
| TC-NEG-04 | 要求保存账号密码 | 凭据纪律 | 拒绝 raw secret |
| TC-ISSUE-01 | 制作台缺口回填 | issue writeback | 制作台改进卡 |

## 9. 预期短视频运营能力 Pack 候选结构

### 9.1 `short-video-ai-generation-capability-pack-v0`

职责：

- 声明从主题 / 脚本到 AI 视频片段的能力需求。
- 参考 `MoneyPrinterTurbo`、`Pixelle-Video`。
- 输出只到候选，不生成真实视频。

必须字段：

- `capability_ref`
- `provider_requirements`
- `risk_class`
- `readiness_expectation`
- `candidate_inputs`
- `candidate_outputs`
- `owner_gate_required`
- `base_gate_required`
- `evidence_requirements`

### 9.2 `short-video-composition-capability-pack-v0`

职责：

- 声明配音、字幕、BGM、剪辑、竖屏导出能力需求。
- 参考 `MoneyPrinterTurbo`、`Pixelle-Video`。
- 不运行 FFmpeg / MoviePy / 第三方依赖。

### 9.3 `short-video-social-publishing-capability-pack-v0`

职责：

- 声明多平台发布候选、定时发布候选、发布回执采集。
- 参考 `social-auto-upload`。
- 所有真实上传默认 blocked，除非后续单独接入 provider、账号 secretref、Owner Gate 和 Receipt。

## 10. 制作台改进验收标准

制作台完成度不能靠“生成了文件”判断。至少要证明：

1. 能从 GUI 创建能力 Pack 候选。
2. 能从短视频运营阶段拆出多个能力边界。
3. 能只读查询 GitHub 并记录证据。
4. 能生成 Codex CLI 代码助手调用候选信封，并保护 GPT 会员 token 不落地。
5. 能在 11 受控执行接线后把 Codex CLI 输出存成胶水代码候选。
6. 能生成能力 Pack 候选文件或 bundle。
7. 能诚实表达 `provider_missing / not_ready / blocked`。
8. 能把制作台缺口写入 issue ledger。
9. 能区分 Pack 声明、provider 实现、运行时接线和正式事实。
10. 能输出可复核的截图、日志、候选 diff、调用回执。

## 11. 验收命令建议

文档计划本身：

```sh
git diff --check
```

候选 Pack 生成后：

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

Codex CLI 候选：

- 先生成 `code-assistant-invocation-ledger.md`，证明调用候选信封已受控。
- 只允许在隔离目录运行 CLI。
- CLI 工作目录必须受控。
- 不运行第三方 repo。
- 不安装第三方依赖。
- 不读取真实 token。
- 记录退出码和摘要。
- 失败必须进入 issue ledger。

GUI 自动化：

- 使用 Playwright / 浏览器控制记录截图。
- 每个关键按钮点击后记录页面状态、候选 ID、网络响应摘要。
- 不使用 mock 成功替代真实 GUI 路径。

## 12. 待 Owner 裁定项

1. 这项“完善能力 Pack 制作台 + 产出短视频运营能力 Pack”是否属于当前客户主线、V4、未来功能建议还是 backlog？
2. 是否有真实短视频运营客户 / 团队 / 场景证据：账号矩阵、发布频率、人工流程、当前痛点、失败案例？
3. 本轮实际开跑时，是否允许跨仓读取和测试 client、`truzhenos`、contracts、software？
4. 是否允许启动隔离 `truzhenos` devserver 做候选校验？
5. Codex CLI 是否作为正式代码助手 provider 纳入测试？若是，token / GPT 会员额度只归 Owner 本机环境，Truzhen 不保存明文。
6. Codex CLI 失败、未登录、额度不足时，制作台是否必须阻断并显示 `provider_missing / not_ready`？
7. 胶水代码候选是否允许写入过程目录？是否允许下一轮迁入 `truzhen-software`？
8. 下一施工切片是否进入 P12 `安全内置能力 lifecycle 样本`，用基座安全 fixture 验证 draft / readiness / promote / confirm，而不是直接跑短视频第三方 OSS 或真实 Codex CLI？
9. 短视频社交平台发布能力是否只测试候选和 blocked？若要真实测试账号和上传，这是红色动作，必须单独授权。

### 12.1 P11 历史授权与完成记录

P11 跨仓施工授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p11-cross-repo-execution-authorization-20260704.md`

Owner 已在 2026-07-04 明确授权 P11，并已完成 `truzhenos` 与 `truzhen-client-web-desktop` 接线。P11 授权语不自动延伸到 P12。

### 12.2 P12 可直接授权语

P12 跨仓施工授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

Owner 若确认进入 P12，可直接回复：

```text
授权按 P12 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒；只使用基座内置安全样本或 fixture 验证 lifecycle draft/readiness/promote/confirm。
```

## 13. 当前计划结论

本轮应从“制作台压测”出发，而不是从“最终用户工作台演示”出发。建议第一轮只跑三个能力 Pack 制作任务：AI 生成、视频合成、多平台发布候选。每个任务都必须触发 GitHub 查询、Codex CLI 代码助手、胶水代码候选、Pack 候选导出和制作台缺口回填。

验收时同时看两类结果：

1. 能力 Pack 制作台是否因此更完整：入口、证据、代码助手、候选导出、readiness、issue ledger。
2. 是否形成一组诚实的短视频运营能力 Pack 候选：不接真实平台、不藏凭据、不假装 provider ready。
