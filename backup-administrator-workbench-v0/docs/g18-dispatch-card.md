# G18 备份管理员 Pack 商品化与影子恢复派活卡

- **版本 / 优先级**：Truzhen v4 发布前当前主线，G18；生命周期起点为 `已实现`，本任务目标是以隔离证据推进至 `已验收`，不代表已发布。
- **真实客户 / 场景证据**：Owner 的发布前目标明确要求“备份管理员 Pack 商品化与影子恢复”；本仓 `FEATURE_LEDGER.md` 记录该 Pack 声明层已实现而隔离 install E2E 未补。没有客户原话或生产使用记录，标记为**缺证据**；本轮仅验证最小受控闭环，不扩张产品范围。
- **最小可交付**：在 `127.0.0.1:18218` 的一次性 OS test-store 中证明 Pack 生命周期与既有 restic 受控备份、一次性影子恢复、双 Receipt 可反查；补齐本 Pack 的可复核 closeout。砍掉云市场上架、生产 bucket / repo、定时调度、外部 Provider 新实现及任何跨仓改动。
- **真相源**：Pack 仅为流程、角色和 ProviderRequirement 声明；项目、Gate、正式 Receipt 与恢复结果的事实归 `truzhenos`；restic Provider 实现归 OS / software；商品、License 与 Entitlement 归 cloud。
- **仓库 / 层归属**：只修改 `truzhen-packs/backup-administrator-workbench-v0/**` 的 Pack 验收资产；对 `/Users/li/Documents/truzhenos` 仅固定 SHA 只读并运行隔离 devserver，不修改任何 OS 文件。
- **风险与契约**：橙色为既有 OS lifecycle / ProviderRequirement 语义的兼容核验；红色为隔离 restic backup / shadow restore。无 schema、DTO、Gateway 或主权链变更；backup 与 restore 均必须分别取得既有 Base Gate 和正式 Receipt。
- **允许上下文**：本 Pack 全部资产、本仓治理与相关 Python 测试；OS 根治理、模块 19、devserver 的既有端点及其测试，仅限读取 / 运行。测试数据只在临时状态目录，完成后清理。
- **禁止边界**：不得改 Pack 之外文件、OS / cloud / software / client 代码、公共工具或 contracts；不得写生产备份目标、真实外部 bucket、正式 restic repo、在线源目录、密码或凭据；不得 push、merge、tag、deploy。
- **验收设计**：先验证 clean registry `0 → 1`、重复安装幂等、enabled version、角色、SlotBinding 与 ProviderRequirement；再对隔离项目执行资产盘点 / RPO-RTO 候选 / 恢复审计质询 / Owner 决策卡、受控 backup、影子 restore、摘要校验和两条 Receipt 反查；再验证 provider missing、错误密码、网络中断、坏 snapshot、非空目标、Gate deny、Receipt fail、restart，最后 disable → rollback → uninstall → reinstall 并证明历史 Receipt 保留。
- **用户验收**：复核运行命令、OS SHA、端口、隔离项目、Gate / Receipt 引用及清理结果；确认源码在线数据未变，且 provider 缺失诚实为 `not_ready`。
- **变更影响 / 监控**：仅新增 Pack 内派活与 closeout 证据；采用既有 `pack_diagnostics` 与 OS Receipt / 状态端点，不新增平行日志、监控或诊断格式。

## 执行状态

`设计中`：等待隔离 test-store 启动与既有端点实测。Owner Gate 仅用于真实生产资源；本卡授权范围内的本地一次性 test-store Gate 由既有 devserver 受控签发，不能外推为生产授权。
