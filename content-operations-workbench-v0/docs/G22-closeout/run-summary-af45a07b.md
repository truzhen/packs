# G22 当前 OS 验收摘要

固定 OS：`af45a07b69c2ea44ed4b9d38612468bde0b2050d`；端口：`127.0.0.1:18222`。

- schema 1.1 候选 bundle 通过校验：方向、中文母内容、抖音/B站/公众号 internal/public 双稿、公开稿清洗和人工发布包均为 `package_ready_for_owner_review` / `pending` / `false` / `not_published`。
- OpenMontage 最小只读投影的 preflight 为 ready。官方 Gateway 经单次 Base Gate 生成了合成素材的 1080×1920 H.264/AAC、15 秒 MP4 候选；SHA-256 为 `df96d0cc7bccdd42608f184890958cc0ab2eeb0bac66f29093405053adb12467`，重启后可复取同 hash。
- 权利拒绝为 HTTP 422，已消费 Gate 重放为 HTTP 403，均未产生 task-store MP4。install/enable、disable/uninstall/reinstall 与幂等重装通过。
- 审计纠正后，Codex Hands 的当前首因不是宿主版本漂移：Docker digest 镜像和本地 os-08 模型均可复核，但当前 source lock 的 Docker 离线物料未就绪，官方 Provider candidate 返回 `controlled_download_stage_required`，无法形成 runtime/model-binding 投影。fresh store 的 Hands decision、run Receipt、session、trace、artifact 和 Provider install row 均为 0；没有调用。详见 [Docker Hands 审计纠正](audit-correction-af45a07b-docker-hands.md)。
- 周复盘因没有真实发布指标保持 blocked。未发生登录、上传、发送或发布。

运行态、合成素材、候选正文、MP4、cookie、日志与最小 registry 投影均在 closeout 后清理，不进入 Git。
