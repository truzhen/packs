# 行业软件调用声明 (software/)

本目录包含本 Scene Pack 对宿主环境专业行业软件（如 CAD, Photoshop, Dendrite 聊天客户端, ERP 等）的调用适配声明、安全策略和使用说明。

## 📋 核心治理原则

1. **三类适配定义**：
   - **Mock Adapter**：包含本地不联网、不调用真实进程、只读取测试文件并输出 Mock 结果的测试代码，用于确权测试。
   - **Provider Hint**：告诉 Truzhen 平台运行时它期望使用什么样的本地 Provider。例如：`preferred_gateway: "execution_gateway"`, `provider_type: "windows_compat_sandbox"`, `real_execution_default: false`（默认必须为 false）。
   - **Real Provider Adapter (默认禁用)**：提供给平台用于真正集成该软件时的配置和安全性声明文件（`safety-policy.json`）。实操中，真实动作必须提交到 `Execution Gateway`，由隔离沙盒 (Sandbox)、本地受控进程或人类接管 (Human Takeover) 执行，并输出不可篡改的 `Artifact` 和回执。
2. **接口声明原则**：
   - 本目录只描述外部软件接口、参数、前置条件和安全限制。
   - 真实软件操作必须由 11 Execution Gateway 生成 ExecutionIntentCandidate 后再受控处理。
   - 真实沟通动作必须由 10 Communication Gateway 生成 Draft / Send Candidate 后再受控处理。
   - 默认 provider 必须是 `disabled` 或 `noop`。

## 🚫 禁止

- 直接启动外部应用。
- 直接读写宿主文件。
- 直接连接微信、邮件、短信、CAD、ERP 或云服务。

## 📝 示例 CAD 调用描述: `cad/adapter-manifest.json`

```json
{
  "software_id": "autocad_2026",
  "action_name": "generate_layout_screenshot",
  "safety": {
    "runs_in_sandbox": true,
    "max_cpu_percent": 50,
    "requires_owner_confirmation": true
  }
}
```
