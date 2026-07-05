# 团队办公室角色包候选集 v0

本目录是角色制作台全链路测试用的 `Role Pack Candidate Set`，用于压测“角色制作、能力 Pack 引用、团队设置替换、秘书长音色 / VRM、云端商品化草稿、下载后安装”的商品化闭环。

当前状态：`设计中`。所有文件都是候选资产，不代表已发布、已上架、已购买、已安装或已启用。

## 边界

- 角色包只作为 Proposer：提出建议、草稿、质询、候选和风险提示。
- 正式裁定、正式任务、正式记忆、正式发送、真实执行必须经 Owner + Base Gate。
- 音色 / VRM 只保存 `asset_ref` 和 `provider_ref`，不保存 raw 音频、raw VRM、真人肖像、账号或凭据。
- 云上传、商品草稿、订单、支付、License、Entitlement、下载分发真相归 `truzhen-cloud`。
- 本地安装、启用、团队设置替换和回执归 `truzhenos`。

## 候选资产

- `candidate-set.json`：候选集总入口。
- `role-packs/*.rolepack.json`：秘书长 + 五顾问角色包候选。
- `role-slots/team-office-role-slots.json`：团队办公室槽位候选。
- `bindings/team-office-role-binding-candidate.json`：团队设置角色绑定候选，覆盖秘书长和五顾问替换、回滚和 Owner Gate。
- `bindings/team-settings-installed-role-catalog-candidate.json`：安装后团队设置可替换角色目录候选，声明已安装角色如何进入团队设置 tab。
- `capability-role-requirements/sample-team-research.role-requirements.json`：能力 Pack 引用角色需求样例。
- `appearance/secretary-appearance-preferences.json`：秘书长音色 / VRM 表现层偏好，声明团队设置页 GUI 选择、清空、恢复默认和 provider readiness 显示规则。
- `appearance/secretary-appearance-asset-rights-candidate.json`：秘书长音色 / VRM 授权证据候选，声明 asset ref、市场审核、安装预检和阻断条件。
- `commerce/cloud-listing-candidate.json`：云端商品草稿候选，含买家可见的六角色商品组件清单、秘书长音色 / VRM asset-ref-only 边界、购买前支持 / 退款 / 撤回 / 候选角色免责声明 surface；不是真实上架记录。
- `commerce/artifact-manifest.json`：上传、下载、安装使用的候选包清单和前置条件。
- `commerce/sandbox-commerce-flow-candidates.json`：sandbox 上传、购买、entitlement、下载、安装流程候选。
- `commerce/product-handoff-candidate.json`：云端商品化与本地安装交接候选，不是真实云 API 或安装实现。
- `commerce/license-entitlement-policy-candidate.json`：License / Entitlement 策略候选，声明购买后下载、安装授权边界和 team-scoped entitlement 必须匹配安装目标团队。
- `commerce/order-payment-state-machine-candidate.json`：订单 / 支付状态机候选，声明 sandbox 订单、支付、失败、退款、chargeback、entitlement 发放和撤销边界。
- `commerce/support-refund-revocation-policy-candidate.json`：售后 / 退款 / 撤回策略候选，声明支持入口、购买前披露、退款、entitlement 撤回、发布撤回通知和历史回执保留边界。
- `commerce/buyer-library-install-state-candidate.json`：买家已购库与安装状态候选，声明已购库、下载、安装、重新安装、团队设置可替换和撤销阻断边界。
- `commerce/publisher-account-settlement-policy-candidate.json`：发布者账号 / 结算 / 发票策略候选，声明发布者身份、定价审批、结算、税务 / 发票和上线门槛。
- `commerce/commercial-terms-privacy-policy-candidate.json`：商业条款 / 隐私 / 数据边界候选，声明购买前条款接受、隐私告知、数据归属、候选输出免责声明和上线阻断条件。
- `commerce/release-candidate-package.json`：发布候选包，声明版本、渠道、签名、升级、撤回边界，并引用云商品草稿里的六角色商品组件清单。
- `commerce/download-install-access-matrix.json`：下载 / 安装访问控制矩阵候选，声明未购买、退款、授权过期、版本下架、hash 不一致和缺 Owner Gate 的阻断。
- `commerce/artifact-bundle-layout-candidate.json`：确定性打包规格候选，声明候选包文件范围、排序、压缩、哈希、签名请求和上传请求。
- `commerce/artifact-bundle-digest-candidate.json`：候选包整体 digest，声明上传、下载、安装共同引用的 payload tree SHA-256；当前只是本仓候选哈希，不是云端上传回执。
- `commerce/install-compatibility-matrix.json`：安装兼容性矩阵候选，声明下载后本地安装、升级、回滚和重新启用前的版本、授权、槽位、签名 / hash 检查。
- `commerce/marketplace-review-submission-candidate.json`：云市场审核提交候选，声明 listing、价格 / 许可、素材授权、支持 / 退款、审核清单和阻断条件。
- `commerce/commercial-go-live-approval-candidate.json`：商品化上线批准候选，汇总 GUI、云端、支付、下载、安装、下载 / 安装访问控制矩阵、P11 证据验收清单、阶段依赖、负例和独立验收证据门；未获 Owner 裁定和跨仓证据前保持不可上线。
- `commerce/commercial-production-promotion-gate-candidate.json`：生产发布晋级门候选，要求 P11 最终证据包、安装后团队办公室运行使用烟测、下载 / 安装访问控制矩阵、独立验收和 Owner go/no-go 齐备后，才可请求生产上架、真实支付和生产签名下载。
- `commerce/commercial-receipt-chain-candidate.json`：商品化回执链候选，声明上传、审核、购买、entitlement、下载、安装、绑定和运行使用之间的关联键与断链阻断；同时要求 `candidate_set_ref`、`bundle_tree_sha256`、6 个角色引用集合、安装后启用版本和团队 slot 映射贯穿到运行使用。
- `commerce/commercial-distribution-receipt-schema-candidate.json`：商业分发回执记录规格候选，声明上传、审核、购买、支付、授权、下载、安装、团队绑定、运行使用、负例阻断和独立验收签收必须保留的回执字段。
- `install/install-preflight-request-candidate.json`：安装预检请求候选，声明下载后安装前的 entitlement、授权团队、目标团队、hash、签名、schema、禁入产物、槽位、asset ref 和 Owner Gate 检查。
- `install/install-runtime-activation-map-candidate.json`：安装到运行使用激活映射候选，把下载回执、安装预检、启用版本、团队设置目录刷新、Owner Gate 绑定和运行使用烟测串成可验链路。
- `tests/gui-user-agent-scenarios.json`：用户视角智能体 GUI 场景候选，约束后续不能绕过前端验收。
- `tests/gui-user-agent-execution-script-candidate.json`：用户视角 GUI 执行脚本候选，按顺序列出从角色制作到上传、购买、下载、安装、团队替换、运行使用和负例阻断的 GUI 步骤。
- `tests/gui-evidence-capture-protocol.json`：GUI 证据采集协议候选，约束截图、页面状态、网络 / 回执摘要、脱敏、链路关联和缺口回填。
- `tests/role-studio-export-provenance-attestation-candidate.json`：角色制作台导出来源证明候选，要求用户视角 GUI 导出证据、candidate bundle export receipt、6 个角色候选引用和 bundle hash 齐备后才可进入云上传。
- `tests/artifact-secret-raw-asset-scan-candidate.json`：上传前 secret / raw asset 扫描候选，约束云上传、市场审核、签名下载、安装预检和生产晋级前不得夹带密钥、凭据、签名 URL、raw 音频或 raw VRM。
- `tests/artifact-manifest-closure-gate-candidate.json`：artifact manifest 闭合门候选，约束当前目录、candidate-set 和 artifact manifest 必须使用同一组文件，防止漏打包或多打包进入上传 / 下载 / 安装链路。
- `tests/product-readiness-evidence-matrix.json`：商品化完成证据矩阵候选，约束后续不能用窄测试冒充全链路完成。
- `tests/role-studio-goal-completion-evidence-map-candidate.json`：目标完成证据地图候选，把活跃目标拆成角色制作、使用、上传、购买、下载、安装、团队设置替换和正常商品化 go/no-go 的权威证据要求，并把商品化上线批准、P11 证据验收清单、问题台账全量 triage、商品阶段前后端收口报告和安装后团队办公室运行使用烟测列为完成声明前置门。
- `tests/commercialization-execution-packet.json`：发布者 / 购买者商品化全链路执行包候选，约束角色制作、上传、购买、下载、安装、团队替换和运行使用必须成对留下 GUI 与后端证据。
- `tests/e2e-evidence-run-record.json`：E2E 证据运行记录包候选，供后续真实跨仓 GUI、云端、安装和负例验收逐项填写证据。
- `tests/commercial-evidence-gate-candidate.json`：商品化证据门槛候选，定义每阶段最小证据记录、跨仓关联键、独立验收签收和负例阻断要求。
- `tests/sandbox-environment-readiness-candidate.json`：sandbox 商品化环境就绪候选，声明云端 sandbox、买卖双方、支付桩、签名下载、本地安装目标和团队设置入口的 ref-only 前置条件。
- `tests/commercial-observability-diagnostics-candidate.json`：商品化监控 / 诊断候选，声明 GUI、云端、安装和团队绑定的 trace、日志、指标、告警、脱敏和回执关联要求。
- `tests/normal-commercialization-completion-audit-candidate.json`：正常商品化完成审计候选，按 Owner 目标逐条映射当前候选证据与仍缺的 GUI、云端、安装和独立验收权威证据。
- `tests/commercial-chain-verifier-candidate.json`：正常商品化链路核验器候选，要求跨仓执行后逐项核对 GUI 证据、云端回执、支付 / entitlement、下载 hash、安装回执、团队绑定回执、P11 执行队列阶段依赖证明、负例阻断和独立验收签收。
- `tests/commercial-go-no-go-gate-candidate.json`：商品化 go/no-go 门禁候选，消费执行队列、跨仓证据台账、P11 阶段依赖证明、P11 证据包、问题台账、商品阶段前后端收口报告、安装后团队办公室运行使用烟测、角色制作台血缘报告、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵和上线批准，当前明确阻断商品化就绪。
- `tests/commercial-readiness-verifier-candidate.json`：商品化 readiness verifier 候选，汇总当前未授权、证据未回写、问题台账未 triage、商品阶段前后端收口报告未验证、安装后团队办公室运行使用烟测未验证、角色制作台血缘报告、秘书长音色 / VRM 证据报告、下载 / 安装访问控制矩阵、商品化上线批准和 go/no-go 阻断状态，当前明确 `blocked_not_commercial_ready`。
- `tests/role-studio-phase-coverage-matrix-candidate.json`：P0-P11 阶段覆盖矩阵候选，把计划阶段逐项映射到用户视角 GUI 证据、后端回执、当前候选资产、缺失权威证据、目标仓和阻断状态。
- `tests/role-studio-test-case-coverage-matrix-candidate.json`：测试用例覆盖矩阵候选，把计划第 8 节 24 个 `TC-*` 逐项映射到 GUI 证据、权威回执 / blocked 证据、目标仓和当前阻断状态。
- `tests/independent-acceptance-signoff-matrix-candidate.json`：独立验收签收矩阵候选，声明验收智能体如何逐阶段复核 GUI 证据、GUI/API traceability、云端回执、本地安装回执、负例阻断和 Owner go / no-go 裁定。
- `tests/p0-p11-commercialization-blocker-register-candidate.json`：P0-P11 商品化阻塞清单候选，把未获授权或缺少权威证据的阶段登记成 issue、目标仓、复现步骤、所需回执和独立验收要求。
- `tests/product-stage-frontend-backend-closure-report-candidate.json`：商品阶段前后端收口报告候选，把 P10 要求的前端 GUI 能力、后端 candidate / Gate / Receipt、字段一致性、P0/P1 剩余缺口和下一轮执行卡集中成可复核入口。
- `tests/p11-normal-commercialization-acceptance-gate-candidate.json`：P11 正常商品化通过门候选，按上传、审核、sandbox 购买 / 支付、entitlement、下载 hash、安装、团队绑定、负例和独立验收判定 go / no-go。
- `tests/p11-normal-commercialization-verification-record-template.json`：P11 正常商品化验证记录模板，定义真实执行后如何填写阶段结果、hash 连续性、P11 阶段依赖结果、云上传商品草稿结果、商品页 / 审核合规结果、支付购买 / 授权结果、已购签名下载交付结果、下载 / 安装访问矩阵结果、角色制作台血缘结果、负例结果、独立验收和最终决策。
- `tests/p11-evidence-ingestion-binder-candidate.json`：P11 证据接入绑定器候选，定义真实 GUI 步骤、云回执、安装回执、团队绑定、P11 阶段依赖证据、云上传商品草稿证据、商品页 / 审核合规证据、支付购买 / 授权证据、已购签名下载交付证据、下载 / 安装访问矩阵证据、角色制作台血缘证据、负例和独立验收如何绑定到 P11 验证记录。
- `tests/p11-sandbox-execution-runbook-candidate.json`：P11 sandbox 商品化执行 runbook 候选，定义 Owner 授权后跨仓 GUI / 云 / 安装执行顺序、证据输出、回滚和禁区。
- `tests/p11-sandbox-preflight-gate-candidate.json`：P11 sandbox 开跑前预检门候选，定义 Owner 授权、sandbox 环境、支付桩、签名下载、本地安装目标和团队设置入口的开跑前阻断条件。
- `tests/p11-sandbox-run-request-candidate.json`：P11 sandbox 执行请求候选，定义授权后由发布者 / 买家用户视角 GUI 智能体执行上传、购买、下载、安装、团队替换、独立验收和生产晋级控制的十阶段请求；商业阶段逐项绑定 GUI 证据协议、权威回执、阶段依赖、前置回执传递和同一 `bundle_tree_sha256` 关联键；当前不可执行。
- `tests/p11-evidence-acceptance-checklist-candidate.json`：P11 证据验收检查清单候选，定义上传、购买、下载、安装、团队设置替换、GUI/API traceability、P11 阶段依赖证明、角色制作台血缘报告、负例、独立验收和 Owner go/no-go 的逐阶段验收条件。
- `tests/p11-commercial-go-live-evidence-package-template.json`：P11 商品化最终证据包模板，定义 GUI 证据索引、GUI/API traceability 报告、云 / 安装 / 团队绑定回执索引、秘书长音色 / VRM 证据报告、安装后能力 Pack 引用报告、hash 连续性、P11 阶段依赖报告、云上传商品草稿报告、商品页 / 市场审核 / 购买前披露合规报告、支付购买 / 授权报告、角色制作台血缘报告、下载 / 安装访问矩阵报告、负例、独立验收、Owner go / no-go 和最终上线裁定的收口结构；GUI 证据索引必须覆盖用户视角执行脚本全部 14 个步骤，GUI/API traceability 必须逐项对上商业 API 端点、回执字段、阻断态和 GUI 证据槽。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `install_catalog_and_slot_mapping_report` 要求团队设置目录刷新回执、六个可替换角色引用、enabled role pack version、slot 映射、团队设置截图和秘书长 / 五顾问映射截图齐备；readiness、go/no-go、产品完成矩阵和目标完成地图必须消费 `six_installed_role_catalog_evidence_verified` 后才可声称安装后可替换角色目录已就绪。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `cloud_upload_listing_report` 要求云上传回执、商品草稿引用、artifact manifest、候选包 digest、发布候选包、禁入产物扫描和 hash 匹配回执齐备；P11 证据验收检查清单、证据接入绑定器和验证记录模板必须把这些证据写回 `cloud_upload_listing_result`，且 readiness、go/no-go、产品完成矩阵和目标完成地图必须消费 `cloud_upload_listing_verified` 后才可计入 P11 通过。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `marketplace_listing_review_compliance_report` 要求商品页草稿、六角色商品组件清单、市场审核候选、购买前支持 / 退款 / 撤回披露、条款隐私数据边界、发布者身份 / 定价审批和生产发布阻断回执齐备；P11 证据验收检查清单、证据接入绑定器和验证记录模板必须把这些证据写回 `marketplace_listing_review_result` 后才可计入 P11 通过。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `purchase_entitlement_report` 要求 sandbox 订单、sandbox 支付、无真实扣款、License / Entitlement、团队授权范围、买家已购库 GUI、失败支付阻断和退款 / chargeback 撤权证据齐备；P11 证据验收检查清单、证据接入绑定器和验证记录模板必须把这些证据写回 `purchase_entitlement_result` 后才可计入 P11 通过。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `download_artifact_delivery_report` 要求已购库下载 GUI、entitlement 校验、签名下载引用、下载回执、下载 artifact hash、候选包 digest 和 hash 匹配回执齐备；P11 证据验收检查清单、证据接入绑定器和验证记录模板必须把这些证据写回 `download_artifact_delivery_result`，且 readiness、go/no-go、产品完成矩阵和目标完成地图必须消费 `download_artifact_delivery_verified` 后才可声称角色包可下载安装。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `download_install_access_matrix_report` 要求未购买、退款撤权、授权过期、版本下架 / 撤回和 artifact hash 不一致 blocked 回执齐备；P11 证据验收检查清单、证据接入绑定器和验证记录模板必须把这些证据写回 `download_install_access_result` 后才可计入 P11 通过。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的 `role_studio_lineage_report` 要求同一 6 个角色引用集合和 `role_pack_ref_set_hash` 从角色制作台导出、云上传、下载、安装、团队设置替换延续到团队办公室运行使用，缺少任一权威回执都不得通过正常商品化。
- 商品化 readiness verifier、go/no-go 门禁和产品完成矩阵已直接消费 `commercial_go_live_approval_verified`；缺少上线批准候选、P11 证据验收清单、最终证据包、Owner go/no-go 和生产发布 / 真实支付阻断证据时，不得通过正常商品化。
- `tests/p11-commercial-go-live-evidence-package-template.json` 的独立验收 section 还必须记录 `gui_api_traceability_matrix_reviewed`，确保验收智能体复核该矩阵后才允许进入 Owner go / no-go。
- `integration/cross-repo-execution-cards.json`：跨仓执行卡候选，拆分 client、`truzhenos`、`truzhen-cloud`、contracts 和生产晋级控制的授权后施工范围、验证命令、证据和禁区。
- `integration/cross-repo-execution-readiness-package.json`：跨仓执行就绪包候选，汇总目标仓路径、允许动作、状态命令、证据输出、阶段顺序、生产晋级阶段和 Owner 授权问题。
- `integration/owner-authorization-evidence-intake-candidate.json`：跨仓 Owner 授权证据接入口候选，定义授权原文、逐仓范围、允许动作、红色禁区、证据输出、过期条件和未授权工作闸门。
- `integration/commercial-cross-repo-execution-queue-candidate.json`：商品化跨仓执行队列候选，按授权 intake、contracts、后端回执、用户视角 GUI、云端 sandbox、安装绑定、负例观测和独立验收顺序串联阶段闸门，并镜像 P11 阶段依赖和前置回执阻断；其中 `frontend_user_view_gui_flow` 必须输出 `gui_api_traceability_matrix_verified`、商业 API 端点追踪、回执字段追踪和 GUI 证据槽追踪。
- `docs/commercial-cross-repo-evidence-ledger.json`：商品化跨仓证据台账候选，为执行队列每个阶段声明 `evidence_id`、目标仓、证据引用、写回目标、P11 阶段依赖、前置回执和阻断原因；其中用户视角 GUI 阶段必须把 GUI/API traceability 写回 P11 最终证据包、P11 验收清单和 P10 收口报告；当前等待 Owner 授权。
- `integration/frontend-backend-contract-map.json`：前端商品阶段 surface 到后端候选、Gate、Receipt 和云端证据的接线映射候选，约束 GUI 成功必须有后端证据承接；其中 `gui_api_traceability_matrix` 把 P11 用户视角 GUI 步骤逐项映射到商业 API 端点、回执字段、阻断态和 GUI 证据槽，并被 P11 最终证据包、验收清单、readiness verifier、go/no-go、商品完成矩阵和目标完成地图共同消费。
- `integration/commercial-api-contract-candidate.json`：商品化前后端 API 接线契约候选，声明上传、审核、购买、支付、entitlement、下载、安装、团队设置刷新、生产发布、真实支付启用和生产签名下载启用请求的 request、response、回执和错误态；生产晋级端点必须区分 P11 最终证据包缺失和安装后运行使用烟测缺失。
- `integration/commercial-api-example-cases-candidate.json`：商品化 API 请求 / 响应样例候选，为每个商品化端点提供非执行示例，便于前端、后端、云端和验收智能体按同一字段核对；生产晋级失败样例必须返回可读 blocked 原因和候选 / 回执引用。
- `usage/team-office-runtime-usage-candidate.json`：团队办公室运行使用候选，声明安装后秘书长与五顾问如何被调用，并被 readiness、go/no-go、产品完成矩阵和目标完成地图消费为 `team_office_runtime_usage_smoke_verified`。
- `docs/role-studio-issue-ledger.md`：角色制作台问题台账，登记 GUI、云端、go/no-go、P11 证据包等缺口；readiness、go/no-go、产品完成矩阵和目标完成地图必须消费 `issue_ledger_all_entries_triaged` 后才能声称商品化完成。
- `docs/*.md`：制作、绑定、导出、云商品化台账占位。

## 验收

```sh
go test ./... -run TestTeamOffice -count=1
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('role-pack-candidates/team-office-v0/**/*.json',recursive=True)];print('JSON 合法')"
```

收尾时还需运行仓库统一的敏感信息扫描；本候选集不得出现 raw secret、账号凭据、真实 License 或支付回调配置。
