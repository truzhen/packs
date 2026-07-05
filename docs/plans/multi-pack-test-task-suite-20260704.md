# 多 Pack 测试任务套件 — 用多 pack 各压测 truzhen 不同路径以广覆盖完善（2026-07-04）

> 状态：`设计中`（供 Owner 逐个下发；每份任务=一次独立实测会话）
> 来源：Owner 2026-07-04「根据此次测试的过程，看一下 pack 仓有哪几个 pack，分别做几份测试任务，我用多个 pack 的测试任务来完善 truzhen」。
> 基线方法论：第三轮环保 pack GUI 实测计划 `env-pack-user-simulation-e2e-round3-gui-plan-20260704.md`（全员 Opus + GUI-only 铁律 + 四画像 + REG 回归 + 深水区 + §6.4 真实性铁律 + 组织者 neuter 独立复核）。**各 pack 任务=沿用此方法论，只换 pack 事实基线 + 换「本 pack 主压测的 truzhen 路径」**，不重述通用纪律。

## 0. 策略：为什么用多 pack

truzhen 是平台，pack 是压测它的载体。**不同 pack 的形态差异会把 truzhen 的不同路径逼出来**——单 pack 测不全。仓内 3 个领域 pack 各自形态互补，组合起来能广覆盖：

| Pack | 角色数 | 知识库 | 真实系统绑定 | 主压测的 truzhen 路径（差异点） | 测试状态 |
| --- | --- | --- | --- | --- | --- |
| environmental-enforcement（环保执法证据链） | 2（执法精英+挑剔律师） | **15 知识域 752 条** | 文书送达（provider_missing） | 知识挂载/召回 · 双角色对照 · PDF 真解析 · 执法证据链流程 | 已测 3 轮（含本轮 GUI） |
| housekeeping-ops（家政客户服务全生命周期） | 2（角色对照） | **0 知识库** | capabilities | **无知识库 pack 路径**（与环保相反）· 客户服务生命周期 · capabilities 调用 · 双角色对照 | **本方法论未测** |
| smart-home-owner（智能家居老板项目经营） | **1（单角色）** | 0 知识库 | **Frappe provider（三态）** | **单角色路径**（与双角色对照）· **Frappe 真实系统绑定** · 项目经营/项目对象工作区 | 已测 2 轮（GUI 早期版） |

**覆盖逻辑**：环保压「知识 + 双角色 + 证据链」；家政压「无知识 + 生命周期 + capabilities」；智能家居压「单角色 + Frappe 真实系统 + 项目工作区」。三者跑完，truzhen 的知识/角色/provider/流程四大路径都被真实用户视角压过。

**优先级建议**：① **家政（最高，全新领地）** → ② 智能家居（GUI 方法论 redux + 验本轮修复） → ③ 环保 round-4（验橙修复 + 处置闭环，最低，已充分覆盖）。

---

## 任务一：家政客户服务全生命周期 Pack（scene_pack://housekeeping-ops）— 最高优先

### 派活卡
- **要做的事**：以第三轮 GUI 方法论对家政 pack 跑一次全链路实测，**重点压测「无知识库 pack + 客户服务生命周期 + 双角色对照 + capabilities 调用」路径**。
- **真实客户/场景证据**：家政 pack 是仓内既有资产（`housekeeping-ops-pack-v0`，manifest v0.1.0，2 角色，0 知识，含 capabilities/flows/role-packs/role-slots）；历史有 housekeeping-pack-e2e（converter 修复）但**非本 GUI 方法论**。**开工前 Owner 确认这是当前客户主线还是 backlog 验证**。
- **真相源**：pack 事实=packs 文件夹+registry；流程=truzhenos SQLite；正式动作=03 回执。
- **本 pack 主压测的 truzhen 路径（差异点，重点看这些）**：
  1. **无知识库 pack 路径**：家政 pack 0 知识域——验证 truzhen 对「无知识 pack」的处理（安装/启用/流程推进时不因缺知识而误报/空指针/假挂载；知识召回节点应诚实「本 pack 无知识库」而非报错）。与环保 15 知识域是对照实验。
  2. **客户服务全生命周期流程**：从客户线索→服务派单→上门→回访→投诉处理→结案的完整生命周期（对比环保执法证据链、智能家居项目经营，是第三种流程形态）。看 06 引擎跑不同领域流程图是否中立正确。
  3. **capabilities 调用**：家政 pack 有 capabilities 目录——验证能力荚在流程中被调用的路径（与环保/智能家居的 capability 差异）。
  4. **双角色对照**（复用本轮 dca2179 修复）：验证家政的 2 角色对照本地能否真产出（本轮已证环保双角色本地无需接云，家政应同样）。
- **最小可交付**：pack GUI 加载（enabled）；一次客户服务生命周期 scene flow run 从线索到结案（与秘书互动）；无知识库路径诚实性验证；双角色对照本地产出验证；问题清单+修复/方案。
- **风险颜色**：整体黄（前端/流程/pack 内容绿~黄；网关/主权橙只方案）。
- **不允许碰**：真实对外派单/通知（保持 provider_missing/blocked）；他会话 WIP；契约只读。
- **验收**：GUI 走通生命周期 + 回执 03 反查 + 行为审计纯 GUI + 无知识路径诚实 + 改哪证哪。
- **先输出**：影响清单（本 pack 事实基线 + 主压测路径）→ Owner 确认主线 → 开工。

---

## 任务二：智能家居老板项目经营 Pack（scene_pack://smart-home-owner-project-ops）— GUI redux

### 派活卡
- **要做的事**：以第三轮 GUI 方法论对智能家居 pack 重跑（二轮是早期 GUI 版），**重点压测「单角色路径 + Frappe 真实系统绑定 + 项目对象工作区」+ 回归验证本轮 3 修复**。
- **真实客户/场景证据**：二轮已实测（报告 `smarthome-pack-...-report-20260703.md`，16 橙项+3 数据处置待裁）；本轮=方法论升级 redux + 验本轮修复在此 pack 是否也生效。**开工前 Owner 确认**。
- **本 pack 主压测的 truzhen 路径（差异点）**：
  1. **单角色路径**：智能家居 pack 只 1 角色（对比环保/家政的双角色对照）——验证 truzhen 单角色 slot 的 advice 路径（本轮 dca2179 冷起重试对单角色是否同样生效；单角色时「对照」UI 如何呈现）。
  2. **Frappe 真实系统绑定（三态）**：智能家居 pack 绑 Frappe provider——验证 ProviderRequirement→Registry→Gateway 的真实系统绑定链路（ready/blocked/not_ready 三态诚实），是环保/家政没有的路径。
  3. **项目对象工作区 / 项目经营流程**：与执法证据链、客户服务生命周期是第三种流程；验证 05 项目对象容器（若涉及）。
  4. **本轮 3 修复回归**：R3-P3（已启用包编辑发布，注意智能家居也会遇 canvas-seeding 橙缺口）、R3-P4-04（双角色→单角色 advice 本地产出）、交互批（侧栏导航/回执容器/账号菜单等通用前端修复在此 pack 是否生效）。
- **继承基线**：二轮 16 橙项遇到直接引用不重复分析；本轮橙-1/橙-2（canvas-seeding/skill 路由）对智能家居同样适用。
- **最小可交付**：pack GUI 加载；一次项目经营 scene flow run；Frappe 三态诚实验证；单角色 advice 本地产出；本轮 3 修复回归确认；问题清单。
- **风险/边界/验收**：同任务一口径。

---

## 任务三：环保执法证据链 Pack round-4（scene_pack://environmental-enforcement-flow）— 收口验证，最低优先

### 派活卡
- **要做的事**：环保 pack 已测 3 轮（含本轮 GUI）。round-4 **只做收口验证**，不重跑全量：① 验本轮橙修复 land 后是否闭合（橙-1 canvas-seeding 修后深水①完整闭环、橙-2 skill 路由修后送达触达 provider_missing）；② 完整验证 R3-P4-03 处置闭环（上传真 PDF 过 screening→真 advice→disposition→RiskLow 卡，本轮因案卡在 screening 未完整验）；③ 知识挂载/召回路径深压（15 知识域的召回质量，本轮未深测）。
- **前置**：需先 land 本轮橙-1/橙-2 修复（Owner 裁 + 施工）后才有验证对象；否则 round-4 只做 R3-P4-03 处置闭环 + 知识召回。
- **本 pack 主压测的 truzhen 路径**：知识挂载/语义召回（15 域，truzhen 独有的知识路径深压）+ 处置闭环 + 已 land 橙修复回归。
- **最小可交付**：橙修复闭环验证（若已 land）+ 处置全链闭环（真 PDF→disposition）+ 知识召回质量；问题清单。
- **风险/边界/验收**：同上；round-4 是最小收口不是全量。

---

## 通用（三份任务共用，不重述）
- 方法论、角色结构（全员 Opus）、GUI-only 铁律、§6.4 真实性铁律、组织者 neuter 独立复核、隔离栈搭法（假 HOME 隔离 os.UserCacheDir + 专属端口 + fresh install，避他会话占用端口）、验收口径、坑速查——**全部沿用第三轮计划**，执行会话先读 `env-pack-user-simulation-e2e-round3-gui-plan-20260704.md` + 本轮报告/台账。
- 每份任务开工前按派活纪律：确认版本/优先级（主线 vs backlog）、要真实客户证据、出影响清单先于施工、橙红只方案。
- 每份任务独立 harness 目录（`/Users/li/Documents/过程文档/<pack>-live-r?-<date>/`）+ 独立主任务登记表 + 各仓账本登记。
- **land 政策**：绿黄验收后 land（EGR 绿 + 组织者 neuter 复核）；橙红只方案；他会话活跃移动 origin/main 时不自动 push、land-ready 呈 Owner（本轮先例）。
