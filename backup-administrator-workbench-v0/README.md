# 备份管理员工作台 Pack（backup-administrator-workbench-v0）

数据安全与备份治理场景荚（Domain Work Pack）。把「什么值得备份 / RPO-RTO 策略 / 受控备份执行 / 影子恢复演练」编排成受主权闸约束的治理流程：**声明策略，不持主权**——正式裁定由 Owner + Base Gate 完成，真实执行经 truzhenos 执行网关驱动 restic 引擎，回执入 03 可反查。

## 流程（flows/backup-lifecycle.flow.json）

盘点数据资产 → 备份管理员出策略方案（advice）→ 恢复审计员对抗质询（challenge）→ Owner 对照确认门 → 备份计划主权确认 → Owner 五要素确认 → **受控执行加密备份**（restic 经 Base Gate，产 `recovery.backup.created` 回执）→ Owner 确认演练 → **影子恢复演练**（一次性隔离目标 + 内容摘要校验，绝不覆盖在线数据，产 `recovery.restore.verified` 回执）→ 复盘归档。

## 角色（全员 Proposer）

- **备份管理员**（advice）：出备份范围/频率/保留/RPO-RTO 候选，绝不自行触发真实备份。
- **恢复审计员**（challenge）：从「真出事那天能不能恢复」质询范围遗漏、可恢复性、演练缺失。

## 真实系统绑定

`software_requirements`: `restic-family >=0.19.0,<1.0.0`（BSD-2，登记于 truzhen-software `providers/restic-backup-engine/source-lock.json`）；缺引擎 fail-closed `not_ready`，不假成功。执行链 = Registry(os-02) → Execution Gateway(os-11/os-19 backupservice) → Base Gate(os-01) → Receipt(os-03)。

## 安装 / 卸载

```sh
# 先起 devserver（默认 127.0.0.1:18080）
python3 backup-administrator-workbench-v0/install.py
python3 backup-administrator-workbench-v0/uninstall.py
```

## 边界

- Pack 内无 runtime / binary / secret；仓库密码归 SecureStore `credential_ref`，Pack 不可见。
- 恢复只进一次性影子目标；覆盖在线数据属禁止边界。
- 生命周期：`已实现`（声明层）；产品级验收以 truzhenos 门禁 + 03 回执为准。
