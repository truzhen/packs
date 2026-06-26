# 智能家居老板项目经营 Pack（v0）

独立文件夹、可加载可卸载的智能家居项目经营场景荚（Domain Work Pack）。单角色（智能家居项目经理）经 advice 出经营建议候选，正式动作（里程碑/采购/对外发送/Frappe 写回）一律经 Owner + Base Gate。Frappe 仅作 provider 需求声明，未接通诚实 provider_missing/blocked。

## 加载 / 卸载

```sh
# 起含本分支改动的 devserver 后：
python3 packs/smart-home-owner-pack-v0/install.py
python3 packs/smart-home-owner-pack-v0/uninstall.py
```

本 pack 无知识库（knowledge/），install.py 会自动跳过知识入库步骤。
产品基座默认不自带本 pack（server.go 已摘除自动 seed），只在 install 后出现、uninstall 后消失。
