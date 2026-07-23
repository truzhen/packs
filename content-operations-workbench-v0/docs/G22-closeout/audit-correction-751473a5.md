# G22 审计更正：当前权威 OS SHA 复核

日期：2026-07-24
范围：仅 G22 已有 Pack worktree 的 closeout 文档；未修改 OS、公共工具或其他 Pack。

## 更正结论

旧证据使用的 OS SHA `3a97db377d1cbbf171b18ff5a98fcaf4c2e417c1` 不是当前权威 `main` SHA `751473a5dd1b0ee2965397b541ef57a72e8ce273` 的祖先，故旧 SHA 的 I04 观察不再作为 G22 当前阻断。G22 已在该权威 SHA 的只读临时副本上、使用任务专属临时 test-store 和 `127.0.0.1:18222` 重跑无真实外部副作用的验证。

## 当前 SHA 复核结果

| 范畴 | 结果 |
| --- | --- |
| execution sidecar 隔离 | 通过：空的软件注册表根，并显式关闭 Android 与 execution sidecar autostart；未发现任务 OS sidecar 子进程或持久 keyfile 加载。 |
| install / enable | 通过：从空 registry 装入 Pack、角色槽和两个 schedule。 |
| Gate / Receipt | 通过：无可信 origin 的 Gate prepare 为 403；受控本地 Owner presence 仅为 Pack lifecycle disable 签发 Gate 并得到正式 lifecycle Receipt。该本地 lifecycle Gate 不授权任何真实 Codex、OpenMontage 或平台动作。 |
| disable / uninstall / restart / reinstall | 通过：disable 后 uninstall 暂停 schedule；重启持久化；reinstall 同版本重新激活并恢复 schedule。 |
| OpenMontage 拒绝路径 | 通过：read-only preflight 为 `provider_missing`；素材权利未确认时为 422，未生成 MP4。 |
| Codex 拒绝路径 | 通过：readiness 为 `not_ready`；没有执行运行记录或 artifact。 |
| 外部平台 | 未登录、未上传、未发送、未发布。 |

## 仍然阻断的正向工作

中文母内容、渠道公开候选、人工发布包、周复盘候选和竖屏 MP4 均未生产。它们需要 Owner 为 Codex Hands / OpenMontage 分别签发逐次 Gate，并提供去敏、可用且已确认权利的证据或素材；本次 G22 不因此获得任何平台发布授权。

## 审计边界

本结论只说明当前 SHA 在本次任务隔离配置下没有启动任务侧 execution sidecar。它不声明 OS 删除了全部 Windows VM 自动探测逻辑，也不构成对任何真实 Provider 或外部平台的可用性证明。
