# 基础能力需求与绑定 (capabilities/)

本目录声明本 Scene Pack 执行时所依赖的底层能力 (Capabilities)、模型分析要求及技能插槽。真实能力调用由 Base 和 Gateway 统一裁定。

## 📋 目录原则

1. **静态能力清单**：Pack 必须明确列出其所需要的一切外部能力（如大模型调用、文档解析、CAD 适配器等），在平台静态审查时全部公开。
2. **拒绝盲盒行为**：严禁未经声明隐藏危险 API，不允许在运行时临时请求未经 manifest 授权的高风险资源。
3. **Skill Bindings**：声明流程节点如何与底层 Capability Management 对接，并将动作映射到对应的端口。
4. **必须声明**：
   - 能力 ID 和用途。
   - 是否必需。
   - 所属网关：Model、Communication、Execution、Memory 或 Capability。
   - 默认模式：`disabled`、`mock` 或 `candidate_only`。
   - 是否需要 Gate、ReceiptCandidate 和 OwnerDecision。

## 📝 示例能力需求: `requirements.json`

```json
{
  "dependencies": [
    {
      "capability_id": "truzhen.model.llm_analysis",
      "gateway": "model",
      "purpose": "分析施工延期风险并生成候选说明",
      "default_mode": "candidate_only",
      "required": true,
      "requires_gate": true
    },
    {
      "capability_id": "truzhen.gateway.communication",
      "gateway": "communication",
      "purpose": "向业主发送水电验收报告和确认选项",
      "default_mode": "candidate_only",
      "required": true,
      "requires_gate": true
    }
  ]
}
```
