# 内容运营工作台 v0.2.0 成片能力实施与验收报告

日期：2026-07-18

## 结论

本轮把 `content-operations-workbench` 从“只生成文档候选”推进为“能生成真实 9:16 MP4 成片候选”。Pack 仍只声明工作模式、对象、能力依赖与治理边界；真实渲染归 11 Execution Gateway，外部软件事实归 `truzhen-software` registry，预览与 Owner 判断归独立客户端仓。未新增社交平台登录、上传或发布能力。

生命周期：`已接线`，专项验收已通过；根 EGR 前方仍有两个其它 worktree 任务占用 / 等待全局锁，本任务为避免无边界占用已退出队列，未宣称 `已验收 / 已发布`。

## 派活卡与边界

- 版本 / 优先级：Owner 已裁定为当前内容运营 Pack 完善主线。
- 场景证据：Owner 明确要求“生成视频，而不是文档”，并要求以真实用户 GUI 操作制作一期 Truzhen 宣传内容。
- 最小可交付：1–12 个 Owner 有权使用的图片 / 视频，可选 Owner 口播，生成 15–90 秒、1080×1920、30fps、H.264/AAC MP4 候选；支持预览、下载、Receipt 与 Artifact 反查。
- 真相源：Pack 声明归 `truzhen-packs`；Provider / License / Resource 事实归 `truzhen-software`；Gate、执行、Receipt、Artifact 归 Truzhen OS；预览和判断面归 client。
- 风险颜色：Pack / GUI 为黄；Provider 接入为橙；Owner evidence、Receipt、真实执行为红。本轮复用现有 OpenMontage Provider 与 Base-issued OwnerActionEvidence，不修改 contracts。
- 禁止边界：不登录抖音 / TikTok / X，不上传、不发布、不读取 raw secret、不把本机路径或渲染产物写入 Git，不允许前端自铸 decision / evidence ref。
- 砍入 backlog：自动配音、自动挑选素材、模板市场、多平台上传与发布、云端渲染。

## 实施结果

### Pack

- 升级到 `0.2.0`，增加 `video_render` 工作模式、`ContentVideoArtifactCandidate`、素材权利确认 Gate、OpenMontage capability/provider/software requirements。
- flow 在母内容候选后增加 Owner 选材、权利确认、受控渲染与 Artifact/Receipt 登记。
- forbidden 明确拒绝社交平台登录、上传与发布。

### Truzhen OS / Execution Gateway

- 新增确定性 FFmpeg + Pillow 竖屏渲染器，输入图片 / 视频和字幕镜头，输出 1080×1920 MP4。
- 有口播音频时使用 Owner narration；没有音频时明确生成 silent AAC，不冒充 TTS。
- multipart 候选接口限制 500 MB 请求、250 MB 单文件、1–12 个素材、15–90 秒；FFprobe 校验真实流格式。
- OwnerActionEvidence 必须由 Base 签发；成功后先写 03 Receipt，再由 ExecutionProductStore 保存隔离快照。
- 对抗复核曾发现“客户端自动 confirm、执行端只看 evidence ref”的主权缺口；已改为 GUI 人工确认卡，确认前真实渲染调用为零，确认后透传完整 Base proof，11 反查动作绑定并受 EmergencyStop 真闸约束。缺 verifier 时 fail closed。
- API 响应和 Receipt 不返回本机路径；预览端点只按 store-owned `artifact_ref` 提供视频，支持 Range。
- 同名并发请求使用唯一工作文件，快照完成后清理，避免覆盖与运行态双份残留。
- 修复执行侧车退出时遗漏导入 `terminate_active_codex_process_groups` 的既有根因；Ctrl-C 退出已无 NameError。

### Client

- 内容候选生成后显示素材、可选口播、使用权确认与“生成 9:16 MP4 候选”。
- 先通过后端 prepare/confirm 获取 Base-issued OwnerActionEvidence，再提交 multipart；不由前端自铸主权引用。
- 结果显示“成片候选 / 未发布 / Owner 口播或静音音轨”，提供 9:16 视频播放器、下载和 Artifact / Receipt / SHA256 技术详情。

### Software registry

- 复用 `resource-openmontage-media-gen-v1`，没有登记第二套平行 Provider。
- 登记 FFmpeg 8.1.2、Pillow 12.2.0、scene pack 消费关系和 adapter v0.2。
- License 结论仍为 review-required：当前只做本地子进程调用；商业分发、云服务或随安装包分发前必须完成 Owner/legal 复核。

## 验收证据

- 真实成片：`/Users/li/Documents/truzhen-software/artifacts/content-ops-v020/truzhen-content-ops-45s-v020.mp4`
- 媒体规格：45.000 秒，1080×1920，30fps，H.264 + AAC。
- SHA256：`61b0b831480509b44ebbdc8ddbd0e90e1e940ec0bad2e8c504548852b7f0bf50`。
- `python3 -m unittest sidecars/test_content_video_renderer.py`：通过，真实调用 FFmpeg。
- `go test ./backend/internal/execution/httpapi ./backend/internal/execution/openmontage`：通过，覆盖真实 HTTP → sidecar → MP4 → 03 Receipt → Artifact → Range 预览链。
- `bash scripts/v3-execution-gateway-acceptance.sh`：通过，含 Go race、coverage 与 forbidden scan。
- `bash scripts/check-openapi-wave7.sh`：通过，含 22 类对抗错误、route parity 与 P1 schema gate。
- Pack `go test ./...`：通过。
- Client 全量：219 个测试文件、1557 项测试通过；`npm run build` 通过。
- Software registry TOML 解析、四仓 `git diff --check`、无产物纳管检查：通过。

## GUI 验收

- 新版 Pack v0.2.0 已安装并启用到隔离 Truzhen OS 实例。
- 从 Chrome 以真实用户路径打开内容运营工作台；旧项目正确保持创建时绑定的 Pack v0.1.1，不自动漂移到新版。
- 隔离实例创建新项目时，Base 因缺少 Owner presence 返回受控阻断；未绕过 Gate。现有主实例能够打开真实内容项目和候选生成面。
- 视频控件与回执展示由组件级 GUI 测试覆盖；隔离实例的真实 45 秒执行链由后端 E2E 覆盖。待 Owner presence 可用时，再补“GUI 点击确认 → 成片播放器”最终截图证据。

## 未完成与发布条件

1. 等待当前两个全局 EGR 任务结束后，重新运行根 `scripts/verify.sh`。
2. Owner presence 可用时补一次完整 GUI 点击渲染证据。
3. 四仓提交后仍不自动 merge / push；需 Owner 复核本报告并明确授权。
4. 合并后再把功能账本状态从 `已接线` 更新到 `已验收`，随后才可进入 `已发布`。
