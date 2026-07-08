# 墅学家大宅装修设计指导 Pack

`shuxuejia-renovation-pack-v0` 是面向大宅 / 别墅装修设计师和项目负责人的 Domain Work Pack。它把旧墅学家 product seed 中的 EPC 原版 457 节点 / 543 边拓扑迁移为可安装的 folder pack，用于组织设计准备、深化设计、材料合同付款、现场质量验收和售后保修的候选工作流。

## 主权边界

- 本 Pack 只声明流程、角色、能力需求、知识域和治理策略。
- 7 个角色包全部是 Proposer，只能提出建议、质询、任务候选、沟通草稿和证据要求。
- 设计确认、合同变更、付款、验收、对外发送和售后关闭必须经 Owner + Base Gate。
- Provider 未接通时必须返回 `provider_missing / not_ready / blocked`，不得假成功。
- 本仓不包含原始 PDF、真实客户资料、provider 实现、数据库、日志或密钥。

## 来源

- 只读来源仓：`/Users/li/Documents/truzhenv3`
- 来源提交：`3c161aff4e2cc9b6ee896d0ef1ec4b37aaf4b062`
- 来源文件：`backend/internal/devserver/pack_asset_seed_shuxuejia_topology.md`

## 安装

先启动隔离 `truzhenos` devserver，然后执行：

```sh
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 shuxuejia-renovation-pack-v0/install.py
```

卸载：

```sh
TRUZHEN_DEVSERVER_BASE=http://127.0.0.1:18099 python3 shuxuejia-renovation-pack-v0/uninstall.py
```

## 验收口径

完成声明必须同时提供 JSON / 脚本 / 结构审计、457/543 拓扑校验、真实 install / uninstall / reinstall 证据、项目运行候选证据、Receipt 引用和前端用户视角验收记录。前端缺入口时必须登记缺口，不能用 API 冒充前端验收。
