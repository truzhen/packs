# AGENTS.md — truzhen-packs

本仓 = Truzhen **包层**（`github.com/truzhen/packs`，开放）。智能体只读本文件即可独立维护本仓。

## 1. 本仓职责
承载可独立**加载 / 卸载 / 分发**的领域工作包（场景荚 = Domain Work Pack）。每个包是一个自包含文件夹：`manifest.json`（六件事声明）+ `flows/`（GateFlowSpec 流程图）+ `role-slots/` + `role-packs/` + `capabilities/` + `knowledge/`（领域知识，权威资料结构化）+ `install.py` / `uninstall.py`（经基座真实 lifecycle 端点装入 / 卸载）。包**面向 `truzhen-contracts` 契约**编写，物理上 import 不到基座 `truzhenos` 内部。

## 2. 禁止事项
- **禁止持有 Base 主权**：包只**声明**领域治理策略需求（判人 / 判事 / 门控流程 / Provider 绑定 / 证据要求 / 回执口径），正式裁定权在 Owner + Base Gate，不在包。AI 永远是 Proposer。
- **禁止真凭据进仓**：任何 API key / token / 口令 / terminal_sn / 激活码绝不进本开放仓。
- **禁止原始资料进 Git**：`_source-materials/` 是 Owner 投放的权威法规 / 标准原文区（体积 + 版权 + forbidden_artifacts），只留 `.gitignore` + `README` 占位；结构化后的内容进 `knowledge/`。
- **禁止 forbidden_artifacts**：`node_modules` / `dist` / `build` / `.vite` / `*.db` / `*.sqlite` / `*.log` / `*.jsonl` / `__pycache__` 不进 Git。
- **禁止发明第四种 Pack**：只有场景荚（Domain Work Pack）/ 能力荚（Capability Pack）/ 角色荚（Role Pack）三类。
- **禁止伪造护城河**：每个场景荚必过护城河测试（见 §3），不合格降级为 prompt / 模板，不冒充。

## 3. 必读文件
- `README.md`（依赖方向、加载方式、三类 Pack）
- `MODULES.md`（包清单 + 文件夹包标准结构）
- `CLAUDE.md`（铁律速记）
- 基座治理原文 `AGENTS.md` / `V3_GOVERNANCE.md`（在基座仓 truzhenos）：三类 Pack 定义、护城河测试、场景荚六件事、术语字典。

## 4. 常用验证命令
```sh
# 所有 JSON 合法性（manifest / flow / knowledge-index 等）
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
# install / uninstall 脚本语法（仅有脚本的包）
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile && echo "脚本语法 OK"
# 端到端装入（隔离基座 devserver；truzhen-packs 须与 truzhenos 平级）：
#   起 truzhenos devserver(隔离 DB, -addr 127.0.0.1:18099)
#   TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 <pack>/install.py
#   铁证：registry 0→1，知识入 09 FormalKnowledge，正式动作经 Owner+Base Gate 产 Receipt
```

## 5. 出错时先看哪里
- **install.py 报错** → 多半是基座 devserver 未起 / `TRUZHEN_DEVSERVER_BASE` 未指向 / 端点契约漂移：核对基座 lifecycle 端点。
- **JSON 校验失败** → 改坏了 manifest / flow / knowledge-index 的 JSON 结构。
- **manifest 缺字段** → 场景荚六件事（person_strategy / formalization_requirement / 门控 gates / provider_requirements / notification_command_report_routes / multi_role_comparison）+ pack 标识 + template_family 应齐全。
- **跨仓路径找不到 pack** → truzhen-packs 必须与 truzhenos **平级**放在同一父目录（基座加载器 `runtime.Caller` 上溯到三仓父目录再进 truzhen-packs）。

## 6. 哪些变更必须回 Owner
- 新增 / 删除一个 pack，或改 pack 的主权策略声明（判人 / 判事 / 门控）。
- 删 pack 核心能力 / 简化护城河 / 改 forbidden_artifacts。
- `git push` / 打 tag / 发版 / 上架云市场等外向动作（Owner 在场授权）。
- 依赖方向、仓边界相关的任何调整。
