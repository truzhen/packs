# truzhen-packs 模块与包清单

本文件记录 `/Users/li/Documents/truzhen-packs` 当前包模块、职责边界、成熟度和协作关系。完整纪律见 `AGENTS.md`。

## 1. 仓库总览

`truzhen-packs` 是 Truzhen 包层资产仓。每个正式 pack 是一个自包含文件夹，通过 `manifest.json`、flow、role-slots、role-packs、capabilities、knowledge 和 install / uninstall 脚本声明领域治理资产，并经基座真实 lifecycle 装入。

核心链路：

```text
Pack 声明 -> 基座 lifecycle enabled -> Scene Runtime 解释 -> Candidate -> Owner + Base Gate -> Gateway -> Receipt
```

Pack 不直接生成 Formal Record，不直接执行 provider，不直接写基座数据库。

## Truzhen 底层逻辑指向

权威总纲：`/Users/li/Documents/truzhen-contracts/TRUZHEN_PHILOSOPHY.md`（远端 `github.com/truzhen/contracts` 根目录同名文件）。本仓是总纲里的“人类经验蒸馏与售卖平台”落点：Pack 把行业高手的经验、红线、流程和能力需求写成可安装边界，让 AI 在边界内生成候选，而不是让 AI 自己拥有主权。

### 本仓映射表

| 总纲原则 | 本仓落点 | 边界 |
| --- | --- | --- |
| Pack 是人类提供主权的工作流边界 | `manifest.json`、flow、role-slots、role-packs、capabilities | Pack 只声明和编排，不裁定、不执行、不写正式事实。 |
| 红黄绿按责任后果划，不按技术难度划 | `person_strategy`、`formalization_requirement`、`gates`、risk / fallback policy | 行业红线由真实场景和行业作者经验决定，最终仍受 Owner + Base Gate 约束。 |
| Pack 边界从真实客户中长出来 | Pack README、manifest、验收报告、`FEATURE_LEDGER.md` | 不先在纸面设计完整职业宇宙，再找客户验证。 |
| Provider 是能力需求，不是 Pack 内实现 | `provider_requirements`、`capabilities/capabilities.json` | Frappe、OCR、IM、Baserow 等只能声明需求；真实 provider 归 `truzhen-software` / 基座 Gateway。 |
| 底层软件是可解析需求，不是 Pack 资产 | `software_requirements` | Pack 只声明 Baserow / OCR 等软件 family、版本范围、capability、reuse / isolation / fallback / license 策略；复用、安装、版本冲突、隔离和 lock 归 `truzhenos` resolver。 |
| AI / Role Pack 永远是 Proposer | `role-packs/*.json`、多角色对照、质询节点 | 角色只提建议、草稿、质询和风险，不批准、不发送、不执行。 |
| 可售卖的是可复用经验资产 | folder pack、知识域、护城河、市场元数据候选 | 正式商品、价格、订单、License / Entitlement 和分发状态归 `truzhen-cloud`。 |
| 行业语义从真实场景长出来 | 业务对象 schema / relations、`knowledge/` glossary / SOP / checklist / index、flow / gates、role / capability refs | 不建中央语义层；不从数据表或 Provider 反推行业定义；缺客户证据只登记缺口，不升级成熟度。 |

### 禁止误读清单

- 不得把总纲读成“Pack 越大越完整”；抽象层级高不是成就，真实客户用起来才是成功标准。
- 不得发明第四种 Pack；新需求先判断归 Domain Work Pack、Capability Pack、Role Pack，还是应迁往基座、contracts、software、cloud 或 client。
- 不得把 Pack manifest、flow、角色包或 install 脚本写成 OwnerDecision、Base Gate、Receipt 或正式执行结果。
- 不得让 AI 未经 Owner 确认就使自己起草的风险分级、治理规则或 manifest 变成有效规则。
- 不得把 provider `ready`、脚本存在、mock 成功或说明文字当成真实产品接线；未接通只能写 `blocked / provider_missing / not_ready`。
- 不得把高风险法律、财务、医疗、合同等知识写成正式适用结论；默认要有来源、人工核验和责任边界。

## 2. 当前包清单

| 包 | pack 标识 | template_family | 成熟度 | 职责 |
|---|---|---|---|---|
| `environmental-enforcement-pack-v0/` | `scene_pack://environmental-enforcement-flow` | 合规审查执法证据链型 | 完整文件夹包；材料硬门、事实锚定、检索相关性与 Flow/Transaction 完成语义为`已实现 -> 已接线`，待 v18 按三标准单项目复验（未发布） | 生态环境执法证据链：线索、立案、取证、证据三性、执法精英 / 挑剔律师对照、处置、文书、Owner/Base 裁定、Receipt。 |
| `smart-home-owner-pack-v0/` | `scene_pack://smart-home-owner-project-ops` | 长周期项目交付型 | v1.1.0 完整文件夹包；v17 后端/Provider 单项目 11 次闭环通过，GUI 待 v18 在允许的受控会话按三标准补验（未发布） | 智能家居老板项目经营：商机、立项、进度、物料、交付、历史查询、Frappe 受控写回；硬件仅可选复用 Home Assistant L2 Provider，不自造程序且不作为项目主链放行前提。 |
| `housekeeping-ops-pack-v0/` | `scene_pack://housekeeping-ops`，兼容 `pack_housekeeping_ops_v0` | 客户服务全生命周期型 | 可装入 / 可卸载文件夹包：manifest、flow、capabilities、2 角色包、role-slots、install、uninstall；`knowledge/` 待补 | 家政客户服务全生命周期：受理咨询、顾问出方案、质检质询、对照确认门、排期报价、派工确认、通知客户、上门服务、服务回执、归档。 |
| `content-operations-workbench-v0/` | `scene_pack://content-operations-workbench` | Founder 自营内容候选与复盘型 | v0.1.1 已完成打包前验收：业务 Skill、输出契约和按三种工作模式选择的模型 Schema 均由 framed hash 锁定；真实 GUI 已生成完整 45 秒抖音候选，并可反查 T06、动态 Gate、08 usage、11 执行与 03 Receipt；v0.2.0 真实视频生成升级处于`设计中`，尚未接 Provider、OS 或 Client；代码发布不等于产品已安装、启用或上架 | 把真实产品证据变成方向候选、母内容、渠道候选、人工发布包与周复盘；0.2.0 计划增加本地可播放 MP4 候选，但 Pack 仍不含 CLI / Provider / 平台登录 / 自动发布，真实渲染由 software + 11 Gateway 供给。 |
| `shuxuejia-renovation-pack-v0/` | `scene_pack://shuxuejia-large-home-renovation` | 长周期项目交付型 | 完整文件夹包；2026-07-08 已完成隔离 install / uninstall / reinstall、457 / 457 SceneFlowRun 与前端用户视角检查，生命周期`已接线`（未发布）；任务、知识、沟通与 ProviderRequirement 仍有前端投影缺口 | 墅学家大宅装修设计指导：设计准备、深化设计、材料合同付款、现场质量验收、售后保修，多角色候选协作，Owner/Base Gate 和 Receipt 回放。 |
| `capability-pack-candidates/short-video-ops-v0/` | `capability-pack-candidate-set://short-video-ops-v0` | 短视频运营能力候选集 | 候选资产：3 个 Capability Pack 候选 JSON、OSS 证据矩阵、Code Assistant 调用候选台账、运行请求候选台账、candidate bundle 导出台账、candidate bundle dry-run 台账、PatchCandidate 承接台账、PatchCandidate 复核台账、P11 lifecycle preflight 执行规格、P12 开跑前门禁候选、P12 / P13 / P15 / P16 / P17 / P18 运行后证据验收门候选、P12-P18 执行规格 / 授权卡 / 机器证据契约 / 覆盖 `evidence_id` 的证据台账骨架 / 执行就绪证据写回计划、带写回摘要的商用跨仓执行队列、带写回总计并消费 P12 / P13 / P15 / P16 / P17 / P18 后验收门的商用 readiness verifier、消费写回门禁和 P12 / P13 / P15 / P16 / P17 / P18 后验收门的商用 go/no-go 候选门禁、商用后验收门覆盖 verifier、商用禁入动作覆盖 verifier、商用执行就绪守卫覆盖 verifier、消费商用后验收门 / 禁入动作 / 执行就绪守卫覆盖 verifier 的目标完成证据地图候选、P12-P18 商用证据契约总索引和逐分片签字矩阵、前后端证据交接 runbook / 命令计划、GUI/API/Receipt 追踪矩阵、商用独立验收签收矩阵、P14 商用 readiness 审计、商用改进清单、商用最终证据包模板、授权路线图、商用缺口台账、源材料占位；不安装、不启用、不发布 | 用 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 压测能力包制作台的 GitHub OSS 证据、代码助手胶水边界、11 run 请求候选、候选 bundle 导出、candidate bundle dry-run 阻断、PatchCandidate 结果承接与复核、lifecycle preflight 商用缺口、P12 安全样本 lifecycle 待授权项、P13 GUI lifecycle 面板待授权项、P14 商用缺口审计、P15 GUI 实操待授权项、P16 受控 run 待授权项、P17 provider 归属待授权项、P18 云市场 sandbox 待授权项和候选 Pack 产出；当前明确 not commercial ready。 |
| `role-pack-candidates/team-office-v0/` | `role-pack-candidate-set://team-office-v0` | 团队办公室角色候选集 | 候选资产：秘书长 + 五顾问 6 个 Role Pack 候选 JSON、团队角色槽、团队设置绑定候选、安装后团队设置角色目录候选、六角色能力引用样例、秘书长音色 / VRM asset ref、秘书长音色 / VRM 授权证据候选、云商品草稿候选（含买家可见六角色商品组件清单）、artifact manifest、sandbox 商品化流程候选、商品化交接候选、License / Entitlement 策略候选、订单 / 支付状态机候选、售后 / 退款 / 撤回策略候选、买家已购库 / 安装状态候选、发布者账号 / 结算 / 发票策略候选、商业条款 / 隐私 / 数据边界候选、发布候选包、下载 / 安装访问控制矩阵候选、确定性打包规格候选、候选包整体 digest、安装兼容性矩阵、云市场审核提交候选、商品化上线批准候选、生产发布晋级门候选（含下载 / 安装访问矩阵硬门）、商品化回执链候选、商业分发回执记录规格候选、安装预检请求候选、运行使用候选、用户视角 GUI 场景候选、用户视角 GUI 执行脚本候选、GUI 证据采集协议候选、上传前 secret / raw asset 扫描候选、商品化完成证据矩阵、目标完成证据地图候选、角色制作台问题台账回填门槛、商品化执行包、E2E 证据运行记录包、商品化证据门槛候选、sandbox 商品化环境就绪候选、商品化监控 / 诊断候选、正常商品化完成审计候选、正常商品化链路核验器候选、商品化 go/no-go 门禁候选、商品化 readiness verifier 候选、P0-P11 阶段覆盖矩阵候选、计划 TC 测试用例覆盖矩阵候选、独立验收签收矩阵候选、P0-P11 商品化阻塞清单候选、商品阶段前后端收口报告候选、P11 正常商品化通过门候选、P11 正常商品化验证记录模板、P11 证据接入绑定器候选、P11 sandbox 执行 runbook 候选、P11 sandbox 开跑前预检门候选、P11 sandbox 执行请求候选、P11 证据验收检查清单候选、P11 商品化最终证据包模板（含云上传商品草稿报告、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告）、商品化前后端 API 契约候选、跨仓执行卡（含生产晋级控制施工卡）、跨仓执行就绪包（含生产晋级阶段）、Owner 授权证据接入口候选、商品化跨仓执行队列候选（含生产晋级控制阶段）、商品化跨仓证据台账候选（含生产晋级回执写回）、前后端接线映射（含 P11 GUI/API traceability）和台账；不安装、不启用、不发布 | 用团队办公室角色制作、能力引用、团队设置替换、云端商品草稿、sandbox 购买 / 下载 / 安装设计压测角色制作台与商品化链路；真实 GUI、云端、支付、安装均待跨仓接线。 |
| `templates/` | — | — | 作者端脚手架 | Pack 作者工程模板，不参与 enabled pack 分发。 |

短视频运营 Capability Pack 候选集补充包含商用授权尝试覆盖 verifier，用于证明被拒绝的 P12 授权尝试已被商用门禁文档消费且不能当作 Owner 授权；另含商用下一授权启动守卫，用于集中校验 P12 授权卡原话、允许仓、禁改仓和禁入动作，未通过前不得开始跨仓执行，也不得通过 P12-P18 分片级 execution readiness gate、商用证据索引完成门、跨仓队列首跑门、readiness、go/no-go 或 Owner 签收；另含 Pack Studio issue 回填台账候选，用于把 P10-P18 制作台缺口、阻断原因、归属、验收证据、权威关闭证据槽和禁止绕过方式集中登记，并被 P12 开跑前门禁、P12-P18 execution readiness、商用跨仓执行队列首跑门和各分片队列镜像 gate、商用执行就绪守卫覆盖 verifier、商用当前状态审计、前后端验收契约、前后端交接、backlog、readiness、go/no-go、最终机器证据包和目标完成图消费；每个 issue 关闭前必须有实现 / 授权证据、验证命令、GUI/API/ReadModel 证据、Receipt / Candidate 引用和禁入方案缺席证明，当前仍 `pack_studio_issue_ledger_open`；另含 `capability-pack-candidates/short-video-ops-v0/docs/commercial-go-live-evidence-package.json` 机器摘要候选，用于把商用最终证据包模板的 P12-P18 分片矩阵、硬门来源、证据写回门、Owner/Base Receipt 门和非充分证据清单转成可校验资产；另含 `capability-pack-candidates/short-video-ops-v0/docs/commercial-machine-go-live-evidence-package-coverage-verifier.json` 覆盖 verifier，用于证明该机器摘要候选已被 readiness、go/no-go、证据索引、当前状态审计和目标完成图消费；当前仍 `not_commercial_ready`。

短视频运营 Capability Pack 候选集补充：目标完成证据地图、商用最终机器证据包、readiness verifier 和 go/no-go 门禁已明确拒绝覆盖 verifier、机器证据包摘要、go/no-go 候选或前后端验收契约替代真实 runtime receipts；目标完成前必须绑定前后端 runtime 回执、P12-P18 分片回执、覆盖 verifier 背后的终态证据和 Owner/Base Gate 回执。

短视频运营 Capability Pack 候选集补充：前后端商用验收契约已新增 GUI/API/Receipt 追踪矩阵，候选 bundle、delivery 状态、readiness issue、promote/confirm 控件、商用 verifier 面板和三候选 GUI walkthrough 必须逐项对上后端 endpoint / ReadModel 与 Receipt / Candidate 引用；商用 readiness、go/no-go、最终机器证据包和目标完成图均已消费 `gui_api_receipt_traceability_verified`，截图-only、后端测试-only 或网络摘要-only 均不得计入前后端验收或商用完成。

短视频运营 Capability Pack 候选集补充：新增商用独立验收签收矩阵候选，用于把 P12-P18 分片、GUI/API/Receipt 追踪、Pack Studio issue 关闭证据、禁入动作终态、证据写回、Owner/Base Gate 回执和最终机器证据包逐项交给独立 reviewer 复核；商用 readiness、go/no-go、最终机器证据包和目标完成图均已消费 `independent_acceptance_signoff_matrix_passed`，自测通过、绿色测试、人工摘要或 Owner 授权本身都不能替代独立验收证据。

短视频运营 Capability Pack 候选集补充：P16 受控 Code Assistant run 后验收门已明确拒绝静态胶水代码、本地 Codex CLI 登录态、手工 PatchCandidate 文件或 P8 run request candidate 替代 11 Gateway run receipt；P16 完成和 P17 解锁必须绑定 `DecisionRef`、`RunID`、`Nonce`、`ReceiptRef`、隔离 PatchCandidate artifact 和 no-auto-apply 证据。

短视频运营 Capability Pack 候选集补充：P17 provider / adapter candidate 后验收门已明确拒绝 packs 内 scaffold、未绑定 readiness receipt 的 provider manifest、一句 provider ready 声明或 vendor 第三方 OSS 源码替代外部 provider 仓证据；P17 完成和 P18 解锁必须绑定外部 provider 仓 diff / artifact、provider readiness receipt、packs 无 provider runtime 扫描和 Owner/Base 回执。

短视频运营 Capability Pack 候选集补充：P18 云市场 sandbox 后验收门已明确拒绝 packs listing draft、sandbox runbook、License / Entitlement 文案、人工订单状态或生产发布声称替代 `truzhen-cloud` receipt；P18 完成、商用签字和 go-live 解锁必须绑定 cloud sandbox receipt、entitlement receipt、packs 无云市场真相扫描和 Owner/Base 回执。

团队办公室 Role Pack 候选集补充包含商品化 API 请求 / 响应样例候选，用于前端、后端、云端和验收智能体按相同端点字段核对；前后端接线映射已新增 `gui_api_traceability_matrix`，把 P11 用户视角 GUI 脚本中的导出、云上传、审核、购买、下载、安装、团队替换、运行使用和负例逐项映射到商业 API 端点、回执字段、阻断态和 GUI 证据槽，P10 收口报告、P11 最终证据包、P11 验收清单、readiness verifier、go/no-go、商品完成矩阵和目标完成地图均要求 `gui_api_traceability_matrix_verified`；商品化跨仓执行队列和证据台账已要求 `frontend_user_view_gui_flow` 阶段交付该证明并写回 P11 最终证据包、P11 验收清单和 P10 收口报告；另含 artifact manifest 闭合门候选，用于上传、签名下载和本地安装前核对当前目录、候选集和 artifact manifest 文件集合一致；另含安装到运行使用激活映射候选，用于把下载回执、安装预检、启用版本、团队设置目录刷新、Owner Gate 绑定和运行使用烟测串成可验链路；能力 Pack 角色引用样例覆盖秘书长和五顾问全部 6 个制作角色，并要求安装后引用证据门，只有 entitlement 校验、安装回执和 enabled role version 齐全后才允许引用已购角色；云商品草稿含买家可见的六角色商品组件清单，发布包和云市场审核提交均引用同一清单，秘书长音色 / VRM 保持 asset-ref-only；秘书长表现层偏好另声明团队设置页 GUI 选择、清空、恢复默认和 provider readiness 显示规则，并把 `secretary_appearance_gui_controls_verified` 纳入完成声明前置证据；License / Entitlement 策略和安装预检要求 team-scoped entitlement 的授权团队必须匹配安装目标团队，跨团队安装必须 blocked；购买前必须展示支持入口、退款条款、撤回说明和候选角色免责声明，缺少披露或买家确认时 sandbox 购买必须 blocked；P11 商品化最终证据包模板要求六角色安装目录、启用版本、slot 映射、团队设置截图、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告和下载 / 安装访问矩阵报告齐全后，才允许进入 go-live 判断；该证据包还要求用户视角 GUI 执行脚本全部 14 个步骤都有截图、页面状态、操作日志和候选 / 回执引用，缺一步都不得通过 `all_user_view_gui_steps_covered`；商品化 readiness verifier 和 go/no-go 门禁已直接消费秘书长音色 / VRM 证据报告，完成声明前必须证明 GUI 选择、asset ref、授权回执和 raw asset 缺席扫描齐备；商品化 readiness verifier、go/no-go 门禁、产品完成矩阵和目标完成地图已直接消费角色制作台问题台账，完成声明前必须证明 `issue_ledger_all_entries_triaged`；商品化 readiness verifier、go/no-go 门禁和目标完成地图已直接消费商品阶段前后端收口报告，完成声明前必须证明 `product_stage_frontend_backend_closure_report_verified`；商品化 readiness verifier、go/no-go 门禁、产品完成矩阵和目标完成地图已直接消费安装后团队办公室运行使用烟测，完成声明前必须证明 `team_office_runtime_usage_smoke_verified`；并新增阶段依赖报告 section，把执行队列 phase handoff、前置回执和缺失阻断状态交给 Owner 与独立验收主体复核；P11 证据接入绑定器和验证记录模板已新增阶段依赖接入 / 填写字段，真实执行证据必须落到 `p11_phase_dependency_result` 后才能进入最终证据包；P11 证据验收检查清单已新增阶段依赖验收门，完成声明前必须有执行队列 phase link、前置回执、`bundle_tree_sha256` 和 `correlation_id` 的连续证明；商品化上线批准候选已直接消费 P11 证据验收清单、最终证据包、阶段依赖验收门和下载 / 安装访问控制矩阵，缺少清单 verified 证据或未购买 / 退款撤销 / 版本下架 / hash 不一致 blocked 回执时不能进入 go-live approval；目标完成证据地图已直接消费商品化上线批准和 P11 证据验收清单，完成声明前必须同时有 `commercial_go_live_approval_verified` 与 `p11_evidence_acceptance_checklist_verified`；商品化 readiness verifier、go/no-go 门禁和产品完成矩阵也已直接消费 `commercial_go_live_approval_verified`，缺少上线批准候选、P11 证据验收清单、最终证据包、Owner go/no-go 和生产发布 / 真实支付阻断证据时不得通过正常商品化；P11 sandbox 执行请求要求导出、云上传、审核、sandbox 购买、下载、安装、团队绑定和负例阶段逐项绑定 GUI 证据协议、权威回执和同一 `bundle_tree_sha256` 关联键，并把 `p11_gui_receipt_continuity_verified` 纳入完成声明前置证据；P11 sandbox 执行请求另声明阶段依赖和前置回执传递契约，后一阶段必须引用前一阶段回执，不能跳过上传、审核、entitlement、下载或安装直接声称商品化完成；商品化跨仓证据台账镜像 P11 阶段依赖，真实写回时每个 phase link 必须记录依赖 phase、前置回执和缺失阻断状态；商品化跨仓执行队列也镜像 P11 阶段依赖，授权后执行时缺少前置回执的 phase 必须保持 `blocked_previous_phase_evidence_missing`；正常商品化链路核验器把 P11 执行队列阶段依赖证明列为通过条件，不能只凭最终支付、下载或安装回执跳过中间阶段；商品化 go/no-go 门禁把 P11 执行队列阶段依赖证明列为终局门槛，Owner go/no-go 前必须证明上传、审核、支付授权、下载、安装、绑定连续；P11 正常商品化通过门、链路核验器和 go/no-go 门禁要求版本下架 / release revoked 时下载与安装均有 blocked 证据；商品化 readiness verifier、go/no-go 门禁、上线批准候选和生产发布晋级门已直接消费下载 / 安装访问控制矩阵，完成声明、go-live approval 和生产签名下载启用前必须证明未购买、退款撤销、授权过期、版本下架 / 撤回和 artifact hash 不一致均有权威 blocked 回执；生产发布晋级门另已直接消费 P11 最终证据包和安装后团队办公室运行使用烟测，真实支付、生产签名下载、正式上架和安装分发前必须同时证明 `p11_go_live_evidence_package_verified` 与 `team_office_runtime_usage_smoke_verified`；商品化 API 契约和请求 / 响应样例已暴露 `blocked_p11_evidence_package_incomplete` 与 `blocked_team_office_runtime_usage_smoke_missing`，前端必须展示可读阻断原因且后端必须返回候选回执；商品化回执链已要求 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、安装后启用版本和团队 slot 映射从制作台导出贯穿到团队办公室运行使用；另含角色制作台导出来源证明候选，用于阻断手工 JSON、后端直调或无 GUI 导出回执的云上传。

团队办公室 Role Pack 候选集补充：P11 最终证据包的 `install_catalog_and_slot_mapping_report` 已被 readiness、go/no-go、产品完成矩阵和目标完成地图直接消费为 `six_installed_role_catalog_evidence_verified`；该报告把团队设置目录刷新回执、六个可替换角色引用、enabled role pack version、slot 映射、团队设置截图和秘书长 / 五顾问映射截图收成安装后可替换角色目录前置门。真实 `truzhenos` 安装、启用版本、团队设置刷新和 TeamRoleSlotBinding 回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：独立验收签收矩阵也要求 P11 复核 `gui_api_traceability_matrix_verified`、商业 API 端点、回执字段和 GUI 证据槽，并在 P11 最终证据包独立验收 section 记录 `gui_api_traceability_matrix_reviewed`；真实复核仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包已新增 `cloud_upload_listing_report`，P11 证据验收检查清单、接入绑定器和验证记录模板已新增 `cloud_upload_listing_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `cloud_upload_listing_verified`；go-live 前必须证明云上传回执、商品草稿引用、artifact manifest、候选包 digest、发布候选包、禁入产物扫描和 hash 匹配回执齐备。真实云上传与商品草稿回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包已新增 `marketplace_listing_review_compliance_report`，P11 证据验收检查清单、接入绑定器和验证记录模板已新增 `marketplace_listing_review_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `marketplace_listing_review_compliance_verified`；该报告把商品页草稿、六角色商品组件清单、市场审核候选、购买前支持 / 退款 / 撤回披露、条款隐私数据边界、发布者身份 / 定价审批和生产发布阻断回执收成正常商品化前置门。真实云端审核与回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包已新增 `purchase_entitlement_report`，P11 证据验收检查清单、接入绑定器和验证记录模板已新增 `purchase_entitlement_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `purchase_entitlement_verified`；该报告把 sandbox 订单、sandbox 支付、无真实扣款、License / Entitlement、团队授权范围、买家已购库 GUI、失败支付阻断和退款 / chargeback 撤权证据收成“可以支付购买”的正常商品化前置门。真实云端订单、支付、授权和已购库回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包已新增 `download_artifact_delivery_report`，P11 证据验收检查清单、接入绑定器和验证记录模板已新增 `download_artifact_delivery_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `download_artifact_delivery_verified`；该报告把已购库下载 GUI、entitlement 校验、签名下载引用、下载回执、下载 artifact hash、候选包 digest 和 hash 匹配回执收成“可以下载安装”的正常商品化前置门。真实云端签名下载、artifact 分发和本地下载回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包的 `download_install_access_matrix_report` 已接入 P11 证据验收检查清单、接入绑定器和验证记录模板的 `download_install_access_result` 写回链；未购买、退款撤权、授权过期、版本下架 / 撤回和 artifact hash 不一致的 blocked 回执必须进入该结果后，才能计入正常商品化 P11 通过。真实云端下载、安装和阻断回执仍待跨仓授权后执行。

团队办公室 Role Pack 候选集补充：P11 最终证据包已新增 `role_studio_lineage_report`，P11 证据验收检查清单、接入绑定器和验证记录模板已新增 `role_studio_lineage_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `role_studio_lineage_verified`；该报告把 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、`role_pack_ref_set_hash`、导出 / 上传 / 下载 / 安装 / 团队绑定回执和运行会话引用收成同一血缘门槛，防止角色制作、商品分发和团队办公室运行使用之间断链。真实 GUI、云端、安装和运行回执仍待跨仓授权后执行。

## 3. 标准文件夹包结构

| 路径 | 职责 | 真相源归属 |
|---|---|---|
| `README.md` | 包说明、主权红线、加载 / 卸载和验收入口 | 本仓文档 |
| `manifest.json` | pack 标识、版本、template_family、治理六件事、ProviderRequirement、护城河、知识域 | Pack 声明；schema 归 contracts / 基座 |
| `flows/*.flow.json` | GateFlowSpec / Scene Flow 图纸，声明节点、边、等待点和门控 | Pack 声明；运行解释归基座 06 |
| `role-slots/role-slots.json` | 场景角色槽、责任、默认角色包、绑定期望 | Pack 声明；绑定执行归基座 13 |
| `role-packs/*.json` | 角色人格、口吻、决策习惯、模型策略和权限边界 | Pack 声明；启用归基座 13 |
| `capabilities/capabilities.json` | 能力需求、ProviderRequirement、gateway class、risk class、fallback policy | Pack 声明；provider 真相归基座 / 外部 provider |
| `knowledge/knowledge-scopes.json` | 知识域、挂载策略、知识 kind | Pack 知识索引；正式知识真相归基座 09 |
| `knowledge/knowledge-index.json` | 知识条目、source_ref、scope、kind、verification status | Pack 知识索引；正式知识真相归基座 09 |
| `knowledge/**/*.md` | 结构化知识内容 | 本仓可分发资料；正式适用需人工核验 |
| `_source-materials/` | Owner 本地投放原始资料入口 | 原始资料不进 Git |
| `install.py` | 调基座 lifecycle 装入 pack、角色、槽位、知识 | 胶水脚本；正式状态归基座 |
| `uninstall.py` | 调基座 lifecycle disable 停用 pack（可逆，只改 KnowledgeMount 可见性/运行访问权，自动留痕） | 胶水脚本；不物理删除 FormalKnowledge、不破坏历史 Receipt——可逆停用属留痕类，不走 Base Gate（Owner 2026-07-03「门禁 vs 留痕」裁定；与 `install.py` 对称，见 AGENTS §5） |
| `docs/` | 测试报告、验收证据、设计说明 | 本仓文档 |

`knowledge/` 是领域知识资产，不是可有可无的附件：可承载术语、口径、规则、例外、SOP、案例和索引（优先使用既有 `glossary / sop / checklist / index` kind），但不替代 FormalKnowledge 真相源（归基座 09）。

## 4. Pack 治理六件事

正式 Domain Work Pack 必须声明：

1. `person_strategy`：谁只是 Proposer，高风险动作如何回 Owner + Base。
2. `formalization_requirement`：哪些候选可进入正式态，裁定权归谁。
3. `gates`：Owner Gate、Base Gate、角色对照门、执行门、发送门等。
4. `provider_requirements`：需要哪些能力和 provider readiness，不实现 provider。
5. `software_requirements`：需要哪些底层软件 family 和版本范围，不封装软件本体、镜像、模型、DB、路径、端口、账号或 runtime state。
5. `notification_command_report_routes`：通知、命令候选、回报和回执路由。
6. `multi_role_comparison`：多角色建议 / 质询 / 对照确认门，禁止隐藏 agent 回路。

此外必须声明 `moat_justification`，回答“用户能否一句话让最强模型直接做出同样结果”。如果能，则应降级为 prompt / 模板 / 能力说明，不应冒充场景荚。

## 5. 跨仓协作边界

| 目标 | 本仓规则 |
|---|---|
| `truzhen-contracts` | Pack manifest / flow / role / capability 引用必须面向契约；schema 变更先改 contracts，并列影响清单。 |
| `truzhenos`（旧 `truzhenv3` 已冻结） | 本仓只放声明数据和装载脚本；Base Gate、Receipt、Gateway、runtime、loader、ReadModel、Pack Studio 归基座。 |
| `truzhen-software` / provider 仓 | Provider、sidecar、外部软件安装、runtime profile、端口、健康检查归外部 provider / software 线。 |
| client repo / frontend | Pack UI 只声明 Surface / visual unit 意图；具体 Web / Desktop / Mobile 渲染归 client repo 或基座前端。 |
| `truzhen-cloud` | 本仓可准备可分发资产和市场展示候选元数据；正式 `PackListing`、审核发布、价格、订单、支付、License / Entitlement、Release、下载分发状态、作者后台和云端分发服务归 `truzhen-cloud`；发布、上架、授权动作必须 Owner 明确授权。 |

### 5.1 云边界红线

- 不保存 `payment_webhook`、`license_server`、`entitlement_db`、`cloud_secret`、`admin_secret`。
- 不把订单、支付状态、License、Entitlement 写成 Pack 真相。
- 不承载官方云端 server、官方云端网页或云端部署脚本。
- 需要云市场发布时，只写 Pack 元数据和发布引用，由 `truzhen-cloud` 持有发布状态和回执。

## 6. 当前包细节

### 6.1 `environmental-enforcement-pack-v0/`

**定位**：生态环境执法全领域证据链 Pack。

**核心资产**：

- `manifest.json`
- `flows/environmental-enforcement-flow.flow.json`
- 角色包：`role_pack://enforcement-elite`、`role_pack://critical-lawyer`
- 能力需求：OCR、文书生成、文书送达、执法执行、在线监测读取
- 知识库：`knowledge-scopes.json`、`knowledge-index.json`、法规 / SOP / 案例 / 索引 Markdown
- `install.py` / `uninstall.py`

**纪律**：

- 法律和执法知识默认 `pending_human_review`。
- 文书、处罚、送达、查封扣押、移送公安等只能是候选，正式动作需 Owner + Base Gate。
- Provider 未接通时必须 `blocked / provider_missing / not_ready`。
- `manifest.json` 的知识域声明必须与 `knowledge/knowledge-scopes.json` 和实际目录一致。

### 6.2 `smart-home-owner-pack-v0/`

**定位**：智能家居老板项目经营 Pack。

**核心资产**：

- `manifest.json`
- `flows/smart-home-owner-project-ops-flow.flow.json`
- 角色包：`role_pack://smart-home-project-manager`
- 能力需求：Frappe 项目快照、客户快照、项目写回候选
- `install.py` / `uninstall.py`

**纪律**：

- Frappe 只是 ProviderRequirement，不是真相源。
- 里程碑、采购、施工、对外承诺、Frappe 写回必须经 Owner + Base Gate。
- 本 pack 当前无知识库；不得声称有 FormalKnowledge 装入。

### 6.3 `housekeeping-ops-pack-v0/`

**定位**：家政 / 保洁客户服务全生命周期 Pack。

**核心资产**：

- `manifest.json`
- `flows/customer-service-lifecycle.flow.json`
- 角色包：`role_pack://housekeeping-consultant`、`role_pack://quality-auditor`
- 能力需求：排期、报价草稿、上门执行意图
- `install.py`
- `uninstall.py`

**纪律**：

- 当前 `knowledge/` 缺失，不能标成完整知识包。
- `uninstall.py` 已补齐，卸载仍只停用 Pack，不删除历史事务对象、候选或 03 回执。
- 报价发送、派工、上门执行和归档正式化必须经 Owner + Base Gate。
- legacy seed 脚本如继续保留，必须标明不替代正式 install lifecycle。

### 6.4 `shuxuejia-renovation-pack-v0/`

**定位**：墅学家大宅装修设计指导 Pack。

**核心资产**：

- `manifest.json`
- `flows/shuxuejia-epc-executable-projection.flow.json`：从旧本地分支 `feature/v3-mod-14-pack-studio-shuxuejia-product-seed-20260624` 的提交 `3c161aff4e2cc9b6ee896d0ef1ec4b37aaf4b062` 只读迁移，保留 457 节点 / 543 边 EPC 拓扑。
- 角色包：业主确认参与者、大宅装修设计负责人、大宅项目交付经理、现场质量监理、材料供应协调员、第三方专项顾问、售后保修服务负责人。
- 能力需求：文档草稿、项目历史知识切片、沟通草稿、受控 artifact 导出、长周期项目 ReadModel、证据与回执回放。
- 知识库：设计指导、材料合同付款、现场质量验收、售后保修，均为 `pending_human_review`。
- `install.py` / `uninstall.py`
- 运行证据：`docs/测试报告.md` 记录隔离安装 / 卸载 / 重装、457 / 457 步完成、665 个候选预览、112 个回执引用和前端用户视角检查；前端缺口仍保留，当前仅声明`已接线`，不声明`已验收 / 已发布`。

**纪律**：

- 旧 seed 只作为来源，不迁移 Go seed 代码，不把基座 product seed 逻辑塞进本仓。
- 大宅装修设计、合同、付款、验收、对外发送和售后关闭均只能先产候选，正式化必须经 Owner + Base Gate。
- Provider 未接通时必须 `provider_missing / not_ready / blocked`。
- 前端用户视角验收必须真实使用 `truzhen-client-web-desktop` 页面；若缺入口，只能登记缺口，不能用 API 冒充通过。

### 6.5 `templates/scene-pack-software-template/`

**定位**：作者端软件工程模板。

**负责**：

- 给 Pack 作者展示可登记、可测试、可说明的软件工程材料结构。
- 提供 objects、flows、views、policies、capabilities、adapters、software、tests、docs 等目录范式。

**不负责**：

- 不作为 enabled pack。
- 不接真实 provider。
- 不保存构建产物、真实密钥或运行数据库。

## 7. 验收映射

| 改动类型 | 必跑验证 |
|---|---|
| Markdown 文档 | `git diff --check` |
| JSON / manifest / flow / role / capability | JSON 合法性 + 结构审计 |
| install / uninstall 脚本 | `python3 -m py_compile` |
| knowledge | JSON 合法性 + scope / index / 文件一致性 + 来源与 `pending_human_review` 检查 |
| 领域语义 / 隐性业务知识 | 概念—对象—关系—规则—来源映射检查；歧义术语不得无来源；knowledge scope / index / 文件一致；高风险内容保持人工核验；只验证声明，不宣称运行时已预加载 |
| forbidden artifacts | `git ls-files` 禁入产物扫描 |
| lifecycle 行为 | 隔离基座 devserver E2E 装入 / 卸载，提供 registry、enabled version、Receipt 证据 |
| ProviderRequirement 语义 | 影响清单 + readiness / blocked 行为说明；必要时跨仓验收 |

## 8. 待统一项

- `environmental-enforcement-pack-v0/` 无固体废物（solid-waste）专属知识域：15 个知识域中无该域目录，固废内容仅散见 penalty / eia-permit / code 等域文件内。2026-07-02 已从 manifest 删除该 scope 声明对齐事实；如 Owner 权威资料含固废专章，后续按知识导入流程真实建域后再恢复声明，不得空声明。
- `housekeeping-ops-pack-v0/` 仍缺 `knowledge/`；后续补齐知识库时不得改变既有 lifecycle / role slot / provider requirement 行为。
- `environmental-enforcement-pack-v0/` 的知识域声明需要持续与实际 `knowledge-scopes.json` 同步。
- `environmental-enforcement-pack-v0/` 当前权威资产口径为 45 份源文档、15 个 scope；09 运行态会按实现切分 FormalKnowledge。历史前台全局 752 条不得写成本 Pack 源数量或固定分片契约，验收须按 Pack/version/source 全量分页覆盖。
- 模板 manifest 与实际 pack manifest 目前不是同一 schema，后续如要作为新包脚手架，应补一个更贴近当前 Domain Work Pack 文件夹结构的模板。

## 9. 完成口径

Pack 商品化或成熟度升级至少要能证明：

- 可安装、可启用、可停用或明确说明卸载缺口。
- Pack Readiness / 护城河理由明确。
- Role Pack 和 SlotBinding 能闭环或明确 `not_ready`。
- Provider readiness 诚实显示 `ready / degraded / provider_missing / blocked`。
- 正式动作不绕过 Owner + Base Gate。
- 关键动作有 Evidence / Receipt 可反查。
- 缺口写入 README、MODULES 或测试报告，不能静默包装成完成。
