# G23 closeout

状态：`blocked`；生命周期：`Pack owned uninstall 已验收；等待共享扫描器集成`。

旧 lifecycle blocked 已由固定 OS `af45a07b69c2ea44ed4b9d38612468bde0b2050d` 的运行态复验取代。唯一恢复的任务专属 store 完成 Pack lifecycle、457/543 parity、重启/Receipt 回放、R23D retained-vs-created 与 recovery 定向证明、去敏 457 长链，以及所有既定负向矩阵。

五类代表投影（合同变更、材料延期、材料替代、验收失败整改、售后重开）全部仅为 candidate-only/non-formal：首次 `5×201`，稳定键重放 `5×200 deduplicated`，Receipt `5/5` 反查；重启后仍为 candidate，FormalTask 与 finance snapshot 均为 0。真实合同、付款/退款、Provider、发送/执行、通知、验收/售后关闭和外发均为 0。

本次撤销了为迎合旧扫描器而存在的 legacy `disable` handoff。`uninstall.py` 现在只接受 exact `14.pack-studio.lifecycle.uninstall` proof，并只 POST 正式 `/v3/pack-studio/lifecycle/uninstall`；本地 7/7 proof 测试同时拒绝 legacy action、prepare、confirm 与 self-mint。

当前唯一 blocker 是公共 `pack_forbidden_artifacts_test.go`：它仍硬编码要求 uninstall.py 含旧 disable action 字面量，故 `GOWORK=off go test ./...` exit 1（SHA256 `dedb85cd8089ce2ca3560770c52903cb8543c6e7fdbae5fcfff786d71ba51b33`）。本 Pack 不会放回双语义死代码；精确共享测试变更请求见 `integration-test-debt.md`。

完整命令、退出码、耗时及哈希摘要见 `goal-result.json` 与 `evidence-envelope.json`；本机临时运行态证据未写入 Git。

push/merge/tag/deploy: 未执行，等待 Owner。
