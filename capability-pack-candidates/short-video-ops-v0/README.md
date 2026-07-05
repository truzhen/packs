# 短视频运营能力包候选集 v0

本目录是能力包制作台压测产物，不是已发布 Pack，也不是 provider 实现。目标是用短视频运营过程中的真实能力需求，验证 Truzhen 能否在能力包制作台中接收 GitHub OSS 候选、生成受控胶水计划，并沉淀一组可继续加工的 Capability Pack 候选。

生命周期档位：`设计中`。分片进度单独记录：P3-P11 已接线到候选信封、11 运行请求候选、candidate bundle 导出 / dry-run、P11 lifecycle preflight 阻断、PatchCandidate 结果承接和复核候选；P12-P18 仍处于设计中 / 待授权，执行就绪包已内置 `evidence_writeback_plan`，P12-P18 分片级 `cross_repo_work_gate` 已统一消费授权尝试覆盖 verifier、下一授权启动守卫和 Pack Studio issue 台账关闭门，商用 execution readiness guard coverage verifier 已核对每个分片执行就绪包和跨仓队列镜像 gate 的授权与 issue 台账守卫覆盖，商用跨仓执行队列已汇总每个分片的写回摘要，商用 readiness verifier 已汇总写回总数并消费 P12 / P13 / P15 / P16 / P17 / P18 运行后证据验收门，商用 go/no-go 候选门禁已把写回未完成和 P12 / P13 / P15 / P16 / P17 / P18 后验收未通过列为阻断，商用后验收门覆盖 verifier 已核对每个必经分片在候选集、证据索引、readiness verifier 和 go/no-go 门中的后验收门引用，商用禁入动作覆盖 verifier 已核对 9 个禁入动作终态检查在 candidate-set、证据索引、readiness verifier、go/no-go 门和目标完成图中的引用，商用独立验收签收矩阵已把 P12-P18、GUI/API/Receipt、issue 关闭证据、禁入动作、证据写回、Owner/Base Gate 和最终机器证据包纳入独立 reviewer 复核门，用于把授权执行后的每个 `evidence_id`、终态 false 证据和验收签收写回对应台账。P14 商用 readiness 审计已起草，但未获跨仓施工授权，没有安装、启用、发布或正式执行短视频能力。

## 候选能力包

| 候选 | pack ref | 样本 OSS | 当前状态 |
| --- | --- | --- | --- |
| 短视频草稿生成能力包 | `capability-pack://short-video-draft-generation` | MoneyPrinterTurbo | candidate-only；允许胶水补丁候选；provider 未接通时 `provider_missing` |
| 短视频合成编排能力包 | `capability-pack://short-video-composition-orchestration` | Pixelle-Video | candidate-only；允许胶水补丁候选；模型 / 媒体 provider 未接通时 `provider_missing` |
| 短视频发布草稿能力包 | `capability-pack://short-video-social-publish-draft` | social-auto-upload | blocked-by-default；只能生成发布草稿 / 上传意图候选，不真实发布 |

## 主权边界

- 本目录不保存账号、密码、API key、GPT/Codex 凭据原文或社媒登录态。
- 本目录不包含 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload 的源码、依赖、构建产物或运行脚本。
- 制作台现在能生成 Code Assistant 调用候选信封；代码助手只能引用 `capability-pack://software-code-assistant`，真实 Codex CLI 运行归 11 Execution Gateway，必须经 Owner + Base Gate + Receipt。
- 制作台现在能生成 Code Assistant 运行请求候选；该候选只记录 11 endpoint、Gate 动作、provider、skill、范围和证据，`DecisionRef` / `RunID` / `Nonce` / `ReceiptRef` 必须等待 Base 签发，不能直接执行。
- 制作台现在能导出 candidate-only Pack bundle 文件清单；该 bundle 不可安装、不可启用、不含 provider 实现。
- 制作台现在能对 candidate bundle 做 dry-run 静态校验；独立 Capability Pack loader 未接线时必须 `capability_pack_loader_missing` 阻断，不能借其它 loader 冒充安装成功。
- P11 lifecycle preflight 已接线：候选 bundle 必须先被识别为“不是正式 delivery”，并展示 delivery、黄金用例、evaluation ready、provider dependency、Owner/Base Gate 和 Receipt 缺口；不得直接启用。
- P12 只允许用基座内置安全样本或 fixture 验证 lifecycle draft / readiness / promote / confirm；不能把 MoneyPrinterTurbo、Pixelle-Video 或 social-auto-upload 当成可运行 provider。
- 制作台现在能承接 11 Code Assistant run 返回的 PatchCandidate 结果摘要；该承接记录不可自动应用、必须复核、不能写正式事实。
- 制作台现在能生成 PatchCandidate 复核候选；复核通过只允许进入 apply 候选 Gate，不代表补丁已应用。
- 发布、上传、发送、登录平台等外部副作用只能是候选，正式动作归 Communication / Execution Gateway。

## 文件说明

- `candidate-set.json`：候选集索引和测试边界。
- `capability-packs/*.capability-pack.json`：三个 Capability Pack 候选声明。
- `docs/code-assistant-invocation-ledger.md`：P3 调用候选信封台账；记录隔离目录、模型额度来源、凭据策略、禁止动作和未真实运行状态。
- `docs/code-assistant-run-request-ledger.md`：P8 运行请求候选台账；记录 11 run envelope 候选、Base 待签发字段和不可执行状态。
- `docs/candidate-bundle-ledger.md`：P4-P6 候选 bundle 导出台账；记录文件清单、readiness、不可安装 / 不含 provider 边界。
- `docs/candidate-bundle-dry-run-ledger.md`：P10 候选 bundle dry-run 台账；记录静态校验、`capability_pack_loader_missing`、不可安装 / 不可启用边界。
- `docs/lifecycle-preflight-commercial-gap-ledger.md`：P11 lifecycle preflight 商用缺口台账；记录三类短视频能力候选进入正式 lifecycle 前的阻断项。
- `docs/p12-safe-lifecycle-evidence-contract.json`：P12 安全内置 lifecycle 样本机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p12-cross-repo-authorization-scope.json`：P12 跨仓施工机器可验授权范围；列出允许仓、禁改仓、禁入动作和授权语必含项；当前 `pending_owner_authorization`。
- `docs/safe-lifecycle-sample-ledger.md`：P12 安全内置 lifecycle 样本人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_authorization`。
- `docs/p12-execution-readiness-package.json`：P12 跨仓执行就绪包候选，绑定授权摄取契约、前后端验收契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization`。
- `docs/p12-pre-run-gate.json`：P12 开跑前门禁候选，合并授权摄取、授权范围、P11 复验、执行就绪包、Pack Studio issue 台账关闭门、禁入动作和证据写回目标；当前 `blocked_pending_owner_authorization`，不得据此开始跨仓施工。
- `docs/p12-post-run-evidence-acceptance-gate.json`：P12 运行后证据验收门候选，覆盖 P12 后端、前端和禁入动作全部 `evidence_id`；当前 `blocked_pending_p12_execution_evidence`，已被商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P12 完成、进入 P13 或请求商用签字。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-safe-lifecycle-sample-execution-spec-20260704.md`：P12 安全内置能力 lifecycle 样本执行规格；不执行第三方 OSS。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`：P12 跨仓施工授权卡；待 Owner 明确授权。
- `docs/p12-owner-authorization-evidence-intake.json`：P12 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p13-gui-lifecycle-panel-evidence-contract.json`：P13 GUI lifecycle 面板机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p13-cross-repo-authorization-scope.json`：P13 GUI lifecycle 面板机器可验授权范围；当前 `pending_owner_authorization`。
- `docs/p13-owner-authorization-evidence-intake.json`：P13 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p13-execution-readiness-package.json`：P13 GUI lifecycle 面板跨仓执行就绪包候选，绑定 P12 证据前置条件、授权摄取契约、前后端验收契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization_and_p12_evidence`。
- `docs/p13-post-run-evidence-acceptance-gate.json`：P13 运行后证据验收门候选，覆盖 GUI lifecycle 面板前端、后端和禁入动作全部 `evidence_id`；当前 `blocked_pending_p13_execution_evidence`，已被 P15 执行就绪包、商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P13 完成、进入 P15 或请求商用签字。
- `docs/gui-lifecycle-panel-ledger.md`：P13 GUI lifecycle 面板人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_authorization`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-gui-lifecycle-panel-execution-spec-20260704.md`：P13 GUI lifecycle 面板执行规格；防止绕过 readiness / Gate / Receipt。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-cross-repo-execution-authorization-20260704.md`：P13 跨仓施工授权卡；待 Owner 明确授权。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-commercial-readiness-audit-20260704.md`：P14 商用 readiness 审计；当前结论为 `not_commercial_ready`。
- `docs/commercial-readiness-current-state-audit.json`：P14 商用 readiness 当前状态机器审计；绑定 P12-P18 授权缺口、禁入动作终态检查、Pack Studio issue 台账关闭门和完成门，当前 `not_commercial_ready`。
- `docs/p15-three-candidate-gui-walkthrough-evidence-contract.json`：P15 三候选 GUI 实操机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p15-cross-repo-authorization-scope.json`：P15 三候选 GUI 实操机器可验授权范围；当前 `pending_owner_authorization`。
- `docs/p15-owner-authorization-evidence-intake.json`：P15 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p15-execution-readiness-package.json`：P15 三候选 GUI 实操跨仓执行就绪包候选，绑定 P13 证据前置条件、三候选 Pack refs、授权摄取契约、前后端验收契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization_and_p13_evidence`。
- `docs/p15-post-run-evidence-acceptance-gate.json`：P15 运行后证据验收门候选，覆盖三候选 GUI walkthrough 前端、后端和禁入动作全部 `evidence_id`；当前 `blocked_pending_p15_execution_evidence`，已被 P16 执行就绪包、商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P15 完成、进入 P16 或请求商用签字。
- `docs/gui-walkthrough-evidence-ledger.md`：P15 三候选 GUI 实操人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_gui_walkthrough`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-three-candidate-gui-walkthrough-spec-20260704.md`：P15 三候选 GUI 实操验收规格。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-cross-repo-execution-authorization-20260704.md`：P15 跨仓测试授权卡；待 Owner 明确授权。
- `docs/p16-controlled-code-assistant-run-evidence-contract.json`：P16 受控 Code Assistant run 机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p16-cross-repo-authorization-scope.json`：P16 受控 Code Assistant run 机器可验授权范围；当前 `pending_owner_authorization`。
- `docs/p16-owner-authorization-evidence-intake.json`：P16 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p16-execution-readiness-package.json`：P16 受控 Code Assistant run 跨仓执行就绪包候选，绑定 P15 证据和 P15 运行后证据验收门前置条件、11 Gateway / Owner Gate / Receipt 约束、授权摄取契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization_and_p15_evidence`。
- `docs/p16-post-run-evidence-acceptance-gate.json`：P16 运行后证据验收门候选，覆盖受控 Code Assistant run、PatchCandidate、GUI 展示和禁入动作全部 `evidence_id`，并拒绝用静态胶水代码、本地 Codex CLI 登录态、手工 PatchCandidate 文件或 P8 run request candidate 替代 11 Gateway run receipt；当前 `blocked_pending_p16_execution_evidence`，已被 P17 执行就绪包、商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P16 完成、进入 P17 或请求商用签字。
- `docs/code-assistant-controlled-run-ledger.md`：P16 受控 Code Assistant run 人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_authorization`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-controlled-code-assistant-run-spec-20260704.md`：P16 受控 Code Assistant 最小 run 执行规格。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-cross-repo-execution-authorization-20260704.md`：P16 跨仓授权卡；待 Owner 明确授权。
- `docs/p17-provider-adapter-candidate-evidence-contract.json`：P17 provider / adapter candidate 机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p17-cross-repo-authorization-scope.json`：P17 provider / adapter candidate 机器可验授权范围；当前 `pending_owner_authorization`。
- `docs/p17-owner-authorization-evidence-intake.json`：P17 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p17-execution-readiness-package.json`：P17 provider / adapter candidate 跨仓执行就绪包候选，绑定 P16 证据和 P16 运行后证据验收门前置条件、provider 候选 refs、默认 readiness、授权摄取契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization_and_p16_evidence`。
- `docs/p17-post-run-evidence-acceptance-gate.json`：P17 运行后证据验收门候选，覆盖 provider / adapter 归属、默认 readiness、GUI 展示和禁入动作全部 `evidence_id`，并拒绝用 packs 内 scaffold、未绑定 readiness receipt 的 provider manifest、一句 provider ready 声明或 vendor 第三方 OSS 源码替代外部 provider 仓证据；当前 `blocked_pending_p17_provider_adapter_evidence`，已被 P18 执行就绪包、商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P17 完成、进入 P18 或请求商用签字。
- `docs/provider-adapter-candidate-ledger.md`：P17 provider / adapter candidate 人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_authorization`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-provider-adapter-candidate-spec-20260704.md`：P17 provider / adapter candidate 执行规格。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-cross-repo-execution-authorization-20260704.md`：P17 跨仓授权卡；待 Owner 明确授权。
- `docs/p18-cloud-market-sandbox-evidence-contract.json`：P18 云市场 sandbox 机器可验必填证据契约；当前 `pending_authorization`。
- `docs/p18-cross-repo-authorization-scope.json`：P18 云市场 sandbox 机器可验授权范围；当前 `pending_owner_authorization`。
- `docs/p18-owner-authorization-evidence-intake.json`：P18 Owner 授权证据摄取契约；要求记录授权原话、授权卡、允许仓、禁改仓和禁入动作，当前 `missing_owner_authorization`。
- `docs/p18-execution-readiness-package.json`：P18 云市场 sandbox 跨仓执行就绪包候选，绑定 P17 证据和 P17 运行后证据验收门前置条件、sandbox-only、cloud 真相源、授权摄取契约、目标仓、验证命令、禁入动作、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门和证据写回计划；当前 `blocked_pending_owner_authorization_and_p17_evidence`。
- `docs/p18-post-run-evidence-acceptance-gate.json`：P18 运行后证据验收门候选，覆盖 sandbox listing、License / Entitlement、download、install preflight、GUI 展示和禁入动作全部 `evidence_id`，并拒绝用 packs listing draft、sandbox runbook、License / Entitlement 文案、人工订单状态或生产发布声称替代 `truzhen-cloud` receipt；当前 `blocked_pending_p18_cloud_sandbox_evidence`，已被商用 readiness verifier 和 go/no-go 门禁消费，未通过前不得标记 P18 完成或请求商用签字。
- `docs/cloud-market-sandbox-ledger.md`：P18 云市场 sandbox 人工证据台账骨架，覆盖机器契约 `evidence_id`；当前 `pending_authorization`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cloud-market-sandbox-spec-20260704.md`：P18 云市场 sandbox 执行规格。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cross-repo-execution-authorization-20260704.md`：P18 跨仓授权卡；待 Owner 明确授权。
- `docs/commercial-evidence-contract-index.json`：P12-P18 商用证据契约总索引，绑定授权卡、台账、执行顺序、分片依赖、逐分片签字矩阵、禁入动作终态检查、授权尝试覆盖 verifier、下一授权启动守卫、签字门和完成门；当前 `not_commercial_ready`。
- `docs/commercial-readiness-verifier.json`：商用 readiness 机器判定候选，汇总 P12-P18 授权 / 证据阻断、P12 / P13 / P15 / P16 / P17 / P18 运行后证据验收门、商用后验收门覆盖 verifier、商用禁入动作覆盖 verifier、证据写回总数、禁入动作终态检查、下一授权卡和下一授权启动守卫；当前 `blocked_not_commercial_ready`，不写正式事实。
- `docs/commercial-frontend-backend-acceptance-contract.json`：能力 Pack 制作台前后端商用验收契约候选，绑定 P12/P13/P15 的后端 endpoint、前端 surface、E2E 证据、GUI/API/Receipt 追踪矩阵和 Pack Studio issue 台账关闭门；当前 `blocked_pending_cross_repo_authorization_and_evidence`。
- `docs/commercial-frontend-backend-handoff-runbook.json`：能力 Pack 制作台前后端证据交接 runbook 候选，按 P12 后端、P13 前端、P15 E2E、GUI/API/Receipt 追踪矩阵、Pack Studio issue 台账关闭门和禁入终态检查拆分执行顺序、验收命令、写回台账和阻断动作；当前 `blocked_pending_cross_repo_authorization_and_evidence`。
- `docs/commercial-cross-repo-execution-queue.json`：P12-P18 商用跨仓执行队列候选，按 required_slices 顺序汇总授权、执行就绪包、证据契约、目标仓、授权尝试覆盖 verifier、下一授权启动守卫、Pack Studio issue 台账关闭门、分片级阻断门和证据写回摘要；当前 `blocked_pending_owner_authorization_and_evidence`。
- `docs/commercial-improvement-backlog.json`：能力 Pack 制作台商用改进清单候选，按 P12/P13/P15/P16/P17/P18 拆出改进卡、目标仓、授权卡、证据要求和禁入动作；当前 `blocked_pending_cross_repo_authorization_and_evidence`，不授权施工。
- `docs/pack-studio-issue-ledger.json`：能力 Pack 制作台 issue 回填台账候选，按 M7 要求集中登记 P10-P18 制作台缺口、阻断原因、归属、验收证据、权威关闭证据槽和禁止绕过方式；每个 issue 必须绑定实现 / 授权证据、验证命令、GUI/API/ReadModel 证据、Receipt / Candidate 引用和禁入方案缺席证明后才可关闭；当前 `blocked_pending_issue_resolution_evidence`，不代表问题已解决。
- `docs/commercial-go-no-go-gate.json`：商用 go/no-go 候选门禁，汇总 P12-P18 分片门、P12 / P13 / P15 / P16 / P17 / P18 运行后证据验收门、商用后验收门覆盖 verifier、商用禁入动作覆盖 verifier、禁入动作终端门、证据写回门、Owner/Base Receipt 前置条件、下一授权卡和下一授权启动守卫；当前 `blocked_not_ready_for_commercial_go_live`，不能请求商用签字。
- `docs/commercial-post-run-gate-coverage-verifier.json`：商用后验收门覆盖 verifier 候选，核对 P12-P18 每个必经切片的 `post_run_evidence_acceptance_gate` 已同步到 candidate-set、商用证据索引、readiness verifier 和 go/no-go 门；当前 `blocked_pending_all_post_run_gates_passed`，只证明覆盖关系，不代表执行证据已存在。
- `docs/commercial-forbidden-action-coverage-verifier.json`：商用禁入动作覆盖 verifier 候选，核对 9 个禁入动作终态检查已同步到 candidate-set、商用证据索引、readiness verifier、go/no-go 门和目标完成图；当前 `blocked_pending_forbidden_action_terminal_evidence`，只证明覆盖关系，不代表禁入动作终态 false 证据已存在。
- `docs/commercial-authorization-attempt-coverage-verifier.json`：商用授权尝试覆盖 verifier 候选，核对被拒绝的 P12 授权尝试已同步到商用证据索引、当前状态审计、改进清单和前后端交接文档，并明确不能当作 Owner 授权。
- `docs/commercial-next-authorization-start-guard.json`：商用下一授权启动守卫候选，集中校验 P12 授权卡原话、允许仓、禁改仓、禁入动作和安全 fixture 范围；当前 `blocked_missing_exact_p12_owner_authorization`，未打开前不得开始跨仓执行，也不得通过 P12-P18 分片级 execution readiness gate、证据索引完成门、跨仓队列首跑门、readiness、go/no-go 或 Owner 签收。
- `docs/commercial-execution-readiness-guard-coverage-verifier.json`：商用执行就绪守卫覆盖 verifier 候选，核对 P12-P18 每个 execution readiness package、商用跨仓队列镜像 gate 和商用总结文档完成门已消费授权尝试覆盖 verifier、下一授权启动守卫和 Pack Studio issue 台账关闭门；当前 `blocked_pending_owner_authorization_and_execution_evidence`，只证明覆盖关系，不代表已获授权或已执行。
- `docs/pack-studio-goal-completion-evidence-map.json`：活跃目标完成证据地图候选，汇总前端、后端、lifecycle、三候选 GUI、Code Assistant、provider、cloud、禁入动作、Owner/Base Gate、go/no-go、商用后验收门覆盖 verifier、商用禁入动作覆盖 verifier、商用执行就绪守卫覆盖 verifier、授权尝试覆盖 verifier 和下一授权启动守卫；当前 `not_achieved_requires_p12_p18_cross_repo_evidence`，不能据此标记目标完成。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-commercial-go-live-evidence-package-20260704.md`：商用最终证据包模板，要求 P12-P18 证据、后验收门、证据写回门、授权尝试覆盖、下一授权启动守卫、执行就绪覆盖 verifier 和 Pack Studio issue 台账关闭门全部满足；当前 `template_ready_not_commercial_ready`。
- `docs/commercial-go-live-evidence-package.json`：商用最终证据包机器摘要候选，镜像人工模板的 P12-P18 分片证据矩阵、硬门来源、写回门、Owner/Base Receipt 门和非充分证据清单；当前 `blocked_not_commercial_ready`，不能请求 Owner 商用签字。
- `docs/commercial-machine-go-live-evidence-package-coverage-verifier.json`：商用最终机器证据包覆盖 verifier 候选，核对机器证据包已被 readiness、go/no-go、证据索引、当前状态审计和目标完成图消费；当前 `blocked_pending_machine_go_live_evidence_package_verified`，不代表已获授权或已写回证据。
- `docs/commercial-independent-acceptance-signoff-matrix.json`：商用独立验收签收矩阵候选，要求 P12-P18、GUI/API/Receipt、issue 关闭证据、禁入动作、证据写回、Owner/Base Gate 和最终机器证据包逐行绑定独立 reviewer、artifact refs 和权威证据；当前 `blocked_pending_independent_acceptance_evidence`，不能请求 Owner 商用签字。
- `docs/pack-studio-goal-completion-evidence-map.json`：目标完成证据地图候选，已明确拒绝用覆盖 verifier、机器证据包摘要、go/no-go 候选或前后端验收契约替代真实 runtime receipts、P12-P18 分片回执、前后端回执和 Owner/Base Gate 回执。
- 前后端商用验收已新增 `gui_api_receipt_traceability_matrix`，要求候选 bundle、delivery 状态、readiness issue、promote/confirm 控件、商用 verifier 面板和三候选 GUI walkthrough 都能对上后端 endpoint / ReadModel 与 Receipt / Candidate 引用；仅有截图、仅有后端测试或仅有网络摘要都不能计入前后端验收，且商用 readiness、go/no-go、最终机器证据包和目标完成图均已消费 `gui_api_receipt_traceability_verified`。
- 独立验收签收矩阵已被商用 readiness、go/no-go、最终机器证据包和目标完成图消费；自测通过、绿色测试、人工摘要或 Owner 授权本身都不能替代 `independent_acceptance_signoff_matrix_passed`。
- `/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-authorization-roadmap-20260704.md`：P12-P18 授权路线图和授权短语索引。
- `docs/patch-candidate-intake-ledger.md`：P7 PatchCandidate 结果承接台账；记录 run 结果摘要、11 回执、不可自动应用和需复核边界。
- `docs/patch-candidate-review-ledger.md`：P9 PatchCandidate 复核候选台账；记录 reviewer、复核证据、发现项和 apply Gate pending 边界。
- `docs/oss-evidence-matrix.md`：本轮 GitHub 样本的证据矩阵与风险口径。
- `_source-materials/`：Owner 本地投放原始资料入口；原始资料不进 Git。
