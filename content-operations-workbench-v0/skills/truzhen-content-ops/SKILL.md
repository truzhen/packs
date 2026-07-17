---
name: truzhen-content-ops
description: 将 Truzhen 的真实产品证据转成选题方向、中文母内容、抖音/B站/微信公众号等渠道候选、人工发布包和周复盘候选。用于内容运营工作台的选题雷达、内容生产和复盘；只生成内部审计稿与已清洗公开候选，不登录账号、不上传、不自动发布。
---

# Truzhen 内容运营

仅在 Pack 声明的隔离候选工作区内工作。先读取本次 run 提供的证据引用、Owner 判断和渠道约束；不得自行扩大读取范围。

## 选择模式

- `direction_radar`：提出最多三个方向，逐项写明受众问题、证据、风险、渠道与 Owner 待判断项。
- `content_production`：只对 Owner 已选择方向生成中文母内容、内部审计稿、已清洗公开候选和人工发布包。
- `weekly_review`：只依据真实指标输入生成 `keep/change/stop` 候选；缺指标时停止并列出最小补数清单。

## 渠道边界

- 只允许从本次 `channel_targets` 中选择渠道；未提供时默认只用：新版官网、GitHub、微信公众号、B站、小红书、知乎、视频号、抖音、YouTube、LinkedIn、X、TikTok。
- 不得自行增加 Wikipedia、媒体投稿、论坛、邮件、私信或其它未授权渠道。
- 方向雷达只说明“适合哪个渠道及原因”，不得把候选写成已经发布或已经获得平台数据。

## 固定流程

1. 核对输入是否包含 `transaction_ref`、`pack_version_ref`、`skill_id`、证据引用和允许写入的候选目录。
2. 区分 `已实现`、`已接线`、`已验收`、`已发布`；不得把计划、声明、候选或基础能力写成更高状态。
3. 生成内部审计稿，保留证据映射、限制、未知项和风险。
4. 生成公开候选时删除本机绝对路径、内部 ref、UUID、秘密、凭据、未公开客户信息和不可核验主张；证据引用只能放在结构化 `evidence_refs` 字段，不得进入可发布标题或正文。
5. 按 Execution Gateway 提供的 JSON Schema 输出模型候选；Gateway 负责计算候选 ref、内容哈希、泄漏扫描结果和 receipt candidate，并封装为 `output-contract.json` 的最终运行输出。模型不得自铸这些引用或哈希。
6. 发现证据不足、Provider 未就绪、T06/Base 证明不全或 Receipt 无法落账时，返回 `blocked`；不得用模板或本地假数据伪装成功。

## 硬边界

- 不登录任何社交平台，不读取浏览器 cookie，不上传，不点击发布，不发送私信，不抓取联系人。
- 不自铸 `OwnerDecision`、`decision_ref`、`owner_action_evidence_ref`、`run_id`、`nonce` 或 `receipt_ref`。
- 不把候选产物写成正式发布事实；人工发布是否发生由 Owner 和平台外部事实决定。
- 不把“免 Codex OAuth”改写成“免费 OAuth”，不从技术状态推导价格、许可、上线或客户采用等未提供事实。
- 不读取 `~/.codex` 共享运行态、全局 `auth.json` 或 Pack 目录外未显式授权的文件。
- 不在 Pack 内启动 CLI、Provider、sidecar、定时器或常驻进程；执行与调度由 Truzhen 基座通过 Gateway 提供。

详细输入输出字段以 [output-contract.json](output-contract.json) 为准。
