# I04 交接：G23 shared scanner contract drift

状态：`blocked`；本 Pack formal uninstall 语义：`已验收`。

固定 OS `af45a07b69c2ea44ed4b9d38612468bde0b2050d` 已取代旧 SHA 与编排型 lifecycle blocked。G23 在唯一恢复的隔离 store 中完成 governed lifecycle、457/543 flow、重启/Receipt replay、去敏代表投影和运行态负向矩阵。

OS archive 的额外全仓 Go 尝试并非 OS 主仓版本化 pre-push/EGR 环境，不作为本次 blocker。OS 的权威全量门已在 R23D 集成完成；本目标复用其固定 SHA，并完成定向 normal/race、负向矩阵和隔离 E2E。

Packs worktree 的 `GOWORK=off go test ./...` 当前唯一失败是公共 `pack_forbidden_artifacts_test.go` 要求 legacy disable action 字面量。G23 已删除该死标记，补齐 `UNINSTALL_CONNECTIVITY`，并用本地 7/7 proof 测试锁定 exact uninstall、正式 endpoint、拒绝 disable/prepare/confirm/self-mint。请协调线程按 `integration-test-debt.md` 更新共享扫描器后复验；不得由 Pack 恢复双语义标记。

本次没有真实合同变更、付款、供应商或客户通知、正式验收/售后关闭、外发、OpenShip 或真实 Provider 动作；所有去敏场景均停在 candidate-only/non-formal。

push/merge/tag/deploy: 未执行，等待 Owner。
