# G22 审计纠正：R22A/R22B 新基线的离线 Hands 契约缺口

日期：2026-07-25。唯一基线：OS
`4784c2dc9048968542262349844d413c252f0e51`；Software
`346018276e5cc970f10b513ff9413520763f1323`。

本轮只改 Pack closeout。OS 与 Software 均只读；没有改写忽略的物料目录、registry、
source lock、image、key 或运行时。

## R22A 安全开关复核

隔离 devserver 使用全新 task-owned `HOME`、store、workroot 和 keys，在
`127.0.0.1:18222` 启动时明确设置：

- `TRUZHEN_WINDOWS_VM_AUTODETECT=0`；
- Android 与所有 execution sidecar autostart 均为 `0`；
- 无 dummy key、caller marker 或替代密钥；模型 refresh 禁用。

R22A 的 `ResolveWindowsVMAutodetect` 在 `MaybeLoadWindowsGuestTicketPrivateKey`
之前短路；本轮启动日志没有 guest-agent key load、sidecar 或 Android 自启动。服务仅做
只读 Provider candidate 查询后停止，端口已释放。

## R22B 物料复核与当前失败

Software 的官方离线验证器对现有 ignored 目录返回 `codex material verification: OK`；
`material-set.json` SHA-256 为
`a1a92b054b8a94c2f660320a5e35813106e7b580e194ee3f1c35589d5b585485`，其绑定 image digest
为 `sha256:08b687b435e7bb00bf803e2e2f099bfb05d9701ae7900f60a7fbde256aae97f5`。

但当前 OS `GET /v3/execution/code-assistant/providers/install-candidates?provider_ref=provider-codex-cli-docker`
返回 HTTP 503：`providerinstall: source registry or provider assets invalid: Docker material directory is not owner-only`。
响应正文 SHA-256：`948e8eaf66c4944a16847f025cb01bc88db4f142624a7921a854e6d750b1f77f`。

根因是可复现的跨仓契约不一致：

1. Software R22B 产出 `materials` 为 `0700`，但其 `tarballs` 子目录为 `0755`；OS
   `dockerOfflineMaterialsReady` 要求两者都是 owner-only `0700`。
2. Software R22B 的受控证据为 `material-set.json`；当前 OS resolver/runtime binding
   只认可由其旧 fetch 链、receipt 和 `MarkReady` 生成的 `material-ready.json`。
   R22B 的 existing-image export 没有被当前 OS 正式采纳，不能由 G22 手工补 marker。

因此执行停在正式 Provider candidate 之前：没有 Base Gate、T06、Provider install、Runtime
投影、model usage、preflight/final Receipt、Task Hands Receipt 或 `pack.candidate.generate`。
fresh store 的 `provider_install_candidates=0`、`base_gated_action_issues=0`；execution
product 的 session/trace/artifact 均为空。所有平台动作仍为 0。

## 交接边界

这是 OS/Software 的公共 root cause，不可在 Pack 中 TDD 修复：需要由两仓定义并测试同一份
离线 material adoption 契约（权限、`material-set.json` 验真、Owner Gate/Receipt、runtime
projection），或由 Software 输出当前 OS 已签名认可的物料状态。不得通过 chmod、手写
`material-ready.json`、登记 ready 或联网 fetch/pull 绕过。

旧 OpenMontage 视频验收继续作为历史证据，本轮未重复其长 lane。周复盘没有真实手工发布
指标，保持计划内 `not_ready`；候选包不等于发布。
