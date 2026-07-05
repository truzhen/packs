# 云端商品化台账

当前未接真实 `truzhen-cloud`。本文件只登记商品化预期和阻断条件。

| 阶段 | 状态 | 说明 |
|---|---|---|
| artifact manifest | candidate ready | `commerce/artifact-manifest.json` 已声明上传、下载和安装前置条件，并记录文件级 SHA-256；真实云端 artifact hash 由上传 / 打包环节生成。 |
| sandbox flow candidates | candidate ready | `commerce/sandbox-commerce-flow-candidates.json` 已声明 sandbox 上传、购买、entitlement、下载、安装阶段和阻断负例。 |
| product handoff candidate | candidate ready | `commerce/product-handoff-candidate.json` 已声明云端上传、listing 草稿、sandbox 购买、entitlement、下载和本地安装交接证据；真实 API 和状态仍归 `truzhen-cloud` / `truzhenos`。 |
| license / entitlement policy candidate | candidate ready | `commerce/license-entitlement-policy-candidate.json` 已声明购买后下载、entitlement 校验、安装授权、退款 / 失效阻断和真实支付红线；真实订单和授权真相仍归 `truzhen-cloud`。 |
| order / payment state machine | candidate ready | `commerce/order-payment-state-machine-candidate.json` 已声明 sandbox 订单、支付、失败、退款、chargeback、entitlement 发放和撤销状态机；真实订单和支付状态仍归 `truzhen-cloud`。 |
| support / refund / revocation policy | candidate ready | `commerce/support-refund-revocation-policy-candidate.json` 已声明售后入口、退款、entitlement 撤回、发布撤回通知、历史回执保留和负例阻断；真实支持工单、退款和授权撤回真相仍归 `truzhen-cloud`。 |
| buyer library / install state | candidate ready | `commerce/buyer-library-install-state-candidate.json` 已声明买家已购库、下载、安装、重新安装、团队设置可替换和撤销阻断状态；真实已购库、下载链接、安装和团队设置可替换列表仍归 `truzhen-cloud` / `truzhenos`。 |
| publisher account / settlement policy | candidate ready | `commerce/publisher-account-settlement-policy-candidate.json` 已声明发布者身份、定价审批、结算、税务 / 发票和上线门槛；真实发布者身份、收款资料、结算、发票和支付状态仍归 `truzhen-cloud`。 |
| commercial terms / privacy / data boundary | candidate ready | `commerce/commercial-terms-privacy-policy-candidate.json` 已声明购买前条款接受、隐私告知、数据归属、角色候选免责声明和上线阻断条件；真实条款接受、隐私展示、订单个人数据和购买回执仍归 `truzhen-cloud`，运行态数据仍归 `truzhenos`。 |
| release candidate package | candidate ready | `commerce/release-candidate-package.json` 已声明版本、渠道、签名、升级和撤回策略；真实 Release、签名和生产发布仍归 `truzhen-cloud`。 |
| download / install access matrix | candidate ready | `commerce/download-install-access-matrix.json` 已声明未购买、sandbox 已购、退款撤销、授权过期、版本下架 / 撤回、artifact hash 不一致和缺 Owner Gate 的下载 / 安装访问控制；真实阻断回执归 `truzhen-cloud` / `truzhenos`。 |
| secretary appearance asset rights | candidate ready | `appearance/secretary-appearance-asset-rights-candidate.json` 已声明秘书长音色 / VRM 授权证据、市场审核、安装预检和阻断条件；真实授权与 provider readiness 不归本仓。 |
| artifact bundle layout | candidate ready | `commerce/artifact-bundle-layout-candidate.json` 已声明确定性打包、哈希、签名请求和上传请求；真实 bundle sha256、签名和上传回执仍归 `truzhenos` / `truzhen-cloud`。 |
| artifact bundle digest | candidate ready | `commerce/artifact-bundle-digest-candidate.json` 已声明候选包 payload tree SHA-256 和上传 / 下载 / 安装 hash 连续性引用；真实云端上传、签名下载和安装回执仍待跨仓执行。 |
| install compatibility matrix | candidate ready | `commerce/install-compatibility-matrix.json` 已声明下载后安装、升级、回滚和重新启用的兼容检查；真实安装兼容性验证归 `truzhenos` / `truzhen-cloud` / client GUI。 |
| marketplace review submission | candidate ready | `commerce/marketplace-review-submission-candidate.json` 已声明 listing、价格 / 许可、素材授权、支持 / 退款、审核清单和阻断条件；真实审核状态仍归 `truzhen-cloud`。 |
| commercial go-live approval | candidate ready | `commerce/commercial-go-live-approval-candidate.json` 已声明商品化上线批准总门，要求 GUI、云端、支付、下载、安装、下载 / 安装访问控制矩阵、负例阻断、独立验收和 Owner 裁定证据齐备；当前状态为 `not_approved_requires_cross_repo_evidence`。 |
| commercial production promotion gate | candidate ready | `commerce/commercial-production-promotion-gate-candidate.json` 已声明从 P11 sandbox 证据包、P11 最终证据包、安装后团队办公室运行使用烟测和下载 / 安装访问控制矩阵晋级到生产发布、真实支付和生产签名下载的阻断门；当前状态为 `blocked_requires_p11_verified_evidence_and_owner_decision`。 |
| commercial receipt chain | candidate ready | `commerce/commercial-receipt-chain-candidate.json` 已声明上传、审核、购买、entitlement、下载、安装、绑定和运行使用的回执关联键与断链阻断条件，并要求 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、安装后启用版本和团队 slot 映射从导出贯穿到运行使用；真实回执链待跨仓执行后采集。 |
| commercial distribution receipt schema | candidate ready | `commerce/commercial-distribution-receipt-schema-candidate.json` 已声明上传、审核、sandbox 下单 / 支付、entitlement、签名下载、安装、团队绑定、运行使用、负例阻断和独立验收签收的必填回执字段；真实回执仍归 `truzhen-cloud` / `truzhenos` / client / 独立验收主体。 |
| install preflight request | candidate ready | `install/install-preflight-request-candidate.json` 已声明下载后安装前 entitlement、hash、签名、schema、禁入产物、槽位、asset ref 和 Owner Gate 预检要求；真实预检结果和安装回执归 `truzhenos`。 |
| commercialization execution packet | candidate ready | `tests/commercialization-execution-packet.json` 已声明发布者制作 / 上传与买家购买 / 下载 / 安装 / 团队替换 / 使用的全链路证据包；真实 GUI、云回执和安装回执待跨仓执行。 |
| E2E evidence run record | candidate ready | `tests/e2e-evidence-run-record.json` 已声明真实执行时要填写的阶段证据槽、负例证据槽和完成门；当前全部为 `not_run`，待跨仓授权后采集。 |
| commercial evidence gate | candidate ready | `tests/commercial-evidence-gate-candidate.json` 已声明每阶段最小证据记录、跨仓关联键、独立验收签收和负例阻断要求；无真实 GUI、云端、安装和团队绑定回执时仍不得声称商品化完成。 |
| sandbox environment readiness | candidate ready | `tests/sandbox-environment-readiness-candidate.json` 已声明云端 sandbox、买卖双方、支付桩、签名下载、本地安装目标和团队设置入口的 ref-only 前置条件；真实环境就绪仍需 client、`truzhen-cloud`、`truzhenos` 授权后验证。 |
| commercial observability diagnostics | candidate ready | `tests/commercial-observability-diagnostics-candidate.json` 已声明 GUI、云端、安装和团队绑定的 trace、日志、指标、告警、脱敏和回执关联要求；真实诊断接线仍归 client、`truzhen-cloud`、`truzhenos`。 |
| normal commercialization completion audit | candidate ready | `tests/normal-commercialization-completion-audit-candidate.json` 已按 Owner 目标逐条映射当前候选证据与仍缺的 GUI 截图、云端回执、安装回执、团队绑定回执和独立验收签收；当前明确 `not_achieved_requires_cross_repo_execution`。 |
| commercial chain verifier | candidate ready | `tests/commercial-chain-verifier-candidate.json` 已声明正常商品化链路的逐阶段核验门：GUI 证据、云上传、审核、sandbox 订单 / 支付、entitlement、下载 hash、安装、团队绑定、运行使用、负例阻断和独立验收；当前未运行，不能替代真实回执。 |
| commercial go/no-go gate | candidate ready | `tests/commercial-go-no-go-gate-candidate.json` 已把执行队列、跨仓证据台账、P11 证据包、问题台账 triage、商品阶段前后端收口报告、安装后团队办公室运行使用烟测、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵、链路核验器和上线批准收束成可机读 go/no-go 门；当前 `decision_status=blocked_not_ready_for_commercial_go_live`。 |
| commercial readiness verifier | candidate ready | `tests/commercial-readiness-verifier-candidate.json` 已汇总当前商品化阻断、证据回写计数、问题台账写回摘要、商品阶段前后端收口报告、安装后团队办公室运行使用烟测、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵和终端门状态；当前 `verification_status=blocked_not_commercial_ready`，不能替代云端上传、支付、下载、安装或 Owner 裁定回执。 |
| role studio phase coverage matrix | candidate ready | `tests/role-studio-phase-coverage-matrix-candidate.json` 已把 P0-P11 计划阶段逐项映射到用户视角 GUI 证据、后端回执、当前候选资产、缺失权威证据、目标仓和阻断状态；当前所有跨仓阶段仍待授权执行。 |
| role studio test case coverage matrix | candidate ready | `tests/role-studio-test-case-coverage-matrix-candidate.json` 已把计划第 8 节 24 个 `TC-*` 测试任务映射到 GUI 证据、权威回执 / blocked 证据、目标仓和当前阻断状态；当前所有行仍待跨仓授权执行。 |
| independent acceptance signoff matrix | candidate ready | `tests/independent-acceptance-signoff-matrix-candidate.json` 已声明独立验收智能体逐阶段签收规则，要求 GUI 证据、后端回执、云端购买 / 下载回执、本地安装 / 团队绑定回执、负例阻断和 Owner go / no-go 裁定齐备；组织者自述不能作为完成证据。 |
| P0-P11 commercialization blocker register | candidate ready | `tests/p0-p11-commercialization-blocker-register-candidate.json` 已把上传、购买、下载、安装和团队绑定等缺口登记为 issue、目标仓、复现步骤、权威证据和独立验收要求；未获 Owner 授权和真实回执前仍不得声明商品化完成。 |
| product-stage frontend backend closure report | candidate ready | `tests/product-stage-frontend-backend-closure-report-candidate.json` 已把商品阶段前端 GUI 能力、后端 candidate / Gate / Receipt、字段一致性、P0/P1 剩余缺口和下一轮执行卡集中登记；真实前端 smoke、后端 receipt lookup、云端回执和安装回执仍待跨仓授权后验证。 |
| P11 normal commercialization acceptance gate | candidate ready | `tests/p11-normal-commercialization-acceptance-gate-candidate.json` 已把上传、审核、sandbox 购买 / 支付、entitlement、下载 hash、安装、团队绑定、负例和独立验收拆成 go / no-go 判定门；真实通过仍需 truzhen-cloud、truzhenos 和用户视角 GUI 回执。 |
| P11 normal commercialization verification record template | candidate ready | `tests/p11-normal-commercialization-verification-record-template.json` 已定义真实执行后如何填写阶段结果、hash 连续性、负例结果、独立验收和最终决策；当前仍为 `not_run` 模板，不是云上传、购买、下载或安装成功记录。 |
| P11 evidence ingestion binder | candidate ready | `tests/p11-evidence-ingestion-binder-candidate.json` 已定义真实 GUI 步骤、云端回执、本地安装回执、团队绑定、负例阻断、独立验收和 Owner 裁定如何接入 P11 验证记录；当前仍为 `not_run` 规则，不是真实执行记录。 |
| P11 sandbox execution runbook | candidate ready | `tests/p11-sandbox-execution-runbook-candidate.json` 已定义 Owner 授权后跨仓 GUI / 云 / 安装执行顺序、证据输出、hash 连续性、回滚和禁区；当前仍为 `pending_owner_authorization`，不是授权或执行记录。 |
| P11 sandbox preflight gate | candidate ready | `tests/p11-sandbox-preflight-gate-candidate.json` 已定义 P11 开跑前 Owner 授权、sandbox 环境、支付桩、签名下载、本地安装目标和团队设置入口的阻断条件；当前仍为 `not_ready` 预检门，不是授权或执行记录。 |
| P11 sandbox run request | candidate ready | `tests/p11-sandbox-run-request-candidate.json` 已定义授权后 P11 sandbox 十阶段执行请求，要求发布者 / 买家用户视角 GUI 智能体产出上传、购买、下载、安装、团队替换、负例、独立验收和生产晋级控制证据；当前仍为 `not_ready` 请求，不是真实执行记录。 |
| P11 evidence acceptance checklist | candidate ready | `tests/p11-evidence-acceptance-checklist-candidate.json` 已定义 P11 sandbox 实操后逐阶段验收条件，覆盖上传、购买、下载、安装、团队设置替换、负例、独立验收和 Owner go/no-go；当前仍为 `not_run` 检查清单，不是真实通过记录。 |
| P11 commercial go-live evidence package template | candidate ready | `tests/p11-commercial-go-live-evidence-package-template.json` 已定义最终证据包的 GUI 证据索引、云 / 安装 / 团队绑定回执索引、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告、hash 连续性、下载 / 安装访问矩阵报告、负例、独立验收、Owner go / no-go 和最终上线裁定；GUI 证据索引必须覆盖用户视角脚本全部 14 个步骤。当前仍为 `not_run` 模板，不是商品化通过记录。 |
| role studio goal completion evidence map | candidate ready | `tests/role-studio-goal-completion-evidence-map-candidate.json` 已把活跃目标拆成角色制作、使用、上传、购买、下载、安装、团队设置替换和正常商品化 go/no-go 的权威证据要求，并强制消费商品化上线批准、P11 证据验收清单、问题台账全量 triage、商品阶段前后端收口报告和安装后团队办公室运行使用烟测；当前全部为 `missing_authoritative_evidence`。 |
| GUI evidence capture protocol | candidate ready | `tests/gui-evidence-capture-protocol.json` 已声明用户视角 GUI 证据字段、截图 / 页面状态 / 回执配对、脱敏和链路关联；真实证据待跨仓 GUI 执行采集。 |
| GUI user agent execution script | candidate ready | `tests/gui-user-agent-execution-script-candidate.json` 已声明用户视角智能体按 GUI 执行的有序步骤、动作类型、页面状态、证据槽和负例脚本；当前未运行。 |
| cross repo execution cards | candidate ready | `integration/cross-repo-execution-cards.json` 已声明 cloud / client / `truzhenos` / contracts 授权后施工卡、验证命令和证据输出；当前为 `pending_owner_authorization`。 |
| Owner authorization evidence intake | candidate ready | `integration/owner-authorization-evidence-intake-candidate.json` 已声明跨仓执行前必须记录的 Owner 授权原文、逐仓范围、允许动作、红色禁区、证据输出和过期条件；当前为 `missing_owner_authorization`，不能启动跨仓执行。 |
| commercial cross repo execution queue | candidate ready | `integration/commercial-cross-repo-execution-queue-candidate.json` 已按授权 intake、contracts、后端回执、用户视角 GUI、云端 sandbox、安装绑定、负例观测和独立验收拆成阶段队列；当前为 `blocked_pending_owner_authorization_and_stage_evidence`。 |
| commercial cross repo evidence ledger | candidate ready | `docs/commercial-cross-repo-evidence-ledger.json` 已按 8 个商品化跨仓阶段固定 `evidence_id`、目标仓、待填证据位置、写回目标和未授权阻断；当前所有行均为 `pending_authorization`，不能替代真实 GUI、云端、安装或独立验收回执。 |
| role studio issue ledger | candidate ready | `docs/role-studio-issue-ledger.md` 已登记 GUI、云端、go/no-go、P11 证据包等角色制作台缺口；readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `issue_ledger_all_entries_triaged`，真实关闭仍待跨仓授权和权威回执。 |
| cross repo execution readiness package | candidate ready | `integration/cross-repo-execution-readiness-package.json` 已汇总目标仓路径、状态命令、允许动作、证据输出、阶段顺序和 Owner 授权问题；当前仍不是授权或执行记录。 |
| commercial API contract | candidate ready | `integration/commercial-api-contract-candidate.json` 已声明上传、审核、购买、支付、entitlement、下载、安装、团队设置刷新、生产发布、真实支付启用和生产签名下载启用请求的候选 request / response / 回执 / 幂等键 / 错误态；生产晋级端点必须暴露 P11 最终证据包缺失和安装后运行使用烟测缺失的独立 blocked 状态；真实 API 实现仍归 client、`truzhen-cloud` 和 `truzhenos`，生产动作仍需生产晋级门和 Owner go/no-go。 |
| commercial API example cases | candidate ready | `integration/commercial-api-example-cases-candidate.json` 已为生产 go-live、真实支付启用、生产签名下载启用和正式上架发布请求提供 P11 最终证据包缺失与安装后运行使用烟测缺失的 failure response 样例；真实前后端错误展示和回执仍待跨仓执行。 |
| GUI / API traceability matrix | candidate ready | `integration/frontend-backend-contract-map.json` 已把 P11 导出、云上传、审核、sandbox 购买、下载、安装、团队替换、运行使用和负例 GUI 步骤映射到商业 API 端点、回执字段、阻断态和证据槽；真实 GUI / API 对账仍待跨仓执行。 |
| upload draft | pending | 需 `truzhen-cloud` 授权后由 GUI 上传候选包。 |
| listing review | pending | 正式审核 / 上架必须 Owner 单独授权。 |
| sandbox purchase | pending | 可在 sandbox 中测试订单、支付、License / Entitlement。 |
| download | pending | 已购下载需 entitlement 校验和 artifact hash 校验。 |
| install | pending | 下载后安装归 `truzhenos`，需安装回执和 enabled version。 |
| real payment | blocked | 红色动作，未授权前不得执行。 |
| production publish | blocked | 红色动作，未授权前不得执行。 |
