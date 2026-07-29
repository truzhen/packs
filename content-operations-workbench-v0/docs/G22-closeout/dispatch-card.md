# G22 内容运营工作台候选生产与绝不自动发布：派活卡

日期：2026-07-29。状态：`blocked`；生命周期：候选包已验收；未发布。

| 维度 | 本轮裁定 |
| --- | --- |
| 固定依据 | OS `be7fc745…`、Software `0768faa1…`、Cloud `e35df2da…`、Client `614f843f…`、Contracts `07854fef…`、Packs `7289f836…`。 |
| 真实场景证据 | Owner 授权仅使用 Truzhen 自有、去敏或合成素材并交付候选；发布事实归 Owner 与外部平台。 |
| 归属与影响 | 仅写 `content-operations-workbench-v0/**`；不改 OS、Software、Provider、Gate、Receipt 或协调现场。 |
| 风险 | 黄：候选文本与清洗；橙：Provider/Gateway；红：逐次 Owner/Base Gate、真实 Hands、素材权利与任何发布动作。 |
| 禁止边界 | 不登录、上传、发送、发布、私信、抓取联系人；不读取凭据、不伪造 ready、不改 marker/registry/source-lock；不提交运行态、素材或视频。 |
| 验收 | Provider 安装链必须先正式 ready；再以 schema 1.1 双稿、bundle-index、清洗、媒体、Gate/Receipt 和 lifecycle 证明候选。 |

## 2026-07-29 受控阻断

协调现场以全新 runtime 完成两阶段 Provider 安装、Owner/Base Gate、FormalReceipt、materials ready、离线安装、重启和幂等；外部平台动作、pull/fetch、容器运行、Provider/model 调用均为 0。该安装事实不等于可调用能力：当前 runtime/auth probe 与模型绑定仍未就绪，ReadModel 严格为 `sandbox_not_ready` / `install_completed_runtime_or_auth_probe_still_required`。

因此本轮没有 Hands、模型或 `pack.candidate.generate`，也不重复 OpenMontage 长渲染。下一步只能由 OS/Provider Owner 受治理地完成 runtime/auth probe 和本地 OMLX 模型绑定；ReadModel 同时 ready 后，再从新的空 store 进入一次 candidate-only 内容生成。

详见 [当前安装成功与 runtime 阻断审计](audit-correction-be7fc745-provider-install-pass-runtime-blocked.md)。
