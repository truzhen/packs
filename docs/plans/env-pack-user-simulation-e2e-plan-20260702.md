# 环保执法 Pack 用户视角全链路实测与完善计划（2026-07-02）

> 状态：`设计中`（待 Owner 审批后进入执行）
> 计划文件：`/Users/li/Documents/truzhen-packs/docs/plans/env-pack-user-simulation-e2e-plan-20260702.md`
> 修订：2026-07-02 二稿——按 Owner 追加指令补入 专项功能实测（§6A、P4A）与 截图驱动界面成品化（§7A、P4B、§12.9），并补 §3 专项功能代码坐标实测基线。
> 修订：2026-07-02 三稿——按 Owner 追加指令补入 §5.0 模型档位 / 主任务登记 / 上下文隔离纪律：主任务全程登记；问题一律外派独立智能体（尽量 Sonnet 档）分析解决，不占主线上下文；用户测试智能体用最高档模型并模拟用户全部行为。

## 0. 一句话目标

组建「用户智能体（只走前端 UI / computer use）+ 组织者（我，前后端协调）+ 监控程序（truzhen-monitor）」三方结构，以真实用户视角走通：三制作台使用 → 加载环保执法 pack → 在制作台修改完善它 → 用它创建一个执法行动 → 与秘书长互动完成执法闭环；过程中**发现问题→修问题**，同时完善 Truzhen 程序与环保执法 pack。

本轮同时承担两项 Owner 追加目标（2026-07-02 追加指令）：

1. **专项功能实测**：三荚卡片封面开关（场景包 / 能力包 / 角色包管理页）+ 事务对象卡片封面开关（前端）+ 场景项目菜单启动/暂停，逐项经 UI 实测并由组织者后端复核（细则见 §6A）。
2. **截图驱动界面成品化**：测试沿途用到的每个界面逐屏截图审计——**去掉工程信息、前端只显示中文、排版布局优化到成品标准**——整改后出前后对比截图（标准与流程见 §7A，界面清单见 §12.9）。

## 1. 版本 / 优先级（需 Owner 确认）

- Agent 建议定位：**当前主线**——收尾阶段「全链路真实连接」的用户视角验收测试，不是新功能开发。环保执法 pack 是 Owner 裁定的护城河标杆样板（非第一现金流），用它做全链路实测正是标杆样板的用途。
- 非目标：不做通用测试平台、不做新制作台功能、不建 Governance Console。

## 2. 真实客户 / 场景证据

- 证据 = Owner 本轮原话直接指令（2026-07-02）：「使用场景制作台、能力制作台、角色制作台；读取当前 pack 仓库的环保执法 pack，对他加载并作修改……创建一个执法行动，通过和秘书长的互动，完成执法……目的是用这个测试遇到问题解决问题」。
- 环保执法领域证据：pack 知识库来自 Owner 权威资料（`/Users/li/Documents/trae/knowledge_base`，752 条已真实导入），此前已有污水厂秘书对话 E2E 夜跑（已合 main）历史。

## 3. 现状基线（2026-07-02 实测）

- 环保执法 pack：`/Users/li/Documents/truzhen-packs/environmental-enforcement-pack-v0/`（manifest v3、15 知识域、双角色包执法精英+挑剔律师、flow 9 阶段、install.py/uninstall.py）。
- devserver 18080 registry 里 `scene_pack://environmental-enforcement-flow` 已注册但 **disabled**（enabled_pointer.current_version=""，previous 1.0.2）——第一步「加载」就是真实用户任务。
- 运行环境：vite :5173 ✅、devserver :18080 ✅、oMLX :8000 ✅（均已在跑）。
- client 前端 R1-R19 批次已 land（f66a999），三制作台 v2 是全新表面，正需要用户视角实测。
- Playwright chromium 本机已有缓存（`~/Library/Caches/ms-playwright`），可搭 computer-use harness。

### 3.1 专项功能代码坐标（2026-07-02 代码实测基线，client main f66a999 / truzhenos main b76e0eb）

以下为「代码里现在有什么」的实测事实，执行会话不必重查；但**代码里有类名 ≠ 功能可用**，P4A 仍须 UI 实测：

- **卡片封面开关 `pack-cover-toggle`**：
  - 场景包管理页 `src/pages/ScenePackManagePage.tsx:600` 已有（停用走 `apiClient.disableScenePackLifecycleVersion` 主权动作，守卫测试 `src/pages/__tests__/scenePlatformPackLifecycle.test.tsx`）。
  - 能力包（场景平台页）`src/pages/ScenePlatformPage.tsx:1069` 已有（`isCapabilityPackLoaded` 态）。
  - 角色包管理页 `src/pages/RolePackManagePage.tsx` **未发现封面开关**，页内明示「角色包不是新智能体，也不是权限开关」（:800）——语义可能是绑定槽位而非启停，实测时先确认产品语义再定是否缺口。
- **事务对象卡片封面**：封面流程节点角标真接 06 场景流程 run ReadModel（`src/pages/BusinessObjectWorkbenchPage.tsx:2128-2226`），卡片组件 `src/components/business-objects/TransactionCapsuleCard.tsx`；**卡片与组件内均未发现「开关」控件**——若 UI 实测确无，即为专项缺口（处理见 §6A 第 4 条）。
- **场景项目侧栏菜单**：侧栏「场景项目」面板在 `src/components/layout/AppShell.tsx:616` 起；条目下拉菜单现仅有 重命名 / 归档 / 删除，**无启动/暂停**。「启动流程」入口现只在事务对象工作台（`BusinessObjectWorkbenchPage.tsx:847-965`，对任何已激活场景包通用、经 Base Gate 发起 SceneFlowRun）。
- **暂停能力**：06 场景流程后端 `truzhenos backend/internal/devserver/sceneflowdev/` **全无 pause/paused 实现**；前端 `FLOW_STAGE_LABELS` 里的 `paused: '已暂停'` 只是防御性映射。resume 已有三证主权链（owner_ref / gate_decision_ref / evidence_refs，`src/api/client.ts:5832-5862`）；task-governance schedules 的 pause/resume/cancel 前端只读暂不接真（`client.ts:5763-5772` 注释，缺 server-issued 包装）。

## 4. 最小可交付

即使修复循环没跑满，本轮最少交付：
1. 环保执法 pack 经前端 UI 完成加载（disabled→enabled，有 receipt）；
2. 一次执法行动 scene flow run 从线索受理推进到 Receipt 归档（与秘书长互动完成）；
3. 问题清单（每项：现象+截图 → 根因 → 归属仓/风险颜色 → 修复或方案）；
4. 专项功能实测结论（§6A 六条逐项）：每项 = 可用证明（截图 + 组织者后端复核）或诚实缺口记录（归属仓 + 风险颜色 + 修复 commit 或方案）；
5. 测试主路径界面（至少 P2 / P4 沿途各屏）完成 §7A 成品化整改：去工程信息、全中文、成品排版布局，每屏前后对比截图归档，§12.9 清单打勾。

砍掉项 / backlog：三制作台的功能增强建议（只记录不施工）；执法行动的多案并行；移动端视角；**非测试沿途界面的成品化**（只记 backlog，防 scope 爆炸）；场景流程「暂停」的后端实现（橙，默认只出方案，见 §9 第 5 项）。

## 5. 角色结构与派活卡

### 5.0 模型档位、主任务登记与上下文隔离纪律（Owner 2026-07-02 追加，全阶段生效）

**模型档位（派子代理时显式指定，写进每张派活卡）：**

| 角色 | 模型档位 | 理由 |
| --- | --- | --- |
| 用户智能体（测试） | **最高档模型（Fable 5，`model` 缺省继承主会话即是）** | Owner 指令：测试智能体用最高模型，模拟用户全部行为 |
| 问题分析智能体 / 修复工 | **Sonnet 档（`model: "sonnet"`，尽量）** | Owner 指令：分析解决问题的智能体尽量用 Sonnet，不用 Fable |
| 组织者（主线） | Fable 5（主会话本体） | 只做派活/验收/登记/简报，不做长篇分析 |

- Sonnet 档确实啃不动的问题（如橙色主权方案设计、跨仓根因），允许升档，但**升档必须在主任务登记表记一行原因**，不得默认升。

**主任务登记（两级，开工即建，随做随更）：**

1. 各仓 `FEATURE_LEDGER.md` §0 在途区登记本测试任务（开工登记、收尾按 §12.8 转正/销账）；
2. 运行期主任务登记表 `logs/main-task-register.md`（harness 目录内，见 §12.5）：逐行记录 阶段状态（P0-P7 各档位）、问题编号（ISSUE-NNN）、派出的智能体（角色/模型档/派出时间）、返回结论摘要（≤3 行）、修复 commit 或方案路径、land 状态。**每次派活与收活即更新，主线恢复上下文以此表为准。**

**上下文隔离（防主任务上下文污染）：**

- 主线（组织者）**不在自己上下文里做问题分析**：测试中每遇到一个问题，先在主任务登记表立 ISSUE-NNN，然后**单独派一个分析智能体**（Sonnet 档）带独立上下文包去分析定位；需要施工再派修复工（Sonnet 档，独立 worktree）。
- 分析/修复智能体的上下文包必须自足：问题现象 + 截图绝对路径 + 复现步骤 + §3.1 相关坐标 + 归属仓与风险颜色初判 + 验收命令（§12.6）+ 禁止边界（§8）。它们**不共享主线对话历史**。
- 返回格式固定：`根因 / 归属仓 / 风险颜色 / 修复 commit 或方案文件绝对路径 / 验证证据`，主线只收这份结构化结论并回填登记表；分析过程长文一律落该智能体自己的 markdown（`logs/issues/ISSUE-NNN.md`），不进主线对话。
- 用户智能体不因某个问题停摆：问题登记外派后，用户智能体绕开继续主旅程（该问题阻断主路径时才等修复）。

### 5.1 用户智能体（User Agent，子代理，长期会话可续）

- **模型档位**：最高档模型（Fable 5，§5.0）。
- **只允许**：通过 Playwright 驱动真 Chromium 操作 `http://127.0.0.1:5173`，每步截图；用 Read 看截图理解界面（computer-use 循环：截图→理解→操作）。
- **禁止**：curl / fetch 后端、读产品源码、读数据库、用开发者知识绕过 UI。UI 走不通 = 记录产品问题，不许绕。
- **行为全覆盖（Owner 追加：模拟用户的所有行为）**：主路径沿途每个界面，对**每个可见交互控件**至少各触发一次真实用户行为——点击、悬浮（悬浮卡/tooltip）、文本输入（含空输入、超长输入、中文输入）、下拉/展开收起、滚动（页滚与卡片区）、拖拽（若有）、键盘（Tab 焦点走查、Enter 提交、Esc 关弹层）、刷新后状态复看、返回/切页再进；每屏产出一份**行为覆盖矩阵**（控件 × 行为 × 实际结果 × 截图路径），未能触发的控件诚实标注原因。误操作路径也要模拟（取消、重复点击、必填漏填），验证产品的容错与提示。
- 选择器纪律：优先用户可见文本 / role；仅在可见文本无法唯一定位时才允许 data-testid，且每次使用记为一条潜在 UX 发现（真实用户无法凭可见信息区分控件）。
- 产出：每阶段一份用户日志（做了什么、看到什么、卡在哪、截图路径）。
- Harness 与截图目录（稳定路径，跨会话可复用，不入任何仓）：`/Users/li/Documents/过程文档/envpack-live-20260702/`（`harness/` 放 playwright 脚本，`shots/` 放截图，`logs/` 放用户日志）。

### 5.2 组织者 / 协调者（我，主线）

- 派活、验收用户日志、前后端协调检查（后端日志、SQLite、03 回执反查、监控事件）、问题分诊、修复派工、阶段简报、主任务登记表维护（§5.0）。
- 对每个用户报告的问题先独立复核（不轻信自报），再定根因与归属——**复核与根因分析本身也外派分析智能体（Sonnet 档）**，主线只读结论回填登记表，不在主线上下文展开分析（§5.0 上下文隔离）。

### 5.3 问题分析智能体 + 修复工（按问题派，Sonnet 档）

- 每个 ISSUE-NNN 先派分析智能体（Sonnet 档，独立上下文包，见 §5.0）定根因/归属/颜色；确需施工再派修复工。
- 每个修复在**归属仓**独立 worktree 分支，TDD（先加失败测试再修），改哪证哪；修复工同为 Sonnet 档，验收由组织者按 §10 复核，不采信修复工自报。

### 5.4 监控程序

- 起 `truzhen-monitor`（`/Users/li/Documents/truzhenos/backend/cmd/truzhen-monitor`），全程收集 `/v3/monitoring/*` 事件；测试收尾导出诊断包。发现监控缺口（该有事件没有）→ 补回现有体系，不另起格式。

## 6. 阶段设计

| 阶段 | 视角 | 内容 | 验收 |
| --- | --- | --- | --- |
| P0 环境基线 | 组织者 | 确认三服务、起 truzhen-monitor、搭 Playwright harness、user agent 首页可达冒烟 | 首页截图 + monitor 收到事件 |
| P1 三制作台巡检 | 用户 | 场景/能力/角色制作台各走一遍最小用户旅程（进入→看引导→做一次最小操作→退出） | 三份用户日志 + 问题记录 |
| P2 加载 pack | 用户 | 通过 UI 把环保执法 pack disabled→enabled（lifecycle 主权确认路径）；若 UI 无此入口即为 P0 级产品缺口 | enabled 状态 + receipt 可在 03 反查（组织者核） |
| P3 修改完善 pack | 用户+组织者 | 场景制作台加载 env pack 规格，做真实完善（候选见 §7），走 draft→readiness→promote→confirm→enable 新版本 | 新版本 enabled + pack 仓文件同步更新 |
| P4 执法行动 | 用户 | 用 pack 创建执法行动（scene flow run）：与秘书长对话推进 线索受理→立案→现场检查取证→证据三性审查→执法精英×挑剔律师对照→处置推演→文书草稿→Owner 门→Base Gate→Receipt | run 全节点 evidence + 最终 receipt 反查 |
| P4A 专项功能实测 | 用户+组织者 | §6A 六条逐项：三荚卡片封面开关（场景包/能力包/角色包）、事务对象卡片封面开关、场景项目菜单启动/暂停；每项 UI 实测 + 组织者后端复核，缺口按 §6A 判定处理 | 每项结论 = 可用（截图+后端事实）或缺口闭环记录 |
| P4B 界面成品化 | 用户+组织者+修复工 | 对 P0-P4A 沿途每个界面按 §7A 三项标准（去工程信息/全中文/成品排版布局）截图审计 → client 仓整改 → 复截图对比 | §12.9 界面清单全屏过审：每屏 审计记录+整改 commit+前后截图 |
| P5 验收对账 | 组织者 | 回执 03 反查、监控事件对账、诊断包导出、用户智能体行为审计（证其全程只走 UI） | 对账清单全绿或诚实标注 |
| P6 修复循环 | 贯穿 | 每问题：登记 ISSUE-NNN → 外派分析智能体（Sonnet 档）定根因 → 分诊：绿=派修复工直接修+自测；黄=修后阶段简报确认；橙/红=只出方案待 Owner；全程主线只收结构化结论（§5.0） | 每问题在主任务登记表闭环记录 |
| P7 收尾 | 组织者 | 测试报告、pack 完善提交、FEATURE_LEDGER 登记、land（按 §9 裁定） | 报告文件 + 账本更新 |

## 6A. 专项功能实测细则（Owner 2026-07-02 追加；P4A 执行脚本）

原则：逐项「先实测存在性 → 再测行为 → 组织者后端复核」；禁止把「代码里有类名/端点」当作「功能可用」，也禁止把缺口默默补成想当然的样子——语义不清先按下述判定处理。

1. **场景包卡片封面开关**（场景包管理页）：对环保执法 pack 卡片执行 开→关→开 各一次；组织者复核 lifecycle 端点状态真实翻转 + receipt 在 03 反查出原文；停用后卡片必须仍保留在列表（不消失）。预期已可用（§3.1 有守卫测试），实测重点是真实后端下的行为与 UI 反馈。
2. **能力包卡片封面开关**（场景平台/能力包管理页）：任选一个已装能力包同上实测；若无已装能力包，先经 UI 装一个（本身属 P1 用户旅程的一部分）。
3. **角色包卡片封面开关**：先实测角色包管理页卡片有无开关。若无——先确认产品语义（§3.1：页内明示角色包非权限开关，可能是绑定槽位语义）：语义如此则记为「用户预期与产品语义差异」的 UX 发现（前端应把语义讲清楚），不立缺陷；语义应有开关而缺失才立缺陷走修复循环。
4. **事务对象卡片封面开关（前端）**：实测事务对象工作台卡片封面（流程节点角标区）有无开关控件。§3.1 已证代码中无 → 预期结论为缺口，本轮补做（黄，client 仓，前端交互）。**Agent 建议的开关语义（Owner 批计划即认可，或在 §12.1 回填改判）**：开关 = 卡片封面（流程节点角标区）展开/收起的本地展示偏好，纯前端展示态；封面数据仍真接 06 run ReadModel，「未启动流程 / 流程状态不可用」等诚实态文案不得因开关伪装或消失；开关不是主权动作、不产生 receipt、不得伪装成启停流程。
5. **场景项目菜单「启动」**：实测侧栏「场景项目」条目菜单有无启动入口。§3.1 已证现无 → 预期结论为缺口，本轮补做（黄，client 仓）：**复用**事务对象工作台既有启动链（已激活场景包 → Base Gate 发起 SceneFlowRun，`BusinessObjectWorkbenchPage.tsx:847-965`），把入口接到菜单；**不得另铸第二条启动链路**（防重复创建口径与页内注释一致）。已有在途 run 的项目，菜单入口须给出诚实提示而非重复发起。
6. **场景项目菜单「暂停」**：§3.1 已证 06 后端全无暂停能力 → 本项为**橙色缺口**：暂停一个 scene flow run 是主权动作，须 server-issued decision、与 resume 三证口径对齐。默认处理 = 本轮只出《scene flow run 暂停能力影响清单与方案》（落 `/Users/li/Documents/truzhen-packs/docs/reports/` 同名前缀），**前端不得先做假暂停按钮**（点了没有真效果即假成功，违反收尾红线）；Owner 若在 §9 第 5 项授权橙色施工，另派修复工按方案 TDD 实现。

## 7. Pack 完善候选（P3 具体改什么，按测试发现动态增减）

1. 把 install.py 里散落的门槛信息（模型/provider 要求）对齐进 manifest 声明；
2. 流程节点的用户可读提示语完善（9 阶段 stage guide 口径与前端一致）；
3. 上轮已删 solid-waste 假声明——本轮核对 15 域 knowledge-scopes 与实际文件一一对应；
4. role-packs 双角色（执法精英/挑剔律师）的 steering 提示按实测对话质量调优；
5. 测试报告 `docs/测试报告.md` 更新为本轮真实结果。

## 7A. 界面成品化标准与流程（Owner 2026-07-02 追加；P4B 执行脚本）

### 7A.1 三项标准（每屏逐条过）

1. **去掉工程信息**——用户可见区域不得出现：
   - 内部 ref 原文直排（`scene_pack://…`、`transaction://…`、`receipt://03/…`、`run-gated-…` 等 URI/句柄裸露在正文）；
   - 模块编号口径（「05 事务对象 ReadModel」「06 场景流程」「03 回执账本」等对用户无意义的内部编号叙述）；
   - 英文状态枚举原文（running / blocked / candidate_only / provider_missing…）；
   - 开发批次/任务代号（FE-1、R2-P2 等）、data-testid 式文案、面向开发者的调试性长句（如「数据源限定为 05 事务对象 ReadModel；不可用时不显示演示对象」应改写为用户语言）。
   - 改法 = 翻译成用户语言（例：「数据来源：真实事务对象」「后端服务未连接，暂时无法读取」），或折叠进「详情」二级展示。
   - **主权红线（只能改呈现，不得删除）**：高风险动作确认面展示 目标/内容/影响/run_id/nonce 是受控真实 E2E 的主权要求；回执/凭证 ref 必须保持可反查。整改只允许成品化呈现（折叠「凭证详情」、等宽字体、一键复制），禁止删除、禁止隐藏到不可达；涉及确认面要素**增删**的改动一律先出方案不施工（红）。
2. **前端只显示中文**：所有用户可见文案中文；英文仅限专有名词（Truzhen、pack 英文名等 Owner 认可者）与必须原样保留的错误原文（收进「详情」内）。状态枚举一律走中文映射：补齐 `FLOW_STAGE_LABELS` 一类映射的漏网值；未知枚举兜底显示中文（如「状态：未知（原文见详情）」），不得直排英文原文。
3. **成品排版布局**：对齐既有主题 token 体系（R1 主题批次成果），不另起平行样式；卡片无内部小滚动条、滚轮不滚单卡（C9 口径回归）；同屏间距/字号/层级一致；空态/错误态/加载态三态文案完整且中文；长文本截断 + 悬浮/详情可见全文；散落 inline style 造成的错位改用主题 class 收敛。

### 7A.2 流程（每屏一循环，截图是唯一驱动）

1. 用户智能体进入该屏即存基线截图 `shots/audit/<screen>-before.png`（多状态屏按状态各存一张：空态/有数据/错误态尽量都留）；
2. 组织者按 7A.1 逐条出该屏审计记录（违反项：截图内位置 + 违反哪条 + 建议改法），落 `logs/ui-audit-<screen>.md`；
3. 修复工在 client 仓整改（绿/黄；纯文案改动可攒批一个 commit，布局改动按屏拆 commit）；
4. 复截图 `shots/audit/<screen>-after.png`，组织者前后对比确认；
5. §12.9 界面清单该屏打勾。发现某屏问题属后端（如后端把英文枚举直接当 message 返回）→ 按归属仓走修复循环，不在前端硬编码遮丑。

### 7A.3 防回潮与边界

- 每改完一屏随手加守卫：对已成品化组件补「用户可见文案不得出现英文枚举原文 / 内部 ref 直排」断言，沿用 `src/components/__tests__/frontendStrict.test.tsx` 一类真实性测试口径（真渲染真断言，不 mock fetch 编假数据）。
- 只收测试沿途界面（§12.9 清单为准，执行时按实际路由增删）；非沿途屏发现问题记 backlog，不做全仓一次性文案大扫除。
- 成品化不改行为：本节所有整改不得改变任何主权链路、端点调用与门控行为；行为缺陷一律走 §6 修复循环而非借成品化夹带。

## 8. 治理维度

- **真相源**：pack 事实 = truzhen-packs 文件夹 + devserver registry；运行态/流程 = truzhenos SQLite；正式动作事实 = 03 回执账本；监控 = truzhen-monitor。前端只是投影，用户日志只是观察记录。
- **仓库归属**：pack 完善→`truzhen-packs`；前端问题→`truzhen-client-web-desktop`；后端/网关/流程问题→`truzhenos`；契约问题→`truzhen-contracts`（**只出方案不施工**）。本计划即为跨仓授权申请，Owner 批准计划 = 授权在上述四仓（contracts 仅读）作业。
- **风险颜色**：整体黄。前端/文档/pack 内容=绿~黄；devserver 业务逻辑=黄；网关/主权闸/契约=橙红（只出方案）。专项与成品化细分：封面开关实测=绿（只读观察+既有主权动作）；事务对象封面开关补做、场景项目菜单「启动」接线（复用既有主权端点）=黄；场景流程「暂停」（需 06 新主权动作端点）=**橙，默认只出方案**；UI 成品化=绿~黄，但凡涉及高风险确认面要素增删=**红，只出方案**（§7A.1 主权红线）。
- **契约影响**：默认零契约改动。若测试暴露 contracts 缺口，单独出影响清单待裁。
- **上下文边界**：允许读四仓代码/文档、devserver 日志与本地测试库；用户智能体上下文仅限 UI 与自己的截图。
- **禁止边界**：
  - 不做任何真实对外发送/真实执法文书投递（`outbound_legal_document_delivery` 等保持 provider_missing/blocked 诚实态）；
  - 不碰真实客户数据、不碰他会话 WIP（C8 codex worktree、superpowers worktree）；
  - `_source-materials` 权威资料只读；
  - 不自动 merge 橙/红改动；不删任何 worktree/分支。
- **生命周期档位**：本计划=`设计中`；执行后测试结论=`已验收`（有 receipt/截图/对账证据）；修复项逐项走 `已实现→已接线→已验收`。

## 9. 待 Owner 裁定项（批计划时一并回答）

1. **修复 land 政策**：建议——绿色修复攒批、每阶段简报后由我统一 land 对应仓 main；黄色修复简报确认后 land；橙/红只出方案。是否同意？（另选项：全部攒到 P7 一次性等 Owner 逐项裁。）
2. **秘书长对话模型**：建议本地 oMLX（Qwen3.5 系）优先、零成本可长跑；若对话质量不足以推进执法节点，再请示切云模型（ark key 待轮换，需 Owner 点头）。
3. **Owner 门代点授权**：执法行动推进会遇到 Owner 主权确认点。建议——本机测试数据范围内，授权用户智能体代点 Owner 确认（构成测试态 OwnerActionEvidence，全部留 receipt），但任何真实对外动作保持 blocked。是否授权？
4. **时长/预算**：预计多轮长任务（P1-P7 数小时级）。无人值守跑到底，还是每阶段暂停等你看简报？
5. **场景流程「暂停」缺口处理**：06 后端现无暂停能力（§3.1 实测）。建议——本轮只出影响清单与方案（橙），前端不做假暂停按钮；是否同意？（另选项：授权本轮橙色施工，暂停按与 resume 对齐的三证主权链 TDD 实现。）
6. **事务对象卡片封面开关语义**：Agent 建议 = 纯前端展示偏好（封面展开/收起），不产生 receipt、不伪装启停流程（§6A 第 4 条）。按此实现是否同意？（另选项：Owner 另定语义，回填后再施工。）

## 10. 验收设计（严苛复核，非自报）

- 用户智能体自报不算数：P2/P3/P4 每个关键声明由组织者用后端事实复核（registry 状态、SQLite 行、03 回执原文反查、监控事件流）。
- 反伪造检查：receipt_ref 必须能在 03 反查出原文；scene flow run 节点 evidence 必须真实存在；禁止把 candidate 态冒充 formal 完成。
- 用户智能体行为审计：抽查其 harness 脚本与 transcript，确认无后端直连。
- 行为全覆盖验收：每屏行为覆盖矩阵（§5.1）齐全——控件清单与 after 截图逐一对照，无「代码里有但从未被触发」的沿途控件；未触发项有诚实原因标注。
- 主任务登记验收：`logs/main-task-register.md` 与实际派活/收活一致（抽查 3 个 ISSUE：登记时间、智能体模型档、结论、commit/方案路径可复核），FEATURE_LEDGER §0 在途登记与收尾登记齐。
- 专项功能验收（P4A）：每项按 §6A 结论口径——「可用」必须有 截图 + 组织者后端事实复核（状态真实翻转 / receipt 03 反查 / SQLite 行）；「缺口」必须闭环记录（归属仓 + 风险颜色 + 修复 commit 或方案文件路径）。补做的封面开关 / 菜单启动接线，验收须证明：真实端点被调用（非本地 state 假翻转）、诚实态不因整改消失、既有守卫测试与新增测试全绿。
- 界面成品化验收（P4B）：§12.9 清单全屏「审计记录 + 整改 commit + 前后截图」三件套齐；每屏 after 截图逐条对照 §7A.1 三项标准复核；新增文案守卫测试为真渲染真断言；抽 3 屏由 Owner 亲眼比对前后截图。
- 回归验收：改哪证哪——client 改动跑 typecheck+test+smoke；truzhenos 改动跑 `go test`+EGR `bash scripts/verify.sh`；packs 改动跑隔离 devserver install/uninstall E2E。
- Owner 验收：刷新 :5173 亲手复走 P2/P4 关键路径 + 亲手点一遍 §6A 六项 + 读最终测试报告。

## 11. 变更影响预览

- 会碰：client UI（修 bug + §7A 成品化文案/排版整改 + §6A 事务对象封面开关补做与场景项目菜单启动接线）、truzhenos devserver（修 bug）、truzhen-packs env pack（内容完善）、各仓 FEATURE_LEDGER、监控事件流（新增测试期事件）。
- 不碰：contracts schema、主权闸/Base Gate 逻辑（发现问题只出方案）、高风险确认面的主权要素集合（run_id/nonce/目标/影响，只改呈现不增删）、场景流程暂停的后端实现（默认只出方案，见 §9 第 5 项）、支付/市场、他会话 WIP、生产凭据。

 
---

## 12. 新会话无缝接手背景（执行会话开工必读）

> 本节写给执行本计划的新对话：读完本节 + §9 裁定回填区即可直接开工，无需回溯旧会话。

### 12.1 Owner 裁定回填区（开工前先看这里）

| §9 裁定项 | Owner 答复（审批时回填） |
| --- | --- |
| 1. 修复 land 政策 | （待回填） |
| 2. 秘书长对话模型 | （待回填） |
| 3. Owner 门代点授权 | （待回填） |
| 4. 节奏（无人值守 / 阶段暂停） | （待回填） |
| 5. 场景流程「暂停」缺口处理 | （待回填） |
| 6. 事务对象卡片封面开关语义 | （待回填） |

若 Owner 审批时未逐项回填，按建议值执行：绿/黄阶段简报后统一 land、本地 oMLX 优先、测试数据内代点 Owner 确认（真实对外动作保持 blocked）、每阶段出简报但不阻塞继续、暂停只出方案不施工、封面开关按纯前端展示偏好语义实现。

### 12.2 六仓布局与基线（2026-07-02 实测）

| 仓 | 路径 | main 基线 | 本计划角色 |
| --- | --- | --- | --- |
| truzhenos（基座） | `/Users/li/Documents/truzhenos` | `b76e0eb` | 后端问题修复；devserver/monitor 宿主 |
| truzhen-packs | `/Users/li/Documents/truzhen-packs` | `bdcafce` | 环保执法 pack 完善（主交付） |
| truzhen-client-web-desktop | `/Users/li/Documents/truzhen-client-web-desktop` | `f66a999` | 前端问题修复；用户智能体操作对象 |
| truzhen-contracts | `/Users/li/Documents/truzhen-contracts` | — | **只读**；缺口只出方案 |
| truzhenv3（旧主仓） | `/Users/li/Documents/truzhenv3` | 已封棺 | 只读历史参考 |

- client main `f66a999` = 刚 land 的 R1-R19 前端调整批次（三制作台 v2 / StudioShell / 能力包对话式调用 / 主题管理程序 / 授权同步状态机）——**本测试的对象正是这批新表面**。批次文档：`/Users/li/Documents/truzhen-client-web-desktop/docs/agent-task-briefs/six-repo-task-breakdown-frontend-adjustments-20260702.md` 与同目录 `execution-progress-frontend-adjustments-20260702.md`。
- 注意：client 主仓目录 checkout 可能停在他会话分支（此前是 `codex/p0-p1-optimization-20260701`）；本地 `main` 分支可能被 superpowers worktree（`/Users/li/.config/superpowers/worktrees/truzhen-client-web-desktop/merge-main-frontend-adjustment-20260702`）checkout——同步 main 用 `git -C <该worktree> merge --ff-only origin/main`，不要 `branch -f`。
- **任何修复开发一律新建独立 worktree 到 `/Users/li/Documents/truzhenv3worktree/`（truzhenos）或各仓约定位置，绝不在主仓目录直接开发。**

### 12.3 运行环境（先探活再起，三服务此前均已在跑）

```sh
# 探活
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:18080/   # devserver（404=活着）
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:5173/    # vite（200=活着）
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8000/    # oMLX 模型侧车（404=活着）

# 若没起：
cd /Users/li/Documents/truzhenos && go run ./backend/cmd/devserver          # :18080
cd /Users/li/Documents/truzhen-client-web-desktop && npm run dev -- --host 127.0.0.1   # :5173（vite 默认绑 IPv6，必须 --host 127.0.0.1）
# vite 代理后端：env TRUZHEN_DEV_BACKEND_PROXY_TARGET=http://127.0.0.1:18080（默认即此）
# oMLX 由 Owner 侧常驻（:8000）；没起时秘书回复会失败——先报告再说，不得假成功
# 监控：/Users/li/Documents/truzhenos/backend/cmd/truzhen-monitor（P0 起，统一走 /v3/monitoring/*）
```

### 12.4 环保执法 pack 事实基线

- Pack 目录：`/Users/li/Documents/truzhen-packs/environmental-enforcement-pack-v0/`（manifest v3、`pack_ref=scene_pack://environmental-enforcement-flow`、15 知识域、role-packs 执法精英+挑剔律师、flow 文件 `flows/environmental-enforcement-flow.flow.json`、install.py/uninstall.py；`_source-materials` 权威资料只读）。
- devserver registry 现状（实测）：pack 已注册但 **disabled**（`enabled_pointer.current_version=""`，previous `1.0.2`，default `1.0.0`）——P2「加载」就是把它经 UI 主权路径启用。
- 关键后端端点（组织者复核用，用户智能体禁用）：
  - `GET /v3/pack-studio/lifecycle/packs?pack_ref=scene_pack://environmental-enforcement-flow`（pack 状态）
  - lifecycle：`POST /v3/pack-studio/lifecycle/{draft|readiness|promote|confirm|reactivate}`
  - 角色包：`GET /v3/agent-orchestration/role-packs/readmodel`、agent-slots 同前缀
  - 知识：`POST /v3/memory/knowledge/batches`
  - install.py 用法：`TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:<port> python3 install.py`（隔离验证时指向独立端口的隔离 devserver，勿对着 18080 重装除非计划明确）
- 秘书长闭环历史证据：污水厂秘书对话 E2E 已证真（`start_scene_flow_run` skill + VerifyIssuedDecision 主权闸 + 云模型对话办案，已在旧主仓合并线）；本轮前端入口 = 沟通中心「新对话」。
- 专项功能（三荚封面开关 / 事务对象封面开关 / 场景项目菜单启动暂停）的代码坐标实测基线见 §3.1，执行细则见 §6A，不必重查。

### 12.5 Computer-use harness（用户智能体）

- 稳定目录：`/Users/li/Documents/过程文档/envpack-live-20260702/`（`harness/`、`shots/`、`logs/`；成品化审计截图放 `shots/audit/`，每屏 `-before.png` / `-after.png` 成对，审计记录放 `logs/ui-audit-<screen>.md`；**主任务登记表 `logs/main-task-register.md`、问题分析长文 `logs/issues/ISSUE-NNN.md`、行为覆盖矩阵 `logs/behavior-matrix-<screen>.md`**）。
- 模型档位纪律见 §5.0：用户智能体=最高档（Fable 5）；问题分析/修复智能体=Sonnet 档（派活时显式 `model: "sonnet"`）；升档须登记原因。
- 初始化：`cd harness && npm init -y && npm i playwright`（chromium 本机已有缓存 `~/Library/Caches/ms-playwright`，通常免下载；缺则 `npx playwright install chromium`）。
- 循环：脚本操作 `http://127.0.0.1:5173` → 每步截图到 `shots/` → 智能体 Read 截图理解 → 决定下一步。选择器纪律见 §5.1。
- 用户智能体建议用可续会话（组织者记住 agent 名，跨阶段 SendMessage 续用，保持「同一个用户」的连续记忆）。

### 12.6 验证命令映射（改哪证哪）

- client：`npm run typecheck && npm run test && npm run build && npm run smoke:frontend-shell`
- truzhenos：`go test ./... && go test -race ./...`（范围见 `scripts/go-test-packages.txt`）；land 前 EGR `bash scripts/verify.sh`——**整体委托子代理后台跑（约 24 分钟），日志落固定文件，读真实 `VERIFY_EXIT` 值，禁双层后台化（`&` 嵌 `run_in_background`）**
- packs：隔离 devserver 起独立端口 → install.py/uninstall.py 全程 E2E → registry 0→1→0 铁证

### 12.7 已知坑速查

- zsh：不做无引号变量分词（for 失效改 while-read）；`echo ===`/裸 glob 会炸；`grep|head && echo` 退出码恒 0 假绿。
- 长命令 `run_in_background: true` 即可，命令里**不得再加 `&`**（双层后台化会丢完成通知）。
- 前端截图/文件写入后用读回（read-back）确认，工具可能静默写失败报 success。
- 秘书回复失败先查 oMLX(:8000) 与 devserver 是否都活，再查 08 网关——历史上「预览生成失败」根因就是后端没起。
- devserver 可能由别的 checkout 起着——复核 `/v3/` 行为异常时先确认它跑的是哪个工作树的代码。
- FEATURE_LEDGER 的 CJK 长行用 Edit 易 unicode 匹配失败——改用 python utf-8 读写插入。

### 12.8 收尾要求

- 测试报告写入 `/Users/li/Documents/truzhen-packs/docs/reports/`（本计划同名前缀），pack 的 `docs/测试报告.md` 同步更新；专项功能六项结论、§12.9 界面清单终态、主任务登记表终态（全 ISSUE 闭环状态）并入报告。
- packs / truzhenos / client 各自 FEATURE_LEDGER 开工登记 §0 在途、收尾登记功能域行。
- 对话简报 ≤500 字：结论、报告绝对路径、关键验证、待 Owner 决策项。

### 12.9 界面成品化清单（P4B 逐屏打勾；按执行时实际路由增删，增删须在清单中留痕）

测试主路径预计经过的界面（初版，B=before 截图 / A=审计记录 / F=整改 / T=after 截图）：

| # | 界面 | 入口 | B | A | F | T |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | AppShell 框架（顶栏 + 主侧栏 + 侧栏「场景项目」面板，含条目菜单与悬浮卡） | `http://127.0.0.1:5173/` | ☐ | ☐ | ☐ | ☐ |
| 2 | 沟通中心（秘书长对话，含新对话入口） | 侧栏 → 沟通 | ☐ | ☐ | ☐ | ☐ |
| 3 | 事务对象工作台（卡片列表 + 卡片封面 + 启动流程入口 + 详情/时间线） | 侧栏 → 事务对象 | ☐ | ☐ | ☐ | ☐ |
| 4 | 场景包管理页（pack 卡片 + 封面开关 + lifecycle 各态） | 侧栏 → 场景包管理 | ☐ | ☐ | ☐ | ☐ |
| 5 | 能力包管理 / 场景平台页（能力包卡片 + 封面开关） | 侧栏 → 对应入口 | ☐ | ☐ | ☐ | ☐ |
| 6 | 角色包管理页（角色包卡片 + 槽位绑定） | 侧栏 → 角色管理 | ☐ | ☐ | ☐ | ☐ |
| 7 | 场景制作台（含 env pack 规格加载与 draft→readiness→promote 各面） | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 8 | 能力制作台 | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 9 | 角色制作台 | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 10 | 执法行动 run 视图（流程推进 + Owner 确认面 + 回执时间线） | P4 路径内 | ☐ | ☐ | ☐ | ☐ |
| 11 | 高风险动作确认面（Owner 门弹层——只改呈现，主权要素不增删） | P4 路径内 | ☐ | ☐ | ☐ | ☐ |

- 打勾标准见 §7A.2；F 列填 client 仓整改 commit 短哈希，未整改的屏诚实留空并在报告说明原因。
- 多状态屏（空态/错误态/加载态）截图尽量齐；确无法复现的状态在审计记录注明「未覆盖」。
