# CLAUDE.md — truzhen-packs

本仓是 Truzhen 六仓协同架构的**包层**（`github.com/truzhen/packs`，开放）。本文件供 Claude Code / Agent 在 `/Users/li/Documents/truzhen-packs` 工作时快速加载。**完整纪律以 `AGENTS.md` 为准**；本文件只做提炼和导航。

> 文档语言：中文为主；英文只用于专有名词、命令、文件名、路径、协议名、API 字段和代码标识。

## 1. 仓库定位

`truzhen-packs` 是 Truzhen 的**包层 / Pack 资产仓**，开放仓，面向 `truzhen-contracts` 契约，不 import 基座内部实现。包把一个行业的高手经验编译成受主权闸门约束的领域治理资产；官方云市场的商品、支付、License / Entitlement 与下载分发真相归 `truzhen-cloud`。

本仓承载：

- 场景荚 / `Domain Work Pack` / `Scene Pack`
- 角色荚 / `Role Pack`
- 能力荚引用 / `Capability Pack` 的 ProviderRequirement
- 结构化知识与 Pack 装入 / 卸载脚本
- 作者端模板

本仓不承载：

- Base Gate、Receipt、Gateway、Memory、Model、Communication、Execution 的基座实现
- Provider / sidecar / MCP Server / 外部软件执行体
- 前端组件、AppShell 或真实 UI 渲染
- raw secret、运行态数据库、日志、构建产物

## 2. 三条铁律

1. **Pack 不持 Base 主权**：Pack 只声明判人、判事、门控、ProviderRequirement、证据和回执口径；正式裁定权在 Owner + Base Gate。
2. **AI 永远是 Proposer**：Role Pack、Agent、模型只产候选、草稿、质询和风险提示；不得声称已批准、已执行、已发送或已写正式记忆。
3. **真实动作必须穿 Gateway + Receipt**：真实模型、真实记忆、真实沟通、真实执行、真实市场动作归基座受控链路；未接通只能 `blocked / provider_missing / not_ready`，不得假成功。

附则——**概念先于数据和动作**：行业术语、口径、关系、判断约束必须来自真实客户 / 行业作者，并显式落入既有对象、知识（`glossary / sop / checklist / index`）、流程、角色或能力声明；Agent 不得临场猜口径。领域语义不是第四种 Pack，也不新增 Pack 内运行时。细则见 `AGENTS.md` §5.1。

## 3. 三类 Pack 封顶

- `Domain Work Pack / Scene Pack`：垂直职业工作台定义，组织对象、流程、角色、能力、Surface 和治理策略。
- `Capability Pack`：能力契约和风险声明；本仓只声明需求，不实现 provider。
- `Role Pack`：行业角色人格、口吻、决策习惯、模型策略和权限边界。

不得发明第四种 Pack。新需求先判断是现有 Pack 的 mode、template、slot、provider requirement，还是应归基座 / contracts / provider 仓。

## 4. 开工必读

每个新任务先读：

1. `AGENTS.md`
2. `README.md`
3. `MODULES.md`
4. 当前 pack 的 `README.md`、`manifest.json`、flow、role-slots、role-packs、capabilities、install / uninstall 脚本

需要参考基座范式时，只读 `/Users/li/Documents/truzhenos` 根治理文件（旧 `truzhenv3` 已封棺冻结）；需要修改或测试基座仓时，必须先获 Owner 重新授权。

## 4.1 固定主仓与 main 基准

六仓固定主仓目录为 `/Users/li/Documents/truzhen-client-web-desktop`、`/Users/li/Documents/truzhen-cloud`、`/Users/li/Documents/truzhen-software`、`/Users/li/Documents/truzhen-contracts`、`/Users/li/Documents/truzhenos`、`/Users/li/Documents/truzhen-packs`。拉新分支、合并后同步和主线核查都以这些目录为准。

若 `git worktree list --porcelain` 显示某仓 `refs/heads/main` 检出在旁路 worktree，或固定主仓不在 `main` / 落后 `origin/main`，先停工记录，不得从旧 feature worktree 继续派生。Owner 确认合并后，必须逐仓把 `origin/main` 和固定主仓同步到同一提交。

## 5. 当前包导航

见 `MODULES.md`。当前主要资产：

- `environmental-enforcement-pack-v0/`：环保执法证据链 Pack，知识型完整包。
- `smart-home-owner-pack-v0/`：智能家居老板项目经营 Pack，无知识库。
- `housekeeping-ops-pack-v0/`：家政客户服务全生命周期 Pack，可装入 / 卸载，`knowledge/` 待补。
- `templates/`：作者端模板，不参与 enabled pack 分发。

## 6. 验证命令

```sh
git status --short --branch
GOWORK=off go test ./...            # 含 forbidden-artifact 扫描 + Y4 错误码细分 guard（pack_error_code_taxonomy_test）
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
python3 -m unittest test_pack_diagnostics    # Y4 错误码行为级验证（连通性→002、readiness→004）
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

端到端装入必须有隔离基座 devserver：

```sh
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 <pack>/install.py
```

没有 devserver 铁证时，不得声称 E2E 通过。

## 7. P6 防回潮纪律

- Pack 仓只保存行业工作台声明、流程、角色、知识、能力引用、manifest、示例和可复核验收材料；不得保存业务数据、正式 `decision_ref` / `receipt_ref` / `pack_version_ref`、License / Entitlement 真相或内部运营计划。
- Pack 只能声明需求和引用能力，不持主权、不持凭据、不保存客户生产数据、不执行真实动作；真实可用性由 truzhenos / cloud / software 的受控链路证明。
- 提交前必须运行 forbidden artifact / business data 静态扫描，确认构建产物、数据库、日志、密钥和本地临时文件未进入 Git。
