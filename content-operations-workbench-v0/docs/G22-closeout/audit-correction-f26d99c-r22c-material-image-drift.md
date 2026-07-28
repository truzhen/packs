# G22 R22C 临时离线物料路径审计更正

日期：2026-07-28。范围仅为 G22 的去敏 closeout；不改变 OS、Software 或镜像状态。

## 更正

先前把固定 Software 工作区中被忽略的 `providers/codex-hands/docker/materials` ACL 当作阻断，是错误配置下的 preflight，不是最终根因。该请求在 candidate 创建前被拒绝，隔离库候选数为 0，未产生 Gate 或 Receipt，现作废。

## 当前受控复验

- 新建 task-owned 临时 Software 副本，固定 `614e771086fa3da9bf495a0ba7b89aca50170015`；OS 以该临时根作为正式 source-lock 消费根。
- 官方 `prepare-offline-materials.py --materials <task-owned path>` 实际退出 1，未写入材料。它在本机 image inspect 阶段拒绝 `local image RepoDigest does not match source lock`。
- source-lock 锁定 `truzhen/codex-hands:0.144.1@sha256:08b687b435e7bb00bf803e2e2f099bfb05d9701ae7900f60a7fbde256aae97f5`；当前同标签本机镜像 RepoDigest / image ID 均为 `sha256:7a86b09e2bda34d17531cc5fe0aacd5527e49d1b3f0bf12ed49a1bf21bb3f926`；按 digest inspect `08b687…` 返回不存在。
- 官方 verifier 随后因材料目录不存在而拒绝；这不是可接受的 fallback，也没有手工 marker、registry、source-lock 或镜像标签改写。

结论：Provider 正向路径未达到可创建候选的阶段。无 Provider candidate、Owner/Base Gate、Receipt、Docker build/run、模型调用、Hands session、素材产物或平台动作。需由 OS/Software Owner 恢复 source-lock 所需的 exact local image 或做受治理的供给版本化裁定后，才可重新开始新的空 store 验收。
