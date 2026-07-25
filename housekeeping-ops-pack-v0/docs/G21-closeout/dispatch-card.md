# G21 派活卡：家政日常运营 Pack 商品化与派工/报价隔离

- **版本 / 生命周期**：Truzhen v4 发布前目标 G21；已在当前权威 OS SHA 的隔离环境完成`已验收`，不代表发布或上架。
- **真实客户 / 场景证据**：本目标委派指定“去敏服务项目”的咨询、方案、质检、排期报价、派工候选、执行候选、服务回执与历史归档路径。未提供客户原话、真实订单或现场观察，故行业口径与报价数值均为**缺证据**；只验证声明、隔离、幂等与 Gate，不写入行业规则或生产数据。
- **最小可交付**：使用一个合成、去敏服务项目，在固定 OS SHA 的隔离 devserver 完成 install/enable → 双角色与对照门 → 报价、派工、通知、执行各自 Candidate/Gate/Receipt 断言 → restart/失败/幂等 → disable/rollback/uninstall/reinstall；提交紧凑 evidence envelope、closeout 与必要 Pack 文档修正。砍掉真实外发、真实排期派工、上门执行、客户联系、Provider 实现、云市场上架和任何跨仓改动。

## 真相源、归属与边界

| 项目 | 裁定 |
| --- | --- |
| 正式客户、订单、排期、派工与回执事实 | `truzhenos` 05/07/03 与 Gateway；Pack 不保存、不铸造 |
| Pack 声明、角色、能力需求、流程与安装胶水 | `truzhen-packs/housekeeping-ops-pack-v0` |
| OS 使用范围 | 仅只读 K03R 固定 SHA `0c218143de84ad60df9c5d78fdebd56cb2cb27da` 的临时副本，启动本机 `127.0.0.1:18221` 与临时 test-store；不修改 OS、contracts、provider 或公共工具 |
| 禁止边界 | 不连接 OpenShip；不使用真实 PII、凭据、客户或生产端点；不改公共工具、其他 Pack、Provider、契约、Gateway 或 UI；不 push/merge/tag/deploy |

## 风险与 Gate

- **黄色**：Pack 生命周期与声明验证；可在授权 worktree 内执行并留证。
- **红色**：报价对外发送、排期派工、上门执行、真实客户联系、正式化与删除。测试仅验证独立 Candidate/Gate/Receipt 或 `not_ready`；真实动作一律暂停等待新的 Owner Gate。
- **契约影响**：目标只验证/澄清既有 Pack 声明；若发现需要修改 contracts、Gateway、ProviderRequirement 语义或主权链，停止 Pack 施工并交 I04。

## 验收与影响清单

| 验收项 | 证据 |
| --- | --- |
| 生命周期与多版本共存 | registry/enabled version、resolver 结论、disable/uninstall/reinstall 的摘要 |
| 三动作隔离 | 报价草稿、派工候选、执行候选分别存在；没有真实 send/dispatch/execute；各自 Gate/Receipt 或诚实 `not_ready` |
| 异常与恢复 | 重复派工/发送幂等、人员不可用、客户改期、Receipt fail、restart 的断言与日志摘要 |
| Pack 静态门 | JSON、Python、forbidden、结构审计、Pack 测试、`GOWORK=off go test ./...`、`git diff --check` |

用户验收方式：复核 `evidence-envelope.json` 与 `goal-result.json` 中的固定 SHA、端口、命令、退出码、哈希、Gate/Receipt 摘要和清理状态；任何真实动作都应显示为未执行并要求 Owner Gate。
