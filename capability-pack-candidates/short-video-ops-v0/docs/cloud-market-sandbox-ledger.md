# P18 云市场 Sandbox 证据台账

> 状态：`pending_authorization`
> 本文件只记录 P18 sandbox 证据。当前未获 P18 授权，未修改 `truzhen-cloud`，未创建 sandbox listing，未生成 License / Entitlement，未下载或安装。

## Sandbox 边界

| 字段 | 当前值 |
|---|---|
| listing_ref | `pending_authorization` |
| artifact_ref | `capability-pack-candidate-set://short-video-ops-v0` |
| entitlement_status | `pending_authorization` |
| download_status | `pending_authorization` |
| install_preflight_status | `pending_authorization` |
| production_release | false |
| real_payment | false |
| production_license | false |

## 待填证据

P18 获授权并执行后，必须补齐：

evidence_contract：`docs/p18-cloud-market-sandbox-evidence-contract.json`

## 非充分证据

以下内容只能作为候选说明或待执行材料，不能标记 P18 完成，不能请求商用签字，也不能解锁 go-live：

- packs 内的 listing draft，但没有 `truzhen-cloud` sandbox listing receipt。
- sandbox runbook / 执行说明，但没有 cloud sandbox receipt。
- License / Entitlement 文案或字段样例，但没有 cloud truth 和 entitlement receipt。
- 人工填写的订单 / 授权状态，但没有 cloud receipt。
- 生产发布声称，但没有 cloud receipt；本批也禁止生产发布。

### 机器证据 ID 覆盖表

| evidence_id | 当前状态 | 待填证据位置 |
|---|---|---|
| `sandbox_listing_created_in_cloud` | `pending_authorization` | truzhen-cloud sandbox receipt |
| `sandbox_entitlement_download_install_preflight_receipts` | `pending_authorization` | entitlement / download / install preflight receipts |
| `gui_cloud_sandbox_status_visible` | `pending_authorization` | GUI test / screenshot |
| `no_real_payment` | `pending_authorization` | cloud sandbox no-production-payment evidence |
| `no_production_release_or_license` | `pending_authorization` | cloud no-production-release/license evidence |
| `packs_do_not_store_cloud_truth` | `pending_authorization` | packs artifact scan |

- sandbox PackListingDraft ref。
- sandbox review result。
- sandbox License / Entitlement ref。
- download preflight response。
- install preflight response。
- blocked reason for missing entitlement。
- 生产支付 / 生产发布未发生的证据。

任一证据缺失时，P18 不得标为完成。
