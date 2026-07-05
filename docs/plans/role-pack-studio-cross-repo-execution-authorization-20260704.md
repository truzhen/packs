# 角色制作台全链路跨仓执行授权清单

> 日期：2026-07-04
> 来源计划：`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/role-pack-studio-team-office-test-plan-20260704.md`
> 当前状态：待 Owner 明确授权后执行跨仓读取 / 修改 / 测试。
> 本文件不是授权本身；它是执行前影响清单。

## 1. 本次真实目标

让角色制作台的产出从“本仓候选资产”走到商品阶段闭环：

1. 用户视角智能体只通过 GUI 制作秘书长 + 五顾问角色。
2. 角色可被能力 Pack 引用。
3. 角色可在团队设置 tab 替换。
4. 秘书长可选音色和 VRM，且只保存 asset ref。
5. 角色包可上传云服务器、sandbox 支付购买、下载、安装。
6. 安装后角色可用于团队办公室，正式绑定仍经 Owner + Base Gate。

## 2. 当前已完成的本仓部分

| 项 | 状态 | 证据 |
|---|---|---|
| 团队办公室 Role Pack 候选资产 | 已建立 | `role-pack-candidates/team-office-v0/` |
| 6 个角色包候选 | 已建立 | `role-pack-candidates/team-office-v0/role-packs/*.json` |
| 团队角色槽 | 已建立 | `role-pack-candidates/team-office-v0/role-slots/team-office-role-slots.json` |
| 能力 Pack 角色引用样例 | 已建立 | `role-pack-candidates/team-office-v0/capability-role-requirements/sample-team-research.role-requirements.json` |
| 秘书长音色 / VRM asset ref | 已建立 | `role-pack-candidates/team-office-v0/appearance/secretary-appearance-preferences.json` |
| 云商品草稿候选 | 已建立 | `role-pack-candidates/team-office-v0/commerce/cloud-listing-candidate.json` |
| 上传 / 下载 / 安装 artifact manifest | 已建立，含文件级 SHA-256 | `role-pack-candidates/team-office-v0/commerce/artifact-manifest.json` |
| sandbox 商品化流程候选 | 已建立，含上传、购买、entitlement、下载、安装阶段和阻断负例 | `role-pack-candidates/team-office-v0/commerce/sandbox-commerce-flow-candidates.json` |
| 生产发布晋级门候选 | 已建立，仍阻断真实生产发布 / 真实支付 / 生产签名下载 | `role-pack-candidates/team-office-v0/commerce/commercial-production-promotion-gate-candidate.json` |
| 本仓结构测试 | 已建立 | `role_pack_candidates_test.go` |

当前仍不能声称 GUI、后端、云上传、支付、下载或安装已通过。

## 3. 需要 Owner 授权的仓库

| 仓库 | 职责 | 需要动作 | 目标文件 / 模块范围 | 验证命令 | 禁止边界 |
|---|---|---|---|---|---|
| `/Users/li/Documents/truzhen-contracts` | 跨仓 DTO / schema 真相源 | 只读影响审计；若缺契约再申请修改 | RolePackCandidate、CapabilityRoleRequirement、TeamRoleSlotBinding、SecretaryAppearancePreference、PackListing / Entitlement 相关 schema | schema 测试、反向依赖检查 | 不让消费方先行固化新字段；不破坏兼容 |
| `/Users/li/Documents/truzhenos` | 角色候选、Owner Gate、Receipt、安装启用、团队绑定真相 | 读取、修改、测试 | Role Studio candidate、TeamRoleSlotBinding、Pack install、asset ref readiness、Receipt 反查 | Go 后端相关测试、隔离 devserver E2E | 不绕 Base Gate；不直连 cloud 支付；不写正式事实 |
| `/Users/li/Documents/truzhen-client-web-desktop` | 用户视角 GUI | 读取、修改、测试 | 角色制作台、能力 Pack 引用入口、团队设置 tab、音色 / VRM 选择、购买 / 下载 / 安装入口 | typecheck、相关单测、Playwright / browser smoke | 用户智能体只能走 GUI；不伪造后端成功 |
| `/Users/li/Documents/truzhen-cloud` | 云上传、商品、订单、支付、License、Entitlement、下载分发真相 | 读取、sandbox 修改、sandbox 测试 | upload draft、listing draft、sandbox payment、license / entitlement、download artifact hash | cloud 单测 / sandbox flow | 不真实扣款；不正式上架；不推生产云；不保存本仓资产真相 |
| `/Users/li/Documents/truzhen-packs` 当前 worktree | 角色包候选资产 | 已修改、继续验证 | `role-pack-candidates/team-office-v0/`、`role_pack_candidates_test.go`、治理账本 | `go test ./... -run TestTeamOffice -count=1`、JSON 合法性、禁入产物扫描 | 不把候选冒充发布；不保存订单 / License 真相 |

## 4. 执行顺序

1. 每个目标仓独立运行 `git status --short --branch`。
2. 先审计 `truzhen-contracts` 是否已有足够 schema。若缺口影响跨仓边界，先出 contracts 影响清单。
3. 接 `truzhenos`：RolePackCandidate、TeamRoleSlotBindingCandidate、Owner Gate、Receipt、安装校验。
4. 接 client：用户视角智能体只走 GUI 制作、引用、替换、购买、下载、安装。
5. 接 `truzhen-cloud`：upload draft、listing draft、sandbox purchase、entitlement、download hash。
6. 跑用户视角 GUI E2E：创建角色 -> 能力引用 -> 团队替换 -> 上传 -> sandbox 购买 -> 下载 -> 安装 -> 替换启用。
7. 独立验收后再评估生产晋级门；真实支付、正式上架、生产签名下载必须另有 Owner go/no-go 裁定。

## 5. 必须产出的证据

- 用户视角智能体操作日志和截图。
- 前端网络响应摘要。
- `RolePackCandidate` ID。
- `TeamRoleSlotBindingCandidate` ID。
- Owner Gate / Receipt ref。
- cloud upload receipt。
- sandbox payment receipt。
- entitlement ID。
- download receipt 和 artifact hash。
- install receipt 和 enabled role pack version。
- 负例 blocked 证据：未购买、hash 不一致、provider_missing、raw asset、主权越权。
- `commercial-production-promotion-gate` 逐项通过证据和 Owner go/no-go 裁定；没有该证据不得声称可生产发布或正常商品化完成。

## 6. 请求 Owner 明确授权

请 Owner 明确授权以下范围后，才能继续跨仓施工：

1. 是否允许读取、修改、测试 `/Users/li/Documents/truzhen-contracts`、`/Users/li/Documents/truzhenos`、`/Users/li/Documents/truzhen-client-web-desktop`、`/Users/li/Documents/truzhen-cloud`？
2. 是否允许在 `truzhen-cloud` 执行 sandbox 支付 / entitlement / download flow？
3. 是否仍禁止真实支付扣款、正式上架、生产云发布和推送？建议默认禁止。
4. 是否确认生产发布、真实支付和生产签名下载必须另经 `commerce/commercial-production-promotion-gate-candidate.json` 与 Owner go/no-go 裁定？
5. 是否允许启动隔离 devserver 和前端 devserver 做用户视角 GUI E2E？
6. 是否允许后续在各仓创建独立 worktree / 分支？默认使用 `codex/role-pack-studio-commercial-flow` 前缀。
