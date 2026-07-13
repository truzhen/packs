# 智能家居服务商 Pack 测试计划 v6（R3 剩余覆盖收口版）

日期：2026-07-13
状态：`契约已定 / 待执行`
定位：当前客户主线。v5 已真实 GUI 通过 Pack 安装、项目/Run、六类 Candidate、EmergencyStop 四方一致和急停阻断；本计划只收口未覆盖的 R3，不重复冒充整轮已验收。

## 1. 派活卡

| 维度 | v6 裁定 |
|---|---|
| 真实证据 | 顶部查看急停曾误创建 enable candidate，现已改为只读；`emergency_stop_disabled` 已中文映射；Base/两个 ReadModel/GUI enabled、disabled、unavailable 已真实复验。 |
| 最小可交付 | 在新鲜 R3 运行态复核急停烟测后，完成伪 ref、白名单发送、双 Tab/多视口、Dendrite、历史回放、Provider/PDF/Frappe 对抗、增量统计和 Owner UAT。 |
| 真相源 | Base 持急停/Gate；06 持 Run；07 持任务；03 持 Receipt；Provider 持外部能力；client 只读投影。 |
| 仓库归属 | 现有 client/OS 修复 worktree；packs 持计划；contracts/software 只读。 |
| 风险/契约 | 对抗测试含红级动作，只允许验证 fail-closed；不改 contracts、Gate/Receipt、权限、ProviderRequirement。 |
| 禁止边界 | 不真实发送、付款、Frappe 写回或客户数据；不把 provider_missing 包装成 ready；不提交/合并/推送。 |
| 生命周期 | v5 为 `部分验收`；本计划全绿和 Owner UAT 后才可 `已验收（打包前）`。 |

新过程目录：`/Users/li/Documents/过程文档/smart-home-service-provider-v6-20260713/`。使用新 HOME/DB/profile/origin、Pack 安装、新项目/Run/Receipt；v5 证据只读。记录四仓 SHA、WIP、diff hash与服务端口。

## 2. 开跑烟测

1. 顶部紧急停止按钮只读 Base ReadModel：点击前后 Candidate/Receipt 计数不变；不得调用 `enableEmergencyStopCandidate`。
2. disabled 显示“急停未启用”，enabled 显示“急停已生效”，Base unavailable 显示不可用；不得出现“未知状态”或 preview。
3. 从安全核心显式启用后，Base、`/v3/base/readmodel/security-core`、`/v3/security-core/readmodel`、顶部与安全核心 GUI 五方一致；急停下 advice 的模型/Candidate/Receipt 增量为 0。

## 3. 剩余 R3 必测矩阵

### 3.1 主权与引用攻击

- 对 resume、候选决定、发送、Frappe 写回分别提交伪造/过期/错 run/重复使用的 decision ref、candidate ref、receipt ref；必须 fail-closed，不能推进 06、不能产生 Formal 或外部副作用。
- 非白名单接收方只能生成 CommunicationDraftCandidate；真实发送入口必须 blocked。白名单测试也只验证 Gateway readiness/候选链，除非 Owner 另行授权真实发送。

### 3.2 Provider、PDF 与 Frappe

- provider 未配置、超时、断连、恢复四态均诚实；项目/客户快照不得假 ready。
- PDF 覆盖正常、坏文件、加密、伪装、超限与重复提交；只产 Candidate/Receipt，不写 FormalEvidence。
- Frappe query/write 在未接通时为 provider_missing/blocked；不得用 fixture 声称真实写回。若本轮没有受控真实 Frappe，结论明确收窄。

### 3.3 双 Tab、多视口与历史

- 双 Tab 同项目并发推进验证 OCC/幂等；双项目验证 current run、正文、task、receipt 不串案。
- 390×844、1024×768、1440×900 覆盖项目状态、任务、Provider、急停、Gate 入口可达。
- disable/re-enable/reinstall 和 OS 重启后，历史 Receipt 可回放；停用不删历史、不允许新 Provider 使用。

### 3.4 Dendrite 与增量统计

- 运行 Dendrite readiness/外部 E2E。未就绪保留原始日志并标 `blocked/not_ready`，只收窄联邦结论。
- 每个急停/攻击 lane 记录前后模型 usage、Candidate、Receipt、FormalTask/FormalMemory、发送、执行、Frappe 写入计数；禁止只凭文案断言零副作用。

### 3.5 Owner UAT

Owner 亲手复走安装→项目→current run→六类候选→任务→Provider→Receipt→急停→历史；抽查至少三条 `transaction_ref → run_ref → candidate_ref → receipt_ref`，确认不存在跨项目泄漏。

## 4. 问题处理与放行

中小问题自行红绿修复并回到原 GUI 点继续。大问题登记污染和依赖；安全独立 lane 可继续，身份、Receipt、主权或外部副作用证据不可信时暂停。完整矩阵、增量统计、Dendrite 收窄结论和 Owner UAT 全部有新证据，且开放 P0/阻塞 P1 为 0，才允许打包；否则保持 `未验收 / 禁止打包`。
