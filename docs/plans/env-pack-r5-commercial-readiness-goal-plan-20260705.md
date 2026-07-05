# 环保执法 Pack · R5 商用就绪轮 — GOAL 定义 + 建设/测试/收尾计划（2026-07-05）

> **计划-only + /goal 目标定义**。本文件供 Owner 在**新会话 /goal 模式**驱动；「写清目标」= 给出 /goal 的**可判定退出条件**。**执行前仍须 Owner 在新会话确认「开工」并逐条裁定 §11 待裁**（本轮含橙/红建设，不是纯测试）。
> **来源**：Owner 2026-07-05「完善下一轮环保执法测试计划文档，我将在新会话 /goal 模式开始；目标 = 这一轮后 **truzhen + 环保执法 pack 都能正式商用**。要求 ① pack 能产生**即时任务 + 阶段任务**，所有任务经**任务管理驱动**；② **codex CLI 作为第一优先 hands 能在 docker 中执行各类任务**」。
> **承接**：R4 报告 `/Users/li/Documents/过程文档/environmentalpack-round4-20260704/environmentalpack-round4-report.md` + 本会话已 land 的 pack-suite hardening（os `b4ee07b` / client `966c2f1`）+ 下方 REG 回归基线。

---

## 0. GOAL（/goal 模式主目标 + 子目标，可判定）

**主目标（GOAL）**：环保执法 Pack 在本地隔离栈达到「**受控商用就绪**」，并产出可呈 Owner 的**正式商用裁定证据包**。达标 = 下列 **G1–G6 全绿 + 组织者 neuter 独立复核无伪造**。

| 子目标 | 达标判据（/goal 判绿条件） |
| --- | --- |
| **G1 任务驱动骨架** | 环保 9 阶段工作**全部**经 **07 任务管理**驱动：SceneFlowRun 每个正式动作节点生成 TaskCandidate（区分**即时任务 immediate** + **阶段任务 phase**）→ 07 队列 → Owner/agent 裁定 → dispatch → execution → 03 receipt → 回挂 timeline；**零绕过任务管理的执行**（静态守卫禁 06→11 直连 + 运行时每次执行带 `task_ref`）。 |
| **G2 codex-cli-docker hands 第一优先** | 任务 dispatch 到 **11 Execution Gateway** 时，**codex CLI sidecar（Docker `net=none` 沙箱）为第一优先 hands** 真实执行任务（草拟/分析/生成），产候选 + 03 回执；未接通 → `provider_missing`/`not_ready` **诚实态不伪造**；fallback 到其它 hands 有明确路由证据。 |
| **G3 主权闸零绕过** | 所有「**对外正式生效**」动作穿 **Owner + Base Gate** 产 03 可反查 receipt；codex 沙箱内产物 = **候选留痕（免闸）**，跨到真实对外 = **必过闸**；留痕 vs 过闸边界静态 + 运行时双证。codex 永远 Proposer。 |
| **G4 商用闭环 GUI 实测** | 购买 → 授权（License/Entitlement）→ 安装 → 启用 → **任务驱动执行**（即时 + 阶段）→ 回执 全链路 GUI 走通；真实 paid 订单**或** Owner 明确商用验收替代。 |
| **G5 REG 无回归** | 本会话及此前 land 的所有修复（REG-R5-01~09）现状无回归。 |
| **G6 诚实证据包** | 产出 `环保pack-商用就绪证据包`（每 G 双证 + neuter 复核 + 诚实上限声明），供 Owner 做**正式商用裁定**。 |

**退出条件（/goal 判绿）**：G1–G6 全绿且证据包完成 → 生命周期置「**已接线 → 已验收（本地隔离）+ 待 Owner 商用裁定**」。
**★ /goal 反造假铁律**：G1–G4 任一以 **mock / demo / 假成功 / 自铸 ref / 沙箱假执行 / 绕过任务管理或主权闸** 达成 = **不达标，判红**。「正式商用」是 Owner **裁定态**，不是 /goal 或 Agent 的**声明态**——本轮把 pack 推到「就绪 + 证据齐」，最终商用由 Owner 依证据包拍板。

---

## 1. 战略定位校准 + 诚实上限（先对齐，防「一轮宣称商用」腐化）

- **Owner 本轮裁定（需 §11.1 确认）**：环保执法 pack 从「护城河标杆样板（非第一现金流，2026-06-21 战略）」**提升为本轮商用就绪目标**。此提升改变既有战略排序，须 Owner 明确确认。
- **诚实上限（不可绕）**：缺**真实客户原话**时，「正式商用」由 **Owner 商用验收替代**，但**不得由 Agent 自我宣称「已发布 / 客户生产全面可用」**。证据包必须写明「本地隔离验收 + 待 Owner 商用裁定」。
- **truzhen 本体 vs env pack**：叠加关系——truzhen 基座提供**通用能力**（任务管理驱动 + codex-docker hands + 主权闸 + 商用闭环）；env pack 是**首个跑通商用闭环的场景荚样板**。**codex-cli-docker 是通用 hands**，env pack 是测试载体（§11.2 确认此定位）。

## 2. 需求验证 + 真相源 + 归属 + 风险（两新能力落位，开工前必对齐）

### 2.1 需求验证（第一性）
- **需求 1（任务驱动）**：证据 = Owner 2026-07-05 直述产品要求。**缺真实客户原话** → 标记「Owner 产品意志驱动」（非客户压力驱动）；本轮以 env 样板证明架构，商用宣称待真实 pilot / Owner 验收。
- **需求 2（codex-docker hands）**：**部分已存在** —— 11 SidecarBackedProvider `kind=codex` + `codex_provider.py` `codex exec`（docker `net=none`/ro，Go core 不碰 docker）已在主线（见 memory `code-assistant-capability-pack-codex-cli-host-docker` / `docker-sandbox-openclaw-three-env`）。本轮 = **验证 / 硬化 / 置为第一优先 + 接入任务驱动**，非从零。缺真实客户原话同上。

### 2.2 真相源（谁拥有事实）
| 事实 | 真相源 | 非真相源 |
| --- | --- | --- |
| 任务（即时/阶段）状态、队列、裁定 | **07 task governance**（taskgovernance SQLite） | 06 只产 TaskCandidate；前端只投影 |
| 执行结果、hands 路由、provider 可达性 | **11 Execution Gateway** + provider registry（02） | 前端 / pack 只声明 ProviderRequirement |
| codex 沙箱执行产物 | 11 codex sidecar 回执 + 03 ledger | codex CLI stdout 非正式事实，须落候选/回执 |
| 主权裁定、对外正式生效 | **Owner + Base Gate** + 03 receipt | pack / flow / codex 都不持主权 |
| License / Entitlement / 购买 / 支付 | **truzhen-cloud**（市场/支付/授权真相源） | 前端 / 基座只经受控 API 消费 |

### 2.3 仓库 / 层归属
- **需求 1 任务驱动**：07 taskgovernance（任务真相 + 即时/阶段类型 + 队列驱动）【truzhenos】；06 sceneflow 产 TaskCandidate 接线【truzhenos】；TaskCandidate schema 若加 `task_kind: immediate|phase`【**truzhen-contracts，橙，改 schema 必 bump VERSION 铁律 7**】；env pack flow 声明各节点任务类型【truzhen-packs 数据】；任务管理前端视图【client】。
- **需求 2 codex-docker hands**：11 Execution Gateway 路由 + codex sidecar provider（Docker 执行，**Go core 不碰 docker，经 TCP sidecar**）【truzhenos + truzhen-software sidecar】；04 能力荚描述符 `software-code-assistant`（W5 裁定 = 基座常驻执行 provider）【truzhenos】；execution/任务执行前端【client】。
- **禁止**：不把任务真相塞进 06 / 前端；不把 codex docker 执行写进 Go core；不默认全塞主仓。

### 2.4 风险颜色（前置，决定本轮不是纯测试轮）
- 需求 1 = **橙 → 红**：07 驱动 → 11 真实执行；改 TaskCandidate 契约（橙）；真实对外动作过 Base Gate（红）。
- 需求 2 = **红**：真实执行引擎。**沙箱内产候选 = 留痕（免闸）**；**跨到真实对外生效 = 必过 Owner + Base Gate + Receipt**。codex 永远 Proposer——它执行「生成/草稿/沙箱内动作」，对外生效仍须 Owner 确认。
- ∴ 本轮 = **橙/红建设 + 测试**：先出契约影响 + 架构方案 → **Owner 裁** → TDD 建 → 测。**红线动作（真实对外）只在 Owner 逐次确认 + Base Gate + Receipt 下发生**。

## 3. 两新能力架构方案（P0 Owner 先审，本轮要建的骨架）

### 3.1 需求 1：任务驱动骨架（即时任务 + 阶段任务，全经 07）
- **即时任务（immediate）**：一次性、立即执行的原子任务（PDF 解析 / 监测数据拉取 / 单份文书草拟）。SceneFlow capability/execution 节点 → 生成 immediate TaskCandidate → 07 → 裁定 → dispatch → hands 执行 → receipt → 回挂。
- **阶段任务（phase）**：跨阶段里程碑任务（「完成 screening 阶段」「完成 enforcement 阶段」），可含子任务聚合、有**阶段门**（阶段完成需 Owner/Base 确认才进下一阶段）。映射 env 9 阶段。
- **驱动铁律**：**所有**执行从 07 任务出发（no ad-hoc execution）；06 不直接调 11，必经 07 TaskCandidate。静态守卫（禁 06→11 直连）+ 运行时证（每次执行带 `task_ref`）。
- **契约影响**：TaskCandidate 加 `task_kind`（immediate|phase）+ 可选 `phase_ref`/`parent_task_ref`。**additive**（老 pack 默认 immediate）；先出 contracts 影响清单 + 兼容策略，Owner 批后再改。

### 3.2 需求 2：codex-cli-docker hands（第一优先）
- **路由**：07 task dispatch → 11 Execution Gateway → hands 选择器**优先 codex CLI sidecar**（真实 readiness 探测，可达则用）；不可达 → 明确 fallback 路由证据 **或** `provider_missing`/`not_ready`，不静默降级。
- **沙箱**：codex CLI 在 Docker（**`net=none` 默认、ro 挂载、资源限**）真实 `codex exec` 执行任务；Go core 不碰 docker，经 TCP sidecar（kind=codex，复用现有范式）。
- **主权**：沙箱内产物 = 候选/草稿（**留痕、免闸**）；codex 建议的**对外正式动作**（发文书、写真实系统）= **必过 Owner + Base Gate + 11 受控执行 + Receipt**。
- **「各类任务」边界（本轮限定）**：codex 执行**沙箱内可完成的生成/分析/草拟类任务**（法律文书草拟、监测数据分析、报告生成）；**不**授权 codex 直接触发真实对外发送/写库（走网关受控执行）。§11.3 由 Owner 确认边界。

## 4. 分阶段目标（build → integrate → test → commercialize；每阶段 /goal 可判定）

| 阶段 | 内容 | 退出判据 |
| --- | --- | --- |
| **P0 契约 + 架构闸（橙/红，Owner 先审）** | 出 TaskCandidate `task_kind` 契约影响清单 + codex-first 路由方案 + 留痕/过闸边界图 | Owner 批架构 + 契约兼容策略冻结 |
| **P1 任务驱动骨架（建 + TDD）** | 07 即时/阶段任务 + 06→07 接线 + 禁 06→11 直连守卫 | TDD 绿 + 静态守卫 + env flow 每节点产 task（G1 雏形） |
| **P2 codex-docker hands（硬化 + TDD）** | 11 codex-first 路由 + docker 沙箱真执行 + readiness 三态 | 真 `codex exec` 在 `net=none` 沙箱产候选+回执；不可达诚实态（G2） |
| **P3 主权闸贯通** | 留痕/过闸边界 + Base Gate + Receipt 全链 | 对外动作零绕过、03 可反查（G3） |
| **P4 商用闭环 E2E（GUI 实测）** | 购买→授权→安装→启用→任务驱动执行→回执 | G4 绿（真 paid 或 Owner 商用验收替代） |
| **P5 REG + 证据包 + 诚实裁定** | REG-R5-01~09 无回归 + 证据包 + neuter 复核 | G5 + G6 绿，呈 Owner 商用裁定 |

## 5. 测试与验证矩阵

### 5.1 新能力（本轮建设 → 验证）
| 编号 | 对象 | 判绿标准 |
| --- | --- | --- |
| TASK-R5-01 | env flow 每正式动作节点产 TaskCandidate | 即时/阶段正确分类；07 队列可见；**零 06→11 直连**（静态守卫 + 运行时 `task_ref`） |
| TASK-R5-02 | 即时任务全生命周期 | 产→裁→dispatch→执行→03 receipt→回挂 timeline，可反查 |
| TASK-R5-03 | 阶段任务 + 阶段门 | 阶段任务聚合子任务；阶段完成需 Owner/Base 确认才进下一阶段（不自动越阶段） |
| HANDS-R5-01 | codex docker readiness 三态 | 可达→用；不可达→`provider_missing`/`not_ready` 诚实，不伪造执行 |
| HANDS-R5-02 | codex-first 路由 | codex 可达时**优先选中**且有路由证据；fallback 明确留痕 |
| HANDS-R5-03 | codex `net=none` docker 真执行 | 真 `codex exec` 产候选 + 03 回执可反查；沙箱隔离证（`net=none`/ro） |
| HANDS-R5-04 | 留痕 vs 过闸边界 | 沙箱内产物=候选免闸；codex 建议对外动作**必过 Owner+Base Gate**（不绕，双证） |
| COMM-R5-01 | 商用闭环 GUI | 购买→授权→安装→启用→任务驱动执行→回执 全链（真 paid 或 Owner 商用验收替代） |

### 5.2 REG 回归区（已 land 修复 → origin/main 现状无回归）
| 编号 | 回归对象 | land 锚点 |
| --- | --- | --- |
| REG-R5-01 | ENV-R4-01 知识 total_count / 不截断 | truzhenos `37fdd9e` |
| REG-R5-02 | ENV-R4-02 provider_family 透传 | packs `d140723` |
| REG-R5-03 | flow-stall（run-start 不被请求 deadline 截断） | truzhenos `658e495` |
| REG-R5-04 | advice node_type 四层透传（不误判 challenge） | node_type 链 |
| REG-R5-05 | R3 三修现状无回归 | `bc6f2e0`/`dca2179`/`d2b5e58` |
| REG-R5-06 | provider 机器态投影贯通 | os `b4ee07b` + client `966c2f1` |
| REG-R5-07 | 「在 03 反查」按钮默认可见 | client `966c2f1` |
| REG-R5-08 | 详情页 flow/knowledge 错误分域 | client `966c2f1` |
| REG-R5-09 | Base Gate 冷起忙态提示 | client `966c2f1` |

### 5.3 深水 + 红线（沿用 R5）
| 编号 | 对象 | 判定 |
| --- | --- | --- |
| DW-R5-01 | `screening` 缺 `pdf_parse_status` blocker（归 05） | 受控 blocked（诚实）vs 死卡（缺口）；有 PDF 样本则验 parse→解阻闭环，否则止于定性 |
| DW-R5-02 | 九阶段全链路诚实态（含 provider 机器态一致性 REG-R5-06） | `provider_missing`/`blocked`/`manual_handoff`/`not_ready` 诚实，无假成功/自铸 ref/假读 |
| RED-R5-01 | P4-09 引擎不消费 GatePolicy 越 owner_gate | **最高优先**，复现出证据 + 影响清单，**不修**（改主权链=红，归 Owner） |

## 6. 四画像（沿用 R3，GUI-only）
一线执法员（主流程 + 任务驱动）/ 挑剔律师（双角色 challenge 真实性）/ 监管审计（03 反查 + evidence + 任务链完整性）/ 门外汉 Owner（主权门控 + 阶段门是否挡得住误操作 + codex 越界是否被 Base Gate 拦）。

## 7. 治理必填（真相源/归属/风险/契约/上下文/禁止边界/验收/生命周期）
- **真相源**：任务=07；执行/路由=11；codex 产物=11 sidecar+03；主权=Owner+Base Gate+03；商用=truzhen-cloud。ReadModel/GUI/codex stdout ≠ 真相。
- **仓库/层归属**：见 §2.3（07/06/11/04 truzhenos + contracts task_kind + packs 数据 + software sidecar + client 前端 + cloud 商用）。
- **风险维度**：AI = 施工者（建骨架）+ 施工者（测）+ neuter 复核（验伪）；**红线真实对外动作只在 Owner 逐次确认 + Base Gate + Receipt 下发生**；codex 永远 Proposer。
- **契约影响**：TaskCandidate `task_kind`（additive，橙）先出影响清单 + 兼容 + bump VERSION（铁律 7）；不改 Receipt/Gate 语义。
- **上下文维度**：允许读 = env pack 全目录、06/07/11/04 相关后端、codex sidecar、client 任务/执行前端、R3/R4/本计划、REG land commit。**禁无边界扫全仓后直接施工**。
- **禁止边界**：① 不碰非本 lane 的其它会话 worktree（O9/mobile/cloud/memory/software 已移交）；② codex 不直接触发真实对外发送/写库（走网关受控执行）；③ 不把任务真相/执行写进 06 或前端或 Go core docker；④ 不改主权链语义；⑤ 不自我宣称「已发布/客户生产可用」；⑥ 撞模块体量红区（07/09/10/11）不自行 bump（呈 Owner，同 B1 先例）。
- **验收设计**：AI 设计严苛验收 + **组织者 neuter 独立复核**（revert 必 FAIL）；每 G / 每 TASK-*/HANDS-*/COMM-* 双证映射实现 + 运行时证据。
- **验收维度**：改什么证什么——任务驱动证 07 链、codex hands 证 docker 真执行 + 回执、主权证过闸零绕过、商用证闭环 E2E。不得用无关测试替代。
- **生命周期档位**：`设计中 → 契约已定（P0）→ 已实现（P1-P2）→ 已接线（P3-P4）→ 已验收（本地隔离，P5）→ 待 Owner 商用裁定`。**不由 Agent 置「已发布」**。

## 8. 隔离栈 + 证据产物（执行时填）
- 假 HOME + 专属端口 devserver（避让在用 :18095/:18096/:18109）+ 独立 vite + docker 沙箱（codex sidecar）+ env pack fresh install。
- 产物落 `/Users/li/Documents/过程文档/envpack-r5-commercial-<date>/`：user-log / behavior-matrix / ui-audit / issues / shots + 每 G 双证 + 任务链 JSON + codex 沙箱回执 + neuter 复核记录。
- **证据包** `环保pack-商用就绪证据包-<date>.md`：G1–G6 逐项证据 + 诚实上限声明 + 待 Owner 商用裁定项。报告落同目录，对话只给路径 + 结论 + 橙/红待裁。

## 9. 依赖 / 前置（/goal 开工前需就位）
- codex CLI 在本机可用 + Docker 可用（`net=none` 沙箱能起）；不可用则 HANDS-* 走 `not_ready` 诚实态并标依赖缺口。
- truzhen-cloud 市场/支付/授权受控 API 可达（COMM-R5-01 商用闭环需要）；真实 paid 走收钱吧正式商户（见 memory `shouqianba-payment-impl`）或 Owner 商用验收替代。
- env pack 知识库 752 条 FormalKnowledge（Owner 权威资料）在位。

## 10. 最小可交付（砍到多小仍有用）
- **必进本轮**：G1（任务驱动骨架，即时+阶段）+ G2（codex-docker hands 第一优先真执行）+ G3（主权闸零绕过）+ G5（REG 无回归）+ G6（证据包）。
- **可砍 / 依赖 Owner**：G4 商用闭环的**真实 paid 订单**（可先 Owner 商用验收替代，真支付另轮）；DW-R5-01 真实 PDF 闭环（无样本则止于定性）；RED-R5-01 只复现不修。
- **砍不掉的底线**：任何执行都经任务管理 + 任何对外生效都过主权闸 + 任何未接通都诚实态。**这三条破了本轮即不达标**。

## 11. 待 Owner 裁定（/goal 开工前，逐条）
1. **§11.1 战略提升**：确认环保 pack 从「护城河标杆样板（非第一现金流）」提升为**本轮商用就绪目标**？（改 2026-06-21 战略排序）
2. **§11.2 codex-hands 定位**：codex-cli-docker 作为 **truzhen 通用第一优先 hands**、env pack 作**测试载体**——确认？还是限定 codex 仅服务 env pack？
3. **§11.3 codex 执行边界**：本轮 codex 只在沙箱做**生成/分析/草拟类任务**（对外生效走网关受控执行）——确认？是否授权更宽的沙箱任务类型？
4. **§11.4 商用证据替代**：缺真实客户原话时，「正式商用」由 **Owner 商用验收裁定替代**（不自我宣称已发布）——确认？是否本轮就要真实 paid 订单闭环？
5. **§11.5 契约变更**：批准 TaskCandidate additive `task_kind`（immediate|phase）改 contracts + bump VERSION？
6. **§11.6 范围/资源**：本轮是**建设 + 测试合并的大轮**（P0–P5），预计跨多阶段；确认按 /goal 分阶段推进，每阶段 Owner 可中断校准？

---

> **本文件写完即停。Owner 将在新会话 /goal 模式开工。/goal 的成功 = G1–G6 证据齐、呈 Owner 商用裁定；不得以 mock/demo/假成功/绕过任务管理或主权闸判绿。「正式商用」是 Owner 裁定态，非 Agent 声明态。**
