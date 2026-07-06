# 贡献指南

欢迎开源作者、工具开发者和行业专家提交徒真（truzhen）Pack。提交前请先读：

- [三类 Pack 怎么选](docs/pack-types.md)
- [作者指南](docs/author-guide.md)
- [提交 Pack](docs/submit-pack.md)
- [审核与上架](docs/review-and-listing.md)

## 基本要求

- 文档中文为主；英文只用于专有名词、命令、文件名、路径、协议名、API 字段和代码标识。
- Pack 只声明和编排，不持 Base 主权，不实现 Provider。
- 真实发送、真实执行、正式记忆、正式任务和正式业务对象必须经 Owner + Base Gate、Gateway 和 Receipt。
- 未接通能力必须标 `blocked / provider_missing / not_ready`。
- 收益文案必须诚实，不承诺固定收入，也不暗示可以跳过审核直接上架。

## 提交前检查

```sh
git diff --check
python3 -c "import json,glob;[json.load(open(f)) for f in glob.glob('**/*.json',recursive=True)];print('JSON 合法')"
find . -name install.py -o -name uninstall.py | xargs -r python3 -m py_compile
git ls-files | rg '(^|/)(__pycache__|node_modules|dist|build|\.vite)(/|$)|\.(db|sqlite|log|jsonl|pyc)$' || true
```

不要提交 raw secret、token、cookie、private key、账号密码、运行态数据库、日志、构建产物或真实客户隐私材料。
