#!/usr/bin/env python3
"""环保 Pack installer Owner handoff 与组合 readiness 行为测试（不触网）。"""

import contextlib
import importlib.util
import io
import json
import os
import tempfile
import unittest
import urllib.parse
from pathlib import Path
from unittest import mock


REPO = Path(os.environ.get("TRUZHEN_PACKS_REPO_UNDER_TEST", Path(__file__).resolve().parents[2]))
PACK = REPO / "environmental-enforcement-pack-v0"
MANIFEST = json.loads((PACK / "manifest.json").read_text(encoding="utf-8"))
SCOPES = json.loads((PACK / "knowledge/knowledge-scopes.json").read_text(encoding="utf-8"))
INDEX = json.loads((PACK / "knowledge/knowledge-index.json").read_text(encoding="utf-8"))
PACK_REF = MANIFEST["pack_ref"]
VERSION = MANIFEST["version"]
PACK_VERSION_REF = PACK_REF + "@" + VERSION
REQUIRED_SCOPES = sorted(
    scope["scope_ref"] for scope in SCOPES["scopes"] if scope.get("required", True)
)


def load_install_module():
    spec = importlib.util.spec_from_file_location(
        "environmental_install_owner_handoff_under_test",
        PACK / "install.py",
    )
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def parse_handoff(stdout):
    lines = [
        line.split("=", 1)[1]
        for line in stdout.splitlines()
        if line.startswith("TRUZHEN_PACK_HANDOFF=")
    ]
    if len(lines) != 1:
        raise AssertionError("handoff 输出数量不为 1: %s" % stdout)
    return json.loads(lines[0])


def active_mount(scope_ref):
    suffix = scope_ref.rsplit("/", 1)[-1]
    receipt_ref = "receipt://knowledge-mount-enable/" + suffix
    return {
        "mount_ref": "knowledge_mount://" + suffix,
        "owner_id": "owner://local/default",
        "pack_ref": PACK_REF,
        "pack_version_ref": PACK_VERSION_REF,
        "scene_ref": "scene://environmental-enforcement",
        "knowledge_scope_ref": scope_ref,
        "status": "active",
        "enabled_receipt_ref": receipt_ref,
        "last_receipt_ref": receipt_ref,
    }


class FakeOS:
    def __init__(self, *, lifecycle_mode, mounts=None, missing_receipt_scope="",
                 lifecycle_body_override=None):
        self.lifecycle_mode = lifecycle_mode
        self.mounts = list(mounts or [])
        self.missing_receipt_scope = missing_receipt_scope
        self.lifecycle_body_override = lifecycle_body_override
        self.calls = []
        self.roles = set()
        self.bindings = set()
        self.knowledge_candidates = {}
        self.knowledge_batch_source_sets = []

    def lifecycle_body(self):
        if self.lifecycle_body_override is not None:
            return self.lifecycle_body_override
        if self.lifecycle_mode == "empty":
            return {"readmodel": True, "packs": []}
        entry = {
            "pack_ref": PACK_REF,
            "records": [{"version": VERSION, "state": self.lifecycle_mode}],
        }
        if self.lifecycle_mode == "enabled":
            entry["enabled_pointer"] = {"current_version": VERSION}
            entry["records"][0]["state"] = "pack_enabled_version"
        return {"readmodel": True, "packs": [entry]}

    def __call__(self, method, path, body=None):
        self.calls.append((method, path, body))
        if path == "/v3/pack-studio/lifecycle/confirm":
            raise AssertionError("installer 不得调用 lifecycle confirm")
        if path in ("/v3/base/gated-actions/prepare", "/v3/base/gated-actions/confirm"):
            raise AssertionError("installer 不得调用 Base prepare/confirm")
        if "/memory/knowledge/candidates/" in path and path.endswith("/approve"):
            raise AssertionError("installer 不得自动 approve 知识候选")

        if path.startswith("/v3/pack-studio/lifecycle/packs?"):
            query = urllib.parse.parse_qs(urllib.parse.urlsplit(path).query)
            if query != {"pack_ref": [PACK_REF]}:
                raise AssertionError("lifecycle query 必须 exact: %s" % query)
            return 200, self.lifecycle_body()
        if path == "/v3/pack-studio/canvas":
            return 200, {"engine_sync": {"synced": True}}
        if path == "/v3/pack-studio/lifecycle/draft":
            return 200, {}
        if path == "/v3/pack-studio/lifecycle/readiness":
            return 200, {"record": {"readiness_report": {"ready": True}}}
        if path == "/v3/pack-studio/lifecycle/promote":
            return 200, {}

        if path.startswith("/v3/memory/knowledge/mounts?"):
            query = urllib.parse.parse_qs(urllib.parse.urlsplit(path).query)
            expected = {
                "owner_id": ["owner://local/default"],
                "pack_ref": [PACK_REF],
                "pack_version_ref": [PACK_VERSION_REF],
            }
            for key, value in expected.items():
                if query.get(key) != value:
                    raise AssertionError("mount query 缺 exact %s: %s" % (key, query))
            selected = list(self.mounts)
            if "knowledge_scope_ref" in query:
                scope_ref = query["knowledge_scope_ref"][0]
                if query.get("status") != ["active"]:
                    raise AssertionError("required scope query 必须 status=active: %s" % query)
                selected = [
                    mount for mount in selected
                    if mount["knowledge_scope_ref"] == scope_ref
                    and mount["status"] == "active"
                ]
            return 200, {"mounts": selected}
        if path.startswith("/v3/receipts/"):
            receipt_ref = urllib.parse.unquote(path.split("/v3/receipts/", 1)[1])
            if self.missing_receipt_scope and receipt_ref.endswith(
                    self.missing_receipt_scope.rsplit("/", 1)[-1]):
                return 404, {"status": "blocked", "error": "receipt_not_found"}
            return 200, {
                "receipt_ref": receipt_ref,
                "schema_version": "v1",
                "status": "recorded",
            }

        if path == "/v3/agent-orchestration/role-packs/readmodel":
            return 200, {
                "enabled_versions": [
                    {"role_pack_id": role_ref + "@1.0.0"} for role_ref in sorted(self.roles)
                ],
            }
        if path == "/v3/agent-orchestration/role-packs/drafts":
            return 200, {"draft_id": body["draft_id"]}
        if path.endswith("/readiness-check"):
            return 200, {}
        if path.endswith("/promote-candidate"):
            return 200, {"ok": True}
        if path.endswith("/enable-candidate"):
            return 200, {"ok": True}
        if path.endswith("/enable-confirm"):
            self.roles.add(body["role_pack_id"])
            return 200, {"status": "enabled"}

        if path == "/v3/agent-orchestration/agent-slots/readmodel":
            return 200, {
                "agent_slot_bindings": [
                    {
                        "slot_ref": slot_ref,
                        "scope_ref": PACK_VERSION_REF,
                        "enabled_state": "enabled",
                    }
                    for slot_ref in sorted(self.bindings)
                ],
            }
        if path == "/v3/agent-orchestration/agent-slots/bind-candidate":
            return 200, {"ok": True, "binding_ref": "binding://" + body["slot_ref"]}
        if path == "/v3/agent-orchestration/agent-slots/confirm":
            self.bindings.add(body["binding_ref"].split("binding://", 1)[1])
            return 200, {"status": "enabled"}

        if path == "/v3/memory/knowledge/batches":
            source_refs = tuple(sorted(item["source_ref"] for item in body["source_files"]))
            self.knowledge_batch_source_sets.append(source_refs)
            candidates = []
            for source_ref in source_refs:
                candidate_ref = self.knowledge_candidates.setdefault(
                    source_ref, "candidate://knowledge/" + str(len(self.knowledge_candidates) + 1)
                )
                candidates.append({"candidate_ref": candidate_ref, "status": "pending"})
            return 200, {"candidates": candidates}
        raise AssertionError("fake OS 未覆盖：%s %s" % (method, path))


class EnvironmentalInstallOwnerHandoffTest(unittest.TestCase):
    def run_main(self, fake, state_dir):
        module = load_install_module()
        module.call = fake
        stdout = io.StringIO()
        stderr = io.StringIO()
        argv = [str(PACK / "install.py"), "--devserver-base", "http://127.0.0.1:18099"]
        with mock.patch.dict(
                os.environ, {"TRUZHEN_PACK_INSTALL_STATE_DIR": state_dir}, clear=False):
            with mock.patch("sys.argv", argv):
                with contextlib.redirect_stdout(stdout), contextlib.redirect_stderr(stderr):
                    module.main()
        self.assertEqual("", stderr.getvalue())
        return parse_handoff(stdout.getvalue())

    def test_source_forbids_installer_owned_gate_confirmation(self):
        source = (PACK / "install.py").read_text(encoding="utf-8")
        for forbidden in (
            '"/v3/pack-studio/lifecycle/confirm"',
            '"/v3/base/gated-actions/prepare"',
            '"/v3/base/gated-actions/confirm"',
            '"/v3/pack-studio/lifecycle/reactivate"',
            '"/v3/memory/knowledge/candidates/"',
            '"/v3/memory/knowledge/batches"',
            '"/v3/agent-orchestration/role-packs/enable-confirm"',
            '"/v3/agent-orchestration/agent-slots/confirm"',
            '"approve": True',
            "evidence://",
        ):
            self.assertNotIn(forbidden, source)

    def test_first_run_stages_candidate_then_hands_off(self):
        fake = FakeOS(lifecycle_mode="empty")
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("awaiting_owner_confirmation", handoff["status"])
        self.assertEqual("lifecycle_candidate_staged", handoff["reason"])
        paths = [path for _, path, _ in fake.calls]
        self.assertIn("/v3/pack-studio/lifecycle/draft", paths)
        self.assertIn("/v3/pack-studio/lifecycle/readiness", paths)
        self.assertIn("/v3/pack-studio/lifecycle/promote", paths)
        self.assertFalse(any("/agent-orchestration/" in path for path in paths))
        self.assertFalse(any("/memory/knowledge/batches" in path for path in paths))

    def test_existing_candidate_retry_does_not_restage(self):
        fake = FakeOS(lifecycle_mode="pack_spec_candidate")
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("awaiting_owner_confirmation", handoff["status"])
        self.assertEqual(
            ["GET"],
            [method for method, _, _ in fake.calls],
            "已有 candidate 时只能只读等待 Owner，不能复制 staging 写入",
        )

    def test_target_record_missing_state_fails_closed(self):
        fake = FakeOS(
            lifecycle_mode="unused",
            lifecycle_body_override={
                "readmodel": True,
                "packs": [{
                    "pack_ref": PACK_REF,
                    "records": [{"version": VERSION}],
                }],
            },
        )
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("not_ready", handoff["status"])
        self.assertEqual(
            "malformed_lifecycle_readmodel_fail_closed",
            handoff["reason"],
        )
        self.assertEqual(["GET"], [method for method, _, _ in fake.calls])

    def test_pointer_only_partial_mounts_stop_all_downstream(self):
        fake = FakeOS(
            lifecycle_mode="enabled",
            mounts=[active_mount(scope_ref) for scope_ref in REQUIRED_SCOPES[:-1]],
        )
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("not_ready", handoff["status"])
        self.assertEqual(14, handoff["readiness"]["active_scope_count"])
        self.assertTrue(any(
            reason.startswith("required_scope_not_exactly_active:")
            for reason in handoff["readiness"]["reason_codes"]
        ))
        self.assertFalse(any(
            "/agent-orchestration/" in path or "/memory/knowledge/batches" in path
            for _, path, _ in fake.calls
        ))

    def test_recovery_state_preserves_existing_audit_refs(self):
        mounts = [active_mount(scope_ref) for scope_ref in REQUIRED_SCOPES[:-1]]
        blocked_scope = REQUIRED_SCOPES[-1]
        mounts.append({
            "mount_ref": "knowledge_mount://blocked",
            "owner_id": "owner://local/default",
            "pack_ref": PACK_REF,
            "pack_version_ref": PACK_VERSION_REF,
            "knowledge_scope_ref": blocked_scope,
            "status": "blocked",
            "blocked_reason": "recovery_required",
            "last_receipt_ref": "receipt://pack-enable/failure-existing",
        })
        fake = FakeOS(lifecycle_mode="enabled", mounts=mounts)
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("recovery", handoff["status"])
        self.assertIn("receipt://pack-enable/failure-existing", handoff["audit_refs"])
        blocked = [
            row for row in handoff["readiness"]["mounts"]
            if row["knowledge_scope_ref"] == blocked_scope
        ]
        self.assertEqual("recovery_required", blocked[0]["blocked_reason"])

    def test_active_mount_with_blocked_duplicate_fails_closed(self):
        duplicate_scope = REQUIRED_SCOPES[-1]
        mounts = [active_mount(scope_ref) for scope_ref in REQUIRED_SCOPES]
        mounts.append({
            "mount_ref": "knowledge_mount://duplicate-blocked",
            "owner_id": "owner://local/default",
            "pack_ref": PACK_REF,
            "pack_version_ref": PACK_VERSION_REF,
            "knowledge_scope_ref": duplicate_scope,
            "status": "blocked",
            "blocked_reason": "recovery_required",
            "last_receipt_ref": "receipt://pack-enable/duplicate-failure",
        })
        fake = FakeOS(lifecycle_mode="enabled", mounts=mounts)
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("recovery", handoff["status"])
        self.assertIn(
            "required_scope_has_non_active_mount:" + duplicate_scope,
            handoff["readiness"]["reason_codes"],
        )
        self.assertIn("receipt://pack-enable/duplicate-failure", handoff["audit_refs"])
        self.assertFalse(any(
            "/agent-orchestration/" in path or "/memory/knowledge/batches" in path
            for _, path, _ in fake.calls
        ))

    def test_missing_formal_receipt_fails_closed(self):
        missing_scope = REQUIRED_SCOPES[3]
        fake = FakeOS(
            lifecycle_mode="enabled",
            mounts=[active_mount(scope_ref) for scope_ref in REQUIRED_SCOPES],
            missing_receipt_scope=missing_scope,
        )
        with tempfile.TemporaryDirectory() as state_dir:
            handoff = self.run_main(fake, state_dir)
        self.assertEqual("not_ready", handoff["status"])
        self.assertIn(
            "formal_receipt_lookup_failed:" + missing_scope,
            handoff["readiness"]["reason_codes"],
        )
        self.assertFalse(any(
            "/agent-orchestration/" in path or "/memory/knowledge/batches" in path
            for _, path, _ in fake.calls
        ))

    def test_ready_resume_only_hands_off_without_downstream_writes(self):
        fake = FakeOS(
            lifecycle_mode="enabled",
            mounts=[active_mount(scope_ref) for scope_ref in REQUIRED_SCOPES],
        )
        with tempfile.TemporaryDirectory() as state_dir:
            first = self.run_main(fake, state_dir)
            first_calls = list(fake.calls)
            second = self.run_main(fake, state_dir)
        self.assertEqual("awaiting_owner_confirmation", first["status"])
        self.assertEqual("downstream_owner_confirmation_required", first["reason"])
        self.assertEqual(first, second)
        self.assertEqual([], first["candidate_refs"])
        self.assertEqual(
            [
                ("role_pack_candidate_and_enable", 2),
                ("agent_slot_candidate_and_confirm", 2),
                ("knowledge_candidate_review_and_formalize", 45),
            ],
            [
                (step["step_id"], step["target_count"])
                for step in first["owner_steps"]
            ],
        )
        self.assertEqual(
            first_calls,
            fake.calls[len(first_calls):],
            "retry 必须只重读相同 readiness surfaces",
        )
        self.assertTrue(all(method == "GET" for method, _, _ in fake.calls))
        self.assertEqual({}, fake.knowledge_candidates)
        self.assertEqual([], fake.knowledge_batch_source_sets)
        self.assertEqual(set(), fake.roles)
        self.assertEqual(set(), fake.bindings)


if __name__ == "__main__":
    unittest.main()
