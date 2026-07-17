# 内容运营工作台 Pack 0.1.1 发布说明

## 本版内容

- 将统一模型输出 Schema 拆为选题雷达、内容生产、周复盘三套按 `skill_id` 选择的 Schema。
- 内容生产候选强制包含标题、30–60 秒口播稿、镜头清单、字幕重点、封面文案、置顶评论、证据引用和 Owner 判断项。
- Skill bundle 以 framed SHA-256 覆盖入口、输出契约、兼容 Schema 和三套业务 Schema。
- 修复同版本停用后重启与新版本装入的 lifecycle 判定，旧版本 disabled 记录不再阻断 0.1.1。
- 补齐 canonical `scene_pack` 类型、最低 Truzhen 版本和 Provider 复用策略，使交付打包器可以 fail-closed 验证。

## 封装与安装

从 `truzhen-packs` 根目录执行：

```bash
python3 -B build_pack_bundle.py content-operations-workbench-v0 <仓外输出目录>
```

解压后，面向已运行的 Truzhen 基座执行：

```bash
TRUZHEN_DEVSERVER_BASE=http://<基座地址> python3 content-operations-workbench-v0/install.py
```

装入、角色绑定、计划启用、每次 Hands 运行、停用与重启继续由基座 lifecycle、01 Base、07 Scheduler、08 Model Gateway、11 Execution Gateway 和 03 Receipt Ledger 持有事实。

## 非目标

本 Pack 不包含 Codex CLI、模型、Provider、账号、Cookie、token 或平台发布器；不登录、不上传、不发布、不评论、不私信。代码合并与 Git 发布不代表 Pack 已在产品中安装、启用或在云市场上架。
