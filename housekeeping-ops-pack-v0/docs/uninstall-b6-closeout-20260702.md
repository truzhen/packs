# 家政运营 Pack 卸载闭环 B6 收口报告（2026-07-02）

## 结论

`housekeeping-ops-pack-v0` 已补齐 `uninstall.py`。脚本经 Base gated-action prepare / confirm 获取真签发 `decision_ref/run_id/nonce` 后调用 Pack Studio lifecycle disable。卸载只停用 Pack 当前版本，不删除历史事务对象、候选或 03 回执。

## 改动范围

- 新增 `housekeeping-ops-pack-v0/uninstall.py`。
- 更新根 `README.md`、`CLAUDE.md`、`MODULES.md`、`FEATURE_LEDGER.md`。
- 更新 `housekeeping-ops-pack-v0/README.md`，补充加载 / 卸载命令与卸载边界。

## 验收证据

- `git diff --check` 通过。
- 全仓 JSON 解析通过，输出 `JSON 合法`。
- `install.py` / `uninstall.py` 语法检查通过，输出 `脚本语法 OK`。
- Pack 结构审计通过，输出 `Pack 结构审计 OK`。
- 隔离本地 devserver：`GOWORK=off go run ./backend/cmd/devserver -addr 127.0.0.1:18109`。
- 家政包重装 / 卸载闭环通过：
  - `after_reinstall_current_version='0.1.0'`
  - `after_second_uninstall_current_version=''`
  - `after_second_uninstall_state='disabled'`

## 仍未完成

- `knowledge/` 仍未补齐；本包仍不能标成完整知识包。
- 本轮未修改 ProviderRequirement、flow、role slot、role pack 或 capabilities 行为。

## 新发现

- `smart-home-owner-pack-v0/uninstall.py` 可读文案仍是环保执法 Pack 复制件。逻辑按 manifest 读取 `pack_ref`，但 `content_summary`、`reason`、说明文本会误导验收。建议另开小清理卡，不在 B6 顺手修改。
