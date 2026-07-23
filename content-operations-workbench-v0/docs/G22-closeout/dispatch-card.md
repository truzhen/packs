# G22 内容运营工作台候选生产与绝不自动发布：派活卡

日期：2026-07-24
状态：`blocked`；生命周期：隔离候选包 lifecycle 已验收，正向候选生产未获逐次授权，未发布。

| 维度 | 本次裁定 |
| --- | --- |
| 要做的事 | 在隔离运行态中验证内容运营 Pack 的 install / enable、选题、母内容、渠道候选、公开稿清洗、人工发布包、周复盘与竖屏 MP4 候选；任何平台发布永远由 Owner 手工完成。 |
| 真实场景证据 | Owner 的 G22 指令：用已明确可用、去敏的 Truzhen 产品证据生产候选；不得登录、上传、发送或发布。该指令是边界与验收授权，不是 Codex Hands、OpenMontage 或素材权利的逐次 Owner Gate。 |
| 真相源 | 产品事实归获准产品证据；运行、Gate 和 Receipt 归隔离 Truzhen OS；发布事实归 Owner 与外部平台。本 Pack 只持候选声明。 |
| 仓库 / 层归属 | 仅 `content-operations-workbench-v0/**`；不改 OS、contracts、software、client、cloud 或公共 Pack 工具。 |
| 风险 | 黄：候选文本与人工发布包；橙：Gateway、Provider 和隔离 E2E；红：素材权利确认、真实 Codex / OpenMontage 执行、账号登录、上传、发送、发布。 |
| 契约影响 | 无。沿用 Pack 的 `candidate_only`、`never_formalizes`、`owner_manual_only` 与 Gateway-issued ref 约束。 |
| 禁止边界 | 不访问平台账号，不读取 Cookie/凭据，不上传、不发送、不发布；不自铸 OwnerDecision、Gate、Receipt 或执行证明；不提交源素材、MP4、运行日志或 test-store。 |
| 最小验收 | Pack 静态门禁全绿；以权威 OS SHA 的隔离 test-store 跑通 install / enable / disable / uninstall / restart / reinstall，并复核 Gate、素材权利与 Provider 拒绝路径；正向产出仅在逐次 Gate、素材授权和 Provider 允许后才可开始。 |

## 本次阻断

1. 无可信 Owner origin 的 Gate prepare 返回 HTTP 403 `trusted_owner_origin_required`；这是安全拒绝路径，未形成正式 Gate 或 Receipt。
2. Codex Hands 的真实内容生产没有逐次 Owner Gate；未调用。
3. OpenMontage 只做 read-only readiness，返回 `provider_missing`；视频请求在 `source_rights_confirmed=false` 时返回 HTTP 422，未渲染、未生成 MP4。没有素材权利确认与独立 Gate，不能进入正向视频路径。
4. 当前权威 OS SHA `751473a5dd1b0ee2965397b541ef57a72e8ce273` 下，G22 以空的软件注册表根和 `TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0` 隔离启动；未发现任务 OS 的 execution sidecar 或持久 keyfile 加载。旧 SHA 的 I04 观察已由审计更正记录取代，不能作为当前阻断。

前三项是主权与能力边界正常生效，不是可由本 Pack 文档或脚本修复的缺陷。隔离 lifecycle 已以受控本地 Owner presence 完成 Gate / Receipt 验证，但该 Gate 不授予 Codex Hands、OpenMontage 或任何平台操作的逐次授权。详见 [审计更正](audit-correction-751473a5.md)。
