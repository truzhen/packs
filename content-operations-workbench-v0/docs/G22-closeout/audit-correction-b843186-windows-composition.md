# G22 审计纠正：b843186 离线 Hands 安装与组合根阻断

日期：2026-07-25。固定 OS：`b843186cd2a7e93682f5d67e62d79a228376e368`。

本轮只读核对了当前 source lock、正式 Provider install 实现和组合根；不修改
`truzhenos` 或 `truzhen-software`，不联网拉取、不登录、不调用 Codex Hands，也不访问
平台。

## 已核对的离线安装边界

- 本机存在并可只读核验
  `truzhen/codex-hands@sha256:08b687b435e7bb00bf803e2e2f099bfb05d9701ae7900f60a7fbde256aae97f5`
  （`linux/arm64`、`10001:10001`）。这不是 Provider ready 或 Runtime 投影。
- 当前 OS `RegistrySource.ResolveMaterialPlan` 将 Docker 物料输出固定为
  `truzhen-software/providers/codex-hands/docker/materials`；source lock 要求的目录当前
  不存在。官方 `SecureMaterialFetcher.FetchAndVerify` 对每个锁定 image 执行
  `docker pull`，随后才可写 material-ready 标记和进入 Owner-gated install。
- 因而没有“任务临时离线物料目录 + 已有成品 image 采纳”的正式路径：临时目录不被当前
  source lock 接受，已有成品 image 也不能替代受锁 archive、基础 image 验真、material
  receipt 或 Provider/Runtime projection。继续需要跨仓修改，G22 无权实施。

## 新的组合根 formal blocker

在 `b843186`，`backend/cmd/devserver/main.go:45` 先调用
`EnsureAndroidExecutionRuntime`，随后 `main.go:47` 调用 `NewHostServer`。后者在
`backend/internal/devserver/server.go:373` 无条件调用
`MaybeLoadWindowsGuestTicketPrivateKey`；该函数在
`executionruntime/windows_vm_oneclick.go:304-316` 读取持久
`ticket_private_key.b64`，发生在 `server.go:378-389` 解析正式
`TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0` 之前。

`TRUZHEN_WINDOWS_VM_AUTODETECT_IN_TESTS=0` 仅影响 Go 测试进程，不能关闭真实 devserver
的常规路径；未发现读取前可生效的正式 product 开关。按安全纠正，本轮已停止动态服务，
未再以 dummy key、替代 key 或伪造变量绕过该调用。

安全纠正后的证据：18222 未监听；没有运行中的本任务 devserver/sidecar；没有发起 Provider
prepare/execute、模型调用或 Hands 请求；没有新增 task-owned store、session、trace、artifact、
Provider install candidate 或 Receipt。先前两次隔离尝试的 task 临时目录均已移入本机废纸篓，
不进入 Git；第一次尝试的持久 key 自动加载风险已透明记录，不能当作零读取。此后的阻断复核
完全为只读源码审计。

## 需要 Owner 协调的跨仓影响（未实施）

1. `truzhenos`：在组合根中提供读取持久 guest-agent key **之前**生效、可审计的真实
   devserver product disable 开关，并保留 `autostart=0` 的 fail-closed 行为。
2. `truzhenos` + `truzhen-software`：正式定义并审核“任务隔离的既有 digest image 离线物料
   采纳”契约，包含 source-lock 输出边界、archive/base-image attestation、Owner Gate、Receipt
   与 Runtime projection；不能以手工 registry/source-lock 改写代替。

在这两项由相应 Owner/仓库完成前，G22 的 Codex Hands 正向 candidate-only 运行保持
`blocked`。OpenMontage 已通过的历史 lane 不在本轮重复。平台动作始终为 0；
`publication_authorized=false`、`publication_status=not_published`。
