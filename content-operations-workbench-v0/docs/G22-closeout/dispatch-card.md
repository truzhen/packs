# G22 内容运营工作台候选生产与绝不自动发布：派活卡

日期：2026-07-24
状态：`blocked`；生命周期：`已接线`，未达到“候选包已验收”，未发布。

| 维度 | 本次裁定 |
| --- | --- |
| 要做的事 | 在隔离运行态中验证内容运营 Pack 的 install / enable、选题、母内容、渠道候选、公开稿清洗、人工发布包、周复盘与竖屏 MP4 候选；任何平台发布永远由 Owner 手工完成。 |
| 真实场景证据 | Owner 的 G22 指令：用已明确可用、去敏的 Truzhen 产品证据生产候选；不得登录、上传、发送或发布。该指令是边界与验收授权，不是 Codex Hands、OpenMontage 或素材权利的逐次 Owner Gate。 |
| 真相源 | 产品事实归获准产品证据；运行、Gate 和 Receipt 归隔离 Truzhen OS；发布事实归 Owner 与外部平台。本 Pack 只持候选声明。 |
| 仓库 / 层归属 | 仅 `content-operations-workbench-v0/**`；不改 OS、contracts、software、client、cloud 或公共 Pack 工具。 |
| 风险 | 黄：候选文本与人工发布包；橙：Gateway、Provider 和隔离 E2E；红：素材权利确认、真实 Codex / OpenMontage 执行、账号登录、上传、发送、发布。 |
| 契约影响 | 无。沿用 Pack 的 `candidate_only`、`never_formalizes`、`owner_manual_only` 与 Gateway-issued ref 约束。 |
| 禁止边界 | 不访问平台账号，不读取 Cookie/凭据，不上传、不发送、不发布；不自铸 OwnerDecision、Gate、Receipt 或执行证明；不提交源素材、MP4、运行日志或 test-store。 |
| 最小验收 | Pack 静态门禁全绿；隔离 OS 的 Gate / 素材权利拒绝路径可复核；正向产出仅在逐次 Gate、素材授权和 Provider 允许后才可开始。 |

## 本次阻断

1. 隔离 OS 对无可信 Owner origin 的 `11.execution.openmontage.content_video_candidate` Gate prepare 返回 HTTP 403 `trusted_owner_origin_required`。没有 Owner presence，不能 install / enable 或签发执行证明。
2. Codex Hands 的真实内容生产没有逐次 Owner Gate；未调用。
3. OpenMontage 只做 read-only readiness；视频请求在 `source_rights_confirmed=false` 时返回 HTTP 422，未渲染、未生成 MP4。没有素材权利确认与独立 Gate，不能进入正向视频路径。
4. 固定 SHA OS 在 `TRUZHEN_ANDROID_AUTO_RUNTIME=0` 时仍因 Windows VM 自动探测启动执行侧车，并加载本机持久 keyfile。进程随 devserver 停止；没有调用 Codex/OpenMontage、没有渲染。该环境越界自启动应移交 I04，未修复前不应重启本 lane。

前三项阻断是主权边界正常生效，不是可由本 Pack 文档或脚本修复的缺陷；第四项是公共 OS 组合根因。I04 关闭越界自启动后，解锁方应在同一隔离端口重新执行 install / enable、三种内容技能、视频探针、disable / uninstall / reinstall 和失败恢复验证。
