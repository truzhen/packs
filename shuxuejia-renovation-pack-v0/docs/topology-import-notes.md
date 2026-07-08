# 墅学家拓扑迁移说明

- 来源：`/Users/li/Documents/truzhenv3` 本地提交 `3c161aff4e2cc9b6ee896d0ef1ec4b37aaf4b062`。
- 原始资产：`backend/internal/devserver/pack_asset_seed_shuxuejia_topology.md`。
- 迁移结果：`flows/shuxuejia-epc-executable-projection.flow.json` 保留 457 个节点和 543 条显式边。
- 节点类型映射是保守声明：除起点、终点和高风险确认节点外，均保持候选型流程节点；正式任务、沟通、付款、合同、验收和售后关闭仍必须经 Owner + Base Gate。
- 原文标注的 `needs_review` 继续保留：跨段长回流、YES / NO 标注和复杂施工段总线仍需设计师在前端制作台中做人工 parity 校正。

## 可执行投影调整

历史拓扑是 EPC 参考图，不是可直接进入 Scene Runtime 的有向无环执行图。迁移为可执行投影时做了以下最小调整：

- 为缺少角色槽的协作节点补齐 `slot_ref`，统一落到 `design_lead`，避免运行时无法生成角色候选。
- 保留节点数 457 和边数 543，不删除历史节点。
- 调整 6 条形成回流环或孤点的边，使流程可以被基座 SceneFlowRun 完整拓扑排序：
  - `shuxuejia-edge-227`: `S4_015 -> S4_019`
  - `shuxuejia-edge-244`: `S4_024 -> S4_035`
  - `shuxuejia-edge-272`: `S4_046 -> S4_086`
  - `shuxuejia-edge-275`: `S4_047 -> S4_067`
  - `shuxuejia-edge-501`: `S7_025 -> S7_036`
  - `shuxuejia-edge-537`: `S8_008 -> S8_010`
- 调整后校验结果：457 节点、543 边、边引用完整、无重复节点、无孤点、无缺失协作角色槽、无环。

这些调整只改变可执行投影的运行顺序，不改变 Pack 主权边界：所有输出仍为候选，正式化必须经过 Owner + Base Gate。
