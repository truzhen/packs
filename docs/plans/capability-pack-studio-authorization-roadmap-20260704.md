# 能力 Pack 制作台后续授权路线图

> 日期：2026-07-04
> 状态：`authorization_roadmap_ready`
> 归属：`truzhen-packs` 授权索引。本文件只列授权顺序和边界，不代表任何跨仓施工已获授权。

## 1. 推荐顺序

| 顺序 | 切片 | 为什么排这里 | 风险 | 当前授权状态 |
|---|---|---|---|---|
| 1 | P12 安全内置能力 lifecycle 样本 | 不运行第三方 OSS，不跑真实 Codex，是验证 lifecycle 底座的最小安全闭环 | 黄 / 橙 | 待授权 |
| 2 | P13 GUI lifecycle 面板 | P12 有状态后，前端才能防止用户绕过 readiness / Gate / Receipt | 黄 | 待授权 |
| 3 | P15 三候选 GUI 实操验收 | 有面板后再跑用户视角证据，避免截图证明不了 lifecycle | 黄 | 待授权 |
| 4 | P16 受控 Code Assistant 最小 run | 进入真实执行链，必须晚于 GUI 和 Gate 基础 | 橙 | 待单独授权 |
| 5 | P17 provider / adapter candidate | 有 PatchCandidate 后再落 provider 归属，避免把 provider 塞进 packs | 橙 / 红 | 待单独授权 |
| 6 | P18 云市场 sandbox | 最后接商品 / License / Entitlement / 下载链，避免未成熟资产上架 | 橙 / 红 | 待单独授权 |

## 2. 可直接回复的授权语

### P12

```text
授权按 P12 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒；只使用基座内置安全样本或 fixture 验证 lifecycle draft/readiness/promote/confirm。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p12-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p12-cross-repo-authorization-scope.json`

### P13

```text
授权按 P13 授权卡修改和测试 truzhen-client-web-desktop；必要时只为 P13 只读聚合字段修改和测试 truzhenos；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p13-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p13-cross-repo-authorization-scope.json`

### P15

```text
授权按 P15 授权卡测试 truzhen-client-web-desktop 与 truzhenos，并在 truzhen-packs 写入 GUI 实操证据；本批不改 contracts、software、cloud，不运行真实 Codex CLI，不执行第三方 OSS，不登录或上传社媒。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p15-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p15-cross-repo-authorization-scope.json`

### P16

```text
授权按 P16 授权卡修改和测试 truzhenos 与 truzhen-client-web-desktop，并允许通过 11 Gateway 执行一次受控 Code Assistant 最小 run；本批不改 contracts、software、cloud，不执行第三方 OSS，不登录或上传社媒，不应用补丁。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p16-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p16-cross-repo-authorization-scope.json`

### P17

```text
授权按 P17 授权卡修改和测试 truzhen-software；必要时只为 provider readiness 展示修改和测试 truzhenos 与 truzhen-client-web-desktop；本批不改 contracts、cloud，不执行第三方 OSS，不登录或上传社媒，不保存 raw secret。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p17-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p17-cross-repo-authorization-scope.json`

### P18

```text
授权按 P18 授权卡修改和测试 truzhen-cloud、truzhen-client-web-desktop 与 truzhenos，并在 truzhen-packs 写入云市场 sandbox 证据；本批不改 contracts、software，不真实支付、不生产发布、不执行第三方 OSS、不登录或上传社媒。
```

授权卡：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/capability-pack-studio-p18-cross-repo-execution-authorization-20260704.md`

机器可验授权范围：

`/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/capability-pack-candidates/short-video-ops-v0/docs/p18-cross-repo-authorization-scope.json`

## 3. 不建议跳过的依赖

- 不建议先做 P16：没有 P12/P13 时，真实 run 产物无法证明能安全进入 lifecycle 和 GUI。
- 不建议先做 P17：没有 P16 时，adapter candidate 容易变成手写 provider，而不是受控 PatchCandidate 后续。
- 不建议先做 P18：没有 P12-P17 时，云市场只能卖未接通候选资产，容易腐化商品真相。
- 不建议把 P15 放到 P13 前：截图能证明点击路径，但不能证明启用按钮不会绕过 Gate。

## 4. 可并行项

仅在 Owner 明确授权多仓后，以下可以并行：

| 并行组 | 条件 |
|---|---|
| P13 前端面板 + P15 evidence ledger 完善 | P15 不跑 GUI，只补 ledger 字段时可并行。 |
| P17 provider scaffold + P18 cloud sandbox 设计 | 只写设计 / fixture，不执行、不发布时可并行。 |

涉及真实 run、provider readiness、云 entitlement 的实证不能并行跳过依赖。

## 5. 授权前默认边界

在未收到明确授权语前，默认：

- 不改 `truzhenos`。
- 不改 `truzhen-client-web-desktop`。
- 不改 `truzhen-software`。
- 不改 `truzhen-cloud`。
- 不改 `truzhen-contracts`。
- 不运行 Codex CLI。
- 不执行第三方 OSS。
- 不登录或上传社媒。
- 不生成生产 License / Entitlement。
