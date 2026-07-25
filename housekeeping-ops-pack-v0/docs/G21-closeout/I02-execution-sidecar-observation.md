# G21 → I02：隔离 devserver 的 execution sidecar 观察

状态：已被 R21A 修复取代；保留 af45 观察为历史，不修改 OS。

## 历史与复验

历史观察：OS SHA af45a07b69c2ea44ed4b9d38612468bde0b2050d 曾在隔离
devserver 启动时自动启动临时副本的 execution_sidecar.py；该运行态已按精确 PID
清理，不能作为本次结论。

独立复验：OS SHA b843186cd2a7e93682f5d67e62d79a228376e368 使用全新任务专属
store/key 目录，在 127.0.0.1:18221 启动，并显式设置：

- TRUZHEN_ANDROID_AUTO_RUNTIME=0
- TRUZHEN_ANDROID_AUTO_SIDECAR=0
- TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0
- TRUZHEN_REAL_E2E=0

devserver 正常监听 18221，但未产生本任务 execution_sidecar.py PID 或监听。
governed local-source install、formal uninstall、devserver restart 与同版本 governed reinstall
后复查仍为 absent。

## G21 影响判断

G21 没有调用 execution sidecar、真实 Provider、真实报价发送、派工、上门执行或客户
联系。b843 的 hermetic E2E 覆盖 ExecutionIntentCandidate、Receipt、重启、打回分支与
Candidate/Formal 隔离；缺 Provider 时仍为 not_ready 或 provider_missing。因此 R21A
已解除本目标的组合根阻断。

R21A 的修复语义为：显式 TRUZHEN_EXECUTION_SIDECAR_AUTOSTART=0 优先于 Windows VM
产品 fallback。该修复归 OS 组合根；G21 只记录对固定 SHA 的消费性验证，不授权 Pack 修改。
