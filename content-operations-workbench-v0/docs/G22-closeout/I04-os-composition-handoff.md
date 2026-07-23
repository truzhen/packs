# I04 移交：已被 G22 当前权威 SHA 审计更正取代

原观察来自 G22 的旧候选 SHA `3a97db377d1cbbf171b18ff5a98fcaf4c2e417c1`；该提交不是当前权威 `main` SHA `751473a5dd1b0ee2965397b541ef57a72e8ce273` 的祖先。因此本文件不再构成当前 G22 的公共根因移交或阻断依据。

## 观察

旧候选 SHA 曾在未隔离的软件注册表环境中观察到 Windows VM 自动探测路径。此观察不能外推为当前 `main` 的缺陷。

按 Owner 审计纠正，G22 在当前权威 SHA 上以空的软件注册表根、`TRUZHEN_ANDROID_AUTO_RUNTIME=0`、`TRUZHEN_ANDROID_AUTO_SIDECAR=0` 与 `TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0` 重跑。结果未发现任务 OS 的 execution sidecar 子进程、持久 keyfile 加载或执行运行记录；OpenMontage 保持 `provider_missing`，Codex 保持 `not_ready`。详见 [audit-correction-751473a5.md](audit-correction-751473a5.md)。

## 历史建议（非当前 G22 blocker）

- 风险：调用方已明确关闭 Android 自动运行时，组合根仍可启动执行侧车并接触持久密钥；这会破坏“未获逐次 Gate 不启动真实执行面”的测试隔离预期。
- 归属：OS 组合根 / execution runtime，不属于 Pack。
- 期望：提供统一、显式且 fail-closed 的总禁用开关，覆盖 Android、Windows VM 自动探测和所有 execution sidecar 自启动路径；禁用时不得读取持久密钥、不得监听 loopback、不得声称 provider ready。
- 验收：在干净临时数据根启动 devserver，设置总禁用后确认无 execution sidecar 子进程、无 keyfile 加载、OpenMontage/Codex 仅返回 `provider_missing` 或 `not_ready`，并保留现有逐次 Gate / Receipt 测试。

G22 不修改 OS。当前 G22 已完成所需的隔离 restart 验证；正向执行 lane 仍仅因缺少逐次 Owner Gate、素材权利确认和 Provider 授权而阻断。
