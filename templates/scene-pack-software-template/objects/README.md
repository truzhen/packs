# Business Object 定义层 (objects/)

放置 Business Object schema，例如客户、项目、材料、证据和任务引用。业务对象是流程中的主角、道具、凭证或状态实例，而非流程本身。

## 📋 目录原则

1. **强类型与 Schema 约束**：每个业务对象均有一个独立的 `.object.json` 描述文件，包含字段定义、数据类型、校验规则及关系描述。
2. **实例独立性与乐观锁机制 (OCC)**：同一个 Schema 可以对应系统中成百上千个独立的业务对象实例。这些实例共享 Schema（图纸），但数据、生命周期及运行状态绝对隔离。任何写入操作必须校验 `state_version` 乐观锁，保证并发安全。
3. **对象规则**：
   - schema 只描述业务对象结构，不写正式业务数据。
   - 示例数据必须脱敏，不能包含真实客户、地址、电话、账号或合同。
   - 运行态对象由 05 Business Object Workbench 管理，Pack 只能提供 schema 或候选初始化建议。
   - 对象状态变化必须通过 Candidate / Gate / Receipt 链路沉淀。

## 📝 示例 Schema: `renovation_project.object.json`

```json
{
  "object_type": "renovation_project",
  "display_name": "装修项目",
  "fields": {
    "project_id": { "type": "string", "required": true, "unique": true },
    "customer_name": { "type": "string", "required": true },
    "budget": { "type": "decimal", "min": 0 },
    "current_state": { "type": "string", "default": "咨询建档" },
    "state_version": { "type": "integer", "default": 1 }
  },
  "relations": {
    "contracts": { "type": "has_many", "target": "contract" },
    "materials": { "type": "has_many", "target": "material" },
    "receipts": { "type": "has_many", "target": "receipt" }
  }
}
```
