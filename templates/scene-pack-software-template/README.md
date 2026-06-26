# Scene Pack Software Template (作者端软件工程模板)

本项目是 Truzhen 平台的 **Scene Pack Software Project** 标准工程模板。它作为一个独立的软件作品，致力于帮助 Pack 开发者对业务流程、数据模式、软件适配及用户体验进行完整的工程化定义、测试和文档化。运行态仍由 Truzhen Runtime 解释安全规格；作者工程只作为交付、审查和本地仿真的材料包。

## 💡 核心设计理念与边界

1. **运行态与确权态分离**：
   - **运行态 (Truzhen Engine)**：仅读取安全的、无副作用的 `Scene Pack Spec`（声明式规格），确保宿主系统及网络的主权控制安全。
   - **确权态 (Software Project)**：提供完整的适配代码、单元测试、测试报告以及规范的用户手册与设计说明书，构成可独立登记和售卖的独立软件作品包。
2. **三层设计架构**：
   - **剧本规格层 (Scene Pack Spec)**：核心的流程 DAG、对象 Schema、门禁策略定义，必须保持声明式、候选态和无真实副作用。
   - **适配器层 (Pack Adapter Layer)**：声明外部软件（如 Dendrite, CAD, PS, PS, ERP 等）的接口契网与 Mock 行为。
   - **视图层 (Pack UI Surface)**：声明式配置标准 Slot 内的卡片、看板和仪表盘。Pack UI 只能声明到 Truzhen 标准 UI slot，由 `modules/12-frontend-shell` 投影为只读页面、候选卡或事务胶囊。
3. **主权控制与 Gateway Candidate**：
   - 任何真实发送、真实执行、正式写入或真实模型调用，都必须转成 Gateway Candidate，经 Base Gate / OwnerDecision 裁定后，再由对应网关受控处理。

---

## 📂 目录结构与说明

```text
my-pack/
├── README.md                          # 项目说明书
├── LICENSE.md                         # 软件授权协议
├── AUTHORS.md                         # 作者与所有权声明
├── pack.toml                          # 工程配置与平台要求
├── manifest.json                      # 静态审查清单与权限声明
│
├── objects/                           # 业务对象定义层 (Business Objects)
│   └── renovation_project.object.json # 例如：装修项目对象 schema
│
├── flows/                             # 场景流程定义层 (Scene Flows)
│   └── main.flow.json                 # DAG 流程设计图 spec，节点只生成 Candidate
│
├── views/                             # 声明式视图配置 (UI Surfaces)
│   ├── dashboard.ui.json              # Slot: pack.dashboard 视图配置
│   └── object-detail.ui.json          # Slot: pack.object.detail 视图配置
│
├── policies/                          # 核心治理策略 (Policies)
│   ├── gate-policy.json               # 主权门禁 Gate 策略
│   ├── receipt-policy.json            # 证据账本 Receipt 策略
│   └── permission-policy.json         # 细粒度资源访问许可策略
│
├── capabilities/                      # 基础能力需求与绑定 (Capabilities)
│   └── requirements.json              # 运行所需的模型/通信/执行等能力需求声明
│
├── adapters/                          # 适配器定义层 (Pack Adapters)
│   ├── interface-manifest.json        # 适配器外部接口清单
│   └── mock-adapter.ts                # Mock 适配器与接口契约，禁止包含 raw secret
│
├── software/                          # 行业软件调用声明 (Software Integrations)
│   ├── cad/
│   │   ├── adapter-manifest.json      # CAD 调用适配器声明
│   │   ├── usage.md                   # CAD 调用说明书
│   │   └── safety-policy.json         # CAD 安全隔离策略
│   └── ps/
│       ├── adapter-manifest.json      # PS 接口调用声明
│       └── safety-policy.json         # PS 运行安全定义
│
├── src/                               # 辅助业务逻辑 (Validators/Transforms)
│   └── pack-entry.ts                  # Pack 自写 validator 与 transform 辅助代码
│
├── generated/                         # 由 Pack SDK 自动生成的强类型代码
│   └── schema.gen.ts                  # 授权 Pack 作者在本项目内免费使用
│
├── tests/                             # 严苛单元与仿真测试 (Test Suites)
│   ├── pack-smoke.test.ts             # 冒烟测试：静态审查与 schema 校验
│   ├── flow-simulate.test.ts          # 流程模拟：无真实发送/真实执行的 DAG 推进
│   └── redteam.test.ts                # 红队测试：拦截真实密钥、发送与覆盖操作
│
├── docs/                              # 独立软件确权文档 (Documentation)
│   ├── 软件说明书.md                  # 用于软件著作权登记的标准材料之一
│   ├── 设计说明书.md                  # 整体业务架构与数据架构详细设计
│   ├── 用户手册.md                    # 针对最终用户或执导 Agent 的使用手册
│   ├── 测试报告.md                    # 冒烟测试、仿真测试及红队测试的完整通过记录
│   └── 版本说明.md                    # 更新日志、版本兼容性声明及 Pack 升级建议
│
└── build/                             # 构建与发布输出
    ├── pack.lock.json                 # 版本依赖锁
    └── export-manifest.json           # 发布到 Truzhen 平台的静态包描述
```

---

## 🛡️ 安全合规规则

1. **默认禁用与受控动作**：本工程中的所有 Adapter 与外部软件声明，在 Truzhen 运行态中默认为禁用 (disabled) 或 Mock。任何真实发送与执行动作，必须单独安装 Real Provider，并通过 `Base Gate` 获得主人授权，产生追加式回执 (Receipt)。
2. **严禁包含敏感凭证**：不得在本工程任何位置（包括代码、配置、文档中）写入真实 API 密钥、数据库连接串或账号 Token。
3. **禁止伪造及低俗凑数**：我们不支持且不建议生成任何无意义代码来凑数。本工程的所有代码、测试、适配逻辑都是真实可用且完全可执行的软件组成部分。
4. **接收自实验资产的设计原则**：
   - 保留运行态规格与作者工程分离。
   - 保留 Pack UI Surface / Capsule / Candidate Card 的声明式投影方式。
   - 保留 flow simulate、adapter mock、redteam 和 no-real-send/no-real-exec 的测试分层。
   - 不接收实验运行日志、数据库、构建产物、真实 provider、真实 secret 或外部执行脚本。
