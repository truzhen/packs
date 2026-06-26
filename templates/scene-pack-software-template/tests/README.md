# 单元与仿真测试层 (tests/)

本目录包含保障本 Scene Pack 品质与主权安全的测试用例 (Test Suites)。这些测试可在开发阶段及导出审查时全量执行，以证明工程的完整性。

## 📋 测试套件类型

1. **Pack Smoke Test (`pack-smoke.test.ts`)**：
   - 冒烟测试：静态审查包 manifest、权限清单、确保不包含危险 of 未定义行为，验证所有的 Business Object Schemas 以及 flows 满足语法规范。
2. **Flow Simulation Test (`flow-simulate.test.ts`)**：
   - 仿真测试：在无真实网络、无真实执行的高速 Mock 环境中推进 DAG 节点，模拟状态转移，验证优先级评分与风险规则是否正确生效。
3. **Red-team Safety Audit (`redteam.test.ts`)**：
   - 红队测试：极度严苛的安全门禁防线测试。验证若有恶意篡改的节点企图直接真实发送、真实执行或在 `generated/` 之外读取 raw secret、覆盖只读文件、窃取宿主 Token 时，能否被平台的 `Base Gate` 完美拦截和阻断。

## 🛡️ 必备验证防线 (Static & Dynamic Gates)

- manifest 和 schema 静态审查。
- DAG 校验、cycle rejection、危险权限拦截。
- no-real-send、no-real-exec、no-real-model-call 断言。
- raw secret scan 和明文敏感信息扫描。
- Candidate / Formal 视觉和语义隔离。
- ReceiptCandidate 生成与持久化回放路径。
