# 角色包导出台账

当前导出方式：本仓候选资产文件。正式导出应由角色制作台 GUI 触发，并由 `truzhenos` 生成 candidate bundle、hash 和回执。

| artifact | status | truth source |
|---|---|---|
| `candidate-set.json` | candidate asset ready | `truzhen-packs` |
| `role-packs/*.rolepack.json` | candidate asset ready | `truzhen-packs` |
| `role-slots/team-office-role-slots.json` | candidate asset ready | `truzhen-packs` |
| `bindings/team-office-role-binding-candidate.json` | candidate asset ready，含团队设置替换、回滚和 Owner Gate 绑定要求 | `truzhen-packs` |
| `bindings/team-settings-installed-role-catalog-candidate.json` | candidate asset ready，含安装后团队设置 tab 可替换角色目录、刷新策略和阻断条件 | `truzhen-packs` |
| `appearance/secretary-appearance-asset-rights-candidate.json` | candidate asset ready，含秘书长音色 / VRM 授权证据、市场审核、安装预检和阻断条件 | `truzhen-packs` |
| `commerce/cloud-listing-candidate.json` | candidate asset ready | `truzhen-packs` |
| `commerce/artifact-manifest.json` | candidate asset ready，含文件级 SHA-256 | `truzhen-packs` |
| `commerce/sandbox-commerce-flow-candidates.json` | candidate asset ready，含 sandbox 商品化阶段和负例 | `truzhen-packs` |
| `commerce/product-handoff-candidate.json` | candidate asset ready，含云端商品化和本地安装交接证据要求 | `truzhen-packs` |
| `commerce/license-entitlement-policy-candidate.json` | candidate asset ready，含 License / Entitlement 策略和授权阻断要求 | `truzhen-packs` |
| `commerce/order-payment-state-machine-candidate.json` | candidate asset ready，含 sandbox 订单、支付、失败、退款、chargeback 和 entitlement 发放 / 撤销状态机 | `truzhen-packs` |
| `commerce/support-refund-revocation-policy-candidate.json` | candidate asset ready，含售后入口、退款、entitlement 撤回、发布撤回通知和历史回执保留边界 | `truzhen-packs` |
| `commerce/buyer-library-install-state-candidate.json` | candidate asset ready，含买家已购库、下载、安装、重新安装、团队设置可替换和撤销阻断状态 | `truzhen-packs` |
| `commerce/publisher-account-settlement-policy-candidate.json` | candidate asset ready，含发布者身份、定价审批、结算、税务 / 发票和上线门槛 | `truzhen-packs` |
| `commerce/commercial-terms-privacy-policy-candidate.json` | candidate asset ready，含购买前条款接受、隐私告知、数据归属、候选输出免责声明和上线阻断条件 | `truzhen-packs` |
| `commerce/release-candidate-package.json` | candidate asset ready，含发布候选、签名、升级和撤回边界 | `truzhen-packs` |
| `commerce/download-install-access-matrix.json` | candidate asset ready，含未购买、sandbox 已购、退款、过期、下架、hash 不一致和缺 Owner Gate 的下载 / 安装访问控制 | `truzhen-packs` |
| `commerce/artifact-bundle-layout-candidate.json` | candidate asset ready，含候选包文件范围、排序、压缩、哈希、签名请求和上传请求 | `truzhen-packs` |
| `commerce/artifact-bundle-digest-candidate.json` | candidate asset ready，含候选包 payload tree SHA-256、上传 / 下载 / 安装 hash 连续性引用和云端回执前置条件 | `truzhen-packs` |
| `commerce/install-compatibility-matrix.json` | candidate asset ready，含下载后安装、升级、回滚和重新启用的版本 / 授权 / 槽位 / 签名 / hash 检查 | `truzhen-packs` |
| `commerce/marketplace-review-submission-candidate.json` | candidate asset ready，含云市场审核提交所需 listing、价格 / 许可、素材授权、支持 / 退款和阻断条件 | `truzhen-packs` |
| `commerce/commercial-go-live-approval-candidate.json` | candidate asset ready，含商品化上线批准总门、下载 / 安装访问控制矩阵硬门、跨仓证据、Owner 裁定、真实支付 / 生产发布阻断和回滚计划 | `truzhen-packs` |
| `commerce/commercial-production-promotion-gate-candidate.json` | candidate asset ready，含 P11 最终证据包、安装后团队办公室运行使用烟测、下载 / 安装访问控制矩阵、独立验收、Owner go/no-go、生产发布、真实支付和生产签名下载的晋级门与阻断口径 | `truzhen-packs` |
| `commerce/commercial-receipt-chain-candidate.json` | candidate asset ready，含上传、审核、购买、entitlement、下载、安装、绑定和运行使用的回执关联键、角色制作血缘键与断链阻断条件 | `truzhen-packs` |
| `commerce/commercial-distribution-receipt-schema-candidate.json` | candidate asset ready，含上传、审核、sandbox 下单 / 支付、entitlement、签名下载、安装、团队绑定、运行使用、负例阻断和独立验收签收的统一回执字段 | `truzhen-packs` |
| `install/install-preflight-request-candidate.json` | candidate asset ready，含下载后安装前 entitlement、hash、签名、schema、禁入产物、槽位、asset ref 和 Owner Gate 预检要求 | `truzhen-packs` |
| `install/install-runtime-activation-map-candidate.json` | candidate asset ready，含下载回执、安装预检、启用版本、团队设置目录刷新、Owner Gate 绑定和运行使用烟测的激活映射 | `truzhen-packs` |
| `tests/gui-user-agent-scenarios.json` | candidate asset ready，含用户视角智能体 GUI 操作场景和证据要求 | `truzhen-packs` |
| `tests/gui-user-agent-execution-script-candidate.json` | candidate asset ready，含用户视角 GUI 有序执行步骤、动作类型、页面状态、证据槽和负例脚本 | `truzhen-packs` |
| `tests/gui-evidence-capture-protocol.json` | candidate asset ready，含用户视角 GUI 证据采集、脱敏、链路关联和缺口回填要求 | `truzhen-packs` |
| `tests/role-studio-export-provenance-attestation-candidate.json` | candidate asset ready，含制作台 GUI 导出来源证明、6 个角色候选引用、candidate bundle export receipt、bundle hash 和手工 JSON / 后端直调负例阻断 | `truzhen-packs` |
| `tests/artifact-secret-raw-asset-scan-candidate.json` | candidate asset ready，含上传前 secret / raw asset 扫描规则、禁入文件扩展名、真实密钥形态、写回目标和阻断条件 | `truzhen-packs` |
| `tests/artifact-manifest-closure-gate-candidate.json` | candidate asset ready，含当前目录、candidate-set 和 artifact manifest 三方文件闭合检查，作为上传 / 下载 / 安装前置门 | `truzhen-packs` |
| `tests/product-readiness-evidence-matrix.json` | candidate asset ready，含商品化完成证据 gates 和未验证阻断口径 | `truzhen-packs` |
| `tests/role-studio-goal-completion-evidence-map-candidate.json` | candidate asset ready，含活跃目标到角色制作、使用、上传、购买、下载、安装、团队设置替换、安装后团队办公室运行使用烟测、商品化上线批准、P11 证据验收清单和正常商品化 go/no-go 权威证据的映射 | `truzhen-packs` |
| `tests/commercialization-execution-packet.json` | candidate asset ready，含发布者 / 购买者商品化全链路执行阶段、证据要求和负例阻断 | `truzhen-packs` |
| `tests/e2e-evidence-run-record.json` | candidate asset ready，含真实跨仓 E2E 执行时需填写的阶段证据槽、负例证据槽和完成判定门 | `truzhen-packs` |
| `tests/commercial-evidence-gate-candidate.json` | candidate asset ready，含每阶段最小证据记录、跨仓关联键、独立验收签收和负例阻断要求 | `truzhen-packs` |
| `tests/sandbox-environment-readiness-candidate.json` | candidate asset ready，含云端 sandbox、买卖双方、支付桩、签名下载、本地安装目标和团队设置入口的 ref-only 前置条件 | `truzhen-packs` |
| `tests/commercial-observability-diagnostics-candidate.json` | candidate asset ready，含 GUI、云端、安装和团队绑定的 trace、日志、指标、告警、脱敏和回执关联要求 | `truzhen-packs` |
| `tests/normal-commercialization-completion-audit-candidate.json` | candidate asset ready，含 Owner 目标逐条完成审计、当前候选证据、缺失权威证据和下一步授权执行顺序 | `truzhen-packs` |
| `tests/commercial-chain-verifier-candidate.json` | candidate asset ready，含正常商品化链路逐阶段核验门、跨仓回执关联键、hash 连续性、负例阻断和独立验收签收要求 | `truzhen-packs` |
| `tests/commercial-go-no-go-gate-candidate.json` | candidate asset ready，含商品化 go/no-go 阶段门、问题台账 triage 终端门、商品阶段前后端收口报告终端门、安装后团队办公室运行使用烟测终端门、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵终端门、完成规则、非充分证据和未授权阻断口径 | `truzhen-packs` |
| `tests/commercial-readiness-verifier-candidate.json` | candidate asset ready，含当前商品化 readiness 阻断汇总、证据回写计数、问题台账写回摘要、商品阶段前后端收口报告、安装后团队办公室运行使用烟测、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵和 go/no-go 前置项 | `truzhen-packs` |
| `tests/role-studio-phase-coverage-matrix-candidate.json` | candidate asset ready，含 P0-P11 阶段到用户视角 GUI 证据、后端回执、候选资产、缺失权威证据、目标仓和阻断状态的覆盖映射 | `truzhen-packs` |
| `tests/role-studio-test-case-coverage-matrix-candidate.json` | candidate asset ready，含计划第 8 节 24 个 `TC-*` 到 GUI 证据、权威回执 / blocked 证据、目标仓和当前阻断状态的覆盖映射 | `truzhen-packs` |
| `tests/independent-acceptance-signoff-matrix-candidate.json` | candidate asset ready，含独立验收智能体对 P0-P11、负例、hash 连续性、云端回执、本地安装回执和 Owner go / no-go 裁定的签收要求 | `truzhen-packs` |
| `tests/p0-p11-commercialization-blocker-register-candidate.json` | candidate asset ready，含 P0-P11 每阶段 issue、目标仓、复现步骤、权威证据、回执要求和不得声明完成的阻断口径 | `truzhen-packs` |
| `tests/product-stage-frontend-backend-closure-report-candidate.json` | candidate asset ready，含 P10 商品阶段前端 GUI 能力、后端 candidate / Gate / Receipt、字段一致性、P0/P1 剩余缺口和下一轮执行卡 | `truzhen-packs` |
| `tests/p11-normal-commercialization-acceptance-gate-candidate.json` | candidate asset ready，含 P11 上传、审核、sandbox 购买 / 支付、entitlement、下载 hash、安装、团队绑定、负例和独立验收的 go / no-go 判定门 | `truzhen-packs` |
| `tests/p11-normal-commercialization-verification-record-template.json` | candidate asset ready，含 P11 真实执行后阶段结果、hash 连续性、负例结果、独立验收和最终决策的验证记录模板 | `truzhen-packs` |
| `tests/p11-evidence-ingestion-binder-candidate.json` | candidate asset ready，含真实 GUI 步骤、云回执、安装回执、团队绑定、负例和独立验收到 P11 验证记录的证据接入绑定规则 | `truzhen-packs` |
| `tests/p11-sandbox-execution-runbook-candidate.json` | candidate asset ready，含 Owner 授权后 P11 sandbox 商品化跨仓 GUI / 云 / 安装执行顺序、证据输出、回滚和禁区 | `truzhen-packs` |
| `tests/p11-sandbox-preflight-gate-candidate.json` | candidate asset ready，含 P11 sandbox 开跑前 Owner 授权、云端 sandbox、买卖双方、支付桩、签名下载、本地安装目标和团队设置入口预检门 | `truzhen-packs` |
| `tests/p11-sandbox-run-request-candidate.json` | candidate asset ready，含 P11 sandbox 授权后由用户视角 GUI 智能体执行上传、购买、下载、安装、团队替换、独立验收和生产晋级控制的十阶段请求入口；当前不可执行 | `truzhen-packs` |
| `tests/p11-evidence-acceptance-checklist-candidate.json` | candidate asset ready，含 P11 sandbox 实操后的上传、购买、下载、安装、团队设置替换、负例、独立验收和 Owner go/no-go 逐阶段验收检查项 | `truzhen-packs` |
| `tests/p11-commercial-go-live-evidence-package-template.json` | candidate asset ready，含 P11 最终证据包的 GUI 证据索引、回执索引、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告、hash 连续性、下载 / 安装访问矩阵报告、负例、独立验收、Owner go / no-go 和最终上线裁定结构；GUI 证据索引必须覆盖用户视角脚本全部 14 个步骤 | `truzhen-packs` |
| `integration/cross-repo-execution-cards.json` | candidate asset ready，含 client / `truzhenos` / `truzhen-cloud` / contracts 授权后施工范围、验证命令、证据和禁区 | `truzhen-packs` |
| `integration/cross-repo-execution-readiness-package.json` | candidate asset ready，含目标仓路径、状态命令、允许动作、证据输出、阶段顺序和 Owner 授权问题 | `truzhen-packs` |
| `integration/owner-authorization-evidence-intake-candidate.json` | candidate asset ready，含 Owner 授权原文、逐仓范围、允许动作、红色禁区、证据输出、过期条件和未授权工作闸门 | `truzhen-packs` |
| `integration/commercial-cross-repo-execution-queue-candidate.json` | candidate asset ready，含商品化跨仓执行队列、阶段依赖、目标仓、证据输出、未授权阻断和完成门 | `truzhen-packs` |
| `docs/commercial-cross-repo-evidence-ledger.json` | candidate asset ready，含 8 个商品化跨仓阶段的 `evidence_id`、目标仓、待填证据位置、写回目标、负例阻断和完成前置门 | `truzhen-packs` |
| `docs/role-studio-issue-ledger.md` | candidate asset ready，含 GUI、云端、go/no-go、P11 证据包等角色制作台缺口；readiness、go/no-go、产品完成矩阵和目标完成地图均要求 `issue_ledger_all_entries_triaged` | `truzhen-packs` |
| `integration/frontend-backend-contract-map.json` | candidate asset ready，含角色制作台商品阶段前端 surface 到后端候选 / Gate / Receipt / 云端证据的接线映射，并含 P11 用户视角 GUI 步骤到商业 API 端点、回执字段、阻断态和 GUI 证据槽的 traceability 矩阵 | `truzhen-packs` |
| `integration/commercial-api-contract-candidate.json` | candidate asset ready，含上传、审核、购买、支付、entitlement、下载、安装、团队设置刷新、生产发布、真实支付启用和生产签名下载启用请求的 request、response、回执、幂等键和错误态；生产晋级端点必须声明 `blocked_p11_evidence_package_incomplete` 与 `blocked_team_office_runtime_usage_smoke_missing` | `truzhen-packs` |
| `integration/commercial-api-example-cases-candidate.json` | candidate asset ready，含每个商品化端点的非执行 request / success response / failure response 样例，用于前后端接线和验收智能体核对；生产晋级失败样例必须分别覆盖 P11 最终证据包缺失和安装后运行使用烟测缺失 | `truzhen-packs` |
| `usage/team-office-runtime-usage-candidate.json` | candidate asset ready，含团队办公室运行使用场景、负例阻断要求和 `team_office_runtime_usage_smoke_verified` 完成前置门 | `truzhen-packs` |
| candidate bundle hash | pending | `truzhenos` / `truzhen-cloud` 打包上传环节 |
| cloud upload receipt | pending | `truzhen-cloud` |
