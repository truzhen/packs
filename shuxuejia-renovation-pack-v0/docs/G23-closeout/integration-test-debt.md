# G23 集成测试债：R23C governed uninstall

状态：已解决；R21B Packs 集成提交 `38dd98a` 已更新仓根共享测试。

## 需要的共享测试变更

`test_pack_issued_binding.py` 的墅学家子用例仍模拟旧的 `disable` 路径和脚本自行
Prepare → Confirm。应改为：向 `TRUZHEN_PACK_UNINSTALL_PROOF_JSON` 注入测试签发 proof，
其中 `action_type` 必须为 `14.pack-studio.lifecycle.uninstall`、`target_ref` 为本 Pack、
`transaction_ref` 为 `transaction://pack-uninstall:<pack_ref>@<version>`，并断言脚本仅
POST `/v3/pack-studio/lifecycle/uninstall`，原样转发 decision/run/nonce/evidence；缺 proof
或任一绑定错误必须在网络调用前退出。不得允许脚本自行 Prepare、Confirm 或调用 disable。

`pack_forbidden_artifacts_test.go` 的 OS-14 扫描应接受并要求 canonical
`14.pack-studio.lifecycle.uninstall` 常量，同时禁止该 governed-uninstall 脚本出现
`/v3/base/gated-actions/prepare`、`/v3/base/gated-actions/confirm` 与
`/v3/pack-studio/lifecycle/disable` 调用；现有仅接受旧 disable 字符串的规则已过期。

## 本 Pack 已有证明

`uninstall.py` 只消费 external formal uninstall proof，且
`python3 shuxuejia-renovation-pack-v0/test_uninstall_proof.py -v` 覆盖正向、缺失、畸形、
错 action、错 target、错 transaction，以及脚本中 legacy disable/prepare/confirm/self-mint
路径均不存在，共七项。该测试不连接 devserver、不产生外部副作用。

当前复现命令为 `GOWORK=off go test ./...`，唯一失败日志 SHA256 为
`dedb85cd8089ce2ca3560770c52903cb8543c6e7fdbae5fcfff786d71ba51b33`。

## 解决证据

R21B 将家政与墅学家登记为 formal uninstall Pack：共享 scanner 同时要求 exact
uninstall action/endpoint 并拒绝 disable 双语义；issued-binding 测试注入外部签发 proof，
断言只调用 formal uninstall 且原样转发治理绑定。递归 JSON、Python compile/discovery
28/28、Pack 结构、禁品、`GOWORK=off go test ./...` 与 `git diff --check` 全绿。
