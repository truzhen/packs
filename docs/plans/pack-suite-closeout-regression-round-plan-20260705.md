# Pack 套件 · 收口+回归轮 测试任务计划（2026-07-05）

> **计划-only**。本文件写完即停，等 Owner 明确「开工」再执行；对话中只给路径/结论/待裁项。
> **来源**：Owner 2026-07-05「吸收 codex pack-suite-hardening 三份报告资产，做下一轮测试任务计划」+ 问卷三裁定：① 焦点=**收口+回归轮**；② codex hardening=**测 origin/main 现状**（codex 自 land+对账，我不碰其 worktree）；③ 证据=**续多包能力验证口径**（标「缺真实客户原话」）。
> **基线方法论**：沿用第三轮环保 GUI 计划（全员 Opus + GUI-only 铁律 + §6.4 真实性铁律 + 组织者 neuter 独立复核 + 假 HOME 隔离栈 + 专属端口 fresh install）。承接 `multi-pack-test-task-suite-20260704.md`（3 包各已测一轮）。通用纪律不重述。

---

## 0. 版本 / 优先级 / 生命周期档位

- **版本 / 优先级**：平台能力验证**收口轮**（非特定客户主线，非 V4，非 backlog）。收口 = 把 codex 三轮 + 我三轮共同发现的遗留收干净 + 回归验证我已 land 的两修复。
- **生命周期档位**：本轮把「遗留问题」从 `想法/设计中` 推到 `已实现 -> 已接线 -> 已验收（本地隔离）`；把「我已 land 的两修复」从 `已发布`（origin/main）加一道 **GUI 生效 + 无回归** 的独立验收。**不宣称已发布 / 客户生产可用**。
- **证据口径（Owner 裁）**：续多包能力验证。证据 = Owner 派测记录 + 仓内既有 pack + GUI / registry / 03 receipt / monitor + codex/我 prior 轮资产。**全程标「缺真实客户原话」**，不纳入真实客户/凭据/发送/PDF 样本。

---

## 1. 吸收的 codex 资产（pack-suite-hardening 三份报告）

| 报告 | 路径 | 关键结论 |
|---|---|---|
| 测试汇总 | `/Users/li/Documents/过程文档/pack-suite-hardening-20260704/three-round-summary-test-report-20260704.md` | 3 包本地 live 闭环、主权链未绕；遗留 ENV-R4-01/02、monitor warning、缺客户证据 |
| 修复汇总 | `.../three-round-summary-fix-report-20260704.md` | R1/R2 缺口在 hardening worktree 修复（**未提交未 land**）：provider 机器态投影 / 03 反查按钮 / Frappe requirement 展示 / 服务履约文案 / TRUZHEN_PACKS_ROOT |
| 资产汇总 | `.../three-round-summary-asset-report-20260704.md` | 截图/registry/03/run-readmodel/monitor/运行态 SQLite 全登记；关键 ID（tx/run/receipt）表 |

**codex 已修但未 land 的问题（在 `codex/pack-suite-hardening-20260704-{os,client}` worktree，origin/main 上仍存在）**：
- truzhenos：`TRUZHEN_PACKS_ROOT` 覆盖、capability 节点 enrich `provider_family/requirement_id/fallback_policy/provider_status`、smart-home snapshot `provider_missing` 非阻断投影。
- client：新建长等待提示、时间线 03 反查按钮默认可见+聚焦、能力卡 provider/fallback 展示、家政服务履约文案、`ScenePackManagePage` 真实系统绑定展示、canvas 错误与知识错误分离、smoke 断言 `StudioShell`。

**⚠️ 交叉点（land 协调，非本轮施工）**：codex hardening 的 `client/src/pages/BusinessObjectWorkbenchPage.tsx`（candidate 面板 `candidate_previews` 投影）与**我已 land 的 SH-R1-01**（同文件 `roleCandidateType` node_type 权威判定，origin/main `8c9a493`）**同文件重叠**；codex land hardening 时须 rebase 到 `8c9a493` 并合并 candidate 区。**本轮不碰 codex worktree/分支，只在报告中标注该协调点。**

**codex 遗留问题账本（本轮收口对象）**：

| ID | 风险 | 归属 | 说明 |
|---|---|---|---|
| ENV-R4-01 | 橙 | truzhenos / 09 knowledge readmodel | SQLite 752 条环保 FormalKnowledge，API 单次仅返 100/GUI 见 366；`knowledge_scope://environmental/cases` SQLite 486 API 100。需 pagination/cursor/total_count |
| ENV-R4-02 | 黄 | truzhen-packs install.py | pack 声明 `provider_family`，install.py 请求体白名单丢弃该字段（**与我 SH-R1-01 install.py 丢 node_type 同类根因**） |
| Frappe not_ready | 黄 | truzhenos executor/registry | 当前 live 只覆盖 `provider_missing`/`blocked`，未覆盖「provider 已登记但 endpoint/运行体不可用」的专门 `not_ready` 态 |
| MON-HARDENING-01 | 绿/黄 | monitor | hardening final `degraded`，`security=0 error=0 warning=1`，需消化确认非阻断 |

---

## 2. 本轮目标 / 范围 / 非目标

- **目标 A（回归验证，GUI-only）**：确认我已 land 的 **HK-R1-P3-01**（`658e495`+`d200065`）与 **SH-R1-01**（4 仓 `528f2df/83706eb/aec95f6/8c9a493`）在**最新 origin/main** 上 GUI 生效 + 无回归。
- **目标 B（闭合遗留）**：ENV-R4-01（橙，只出 09 API 影响清单）、ENV-R4-02（黄，修 packs install.py 白名单）、Frappe not_ready（黄，构造专门态复测）、MON-HARDENING-01（消化）。
- **非目标**：不建新包（第 4 包另立）；不碰 codex hardening worktree/分支/端口；**不改 contracts / 09 public API/DTO**（ENV-R4-01 只影响评估）；不真实发送/执行/Frappe 写/派工/执法；不重复修 codex 已修未 land 的问题（标「待 codex hardening land」）。

---

## 3. 工作项（每项含 真相源 / 归属 / 风险色 / 契约影响 / 验收）

### A. 回归验证（origin/main GUI-only）

- **A1 · HK-R1-P3-01 回归**（真相源=truzhenos SQLite run readmodel / 03；归属=truzhenos 06；风险=绿验证）：家政（双角色）+ 智能家居（单角色）GUI 起 auto-advice flow run，冷起 advice 后 flow **不再 stall**、推进到主权 wait（compare_gate/owner_gate）或终态。验收=readmodel 状态非半途 `running` 卡死、GUI 无「卡住无恢复」。
- **A2 · SH-R1-01 回归**（真相源=scene-pack role_slot node_type / bow readmodel；归属=client 05 + truzhenos packstudio；风险=绿验证）：智能家居单角色 `smart_home_project_manager`（责任含「风险」）bow 角色面板生成按钮显「生成**建议**候选」（advice），产出**项目管理**建议（里程碑/采购/排期）而非法律质询；readmodel `node_type='advice'`。验收=候选类型=advice_candidate、产出领域贴切。
- **A3 · 无回归基线**（风险=绿）：3 包 fresh install→enabled、GUI 起 run、03 反查基本链路仍绿；主权链未绕（真实发送/执行/Frappe 写=0）。

### B. 闭合遗留

- **B1 · ENV-R4-01 知识分页**（真相源=09 knowledge SQLite；归属=**truzhenos / 09 knowledge readmodel**；**风险=橙**；契约影响=**触 09 public API/查询契约**）：**本轮只出 API 影响清单**（pagination/cursor/total_count 方案 + 兼容策略 + 反向依赖核查），**不改 public API、不施工**。交付=影响评估文档。
- **B2 · ENV-R4-02 provider_family 透传**（真相源=packs manifest/capabilities.json；归属=**truzhen-packs install.py**；**风险=黄**；契约影响=无，纯 install.py 请求体字段）：install.py 请求体白名单加 `provider_family` 透传（**与 SH-R1-01 node_type 同一类修法**；先核查是否同一白名单代码处）。**注意**：GUI 展示该字段依赖 codex hardening 的后端 enrich + client 展示（未 land）——本轮 packs 侧修根因 + 登记「GUI 展示待 codex hardening land」。验收=install 后 registry/draft 含 provider_family、py_compile+结构审计绿。
- **B3 · Frappe not_ready 专门态**（真相源=registry provider 登记态 + gateway readiness；归属=**truzhenos executor/registry**；**风险=黄**）：构造「provider 已登记但 endpoint/运行体不可用」环境，验证 run/readmodel 显 `not_ready`（区别于 `provider_missing`）、不阻断 advice、写回仍候选/blocked。验收=三态齐（provider_missing/not_ready/blocked）诚实、无假 ready。**若构造成本高 → 降级为影响清单，待 Owner 裁。**
- **B4 · MON-HARDENING-01 消化**（归属=monitor；风险=绿/黄）：复核 hardening final `warning=1` 来源，确认 `security=0 error=0`、非三包阻断项；接入统一 `truzhen-monitor` 体系登记。

### C. codex hardening 对待（Owner 裁：测 origin/main 现状）

- 只测 origin/main 真相（我已 land 的修复）。**codex 已修未 land 的问题**（§1 清单：provider 机器态投影 / 03 反查按钮默认可见 / Frappe requirement 展示 / 服务履约文案 / 长等待提示）在 origin/main 上**仍存在**——遇到**标注「待 codex hardening land」，不重复修、不误报为新 bug**。
- 报告列一张「codex-已修-待land」对照表，明确本轮**不认领**这些，避免与 codex 双改冲突。全程不碰 codex worktree/分支/端口。

---

## 4. 归属 / 风险 / 契约 汇总

| 工作项 | 真相源 | 仓库/层 | 风险色 | 契约影响 | 本轮动作 |
|---|---|---|---|---|---|
| A1/A2/A3 回归 | truzhenos SQLite / bow readmodel | 只读观察 client+truzhenos | 绿 | 无 | GUI 验证 |
| B1 ENV-R4-01 | 09 knowledge SQLite | truzhenos/09 | **橙** | **触 09 public API** | 只出影响清单 |
| B2 ENV-R4-02 | packs manifest | truzhen-packs install.py | 黄 | 无 | 修 install.py 白名单 |
| B3 Frappe not_ready | registry/gateway | truzhenos executor/registry | 黄 | 无 | 构造专门态复测（或降影响清单） |
| B4 monitor warning | monitor 导出 | truzhen-monitor | 绿/黄 | 无 | 消化登记 |
| C codex 对待 | origin/main | — | 绿 | 无 | 标注不重复修 |

---

## 5. 上下文包 / 禁止边界

- **允许读**：3 pack 目录（manifest/flows/role-slots/role-packs/capabilities/install.py）、truzhenos `packstudio/lifecycle`+`sceneflowdev`+`09 knowledge` 只读、client `BusinessObjectWorkbenchPage`/`ScenePackManagePage`/`TransactionTimeline` 只读、codex 三份报告 + 原始轮报告、`env-pack-...-round3-gui-plan-20260704.md`、本轮各台账。
- **禁止边界**：不碰 codex hardening worktree/分支/端口；不改 contracts / 09 public API/DTO；不真实发送/执行/Frappe 写/派工/执法/PDF 真闭环；不覆盖他会话 WIP；不 push/merge（本轮结束呈 Owner 裁 land）；不整仓复制/导入 raw artifacts。

---

## 6. 验收设计（严苛 + 组织者 neuter 独立复核）

- **AI 设计的严苛验收**：A1/A2 回归必须以 **run readmodel + 03 反查 + bow readmodel node_type** 三重证据坐实，非截图自证；B2 install.py 修复必须 neuter 复核（revert install.py → registry/draft 缺 provider_family 真红 → 恢复 → 绿）；B1 影响清单必须列反向依赖 + 兼容策略；主权链未绕以行为审计（真实发送/执行/Frappe=0）坐实。
- **交验收主体复核**：组织者（独立会话/子代理）对每个 PASS 回退源码验真红→绿，禁自报。
- **用户（Owner）如何验收**：读本轮汇总报告 + 台账；抽查 A2 智能家居单角色 advice 产出（应是项目建议非法律质询）；抽查 B2 install.py diff + registry；确认 B1 只出清单未改 API；确认 codex-已修-待land 对照表无重复修。
- **验收维度（改什么证什么）**：回归=GUI+readmodel；packs install.py=py_compile+结构审计+registry+neuter；09 影响清单=反向依赖+兼容；监控=doctor `status/events/security/error/warning`。

---

## 7. 变更影响 / 隔离栈 / 坑

- **变更影响**：B2 改 truzhen-packs 1 处 install.py 白名单（3 包或仅环保，视核查）；B1/B3/B4 主要产文档/复测，B3 若施工触 truzhenos executor（黄）。前端只读观察（不改，除非回归揪出新绿区 bug 再单议）。**均不改契约**。
- **隔离栈**：假 `HOME=/tmp/truzhen-closeout-home` fresh SQLite；专属端口（避 :18092-18094/:5183-5186 我 prior 轮 + codex vite/tsserver 占用，取新段如 :18095/:5187）；devserver 从最新 origin/main 干净 detached worktree 构建；client 从最新 origin/main worktree（`npm ci`）。
- **坑速查（沿用 + 本轮）**：后台 EGR/命令重定向读**自己的 log 目标**非 harness task 输出；client vitest 高负载 flaky（codex 并发 vite/tsserver 抢 CPU）——隔离重跑坐实非真错；`node_type`/`provider_family` 同类白名单丢字段先核查是否同一代码处。

---

## 8. 待 Owner 裁定项（开工前）

1. **B2 ENV-R4-02**：本轮**修** packs install.py provider_family（黄，与 SH-R1-01 同类根因）✅建议修；还是只登记待 codex hardening 一起？
2. **B3 Frappe not_ready**：本轮**构造专门态复测**（黄，需造 provider 已登记但不可用环境）；还是降为影响清单（构造成本高时）？
3. **B1 ENV-R4-01**：确认**只出 09 API 影响清单不施工**（橙）。
4. **land 口径**：本轮若产 B2 packs 修复，结束呈 Owner 裁 land（不自动 push），还是本轮纯验证/文档不产代码？
5. **隔离栈端口段**：确认取 :18095/:5187（或 Owner 指定），避他会话/codex 占用。

---

## 9. 最小可交付

- A1/A2 回归 PASS（我两修复 GUI 生效 + 无回归）+ A3 基线绿。
- B1 ENV-R4-01 09 API 影响清单（文档）。
- B2 ENV-R4-02 packs install.py 修复 land-ready（若 Owner 裁修）+ neuter 复核。
- B3 Frappe not_ready 复测结论（或影响清单）。
- codex-已修-待land 对照表 + 交叉协调点（BusinessObjectWorkbenchPage）登记。
- 汇总报告 + 台账 + memory；缺客户证据显式标注；生命周期止于「已验收（本地隔离）」。
