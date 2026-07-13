# 环保执法 Pack 打包前真实使用测试战役计划 v4（三轮自主修复与条件续跑版）

日期：2026-07-12
状态：`契约已定 / 待执行`
版本定位：当前客户主线、打包前放行门
生命周期上限：本计划执行前为 `已接线`；R1-R3 与 Owner UAT 全绿后才可标记 `已验收（打包前）`。

## 1. 真实证据与本版目标

权威输入：

- R1 final 报告：`/Users/li/Documents/过程文档/env-prepack-v3-20260711/R1-final-20260712/R1-round-report.md`
- R1 final2 GUI 记录：`/Users/li/Documents/过程文档/env-prepack-v3-20260711/R1/lanes/R1-final2-GUI-retest-20260712.md`
- 大问题台账：`/Users/li/Documents/过程文档/env-prepack-v3-20260711/issues/big-issue-register.md`
- v3 用例权威：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs/docs/plans/environmental-enforcement-prepackage-real-use-test-campaign-v3-20260711.md`

已确认通过：真实 GUI 市场安装、registry 0→1、2 角色槽、15 KnowledgeMount、752 条知识、C01 PDF Candidate、compare Gate、EmergencyStop 下模型/Candidate/Receipt 零增量。新增真实风险是 `ENV-V3-P0-004`：C02 曾显示 C01 的角色正文与引用；另有 `ENV-V3-P1-002` 冷启动后手工候选回灌的 GUI 专项覆盖缺口。

本版最小可交付：先以全新环境验证跨事务隔离与 P1 专项，再完成 R1；R1 通过后继续 R2 恢复/竞态和 R3 对抗/UAT。砍掉项：生产支付、真实执法、真实外发、客户数据、contracts 变更、完整治理控制台，均不进入本轮。

## 2. 派活卡与真相源

| 维度 | 裁定 |
|---|---|
| 我要做的事 | 用真实 GUI 连续跑 R1-R3；遇到中小问题由执行者在隔离分支修复并推进；大问题登记并判断能否安全续跑。 |
| 客户/场景证据 | 同一 GUI profile 新建 C02 后出现 C01 正文与 transaction/run 引用；冷启动分支缺真实 GUI 终轮证据。 |
| 真相源 | 05 持事务对象；06 持 Run；03 持 Receipt；Base 持 EmergencyStop/Gate；cloud 持商品与 entitlement；client 只做事务级投影；packs 持声明。 |
| 仓库归属 | client：事务切换与投影隔离；OS：05/06/03、急停与回灌；cloud：本地市场链；packs：Pack 与测试计划。contracts/software 只读。 |
| 风险颜色 | 绿：文档/只读；黄：UI/局部业务；橙：Gateway/跨仓边界；红：权限、Gate、Receipt、急停、真实执行、跨事务泄漏。 |
| 契约影响 | 本轮不改 contracts/schema/DTO，不改变 Gate、Receipt、Candidate/Formal 或主权语义。 |
| 上下文 | 只读本计划列出的报告、台账与四个修复 worktree；运行态必须是本轮独立 HOME/DB/profile/project。 |
| 禁止边界 | 不触生产、不用客户数据、不真实支付/处罚/发送/执行，不写 main，不提交/合并/推送，不用 API/DB 伪造 GUI PASS。 |
| 用户验收 | Owner 从 GUI 复走市场→安装→C01/C02 隔离→PDF→双角色→compare→Receipt→急停→生命周期回放。 |

## 3. 固定修复基线与新过程目录

测试必须直接使用以下含未提交修复的隔离 worktree，并在每轮记录 SHA、branch、`git status`、diff hash：

- client：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-client-web-desktop`
- OS：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhenos`
- cloud：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-cloud`
- packs：`/Users/li/Documents/truzhenv3worktree/two-packs-product-repair-20260711/truzhen-packs`

新过程目录固定为：

`/Users/li/Documents/过程文档/env-prepack-v4-20260712/`

不得继续写旧 v3 轮目录。每轮使用全新 OS SQLite、cloud runtime、Docker project、浏览器 profile、端口、C01/C02、Run、Candidate、Receipt 与截图；旧材料只读参考。

## 4. 问题分级、自主修复与续跑决策

### 4.1 小问题：自行修复，修通后继续

P2/P3 或局部展示、选择、刷新、文案、测试脚本问题，且真相源明确、可逆、不改契约/权限/Gate/Receipt/认证、不产生外部副作用时：

`登记 SF → 最短失败复现 → 红灯测试 → 根因 → 最小修改 → 定向测试 → 邻接回归 → 从失败 GUI 检查点重跑 → PASS → 继续本轮`

执行者无需等待 Owner 再授权。每项写入 `R<N>/small-fixes/SF-NN.md`，不得口头略过。

### 4.2 中等问题：在既有授权内自行修复，复验后推进

P1 或跨 client/OS 的普通接线问题，只要同时满足：根因已确认；修改限于上述四仓隔离分支；不改 contracts/schema/DTO、Gateway 信任边界、Gate/Receipt/权限/认证；不存在客户数据或不可逆状态；有自动红灯和可重复 GUI 检查点。执行循环与小问题相同，但另写 `R<N>/medium-fixes/MF-NN.md`，包含影响清单、兼容性和回滚点。

中等问题修复后，定向测试、邻接回归和原 GUI lane 均 PASS 即继续，不因曾出现 P1 自动终止整轮。连续两次最小修复仍失败、根因扩散或需要改契约时，立即升级为大问题。

### 4.3 大问题：必须登记；能安全继续则继续，不能则暂停

以下任一项为大问题：P0；跨事务/跨用户数据暴露；权限、认证、急停、Gateway、Gate、Receipt、Candidate/Formal 穿透；生产或真实外部副作用；真相源迁移；contracts/schema/DTO；数据迁移/删除；根因不明；超出四仓授权。

立即写入 `issues/big-issue-register.md` 并保存现场。然后执行“可继续性判定”：

- **继续跑**：问题只污染单一 lane；未污染共享 DB、浏览器身份、Receipt 计数或证据归属；可完全冻结危险入口；后续 lane 与问题无依赖且不会读取受污染数据。此时大问题保持开放，该轮不得判总 PASS，但继续完成独立安全 lane 取证。
- **暂停当前轮**：存在信息暴露、主权/权限/急停穿透、正式副作用、证据归属不可信、共享运行态污染；后续 lane 依赖故障链；无法物理冻结危险入口；或继续操作会破坏现场。暂停时停止危险动作、导出只读证据、清理本轮服务。
- **禁止事项**：不得刷新/清缓存掩盖问题，不得以 API/DB 写入绕过，不得降低断言，不得用后续 PASS 稀释开放的 P0/P1。

每条大问题记录必须包含：影响面、已发生/未发生副作用、污染范围、可继续性结论、允许继续的 lane 白名单、必须暂停的 lane、恢复门槛。

## 5. R1：新鲜安装、跨事务隔离与 P1 专项

### 5.1 开跑硬门

1. cloud/node/market/审批桥 health 全绿；paid sandbox 商品、模拟 entitlement、真实 Pack artifact digest 一致。
2. OS 未登录市场 401，GUI 登录后 200；禁止浏览器参数选择生产上游。
3. registry、Role Pack、KnowledgeMount 与本轮业务数据均为 0；仅 GUI 点击安装后允许 0→1。
4. 对账 2 角色槽、15 mount、752 knowledge、lifecycle Receipt；所有法律知识保持 `pending_human_review`。

### 5.2 P0-004 强制首测

1. GUI 创建 C01，生成执法精英与挑剔律师候选并记录正文、transaction/run/Candidate/Receipt refs。
2. 保持同一登录态与浏览器 profile，新建或切换到 C02；C02 未点击生成前，DOM、可见文本、技术详情、无障碍树均不得包含任何 C01 正文或 C01 transaction/run/Candidate/Receipt ref。
3. 在 C01 请求仍 pending 时切换到 C02，待 C01 迟到响应返回；C02 仍不得出现 C01 数据。
4. C02 独立生成候选；请求体、响应、DOM、ReadModel、Receipt 必须全部只属于 C02。再切回 C01，C01 数据仍正确，不被 C02 覆盖。
5. 双 Tab 同时打开 C01/C02 重跑一次。任一跨事务可见即 P0，按 4.3 判断；由于证据归属被污染，默认暂停当前轮。

### 5.3 P1 冷启动后手工回灌专项

不得停止或改动共享真实模型。为本轮 OS 配置独立、可审计的本地 provider fault proxy：初始仅对本轮 OS 返回 `not_ready/timeout`，不伪造候选；自动协作节点诚实产生 `model_not_ready` 后，将同一 proxy 切到透传真实模型。全程由 GUI 操作，记录 proxy 状态变化与模型 usage。

随后在同一 C02 GUI 手工生成执法精英和挑剔律师真实 `live=true` 候选，必须证明：

- Candidate 与真实 Receipt 归属 C02；
- `model_not_ready`/pending 协作节点被受控恢复；
- 同槽多 advice 节点一致，重复点击幂等；
- compare Gate 可达但保持 waiting；Formal、发送、执行均为 0；
- 不覆盖已有 live 节点，不影响 C01 或其他模型使用者。

无法建立独立 fault proxy 时，不得干扰共享模型；将本专项记为“环境覆盖阻断”。若其余 R1 lane 独立，可继续取证，但 R1 不判总 PASS、不进入依赖该门槛的 R2。

### 5.4 R1 其余门槛

完成 PDF Candidate/replay、知识引用、compare Gate、EmergencyStop 下模型/PDF/Candidate/Receipt/usage 零增量，以及 Candidate/Formal 隔离。P0-004、P1 专项和既有核心链均 PASS 才允许进入 R2。

## 6. R2：恢复、幂等、生命周期与竞态

使用新事务，不复用 R1。覆盖刷新、浏览器重开、OS/cloud 重启、断网重试、双击、双 Tab、disable/re-enable/reinstall、OCC、Receipt replay、版本固定和历史回放。每次切换事务均抽检“旧正文/旧 ref 不可见”。中小问题按 4.1/4.2 自行修复并继续；大问题按 4.3 决定局部续跑或暂停。

## 7. R3：对抗、法律时点、多视口与 Owner UAT

使用新事务。覆盖坏/加密/伪装/超限 PDF、无命中、冲突时点、伪 ref、非白名单发送、停用后访问、无效急停 body、390×844/1024×768/1440×900、双 Tab C01/C02 隔离、免费商品 Gate 观察项。Owner 最终复走：市场→安装→C01/C02 隔离→PDF→双角色→compare→Receipt→急停→停用历史回放。

## 8. 每轮证据、清理与放行

每轮必须保存 baseline、GUI 截图/读屏树、network、Candidate/Receipt/ReadModel/SQLite、模型 usage、monitor doctor、small/medium fix、big issue 与可继续性判断、端口/容器释放记录。API-only、fixture、旧 transaction/Receipt 或静态测试不能替代 GUI PASS。

放行标准：R1-R3 全部硬门有新证据；P0=0、阻塞性 P1=0；跨事务隔离、Candidate/Formal、急停、Gate、Receipt 全程成立；中小修复均完成红绿与 GUI 重跑；开放大问题为 0 或有 Owner 书面豁免、替代控制和到期日；Owner UAT 通过。否则保持 `已接线/未验收`，禁止打包。
