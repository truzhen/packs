# 环保执法 Pack 测试计划 v8（停用一致性、知识卸载与跨案隔离终验版）

日期：2026-07-13
状态：`契约已定 / 修复已实现 / 待全新三轮 GUI 验收`
版本定位：当前客户主线；承接 v7 的真实登录、市场、C01/PDF/双角色证据，不沿用失败运行态判 PASS。

## 1. 派活卡

| 维度 | v8 裁定 |
|---|---|
| 真实证据 | v7 停用 Receipt 与 lifecycle 已落 disabled，但 GUI 报失败且 15 个 KnowledgeMount 仍 active；受控写入又产生 4 次真实模型调用与 6 条 Receipt，13 新 turn 未落，形成 P0 分叉。另有 C01 材料路径暂存泄漏到 C02。 |
| 最小可交付 | 证明 lifecycle 与 15 mounts 原子收口；停用后秘书所有入口零模型/零候选/零 Receipt；C01 材料表单不进入 C02；随后继续 C01-C08 安全 lane。 |
| 真相源 | 14 持 Pack lifecycle，09 持 KnowledgeMount，05 持案件/Pack 绑定，13 持秘书会话，08/03 持模型与回执，client 不持真相。 |
| 仓库归属 | OS 修 14→09 顺序、实际 mount_ref 定位和秘书 05→14 门；client 修结构化上下文、停用阻断与项目表单隔离；packs 不实现门或 Provider。 |
| 风险/契约 | 法律/模型/Receipt/急停为红；表单隔离为黄。不改 contracts/schema/DTO，不降低 Gate。 |
| 生命周期 | 修复为`已实现`；全新 GUI/DB 四路验收前保持`未验收`。 |

非目标：生产云、真实付款、真实执法/外发、客户材料、contracts、提交/合并/推送/发布。

## 2. 权限与连续推进

延续 Owner v7 的隔离全权限：四仓 worktree 与本轮专用 HOME/DB/Docker/browser profile 内可直接修改、构建、测试、启停；必须真实 GUI 登录本地 cloud、真实市场 sandbox Entitlement、真实制品下载/digest/OS 安装。中小问题自行修复并回原 GUI 点继续；大问题登记污染范围、可继续 lane、暂停 lane和恢复门，安全独立 lane 可继续就继续。

不授权生产、客户数据、真实资金、真实发送/执行、真实 secret、contracts、覆盖他人 WIP、提交、合并、推送或发布。

## 3. R1 首门：14/09/GUI 一致性

使用 `/Users/li/Documents/过程文档/env-prepack-v8-20260713/`，全新 HOME、DB、profile、cloud project、C01/C02、Run、Candidate、Receipt。

1. GUI 真实登录/市场安装并核对 2 Role Slots、15 scopes、752 knowledge；创建真实 05 C01。
2. enabled 下 15 mounts 全 active，C01 秘书请求只允许与 05/14 同源的一次模型与候选链。
3. 停用顺序硬门：09 全部目标 mount 成功软停用后，14 才可提交 disabled/pointer 清空；安装 owner 与 GUI actor 不同也必须按 packVersion+scope 找到原 mount_ref。
4. 注入任一 mount disable 失败：接口失败、14 仍 enabled、GUI 仍显示 enabled、未产生 `pack_version_disabled`；不得再出现半提交。
5. 正常停用：GUI、14、15 mounts、03 四方一致；历史 FormalKnowledge/Receipt 保留但不能进入新 ContextSlice。
6. disabled 后攻击流式、非流式、动作交接、Owner Confirm：必须在 13/08/03/07 写入前 blocked。模型 usage、Candidate、Receipt、FormalTask、发送、执行增量均为 0。
7. 无 05 项目但文本提及环保 Pack：停用态直接 blocked；启用态只给创建/选择案件引导，不调用模型或合成 transaction。

## 4. 跨案材料与法律链

- C01 输入 PDF 路径后切 C02，表单必须按 transaction key 重置/隔离；C02 DOM、截图、请求体、06/SQLite 均不得出现 C01 路径。
- C01/C02 PDF、双角色、compare Gate、法律引用和 Receipt 四路引用分别同源；切案、双 Tab、刷新、重启均不串案。
- 急停下 PDF、角色、秘书、candidate-input、Receipt、usage 零副作用；法律知识继续 `pending_human_review`，不得自动正式化。

## 5. R2/R3 与放行

R2：disable/reactivate/reinstall、不同 owner/actor、cloud/OS 重启、session 过期、15 mounts 状态与历史回放。
R3：C01-C08、坏/加密/伪装/超限 PDF、Provider 四态、双 Tab、多视口、跨事务 refs、Dendrite 收窄与 Owner UAT。

开放 P0/阻塞 P1 为 0，生命周期/知识卸载、停用零副作用、跨案材料隔离、真实登录/市场、R1-R3 与 Owner UAT 全有新证据，才可判`已验收（打包前）`；否则`未验收 / 禁止打包放行`。
