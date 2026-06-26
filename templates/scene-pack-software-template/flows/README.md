# Scene Flow 定义层 (flows/)

本目录包含本 Scene Pack 的场景流程 (Scene Flow) 规格设计图，它们是指导 Agent 推进剧情的 DAG（有向无环图）剧本。流程只能生成候选，不能直接执行副作用。

## 📋 目录原则

1. **版本不可变性 (Immutable Spec)**：流程文件发布后具有不可变版本（如 `main.flow.json@1.0.0`）。正在运行的流程实例必须强绑定到某个特定发布版，不得中途受 Pack 升级影响。
2. **纯声明式无副作用**：DAG 节点仅声明动作意图（Candidate）、所需的 Gateway、门禁策略，不得直接编写会导致外部系统发生改变 of 任何数据。
3. **DAG 校验门禁**：必须是 DAG，并包含明确的 `start` 和 `end` 节点。
4. **Flow 规则**：
   - Canvas 导入可由 06 Scene Flow 转换器兼容已知别名，但最终规格必须通过 06 验证器。
   - 任务、沟通、执行、记忆、能力节点只能生成对应 Candidate。
   - 高风险节点必须声明 `gate_policy`。
   - 危险权限必须进入 GateCandidate 路径，不能直接进入执行节点。

## 📝 示例流程定义: `main.flow.json`

```json
{
  "flow_id": "hydro_acceptance_flow",
  "title": "水电验收流程",
  "version": "1.0.0",
  "nodes": [
    { "id": "n1", "type": "start", "title": "水电施工完成" },
    {
      "id": "n2",
      "type": "human_approval",
      "title": "现场水电验收确认",
      "gateway": "communication_gateway",
      "gate_policy": { "required_gate": "base_gate", "pending_owner_confirmation": true }
    },
    { "id": "n3", "type": "end", "title": "项目进入下一阶段" }
  ],
  "edges": [
    { "id": "e1", "source": "n1", "target": "n2" },
    { "id": "e2", "source": "n2", "target": "n3" }
  ]
}
```
