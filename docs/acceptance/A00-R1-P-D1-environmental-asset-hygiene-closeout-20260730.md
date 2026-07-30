# A00-R1-P-D1 环保执法资产卫生 closeout

## 结论

P-D1 已在 Packs current-main 精确基线 `511a72c46c38bfd6e458049ab14e3549752f1a29` 上完成最小修复：

- 两个索引共 116 个 Markdown 跳转全部保留并改指 current-main 现存目录，失效链接从 116 降为 0。
- README 案例口径从 31 校正为 Git 与 `knowledge-index.json` 共同证明的 30。
- 两个索引内容变化只同步其 `knowledge-index.json` SHA-256；15 scope、45 documents、30 cases、所有其它 checksum、`pending_human_review` 与 `reference_only` 均未改变。
- 已新增确定性回归门，后续旧目录回潮、链接断裂、15/45/30 漂移、checksum 漂移或法律状态升级都会失败。

生命周期为 `已实现 -> 待 P-E2 / 独立验收`，未发布。

## 派活卡与边界

| 字段 | 本节点裁定 |
| --- | --- |
| 版本 / 优先级 | A00-R1 / P0 runtime 修复的 Packs P-D1 前置卫生节点 |
| 真实场景证据 | A00-G0 在 current main 复现两个索引 116 个旧目录失效链接，并确认 README 写 31、实际案例与 index 均为 30 |
| 最小可交付 | 只修 116 个链接、31→30、对应 checksum、测试、账本与 closeout |
| 真相源 | 当前 Packs Git 资产：manifest 与 scopes 的 15 集合、index 的 45 entries、`knowledge/cases/*.md` 的 30 文件 |
| 仓 / 层归属 | 仅 `truzhen-packs` 声明与索引层；不读取或修改 Contracts、Client、Cloud、Software 产品状态 |
| 风险颜色 | 绿：确定性路径与文档口径修复；未改变法律内容、主权链、运行时或契约 |
| 契约影响 | 无 route、DTO、schema、manifest scope、ProviderRequirement 或 Client consumer 变化 |
| 禁止边界 | 未改 `install.py`、`uninstall.py`、manifest、其它知识正文、法律状态、历史 16/53；未启动动态实例、登录、Provider、生产动作、完整 EGR 或 pre-push |
| 外部动作 | `external_actions=0` |
| 任务图 | 复用密封图 P-D1 单写 Loop；P-E2 等待 O-T1 交接后由同一 Packs writer 继续，本节点现在停住 |
| 回边预算 | 初次实现；Packs 共享修复回边计数仍为 0/2 |
| 停止条件 | P-D1 clean local commit 后停止，不执行 P-E2、不 merge、不 push |

## 变更

1. `knowledge/index/法条速查索引.md`：66 个链接迁到 current-main scope 目录。
2. `knowledge/index/违法行为分类索引.md`：50 个链接迁到 current-main scope 目录。
3. 当前 15 scope 没有固废专属指南；原 3 个失效固废指南链接改指现存且包含固废检查内容的 `pollution-source-overview/01_污染源执法概述.md`，没有恢复 `solid-waste` scope。
4. `knowledge/knowledge-index.json`：只更新上述两份索引的 SHA-256。
5. `environmental-enforcement-pack-v0/README.md`：案例数校正为 30。
6. `environmental_knowledge_asset_hygiene_test.go`：锁定 15/45/30、116 links、全部链接存在、全部 Markdown 已登记、checksum、`pending_human_review + reference_only`。
7. `FEATURE_LEDGER.md`：登记 P-D1 生命周期和边界。

## TDD 与验证证据

初始红灯：

```text
README 未按实际资产声明 30 个真实案例
README 仍残留 31 个真实案例漂移
两个索引仍有旧目录链接 116 个
两个索引仍有失效链接 116 个
```

绿色验证：

```text
GOWORK=off go test ./... -run '^(TestEnvironmentalKnowledgeAssetHygiene|TestForbiddenPatternsCatchBusinessPII|TestForbiddenDatabaseArtifact|TestPackAssetsDoNotCarryBusinessDataFormalRefsOrRawSecrets)$' -count=1
ok github.com/truzhen/packs

python3 -m unittest test_knowledge_checksums
Ran 7 tests
OK

递归 JSON：合法
forbidden artifacts：0
git diff --check：通过
```

精确资产结论：

| 检查 | 结果 |
| --- | ---: |
| manifest / scope manifest | 15 / 15，集合相同 |
| knowledge index | 45 / 45 |
| case entries / case Markdown | 30 / 30 |
| 两个索引链接 | 66 + 50 = 116 |
| 失效链接 / 旧目录链接 | 0 / 0 |
| checksum | 45 / 45 通过 |
| `pending_human_review` | 45 / 45 |
| `reference_only` | 45 / 45 |
| forbidden artifacts | 0 |

本 closeout、修复、测试与账本必须位于同一 clean local commit；确切 commit SHA 由节点回执记录，以避免在提交正文内自引用不稳定 SHA。
