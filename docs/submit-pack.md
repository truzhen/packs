# 提交 Pack

本仓接受公开 Pack 资产、模板和候选说明。提交前请先确认：你提交的是 Pack 声明和可审核材料，不是 provider 实现、用户数据、运行态数据库或真实凭据。

## 提交流程

1. Fork 本仓或在授权分支中创建目录。
2. 选择 Pack 类型：能力包、角色包或场景包。
3. 从 `templates/scene-pack-software-template/` 或已有 pack 复制最小结构。
4. 写清 README：目标客户、场景、能力、边界、风险、ProviderRequirement、门控和 Receipt。
5. 写清 manifest / JSON：不要把订单、License、支付状态、正式执行结果写成本仓事实。
6. 运行本仓检查。
7. 提交 PR，并说明目标客户、场景证据、风险说明和验证结果。

## 本地检查

```sh
git diff --check
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

## 禁止提交

- raw secret、token、cookie、private key、账号密码、terminal_sn、激活码。
- `node_modules`、`dist`、`build`、`.vite`、数据库、日志、运行态文件。
- Provider / sidecar / MCP Server 的真实运行实现。
- 真实客户资料、未授权法规原文、扫描件、合同、财务或医疗隐私材料。
- 声称 Pack 可以绕过 Owner + Base Gate、Gateway 或 Receipt 的说明。

## PR 应包含

- 目标作者 / 目标客户是谁。
- Pack 属于能力包、角色包还是场景包。
- 真实场景证据或明确标注“缺证据，待验证”。
- Provider readiness 与未接通时的 `blocked / provider_missing / not_ready` 口径。
- 文档级或结构级验证命令输出。
