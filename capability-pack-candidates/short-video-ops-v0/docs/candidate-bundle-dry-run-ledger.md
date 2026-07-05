# Candidate Bundle Dry-Run Ledger

> P10 台账。本文只记录候选 Pack bundle 的 dry-run 校验，不代表安装、启用、发布或真实 provider 接通。

## 当前结论

- Endpoint：`/v3/capability/studio/candidate_bundle_dry_run`
- 输入来源：已保存的 `/v3/capability/studio/candidate_bundle` 候选产物。
- 期望状态：`candidate_bundle_dry_run_blocked`
- 阻断原因：`capability_pack_loader_missing`
- 静态校验：`passed`
- 生命周期：`NoInstall=true`、`NoEnable=true`、`InstallSupported=false`、`EnableSupported=false`
- 执行边界：`NoRealExecution=true`、`FormalWrite=false`、`CandidateOnly=true`

## 必需文件清单

```text
README.md
manifest/capability-pack.json
capabilities/capabilities.json
docs/oss-evidence.md
docs/code-assistant-invocation-ledger.md
docs/code-assistant-run-request-ledger.md
docs/patch-candidate-review-ledger.md
docs/candidate-bundle-dry-run-ledger.md
```

## 禁止冒充成功

- 不使用 Domain Work Pack loader 代替独立 Capability Pack loader。
- 不调用 `install.py` / `uninstall.py`。
- 不写 enabled pack registry。
- 不生成正式 `receipt://`。
- 不执行 MoneyPrinterTurbo、Pixelle-Video、social-auto-upload。
- 不运行 Codex CLI，也不应用 PatchCandidate。

## 后续商用缺口

独立 Capability Pack loader、Capability Pack lifecycle、安装 / 启用 Gate、正式 Receipt、provider readiness 和前端 E2E 录像仍未完成。完成这些之前，本候选集只能保持 `blocked / provider_missing / not_ready` 的诚实口径。
