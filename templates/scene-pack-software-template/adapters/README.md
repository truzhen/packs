# 适配器定义层 (adapters/)

本目录用于定义本 Scene Pack 内部使用的适配器 (Adapters) 和转换代码。这些适配器提供外部业务接口的接入规格与安全模拟层。默认禁用真实发送和真实执行。

## 📋 目录原则

1. **Mock 优先与安全防线**：此目录下的适配器仅在测试和确权时执行。Truzhen 运行态默认禁止 Pack 直接访问外部 API。
2. **严禁凭证硬编码**：禁止携带任何真实的账号 Token、API Key，一切鉴权凭据及网络调用应声明给平台的 Communication / Execution Gateway 并由其统一管理。
3. **接口契约描述 (`interface-manifest.json`)**：明确向平台表明该适配器的数据输入、输出格式与转换函数（Transformers）。
4. **可放置内容**：
   - `interface-manifest.json`：声明输入、输出、候选类型和需要的 Gateway。
   - mock adapter：只用于本地仿真和测试，不访问外网、不读宿主敏感文件。
   - transform：把作者工程数据转换为平台可识别的 Candidate payload。

## 🚫 禁止内容

- 真实 API key、账号 token、数据库连接串。
- 真实发送、真实执行、宿主全盘访问或云同步代码。
- 绕过 Base Gate、Communication Gateway、Execution Gateway 或 Model Gateway 的直接调用。

## 📝 示例适配接口定义: `interface-manifest.json`

```json
{
  "adapter_name": "renovation_mock_adapter",
  "mode": "mock",
  "interfaces": [
    {
      "name": "emit_acceptance_report_candidate",
      "input_schema": { "project_id": "string", "report_ref": "string" },
      "output_candidate": "CommunicationDraftCandidate",
      "security": { "requires_gate": "base_gate" }
    }
  ]
}
```
