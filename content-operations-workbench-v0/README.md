# 内容运营工作台 Pack 0.1.1

这是一个面向 Founder / Owner 自营内容的 Domain Work Pack。它把“选题雷达 → Owner 选方向 → 候选内容包 → 事实审校 → Owner 发布前判断 → 人工发布后周复盘”声明为可安装工作台。

生命周期：0.1.1 已完成打包前验收：Pack lifecycle / Scheduler、方案 A Host Codex → 本地 08 → 35B candidate-only Hands 与 Receipt 已实跑，并已从真实 GUI 生成完整 45 秒抖音候选。代码发布不等于产品已安装、启用或上架；平台发布仍只能由 Owner 人工完成。

## 首个真实场景

Owner 已有一个约 13 万粉、以别墅装修效果和灯光为主的抖音账号，每天可投入约 2 小时，希望自己只负责方向、事实和发布判断，把搜集、选题、改写、整理和复盘候选交给 Truzhen。

## 三种工作模式

1. 选题雷达：读取允许范围内的真实证据，最多生成 3 个方向候选；Owner 未选择时停止。
2. 内容生产：只在 Owner 选定方向后生成内部审计稿、公开候选稿和人工发布包；事实预检未通过时停止。
3. 周复盘：只消费人工录入或可核验导出的真实指标；缺指标时返回 blocked 和最小补数清单。

## 仓库边界

- 本目录持有工作模式、流程、角色、候选对象、Skill 声明和 ProviderRequirement。
- Codex CLI、运行 wrapper、Skill mount 和本机 Provider 登记归 `truzhen-software` 的既有 `codex-hands`。
- lifecycle、Scheduler、FormalTask、Owner/Base Gate、Execution Gateway、Provider readiness 和 03 Receipt 归 `truzhenos`。
- 本 Pack 不保存社媒账号、Cookie、token、客户隐私、登录态或平台发布事实。

## 发布边界

本 Pack 没有“自动发布”能力。它只能生成 `ContentPublicationPackageCandidate` 和人工发布检查项；登录、上传、发布、评论、私信、群发、联系人抓取均被禁止。发布事实只能由 Owner 人工录入或未来经独立 Communication / Execution Gateway 能力接入后产生。

## 结构

- `manifest.json`：工作台、Gate、Provider 和软件需求声明。
- `flows/content-operations.flow.json`：三模式候选流。
- `objects/content-candidates.json`：三种候选对象与发布状态模型。
- `role-slots/`、`role-packs/`：内容策略与事实审校 Proposer。
- `capabilities/capabilities.json`：Codex Hands 候选生成能力需求；无发布能力。
- `skills/truzhen-content-ops/`：由 Codex Hands 解释的声明式 Skill bundle，不含 CLI 或凭据。
- `docs/release-notes-0.1.1.md`：本版内容、封装方式和非目标。
- `install.py` / `uninstall.py`：只调用 Truzhen 现有 lifecycle；不直连 Provider。

## 封装

在 `truzhen-packs` 根目录执行：

```bash
python3 -B build_pack_bundle.py content-operations-workbench-v0 <仓外输出目录>
```

生成的 `.bundle.zip` 包含父层 `pack_diagnostics.py` 和完整 Pack 目录，可对已运行的 Truzhen 基座执行 `install.py`。制品与 sidecar manifest 属于构建产物，不进入 Git。

## 验收口径

- Pack 能安装、启用、停用和同版本重载，并留下可反查 Receipt。
- 缺 Codex Provider 或每次运行证明时必须 `provider_missing/not_ready/blocked`。
- 所有生成物保持 candidate-only；`receipt_ref` 只能由 03 返回。
- 公开候选不得包含内部路径、`local_candidate://` 或本机工作区信息。
- 任何测试都不得登录、上传或发布到社媒。
