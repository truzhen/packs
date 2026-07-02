# FEATURE_LEDGER — truzhen-packs 开发过程登记清单

> 本文件记录 `/Users/li/Documents/truzhen-packs` 当前“已经有什么、体量多大、还缺什么”。
> 它是本仓 Pack 资产进度账本，不是运行态事实源；正式启用、回执、知识挂载和执行结果归基座 Owner + Base + Gateway + Receipt 链路。

## 0. 本轮派活卡

| 维度 | 结论 |
|---|---|
| 我要做的事 | 建立根目录 `FEATURE_LEDGER.md`，登记当前 Pack 模块、体量、进度口径和待完善项 |
| 真实客户 / 场景证据 | Owner 明确要求更新 truzhen-packs 根目录开发过程登记清单 |
| 真相源 | 本仓文件结构、`AGENTS.md`、`README.md`、`MODULES.md`、各 pack `manifest.json` 与 tracked 文件统计 |
| 仓库 / 层归属 | 仅 `truzhen-packs` 包层资产仓 |
| 风险颜色 | 绿：文档账本新增，不改契约、不改 manifest、不改脚本、不改 Pack 行为 |
| 是否改契约 | 否 |
| 不允许碰的边界 | 不改基座仓、不改 provider、不改真实装入 / 卸载、不提交 raw source materials / 运行产物 |
| 验收方式 | Markdown diff、JSON 合法性、脚本语法、禁入产物扫描、Pack 结构审计 |

## 1. 当前基线

| 项 | 值 |
|---|---|
| 维护日期 | 2026-06-30（06-30 快照 salvage 于 2026-07-02 land；此后主线新增：六仓治理对齐、solid-waste scope 漂移已修、云边界红线，均见 `MODULES.md`） |
| 当前主线 | `origin/main` / `main` at `11d6434`（`docs: strengthen packs governance docs`） |
| 当前分支 | `codex/feature-ledger-refresh-20260630` |
| 本仓定位 | Pack 资产仓：只声明和编排，不持 Base 主权，不实现 Provider，不绕过 Gateway |
| 生命周期档位 | 整仓处于 `已实现 -> 已接线` 之间；单个 Pack 按下表区分成熟度 |

## 2. 状态图例

| 标记 | 含义 |
|---|---|
| ✅ | 已在本仓主线具备对应资产，且有基本验证入口 |
| 🟡 | 可用但存在明确缺口，不能包装成完整成品 |
| ⚪ | 尚未建立 |
| ⏸️ | 暂缓 / 待 Owner 裁定 |

## 3. 当前模块功能与体积

统计口径：`git -c core.quotePath=false ls-files` tracked 文件；不计未跟踪文件、本地运行缓存和基座运行态数据。

| 模块 / 目录 | 当前职能 | tracked 文件 | Markdown | JSON | 脚本 | 行数 | tracked 体积 | 状态 |
|---|---|---:|---:|---:|---:|---:|---:|---|
| `environmental-enforcement-pack-v0/` | 生态环境执法证据链 Domain Work Pack；含线索、立案、取证、证据三性、双角色对照、文书 / 处置候选、知识库和装卸载脚本 | 61 | 48 | 8 | 3 | 8419 | 1.2 MiB | 🟡 完整文件夹包主体已具备；知识域存在 `solid-waste` 索引漂移，法律知识仍需人工核验 |
| `smart-home-owner-pack-v0/` | 智能家居老板项目经营 Pack；以项目经理角色和 Frappe ProviderRequirement 组织项目快照、里程碑、采购、施工和写回候选 | 9 | 1 | 5 | 2 | 890 | 37.3 KiB | ✅ 完整文件夹包，无知识库 |
| `housekeeping-ops-pack-v0/` | 家政 / 保洁客户服务全生命周期 Pack；含咨询受理、顾问方案、质检质询、排期报价、派工、通知和归档候选 | 11 | 2 | 6 | 3 | 1237 | 72.9 KiB | 🟡 可装入 / 可卸载，仍缺 `knowledge/` |
| `templates/` | Pack 作者端软件工程模板，提供 manifest、docs、objects、flows、views、policies、software、tests 等骨架 | 18 | 16 | 1 | 0 | 650 | 32.5 KiB | ✅ 模板，不参与 enabled pack 分发 |
| 根治理文件 | `AGENTS.md`、`CLAUDE.md`、`README.md`、`MODULES.md`、`FEATURE_LEDGER.md`、CI / license / go.mod / ignore 规则 | 9 | 5 | 0 | 0 | 987 | 55.0 KiB | ✅ 治理入口已建立 |

## 4. 当前功能登记

| 功能 / 资产 | 职责 | 状态 | 位置 / 证据 | 当前口径 |
|---|---|---|---|---|
| Pack 仓治理边界 | 明确本仓是包层，不是基座、运行时、Provider、前端产品仓 | ✅ | `AGENTS.md`、`CLAUDE.md`、`README.md`、`MODULES.md` | Pack 只声明和编排 |
| 三类 Pack 封顶 | 固定 `Domain Work Pack`、`Capability Pack`、`Role Pack` 三类 | ✅ | `AGENTS.md`、`README.md` | 不发明第四种 Pack |
| 环保执法 Pack | 高风险合规 / 执法证据链样板，带知识库、双角色对照和 install / uninstall | 🟡 | `environmental-enforcement-pack-v0/` | 法律 / 监管知识默认 `pending_human_review`；当前 `solid-waste` scope 声明与 knowledge scopes 索引不一致，不能声称结构完全闭环 |
| 智能家居老板 Pack | 长周期项目交付样板，声明 Frappe 读写候选与项目经理角色 | ✅ | `smart-home-owner-pack-v0/` | Frappe 只是 ProviderRequirement，不是真相源 |
| 家政运营 Pack | 客服全生命周期样板，可装入 / 卸载、可声明角色 / 能力 / flow | 🟡 | `housekeeping-ops-pack-v0/` | 不能标完整知识包；卸载脚本只停用 Pack，不删除历史对象、候选或回执 |
| Pack lifecycle 胶水脚本 | 通过基座 devserver lifecycle 端点装入 / 卸载 | 🟡 | 各 pack `install.py` / `uninstall.py` | 无隔离 devserver 铁证时，不声称 E2E 通过 |
| 结构化知识 | 将可分发资料组织为 `knowledge-scopes.json`、`knowledge-index.json` 和 Markdown | 🟡 | 当前仅环保执法 Pack 完整具备 | 知识启用归基座 09；停用只改变可见性，不删除正式知识 |
| 作者端模板 | 给未来 Pack 作者提供工程骨架 | ✅ | `templates/scene-pack-software-template/` | 不是已发布 Pack，不参与 enabled 分发 |

## 5. 已校正的进度口径

本文件首次建立，不保留旧账本的历史表述。当前明确采用以下口径，避免与开发进度不一致：

| 旧风险口径 | 当前口径 |
|---|---|
| 把所有目录都称为完整 Pack | 只有环保执法、智能家居是完整文件夹包；家政 Pack 是可装入但未完整 |
| 把 ProviderRequirement 写成 provider 已接通 | 本仓只声明 provider 需求；真实 provider readiness 归基座 / provider 仓 |
| 把模板当成可启用 Pack | `templates/` 只是作者端脚手架 |
| 把高风险知识当成正式适用依据 | 法律、监管、合同等知识默认待人工核验 |
| 把 install 脚本存在等同于 E2E 通过 | 只有隔离基座 devserver 装入、enabled version、SlotBinding、KnowledgeMount、Receipt 证据齐全时才算 E2E |
| 把候选写成正式动作 | Pack 只能产候选、声明约束或提供上下文；正式动作必须 Owner + Base Gate |

## 6. 待开发 / 待完善 / 已收口项

| 优先级 | 归属 | 待完善项 | 当前缺口 | 验收口径 |
|---|---|---|---|---|
| ✅ | `housekeeping-ops-pack-v0/` | 补 `uninstall.py` | 2026-07-02 已补齐；脚本经 Base gated-action prepare / confirm 后调用 lifecycle disable | 卸载只停用 Pack，不删除历史对象、候选或 Receipt；脚本语法通过 |
| P0 | `housekeeping-ops-pack-v0/` | 补 `knowledge/` 或明确无知识库版本 | 当前无 `knowledge-scopes.json` / `knowledge-index.json` | README / manifest / MODULES 口径一致；若补知识，scope / index / Markdown 一致 |
| P0 | `environmental-enforcement-pack-v0/` | 收口 `solid-waste` 知识域漂移 | `manifest.json` 声明 `knowledge_scope://environmental/solid-waste`，但 `knowledge/knowledge-scopes.json` 缺该 scope，且当前无 `knowledge/solid-waste/` 目录 | Owner 裁定是补 `solid-waste` scope + 知识文件，还是从 manifest 移除该知识域；结构审计必须转绿 |
| P1 | `templates/scene-pack-software-template/` | 模板 manifest 贴近当前 Domain Work Pack schema | 当前模板 manifest 与真实 pack manifest 不是同一成熟结构 | 模板包含六件事、ProviderRequirement、moat、role slots、flow、knowledge 占位说明 |
| P1 | 全仓 | 建立自动结构审计脚本 | 当前验证命令写在治理文档中，未形成仓内统一脚本 | 一条命令覆盖 JSON、脚本语法、forbidden artifacts、manifest 结构审计 |
| P1 | 全仓 | 补 Pack 体量 / 成熟度变更登记纪律 | 本文件刚建立，后续新增 / 删除 / 重命名 pack 时需同步维护 | README、MODULES、FEATURE_LEDGER 与目标 pack README / manifest 同步更新 |
| P2 | `environmental-enforcement-pack-v0/` | 知识来源与人工核验持续维护 | 高风险法律知识体量最大，需防止 scope / index 漂移 | knowledge scopes、index、文件路径一致；高风险知识保持 `pending_human_review` |
| P2 | 全仓 | devserver E2E 验收记录 | 当前仓可静态验证；真实装入需相邻基座 devserver | 记录 registry 0→1、enabled version、Role Pack enabled、SlotBinding、KnowledgeMount、Receipt |

## 7. 维护规程

1. 新增、删除、重命名 pack，必须同步更新 `README.md`、`MODULES.md`、本文件和目标 pack README / manifest。
2. `✅` 只能用于本仓资产已具备且有可复核验证入口的功能；缺 provider、缺 devserver、缺卸载、缺知识库时必须写 `🟡`。
3. 本文件只登记 Pack 仓资产事实；基座运行结果、正式回执、正式知识和正式执行不在本文件伪造。
4. 每次收尾至少运行：`git diff --check`、JSON 合法性、脚本语法、禁入产物扫描；涉及 manifest / knowledge 时加跑结构审计。
5. 对话汇报超过 500 字时写入 Markdown 文件；对话只给路径、验证和待 Owner 裁定项。
