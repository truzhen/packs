# 环保执法 Pack · R5 GUI 续测计划（2026-07-05）

> **计划-only**。本文件写完即停，等 Owner 明确「开工」再执行；对话中只给路径 / 结论 / 待裁项。
> **来源**：Owner 2026-07-05「（吸收资产后）做下一轮的环保执法的测试计划」。
> **基线方法论**：沿用 R3 GUI 计划 + R4 续测——全员 Opus + **GUI-only 铁律** + 四画像 + REG 回归 + 深水区 + §6.4 真实性铁律 + **组织者 neuter 独立复核** + 假 HOME 隔离栈 + 专属端口 fresh install。通用纪律不重述。
> **承接**：R4 续测报告 `/Users/li/Documents/过程文档/environmentalpack-round4-20260704/environmentalpack-round4-report.md`（R4 发现 ENV-R4-01/02，本会话已 land；R4 卡在 `screening` 缺 `pdf_parse_status` blocker、真实 PDF 处置闭环未测）。

---

## 0. 为什么是 R5（本轮存在的理由）

R4 是「发现问题」轮：它挖出 ENV-R4-01（知识大域召回被 09 单查上限截断）+ ENV-R4-02（install.py 丢 provider_family），但当时**只出影响清单未修**，且 run 卡在 `screening` 缺 `pdf_parse_status`。

本会话已把 R4 的两个发现 **land 到 origin/main**：
- ENV-R4-01 → truzhenos `37fdd9e`（knowledge readmodel `total_count`/`offset` + ContextCompiler `knowledge_recall_truncated` warning，Owner 裁 09 基线 bump 20660→20697）。
- ENV-R4-02 → packs `d140723`（install.py 透传 `provider_family`）。

∴ **R5 的第一性 = 回归证明 R4 发现的修复在 origin/main 干净现状真生效无回归**（不是重新发现，是收口验证）；第二性 = 打通 / 定性 R4 被 blocker 挡住的 `screening → pdf_parse` 深水。

---

## 1. 版本 / 优先级（待 Owner 裁定）

- 定位：**护城河标杆样板（环保执法）的回归 / 硬化线**，非第一现金流主线（第一现金流=销售 / CRM / 私域 / 客服 / 售后，见战略转向 2026-06-21）。
- 档位判断：**backlog / 标杆样板持续硬化**，不是 V4 新功能，也不是当前客户主线。
- **Owner 需裁**：是否值得现在投这一轮 Opus 全员实测？还是收口到 R4-01/02 已 land 即止、R5 降级为「只跑 REG 回归轻验」不做深水？

## 2. 真实客户 / 场景证据

- **缺真实客户原话**（与 R1–R4 一致）。证据口径 = Owner 派测、仓内环保 pack、当前 GUI 实测。
- 生命周期止于「**已验收（本地隔离）**」，**不宣称客户生产可用**。

## 3. 最小可交付（本轮最小闭环）

**必须进本轮**：
- REG 回归区（§5）全绿：R4-01 / R4-02 + flow-stall + advice node_type + R3 三修，在 origin/main 现状**无回归**，GUI + readmodel 双证 + 组织者 neuter 独立 curl 复核。
- 深水一项 DW-R5-01（§6）：`screening` 缺 `pdf_parse_status` 的**受控性定性**（诚实 blocked vs 死卡缺口）。

**明确砍掉 / backlog**：
- 云 provider 真实接管、真实文书送达、真实执法执行（受控边界，只允许 `provider_missing`/`blocked`/`manual_handoff`）。
- 真实 PDF 处置闭环（upload→parse→screening 解阻）——**需 Owner 提供可用 PDF 样本或明确授权用测试 PDF**，否则本轮止于 blocker 诚实性定性。
- P4-09 引擎越 owner_gate 修复（红线，本轮**只复现出证据不修**，归 Owner 裁）。

## 4. 要做的事（目标 / 范围 / 非目标）

- **目标**：证明 R4 已 land 修复在 origin/main 干净现状生效无回归；打通 / 定性 R4 卡住的 `screening→pdf_parse` 深水。
- **范围**：env pack GUI 全链路（九阶段）+ 知识页 + 03 反查 + 生命周期停用 / 启用回环。
- **非目标**：新增 pack 能力、改契约、真实外部动作、修 P4-09。

## 5. REG 回归区（已 land 修复 → origin/main 现状验证）

> 口径：`run status 权威`（09 readmodel + 06 run + 03 receipt）；**GUI pod cover 静态标签不作准**（R3 教训）。每项双证 + 组织者 neuter 独立复核（防伪造）。

| 编号 | 回归对象 | land 锚点 | 判绿标准 |
| --- | --- | --- | --- |
| REG-R5-01 | ENV-R4-01 知识 total_count / 不截断 | truzhenos `37fdd9e` | `/v3/memory/knowledge/formal` 返回 `total_count`；`environmental/cases` 域 **total_count=486**（非 100 / 366）；分页 `offset` 生效；ContextCompiler 召回命中>fetch 上限时出 `knowledge_recall_truncated` warning。GUI 知识页 + API 双证。 |
| REG-R5-02 | ENV-R4-02 provider_family 透传 | packs `d140723` | install 后 lifecycle draft / readmodel 带 `provider_family`（env pack 自身无 frappe 声明则 no-op，取 smart-home 作对照证透传链路活）。 |
| REG-R5-03 | flow-stall（run-start 不被请求 deadline 截断） | truzhenos `658e495`（以 git log 为准） | 冷起 advice→主权 wait **停稳**不卡 running；请求 ctx 超时不打断 run 编排。 |
| REG-R5-04 | advice node_type 四层透传（不误判 challenge） | contracts/os/packs/client node_type 链（以 git log 为准） | 单角色 advice 节点 `node_type='advice'`（责任含"风险"仍 advice），双角色对照节点正确；GUI 按钮文案与 node_type 一致。 |
| REG-R5-05 | R3 三修现状无回归 | `bc6f2e0`/`dca2179`/`d2b5e58`（以 git log 为准） | 前端 flowId 正确、双角色本地冷起不超时、交互批 6 项 GUI 可操作。 |

## 6. 深水区（R4 未竟）

| 编号 | 深水对象 | 归属 | 判定 |
| --- | --- | --- | --- |
| DW-R5-01 | `screening` 缺 `pdf_parse_status` blocker | **05-business-object-workbench** | 定性：run 停 screening 是**受控 blocked（诚实留痕，可解释）**还是**无 GUI 恢复的死卡（缺口）**。若 Owner 给测试 PDF 样本→验 upload→parse→screening 解阻→下一 stage 闭环；不给样本→止于 blocker 诚实性定性 + 出影响清单。 |
| DW-R5-02 | 九阶段全链路诚实态 | 06/09/10/11 | `monitoring_data_fetch`→`provider_missing`、`legal_doc_draft`/`enforcement_exec`→`blocked`、文书送达→`manual_handoff`，**全程无假成功、无自铸 ref、无假读**。 |

## 7. 四画像（沿用 R3，GUI-only）

一线执法员（主流程）/ 挑剔律师（质检对照角色，验双角色 challenge 真实性）/ 监管审计（03 反查 + evidence 完整性）/ 门外汉 Owner（主权门控是否挡得住误操作）。

## 8. 主权红线复现（只复现不修）

| 编号 | 红线 | 处置 |
| --- | --- | --- |
| RED-R5-01 | P4-09：引擎不消费 GatePolicy → 越 `owner_gate` | **最高优先**。复现验证 origin/main 现状是否仍存在；出复现证据 + 影响清单，**不修**（改主权链=红，归 Owner 裁）。 |

## 9. 真相源 / 归属 / 风险 / 契约 / 上下文 / 禁止边界 / 验收（治理必填）

- **真相源**：run status = 06 run + 09 readmodel + 03 receipt；知识真相 = 09 SQLite `formal_rows`。ReadModel / GUI 静态标签≠真相。
- **仓库 / 层归属**：跨 truzhenos（06/09/05 后端）+ client（12 前端 GUI）+ packs（env pack 事实基线）。screening blocker 归 05。**测试轮默认不改代码**。
- **风险维度**：AI = 施工者（执行测试）+ 组织者 neuter 独立复核（对抗验伪）。高风险动作受控，只允许 `provider_missing`/`blocked`/`manual_handoff`。
- **契约影响**：纯测试轮**不改契约**。深水若发现需改 05 screening/pdf 契约 = **橙，只出影响清单不施工**（改 schema 先查反向依赖 + bump VERSION 铁律 7）。
- **上下文维度**：允许读 = env pack 全目录、06/09/05 相关后端、12 前端 env 页面、R3/R4 报告 + issue 台账。**禁止无边界扫全仓后直接施工**。
- **禁止边界**：① 不碰 codex worktree（`*pack-suite-hardening*` 等他会话 WIP）；② 不真实发送 / 执行 / 写外部系统；③ 不真实上传 PDF（除非 Owner 给样本）；④ 不改主权链（Owner/Base Gate/Receipt）；⑤ 不修 P4-09（红线只复现）。
- **验收设计**：AI 设计严苛验收 + **组织者 neuter 独立 curl 复核**（防伪造，revert 必 FAIL）。每 REG 项双证映射到 land commit；深水项映射到 R4 blocker。
- **验收维度**：改了什么证明什么——REG-* 逐项映射已 land commit，DW-* 映射 R4 未竟，RED-* 只出复现证据。**不得用无关测试替代**。
- **变更影响**：纯测试默认零代码；发现新缺口按颜色分流（绿自修自证 / 黄改+Owner 确认 / 橙影响清单 / 红只方案）。
- **生命周期档位**：`已实现 → 已接线 → 已验收（本地隔离）`；**不进「已发布」**（缺真实客户证据）。

## 10. 隔离栈与验收产物（执行时填）

- 假 HOME + 专属端口 devserver（避让在用 :18095/:18096/:18109）+ 独立 vite + env pack fresh install。
- 产物落 `/Users/li/Documents/过程文档/envpack-live-r5-<date>/`：user-log / behavior-matrix / ui-audit / issues / shots + REG 双证 JSON。
- 报告落同目录 `environmentalpack-round5-report.md`；对话只给路径 + 结论 + 橙 / 红待裁项。

## 11. 待 Owner 裁定（开工前）

1. **§1 优先级**：值得投 R5 全员 Opus 实测，还是降级为「只 REG 轻回归」不做深水？
2. **DW-R5-01 PDF 样本**：是否提供真实 PDF 样本 / 授权用测试 PDF 跑 screening→parse 闭环？不给则止于 blocker 定性。
3. **是否纳入 R3 §9 遗留橙项**（云 provider 真实接管 / 服务接管 / 多案对照）——建议**本轮不纳入**（需 Owner 逐项裁 + 缺客户证据）。

> 本计划写完即停。**未获 Owner「开工」不执行任何实测。**
