# G22 审计纠正：当前 SHA 的 Docker Hands 复核

日期：2026-07-25。固定 OS：`af45a07b69c2ea44ed4b9d38612468bde0b2050d`。

本记录取代“仅因宿主 Codex 版本漂移而阻断”的旧结论。隔离 lane 使用
`127.0.0.1:18222`、全新 task-owned store、空软件 registry、禁 Android 与 execution
sidecar 自动启动；没有复用旧 `751473a5` 的阻断结论。

## 已复核事实

- 首次隔离启动在 `go` 自动工具链获取提示处立即停止，未监听端口、未启动 devserver，也未产生 Hands 会话；后续复核强制 `GOTOOLCHAIN=local`、`GOPROXY=off`、`GOSUMDB=off`，只复用本机 Go 缓存。该启动前工具链事件不属于模型、Provider 或平台访问。
- 本地 Docker 可检查到 digest 固定的 `truzhen/codex-hands@sha256:08b687b435e7bb00bf803e2e2f099bfb05d9701ae7900f60a7fbde256aae97f5`（linux/arm64、非 root `10001:10001`）。
- 当前 OS 的 `provider-codex-cli-docker`、`truzhen_gateway`、`codeAssistantModelBindingProbe` 和通用 `pack.candidate.generate` 路径均存在；本地 os-08 readmodel 识别 `local-omlx / Qwen3.6-35B-A3B-4bit` 为 ready。
- 动态 binding probe 在写入模型 binding 证据后，因 Docker Provider/Runtime 尚未有可投影的已安装回执而返回 `code-assistant provider projection unavailable`；这不是模型不可用，也不是宿主 CLI 版本判断。
- 当前 `truzhen-software` source lock 的 Docker 离线物料目录缺失。官方 Provider install candidate 因而是 `not_ready` / `controlled_download_stage_required`；本轮禁止外网和 material fetch，不能安装、构建或把现有镜像旁路登记为 ready。

## 零会话结论

没有伪造 T06、没有调用 Codex CLI/Docker、没有模型 usage、没有 Pack 候选输出。fresh store 中：

- `code_assistant_run` decision：0；
- `code_assistant_run_*` Receipt：0；
- execution product session / trace / artifact：均为空；
- Provider install candidate：0。

本地模型 binding probe 是 readiness 证据而非 Hands 会话；其 formal ref 不写入 Pack。
运行摘要哈希：`7f8eaaa09552936062ba2b18c22add11f92cb56d7f1325834c223766647c6a78`。
临时 store、响应摘要和构建目录已移入本机废纸篓，18222 服务已停止；未提交运行态、日志、DB、凭据或产物。

## 当前阻断与后续

G22 仍为 `blocked`，首因是当前 source lock 的受控 Docker Provider 未 material-ready，无法完成官方安装、runtime 投影和动态 binding 投影。若 Owner 另行授权受控物料取得/安装，必须重新从空 store 做逐次 Gate；在此之前不得以本地已存在镜像、宿主 CLI 或旧 E2E 代替 readiness。

平台登录、上传、发送、发布、私信、联系人访问均为 0；`publication_authorized=false`、`publication_status=not_published`。
