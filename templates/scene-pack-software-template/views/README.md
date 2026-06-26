# 声明式视图配置 (views/)

本目录包含本 Scene Pack 的声明式 UI 视图定义，用于给 Truzhen 平台的前端 AppShell 提供独特且个性的界面结构展示，但必须在 Truzhen 标准渲染规范中运行。只能映射到 Truzhen 标准 UI slots。

## 📋 目录原则

1. **安全性隔离**：禁止 Pack UI 包含直接写 DB、直接连外网、直接读文件或修改全局 AppShell 的代码。
2. **只读 ReadModel 数据流**：前端视图只能映射到 standard UI slots 并读取 ReadModel 投影，一切修改意图必须作为 `GatedAction` 发送给后端，由 Base Core 门禁审核后触发。
3. **视觉隔离**：卡片与表单的“候选态 (Candidate)”与“正式态 (Formal)”必须使用明显的边框、底色及标记做强视觉区隔。

## 🛡️ 可声明的 UI Surface Slots

- `pack.dashboard`：Pack 总览、风险、候选和回执摘要。
- `pack.object.detail`：业务对象详情页局部 surface。
- `candidate.card`：候选卡展示，不代表正式完成。
- `transaction.capsule`：事务胶囊投影，只读展示当前对象、任务、风险和回执。

## 📝 示例声明视图: `dashboard.ui.json`

```json
{
  "slot": "pack.dashboard",
  "layout": "grid",
  "components": [
    {
      "id": "progress_chart",
      "type": "Timeline",
      "title": "项目施工回执时间线",
      "data_source": "/api/readmodel/renovation/receipts"
    },
    {
      "id": "task_kanban",
      "type": "Kanban",
      "title": "任务看板 (Gated)",
      "data_source": "/api/readmodel/renovation/tasks",
      "action_trigger": "GatedAction::TaskTransition"
    }
  ]
}
```
