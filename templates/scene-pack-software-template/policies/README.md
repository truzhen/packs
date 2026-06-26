# 核心治理策略 (policies/)

本目录包含本 Scene Pack 内部声明 of 门禁限制、回执保留、防抖回滚以及防风险授权策略。

## 📋 核心策略类型

1. **Gate Policy (`gate-policy.json`)**：
   - 声明哪些状态转移、对外发送、以及文件修改动作需要进入 `Base Gate` 并请求 `OwnerDecision` 终裁。
   - 预算变更或外部执行默认强制要求确认。
2. **Receipt Policy (`receipt-policy.json`)**：
   - 声明节点执行的录像和签字要求。
   - 遵循 `append-only` 记账原则，任何已经入账的回执绝对无法被删除、覆写或隐匿。
3. **Permission Policy (`permission-policy.json`)**：
   - 声明模型、沟通、执行、记忆、文件和网络能力边界。
4. **Rollback & Rollforward Policy (`rollback-policy.json`)**：
   - 声明失败、撤销、void、supersede 和 correction 的证据要求。
   - 定义在流程发生意外阻塞时，各个事务节点如何安全退回及补偿。

## 🚫 禁止

- 用 policy 直接授权真实发送或真实执行。
- 在 policy 中写 raw secret。
- 把 mock 结果伪装成正式回执。

## 📝 示例 Gate 声明: `gate-policy.json`

```json
{
  "rules": [
    {
      "action": "business_object.write",
      "target_field": "budget",
      "requires_owner_confirmation": true,
      "audit_level": "critical"
    },
    {
      "action": "gateway.communication.send",
      "requires_owner_confirmation": true,
      "gate_id": "communication_send_gate"
    }
  ]
}
```
