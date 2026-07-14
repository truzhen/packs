# 智能家居服务商 Pack 单项目完整 lifecycle 测试计划 v15

日期：2026-07-14  
状态：`当前客户主线；v14 流式阻断终态 P2 已验收，完整 flow/Provider/作者台/按钮覆盖待验；未发布`

前序证据：

- `/Users/li/Documents/过程文档/smart-home-service-provider-v14-20260714/测试报告.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v14-20260714/问题台账.md`
- `/Users/li/Documents/过程文档/smart-home-service-provider-v14-20260714/按钮覆盖表.md`

被替代计划：`/Users/li/Documents/truzhen-packs/docs/plans/smart-home-service-provider-test-campaign-v14-20260714.md`

## 1. 派活卡与 Owner 裁定

| 维度 | v15 裁定与建议 |
| --- | --- |
| 版本 / 优先级 | 当前客户主线，只验 `smart-home-owner-pack-v0`；环保结果不冻结本计划。 |
| 我要做的事 | 从最新三仓 `main` 拉全新隔离分支，以一个智能家居项目完成直接安装、启用、全 flow、停用、历史查询、重新启用、重启、卸载、卸载后历史查询和同版本重装，并补完 v14 未覆盖的 Provider、作者工作台和所有可见按钮。 |
| 真实客户 / 场景证据 | v14 已在真实 GUI 通过 lifecycle、停用/卸载/急停零副作用、双 Tab、三视口、断连恢复，并验收急停阻断后流式占位终止；未覆盖全部 flow 节点、Provider 各诚实态、作者工作台和全部按钮。 |
| 最小可交付 | 同一 Pack v1.0.0、digest、安装实例和 `SH-C01` 贯穿三轮；通过 manifest/schema、Candidate/Formal、最小权限、Gate/Gateway/Receipt、幂等重试、断网恢复与卸载隔离硬门；在 L2 隔离测试 Provider/测试设备上连续完成至少 10 次“Candidate→确认→执行→Receipt”闭环；disabled/uninstalled/emergency 零业务副作用；历史项目始终可查询且版本固定。 |
| 本轮砍掉项 | `提交审核`、`模拟支付`、`Entitlement`、`下载`直接跳过并记 `not_tested_owner_directed`；不启动 cloud、不登录市场、不建商品/订单/授权。 |
| 真相源 | Pack/digest 归本仓；14 持 lifecycle；05/06 持项目/Run；13 持会话；01 持急停；07/08/03 持任务候选、模型、回执；11 持执行/Frappe；client 只展示。 |
| 仓库 / 层归属 | OS 负责 lifecycle、停用/卸载/急停硬门；client 负责流式终态与投影；本仓负责 Pack 与测试计划。contracts、cloud、software 不改。 |
| 风险 / AI 角色 | 黄：安装、启停、候选；橙：Provider；红：Gate、Receipt、真实发送/执行/Frappe。AI 只产 Candidate，Owner + Base 保留正式裁定。 |
| 契约影响 | 零 contracts/DTO/schema/Gate/Receipt/Candidate/ReadModel/Surface 变更；本轮只验收现有实现。 |
| 上下文 | 只读本计划、Pack、v14 证据、隔离运行库、合成项目数据和相关 GUI；不读真实客户数据或 secret。 |
| 禁止边界 | 不改库、不正式化、不真实发送/执行/Frappe/付款、不用 mock、旧截图、API-only 或候选态冒充最终 PASS。 |
| 变更影响 | 单 Pack lifecycle、项目上下文对话、flow/Provider/作者工作台/按钮覆盖和证据文档；不改变业务主权。 |
| 生命周期档位 | 流式 P2 为`已验收`；Pack 全程须经 v15 与独立复核才为`已验收（打包前）`；发布另行裁定。 |

## 2. 执行权限、隔离和自治

执行者获本计划范围内全部本地权限：隔离 worktree、专用 HOME/SQLite/浏览器 profile、真实本地模型、服务和端口、直接 lifecycle、合成项目数据、截图、只读对账、中小问题根因修复和回归，以及经 Owner + Base Gate 后在 L2 隔离测试 Provider/明确标识测试设备上执行本计划的 10 次闭环。不得触发真实市场、付款、对外发送、生产 Frappe、生产设备、客户设备或生产数据。

隔离分级必须写入每条测试证据：纯知识与 manifest 声明按 L0 只读校验；秘书/项目经理等 Agent 推理按 L1 sandbox，不能持设备凭据；设备驱动、网络、本地软件和 Frappe 等调用按 L2 隔离 Provider，经 Gateway 调用并留 Receipt。容器、进程或端口“启动成功”只证明载体存活，不能替代 readiness、Gate、Gateway、Receipt 或设备结果验收。

每仓从固定主仓最新 `origin/main` 建分支并独立记录 status/HEAD/WIP/验证。P2/P3 和影响清晰、可逆、不改契约/主权链的中等问题必须现场修复、补测试、从最近可靠点继续。P0/P1 记录并冻结受影响路径，但其余安全 lane 继续；只有继续会造成真实副作用或全部剩余步骤依赖失真状态时暂停。不得因单点问题笼统停止 R2/R3。测试任务不自行 merge/push，除非 Owner 另行授权。

## 3. 唯一项目和直接安装口径

唯一主项目：`SH-C01-v15-王先生灯组离线售后`。同一 `transaction_ref`、Pack v1.0.0、digest、安装实例、浏览器 profile 贯穿启停、重启、断网、急停、卸载、重装和历史查询。只有 05/06 真相被不可逆污染时可建立一次 `SH-C01-RECOVERY-01`，并登记原因与 refs。

1. 从 `/Users/li/Documents/truzhen-packs/smart-home-owner-pack-v0` 构建净制品并登记清单、SHA-256 和禁品扫描。
2. 跳过市场四阶段，直接用 `install.py` 调隔离 OS 真实 lifecycle；如有本地文件安装 GUI，优先 GUI 触发。
3. `uninstall.py` 只撤销当前运行访问，不删除项目、Run、Candidate、Receipt 或历史版本引用。
4. lifecycle 结果必须由 GUI、14 ReadModel 和 03 Receipt 三方对账。

## 4. 共通治理硬门：Pack 不得腐化基座

以下为本 Pack 的放行前置条件，任一项不满足均不得用后续 GUI 成功抵消：

1. `manifest.json`、flow、role、capabilities 和 software/provider requirements 必须通过权威 schema 与结构校验；不得在 Pack 内复制、放宽或私造 contracts schema，也不得以解析成功代替语义校验。
2. Candidate 与 Formal 必须物理、接口和 UI 三层隔离。Role Pack、模型和 Pack flow 只能先产 Candidate；只有 Owner 明确确认并通过 Base Gate 后，才可生成本轮受控 Formal/ExecutionIntent。Candidate、草稿和候选回执不得显示为已执行或正式完成。
3. 权限取 Owner 授权、角色白名单、Pack 声明、项目 ContextSlice 与 Provider 能力的最小交集。Pack 不得申请通配设备/网络/文件权限、持设备密钥、跨项目读写或绕开项目绑定。
4. 缺配置、缺证据、缺 Provider、Provider 未 ready、Pack/Provider/设备协议版本冲突时，必须明确返回 `blocked`、`not_ready` 或 `provider_missing`，并给出稳定 reason/missing fields；版本冲突统一为 `blocked(reason=version_conflict)`。不得静默降级、换设备、换 Provider、吞错、补假状态或返回空成功；`degraded` 不得掩盖上述硬缺口。
5. 每次受控设备动作必须证明 Candidate → Owner 确认 → Base Gate 裁定 → Execution Gateway → L2 Provider → Receipt；Provider `ready`、容器存活、前端提示或 Pack 声明均不能代替 Gate 和实际设备结果。
6. 每个接受、阻断和执行尝试都要有可反查、可回放、可关联项目/设备/动作/Gate 的 03 Receipt。Receipt 回放必须复现当时 Pack/Provider/设备版本、输入摘要、裁定、实际结果和副作用计数，不得从当前状态重算历史。
7. 安装、启停、命令提交、确认、执行、重试、重复点击、旧请求重放和卸载必须使用稳定幂等键。断网前已提交/未提交边界要可判定；恢复后不得重复 Candidate、FormalTask、设备动作或 Receipt。
8. 卸载后本 Pack 的 RoleBinding、Provider session、订阅、轮询、缓存与临时资源必须释放；项目和 Receipt 只读保留。不得控制其它 Pack/项目的设备，不得串用设备绑定、Provider session、端口、缓存、任务或回执。
9. 对基座做静态防腐审计：`truzhenos` 不得出现 smart-home/pack ref/具体设备型号/flow node 专用分支、专用 Gate/Gateway/Receipt、行业 schema 副本或 Pack seed 回潮；client 不得成为设备状态真相源；本仓不得包含驱动或 Provider 实现。若只能靠 Pack 专用基座代码跑通，登记跨仓大问题并停止该 lane，不得现场焊死。

## 5. 开跑恢复硬门

1. Pack JSON/Python/结构/禁品/digest 与 `git diff --check` 全绿；registry、项目、任务、模型、Receipt 基线有快照。
2. client 定向回归证明急停/停用失败终态原位替换“正在流式生成回复...”，不留 loading/streaming 卡；旧占位不污染下一请求。
3. OS 定向回归证明 disabled/uninstalled/emergency 在 dialogue/stream/07/03/08/turn 前 fail-closed；伪 receipt、双击和旧请求不得写回。
4. client 全量测试、typecheck、build，OS 相关 test/race/vet/EGR 和 Pack 文档/结构审计全绿后启动 GUI。

## 6. R1：安装、启用和唯一项目全 flow

过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v15-20260714/`。

1. registry=0 起步；直接安装后验证 registry 0→1、enabled v1.0.0、Role Pack/SlotBinding、安装 Receipt 和 ProviderRequirement 投影。
2. GUI 创建唯一 SH-C01 和同一 06 Run，绑定两台名称相近但物理标识不同的 L2 隔离测试设备；逐项覆盖 `intake`、项目对象、Frappe 项目/客户快照、项目经理 advice、里程碑/采购/施工任务候选、客户沟通草稿、Frappe 写回候选、Owner Gate、Receipt archive、done。每个节点标 `PASS / 受控阻断 / provider_missing / not_ready / not_tested`。
3. 项目上下文页分别执行一条流式和一条非流式新请求；每条最多一次真实本地模型、一个待 Owner Candidate。完成或阻断后不得残留“正在流式生成回复...”或“秘书正在思考中”。
4. 对全部可见项目按钮逐个操作：正常入口、重复点击、取消、刷新、返回、技术详情、Receipt 反查、候选重试；不能操作的写明权限/前置条件，不得静默漏测。
5. 执行前确认页必须清楚展示设备名称、不可混淆的设备 ID/位置、动作、参数、预期影响、影响范围、风险、超时、回滚/停止方式和 Provider 状态；Owner 未确认、设备绑定不唯一或展示信息不全时必须 `blocked`，不得发出设备指令。
6. 在同一 SH-C01 内连续完成至少 10 次“Candidate → Owner 明确确认 → Base Gate → Execution Gateway → L2 测试 Provider/测试设备 → Receipt”闭环。每次使用唯一业务意图和稳定幂等键，保存 GUI、Gate decision、Provider request/result、设备前后状态和 Receipt 回放；10 次中任何一次靠脚本直调、mock、容器状态或手改数据库完成，整项不通过。
7. 10 次闭环必须覆盖：设备离线、云/本地状态漂移、重复指令、重复确认、错误设备绑定、Provider 超时、执行结果迟到、断网后恢复和人工急停。离线/错误绑定/超时应受控阻断或得到诚实失败 Receipt；状态漂移必须在执行前重新确认，不得以陈旧 UI 状态继续。
8. 重放同一幂等键不得二次执行；使用新幂等键的同动作必须重新走确认和 Gate。人工急停须在指令进入 Provider 前阻断；若执行中触发，必须停止后续动作、记录已发生影响并生成可回放 Receipt，解除急停后旧命令不得自动续跑。
9. 除上述 10 次 L2 隔离闭环外，后段只验证 Candidate、诚实 Provider 态和 Gate 阻断；保持 0 对外 Send、0 生产 Execution、0 Frappe 写回、0 正式 Memory。正式任务/执行记录仅允许属于这 10 次且必须逐条有 Owner/Base/Gateway/Receipt 全链证据。

## 7. R2：停用、历史查询、重新启用和重启

1. 停用 Pack 后从同一 SH-C01 攻击流式、非流式、tool confirm、伪 receipt、旧 request、双击和双 Tab；模型、turn、Candidate、Receipt、Task、Run、Memory、Send、Execution/Frappe 业务增量全部为 0。
2. 停用态从项目清单、详情、对话历史、时间线、候选和 Receipt 入口查询 SH-C01；项目、Run、候选、Receipt 和 v1.0.0 只读可见，不丢失、不串项目、不自动迁移。
3. 重新启用并重启 OS/client；刷新后仍恢复同一当前项目，只有全新请求恢复一次模型/候选，停用期请求不补写。
4. 若全局沟通中心没有项目上下文，必须明确显示无项目上下文；Pack 主链固定从项目上下文对话页进入，不把错误上下文产生的阻断当作停用 PASS。

## 8. R3：卸载、重装、急停、Provider 和作者台

1. 卸载后 enabled 指针、SlotBinding 和运行访问撤销；重复卸载幂等。卸载态继续查询同一 SH-C01 历史，新请求全链零副作用阻断。
2. 同 digest 重装后 registry 仅一个版本，不复制项目/候选/Receipt；SH-C01 保持 v1.0.0，新请求恢复，卸载期请求不补写。
3. 顶部急停浮层保持打开：启用急停后单/双 Tab、流式/非流式、旧 request、伪 receipt、重复点击均零副作用且流式卡立即终止；显式解除后只有新请求恢复。
4. 对 Frappe snapshot/read/write 和其它 ProviderRequirement 分别验证可安全诱发的 `ready/not_ready/provider_missing/blocked/degraded`；不能诱发的写 `not_tested`，不伪造状态、不真实写回。
5. 进入 Pack 作者工作台只读核对 manifest、flow、角色槽、capabilities、ProviderRequirement、版本和导出/预检状态；逐个覆盖全部可见按钮，但跳过提交审核和发布。
6. 完成断网/恢复、浏览器刷新、390×844 / 1024×768 / 1440×900、历史回放、双 Tab 和重复点击竞态；逐项更新 v14 按钮覆盖表，不因已通过旧轮而免测关键回归。

## 9. 证据、验收和收尾

每个关键动作保存 GUI、网络状态、01/03/05/06/07/08/10/11/13/14 前后计数、refs、截图和时间戳；API/SQLite 只作对账。交付 lifecycle 时间线、全 flow 节点表、按钮覆盖表、Provider 状态矩阵、三视口/双 Tab/断网记录、问题台账、截图/Receipt 索引和端口清理证据。

独立验收主体确认开放 P0=0、阻塞 P1=0；共通治理硬门全部有证据；流式阻断终态不回潮；设备离线、状态漂移、重复指令、错误绑定、超时和人工急停均符合故障断言；10/10 受控执行闭环可逐条回放；单项目安装→停用→启用→卸载→重装及历史查询全程通过；全部安全可继续 lane、所有可见按钮、作者工作台和 Provider 诚实态已覆盖或明确 `not_tested`。达到后才标`已验收（打包前）`；四个市场阶段始终为 `not_tested_owner_directed`。
