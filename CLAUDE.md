# CLAUDE.md — truzhen-packs

本仓是 Truzhen 三仓架构的**包层**（`github.com/truzhen/packs`，开放）。供 Claude Code / Agent 在本仓工作时加载。**完整纪律见 `AGENTS.md`。**

> 文档纪律：中文为主，英文只用于专有名词、命令、文件名、路径、代码标识。

## 这是什么

可独立加载 / 卸载 / 分发的领域工作包（场景荚）。每个包是自包含文件夹（manifest + flow + role-packs + knowledge + install.py），**面向 `truzhen-contracts` 契约**，物理上 import 不到基座内部。包把一个行业的高手经验编译成受主权闸门约束的领域治理资产，由作者在本行业自分销。

## 铁律（改任何文件前先记住）

1. **包不持 Base 主权**：包只**声明**领域治理策略需求（判人 / 判事 / 门控流程 / Provider 绑定 / 证据 / 回执），正式裁定权在 Owner + Base Gate。AI 永远是 Proposer。
2. **真凭据绝不进仓**：API key / token / 口令 / terminal_sn / 激活码绝不进本开放仓。
3. **原始资料不进 Git**：`_source-materials/` 只留 `.gitignore` + `README` 占位；结构化内容进 `knowledge/`。
4. **三类 Pack 封顶**：场景荚 / 能力荚 / 角色荚，不发明第四种。
5. **护城河必过**：「用户能否一句话让最强模型直接做出同样结果？」能 → 不合格降级为 prompt / 模板；不能（需真实系统门控 / 领域审批合规 / 可审计证据回执 / 显式多角色对照 / 事务生命周期回放）→ 合格。

## 验证命令

```sh
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile
# 端到端装入见 AGENTS.md §4
```

## 子包导航

见 `MODULES.md`。
