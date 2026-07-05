# Truzhen Packs

> Truzhen 主权事务操作层的**领域工作包仓**：可独立加载 / 卸载 / 分发的场景荚（Domain Work Pack），面向契约、不含基座实现。

`github.com/truzhen/packs` 是 Truzhen 六仓协同架构的**包层**。每个包把一个行业的高手经验编译成受主权闸门约束的领域治理资产——**声明**判人 / 判事策略、门控流程、真实系统 Provider 绑定、证据与回执要求，由作者在本行业自分销或提交官方云市场。正式裁定权在 Owner + Base Gate，**包不持有主权**，AI 永远是 Proposer。

## 依赖方向（单向不可逆）

```
┌─────────────────────┐   implements   ┌──────────────────────┐   faces   ┌─────────────────────┐
│ truzhenos (基座·私有) │ ─────────────▶ │ truzhen-contracts     │ ◀──────── │ truzhen-packs (本仓) │
│ 实现契约             │                │ 纯接口/类型/Schema     │           │ 面向契约            │
└─────────────────────┘                └──────────────────────┘           └─────────────────────┘
```

包**面向契约**编写，物理上 import 不到基座内部；基座通过文件夹包加载器 / 各包 `install.py` 经真实 lifecycle 端点装入。`truzhenos`、`truzhen-contracts`、`truzhen-packs`、`truzhen-software`、`truzhen-cloud` 与 client repo 默认**平级**放在 `/Users/li/Documents/`。

官方云市场的 `PackListing`、审核发布、价格、支付、License / Entitlement、Release 与下载分发服务归 `truzhen-cloud`。本仓 manifest 可以提供展示候选字段，但不得把商品状态、授权状态或支付状态当成本仓事实。

官方云服务、云市场商品运营状态、订单、支付回调、License / Entitlement 真相、Pack 文件分发运行面和官方云端网页归 `/Users/li/Documents/truzhen-cloud`。本仓只声明 Pack 本体、工作台、能力需求和可被商品化引用的元数据，不保存订单、License、支付状态或云端 server 实现。

## 三类 Pack（不得发明第四种）

- **场景荚（Domain Work Pack / Scene Pack）**：可交付的领域治理资产，声明六件事（判事 / 判人 / 门控流程 / Provider 绑定 / 通知-命令-回报路由 / 多角色对照），**不持 Base 主权**。本仓主体。
- **能力荚（Capability Pack）**：能力描述符 / ProviderRequirement。执行 provider 本体归 `truzhenos` 或外部 provider / `truzhen-software`，不写进本仓。
- **角色荚（Role Pack）**：智能体人格 / 口吻 / 决策习惯 / 模型策略，绑定到任意 Role Slot。本仓可随场景荚携带角色包数据。
- **商品化引用**：Pack 可以声明市场展示元数据和商品化引用，但订单、支付、License、Entitlement 和分发运行状态归 `truzhen-cloud`，不归本仓。

## 当前包状态（2026-07-04）

- `environmental-enforcement-pack-v0/`：完整文件夹包，含 install/uninstall、15 知识域、角色包。
- `smart-home-owner-pack-v0/`：完整文件夹包，含 install/uninstall、项目经理角色包、Frappe ProviderRequirement。
- `housekeeping-ops-pack-v0/`：已升级为可装入 / 可卸载文件夹包，含 `pack_ref` / `template_family` / flow / capabilities / role-slots / 2 角色包 / `install.py` / `uninstall.py`；`knowledge/` 仍待补，不得标成完整知识包。
- `capability-pack-candidates/short-video-ops-v0/`：短视频运营能力包候选集，含 3 个 Capability Pack 候选声明、Code Assistant 调用候选台账、运行请求候选台账、candidate bundle 导出台账、candidate bundle dry-run 台账、PatchCandidate 承接台账、PatchCandidate 复核台账、P11 lifecycle preflight 执行规格和商用缺口台账、P12 开跑前门禁候选、P12 / P13 / P15 / P16 / P17 / P18 运行后证据验收门候选、P12-P18 执行规格 / 授权卡 / 机器证据契约 / 覆盖 `evidence_id` 的证据台账骨架 / 执行就绪证据写回计划、带写回摘要的商用跨仓执行队列、带写回总计并消费 P12 / P13 / P15 / P16 / P17 / P18 后验收门的商用 readiness verifier、消费写回门禁和 P12 / P13 / P15 / P16 / P17 / P18 后验收门的商用 go/no-go 候选门禁、商用后验收门覆盖 verifier、商用禁入动作覆盖 verifier、商用执行就绪守卫覆盖 verifier、消费商用后验收门 / 禁入动作 / 执行就绪守卫覆盖 verifier 的目标完成证据地图候选、P12-P18 商用证据契约总索引和逐分片签字矩阵、前后端证据交接 runbook / 命令计划、GUI/API/Receipt 追踪矩阵、P14 商用 readiness 审计、商用改进清单、商用最终证据包模板和授权路线图；只用于能力包制作台压测和候选资产沉淀，不参与 enabled pack 分发。
  - 补充：短视频候选集另含 `capability-pack-candidates/short-video-ops-v0/docs/commercial-go-live-evidence-package.json`，作为商用最终证据包的机器摘要候选；P12-P18 未授权和证据未写回前仍为 `not_commercial_ready`。
  - 补充：短视频候选集另含 `capability-pack-candidates/short-video-ops-v0/docs/commercial-machine-go-live-evidence-package-coverage-verifier.json`，用于证明最终机器证据包已被商用签收链路消费；当前仍 blocked。
  - 补充：目标完成证据地图、机器证据包、readiness verifier 和 go/no-go 门禁已明确拒绝用覆盖 verifier、机器证据包摘要、go/no-go 候选或前后端验收契约替代真实 runtime receipts、P12-P18 分片回执、前后端回执和 Owner/Base Gate 回执。
  - 补充：前后端商用验收契约已新增 GUI/API/Receipt 追踪矩阵，并由商用 readiness、go/no-go、最终机器证据包和目标完成图消费；截图-only、后端测试-only 或网络摘要-only 都不能计入前后端验收或商用完成。
  - 补充：该候选集另含商用授权尝试覆盖 verifier，用于证明被拒绝的 P12 授权尝试已被商用门禁文档消费且不能当作 Owner 授权。
  - 补充：该候选集另含商用下一授权启动守卫，用于集中校验 P12 授权卡原话、允许仓、禁改仓和禁入动作，未通过前不得开始跨仓执行，也不得通过 P12-P18 分片级 execution readiness gate、商用证据索引完成门、跨仓队列首跑门、readiness、go/no-go 或 Owner 签收。
  - 补充：该候选集另含 Pack Studio issue 回填台账候选，用于把 P10-P18 制作台缺口、阻断原因、归属、验收证据、权威关闭证据槽和禁止绕过方式集中登记，并被 P12 开跑前门禁、P12-P18 execution readiness、商用跨仓执行队列首跑门和各分片队列镜像 gate、商用执行就绪守卫覆盖 verifier、商用当前状态审计、前后端验收契约、前后端交接、backlog、readiness、go/no-go、最终机器证据包和目标完成图消费；每个 issue 关闭前必须有实现 / 授权证据、验证命令、GUI/API/ReadModel 证据、Receipt / Candidate 引用和禁入方案缺席证明，当前仍 `pack_studio_issue_ledger_open`。
  - 补充：该候选集另含商用独立验收签收矩阵，要求 P12-P18、GUI/API/Receipt、issue 关闭证据、禁入动作、证据写回、Owner/Base Gate 和最终机器证据包逐项绑定独立 reviewer、artifact refs 和权威证据；readiness、go/no-go、最终机器证据包和目标完成图均要求 `independent_acceptance_signoff_matrix_passed`，当前仍 blocked。
  - 补充：P18 云市场 sandbox 后验收门已明确拒绝 packs listing draft、sandbox runbook、License / Entitlement 文案、人工订单状态或生产发布声称替代 `truzhen-cloud` receipt；P18 完成和商用签字必须绑定 cloud sandbox receipt、entitlement receipt 和 packs 无云市场真相扫描。
- `role-pack-candidates/team-office-v0/`：团队办公室 Role Pack 候选集，含秘书长 + 五顾问 6 个角色包候选、团队角色槽、团队设置绑定候选、安装后团队设置角色目录候选、六角色能力引用样例、秘书长音色 / VRM asset ref、秘书长音色 / VRM 授权证据候选、云商品草稿候选（含买家可见六角色商品组件清单）、artifact manifest、sandbox 商品化流程候选、商品化交接候选、License / Entitlement 策略候选、订单 / 支付状态机候选、售后 / 退款 / 撤回策略候选、买家已购库 / 安装状态候选、发布者账号 / 结算 / 发票策略候选、商业条款 / 隐私 / 数据边界候选、发布候选包、下载 / 安装访问控制矩阵候选、确定性打包规格候选、候选包整体 digest、安装兼容性矩阵、云市场审核提交候选、商品化上线批准候选、生产发布晋级门候选（含下载 / 安装访问矩阵硬门）、商品化回执链候选、商业分发回执记录规格候选、安装预检请求候选、运行使用候选、用户视角 GUI 场景候选、用户视角 GUI 执行脚本候选、GUI 证据采集协议候选、上传前 secret / raw asset 扫描候选、商品化完成证据矩阵、目标完成证据地图候选、商品化执行包、E2E 证据运行记录包、商品化证据门槛候选、sandbox 商品化环境就绪候选、商品化监控 / 诊断候选、正常商品化完成审计候选、正常商品化链路核验器候选、商品化 go/no-go 门禁候选、商品化 readiness verifier 候选、P0-P11 阶段覆盖矩阵候选、计划 TC 测试用例覆盖矩阵候选、独立验收签收矩阵候选、P0-P11 商品化阻塞清单候选、商品阶段前后端收口报告候选、P11 正常商品化通过门候选、P11 正常商品化验证记录模板、P11 证据接入绑定器候选、P11 sandbox 执行 runbook 候选、P11 sandbox 开跑前预检门候选、P11 sandbox 执行请求候选、P11 证据验收检查清单候选、P11 商品化最终证据包模板（含云上传商品草稿报告、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告）、商品化前后端 API 契约候选、跨仓执行卡（含生产晋级控制施工卡）、跨仓执行就绪包（含生产晋级阶段）、Owner 授权证据接入口候选、商品化跨仓执行队列候选（含生产晋级控制阶段）、商品化跨仓证据台账候选（含生产晋级回执写回）、前后端接线映射（含 P11 GUI/API traceability）和台账；只用于角色制作台与商品化链路压测，不参与 enabled pack 分发。
  - 补充：该候选集另含商品化 API 请求 / 响应样例候选，用于前端、后端、云端和验收智能体按同一端点字段核对。
  - 补充：该候选集另含 artifact manifest 闭合门候选，用于上传、签名下载和本地安装前核对当前目录、候选集和 artifact manifest 文件集合一致。
  - 补充：该候选集另含安装到运行使用激活映射候选，用于把下载回执、安装预检、启用版本、团队设置目录刷新、Owner Gate 绑定和运行使用烟测串成可验链路。
  - 补充：该候选集的能力 Pack 角色引用样例覆盖秘书长和五顾问全部 6 个制作角色，并要求安装后引用证据门；能力 Pack 只有在 entitlement 校验、安装回执和 enabled role version 齐全后，才能引用已购角色。
  - 补充：云商品草稿含买家可见的六角色商品组件清单，发布包和云市场审核提交均引用同一清单；秘书长音色 / VRM 只允许 asset ref，不随商品包夹带 raw 音频或 raw VRM。
  - 补充：License / Entitlement 策略和安装预检要求 team-scoped entitlement 的授权团队必须匹配安装目标团队，跨团队安装必须 blocked。
  - 补充：购买前必须展示支持入口、退款条款、撤回说明和候选角色免责声明；缺少披露或买家确认时，sandbox 购买必须 blocked。
  - 补充：秘书长音色 / VRM 表现层候选新增团队设置页 GUI 选择、清空、恢复默认和 provider readiness 显示规则；完成声明前必须有 `secretary_appearance_gui_controls_verified` 证据门。
  - 补充：P11 商品化最终证据包模板要求六角色安装目录、启用版本、slot 映射、团队设置截图、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告和下载 / 安装访问矩阵报告齐全后，才允许进入 go-live 判断。
  - 补充：P11 商品化最终证据包模板的 `install_catalog_and_slot_mapping_report` 已被 readiness、go/no-go、产品完成矩阵和目标完成地图直接消费为 `six_installed_role_catalog_evidence_verified`；缺少团队设置目录刷新回执、六个可替换角色引用、enabled role pack version、slot 映射和秘书长 / 五顾问映射截图时，不得声称安装后可替换角色已就绪。
  - 补充：P11 商品化最终证据包模板已新增 `cloud_upload_listing_report`，要求云上传回执、商品草稿引用、artifact manifest、候选包 digest、发布候选包、禁入产物扫描和 hash 匹配回执齐全；P11 证据验收检查清单、证据接入绑定器和验证记录模板已新增 `cloud_upload_listing_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `cloud_upload_listing_verified`，缺少该证明时不得进入 go-live 判断。
  - 补充：P11 商品化最终证据包模板已新增 `purchase_entitlement_report`，要求 sandbox 订单、sandbox 支付、无真实扣款、License / Entitlement、团队授权范围、买家已购库 GUI、失败支付阻断和退款 / chargeback 撤权证据齐全；P11 证据验收检查清单、证据接入绑定器和验证记录模板已新增 `purchase_entitlement_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `purchase_entitlement_verified`。
  - 补充：P11 商品化最终证据包模板已新增 `download_artifact_delivery_report`，要求已购库下载 GUI、entitlement 校验、签名下载引用、下载回执、下载 artifact hash、候选包 digest 和 hash 匹配回执齐全；P11 证据验收检查清单、证据接入绑定器和验证记录模板已新增 `download_artifact_delivery_result` 写回链，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `download_artifact_delivery_verified`，缺少该证明时不得声称“可以下载安装”。
  - 补充：P11 商品化最终证据包模板的 `download_install_access_matrix_report` 已接入 P11 证据验收检查清单、证据接入绑定器和验证记录模板的 `download_install_access_result` 写回链；未购买、退款撤权、授权过期、版本下架 / 撤回和 artifact hash 不一致的 blocked 回执必须写回后才可计入 P11 通过。
  - 补充：P11 商品化最终证据包模板已新增 `role_studio_lineage_report`，要求同一 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、`role_pack_ref_set_hash`、导出 / 上传 / 下载 / 安装 / 团队绑定回执和运行会话引用贯穿，readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `role_studio_lineage_verified`。
  - 补充：P11 证据验收检查清单、证据接入绑定器和验证记录模板已新增 `role_studio_lineage_result` 写回链，独立验收智能体必须复核并写回血缘结果后才可计入 P11 通过。
  - 补充：商品化 readiness verifier 和 go/no-go 门禁已直接消费秘书长音色 / VRM 证据报告，完成声明前必须证明 GUI 选择、asset ref、授权回执和 raw asset 缺席扫描齐备。
  - 补充：商品化 readiness verifier、go/no-go 门禁、产品完成矩阵和目标完成地图已直接消费 `docs/role-studio-issue-ledger.md`，完成声明前必须证明 `issue_ledger_all_entries_triaged`。
  - 补充：商品化 readiness verifier、go/no-go 门禁和目标完成地图已直接消费商品阶段前后端收口报告，完成声明前必须证明 `product_stage_frontend_backend_closure_report_verified`。
  - 补充：商品化 readiness verifier、go/no-go 门禁、产品完成矩阵和目标完成地图已直接消费安装后团队办公室运行使用烟测，完成声明前必须证明 `team_office_runtime_usage_smoke_verified`。
  - 补充：P11 sandbox 执行请求要求导出、云上传、审核、sandbox 购买、下载、安装、团队绑定和负例阶段逐项绑定 GUI 证据协议、权威回执和同一 `bundle_tree_sha256` 关联键；完成声明前必须有 `p11_gui_receipt_continuity_verified`。
  - 补充：P11 sandbox 执行请求另声明阶段依赖和前置回执传递契约，后一阶段必须引用前一阶段回执，不能跳过上传、审核、entitlement、下载或安装直接声称商品化完成。
  - 补充：商品化跨仓证据台账已镜像 P11 阶段依赖，真实写回时每个 phase link 必须记录依赖 phase、前置回执和缺失阻断状态。
  - 补充：商品化跨仓执行队列已镜像 P11 阶段依赖，授权后执行时缺少前置回执的 phase 必须保持 `blocked_previous_phase_evidence_missing`。
  - 补充：正常商品化链路核验器已把 P11 执行队列阶段依赖证明列为通过条件，不能只凭最终支付、下载或安装回执跳过中间阶段。
  - 补充：商品化 go/no-go 门禁已把 P11 执行队列阶段依赖证明列为终局门槛，Owner go/no-go 前必须证明上传、审核、支付授权、下载、安装、绑定连续。
  - 补充：P11 商品化最终证据包模板已新增阶段依赖报告 section，用于把执行队列 phase handoff、前置回执和缺失阻断状态交给 Owner 与独立验收主体复核。
  - 补充：P11 证据接入绑定器和验证记录模板已新增阶段依赖接入 / 填写字段，真实执行证据必须落到 `p11_phase_dependency_result` 后才能进入最终证据包。
  - 补充：P11 证据验收检查清单已新增阶段依赖验收门，完成声明前必须有执行队列 phase link、前置回执、`bundle_tree_sha256` 和 `correlation_id` 的连续证明。
  - 补充：商品化上线批准候选已直接消费 P11 证据验收清单、最终证据包和阶段依赖验收门，缺少清单 verified 证据时不能进入 go-live approval。
  - 补充：商品化上线批准候选已直接消费下载 / 安装访问控制矩阵，未购买、退款撤销、版本下架 / 撤回或 artifact hash 不一致缺少 blocked 回执时不能进入 go-live approval。
  - 补充：生产发布晋级门已直接消费 P11 最终证据包、下载 / 安装访问控制矩阵和安装后团队办公室运行使用烟测；缺少最终证据包、矩阵 verified 或 `team_office_runtime_usage_smoke_verified` 时，不得启用真实支付、生产签名下载、正式上架或安装分发。
  - 补充：目标完成证据地图已直接消费商品化上线批准和 P11 证据验收清单，完成声明前必须同时有 `commercial_go_live_approval_verified` 与 `p11_evidence_acceptance_checklist_verified`。
  - 补充：商品化 readiness verifier、go/no-go 门禁和产品完成矩阵已直接消费 `commercial_go_live_approval_verified`；缺少上线批准候选、P11 证据验收清单、最终证据包、Owner go/no-go 和生产发布 / 真实支付阻断证据时，不得声称正常商品化可通过。
  - 补充：P11 正常商品化通过门、链路核验器和 go/no-go 门禁要求版本下架 / release revoked 时下载与安装均有 blocked 证据。
  - 补充：商品化 readiness verifier 和 go/no-go 门禁已直接消费下载 / 安装访问控制矩阵，完成声明前必须证明未购买、退款撤销、授权过期、版本下架 / 撤回和 artifact hash 不一致均有权威 blocked 回执。
  - 补充：该候选集另含角色制作台导出来源证明候选，用于阻断手工 JSON、后端直调或无 GUI 导出回执的云上传。
  - 补充：商品化回执链已要求 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、安装后启用版本和团队 slot 映射从制作台导出贯穿到团队办公室运行使用，防止下载安装后的运行角色集被替换成非制作台产物。
  - 补充：P11 商品化最终证据包模板已要求用户视角 GUI 执行脚本全部 14 个步骤都有截图、页面状态、操作日志和候选 / 回执引用，缺一步都不得进入 go-live 判断。
  - 补充：商品化 API 契约和请求 / 响应样例已暴露生产晋级失败态，真实支付启用、生产签名下载和正式上架缺少 P11 最终证据包或安装后运行使用烟测时必须返回可读 blocked 回执。
  - 补充：前后端接线映射已新增 `gui_api_traceability_matrix`，把 P11 用户视角 GUI 脚本中的导出、云上传、审核、购买、下载、安装、团队替换、运行使用和负例逐项映射到商业 API 端点、回执字段、阻断态和 GUI 证据槽；P10 收口报告、P11 最终证据包、P11 验收清单、readiness verifier、go/no-go、商品完成矩阵和目标完成地图均要求 `gui_api_traceability_matrix_verified` 后才允许完成声明。
  - 补充：商品化跨仓执行队列和证据台账已要求 `frontend_user_view_gui_flow` 阶段交付 `gui_api_traceability_matrix_verified`，并写回 P11 最终证据包、P11 验收清单和 P10 收口报告；真实跨仓 GUI 执行时不能只交截图和网络摘要。
  - 补充：独立验收签收矩阵已把 `gui_api_traceability_matrix_verified` 纳入 P11 签收前置证据，最终证据包的独立验收 section 必须记录 `gui_api_traceability_matrix_reviewed`；没有逐步对上 GUI 操作、商业 API 端点、回执字段和 GUI 证据槽，不得通过商品化完成声明。
  - 补充：P11 商品化最终证据包已新增 `marketplace_listing_review_compliance_report`，要求商品页草稿、六角色商品组件清单、市场审核候选、购买前支持 / 退款 / 撤回披露、条款隐私数据边界、发布者身份 / 定价审批和生产发布阻断回执齐备后，才允许进入正常商品化 go/no-go。
  - 补充：P11 证据验收检查清单、证据接入绑定器和验证记录模板已新增 `marketplace_listing_review_result` 写回链；商品页草稿、市场审核、购买前披露、条款隐私、发布者身份 / 定价和生产发布阻断证据必须写回后，才可计入 P11 通过。

## 文件夹包标准结构

```
<pack>-v0/
  manifest.json          # 场景荚规格：pack_ref / version / template_family + 六件事 + provider_requirements + moat_justification
  flows/*.flow.json      # GateFlowSpec 门控流程图
  role-slots/            # Role Slot 声明
  role-packs/            # 绑定的角色包
  capabilities/          # 能力需求 / ProviderRequirement 引用，不放 provider 实现
  knowledge/             # 领域知识（knowledge-scopes.json + knowledge-index.json + 各 .md，权威资料结构化）
  _source-materials/     # Owner 投放的原始资料区（不进 Git，只留 .gitignore + README 占位）
  install.py / uninstall.py  # 经基座真实 lifecycle 端点装入 / 卸载
```

## 加载 / 卸载（经基座真实主权链）

```sh
# 起基座 truzhenos devserver（隔离 DB），与本仓平级放在同一父目录
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 environmental-enforcement-pack-v0/install.py
# 铁证：registry 0→1；知识入 09 FormalKnowledge；正式动作经 Owner + Base Gate 裁定并产 Receipt
```

## 护城河测试

每个场景荚必过：「用户能否一句话让最强模型直接做出同样结果？」能 → 不合格降级为 prompt / 模板；不能（因需真实系统门控 / 领域审批合规 / 可审计证据回执 / 显式多角色对照 / 事务对象生命周期回放）→ 合格。写入 manifest 的 `moat_justification`。

## 子包清单

见 [MODULES.md](MODULES.md)。

## License

[Apache-2.0](LICENSE)。
