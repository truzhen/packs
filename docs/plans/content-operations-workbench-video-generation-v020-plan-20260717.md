# 内容运营工作台 Pack 0.2.0 真实视频生成计划

日期：2026-07-17  
目标版本：`scene_pack://content-operations-workbench@0.2.0`  
当前状态：`设计中`（仅完成计划；0.1.1 的既有验收状态不变）  
Owner 投入目标：每天只做方向、素材权利确认、预览与最终发布判断，单日尽量不超过 30 分钟。

## 1. 派活卡

| 项目 | 裁定 |
| --- | --- |
| 我要做的事 | 把已批准的内容候选与 Owner 授权素材渲染为真正可播放的竖屏 MP4 候选；文档、JSON 和脚本只作为审计伴随物，不能再冒充最终交付。 |
| 版本 / 优先级 | 当前 Owner 自媒体 dogfood 主线，P0；不是 V4 Governance Console，也不是“全自动媒体公司”大工程。 |
| 真实客户 / 场景证据 | Owner 已用 13 万粉的别墅装修与灯光抖音账号测试转型；0.1.1 已经通过真实 GUI 生成 45 秒抖音候选，但用户明确反馈“要求能够生成视频，而不是文档”。这是当前直接使用证据，不是推演需求。 |
| 最小可交付 | 1 个已批准的 45 秒内容候选 + 一组 Owner 授权图片/视频/音频，生成本地 9:16 MP4、封面、字幕、音轨、来源清单、哈希与可反查 Receipt；Client 可预览、重生成和下载。 |
| 真相源 | 内容事实归产品 Evidence；素材与权利声明归 Owner 项目 `ArtifactBinding`；软件安装与 readiness 归 02 Registry；渲染运行归 11 Execution Gateway；模型/TTS usage 归 08 Model Gateway；视频文件归本地 Artifact Store；正式 Gate 与 Receipt 分别归 01 / 03。Pack 只声明工作台与需求。 |
| 仓库 / 层归属 | `truzhen-packs` 只改 Pack 声明；`truzhenv3-software` 登记和封装外部软件；`truzhenos` 完成 Provider、Gateway、Artifact、Gate、Receipt、监控；`truzhen-client-web-desktop` 完成预览和判断面。 |
| 风险颜色 | Pack 声明为黄；Provider、网络、模型与 Gateway 接线为橙；真实执行、凭据、下载、删除和平台发布为红。本轮不做平台登录、上传或发布。 |
| 契约影响 | 首选复用现有 CandidateArtifact / ProviderRequirement / GatedAction / Receipt 形状；若无法表达视频探针、素材权利或渲染进度，必须先做 contracts 影响清单与向后兼容方案，不能在 Packs 私造隐藏契约。 |
| 禁止边界 | Pack 内不得放 FFmpeg 二进制、Provider runtime、模型、密钥、账号、Cookie、生成视频、缓存或发布脚本；前端不得直调渲染器；不得自动抓取无权利素材、安装任意 ComfyUI 节点、上传社媒或把 `provider_missing` 包装成成功。 |
| 用户如何验收 | 在 GUI 选择一期 Truzhen 宣传候选和授权素材，确认一次渲染，看到进度并直接播放 MP4；字幕、画面、声音、时长正确，可下载；Receipt 能反查输入、执行、artifact hash；全程无平台发布。 |
| 先输出什么 | 本文件先给影响清单、选型与分阶段计划。进入跨仓实现前，Owner 需按仓重新授权；计划本身不授权修改 `truzhen-software`、`truzhenos`、Client 或 contracts。 |

## 2. 本轮最小闭环与砍项

### 2.1 必须进入 0.2.0

1. 输入必须是已通过事实审校和 Owner 方向确认的 `ContentPublicationPackageCandidate`，渲染器不得自行重写产品主张。
2. 素材优先使用 Owner 自有或明确授权的别墅、灯光、Truzhen 演示图片/视频；每项保存 `source_ref`、`content_hash`、权利声明与允许用途。
3. 输出为本地 1080×1920、H.264/AAC、30–60 秒的 MP4 候选；同时生成封面图、烧录字幕、音频响度结果、`ffprobe` 元数据与 SHA-256。
4. TTS 只走已登记且 ready 的 08 Provider；Owner 自录音可作为替代输入。缺 Provider 时必须 `provider_missing / not_ready`，不能用静音假片冒充完成。
5. 11 负责受控渲染执行、工作目录隔离、资源上限、artifact capture、幂等与失败收口；01 / 03 提供逐次 Gate 与可反查 Receipt。
6. Client 提供输入摘要、素材权利确认、渲染确认、进度、视频播放、A/B 候选对比、下载与“退回修改”；默认没有“发布到抖音”按钮。
7. `truzhen-monitor` 记录队列时间、渲染时长、CPU/GPU、失败原因、Provider readiness、artifact 大小与 Receipt 连续性，不另造监控格式。

### 2.2 明确砍入 backlog

- 自动登录、上传、定时发布、私信、评论、数据抓取和多账号矩阵。
- 自动搜索或下载版权不明的库存视频；首版仅用 Owner 授权素材。
- 数字人、口型驱动、声音克隆、完整时间线编辑器、多人协作、云渲染。
- 一开始同时适配抖音、TikTok、YouTube Shorts、B站和视频号；0.2.0 只交付通用 9:16 MP4。
- 把文生视频作为首版硬依赖。AI B-roll 属增强轨，缺 GPU 或模型时不阻塞确定性成片。

## 3. 技术选型

### 3.1 推荐主链：FFmpeg 确定性渲染

`FFmpeg + ffprobe` 负责裁剪、缩放、转场、音画合成、字幕烧录、响度处理、编码与产物探针。它是首版最小、稳定且可核验的渲染骨架；官方文档覆盖编码器、滤镜、格式与 `ffprobe`：[FFmpeg Documentation](https://www.ffmpeg.org/documentation.html)。

Pack 只声明 `capability://video-candidate-render` 和软件/Provider requirements；薄适配器、二进制版本锁、进程监督与沙箱必须放在 `truzhenv3-software` / 11，不得进入 Pack。

### 3.2 可复用开源项目

| 软件 | 建议档位 | 用途与边界 |
| --- | --- | --- |
| FFmpeg | 0.2.0 必选 | 确定性成片与探针；固定版本/哈希，限制输入路径、协议、时长、分辨率和资源。 |
| MoneyPrinterTurbo | 0.2.x 可选首适配器 | 已覆盖脚本、素材匹配、字幕、BGM 和短视频合成，可复用其编排思路或受控 Provider；必须固定 release/commit、审计依赖，且只接受已批准脚本和 Owner 素材，不能绕过 Gateway 抓取外部素材。官方项目：[MoneyPrinterTurbo](https://github.com/harry0703/MoneyPrinterTurbo)。 |
| Pixelle-Video | 0.3 增强候选 | 适合图像/视频生成、TTS 与模块化 Provider 编排；本仓旧证据仍指向旧组织地址，实施前必须复核仓库迁移、许可证、维护者与依赖，不得静默沿用旧证据。当前项目：[Pixelle-Video](https://github.com/ATH-MaaS/Pixelle-Video)。 |
| ComfyUI | 0.3 增强底座 | 只用于经白名单锁定的图/视频生成工作流；禁止运行时任意安装 custom nodes，节点、模型与 workflow 必须锁版本和哈希。官方项目：[ComfyUI](https://github.com/comfy-org/comfyui)。 |
| LTX-2 / LTX-Video | 0.3 GPU 候选 | 适合 image-to-video / B-roll；硬件和磁盘成本较高，必须先做本机或云 Provider benchmark，不作为 0.2.0 放行条件。官方项目：[LTX-2](https://github.com/Lightricks/LTX-2)、[ComfyUI-LTXVideo](https://github.com/Lightricks/ComfyUI-LTXVideo)。 |

选型决策：0.2.0 先落 `FFmpeg + 已有 08 TTS + Owner 素材`；MoneyPrinterTurbo 做同一 Provider 接口下的可替换适配器验证；Pixelle/ComfyUI/LTX-2 仅在“实际成片价值明显高于安装与 GPU 成本”的 benchmark 通过后进入 0.3。首版不新增 Remotion/MoviePy 平行实现，除非 FFmpeg 薄适配器的可维护性实测不达标。

## 4. 目标链路

1. Pack 的“内容生产”模式生成并事实审校 45 秒脚本、镜头单、字幕与 CTA。
2. Owner 选方向，挂载本地素材并确认权利、声音和 BGM 策略。
3. 04 根据 02 ProviderResource 投影视频渲染与 TTS readiness；缺项即 fail-closed。
4. Client 提交 `VideoRenderCandidate`，展示输入摘要、成本、网络、素材来源和风险。
5. Owner 确认后，由 01 签发逐次 T06；11 运行受控渲染 Provider，08 只承担需要的 TTS/生成模型调用。
6. FFmpeg 产出 MP4、封面与探针结果；Artifact Store 保存本地文件，03 记录可反查 Receipt，`truzhen-monitor` 收集运行指标。
7. Client 播放真实 MP4，Owner 选择接受、退回重生成或下载后人工发布。接受视频不等于已发布。

## 5. 分仓影响清单

| 仓库 | 计划改动 | 验收映射 |
| --- | --- | --- |
| `truzhen-packs` | 升级 manifest 至 0.2.0；新增视频候选工作模式/阶段、能力与软件需求、视频候选对象字段、素材权利门、Skill 输出 Schema、安装/卸载与测试；不放 Provider。 | JSON/Schema/hash、bundle、install/uninstall、forbidden artifact、`go test ./...`、`git diff --check`。 |
| `truzhenv3-software` | 登记 FFmpeg 与可选 MoneyPrinterTurbo；锁版本/源码/许可证/哈希；提供受监督 adapter、probe、readiness、资源配额和清理策略。 | 安装与 source lock、离线 probe、恶意输入、超时/取消、路径逃逸、无网络 FFmpeg smoke。 |
| `truzhenos` | 02 ProviderResource；04 readiness；08 TTS/生成模型 usage；11 render execution 与 artifact capture；01 Gate；03 Receipt；统一监控事件。 | Provider 专项 E2E、Gate/Receipt 反查、幂等/急停、artifact hash、真实 FFmpeg 渲染与 `truzhen-monitor` 诊断包。 |
| `truzhen-client-web-desktop` | 素材选择与权利确认、渲染预检、进度、播放器、候选对比、下载、失败原因与 Receipt/Evidence 展示。 | typecheck/test/build、Candidate/Formal 隔离、GatedAction、真实 MP4 播放 GUI smoke 与截图/录屏证据。 |
| `truzhen-contracts` | 默认不改；只有现有形状不能表达 render spec、artifact probe 或 rights refs 时才提出兼容增量。 | 先做反向依赖清单、schema/DTO 双向测试、消费者固定 SHA EGR；未经新授权不得改。 |

## 6. 实施阶段与工期

| 阶段 | 产物 | 预计工程时间 | Owner 判断点 |
| --- | --- | --- | --- |
| P0 基线与 benchmark | 选一份真实候选、10–20 个授权素材，生成 15 秒 FFmpeg 样片；盘点本机硬件、已有 TTS 与 software registry | 0.5–1 天 | 确认素材可用与样片风格 |
| P1 契约/安全设计 | 现有契约可复用结论、render spec、权利清单、资源/网络/失败策略 | 1 天 | 只裁定有争议的字段或风险 |
| P2 Pack 0.2.0 | flow、capability、object、Skill Schema、install/uninstall、测试 | 1–2 天 | 确认工作流和默认 9:16 模板 |
| P3 Software Provider | FFmpeg adapter + probe + supervisor；MoneyPrinterTurbo 可选 spike | 2–4 天 | 决定是否保留 MoneyPrinterTurbo 轨 |
| P4 OS 接线 | 02/04/08/11、Artifact、Gate、Receipt、监控真实闭环 | 2–3 天 | 确认每次渲染 Gate 摘要 |
| P5 Client 成片面 | 进度、播放、A/B、退回、下载、证据与错误状态 | 1–2 天 | 从 GUI 选出一期候选 |
| P6 独立验收 | 生成一期 45 秒 Truzhen 宣传 MP4，对抗测试与跨仓固定 SHA EGR | 1–2 天 | 观看、批准或退回；仍由 Owner 人工发布 |

总工期建议：8–12 个工程日；若 Owner 只投入 2 小时/天做开发，则约 3–5 周。若 Codex/Truzhen 承担施工、Owner 每日只做 15–30 分钟判断，可按 2 周左右完成首个本地成片闭环。生成式 B-roll 不计入该工期。

## 7. 严格验收

### 7.1 机器验收

- `ffprobe` 必须证明产物存在非零视频轨与音频轨，分辨率 1080×1920、时长 30–60 秒、编码 H.264/AAC；字幕与封面存在且 hash 与 Receipt 一致。
- 同一 render request / input hash 重放不得重复执行或产生两个正式产物；素材、脚本、声音或 render spec 变化必须产生新候选。
- 缺素材权利、Provider、TTS、T06、run_id、nonce、artifact store 或 03 ledger 时必须 `blocked / provider_missing / not_ready`，不得写空 MP4、复制 fixture 或返回模板成功。
- 路径穿越、危险输入协议、超大文件、不支持 codec、超时、取消、EmergencyStop 与磁盘不足必须 fail-closed，并留下失败 Receipt / monitoring event。
- 本地确定性渲染轨默认不得联网；任何模型或素材网络调用必须走动态网络 Gate，并能由 usage/Receipt 回答“哪一期视频消耗了什么模型与资源”。
- Git 扫描必须证明没有 MP4、MOV、模型、缓存、二进制、密钥、Cookie、日志和本地运行 state 被提交。

### 7.2 真实用户 GUI 验收

Owner 以普通用户完成：选择“做一期 Truzhen 抖音宣传”→查看脚本事实边界→选择/挂载授权素材→确认渲染→看到真实进度→直接播放成片→核对字幕、声音、画面和 CTA→退回一次修改并生成第二候选→下载选中 MP4。验收记录必须包含 GUI 证据、artifact hash、T06、08 usage（如有）、11 final Receipt 和 03 原文反查。

执行方的 DONE、编译通过或文档说明不算验收。必须由独立验收主体进行对抗抽查，并运行 Packs、Software、OS、Client 各自专项门禁及固定提交组合的全链路 EGR。

## 8. 生命周期与完成口径

- 0.1.1：保持`已验收（打包前）`，能力上限仍是内容候选文档。
- 0.2.0 视频升级：当前为`设计中`。
- Pack 声明、Schema 和测试完成后最多到`已实现`；Provider/OS/Client 接通后到`已接线`；真实 GUI 成片与独立验收通过后到`已验收`；只有安装启用并完成发布流程后才可标`已发布`。
- “生成了脚本、镜头表、渲染命令、占位链接或 mock MP4”都不等于生成视频；完成判据是用户能播放、下载并由 Receipt 反查的真实 MP4 候选。

## 9. 开工前 Owner 授权卡

建议下一步按以下顺序单独授权：

1. `truzhen-packs`：修改、测试、提交视频升级声明与 Pack 测试。
2. `/Users/li/Documents/truzhenv3-software`：读取、修改、测试、提交 FFmpeg Provider 与软件 registry；不得读取共享登录态或写入凭据。
3. `/Users/li/Documents/truzhenos`：读取、修改、测试、提交 02/04/08/11、Artifact、Gate、Receipt 与监控接线；不得改变 Base 硬地板。
4. `/Users/li/Documents/truzhen-client-web-desktop`：读取、修改、测试、提交 GUI 视频候选面；不得直连 Provider 或发布平台。
5. 仅在影响评估证明必要时，再单独授权 `truzhen-contracts`。

首轮应继续禁止社媒发布、真实客户触达、自动外部素材抓取、跨仓 merge/push 与生产部署；这些动作需要另行裁定。
