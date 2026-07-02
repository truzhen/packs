# 智能家居服务商 Pack 用户视角全链路实测最终报告（第二轮，2026-07-03）

> 生命周期档位：**本轮测试 = 已验收**（证据在案，见 §3/§4）。
> 计划：`/Users/li/Documents/truzhen-packs/docs/plans/smarthome-pack-user-simulation-e2e-plan-20260703.md`（P2 语义翻转 / SH-xxx 编号 / §6A 专项 / §9 沿用第一轮默认值）。
> 权威事实源：`/Users/li/Documents/过程文档/smarthomepack-live-20260703/logs/main-task-register.md`（阶段状态 SH-P0~P7 / SH-ISSUE-01~10 台账 / Land 记录）。
> 证据目录：同上 `logs/`（各阶段 `sh-p1~p4`/`sh-final` user-log、`behavior-matrix-*`、`issues/SH-ISSUE-*`）与 `shots/`（含 `shots/audit/*-sh-after.png` 巡屏 after）。
> 第一轮衔接报告：`/Users/li/Documents/truzhen-packs/docs/reports/env-pack-user-simulation-e2e-report-20260703.md`（环保执法 pack，护城河标杆样板）。本轮为智能家居 pack（第一现金流样板），与第一轮形成两类 pack 方法论对照；第一轮橙色项遇到即引用其 §5，不重复分析。

---

## 1. 执行摘要

- 沿用第一轮三方结构（用户智能体走前端 UI / computer use + 组织者[前后端协调] + 监控程序 truzhen-monitor），**一日之内跑完 SH-P0~P7**（含专项 SH-P4A、成品化 SH-P4B），以智能家居/智能装修公司工程主管的真实用户视角走通三制作台巡检 → 智能家居 pack 已启用态验证 → 制作台完善 → 王先生灯组离线售后工单经营行动 → 与项目经理（秘书长）互动至对外送达。
- **P2 语义翻转成立**：本 pack 与第一轮相反，registry 一开始就是 `enabled@1.0.0`（`occ_version=0` 从未 toggle）。故 P2 不是「加载 disabled→enabled」，而是「**验证已启用态真实性 + 停用→再启用回环回归确认（018 唯一幂等键 / 013 枚举） + 首撞已启用包编辑入口（继承第一轮 017a 橙）**」。回环 UI PASS：确认卡双向标题正确（停用确认 / 启用确认，批 8 回归 ✓），启用真实新回执 **`5f10877e`**（≠计划基线原启用号 `dd59131c`，`replayed=false`，018 反假成功回归 ✓），刷新终态还原 `enabled@1.0.0`。
- **服务单闭环结论：对话层闭上、系统层因诚实闸卡在正式化**。王先生工单五阶段秘书长对话全通（真模型 Qwen3.6-35B-A3B-4bit local oMLX live 流式，8 秒/轮 attempt=1，诊断/备件/质保口径专业度高，Proposer 边界稳定）；时间线真实回挂 task_ref；Frappe provider 三态诚实（frappe_snapshot 无假读、对外送达 `provider_missing` 受控拒绝无假发送）；**020 P4-09 越 owner_gate 未复现**（run 停在 frappe_snapshot 上游，未穿过 owner_gate，好消息）；系统层正式化受批 2 修复后终验走到「确认面弹出 + Base 签发凭据」，但被 Base 以 `owner decision target mismatch` 受控拒绝（fail-closed 零写入 = SH-F-01）。
- **修复 land 统计**（3 批 17 commits）：client main `8dd8de3 → 538c55a → 82bbc71 → 3f2b91c`（批 1-3，最终 **779 测全绿**）；truzhenos main `2ea6f59 → 764bc97`（批 2-后端：新增按 transaction_ref 拉候选的只读 GET 端点，EGR `verify` VERIFY_EXIT=0）。packs 本轮无 pack 内容修改 land（install.py 残留清理为计划 §7 列的完善候选，本轮**未做**，如实记 backlog，见 §5/§7）。
- **第二轮终验结论**：4 项确认修复（候选回填 / 保存失败反馈 / 发布拒绝落终态 / 停启回环双向确认卡）、2 项部分修复（三口径统一、停用回执就近 → 补漏进批 3）、1 项复现（受控拒绝人话动作流份 → 批 3）、SH-F-01 正式化 target mismatch（批 3 修——正式化写入 target 已改由 idempotency_key 派生候选 ref 对齐）。
- **验收对账**：行为审计通过（99 harness 脚本零后端直连，符合计划 §10「用户智能体无后端直连」）；候选回执 03 真反查（展示前缀 `7483a497`，完整 ref 由前端可复制 `4662a9d3-e24c-433e-92cb-bbe7d0f33d5d` 在 03 反查出 recorded 回执，事件类型 `gated_action_owner_confirmed` / 链序号 9368 / 证据引用 owner_action_evidence/run-gated/nonce）；停启回环生命周期新回执 `5f10877e`（P2）、`762c15d1`（终验停用）、`608cea4d`（终验启用）等真号；巡屏 after 齐（`shots/audit/*-sh-after.png`）。
- **遗留橙色 / 待 Owner 一览**（详见 §5，两轮合并去重，本轮**只出方案不施工**）：第一轮 10 项橙色继承 + 第二轮新增 6 类（SH-04 后端 secret 正则词边界收严、SH-05 画布草稿键 pack 隔离、SH-06 创建新版本入口[并 017a]、SH-09c 引擎不消费 blocked 态[引 020]、单工作区制作台草稿名互踩、测试探针 probe_A 清理）。

---

## 2. 阶段结果表（SH-P0~P7）

| 阶段 | 目标 | 结果 | 证据指针 |
| --- | --- | --- | --- |
| SH-P0 环境基线 | devserver:18080 / vite:5173 / oMLX:8000 / monitor（换新 run_id）/ harness 就绪 | ✅ 完成；首页零 console 错误 | `logs/main-task-register.md` SH-P0 行、`shots/sh-p0-home.png`；服务复用第一轮全活（superpowers worktree 已 ff 到 8dd8de3） |
| SH-P1 三制作台巡检 | 场景/能力/角色制作台最小旅程；**回归确认第一轮批 1-6 修复未回潮** | ✅ 完成；第一轮修复 11 项全回归确认（001/005a/006/009/010/011/012/013/014/015/017c）+ 新问题 4（SH-P1-01~04，合并 SH-ISSUE-01）+ 继承复现 003b | `logs/sh-p1-user-log.md`、3 份 `behavior-matrix-sh-p1-*`、53 截图 + 3 矩阵 |
| SH-P2 停用→启用回归回环 | 已启用态真实性 + 停用→再启用回环（核 receipt 新号且 replayed=false） | ✅ 核心 PASS；确认卡双向标题正确 + 启用真新回执 `5f10877e`（018 回归✓）+ 刷新还原 `enabled@1.0.0`；瑕疵 SH-ISSUE-02（停用反馈无回执号 + 详情回执投影只显 1 条） | `logs/sh-p2-user-log.md`、`sh-p2-*-confirmcard-text.txt`、`sh-p2-receipt-tab-innertext.txt` |
| SH-P3 制作台完善 pack | 制作台内改 pack；install.py 残留清理（绿，可 land） | 🔴 主体受阻；提示语行业口径真实编辑到规格草案；017a 无编辑入口继承确认；**草稿保存被后端 secret-like 误拒 + 前端假成功（SH-ISSUE-04 高）**；画布撞键灌 457 历史节点污染（SH-ISSUE-05）；发布 3 次被拒不落终态死路（SH-ISSUE-06）；「新增阶段」被 rail 遮挡（SH-ISSUE-03）；**install.py 残留清理未做（留 backlog）** | `logs/sh-p3-user-log.md`、`issues/SH-ISSUE-04-05-06.md`、`sh-p3-addstage-occlusion-evidence.json` |
| SH-P4 服务单闭环 | 用 pack 起客户项目经营 run，推进至 Receipt 归档；验证越门/任务黑洞/对话丢失 | 🟢 对话层完成 / 系统层卡正式化；王先生工单五阶段对话全通（真模型 Qwen3.6-35B live 流式）；时间线 task_ref 真挂；主权闸全 fail-closed（**020 P4-09 越门未复现**）；对外送达 `provider_missing` 无假发送✓；卡点=正式化死入口（SH-P4-05）+ owner_gate 需真 11 回执回填永不可达；新问题 SH-P4-01~09 | `logs/sh-p4-user-log.md`、`behavior-matrix-sh-p4-服务单.md`、`sh-p4-chat-transcript.txt` |
| SH-P4A 专项回归 | §6A 六条（三荚封面开关 / 事务对象封面开关 / 场景项目启动·暂停）回归确认 | ✅ 终验 8 项：候选回填 / 保存反馈 / 发布终态 / 停启回环✅；正式化确认面已通但 confirm target mismatch（SH-F-01）；三口径部分（SH-F-02）；动作流 provider_missing 复现（SH-F-03）；回执 tab 只列 1 单（SH-F-04） | `logs/sh-final-user-log.md`、`sh-f-*.txt`（逐项 dump） |
| SH-P4B 成品化增量 | 沿途各屏 §7A 三标准截图审计 + 智能家居专属新文案审计 | ✅ 4 屏 after 巡屏入 `shots/audit/`；角色徽标 / role_pack:// 清零确认（SH-P1-01/02 修复✓）；场景制作台扫描零命中；残留=沟通中心右栏 scene_flow_pending（SH-F-05）→批 3 | `shots/audit/{bow,bowdetail,scenestudio,rolepack}-sh-after.png`、`sh-f-28-scan-findings.txt` |
| SH-P5 验收对账 | 回执 03 反查 / 行为审计 / 监控事件 | 🟢 主体完成；候选完整 ref `4662a9d3` 在 03 反查 recorded✓；99 脚本零后端直连；生命周期新回执 `5f10877e`/`762c15d1`/`608cea4d` 真号 | `logs/main-task-register.md` SH-P5、`sh-f-03-lookup-created.txt` |
| SH-P6 修复循环 | 每问题登记 SH-NNN → 外派 Opus 分析定根因 → 分诊 land / 只出方案 | 🟡 贯穿全程；SH-ISSUE 台账 01-10（§4）；3 批修复 land | `issues/SH-ISSUE-04-05-06.md`、`issues/SH-ISSUE-09.md` |
| SH-P7 收尾 | 最终报告 + packs 账本收口 | ✅ 本报告 + packs FEATURE_LEDGER §0/§4 更新 | 本文件 + packs `docs/sh-round-closeout-20260703` |

---

## 3. 核心验收证据

### 3.1 第一轮修复 11 项全回归确认（SH-P1，未回潮）

`logs/sh-p1-user-log.md` 与 `logs/sh-p2-user-log.md` 铁证——第一轮已 land 修复在智能家居 pack 上依然生效（回潮即为回归 bug，本轮无回潮）：三台步骤导航点击有反馈（001）；表单接收后端裁定、编辑清陈旧 readiness、脏态跟踪（002/017c）；荚卡开关 z-index 可点、开关在 footer、就近提示（005a/006/014/015，SH-P2「frame 内 BUTTON 直达」实证）；角色包「当前生效」区改读绑定轴（009）；多斜杠 ref 子资源 405 修复（010）；事务对象封面开关=纯前端展示偏好（011）；场景项目菜单「启动」复用既有启动链（012）；启用重放 replayed 处理 + 枚举中文（013，「已启用版本 X / 未启用 · 最近版本 X」）。

### 3.2 SH-P2 停用→启用回环 UI PASS（018 唯一幂等键 + 反假成功回归）

`logs/sh-p2-user-log.md` 铁证，回环双向完整：

1. **停用 → 确认卡**：点 footer 开关（aria-label「已启用 智能家居老板项目经营 Pack，点击停用」）出确认卡，**标题加粗「停用确认」**（批 8 方向✓）；四要素齐（名称 / 版本 1.0.0 / 风险「中（停用全队在用的场景包）」/ 影响「知识挂载级联下线、全队无法使用直到重新启用」）+ 凭证详情折叠 + 强制勾核对 + 取消 / 暂不处理退路；先点「取消」验证卡片保持已启用（无副作用✓）。
2. **启用 → 确认卡 → 真成功**：标题「启用确认」✓；凭证详情展开操作编号 `scene-pack-toggle-reactivate:scene_pack://smart-home-owner-project-ops@1.0.0:38ce18e6d9714741a8440e5f0822c7f6`——**幂等键带唯一 nonce 后缀（018 唯一幂等键回归✓）**；确认后回执号 **`5f10877e-56c9-4a26-8357-137504552b8d` 为新号**（≠计划基线原启用回执 `dd59131c`，`replayed=false`，018 反假成功回归✓）。
3. **刷新持久**：重开页面卡片仍「已启用版本 1.0.0」；终态 `enabled@1.0.0`（回环还原）。
4. **旁证隔离**：环保执法证据链场景包全程保持「已启用版本 1.0.2」未被误碰；未碰任何其它包开关 / 归档 / 删除。

（终验 SH-P4A 复走第二次停启回环，另产停用单 `762c15d1-cb1e-4fa4-91a0-26a40b396038` / 启用单 `608cea4d-40f6-42bd-a90e-1c221ad0594b`，见 `logs/sh-final-user-log.md` §7；终态仍还原 `enabled@1.0.0`。）

### 3.3 SH-P4 王先生工单五阶段真模型对话（Qwen3.6-35B 流式）

`logs/sh-p4-user-log.md` + `sh-p4-chat-transcript.txt` 铁证（客户王先生报修「全屋智能装完两周，客厅灯组集体离线，网关红灯常亮，App 里设备全部失联」）：

- 五阶段全通：受理 → 远程诊断（补充 TZ-GW200 / 红灯=广域网断连 / 客户昨天换过路由器）→ 上门检修安排 → 现场处置（重新配网 / 灯组固件 2.1.3 批量升级 / 47 台设备恢复在线）→ 客户确认与回访。五轮全部 **8 秒/轮、attempt=1**，运行时详情 = `Qwen3.6-35B-A3B-4bit · local omlx · mode live · Model Gateway 已调用（流式）`。
- **秘书长质量**：诊断推理正确用足线索（新路由器→网关未注册广域网→重启 / 查光猫 / 上门条件）、备件清单与质保口径（配置软件问题非硬件、标准质保）像真项目经理；每轮主动声明「以上为待确认建议，未正式派工 / 未正式归档」，**Proposer 边界口径稳定**。三处瑕疵（诚实记录，非系统假成功）：称呼错位（把工程主管当成王先生本人）；候选卡 summary 逐字复读用户原话非任务概括（SH-P4-04）；第 5 轮口头宣称「工单状态已更新为待归档」但读回详情页 run 状态纹丝未动——**系统层没有假成功，是模型话术过头**。

### 3.4 SH-P4-06 候选回填修复（批 2）终验 PASS + 018 回归 + 新回执 03 反查

`logs/sh-final-user-log.md` §1 铁证——SH-P4-06 首现「售后候选卡刷新蒸发」（SH-ISSUE-09a 根因=候选/回执真落 SQLite+03，但后端无按 tx 拉候选的 GET、前端候选卡纯 useState 无回填、时间线映射丢 EventRef），批 2 双仓修复后终验回归：

- 冷会话进王先生工单详情：售后方案候选卡**回填在场**——`project_control_candidate · task_cand_40c485e97067e200`，状态行显示**完整回执 UUID**「状态：candidate · receipt 7483a497-a009-499e-b612-a847e688e80e」，带确认 / 修改 / 驳回三按钮。
- 时间线 8 条事件全部带「回执号（点开查看完整并可反查）」折叠：展开出完整 ref +「复制回执号」（剪贴板实测 **`4662a9d3-e24c-433e-92cb-bbe7d0f33d5d`**）+「在 03 反查」——反查真返回 **recorded 回执**（事件类型 `gated_action_owner_confirmed`、链序号 9368、链哈希、证据引用 `owner_action_evidence/run-gated/nonce`）（`sh-f-03-lookup-created.txt`）。
- 反内存态铁证（引 SH-ISSUE-09a）：`GET /v3/receipts/7483a497` 返 404 `receipt_not_found`——`7483a497` 是 8 位展示前缀非完整 ref，恰印证根因（完整 receipt_ref 此前被前端映射丢弃，批 2 前端保 ref 后可复制反查）。

### 3.5 SH-P4-05 正式化死入口修复（批 2 两段式接 gated 链）+ SH-F-01 target mismatch（批 3）

`logs/sh-p4-user-log.md` §3 首现「正式化登记按钮永远撞墙」（`formalizeOpenObject` 裸提交无 proof，后端 `owner_decision_ref is required`）。批 2 前端接既有 `issueOwnerDecisionGrant`（prepare→confirm）两段式链、零自铸凭据后，终验（`logs/sh-final-user-log.md` §2）走到：

- 点「正式化登记（需主人主权确认凭据）」→**真弹确认卡**（不再死入口）：标题「高风险动作确认：business_object_create_with_scene」+「确认即构成 OwnerActionEvidence；decision_ref 由后端签发，执行面拒绝任何前端自带决策引用」+ 目标（tx_39a2c9d0…）/ 内容 / 影响三要素 + 技术凭据折叠（run_id `run-gated-b94a45ed-…`、nonce `nonce-c4ee0fbe-…` 均标注后端签发）；勾核对前确认按钮 disabled✓。
- 确认 →6 秒后**受控拒绝**：「业务对象正式写入 proof 被拒：owner decision target mismatch（Owner 主权链 proof 由 Base 签发；未通过即不写入）」，2/2 复现，对象仍 candidate_only / non_formal，**零写入 fail-closed 正确**——卡点从「无凭据路径」变为「凭据目标错配」（SH-F-01）。批 3 修复：正式化写入 target 已改由 `idempotency_key` 派生候选 ref 对齐（确认卡目标 tx_ 与写入侧校验口径统一）。

### 3.6 停启回环双向 + 候选回填 + 正式化 gated 链接通 + 送达 provider_missing 受控拒绝 + 99 脚本零直连（终验汇总）

`logs/sh-final-user-log.md` 8 项终验结论：候选回填✅（1→8 条回执投影，见 §3.4）；正式化 gated 链接通（确认卡 + Base 签发凭据，target mismatch 由批 3 修，见 §3.5）；停用启用双向确认卡标题正确 + 双向回执号就近✅（停用单 `762c15d1` / 启用单 `608cea4d`）；对外送达 `provider_missing` 受控拒绝无假发送✅；frappe_snapshot 缺口诚实人话已生效（「此步骤需要外部系统连接才能自动读取结果……秘书不会伪造结果」，SH-P4-07 修复确认）；巡屏 after 完成（4 屏）；角色包管理扫描干净（徽标全改「角色包」、role_pack:// URI 清零）。**行为审计**：99 harness 脚本抽查零后端直连，用户所有关键动作经真实前端 UI / computer use 触发。

---

## 4. SH-ISSUE-01~10 台账终态

| ISSUE | 阶段 | 一句根因 | 归属 / 颜色 | 解决 commit 或方案路径 | 状态 |
| --- | --- | --- | --- | --- | --- |
| SH-ISSUE-01 | SH-P1 | 角色包 10 卡类型徽标全误标「场景包」+ role_pack:// URI 裸排 + 悬浮卡 active/watch 英文 + Esc 关不掉浮层×3（SH-P1-01~04 合并；类型映射 / 成品化漏网 / 键盘可达性） | client / 绿~黄 | SH 批 1 | ✅ land `538c55a` |
| SH-ISSUE-02 | SH-P2 | 停用成功反馈无回执号（全 UI 无处查）+ 包详情回执投影只显 1 条（停用单不可见）；投影过滤 / 数据源覆盖不全 | client / 绿~黄 | SH 批 1 | ✅ land `538c55a`（回执号就近已补；投影覆盖份终验仍复现 SH-F-04，见 §5） |
| SH-ISSUE-03 | SH-P3 | 制作台「新增阶段」按钮被 flow-product-rail（z=30 fixed 280px）遮挡不可点；CSS 层级遮挡 | client / 黄 | SH 批 1（视觉自证 occlusion json） | ✅ land `538c55a`（终验命中 top=BUTTON 本体确认修复） |
| SH-ISSUE-04 | SH-P3 | 草稿保存 secret 误报 + 前端假成功：`store.go:20` 正则 `sk-[a-z0-9_-]{12,}` 无词边界，系统自造节点 id `node-task-<13位时间戳>` 内含 `sk-1783019102407` 命中（重放 A 纯中文=saved / B 带 id=拒，铁证）；前端 `SceneStudioPage.tsx:3078` `void` 吞失败无条件推进阶段 | 后端正则=橙（安全语义）；前端 fail 反馈=黄 | `issues/SH-ISSUE-04-05-06.md`；前端份 SH 批 1；正则收严入 §5 | 🟠 前端 land `538c55a` / 后端橙待裁 |
| SH-ISSUE-05 | SH-P3 | 画布撞键灌历史节点：`SceneStudioPage.tsx:238` flowId 仅由模板族 slug 派生无 pack 隔离，撞上同族「装修项目流程包」历史草稿（flowId `long_term_project_delivery_scene_flow`，**457 节点 258KB**，6-15 起累积；**非环保包**，用户猜证伪）；restore :1706 载回灌画布 | truzhenos `canvas_store.go:30`（FlowID 是 primaryKey）+ client 跨仓 / 橙 | `issues/SH-ISSUE-04-05-06.md`（并 017a 合并设计）；污染草稿**保留**（该 flowId enabled=1.0.1，删会丢装修包） | 🟠 橙待裁 |
| SH-ISSUE-06 | SH-P3 | 发布死路 + 不落终态：拒因 `ErrDraftFrozen`（1.0.0 已 promote、enabled=1.0.1）设计正确；前端无「创建新版本」入口（grep 零）+ catch 不落终态 + 丢 detail（ISSUE-017 P3-03 第三现场） | 落态 / detail=黄；新版本入口=橙跨仓（并 017a） | `issues/SH-ISSUE-04-05-06.md`；前端份 SH 批 1；新版本入口入 §5 | 🟠 前端 land `538c55a` / 橙待裁 |
| SH-ISSUE-07 | SH-P4 | 状态三口径矛盾：封面「待启动流程」vs 药丸「进行中·customer_draft」vs 悬浮卡「scene_flow_pending」+ 悬浮卡英文枚举；三处各读不同字段无统一映射 | client / 黄 | SH 批 2 | ✅ land `82bbc71`（卡片内两处已统一；悬浮卡数据源份终验仍矛盾 SH-F-02，见 §5） |
| SH-ISSUE-08 | SH-P4 | 受控拒绝呈现无人话：动作流裸英文 `provider_missing` 无解释 + frappe_snapshot 缺口全程无诚实人话；前端未接中文映射 / 缺口说明 | client / 绿~黄 | SH 批 2 | ✅ land `82bbc71`（frappe_snapshot 缺口人话已生效；动作流份终验仍复现 SH-F-03，见 §5） |
| SH-ISSUE-09 | SH-P4 | 系统层正式化簇三子问题：(a) 候选蒸发=读侧缺口（候选 / 回执真落 SQLite+03，后端无按 tx 拉候选 GET，前端候选卡纯 useState 无回填，时间线映射丢 EventRef）；(b) 死入口=前端 `formalizeOpenObject` 裸提交无 proof（后端 gated-action 链已存在被 `startSceneFlowRun` 跑通，不用改后端）；(c) 阻断误标=`candidate_emitter.go:31` 恒 Succeeded，**本例 owner_gate 未被越过**（轻 020 一档，非主权红线触发） | (a) 后端只读 GET=橙轻 + 前端=黄；(b) client / 黄；(c) truzhenos 06=橙引 020 | `issues/SH-ISSUE-09.md`；(a) 后端 + 前端、(b) client → SH 批 2；(c) 引 020 入 §5 | ✅ (a)(b) 双仓 land（`764bc97`+`82bbc71`）；(c) 橙引 020 待裁 |
| SH-ISSUE-10 | 终验 | SH-F-01 正式化 confirm 被拒 `owner decision target mismatch`（确认面 / 凭据已通，target 拼接不一致，fail-closed 零写入）；F-02 悬浮卡「待启动流程」仍矛盾；F-03 动作流裸 provider_missing 第二处；F-04 详情回执 tab 只列 1 单；F-05 沟通中心右栏 scene_flow_pending | client / 黄 | SH 批 3（末批） | ✅ land `3f2b91c`（SH-F-01 正式化 target 由 idempotency_key 派生候选 ref 已修；F-02~F-05 投影 / 枚举残留随批 3 收敛） |

> 说明：SH-P1-01~04 合并为 SH-ISSUE-01；SH-P4-01~09 中系统层三项归 SH-ISSUE-09、呈现三项归 SH-ISSUE-07/08、终验残留归 SH-ISSUE-10；SH-P4-01（阻断误标完成）= SH-ISSUE-09c，橙引 020 只出方案。台账主编号以 SH-ISSUE-01~10 为准。

---

## 5. 橙色 / 待 Owner 裁定清单（两轮合并去重）

以下各项本轮**只出方案、不施工**，均需 Owner 裁定。第一轮 10 项橙色（引第一轮报告 §5，本轮遇到即引用不重复分析，标注本轮复现情况）+ 第二轮新增 6 类。

### 5.A 第一轮橙色继承（本轮复现情况标注）

1. **ISSUE-017a — 已启用包编辑入口（跨仓 / 橙）**：后端无 `enabled → draft` 派生端点，flowId 只由模板 slug 派生。**本轮 SH-P3 首撞实证**：本 pack `enabled@1.0.0`，制作台无路径进入编辑改稿。方案：新增 derive-draft 端点。与 SH-ISSUE-05/06 的「新版本入口」是同一能力缺口两个面，建议合并设计。
2. **ISSUE-020 P4-09 — 引擎不消费 GatePolicy 真越门【主权红线，橙红】**：`runtime.go:437` 只 park `human_approval` / `wait_event`，其余 `candidate_emitter` 自动 `succeeded`，从不读 `GatePolicy.PendingOwnerConfirmation`。**本轮 P4 越门未复现**（run 停在 frappe_snapshot 上游，未穿过 owner_gate），但同族 `candidate_emitter.go:31` 恒 Succeeded 机制在 SH-ISSUE-09c 复现（阻断态被伪装已完成，见下 §5.B 第 4 项）。第一轮红线定性不变，最高优先。
3. **ISSUE-020 P4-07 — 处置决定断链（橙）**：candidate-input 回灌等值条件不成立不解锁；RiskLow task 不入 queue 四 tab 查无。**本轮 P4 复现**（任务系统主区 6 tab 零命中）。方案引第一轮 020。
4. **ISSUE-020 P4-13 — 对话零持久化（橙）**：对话前端 useState + 后端内存 map 跨会话丢失。**本轮 P4 复现**（且连带候选评定 / 动手执行入口一次性丢失，SH-P4-09）。方案：对话入持久化存储。
5. **ISSUE-020 P4-11B — 角色建议需云模型裁定（橙）**：本 pack 单角色项目经理，无双角色对照——本轮对话路径经同一模型网关真产出（Qwen3.6-35B live），此橙项本轮不构成障碍；但流程节点 manager_advice 仍 provider_fallback（口径矛盾归 SH-ISSUE-09c）。
6. **ISSUE-018 — 后端幂等回放状态校验（后端侧橙）**：前端已修（唯一幂等键 + 读 replayed），后端缺回放态与实际状态一致性校验。**本轮 SH-P2 回环未触发假成功**（新号 replayed=false），后端侧校验仍只出方案。
7. **ISSUE-003b — smoke 残留数据清理（truzhenos 数据卫生 / 橙）**：测试污染正式库残留。**本轮复现**（侧栏 7 个「未命名场景项目」组，SH-P4-08）。清理方案只出不执行。
8. **ISSUE-022 — 监控 run_id 复用（truzhenos monitoring / 黄，backlog）**：watch 重启复用 run_id 致 `(run_id,sequence)` UNIQUE 冲突刷屏（本轮 P0 起已按教训换新 run_id 规避）。须按纪律补回现有监控体系。
9. **启停主权对称呈现【005/013 橙】**：停用 / 启用凭证呈现口径尚未统一为对称主权面。本轮 SH-ISSUE-02 的「停用回执号就近 / 详情投影覆盖」是其表现，批 1 补回执号就近，投影覆盖份终验仍复现（SH-F-04）。
10. **错误契约统一 message/detail【后端侧橙】**：多处受控拒绝仅返回裸英文码（`provider_missing` / `owner decision target mismatch` 等），缺统一 `message`（中文人话）+ `detail`。本轮 SH-F-03 / SH-F-01 顶栏裸英文是其表现。根治须后端统一错误契约。

### 5.B 第二轮新增橙色 / 待完善

1. **SH-04 后端 secret 正则词边界收严（truzhenos / 橙 — 安全语义）**：`packstudio/draft/store.go:20` 正则 `sk-[a-z0-9_-]{12,}` 无词边界，误吞系统自造节点 id `node-task-<13位时间戳>`（=`sk-1783019102407`）。方案：加词边界 / 前缀锚定（如要求 `sk-(proj|live|test|ant)-` 已知真 key 前缀），或**更稳方案**——把检测作用域从「整份 JSON 序列化」改为「只扫 description/label 等用户自由文本字段，跳过系统生成的 id/ref/hash 字段」（根因是把机器生成 id 当用户内容扫）。**橙**因该门是安全检测门，收严须评估不放过真 key，且正则被 SaveDraft 与 SaveCanvasDraft 共用、须对齐 `monitoring/redactor.go`、`storage/testutil/plaintext_scanner.go` 同类口径防漂移。方案：`issues/SH-ISSUE-04-05-06.md`。
2. **SH-05 画布草稿键 pack 隔离（truzhenos + client 跨仓 / 橙 — 键空间语义）**：`CanvasDraftModel.FlowID` 是 primaryKey，flowId 仅由模板族 slug 派生无 pack 维度，撞同族历史草稿。方案：草稿键加 pack 隔离维度（`<slug>__<packSlug>_scene_flow` 或复合键 `(flowId, pack_draft_id)`）。**橙**因改键空间牵动 06 引擎 SaveSpec/GetSpec 一致性、历史草稿迁移、lifecycle record 的 flow_id 引用，须出兼容清单。污染草稿（`long_term_project_delivery_scene_flow` occ=70 258KB）**保留不删**（其 enabled pointer 已 1.0.1，是装修包真实来源资产，删会丢草稿），靠键隔离治本。方案：`issues/SH-ISSUE-04-05-06.md`。
3. **SH-06 创建新版本入口（truzhenos + client 跨仓 / 橙 — 状态机 / 能力新增，并 017a）**：冻结 record（`pack_spec_candidate`）SaveDraft 被拒正确，但 UI 无「以当前草稿开新版本号走全新 draft→readiness→promote」入口 = 死路。方案：后端加 `/v3/pack-studio/lifecycle/new-version { pack_ref, base_version }`（从冻结 record 派生新版本号草稿，candidate_only，产 receipt），前端在发布被拒面加「创建 1.0.x 新版本继续编辑」按钮。**橙**因触及版本派生语义 / 版本号冲突策略 / 与云市场版本口径关系。与 017a「反投影可编辑草稿」合并设计。方案：`issues/SH-ISSUE-04-05-06.md`。
4. **SH-09c 引擎不消费 blocked 态（truzhenos 06 / 橙 — 引 020）**：`candidate_emitter.go:31` 恒 `StepStatusSucceeded`；capability 节点 `executor.go:109-110` 不查 provider、advice 节点 `executors.go:453` BlockedReason 不映射 StepResult.Status，导致 frappe_snapshot（无 provider）/ manager_advice（model_not_ready）被误标「已完成」流程照走。**本例未越 owner_gate**（run 停在其上游），比 020 P4-09 真越门轻一档 = 呈现 / 状态机诚实性缺陷。方案：让 runtime 消费 blocked 信号（BlockedReason 非空 ⇒ 节点非 succeeded ⇒ 触发已存在 blocked 分支），与 020 P4-09 GatePolicy-park 方案同批收敛（一处改状态机消费 blocked/gate 态即同时修 P4-09 + P4-10 + 本例）。方案：`issues/SH-ISSUE-09.md`。
5. **单工作区制作台草稿名互踩（client / 黄~橙 — 产品模型缺口）**：终验因单工作区模型，纯中文草稿保存真生效后，制作台工作区草稿名被本轮改为「终验测试草稿（保存失败反馈回归）」（原「项目交付流程包」）。此为 SH-05 键空间无 pack 隔离在草稿名层面的连带表现——单工作区下不同 pack 编辑不可避免互踩。随 SH-05 键隔离方案一并设计。见 §7 数据残局。
6. **测试探针 probe_A 清理（truzhenos 数据卫生 / 黄，backlog）**：SH-ISSUE-04 动态复核向 `pack_studio_dev.sqlite` 写入测试草稿探针 `flow_id=SHISSUE04_TEST_probe_A_flowid`（中文业务文本 saved occ=1；probe_B 被 secret 拒未落盘）。清理方案（不执行，供 Owner 授权）：`DELETE FROM canvas_draft_models WHERE flow_id LIKE 'SHISSUE04_TEST_%';` + 对应 version_models（仅此两键，不碰真实 pack 草稿）。

---

## 6. 方法论教训（增量，第一轮 8 条见其报告 §6）

1. **waiter 匹配 echo 落点**：后台长命令（EGR verify）的完成信号需读真实 `VERIFY_EXIT` 值，`echo` 落在命令末尾会让「exit 0」误当验证通过——waiter/until 循环须 grep `VERIFY_EXIT=<值>` 而非命令退出码（第一轮 §12.7 强化）。
2. **glob null 坑**：zsh 裸 glob / `echo ===` 在无匹配时会炸；截图与 dump 文件批量处理用 `while-read` + 显式路径，不依赖 shell 分词。
3. **断线接力 = 先查 git log 定损**：本轮 client 批 1 含前任断线接力，恢复时**先 `git log --oneline` 定位真实 HEAD 与已 land 增量**，再叠改动，避免重复施工或覆盖他人 WIP（服务 worktree 已 ff 到基线后再动）。
4. **`"candidate".slice(-4)="date"` 式合成 id 陷阱**：SH-04 揭示——凡把「系统生成的 id / 时间戳 / hash」当用户自由文本参与安全扫描 / 语义判断，都会撞出假阳（`node-task-<timestamp>` 命中 `sk-…`）。教训：安全检测 / 内容判定必须区分「用户自由文本字段」与「机器生成标识字段」，只扫前者。
5. **反内存态论证靠「反查 404 恰印证」**：SH-P4-06 候选蒸发的根因证明——`GET /v3/receipts/<8位前缀>` 返 404 不证伪「回执在 03」，恰印证「完整 ref 被前端映射丢弃」（前缀查不到 ≠ 数据不在库）。复核缺口时须用完整 ref 形状，短前缀掩盖真根因（同第一轮 ISSUE-010「用错 ref 形状翻案」一脉）。

---

## 7. 数据留痕与残局

- **pack 数据终态**：`scene_pack://smart-home-owner-project-ops` **已启用 1.0.0**（回环还原✓）；本轮真实动作痕迹=停用×2 / 启用×2（P2 一次 + 终验一次），新增生命周期回执 `5f10877e` / `762c15d1` / `608cea4d` / `a9b769a6` 等真号。环保执法包、李女士数据未触碰（李女士条目仅 hover 对照，零写入）。
- **王先生工单（本轮测试单）保留**：创建 1 单（含自动 run）+ 六轮对话消息 + 售后工单候选 1 张（ready，未采纳未驳回，回执 `7483a497` 完整 `4662a9d3`）+「让秘书动手执行」2 次（均死于 provider_missing 零副作用）+ 正式化登记确认 2 次（均受控拒绝零写入，Base 侧有 gated 拒绝痕迹）。侧栏落「未命名场景项目」组（SH-P4-08 / 003b 复现）。**请组织者 / Owner 对账是否清理测试单**。
- **制作台草稿名被改**：单工作区模型下，制作台工作区草稿名被本轮改为「终验测试草稿（保存失败反馈回归）」（原「项目交付流程包」）；画布草稿写入全部被 secret 拒未落盘，发布 2 次点击均被服务端拒绝无落地物。**请 Owner 裁定是否还原草稿名**（根因随 SH-05 键隔离方案治本）。
- **测试探针 probe_A**：`pack_studio_dev.sqlite` 残留 SH-ISSUE-04 复核探针 `flow_id=SHISSUE04_TEST_probe_A_flowid`（明显测试标记）。清理 SQL 见 §5.B 第 6 项，待 Owner 授权执行（碰真相表，不自行删）。

---

## 8. 生命周期档位声明

- **本轮测试** = `已验收`（证据在案：SH-P2 停启回环 UI PASS + 候选回填终验 + 完整 ref `4662a9d3` 03 反查 recorded + 行为审计 99 脚本零直连 + 巡屏 4 屏 after，符合计划 §10 严苛复核口径）。
- **各修复** = `已接线 -> 已验收`（land 事实：client `8dd8de3→538c55a→82bbc71→3f2b91c` 批 1-3 共 17 commits / 779 测全绿，各批组织者独立复核 VERIFY_EXIT=0；truzhenos `2ea6f59→764bc97` 批 2-后端只读 GET，EGR verify ok，devserver 换二进制线上验证通过）。
- **橙色项**（§5 第一轮 10 项 + 第二轮 6 类）= `设计中`，只出方案待 Owner 裁定；其中 ISSUE-020 P4-09 触主权红线（本轮越门未复现但机制在 SH-09c 同族复现），最高优先。
- **pack 内容完善**（install.py 残留清理，计划 §7 第 1 项）= `未做 / backlog`——本轮 SH-P3 主体受阻于制作台缺陷（SH-ISSUE-04/05/06），install.py docstring 首行「环保执法 Pack」残留与 knowledge batch tags 硬编码「环保执法」未清理（本 pack 无知识库不走该行，属残留），如实记 backlog（见 packs FEATURE_LEDGER §6）。

---

*报告作者：收尾报告智能体（Opus 4.8）。事实源见页首；所有回执号、commit 哈希、run 阶段、节点计数、脚本计数均取自 `logs/main-task-register.md` 与各阶段 user-log / issues，未作推演补全。*
