# AGENTS.md — truzhen-packs

本文件是 `/Users/li/Documents/truzhen-packs` 的 Agent 工作纪律入口。本仓是 Truzhen 六仓协同架构中的**包层 / Pack 仓**（`github.com/truzhen/packs`，开放仓），承载可独立加载、卸载、参与分发的行业工作台与治理资产；官方云市场的商品状态、审核、支付、License / Entitlement 和分发服务真相归 `truzhen-cloud`。

根原则：**Pack 只声明和编排，不持 Base 主权，不实现 Provider，不绕过 Gateway，不直接产生正式事实。**

## 0. 开工首读

每个新任务开始时，必须先读本仓：

1. `AGENTS.md`：本文件，工作纪律入口。
2. `CLAUDE.md`：速记和导航。
3. `README.md`：仓库定位、依赖方向和加载方式。
4. `MODULES.md`：当前包清单、成熟度、标准结构和边界。
5. 任务涉及的 pack 目录下 `README.md`、`manifest.json`、`flows/*.flow.json`、`role-slots/role-slots.json`、`role-packs/*.json`、`capabilities/capabilities.json`、`install.py` / `uninstall.py`。

需要核对基座治理范式、契约或运行端点时，只读参考相邻基座仓根治理文件，例如 `/Users/li/Documents/truzhenos/AGENTS.md`、`/Users/li/Documents/truzhenos/V3_GOVERNANCE.md`、`/Users/li/Documents/truzhenos/MODULES.md`（旧 `truzhenv3` 已封棺冻结，只作历史参考）。凡需修改、测试、提交或推送当前仓以外的仓库，必须先说明目标仓、原因、影响范围，并获得 Owner 重新明确授权。

## 1. 本仓定位

`truzhen-packs` 是 Pack 资产仓，不是基座、不是运行时、不是 Provider 仓、不是前端产品仓。

### 1.1 六仓协作关系

```text
truzhenos（基座私有实现；旧 truzhenv3 已冻结）
        implements
truzhen-contracts（纯接口 / 类型 / Schema）
        faces
truzhen-packs（本仓：Pack 资产，面向契约）
```

本仓只面向契约和声明式规格编写。基座通过文件夹包加载器或各 pack 的 `install.py` 经真实 lifecycle 端点装入；本仓物理上不得 import 基座内部代码。

### 1.2 本仓负责

- Domain Work Pack / Scene Pack：垂直职业工作台定义、事务流程、角色槽、能力引用、证据与回执口径、主权门控声明。
- Role Pack：行业角色的人格、口吻、决策习惯、模型策略和权限边界，绑定到 Role Slot。
- Capability Pack 引用：能力需求、ProviderRequirement、readiness 期望和风险声明。
- 结构化知识：权威资料经整理后的 `knowledge/`，以及 `knowledge-scopes.json`、`knowledge-index.json`。
- Pack lifecycle 胶水脚本：通过基座真实端点进行安装、启用、停用、角色绑定、知识入库或卸载。
- Pack 作者端模板：`templates/` 下的工程模板或说明，不直接参与运行分发。

### 1.3 本仓不负责

- 不实现 Base Gate、Receipt Ledger、Memory Gateway、Model Gateway、Communication Gateway、Execution Gateway。
- 不实现 Provider、sidecar、MCP Server、外部软件运行时、Frappe / ERP / OCR / IM / PDF / Codex 等真实执行体。
- 不保存正式业务对象、正式任务、正式记忆、正式回执或用户运行态数据库。
- 不实现前端组件、AppShell、移动端、Tauri 壳或真实 UI 渲染。
- 不保存 raw secret、真实凭据、token、账号密码、terminal_sn、激活码或生产端点密钥。
- 不保存云市场订单、支付状态、正式 `PackListing`、License / Entitlement 或下载分发运行状态（真相归 `truzhen-cloud`）。

## 2. 三类 Pack 固定边界

只允许三类 Pack，不得发明第四种。

| 类型 | 中文名 | 职责 | 本仓形态 |
|---|---|---|---|
| `Domain Work Pack` / `Scene Pack` | 场景荚 / 领域工作包 | 在某个业务情境下组织对象、流程、角色、能力、工作台 UI Surface 和治理策略 | 本仓主体，每个为独立文件夹 |
| `Capability Pack` | 能力荚 | 描述能做什么、风险、门控、回执和 ProviderRequirement | 本仓只放能力需求 / 引用，不放执行实现 |
| `Role Pack` | 角色荚 | 描述 Agent 如何以某种行业角色提出候选、质询或复核 | 可随场景包携带 `role-packs/*.json` |

Pack、Role Pack、Capability Pack 都只能生成候选、声明约束或提供上下文，不拥有审批权、执行权、发送权、正式记忆权、正式任务权或真相权。

## 3. Domain Work Pack 完整结构

一个合格的文件夹包至少应说明以下内容，允许按成熟度渐进补齐，但不得把缺口伪装为完成。

```text
<pack>-v0/
  README.md
  manifest.json
  flows/*.flow.json
  role-slots/role-slots.json
  role-packs/*.json
  capabilities/capabilities.json
  knowledge/
    knowledge-scopes.json
    knowledge-index.json
    **/*.md
  _source-materials/
    .gitignore
    README.md
  install.py
  uninstall.py
  docs/
```

### 3.1 manifest 必须声明

- `pack_ref` / `pack_id` / `version` / `template_family`。
- 判人策略：`person_strategy`，说明哪些角色只是 Proposer，高风险动作如何回 Owner + Base。
- 判事策略：`formalization_requirement`，说明哪些对象仍是 candidate，哪些动作需要正式化。
- 门控：`gates`，说明 Owner Gate、Base Gate、角色对照门、执行门、发送门等。
- Provider 需求：`provider_requirements`，只声明能力和风险，不绑定死具体实现。
- 底层软件需求：`software_requirements`，只声明 Baserow / OCR 等软件 family、版本范围、能力标签、reuse / isolation / fallback / license 策略；复用、需安装、版本冲突、需隔离和用户侧 lock 由 `truzhenos` resolver 产出。Pack 仓不得包含软件本体、镜像、模型权重、真实 DB、路径、端口、账号、secret 或 runtime state。
- 通知 / 命令候选 / 回报路由：`notification_command_report_routes`。
- 多角色对照：`multi_role_comparison`，必须显式节点编排，禁止 runtime 内隐藏 agent 回路。
- 护城河：`moat_justification`。
- 知识域：如有知识库，`knowledge_scopes` 必须与 `knowledge/knowledge-scopes.json` 保持一致。

### 3.2 完整垂直职业工作台六部分

基座范式中，Domain Work Pack 的最新定性是“一个垂直职业的工作台定义”。完整 pack 应能逐步表达：

1. 工作模式集：职业有哪些工作区 / 模式。
2. 事务流程：一笔业务如何从候选、确认、闸门、执行到回执。
3. 业务对象与领域语义：对象 schema、生命周期、状态覆盖和外部映射意图；并把行业概念、别名、口径边界、从属关系、判断 / 计算约束与例外映射到既有对象 schema、`knowledge/`、flow / gates、角色和能力引用。Pack 只定义模板和行业声明，运行实例与正式对象真相归基座 05。该内容仍属于“业务对象”部分，不新增第七部分。
4. 能力引用：能力域、能力操作、ProviderRequirement 和 readiness；不实现能力。
5. 角色引用：Role Slots、Role Packs、作用范围和模型策略；全员 Proposer。
6. 工作台 UI 声明：tab、工作区、标准视觉单元和 Surface 意图；不实现前端组件。

当前仓内既有 pack 可以从事务流程 + 角色 + 能力 + 知识起步；新增或升级 pack 时，应在 README / manifest 中诚实标注哪些部分已实现、哪些仍是 backlog。

## 4. 主权链纪律

### 4.1 AI 永远是 Proposer

Agent、模型、角色包、场景包只能提出候选、草稿、建议、质询、风险提示、能力调用候选或执行意图候选。正式裁定权只在 Owner + Base Gate。

禁止任何 pack 文案、脚本或测试声称：

- AI 已经批准。
- AI 已经正式执行。
- AI 已经正式发送。
- AI 已经写入正式记忆。
- pack 自己拥有主权裁定权。

### 4.2 Candidate 与 Formal 必须隔离

Pack 只能声明或生成：

- `TaskCandidate`
- `MemoryRequestCandidate`
- `CommunicationDraftCandidate`
- `ExecutionIntentCandidate`
- `BusinessObjectCandidate`
- `SceneFlowRunCandidate`
- `CapabilityInvocationCandidate`
- `PackCandidate`

候选成为 Formal Record 前必须经 Owner + Base Gate。Pack 内不得把 `candidate_only` 写成完成态，也不得把 `receipt_candidate://` 冒充正式 `receipt://`。

### 4.3 真实动作必须经 Gateway + Receipt

真实模型、真实记忆、真实沟通、真实执行、真实登录、真实市场动作必须由基座对应 Gateway / Auth / Market 链路处理，并产生 Evidence / Receipt。Pack 不得直连模型、文件系统执行器、浏览器、IM、ERP、CRM、OCR、PDF、MCP `tools/call` 或其它外部系统。

未接通的 provider 只能诚实返回或声明：

- `blocked`
- `provider_missing`
- `not_ready`
- `degraded`

不得用 mock 成功、disabled 成功、fixture 成功、说明文字或静态字段冒充真实接通。

## 5. 知识与资料纪律

- `_source-materials/` 是 Owner 投放权威原文的本地入口，只允许提交 `.gitignore` 和 `README.md` 占位；原始法规、标准、PDF、扫描件、商业资料、受版权限制材料不得进 Git。
- 结构化后可分发的资料进入 `knowledge/`，每条必须有来源引用、知识域、kind、verification status。
- 涉及法律、监管、合同、医疗、财务等高风险知识时，默认标 `pending_human_review`；正式适用必须以现行有效官方原文和人工核验为准。
- `knowledge_scopes[]` 是所有场景包通用平台能力，不得写成某个行业专用分支。
- Pack 启用 / 停用只改变 KnowledgeMount 可见性和运行访问权；不物理删除 FormalKnowledge，不破坏历史 Receipt。
- `manifest.json`、`knowledge/knowledge-scopes.json`、`knowledge/knowledge-index.json` 和实际 Markdown 文件必须保持一致。

### 5.1 领域语义与隐性业务知识纪律

- “领域语义资产”（行业概念、口径、关系、判断约束等隐性业务知识，俗称“暗知识”）是三类 Pack 的横向内容维度，不是第四种 Pack、独立真相源、中央语义层或 Pack 私有数据库。Domain Work Pack 是行业语义主载体；Role Pack 只携带专业观察与判断边界；Capability Pack 在本仓只体现能力含义、风险、Gate、Receipt 与 ProviderRequirement 引用，不实现能力。
- 设计顺序必须从真实业务概念和场景证据开始，再映射到对象字段、外部数据、能力、流程和 UI；不得从现有数据表、Provider 字段或某个软件界面反推并固化行业世界。
- 对会影响候选质量或正式动作的行业术语，至少说明：标准名称与别名、定义、适用范围 / 边界、对象关系或上下级关系、判断 / 计算约束、例外、来源与核验状态。无真实客户 / 作者证据时标记“缺证据 / 待人工核验”，不得由模型临场猜测后写成 Pack 规则。
- 对象、字段和关系进入对象 schema；术语、口径、SOP、清单和索引进入 `knowledge/`（优先使用已有 `glossary / sop / checklist / index` kind）；流程条件和门控进入声明式 flow / gates；专业观察角度进入 Role Pack；动作含义、风险、Gateway 和回执要求进入 capability requirement。不得新增可执行脚本、条件分支实现或平行 schema 冒充语义资产。
- `knowledge_scopes[]` / `KnowledgeMount` 是所有 Scene Pack 共用的挂载机制；scope 的内容可以且应按行业专用。禁止的是在 truzhenos 中新增某行业专用的加载、权限或运行分支，不是禁止行业知识域。
- Agent / Model 只能消费 truzhenos os-09 ContextCompiler 裁剪出的 ContextSlice。存在歧义、冲突、过期或缺来源的概念时，应请求 Owner 澄清，或输出带不确定性与引用的 Candidate；不得静默选择一种口径，更不得把推测升级为 Formal。
- 语义声明本身不授予查看、修改、调用或执行权限。有效范围仍取 Owner 授权、Role Pack 白名单、Pack 声明、KnowledgeMount / ContextSlice 与 Provider 能力的交集；正式动作仍经 Owner + Base Gate + Gateway + Receipt。
- 新增概念要有来源和复核；废弃概念保留迁移 / 别名说明；冲突口径按 owner / team / scene / pack version / effective time 隔离，不得用新 Pack 静默覆盖 FormalKnowledge 或历史 Receipt。
- 最小 Pack 可暂不补齐全部领域语义，但必须在 README / MODULES 诚实列出覆盖范围和缺口。存在行业歧义术语、统计 / 判断口径或高风险规则，却没有 glossary / 知识条目、来源和核验状态的 Pack，不得宣称“语义完整”或据此升级成熟度。

## 6. Provider / 能力纪律

- `capabilities/capabilities.json` 只声明能力需求、gateway class、risk class、fallback policy 和 optional 状态。
- Provider 是真实软件 / 外部系统 / sidecar / MCP Server / 云服务 / 本机工具 / 用户自装系统，归基座或外部 provider 仓，不归本仓。
- 一个 Provider 可实现多个能力；一个能力可由多个 Provider 实现。Pack 不得默认写死“某能力 = 某软件”。
- Frappe、ERPNext、Odoo、金蝶、Salesforce、OCR、PDF、IM、Dendrite、Codex 等只能作为 ProviderRequirement 或参考 Provider 出现；不得把全家桶包装成巨大 Capability Pack。
- Provider health / readiness 不能替代执行授权；真实调用仍需 Owner + Base Gate + Gateway + Receipt。

## 7. 模块与目录纪律

### 7.1 根目录治理文件

- `AGENTS.md`：完整工作纪律入口。
- `CLAUDE.md`：短上下文速记，必须指向 `AGENTS.md`。
- `README.md`：对外说明仓库定位、依赖方向、三类 Pack 和加载方式。
- `MODULES.md`：包清单、成熟度、标准结构、边界与验收。
- `.github/workflows/ci.yml`：开放仓基础校验。

### 7.2 Pack 目录

每个 pack 是自包含单元。新增、删除、重命名或改变 pack 主权策略，必须同步更新：

- 根 `README.md`
- 根 `MODULES.md`
- 该 pack `README.md`
- `manifest.json`
- 必要时更新测试报告 / docs。

### 7.3 templates 目录

`templates/` 是作者工程模板，不是已发布 pack。模板可以包含 README、manifest 示例、文档骨架和安全说明；不得包含真实 secret、真实 provider 实现、构建产物或运行时数据库。

## 8. 开发与变更纪律

### 8.1 开工状态

修改文件前必须运行：

```sh
git status --short --branch
```

新任务默认使用新分支和独立 worktree。当前主仓目录只做主线和只读核对，避免多个会话抢同一目录。

### 8.1.1 固定主仓与 main 基准纪律

六仓固定主仓目录如下，均视为各仓 `main` 的本机主线工作区：`/Users/li/Documents/truzhen-client-web-desktop`、`/Users/li/Documents/truzhen-cloud`、`/Users/li/Documents/truzhen-software`、`/Users/li/Documents/truzhen-contracts`、`/Users/li/Documents/truzhenos`、`/Users/li/Documents/truzhen-packs`。

- `main` 不得长期检出在 `/Users/li/Documents/truzhenv3worktree/`、`~/.config/superpowers/worktrees/` 或其它旁路 worktree。若 `git worktree list --porcelain` 显示 `refs/heads/main` 不在对应固定主仓目录，先停工，记录偏离情况并请 Owner 裁定迁回或清理。
- 拉新任务分支前，必须在对应固定主仓目录执行 `git status --short --branch`、`git fetch origin`，确认本地 `main` 与 `origin/main` 对齐；若领先、落后或分叉，必须先记录并等 Owner 裁定，不得从旧 feature worktree、历史 merge worktree 或旁路 worktree 派生新分支。
- 功能经 Owner 确认合并后，必须逐仓确保 `origin/main` 与上述固定主仓目录同步到同一提交；不得只合到旁路 worktree、只推远端、或让固定主仓停留在旧分支，避免出现“主线已完成但本机主目录看不到”的假故障。
- 多仓合并和拉分支必须逐仓执行：每仓独立 status、验证、合并、主仓同步和账本登记；不得用一个仓的成功替代另一个仓。

### 8.2 变更分级

| 风险 | 例子 | 处理方式 |
|---|---|---|
| 绿 | 文档说明、README、测试报告、只读清单 | 可直接改后自测 |
| 黄 | manifest 字段、flow、role slot、capabilities、pack lifecycle 脚本 | 改前说明影响，改后跑结构校验 |
| 橙 | ProviderRequirement 语义、contracts schema、跨仓加载协议、Gateway 对接 | 先出影响清单和兼容策略，获 Owner 明确确认 |
| 红 | Base Gate、Receipt、权限、真实执行、真实发送、数据删除、凭据处理 | 本仓不得直接施工，只能提方案并转基座 / provider 仓 |

### 8.3 禁止事项

- 禁止提交 `node_modules`、`dist`、`build`、`.vite`、数据库、日志、`*.jsonl`、`__pycache__`、`*.pyc`、raw secrets。
- 禁止把 provider 执行代码塞进本仓。
- 禁止把 pack 写成绕过基座的独立 App。
- 禁止把 ReadModel、前端 Projection、Capsule、模板或说明文字写成真相源。
- 禁止用 fallback、mock、demo、seed-only、candidate-only 冒充完成。
- 禁止自动 merge、push、打 tag、发版或上架云市场，除非 Owner 明确授权。
- 禁止回滚、覆盖、删除非本人改动；遇到他人 WIP，先保护、绕开或等待 Owner 裁定。

### 8.4 出错处理

遇到 JSON 校验失败、脚本报错、端点契约漂移、测试失败或异常，必须定位根因。禁止吞异常、假成功、调参重试掩盖问题或新增平行 fallback 冒充完成。确需防御性处理时，必须说明防什么、为什么在这里防、根因是否已修。

## 9. 验证命令

基础验证：

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

推荐结构审计：

```sh
python3 - <<'PY'
import json, pathlib, sys
bad = []
for manifest_path in sorted(pathlib.Path('.').glob('*/manifest.json')):
    if str(manifest_path).startswith('templates/'):
        continue
    m = json.load(open(manifest_path, encoding='utf-8'))
    base = manifest_path.parent
    for key in ['flow_file', 'role_slots_file', 'capabilities_file']:
        rel = m.get(key)
        if not rel or not (base / rel).exists():
            bad.append(f'{manifest_path}: {key} missing or not found: {rel}')
    required = [
        'person_strategy', 'formalization_requirement', 'gates',
        'provider_requirements', 'notification_command_report_routes',
        'multi_role_comparison', 'moat_justification'
    ]
    miss = [k for k in required if k not in m]
    if miss:
        bad.append(f'{manifest_path}: missing governance fields {miss}')
    if m.get('knowledge_scopes') and m.get('knowledge_scopes_manifest'):
        declared = set(m.get('knowledge_scopes') or [])
        scopes = json.load(open(base / m['knowledge_scopes_manifest'], encoding='utf-8')).get('scopes', [])
        actual = {s.get('scope_ref') for s in scopes}
        if declared != actual:
            bad.append(f'{manifest_path}: knowledge scopes mismatch declared-only={sorted(declared-actual)} actual-only={sorted(actual-declared)}')
if bad:
    print('\n'.join(bad))
    sys.exit(1)
print('Pack 结构审计 OK')
PY
```

端到端装入需隔离基座 devserver，且本仓与基座仓平级：

```sh
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 <pack>/install.py
```

验收铁证应包括：registry 从 0 到 1、Pack enabled version、Role Pack enabled、SlotBinding 生效、KnowledgeMount / FormalKnowledge 入库、正式动作经 Owner + Base Gate 并产生可反查 Receipt。无基座 devserver 时不得声称 E2E 通过。

## 10. 当前包状态

| 包 | 状态 | 纪律提示 |
|---|---|---|
| `environmental-enforcement-pack-v0/` | 完整文件夹包，含 install / uninstall、flow、2 角色包、capabilities、knowledge | 高风险法律知识默认 `pending_human_review`；知识域声明必须与实际 scopes 保持一致 |
| `smart-home-owner-pack-v0/` | 完整文件夹包，含 install / uninstall、单角色、Frappe ProviderRequirement，无知识库 | Frappe 只能是 provider requirement / write candidate，不是真相源 |
| `housekeeping-ops-pack-v0/` | 可装入文件夹包，含 manifest、flow、capabilities、2 角色包、install；`knowledge/`、`uninstall.py` 待补 | 不得标成完整知识包；补卸载脚本前不得声称完整可卸载闭环 |
| `templates/` | 脚手架 | 不参与分发，不作为 enabled pack |

## 11. 哪些变更必须回 Owner

- 新增、删除、重命名一个 pack。
- 改 pack 的主权策略声明：判人、判事、门控、ProviderRequirement、通知 / 命令候选 / 回报路由、多角色对照。
- 删减核心能力、降低护城河、删除知识域、改变 forbidden artifacts。
- 改 install / uninstall lifecycle 语义。
- 改依赖方向、仓库边界、跨仓加载路径或 contracts schema。
- 执行真实 devserver 装入 / 卸载、推送、打 tag、发布市场、发版。

## 12. 完成口径

声称完成必须提供可复核证据：

- 文件 diff。
- JSON / 脚本语法 / forbidden artifacts / 结构审计结果。
- 涉及运行时则提供隔离 devserver E2E 日志或测试报告。
- 涉及 UI / Surface 只声明时，明确说明“只声明，不实现前端组件”。
- 涉及 provider 未接通时，明确标记 `provider_missing / not_ready / blocked`，不得包装成已接通。

计划、声明、mock-only、demo-only、candidate-only、readmodel-only 或静态文档不能单独算 Pack 成品完成。

## 13. P6 防回潮纪律

- Pack 仓只保存行业工作台声明、流程、角色、知识、能力引用、manifest、示例和可复核验收材料；不得保存业务数据、正式 `decision_ref` / `receipt_ref` / `pack_version_ref`、License / Entitlement 真相或内部运营计划。
- Pack 只能声明需求和引用能力，不持主权、不持凭据、不保存客户生产数据、不执行真实动作；真实可用性由 truzhenos / cloud / software 的受控链路证明。
- 提交前必须运行 forbidden artifact / business data 静态扫描，确认构建产物、数据库、日志、密钥和本地临时文件未进入 Git。
