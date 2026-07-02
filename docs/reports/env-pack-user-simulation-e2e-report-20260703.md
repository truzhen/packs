# 环保执法 Pack 用户视角全链路实测最终报告（2026-07-02～03）

> 生命周期档位：**本轮测试 = 已验收**（证据在案，见 §3/§4）。
> 计划：`/Users/li/Documents/truzhen-packs/docs/plans/env-pack-user-simulation-e2e-plan-20260702.md`（验收口径 §10 / 界面清单 §12.9）。
> 权威事实源：`/Users/li/Documents/过程文档/envpack-live-20260702/logs/main-task-register.md`（Owner 裁定 / 阶段状态 / ISSUE 台账 001-022 / Land 记录）。
> 证据目录：同上 `logs/`（各阶段 user-log、behavior-matrix、ui-audit、issues/ISSUE-*）与 `shots/`（含 `shots/audit/` 10 屏 after）。

---

## 1. 执行摘要

- 以「用户智能体（只走前端 UI / computer use）+ 组织者（前后端协调）+ 监控程序（truzhen-monitor）」三方结构，**一日之内跑完 P0-P7**（含专项 P4A、成品化 P4B），真实用户视角走通三制作台巡检 → 环保执法 pack 加载 → 制作台完善 → 发起执法行动 → 与秘书长互动至对外送达。
- **主路径闭环结论：闭上了。** P2 加载 pack 达成 UI 全回环 PASS（停用确认卡 → 真停用 → 启用确认卡 → 真启用，新回执 `608291f0` ≠ 旧号 `57b65f82` 非复读 → 刷新持久 → 终态 `enabled@1.0.2`；误点探针包被主权闸受控驳回、零副作用）。
- **执法行动 run** 从线索受理推进 7 阶段（受理 / 立案[PDF 真解析]/取证 / 三性 / 推演 / 文书草稿），run 最远推进至 `lawyer_challenge`；文书对外送达命中 `provider_missing` 受控拒绝、无假成功、无真送达——**fail-closed 方向正确**。
- **修复 land 统计**：client main `6cd6084 → a9a4c87 → 29670e0 → 044450d → 8dd8de3`（批 1-6、8 共 37 commits，最终 730 测试全绿）；truzhenos main `c8c3702 → 2ea6f59`（多斜杠 ref 405 修复，EGR VERIFY_EXIT=0）；packs main `caa034d → 283849c`（批 7 内容完善，隔离 E2E registry 0→1→0）。
- **验收对账**：行为审计通过（131 脚本零后端直连，唯一 5174 为修复工自证）；五张回执 03 真反查（`57b65f82` / `05eef3cb` / `9144865e` / `608291f0` / `718b0baa`）；界面成品化 10 屏 after 齐（`shots/audit/`）。
- **遗留橙色一览**（详见 §5，均只出方案不施工，待 Owner 裁定）：ISSUE-017a 已启用包编辑入口（需 derive-draft 端点）；ISSUE-020 P4-09 引擎不消费 GatePolicy 真越门[主权红线]/P4-07 处置断链/P4-13 对话零持久化/P4-11B 角色建议需云模型裁定；ISSUE-018 后端幂等回放状态校验；ISSUE-003b smoke 残留数据清理；ISSUE-022 监控 run_id 复用+诊断链校验拦截；启停主权对称呈现[005/013 橙]；错误契约统一 message/detail[后端侧橙]。

---

## 2. 阶段结果表

| 阶段 | 目标 | 结果 | 证据指针 |
| --- | --- | --- | --- |
| P0 环境基线 | devserver:18080 / vite:5173 / oMLX:8000 / monitor / harness 就绪 | ✅ 完成 | `logs/monitor.log`、`shots/p0-home-smoke.png`（首屏已见成品化素材；li_customer 未设显示名） |
| P1 三制作台巡检 | 场景/能力/角色制作台最小旅程走通 | ✅ 完成；16 条问题（0 阻断/12 明显/4 轻微） | `logs/p1-user-log.md`、3 份 `behavior-matrix-p1-*`、71 截图；P1 成品化项归 P4B |
| P2 加载 pack | disabled→enabled，有 receipt | ✅ 完成（UI 全回环 PASS） | `logs/p2-user-log.md`、`p2-retest-user-log.md`、`p2f-user-log.md`、`p2z-user-log.md`；回执 `608291f0`；`shots/p2z-01~29` |
| P3 修改完善 pack | 制作台内改 pack；packs 内容完善 | 🟢 主体完成（批 7 已 land 283849c；批 5-client 草稿恢复+持久化已 land） | `logs/p3-user-log.md`、`behavior-matrix-p3-*`；隔离 E2E registry 0→1→0 铁证 `718b0baa`；已启用包编辑入口=ISSUE-017a 橙待裁 |
| P4 执法行动 | 线索受理→Receipt 归档，与秘书长互动 | 🟢 主旅程完成；run 最远至 `lawyer_challenge`，送达 `provider_missing` 受控拒绝 | `logs/p4-user-log.md`、`behavior-matrix-p4-执法行动.md`；17 新问题最高 3 项入 ISSUE-020 |
| P4A 专项功能实测 | 三荚封面开关+事务卡封面+菜单启动/暂停逐项 UI 实测 | 🟢 终验大半通过：暂停/恢复双路径全通（010 修复实证）、封面折叠真实、开关可达、拒绝就近 | `logs/p4a-user-log.md`、`behavior-matrix-p4a-*`；余：启动落点无定位（019）、能力荚关→开被拒为受控正常 |
| P4B 界面成品化 | 沿途各屏去工程信息/全中文/成品排版，前后对比 | 🟢 整改主体完成；正文近零工程话（用户评「大进步」） | `logs/ui-audit-*`（7 份审计）、`shots/audit/` 10 屏 after；批 3 已 land 29670e0 |
| P5 验收对账 | 后端事实复核/回执反查/行为审计/监控 | 🟢 主体完成 | 行为审计 131 脚本零直连；五单 03 反查过；监控 backend_health 探针 200 留痕；诊断链校验拦截原 run（sequence 68≠66，fail-closed 正确）→ ISSUE-022 |
| P6 修复循环 | 发现问题→修问题 | 🟡 贯穿全程 | ISSUE 台账 001-022（§4）；6 批修复 land |
| P7 收尾 | 出最终报告 + 各仓账本转正 | ✅ 本报告 + packs 账本转正 | 本文件 + packs FEATURE_LEDGER §0/§4 更新 |

---

## 3. 核心验收证据

### 3.1 P2 UI 全回环 PASS（主权闸+真回执）

`logs/p2z-user-log.md` 铁证，四段完整闭环：

1. **停用 → 确认卡**：点开关不再一点即停，先出贴卡确认卡，四要素齐全（名称=环保执法证据链场景包 / 版本=1.0.2 / 风险=中 / 影响=退出运行、知识挂载级联下线），凭证折叠可展开（操作编号 `scene-pack-toggle-disable:…@1.0.2:b91c…`），强制勾「我已核对…」+ 取消/暂不处理退路（`shots/p2z-05/06`）。
2. **启用 → 确认卡 → 真成功**：勾核对→确认，2 秒内卡片本体变「已启用版本 1.0.2」、开关变「开 已启用·点击停用」，轮询 12 秒稳定不回跳；绿横幅回执号展开为 **`608291f0-e872-445b-84f9-fee995736576`**——与上轮复读的 `57b65f82` **完全不同**，是本次新动作的号（`shots/p2z-16/17/18`）。这是「假成功/复读回执」（ISSUE-018）批 6 修复后的真实态。
3. **刷新持久**：重开页面卡片仍「已启用版本 1.0.2」（`shots/p2z-19`）；终态 `enabled@1.0.2`。
4. **误点探针包被主权闸受控驳回**：harness 曾按「距离最近」误点隔壁「护城河探针测试包」启用开关并确认，系统受控驳回（`record is not a confirmable pack_spec_candidate: scene_pack://moat-probe-test@0.0.1 state=draft`），探针包保持未启用、**零副作用**——主权闸实证。确认卡白纸黑字包名可拦误点（设计生效）。

### 3.2 执法行动 run（至 lawyer_challenge，送达受控拒绝）

`logs/p4-user-log.md` 铁证（鑫源电镀厂夜间偷排废水执法案）：

- 7 阶段到达：发起 / 受理 / 立案（PDF 真解析，`pdf_parse_status == parsed`）/ 取证 / 证据三性（最佳回复）/ 推演 / 文书草稿；run 最远推进至 `lawyer_challenge`。
- **fail-closed 诚实**：`fact_advice`（违法事实候选）节点明示「角色槽位 enforcement_elite 已绑定，但模型网关未真实产出（provider_fallback），受控阻断（model_not_ready）：不伪造角色文本」——未编造角色发言。
- **文书草稿**：秘书长产出候选草稿摘要（当事人、pH 3.2/总铬超标 12 倍/记录仪 DSJ-7A-0233、《水污染防治法》、责令 3 日内拆除暗管、复查安排），结尾「请审核内容，确认无误后我将通过网关流转」——草稿态诚实。
- **对外送达 = 受控拒绝**：点「让秘书动手执行」，右栏秘书动作流面板出现 **`provider_missing`**——没有假成功、没有真送达，受控拒绝真实发生（`shots/p4-51/52`）。方向对（对照治理声明 `outbound_legal_document_delivery` 真实系统绑定未接入），呈现待改（裸英文码=ISSUE-020 P4-17）。
- **案件时间线真实增长**：created / 场景荚加载候选已生成 / scene_flow_run_ref 回挂 / task_ref 回挂 ×4（00:20:09～00:33:55），事件真实随动作增长，是全案唯一可靠留痕。

### 3.3 五张回执 03 反查（反伪造铁证）

P5 验收对账中，五张回执均在 03 回执账本真反查出原文（`logs/main-task-register.md` P5 行）：

| 回执号 | 来源动作 | 反查结果 |
| --- | --- | --- |
| `57b65f82` | 环保执法包首次真实启用（enabled@1.0.2） | 候选→gate_decision_712868e6→owner approved→base gate allow→pack_version_reactivated，链完整 |
| `05eef3cb` | 场景包停用 gated 链（prepare→confirm，Base 真签发） | 03 反查 200 |
| `9144865e` | pack 恢复 enabled@1.0.2（occ 6，chain_seq 9036） | 03 真存在 |
| `608291f0` | P2Z 最终回环本次新启用动作 | 与 `57b65f82` 号不同，非复读，本次新执行 |
| `718b0baa` | packs 隔离 devserver install/uninstall E2E（registry 0→1→0） | 隔离基座铁证 |

### 3.4 行为审计（131 脚本零后端直连）

用户智能体 harness 脚本与 transcript 抽查：**131 脚本零后端直连**，唯一命中 5174 端口的为修复工自证（非用户智能体绕过 UI）。用户所有关键动作经真实前端 UI / computer use 触发，符合计划 §10「用户智能体行为审计：确认无后端直连」。

### 3.5 界面成品化 before/after（10 屏）

`shots/audit/` 齐 10 屏 after：`appshell-after` / `bow-after` / `cappack-after` / `capstudio-after` / `platform-after` / `rolepack-after` / `rolestudio-after` / `scenepack-after` / `scenestudio-after` / `tasksystem-after`。配套 7 份 `logs/ui-audit-*` 审计记录（去工程信息 / 全中文 / 成品排版三项标准逐屏）。批 3 成品化 land 29670e0，正文近零工程话（用户评「大进步」）。

---

## 4. ISSUE 台账终态表（001-022）

| ISSUE | 阶段 | 根因（一句话） | 归属/颜色 | 解决 commit 或方案路径 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 001 | P1 | 三台步骤导航无反馈：能力台门控幕空 body 缺反馈 + 场景台全屏画布 CSS 滑出步骤栏（角色台/场景 3/4/9 无缺陷=点击未命中伪影） | client/绿+黄 | `logs/issues/ISSUE-001.md`；批 1 | ✅ land a9a4c87 |
| 002 | P1 | 表单治理缺失：前端丢弃 422 裁定置确认态/编辑不清陈旧 readiness/零脏态跟踪（后端三链校验正确） | client/绿~黄 | `logs/issues/ISSUE-002.md`；批 1 | ✅ land a9a4c87 |
| 003 | P1 | 能力台数据簇：onOpenReceipt 四处未接线（round 级 ref 候选态）+ smoke 残留污染正式库 + slice(0,6) 截断 | client 主 / truzhenos 数据卫生 | `logs/issues/ISSUE-003.md`（a/c 前端修；003b 清理只出方案） | ✅ land a9a4c87（003b 待 Owner） |
| 004 | P1 | 成品化簇：英文管线代号/UUID/裸哈希/错误直排/画布不可读/低对比度/工程话文案 | client 绿~黄 | 归 P4B 整改批 | ✅ 随 P4B land 29670e0 |
| 005 | P2 | 启用调错端点（confirm 只收 candidate，disabled 撞 ErrNotConfirmable=422，正路 reactivate 缺 client 包装）+ 开关被 z-10 覆盖层截获 | client/黄（对称呈现橙） | `logs/issues/ISSUE-005.md`；批 2 | ✅ land a9a4c87（对称呈现橙待裁） |
| 006 | P4A | 能力/角色荚开关同 z-index 点不到 + 受控拒绝提示远离卡片 | client/黄 | 批 2（z-index+开关移 footer+提示就近） | ✅ land a9a4c87 |
| 009 | P4A | 角色包「当前生效」区读绑定轴、开关读生命周期轴，两正交轴投成一个开关（数据自洽，投影错配） | client/黄 | `logs/issues/ISSUE-009-010.md`（开关改读绑定轴）；批 3 | ✅ land 29670e0（a558b41） |
| 010 | P4A | 多斜杠 ref 经 encodeURIComponent→Go net/http 解码 %2F→strings.Split 切碎→subresource 空→GET-only 405（一审用无斜杠 ref 得 200=用错 ref 形状） | truzhenos 05+client 双仓/黄偏橙 | `logs/issues/ISSUE-010-second.md`；兼容加法（新增 transaction_ref 变体不删旧路径） | ✅ land truzhenos 2ea6f59 + client 批 4 |
| 011 | P4A | 事务对象卡片封面开关三处皆无（§6A-4 缺口）；按 Owner 裁定=纯前端展示偏好 | client/黄 | 批 2 实现 | ✅ land a9a4c87 |
| 012 | P4A | 场景项目条目菜单无「启动」入口（§6A-5） | client/黄 | 批 2 实现（复用 FE-1 启动链） | ✅ land a9a4c87 |
| 013 | P2R | 前端状态枚举 `pack_disabled_version` 对不上后端真值 `disabled`→disabledRecords 恒空→仍路由 confirm（批 2 单测绿=fixture 用同套假枚举的自证陷阱）+ 丢 detail | client/黄（错误契约后端侧橙） | `logs/issues/ISSUE-013.md`；枚举对齐+fixture 真值+detail 传递并入批 4 | ✅ land 29670e0（UI 级批 4 回环） |
| 014 | P2R | 荚卡 footer 开关整行被卡片固定高度 overflow 裁死（双视口真人不可见不可点）——批 2 结构修复的视觉回归 | client/黄 | 批 4（vite:5174+playwright 视觉自证） | ✅ land 29670e0 |
| 015 | P2R | 「收起封面」除按钮标签外视觉零变化（折叠只切文案未收布局） | client/绿 | 批 4（视觉自证） | ✅ land 29670e0 |
| 016 | P2R | 「启动流程」落点断头（导航到网格但无启动区/无高亮/预选事件未被工作台消费） | client/黄 | 批 4（落点锚定+滚动+高亮） | ✅ land 29670e0 |
| 017 | P3 | 发布链簇：(a) 编辑入口全链缺失（后端无 enabled→draft 派生端点）；(b) 草稿蒸发（restore 门早退不调 getSceneCanvasDraft + 工作区纯 useState 零持久化）；(c) readiness 拒因丢 detail；(d) 主权面在闸后故 422 前抛出未走到 | (a) 跨仓/橙；(b)(c) client/黄；(d) 随 c 解 | `logs/issues/ISSUE-017.md`；(b)(c) 批 5；(a)=ISSUE-017a 橙待裁 | ✅ (b)(c)(d) land 044450d；(a) 橙待裁 |
| 018 | P2F | 假成功=前端固定幂等键被历史成功消费→后端 reactivate 幂等重放旧结果 replayed=true→前端不读 replayed 把回放当本次成功复读旧回执（同键 replayed:true/换新键 false 铁证）；停用 gated 链完整但静默自动 confirm=UX 缺口 | client/黄（后端幂等回放状态校验橙） | `logs/issues/ISSUE-018.md`；F1 唯一幂等键+F2 replayed 防假成功+真打 devserver 集成断言 → 批 6 | ✅ land 044450d（后端状态校验橙待裁） |
| 019 | P2F | 终验次级：反馈远离操作点/暂停回环状态漂移/启动落点无定位/命名「启动」混淆/成品化残留（英文节点码/li_customer/tx_ ref/裸 receipt-candidate:// 地址/ADOPT DATE 占位）/拒绝病句 | client/绿~黄 | 批 6（与 018 同批）+ 批 8 残留 | ✅ land 044450d + 8dd8de3 |
| 020 | P4 | P4-09=真越门（runtime.go:437 只 park human_approval/wait_event，其余 candidate_emitter 自动 succeeded，从不读 GatePolicy.PendingOwnerConfirmation 死元数据）；P4-07=candidate-input 回灌等值条件不成立不解锁+RiskLow task 不入 queue；P4-13=对话前端 useState+后端内存 map 零持久化；P4-11=A 前端无超时同步慢生成挂死（黄）+B slot 绑定齐但无 model_ref 本地 oMLX 产不出 advice（橙·需云 provider） | 06/13/05/client | `logs/issues/ISSUE-020.md`；橙红全部只出方案；P4-11A 超时修入批 8 | 🟠 橙待裁；11A ✅ land 8dd8de3 |
| 021 | P2Z | 终验回环残留：停用确认卡标题错写「启用确认」/包详情「13 回执」查不到 lifecycle 新回执（半修）/驳回文案裸英文/side_effect_class 标签英文/任务系统「编号 DATE」×6/启动区高亮过淡/停用成功无回执号就近/详情页工程话密 | client/绿~黄 | 批 8（含 P4-11A 前端超时） | ✅ land 8dd8de3 |
| 022 | P5 | 监控缺口：truzhen-monitor watch 重启复用 run_id 20260702-114448-231c833a，`monitor_event_records (run_id,sequence)` UNIQUE 冲突刷屏（新 watch 事件写入受阻，探针事件本体仍产生）；原 run 诊断包导出被链校验拦截（sequence 68≠66，fail-closed 正确） | truzhenos monitoring/黄 | 收尾报告 backlog（按纪律补回现有监控体系）；已改出 snapshot 新 run 20260702-171454-80675718 作诊断基线 | 🟠 backlog |

> 说明：ISSUE-007/008 为 P1 表单簇子项，已并入 ISSUE-002/批 2 处理（008 land a9a4c87）；台账主编号以 001-022 为准。

---

## 5. 橙色 / 待 Owner 裁定清单

以下各项本轮**只出方案、不施工**，均为需 Owner 裁定的橙色（含一项主权红线）。逐项：

1. **ISSUE-017a — 已启用包编辑入口（跨仓/橙）**：后端无 `enabled → draft` 派生端点，flowId 只由模板 slug 派生，导致已启用 pack 无法在制作台进入编辑改稿。方案：新增 derive-draft 端点（从 enabled version 派生可编辑草稿），走契约影响评估。方案文件：`logs/issues/ISSUE-017.md`。

2. **ISSUE-020 — P4-09 引擎不消费 GatePolicy 真越门【主权红线，橙红】**：`runtime.go:437` 只 park `human_approval` / `wait_event`，其余 `candidate_emitter` 节点自动 `succeeded`，**从不读 `GatePolicy.PendingOwnerConfirmation`（死元数据）**——4 张 Owner 待确认卡在用户离开后消失、run 自动推进六节点为「已完成」，未停在任何 Owner 门上。这触碰「正式动作必须受控」红线，须最高优先裁定引擎消费 GatePolicy 的整改路径。方案：`logs/issues/ISSUE-020.md`。

3. **ISSUE-020 — P4-07 处置决定断链（橙）**：candidate-input 回灌等值条件不成立不解锁；填「不适用」提交处置决定后流程不动、无反馈，生成的 BlockerTaskCandidate 在任务系统四 tab 全查无（RiskLow task 不入 queue=黄）。方案：`logs/issues/ISSUE-020.md`。

4. **ISSUE-020 — P4-13 对话零持久化（橙）**：对话历史前端 `useState`+后端内存 map，跨会话丢失（新会话重进本案沟通历史归零）。方案：对话入持久化存储。

5. **ISSUE-020 — P4-11B 角色建议需云模型裁定（橙）**：双角色 slot 绑定齐全但无 `model_ref`，本地 oMLX 产不出 advice，双角色对照「有戏台无演员」。须 Owner 裁定是否接云 provider 产角色建议（本地 fail-closed model_not_ready 已诚实）。（P4-11A 前端无超时挂死=黄，已 land 8dd8de3。）

6. **ISSUE-018 — 后端幂等回放状态校验（后端侧橙）**：`reactivate_flow.go` 幂等重放旧结果 `replayed=true` 属正确幂等语义，但缺回放态与当前实际状态的一致性校验；前端已修（F1 唯一幂等键+F2 读 replayed 防假成功），后端侧状态校验只出方案。方案：`logs/issues/ISSUE-018.md`。

7. **ISSUE-003b — smoke 残留数据清理方案（truzhenos 数据卫生/橙）**：`capability_studio_routes_test.go:51`（2026-06-12 写入）测试污染正式库残留；清理方案只出不执行 + 测试隔离防再发。方案：`logs/issues/ISSUE-003.md`。

8. **ISSUE-022 — 监控 run_id 复用 + 诊断链校验拦截（truzhenos monitoring/黄，backlog）**：watch 重启复用 run_id 致 `(run_id,sequence)` UNIQUE 冲突刷屏；诊断包导出被链校验 fail-closed 拦截（sequence 68≠66，行为正确）。须按纪律补回现有监控体系（不另起格式）：watch 重启换 run_id 或从库恢复 sequence 计数。

9. **启停主权对称呈现【005/013 橙】**：停用有确认卡+回执，启用回环亦已补确认卡+回执，但两方向的凭证呈现口径尚未统一为对称的主权面。方案随 005/013 出，呈现层统一待裁。

10. **错误契约统一 message/detail【后端侧橙】**：多处受控拒绝仅返回裸英文码（`provider_missing` / `method_not_allowed` / `record is not a confirmable...`），缺统一的 `message`（中文人话）+ `detail`（原因/指引）契约。前端已尽力透传，根治须后端统一错误契约，待裁。

---

## 6. 方法论教训

1. **固定幂等键陷阱**（ISSUE-018）：前端用固定幂等键 `scene-pack-toggle-reactivate:${ref}@${ver}`，被历史成功消费后触发后端幂等重放 `replayed=true`，前端不读 replayed 即把回放当本次成功复读旧回执——**幂等键必须每次动作唯一，且前端必须读 replayed 防假成功**。
2. **自证 fixture 禁令**（ISSUE-013）：批 2 单测绿是因为 fixture 用了与产品同一套假枚举（`pack_disabled_version` 而非后端真值 `disabled`），测试和 bug 同源——**fixture 必须取后端真值，禁止用被测代码的假设自证**。
3. **涉布局必视觉自证**（ISSUE-014/015）：footer 开关被卡片固定高度 overflow 裁死，单测断结构断不了像素——**凡涉排版/裁剪/折叠布局，必须 vite:5174 + playwright 真渲染视觉自证后才算修完**。
4. **detail 透传全链**（ISSUE-013/017c）：`packStudioCanvasJSON` 丢 detail 致无原因报错——**错误 detail 必须从后端一路透传到用户可见文案**。
5. **一审用错 ref 形状的翻案教训**（ISSUE-010）：一审用无斜杠 ref `tx_…` 得 200 判「vite 代理无罪/无 bug」，二审用真实多斜杠 ref `transaction://…` 复现 405——**复核必须用真实数据形状，无斜杠 ref 掩盖了 %2F 解码切碎的真根因**。
6. **CJK Edit 用 python**：中文长行用 Edit 工具做 unicode 精确匹配反复失败，改用 python utf-8 读写插入更可靠。
7. **cwd 重置显式 cd**：zsh 每次 Bash 调用 cwd 重置，冻结/隔离仓操作必须显式 `cd` 到目标 worktree 绝对路径。
8. **主检出他会话冲突绕行法**：truzhenos 主检出存在他会话遗留未解决 merge 冲突（UU FEATURE_LEDGER），按纪律未触碰；devserver 改为「fix worktree 编译的 2ea6f59 二进制 + 主检出 cwd（数据连续）」运行——**他会话 WIP 不覆盖，绕行保连续性**。

---

## 7. 第二轮衔接

- **第二轮授权**（Owner 2026-07-02）：本轮收口后自动启动第二轮，用 `truzhen-packs` 智能家居服务商 pack（目录=`/Users/li/Documents/truzhen-packs/smart-home-owner-pack-v0/`）按同一方法论测试；Owner 外出，登记后自主推进。
- **衍生计划已就绪**：`/Users/li/Documents/truzhen-packs/docs/plans/smarthome-pack-user-simulation-e2e-plan-20260703.md`（369 行，已批准态；§2A 结转基线 + §2B 方法论教训 + P2 语义翻转 + SH-xxx 编号）。
- **pack 事实基线**：`scene_pack://smart-home-owner-project-ops` `enabled@1.0.0`，从未 toggle（与本轮 disabled→enabled 语义相反，P2 语义翻转）；无知识库（与环保执法包 15 知识域相反）。
- **结转基线**：本轮已 land 的修复（多斜杠 ref 405、确认卡+真回执、成品化 10 屏、草稿持久化、幂等键）为第二轮共享地基；本轮橙色项（引擎消费 GatePolicy、对话持久化、对外 provider）对第二轮同样适用，第二轮遇到即引用本报告 §5，不重复分析。

---

## 8. 生命周期档位声明

- **本轮测试** = `已验收`（证据在案：P2 UI 全回环 PASS + 五单 03 反查 + 行为审计 131 脚本零直连 + 10 屏 after，符合计划 §10 严苛复核口径）。
- **各修复** = `已接线 -> 已验收`（land 事实：client `6cd6084→8dd8de3` 37 commits/730 测全绿；truzhenos `c8c3702→2ea6f59` EGR VERIFY_EXIT=0；packs `caa034d→283849c` 隔离 E2E 0→1→0；各批独立复核回归通过）。
- **橙色项**（§5 十项）= `设计中`，只出方案待 Owner 裁定；其中 ISSUE-020 P4-09 触主权红线，最高优先。

---

*报告作者：收尾报告智能体（Opus 4.8）。事实源见页首；所有回执号、commit 哈希、run 阶段、脚本计数均取自 `logs/main-task-register.md` 与各阶段 user-log，未作推演补全。*
