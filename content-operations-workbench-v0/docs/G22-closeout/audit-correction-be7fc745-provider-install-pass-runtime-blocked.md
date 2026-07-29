# G22 当前 Provider 安装成功、runtime/auth 未就绪审计更正

日期：2026-07-29。范围仅为 G22 去敏 closeout；不修改协调会话的 live root、Gate、Receipt、registry 或容器。

旧 `08b687…` source-image 缺失结论属于旧 Software 基线，现作废。当前固定 OS `be7fc745…`、Software `0768faa1…` 的协调现场以全新 task-owned runtime 完成：

1. materials candidate 终态为 `material_ready`，install candidate 终态为 `succeeded`；两阶段均经 Owner/Base Gate 和正式 Receipt。
2. Receipt 的 candidate、transaction、event、source、decision 与 evidence 已反查；重启后终态引用/状态/Receipt 稳定，install 和 material-ready Receipt 各只出现一次。
3. runtime image 为 `sha256:5c727701f94b82ddc1c70b20a77fba3bfcd82acb1aaa1973879f116f8e132f1b`，用户 `10001:10001`，默认网络 `none`。
4. 该现场没有 pull/fetch、容器运行、Provider/model 调用或平台登录、上传、发送、发布。

当前 ReadModel 仍返回 `sandbox_not_ready`，根因为 `install_completed_runtime_or_auth_probe_still_required`，且模型绑定未就绪。这是受控 fail-closed，不得将“已安装”写成 Provider ready，亦不得调用 Hands、模型或生成内容候选。下一步归 OS/Provider Owner 的 runtime/auth probe 与本地 OMLX 绑定；本 Pack 无权绕过或修复。
