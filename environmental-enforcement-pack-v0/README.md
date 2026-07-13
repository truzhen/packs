# 环保执法证据链 Pack（v0，全领域；v10 修复已实现、待单 Pack 复验）

一个**独立文件夹、可加载、可卸载、不与基座主程序混**的环保执法场景荚（Domain Work Pack）。
知识库基于 Owner 提供的**权威资料**真实导入（生态环境法典 + 生态环境保护综合行政执法指南 +
当前生态环境领域热点难点问题案例精析），覆盖水/大气/固废/危废/噪声/土壤/生态/核辐射污染领域
＋行政处罚程序、证据规则、行刑衔接与执法风险防范。

## 目录结构

```
environmental-enforcement-pack-v0/
├── manifest.json                  # 总纲：六件事 + 护城河 + 知识域 + provider 需求 + namespace 隔离
├── flows/
│   └── environmental-enforcement-flow.flow.json   # 19 节点执法流程画布（中性 FlowSpecDraft）
├── role-slots/role-slots.json     # 角色槽（执法精英 / 挑剔律师）+ 绑定
├── role-packs/
│   ├── enforcement-elite.rolepack.json            # 执法精英 persona
│   └── critical-lawyer.rolepack.json              # 挑剔律师 persona
├── knowledge/
│   ├── knowledge-scopes.json      # 15 个知识域声明（mount_on_pack_enable）
│   ├── knowledge-index.json       # 45 份源文档 → scope/kind/source_ref/title/生效日期
│   ├── code/ legal-basis/ water/ air/ radiation/ noise/
│   ├── eia-permit/ ecology/ penalty/ criminal/ risk/
│   ├── guide-overview/ pollution-source-overview/
│   ├── cases/                     # 31 个真实案例（行政篇 + 公益诉讼篇）
│   └── index/                     # 法条速查索引 + 违法行为分类索引
├── capabilities/capabilities.json # 5 个 provider 需求（OCR/文书生成/送达/执法动作/在线监测）
├── install.py                     # 装入（加载）到正在运行的 devserver
├── uninstall.py                   # 卸载
├── tools/build-knowledge-from-source.py   # 从权威资料重建 knowledge/（可重跑）
└── _source-materials/             # 原始权威资料投放区（.gitignore，不进库）
```

## 主权与红线

- **AI 全程 Proposer**：执法精英/挑剔律师只出违法事实候选、法律依据候选、质询候选，无正式裁定权。
- **正式动作必经 Owner + Base Gate**：行政处罚、责令改正、查封扣押、文书外发、移送公安须主权确认 + Base 签发 `decision_ref/run_id/nonce`，产 03 Receipt。
- **provider 未接通诚实降级**：返回 `blocked / provider_missing / not_ready`，不假成功。
- **namespace 隔离**：`environmental_enforcement_ns`，与其它 pack 不串。
- **知识 pending_human_review**：knowledge 内容来自权威资料导入，每条标 `verification_status=pending_human_review`；正式适用前须经法务/业务核验，以现行有效官方法规标准原文为准。

## 加载 / 卸载（不与基座混；产品默认不自带本 pack）

> 前置：本 pack 已从基座 `server.go` 的产品自动 seed 中摘除（`defaultPlatformPackAssetSeeds()` 返回空）。
> 即产品启动后「场景包管理」为空，本 pack **只在跑 install.py 后出现、跑 uninstall.py 后消失**。
> 因此需用**含本分支改动的 devserver**（旧二进制仍会自带 1.0.2 旧种子）。

```sh
# 1) 起 devserver（含本分支改动）
go run ./backend/cmd/devserver            # 默认 127.0.0.1:18080

# 2) 加载本 pack（走真实 lifecycle 端点 + Base Gate + 03 receipt）
python3 packs/environmental-enforcement-pack-v0/install.py
#   或指定地址：TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 .../install.py

# 3) 前端「场景包管理」刷新 → 出现「环保执法证据链 Pack（全领域）」

# 4) 卸载
python3 packs/environmental-enforcement-pack-v0/uninstall.py
```

加载会装入：场景包 EnabledVersion + 2 角色包 + 槽位绑定 + 15 知识域 + 45/45 源文档切分出的全部 FormalKnowledge。源文档数与 09 运行态分片数是两个指标；验收必须按 `pack_ref`、`source_ref`、scope 全量分页对账，不得把前台全局知识总数（历史曾显示 752）冒充本 Pack 的源文档数或固定分片契约。
（大文件按章节/条款切块，约数百条，均 pending_human_review）。卸载经真实 `lifecycle/disable`
（Base 签发 disable 决议）级联卸载知识域；**已产生的案件对象与 03 回执仍可反查——卸载不删历史**。

## 重建知识库（资料更新后）

```sh
python3 tools/build-knowledge-from-source.py [权威资料知识库路径]
# 默认源不写入 Pack；由本机过程文档或受控导入配置提供。
# 会重写 knowledge/ 下各域文件 + knowledge-scopes.json + knowledge-index.json
```

## 历史验证与当前状态

- 2026-06-25 的隔离 devserver 历史记录曾验证 install、角色绑定、KnowledgeMount 与 disable/re-enable 生命周期。
- 2026-07-11 R1 修复复验曾记录前台全局 752 条知识；该数字不再作为本 Pack 数量契约。当前 Pack 权威资产为 45 份源文档、15 个 scope，运行态必须以本 Pack 归属全量分页覆盖为准。法律知识仍为 `pending_human_review`，Pack 尚未发布。
- C01 法律金标见 `docs/C01-法律金标-20260711.md`。该金标只约束候选输出和人工复核，不构成正式法律意见。
