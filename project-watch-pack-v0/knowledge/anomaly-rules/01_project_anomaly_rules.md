---
title: 项目异常三条规则与分级说明
verification_status: pending_human_review
source_ref: source://project-watch/01_project_anomaly_rules
---

# 项目异常三条规则与分级说明

本文只解释语义与口径，供角色引用和 Owner 复核。**判定实现不在本 Pack**：三条规则的计算、分级、缺字段处理与快照新鲜度全部由 truzhenos 05 的项目异常清单只读投影（`GET /v3/business-object/project-anomalies`）执行；Pack 与本知识条目都不写条件分支，也不构成第二真相源。真相源是 ERPNext 外部系统，05 只持快照与投影。

## 一、三条规则（第一版，纯确定性）

| 规则 | 语义 | 输入字段（ERPNext 标准 doctype） | 触发口径 |
| --- | --- | --- | --- |
| R1 延期 `delayed` | 项目预计完工日已过而项目仍未结束 | `Project.expected_end_date`、`Project.status` | 项目状态不属于已完成 / 已取消，且 `as_of` 日期减去预计完工日超过「延期天数阈值」 |
| R2 停滞 `stalled` | 项目未完成且长期没有任何更新动作 | `Project.modified`、`Project.percent_complete`、同项目 `Task.modified` | 完成度未达 100%，且 `as_of` 距离项目与其关联任务的最近一次更新时间超过「停滞天数阈值」 |
| R3 付款节点逾期 `payment_overdue` | 项目关联发票到期未收 | `Sales Invoice.project`、`due_date`、`outstanding_amount`、`status`（排除已取消 / 退回） | 未收金额大于 0，且 `as_of` 超过到期日达到「付款逾期天数阈值」 |

## 二、分级说明

- 两级：`warning`（关注）与 `critical`（严重）。
- R1 / R2：超期天数达到对应阈值的 2 倍时升为 `critical`，否则 `warning`。
- R3：超期天数达到「付款严重天数阈值」时升为 `critical`，否则 `warning`。
- 分级只表示「需要 Owner 更早看到」，不表示任何已裁定结论，也不自动触发任何执行。

## 三、默认阈值是可改的默认值，不是行业结论

以下数值是第一版的**默认值**（协调线 2026-09-02 裁定 D3），随策略对象可改，试点第一个月按真实数据回调；本 Pack 只以 `policy_ref: project_watch_policy://default` 引用，不复制数值到 flow 或角色提示词：

| 阈值项 | 默认值 | 含义 |
| --- | --- | --- |
| 延期天数 | 3 天 | 超过预计完工日多少天算延期 |
| 停滞天数 | 7 天 | 多久没有任何更新算停滞 |
| 付款逾期天数 | 0 天 | 到期未收即算逾期 |
| 付款严重天数 | 30 天 | 逾期多少天升为 `critical` |
| 快照最大陈旧小时数 | 24 小时 | 超过即标 `stale`，不静默用旧快照 |
| 负责人字段键 | `owner` | ERPNext 标准元字段，可配；缺失时该项目标 `not_ready`，不静默归入 `unassigned` |

阈值的每次修改由 05 的策略端点持久化，并要求记录旧值到新值与回执；阈值是配置不是业务动作，但改动必须可反查。

## 四、诚实边界（角色引用时必须遵守）

- 缺关键字段时该规则对该项目输出 `skipped` 并注明缺哪个字段，不猜、不用默认值补齐、不当作「无异常」。
- 零项目快照时读模型返回 `not_ready`，不得用空清单表述为「所有项目正常」。
- 快照过期（`stale: true`）时必须在候选与通报里显式标注新鲜度，不得当作实时状态。
- 每条异常都带证据（外部链接引用、快照哈希、快照版本、命中字段原值）；无证据的判断不得进入候选。
- 本知识域整体 `verification_status: pending_human_review`：规则语义与默认阈值来自设计裁定与通用工程常识，尚无真实客户样本核验，不得据此宣称语义完整或直接对外承诺。
