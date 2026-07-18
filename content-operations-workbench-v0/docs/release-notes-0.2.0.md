# 内容运营工作台 Pack 0.2.0

## 本版新增

- 新增 `video_render` 工作模式与 `ContentVideoArtifactCandidate`。
- 新增 `truzhen.content.video_render` 能力需求，绑定既有 `openmontage-family`，不在 Pack 内携带 Provider/runtime。
- 流程增加素材使用权确认、Base Execution Gate、竖屏 MP4 渲染和 Artifact Receipt 节点。
- 目标产物为 1080×1920、H.264/AAC MP4；口播音频由 Owner 可选提供，缺失时明确生成静音 AAC，不冒充 TTS。

## 边界

- 素材、视频产物、平台登录态和凭据不进入 Pack 或 Git。
- Pack 不登录、不上传、不发布、不抓联系人。
- 视频必须由 11 Execution Gateway 调用受控 Provider，经 OwnerActionEvidence、Base Gate、03 Receipt 和 Execution Artifact 快照后才可显示为成片候选。
- OpenMontage/FFmpeg 商业分发与网络服务仍受 AGPL/GPL review gate 约束。
