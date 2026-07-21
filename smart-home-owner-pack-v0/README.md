# 智能家居老板项目经营 Pack（v0）

独立文件夹、可加载可卸载的智能家居项目经营场景荚（Domain Work Pack）。核心主链管理一个项目从商机、立项、进度、物料到交付、历史查询的完整周期。单角色（智能家居项目经理）只提出候选；Frappe 写回一律经 Owner + Base Gate、Execution Gateway 与 Receipt，未接通时诚实返回 `provider_missing`、`not_ready` 或 `blocked`。

硬件控制不是本 Pack 自建程序，也不是项目主链放行前提。只有 Home Assistant Provider 已登记、就绪且 Owner 明确请求时，Pack 才可提出设备动作候选；确认前必须展示设备实体、动作和影响，执行仍经 Base Gate、Gateway 与 Receipt。缺 Provider 时明确 `not_ready`，不得静默降级或绕过。

## 加载 / 卸载

```sh
# 起含本分支改动的 devserver 后：
python3 packs/smart-home-owner-pack-v0/install.py
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18080 \
TRUZHEN_CLIENT_URL=http://127.0.0.1:5197 \
  python3 packs/smart-home-owner-pack-v0/uninstall.py --open-gui
```

`uninstall.py` 不在后台伪造 Owner：它只打开/提示可信前台的「场景包管理」，由 Owner
显式确认停用，并只读等待 os-14 ReadModel 证明已停用。没有 Owner presence 或超时会明确
失败，不会冒充卸载成功；历史项目与 Receipt 始终保留。

本 pack 无知识库（knowledge/），install.py 会自动跳过知识入库步骤。
产品基座默认不自带本 pack（server.go 已摘除自动 seed），只在 install 后出现；Owner 前台停用后从可运行列表消失。

当前生命周期：`已实现 -> 已接线`。v1.1.0 声明、flow 与角色边界已更新，不代表真实 GUI、Frappe 或 Home Assistant 路径已经独立验收或发布。

## 离线 Pack 契约验证

```sh
python3 -m unittest discover -s smart-home-owner-pack-v0/tests -v
```

该测试只验证 Pack 的五阶段 Candidate→Gate→Gateway→Receipt 声明、Frappe 写回门控、
Home Assistant 的可选 `not_ready` 降级，以及 Frappe 16 / Baserow / Home Assistant 的软件声明；
不连接或控制任何 Provider。
