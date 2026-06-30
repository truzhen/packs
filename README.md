# Truzhen Packs

> Truzhen 主权事务操作层的**领域工作包仓**：可独立加载 / 卸载 / 分发的场景荚（Domain Work Pack），面向契约、不含基座实现。

`github.com/truzhen/packs` 是 Truzhen 五落点架构的**包层**。每个包把一个行业的高手经验编译成受主权闸门约束的领域治理资产——**声明**判人 / 判事策略、门控流程、真实系统 Provider 绑定、证据与回执要求，由作者在本行业自分销。正式裁定权在 Owner + Base Gate，**包不持有主权**，AI 永远是 Proposer。

## 依赖方向（单向不可逆）

```
┌─────────────────────┐   implements   ┌──────────────────────┐   faces   ┌─────────────────────┐
│ truzhenos (基座·私有) │ ─────────────▶ │ truzhen-contracts     │ ◀──────── │ truzhen-packs (本仓) │
│ 实现契约             │                │ 纯接口/类型/Schema     │           │ 面向契约            │
└─────────────────────┘                └──────────────────────┘           └─────────────────────┘
```

包**面向契约**编写，物理上 import 不到基座内部；基座通过文件夹包加载器 / 各包 `install.py` 经真实 lifecycle 端点装入。`truzhenos`、`truzhen-contracts`、`truzhen-packs`、`truzhen-software` 与 client repo 默认**平级**放在 `/Users/li/Documents/`。

官方云服务、云市场商品运营状态、订单、支付回调、License / Entitlement 真相、Pack 文件分发运行面和官方云端网页归 `/Users/li/Documents/truzhen-cloud`。本仓只声明 Pack 本体、工作台、能力需求和可被商品化引用的元数据，不保存订单、License、支付状态或云端 server 实现。

## 三类 Pack（不得发明第四种）

- **场景荚（Domain Work Pack / Scene Pack）**：可交付的领域治理资产，声明六件事（判事 / 判人 / 门控流程 / Provider 绑定 / 通知-命令-回报路由 / 多角色对照），**不持 Base 主权**。本仓主体。
- **能力荚（Capability Pack）**：能力描述符 / ProviderRequirement。执行 provider 本体归 `truzhenos` 或外部 provider / `truzhen-software`，不写进本仓。
- **角色荚（Role Pack）**：智能体人格 / 口吻 / 决策习惯 / 模型策略，绑定到任意 Role Slot。本仓可随场景荚携带角色包数据。
- **商品化引用**：Pack 可以声明市场展示元数据和商品化引用，但订单、支付、License、Entitlement 和分发运行状态归 `truzhen-cloud`，不归本仓。

## 当前包状态（2026-06-29）

- `environmental-enforcement-pack-v0/`：完整文件夹包，含 install/uninstall、15 知识域、角色包。
- `smart-home-owner-pack-v0/`：完整文件夹包，含 install/uninstall、项目经理角色包、Frappe ProviderRequirement。
- `housekeeping-ops-pack-v0/`：已升级为可装入文件夹包，含 `pack_ref` / `template_family` / flow / capabilities / role-slots / 2 角色包 / `install.py`；`knowledge/` 与 `uninstall.py` 仍待补，不得标成完整知识包。

## 文件夹包标准结构

```
<pack>-v0/
  manifest.json          # 场景荚规格：pack_ref / version / template_family + 六件事 + provider_requirements + moat_justification
  flows/*.flow.json      # GateFlowSpec 门控流程图
  role-slots/            # Role Slot 声明
  role-packs/            # 绑定的角色包
  capabilities/          # 能力需求 / ProviderRequirement 引用，不放 provider 实现
  knowledge/             # 领域知识（knowledge-scopes.json + knowledge-index.json + 各 .md，权威资料结构化）
  _source-materials/     # Owner 投放的原始资料区（不进 Git，只留 .gitignore + README 占位）
  install.py / uninstall.py  # 经基座真实 lifecycle 端点装入 / 卸载
```

## 加载 / 卸载（经基座真实主权链）

```sh
# 起基座 truzhenos devserver（隔离 DB），与本仓平级放在同一父目录
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 environmental-enforcement-pack-v0/install.py
# 铁证：registry 0→1；知识入 09 FormalKnowledge；正式动作经 Owner + Base Gate 裁定并产 Receipt
```

## 护城河测试

每个场景荚必过：「用户能否一句话让最强模型直接做出同样结果？」能 → 不合格降级为 prompt / 模板；不能（因需真实系统门控 / 领域审批合规 / 可审计证据回执 / 显式多角色对照 / 事务对象生命周期回放）→ 合格。写入 manifest 的 `moat_justification`。

## 子包清单

见 [MODULES.md](MODULES.md)。

## License

[Apache-2.0](LICENSE)。
