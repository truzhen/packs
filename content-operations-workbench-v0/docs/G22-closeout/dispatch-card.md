# G22 内容运营工作台候选生产与绝不自动发布：派活卡

日期：2026-07-25
状态：`blocked`；生命周期：候选包已验收；未发布。

## 2026-07-25 恢复轮（已收口）

| 维度 | 本轮裁定 |
| --- | --- |
| 固定依据 | OS `b843186cd2a7e93682f5d67e62d79a228376e368`；`af45a07b` 的 OpenMontage 验收仅为历史候选证据。 |
| 真实证据 | Owner 仅授权 Truzhen 自有、已去敏或合成素材；不得把合成素材表述为客户、发布或运营事实。 |
| 可做的动作 | 可在 Provider readiness、版本与 license 可反查且逐次本地 Base Gate 生效时，受控调用本机 Codex Hands / OpenMontage / ffmpeg 生成本地候选。 |
| 禁止边界 | 不登录、上传、发送、发布、私信、评论或抓取联系人；`publication_authorized=false`、`publication_status=not_published` 不得改变。 |
| 验收 | schema 1.1 bundle 的 internal/public 隔离、bundle-index 一一映射、公开稿清洗、候选 MP4 媒体检查、Gate/Receipt、生命周期重启与幂等、失败不留假产物。 |

结果：OpenMontage 本地渲染的历史候选证据保持有效，本轮不重复。`b843186` 审计确认固定 digest 镜像不能旁路官方物料、Provider/Runtime 投影；且 devserver 在正式 `autostart=0` 判定前读取持久 Windows guest-agent key，无法满足无凭据读取的隔离要求。没有继续 Hands、模型或平台动作。候选文本包已按 schema 1.1 通过校验，周复盘因没有真实发布指标保持 `not_ready`。本恢复轮的结论是“候选包已验收；未发布；整体 G22 有证据 blocked”。

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
2. 当前 SHA 的 source lock 不支持 task-owned 临时离线物料或已有 image 采纳，官方物料 fetch 会联网 `docker pull`，而其受锁输出目录位于未授权的 `truzhen-software`。未执行 Provider candidate、Hands 或模型调用。
3. OpenMontage 只做 read-only readiness，返回 `provider_missing`；视频请求在 `source_rights_confirmed=false` 时返回 HTTP 422，未渲染、未生成 MP4。没有素材权利确认与独立 Gate，不能进入正向视频路径。
4. 当前权威 OS SHA 是 `b843186cd2a7e93682f5d67e62d79a228376e368`。`MaybeLoadWindowsGuestTicketPrivateKey` 在 `TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0` 生效前无条件读取持久 key，且没有真实 devserver 的读取前禁用开关；按安全纠正已停止服务并释放 18222。旧 SHA 仅作历史证据，不能作为当前阻断。

前三项是主权与能力边界正常生效，不是可由本 Pack 文档或脚本修复的缺陷。隔离 lifecycle 已以受控本地 Owner presence 完成 Gate / Receipt 验证，但该 Gate 不授予 Codex Hands、OpenMontage 或任何平台操作的逐次授权。详见 [b843186 审计更正](audit-correction-b843186-windows-composition.md)。
