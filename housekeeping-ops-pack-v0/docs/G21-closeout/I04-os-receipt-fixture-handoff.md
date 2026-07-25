# G21 → I04：固定 OS SHA 家政 E2E 的 Receipt 密钥夹具断链

状态：`已取代（不再阻断 G21）`；不修改 OS。当前权威 SHA 的
`test(receipt): align devserver ledger key provider` 已修复本记录描述的夹具不一致，
G21 已在该 SHA 复验通过。本文件保留旧根因，供审计追溯而非新的 I04 工作项。

## 可复现事实

- 旧固定 OS：`0c218143de84ad60df9c5d78fdebd56cb2cb27da`（K03R goal-result 记录的验收 HEAD）。
- Pack：`scene_pack://housekeeping-ops@0.1.0`，去敏测试，未发送报价、未派工、未上门执行、未联系客户。
- 命令：`GOWORK=off go test -count=1 -run 'TestHousekeepingOpsPackEndToEnd|TestHousekeepingOpsPackOwnerRejectBranch|TestHousekeepingPackManifestDeclaresSixSovereigntyThings' ./backend/tests/devserver`。
- 结果：家政 Run 已走到 completed，候选、Owner Gate、下游 communication/execution/receipt candidate 与历史回放均到达；最终 `VerifyChain` fail-closed，首个根因是 `secure payload inaccessible: secure store get: key "platform-master-…" not found`。

## 根因与影响

`backend/tests/devserver/business_object_routes_test.go` 的 `seedLedgerReceipt` / `openCentralLedger` 用固定内存 `test-key-v1` 打开 receipt ledger；同一测试中的 `devserver.NewHandler` 使用 `devhttpx.NewPlatformModuleStoreFactory` 的文件 KeyProvider。devserver 追加的加密 receipt 因而不能由测试末尾的内存密钥解开。设置新的 `TRUZHEN_DATA_DIR` 与 `TRUZHEN_KEYS_DIR` 后仍复现，说明不是遗留运行态污染。

旧 SHA 因而不能提供“03 Receipt 哈希链可反查”的最终铁证；Pack 不能安全修复这个跨模块测试夹具。后续 OS 已将 devserver ledger key provider 对齐；当前权威 `751473a5dd1b0ee2965397b541ef57a72e8ce273` 的同一套定向测试已通过，故本项不再阻断。
