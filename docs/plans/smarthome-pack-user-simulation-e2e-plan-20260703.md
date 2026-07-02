# 智能家居服务商 Pack 用户视角全链路实测与完善计划（2026-07-03）

> 状态：`已批准`（Owner 2026-07-02 口头授权：第一轮环保执法 pack 测试收口后**自动启动第二轮**，对象=智能家居服务商 pack；Owner 外出，登记后自主推进。授权原文见第一轮主任务登记表 `/Users/li/Documents/过程文档/envpack-live-20260702/logs/main-task-register.md` 的「Owner 裁定记录」节「第二轮授权」行）
> 计划文件：`/Users/li/Documents/truzhen-packs/docs/plans/smarthome-pack-user-simulation-e2e-plan-20260703.md`
> 衍生自：`/Users/li/Documents/truzhen-packs/docs/plans/env-pack-user-simulation-e2e-plan-20260702.md`（第一轮，结构/纪律完全沿用）
> 模型档位改判：分析/修复智能体 = **Opus 4.8 档（`model: "opus"`）**，随第一轮 Owner 口头改判（登记表已记，覆盖计划三稿的 Sonnet 档），非计划文件原文的 Sonnet。

## 0. 一句话目标

沿用第一轮三方结构（用户智能体走前端 UI / computer use + 组织者[我，前后端协调] + 监控程序 truzhen-monitor），以真实用户视角走通：三制作台使用 → 加载/验证智能家居服务商 pack → 在制作台修改完善它 → 用它创建一个客户项目经营行动 → 与项目经理（秘书长）互动完成项目经营闭环；过程中**发现问题→修问题**，同时完善 Truzhen 程序与智能家居 pack。

本轮同时承担两项 Owner 追加目标（沿用第一轮）：

1. **专项功能实测**：三荚卡片封面开关（场景包 / 能力包 / 角色包管理页）+ 事务对象卡片封面开关（前端）+ 场景项目菜单启动/暂停，逐项经 UI 实测并由组织者后端复核（细则见 §6A）。**第一轮已修复的项本轮改为「回归确认」而非首测**（见 §6A 标注与 §2A 结转基线）。
2. **截图驱动界面成品化**：测试沿途每个界面逐屏截图审计——去工程信息、前端只显示中文、排版布局优化到成品标准——整改后出前后对比截图（标准见 §7A，界面清单见 §12.9）。**第一轮成品化整改已 land 的组件本轮改为「回归验收 + 智能家居专属新文案审计」**。

## 1. 版本 / 优先级

- 定位：**当前主线**——收尾阶段「全链路真实连接」的用户视角验收测试第二轮，不是新功能开发。智能家居服务商 pack 属 Owner 2026-06-21 战略转向裁定的第一现金流方向（销售/CRM/私域/客服/售后等「高手经验可 Pack 化」垂直中的长周期项目交付型），用它做全链路实测正是第一现金流样板的用途——与第一轮环保执法（护城河标杆样板，非第一现金流）形成两类 pack 的方法论对照。
- 非目标：不做通用测试平台、不做新制作台功能、不建 Governance Console、不做通用 CRM Copilot。

## 2. 真实客户 / 场景证据

- 证据 = Owner 第一轮授权原话（2026-07-02，登记表已记）：本轮收口后**自动启动第二轮**，用 `truzhen-packs` 的智能家居服务商 pack 按同一方法论测试。
- 领域证据：pack manifest 声明 `target_industry=智能家居 / 智能装修 / 长周期项目交付`、`target_audience=智能家居公司老板、项目经理、客服与交付团队`；moat_justification 已声明护城河成立理由（真实系统集成/可审计证据回执/事务生命周期回放）。
- **本 pack 与环保 pack 的领域差异（决定剧本换新）**：无知识库（README 明示 install.py 跳过知识入库）、单角色（项目经理，无双角色对照）、Frappe provider 需求（项目/客户快照只读 + 里程碑写回）、长周期项目交付型 9 节点 flow。故用户剧本换成智能家居项目经营场景（见 §5.1A）。

## 2A. 第一轮结转基线（本轮开工前必读——已修则回归确认，未裁则继承为已知问题）

> 来源：第一轮主任务登记表 Land 记录 + ISSUE 台账（`/Users/li/Documents/过程文档/envpack-live-20260702/logs/main-task-register.md`）。本轮**不重测已 land 修复的存在性**，只做回归确认；橙色待裁项继承为本轮已知问题基线，不重复分析，遇到直接引用第一轮 ISSUE 号。

### 2A.1 第一轮已 land 修复（本轮回归确认，不首测）

| 仓 | 终态 commit | 已 land 内容（回归确认对象） |
| --- | --- | --- |
| truzhen-client-web-desktop | `044450d` | 批 1（ISSUE-001A/B、002a/b/c、003a/c 三台导航/表单治理/数据簇）+ 批 2（005a/006、005b、008、011 事务对象封面开关、012 场景项目菜单启动、）+ 批 3（成品化 10 commits）+ 批 4（013 枚举/014 卡片裁剪/015 折叠/016 落点/010 测试）+ 批 5（017b/c 草稿恢复门+工作区持久化+发布引导）+ 批 6（018 唯一幂等键+replayed 防假成功+可见确认卡、019 残留全修、回执屏范围文案） |
| truzhenos | `2ea6f59` | 批 4-后端：多斜杠 transaction_ref 子资源 405 修复（ISSUE-010，兼容加法：新增 transaction_ref query/body 变体不删旧路径） |
| truzhen-packs | `283849c`（第一轮 pack 完善，隔离 E2E registry 0→1→0 铁证 718b0baa） | 环保 pack：19 节点 stage_guide 一线口径 / manifest 对齐 / 15 知识域零漂移 / 测试报告增补 |

**本轮回归确认要点**（这些第一轮的修复应在智能家居 pack 上依然生效，若回潮即为回归 bug）：
- 三台步骤导航点击有反馈（001）；表单接收后端裁定、编辑清陈旧 readiness、脏态跟踪（002）；
- 事务对象卡片封面开关=纯前端展示偏好、诚实态不因开关消失（011）；场景项目菜单「启动」复用既有启动链、在途 run 诚实提示（012）；
- 荚卡开关 z-index 可点、开关在 footer、就近提示（005a/006/014/015）；
- 停用/启用回环：唯一幂等键、replayed 不冒充本次成功、可见确认卡（018）；启用重放 replayed 处理（013）；
- 成品化：枚举中文映射、ref 不裸排、工程话清理（批 3 + 019）。

### 2A.2 橙色/红色待 Owner 裁定项（继承为已知问题基线，本轮遇到直接引用不重复分析）

| 第一轮 ISSUE | 颜色 | 缺口摘要 | 本轮处理 |
| --- | --- | --- | --- |
| ISSUE-017a | 橙（跨仓/契约） | **已启用包编辑入口全链缺失**：后端无 `enabled→draft` 派生端点，flowId 只由模板 slug 派生 | 本 pack **已 enabled@1.0.0**（§3），P3「修改完善」正好首撞此缺口 → 引用 017a 方案，本轮只出/补方案不施工，不重新分析 |
| ISSUE-020 P4-09 | 橙红（主权越门） | **流程自动越 Owner 门**：`runtime.go:437` 只 park `human_approval`/`wait_event`，其余 candidate_emitter 自动 succeeded，从不读 `GatePolicy.PendingOwnerConfirmation`（死元数据） | 本 pack flow 的 `owner_gate` 节点 `pending_owner_confirmation:true`（§3）——P4 到达 owner_gate 时**重点验证是否复现越门**；复现即引用 020 方案，橙红只出方案 |
| ISSUE-020 P4-07 | 橙+黄 | 处置决定断链任务黑洞：candidate-input 回灌等值条件不成立不解锁；RiskLow task 不入 queue 四 tab 查无 | 本 pack `milestone_plan` 产任务候选——P4 验证任务是否入 queue/可见；复现引用 020 |
| ISSUE-020 P4-13 | 橙 | 对话历史跨会话丢失：前端 useState + 后端内存 map 零持久化 | P4 与项目经理多轮对话时验证；复现引用 020 |
| ISSUE-018（后端侧） | 橙 | 后端幂等回放状态校验缺口（前端侧已批 6 修，后端 reactivate 幂等重放旧结果不校验状态） | P2 回环若触发引用 018 后端方案 |
| ISSUE-003b | 黄（数据卫生） | smoke 残留测试污染正式库清理方案（只出方案未执行） | 若本 pack registry 见 smoke 污染引用 003b |
| ISSUE-022 | 黄（监控） | truzhen-monitor watch 重启复用 run_id → `(run_id,sequence)` UNIQUE 冲突刷屏 | **本轮 P0 起 monitor 必须换新 run_id**（教训见 §2B），避免复发；仍为 backlog |

## 2B. 第一轮方法论教训（本轮开工即遵守，防重踩）

1. **固定幂等键陷阱**（ISSUE-018）：前端固定幂等键 `...:${ref}@${ver}` 被历史成功消费后，后端幂等重放旧结果 `replayed=true`，前端不读 `replayed` 把回放当本次成功复读旧回执。**验收启用/停用/写回类动作，必须核对 receipt 号是新号且 `replayed=false`，不认「有回执号」= 真成功**。
2. **自证 fixture 禁令**（ISSUE-013）：批 2 单测曾用「同一套假枚举」的 fixture 自证通过，掩盖真值不匹配。**新增守卫测试的 fixture 必须用后端真枚举值（`disabled`/`pack_enabled_version` 等），不得自铸假枚举让测试与被测代码一起错**。
3. **涉布局必视觉自证**（ISSUE-014/015）：卡片固定高度 + overflow 把开关整行裁死，单测断结构断不出像素问题。**任何涉及布局/像素/可见性的修复，必须 vite:5174 + playwright 真截图视觉自证后才算完**，不认单测绿。
4. **detail 透传全链检查**（ISSUE-013/017c）：`packStudioCanvasJSON` 丢 `detail` 致「无原因报错」。**前端出错提示链路要全链核对 detail 是否从后端透传到用户可见文案**。
5. **CJK Edit 用 python**：`FEATURE_LEDGER.md` 等 CJK 长行用 Edit 易 unicode 匹配失败——**改用 python utf-8 读写插入**。
6. **cwd 重置用绝对 cd**：zsh cwd 在工具调用间会重置——**命令内一律用绝对路径 `cd`，不依赖上次 cwd**；`grep|head && echo` 退出码恒 0 假绿；无引号变量分词 for 失效改 while-read。
7. **背景化不套 `&`**：长命令 `run_in_background: true` 即可，命令里不得再加 `&`（双层后台化丢完成通知）；EGR 整体委托子代理后台跑，读真实 `VERIFY_EXIT` 值。
8. **写文件读回确认**：截图/文件写入后 read-back 确认，工具可能静默写失败报 success。

## 3. 现状基线（2026-07-03 实测）

- 智能家居 pack：`/Users/li/Documents/truzhen-packs/smart-home-owner-pack-v0/`
  - manifest：`manifest_version=v3`、`pack_ref=scene_pack://smart-home-owner-project-ops`、`version=1.0.0`、`pack_type=domain_work_pack`、`template_family=长周期项目交付型`。
  - **知识域数=0**（无 `knowledge/` 目录，README 明示 install.py 跳过知识入库；与环保 15 域相反）。
  - **单角色包**：`role_pack://smart-home-project-manager`（智能家居项目经理，`risk_level=medium`）；**无双角色对照**（环保是执法精英+挑剔律师）。role-slots 1 个（`smart_home_project_manager`，node_type=advice）。
  - flow 文件：`flows/smart-home-owner-project-ops-flow.flow.json`，`flow_id=smart-home-owner-project-ops-flow`，**9 节点**：intake（需求录入）→ project_object（项目事务建模）→ frappe_snapshot（读 Frappe 项目/客户快照，capability.invoke）→ manager_advice（项目经理建议，collaboration.advice，slot_ref=smart_home_project_manager）→ milestone_plan（里程碑/采购/施工任务候选）→ customer_draft（客户沟通草稿，gateway.communication_draft）→ frappe_write_candidate（Frappe 写回候选，capability.invoke）→ owner_gate（policy.gate_config，`pending_owner_confirmation:true`）→ receipt_archive → done。
  - **3 个 Frappe provider requirements**：`frappe_project_snapshot`（project_external_snapshot，risk low，fallback `provider_missing`）、`frappe_customer_snapshot`（customer_relationship_snapshot，risk low，fallback `provider_missing`）、`frappe_project_write_candidate`（project_milestone_write_candidate，risk medium，fallback `blocked`）。
  - install.py / uninstall.py 齐（install.py 复用第一轮通用 loader，`load_opt` 兼容无知识库 pack）。
  - **install.py 残留待清（P3 pack 完善候选）**：docstring 首行仍写「环保执法 Pack」（复制自环保未改）；knowledge batch `tags` 硬编码「环保执法」（本 pack 无知识库不会走到该行，但属残留）。
- **devserver 18080 registry 现状（实测，与第一轮相反）**：本 pack 已 **enabled@1.0.0**——`enabled_pointer.current_version=1.0.0`、`previous_version=""`、`default_version=1.0.0`、`occ_version=0`（从未 toggle 过）、`state=pack_enabled_version`、`receipt_ref=dd59131c-cca7-4f58-b602-db3d67a4fa70`、`pack_spec_hash=47a8e8f0b52fb3c4`。
  - **P2 语义因此翻转**：第一轮 P2 是「加载 disabled→enabled」；本轮 pack 已启用，P2 改为「**验证已启用态真实性 + 走停用→再启用回环（回归确认 018/013）+ 首撞已启用包编辑入口（ISSUE-017a 橙）**」。
- 运行环境：vite :5173 ✅（200）、devserver :18080 ✅（404=活着）、oMLX :8000 ✅（404=活着）——三服务在跑。
- client / truzhenos / packs main 基线：见 §12.2（继承第一轮，本轮开工需实测刷新）。
- 新 harness 目录：`/Users/li/Documents/过程文档/smarthomepack-live-20260703/`（尚不存在，P0 建，结构 `harness/` `shots/` `logs/` 同第一轮；不入任何仓）。

### 3.1 专项功能代码坐标（继承第一轮 §3.1 + 已修状态标注）

第一轮 §3.1 代码坐标本轮仍有效（`ScenePackManagePage.tsx` / `ScenePlatformPage.tsx` / `RolePackManagePage.tsx` / `BusinessObjectWorkbenchPage.tsx` / `AppShell.tsx` / sceneflowdev），**但多数已被第一轮批 1-6 修复**——本轮以「回归确认」为主，坐标以 client main `044450d` 为准（开工实测刷新行号）。仍需 UI 实测的专项项（§6A 逐条已标注回归/首测）：
- 三荚封面开关（场景/能力/角色）：第一轮已修 z-index/footer/就近提示（005a/006/014/015）→ 本轮**回归确认**，在智能家居 pack 卡片上复走开→关→开。
- 事务对象卡片封面开关：第一轮已实现（011，纯前端展示偏好）→ 本轮**回归确认**。
- 场景项目菜单「启动」：第一轮已接线复用既有启动链（012）→ 本轮**回归确认 + 智能家居项目验证在途 run 诚实提示**。
- 场景项目菜单「暂停」：第一轮暂停/恢复双路径已通（010 后端多斜杠 ref 修复 + 前端）→ 本轮**回归确认**；若本 pack 触发新暂停缺口再立新 ISSUE。
- 已启用包编辑入口（017a）：**本 pack 已 enabled，P3 首撞此橙色缺口**，本轮只出/补方案（继承 017a 分析，不重分析）。

## 4. 最小可交付

即使修复循环没跑满，本轮最少交付：
1. 智能家居 pack 已启用态经前端 UI 验证真实（有 receipt 可反查）+ 停用→再启用回环 UI PASS（回归确认 018/013）；
2. 一次客户项目经营行动 scene flow run 从需求录入推进到 Receipt 归档（与项目经理互动完成，Frappe 节点 provider_missing/blocked 诚实态）；
3. 问题清单（每项：现象+截图 → 根因 → 归属仓/风险颜色 → 修复或方案）；
4. 专项功能实测结论（§6A）：回归确认项 = 依然可用证明；新缺口 = 诚实缺口记录；
5. 测试主路径界面（至少 P2 / P4 沿途各屏）完成 §7A 成品化回归 + 智能家居专属新文案审计，前后对比截图归档，§12.9 清单打勾。

砍掉项 / backlog：三制作台的功能增强建议（只记录不施工）；多项目并行；移动端视角；非测试沿途界面的成品化；已启用包编辑入口后端实现（橙，017a，只出/补方案）；ISSUE-020 系列橙红项（只出方案）；install.py 残留清理（绿，pack 内容改，可 land）。

## 5. 角色结构与派活卡

### 5.0 模型档位、主任务登记与上下文隔离纪律（沿用第一轮，Owner 改判档位）

**模型档位（派子代理时显式指定，写进每张派活卡）：**

| 角色 | 模型档位 | 理由 |
| --- | --- | --- |
| 用户智能体（测试） | **最高档模型（Fable 5，`model` 缺省继承主会话即是）** | Owner 指令：测试智能体用最高模型，模拟用户全部行为 |
| 问题分析智能体 / 修复工 | **Opus 4.8 档（`model: "opus"`）** | Owner 第一轮口头改判（登记表已记，覆盖计划三稿 Sonnet）；分析/修复用 Opus |
| 组织者（主线） | Fable 5（主会话本体） | 只做派活/验收/登记/简报，不做长篇分析 |

- 档位改判理由：第一轮 Owner 口头把分析/修复档从 Sonnet 升到 Opus 4.8（登记表「分析/修复智能体档位」行）。本轮沿用 Opus 4.8，不再降回 Sonnet。

**主任务登记（两级，开工即建，随做随更）：**

1. 各仓 `FEATURE_LEDGER.md` §0 在途区登记本测试任务（开工登记、收尾按 §12.8 转正/销账）；
2. 运行期主任务登记表 `logs/main-task-register.md`（**新 harness 目录 `/Users/li/Documents/过程文档/smarthomepack-live-20260703/logs/`**）：逐行记录阶段状态（P0-P7 各档位）、问题编号（ISSUE-NNN，本轮编号从 SH-001 起，避免与第一轮 ISSUE 号混淆）、派出的智能体、返回结论摘要、修复 commit 或方案路径、land 状态。每次派活/收活即更新，主线恢复上下文以此表为准。

**上下文隔离（防主任务上下文污染）：** 同第一轮——主线不在自己上下文做问题分析；每遇问题先立 SH-NNN，单独派分析智能体（Opus 档，独立上下文包）；返回格式固定 `根因 / 归属仓 / 风险颜色 / 修复 commit 或方案文件绝对路径 / 验证证据`；分析长文落该智能体自己的 markdown（`logs/issues/SH-NNN.md`）；用户智能体不因某问题停摆，登记外派后绕开继续主旅程。

### 5.1 用户智能体（User Agent，子代理，长期会话可续）

- **模型档位**：最高档模型（Fable 5，§5.0）。
- **只允许**：通过 Playwright 驱动真 Chromium 操作 `http://127.0.0.1:5173`，每步截图；用 Read 看截图理解界面（computer-use 循环：截图→理解→操作）。
- **禁止**：curl / fetch 后端、读产品源码、读数据库、用开发者知识绕过 UI。UI 走不通 = 记录产品问题，不许绕。
- **行为全覆盖**（沿用第一轮 §5.1）：主路径沿途每个界面，对每个可见交互控件至少各触发一次真实用户行为（点击/悬浮/文本输入含空/超长/中文/下拉展开收起/滚动/拖拽/键盘 Tab-Enter-Esc/刷新复看/返回切页再进）；每屏产出行为覆盖矩阵（控件 × 行为 × 实际结果 × 截图路径），未触发控件诚实标注原因；误操作路径也模拟（取消/重复点击/必填漏填）。
- 选择器纪律：优先用户可见文本 / role；仅可见文本无法唯一定位时才用 data-testid，且每次记为一条潜在 UX 发现。
- 产出：每阶段一份用户日志。
- Harness 与截图目录（稳定路径，跨会话可复用，不入任何仓）：`/Users/li/Documents/过程文档/smarthomepack-live-20260703/`（`harness/` playwright 脚本，`shots/` 截图，`logs/` 用户日志）。

### 5.1A 用户人设与场景剧本（本轮换新，贴智能家居领域）

- **用户人设**：智能家居/智能装修公司老板（兼工程主管），公司 3-50 人（Owner 第一现金流目标客群），第一次用 Truzhen 管一个正在交付中的客户项目；不懂后端术语，只认「客户/项目/里程碑/采购/施工/验收/售后」的行业语言。
- **场景剧本（智能家居服务单/项目经营）**：某客户「全屋智能安装」项目已进入交付期，客户报修「客厅灯组离线 + 中控网关失联」——老板要把这单当一个**客户项目经营事务**录入：录入项目机会/客户需求（intake）→ 建项目事务对象（project_object）→ 读 Frappe 项目/客户快照（frappe_snapshot，provider 未接通预期 provider_missing 诚实态）→ 项目经理给下一步经营建议（manager_advice，本地 oMLX 出候选）→ 生成里程碑/采购/施工任务候选（milestone_plan，如「上门排查网关」「补货灯组模块」「约客户复验时间」）→ 客户沟通草稿候选（customer_draft，如给客户的上门时间确认话术）→ Frappe 写回候选（frappe_write_candidate，预期 blocked 诚实态）→ Owner/Base 经营确认门（owner_gate）→ 回执归档（receipt_archive）。
- **与环保剧本的差异**：环保是「污水厂偷排执法办案」（双角色对照 执法精英×挑剔律师、15 知识域法条、文书送达）；本轮是「客户项目交付经营」（单角色项目经理、无知识库、Frappe 快照/写回、无对照节点）。因此本轮 **P4 不含 lawyer_challenge 对照节点**（第一轮 P4 卡在 lawyer_challenge model_not_ready 的问题本轮不复现），**新增验证点 = Frappe provider 三态诚实性（provider_missing / blocked 不假成功）**。

### 5.2 组织者 / 协调者（我，主线）

- 派活、验收用户日志、前后端协调检查（后端日志、SQLite、03 回执反查、监控事件）、问题分诊、修复派工、阶段简报、主任务登记表维护（§5.0）。
- 对每个用户报告的问题先独立复核（不轻信自报），再定根因与归属——复核与根因分析本身也外派分析智能体（Opus 档），主线只读结论回填登记表。

### 5.3 问题分析智能体 + 修复工（按问题派，Opus 4.8 档）

- 每个 SH-NNN 先派分析智能体（Opus 档，独立上下文包）定根因/归属/颜色；确需施工再派修复工。
- 每个修复在归属仓独立 worktree 分支，TDD（先加失败测试再修），改哪证哪；验收由组织者按 §10 复核，不采信修复工自报。
- **遇第一轮橙色继承项（017a/020 系列/018 后端）直接引用第一轮方案，不重新分析**——分析智能体只在新缺口上花上下文。

### 5.4 监控程序

- 起 `truzhen-monitor`（`/Users/li/Documents/truzhenos/backend/cmd/truzhen-monitor`），全程收集 `/v3/monitoring/*` 事件；测试收尾导出诊断包。**开工即换新 run_id，避免 ISSUE-022 复用冲突**（§2B 第 1 教训的监控对应项）。发现监控缺口 → 补回现有体系。

## 6. 阶段设计

| 阶段 | 视角 | 内容 | 验收 |
| --- | --- | --- | --- |
| P0 环境基线 | 组织者 | 确认三服务、起 truzhen-monitor（换新 run_id）、建新 harness 目录、搭 Playwright harness、user agent 首页可达冒烟 | 首页截图 + monitor 新 run 收到事件 |
| P1 三制作台巡检 | 用户 | 场景/能力/角色制作台各走最小用户旅程；**回归确认第一轮批 1-6 三台修复未回潮** | 三份用户日志 + 回归确认 + 新问题记录 |
| P2 加载/验证 pack | 用户 | 本 pack 已 enabled@1.0.0——UI 验证已启用态真实性；走停用→再启用回环（**回归确认 018 唯一幂等键/013 枚举，核 receipt 新号且 replayed=false**）；已启用包编辑入口首撞（017a） | enabled 真实性 + 回环 receipt 可 03 反查（新号）+ 017a 缺口记录 |
| P3 修改完善 pack | 用户+组织者 | 场景制作台加载智能家居 pack 规格，做真实完善（候选见 §7）；已启用包编辑入口=橙（017a）只出/补方案；install.py 残留清理（绿，可 land） | pack 仓文件更新（install.py 残留修）+ 017a 方案 + 隔离 E2E |
| P4 项目经营行动 | 用户 | 用 pack 创建客户项目经营行动（scene flow run）：与项目经理对话推进 intake→project_object→frappe_snapshot（provider_missing 诚实）→manager_advice→milestone_plan→customer_draft→frappe_write_candidate（blocked 诚实）→owner_gate→receipt_archive；**验证是否复现 P4-09 越门 / P4-07 任务黑洞 / P4-13 对话丢失** | run 全节点 evidence + 最终 receipt 反查 + Frappe 三态诚实 + 越门复现结论 |
| P4A 专项功能实测 | 用户+组织者 | §6A 六条：三荚封面开关（回归）、事务对象封面开关（回归）、场景项目菜单启动/暂停（回归）；每项 UI 实测 + 组织者后端复核，新缺口按 §6A 判定 | 每项结论 = 回归可用 或 新缺口闭环记录 |
| P4B 界面成品化 | 用户+组织者+修复工 | 对 P0-P4A 沿途每个界面按 §7A 三标准截图审计——第一轮已整改组件回归验收 + 智能家居专属新文案（项目经理/里程碑/Frappe/项目术语）审计整改 | §12.9 界面清单全屏过审 |
| P5 验收对账 | 组织者 | 回执 03 反查、监控事件对账、诊断包导出、用户智能体行为审计 | 对账清单全绿或诚实标注 |
| P6 修复循环 | 贯穿 | 每问题：登记 SH-NNN → 外派分析（Opus 档）定根因 → 分诊：绿=修+自测；黄=修后简报确认；橙/红=只出方案；继承第一轮橙项直接引用不重分析 | 每问题在登记表闭环 |
| P7 收尾 | 组织者 | 测试报告、pack 完善提交、FEATURE_LEDGER 登记、land（按 §9 裁定） | 报告文件 + 账本更新 |

## 6A. 专项功能实测细则（沿用第一轮 §6A，已修项改回归确认）

原则：逐项「先实测存在性 → 再测行为 → 组织者后端复核」；禁止把「代码里有类名/端点」当「功能可用」。第一轮已修项本轮重点是**在智能家居 pack 卡片/项目上确认修复未回潮**。

1. **场景包卡片封面开关**（场景包管理页）：对智能家居 pack 卡片执行 开→关→开 各一次；组织者复核 lifecycle 端点状态真实翻转 + receipt 03 反查（**核新号、replayed=false，回归 018**）；停用后卡片仍保留列表。第一轮已修 z-index/footer（005a/014）→ **回归确认**。
2. **能力包卡片封面开关**：任选一个已装能力包同上实测（回归确认 006）；若无已装能力包先经 UI 装一个。
3. **角色包卡片封面开关**：智能家居 pack 单角色（项目经理）——实测角色包管理页卡片。第一轮已修「当前生效区与开关状态矛盾」（009，开关改读绑定轴）→ **回归确认**；语义仍按「角色包非权限开关」的产品口径（页内文案讲清）。
4. **事务对象卡片封面开关**（前端）：第一轮已实现纯前端展示偏好（011）→ **回归确认**：在智能家居项目事务卡片上验证封面展开/收起，诚实态（「未启动流程/流程状态不可用」）不因开关消失、不产 receipt、不伪装启停。
5. **场景项目菜单「启动」**：第一轮已接线复用既有启动链（012）→ **回归确认**：在智能家居项目侧栏条目菜单启动一个项目经营 run，验证复用既有启动链（不铸第二链）、在途 run 诚实提示。
6. **场景项目菜单「暂停」**：第一轮暂停/恢复双路径已通（010 后端多斜杠 ref 修复）→ **回归确认**：暂停一个智能家居项目 run 并恢复，验证三证主权链（owner_ref/gate_decision_ref/evidence_refs）；若本 pack 触发新暂停缺口（如 ISSUE-019 P2F-05 状态漂移「待启动→进行中」）再立新 SH-NNN。

## 7. Pack 完善候选（P3 具体改什么，按测试发现动态增减）

1. **install.py 残留清理（绿，可 land）**：docstring 首行「环保执法 Pack」→ 改「智能家居老板项目经营 Pack」；knowledge batch `tags` 硬编码「环保执法」→ 改通用或按 pack 声明取值（本 pack 无知识库虽不走该行，仍清残留防复制误导）。
2. 把 install.py 里散落的门槛信息（Frappe provider/model 要求）对齐进 manifest 声明；
3. 流程节点用户可读提示语完善（9 节点 stage guide 口径与前端一致，智能家居项目术语）；
4. 单角色（项目经理）role-pack steering 提示按实测对话质量调优（communication_style forbidden_phrases 已声明，实测是否生效）；
5. 补 pack `docs/` 目录与测试报告（本 pack 现无 docs/ 目录，README 已有）；测试报告写入本轮真实结果。

## 7A. 界面成品化标准与流程（沿用第一轮 §7A，全文照搬）

### 7A.1 三项标准（每屏逐条过）

1. **去掉工程信息**——用户可见区域不得出现：内部 ref 原文直排（`scene_pack://…`、`transaction://…`、`receipt://03/…`、`run-gated-…` 等 URI/句柄裸露正文）；模块编号口径（「05 事务对象 ReadModel」「06 场景流程」「03 回执账本」等内部编号叙述）；英文状态枚举原文（running/blocked/candidate_only/provider_missing…）；开发批次/任务代号、data-testid 式文案、面向开发者的调试性长句。改法=翻译成用户语言或折叠进「详情」二级展示。**主权红线（只能改呈现，不得删除）**：高风险动作确认面展示 目标/内容/影响/run_id/nonce 是受控真实 E2E 的主权要求；回执/凭证 ref 必须保持可反查。整改只允许成品化呈现（折叠「凭证详情」、等宽字体、一键复制），禁止删除、禁止隐藏到不可达；确认面要素**增删**一律先出方案不施工（红）。
2. **前端只显示中文**：所有用户可见文案中文；英文仅限专有名词（Truzhen、Frappe、pack 英文名等 Owner 认可者）与必须原样保留的错误原文（收进「详情」）。状态枚举走中文映射，补齐漏网值，未知枚举兜底中文（「状态：未知（原文见详情）」），不直排英文原文。
3. **成品排版布局**：对齐既有主题 token 体系（R1 主题批次成果），不另起平行样式；卡片无内部小滚动条、滚轮不滚单卡（C9 口径回归）；同屏间距/字号/层级一致；空态/错误态/加载态三态文案完整且中文；长文本截断 + 悬浮/详情可见全文；散落 inline style 错位改用主题 class 收敛。

### 7A.2 流程（每屏一循环，截图是唯一驱动）

1. 用户智能体进入该屏即存基线截图 `shots/audit/<screen>-before.png`（多状态屏各存一张）；
2. 组织者按 7A.1 逐条出该屏审计记录，落 `logs/ui-audit-<screen>.md`；
3. 修复工在 client 仓整改（绿/黄；纯文案改动可攒批，布局改动按屏拆 commit）；
4. 复截图 `shots/audit/<screen>-after.png`，前后对比确认；
5. §12.9 界面清单该屏打勾。发现某屏问题属后端 → 按归属仓走修复循环，不前端硬编码遮丑。

### 7A.3 防回潮与边界

- 每改完一屏随手加守卫：对已成品化组件补「用户可见文案不得出现英文枚举原文 / 内部 ref 直排」断言，沿用 `frontendStrict.test.tsx` 真实性测试口径（真渲染真断言，不 mock fetch 编假数据；守卫 fixture 用真值，§2B 第 2 教训）。
- 只收测试沿途界面（§12.9 为准，按实际路由增删）；非沿途屏问题记 backlog。
- 成品化不改行为：整改不得改变任何主权链路、端点调用与门控行为；行为缺陷走 §6 修复循环而非借成品化夹带。

## 8. 治理维度

- **真相源**：pack 事实 = truzhen-packs 文件夹 + devserver registry；运行态/流程 = truzhenos SQLite；正式动作事实 = 03 回执账本；监控 = truzhen-monitor。前端只是投影，用户日志只是观察记录。
- **仓库归属**：pack 完善→`truzhen-packs`；前端问题→`truzhen-client-web-desktop`；后端/网关/流程问题→`truzhenos`；契约问题→`truzhen-contracts`（**只出方案不施工**）。本计划即为跨仓授权申请，Owner 批准计划（第一轮授权自动延续本轮）= 授权在上述四仓（contracts 仅读）作业。
- **风险颜色**：整体黄。前端/文档/pack 内容=绿~黄；devserver 业务逻辑=黄；网关/主权闸/契约=橙红（只出方案）。专项与成品化细分同第一轮。**本轮特有橙项**：已启用包编辑入口（017a 后端派生端点缺失）=橙只出/补方案；owner_gate 越门（020 P4-09）若复现=橙红只出方案；Frappe write 候选走 blocked=预期诚实态非缺陷。
- **契约影响**：默认零契约改动。若测试暴露 contracts 缺口（如 017a 可能涉 lifecycle 契约），单独出影响清单待裁。
- **上下文边界**：允许读四仓代码/文档、devserver 日志与本地测试库；用户智能体上下文仅限 UI 与自己的截图。
- **禁止边界**：
  - 不做任何真实对外发送/真实 Frappe 写回（`frappe_project_write_candidate` 保持 blocked、customer_draft 不真发，诚实态）；
  - 不碰真实客户数据、不碰他会话 WIP；
  - pack 的 `_source-materials`（若有）只读；
  - 不自动 merge 橙/红改动；不删任何 worktree/分支。
- **生命周期档位**：本计划=`已批准`；执行后测试结论=`已验收`（有 receipt/截图/对账证据）；修复项逐项走 `已实现→已接线→已验收`。

## 9. 待 Owner 裁定项（沿用第一轮默认值，无需再问，列出即可开工）

> Owner 第一轮已按 §12.1 建议值裁定并授权自动启动第二轮；本轮全部沿用第一轮默认值，无需再问，开工即按下列执行：

1. **修复 land 政策**：绿色修复攒批、每阶段简报后由组织者统一 land 对应仓 main；黄色修复简报确认后 land；橙/红只出方案。（第一轮默认值，沿用）
2. **秘书长/项目经理对话模型**：本地 oMLX（Qwen3.5 系）优先、零成本可长跑；若对话质量不足以推进项目经营节点，报告后再议切云模型（ark key 待轮换，需 Owner 点头）。（第一轮默认值，沿用）
3. **Owner 门代点授权**：项目经营行动推进会遇 Owner 主权确认点（owner_gate）。本机测试数据范围内，授权用户智能体代点 Owner 确认（构成测试态 OwnerActionEvidence，全部留 receipt），任何真实对外动作/Frappe 写回保持 blocked。（第一轮默认值，沿用）
4. **节奏**：无人值守跑到底，每阶段出简报但不阻塞继续（Owner 外出，登记后自主推进）。（第一轮默认值 + 本轮授权，沿用）
5. **已启用包编辑入口缺口处理（本轮首撞 017a）**：橙色，本轮只出/补方案（继承第一轮 017a 分析），前端不做假编辑入口。（对齐第一轮橙项处理默认）
6. **owner_gate 越门（020 P4-09）若复现**：橙红，只出方案（继承第一轮 020 分析）。（对齐第一轮橙红项处理默认）

## 10. 验收设计（严苛复核，非自报，沿用第一轮 §10）

- 用户智能体自报不算数：P2/P3/P4 每个关键声明由组织者用后端事实复核（registry 状态、SQLite 行、03 回执原文反查、监控事件流）。
- **反伪造 + 反幂等假成功检查**（强化自第一轮 ISSUE-018）：receipt_ref 必须能在 03 反查出原文；启用/停用/写回类动作核对 **receipt 号是新号且 `replayed=false`**，不认「有回执号」=真成功；scene flow run 节点 evidence 必须真实存在；禁止 candidate 冒充 formal。
- 用户智能体行为审计：抽查 harness 脚本与 transcript，确认无后端直连。
- 行为全覆盖验收：每屏行为覆盖矩阵齐全，无「代码里有但从未被触发」的沿途控件；未触发项诚实原因标注。
- 主任务登记验收：`logs/main-task-register.md` 与实际派活/收活一致（抽查 3 个 SH-NNN），FEATURE_LEDGER §0 登记齐。
- 专项功能验收（P4A）：回归确认项须证明修复未回潮（截图 + 后端事实）；新缺口须闭环记录（归属仓 + 风险颜色 + 修复 commit 或方案路径）。
- 界面成品化验收（P4B）：§12.9 清单三件套齐；after 截图逐条对照 §7A.1；新增文案守卫测试为真渲染真断言（fixture 用真值）；抽 3 屏由 Owner 亲眼比对。
- **回归验收（强化自第一轮 §2B 第 3 教训）**：改哪证哪——client 改动跑 typecheck+test+smoke，**涉布局/像素/可见性的修复必须 playwright 视觉自证**；truzhenos 改动跑 `go test` + EGR `bash scripts/verify.sh`（子代理后台跑，读真实 VERIFY_EXIT）；packs 改动跑隔离 devserver install/uninstall E2E（registry 0→1→0 铁证）。
- Owner 验收：刷新 :5173 亲手复走 P2/P4 关键路径 + 亲手点一遍 §6A 六项 + 读最终测试报告。

## 11. 变更影响预览

- 会碰：client UI（回归确认 + 新 bug 修 + §7A 智能家居专属新文案成品化整改）、truzhenos devserver（新 bug 修）、truzhen-packs 智能家居 pack（install.py 残留清理 + stage guide + docs/ 测试报告）、各仓 FEATURE_LEDGER、监控事件流（新测试期事件）。
- 不碰：contracts schema、主权闸/Base Gate 逻辑（发现问题只出方案）、高风险确认面主权要素集合（run_id/nonce/目标/影响，只改呈现不增删）、已启用包编辑入口后端实现（017a 只出方案）、owner_gate 越门后端实现（020 只出方案）、Frappe 真实写回（保持 blocked）、支付/市场、他会话 WIP、生产凭据。

---

## 12. 新会话无缝接手背景（执行会话开工必读）

> 本节写给执行本计划的新对话：读完本节 + §2A 结转基线 + §9 默认值即可直接开工，无需回溯旧会话。

### 12.1 Owner 裁定回填区（第一轮已裁，本轮沿用，无需再问）

| §9 裁定项 | Owner 裁定（第一轮已裁，沿用） |
| --- | --- |
| 1. 修复 land 政策 | 绿/黄阶段简报后统一 land，橙/红只出方案 |
| 2. 对话模型 | 本地 oMLX 优先，不足再议云模型 |
| 3. Owner 门代点授权 | 测试数据内代点，真实对外动作/Frappe 写回保持 blocked |
| 4. 节奏 | 无人值守跑到底，每阶段简报不阻塞（Owner 外出自主推进） |
| 5. 已启用包编辑入口（017a） | 橙，只出/补方案不施工 |
| 6. owner_gate 越门（020 P4-09）若复现 | 橙红，只出方案 |
| 分析/修复模型档 | **Opus 4.8（`model: "opus"`）**，第一轮口头改判 |
| 测试模型档 | 最高档（Fable 5） |

### 12.2 六仓布局与基线（继承第一轮，开工需实测刷新 main 基线）

| 仓 | 路径 | 第一轮终态 main | 本计划角色 |
| --- | --- | --- | --- |
| truzhenos（基座） | `/Users/li/Documents/truzhenos` | `2ea6f59`（第一轮 land；主检出可能仍有他会话遗留 merge 冲突，见第一轮登记表「环境事实」，未解决即隔离 worktree 编译跑） | 后端问题修复；devserver/monitor 宿主 |
| truzhen-packs | `/Users/li/Documents/truzhen-packs` | `283849c`（第一轮 pack 完善 land） | 智能家居 pack 完善（主交付） |
| truzhen-client-web-desktop | `/Users/li/Documents/truzhen-client-web-desktop` | `044450d`（第一轮批 1-6 land） | 前端问题修复；用户智能体操作对象 |
| truzhen-contracts | `/Users/li/Documents/truzhen-contracts` | — | **只读**；缺口只出方案 |
| truzhenv3（旧主仓） | `/Users/li/Documents/truzhenv3` | 已封棺 | 只读历史参考 |

- **开工第一步实测刷新**：`git -C <仓> rev-parse HEAD` 确认三仓真实 main（第一轮 land 后可能有他会话新提交）；client 主仓 checkout 可能停他会话分支，vite 从他会话 worktree 起（只运行不修改），同步 main 用 `git -C <worktree> merge --ff-only origin/main`。
- **任何修复开发一律新建独立 worktree 到 `/Users/li/Documents/truzhenv3worktree/`（truzhenos）或各仓约定位置，绝不在主仓目录直接开发。**

### 12.3 运行环境（先探活再起，三服务此前均在跑）

```sh
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:18080/   # devserver（404=活着）
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:5173/    # vite（200=活着）
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8000/    # oMLX（404=活着）
# 若没起：
cd /Users/li/Documents/truzhenos && go run ./backend/cmd/devserver          # :18080
cd /Users/li/Documents/truzhen-client-web-desktop && npm run dev -- --host 127.0.0.1   # :5173（必须 --host 127.0.0.1）
# oMLX 由 Owner 侧常驻（:8000）；没起时对话会失败——先报告不得假成功
# 监控：/Users/li/Documents/truzhenos/backend/cmd/truzhen-monitor（P0 起，换新 run_id 避免 ISSUE-022，统一走 /v3/monitoring/*）
```

### 12.4 智能家居 pack 事实基线（2026-07-03 实测，执行会话不必重查）

- Pack 目录：`/Users/li/Documents/truzhen-packs/smart-home-owner-pack-v0/`
  - `pack_ref=scene_pack://smart-home-owner-project-ops`、`version=1.0.0`、`pack_id=pack_smart_home_owner_v0`、manifest v3、`template_family=长周期项目交付型`。
  - **无知识库**（无 `knowledge/` 目录，install.py 自动跳过）；**单角色**（`role_pack://smart-home-project-manager` 项目经理，risk medium）；role-slots 1 个（advice）。
  - flow：`flows/smart-home-owner-project-ops-flow.flow.json`，9 节点（intake→project_object→frappe_snapshot→manager_advice→milestone_plan→customer_draft→frappe_write_candidate→owner_gate[pending_owner_confirmation:true]→receipt_archive→done）。
  - 3 Frappe provider requirements：project_snapshot（provider_missing）、customer_snapshot（provider_missing）、project_write_candidate（blocked）。
  - install.py/uninstall.py 齐（复用第一轮通用 loader，`load_opt` 兼容无知识库）；**install.py docstring 首行/tags 有「环保执法」残留待清（§7 第 1 项）**。
- **devserver registry 现状（实测）**：本 pack 已 **enabled@1.0.0**——`current_version=1.0.0`、`previous_version=""`、`occ_version=0`、`state=pack_enabled_version`、`receipt_ref=dd59131c-cca7-4f58-b602-db3d67a4fa70`、`pack_spec_hash=47a8e8f0b52fb3c4`。**P2 语义翻转**：不是「加载」，而是「验证已启用态 + 停用→再启用回环 + 首撞已启用包编辑入口 017a」。
- 关键后端端点（组织者复核用，用户智能体禁用）：
  - `GET /v3/pack-studio/lifecycle/packs?pack_ref=scene_pack://smart-home-owner-project-ops`（pack 状态）
  - lifecycle：`POST /v3/pack-studio/lifecycle/{draft|readiness|promote|confirm|reactivate}`
  - 角色包：`GET /v3/agent-orchestration/role-packs/readmodel`、agent-slots 同前缀
  - Frappe provider：走 04/11 execution 网关候选，未接通返回 provider_missing/blocked
  - install.py 用法：`TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:<port> python3 install.py`（隔离验证指向独立端口隔离 devserver，勿对 18080 重装除非计划明确）
- 项目经理对话闭环：本地 oMLX（:8000）出 advice 候选；前端入口 = 沟通中心「新对话」或事务对象工作台启动流程。
- 专项功能代码坐标见第一轮 §3.1（多数已修，本轮回归确认），本轮 §3.1 已标注每项回归/首测状态。

### 12.5 Computer-use harness（用户智能体）

- **新稳定目录**：`/Users/li/Documents/过程文档/smarthomepack-live-20260703/`（`harness/`、`shots/`、`logs/`；成品化审计截图放 `shots/audit/`，每屏 `-before.png`/`-after.png` 成对，审计记录 `logs/ui-audit-<screen>.md`；主任务登记表 `logs/main-task-register.md`、问题分析长文 `logs/issues/SH-NNN.md`、行为覆盖矩阵 `logs/behavior-matrix-<screen>.md`）。P0 建目录。
- 模型档位纪律见 §5.0：用户智能体=最高档（Fable 5）；分析/修复=Opus 4.8（派活时显式 `model: "opus"`）。
- 初始化：`cd harness && npm init -y && npm i playwright`（chromium 本机已有缓存，缺则 `npx playwright install chromium`）；可复用第一轮 harness 脚本（`/Users/li/Documents/过程文档/envpack-live-20260702/harness/`）作模板。
- 循环：脚本操作 `http://127.0.0.1:5173` → 每步截图到 `shots/` → 智能体 Read 截图理解 → 决定下一步。选择器纪律见 §5.1。
- 用户智能体建议用可续会话（组织者记住 agent 名，跨阶段 SendMessage 续用）。

### 12.6 验证命令映射（改哪证哪）

- client：`npm run typecheck && npm run test && npm run build && npm run smoke:frontend-shell`；**涉布局/像素修复加 playwright 视觉自证（§2B 第 3）**。
- truzhenos：`go test ./... && go test -race ./...`（范围见 `scripts/go-test-packages.txt`）；land 前 EGR `bash scripts/verify.sh`——**整体委托子代理后台跑（约 24 分钟），日志落固定文件，读真实 `VERIFY_EXIT`，禁双层后台化**。
- packs：隔离 devserver 起独立端口 → install.py/uninstall.py 全程 E2E → registry 0→1→0 铁证。

### 12.7 已知坑速查（合并第一轮 §2B）

- zsh：不做无引号变量分词（for 失效改 while-read）；`echo ===`/裸 glob 会炸；`grep|head && echo` 退出码恒 0 假绿。
- 长命令 `run_in_background: true` 即可，命令里不得再加 `&`。
- 前端截图/文件写入后 read-back 确认（工具可能静默写失败报 success）。
- 对话回复失败先查 oMLX(:8000) 与 devserver 是否都活，再查 08 网关。
- devserver 可能由别的 checkout 起着——复核 `/v3/` 行为异常时先确认它跑的是哪个工作树代码。
- FEATURE_LEDGER CJK 长行用 python utf-8 读写插入，不用 Edit。
- **固定幂等键陷阱/自证 fixture 禁令/涉布局必视觉自证/detail 透传全链**（§2B 1-4 教训）。
- cwd 重置——命令内一律用绝对路径 cd。

### 12.8 收尾要求

- 测试报告写入 `/Users/li/Documents/truzhen-packs/docs/reports/`（本计划同名前缀 `smarthome-pack-user-simulation-e2e-*`），pack 的 `docs/测试报告.md`（本 pack 现无 docs/，P3 建）同步；专项功能六项回归结论、§12.9 界面清单终态、主任务登记表终态（全 SH-NNN 闭环）并入报告。
- packs / truzhenos / client 各自 FEATURE_LEDGER 开工登记 §0 在途、收尾登记功能域行。
- 对话简报 ≤500 字：结论、报告绝对路径、关键验证、待 Owner 决策项。

### 12.9 界面成品化清单（P4B 逐屏打勾；按执行时实际路由增删，增删须留痕）

测试主路径预计经过的界面（智能家居版，B=before 截图 / A=审计记录 / F=整改 / T=after 截图；第一轮已整改屏本轮=回归验收 + 智能家居专属新文案）：

| # | 界面 | 入口 | B | A | F | T |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | AppShell 框架（顶栏 + 主侧栏 + 侧栏「场景项目」面板，含条目菜单与悬浮卡） | `http://127.0.0.1:5173/` | ☐ | ☐ | ☐ | ☐ |
| 2 | 沟通中心（项目经理对话，含新对话入口） | 侧栏 → 沟通 | ☐ | ☐ | ☐ | ☐ |
| 3 | 事务对象工作台（智能家居项目卡片列表 + 卡片封面 + 启动流程入口 + 详情/时间线） | 侧栏 → 事务对象 | ☐ | ☐ | ☐ | ☐ |
| 4 | 场景包管理页（智能家居 pack 卡片 + 封面开关 + lifecycle 各态） | 侧栏 → 场景包管理 | ☐ | ☐ | ☐ | ☐ |
| 5 | 能力包管理 / 场景平台页（能力包卡片 + 封面开关） | 侧栏 → 对应入口 | ☐ | ☐ | ☐ | ☐ |
| 6 | 角色包管理页（智能家居项目经理角色包卡片 + 槽位绑定） | 侧栏 → 角色管理 | ☐ | ☐ | ☐ | ☐ |
| 7 | 场景制作台（含智能家居 pack 规格加载与 draft→readiness→promote 各面 + 已启用包编辑入口 017a） | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 8 | 能力制作台 | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 9 | 角色制作台 | 侧栏 → 制作台 | ☐ | ☐ | ☐ | ☐ |
| 10 | 项目经营 run 视图（流程推进 + Frappe provider 诚实态 + Owner 确认面 + 回执时间线） | P4 路径内 | ☐ | ☐ | ☐ | ☐ |
| 11 | 高风险动作确认面（Owner 门弹层——只改呈现，主权要素不增删） | P4 路径内 | ☐ | ☐ | ☐ | ☐ |

- 打勾标准见 §7A.2；F 列填 client 仓整改 commit 短哈希，未整改屏诚实留空并说明原因。
- 多状态屏（空态/错误态/加载态，含 Frappe provider_missing/blocked 诚实态）截图尽量齐；确无法复现的状态注明「未覆盖」。
