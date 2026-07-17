# 内容运营 Pack 第二阶段嵌入证据与剩余影响清单

日期：2026-07-16
结论：Pack、角色槽、两个 FormalSchedule、software Skill mount 与通用 Task Hands 已闭环；第二阶段最小交付达到`已验收（任务分支，未发布）`。真实运行只生成 candidate artifact，未登录、未上传、未发布。

## 1. 仓库与边界

| 仓库 | 本轮动作 | 结果 |
| --- | --- | --- |
| `truzhen-packs` | 新增 `content-operations-workbench-v0`、Pack-local Skill bundle、flow、对象、角色、能力、计划、install/uninstall 和对抗测试 | 已实现并通过隔离 lifecycle 接线 |
| `truzhen-software` | 复用既有 Codex Host/Docker/Squid ProviderResource，增加 Pack/能力引用和 Skill mount source lock | 静态 registry 仍诚实保持 Host=`installed_not_authenticated`、Docker=`sandbox_not_ready`；产品运行时由实时 os-08 binding 投影方案 A 可用性 |
| `truzhenos` | 经 Owner 单独授权，补通用 07→11 Hands、Pack bundle resolver、os-08 binding/实际模型归因、终态 Receipt 收口和测试 | 没有内容运营专用表、路由、渠道模板或业务 Schema；跨仓 contracts/Gate/Receipt schema 零修改 |
| contracts / client / cloud | 未改 | 无 schema、UI、市场或发布扩张 |

Pack 不含 CLI、Provider、sidecar、账号、Cookie、token、平台登录、上传、发布、私信或联系人抓取。software 不持运营模板或候选内容，只持 ref/hash/mount policy 与既有 Provider 供给事实。

## 2. 已接线证据

隔离 devserver 上真实执行并通过：

1. fresh install：canvas 同步 06，lifecycle draft/readiness/promote/confirm，Pack enabled@0.1.0。
2. 两个 Role Pack 启用并绑定 `content_strategist`、`fact_boundary_reviewer` 槽位。
3. 07 通过 intake → submit-review → schedules/approve 创建两个 active FormalSchedule：
   - `content_ops.direction_radar`：`30 0 * * 1-5` UTC，即 Asia/Taipei 工作日 08:30。
   - `content_ops.weekly_review`：`30 12 * * 0` UTC，即 Asia/Taipei 周日 20:30。
4. 重复 install 不新增场景、角色、槽位或计划。
5. uninstall 先经 07 + 01 + 03 暂停 Pack-owned schedules，再经 Base gated action 停用 Pack；重新 install 使用最新 lifecycle OCC reactivate，并恢复两个计划。
6. 最终 Pack lifecycle receipt `5b899f4f-c2c7-4230-afad-b2de1594acf4` 由 `/v3/receipts/{ref}` 反查为 `recorded / 14-pack-studio / pack_version_reactivated`。
7. 最终两个计划均为 `active`，下一次触发分别为 `2026-07-17T00:30:00Z` 与 `2026-07-19T12:30:00Z`。

计划触发只产生请求候选工作的任务，不等于内容已生成，更不等于平台已发布。`content_ops.content_production` 明确不设 standing schedule，必须先等 Owner 选择方向与事实边界。

## 3. 静态与跨仓验证

- Packs：`go test ./...` 通过；JSON 合法；`install.py` / `uninstall.py` 可编译；`git diff --check` 通过。
- software：`python3 scripts/check-codex-provider-registry.py` 通过；`python3 scripts/test-codex-provider-build-assets.py` 9 项通过；全部 registry/source-lock TOML 可解析。
- 配对 truzhenos 全量 EGR：`GOWORK=off bash scripts/verify.sh` 最终 `exit=0 / verify ok`；G2 契约副本检查干净，Go 常规/race、前端 208 个测试文件/1496 项测试、构建与壳层 smoke 全绿。
- 07 Scheduler：审批、触发、急停、pause/resume/cancel、重启恢复、misfire、risk ceiling、后台 loop 与非法配置测试全绿。
- 11 Hands：`TestCodexRunner_MissingPerRunProofBlocksBeforeSession`、`TestBuildTaskHandsBridgeUsesRegistryGatedReadinessWithoutProxy`、`TestExecuteTaskHands_MissingBaseDecision_BlockedNoRun` 全绿，证明缺证明时不会启动 Runner。
- Pack bundle 将模型输出 Schema 收归 Pack 并纳入完整性边界后，content hash 更新为 `10473fa4912cb75f152b9f0c2a8271e45cd8daaaa89121ce1c9e80442c8d2fff`；software lock 与最终真实 Hands Receipt 必须使用该值。第一阶段 source Skill hash 仍为 `3a06fdaecbc072447a4ce7520e56c4ba403200948f213912c4edfced7f927300`。
- `skill-creator` 的 `quick_validate.py` 因系统 Python和 Codex bundled Python 均缺 `PyYAML` 未运行；未安装依赖或吞错。相同 frontmatter、文件、hash 与边界由仓库 Go 测试覆盖。

### 3.1 Codex Provider 真实供给增量

在隔离 devserver 中只通过既有 `InstallProviderCandidate -> 01 Base Gate -> 03 Receipt -> 02 readiness` 完成：

1. Docker 三个锁定 npm tarball 通过 SHA-512 验真，Node 基础镜像 digest 匹配；材料回执 `f9f860d6-c430-4f55-be2d-19bb04c2b63a`。
2. `--pull=false --network=none` 构建 Codex Hands `0.144.1`；镜像 digest `sha256:08b687b435e7bb00bf803e2e2f099bfb05d9701ae7900f60a7fbde256aae97f5`，non-root `10001:10001`，安装回执 `c26a5313-1640-4f07-8d49-9d9445b923b1`。
3. Host `@openai/codex@0.144.1` 以专用前缀安装，SRI/版本探针通过；安装回执 `f4e300bb-4816-4633-a764-7beb98baab3a`。
4. 08 Model Gateway 的 `local-omlx / Qwen3.6-35B-A3B-4bit` ReadModel 探针成功并生成限时 binding Receipt；该临时运行证据没有写成 Pack/software 静态 ready。
5. 未读取或复制共享 ChatGPT/Codex 登录态；方案 A 只在实时 binding 存在时以 `requires_openai_auth=false` 派生 Host route，binding 缺失即 fail-closed。

## 4. 真实 Hands 验收

基座中立化后的 fresh transaction `transaction://content-ops/pack-neutral-final-20260716` 已跑通：

1. os-08 binding Receipt `7f2084cf-edb6-4507-b650-8a8f4a88a8a5`，os-02 Provider/Runtime 投影 `f61196a7-f2d4-4fe3-ad00-e5dc04296b01` / `d18b2384-47ef-4690-9a3d-b0a81fa8c977`。
2. fresh T06 decision `gated_owner_decision_1023914aa1700dd4`，confirmation Receipt `1f6a9038-59ab-412f-ac5b-a915e191d137`，且 `consumed_at` 已落账。
3. os-11 preflight `21356259-a1e0-4d5e-a503-caa1f7975933` 明确通用 `pack.candidate.generate`、Host、隔离 workroot 与动态 `gated_bridge`。
4. os-08 usage `58bfe3fe-a606-4c53-9936-dd5a7a284d62` 明确实际模型 `Qwen3.6-35B-A3B-4bit`、`real_call_succeeded`、`code_assistant_model_call`。
5. os-11 final `f310465e-bb9c-479f-a7ac-14c5dc0d85a0` 与 Hands wrapper `9a4dac52-82ee-40f8-8cbb-4ded5cdd3bee` 均成功；wrapper 反查 Pack bundle hash `10473f…d2fff`、Schema hash `e94859…b3be`、preflight/final 和四个 artifact。
6. 候选文件 hash `7dbd96…67e4`，结构为 `candidate_ready / candidate_only=true / publish_ready=false`；只使用本次允许渠道，确定性泄漏扫描零命中。模型从技术 Receipt 推演了未被证据原文支持的市场/演示主张，因此 `unsupported_product_claim` 正确保留人工复核 pending，不能进入人工发布包。

## 5. 根因修复与剩余边界

- 修复 os-08 把真实 35B 调用错记为 resolver 默认模型的问题；测试直接核对 central Receipt 原文。
- 修复 10 分钟截止后执行取消连带取消终态持久化/Receipt 的问题；终态、Hands wrapper 与 07 故障投影改用独立 5 秒有界收口 context，真实状态不被伪装成 success。
- 修复内容运营 JSON Schema 和业务字段校验误入 os-11 的架构问题：`model-output.schema.json` 现归 Pack 并纳入 bundle hash；OS 仅消费 hash 校验后的 Schema data，使用通用候选边界，静态测试防业务标识和 Schema 副本回潮。
- software 静态 registry 没有改成产品 ready；每次运行仍要求新的 os-08 binding 与 T06。Docker lane 仍为 `sandbox_not_ready`，不影响本次 Owner 裁定的方案 A Host lane。
- 模型本轮只返回 1 个方向，但证据语义不足；后续真实运营前必须挂载产品 Demo、Release、客户问题或访谈证据，不能用 Provider 安装 Receipt 代替内容事实。
- contracts、client、cloud、社媒平台均未改；没有登录、上传、发布、私信、群发或联系人抓取。生命周期为`已验收（任务分支）`，不是`已发布`。
