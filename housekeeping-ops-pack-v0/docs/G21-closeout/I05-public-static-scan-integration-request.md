# G21 → I05：公共静态扫描与 formal uninstall 契约漂移

状态：已解决；R21B 集成提交 `38dd98a` 已修复共享扫描与 issued-binding 测试。

## 事实

固定 OS SHA b843186cd2a7e93682f5d67e62d79a228376e368 提供 canonical
14.pack-studio.lifecycle.uninstall 和 Pack Studio formal uninstall endpoint。G21 已用
OS 签发的测试证明完成 formal uninstall、重启、同版本 governed reinstall，并保留
Receipt、角色包与槽绑定。

仓根公共测试 pack_forbidden_artifacts_test.go 的
TestPackGlueDoesNotMintOwnerActionEvidence 仍要求每个 uninstall.py 含旧
14.pack-studio.lifecycle.disable 字符串。G21 的正式脚本只含 uninstall action 与
uninstall endpoint；Pack 本地测试明确禁止 disable action/endpoint。

## 请求

请公共 Pack 治理 Owner 更新该公共静态契约，使其按 Pack 的卸载语义接受 canonical
formal uninstall，且不允许用无效字符串、死代码或注释满足扫描。该文件不在
housekeeping-ops-pack-v0 owned path，G21 未修改。

## 已撤销尝试

曾短暂加入的 LEGACY_DISABLE_ACTION_COMPAT 兼容标记已删除，并将从未推送的本地提交
历史中移除；它从未被提交到 OS 请求路径，也未触发任何 disable 或外部动作。

## 集成解决

共享测试现在按明确 Pack 清单区分 formal uninstall 与 legacy disable：家政与墅学家必须
同时包含 canonical uninstall action 和 endpoint，并且出现任意 disable action/endpoint
都会失败；其余 Pack 保持既有语义。`test_pack_issued_binding.py` 同时验证 external proof
原样转发、禁止脚本自行 Prepare/Confirm。递归 JSON、Python compile/discovery 28/28、
Pack 结构、禁品、`GOWORK=off go test ./...` 与 `git diff --check` 全绿。
