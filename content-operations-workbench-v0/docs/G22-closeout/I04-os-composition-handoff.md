# I04 移交：runtime-disable 未阻断执行侧车自动启动

来源：G22 隔离 OS 验收（固定 SHA `3a97db377d1cbbf171b18ff5a98fcaf4c2e417c1`）。

## 观察

以 `TRUZHEN_ANDROID_AUTO_RUNTIME=0` 和 `TRUZHEN_ANDROID_AUTO_SIDECAR=0` 启动固定 SHA devserver 后，日志仍显示 Windows VM 自动探测路径启动 execution sidecar，并从本机持久 keyfile 加载运行时密钥。该 lane 随后只执行了 OpenMontage 的 read-only preflight、Gate 拒绝和素材权利拒绝；没有执行 Codex、OpenMontage render、登录、发送、上传或发布。停止 devserver 后，任务端口与该 loopback sidecar 均不再监听。

## 风险与期望修复

- 风险：调用方已明确关闭 Android 自动运行时，组合根仍可启动执行侧车并接触持久密钥；这会破坏“未获逐次 Gate 不启动真实执行面”的测试隔离预期。
- 归属：OS 组合根 / execution runtime，不属于 Pack。
- 期望：提供统一、显式且 fail-closed 的总禁用开关，覆盖 Android、Windows VM 自动探测和所有 execution sidecar 自启动路径；禁用时不得读取持久密钥、不得监听 loopback、不得声称 provider ready。
- 验收：在干净临时数据根启动 devserver，设置总禁用后确认无 execution sidecar 子进程、无 keyfile 加载、OpenMontage/Codex 仅返回 `provider_missing` 或 `not_ready`，并保留现有逐次 Gate / Receipt 测试。

G22 不修改 OS；在 I04 给出修复与复核证据前，不应重启本 G22 正向执行 lane。
