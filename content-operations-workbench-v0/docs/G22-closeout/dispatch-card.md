# G22 内容运营工作台候选生产与绝不自动发布：派活卡

日期：2026-07-28。状态：`blocked`；生命周期：候选包已验收；未发布。

| 维度 | 本轮裁定 |
| --- | --- |
| 固定依据 | OS `f26d99c238a7d28af45ed96d5c5794550c313de3`、Software `614e771086fa3da9bf495a0ba7b89aca50170015`、Packs `7289f83684d697756be7f51034138fd0afa1d568`；其余四仓 SHA 见 evidence envelope。 |
| 真实场景证据 | Owner 授权仅使用 Truzhen 自有、去敏或合成素材并以候选包交付；发布事实归 Owner 与外部平台。 |
| 最小可交付 | schema 1.1 的 internal/public 双层候选包、人工发布包与历史本地视频候选证据；永不自动正式化或发布。 |
| 归属与影响 | 仅写 `content-operations-workbench-v0/**`；OS、Software、Provider、Gate 与 Receipt 均为只读消费，未改变契约。 |
| 风险 | 黄：候选文本与清洗；橙：Provider/Gateway；红：逐次 Owner/Base Gate、真实本地 Hands、素材权利与任何发布动作。 |
| 禁止边界 | 不登录、上传、发送、发布、私信、抓取联系人；不读取凭据、不改镜像标签、source-lock、registry 或 marker；不提交运行态、素材或视频。 |
| 验收 | 28 项 Pack 测试、JSON/Python/结构/禁品、敏感扫描、Go 门与 diff check；正向 Hands 仅在官方 Provider/模型均 ready 后开始。 |

## 2026-07-28 动态前置结果

1. 首次隔离服务错误消费固定 Software 忽略目录，因 ACL fail-closed；该请求在 candidate 前返回，`provider_install_candidates=0`，已作废且不作为最终根因。
2. 纠正后，以 task-owned、固定 `614e771…` 的临时 Software 副本运行官方 `prepare-offline-materials.py`；脚本仅允许 `image inspect/save`，没有 pull/fetch/build/run/login。
3. 官方预备在创建任何材料前拒绝：锁定 source image 为 `sha256:08b687…`，本机同标签镜像实际为 `sha256:7a86b0…`，且 `08b687…` 不存在。官方 verifier 因未生成材料而拒绝。
4. 所以没有 Provider candidate、Owner/Base Gate、Receipt、runtime、模型、Hands 或内容生成；平台动作始终为 0。该供给漂移属于 OS/Software 公共根因，本 Pack 无权修复。

详见 [当前基线审计更正](audit-correction-f26d99c-r22c-material-image-drift.md)。历史 OpenMontage 本地候选仅作既有媒体技术证据，本轮未重复长渲染。
