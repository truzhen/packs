# 审核与上架

徒真（truzhen）Pack 市场交易的是可治理工作能力，不是普通插件或 prompt。审核目标不是阻碍作者，而是让企业敢安装、敢付费、敢在真实工作中使用陌生作者的 Pack。

## 审核会看什么

- 是否说清目标客户、业务场景和真实价值。
- 是否属于能力包、角色包或场景包之一。
- 是否声明 ProviderRequirement、风险等级、门控和 Receipt 口径。
- 是否存在 raw secret、真实凭据、运行态数据库、构建产物或未授权材料。
- 是否把候选、说明、fixture 或 mock 成功冒充真实接通。
- 是否承诺固定收益、跳过审核直接上架、直接执行或直接写正式数据。

## 试装会验证什么

- Pack 结构是否可解析。
- manifest、flow、role slots、role packs、capabilities 和知识索引是否一致。
- Provider 缺失时是否诚实显示 `provider_missing / not_ready / blocked`。
- 安装、启用、停用、候选生成和回执引用是否遵守主权链。

## 常见阻断原因

- README 只写愿景，没有目标客户和使用场景。
- 把 Provider 实现或外部软件运行代码放进本仓。
- 没有说明高风险动作如何回 Owner + Base Gate。
- 用“已接通”“可执行”“可发送”描述尚未验证的能力。
- 收益文案承诺固定收入。
- 缺少许可证、来源、知识核验或隐私边界说明。

## 上架后

上架后的商品展示、订单、支付、License / Entitlement、下载分发和运营状态归 `truzhen-cloud`。本仓仍只保存 Pack 资产、声明和可被商品化引用的元数据。
