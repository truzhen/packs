# 内容运营 Pack 嵌入 Truzhen 第二阶段执行计划

日期：2026-07-16
Owner 指令：开始第二阶段嵌入 Truzhen；与基座分离，作为 Pack 存在，软件在 software 仓。
当前状态：第二阶段最小闭环`已验收（任务分支，未发布）`；方案 A 免 Codex OAuth，但逐 run T06、动态网络 Gate、Provider/Runtime 投影和 Receipt 全部保留

## 1. 派活卡

| 维度 | 本轮裁定 |
| --- | --- |
| 我要做的事 | 把第一阶段已验收的 `truzhen-content-ops 0.1.1` 迁成一个可安装、可启用、可停用、可回滚的内容运营 Domain Work Pack；Pack 只声明三种工作模式、候选对象、角色、能力、Gate 和 Surface。Codex CLI 继续作为外部 Hands，由 `truzhen-software` 既有 `codex-hands` Provider 供应。 |
| 真实客户 / 场景证据 | Owner 本人是首个 dogfood 用户：已有 13 万粉别墅装修/灯光抖音账号；每天只投入 2 小时；希望重复劳动交给 Truzhen，自己只做方向、事实和发布判断。第一阶段已有两个真实可复跑 run 和一次后台计划任务自然触发证据。 |
| 当前主线 / 优先级 | Owner 自用当前主线的最小产品验证，不是通用自媒体平台，也不是自动发布工程。 |
| 真相源 | Pack 声明归 `truzhen-packs`；Codex CLI/Skill mount/Provider 本机登记归 `truzhen-software`；Pack lifecycle、SceneFlowRun、FormalTask、Gate、Receipt、Provider readiness 归 `truzhenos`；跨仓 DTO 归 `truzhen-contracts`；Owner 的社媒账号和真实平台发布事实不进入任何仓库。 |
| 仓库 / 层归属 | 修改 `truzhen-packs` 与 `truzhen-software`；经 Owner 后续明确授权，`truzhenos` 只补通用 Task Hands、模型 binding、终态收口与测试，不写内容领域模板；`truzhen-contracts`、client 零修改。 |
| 风险颜色 | Pack manifest/flow/lifecycle 为黄；Provider registry、Skill mount 和 readiness 为橙；Base Gate、Receipt、真实 Hands、账号登录/上传/发布为红。红区本轮只使用既有主权链做受控验收，不修改其协议。 |
| 契约影响 | 目标是不改 contracts、不改 Gate/Receipt schema。Pack 复用 manifest v3、SceneFlow、FormalTask、ProviderRequirement、software_requirements 和现有 Codex Hands。若无法表达，先列兼容影响，另获 Owner 裁定。 |
| 不允许碰的边界 | 不把 Skill/内容规则硬编码进 `truzhenos`；不把 Codex CLI、wrapper 或凭据塞进 Pack；不读取或保存社媒登录态；不自动上传/发布/私信/群发；不自铸 Decision/Run/Nonce/Receipt；不把 candidate-only 写成已发布。 |
| 先输出什么 | Pack/software 影响清单、冻结版本映射、结构与 fail-closed 测试；随后才做 lifecycle 和受控 Hands E2E。 |
| 用户如何验收 | Owner 能在 Truzhen 中看到并启停内容运营工作台；Scheduler 只生成请求选题雷达 / 周复盘的低风险 FormalTask，不直接生成或发布内容；Owner 选择方向后才允许独立门控的 Hands run 生成公开候选包；缺 Provider / T06 proof 时明确阻断；每次真实 Hands run 必须可反查 03 Receipt；始终没有社媒自动发布。 |

## 2. 最小可交付

本轮必须闭环：

1. 一个正式文件夹 Pack：`content-operations-workbench-v0/`。
2. 三个工作模式：选题雷达、内容生产、周复盘。
3. 三个核心候选对象：`ContentDirectionCandidate`、`ContentPublicationPackageCandidate`、`ContentReviewCandidate`。
4. 一个 Owner 方向角色槽和一个事实审校角色槽；角色均为 Proposer。
5. 一个内容候选生成能力需求，绑定 `codex-hands` Provider family；发布动作不进入本轮能力集。
6. Pack install / disable / reactivate / uninstall 使用既有 lifecycle，保留 Gate 与 Receipt。
7. software registry 将现有 `codex-hands` 显式声明为本 Pack 的供应候选，并登记 Skill bundle/version/hash/mount 边界；不复制 Codex CLI。
8. 结构测试、registry 测试、负向边界测试，以及至少一次真实 lifecycle E2E。

砍入 backlog：社媒登录/上传/发布/评论/私信/群发；第三方视频工具真实运行；平台指标 API；云市场和支付；内容专用基座 UI/DB/Gateway；任何降低 Base Gate、Receipt 或权限地板的改造。

## 3. 目标结构

```text
truzhen-packs/content-operations-workbench-v0
  manifest / flow / objects / roles / capabilities / skill-bundle refs
                   |
                   v
truzhenos 现有 lifecycle -> SceneFlow / FormalTask -> Owner + Base Gate
                   |
                   v
truzhenos 11 Execution Gateway -> truzhen-software codex-hands
                   |
                   v
候选 artifacts + 03 Receipt（不触达社媒）
```

Pack 拥有业务声明，不拥有运行事实；software 拥有 Provider 供给事实，不拥有业务语义；`truzhenos` 只做中立解释、门控和回执。

## 4. 实施顺序

### P0：冻结输入与契约映射

- 以个人 Skill `0.1.1`、schema `1.1`、hash `3a06fdaecbc072447a4ce7520e56c4ba403200948f213912c4edfced7f927300` 为迁移输入。
- 把 `local_candidate://` 映射为 Truzhen Candidate/Artifact 引用策略，但不伪造正式 ref。
- 复用既有短视频能力候选的风险与发布阻断语义，不继承其 P12-P18 商用压测巨型文档链。

### P1：Pack 自包含实现

- 创建 manifest、flow、role slots、role packs、capabilities、对象说明、README。
- Skill 工作流和模板作为 Pack 受版本管理的声明资产；执行入口只写 Skill bundle ref，不直调 CLI。
- install/uninstall 复用仓内现有 lifecycle 端点与诊断协议。
- 静态验证 Candidate/Formal、Provider、账号发布边界。

### P2：software 供给扩展

- 优化既有 `codex-hands` registry，不新增平行 Codex Provider。
- 增加本 Pack与 content candidate skill 的使用关系、Skill bundle hash/mount policy、allowed workspace/ref-only 口径。
- readiness 仍由 `truzhenos` resolver 裁定；registry 不把本机 Codex 登录态写成 ready。

### P3：受控接线与验收

- 运行 Pack JSON/Python/forbidden artifacts/结构测试和 software TOML/registry/敏感扫描。
- 在隔离 devserver 中跑 install -> enabled -> disable -> reactivate/uninstall，核对 03 Receipt。
- 对现有通用 Task Hands 做可达性测试：缺每次 T06 proof 必须 blocked；若现有链可完整签发，则运行一个只生成候选的真实 Codex run。
- 通用 Hands 缺口经 Owner 重新授权后在 `truzhenos` 独立 worktree 修复；不改 Base/Gate/Receipt schema，不加入内容专用基座路由或表。

## 5. 上下文与禁止边界

允许读取：Pack 根治理、正式 Pack 样例、短视频能力候选；software 根治理、`codex-hands` registry/source-lock/profile/测试；`truzhenos` Pack lifecycle、07 Task、11 Hands、02 Provider resolver、03 Receipt 代码与测试；个人 Skill `0.1.1` 冻结 release/run。

禁止读取或复制：`~/.codex/auth.json`、真实 token/Cookie/社媒凭据、联系人和客户私密数据、software runtime/DB/日志/模型权重/构建产物、其它仓库无关 WIP。

## 6. 验收设计

执行方验证：

- Pack：JSON 全量解析、Python 编译、结构审计、forbidden artifacts、Pack 专项测试、`git diff --check`。
- software：TOML 全量解析、Codex registry 专项测试、敏感路径/禁品扫描、`git diff --check`。
- 产品链：隔离 devserver lifecycle smoke、readiness fail-closed、Receipt 反查；仅在既有链完整时跑 candidate-only Codex Hands。

独立验收主体重新确认：Pack 不含 Provider/runtime/凭据；software 不含内容运营领域规则；`truzhenos` 无内容运营专名或专用分支；缺 Provider、T06 proof、Receipt ledger 任一条件均不执行；没有登录、上传或发布事实。

## 7. 变更影响与生命周期

- Pack：新增一个 Domain Work Pack，更新 README、MODULES、FEATURE_LEDGER 和测试。
- software：更新现有 Codex registry/source-lock/profile/ledger/测试，不新增 runtime 数据。
- contracts、client、cloud：零变更；`truzhenos` 仅通用中立接线与收口修复。
- 监控：沿用 `truzhen-monitor` 的 execution/task/receipt 事件，不建平行日志格式。
- 计划与声明完成：`设计中 -> 契约已定`；仓内实现通过：`已实现`；lifecycle/readiness 投影通过：`已接线`；独立复核和真实 candidate-only Hands receipt 通过：`已验收`；本轮禁止标记`已发布`。

最终结果：fresh T06 的 `content_ops.direction_radar` 已经由 Pack bundle `be13c3b9…48d1`、Host Codex `0.144.1`、本地 os-08 `Qwen3.6-35B-A3B-4bit` 和动态 `gated_bridge` 生成候选；os-08 实际模型 usage、os-11 final 与 Task Hands wrapper Receipt 均可反查。证据见 `docs/plans/content-ops-pack-phase2-embedding-evidence-20260716.md`。
