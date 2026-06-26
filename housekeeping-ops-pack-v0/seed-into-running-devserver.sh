#!/usr/bin/env bash
# 把「家政公司日常运营场景荚」真造进**正在运行的 devserver**，使前端「场景包管理」能看到。
#
# 走的是前端画布实际调的同一批端点（canvas 存图 → 06 同步 → lifecycle 上架），
# 不是手写裸存。落库到该 devserver 的默认 pack_studio_dev.sqlite。
#
# 前置（重要）：devserver 必须跑**本分支的新代码**——本荚含「能力节点」
# (capability.invoke)，旧二进制（无 nodeinfo 修复）会拒 UNKNOWN_STUDIO_NODE_TYPE。
#   1) 停掉旧 devserver（占用 127.0.0.1:18080 的那个）
#   2) 在本 worktree 起新的：
#        go run ./backend/cmd/devserver        # 默认 127.0.0.1:18080
#   3) 另开终端跑本脚本：
#        bash packs/housekeeping-ops-pack-v0/seed-into-running-devserver.sh
#   4) 前端刷新「场景包管理」→ 看到「家政公司日常运营 Pack」
#
# 注：本脚本只上架场景荚本身（前端可见）。要让 advice 跑火山，另需：
#   - 前端「添加云模型」填火山 base_url+api_key，或脚本调 /v3/model-gateway/cloud/save
#   - export TRUZHEN_SCENE_FLOW_ADVICE_CLOUD_PROVIDER_ID=cloud-deepseek-DeepSeek-V4-Pro（起 devserver 前）

set -euo pipefail
BASE="${TRUZHEN_DEVSERVER_BASE:-http://127.0.0.1:18080}"
FLOW_ID="housekeeping-cs-lifecycle"
PACK_REF="scene_pack://housekeeping-ops"
VERSION="1.0.0"

echo "[1/5] 画布存图 → 06 同步 (POST /v3/pack-studio/canvas)"
SYNC=$(curl -sS -X POST "$BASE/v3/pack-studio/canvas" -H 'Content-Type: application/json' -d @- <<'JSON'
{
  "flow_id": "housekeeping-cs-lifecycle",
  "title": "家政公司日常运营场景荚",
  "occ_version": 0,
  "save_source": "manual",
  "flow_spec_draft": {
    "flow_id": "housekeeping-cs-lifecycle",
    "title": "家政公司日常运营场景荚",
    "version": "1.0.0",
    "nodes": [
      {"id": "start", "type": "input.user_need", "title": "受理客户咨询"},
      {"id": "intake_inquiry", "type": "object.business_schema", "title": "建客户服务对象"},
      {"id": "consultant_advice", "type": "collaboration.advice", "title": "家政顾问建议", "slot_ref": "housekeeping_consultant"},
      {"id": "quality_challenge", "type": "collaboration.challenge", "title": "质检质询", "slot_ref": "quality_auditor"},
      {"id": "compare_gate", "type": "collaboration.compare_gate", "title": "对照确认门", "gate_policy": {"required_gate": "sovereign_gate", "pending_owner_confirmation": true}},
      {"id": "quote_capability", "type": "capability.invoke", "title": "排期报价能力"},
      {"id": "dispatch_task", "type": "flow.stage_task_candidate", "title": "派工任务", "gate_policy": {"required_gate": "sovereign_gate", "pending_owner_confirmation": true}},
      {"id": "owner_approval", "type": "policy.gate_config", "title": "Owner 派工审批", "gate_policy": {"required_gate": "sovereign_gate", "pending_owner_confirmation": true}},
      {"id": "notify_customer", "type": "gateway.communication_draft", "title": "通知客户"},
      {"id": "onsite_service", "type": "gateway.execution_intent", "title": "上门服务"},
      {"id": "close_receipt", "type": "receipt.link", "title": "服务回执"},
      {"id": "end", "type": "flow.end", "title": "完成归档"},
      {"id": "close_rejected", "type": "flow.end", "title": "退回归档"}
    ],
    "edges": [
      {"id": "e1", "source": "start", "target": "intake_inquiry"},
      {"id": "e2", "source": "intake_inquiry", "target": "consultant_advice"},
      {"id": "e3", "source": "consultant_advice", "target": "quality_challenge"},
      {"id": "e4", "source": "quality_challenge", "target": "compare_gate"},
      {"id": "e5", "source": "compare_gate", "target": "quote_capability", "condition": "approved"},
      {"id": "e6", "source": "quote_capability", "target": "dispatch_task"},
      {"id": "e7", "source": "dispatch_task", "target": "owner_approval"},
      {"id": "e8", "source": "owner_approval", "target": "notify_customer", "condition": "approved"},
      {"id": "e9", "source": "owner_approval", "target": "close_rejected", "condition": "rejected"},
      {"id": "e10", "source": "notify_customer", "target": "onsite_service"},
      {"id": "e11", "source": "onsite_service", "target": "close_receipt"},
      {"id": "e12", "source": "close_receipt", "target": "end"}
    ]
  }
}
JSON
)
echo "$SYNC" | grep -q '"synced":true' || { echo "❌ 06 同步失败（旧二进制不认能力节点？）：$SYNC"; exit 1; }
echo "    ✓ synced"

echo "[2/5] lifecycle/draft（声明 role_slots）"
curl -sS -X POST "$BASE/v3/pack-studio/lifecycle/draft" -H 'Content-Type: application/json' -d "{
  \"pack_ref\": \"$PACK_REF\", \"version\": \"$VERSION\",
  \"title\": \"家政公司日常运营 Pack\", \"template_family\": \"家政服务履约型\",
  \"flow_id\": \"$FLOW_ID\", \"idempotency_key\": \"seed-draft-$PACK_REF\", \"actor_ref\": \"owner:li\",
  \"role_slots\": [
    {\"slot_id\": \"housekeeping_consultant\", \"responsibility\": \"家政顾问：依据客户需求出服务方案候选\", \"required_role\": \"家政顾问\", \"slice_scope_policy\": \"transaction_scoped\"},
    {\"slot_id\": \"quality_auditor\", \"responsibility\": \"质检挑剔：质询方案与报价的遗漏与风险\", \"required_role\": \"质检挑剔\", \"slice_scope_policy\": \"transaction_scoped\"}
  ]
}" >/dev/null && echo "    ✓ draft"

echo "[3/5] lifecycle/readiness"
curl -sS -X POST "$BASE/v3/pack-studio/lifecycle/readiness" -H 'Content-Type: application/json' \
  -d "{\"pack_ref\": \"$PACK_REF\", \"version\": \"$VERSION\", \"actor_ref\": \"owner:li\"}" >/dev/null && echo "    ✓ readiness"

echo "[4/5] lifecycle/promote"
curl -sS -X POST "$BASE/v3/pack-studio/lifecycle/promote" -H 'Content-Type: application/json' \
  -d "{\"pack_ref\": \"$PACK_REF\", \"version\": \"$VERSION\", \"actor_ref\": \"owner:li\"}" >/dev/null && echo "    ✓ promote"

echo "[5/5] lifecycle/confirm（真 01 Base Gate + 03 receipt）"
CONFIRM=$(curl -sS -X POST "$BASE/v3/pack-studio/lifecycle/confirm" -H 'Content-Type: application/json' -d "{
  \"pack_ref\": \"$PACK_REF\", \"version\": \"$VERSION\",
  \"idempotency_key\": \"seed-confirm-$PACK_REF\", \"owner_ref\": \"owner:li\", \"approve\": true,
  \"comment\": \"家政公司日常运营 Pack 启用\"
}")
echo "$CONFIRM" | grep -q '"status":"enabled"' || { echo "❌ 上架失败：$CONFIRM"; exit 1; }
echo "    ✓ enabled"
echo
echo "✅ 家政荚已上架到 $BASE —— 前端刷新「场景包管理」即可看到「家政公司日常运营 Pack」"
