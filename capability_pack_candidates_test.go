package packs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShortVideoCandidateSetUsesSingleLifecycleStage(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))

	if got := requireString(t, candidateSet, "candidate_set_ref"); got != "capability-pack-candidate-set://short-video-ops-v0" {
		t.Fatalf("candidate_set_ref = %s", got)
	}
	requireBool(t, candidateSet, "candidate_only", true)
	requireBool(t, candidateSet, "non_formal", true)

	stage := requireString(t, candidateSet, "lifecycle_status")
	if strings.Contains(stage, "->") {
		t.Fatalf("lifecycle_status must be one standard stage, got %q", stage)
	}
	requireStringIn(t, stage, "想法", "设计中", "契约已定", "已实现", "已接线", "已验收", "已发布", "已弃用")

	summary := requireObject(t, candidateSet, "slice_status_summary")
	if got := requireString(t, summary, "p3_p11"); got != "已接线" {
		t.Fatalf("slice_status_summary.p3_p11 = %s, want 已接线", got)
	}
	if got := requireString(t, summary, "p12_p18"); got != "设计中" {
		t.Fatalf("slice_status_summary.p12_p18 = %s, want 设计中", got)
	}
	if got := requireString(t, summary, "commercial_go_live"); got != "not_commercial_ready" {
		t.Fatalf("slice_status_summary.commercial_go_live = %s, want not_commercial_ready", got)
	}
}

func TestShortVideoCapabilityPacksStayCandidateOnlyAndProviderFree(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0", "capability-packs")
	files := []string{
		"short-video-draft-generation.capability-pack.json",
		"short-video-composition-orchestration.capability-pack.json",
		"short-video-social-publish-draft.capability-pack.json",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			doc := readJSON(t, filepath.Join(base, file))
			if got := requireString(t, doc, "pack_type"); got != "Capability Pack" {
				t.Fatalf("pack_type = %s, want Capability Pack", got)
			}
			requireBool(t, doc, "candidate_only", true)
			requireBool(t, doc, "non_formal", true)
			requireStringIn(t, requireString(t, doc, "lifecycle_status"), "想法", "设计中", "契约已定", "已实现", "已接线", "已验收", "已发布", "已弃用")

			integrationPlan := requireObject(t, doc, "integration_plan")
			requireBool(t, integrationPlan, "candidate_only", true)
			requireBool(t, integrationPlan, "non_formal", true)
			requireBool(t, integrationPlan, "no_direct_execution", true)

			requireCapabilityPackHasProviderRequirement(t, doc)
			requireStringSliceContains(t, asStringSlice(t, doc["forbidden"]), "third_party_repo_execution_from_pack")
		})
	}
}

func TestShortVideoFutureCommercialSlicesRemainAuthorizationGated(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))

	for _, key := range []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
		"p16_controlled_code_assistant_run",
		"p17_provider_adapter_candidate",
		"p18_cloud_market_sandbox",
	} {
		slice := requireObject(t, candidateSet, key)
		status := requireString(t, slice, "status")
		if !strings.Contains(status, "pending_authorization") {
			t.Fatalf("%s.status = %s, want pending_authorization", key, status)
		}
	}

	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	if got := requireString(t, goLive, "current_go_live_status"); got != "not_commercial_ready" {
		t.Fatalf("current_go_live_status = %s, want not_commercial_ready", got)
	}
	requireBool(t, goLive, "no_completion_claim", true)
}

func TestShortVideoCommercialSlicesHaveSpecsAuthCardsAndEvidenceLedgers(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, key := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(key, func(t *testing.T) {
			slice := requireObject(t, candidateSet, key)
			requireExistingPath(t, requireString(t, slice, "implementation_spec"), "")
			requireExistingPath(t, requireString(t, slice, "cross_repo_authorization_card"), "")
			requireExistingPath(t, requireString(t, slice, "evidence_ledger"), base)
		})
	}
}

func TestShortVideoCommercialGoLiveHasMachineCheckableForbiddenActions(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	checks := requireObject(t, goLive, "forbidden_action_terminal_checks")

	for _, key := range []string{
		"raw_token_cookie_password_saved",
		"third_party_oss_vendored_into_packs",
		"third_party_oss_executed_without_authorization",
		"social_login_happened",
		"social_upload_or_publish_happened",
		"candidate_only_written_as_enabled",
		"cloud_listing_truth_stored_in_packs",
		"production_payment_captured",
		"production_license_issued",
	} {
		t.Run(key, func(t *testing.T) {
			check := requireObject(t, checks, key)
			if got := requireString(t, check, "current_result"); got != "pending" {
				t.Fatalf("%s.current_result = %s, want pending before commercial verification", key, got)
			}
			requireBool(t, check, "required_final_value", false)
			if got := requireString(t, check, "evidence_required"); got == "" {
				t.Fatalf("%s.evidence_required missing", key)
			}
		})
	}
}

func TestShortVideoCommercialForbiddenActionCoverageVerifierCoversAllTerminalChecks(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	verifierPath := requireString(t, goLive, "forbidden_action_coverage_verifier")
	requireExistingPath(t, verifierPath, base)
	verifier := readJSON(t, filepath.Join(base, verifierPath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	readinessVerifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	if got := requireString(t, verifier, "verifier_ref"); got != "commercial-forbidden-action-coverage-verifier://short-video-ops-v0" {
		t.Fatalf("verifier_ref = %s, want commercial-forbidden-action-coverage-verifier://short-video-ops-v0", got)
	}
	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	if got := requireString(t, verifier, "coverage_status"); got != "blocked_pending_forbidden_action_terminal_evidence" {
		t.Fatalf("coverage_status = %s, want blocked_pending_forbidden_action_terminal_evidence", got)
	}

	for key, want := range map[string]string{
		"candidate_set_ref":                       requireString(t, candidateSet, "candidate_set_ref"),
		"source_candidate_set":                    "candidate-set.json#commercial_go_live_evidence_package.forbidden_action_terminal_checks",
		"source_evidence_contract_index":          requireString(t, goLive, "evidence_contract_index") + "#forbidden_action_terminal_checks",
		"source_commercial_readiness_verifier":    requireString(t, goLive, "commercial_readiness_verifier") + "#current_blockers.forbidden_action_terminal_checks",
		"source_commercial_go_no_go_gate":         requireString(t, goLive, "commercial_go_no_go_gate") + "#forbidden_action_terminal_gates",
		"source_goal_completion_evidence_map":     requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map") + "#active_goal_requirements.forbidden_actions_proven_false",
		"required_terminal_check_source":          "commercial_go_live_evidence_package.forbidden_action_terminal_checks",
		"required_terminal_final_value":           "false",
		"required_current_result_before_pass":     "pending",
		"required_result_before_commercial_ready": "all_terminal_checks_false",
	} {
		if got := requireString(t, verifier, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), "forbidden_action_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "owner_signoff_gate")["blocking_required_before_owner_signoff"]), "forbidden_action_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "commercial_signoff_matrix")["non_sufficient_evidence"]), "forbidden_action_coverage_verifier_pending")

	if got := requireString(t, readinessVerifier, "forbidden_action_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness forbidden_action_coverage_verifier = %s, want %s", got, verifierPath)
	}
	readinessBlocker := requireObject(t, requireObject(t, readinessVerifier, "current_blockers"), "forbidden_action_coverage_verifier")
	if got := requireString(t, readinessBlocker, "source_forbidden_action_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness blocker source_forbidden_action_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireBool(t, readinessBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, readinessBlocker["blocked_by"]), "forbidden_action_coverage_verifier_pending")

	if got := requireString(t, goNoGo, "source_forbidden_action_coverage_verifier"); got != verifierPath {
		t.Fatalf("go-no-go source_forbidden_action_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "forbidden_action_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "forbidden_action_coverage_verifier_pending")

	if got := requireString(t, goalMap, "source_forbidden_action_coverage_verifier"); got != verifierPath {
		t.Fatalf("goal map source_forbidden_action_coverage_verifier = %s, want %s", got, verifierPath)
	}
	findObjectByString(t, asObjectSlice(t, goalMap["active_goal_requirements"]), "requirement_id", "forbidden_action_coverage_verified")
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "forbidden_action_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "forbidden_action_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), "forbidden_action_coverage_verifier_pending")

	expectedKeys := []string{
		"raw_token_cookie_password_saved",
		"third_party_oss_vendored_into_packs",
		"third_party_oss_executed_without_authorization",
		"social_login_happened",
		"social_upload_or_publish_happened",
		"candidate_only_written_as_enabled",
		"cloud_listing_truth_stored_in_packs",
		"production_payment_captured",
		"production_license_issued",
	}
	goLiveChecks := requireObject(t, goLive, "forbidden_action_terminal_checks")
	indexChecks := requireObject(t, index, "forbidden_action_terminal_checks")
	readinessChecks := requireObject(t, requireObject(t, readinessVerifier, "current_blockers"), "forbidden_action_terminal_checks")
	goNoGoGates := asObjectSlice(t, goNoGo["forbidden_action_terminal_gates"])
	coverageChecks := asObjectSlice(t, verifier["coverage_checks"])
	terminalKeys := asStringSlice(t, verifier["terminal_check_keys"])
	if len(goLiveChecks) != len(expectedKeys) || len(indexChecks) != len(expectedKeys) || len(coverageChecks) != len(expectedKeys) || len(terminalKeys) != len(expectedKeys) {
		t.Fatalf("forbidden action coverage counts goLive=%d index=%d checks=%d terminalKeys=%d want %d", len(goLiveChecks), len(indexChecks), len(coverageChecks), len(terminalKeys), len(expectedKeys))
	}

	for _, key := range expectedKeys {
		t.Run(key, func(t *testing.T) {
			requireStringSliceContains(t, terminalKeys, key)
			source := requireObject(t, goLiveChecks, key)
			indexCheck := requireObject(t, indexChecks, key)
			readinessCheck := requireObject(t, readinessChecks, key)
			goNoGoGate := findObjectByString(t, goNoGoGates, "check_key", key)
			coverageCheck := findObjectByString(t, coverageChecks, "check_key", key)

			for label, check := range map[string]map[string]any{
				"index":     indexCheck,
				"readiness": readinessCheck,
				"go_no_go":  goNoGoGate,
			} {
				if got, want := requireString(t, check, "current_result"), requireString(t, source, "current_result"); got != want {
					t.Fatalf("%s current_result = %s, want %s", label, got, want)
				}
				if got, want := requireBoolValue(t, check, "required_final_value"), requireBoolValue(t, source, "required_final_value"); got != want {
					t.Fatalf("%s required_final_value = %v, want %v", label, got, want)
				}
			}
			if got, want := requireString(t, indexCheck, "evidence_required"), requireString(t, source, "evidence_required"); got != want {
				t.Fatalf("index evidence_required = %s, want %s", got, want)
			}
			if got, want := requireString(t, goNoGoGate, "evidence_required"), requireString(t, source, "evidence_required"); got != want {
				t.Fatalf("go-no-go evidence_required = %s, want %s", got, want)
			}
			requireBool(t, readinessCheck, "can_count_toward_commercial_ready", false)

			if got := requireString(t, coverageCheck, "current_result"); got != "pending" {
				t.Fatalf("%s coverage current_result = %s, want pending", key, got)
			}
			requireBool(t, coverageCheck, "required_final_value", false)
			if got := requireString(t, coverageCheck, "coverage_result"); got != "blocked_pending_terminal_evidence" {
				t.Fatalf("%s coverage_result = %s, want blocked_pending_terminal_evidence", key, got)
			}
			for _, flag := range []string{
				"candidate_set_has_check",
				"evidence_index_has_check",
				"readiness_verifier_has_check",
				"go_no_go_has_terminal_gate",
				"candidate_set_required_final_value_false",
				"go_no_go_required_final_value_false",
			} {
				requireBool(t, coverageCheck, flag, true)
			}
		})
	}

	completionGate := requireObject(t, verifier, "completion_gate")
	requireBool(t, completionGate, "can_pass_coverage_verifier", false)
	for _, proof := range []string{
		"all_forbidden_action_sources_present",
		"all_forbidden_action_terminal_gates_present",
		"all_required_final_values_false",
		"all_current_results_false",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_commercial_ready"]), proof)
	}
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "forbidden_action_coverage_verifier_without_terminal_false_evidence")
}

func TestShortVideoCommercialSlicesDeclareExecutionDependencies(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	expected := map[string][]string{
		"p12_safe_lifecycle_sample":            {},
		"p13_gui_lifecycle_panel":              {"p12_safe_lifecycle_sample"},
		"p15_gui_walkthrough_three_candidates": {"p13_gui_lifecycle_panel"},
		"p16_controlled_code_assistant_run":    {"p12_safe_lifecycle_sample", "p13_gui_lifecycle_panel", "p15_gui_walkthrough_three_candidates"},
		"p17_provider_adapter_candidate":       {"p16_controlled_code_assistant_run"},
		"p18_cloud_market_sandbox":             {"p12_safe_lifecycle_sample", "p13_gui_lifecycle_panel", "p15_gui_walkthrough_three_candidates", "p16_controlled_code_assistant_run", "p17_provider_adapter_candidate"},
	}

	for sliceKey, expectedDeps := range expected {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			gotDeps := []string{}
			if raw, ok := slice["depends_on"]; ok {
				gotDeps = asStringSlice(t, raw)
			}
			if len(gotDeps) != len(expectedDeps) {
				t.Fatalf("%s depends_on = %v, want %v", sliceKey, gotDeps, expectedDeps)
			}
			for _, expectedDep := range expectedDeps {
				requireStringSliceContains(t, gotDeps, expectedDep)
			}
		})
	}
}

func TestShortVideoCommercialGoLiveBlocksOwnerSignoffUntilEvidenceComplete(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gate := requireObject(t, goLive, "owner_signoff_gate")

	requireBool(t, gate, "can_request_owner_signoff", false)
	if got := requireString(t, gate, "blocked_reason"); got == "" {
		t.Fatalf("owner_signoff_gate.blocked_reason missing")
	}

	for _, key := range asStringSlice(t, goLive["required_slices"]) {
		requireStringSliceContains(t, asStringSlice(t, gate["blocking_required_slices"]), key)
	}
	for _, key := range []string{
		"raw_token_cookie_password_saved",
		"third_party_oss_executed_without_authorization",
		"social_upload_or_publish_happened",
		"candidate_only_written_as_enabled",
		"production_payment_captured",
		"production_license_issued",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["blocking_terminal_checks"]), key)
	}
}

func TestShortVideoCommercialOwnerSignoffRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	writebackGate := requireObject(t, goNoGo, "evidence_writeback_gate")

	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	for label, ownerGate := range map[string]map[string]any{
		"candidate_set":           requireObject(t, goLive, "owner_signoff_gate"),
		"evidence_contract_index": requireObject(t, index, "owner_signoff_gate"),
	} {
		if got := requireString(t, ownerGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
			t.Fatalf("%s owner_signoff_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", label, got)
		}
		requireBool(t, ownerGate, "evidence_writeback_gate_must_pass", true)
		if got := requireString(t, ownerGate, "blocking_evidence_writeback_gate"); got != "evidence_writeback_incomplete" {
			t.Fatalf("%s owner_signoff_gate.blocking_evidence_writeback_gate = %s, want evidence_writeback_incomplete", label, got)
		}
		requireStringSliceContains(t, asStringSlice(t, ownerGate["blocking_required_before_owner_signoff"]), "all_evidence_writebacks_completed_and_verified")
	}

	for label, signoffMatrix := range map[string]map[string]any{
		"candidate_set":           requireObject(t, goLive, "commercial_signoff_matrix"),
		"evidence_contract_index": requireObject(t, index, "commercial_signoff_matrix"),
	} {
		if got := requireString(t, signoffMatrix, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
			t.Fatalf("%s commercial_signoff_matrix.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", label, got)
		}
		requireBool(t, signoffMatrix, "evidence_writeback_gate_must_pass", true)
		if !strings.Contains(requireString(t, signoffMatrix, "required_before_owner_signoff"), "all_evidence_writebacks_completed_and_verified") {
			t.Fatalf("%s commercial_signoff_matrix.required_before_owner_signoff must require all_evidence_writebacks_completed_and_verified", label)
		}
		requireStringSliceContains(t, asStringSlice(t, signoffMatrix["non_sufficient_evidence"]), "evidence_writeback_incomplete")
	}
}

func TestShortVideoCommercialOwnerSignoffRequiresOwnerBaseGateReceipts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	for label, ownerGate := range map[string]map[string]any{
		"candidate_set":           requireObject(t, goLive, "owner_signoff_gate"),
		"evidence_contract_index": requireObject(t, index, "owner_signoff_gate"),
	} {
		if got := requireString(t, ownerGate, "required_owner_base_gate_source"); got != "Owner + Base Gate + Gateway + Receipt" {
			t.Fatalf("%s owner_signoff_gate.required_owner_base_gate_source = %s, want Owner + Base Gate + Gateway + Receipt", label, got)
		}
		requireBool(t, ownerGate, "owner_base_gate_receipts_must_exist", true)
		if got := requireString(t, ownerGate, "blocking_owner_base_gate"); got != "owner_base_gate_receipts_missing" {
			t.Fatalf("%s owner_signoff_gate.blocking_owner_base_gate = %s, want owner_base_gate_receipts_missing", label, got)
		}
		requireStringSliceContains(t, asStringSlice(t, ownerGate["blocking_required_before_owner_signoff"]), "owner_base_gate_receipts_bound")
	}

	for label, signoffMatrix := range map[string]map[string]any{
		"candidate_set":           requireObject(t, goLive, "commercial_signoff_matrix"),
		"evidence_contract_index": requireObject(t, index, "commercial_signoff_matrix"),
	} {
		if got := requireString(t, signoffMatrix, "required_owner_base_gate_source"); got != "Owner + Base Gate + Gateway + Receipt" {
			t.Fatalf("%s commercial_signoff_matrix.required_owner_base_gate_source = %s, want Owner + Base Gate + Gateway + Receipt", label, got)
		}
		requireBool(t, signoffMatrix, "owner_base_gate_receipts_must_exist", true)
		if !strings.Contains(requireString(t, signoffMatrix, "required_before_owner_signoff"), "owner_base_gate_receipts_bound") {
			t.Fatalf("%s commercial_signoff_matrix.required_before_owner_signoff must require owner_base_gate_receipts_bound", label)
		}
		requireStringSliceContains(t, asStringSlice(t, signoffMatrix["non_sufficient_evidence"]), "owner_base_gate_receipts_missing")
	}
}

func TestShortVideoOwnerSignoffGatePointsToNextAuthorizationCard(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gate := requireObject(t, goLive, "owner_signoff_gate")

	nextSliceKey := requireString(t, gate, "required_next_authorization")
	nextSlice := requireObject(t, candidateSet, nextSliceKey)
	if status := requireString(t, nextSlice, "status"); !strings.Contains(status, "pending_authorization") {
		t.Fatalf("%s.status = %s, want pending_authorization", nextSliceKey, status)
	}

	nextAuthorizationCard := requireString(t, gate, "next_authorization_card")
	if want := requireString(t, nextSlice, "cross_repo_authorization_card"); nextAuthorizationCard != want {
		t.Fatalf("owner_signoff_gate.next_authorization_card = %s, want %s", nextAuthorizationCard, want)
	}
	requireExistingPath(t, nextAuthorizationCard, "")
}

func TestShortVideoCommercialGoLiveTracksPerSliceOwnerAuthorizationAndEvidenceSignoff(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gate := requireObject(t, goLive, "owner_signoff_gate")

	requiredSlices := asStringSlice(t, goLive["required_slices"])
	signoffMatrix := requireObject(t, goLive, "commercial_signoff_matrix")
	requireCommercialSignoffMatrix(t, candidateSet, signoffMatrix, requiredSlices, gate)

	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexMatrix := requireObject(t, index, "commercial_signoff_matrix")
	requireCommercialSignoffMatrix(t, candidateSet, indexMatrix, requiredSlices, gate)

	completionGate := requireObject(t, index, "completion_gate")
	if got := requireString(t, completionGate, "required_signoff_matrix_source"); got != "commercial_go_live_evidence_package.commercial_signoff_matrix" {
		t.Fatalf("completion_gate.required_signoff_matrix_source = %s, want commercial_go_live_evidence_package.commercial_signoff_matrix", got)
	}
}

func TestShortVideoP12SafeLifecycleSampleHasMachineReadableEvidenceContract(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")

	contractPath := requireString(t, p12, "evidence_contract")
	contractFile := contractPath
	if !filepath.IsAbs(contractFile) {
		contractFile = filepath.Join(base, contractPath)
	}
	requireExistingPath(t, contractFile, "")
	contract := readJSON(t, contractFile)

	requireBool(t, contract, "candidate_only", true)
	requireBool(t, contract, "non_formal", true)
	if got := requireString(t, contract, "slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if status := requireString(t, contract, "status"); !strings.Contains(status, "pending_authorization") {
		t.Fatalf("status = %s, want pending_authorization before P12 execution", status)
	}

	for _, repo := range asStringSlice(t, p12["target_repositories"]) {
		requireStringSliceContains(t, asStringSlice(t, contract["required_repositories"]), repo)
	}
	for _, endpoint := range asStringSlice(t, p12["lifecycle_endpoints"]) {
		requireStringSliceContains(t, asStringSlice(t, contract["lifecycle_endpoints"]), endpoint)
	}

	requireEvidenceItems(t, contract, "required_backend_evidence", []string{
		"server_derived_draft",
		"readiness_issues_or_ready_state",
		"promote_gate_preserved",
		"confirm_requires_owner_base_gate",
		"receipt_ref_bound_to_confirm",
		"enabled_pointer_visible_after_confirm",
	})
	requireEvidenceItems(t, contract, "required_frontend_evidence", []string{
		"gui_safe_sample_lifecycle_controls",
		"confirm_disabled_until_gate_receipt",
		"enabled_pointer_receipt_after_confirm",
		"third_party_oss_not_runnable",
	})
	requireEvidenceItems(t, contract, "required_forbidden_action_checks", []string{
		"codex_cli_not_run",
		"third_party_oss_not_executed",
		"social_login_upload_not_happened",
		"raw_secret_not_saved",
		"contracts_not_changed",
	})

	gate := requireObject(t, contract, "completion_gate")
	requireBool(t, gate, "can_mark_p12_complete", false)
	if got := requireString(t, gate, "required_status_before_completion"); got == "" {
		t.Fatalf("completion_gate.required_status_before_completion missing")
	}
}

func TestShortVideoP12AuthorizationScopeIsMachineCheckableBeforeCrossRepoWork(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")

	scopePath := requireString(t, p12, "authorization_scope_contract")
	requireExistingPath(t, scopePath, base)
	scope := readJSON(t, filepath.Join(base, scopePath))

	requireBool(t, scope, "candidate_only", true)
	requireBool(t, scope, "non_formal", true)
	if got := requireString(t, scope, "slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if got := requireString(t, scope, "status"); got != "pending_owner_authorization" {
		t.Fatalf("status = %s, want pending_owner_authorization", got)
	}

	authorizedRepos := asObjectSlice(t, scope["authorized_repositories"])
	for _, repo := range asStringSlice(t, p12["target_repositories"]) {
		entry := findObjectByString(t, authorizedRepos, "repo", repo)
		requireStringSliceContains(t, asStringSlice(t, entry["allowed_actions"]), "modify")
		requireStringSliceContains(t, asStringSlice(t, entry["allowed_actions"]), "test")
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, scope["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"run_real_codex_cli",
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
		"production_cloud_market_action",
	} {
		requireStringSliceContains(t, asStringSlice(t, scope["forbidden_actions"]), action)
	}
	for _, phrasePart := range []string{
		"授权按 P12 授权卡",
		"truzhenos 与 truzhen-client-web-desktop",
		"不改 contracts、software、cloud",
		"不运行真实 Codex CLI",
		"不执行第三方 OSS",
		"不登录或上传社媒",
		"只使用基座内置安全样本或 fixture",
	} {
		requireStringSliceContains(t, asStringSlice(t, scope["required_authorization_phrase_contains"]), phrasePart)
	}

	evidence := requireObject(t, scope, "current_authorization_evidence")
	if got := requireString(t, evidence, "status"); got != "missing" {
		t.Fatalf("current_authorization_evidence.status = %s, want missing before Owner authorization", got)
	}
	requireBool(t, evidence, "owner_thread_quote_required", true)
}

func TestShortVideoP12AuthorizationEvidenceIntakeContractIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")

	scopePath := requireString(t, p12, "authorization_scope_contract")
	scope := readJSON(t, filepath.Join(base, scopePath))
	scopeEvidence := requireObject(t, scope, "current_authorization_evidence")
	intakePath := requireString(t, p12, "authorization_evidence_intake_contract")
	if got := requireString(t, scopeEvidence, "evidence_intake_contract"); got != intakePath {
		t.Fatalf("current_authorization_evidence.evidence_intake_contract = %s, want %s", got, intakePath)
	}
	requireExistingPath(t, intakePath, base)
	intake := readJSON(t, filepath.Join(base, intakePath))

	requireBool(t, intake, "candidate_only", true)
	requireBool(t, intake, "non_formal", true)
	if got := requireString(t, intake, "slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if got := requireString(t, intake, "status"); got != "missing_owner_authorization" {
		t.Fatalf("status = %s, want missing_owner_authorization", got)
	}
	if got := requireString(t, intake, "authorization_scope_contract"); got != scopePath {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, scopePath)
	}
	if got := requireString(t, intake, "accepted_authorization_card"); got != requireString(t, scopeEvidence, "accepted_authorization_card") {
		t.Fatalf("accepted_authorization_card = %s, want %s", got, requireString(t, scopeEvidence, "accepted_authorization_card"))
	}

	for _, phrasePart := range asStringSlice(t, scope["required_authorization_phrase_contains"]) {
		requireStringSliceContains(t, asStringSlice(t, intake["required_owner_quote_contains"]), phrasePart)
	}
	for _, repo := range asStringSlice(t, p12["target_repositories"]) {
		requireStringSliceContains(t, asStringSlice(t, intake["authorized_repositories"]), repo)
	}
	for _, repo := range asStringSlice(t, scope["disallowed_repositories"]) {
		requireStringSliceContains(t, asStringSlice(t, intake["disallowed_repositories"]), repo)
	}
	for _, action := range asStringSlice(t, scope["forbidden_actions"]) {
		requireStringSliceContains(t, asStringSlice(t, intake["forbidden_actions"]), action)
	}

	gate := requireObject(t, intake, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); got != "owner_authorization_recorded_and_scope_matched" {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_scope_matched", got)
	}
	requireBool(t, gate, "must_record_before_cross_repo_work", true)
}

func TestShortVideoP12AuthorizationIntakeRejectsP11AuthorizationCardMismatch(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	intake := readJSON(t, filepath.Join(base, requireString(t, p12, "authorization_evidence_intake_contract")))

	attempt := findObjectByString(t, asObjectSlice(t, intake["rejected_authorization_attempts"]), "attempt_ref", "owner-message://2026-07-05-p11-card-mismatch-for-p12")
	if got := requireString(t, attempt, "received_authorization_card"); got != "P11" {
		t.Fatalf("received_authorization_card = %s, want P11", got)
	}
	if got := requireString(t, attempt, "target_slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("target_slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if got := requireString(t, attempt, "rejection_reason"); got != "authorization_card_mismatch_p11_not_p12" {
		t.Fatalf("rejection_reason = %s, want authorization_card_mismatch_p11_not_p12", got)
	}
	requireBool(t, attempt, "required_p12_phrase_matched", false)
	requireBool(t, attempt, "can_start_cross_repo_work_after_attempt", false)
	requireStringSliceContains(t, asStringSlice(t, attempt["missing_required_phrase_parts"]), "授权按 P12 授权卡")
	requireStringSliceContains(t, asStringSlice(t, attempt["owner_quote_contains"]), "授权按 P11 授权卡")

	gate := requireObject(t, intake, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, intake, "status"); got != "missing_owner_authorization" {
		t.Fatalf("status = %s, want missing_owner_authorization", got)
	}
}

func TestShortVideoCommercialGatesConsumeRejectedP12AuthorizationAttempt(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	intakePath := requireString(t, p12, "authorization_evidence_intake_contract")
	intake := readJSON(t, filepath.Join(base, intakePath))
	rejected := findObjectByString(t, asObjectSlice(t, intake["rejected_authorization_attempts"]), "attempt_ref", "owner-message://2026-07-05-p11-card-mismatch-for-p12")

	docs := map[string]struct {
		path     string
		boolKey  string
		wantBool bool
	}{
		"cross repo execution queue": {
			path:     requireString(t, goLive, "cross_repo_execution_queue"),
			boolKey:  "can_start_cross_repo_execution_after_attempt",
			wantBool: false,
		},
		"commercial readiness verifier": {
			path:     requireString(t, goLive, "commercial_readiness_verifier"),
			boolKey:  "can_mark_commercial_ready_after_attempt",
			wantBool: false,
		},
		"commercial go/no-go gate": {
			path:     requireString(t, goLive, "commercial_go_no_go_gate"),
			boolKey:  "can_request_owner_go_live_signoff_after_attempt",
			wantBool: false,
		},
		"goal completion evidence map": {
			path:     requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map"),
			boolKey:  "can_mark_goal_complete_after_attempt",
			wantBool: false,
		},
	}
	for label, spec := range docs {
		doc := readJSON(t, filepath.Join(base, spec.path))
		attempt := requireObject(t, doc, "latest_rejected_authorization_attempt")
		if got := requireString(t, attempt, "source_authorization_evidence_intake_contract"); got != intakePath {
			t.Fatalf("%s source_authorization_evidence_intake_contract = %s, want %s", label, got, intakePath)
		}
		if got, want := requireString(t, attempt, "attempt_ref"), requireString(t, rejected, "attempt_ref"); got != want {
			t.Fatalf("%s attempt_ref = %s, want %s", label, got, want)
		}
		if got, want := requireString(t, attempt, "rejection_reason"), requireString(t, rejected, "rejection_reason"); got != want {
			t.Fatalf("%s rejection_reason = %s, want %s", label, got, want)
		}
		if got, want := requireString(t, attempt, "target_slice_key"), requireString(t, rejected, "target_slice_key"); got != want {
			t.Fatalf("%s target_slice_key = %s, want %s", label, got, want)
		}
		requireBool(t, attempt, "required_p12_phrase_matched", false)
		requireBool(t, attempt, spec.boolKey, spec.wantBool)
	}
}

func TestShortVideoCommercialSourceDocsConsumeRejectedP12AuthorizationAttempt(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	intakePath := requireString(t, p12, "authorization_evidence_intake_contract")
	intake := readJSON(t, filepath.Join(base, intakePath))
	rejected := findObjectByString(t, asObjectSlice(t, intake["rejected_authorization_attempts"]), "attempt_ref", "owner-message://2026-07-05-p11-card-mismatch-for-p12")

	docs := map[string]string{
		"commercial evidence index":            requireString(t, goLive, "evidence_contract_index"),
		"commercial current state audit":       requireString(t, goLive, "current_state_audit"),
		"commercial improvement backlog":       requireString(t, goLive, "improvement_backlog"),
		"frontend backend acceptance contract": requireString(t, goLive, "frontend_backend_acceptance_contract"),
		"frontend backend handoff runbook":     requireString(t, goLive, "frontend_backend_handoff_runbook"),
	}
	for label, path := range docs {
		doc := readJSON(t, filepath.Join(base, path))
		attempt := requireObject(t, doc, "latest_rejected_authorization_attempt")
		if got := requireString(t, attempt, "source_authorization_evidence_intake_contract"); got != intakePath {
			t.Fatalf("%s source_authorization_evidence_intake_contract = %s, want %s", label, got, intakePath)
		}
		if got, want := requireString(t, attempt, "attempt_ref"), requireString(t, rejected, "attempt_ref"); got != want {
			t.Fatalf("%s attempt_ref = %s, want %s", label, got, want)
		}
		if got, want := requireString(t, attempt, "rejection_reason"), requireString(t, rejected, "rejection_reason"); got != want {
			t.Fatalf("%s rejection_reason = %s, want %s", label, got, want)
		}
		if got, want := requireString(t, attempt, "target_slice_key"), requireString(t, rejected, "target_slice_key"); got != want {
			t.Fatalf("%s target_slice_key = %s, want %s", label, got, want)
		}
		requireBool(t, attempt, "required_p12_phrase_matched", false)
		requireBool(t, attempt, "can_treat_as_owner_authorization", false)
		requireStringSliceContains(t, asStringSlice(t, attempt["missing_required_phrase_parts"]), "授权按 P12 授权卡")
	}
}

func TestShortVideoCommercialSourceDocsExposeAuthorizationAttemptCoverageVerifierPath(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	verifierPath := requireString(t, goLive, "authorization_attempt_coverage_verifier")
	requireExistingPath(t, verifierPath, base)

	docs := map[string]string{
		"commercial evidence index":            requireString(t, goLive, "evidence_contract_index"),
		"commercial current state audit":       requireString(t, goLive, "current_state_audit"),
		"commercial improvement backlog":       requireString(t, goLive, "improvement_backlog"),
		"frontend backend acceptance contract": requireString(t, goLive, "frontend_backend_acceptance_contract"),
		"frontend backend handoff runbook":     requireString(t, goLive, "frontend_backend_handoff_runbook"),
	}
	for label, path := range docs {
		doc := readJSON(t, filepath.Join(base, path))
		if got := requireString(t, doc, "authorization_attempt_coverage_verifier"); got != verifierPath {
			t.Fatalf("%s authorization_attempt_coverage_verifier = %s, want %s", label, got, verifierPath)
		}
	}
}

func TestShortVideoCommercialAuthorizationAttemptCoverageVerifierCoversAllCommercialGateDocs(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	intakePath := requireString(t, p12, "authorization_evidence_intake_contract")
	intake := readJSON(t, filepath.Join(base, intakePath))
	rejected := findObjectByString(t, asObjectSlice(t, intake["rejected_authorization_attempts"]), "attempt_ref", "owner-message://2026-07-05-p11-card-mismatch-for-p12")

	verifierPath := requireString(t, goLive, "authorization_attempt_coverage_verifier")
	requireExistingPath(t, verifierPath, base)
	verifier := readJSON(t, filepath.Join(base, verifierPath))
	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	if got := requireString(t, readiness, "authorization_attempt_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness authorization_attempt_coverage_verifier = %s, want %s", got, verifierPath)
	}
	blocker := requireObject(t, requireObject(t, readiness, "current_blockers"), "authorization_attempt_coverage_verifier")
	if got := requireString(t, blocker, "source_authorization_attempt_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness blocker source_authorization_attempt_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireBool(t, blocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, blocker["blocked_by"]), "rejected_p12_authorization_attempt_recorded")

	if got := requireString(t, goNoGo, "source_authorization_attempt_coverage_verifier"); got != verifierPath {
		t.Fatalf("go-no-go source_authorization_attempt_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "authorization_attempt_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_pending")

	if got := requireString(t, goalMap, "source_authorization_attempt_coverage_verifier"); got != verifierPath {
		t.Fatalf("goal map source_authorization_attempt_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "authorization_attempt_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "authorization_attempt_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_pending")

	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	if got := requireString(t, verifier, "coverage_status"); got != "blocked_rejected_p12_authorization_attempt_not_authorization" {
		t.Fatalf("coverage_status = %s, want blocked_rejected_p12_authorization_attempt_not_authorization", got)
	}
	if got := requireString(t, verifier, "source_authorization_evidence_intake_contract"); got != intakePath {
		t.Fatalf("source_authorization_evidence_intake_contract = %s, want %s", got, intakePath)
	}
	if got, want := requireString(t, verifier, "expected_attempt_ref"), requireString(t, rejected, "attempt_ref"); got != want {
		t.Fatalf("expected_attempt_ref = %s, want %s", got, want)
	}
	if got, want := requireString(t, verifier, "expected_rejection_reason"), requireString(t, rejected, "rejection_reason"); got != want {
		t.Fatalf("expected_rejection_reason = %s, want %s", got, want)
	}

	requiredDocs := asStringSlice(t, verifier["required_commercial_gate_docs"])
	coverageChecks := asObjectSlice(t, verifier["coverage_checks"])
	if len(requiredDocs) != len(coverageChecks) {
		t.Fatalf("required_commercial_gate_docs len = %d, coverage_checks len = %d", len(requiredDocs), len(coverageChecks))
	}
	for _, docPath := range requiredDocs {
		check := findObjectByString(t, coverageChecks, "doc_path", docPath)
		requireBool(t, check, "latest_rejected_authorization_attempt_present", true)
		requireBool(t, check, "attempt_ref_matches", true)
		requireBool(t, check, "rejection_reason_matches", true)
		requireBool(t, check, "can_treat_as_owner_authorization", false)
		requireBool(t, check, "authorization_attempt_coverage_verifier_path_present", true)
		requireBool(t, check, "authorization_attempt_coverage_verifier_path_matches", true)

		doc := readJSON(t, filepath.Join(base, docPath))
		docVerifierPath, ok := doc["authorization_attempt_coverage_verifier"].(string)
		if !ok || docVerifierPath == "" {
			docVerifierPath, ok = doc["source_authorization_attempt_coverage_verifier"].(string)
		}
		if !ok || docVerifierPath == "" {
			t.Fatalf("%s missing authorization attempt coverage verifier path", docPath)
		}
		if docVerifierPath != verifierPath {
			t.Fatalf("%s authorization attempt coverage verifier path = %s, want %s", docPath, docVerifierPath, verifierPath)
		}
		attempt := requireObject(t, doc, "latest_rejected_authorization_attempt")
		if got, want := requireString(t, attempt, "attempt_ref"), requireString(t, rejected, "attempt_ref"); got != want {
			t.Fatalf("%s attempt_ref = %s, want %s", docPath, got, want)
		}
		if got, want := requireString(t, attempt, "rejection_reason"), requireString(t, rejected, "rejection_reason"); got != want {
			t.Fatalf("%s rejection_reason = %s, want %s", docPath, got, want)
		}
		requireBool(t, attempt, "can_treat_as_owner_authorization", false)
	}
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_commercial_gate_docs_reference_authorization_attempt_coverage_verifier")
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "rejected_p12_authorization_attempt_is_not_owner_authorization")
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_path_not_independently_verified")
}

func TestShortVideoCommercialNextAuthorizationStartGuardLocksP12BeforeCrossRepoWork(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	scope := readJSON(t, filepath.Join(base, requireString(t, p12, "authorization_scope_contract")))
	intake := readJSON(t, filepath.Join(base, requireString(t, p12, "authorization_evidence_intake_contract")))

	guardPath := requireString(t, goLive, "next_authorization_start_guard")
	requireExistingPath(t, guardPath, base)
	guard := readJSON(t, filepath.Join(base, guardPath))

	if got := requireString(t, guard, "guard_ref"); got != "commercial-next-authorization-start-guard://short-video-ops-v0/p12-safe-lifecycle-sample" {
		t.Fatalf("guard_ref = %s, want commercial-next-authorization-start-guard://short-video-ops-v0/p12-safe-lifecycle-sample", got)
	}
	requireBool(t, guard, "candidate_only", true)
	requireBool(t, guard, "non_formal", true)
	requireBool(t, guard, "can_start_cross_repo_execution", false)
	if got := requireString(t, guard, "guard_status"); got != "blocked_missing_exact_p12_owner_authorization" {
		t.Fatalf("guard_status = %s, want blocked_missing_exact_p12_owner_authorization", got)
	}
	if got := requireString(t, guard, "next_required_slice_key"); got != requireString(t, goLive, "next_required_authorization") {
		t.Fatalf("next_required_slice_key = %s, want %s", got, requireString(t, goLive, "next_required_authorization"))
	}
	if got := requireString(t, guard, "source_authorization_scope_contract"); got != requireString(t, p12, "authorization_scope_contract") {
		t.Fatalf("source_authorization_scope_contract = %s, want %s", got, requireString(t, p12, "authorization_scope_contract"))
	}
	if got := requireString(t, guard, "source_authorization_evidence_intake_contract"); got != requireString(t, p12, "authorization_evidence_intake_contract") {
		t.Fatalf("source_authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p12, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, guard, "accepted_authorization_card"); got != requireString(t, intake, "accepted_authorization_card") {
		t.Fatalf("accepted_authorization_card = %s, want %s", got, requireString(t, intake, "accepted_authorization_card"))
	}
	requireStringSlicesEqual(t, asStringSlice(t, guard["required_owner_quote_contains"]), asStringSlice(t, intake["required_owner_quote_contains"]), "required_owner_quote_contains")
	requireStringSlicesEqual(t, asStringSlice(t, guard["authorized_repositories"]), asStringSlice(t, intake["authorized_repositories"]), "authorized_repositories")
	requireStringSlicesEqual(t, asStringSlice(t, guard["disallowed_repositories"]), asStringSlice(t, intake["disallowed_repositories"]), "disallowed_repositories")
	requireStringSlicesEqual(t, asStringSlice(t, guard["forbidden_actions"]), asStringSlice(t, scope["forbidden_actions"]), "forbidden_actions")
	for _, required := range []string{
		"owner_quote_contains_all_required_p12_parts",
		"authorization_card_matches_P12",
		"authorization_attempt_coverage_verifier_passed",
		"p12_pre_run_gate_open",
	} {
		requireStringSliceContains(t, asStringSlice(t, guard["required_before_guard_open"]), required)
	}
	for _, blocker := range []string{
		"owner_authorization_missing",
		"authorization_card_mismatch_p11_not_p12",
		"p12_safe_fixture_scope_not_authorized",
	} {
		requireStringSliceContains(t, asStringSlice(t, guard["current_blockers"]), blocker)
	}
	if got := requireString(t, guard, "authorization_attempt_coverage_verifier"); got != requireString(t, goLive, "authorization_attempt_coverage_verifier") {
		t.Fatalf("authorization_attempt_coverage_verifier = %s, want %s", got, requireString(t, goLive, "authorization_attempt_coverage_verifier"))
	}

	requiredDocs := asStringSlice(t, guard["required_commercial_gate_docs"])
	coverageChecks := asObjectSlice(t, guard["coverage_checks"])
	if len(requiredDocs) != len(coverageChecks) {
		t.Fatalf("required_commercial_gate_docs len = %d, coverage_checks len = %d", len(requiredDocs), len(coverageChecks))
	}
	for _, docPath := range requiredDocs {
		check := findObjectByString(t, coverageChecks, "doc_path", docPath)
		requireBool(t, check, "next_authorization_start_guard_path_present", true)
		requireBool(t, check, "next_authorization_start_guard_path_matches", true)
		requireBool(t, check, "can_start_cross_repo_execution", false)

		doc := readJSON(t, filepath.Join(base, docPath))
		docGuardPath, ok := doc["next_authorization_start_guard"].(string)
		if !ok || docGuardPath == "" {
			docGuardPath, ok = doc["source_next_authorization_start_guard"].(string)
		}
		if !ok || docGuardPath == "" {
			t.Fatalf("%s missing next authorization start guard path", docPath)
		}
		if docGuardPath != guardPath {
			t.Fatalf("%s next authorization start guard path = %s, want %s", docPath, docGuardPath, guardPath)
		}
	}

	docs := map[string]string{
		"commercial evidence index":      requireString(t, goLive, "evidence_contract_index"),
		"current state audit":            requireString(t, goLive, "current_state_audit"),
		"commercial improvement backlog": requireString(t, goLive, "improvement_backlog"),
		"frontend backend acceptance":    requireString(t, goLive, "frontend_backend_acceptance_contract"),
		"frontend backend handoff":       requireString(t, goLive, "frontend_backend_handoff_runbook"),
		"commercial readiness verifier":  requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial go/no-go gate":       requireString(t, goLive, "commercial_go_no_go_gate"),
		"cross repo execution queue":     requireString(t, goLive, "cross_repo_execution_queue"),
	}
	for label, path := range docs {
		doc := readJSON(t, filepath.Join(base, path))
		if got := requireString(t, doc, "next_authorization_start_guard"); got != guardPath {
			t.Fatalf("%s next_authorization_start_guard = %s, want %s", label, got, guardPath)
		}
	}

	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))
	if got := requireString(t, goalMap, "source_next_authorization_start_guard"); got != guardPath {
		t.Fatalf("goal map source_next_authorization_start_guard = %s, want %s", got, guardPath)
	}
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "next_authorization_start_guard_closed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "next_authorization_start_guard_open")
	requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), "next_authorization_start_guard_pending")
}

func TestShortVideoCommercialNextAuthorizationStartGuardBlocksCommercialSignoff(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	guardPath := requireString(t, goLive, "next_authorization_start_guard")

	requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), "next_authorization_start_guard_pending")
	ownerGate := requireObject(t, goLive, "owner_signoff_gate")
	requireStringSliceContains(t, asStringSlice(t, ownerGate["blocking_required_before_owner_signoff"]), "next_authorization_start_guard_open")
	signoffMatrix := requireObject(t, goLive, "commercial_signoff_matrix")
	if got := requireString(t, signoffMatrix, "required_before_owner_signoff"); !strings.Contains(got, "next_authorization_start_guard_open") {
		t.Fatalf("commercial_signoff_matrix.required_before_owner_signoff missing next_authorization_start_guard_open: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, signoffMatrix["non_sufficient_evidence"]), "next_authorization_start_guard_pending")

	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	if got := requireString(t, readiness, "next_authorization_start_guard"); got != guardPath {
		t.Fatalf("commercial readiness next_authorization_start_guard = %s, want %s", got, guardPath)
	}
	readinessBlocker := requireObject(t, requireObject(t, readiness, "current_blockers"), "next_authorization_start_guard")
	if got := requireString(t, readinessBlocker, "source_next_authorization_start_guard"); got != guardPath {
		t.Fatalf("readiness blocker source_next_authorization_start_guard = %s, want %s", got, guardPath)
	}
	if got := requireString(t, readinessBlocker, "guard_status"); got != "blocked_missing_exact_p12_owner_authorization" {
		t.Fatalf("readiness blocker guard_status = %s, want blocked_missing_exact_p12_owner_authorization", got)
	}
	requireBool(t, readinessBlocker, "can_count_toward_commercial_ready", false)
	for _, blocker := range []string{
		"next_authorization_start_guard_closed",
		"owner_authorization_missing",
		"authorization_card_mismatch_p11_not_p12",
	} {
		requireStringSliceContains(t, asStringSlice(t, readinessBlocker["blocked_by"]), blocker)
	}
	requireStringSliceContains(t, asStringSlice(t, readinessBlocker["required_before_commercial_ready"]), "next_authorization_start_guard_open")
	requireStringSliceContains(t, asStringSlice(t, readiness["required_before_commercial_ready"]), "next_authorization_start_guard_open")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "next_authorization_start_guard_pending")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	if got := requireString(t, goNoGo, "next_authorization_start_guard"); got != guardPath {
		t.Fatalf("go/no-go next_authorization_start_guard = %s, want %s", got, guardPath)
	}
	completionRule := requireObject(t, goNoGo, "completion_rule")
	requireBool(t, completionRule, "next_authorization_start_guard_must_open", true)
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_go_live_signoff"]), "next_authorization_start_guard_open")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "next_authorization_start_guard_pending")
}

func TestShortVideoCommercialCompletionGatesConsumeAuthorizationGuards(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))

	indexCompletionGate := requireObject(t, index, "completion_gate")
	if got := requireString(t, indexCompletionGate, "required_authorization_attempt_coverage_verifier_source"); got != "commercial_go_live_evidence_package.authorization_attempt_coverage_verifier" {
		t.Fatalf("index completion_gate.required_authorization_attempt_coverage_verifier_source = %s, want commercial_go_live_evidence_package.authorization_attempt_coverage_verifier", got)
	}
	requireBool(t, indexCompletionGate, "authorization_attempt_coverage_verifier_must_pass", true)
	if got := requireString(t, indexCompletionGate, "required_next_authorization_start_guard_source"); got != "commercial_go_live_evidence_package.next_authorization_start_guard" {
		t.Fatalf("index completion_gate.required_next_authorization_start_guard_source = %s, want commercial_go_live_evidence_package.next_authorization_start_guard", got)
	}
	requireBool(t, indexCompletionGate, "next_authorization_start_guard_must_open", true)
	requiredStatus := requireString(t, indexCompletionGate, "required_status_before_completion")
	for _, proof := range []string{
		"authorization_attempt_coverage_verifier_passed",
		"next_authorization_start_guard_open",
	} {
		if !strings.Contains(requiredStatus, proof) {
			t.Fatalf("index completion_gate.required_status_before_completion missing %s: %s", proof, requiredStatus)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, indexCompletionGate["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, indexCompletionGate["non_sufficient_evidence"]), "next_authorization_start_guard_pending")

	queueCompletionGate := requireObject(t, queue, "completion_gate")
	if got := requireString(t, queueCompletionGate, "required_next_authorization_start_guard_source"); got != "commercial_go_live_evidence_package.next_authorization_start_guard" {
		t.Fatalf("queue completion_gate.required_next_authorization_start_guard_source = %s, want commercial_go_live_evidence_package.next_authorization_start_guard", got)
	}
	requireBool(t, queueCompletionGate, "next_authorization_start_guard_must_open", true)
	if got := requireString(t, queueCompletionGate, "required_before_first_cross_repo_run"); !strings.Contains(got, "next_authorization_start_guard_open") {
		t.Fatalf("queue completion_gate.required_before_first_cross_repo_run missing next_authorization_start_guard_open: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, queueCompletionGate["non_sufficient_evidence"]), "next_authorization_start_guard_pending")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "next_authorization_start_guard_pending")
}

func TestShortVideoP12ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p12, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p12, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p12, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p12, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p12, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p12, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p12, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP12 := findObjectByString(t, indexSlices, "slice_key", "p12_safe_lifecycle_sample")
	if got := requireString(t, indexP12, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p12 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_scope_matched") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_scope_matched", got)
	}

	for _, repo := range []string{"truzhenos", "truzhen-client-web-desktop"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"run_real_codex_cli",
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"GOWORK=off go test ./backend/tests/capability -run 'TestCapabilityPackLifecycleSafeFixtureDraftIsServerDerived' -count=1",
		"npm run smoke:frontend-shell",
		"TRUZHENOS_BACKEND_ROOT=repo_ref://truzhenos/main GOWORK=off npm run smoke:frontend-behavior",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP12ExecutionReadinessGateConsumesAuthorizationGuards(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	readinessPath := requireString(t, p12, "execution_readiness_package")
	readiness := readJSON(t, filepath.Join(base, readinessPath))
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))

	readinessGate := requireObject(t, readiness, "cross_repo_work_gate")
	if got := requireString(t, readinessGate, "required_authorization_attempt_coverage_verifier_source"); got != "commercial_go_live_evidence_package.authorization_attempt_coverage_verifier" {
		t.Fatalf("readiness cross_repo_work_gate.required_authorization_attempt_coverage_verifier_source = %s, want commercial_go_live_evidence_package.authorization_attempt_coverage_verifier", got)
	}
	requireBool(t, readinessGate, "authorization_attempt_coverage_verifier_must_pass", true)
	if got := requireString(t, readinessGate, "required_next_authorization_start_guard_source"); got != "commercial_go_live_evidence_package.next_authorization_start_guard" {
		t.Fatalf("readiness cross_repo_work_gate.required_next_authorization_start_guard_source = %s, want commercial_go_live_evidence_package.next_authorization_start_guard", got)
	}
	requireBool(t, readinessGate, "next_authorization_start_guard_must_open", true)
	requiredStatus := requireString(t, readinessGate, "required_status_before_cross_repo_work")
	for _, proof := range []string{
		"owner_authorization_recorded_and_scope_matched",
		"authorization_attempt_coverage_verifier_passed",
		"next_authorization_start_guard_open",
	} {
		if !strings.Contains(requiredStatus, proof) {
			t.Fatalf("readiness cross_repo_work_gate.required_status_before_cross_repo_work missing %s: %s", proof, requiredStatus)
		}
	}

	p12Entry := findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "slice_key", "p12_safe_lifecycle_sample")
	queueGate := requireObject(t, p12Entry, "cross_repo_work_gate")
	for _, key := range []string{
		"required_authorization_attempt_coverage_verifier_source",
		"required_next_authorization_start_guard_source",
		"required_status_before_cross_repo_work",
	} {
		if got, want := requireString(t, queueGate, key), requireString(t, readinessGate, key); got != want {
			t.Fatalf("queue p12 cross_repo_work_gate.%s = %s, want %s", key, got, want)
		}
	}
	requireBool(t, queueGate, "authorization_attempt_coverage_verifier_must_pass", true)
	requireBool(t, queueGate, "next_authorization_start_guard_must_open", true)

	writebackSummary := requireObject(t, p12Entry, "evidence_writeback_plan_summary")
	for _, blocker := range []string{
		"authorization_attempt_coverage_verifier_pending",
		"next_authorization_start_guard_closed",
	} {
		requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), blocker)
	}
}

func TestShortVideoP12StartGatesBlockOnOpenPackStudioIssueLedger(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")

	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")
	issueLedger := readJSON(t, filepath.Join(base, issueLedgerPath))
	if got := requireString(t, issueLedger, "issue_ledger_status"); got != "blocked_pending_issue_resolution_evidence" {
		t.Fatalf("issue_ledger_status = %s, want blocked_pending_issue_resolution_evidence", got)
	}

	preRunGate := readJSON(t, filepath.Join(base, requireString(t, p12, "pre_run_gate")))
	if got := requireString(t, preRunGate, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("p12 pre-run gate pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	requireStringSliceContains(t, asStringSlice(t, preRunGate["required_before_start"]), "pack_studio_issue_ledger_closed")
	requireStringSliceContains(t, asStringSlice(t, preRunGate["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
	issueLedgerCheck := findObjectByString(t, asObjectSlice(t, preRunGate["pre_run_checks"]), "check_id", "pack_studio_issue_ledger_closed")
	requireBool(t, issueLedgerCheck, "satisfied", false)
	requireBool(t, issueLedgerCheck, "required_before_start", true)
	requireBool(t, issueLedgerCheck, "failure_blocks_start", true)
	if got := requireString(t, issueLedgerCheck, "evidence_source"); got != issueLedgerPath {
		t.Fatalf("pack_studio_issue_ledger_closed evidence_source = %s, want %s", got, issueLedgerPath)
	}

	readiness := readJSON(t, filepath.Join(base, requireString(t, p12, "execution_readiness_package")))
	if got := requireString(t, readiness, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("p12 readiness pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	readinessGate := requireObject(t, readiness, "cross_repo_work_gate")
	if got := requireString(t, readinessGate, "required_pack_studio_issue_ledger_source"); got != "candidate_set.pack_studio_issue_ledger" {
		t.Fatalf("readiness cross_repo_work_gate.required_pack_studio_issue_ledger_source = %s, want candidate_set.pack_studio_issue_ledger", got)
	}
	requireBool(t, readinessGate, "pack_studio_issue_ledger_must_close", true)
	if got := requireString(t, readinessGate, "required_status_before_cross_repo_work"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("readiness cross_repo_work_gate.required_status_before_cross_repo_work missing pack_studio_issue_ledger_closed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))
	if got := requireString(t, queue, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("cross-repo queue pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	queueCompletionGate := requireObject(t, queue, "completion_gate")
	if got := requireString(t, queueCompletionGate, "required_pack_studio_issue_ledger_source"); got != "candidate_set.pack_studio_issue_ledger" {
		t.Fatalf("queue completion_gate.required_pack_studio_issue_ledger_source = %s, want candidate_set.pack_studio_issue_ledger", got)
	}
	requireBool(t, queueCompletionGate, "pack_studio_issue_ledger_must_close", true)
	if got := requireString(t, queueCompletionGate, "required_before_first_cross_repo_run"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("queue completion_gate.required_before_first_cross_repo_run missing pack_studio_issue_ledger_closed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, queueCompletionGate["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

	p12Entry := findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "slice_key", "p12_safe_lifecycle_sample")
	queueGate := requireObject(t, p12Entry, "cross_repo_work_gate")
	if got := requireString(t, queueGate, "required_pack_studio_issue_ledger_source"); got != requireString(t, readinessGate, "required_pack_studio_issue_ledger_source") {
		t.Fatalf("queue p12 required_pack_studio_issue_ledger_source = %s, want %s", got, requireString(t, readinessGate, "required_pack_studio_issue_ledger_source"))
	}
	requireBool(t, queueGate, "pack_studio_issue_ledger_must_close", true)
	if got := requireString(t, queueGate, "required_status_before_cross_repo_work"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("queue p12 required_status_before_cross_repo_work missing pack_studio_issue_ledger_closed: %s", got)
	}
	writebackSummary := requireObject(t, p12Entry, "evidence_writeback_plan_summary")
	requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "pack_studio_issue_ledger_open")
}

func TestShortVideoFutureExecutionReadinessGatesConsumeAuthorizationGuards(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))
	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")

	for _, spec := range []struct {
		candidateKey string
		sliceKey     string
		priorProof   string
	}{
		{
			candidateKey: "p13_gui_lifecycle_panel",
			sliceKey:     "p13_gui_lifecycle_panel",
			priorProof:   "p12_evidence_complete",
		},
		{
			candidateKey: "p15_gui_walkthrough_three_candidates",
			sliceKey:     "p15_gui_walkthrough_three_candidates",
			priorProof:   "p13_evidence_complete",
		},
		{
			candidateKey: "p16_controlled_code_assistant_run",
			sliceKey:     "p16_controlled_code_assistant_run",
			priorProof:   "p15_evidence_complete",
		},
		{
			candidateKey: "p17_provider_adapter_candidate",
			sliceKey:     "p17_provider_adapter_candidate",
			priorProof:   "p16_evidence_complete",
		},
		{
			candidateKey: "p18_cloud_market_sandbox",
			sliceKey:     "p18_cloud_market_sandbox",
			priorProof:   "p17_evidence_complete",
		},
	} {
		t.Run(spec.sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, spec.candidateKey)
			readiness := readJSON(t, filepath.Join(base, requireString(t, slice, "execution_readiness_package")))
			readinessGate := requireObject(t, readiness, "cross_repo_work_gate")
			if got := requireString(t, readiness, "pack_studio_issue_ledger"); got != issueLedgerPath {
				t.Fatalf("%s readiness pack_studio_issue_ledger = %s, want %s", spec.sliceKey, got, issueLedgerPath)
			}

			if got := requireString(t, readinessGate, "required_authorization_attempt_coverage_verifier_source"); got != "commercial_go_live_evidence_package.authorization_attempt_coverage_verifier" {
				t.Fatalf("readiness cross_repo_work_gate.required_authorization_attempt_coverage_verifier_source = %s, want commercial_go_live_evidence_package.authorization_attempt_coverage_verifier", got)
			}
			requireBool(t, readinessGate, "authorization_attempt_coverage_verifier_must_pass", true)
			if got := requireString(t, readinessGate, "required_next_authorization_start_guard_source"); got != "commercial_go_live_evidence_package.next_authorization_start_guard" {
				t.Fatalf("readiness cross_repo_work_gate.required_next_authorization_start_guard_source = %s, want commercial_go_live_evidence_package.next_authorization_start_guard", got)
			}
			requireBool(t, readinessGate, "next_authorization_start_guard_must_open", true)
			if got := requireString(t, readinessGate, "required_pack_studio_issue_ledger_source"); got != "candidate_set.pack_studio_issue_ledger" {
				t.Fatalf("readiness cross_repo_work_gate.required_pack_studio_issue_ledger_source = %s, want candidate_set.pack_studio_issue_ledger", got)
			}
			requireBool(t, readinessGate, "pack_studio_issue_ledger_must_close", true)
			requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

			requiredStatus := requireString(t, readinessGate, "required_status_before_cross_repo_work")
			for _, proof := range []string{
				"owner_authorization_recorded",
				spec.priorProof,
				"authorization_attempt_coverage_verifier_passed",
				"next_authorization_start_guard_open",
				"pack_studio_issue_ledger_closed",
			} {
				if !strings.Contains(requiredStatus, proof) {
					t.Fatalf("readiness cross_repo_work_gate.required_status_before_cross_repo_work missing %s: %s", proof, requiredStatus)
				}
			}

			queueEntry := findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "slice_key", spec.sliceKey)
			queueGate := requireObject(t, queueEntry, "cross_repo_work_gate")
			for _, key := range []string{
				"required_authorization_attempt_coverage_verifier_source",
				"required_next_authorization_start_guard_source",
				"required_pack_studio_issue_ledger_source",
				"required_status_before_cross_repo_work",
			} {
				if got, want := requireString(t, queueGate, key), requireString(t, readinessGate, key); got != want {
					t.Fatalf("queue %s cross_repo_work_gate.%s = %s, want %s", spec.sliceKey, key, got, want)
				}
			}
			requireBool(t, queueGate, "authorization_attempt_coverage_verifier_must_pass", true)
			requireBool(t, queueGate, "next_authorization_start_guard_must_open", true)
			requireBool(t, queueGate, "pack_studio_issue_ledger_must_close", true)

			writebackSummary := requireObject(t, queueEntry, "evidence_writeback_plan_summary")
			requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "authorization_attempt_coverage_verifier_pending")
			requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "next_authorization_start_guard_closed")
			requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "pack_studio_issue_ledger_open")
		})
	}
}

func TestShortVideoExecutionReadinessGatesRejectAuthorizationAttemptSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))
	rejectedAttempt := requireObject(t, goLive, "latest_rejected_authorization_attempt")
	rejectedAttemptRef := requireString(t, rejectedAttempt, "attempt_ref")

	for _, spec := range []struct {
		candidateKey string
		sliceKey     string
	}{
		{"p12_safe_lifecycle_sample", "p12_safe_lifecycle_sample"},
		{"p13_gui_lifecycle_panel", "p13_gui_lifecycle_panel"},
		{"p15_gui_walkthrough_three_candidates", "p15_gui_walkthrough_three_candidates"},
		{"p16_controlled_code_assistant_run", "p16_controlled_code_assistant_run"},
		{"p17_provider_adapter_candidate", "p17_provider_adapter_candidate"},
		{"p18_cloud_market_sandbox", "p18_cloud_market_sandbox"},
	} {
		t.Run(spec.sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, spec.candidateKey)
			readiness := readJSON(t, filepath.Join(base, requireString(t, slice, "execution_readiness_package")))
			readinessGate := requireObject(t, readiness, "cross_repo_work_gate")

			for label, gate := range map[string]map[string]any{
				"readiness": readinessGate,
				"queue":     requireObject(t, findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "slice_key", spec.sliceKey), "cross_repo_work_gate"),
			} {
				if got := requireString(t, gate, "required_rejected_authorization_attempt_source"); got != "commercial_go_live_evidence_package.latest_rejected_authorization_attempt" {
					t.Fatalf("%s %s required_rejected_authorization_attempt_source = %s, want commercial_go_live_evidence_package.latest_rejected_authorization_attempt", spec.sliceKey, label, got)
				}
				requireBool(t, gate, "rejected_authorization_attempt_can_satisfy_owner_authorization", false)
				requireBool(t, gate, "rejected_authorization_attempt_can_start_cross_repo_work", false)
				requireStringSliceContains(t, asStringSlice(t, gate["blocked_attempt_refs"]), rejectedAttemptRef)
				if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "no_rejected_authorization_attempt_substitution") {
					t.Fatalf("%s %s required_status_before_cross_repo_work missing no_rejected_authorization_attempt_substitution: %s", spec.sliceKey, label, got)
				}
			}

			requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "rejected_authorization_attempt_is_not_owner_authorization")
		})
	}
}

func TestShortVideoExecutionReadinessCoverageVerifierCoversRejectedAuthorizationAttemptSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	verifierPath := requireString(t, goLive, "execution_readiness_guard_coverage_verifier")
	verifier := readJSON(t, filepath.Join(base, verifierPath))

	if got := requireString(t, verifier, "required_rejected_authorization_attempt_source"); got != "commercial_go_live_evidence_package.latest_rejected_authorization_attempt" {
		t.Fatalf("required_rejected_authorization_attempt_source = %s, want commercial_go_live_evidence_package.latest_rejected_authorization_attempt", got)
	}
	for _, required := range []string{
		"all_execution_readiness_packages_reject_authorization_attempt_substitution",
		"all_cross_repo_queue_entries_reject_authorization_attempt_substitution",
		"all_queue_writeback_summaries_block_rejected_authorization_attempt_substitution",
	} {
		requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), required)
	}
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "rejected_authorization_attempt_is_not_owner_authorization")

	for _, sliceKey := range []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
		"p16_controlled_code_assistant_run",
		"p17_provider_adapter_candidate",
		"p18_cloud_market_sandbox",
	} {
		t.Run(sliceKey, func(t *testing.T) {
			check := findObjectByString(t, asObjectSlice(t, verifier["coverage_checks"]), "slice_key", sliceKey)
			requireBool(t, check, "readiness_gate_rejects_authorization_attempt_substitution", true)
			requireBool(t, check, "queue_gate_rejects_authorization_attempt_substitution", true)
			requireBool(t, check, "queue_writeback_blocks_rejected_authorization_attempt_substitution", true)
		})
	}
}

func TestShortVideoExecutionReadinessGuardCoverageVerifierCoversAllCommercialSlices(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	verifierPath := requireString(t, goLive, "execution_readiness_guard_coverage_verifier")
	requireExistingPath(t, verifierPath, base)
	verifier := readJSON(t, filepath.Join(base, verifierPath))
	readinessVerifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	if got := requireString(t, verifier, "verifier_ref"); got != "commercial-execution-readiness-guard-coverage-verifier://short-video-ops-v0" {
		t.Fatalf("verifier_ref = %s, want commercial-execution-readiness-guard-coverage-verifier://short-video-ops-v0", got)
	}
	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	if got := requireString(t, verifier, "coverage_status"); got != "blocked_pending_owner_authorization_and_execution_evidence" {
		t.Fatalf("coverage_status = %s, want blocked_pending_owner_authorization_and_execution_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                                       requireString(t, candidateSet, "candidate_set_ref"),
		"source_candidate_set":                                    "candidate-set.json#commercial_go_live_evidence_package.required_slices",
		"source_cross_repo_execution_queue":                       requireString(t, goLive, "cross_repo_execution_queue") + "#execution_entries.cross_repo_work_gate",
		"source_commercial_readiness_verifier":                    requireString(t, goLive, "commercial_readiness_verifier") + "#current_blockers.execution_readiness_guard_coverage_verifier",
		"source_commercial_go_no_go_gate":                         requireString(t, goLive, "commercial_go_no_go_gate") + "#completion_rule.required_before_go_live_signoff",
		"source_goal_completion_evidence_map":                     requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map") + "#completion_claim_policy.required_before_goal_complete",
		"source_p11_reverification_baseline":                      "candidate-set.json#p11_lifecycle_preflight.latest_reverification",
		"required_authorization_attempt_coverage_verifier_source": "commercial_go_live_evidence_package.authorization_attempt_coverage_verifier",
		"required_next_authorization_start_guard_source":          "commercial_go_live_evidence_package.next_authorization_start_guard",
		"required_pack_studio_issue_ledger_source":                "candidate_set.pack_studio_issue_ledger",
	} {
		if got := requireString(t, verifier, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	expectedSlices := []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
		"p16_controlled_code_assistant_run",
		"p17_provider_adapter_candidate",
		"p18_cloud_market_sandbox",
	}
	checks := asObjectSlice(t, verifier["coverage_checks"])
	if len(checks) != len(expectedSlices) {
		t.Fatalf("coverage_checks len = %d, want %d", len(checks), len(expectedSlices))
	}
	for _, sliceKey := range expectedSlices {
		t.Run(sliceKey, func(t *testing.T) {
			check := findObjectByString(t, checks, "slice_key", sliceKey)
			requireBool(t, check, "readiness_gate_has_authorization_attempt_coverage_verifier", true)
			requireBool(t, check, "readiness_gate_has_next_authorization_start_guard", true)
			requireBool(t, check, "readiness_gate_has_pack_studio_issue_ledger_close_gate", true)
			requireBool(t, check, "queue_gate_mirrors_readiness_gate", true)
			requireBool(t, check, "queue_writeback_blocks_authorization_attempt_coverage_verifier", true)
			requireBool(t, check, "queue_writeback_blocks_next_authorization_start_guard", true)
			requireBool(t, check, "queue_writeback_blocks_pack_studio_issue_ledger", true)
			if got := requireString(t, check, "coverage_result"); got != "blocked_pending_owner_authorization_and_execution_evidence" {
				t.Fatalf("coverage_result = %s, want blocked_pending_owner_authorization_and_execution_evidence", got)
			}
		})
	}
	p12Check := findObjectByString(t, checks, "slice_key", "p12_safe_lifecycle_sample")
	requireBool(t, p12Check, "readiness_gate_has_p11_reverification_baseline", true)
	requireBool(t, p12Check, "queue_entry_has_p11_reverification_baseline", true)
	requireBool(t, p12Check, "queue_gate_requires_p11_reverification_current", true)
	requireBool(t, p12Check, "queue_writeback_blocks_stale_p11_reverification", true)

	readinessBlocker := requireObject(t, requireObject(t, readinessVerifier, "current_blockers"), "execution_readiness_guard_coverage_verifier")
	if got := requireString(t, readinessBlocker, "source_execution_readiness_guard_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness blocker source_execution_readiness_guard_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireBool(t, readinessBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, readinessBlocker["required_before_commercial_ready"]), "execution_readiness_guard_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, readinessVerifier["required_before_commercial_ready"]), "execution_readiness_guard_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, readinessVerifier["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")

	if got := requireString(t, goNoGo, "source_execution_readiness_guard_coverage_verifier"); got != verifierPath {
		t.Fatalf("go/no-go source_execution_readiness_guard_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "execution_readiness_guard_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")

	if got := requireString(t, goalMap, "source_execution_readiness_guard_coverage_verifier"); got != verifierPath {
		t.Fatalf("goal map source_execution_readiness_guard_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "execution_readiness_guard_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "execution_readiness_guard_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
}

func TestShortVideoExecutionReadinessGuardCoverageVerifierFeedsCommercialSummaryDocs(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	verifierPath := requireString(t, goLive, "execution_readiness_guard_coverage_verifier")
	requireExistingPath(t, verifierPath, base)
	verifier := readJSON(t, filepath.Join(base, verifierPath))

	summaryDocs := []struct {
		label                        string
		path                         string
		gateKey                      string
		requiredStatusKey            string
		requireTopLevelNonSufficient bool
	}{
		{
			label:             "commercial evidence index",
			path:              requireString(t, goLive, "evidence_contract_index"),
			gateKey:           "completion_gate",
			requiredStatusKey: "required_status_before_completion",
		},
		{
			label:                        "commercial current state audit",
			path:                         requireString(t, goLive, "current_state_audit"),
			gateKey:                      "commercial_completion_gate",
			requiredStatusKey:            "required_status_before_completion",
			requireTopLevelNonSufficient: true,
		},
		{
			label:                        "commercial improvement backlog",
			path:                         requireString(t, goLive, "improvement_backlog"),
			gateKey:                      "completion_gate",
			requiredStatusKey:            "required_before_backlog_done",
			requireTopLevelNonSufficient: true,
		},
		{
			label:                        "frontend backend acceptance contract",
			path:                         requireString(t, goLive, "frontend_backend_acceptance_contract"),
			gateKey:                      "completion_gate",
			requiredStatusKey:            "required_status_before_completion",
			requireTopLevelNonSufficient: true,
		},
		{
			label:                        "frontend backend handoff runbook",
			path:                         requireString(t, goLive, "frontend_backend_handoff_runbook"),
			gateKey:                      "completion_gate",
			requiredStatusKey:            "required_status_before_acceptance",
			requireTopLevelNonSufficient: true,
		},
	}

	requiredSummaryDocs := asStringSlice(t, verifier["required_commercial_summary_docs"])
	summaryCoverageChecks := asObjectSlice(t, verifier["summary_doc_coverage_checks"])
	if len(requiredSummaryDocs) != len(summaryDocs) || len(summaryCoverageChecks) != len(summaryDocs) {
		t.Fatalf("summary doc coverage count required=%d checks=%d want %d", len(requiredSummaryDocs), len(summaryCoverageChecks), len(summaryDocs))
	}

	for _, spec := range summaryDocs {
		t.Run(spec.label, func(t *testing.T) {
			requireStringSliceContains(t, requiredSummaryDocs, spec.path)
			check := findObjectByString(t, summaryCoverageChecks, "doc_path", spec.path)
			requireBool(t, check, "doc_references_execution_readiness_guard_coverage_verifier", true)
			requireBool(t, check, "completion_gate_requires_execution_readiness_guard_coverage_verifier", true)
			requireBool(t, check, "completion_gate_blocks_pending_execution_readiness_guard_coverage_verifier", true)
			requireBool(t, check, "doc_references_pack_studio_issue_ledger", true)
			requireBool(t, check, "completion_gate_requires_pack_studio_issue_ledger_close_gate", true)
			requireBool(t, check, "completion_gate_blocks_pack_studio_issue_ledger_open", true)
			requireBool(t, check, "can_count_toward_commercial_ready", false)

			doc := readJSON(t, filepath.Join(base, spec.path))
			if got := requireString(t, doc, "execution_readiness_guard_coverage_verifier"); got != verifierPath {
				t.Fatalf("%s execution_readiness_guard_coverage_verifier = %s, want %s", spec.label, got, verifierPath)
			}
			gate := requireObject(t, doc, spec.gateKey)
			if got := requireString(t, gate, "required_execution_readiness_guard_coverage_verifier_source"); got != "commercial_go_live_evidence_package.execution_readiness_guard_coverage_verifier" {
				t.Fatalf("%s %s.required_execution_readiness_guard_coverage_verifier_source = %s, want commercial_go_live_evidence_package.execution_readiness_guard_coverage_verifier", spec.label, spec.gateKey, got)
			}
			requireBool(t, gate, "execution_readiness_guard_coverage_verifier_must_pass", true)
			if got := requireString(t, gate, spec.requiredStatusKey); !strings.Contains(got, "execution_readiness_guard_coverage_verifier_passed") {
				t.Fatalf("%s %s.%s missing execution_readiness_guard_coverage_verifier_passed: %s", spec.label, spec.gateKey, spec.requiredStatusKey, got)
			}
			requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
			if spec.requireTopLevelNonSufficient {
				requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
			}
		})
	}

	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_commercial_summary_docs_reference_execution_readiness_guard_coverage_verifier")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_commercial_summary_completion_gates_require_execution_readiness_guard_coverage_verifier")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "p12_readiness_gate_references_p11_reverification_baseline")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "p12_queue_entry_references_p11_reverification_baseline")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "p12_queue_writeback_blocks_stale_p11_reverification")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_commercial_summary_docs_reference_pack_studio_issue_ledger")
	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_commercial_summary_completion_gates_require_pack_studio_issue_ledger_close_gate")
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
}

func TestShortVideoCrossRepoExecutionQueueRequiresExecutionReadinessGuardCoverageVerifier(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	verifierPath := requireString(t, goLive, "execution_readiness_guard_coverage_verifier")
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))

	if got := requireString(t, queue, "execution_readiness_guard_coverage_verifier"); got != verifierPath {
		t.Fatalf("queue execution_readiness_guard_coverage_verifier = %s, want %s", got, verifierPath)
	}
	completionGate := requireObject(t, queue, "completion_gate")
	if got := requireString(t, completionGate, "required_execution_readiness_guard_coverage_verifier_source"); got != "commercial_go_live_evidence_package.execution_readiness_guard_coverage_verifier" {
		t.Fatalf("queue completion_gate.required_execution_readiness_guard_coverage_verifier_source = %s, want commercial_go_live_evidence_package.execution_readiness_guard_coverage_verifier", got)
	}
	requireBool(t, completionGate, "execution_readiness_guard_coverage_verifier_must_pass", true)
	if got := requireString(t, completionGate, "required_before_first_cross_repo_run"); !strings.Contains(got, "execution_readiness_guard_coverage_verifier_passed") {
		t.Fatalf("queue completion_gate.required_before_first_cross_repo_run missing execution_readiness_guard_coverage_verifier_passed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, completionGate["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "execution_readiness_guard_coverage_verifier_pending")
}

func TestShortVideoCrossRepoExecutionQueueFirstRunRequiresCurrentP11PreflightBaseline(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p11 := requireObject(t, candidateSet, "p11_lifecycle_preflight")
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))

	if got := requireString(t, queue, "source_p11_reverification_baseline"); got != "candidate-set.json#p11_lifecycle_preflight.latest_reverification" {
		t.Fatalf("queue source_p11_reverification_baseline = %s, want candidate-set.json#p11_lifecycle_preflight.latest_reverification", got)
	}
	p12Entry := findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "slice_key", "p12_safe_lifecycle_sample")
	if got := requireString(t, p12Entry, "source_p11_reverification_baseline"); got != "candidate-set.json#p11_lifecycle_preflight.latest_reverification" {
		t.Fatalf("p12 queue entry source_p11_reverification_baseline = %s, want candidate-set.json#p11_lifecycle_preflight.latest_reverification", got)
	}
	p12Baseline := requireObject(t, p12Entry, "p11_reverification_baseline")
	if got := requireString(t, p12Baseline, "verified_blocked_status"); got != requireString(t, requireObject(t, p11, "latest_reverification"), "verified_blocked_status") {
		t.Fatalf("p12 p11 baseline verified_blocked_status = %s, want %s", got, requireString(t, requireObject(t, p11, "latest_reverification"), "verified_blocked_status"))
	}
	requireBool(t, p12Baseline, "must_be_current_before_p12_execution", true)
	requireStringSliceContains(t, asStringSlice(t, p12Baseline["required_before_p12_cross_repo_work"]), "p11_lifecycle_preflight_reverification_current")

	gate := requireObject(t, p12Entry, "cross_repo_work_gate")
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "p11_lifecycle_preflight_reverification_current") {
		t.Fatalf("p12 queue gate required_status_before_cross_repo_work missing p11_lifecycle_preflight_reverification_current: %s", got)
	}
	writebackSummary := requireObject(t, p12Entry, "evidence_writeback_plan_summary")
	requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "stale_or_missing_p11_lifecycle_preflight_reverification")

	completionGate := requireObject(t, queue, "completion_gate")
	if got := requireString(t, completionGate, "required_before_first_cross_repo_run"); !strings.Contains(got, "p11_lifecycle_preflight_reverification_current") {
		t.Fatalf("queue completion gate required_before_first_cross_repo_run missing p11_lifecycle_preflight_reverification_current: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, completionGate["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
}

func TestShortVideoP12PreRunGateBlocksExecutionUntilAuthorizationAndEvidenceInputsRecorded(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p12, "pre_run_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP12 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p12_safe_lifecycle_sample")
	if got := requireString(t, indexP12, "pre_run_gate"); got != gatePath {
		t.Fatalf("index p12 pre_run_gate = %s, want %s", got, gatePath)
	}

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_start_cross_repo_run", false)
	if got := requireString(t, gate, "slice_key"); got != "p12_safe_lifecycle_sample" {
		t.Fatalf("slice_key = %s, want p12_safe_lifecycle_sample", got)
	}
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_owner_authorization" {
		t.Fatalf("gate_status = %s, want blocked_pending_owner_authorization", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                       requireString(t, candidateSet, "candidate_set_ref"),
		"authorization_scope_contract":            requireString(t, p12, "authorization_scope_contract"),
		"authorization_evidence_intake_contract":  requireString(t, p12, "authorization_evidence_intake_contract"),
		"execution_readiness_package":             requireString(t, p12, "execution_readiness_package"),
		"evidence_contract":                       requireString(t, p12, "evidence_contract"),
		"evidence_ledger":                         requireString(t, p12, "evidence_ledger"),
		"commercial_readiness_verifier":           requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                requireString(t, goLive, "commercial_go_no_go_gate"),
		"required_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	intake := readJSON(t, filepath.Join(base, requireString(t, p12, "authorization_evidence_intake_contract")))
	intakeEvidence := requireObject(t, intake, "current_authorization_evidence")
	currentState := requireObject(t, gate, "current_state")
	if got, want := requireString(t, currentState, "owner_authorization_status"), requireString(t, intakeEvidence, "status"); got != want {
		t.Fatalf("current_state.owner_authorization_status = %s, want %s", got, want)
	}
	if got := requireString(t, currentState, "readiness_status"); got != "blocked_pending_owner_authorization" {
		t.Fatalf("current_state.readiness_status = %s, want blocked_pending_owner_authorization", got)
	}

	baseline := requireObject(t, gate, "source_p11_reverification_baseline")
	requireBool(t, baseline, "must_be_current_before_p12_execution", true)
	if got := requireString(t, baseline, "verified_blocked_status"); got != "lifecycle_preflight_blocked" {
		t.Fatalf("source_p11_reverification_baseline.verified_blocked_status = %s, want lifecycle_preflight_blocked", got)
	}
	requireStringSliceContains(t, asStringSlice(t, baseline["required_before_p12_cross_repo_work"]), "p11_lifecycle_preflight_reverification_current")

	expectedChecks := map[string]bool{
		"owner_authorization_recorded":                   false,
		"owner_authorization_scope_matched":              false,
		"p11_lifecycle_preflight_reverification_current": true,
		"p11_forbidden_actions_still_false":              true,
		"execution_readiness_package_open":               false,
		"required_verification_commands_present":         true,
		"evidence_writeback_targets_bound":               true,
		"disallowed_repositories_excluded":               true,
		"forbidden_actions_bound":                        true,
		"commercial_evidence_writeback_gate_bound":       true,
		"pack_studio_issue_ledger_closed":                false,
		"no_rejected_authorization_attempt_substitution": false,
	}
	checks := asObjectSlice(t, gate["pre_run_checks"])
	if len(checks) != len(expectedChecks) {
		t.Fatalf("pre_run_checks len = %d, want %d", len(checks), len(expectedChecks))
	}
	for _, check := range checks {
		checkID := requireString(t, check, "check_id")
		wantSatisfied, ok := expectedChecks[checkID]
		if !ok {
			t.Fatalf("unexpected pre_run check_id %s", checkID)
		}
		if got := requireBoolValue(t, check, "satisfied"); got != wantSatisfied {
			t.Fatalf("%s satisfied = %v, want %v", checkID, got, wantSatisfied)
		}
		requireBool(t, check, "required_before_start", true)
		requireBool(t, check, "failure_blocks_start", true)
		if got := requireString(t, check, "evidence_source"); got == "" {
			t.Fatalf("%s evidence_source missing", checkID)
		}
		if got := requireString(t, check, "current_status"); got == "" {
			t.Fatalf("%s current_status missing", checkID)
		}
	}

	for _, repo := range []string{"truzhenos", "truzhen-client-web-desktop"} {
		requireStringSliceContains(t, asStringSlice(t, gate["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, gate["disallowed_repositories"]), repo)
	}
	for _, action := range []string{"run_real_codex_cli", "execute_third_party_oss", "social_login_or_upload", "store_raw_secret"} {
		requireStringSliceContains(t, asStringSlice(t, gate["forbidden_actions"]), action)
	}
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_start"]), "owner_authorization_recorded")
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_start"]), "execution_readiness_package_open")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "p11_preflight_only")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "pre_run_gate_candidate_only")
}

func TestShortVideoP12PreRunGateRejectsRejectedAuthorizationAttemptSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	intake := readJSON(t, filepath.Join(base, requireString(t, p12, "authorization_evidence_intake_contract")))
	rejected := findObjectByString(t, asObjectSlice(t, intake["rejected_authorization_attempts"]), "attempt_ref", "owner-message://2026-07-05-p11-card-mismatch-for-p12")

	gate := readJSON(t, filepath.Join(base, requireString(t, p12, "pre_run_gate")))
	if got := requireString(t, gate, "authorization_attempt_coverage_verifier"); got != requireString(t, goLive, "authorization_attempt_coverage_verifier") {
		t.Fatalf("authorization_attempt_coverage_verifier = %s, want %s", got, requireString(t, goLive, "authorization_attempt_coverage_verifier"))
	}
	if got := requireString(t, gate, "next_authorization_start_guard"); got != requireString(t, goLive, "next_authorization_start_guard") {
		t.Fatalf("next_authorization_start_guard = %s, want %s", got, requireString(t, goLive, "next_authorization_start_guard"))
	}
	attempt := requireObject(t, gate, "latest_rejected_authorization_attempt")
	if got, want := requireString(t, attempt, "attempt_ref"), requireString(t, rejected, "attempt_ref"); got != want {
		t.Fatalf("pre-run attempt_ref = %s, want %s", got, want)
	}
	if got, want := requireString(t, attempt, "rejection_reason"), requireString(t, rejected, "rejection_reason"); got != want {
		t.Fatalf("pre-run rejection_reason = %s, want %s", got, want)
	}
	requireBool(t, attempt, "can_treat_as_owner_authorization", false)
	requireBool(t, attempt, "can_start_cross_repo_run_after_attempt", false)

	check := findObjectByString(t, asObjectSlice(t, gate["pre_run_checks"]), "check_id", "no_rejected_authorization_attempt_substitution")
	requireBool(t, check, "satisfied", false)
	requireBool(t, check, "required_before_start", true)
	requireBool(t, check, "failure_blocks_start", true)
	if got := requireString(t, check, "evidence_source"); got != requireString(t, p12, "authorization_evidence_intake_contract")+"#rejected_authorization_attempts" {
		t.Fatalf("no_rejected_authorization_attempt_substitution evidence_source = %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_start"]), "no_rejected_authorization_attempt_substitution")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "rejected_authorization_attempt_is_not_owner_authorization")
}

func TestShortVideoP12PostRunEvidenceAcceptanceGateBlocksP12CompletionAndP13UntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	p13 := requireObject(t, candidateSet, "p13_gui_lifecycle_panel")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p12, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP12 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p12_safe_lifecycle_sample")
	if got := requireString(t, indexP12, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p12 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	p13Readiness := readJSON(t, filepath.Join(base, requireString(t, p13, "execution_readiness_package")))
	if got := requireString(t, p13Readiness, "source_p12_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p13 source_p12_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p13Readiness["required_preconditions"]), "p12_post_run_evidence_acceptance_gate_passed")

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_mark_p12_complete", false)
	requireBool(t, gate, "can_unlock_p13", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p12_execution_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p12_execution_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p12_safe_lifecycle_sample",
		"evidence_contract":                         requireString(t, p12, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p12, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p12, "execution_readiness_package"),
		"pre_run_gate":                              requireString(t, p12, "pre_run_gate"),
		"next_slice":                                "p13_gui_lifecycle_panel",
		"next_slice_execution_readiness_package":    requireString(t, p13, "execution_readiness_package"),
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p12, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_backend_evidence", "required_frontend_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_execution_evidence" {
			t.Fatalf("%s current_status = %s, want pending_execution_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p12_complete", true)
		requireBool(t, requirement, "blocks_p13", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	blockers := strings.Join(asStringSlice(t, gate["current_blockers"]), "\n")
	for _, blocker := range []string{
		"p12_cross_repo_execution_not_run",
		"owner_authorization_missing",
		"owner_base_gate_receipts_missing",
		"evidence_writeback_incomplete",
		"independent_verification_missing",
	} {
		if !strings.Contains(blockers, blocker) {
			t.Fatalf("current_blockers missing %s", blocker)
		}
	}
	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p12_complete"]), proof)
	}
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "test_logs_without_receipts")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "post_run_gate_candidate_only")
}

func TestShortVideoP12PostRunGateRejectsAuthorizationOnlyCompletion(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	gate := readJSON(t, filepath.Join(base, requireString(t, p12, "post_run_evidence_acceptance_gate")))

	requireBool(t, gate, "owner_authorization_alone_can_mark_p12_complete", false)
	requireBool(t, gate, "owner_authorization_alone_can_unlock_p13", false)
	requireBool(t, gate, "owner_authorization_alone_can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_p12_complete"]), "p12_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_p12_complete"]), "all_authoritative_execution_evidence_refs_bound")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "owner_authorization_without_execution_evidence")
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "owner_authorization_without_post_run_acceptance")
}

func TestShortVideoCommercialReadinessConsumesP12PostRunEvidenceAcceptanceGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p12, "post_run_evidence_acceptance_gate")
	postRunGate := readJSON(t, filepath.Join(base, gatePath))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	blockers := requireObject(t, verifier, "current_blockers")
	acceptanceBlocker := requireObject(t, blockers, "p12_post_run_evidence_acceptance_gate")
	if got := requireString(t, acceptanceBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p12_post_run_evidence_acceptance_gate.source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	if got := requireString(t, acceptanceBlocker, "gate_status"); got != requireString(t, postRunGate, "gate_status") {
		t.Fatalf("p12_post_run_evidence_acceptance_gate.gate_status = %s, want %s", got, requireString(t, postRunGate, "gate_status"))
	}
	requireBool(t, acceptanceBlocker, "can_mark_p12_complete", false)
	requireBool(t, acceptanceBlocker, "can_unlock_p13", false)
	requireBool(t, acceptanceBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, acceptanceBlocker["required_before_p12_complete"]), "commercial_evidence_writeback_gate_updated")
	requireStringSliceContains(t, asStringSlice(t, acceptanceBlocker["blocked_by"]), "p12_post_run_evidence_acceptance_gate_not_passed")

	p12Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p12_safe_lifecycle_sample")
	if got := requireString(t, p12Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p12 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	if got := requireString(t, p12Gate, "post_run_evidence_acceptance_status"); got != requireString(t, postRunGate, "gate_status") {
		t.Fatalf("go-no-go p12 post_run_evidence_acceptance_status = %s, want %s", got, requireString(t, postRunGate, "gate_status"))
	}
	requireStringSliceContains(t, asStringSlice(t, p12Gate["blocked_by"]), "p12_post_run_evidence_acceptance_gate_not_passed")

	completionRule := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_go_live_signoff"]), "p12_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p12_post_run_evidence_acceptance_gate_pending")

	completionGate := requireObject(t, index, "completion_gate")
	if !strings.Contains(requireString(t, completionGate, "required_status_before_completion"), "p12_post_run_evidence_acceptance_gate_passed") {
		t.Fatalf("completion_gate.required_status_before_completion must require p12_post_run_evidence_acceptance_gate_passed")
	}
}

func TestShortVideoP13ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p13 := requireObject(t, candidateSet, "p13_gui_lifecycle_panel")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p13, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization_and_p12_evidence" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization_and_p12_evidence", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p13_gui_lifecycle_panel" {
		t.Fatalf("slice_key = %s, want p13_gui_lifecycle_panel", got)
	}
	requireStringSliceContains(t, asStringSlice(t, readiness["depends_on"]), "p12_safe_lifecycle_sample")
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p13, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p13, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p13, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p13, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p13, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p13, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP13 := findObjectByString(t, indexSlices, "slice_key", "p13_gui_lifecycle_panel")
	if got := requireString(t, indexP13, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p13 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_p12_evidence_complete") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_p12_evidence_complete", got)
	}

	for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"run_real_codex_cli",
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx",
		"npm run typecheck",
		"npm run smoke:frontend-shell",
		"GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP13PostRunEvidenceAcceptanceGateBlocksP13CompletionAndP15UntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p13 := requireObject(t, candidateSet, "p13_gui_lifecycle_panel")
	p15 := requireObject(t, candidateSet, "p15_gui_walkthrough_three_candidates")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p13, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP13 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p13_gui_lifecycle_panel")
	if got := requireString(t, indexP13, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p13 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	p15Readiness := readJSON(t, filepath.Join(base, requireString(t, p15, "execution_readiness_package")))
	if got := requireString(t, p15Readiness, "source_p13_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p15 source_p13_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p15Readiness["required_preconditions"]), "p13_post_run_evidence_acceptance_gate_passed")

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_mark_p13_complete", false)
	requireBool(t, gate, "can_unlock_p15", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p13_execution_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p13_execution_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p13_gui_lifecycle_panel",
		"evidence_contract":                         requireString(t, p13, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p13, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p13, "execution_readiness_package"),
		"next_slice":                                "p15_gui_walkthrough_three_candidates",
		"next_slice_execution_readiness_package":    requireString(t, p15, "execution_readiness_package"),
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p13, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_frontend_evidence", "required_backend_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_execution_evidence" {
			t.Fatalf("%s current_status = %s, want pending_execution_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p13_complete", true)
		requireBool(t, requirement, "blocks_p15", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p13_complete"]), proof)
	}

	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	verifierBlocker := requireObject(t, requireObject(t, verifier, "current_blockers"), "p13_post_run_evidence_acceptance_gate")
	if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("verifier p13 source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), "p13_post_run_evidence_acceptance_gate_not_passed")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	p13Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p13_gui_lifecycle_panel")
	if got := requireString(t, p13Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p13 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p13Gate["blocked_by"]), "p13_post_run_evidence_acceptance_gate_not_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "p13_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p13_post_run_evidence_acceptance_gate_pending")
}

func TestShortVideoP15ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p15 := requireObject(t, candidateSet, "p15_gui_walkthrough_three_candidates")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p15, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization_and_p13_evidence" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization_and_p13_evidence", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p15_gui_walkthrough_three_candidates" {
		t.Fatalf("slice_key = %s, want p15_gui_walkthrough_three_candidates", got)
	}
	requireStringSliceContains(t, asStringSlice(t, readiness["depends_on"]), "p13_gui_lifecycle_panel")
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p15, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p15, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p15, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p15, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p15, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p15, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	for _, packRef := range asStringSlice(t, p15["candidate_pack_refs"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["candidate_pack_refs"]), packRef)
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP15 := findObjectByString(t, indexSlices, "slice_key", "p15_gui_walkthrough_three_candidates")
	if got := requireString(t, indexP15, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p15 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_p13_evidence_complete") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_p13_evidence_complete", got)
	}

	for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos", "truzhen-packs"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"run_real_codex_cli",
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
		"store_browser_cookie",
		"formal_enable_claim",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx",
		"npm run smoke:frontend-shell",
		"TRUZHENOS_BACKEND_ROOT=repo_ref://truzhenos/main GOWORK=off npm run smoke:frontend-behavior",
		"GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1",
		"go test ./... -run TestShortVideoP15ExecutionReadinessPackageIsMachineCheckable -count=1",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP15PostRunEvidenceAcceptanceGateBlocksP15CompletionAndP16UntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p15 := requireObject(t, candidateSet, "p15_gui_walkthrough_three_candidates")
	p16 := requireObject(t, candidateSet, "p16_controlled_code_assistant_run")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p15, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP15 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p15_gui_walkthrough_three_candidates")
	if got := requireString(t, indexP15, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p15 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	p16Readiness := readJSON(t, filepath.Join(base, requireString(t, p16, "execution_readiness_package")))
	if got := requireString(t, p16Readiness, "source_p15_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p16 source_p15_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p16Readiness["required_preconditions"]), "p15_post_run_evidence_acceptance_gate_passed")

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_mark_p15_complete", false)
	requireBool(t, gate, "can_unlock_p16", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p15_execution_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p15_execution_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p15_gui_walkthrough_three_candidates",
		"evidence_contract":                         requireString(t, p15, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p15, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p15, "execution_readiness_package"),
		"previous_slice":                            "p13_gui_lifecycle_panel",
		"next_slice":                                "p16_controlled_code_assistant_run",
		"next_slice_execution_readiness_package":    requireString(t, p16, "execution_readiness_package"),
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p15, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_frontend_evidence", "required_backend_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_execution_evidence" {
			t.Fatalf("%s current_status = %s, want pending_execution_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p15_complete", true)
		requireBool(t, requirement, "blocks_p16", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p15_complete"]), proof)
	}

	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	verifierBlocker := requireObject(t, requireObject(t, verifier, "current_blockers"), "p15_post_run_evidence_acceptance_gate")
	if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("verifier p15 source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), "p15_post_run_evidence_acceptance_gate_not_passed")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	p15Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p15_gui_walkthrough_three_candidates")
	if got := requireString(t, p15Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p15 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p15Gate["blocked_by"]), "p15_post_run_evidence_acceptance_gate_not_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "p15_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p15_post_run_evidence_acceptance_gate_pending")
}

func TestShortVideoP16ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p16 := requireObject(t, candidateSet, "p16_controlled_code_assistant_run")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p16, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization_and_p15_evidence" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization_and_p15_evidence", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p16_controlled_code_assistant_run" {
		t.Fatalf("slice_key = %s, want p16_controlled_code_assistant_run", got)
	}
	for _, dep := range asStringSlice(t, p16["depends_on"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["depends_on"]), dep)
	}
	for _, output := range asStringSlice(t, p16["allowed_outputs"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["allowed_outputs"]), output)
	}
	requireBool(t, readiness, "controlled_code_assistant_run_allowed_now", false)
	requireBool(t, readiness, "patch_auto_apply_allowed", false)
	requireBool(t, readiness, "requires_11_gateway", true)
	requireBool(t, readiness, "requires_owner_base_gate", true)
	requireBool(t, readiness, "requires_03_receipt", true)
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p16, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p16, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p16, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p16, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p16, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p16, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP16 := findObjectByString(t, indexSlices, "slice_key", "p16_controlled_code_assistant_run")
	if got := requireString(t, indexP16, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p16 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_p15_evidence_complete") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_p15_evidence_complete", got)
	}

	for _, repo := range []string{"truzhenos", "truzhen-client-web-desktop", "truzhen-packs"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
		"uncontrolled_codex_cli_run",
		"auto_apply_patch",
		"commit_or_push_patch_candidate",
		"provider_ready_claim",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1",
		"npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx src/components/pack-lifecycle/CodeAssistantPackPanel.test.tsx src/api/__tests__/executionLiveSmokeScript.test.ts",
		"npm run typecheck",
		"npm run smoke:frontend-shell",
		"TRUZHENOS_BACKEND_ROOT=repo_ref://truzhenos/main GOWORK=off npm run smoke:frontend-behavior",
		"go test ./... -run TestShortVideoP16ExecutionReadinessPackageIsMachineCheckable -count=1",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP16PostRunEvidenceAcceptanceGateBlocksP16CompletionAndP17UntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p16 := requireObject(t, candidateSet, "p16_controlled_code_assistant_run")
	p17 := requireObject(t, candidateSet, "p17_provider_adapter_candidate")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p16, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP16 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p16_controlled_code_assistant_run")
	if got := requireString(t, indexP16, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p16 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	p17Readiness := readJSON(t, filepath.Join(base, requireString(t, p17, "execution_readiness_package")))
	if got := requireString(t, p17Readiness, "source_p16_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p17 source_p16_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p17Readiness["required_preconditions"]), "p16_post_run_evidence_acceptance_gate_passed")

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_mark_p16_complete", false)
	requireBool(t, gate, "can_unlock_p17", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p16_execution_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p16_execution_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p16_controlled_code_assistant_run",
		"evidence_contract":                         requireString(t, p16, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p16, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p16, "execution_readiness_package"),
		"previous_slice":                            "p15_gui_walkthrough_three_candidates",
		"next_slice":                                "p17_provider_adapter_candidate",
		"next_slice_execution_readiness_package":    requireString(t, p17, "execution_readiness_package"),
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p16, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_execution_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_execution_evidence" {
			t.Fatalf("%s current_status = %s, want pending_execution_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p16_complete", true)
		requireBool(t, requirement, "blocks_p17", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p16_complete"]), proof)
	}

	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	verifierBlocker := requireObject(t, requireObject(t, verifier, "current_blockers"), "p16_post_run_evidence_acceptance_gate")
	if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("verifier p16 source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), "p16_post_run_evidence_acceptance_gate_not_passed")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	p16Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p16_controlled_code_assistant_run")
	if got := requireString(t, p16Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p16 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p16Gate["blocked_by"]), "p16_post_run_evidence_acceptance_gate_not_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "p16_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p16_post_run_evidence_acceptance_gate_pending")
}

func TestShortVideoP16PostRunGateRejectsStaticGlueAndLocalCliStateSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p16 := requireObject(t, candidateSet, "p16_controlled_code_assistant_run")
	gate := readJSON(t, filepath.Join(base, requireString(t, p16, "post_run_evidence_acceptance_gate")))
	contract := readJSON(t, filepath.Join(base, requireString(t, p16, "evidence_contract")))
	readiness := readJSON(t, filepath.Join(base, requireString(t, p16, "execution_readiness_package")))

	requireBool(t, gate, "static_glue_code_alone_can_mark_p16_complete", false)
	requireBool(t, gate, "owner_local_codex_cli_login_alone_can_mark_p16_complete", false)
	requireBool(t, gate, "manual_patch_candidate_file_alone_can_unlock_p17", false)

	for _, proof := range []string{
		"p16_post_run_evidence_acceptance_gate_passed",
		"controlled_11_gateway_run_receipt_bound",
		"isolated_patch_candidate_artifacts_bound",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p16_complete"]), proof)
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, contract, "completion_gate")["required_before_completion"]), proof)
	}

	for _, nonSufficient := range []string{
		"static_glue_code_without_11_gateway_receipt",
		"owner_local_codex_cli_login_without_controlled_run",
		"manual_patch_candidate_file_without_gateway_receipt",
		"p8_run_request_candidate_without_11_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), nonSufficient)
	}

	policy := requireObject(t, gate, "substitution_rejection_policy")
	for _, key := range []string{
		"static_glue_code_without_11_gateway_receipt",
		"owner_local_codex_cli_login_without_controlled_run",
		"manual_patch_candidate_file_without_gateway_receipt",
		"p8_run_request_candidate_without_11_receipt",
	} {
		if got := requireString(t, policy, key); got != "not_completion_evidence" {
			t.Fatalf("substitution_rejection_policy.%s = %s, want not_completion_evidence", key, got)
		}
	}
}

func TestShortVideoP17ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p17 := requireObject(t, candidateSet, "p17_provider_adapter_candidate")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p17, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization_and_p16_evidence" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization_and_p16_evidence", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p17_provider_adapter_candidate" {
		t.Fatalf("slice_key = %s, want p17_provider_adapter_candidate", got)
	}
	for _, dep := range asStringSlice(t, p17["depends_on"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["depends_on"]), dep)
	}
	for _, providerRef := range asStringSlice(t, p17["provider_candidate_refs"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["provider_candidate_refs"]), providerRef)
	}
	for _, defaultReadiness := range asStringSlice(t, p17["default_readiness"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["default_readiness"]), defaultReadiness)
	}
	requireBool(t, readiness, "provider_implementation_in_packs_allowed", false)
	requireBool(t, readiness, "third_party_oss_execution_allowed_now", false)
	requireBool(t, readiness, "provider_ready_claim_allowed", false)
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p17, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p17, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p17, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p17, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p17, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p17, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP17 := findObjectByString(t, indexSlices, "slice_key", "p17_provider_adapter_candidate")
	if got := requireString(t, indexP17, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p17 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_p16_evidence_complete") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_p16_evidence_complete", got)
	}

	for _, repo := range []string{"truzhen-software", "truzhenos", "truzhen-client-web-desktop", "truzhen-packs"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
		"vendor_third_party_oss_source",
		"provider_implementation_in_packs",
		"provider_ready_claim",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"git status --short --branch && git diff --check",
		"GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1",
		"npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx",
		"go test ./... -run TestShortVideoP17ExecutionReadinessPackageIsMachineCheckable -count=1",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP17PostRunEvidenceAcceptanceGateBlocksP17CompletionAndP18UntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p17 := requireObject(t, candidateSet, "p17_provider_adapter_candidate")
	p18 := requireObject(t, candidateSet, "p18_cloud_market_sandbox")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p17, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP17 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p17_provider_adapter_candidate")
	if got := requireString(t, indexP17, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p17 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	p18Readiness := readJSON(t, filepath.Join(base, requireString(t, p18, "execution_readiness_package")))
	if got := requireString(t, p18Readiness, "source_p17_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("p18 source_p17_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p18Readiness["required_preconditions"]), "p17_post_run_evidence_acceptance_gate_passed")

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_mark_p17_complete", false)
	requireBool(t, gate, "can_unlock_p18", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p17_provider_adapter_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p17_provider_adapter_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p17_provider_adapter_candidate",
		"evidence_contract":                         requireString(t, p17, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p17, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p17, "execution_readiness_package"),
		"previous_slice":                            "p16_controlled_code_assistant_run",
		"previous_slice_post_run_evidence_gate":     requireString(t, requireObject(t, candidateSet, "p16_controlled_code_assistant_run"), "post_run_evidence_acceptance_gate"),
		"next_slice":                                "p18_cloud_market_sandbox",
		"next_slice_execution_readiness_package":    requireString(t, p18, "execution_readiness_package"),
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p17, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_provider_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_provider_adapter_evidence" {
			t.Fatalf("%s current_status = %s, want pending_provider_adapter_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p17_complete", true)
		requireBool(t, requirement, "blocks_p18", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"provider_candidate_outside_packs_verified",
		"readiness_default_blocked_or_missing_verified",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p17_complete"]), proof)
	}

	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	verifierBlocker := requireObject(t, requireObject(t, verifier, "current_blockers"), "p17_post_run_evidence_acceptance_gate")
	if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("verifier p17 source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), "p17_post_run_evidence_acceptance_gate_not_passed")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	p17Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p17_provider_adapter_candidate")
	if got := requireString(t, p17Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p17 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p17Gate["blocked_by"]), "p17_post_run_evidence_acceptance_gate_not_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "p17_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p17_post_run_evidence_acceptance_gate_pending")
	if !strings.Contains(requireString(t, requireObject(t, index, "completion_gate"), "required_status_before_completion"), "p17_post_run_evidence_acceptance_gate_passed") {
		t.Fatalf("completion_gate.required_status_before_completion must require p17_post_run_evidence_acceptance_gate_passed")
	}
}

func TestShortVideoP17PostRunGateRejectsPackLocalProviderSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p17 := requireObject(t, candidateSet, "p17_provider_adapter_candidate")
	gate := readJSON(t, filepath.Join(base, requireString(t, p17, "post_run_evidence_acceptance_gate")))
	contract := readJSON(t, filepath.Join(base, requireString(t, p17, "evidence_contract")))
	readiness := readJSON(t, filepath.Join(base, requireString(t, p17, "execution_readiness_package")))

	requireBool(t, gate, "packs_scaffold_alone_can_mark_p17_complete", false)
	requireBool(t, gate, "provider_manifest_without_readiness_receipt_can_unlock_p18", false)
	requireBool(t, gate, "provider_ready_claim_without_external_provider_evidence_can_count_toward_commercial_ready", false)

	for _, proof := range []string{
		"p17_post_run_evidence_acceptance_gate_passed",
		"external_provider_repository_evidence_bound",
		"provider_readiness_receipt_bound",
		"packs_provider_runtime_absence_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p17_complete"]), proof)
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, contract, "completion_gate")["required_before_completion"]), proof)
	}

	for _, nonSufficient := range []string{
		"adapter_scaffold_inside_packs_without_external_repo_evidence",
		"provider_manifest_without_readiness_receipt",
		"provider_ready_claim_without_owner_base_receipt",
		"vendored_third_party_oss_source_in_packs",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), nonSufficient)
	}

	policy := requireObject(t, gate, "substitution_rejection_policy")
	for _, key := range []string{
		"adapter_scaffold_inside_packs_without_external_repo_evidence",
		"provider_manifest_without_readiness_receipt",
		"provider_ready_claim_without_owner_base_receipt",
		"vendored_third_party_oss_source_in_packs",
	} {
		if got := requireString(t, policy, key); got != "not_completion_evidence" {
			t.Fatalf("substitution_rejection_policy.%s = %s, want not_completion_evidence", key, got)
		}
	}
}

func TestShortVideoP18ExecutionReadinessPackageIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p18 := requireObject(t, candidateSet, "p18_cloud_market_sandbox")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	readinessPath := requireString(t, p18, "execution_readiness_package")
	requireExistingPath(t, readinessPath, base)
	readiness := readJSON(t, filepath.Join(base, readinessPath))

	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "blocked_pending_owner_authorization_and_p17_evidence" {
		t.Fatalf("readiness_status = %s, want blocked_pending_owner_authorization_and_p17_evidence", got)
	}
	if got := requireString(t, readiness, "slice_key"); got != "p18_cloud_market_sandbox" {
		t.Fatalf("slice_key = %s, want p18_cloud_market_sandbox", got)
	}
	requireBool(t, readiness, "sandbox_only", true)
	requireBool(t, readiness, "no_real_payment", true)
	requireBool(t, readiness, "no_production_release", true)
	requireBool(t, readiness, "no_production_license", true)
	requireBool(t, readiness, "packs_must_not_claim_listing_truth", true)
	if got := requireString(t, readiness, "cloud_truth_source"); got != requireString(t, p18, "cloud_truth_source") {
		t.Fatalf("cloud_truth_source = %s, want %s", got, requireString(t, p18, "cloud_truth_source"))
	}
	for _, dep := range asStringSlice(t, p18["depends_on"]) {
		requireStringSliceContains(t, asStringSlice(t, readiness["depends_on"]), dep)
	}
	if got := requireString(t, readiness, "authorization_evidence_intake_contract"); got != requireString(t, p18, "authorization_evidence_intake_contract") {
		t.Fatalf("authorization_evidence_intake_contract = %s, want %s", got, requireString(t, p18, "authorization_evidence_intake_contract"))
	}
	if got := requireString(t, readiness, "authorization_scope_contract"); got != requireString(t, p18, "authorization_scope_contract") {
		t.Fatalf("authorization_scope_contract = %s, want %s", got, requireString(t, p18, "authorization_scope_contract"))
	}
	if got := requireString(t, readiness, "evidence_contract"); got != requireString(t, p18, "evidence_contract") {
		t.Fatalf("evidence_contract = %s, want %s", got, requireString(t, p18, "evidence_contract"))
	}
	if got := requireString(t, readiness, "frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	indexP18 := findObjectByString(t, indexSlices, "slice_key", "p18_cloud_market_sandbox")
	if got := requireString(t, indexP18, "execution_readiness_package"); got != readinessPath {
		t.Fatalf("index p18 execution_readiness_package = %s, want %s", got, readinessPath)
	}

	gate := requireObject(t, readiness, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	if got := requireString(t, gate, "required_status_before_cross_repo_work"); !strings.Contains(got, "owner_authorization_recorded_and_p17_evidence_complete") {
		t.Fatalf("required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_p17_evidence_complete", got)
	}

	for _, repo := range []string{"truzhen-cloud", "truzhen-client-web-desktop", "truzhenos", "truzhen-packs"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["target_repositories"]), repo)
	}
	for _, repo := range []string{"truzhen-contracts", "truzhen-software"} {
		requireStringSliceContains(t, asStringSlice(t, readiness["disallowed_repositories"]), repo)
	}
	for _, action := range []string{
		"execute_third_party_oss",
		"social_login_or_upload",
		"store_raw_secret",
		"real_payment",
		"production_release",
		"production_license",
		"cloud_listing_truth_stored_in_packs",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["forbidden_actions"]), action)
	}
	for _, command := range []string{
		"git status --short --branch && git diff --check",
		"npm test -- src/pages/__tests__/capabilityStudioWizard.test.tsx",
		"GOWORK=off go test ./backend/internal/capability/... ./backend/tests/capability/... -count=1",
		"go test ./... -run TestShortVideoP18ExecutionReadinessPackageIsMachineCheckable -count=1",
	} {
		requireStringSliceContains(t, asStringSlice(t, readiness["required_verification_commands"]), command)
	}
}

func TestShortVideoP18PostRunEvidenceAcceptanceGateBlocksCommercialSignoffUntilAllEvidenceAccepted(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p18 := requireObject(t, candidateSet, "p18_cloud_market_sandbox")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	gatePath := requireString(t, p18, "post_run_evidence_acceptance_gate")
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexP18 := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", "p18_cloud_market_sandbox")
	if got := requireString(t, indexP18, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("index p18 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "final_slice", true)
	requireBool(t, gate, "can_mark_p18_complete", false)
	requireBool(t, gate, "can_request_commercial_signoff", false)
	requireBool(t, gate, "can_mark_commercial_ready", false)
	if got := requireString(t, gate, "gate_status"); got != "blocked_pending_p18_cloud_sandbox_evidence" {
		t.Fatalf("gate_status = %s, want blocked_pending_p18_cloud_sandbox_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"slice_key":                                 "p18_cloud_market_sandbox",
		"evidence_contract":                         requireString(t, p18, "evidence_contract"),
		"evidence_ledger":                           requireString(t, p18, "evidence_ledger"),
		"execution_readiness_package":               requireString(t, p18, "execution_readiness_package"),
		"previous_slice":                            "p17_provider_adapter_candidate",
		"previous_slice_post_run_evidence_gate":     requireString(t, requireObject(t, candidateSet, "p17_provider_adapter_candidate"), "post_run_evidence_acceptance_gate"),
		"commercial_signoff_target":                 "commercial_ready_candidate",
		"commercial_readiness_verifier":             requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                  requireString(t, goLive, "commercial_go_no_go_gate"),
		"commercial_evidence_writeback_gate_source": requireString(t, goLive, "required_evidence_writeback_gate_source"),
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	contract := readJSON(t, filepath.Join(base, requireString(t, p18, "evidence_contract")))
	expectedEvidence := map[string]string{}
	for _, group := range []string{"required_cloud_evidence", "required_forbidden_action_checks"} {
		for _, item := range asObjectSlice(t, contract[group]) {
			expectedEvidence[requireString(t, item, "evidence_id")] = group
		}
	}
	requirements := asObjectSlice(t, gate["acceptance_requirements"])
	if len(requirements) != len(expectedEvidence) {
		t.Fatalf("acceptance_requirements len = %d, want %d", len(requirements), len(expectedEvidence))
	}
	for _, requirement := range requirements {
		evidenceID := requireString(t, requirement, "evidence_id")
		group, ok := expectedEvidence[evidenceID]
		if !ok {
			t.Fatalf("unexpected acceptance evidence_id %s", evidenceID)
		}
		if got := requireString(t, requirement, "source_group"); got != group {
			t.Fatalf("%s source_group = %s, want %s", evidenceID, got, group)
		}
		if got := requireString(t, requirement, "current_status"); got != "pending_cloud_sandbox_evidence" {
			t.Fatalf("%s current_status = %s, want pending_cloud_sandbox_evidence", evidenceID, got)
		}
		requireBool(t, requirement, "required_before_p18_complete", true)
		requireBool(t, requirement, "blocks_commercial_signoff", true)
		requireBool(t, requirement, "evidence_ref_required", true)
		requireBool(t, requirement, "evidence_summary_required", true)
		requireBool(t, requirement, "independent_review_required", true)
		if got := requireString(t, requirement, "authoritative_source"); got == "" {
			t.Fatalf("%s authoritative_source missing", evidenceID)
		}
		delete(expectedEvidence, evidenceID)
	}
	if len(expectedEvidence) != 0 {
		t.Fatalf("acceptance_requirements missing evidence ids: %v", expectedEvidence)
	}

	for _, proof := range []string{
		"all_evidence_ids_written",
		"all_evidence_refs_recorded",
		"cloud_truth_source_verified",
		"sandbox_entitlement_download_install_receipts_verified",
		"forbidden_action_checks_written",
		"independent_verification_passed",
		"owner_base_gate_receipts_bound",
		"commercial_evidence_writeback_gate_updated",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p18_complete"]), proof)
	}

	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	verifierBlocker := requireObject(t, requireObject(t, verifier, "current_blockers"), "p18_post_run_evidence_acceptance_gate")
	if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("verifier p18 source_post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), "p18_post_run_evidence_acceptance_gate_not_passed")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	p18Gate := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", "p18_cloud_market_sandbox")
	if got := requireString(t, p18Gate, "post_run_evidence_acceptance_gate"); got != gatePath {
		t.Fatalf("go-no-go p18 post_run_evidence_acceptance_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, p18Gate["blocked_by"]), "p18_post_run_evidence_acceptance_gate_not_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "p18_post_run_evidence_acceptance_gate_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p18_post_run_evidence_acceptance_gate_pending")
	if !strings.Contains(requireString(t, requireObject(t, index, "completion_gate"), "required_status_before_completion"), "p18_post_run_evidence_acceptance_gate_passed") {
		t.Fatalf("completion_gate.required_status_before_completion must require p18_post_run_evidence_acceptance_gate_passed")
	}
}

func TestShortVideoP18PostRunGateRejectsPackListingAndSandboxDocSubstitution(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p18 := requireObject(t, candidateSet, "p18_cloud_market_sandbox")
	gate := readJSON(t, filepath.Join(base, requireString(t, p18, "post_run_evidence_acceptance_gate")))
	contract := readJSON(t, filepath.Join(base, requireString(t, p18, "evidence_contract")))
	readiness := readJSON(t, filepath.Join(base, requireString(t, p18, "execution_readiness_package")))

	requireBool(t, gate, "packs_listing_draft_alone_can_mark_p18_complete", false)
	requireBool(t, gate, "sandbox_runbook_without_cloud_receipt_can_mark_commercial_ready", false)
	requireBool(t, gate, "license_or_entitlement_text_without_cloud_truth_can_unlock_go_live", false)

	for _, proof := range []string{
		"p18_post_run_evidence_acceptance_gate_passed",
		"cloud_sandbox_receipts_bound",
		"sandbox_entitlement_receipt_bound",
		"packs_listing_truth_absence_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_p18_complete"]), proof)
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, contract, "completion_gate")["required_before_completion"]), proof)
	}

	for _, nonSufficient := range []string{
		"packs_listing_draft_without_cloud_listing_receipt",
		"sandbox_runbook_without_cloud_receipts",
		"license_or_entitlement_text_without_cloud_truth",
		"manual_order_status_without_cloud_receipt",
		"production_release_claim_without_cloud_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), nonSufficient)
	}

	policy := requireObject(t, gate, "substitution_rejection_policy")
	for _, key := range []string{
		"packs_listing_draft_without_cloud_listing_receipt",
		"sandbox_runbook_without_cloud_receipts",
		"license_or_entitlement_text_without_cloud_truth",
		"manual_order_status_without_cloud_receipt",
		"production_release_claim_without_cloud_receipt",
	} {
		if got := requireString(t, policy, key); got != "not_completion_evidence" {
			t.Fatalf("substitution_rejection_policy.%s = %s, want not_completion_evidence", key, got)
		}
	}
}

func TestShortVideoCommercialPostRunGateCoverageVerifierCoversAllRequiredSlices(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	coveragePath := requireString(t, goLive, "post_run_gate_coverage_verifier")
	requireExistingPath(t, coveragePath, base)
	coverage := readJSON(t, filepath.Join(base, coveragePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	readinessVerifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	if got := requireString(t, coverage, "verifier_ref"); got != "commercial-post-run-gate-coverage-verifier://short-video-ops-v0" {
		t.Fatalf("verifier_ref = %s, want commercial-post-run-gate-coverage-verifier://short-video-ops-v0", got)
	}
	requireBool(t, coverage, "candidate_only", true)
	requireBool(t, coverage, "non_formal", true)
	requireBool(t, coverage, "can_mark_commercial_ready", false)
	if got := requireString(t, coverage, "coverage_status"); got != "blocked_pending_all_post_run_gates_passed" {
		t.Fatalf("coverage_status = %s, want blocked_pending_all_post_run_gates_passed", got)
	}
	requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), "post_run_gate_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "owner_signoff_gate")["blocking_required_before_owner_signoff"]), "post_run_gate_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "commercial_signoff_matrix")["non_sufficient_evidence"]), "post_run_gate_coverage_verifier_pending")

	coverageBlocker := requireObject(t, requireObject(t, readinessVerifier, "current_blockers"), "post_run_gate_coverage_verifier")
	if got := requireString(t, coverageBlocker, "source_post_run_gate_coverage_verifier"); got != coveragePath {
		t.Fatalf("post_run_gate_coverage_verifier source = %s, want %s", got, coveragePath)
	}
	if got := requireString(t, coverageBlocker, "coverage_status"); got != requireString(t, coverage, "coverage_status") {
		t.Fatalf("post_run_gate_coverage_verifier coverage_status = %s, want %s", got, requireString(t, coverage, "coverage_status"))
	}
	requireBool(t, coverageBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, coverageBlocker["blocked_by"]), "post_run_gate_coverage_verifier_pending")

	if got := requireString(t, goNoGo, "source_post_run_gate_coverage_verifier"); got != coveragePath {
		t.Fatalf("source_post_run_gate_coverage_verifier = %s, want %s", got, coveragePath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "post_run_gate_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "post_run_gate_coverage_verifier_pending")
	for key, want := range map[string]string{
		"candidate_set_ref":                    requireString(t, candidateSet, "candidate_set_ref"),
		"source_evidence_contract_index":       requireString(t, goLive, "evidence_contract_index"),
		"source_commercial_readiness_verifier": requireString(t, goLive, "commercial_readiness_verifier"),
		"source_commercial_go_no_go_gate":      requireString(t, goLive, "commercial_go_no_go_gate"),
		"required_order_source":                "commercial_go_live_evidence_package.required_slices",
	} {
		if got := requireString(t, coverage, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	requiredSlices := asStringSlice(t, goLive["required_slices"])
	coverageSlices := asStringSlice(t, coverage["required_slice_keys"])
	if len(coverageSlices) != len(requiredSlices) {
		t.Fatalf("required_slice_keys len = %d, want %d", len(coverageSlices), len(requiredSlices))
	}
	for _, sliceKey := range requiredSlices {
		requireStringSliceContains(t, coverageSlices, sliceKey)
	}

	expectedGateMeta := map[string]struct {
		blockerKey     string
		requiredStatus string
		pendingMarker  string
	}{
		"p12_safe_lifecycle_sample": {
			blockerKey:     "p12_post_run_evidence_acceptance_gate",
			requiredStatus: "p12_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p12_post_run_evidence_acceptance_gate_pending",
		},
		"p13_gui_lifecycle_panel": {
			blockerKey:     "p13_post_run_evidence_acceptance_gate",
			requiredStatus: "p13_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p13_post_run_evidence_acceptance_gate_pending",
		},
		"p15_gui_walkthrough_three_candidates": {
			blockerKey:     "p15_post_run_evidence_acceptance_gate",
			requiredStatus: "p15_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p15_post_run_evidence_acceptance_gate_pending",
		},
		"p16_controlled_code_assistant_run": {
			blockerKey:     "p16_post_run_evidence_acceptance_gate",
			requiredStatus: "p16_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p16_post_run_evidence_acceptance_gate_pending",
		},
		"p17_provider_adapter_candidate": {
			blockerKey:     "p17_post_run_evidence_acceptance_gate",
			requiredStatus: "p17_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p17_post_run_evidence_acceptance_gate_pending",
		},
		"p18_cloud_market_sandbox": {
			blockerKey:     "p18_post_run_evidence_acceptance_gate",
			requiredStatus: "p18_post_run_evidence_acceptance_gate_passed",
			pendingMarker:  "p18_post_run_evidence_acceptance_gate_pending",
		},
	}
	coverageChecks := asObjectSlice(t, coverage["coverage_checks"])
	if len(coverageChecks) != len(requiredSlices) {
		t.Fatalf("coverage_checks len = %d, want %d", len(coverageChecks), len(requiredSlices))
	}

	for _, sliceKey := range requiredSlices {
		meta, ok := expectedGateMeta[sliceKey]
		if !ok {
			t.Fatalf("missing expected gate meta for %s", sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		gatePath := requireString(t, sourceSlice, "post_run_evidence_acceptance_gate")
		requireExistingPath(t, gatePath, base)
		gate := readJSON(t, filepath.Join(base, gatePath))
		requireBool(t, gate, "candidate_only", true)
		requireBool(t, gate, "non_formal", true)
		requireBool(t, gate, "can_mark_commercial_ready", false)

		check := findObjectByString(t, coverageChecks, "slice_key", sliceKey)
		if got := requireString(t, check, "post_run_evidence_acceptance_gate"); got != gatePath {
			t.Fatalf("%s coverage gate = %s, want %s", sliceKey, got, gatePath)
		}
		if got := requireString(t, check, "verifier_blocker_key"); got != meta.blockerKey {
			t.Fatalf("%s verifier_blocker_key = %s, want %s", sliceKey, got, meta.blockerKey)
		}
		if got := requireString(t, check, "go_no_go_required_status"); got != meta.requiredStatus {
			t.Fatalf("%s go_no_go_required_status = %s, want %s", sliceKey, got, meta.requiredStatus)
		}
		if got := requireString(t, check, "go_no_go_pending_marker"); got != meta.pendingMarker {
			t.Fatalf("%s go_no_go_pending_marker = %s, want %s", sliceKey, got, meta.pendingMarker)
		}
		if got := requireString(t, check, "coverage_result"); got != "blocked_pending_execution_evidence" {
			t.Fatalf("%s coverage_result = %s, want blocked_pending_execution_evidence", sliceKey, got)
		}
		for _, flag := range []string{
			"candidate_set_has_gate",
			"evidence_index_has_gate",
			"readiness_verifier_has_blocker",
			"go_no_go_gate_has_gate",
			"go_no_go_completion_status_required",
			"go_no_go_pending_marker_present",
		} {
			requireBool(t, check, flag, true)
		}

		indexSlice := findObjectByString(t, asObjectSlice(t, index["required_slices"]), "slice_key", sliceKey)
		if got := requireString(t, indexSlice, "post_run_evidence_acceptance_gate"); got != gatePath {
			t.Fatalf("%s index gate = %s, want %s", sliceKey, got, gatePath)
		}

		verifierBlocker := requireObject(t, requireObject(t, readinessVerifier, "current_blockers"), meta.blockerKey)
		if got := requireString(t, verifierBlocker, "source_post_run_evidence_acceptance_gate"); got != gatePath {
			t.Fatalf("%s verifier blocker source = %s, want %s", sliceKey, got, gatePath)
		}
		requireBool(t, verifierBlocker, "can_count_toward_commercial_ready", false)
		requireStringSliceContains(t, asStringSlice(t, verifierBlocker["blocked_by"]), meta.blockerKey+"_not_passed")

		goNoGoSlice := findObjectByString(t, asObjectSlice(t, goNoGo["required_slice_gates"]), "slice_key", sliceKey)
		if got := requireString(t, goNoGoSlice, "post_run_evidence_acceptance_gate"); got != gatePath {
			t.Fatalf("%s go-no-go gate = %s, want %s", sliceKey, got, gatePath)
		}
		requireStringSliceContains(t, asStringSlice(t, goNoGoSlice["blocked_by"]), meta.blockerKey+"_not_passed")
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), meta.requiredStatus)
		requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), meta.pendingMarker)
	}

	completionGate := requireObject(t, coverage, "completion_gate")
	requireBool(t, completionGate, "can_pass_coverage_verifier", false)
	for _, proof := range []string{
		"all_post_run_gate_sources_present",
		"all_post_run_gate_blockers_present",
		"all_go_no_go_completion_statuses_required",
		"all_pending_markers_present",
		"all_post_run_gates_passed",
		"owner_base_gate_receipts_bound",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_commercial_ready"]), proof)
	}
	requireStringSliceContains(t, asStringSlice(t, coverage["non_sufficient_evidence"]), "post_run_gate_coverage_verifier_without_passed_gates")
}

func TestShortVideoCommercialExecutionReadinessPackagesDefineEvidenceWritebackPlans(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		readiness := readJSON(t, filepath.Join(base, requireString(t, sourceSlice, "execution_readiness_package")))
		contract := readJSON(t, filepath.Join(base, requireString(t, sourceSlice, "evidence_contract")))
		ledgerPath := requireString(t, sourceSlice, "evidence_ledger")

		plan := requireObject(t, readiness, "evidence_writeback_plan")
		requireBool(t, plan, "candidate_only", true)
		requireBool(t, plan, "can_write_completion_claim", false)
		if got := requireString(t, plan, "write_status"); got != "pending_authorization" {
			t.Fatalf("%s evidence_writeback_plan.write_status = %s, want pending_authorization", sliceKey, got)
		}
		if got := requireString(t, plan, "evidence_contract"); got != requireString(t, sourceSlice, "evidence_contract") {
			t.Fatalf("%s evidence_writeback_plan.evidence_contract = %s, want %s", sliceKey, got, requireString(t, sourceSlice, "evidence_contract"))
		}
		if got := requireString(t, plan, "ledger_path"); got != ledgerPath {
			t.Fatalf("%s evidence_writeback_plan.ledger_path = %s, want %s", sliceKey, got, ledgerPath)
		}
		requireStringSliceContains(t, asStringSlice(t, plan["required_before_marking_slice_complete"]), "all_evidence_ids_written")
		requireStringSliceContains(t, asStringSlice(t, plan["required_before_marking_slice_complete"]), "all_evidence_refs_recorded")
		requireStringSliceContains(t, asStringSlice(t, plan["required_before_marking_slice_complete"]), "forbidden_action_checks_written")

		expected := map[string]string{}
		expectedStatus := map[string]string{}
		for group, raw := range contract {
			items, ok := raw.([]any)
			if !ok {
				continue
			}
			for _, itemRaw := range items {
				item, ok := itemRaw.(map[string]any)
				if !ok {
					continue
				}
				evidenceID, ok := item["evidence_id"].(string)
				if !ok || evidenceID == "" {
					continue
				}
				expected[evidenceID] = group
				expectedStatus[evidenceID] = requireString(t, item, "current_status")
			}
		}
		if len(expected) == 0 {
			t.Fatalf("%s evidence contract has no evidence_id rows", sliceKey)
		}

		entries := asObjectSlice(t, plan["required_entries"])
		if len(entries) != len(expected) {
			t.Fatalf("%s evidence_writeback_plan.required_entries len = %d, want %d", sliceKey, len(entries), len(expected))
		}
		seen := map[string]bool{}
		for _, entry := range entries {
			evidenceID := requireString(t, entry, "evidence_id")
			sourceGroup, ok := expected[evidenceID]
			if !ok {
				t.Fatalf("%s unexpected evidence_writeback_plan evidence_id %s", sliceKey, evidenceID)
			}
			if seen[evidenceID] {
				t.Fatalf("%s duplicate evidence_writeback_plan evidence_id %s", sliceKey, evidenceID)
			}
			seen[evidenceID] = true
			if got := requireString(t, entry, "source_group"); got != sourceGroup {
				t.Fatalf("%s %s source_group = %s, want %s", sliceKey, evidenceID, got, sourceGroup)
			}
			if got := requireString(t, entry, "target_ledger"); got != ledgerPath {
				t.Fatalf("%s %s target_ledger = %s, want %s", sliceKey, evidenceID, got, ledgerPath)
			}
			if got := requireString(t, entry, "current_status"); got != expectedStatus[evidenceID] {
				t.Fatalf("%s %s current_status = %s, want %s", sliceKey, evidenceID, got, expectedStatus[evidenceID])
			}
			requireBool(t, entry, "write_required", true)
			requireBool(t, entry, "evidence_ref_required", true)
		}
	}
}

func TestShortVideoCommercialSlicesHaveAuthorizationEvidenceIntakeContracts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexSlices := asObjectSlice(t, index["required_slices"])

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			scopePath := requireString(t, slice, "authorization_scope_contract")
			scope := readJSON(t, filepath.Join(base, scopePath))
			scopeEvidence := requireObject(t, scope, "current_authorization_evidence")

			intakePath := requireString(t, slice, "authorization_evidence_intake_contract")
			indexSlice := findObjectByString(t, indexSlices, "slice_key", sliceKey)
			if got := requireString(t, indexSlice, "authorization_evidence_intake_contract"); got != intakePath {
				t.Fatalf("%s index authorization_evidence_intake_contract = %s, want %s", sliceKey, got, intakePath)
			}
			if got := requireString(t, scopeEvidence, "evidence_intake_contract"); got != intakePath {
				t.Fatalf("%s scope evidence_intake_contract = %s, want %s", sliceKey, got, intakePath)
			}
			requireExistingPath(t, intakePath, base)
			intake := readJSON(t, filepath.Join(base, intakePath))

			requireBool(t, intake, "candidate_only", true)
			requireBool(t, intake, "non_formal", true)
			if got := requireString(t, intake, "slice_key"); got != sliceKey {
				t.Fatalf("slice_key = %s, want %s", got, sliceKey)
			}
			if got := requireString(t, intake, "status"); got != "missing_owner_authorization" {
				t.Fatalf("%s status = %s, want missing_owner_authorization", sliceKey, got)
			}
			if got := requireString(t, intake, "authorization_scope_contract"); got != scopePath {
				t.Fatalf("%s authorization_scope_contract = %s, want %s", sliceKey, got, scopePath)
			}
			if got := requireString(t, intake, "accepted_authorization_card"); got != requireString(t, scopeEvidence, "accepted_authorization_card") {
				t.Fatalf("%s accepted_authorization_card = %s, want %s", sliceKey, got, requireString(t, scopeEvidence, "accepted_authorization_card"))
			}

			for _, phrasePart := range asStringSlice(t, scope["required_authorization_phrase_contains"]) {
				requireStringSliceContains(t, asStringSlice(t, intake["required_owner_quote_contains"]), phrasePart)
			}
			for _, repo := range asStringSlice(t, slice["target_repositories"]) {
				requireStringSliceContains(t, asStringSlice(t, intake["authorized_repositories"]), repo)
			}
			for _, repo := range asStringSlice(t, scope["disallowed_repositories"]) {
				requireStringSliceContains(t, asStringSlice(t, intake["disallowed_repositories"]), repo)
			}
			for _, action := range asStringSlice(t, scope["forbidden_actions"]) {
				requireStringSliceContains(t, asStringSlice(t, intake["forbidden_actions"]), action)
			}

			evidence := requireObject(t, intake, "current_authorization_evidence")
			if got := requireString(t, evidence, "status"); got != "missing" {
				t.Fatalf("%s current_authorization_evidence.status = %s, want missing", sliceKey, got)
			}
			requireBool(t, evidence, "owner_thread_quote_required", true)

			gate := requireObject(t, intake, "cross_repo_work_gate")
			requireBool(t, gate, "can_start_cross_repo_work", false)
			if got := requireString(t, gate, "required_status_before_cross_repo_work"); got != "owner_authorization_recorded_and_scope_matched" {
				t.Fatalf("%s required_status_before_cross_repo_work = %s, want owner_authorization_recorded_and_scope_matched", sliceKey, got)
			}
			requireBool(t, gate, "must_record_before_cross_repo_work", true)
		})
	}
}

func TestShortVideoCommercialSlicesHaveMachineCheckableAuthorizationScopes(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	indexSlices := asObjectSlice(t, index["required_slices"])

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			scopePath := requireString(t, slice, "authorization_scope_contract")
			requireExistingPath(t, scopePath, base)

			indexSlice := findObjectByString(t, indexSlices, "slice_key", sliceKey)
			if got := requireString(t, indexSlice, "authorization_scope_contract"); got != scopePath {
				t.Fatalf("%s index authorization_scope_contract = %s, want %s", sliceKey, got, scopePath)
			}

			scope := readJSON(t, filepath.Join(base, scopePath))
			requireBool(t, scope, "candidate_only", true)
			requireBool(t, scope, "non_formal", true)
			if got := requireString(t, scope, "slice_key"); got != sliceKey {
				t.Fatalf("slice_key = %s, want %s", got, sliceKey)
			}
			if got := requireString(t, scope, "status"); got != "pending_owner_authorization" {
				t.Fatalf("%s status = %s, want pending_owner_authorization", sliceKey, got)
			}

			authorizedRepos := asObjectSlice(t, scope["authorized_repositories"])
			for _, repo := range asStringSlice(t, slice["target_repositories"]) {
				entry := findObjectByString(t, authorizedRepos, "repo", repo)
				if len(asStringSlice(t, entry["allowed_actions"])) == 0 {
					t.Fatalf("%s authorized repo %s has no allowed_actions", sliceKey, repo)
				}
			}
			for _, repo := range asStringSlice(t, scope["disallowed_repositories"]) {
				for _, targetRepo := range asStringSlice(t, slice["target_repositories"]) {
					if repo == targetRepo {
						t.Fatalf("%s disallows target repository %s", sliceKey, repo)
					}
				}
			}
			for _, action := range []string{
				"execute_third_party_oss",
				"social_login_or_upload",
				"store_raw_secret",
			} {
				requireStringSliceContains(t, asStringSlice(t, scope["forbidden_actions"]), action)
			}
			if len(asStringSlice(t, scope["required_authorization_phrase_contains"])) == 0 {
				t.Fatalf("%s required_authorization_phrase_contains missing", sliceKey)
			}

			evidence := requireObject(t, scope, "current_authorization_evidence")
			if got := requireString(t, evidence, "status"); got != "missing" {
				t.Fatalf("%s current_authorization_evidence.status = %s, want missing before Owner authorization", sliceKey, got)
			}
			requireBool(t, evidence, "owner_thread_quote_required", true)
			requireBool(t, evidence, "must_be_recorded_before_cross_repo_work", true)
		})
	}
}

func TestShortVideoAuthorizationCardsMirrorMachineScopeReferencesAndPhraseGates(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			scopePath := requireString(t, slice, "authorization_scope_contract")
			scope := readJSON(t, filepath.Join(base, scopePath))

			cardPath := requireString(t, slice, "cross_repo_authorization_card")
			if strings.HasPrefix(cardPath, "process_worktree_ref://") {
				// 同 requireExistingPath 的方案 b：授权卡产自外部过程 worktree，
				// 内容级镜像断言在其缺席机器上无法执行，降级为引用登记存在性。
				requireExistingPath(t, cardPath, "")
				return
			}
			cardBytes, err := os.ReadFile(resolveReferencePath(cardPath))
			if err != nil {
				t.Fatalf("read authorization card %s: %v", cardPath, err)
			}
			card := string(cardBytes)
			if !strings.Contains(card, scopePath) {
				t.Fatalf("%s authorization card must reference scope path %s", sliceKey, scopePath)
			}
			scopeRef := requireString(t, scope, "scope_ref")
			if !strings.Contains(card, scopeRef) {
				t.Fatalf("%s authorization card must reference scope_ref %s", sliceKey, scopeRef)
			}
			for _, phrasePart := range asStringSlice(t, scope["required_authorization_phrase_contains"]) {
				if !strings.Contains(card, phrasePart) {
					t.Fatalf("%s authorization card missing required phrase part %q", sliceKey, phrasePart)
				}
			}
		})
	}
}

func TestShortVideoCommercialSlicesHaveMachineReadableEvidenceContracts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			contractPath := requireString(t, slice, "evidence_contract")
			contractFile := contractPath
			if !filepath.IsAbs(contractFile) {
				contractFile = filepath.Join(base, contractPath)
			}
			requireExistingPath(t, contractFile, "")
			contract := readJSON(t, contractFile)

			requireBool(t, contract, "candidate_only", true)
			requireBool(t, contract, "non_formal", true)
			if got := requireString(t, contract, "slice_key"); got != sliceKey {
				t.Fatalf("slice_key = %s, want %s", got, sliceKey)
			}
			if status := requireString(t, contract, "status"); !strings.Contains(status, "pending") {
				t.Fatalf("%s status = %s, want pending before execution", sliceKey, status)
			}
			for _, repo := range asStringSlice(t, slice["target_repositories"]) {
				requireStringSliceContains(t, asStringSlice(t, contract["required_repositories"]), repo)
			}
			requireContractHasEvidenceGroup(t, contract)

			gate := requireObject(t, contract, "completion_gate")
			if requireBoolValue(t, gate, "can_mark_slice_complete") {
				t.Fatalf("%s completion gate must remain false before evidence is collected", sliceKey)
			}
			if got := requireString(t, gate, "required_status_before_completion"); got == "" {
				t.Fatalf("%s completion_gate.required_status_before_completion missing", sliceKey)
			}
			if got := requireString(t, gate, "commercial_status_after_slice_only"); got != "not_commercial_ready" {
				t.Fatalf("%s commercial_status_after_slice_only = %s, want not_commercial_ready", sliceKey, got)
			}
		})
	}
}

func TestShortVideoCommercialEvidenceContractsRequireSliceWritebackPlans(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			contractPath := requireString(t, slice, "evidence_contract")
			readinessPath := requireString(t, slice, "execution_readiness_package")
			contract := readJSON(t, filepath.Join(base, contractPath))
			readiness := readJSON(t, filepath.Join(base, readinessPath))
			writebackPlan := requireObject(t, readiness, "evidence_writeback_plan")

			completionGate := requireObject(t, contract, "completion_gate")
			if got := requireString(t, completionGate, "required_evidence_writeback_plan_source"); got != readinessPath+"#evidence_writeback_plan" {
				t.Fatalf("%s completion_gate.required_evidence_writeback_plan_source = %s, want %s#evidence_writeback_plan", sliceKey, got, readinessPath)
			}
			requireBool(t, completionGate, "evidence_writeback_plan_must_complete", true)
			requiredStatus := requireString(t, completionGate, "required_status_before_completion")
			for _, required := range asStringSlice(t, writebackPlan["required_before_marking_slice_complete"]) {
				if !strings.Contains(requiredStatus, required) {
					t.Fatalf("%s completion_gate.required_status_before_completion must include writeback requirement %q", sliceKey, required)
				}
			}
			requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), "evidence_writeback_incomplete")
		})
	}
}

func TestShortVideoCommercialEvidenceContractsRequireOwnerBaseGateReceipts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			contract := readJSON(t, filepath.Join(base, requireString(t, slice, "evidence_contract")))
			gate := requireObject(t, contract, "completion_gate")

			if got := requireString(t, gate, "required_owner_base_gate_source"); got != "Owner + Base Gate + Gateway + Receipt" {
				t.Fatalf("%s completion_gate.required_owner_base_gate_source = %s, want Owner + Base Gate + Gateway + Receipt", sliceKey, got)
			}
			requireBool(t, gate, "owner_base_gate_receipts_must_exist", true)
			if got := requireString(t, gate, "blocking_owner_base_gate"); got != "owner_base_gate_receipts_missing" {
				t.Fatalf("%s completion_gate.blocking_owner_base_gate = %s, want owner_base_gate_receipts_missing", sliceKey, got)
			}
			if !strings.Contains(requireString(t, gate, "required_status_before_completion"), "owner_base_gate_receipts_bound") {
				t.Fatalf("%s completion_gate.required_status_before_completion must require owner_base_gate_receipts_bound", sliceKey)
			}
			requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), "owner_base_gate_receipts_missing")
		})
	}
}

func TestShortVideoCommercialEvidenceLedgersCoverRequiredContractEvidenceIDs(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	for _, sliceKey := range asStringSlice(t, goLive["required_slices"]) {
		t.Run(sliceKey, func(t *testing.T) {
			slice := requireObject(t, candidateSet, sliceKey)
			contractPath := requireString(t, slice, "evidence_contract")
			ledgerPath := requireString(t, slice, "evidence_ledger")
			contract := readJSON(t, filepath.Join(base, contractPath))
			ledgerBytes, err := os.ReadFile(filepath.Join(base, ledgerPath))
			if err != nil {
				t.Fatalf("read ledger %s: %v", ledgerPath, err)
			}
			ledger := string(ledgerBytes)
			if !strings.Contains(ledger, contractPath) {
				t.Fatalf("%s ledger must reference evidence contract %s", sliceKey, contractPath)
			}
			for _, evidenceID := range collectEvidenceIDsFromContract(t, contract) {
				if !strings.Contains(ledger, evidenceID) {
					t.Fatalf("%s ledger missing evidence_id %s from %s", sliceKey, evidenceID, contractPath)
				}
			}
		})
	}
}

func TestShortVideoCommercialEvidenceContractIndexMatchesCandidateSetAndSignoffGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	indexPath := requireString(t, goLive, "evidence_contract_index")
	indexFile := indexPath
	if !filepath.IsAbs(indexFile) {
		indexFile = filepath.Join(base, indexPath)
	}
	requireExistingPath(t, indexFile, "")
	index := readJSON(t, indexFile)

	requireBool(t, index, "candidate_only", true)
	requireBool(t, index, "non_formal", true)
	if got := requireString(t, index, "current_go_live_status"); got != requireString(t, goLive, "current_go_live_status") {
		t.Fatalf("current_go_live_status = %s, want %s", got, requireString(t, goLive, "current_go_live_status"))
	}

	indexSlices := asObjectSlice(t, index["required_slices"])
	goLiveSlices := asStringSlice(t, goLive["required_slices"])
	if len(indexSlices) != len(goLiveSlices) {
		t.Fatalf("index required_slices len = %d, want %d", len(indexSlices), len(goLiveSlices))
	}
	for _, sliceKey := range goLiveSlices {
		indexItem := findObjectByString(t, indexSlices, "slice_key", sliceKey)
		slice := requireObject(t, candidateSet, sliceKey)
		for _, key := range []string{"evidence_contract", "evidence_ledger", "implementation_spec", "cross_repo_authorization_card"} {
			if got, want := requireString(t, indexItem, key), requireString(t, slice, key); got != want {
				t.Fatalf("%s.%s = %s, want %s", sliceKey, key, got, want)
			}
		}
		requireExistingPath(t, requireString(t, indexItem, "evidence_contract"), base)
		requireExistingPath(t, requireString(t, indexItem, "evidence_ledger"), base)
		requireExistingPath(t, requireString(t, indexItem, "implementation_spec"), "")
		requireExistingPath(t, requireString(t, indexItem, "cross_repo_authorization_card"), "")
		for _, repo := range asStringSlice(t, slice["target_repositories"]) {
			requireStringSliceContains(t, asStringSlice(t, indexItem["target_repositories"]), repo)
		}
	}

	indexGate := requireObject(t, index, "owner_signoff_gate")
	goLiveGate := requireObject(t, goLive, "owner_signoff_gate")
	requireBool(t, indexGate, "can_request_owner_signoff", false)
	for _, key := range []string{"required_next_authorization", "next_authorization_card"} {
		if got, want := requireString(t, indexGate, key), requireString(t, goLiveGate, key); got != want {
			t.Fatalf("owner_signoff_gate.%s = %s, want %s", key, got, want)
		}
	}
	for _, sliceKey := range goLiveSlices {
		requireStringSliceContains(t, asStringSlice(t, indexGate["blocking_required_slices"]), sliceKey)
	}
	for _, checkKey := range asStringSlice(t, goLiveGate["blocking_terminal_checks"]) {
		requireStringSliceContains(t, asStringSlice(t, indexGate["blocking_terminal_checks"]), checkKey)
	}

	completionGate := requireObject(t, index, "completion_gate")
	requireBool(t, completionGate, "can_mark_commercial_ready", false)
	if got := requireString(t, completionGate, "required_status_before_completion"); got == "" {
		t.Fatalf("completion_gate.required_status_before_completion missing")
	}
	if got := requireString(t, completionGate, "final_status_after_success"); got != "commercial_ready_candidate" {
		t.Fatalf("completion_gate.final_status_after_success = %s, want commercial_ready_candidate", got)
	}
}

func TestShortVideoCommercialEvidenceContractIndexMirrorsForbiddenActionTerminalChecks(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	expectedChecks := requireObject(t, goLive, "forbidden_action_terminal_checks")
	indexChecks := requireObject(t, index, "forbidden_action_terminal_checks")
	if len(indexChecks) != len(expectedChecks) {
		t.Fatalf("forbidden_action_terminal_checks len = %d, want %d", len(indexChecks), len(expectedChecks))
	}
	for checkKey := range expectedChecks {
		t.Run(checkKey, func(t *testing.T) {
			expected := requireObject(t, expectedChecks, checkKey)
			actual := requireObject(t, indexChecks, checkKey)
			for _, stringKey := range []string{"current_result", "evidence_required"} {
				if got, want := requireString(t, actual, stringKey), requireString(t, expected, stringKey); got != want {
					t.Fatalf("%s.%s = %s, want %s", checkKey, stringKey, got, want)
				}
			}
			if got, want := requireBoolValue(t, actual, "required_final_value"), requireBoolValue(t, expected, "required_final_value"); got != want {
				t.Fatalf("%s.required_final_value = %v, want %v", checkKey, got, want)
			}
		})
	}

	indexGate := requireObject(t, index, "owner_signoff_gate")
	for checkKey := range expectedChecks {
		requireStringSliceContains(t, asStringSlice(t, indexGate["blocking_terminal_checks"]), checkKey)
	}
}

func TestShortVideoCommercialEvidenceContractIndexPreservesExecutionOrderAndDependencies(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	goLiveOrder := asStringSlice(t, goLive["required_slices"])
	indexOrder := asStringSlice(t, index["execution_order"])
	requireStringSlicesEqual(t, indexOrder, goLiveOrder, "commercial evidence execution_order")

	indexSlices := asObjectSlice(t, index["required_slices"])
	if len(indexSlices) != len(goLiveOrder) {
		t.Fatalf("index required_slices len = %d, want %d", len(indexSlices), len(goLiveOrder))
	}
	for i, sliceKey := range goLiveOrder {
		if got := requireString(t, indexSlices[i], "slice_key"); got != sliceKey {
			t.Fatalf("index required_slices[%d].slice_key = %s, want %s", i, got, sliceKey)
		}

		sourceSlice := requireObject(t, candidateSet, sliceKey)
		indexSlice := findObjectByString(t, indexSlices, "slice_key", sliceKey)
		requireStringSlicesEqual(t, stringSliceOrEmpty(t, indexSlice, "depends_on"), stringSliceOrEmpty(t, sourceSlice, "depends_on"), sliceKey+".depends_on")
	}

	completionGate := requireObject(t, index, "completion_gate")
	if got := requireString(t, completionGate, "required_order_source"); got != "commercial_go_live_evidence_package.required_slices" {
		t.Fatalf("completion_gate.required_order_source = %s, want commercial_go_live_evidence_package.required_slices", got)
	}
}

func TestShortVideoCommercialEvidenceContractIndexRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	gatePath := requireString(t, goLive, "commercial_go_no_go_gate")
	gate := readJSON(t, filepath.Join(base, gatePath))
	writebackGate := requireObject(t, gate, "evidence_writeback_gate")

	completionGate := requireObject(t, index, "completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_status_before_completion"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("completion_gate.required_status_before_completion must require all_evidence_writebacks_completed_and_verified")
	}
	if got := requireString(t, index, "commercial_go_no_go_gate"); got != gatePath {
		t.Fatalf("commercial_go_no_go_gate = %s, want %s", got, gatePath)
	}
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")
	requireBool(t, writebackGate, "can_pass_gate", false)
}

func TestShortVideoCommercialCompletionGatesRequireOwnerBaseGateReceipts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	cases := []struct {
		label        string
		path         string
		gateKey      string
		requiredText string
	}{
		{
			label:        "commercial evidence index",
			path:         requireString(t, goLive, "evidence_contract_index"),
			gateKey:      "completion_gate",
			requiredText: "required_status_before_completion",
		},
		{
			label:        "commercial current state audit",
			path:         requireString(t, goLive, "current_state_audit"),
			gateKey:      "commercial_completion_gate",
			requiredText: "required_status_before_completion",
		},
		{
			label:        "commercial cross repo execution queue",
			path:         requireString(t, goLive, "cross_repo_execution_queue"),
			gateKey:      "completion_gate",
			requiredText: "required_before_commercial_ready",
		},
		{
			label:        "frontend backend acceptance",
			path:         requireString(t, goLive, "frontend_backend_acceptance_contract"),
			gateKey:      "completion_gate",
			requiredText: "required_status_before_completion",
		},
		{
			label:        "frontend backend handoff",
			path:         requireString(t, goLive, "frontend_backend_handoff_runbook"),
			gateKey:      "completion_gate",
			requiredText: "required_status_before_acceptance",
		},
		{
			label:        "commercial improvement backlog",
			path:         requireString(t, goLive, "improvement_backlog"),
			gateKey:      "completion_gate",
			requiredText: "required_before_backlog_done",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			doc := readJSON(t, filepath.Join(base, tc.path))
			gate := requireObject(t, doc, tc.gateKey)
			if got := requireString(t, gate, "required_owner_base_gate_source"); got != "Owner + Base Gate + Gateway + Receipt" {
				t.Fatalf("%s.%s.required_owner_base_gate_source = %s, want Owner + Base Gate + Gateway + Receipt", tc.path, tc.gateKey, got)
			}
			requireBool(t, gate, "owner_base_gate_receipts_must_exist", true)
			if got := requireString(t, gate, "blocking_owner_base_gate"); got != "owner_base_gate_receipts_missing" {
				t.Fatalf("%s.%s.blocking_owner_base_gate = %s, want owner_base_gate_receipts_missing", tc.path, tc.gateKey, got)
			}
			if !strings.Contains(requireString(t, gate, tc.requiredText), "owner_base_gate_receipts_bound") {
				t.Fatalf("%s.%s.%s must require owner_base_gate_receipts_bound", tc.path, tc.gateKey, tc.requiredText)
			}
			requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "owner_base_gate_receipts_missing")
		})
	}
}

func TestShortVideoCommercialReadinessAuditIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p14 := requireObject(t, candidateSet, "p14_commercial_readiness_audit")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gate := requireObject(t, goLive, "owner_signoff_gate")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))

	auditPath := requireString(t, p14, "machine_audit")
	if got := requireString(t, goLive, "current_state_audit"); got != auditPath {
		t.Fatalf("commercial_go_live_evidence_package.current_state_audit = %s, want %s", got, auditPath)
	}
	if got := requireString(t, index, "current_state_audit"); got != auditPath {
		t.Fatalf("commercial evidence index current_state_audit = %s, want %s", got, auditPath)
	}
	requireExistingPath(t, auditPath, base)
	audit := readJSON(t, filepath.Join(base, auditPath))

	requireBool(t, audit, "candidate_only", true)
	requireBool(t, audit, "non_formal", true)
	if got := requireString(t, audit, "completion_status"); got != "not_commercial_ready" {
		t.Fatalf("completion_status = %s, want not_commercial_ready", got)
	}
	if got := requireString(t, audit, "evidence_contract_index"); got != requireString(t, goLive, "evidence_contract_index") {
		t.Fatalf("evidence_contract_index = %s, want %s", got, requireString(t, goLive, "evidence_contract_index"))
	}
	if got := requireString(t, audit, "cross_repo_execution_queue"); got != requireString(t, goLive, "cross_repo_execution_queue") {
		t.Fatalf("cross_repo_execution_queue = %s, want %s", got, requireString(t, goLive, "cross_repo_execution_queue"))
	}
	if got := requireString(t, audit, "markdown_audit"); got != requireString(t, p14, "audit_file") {
		t.Fatalf("markdown_audit = %s, want %s", got, requireString(t, p14, "audit_file"))
	}
	if got := requireString(t, audit, "next_required_authorization"); got != requireString(t, gate, "required_next_authorization") {
		t.Fatalf("next_required_authorization = %s, want %s", got, requireString(t, gate, "required_next_authorization"))
	}
	if got := requireString(t, audit, "next_authorization_card"); got != requireString(t, gate, "next_authorization_card") {
		t.Fatalf("next_authorization_card = %s, want %s", got, requireString(t, gate, "next_authorization_card"))
	}
	if !strings.Contains(requireString(t, audit, "objective_source"), "gui-capability-pack-workbench-github-oss-test-plan-20260704.md") {
		t.Fatalf("objective_source must reference GUI capability pack test plan")
	}

	requiredSlices := asStringSlice(t, goLive["required_slices"])
	auditSlices := asObjectSlice(t, audit["required_slices"])
	if len(auditSlices) != len(requiredSlices) {
		t.Fatalf("audit required_slices len = %d, want %d", len(auditSlices), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		auditSlice := auditSlices[i]
		if got := requireString(t, auditSlice, "slice_key"); got != sliceKey {
			t.Fatalf("audit required_slices[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		if got, want := requireString(t, auditSlice, "authorization_scope_contract"), requireString(t, sourceSlice, "authorization_scope_contract"); got != want {
			t.Fatalf("%s authorization_scope_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, auditSlice, "authorization_evidence_intake_contract"), requireString(t, sourceSlice, "authorization_evidence_intake_contract"); got != want {
			t.Fatalf("%s authorization_evidence_intake_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, auditSlice, "execution_readiness_package"), requireString(t, sourceSlice, "execution_readiness_package"); got != want {
			t.Fatalf("%s execution_readiness_package = %s, want %s", sliceKey, got, want)
		}
		intake := readJSON(t, filepath.Join(base, requireString(t, auditSlice, "authorization_evidence_intake_contract")))
		intakeEvidence := requireObject(t, intake, "current_authorization_evidence")
		if got := requireString(t, intakeEvidence, "status"); got != "missing" {
			t.Fatalf("%s intake current_authorization_evidence.status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, auditSlice, "authorization_status"); got != "missing" {
			t.Fatalf("%s authorization_status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, auditSlice, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireBool(t, auditSlice, "can_count_toward_commercial_ready", false)
	}

	auditChecks := requireObject(t, audit, "forbidden_action_terminal_checks")
	indexChecks := requireObject(t, index, "forbidden_action_terminal_checks")
	if len(auditChecks) != len(indexChecks) {
		t.Fatalf("audit forbidden_action_terminal_checks len = %d, want %d", len(auditChecks), len(indexChecks))
	}
	for checkKey := range indexChecks {
		check := requireObject(t, auditChecks, checkKey)
		if got := requireString(t, check, "current_result"); got != "pending" {
			t.Fatalf("%s current_result = %s, want pending", checkKey, got)
		}
		requireBool(t, check, "required_final_value", false)
		requireBool(t, check, "can_count_toward_commercial_ready", false)
	}

	nonSufficient := asStringSlice(t, audit["non_sufficient_evidence"])
	for _, evidence := range []string{
		"candidate_assets_only",
		"green_packs_tests_only",
		"authorization_card_without_owner_quote",
		"manual_ledger_without_machine_evidence",
	} {
		requireStringSliceContains(t, nonSufficient, evidence)
	}

	completionGate := requireObject(t, audit, "commercial_completion_gate")
	requireBool(t, completionGate, "can_mark_commercial_ready", false)
	if got := requireString(t, completionGate, "required_status_before_completion"); got == "" {
		t.Fatalf("commercial_completion_gate.required_status_before_completion missing")
	}
}

func TestShortVideoCommercialReadinessAuditRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	p14 := requireObject(t, candidateSet, "p14_commercial_readiness_audit")
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gatePath := requireString(t, goLive, "commercial_go_no_go_gate")
	verifierPath := requireString(t, goLive, "commercial_readiness_verifier")
	queuePath := requireString(t, goLive, "cross_repo_execution_queue")

	gate := readJSON(t, filepath.Join(base, gatePath))
	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	if pending := int(requireNumber(t, writebackGate, "total_pending_entry_count")); pending <= 0 {
		t.Fatalf("evidence_writeback_gate.total_pending_entry_count = %d, want pending entries before execution", pending)
	}

	auditPath := requireString(t, p14, "audit_file")
	requireExistingPath(t, auditPath, "")
	if strings.HasPrefix(auditPath, "process_worktree_ref://") {
		// 同 requireExistingPath 的方案 b：外部过程 worktree 文档不随本仓分发，
		// 内容级断言在其缺席机器上无法执行，降级为引用登记存在性（上一行已校验）。
		return
	}
	raw, err := os.ReadFile(resolveReferencePath(auditPath))
	if err != nil {
		t.Fatalf("read %s: %v", auditPath, err)
	}
	content := string(raw)
	for _, required := range []string{
		gatePath,
		gatePath + "#evidence_writeback_gate",
		verifierPath,
		verifierPath + "#current_blockers.evidence_writeback_summary",
		queuePath,
		queuePath + "#execution_entries[].evidence_writeback_plan_summary",
		"evidence_writeback_gate_must_pass=true",
		"all_evidence_writebacks_completed_and_verified",
		"total_pending_entry_count=0",
	} {
		if !strings.Contains(content, required) {
			t.Fatalf("%s missing required readiness audit marker %q", auditPath, required)
		}
	}
}

func TestShortVideoCommercialCurrentStateAuditMirrorsEvidenceWritebackSummary(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))

	verifierBlockers := requireObject(t, verifier, "current_blockers")
	verifierSummary := requireObject(t, verifierBlockers, "evidence_writeback_summary")
	auditSummary := requireObject(t, audit, "evidence_writeback_summary")

	for _, field := range []string{
		"source_cross_repo_execution_queue",
		"writeback_status",
		"non_completion_clause",
	} {
		if got, want := requireString(t, auditSummary, field), requireString(t, verifierSummary, field); got != want {
			t.Fatalf("audit evidence_writeback_summary.%s = %s, want %s", field, got, want)
		}
	}
	for _, field := range []string{
		"total_required_entry_count",
		"total_pending_entry_count",
		"total_completed_entry_count",
	} {
		if got, want := int(requireNumber(t, auditSummary, field)), int(requireNumber(t, verifierSummary, field)); got != want {
			t.Fatalf("audit evidence_writeback_summary.%s = %d, want %d", field, got, want)
		}
	}
	requireBool(t, auditSummary, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, audit["non_sufficient_evidence"]), "evidence_writeback_incomplete")
	if !strings.Contains(requireString(t, requireObject(t, audit, "commercial_completion_gate"), "required_status_before_completion"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("commercial_completion_gate.required_status_before_completion must require all_evidence_writebacks_completed_and_verified")
	}

	verifierSlices := asObjectSlice(t, verifierSummary["per_slice"])
	auditSlices := asObjectSlice(t, auditSummary["per_slice"])
	if len(auditSlices) != len(verifierSlices) {
		t.Fatalf("audit evidence_writeback_summary.per_slice len = %d, want %d", len(auditSlices), len(verifierSlices))
	}
	for i, verifierSlice := range verifierSlices {
		auditSlice := auditSlices[i]
		if got, want := requireString(t, auditSlice, "slice_key"), requireString(t, verifierSlice, "slice_key"); got != want {
			t.Fatalf("audit evidence_writeback_summary.per_slice[%d].slice_key = %s, want %s", i, got, want)
		}
		for _, field := range []string{
			"readiness_package",
			"ledger_path",
			"evidence_contract",
			"write_status",
		} {
			if got, want := requireString(t, auditSlice, field), requireString(t, verifierSlice, field); got != want {
				t.Fatalf("%s audit evidence_writeback_summary.%s = %s, want %s", requireString(t, verifierSlice, "slice_key"), field, got, want)
			}
		}
		for _, field := range []string{
			"required_entry_count",
			"pending_entry_count",
			"completed_entry_count",
		} {
			if got, want := int(requireNumber(t, auditSlice, field)), int(requireNumber(t, verifierSlice, field)); got != want {
				t.Fatalf("%s audit evidence_writeback_summary.%s = %d, want %d", requireString(t, verifierSlice, "slice_key"), field, got, want)
			}
		}
		requireBool(t, auditSlice, "can_count_toward_commercial_ready", false)
	}
}

func TestShortVideoCommercialCurrentStateAuditCompletionGateRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionGate := requireObject(t, audit, "commercial_completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("commercial_completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_status_before_completion"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("commercial_completion_gate.required_status_before_completion must require all_evidence_writebacks_completed_and_verified")
	}
	requireStringSliceContains(t, asStringSlice(t, audit["non_sufficient_evidence"]), "evidence_writeback_incomplete")
}

func TestShortVideoCommercialReadinessVerifierIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	gate := requireObject(t, goLive, "owner_signoff_gate")

	verifierPath := requireString(t, goLive, "commercial_readiness_verifier")
	if got := requireString(t, index, "commercial_readiness_verifier"); got != verifierPath {
		t.Fatalf("commercial evidence index commercial_readiness_verifier = %s, want %s", got, verifierPath)
	}
	requireExistingPath(t, verifierPath, base)
	verifier := readJSON(t, filepath.Join(base, verifierPath))
	queuePath := requireString(t, goLive, "cross_repo_execution_queue")
	queue := readJSON(t, filepath.Join(base, queuePath))

	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	if got := requireString(t, verifier, "verification_status"); got != "blocked_not_commercial_ready" {
		t.Fatalf("verification_status = %s, want blocked_not_commercial_ready", got)
	}
	if got := requireString(t, verifier, "verified_current_go_live_status"); got != requireString(t, goLive, "current_go_live_status") {
		t.Fatalf("verified_current_go_live_status = %s, want %s", got, requireString(t, goLive, "current_go_live_status"))
	}
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	requireBool(t, verifier, "formal_write_allowed", false)
	if got := requireString(t, verifier, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, verifier, "evidence_writeback_gate_must_pass", true)
	requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), "evidence_writeback_incomplete")
	if got := requireString(t, verifier, "source_evidence_contract_index"); got != requireString(t, goLive, "evidence_contract_index") {
		t.Fatalf("source_evidence_contract_index = %s, want %s", got, requireString(t, goLive, "evidence_contract_index"))
	}
	if got := requireString(t, verifier, "source_current_state_audit"); got != requireString(t, goLive, "current_state_audit") {
		t.Fatalf("source_current_state_audit = %s, want %s", got, requireString(t, goLive, "current_state_audit"))
	}
	if got := requireString(t, verifier, "next_required_authorization"); got != requireString(t, gate, "required_next_authorization") {
		t.Fatalf("next_required_authorization = %s, want %s", got, requireString(t, gate, "required_next_authorization"))
	}
	if got := requireString(t, verifier, "next_authorization_card"); got != requireString(t, gate, "next_authorization_card") {
		t.Fatalf("next_authorization_card = %s, want %s", got, requireString(t, gate, "next_authorization_card"))
	}

	blockers := requireObject(t, verifier, "current_blockers")
	sliceBlockers := asObjectSlice(t, blockers["required_slices"])
	requiredSlices := asStringSlice(t, goLive["required_slices"])
	if len(sliceBlockers) != len(requiredSlices) {
		t.Fatalf("verifier current_blockers.required_slices len = %d, want %d", len(sliceBlockers), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		blocker := sliceBlockers[i]
		if got := requireString(t, blocker, "slice_key"); got != sliceKey {
			t.Fatalf("verifier current_blockers.required_slices[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		for _, field := range []string{
			"authorization_evidence_intake_contract",
			"authorization_scope_contract",
			"execution_readiness_package",
			"evidence_contract",
			"evidence_ledger",
		} {
			if got, want := requireString(t, blocker, field), requireString(t, sourceSlice, field); got != want {
				t.Fatalf("%s %s = %s, want %s", sliceKey, field, got, want)
			}
		}
		if got := requireString(t, blocker, "authorization_status"); got != "missing" {
			t.Fatalf("%s authorization_status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, blocker, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireBool(t, blocker, "can_count_toward_commercial_ready", false)
	}

	queueEntries := asObjectSlice(t, queue["execution_entries"])
	if len(queueEntries) != len(requiredSlices) {
		t.Fatalf("queue execution_entries len = %d, want %d", len(queueEntries), len(requiredSlices))
	}
	totalRequired := 0
	totalPending := 0
	totalCompleted := 0
	for _, queueEntry := range queueEntries {
		summary := requireObject(t, queueEntry, "evidence_writeback_plan_summary")
		totalRequired += int(requireNumber(t, summary, "required_entry_count"))
		totalPending += int(requireNumber(t, summary, "pending_entry_count"))
		totalCompleted += int(requireNumber(t, summary, "completed_entry_count"))
	}
	writebackBlocker := requireObject(t, blockers, "evidence_writeback_summary")
	if got := requireString(t, writebackBlocker, "source_cross_repo_execution_queue"); got != queuePath {
		t.Fatalf("evidence_writeback_summary.source_cross_repo_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, writebackBlocker, "writeback_status"); got != "blocked_pending_authorization_and_evidence" {
		t.Fatalf("evidence_writeback_summary.writeback_status = %s, want blocked_pending_authorization_and_evidence", got)
	}
	requireBool(t, writebackBlocker, "can_count_toward_commercial_ready", false)
	if got := int(requireNumber(t, writebackBlocker, "total_required_entry_count")); got != totalRequired {
		t.Fatalf("evidence_writeback_summary.total_required_entry_count = %d, want %d", got, totalRequired)
	}
	if got := int(requireNumber(t, writebackBlocker, "total_pending_entry_count")); got != totalPending {
		t.Fatalf("evidence_writeback_summary.total_pending_entry_count = %d, want %d", got, totalPending)
	}
	if got := int(requireNumber(t, writebackBlocker, "total_completed_entry_count")); got != totalCompleted {
		t.Fatalf("evidence_writeback_summary.total_completed_entry_count = %d, want %d", got, totalCompleted)
	}
	writebackSlices := asObjectSlice(t, writebackBlocker["per_slice"])
	if len(writebackSlices) != len(requiredSlices) {
		t.Fatalf("evidence_writeback_summary.per_slice len = %d, want %d", len(writebackSlices), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		queueEntry := queueEntries[i]
		queueSummary := requireObject(t, queueEntry, "evidence_writeback_plan_summary")
		writebackSlice := writebackSlices[i]
		if got := requireString(t, writebackSlice, "slice_key"); got != sliceKey {
			t.Fatalf("evidence_writeback_summary.per_slice[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		for _, field := range []string{
			"readiness_package",
			"ledger_path",
			"evidence_contract",
			"write_status",
		} {
			if got, want := requireString(t, writebackSlice, field), requireString(t, queueSummary, field); got != want {
				t.Fatalf("%s evidence_writeback_summary.%s = %s, want %s", sliceKey, field, got, want)
			}
		}
		for _, field := range []string{
			"required_entry_count",
			"pending_entry_count",
			"completed_entry_count",
		} {
			if got, want := int(requireNumber(t, writebackSlice, field)), int(requireNumber(t, queueSummary, field)); got != want {
				t.Fatalf("%s evidence_writeback_summary.%s = %d, want %d", sliceKey, field, got, want)
			}
		}
		requireBool(t, writebackSlice, "can_count_toward_commercial_ready", false)
	}

	terminalChecks := requireObject(t, blockers, "forbidden_action_terminal_checks")
	goLiveChecks := requireObject(t, goLive, "forbidden_action_terminal_checks")
	if len(terminalChecks) != len(goLiveChecks) {
		t.Fatalf("verifier forbidden_action_terminal_checks len = %d, want %d", len(terminalChecks), len(goLiveChecks))
	}
	for checkKey := range goLiveChecks {
		check := requireObject(t, terminalChecks, checkKey)
		if got := requireString(t, check, "current_result"); got != "pending" {
			t.Fatalf("%s current_result = %s, want pending", checkKey, got)
		}
		requireBool(t, check, "required_final_value", false)
		requireBool(t, check, "can_count_toward_commercial_ready", false)
	}

	requiredBeforeReady := asStringSlice(t, verifier["required_before_commercial_ready"])
	for _, required := range []string{
		"owner_authorization_recorded_for_all_required_slices",
		"machine_evidence_complete_for_all_required_slices",
		"forbidden_action_terminal_checks_proven_false",
		"owner_base_gate_receipts_bound",
		"frontend_backend_e2e_evidence_present",
		"all_evidence_writebacks_completed_and_verified",
	} {
		requireStringSliceContains(t, requiredBeforeReady, required)
	}
}

func TestShortVideoCommercialFrontendBackendAcceptanceContractIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))

	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	if got := requireString(t, index, "frontend_backend_acceptance_contract"); got != contractPath {
		t.Fatalf("commercial evidence index frontend_backend_acceptance_contract = %s, want %s", got, contractPath)
	}
	if got := requireString(t, verifier, "frontend_backend_acceptance_contract"); got != contractPath {
		t.Fatalf("commercial readiness verifier frontend_backend_acceptance_contract = %s, want %s", got, contractPath)
	}
	requireExistingPath(t, contractPath, base)
	contract := readJSON(t, filepath.Join(base, contractPath))

	requireBool(t, contract, "candidate_only", true)
	requireBool(t, contract, "non_formal", true)
	if got := requireString(t, contract, "acceptance_status"); got != "blocked_pending_cross_repo_authorization_and_evidence" {
		t.Fatalf("acceptance_status = %s, want blocked_pending_cross_repo_authorization_and_evidence", got)
	}
	requireBool(t, contract, "can_mark_frontend_backend_commercial_complete", false)
	requireBool(t, contract, "formal_write_allowed", false)
	if got := requireString(t, contract, "source_evidence_contract_index"); got != requireString(t, goLive, "evidence_contract_index") {
		t.Fatalf("source_evidence_contract_index = %s, want %s", got, requireString(t, goLive, "evidence_contract_index"))
	}
	for _, repo := range []string{"truzhenos", "truzhen-client-web-desktop"} {
		requireStringSliceContains(t, asStringSlice(t, contract["required_repositories"]), repo)
	}
	for _, forbiddenRepo := range []string{"truzhen-contracts", "truzhen-software", "truzhen-cloud"} {
		requireStringSliceContains(t, asStringSlice(t, contract["disallowed_without_new_authorization"]), forbiddenRepo)
	}

	requiredSourceSlices := asStringSlice(t, contract["required_source_slices"])
	for _, sliceKey := range []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
	} {
		requireStringSliceContains(t, requiredSourceSlices, sliceKey)
	}

	requireAcceptanceEvidenceItems(t, contract, "required_backend_endpoints", []string{
		"draft_from_delivery",
		"readiness_aggregate",
		"promote_candidate",
		"confirm_with_gate_receipt",
		"commercial_readiness_verifier_api",
	})
	requireAcceptanceEvidenceItems(t, contract, "required_frontend_surfaces", []string{
		"candidate_bundle_panel",
		"delivery_status_panel",
		"readiness_issue_panel",
		"promote_confirm_controls",
		"commercial_readiness_verifier_panel",
	})
	requireAcceptanceEvidenceItems(t, contract, "required_e2e_scenarios", []string{
		"safe_fixture_lifecycle_round_trip",
		"three_short_video_candidates_gui_walkthrough",
		"blocked_actions_visible_and_disabled",
	})
	requireEvidenceItems(t, contract, "required_forbidden_action_checks", []string{
		"no_real_codex_cli_run",
		"no_third_party_oss_execution",
		"no_social_login_or_upload",
		"no_raw_secret_saved",
	})

	gate := requireObject(t, contract, "completion_gate")
	requireBool(t, gate, "can_mark_complete", false)
	if got := requireString(t, gate, "required_status_before_completion"); got == "" {
		t.Fatalf("completion_gate.required_status_before_completion missing")
	}
	if got := requireString(t, gate, "commercial_status_after_frontend_backend_only"); got != "not_commercial_ready" {
		t.Fatalf("commercial_status_after_frontend_backend_only = %s, want not_commercial_ready", got)
	}
}

func TestShortVideoFrontendBackendAcceptanceRequiresGuiApiReceiptTraceability(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	contract := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_acceptance_contract")))
	runbook := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_handoff_runbook")))

	requireBool(t, contract, "frontend_screenshot_only_can_mark_accepted", false)
	requireBool(t, contract, "backend_test_only_can_mark_accepted", false)
	requireBool(t, contract, "network_summary_without_receipt_lookup_can_mark_accepted", false)

	expectedTraceIDs := []string{
		"candidate_bundle_panel",
		"delivery_status_panel",
		"readiness_issue_panel",
		"promote_confirm_controls",
		"commercial_readiness_verifier_panel",
		"safe_fixture_lifecycle_round_trip",
		"three_short_video_candidates_gui_walkthrough",
		"blocked_actions_visible_and_disabled",
	}
	traceRows := asObjectSlice(t, contract["gui_api_receipt_traceability_matrix"])
	if len(traceRows) != len(expectedTraceIDs) {
		t.Fatalf("gui_api_receipt_traceability_matrix len = %d, want %d", len(traceRows), len(expectedTraceIDs))
	}
	traceByID := map[string]map[string]any{}
	for _, row := range traceRows {
		id := requireString(t, row, "gui_evidence_id")
		traceByID[id] = row
		for _, key := range []string{
			"required_backend_evidence_id",
			"required_api_endpoint_or_readmodel",
			"required_receipt_or_candidate_ref_field",
			"required_field_match",
		} {
			if got := requireString(t, row, key); got == "" {
				t.Fatalf("%s.%s missing", id, key)
			}
		}
		if got := requireString(t, row, "current_status"); got != "pending_authorization" {
			t.Fatalf("%s.current_status = %s, want pending_authorization", id, got)
		}
		requireBool(t, row, "can_count_toward_frontend_backend_acceptance", false)
	}
	for _, id := range expectedTraceIDs {
		if _, ok := traceByID[id]; !ok {
			t.Fatalf("gui_api_receipt_traceability_matrix missing %s", id)
		}
	}

	completionGate := requireObject(t, contract, "completion_gate")
	for _, required := range []string{
		"gui_api_receipt_traceability_verified",
		"all_gui_steps_bound_to_backend_receipts",
		"all_backend_results_visible_in_gui",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_completion"]), required)
	}
	for _, blocked := range []string{
		"frontend_screenshot_without_api_receipt_trace",
		"backend_test_without_gui_api_trace",
		"network_summary_without_receipt_lookup",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["non_sufficient_evidence"]), blocked)
		requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), blocked)
	}

	if got := requireString(t, runbook, "gui_api_receipt_traceability_matrix_source"); got != "docs/commercial-frontend-backend-acceptance-contract.json#gui_api_receipt_traceability_matrix" {
		t.Fatalf("runbook.gui_api_receipt_traceability_matrix_source = %s", got)
	}
	runbookGate := requireObject(t, runbook, "completion_gate")
	requireStringSliceContains(t, asStringSlice(t, runbookGate["required_before_acceptance"]), "gui_api_receipt_traceability_verified")
	for _, blocked := range []string{
		"frontend_screenshot_without_api_receipt_trace",
		"backend_test_without_gui_api_trace",
		"network_summary_without_receipt_lookup",
	} {
		requireStringSliceContains(t, asStringSlice(t, runbookGate["non_sufficient_evidence"]), blocked)
		requireStringSliceContains(t, asStringSlice(t, runbook["non_sufficient_evidence"]), blocked)
	}
}

func TestShortVideoCommercialFrontendBackendAcceptanceRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	contract := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_acceptance_contract")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionGate := requireObject(t, contract, "completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_status_before_completion"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("completion_gate.required_status_before_completion must require all_evidence_writebacks_completed_and_verified")
	}
	requireStringSliceContains(t, asStringSlice(t, contract["non_sufficient_evidence"]), "evidence_writeback_incomplete")
}

func TestShortVideoCommercialFrontendBackendHandoffRunbookMatchesAcceptanceContract(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	contract := readJSON(t, filepath.Join(base, contractPath))
	runbookPath := requireString(t, goLive, "frontend_backend_handoff_runbook")
	for label, doc := range map[string]map[string]any{
		"commercial evidence index":     index,
		"commercial readiness verifier": verifier,
		"current state audit":           audit,
		"commercial go/no-go gate":      gate,
		"frontend/backend contract":     contract,
	} {
		if got := requireString(t, doc, "frontend_backend_handoff_runbook"); got != runbookPath {
			t.Fatalf("%s frontend_backend_handoff_runbook = %s, want %s", label, got, runbookPath)
		}
	}
	requireExistingPath(t, runbookPath, base)
	runbook := readJSON(t, filepath.Join(base, runbookPath))

	requireBool(t, runbook, "candidate_only", true)
	requireBool(t, runbook, "non_formal", true)
	if got := requireString(t, runbook, "runbook_status"); got != "blocked_pending_cross_repo_authorization_and_evidence" {
		t.Fatalf("runbook_status = %s, want blocked_pending_cross_repo_authorization_and_evidence", got)
	}
	requireBool(t, runbook, "can_start_handoff", false)
	if got := requireString(t, runbook, "source_frontend_backend_acceptance_contract"); got != contractPath {
		t.Fatalf("source_frontend_backend_acceptance_contract = %s, want %s", got, contractPath)
	}

	expectedPhases := []struct {
		key         string
		sliceKey    string
		group       string
		evidenceIDs []string
	}{
		{
			key:         "p12_backend_safe_lifecycle",
			sliceKey:    "p12_safe_lifecycle_sample",
			group:       "required_backend_endpoints",
			evidenceIDs: collectEvidenceIDsFromItems(t, contract["required_backend_endpoints"]),
		},
		{
			key:         "p13_frontend_lifecycle_panel",
			sliceKey:    "p13_gui_lifecycle_panel",
			group:       "required_frontend_surfaces",
			evidenceIDs: collectEvidenceIDsFromItems(t, contract["required_frontend_surfaces"]),
		},
		{
			key:         "p15_e2e_walkthrough",
			sliceKey:    "p15_gui_walkthrough_three_candidates",
			group:       "required_e2e_scenarios",
			evidenceIDs: collectEvidenceIDsFromItems(t, contract["required_e2e_scenarios"]),
		},
		{
			key:         "terminal_forbidden_checks",
			sliceKey:    "p15_gui_walkthrough_three_candidates",
			group:       "required_forbidden_action_checks",
			evidenceIDs: collectEvidenceIDsFromItems(t, contract["required_forbidden_action_checks"]),
		},
	}
	phases := asObjectSlice(t, runbook["handoff_phases"])
	if len(phases) != len(expectedPhases) {
		t.Fatalf("handoff_phases len = %d, want %d", len(phases), len(expectedPhases))
	}
	for i, expected := range expectedPhases {
		phase := phases[i]
		if got := requireString(t, phase, "phase_key"); got != expected.key {
			t.Fatalf("handoff_phases[%d].phase_key = %s, want %s", i, got, expected.key)
		}
		if got := requireString(t, phase, "required_source_slice"); got != expected.sliceKey {
			t.Fatalf("%s required_source_slice = %s, want %s", expected.key, got, expected.sliceKey)
		}
		if got := requireString(t, phase, "source_evidence_group"); got != expected.group {
			t.Fatalf("%s source_evidence_group = %s, want %s", expected.key, got, expected.group)
		}
		requireStringSlicesEqual(t, asStringSlice(t, phase["evidence_ids"]), expected.evidenceIDs, expected.key+" evidence_ids")
		if got := requireString(t, phase, "owner_authorization_status"); got != "missing" {
			t.Fatalf("%s owner_authorization_status = %s, want missing", expected.key, got)
		}
		if got := requireString(t, phase, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", expected.key, got)
		}
		requireBool(t, phase, "can_execute", false)
		if len(asStringSlice(t, phase["writeback_ledgers"])) == 0 {
			t.Fatalf("%s writeback_ledgers missing", expected.key)
		}
		requireStringSliceContains(t, asStringSlice(t, phase["forbidden_actions"]), "mark_commercial_ready_from_handoff_runbook")
	}

	completionGate := requireObject(t, runbook, "completion_gate")
	requireBool(t, completionGate, "can_mark_frontend_backend_accepted", false)
	requireStringSliceContains(t, asStringSlice(t, runbook["non_sufficient_evidence"]), "handoff_runbook_without_cross_repo_execution_evidence")
}

func TestShortVideoCommercialFrontendBackendHandoffRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	runbook := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_handoff_runbook")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionGate := requireObject(t, runbook, "completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_status_before_acceptance"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("completion_gate.required_status_before_acceptance must require all_evidence_writebacks_completed_and_verified")
	}
	requireStringSliceContains(t, asStringSlice(t, runbook["non_sufficient_evidence"]), "evidence_writeback_incomplete")
}

func TestShortVideoCommercialFrontendBackendHandoffRunbookCarriesExpectedCommands(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	contract := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_acceptance_contract")))
	runbook := readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_handoff_runbook")))

	expectedGroups := map[string]string{
		"p12_backend_safe_lifecycle":   "required_backend_endpoints",
		"p13_frontend_lifecycle_panel": "required_frontend_surfaces",
		"p15_e2e_walkthrough":          "required_e2e_scenarios",
		"terminal_forbidden_checks":    "required_forbidden_action_checks",
	}
	for phaseKey, groupKey := range expectedGroups {
		phase := findObjectByString(t, asObjectSlice(t, runbook["handoff_phases"]), "phase_key", phaseKey)
		plan := requireObject(t, phase, "command_plan")
		requireBool(t, plan, "requires_owner_authorization", true)
		requireBool(t, plan, "can_run_now", false)
		if got := requireString(t, plan, "source_evidence_group"); got != groupKey {
			t.Fatalf("%s command_plan.source_evidence_group = %s, want %s", phaseKey, got, groupKey)
		}
		commands := asObjectSlice(t, plan["expected_commands"])
		contractItems := asObjectSlice(t, contract[groupKey])
		if len(commands) != len(contractItems) {
			t.Fatalf("%s expected_commands len = %d, want %d", phaseKey, len(commands), len(contractItems))
		}
		for i, source := range contractItems {
			command := commands[i]
			evidenceID := requireString(t, source, "evidence_id")
			if got := requireString(t, command, "evidence_id"); got != evidenceID {
				t.Fatalf("%s command[%d].evidence_id = %s, want %s", phaseKey, i, got, evidenceID)
			}
			if got := requireString(t, command, "expected_command"); got != requireString(t, source, "expected_command") {
				t.Fatalf("%s %s expected_command = %s, want %s", phaseKey, evidenceID, got, requireString(t, source, "expected_command"))
			}
			if got := requireString(t, command, "target_repository"); got != requireString(t, source, "target_repository") {
				t.Fatalf("%s %s target_repository = %s, want %s", phaseKey, evidenceID, got, requireString(t, source, "target_repository"))
			}
			if got := requireString(t, command, "evidence_required"); got != requireString(t, source, "evidence_required") {
				t.Fatalf("%s %s evidence_required = %s, want %s", phaseKey, evidenceID, got, requireString(t, source, "evidence_required"))
			}
			if got := requireString(t, command, "current_status"); got != "pending_authorization" {
				t.Fatalf("%s %s current_status = %s, want pending_authorization", phaseKey, evidenceID, got)
			}
		}
	}
}

func TestShortVideoCommercialCrossRepoExecutionQueueIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	gate := requireObject(t, goLive, "owner_signoff_gate")

	queuePath := requireString(t, goLive, "cross_repo_execution_queue")
	if got := requireString(t, index, "cross_repo_execution_queue"); got != queuePath {
		t.Fatalf("commercial evidence index cross_repo_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, verifier, "cross_repo_execution_queue"); got != queuePath {
		t.Fatalf("commercial readiness verifier cross_repo_execution_queue = %s, want %s", got, queuePath)
	}
	requireExistingPath(t, queuePath, base)
	queue := readJSON(t, filepath.Join(base, queuePath))

	requireBool(t, queue, "candidate_only", true)
	requireBool(t, queue, "non_formal", true)
	if got := requireString(t, queue, "queue_status"); got != "blocked_pending_owner_authorization_and_evidence" {
		t.Fatalf("queue_status = %s, want blocked_pending_owner_authorization_and_evidence", got)
	}
	requireBool(t, queue, "can_start_cross_repo_execution", false)
	if got := requireString(t, queue, "next_executable_slice"); got != requireString(t, gate, "required_next_authorization") {
		t.Fatalf("next_executable_slice = %s, want %s", got, requireString(t, gate, "required_next_authorization"))
	}
	if got := requireString(t, queue, "next_authorization_card"); got != requireString(t, gate, "next_authorization_card") {
		t.Fatalf("next_authorization_card = %s, want %s", got, requireString(t, gate, "next_authorization_card"))
	}
	if got := requireString(t, queue, "source_evidence_contract_index"); got != requireString(t, goLive, "evidence_contract_index") {
		t.Fatalf("source_evidence_contract_index = %s, want %s", got, requireString(t, goLive, "evidence_contract_index"))
	}
	if got := requireString(t, queue, "source_commercial_readiness_verifier"); got != requireString(t, goLive, "commercial_readiness_verifier") {
		t.Fatalf("source_commercial_readiness_verifier = %s, want %s", got, requireString(t, goLive, "commercial_readiness_verifier"))
	}
	if got := requireString(t, queue, "source_current_state_audit"); got != requireString(t, goLive, "current_state_audit") {
		t.Fatalf("source_current_state_audit = %s, want %s", got, requireString(t, goLive, "current_state_audit"))
	}
	authorizationAttemptVerifierPath := requireString(t, goLive, "authorization_attempt_coverage_verifier")
	if got := requireString(t, queue, "authorization_attempt_coverage_verifier"); got != authorizationAttemptVerifierPath {
		t.Fatalf("authorization_attempt_coverage_verifier = %s, want %s", got, authorizationAttemptVerifierPath)
	}
	requireStringSlicesEqual(t, asStringSlice(t, queue["execution_order"]), asStringSlice(t, goLive["required_slices"]), "execution_order")

	entries := asObjectSlice(t, queue["execution_entries"])
	requiredSlices := asStringSlice(t, goLive["required_slices"])
	if len(entries) != len(requiredSlices) {
		t.Fatalf("execution_entries len = %d, want %d", len(entries), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		entry := entries[i]
		if got := requireString(t, entry, "slice_key"); got != sliceKey {
			t.Fatalf("execution_entries[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		readinessPath := requireString(t, sourceSlice, "execution_readiness_package")
		readiness := readJSON(t, filepath.Join(base, readinessPath))
		for _, field := range []string{
			"authorization_scope_contract",
			"authorization_evidence_intake_contract",
			"execution_readiness_package",
			"evidence_contract",
			"evidence_ledger",
			"implementation_spec",
			"cross_repo_authorization_card",
		} {
			if got, want := requireString(t, entry, field), requireString(t, sourceSlice, field); got != want {
				t.Fatalf("%s %s = %s, want %s", sliceKey, field, got, want)
			}
		}
		if got := requireString(t, entry, "readiness_status"); got != requireString(t, readiness, "readiness_status") {
			t.Fatalf("%s readiness_status = %s, want %s", sliceKey, got, requireString(t, readiness, "readiness_status"))
		}
		writebackPlan := requireObject(t, readiness, "evidence_writeback_plan")
		writebackSummary := requireObject(t, entry, "evidence_writeback_plan_summary")
		if got := requireString(t, writebackSummary, "readiness_package"); got != readinessPath {
			t.Fatalf("%s evidence_writeback_plan_summary.readiness_package = %s, want %s", sliceKey, got, readinessPath)
		}
		if got := requireString(t, writebackSummary, "ledger_path"); got != requireString(t, writebackPlan, "ledger_path") {
			t.Fatalf("%s evidence_writeback_plan_summary.ledger_path = %s, want %s", sliceKey, got, requireString(t, writebackPlan, "ledger_path"))
		}
		if got := requireString(t, writebackSummary, "evidence_contract"); got != requireString(t, writebackPlan, "evidence_contract") {
			t.Fatalf("%s evidence_writeback_plan_summary.evidence_contract = %s, want %s", sliceKey, got, requireString(t, writebackPlan, "evidence_contract"))
		}
		if got := requireString(t, writebackSummary, "write_status"); got != requireString(t, writebackPlan, "write_status") {
			t.Fatalf("%s evidence_writeback_plan_summary.write_status = %s, want %s", sliceKey, got, requireString(t, writebackPlan, "write_status"))
		}
		requireBool(t, writebackSummary, "can_write_completion_claim", false)
		requiredEntries := asObjectSlice(t, writebackPlan["required_entries"])
		if got := int(requireNumber(t, writebackSummary, "required_entry_count")); got != len(requiredEntries) {
			t.Fatalf("%s evidence_writeback_plan_summary.required_entry_count = %d, want %d", sliceKey, got, len(requiredEntries))
		}
		if got := int(requireNumber(t, writebackSummary, "pending_entry_count")); got != len(requiredEntries) {
			t.Fatalf("%s evidence_writeback_plan_summary.pending_entry_count = %d, want %d", sliceKey, got, len(requiredEntries))
		}
		if got := int(requireNumber(t, writebackSummary, "completed_entry_count")); got != 0 {
			t.Fatalf("%s evidence_writeback_plan_summary.completed_entry_count = %d, want 0", sliceKey, got)
		}
		requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_cross_repo_run"]), "owner_authorization_missing")
		for _, required := range asStringSlice(t, writebackPlan["required_before_marking_slice_complete"]) {
			requireStringSliceContains(t, asStringSlice(t, writebackSummary["blocking_before_slice_completion"]), required)
		}
		requireStringSlicesEqual(t, asStringSlice(t, entry["target_repositories"]), asStringSlice(t, sourceSlice["target_repositories"]), sliceKey+" target_repositories")
		if rawDeps, ok := sourceSlice["depends_on"]; ok {
			requireStringSlicesEqual(t, asStringSlice(t, entry["depends_on"]), asStringSlice(t, rawDeps), sliceKey+" depends_on")
		}
		requireBool(t, entry, "can_start_now", false)
		if got := requireString(t, entry, "owner_authorization_status"); got != "missing" {
			t.Fatalf("%s owner_authorization_status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, entry, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		entryGate := requireObject(t, entry, "cross_repo_work_gate")
		readinessGate := requireObject(t, readiness, "cross_repo_work_gate")
		requireBool(t, entryGate, "can_start_cross_repo_work", false)
		if got, want := requireString(t, entryGate, "required_status_before_cross_repo_work"), requireString(t, readinessGate, "required_status_before_cross_repo_work"); got != want {
			t.Fatalf("%s required_status_before_cross_repo_work = %s, want %s", sliceKey, got, want)
		}
	}

	completionGate := requireObject(t, queue, "completion_gate")
	requireBool(t, completionGate, "can_mark_commercial_execution_ready", false)
	if got := requireString(t, completionGate, "required_authorization_attempt_coverage_verifier_source"); got != "commercial_go_live_evidence_package.authorization_attempt_coverage_verifier" {
		t.Fatalf("completion_gate.required_authorization_attempt_coverage_verifier_source = %s, want commercial_go_live_evidence_package.authorization_attempt_coverage_verifier", got)
	}
	requireBool(t, completionGate, "authorization_attempt_coverage_verifier_must_pass", true)
	if got := requireString(t, completionGate, "required_before_first_cross_repo_run"); got == "" {
		t.Fatalf("completion_gate.required_before_first_cross_repo_run missing")
	} else if !strings.Contains(got, "authorization_attempt_coverage_verifier_passed") {
		t.Fatalf("completion_gate.required_before_first_cross_repo_run = %s, want authorization_attempt_coverage_verifier_passed", got)
	}
	requireStringSliceContains(t, asStringSlice(t, completionGate["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "queue_defined_without_cross_repo_receipts")
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "authorization_attempt_coverage_verifier_pending")
}

func TestShortVideoCommercialCrossRepoExecutionQueueRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionGate := requireObject(t, queue, "completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_before_commercial_ready"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("completion_gate.required_before_commercial_ready must require all_evidence_writebacks_completed_and_verified")
	}
	requireStringSliceContains(t, asStringSlice(t, queue["non_sufficient_evidence"]), "evidence_writeback_incomplete")
}

func TestShortVideoCommercialImprovementBacklogIsMachineCheckable(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))

	backlogPath := requireString(t, goLive, "improvement_backlog")
	for label, doc := range map[string]map[string]any{
		"commercial evidence index":     index,
		"commercial readiness verifier": verifier,
		"current state audit":           audit,
		"cross repo execution queue":    queue,
	} {
		if got := requireString(t, doc, "improvement_backlog"); got != backlogPath {
			t.Fatalf("%s improvement_backlog = %s, want %s", label, got, backlogPath)
		}
	}
	requireExistingPath(t, backlogPath, base)
	backlog := readJSON(t, filepath.Join(base, backlogPath))

	requireBool(t, backlog, "candidate_only", true)
	requireBool(t, backlog, "non_formal", true)
	if got := requireString(t, backlog, "backlog_status"); got != "blocked_pending_cross_repo_authorization_and_evidence" {
		t.Fatalf("backlog_status = %s, want blocked_pending_cross_repo_authorization_and_evidence", got)
	}
	if got := requireString(t, backlog, "source_cross_repo_execution_queue"); got != requireString(t, goLive, "cross_repo_execution_queue") {
		t.Fatalf("source_cross_repo_execution_queue = %s, want %s", got, requireString(t, goLive, "cross_repo_execution_queue"))
	}
	if got := requireString(t, backlog, "source_frontend_backend_acceptance_contract"); got != requireString(t, goLive, "frontend_backend_acceptance_contract") {
		t.Fatalf("source_frontend_backend_acceptance_contract = %s, want %s", got, requireString(t, goLive, "frontend_backend_acceptance_contract"))
	}
	if got := requireString(t, backlog, "source_commercial_readiness_verifier"); got != requireString(t, goLive, "commercial_readiness_verifier") {
		t.Fatalf("source_commercial_readiness_verifier = %s, want %s", got, requireString(t, goLive, "commercial_readiness_verifier"))
	}

	cardIDsBySlice := map[string]string{
		"p12_safe_lifecycle_sample":            "CAP-STUDIO-P12-LIFECYCLE-SAFE-SAMPLE",
		"p13_gui_lifecycle_panel":              "CAP-STUDIO-P13-GUI-LIFECYCLE-PANEL",
		"p15_gui_walkthrough_three_candidates": "CAP-STUDIO-P15-THREE-CANDIDATE-GUI-WALKTHROUGH",
		"p16_controlled_code_assistant_run":    "CAP-STUDIO-P16-CONTROLLED-CODE-ASSISTANT-RUN",
		"p17_provider_adapter_candidate":       "CAP-STUDIO-P17-PROVIDER-ADAPTER-CANDIDATE",
		"p18_cloud_market_sandbox":             "CAP-STUDIO-P18-CLOUD-MARKET-SANDBOX",
	}
	cards := asObjectSlice(t, backlog["improvement_cards"])
	requiredSlices := asStringSlice(t, goLive["required_slices"])
	if len(cards) != len(requiredSlices) {
		t.Fatalf("improvement_cards len = %d, want %d", len(cards), len(requiredSlices))
	}
	for _, sliceKey := range requiredSlices {
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		card := findObjectByString(t, cards, "source_slice", sliceKey)
		if got := requireString(t, card, "card_id"); got != cardIDsBySlice[sliceKey] {
			t.Fatalf("%s card_id = %s, want %s", sliceKey, got, cardIDsBySlice[sliceKey])
		}
		for _, field := range []string{
			"execution_readiness_package",
			"evidence_contract",
			"cross_repo_authorization_card",
		} {
			if got, want := requireString(t, card, field), requireString(t, sourceSlice, field); got != want {
				t.Fatalf("%s %s = %s, want %s", sliceKey, field, got, want)
			}
		}
		requireStringSlicesEqual(t, asStringSlice(t, card["target_repositories"]), asStringSlice(t, sourceSlice["target_repositories"]), sliceKey+" target_repositories")
		if got := requireString(t, card, "owner_authorization_status"); got != "missing" {
			t.Fatalf("%s owner_authorization_status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, card, "implementation_status"); got != "not_started" {
			t.Fatalf("%s implementation_status = %s, want not_started", sliceKey, got)
		}
		if got := requireString(t, card, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireStringIn(t, requireString(t, card, "risk_color"), "yellow", "orange", "red")
		requireStringSliceContains(t, asStringSlice(t, card["blocked_by"]), "owner_authorization_missing")
		if len(asStringSlice(t, card["acceptance_evidence_required"])) == 0 {
			t.Fatalf("%s acceptance_evidence_required missing", sliceKey)
		}
		requireStringSliceContains(t, asStringSlice(t, card["forbidden_actions"]), "mark_commercial_ready_from_candidate_assets_only")
	}

	completionGate := requireObject(t, backlog, "completion_gate")
	requireBool(t, completionGate, "can_mark_backlog_done", false)
	requireStringSliceContains(t, asStringSlice(t, backlog["non_sufficient_evidence"]), "candidate_cards_without_cross_repo_evidence")
}

func TestShortVideoPackStudioIssueLedgerFeedsCommercialBacklogAndCompletionGates(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")
	requireExistingPath(t, issueLedgerPath, base)
	issueLedger := readJSON(t, filepath.Join(base, issueLedgerPath))

	commercialDocs := map[string]map[string]any{
		"commercial improvement backlog":       readJSON(t, filepath.Join(base, requireString(t, goLive, "improvement_backlog"))),
		"commercial readiness verifier":        readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier"))),
		"commercial current state audit":       readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit"))),
		"commercial go/no-go gate":             readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate"))),
		"goal completion evidence map":         readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map"))),
		"frontend backend handoff runbook":     readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_handoff_runbook"))),
		"frontend backend acceptance contract": readJSON(t, filepath.Join(base, requireString(t, goLive, "frontend_backend_acceptance_contract"))),
	}
	for label, doc := range commercialDocs {
		if got := requireString(t, doc, "pack_studio_issue_ledger"); got != issueLedgerPath {
			t.Fatalf("%s pack_studio_issue_ledger = %s, want %s", label, got, issueLedgerPath)
		}
	}

	requireBool(t, issueLedger, "candidate_only", true)
	requireBool(t, issueLedger, "non_formal", true)
	if got := requireString(t, issueLedger, "issue_ledger_status"); got != "blocked_pending_issue_resolution_evidence" {
		t.Fatalf("issue_ledger_status = %s, want blocked_pending_issue_resolution_evidence", got)
	}
	if got := requireString(t, issueLedger, "source_plan_requirement"); got != "M7_issue_backfill" {
		t.Fatalf("source_plan_requirement = %s, want M7_issue_backfill", got)
	}
	requireBool(t, issueLedger, "can_mark_commercial_ready", false)

	expectedIssues := map[string]string{
		"CAP-STUDIO-ISSUE-P10-CAPABILITY-PACK-LOADER-MISSING":      "capability_pack_loader_missing",
		"CAP-STUDIO-ISSUE-P11-CANDIDATE-BUNDLE-NOT-DELIVERY":       "candidate_bundle_not_delivery",
		"CAP-STUDIO-ISSUE-P12-OWNER-AUTHORIZATION-MISSING":         "owner_authorization_missing",
		"CAP-STUDIO-ISSUE-P13-GUI-LIFECYCLE-EVIDENCE-MISSING":      "frontend_backend_evidence_pending",
		"CAP-STUDIO-ISSUE-P15-THREE-CANDIDATE-WALKTHROUGH-MISSING": "three_candidate_gui_walkthrough_missing",
		"CAP-STUDIO-ISSUE-P16-CODE-ASSISTANT-RUN-NOT-AUTHORIZED":   "code_assistant_run_not_authorized",
		"CAP-STUDIO-ISSUE-P17-PROVIDER-ADAPTER-OWNERSHIP-MISSING":  "provider_adapter_ownership_missing",
		"CAP-STUDIO-ISSUE-P18-CLOUD-SANDBOX-NOT-AUTHORIZED":        "cloud_sandbox_not_authorized",
	}
	issues := asObjectSlice(t, issueLedger["issues"])
	if len(issues) != len(expectedIssues) {
		t.Fatalf("issues len = %d, want %d", len(issues), len(expectedIssues))
	}
	for issueID, blockedReason := range expectedIssues {
		issue := findObjectByString(t, issues, "issue_id", issueID)
		if got := requireString(t, issue, "blocked_reason"); got != blockedReason {
			t.Fatalf("%s blocked_reason = %s, want %s", issueID, got, blockedReason)
		}
		requireStringIn(t, requireString(t, issue, "risk_color"), "yellow", "orange", "red")
		requireBool(t, issue, "can_count_toward_commercial_ready", false)
		if got := requireString(t, issue, "status"); !strings.HasPrefix(got, "open_") && !strings.HasPrefix(got, "blocked_") {
			t.Fatalf("%s status = %s, want open_* or blocked_*", issueID, got)
		}
		sourceArtifact := requireString(t, issue, "source_artifact")
		requireExistingPath(t, sourceArtifact, base)
		if len(asStringSlice(t, issue["required_resolution_evidence"])) == 0 {
			t.Fatalf("%s required_resolution_evidence missing", issueID)
		}
		requireStringSliceContains(t, asStringSlice(t, issue["forbidden_resolutions"]), "manual_edit_without_gui_or_gateway_evidence")
	}

	completionGate := requireObject(t, issueLedger, "completion_gate")
	requireBool(t, completionGate, "can_close_issue_ledger", false)
	requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_close"]), "all_issue_cards_resolved_with_evidence")
	requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_commercial_ready"]), "pack_studio_issue_ledger_closed")
	requireStringSliceContains(t, asStringSlice(t, issueLedger["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
}

func TestShortVideoPackStudioIssueLedgerRequiresAuthoritativeResolutionEvidenceRefs(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	issueLedger := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_issue_ledger")))
	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	machinePackage := readJSON(t, filepath.Join(base, requireString(t, goLive, "machine_evidence_package")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	expectedTypes := []string{
		"implementation_or_authorization_evidence",
		"verification_command_evidence",
		"gui_api_or_readmodel_evidence",
		"receipt_or_candidate_ref_evidence",
		"forbidden_resolution_absence_evidence",
	}

	schema := requireObject(t, issueLedger, "resolution_evidence_schema")
	requireBool(t, schema, "manual_resolution_note_can_close_issue", false)
	requireBool(t, schema, "summary_only_can_close_issue", false)
	for _, evidenceType := range expectedTypes {
		requireStringSliceContains(t, asStringSlice(t, schema["required_evidence_types"]), evidenceType)
	}

	for _, issue := range asObjectSlice(t, issueLedger["issues"]) {
		issueID := requireString(t, issue, "issue_id")
		if got := requireString(t, issue, "resolution_evidence_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("%s resolution_evidence_status = %s, want missing_authoritative_evidence", issueID, got)
		}
		requireStringSliceContains(t, asStringSlice(t, issue["forbidden_resolutions"]), "close_issue_with_manual_summary_only")
		requireStringSliceContains(t, asStringSlice(t, issue["forbidden_resolutions"]), "close_issue_without_receipt_or_candidate_ref")

		requirements := asObjectSlice(t, issue["resolution_evidence_requirements"])
		if len(requirements) < len(expectedTypes) {
			t.Fatalf("%s resolution_evidence_requirements len = %d, want at least %d", issueID, len(requirements), len(expectedTypes))
		}
		for _, evidenceType := range expectedTypes {
			req := findObjectByString(t, requirements, "evidence_type", evidenceType)
			if got := requireString(t, req, "evidence_status"); got != "missing" {
				t.Fatalf("%s %s evidence_status = %s, want missing", issueID, evidenceType, got)
			}
			for _, key := range []string{"required_source", "authority", "required_ref_field"} {
				if got := requireString(t, req, key); got == "" {
					t.Fatalf("%s %s %s missing", issueID, evidenceType, key)
				}
			}
			requireBool(t, req, "can_substitute_with_manual_note", false)
		}
	}

	completionGate := requireObject(t, issueLedger, "completion_gate")
	for _, required := range []string{
		"all_issue_resolution_evidence_schema_verified",
		"all_issue_resolution_evidence_refs_authoritative",
		"manual_resolution_notes_rejected",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_close"]), required)
	}
	for _, nonSufficient := range []string{
		"manual_resolution_note_without_authoritative_refs",
		"closed_status_without_resolution_evidence_refs",
		"tests_without_gui_or_receipt_evidence",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, issueLedger["non_sufficient_evidence"]), nonSufficient)
	}

	for label, values := range map[string][]string{
		"readiness":       asStringSlice(t, readiness["required_before_commercial_ready"]),
		"go/no-go":        asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]),
		"machine package": asStringSlice(t, machinePackage["required_before_commercial_ready"]),
		"goal map":        asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]),
	} {
		if label == "" {
			t.Fatalf("unreachable")
		}
		requireStringSliceContains(t, values, "pack_studio_issue_resolution_evidence_schema_verified")
	}
	for label, doc := range map[string]map[string]any{
		"readiness":       readiness,
		"go/no-go":        goNoGo,
		"machine package": machinePackage,
		"goal map":        goalMap,
	} {
		for _, nonSufficient := range []string{
			"manual_resolution_note_without_authoritative_refs",
			"closed_status_without_resolution_evidence_refs",
			"tests_without_gui_or_receipt_evidence",
		} {
			if label == "" {
				t.Fatalf("unreachable")
			}
			requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), nonSufficient)
		}
	}
}

func TestShortVideoPackStudioIssueLedgerBlocksCommercialSummaryCompletionGates(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")
	specs := []struct {
		label      string
		path       string
		gateKey    string
		statusKey  string
		statusVerb string
	}{
		{
			label:      "commercial current state audit",
			path:       requireString(t, goLive, "current_state_audit"),
			gateKey:    "commercial_completion_gate",
			statusKey:  "required_status_before_completion",
			statusVerb: "commercial ready",
		},
		{
			label:      "frontend backend acceptance contract",
			path:       requireString(t, goLive, "frontend_backend_acceptance_contract"),
			gateKey:    "completion_gate",
			statusKey:  "required_status_before_completion",
			statusVerb: "complete",
		},
		{
			label:      "frontend backend handoff runbook",
			path:       requireString(t, goLive, "frontend_backend_handoff_runbook"),
			gateKey:    "completion_gate",
			statusKey:  "required_status_before_acceptance",
			statusVerb: "accepted",
		},
	}
	for _, spec := range specs {
		t.Run(spec.label, func(t *testing.T) {
			doc := readJSON(t, filepath.Join(base, spec.path))
			if got := requireString(t, doc, "pack_studio_issue_ledger"); got != issueLedgerPath {
				t.Fatalf("%s pack_studio_issue_ledger = %s, want %s", spec.label, got, issueLedgerPath)
			}
			gate := requireObject(t, doc, spec.gateKey)
			if got := requireString(t, gate, "required_pack_studio_issue_ledger_source"); got != "candidate_set.pack_studio_issue_ledger" {
				t.Fatalf("%s %s.required_pack_studio_issue_ledger_source = %s, want candidate_set.pack_studio_issue_ledger", spec.label, spec.gateKey, got)
			}
			requireBool(t, gate, "pack_studio_issue_ledger_must_close", true)
			if got := requireString(t, gate, spec.statusKey); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
				t.Fatalf("%s %s.%s missing pack_studio_issue_ledger_closed before %s: %s", spec.label, spec.gateKey, spec.statusKey, spec.statusVerb, got)
			}
			requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
			requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
		})
	}
}

func TestShortVideoPackStudioIssueLedgerBlocksCommercialSignoffAndEvidenceIndex(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")
	if got := requireString(t, goLive, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("commercial_go_live_evidence_package.pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "owner_signoff_gate")["blocking_required_before_owner_signoff"]), "pack_studio_issue_ledger_closed")
	signoffMatrix := requireObject(t, goLive, "commercial_signoff_matrix")
	if got := requireString(t, signoffMatrix, "required_before_owner_signoff"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("commercial_signoff_matrix.required_before_owner_signoff missing pack_studio_issue_ledger_closed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, signoffMatrix["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	if got := requireString(t, index, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("commercial evidence index pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	indexCompletionGate := requireObject(t, index, "completion_gate")
	if got := requireString(t, indexCompletionGate, "required_pack_studio_issue_ledger_source"); got != "candidate_set.pack_studio_issue_ledger" {
		t.Fatalf("completion_gate.required_pack_studio_issue_ledger_source = %s, want candidate_set.pack_studio_issue_ledger", got)
	}
	requireBool(t, indexCompletionGate, "pack_studio_issue_ledger_must_close", true)
	if got := requireString(t, indexCompletionGate, "required_status_before_completion"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("completion_gate.required_status_before_completion missing pack_studio_issue_ledger_closed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, indexCompletionGate["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, index, "owner_signoff_gate")["blocking_required_before_owner_signoff"]), "pack_studio_issue_ledger_closed")
	indexSignoffMatrix := requireObject(t, index, "commercial_signoff_matrix")
	if got := requireString(t, indexSignoffMatrix, "required_before_owner_signoff"); !strings.Contains(got, "pack_studio_issue_ledger_closed") {
		t.Fatalf("index commercial_signoff_matrix.required_before_owner_signoff missing pack_studio_issue_ledger_closed: %s", got)
	}
	requireStringSliceContains(t, asStringSlice(t, indexSignoffMatrix["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	if got := requireString(t, goNoGo, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("commercial go/no-go pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "pack_studio_issue_ledger_closed")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")

	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))
	if got := requireString(t, goalMap, "pack_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("goal completion evidence map pack_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "pack_studio_issue_ledger_open")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "pack_studio_issue_ledger_closed")
	requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), "pack_studio_issue_ledger_open")
}

func TestShortVideoCommercialImprovementBacklogRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	backlog := readJSON(t, filepath.Join(base, requireString(t, goLive, "improvement_backlog")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionGate := requireObject(t, backlog, "completion_gate")
	if got := requireString(t, completionGate, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("completion_gate.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, completionGate, "evidence_writeback_gate_must_pass", true)
	if !strings.Contains(requireString(t, completionGate, "required_before_backlog_done"), "all_evidence_writebacks_completed_and_verified") {
		t.Fatalf("completion_gate.required_before_backlog_done must require all_evidence_writebacks_completed_and_verified")
	}
	requireStringSliceContains(t, asStringSlice(t, backlog["non_sufficient_evidence"]), "evidence_writeback_incomplete")
}

func TestShortVideoCommercialGoNoGoGateBlocksReadyClaimUntilAllEvidencePasses(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	queue := readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue")))
	backlog := readJSON(t, filepath.Join(base, requireString(t, goLive, "improvement_backlog")))

	gatePath := requireString(t, goLive, "commercial_go_no_go_gate")
	for label, doc := range map[string]map[string]any{
		"commercial evidence index":      index,
		"commercial readiness verifier":  verifier,
		"current state audit":            audit,
		"cross repo execution queue":     queue,
		"commercial improvement backlog": backlog,
	} {
		if got := requireString(t, doc, "commercial_go_no_go_gate"); got != gatePath {
			t.Fatalf("%s commercial_go_no_go_gate = %s, want %s", label, got, gatePath)
		}
	}
	requireExistingPath(t, gatePath, base)
	gate := readJSON(t, filepath.Join(base, gatePath))

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	if got := requireString(t, gate, "decision_status"); got != "blocked_not_ready_for_commercial_go_live" {
		t.Fatalf("decision_status = %s, want blocked_not_ready_for_commercial_go_live", got)
	}
	requireBool(t, gate, "can_mark_commercial_ready", false)
	requireBool(t, gate, "can_request_owner_go_live_signoff", false)
	if got := requireString(t, gate, "source_commercial_readiness_verifier"); got != requireString(t, goLive, "commercial_readiness_verifier") {
		t.Fatalf("source_commercial_readiness_verifier = %s, want %s", got, requireString(t, goLive, "commercial_readiness_verifier"))
	}
	if got := requireString(t, gate, "source_improvement_backlog"); got != requireString(t, goLive, "improvement_backlog") {
		t.Fatalf("source_improvement_backlog = %s, want %s", got, requireString(t, goLive, "improvement_backlog"))
	}
	if got := requireString(t, gate, "next_required_authorization"); got != requireString(t, verifier, "next_required_authorization") {
		t.Fatalf("next_required_authorization = %s, want %s", got, requireString(t, verifier, "next_required_authorization"))
	}
	if got := requireString(t, gate, "next_authorization_card"); got != requireString(t, verifier, "next_authorization_card") {
		t.Fatalf("next_authorization_card = %s, want %s", got, requireString(t, verifier, "next_authorization_card"))
	}

	requiredSlices := asStringSlice(t, goLive["required_slices"])
	sliceGates := asObjectSlice(t, gate["required_slice_gates"])
	if len(sliceGates) != len(requiredSlices) {
		t.Fatalf("required_slice_gates len = %d, want %d", len(sliceGates), len(requiredSlices))
	}
	backlogCards := asObjectSlice(t, backlog["improvement_cards"])
	for i, sliceKey := range requiredSlices {
		sliceGate := sliceGates[i]
		if got := requireString(t, sliceGate, "slice_key"); got != sliceKey {
			t.Fatalf("required_slice_gates[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		card := findObjectByString(t, backlogCards, "source_slice", sliceKey)
		if got, want := requireString(t, sliceGate, "source_improvement_card"), requireString(t, card, "card_id"); got != want {
			t.Fatalf("%s source_improvement_card = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, sliceGate, "evidence_contract"), requireString(t, sourceSlice, "evidence_contract"); got != want {
			t.Fatalf("%s evidence_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, sliceGate, "execution_readiness_package"), requireString(t, sourceSlice, "execution_readiness_package"); got != want {
			t.Fatalf("%s execution_readiness_package = %s, want %s", sliceKey, got, want)
		}
		if got := requireString(t, sliceGate, "owner_authorization_status"); got != "missing" {
			t.Fatalf("%s owner_authorization_status = %s, want missing", sliceKey, got)
		}
		if got := requireString(t, sliceGate, "implementation_status"); got != "not_started" {
			t.Fatalf("%s implementation_status = %s, want not_started", sliceKey, got)
		}
		if got := requireString(t, sliceGate, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireBool(t, sliceGate, "can_pass_gate", false)
		requireStringSliceContains(t, asStringSlice(t, sliceGate["blocked_by"]), "owner_authorization_missing")
	}

	indexChecks := requireObject(t, index, "forbidden_action_terminal_checks")
	terminalGates := asObjectSlice(t, gate["forbidden_action_terminal_gates"])
	if len(terminalGates) != len(indexChecks) {
		t.Fatalf("forbidden_action_terminal_gates len = %d, want %d", len(terminalGates), len(indexChecks))
	}
	for _, terminalGate := range terminalGates {
		checkKey := requireString(t, terminalGate, "check_key")
		sourceCheck := requireObject(t, indexChecks, checkKey)
		if got := requireString(t, terminalGate, "current_result"); got != requireString(t, sourceCheck, "current_result") {
			t.Fatalf("%s current_result = %s, want %s", checkKey, got, requireString(t, sourceCheck, "current_result"))
		}
		requireBool(t, terminalGate, "required_final_value", false)
		requireBool(t, terminalGate, "can_pass_gate", false)
		if got := requireString(t, terminalGate, "evidence_required"); got != requireString(t, sourceCheck, "evidence_required") {
			t.Fatalf("%s evidence_required = %s, want %s", checkKey, got, requireString(t, sourceCheck, "evidence_required"))
		}
	}

	verifierBlockers := requireObject(t, verifier, "current_blockers")
	verifierWriteback := requireObject(t, verifierBlockers, "evidence_writeback_summary")
	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	if got := requireString(t, writebackGate, "source_commercial_readiness_verifier"); got != requireString(t, goLive, "commercial_readiness_verifier") {
		t.Fatalf("evidence_writeback_gate.source_commercial_readiness_verifier = %s, want %s", got, requireString(t, goLive, "commercial_readiness_verifier"))
	}
	if got := requireString(t, writebackGate, "source_cross_repo_execution_queue"); got != requireString(t, verifierWriteback, "source_cross_repo_execution_queue") {
		t.Fatalf("evidence_writeback_gate.source_cross_repo_execution_queue = %s, want %s", got, requireString(t, verifierWriteback, "source_cross_repo_execution_queue"))
	}
	if got := requireString(t, writebackGate, "writeback_status"); got != requireString(t, verifierWriteback, "writeback_status") {
		t.Fatalf("evidence_writeback_gate.writeback_status = %s, want %s", got, requireString(t, verifierWriteback, "writeback_status"))
	}
	requireBool(t, writebackGate, "can_pass_gate", false)
	for _, field := range []string{
		"total_required_entry_count",
		"total_pending_entry_count",
		"total_completed_entry_count",
	} {
		if got, want := int(requireNumber(t, writebackGate, field)), int(requireNumber(t, verifierWriteback, field)); got != want {
			t.Fatalf("evidence_writeback_gate.%s = %d, want %d", field, got, want)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, writebackGate["blocked_by"]), "evidence_writeback_incomplete")
	requireStringSliceContains(t, asStringSlice(t, writebackGate["required_before_pass"]), "all_evidence_writebacks_completed_and_verified")

	completionRule := requireObject(t, gate, "completion_rule")
	requireBool(t, completionRule, "all_slice_gates_must_pass", true)
	requireBool(t, completionRule, "all_terminal_gates_must_pass", true)
	requireBool(t, completionRule, "evidence_writeback_gate_must_pass", true)
	requireBool(t, completionRule, "owner_base_gate_receipts_required", true)
	requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "commercial_go_no_go_gate_without_cross_repo_receipts")
}

func TestShortVideoCommercialGoNoGoRequiresOwnerBaseGateReceipts(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	gate := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	for label, doc := range map[string]map[string]any{
		"commercial_go_live_evidence_package": goLive,
		"commercial_readiness_verifier":       verifier,
		"commercial_go_no_go_gate":            gate,
	} {
		if got := requireString(t, doc, "required_owner_base_gate_source"); got != "Owner + Base Gate + Gateway + Receipt" {
			t.Fatalf("%s.required_owner_base_gate_source = %s, want Owner + Base Gate + Gateway + Receipt", label, got)
		}
		requireBool(t, doc, "owner_base_gate_receipts_must_exist", true)
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "owner_base_gate_receipts_missing")
	}

	completionRule := requireObject(t, gate, "completion_rule")
	requireBool(t, completionRule, "owner_base_gate_receipts_required", true)
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_go_live_signoff"]), "owner_base_gate_receipts_bound")
}

func TestShortVideoCommercialGoLiveEvidencePackageRequiresEvidenceWritebackGate(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	gatePath := requireString(t, goLive, "commercial_go_no_go_gate")
	verifierPath := requireString(t, goLive, "commercial_readiness_verifier")
	queuePath := requireString(t, goLive, "cross_repo_execution_queue")
	authorizationAttemptVerifierPath := requireString(t, goLive, "authorization_attempt_coverage_verifier")
	nextAuthorizationStartGuardPath := requireString(t, goLive, "next_authorization_start_guard")
	executionReadinessGuardCoverageVerifierPath := requireString(t, goLive, "execution_readiness_guard_coverage_verifier")
	issueLedgerPath := requireString(t, candidateSet, "pack_studio_issue_ledger")

	if got := requireString(t, goLive, "required_evidence_writeback_gate_source"); got != "commercial_go_no_go_gate.evidence_writeback_gate" {
		t.Fatalf("commercial_go_live_evidence_package.required_evidence_writeback_gate_source = %s, want commercial_go_no_go_gate.evidence_writeback_gate", got)
	}
	requireBool(t, goLive, "evidence_writeback_gate_must_pass", true)
	requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), "evidence_writeback_incomplete")

	gate := readJSON(t, filepath.Join(base, gatePath))
	writebackGate := requireObject(t, gate, "evidence_writeback_gate")
	requireBool(t, writebackGate, "can_pass_gate", false)
	if pending := int(requireNumber(t, writebackGate, "total_pending_entry_count")); pending <= 0 {
		t.Fatalf("evidence_writeback_gate.total_pending_entry_count = %d, want pending entries before execution", pending)
	}

	packagePath := filepath.Join("docs", "plans", "capability-pack-studio-commercial-go-live-evidence-package-20260704.md")
	requireExistingPath(t, packagePath, "")
	raw, err := os.ReadFile(packagePath)
	if err != nil {
		t.Fatalf("read %s: %v", packagePath, err)
	}
	content := string(raw)
	for _, required := range []string{
		gatePath,
		gatePath + "#evidence_writeback_gate",
		verifierPath,
		verifierPath + "#current_blockers.evidence_writeback_summary",
		queuePath,
		queuePath + "#execution_entries[].evidence_writeback_plan_summary",
		"evidence_writeback_gate_must_pass=true",
		"all_evidence_writebacks_completed_and_verified",
		"total_pending_entry_count=0",
		authorizationAttemptVerifierPath,
		"authorization_attempt_coverage_verifier_passed",
		nextAuthorizationStartGuardPath,
		"next_authorization_start_guard_open",
		executionReadinessGuardCoverageVerifierPath,
		"execution_readiness_guard_coverage_verifier_passed",
		issueLedgerPath,
		"pack_studio_issue_ledger_closed",
		"pack_studio_issue_ledger_open",
	} {
		if !strings.Contains(content, required) {
			t.Fatalf("%s missing required final evidence package marker %q", packagePath, required)
		}
	}
}

func TestShortVideoCommercialGoLiveEvidencePackageHasMachineReadableGateSummary(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	packagePath := requireString(t, goLive, "machine_evidence_package")
	requireExistingPath(t, packagePath, base)

	pkg := readJSON(t, filepath.Join(base, packagePath))
	if got := requireString(t, pkg, "evidence_package_ref"); got != "commercial-go-live-evidence-package://short-video-ops-v0" {
		t.Fatalf("machine evidence package ref = %s, want commercial-go-live-evidence-package://short-video-ops-v0", got)
	}
	requireBool(t, pkg, "candidate_only", true)
	requireBool(t, pkg, "non_formal", true)
	requireBool(t, pkg, "can_mark_commercial_ready", false)
	requireBool(t, pkg, "can_request_owner_signoff", false)
	if got := requireString(t, pkg, "current_go_live_status"); got != requireString(t, goLive, "current_go_live_status") {
		t.Fatalf("machine evidence package current_go_live_status = %s, want %s", got, requireString(t, goLive, "current_go_live_status"))
	}
	if got := requireString(t, pkg, "source_candidate_set"); got != "candidate-set.json#commercial_go_live_evidence_package" {
		t.Fatalf("machine evidence package source_candidate_set = %s, want candidate-set.json#commercial_go_live_evidence_package", got)
	}
	if got := requireString(t, pkg, "source_human_evidence_package"); got != requireString(t, goLive, "evidence_package") {
		t.Fatalf("machine evidence package source_human_evidence_package = %s, want %s", got, requireString(t, goLive, "evidence_package"))
	}

	sources := requireObject(t, pkg, "hard_gate_sources")
	for field, want := range map[string]string{
		"evidence_contract_index":                         requireString(t, goLive, "evidence_contract_index"),
		"commercial_readiness_verifier":                   requireString(t, goLive, "commercial_readiness_verifier"),
		"commercial_go_no_go_gate":                        requireString(t, goLive, "commercial_go_no_go_gate"),
		"post_run_gate_coverage_verifier":                 requireString(t, goLive, "post_run_gate_coverage_verifier"),
		"forbidden_action_coverage_verifier":              requireString(t, goLive, "forbidden_action_coverage_verifier"),
		"authorization_attempt_coverage_verifier":         requireString(t, goLive, "authorization_attempt_coverage_verifier"),
		"next_authorization_start_guard":                  requireString(t, goLive, "next_authorization_start_guard"),
		"execution_readiness_guard_coverage_verifier":     requireString(t, goLive, "execution_readiness_guard_coverage_verifier"),
		"pack_studio_issue_ledger":                        requireString(t, candidateSet, "pack_studio_issue_ledger"),
		"cross_repo_execution_queue":                      requireString(t, goLive, "cross_repo_execution_queue"),
		"required_evidence_writeback_gate_source":         "commercial_go_no_go_gate.evidence_writeback_gate",
		"required_owner_base_gate_source":                 "Owner + Base Gate + Gateway + Receipt",
		"required_commercial_signoff_matrix_source":       "commercial_go_live_evidence_package.commercial_signoff_matrix",
		"required_forbidden_action_terminal_check_source": "commercial_go_live_evidence_package.forbidden_action_terminal_checks",
	} {
		if got := requireString(t, sources, field); got != want {
			t.Fatalf("hard_gate_sources.%s = %s, want %s", field, got, want)
		}
	}

	required := asStringSlice(t, pkg["required_before_commercial_ready"])
	for _, marker := range []string{
		"all_required_slices_authorized_and_evidence_complete_verified",
		"post_run_gate_coverage_verifier_passed",
		"forbidden_action_coverage_verifier_passed",
		"authorization_attempt_coverage_verifier_passed",
		"next_authorization_start_guard_open",
		"execution_readiness_guard_coverage_verifier_passed",
		"pack_studio_issue_ledger_closed",
		"forbidden_action_terminal_checks_proven_false",
		"all_evidence_writebacks_completed_and_verified",
		"owner_base_gate_receipts_bound",
	} {
		requireStringSliceContains(t, required, marker)
	}

	nonSufficient := asStringSlice(t, pkg["non_sufficient_evidence"])
	for _, blocker := range []string{
		"authorization_cards_without_owner_authorization",
		"post_run_gate_coverage_verifier_pending",
		"forbidden_action_coverage_verifier_pending",
		"authorization_attempt_coverage_verifier_pending",
		"next_authorization_start_guard_pending",
		"execution_readiness_guard_coverage_verifier_pending",
		"pack_studio_issue_ledger_open",
		"evidence_writeback_incomplete",
		"owner_base_gate_receipts_missing",
	} {
		requireStringSliceContains(t, nonSufficient, blocker)
	}

	requiredSlices := asStringSlice(t, goLive["required_slices"])
	sliceMatrix := asObjectSlice(t, pkg["slice_evidence_matrix"])
	if len(sliceMatrix) != len(requiredSlices) {
		t.Fatalf("slice_evidence_matrix len = %d, want %d", len(sliceMatrix), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		entry := sliceMatrix[i]
		if got := requireString(t, entry, "slice_key"); got != sliceKey {
			t.Fatalf("slice_evidence_matrix[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		sourceSlice := requireObject(t, candidateSet, sliceKey)
		if got, want := requireString(t, entry, "evidence_contract"), requireString(t, sourceSlice, "evidence_contract"); got != want {
			t.Fatalf("%s evidence_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, entry, "evidence_ledger"), requireString(t, sourceSlice, "evidence_ledger"); got != want {
			t.Fatalf("%s evidence_ledger = %s, want %s", sliceKey, got, want)
		}
		if got := requireString(t, entry, "owner_authorization_status"); got != "pending_authorization" {
			t.Fatalf("%s owner_authorization_status = %s, want pending_authorization", sliceKey, got)
		}
		if got := requireString(t, entry, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireBool(t, entry, "can_count_toward_commercial_ready", false)
	}
}

func TestShortVideoMachineGoLiveEvidencePackageFeedsCommercialCompletionGates(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	machinePackagePath := requireString(t, goLive, "machine_evidence_package")

	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	for label, doc := range map[string]map[string]any{
		"commercial_readiness_verifier":        readiness,
		"commercial_go_no_go_gate":             goNoGo,
		"commercial_evidence_contract_index":   index,
		"commercial_readiness_current_audit":   audit,
		"pack_studio_goal_completion_evidence": goalMap,
	} {
		if got := requireString(t, doc, "machine_go_live_evidence_package"); got != machinePackagePath {
			t.Fatalf("%s machine_go_live_evidence_package = %s, want %s", label, got, machinePackagePath)
		}
	}

	readinessBlocker := requireObject(t, requireObject(t, readiness, "current_blockers"), "machine_go_live_evidence_package")
	if got := requireString(t, readinessBlocker, "source_machine_go_live_evidence_package"); got != machinePackagePath {
		t.Fatalf("readiness blocker source_machine_go_live_evidence_package = %s, want %s", got, machinePackagePath)
	}
	requireBool(t, readinessBlocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, readinessBlocker["blocked_by"]), "machine_go_live_evidence_package_blocked")

	for label, gate := range map[string]map[string]any{
		"go_no_go_completion_rule": requireObject(t, goNoGo, "completion_rule"),
		"index_completion_gate":    requireObject(t, index, "completion_gate"),
		"audit_completion_gate":    requireObject(t, audit, "commercial_completion_gate"),
	} {
		requireBool(t, gate, "machine_go_live_evidence_package_must_pass", true)
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_go_live_signoff"]), "machine_go_live_evidence_package_verified")
		if got := requireString(t, gate, "required_machine_go_live_evidence_package_source"); got != "commercial_go_live_evidence_package.machine_evidence_package" {
			t.Fatalf("%s required_machine_go_live_evidence_package_source = %s, want commercial_go_live_evidence_package.machine_evidence_package", label, got)
		}
	}

	for _, doc := range map[string]map[string]any{
		"commercial_go_no_go_gate":           goNoGo,
		"commercial_evidence_contract_index": index,
		"commercial_readiness_current_audit": audit,
		"pack_studio_goal_completion_map":    goalMap,
	} {
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "machine_go_live_evidence_package_blocked")
	}

	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "machine_go_live_evidence_package_blocked")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "machine_go_live_evidence_package_verified")
	req := findObjectByString(t, asObjectSlice(t, goalMap["active_goal_requirements"]), "requirement_id", "machine_go_live_evidence_package_verified")
	if got := requireString(t, req, "current_status"); got != "blocked_not_commercial_ready" {
		t.Fatalf("machine_go_live_evidence_package_verified current_status = %s, want blocked_not_commercial_ready", got)
	}
	requireBool(t, req, "can_count_toward_goal_completion", false)
}

func TestShortVideoMachineGoLiveEvidencePackageCoverageVerifierBlocksCommercialSignoff(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	verifierPath := requireString(t, goLive, "machine_evidence_package_coverage_verifier")
	machinePackagePath := requireString(t, goLive, "machine_evidence_package")
	requireExistingPath(t, verifierPath, base)

	verifier := readJSON(t, filepath.Join(base, verifierPath))
	if got := requireString(t, verifier, "verifier_ref"); got != "commercial-machine-go-live-evidence-package-coverage-verifier://short-video-ops-v0" {
		t.Fatalf("verifier_ref = %s, want commercial-machine-go-live-evidence-package-coverage-verifier://short-video-ops-v0", got)
	}
	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	if got := requireString(t, verifier, "coverage_status"); got != "blocked_pending_machine_go_live_evidence_package_verified" {
		t.Fatalf("coverage_status = %s, want blocked_pending_machine_go_live_evidence_package_verified", got)
	}

	for key, want := range map[string]string{
		"source_machine_go_live_evidence_package": requireString(t, goLive, "machine_evidence_package"),
		"source_commercial_readiness_verifier":    requireString(t, goLive, "commercial_readiness_verifier") + "#current_blockers.machine_go_live_evidence_package",
		"source_commercial_go_no_go_gate":         requireString(t, goLive, "commercial_go_no_go_gate") + "#completion_rule",
		"source_evidence_contract_index":          requireString(t, goLive, "evidence_contract_index") + "#completion_gate",
		"source_current_state_audit":              requireString(t, goLive, "current_state_audit") + "#commercial_completion_gate",
		"source_goal_completion_evidence_map":     requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map") + "#active_goal_requirements.machine_go_live_evidence_package_verified",
		"required_source_field":                   "commercial_go_live_evidence_package.machine_evidence_package",
		"required_pass_marker":                    "machine_go_live_evidence_package_verified",
		"required_blocker_marker":                 "machine_go_live_evidence_package_blocked",
	} {
		if got := requireString(t, verifier, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	requiredDocs := asStringSlice(t, verifier["required_commercial_gate_docs"])
	for _, docPath := range []string{
		requireString(t, goLive, "commercial_readiness_verifier"),
		requireString(t, goLive, "commercial_go_no_go_gate"),
		requireString(t, goLive, "evidence_contract_index"),
		requireString(t, goLive, "current_state_audit"),
		requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map"),
	} {
		requireStringSliceContains(t, requiredDocs, docPath)
	}

	checks := asObjectSlice(t, verifier["coverage_checks"])
	if len(checks) != len(requiredDocs) {
		t.Fatalf("coverage_checks len = %d, want %d", len(checks), len(requiredDocs))
	}
	for _, check := range checks {
		docPath := requireString(t, check, "doc_path")
		requireStringSliceContains(t, requiredDocs, docPath)
		requireBool(t, check, "doc_references_machine_go_live_evidence_package", true)
		requireBool(t, check, "completion_gate_requires_machine_go_live_evidence_package", true)
		requireBool(t, check, "completion_gate_blocks_machine_go_live_evidence_package", true)
		requireBool(t, check, "can_count_toward_commercial_ready", false)
	}

	for _, marker := range []string{
		"all_commercial_gate_docs_reference_machine_go_live_evidence_package",
		"all_commercial_completion_gates_require_machine_go_live_evidence_package",
		"machine_go_live_evidence_package_verified",
		"all_evidence_writebacks_completed_and_verified",
		"owner_base_gate_receipts_bound",
	} {
		requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), marker)
	}
	for _, blocker := range []string{
		"machine_evidence_package_coverage_verifier_pending",
		"machine_go_live_evidence_package_blocked",
		"evidence_writeback_incomplete",
		"owner_base_gate_receipts_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, verifier["non_sufficient_evidence"]), blocker)
		requireStringSliceContains(t, asStringSlice(t, goLive["non_sufficient_evidence"]), blocker)
	}

	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))
	if got := requireString(t, readiness, "machine_evidence_package_coverage_verifier"); got != verifierPath {
		t.Fatalf("readiness machine_evidence_package_coverage_verifier = %s, want %s", got, verifierPath)
	}
	if got := requireString(t, goNoGo, "machine_evidence_package_coverage_verifier"); got != verifierPath {
		t.Fatalf("go/no-go machine_evidence_package_coverage_verifier = %s, want %s", got, verifierPath)
	}
	if got := requireString(t, goalMap, "source_machine_evidence_package_coverage_verifier"); got != verifierPath {
		t.Fatalf("goal map source_machine_evidence_package_coverage_verifier = %s, want %s", got, verifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "owner_signoff_gate")["blocking_required_before_owner_signoff"]), "machine_evidence_package_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goLive, "commercial_signoff_matrix")["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]), "machine_evidence_package_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, goalMap["goal_completion_barriers"]), "machine_evidence_package_coverage_verifier_pending")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), "machine_evidence_package_coverage_verifier_passed")
	if got := requireString(t, verifier, "covered_machine_evidence_package"); got != machinePackagePath {
		t.Fatalf("covered_machine_evidence_package = %s, want %s", got, machinePackagePath)
	}
}

func TestShortVideoMachineGoLiveEvidencePackageCoverageVerifierFeedsSummaryCompletionGates(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	coverageVerifierPath := requireString(t, goLive, "machine_evidence_package_coverage_verifier")
	machinePackagePath := requireString(t, goLive, "machine_evidence_package")

	machinePackage := readJSON(t, filepath.Join(base, machinePackagePath))
	index := readJSON(t, filepath.Join(base, requireString(t, goLive, "evidence_contract_index")))
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))

	sources := requireObject(t, machinePackage, "hard_gate_sources")
	if got := requireString(t, sources, "machine_evidence_package_coverage_verifier"); got != coverageVerifierPath {
		t.Fatalf("machine package hard_gate_sources.machine_evidence_package_coverage_verifier = %s, want %s", got, coverageVerifierPath)
	}
	requireStringSliceContains(t, asStringSlice(t, machinePackage["required_before_commercial_ready"]), "machine_evidence_package_coverage_verifier_passed")
	requireStringSliceContains(t, asStringSlice(t, machinePackage["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")

	for label, doc := range map[string]map[string]any{
		"commercial_evidence_contract_index": index,
		"commercial_readiness_current_audit": audit,
	} {
		if got := requireString(t, doc, "machine_evidence_package_coverage_verifier"); got != coverageVerifierPath {
			t.Fatalf("%s machine_evidence_package_coverage_verifier = %s, want %s", label, got, coverageVerifierPath)
		}
	}

	for label, gate := range map[string]map[string]any{
		"index_completion_gate": requireObject(t, index, "completion_gate"),
		"audit_completion_gate": requireObject(t, audit, "commercial_completion_gate"),
	} {
		if got := requireString(t, gate, "required_machine_evidence_package_coverage_verifier_source"); got != "commercial_go_live_evidence_package.machine_evidence_package_coverage_verifier" {
			t.Fatalf("%s required_machine_evidence_package_coverage_verifier_source = %s, want commercial_go_live_evidence_package.machine_evidence_package_coverage_verifier", label, got)
		}
		requireBool(t, gate, "machine_evidence_package_coverage_verifier_must_pass", true)
		if got := requireString(t, gate, "required_status_before_completion"); !strings.Contains(got, "machine_evidence_package_coverage_verifier_passed") {
			t.Fatalf("%s required_status_before_completion missing machine_evidence_package_coverage_verifier_passed: %s", label, got)
		}
		requireStringSliceContains(t, asStringSlice(t, gate["required_before_go_live_signoff"]), "machine_evidence_package_coverage_verifier_passed")
		requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")
	}

	for _, doc := range map[string]map[string]any{
		"commercial_evidence_contract_index": index,
		"commercial_readiness_current_audit": audit,
	} {
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")
	}
}

func TestShortVideoFrontendBackendAcceptanceConsumesMachineEvidenceCoverageVerifier(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	coverageVerifierPath := requireString(t, goLive, "machine_evidence_package_coverage_verifier")

	verifier := readJSON(t, filepath.Join(base, coverageVerifierPath))
	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	runbookPath := requireString(t, goLive, "frontend_backend_handoff_runbook")
	contract := readJSON(t, filepath.Join(base, contractPath))
	runbook := readJSON(t, filepath.Join(base, runbookPath))

	requiredDocs := asStringSlice(t, verifier["required_commercial_gate_docs"])
	requireStringSliceContains(t, requiredDocs, contractPath)
	requireStringSliceContains(t, requiredDocs, runbookPath)

	for label, doc := range map[string]map[string]any{
		"frontend_backend_acceptance_contract": contract,
		"frontend_backend_handoff_runbook":     runbook,
	} {
		if got := requireString(t, doc, "machine_evidence_package_coverage_verifier"); got != coverageVerifierPath {
			t.Fatalf("%s machine_evidence_package_coverage_verifier = %s, want %s", label, got, coverageVerifierPath)
		}
	}

	for label, gateSpec := range map[string]struct {
		doc       map[string]any
		statusKey string
	}{
		"frontend_backend_acceptance_contract": {doc: contract, statusKey: "required_status_before_completion"},
		"frontend_backend_handoff_runbook":     {doc: runbook, statusKey: "required_status_before_acceptance"},
	} {
		gate := requireObject(t, gateSpec.doc, "completion_gate")
		if got := requireString(t, gate, "required_machine_evidence_package_coverage_verifier_source"); got != "commercial_go_live_evidence_package.machine_evidence_package_coverage_verifier" {
			t.Fatalf("%s completion gate coverage verifier source = %s, want commercial_go_live_evidence_package.machine_evidence_package_coverage_verifier", label, got)
		}
		requireBool(t, gate, "machine_evidence_package_coverage_verifier_must_pass", true)
		if got := requireString(t, gate, gateSpec.statusKey); !strings.Contains(got, "machine_evidence_package_coverage_verifier_passed") {
			t.Fatalf("%s %s missing machine_evidence_package_coverage_verifier_passed: %s", label, gateSpec.statusKey, got)
		}
		requireStringSliceContains(t, asStringSlice(t, gate["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")
		requireStringSliceContains(t, asStringSlice(t, gateSpec.doc["non_sufficient_evidence"]), "machine_evidence_package_coverage_verifier_pending")
	}

	checks := asObjectSlice(t, verifier["coverage_checks"])
	for _, docPath := range []string{contractPath, runbookPath} {
		check := findObjectByString(t, checks, "doc_path", docPath)
		requireBool(t, check, "doc_references_machine_evidence_package_coverage_verifier", true)
		requireBool(t, check, "completion_gate_requires_machine_evidence_package_coverage_verifier", true)
		requireBool(t, check, "completion_gate_blocks_pending_machine_evidence_package_coverage_verifier", true)
		requireBool(t, check, "can_count_toward_commercial_ready", false)
	}

	requireStringSliceContains(t, asStringSlice(t, verifier["required_before_coverage_pass"]), "all_frontend_backend_completion_gates_require_machine_evidence_package_coverage_verifier")
}

func TestShortVideoCommercialReadinessUsesLatestP11PreflightReverificationAsP12Baseline(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p11 := requireObject(t, candidateSet, "p11_lifecycle_preflight")
	p12 := requireObject(t, candidateSet, "p12_safe_lifecycle_sample")
	latest := requireObject(t, p11, "latest_reverification")
	sourceLedger := requireString(t, p11, "ledger")

	docs := map[string]map[string]any{
		"commercial_readiness_verifier": readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier"))),
		"commercial_go_no_go_gate":      readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate"))),
		"current_state_audit":           readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit"))),
		"cross_repo_execution_queue":    readJSON(t, filepath.Join(base, requireString(t, goLive, "cross_repo_execution_queue"))),
		"p12_readiness_package":         readJSON(t, filepath.Join(base, requireString(t, p12, "execution_readiness_package"))),
	}

	for label, doc := range docs {
		baseline := requireObject(t, doc, "p11_reverification_baseline")
		if got := requireString(t, baseline, "source_candidate_set_slice"); got != "p11_lifecycle_preflight" {
			t.Fatalf("%s p11_reverification_baseline.source_candidate_set_slice = %s, want p11_lifecycle_preflight", label, got)
		}
		if got := requireString(t, baseline, "source_ledger"); got != sourceLedger {
			t.Fatalf("%s p11_reverification_baseline.source_ledger = %s, want %s", label, got, sourceLedger)
		}
		for _, field := range []string{
			"verified_at",
			"verified_blocked_status",
			"verified_blocked_reason",
			"commercial_status_after_reverification",
		} {
			if got, want := requireString(t, baseline, field), requireString(t, latest, field); got != want {
				t.Fatalf("%s p11_reverification_baseline.%s = %s, want %s", label, field, got, want)
			}
		}
		requireBool(t, baseline, "must_be_current_before_p12_execution", true)
		requireStringSliceContains(t, asStringSlice(t, baseline["required_before_p12_cross_repo_work"]), "p11_lifecycle_preflight_reverification_current")
		requireStringSliceContains(t, asStringSlice(t, baseline["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
	}
}

func TestShortVideoNextAuthorizationStartGuardConsumesLatestP11PreflightReverificationBaseline(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	p11 := requireObject(t, candidateSet, "p11_lifecycle_preflight")
	latest := requireObject(t, p11, "latest_reverification")

	guardPath := requireString(t, goLive, "next_authorization_start_guard")
	guard := readJSON(t, filepath.Join(base, guardPath))
	baseline := requireObject(t, guard, "p11_reverification_baseline")

	if got := requireString(t, baseline, "source_candidate_set_slice"); got != "p11_lifecycle_preflight" {
		t.Fatalf("p11_reverification_baseline.source_candidate_set_slice = %s, want p11_lifecycle_preflight", got)
	}
	if got := requireString(t, baseline, "source_ledger"); got != requireString(t, p11, "ledger") {
		t.Fatalf("p11_reverification_baseline.source_ledger = %s, want %s", got, requireString(t, p11, "ledger"))
	}
	for _, field := range []string{
		"verified_at",
		"verified_blocked_status",
		"verified_blocked_reason",
		"commercial_status_after_reverification",
	} {
		if got, want := requireString(t, baseline, field), requireString(t, latest, field); got != want {
			t.Fatalf("p11_reverification_baseline.%s = %s, want %s", field, got, want)
		}
	}
	requireBool(t, baseline, "must_be_current_before_p12_execution", true)
	requireStringSliceContains(t, asStringSlice(t, baseline["required_before_p12_cross_repo_work"]), "p11_lifecycle_preflight_reverification_current")
	requireStringSliceContains(t, asStringSlice(t, baseline["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
	requireStringSliceContains(t, asStringSlice(t, guard["required_before_guard_open"]), "p11_lifecycle_preflight_reverification_current")
	requireStringSliceContains(t, asStringSlice(t, guard["non_sufficient_evidence"]), "stale_or_missing_p11_lifecycle_preflight_reverification")
}

func TestShortVideoPackStudioGoalCompletionEvidenceMapBlocksGoalCompletionUntilCommercialEvidenceExists(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")

	mapPath := requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")
	requireExistingPath(t, mapPath, base)
	evidenceMap := readJSON(t, filepath.Join(base, mapPath))

	requireBool(t, evidenceMap, "candidate_only", true)
	requireBool(t, evidenceMap, "non_formal", true)
	requireBool(t, evidenceMap, "can_mark_goal_complete", false)
	if got := requireString(t, evidenceMap, "completion_status"); got != "not_achieved_requires_p12_p18_cross_repo_evidence" {
		t.Fatalf("completion_status = %s, want not_achieved_requires_p12_p18_cross_repo_evidence", got)
	}
	if got := requireString(t, evidenceMap, "source_plan"); got != requireString(t, candidateSet, "source_plan") {
		t.Fatalf("source_plan = %s, want %s", got, requireString(t, candidateSet, "source_plan"))
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         requireString(t, candidateSet, "candidate_set_ref"),
		"source_commercial_readiness_verifier":      requireString(t, goLive, "commercial_readiness_verifier"),
		"source_commercial_go_no_go_gate":           requireString(t, goLive, "commercial_go_no_go_gate"),
		"source_post_run_gate_coverage_verifier":    requireString(t, goLive, "post_run_gate_coverage_verifier"),
		"source_commercial_evidence_contract_index": requireString(t, goLive, "evidence_contract_index"),
		"source_frontend_backend_acceptance":        requireString(t, goLive, "frontend_backend_acceptance_contract"),
		"source_cross_repo_execution_queue":         requireString(t, goLive, "cross_repo_execution_queue"),
		"next_required_authorization":               "p12_safe_lifecycle_sample",
		"next_authorization_card":                   requireString(t, goLive, "next_authorization_card"),
	} {
		if got := requireString(t, evidenceMap, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	expectedRequirements := map[string]struct {
		keyword string
		status  string
	}{
		"frontend_capability_pack_studio_complete":  {"前端", "missing_authoritative_evidence"},
		"backend_capability_pack_studio_complete":   {"后端", "missing_authoritative_evidence"},
		"safe_lifecycle_enabled_path":               {"lifecycle", "missing_authoritative_evidence"},
		"three_short_video_candidates_gui":          {"三候选", "missing_authoritative_evidence"},
		"controlled_code_assistant_run":             {"Code Assistant", "missing_authoritative_evidence"},
		"provider_adapter_readiness":                {"Provider", "missing_authoritative_evidence"},
		"cloud_market_sandbox":                      {"云市场", "missing_authoritative_evidence"},
		"gui_api_receipt_traceability":              {"GUI/API/Receipt", "missing_authoritative_evidence"},
		"forbidden_actions_proven_false":            {"禁入动作", "missing_authoritative_evidence"},
		"forbidden_action_coverage_verified":        {"禁入动作覆盖", "missing_authoritative_evidence"},
		"owner_base_gate_receipts_bound":            {"Owner/Base Gate", "missing_authoritative_evidence"},
		"commercial_go_no_go_passed":                {"go/no-go", "missing_authoritative_evidence"},
		"post_run_gate_coverage_verified":           {"后验收门覆盖", "missing_authoritative_evidence"},
		"machine_go_live_evidence_package_verified": {"商用最终机器证据包", "blocked_not_commercial_ready"},
		"independent_acceptance_signoff_matrix":     {"独立验收", "missing_authoritative_evidence"},
	}
	requirements := asObjectSlice(t, evidenceMap["active_goal_requirements"])
	if len(requirements) != len(expectedRequirements) {
		t.Fatalf("active_goal_requirements len = %d, want %d", len(requirements), len(expectedRequirements))
	}
	for _, requirement := range requirements {
		reqID := requireString(t, requirement, "requirement_id")
		expected, ok := expectedRequirements[reqID]
		if !ok {
			t.Fatalf("unexpected requirement_id %s", reqID)
		}
		if got := requireString(t, requirement, "goal_keyword"); !strings.Contains(got, expected.keyword) {
			t.Fatalf("%s goal_keyword = %s, want contains %s", reqID, got, expected.keyword)
		}
		if got := requireString(t, requirement, "current_status"); got != expected.status {
			t.Fatalf("%s current_status = %s, want %s", reqID, got, expected.status)
		}
		requireBool(t, requirement, "can_count_toward_goal_completion", false)
		for _, key := range []string{
			"truth_sources",
			"required_authoritative_evidence",
			"current_candidate_evidence",
			"missing_authoritative_evidence",
			"blocking_reasons",
		} {
			if values := asStringSlice(t, requirement[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", reqID, key)
			}
		}
	}

	barriers := strings.Join(asStringSlice(t, evidenceMap["goal_completion_barriers"]), "\n")
	for _, barrier := range []string{
		"p12_p18_cross_repo_authorization_missing",
		"frontend_backend_commercial_evidence_missing",
		"provider_adapter_readiness_evidence_missing",
		"cloud_market_sandbox_evidence_missing",
		"post_run_gate_coverage_verifier_pending",
		"machine_go_live_evidence_package_blocked",
		"owner_base_gate_receipts_missing",
		"commercial_go_no_go_gate_blocked",
	} {
		if !strings.Contains(barriers, barrier) {
			t.Fatalf("goal_completion_barriers missing %s", barrier)
		}
	}

	policy := requireObject(t, evidenceMap, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	requiredBeforeCompletion := strings.Join(asStringSlice(t, policy["required_before_goal_complete"]), "\n")
	for _, proof := range []string{
		"all_active_goal_requirements_verified",
		"commercial_readiness_verifier_passed",
		"commercial_go_no_go_gate_passed",
		"post_run_gate_coverage_verifier_passed",
		"all_p12_p18_evidence_writebacks_completed",
		"owner_base_gate_receipts_bound",
	} {
		if !strings.Contains(requiredBeforeCompletion, proof) {
			t.Fatalf("completion_claim_policy required_before_goal_complete missing %s", proof)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "p11_preflight_only")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "candidate_assets_only")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "post_run_gate_coverage_verifier_pending")
}

func TestShortVideoGoalCompletionRejectsCoverageOnlyCommercialEvidence(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))
	machinePackage := readJSON(t, filepath.Join(base, requireString(t, goLive, "machine_evidence_package")))
	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))

	requireBool(t, goalMap, "coverage_verifiers_without_runtime_receipts_can_mark_goal_complete", false)
	requireBool(t, goalMap, "machine_go_live_summary_without_slice_receipts_can_mark_goal_complete", false)
	requireBool(t, goalMap, "go_no_go_candidate_without_owner_base_receipts_can_mark_goal_complete", false)

	for _, proof := range []string{
		"all_authoritative_runtime_receipts_bound",
		"frontend_backend_runtime_receipts_bound",
		"p12_p18_slice_receipts_bound",
		"coverage_verifiers_backed_by_terminal_evidence",
	} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]), proof)
		requireStringSliceContains(t, asStringSlice(t, machinePackage["required_before_commercial_ready"]), proof)
	}

	for _, nonSufficient := range []string{
		"coverage_verifiers_without_runtime_receipts",
		"machine_go_live_summary_without_slice_receipts",
		"go_no_go_candidate_without_owner_base_receipts",
		"frontend_backend_contract_without_runtime_receipts",
	} {
		requireStringSliceContains(t, asStringSlice(t, goalMap["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, machinePackage["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), nonSufficient)
		requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), nonSufficient)
	}

	policy := requireObject(t, goalMap, "substitution_rejection_policy")
	for _, key := range []string{
		"coverage_verifiers_without_runtime_receipts",
		"machine_go_live_summary_without_slice_receipts",
		"go_no_go_candidate_without_owner_base_receipts",
		"frontend_backend_contract_without_runtime_receipts",
	} {
		if got := requireString(t, policy, key); got != "not_completion_evidence" {
			t.Fatalf("substitution_rejection_policy.%s = %s, want not_completion_evidence", key, got)
		}
	}
}

func TestShortVideoCommercialGatesRequireGuiApiReceiptTraceability(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	contract := readJSON(t, filepath.Join(base, contractPath))
	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	machinePackage := readJSON(t, filepath.Join(base, requireString(t, goLive, "machine_evidence_package")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	traceRows := asObjectSlice(t, contract["gui_api_receipt_traceability_matrix"])
	if len(traceRows) == 0 {
		t.Fatalf("gui_api_receipt_traceability_matrix must not be empty")
	}

	frontendBackend := requireObject(t, requireObject(t, readiness, "current_blockers"), "frontend_backend_acceptance")
	if got := requireString(t, frontendBackend, "gui_api_receipt_traceability_matrix_source"); got != contractPath+"#gui_api_receipt_traceability_matrix" {
		t.Fatalf("readiness frontend_backend_acceptance.gui_api_receipt_traceability_matrix_source = %s", got)
	}
	traceGroup := requireObject(t, requireObject(t, frontendBackend, "required_evidence_groups"), "gui_api_receipt_traceability_matrix")
	if got := int(requireNumber(t, traceGroup, "total")); got != len(traceRows) {
		t.Fatalf("gui_api_receipt_traceability_matrix total = %d, want %d", got, len(traceRows))
	}
	if got := int(requireNumber(t, traceGroup, "pending")); got != len(traceRows) {
		t.Fatalf("gui_api_receipt_traceability_matrix pending = %d, want %d", got, len(traceRows))
	}
	if got := int(requireNumber(t, traceGroup, "passed")); got != 0 {
		t.Fatalf("gui_api_receipt_traceability_matrix passed = %d, want 0", got)
	}
	requireBool(t, frontendBackend, "gui_api_receipt_traceability_can_count", false)

	for label, values := range map[string][]string{
		"readiness required_before_commercial_ready": asStringSlice(t, readiness["required_before_commercial_ready"]),
		"go/no-go required_before_go_live_signoff":   asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]),
		"machine package required_before_ready":      asStringSlice(t, machinePackage["required_before_commercial_ready"]),
		"goal map required_before_goal_complete":     asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]),
	} {
		if label == "" {
			t.Fatalf("unreachable")
		}
		requireStringSliceContains(t, values, "gui_api_receipt_traceability_verified")
	}
	for label, doc := range map[string]map[string]any{
		"readiness":       readiness,
		"go/no-go":        goNoGo,
		"machine package": machinePackage,
		"goal map":        goalMap,
	} {
		for _, nonSufficient := range []string{
			"frontend_screenshot_without_api_receipt_trace",
			"backend_test_without_gui_api_trace",
			"network_summary_without_receipt_lookup",
		} {
			if label == "" {
				t.Fatalf("unreachable")
			}
			requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), nonSufficient)
		}
	}

	requirementFound := false
	for _, requirement := range asObjectSlice(t, goalMap["active_goal_requirements"]) {
		if requireString(t, requirement, "requirement_id") != "gui_api_receipt_traceability" {
			continue
		}
		requirementFound = true
		if got := requireString(t, requirement, "current_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("gui_api_receipt_traceability.current_status = %s, want missing_authoritative_evidence", got)
		}
		requireBool(t, requirement, "can_count_toward_goal_completion", false)
		requireStringSliceContains(t, asStringSlice(t, requirement["current_candidate_evidence"]), contractPath+"#gui_api_receipt_traceability_matrix")
		requireStringSliceContains(t, asStringSlice(t, requirement["blocking_reasons"]), "gui_api_receipt_traceability_evidence_missing")
	}
	if !requirementFound {
		t.Fatalf("active_goal_requirements missing gui_api_receipt_traceability")
	}

	policy := requireObject(t, goalMap, "substitution_rejection_policy")
	for _, key := range []string{
		"frontend_screenshot_without_api_receipt_trace",
		"backend_test_without_gui_api_trace",
		"network_summary_without_receipt_lookup",
	} {
		if got := requireString(t, policy, key); got != "not_completion_evidence" {
			t.Fatalf("substitution_rejection_policy.%s = %s, want not_completion_evidence", key, got)
		}
	}
}

func TestShortVideoCommercialReadinessVerifierIncludesFrontendBackendAcceptanceBlocker(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	verifier := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	contract := readJSON(t, filepath.Join(base, contractPath))

	blockers := requireObject(t, verifier, "current_blockers")
	frontendBackend := requireObject(t, blockers, "frontend_backend_acceptance")
	if got := requireString(t, frontendBackend, "acceptance_contract"); got != contractPath {
		t.Fatalf("frontend_backend_acceptance.acceptance_contract = %s, want %s", got, contractPath)
	}
	if got := requireString(t, frontendBackend, "acceptance_status"); got != requireString(t, contract, "acceptance_status") {
		t.Fatalf("frontend_backend_acceptance.acceptance_status = %s, want %s", got, requireString(t, contract, "acceptance_status"))
	}
	requireBool(t, frontendBackend, "can_count_toward_commercial_ready", false)
	if got := requireString(t, frontendBackend, "blocked_reason"); got == "" {
		t.Fatalf("frontend_backend_acceptance.blocked_reason missing")
	}
	for _, repo := range []string{"truzhenos", "truzhen-client-web-desktop"} {
		requireStringSliceContains(t, asStringSlice(t, frontendBackend["required_repositories"]), repo)
	}
	for _, sliceKey := range []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
	} {
		requireStringSliceContains(t, asStringSlice(t, frontendBackend["required_source_slices"]), sliceKey)
	}

	evidenceGroups := requireObject(t, frontendBackend, "required_evidence_groups")
	expectedCounts := map[string]int{
		"required_backend_endpoints":       len(asObjectSlice(t, contract["required_backend_endpoints"])),
		"required_frontend_surfaces":       len(asObjectSlice(t, contract["required_frontend_surfaces"])),
		"required_e2e_scenarios":           len(asObjectSlice(t, contract["required_e2e_scenarios"])),
		"required_forbidden_action_checks": len(asObjectSlice(t, contract["required_forbidden_action_checks"])),
	}
	for group, wantCount := range expectedCounts {
		groupSummary := requireObject(t, evidenceGroups, group)
		if got := int(requireNumber(t, groupSummary, "total")); got != wantCount {
			t.Fatalf("%s total = %d, want %d", group, got, wantCount)
		}
		if got := int(requireNumber(t, groupSummary, "passed")); got != 0 {
			t.Fatalf("%s passed = %d, want 0 before authorization", group, got)
		}
		if got := int(requireNumber(t, groupSummary, "pending")); got != wantCount {
			t.Fatalf("%s pending = %d, want %d before authorization", group, got, wantCount)
		}
	}

	requiredBeforeReady := asStringSlice(t, verifier["required_before_commercial_ready"])
	requireStringSliceContains(t, requiredBeforeReady, "frontend_backend_acceptance_contract_passed")
}

func TestShortVideoCommercialReadinessAuditIncludesFrontendBackendAcceptanceBlocker(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	audit := readJSON(t, filepath.Join(base, requireString(t, goLive, "current_state_audit")))
	contractPath := requireString(t, goLive, "frontend_backend_acceptance_contract")
	contract := readJSON(t, filepath.Join(base, contractPath))

	if got := requireString(t, audit, "frontend_backend_acceptance_contract"); got != contractPath {
		t.Fatalf("audit frontend_backend_acceptance_contract = %s, want %s", got, contractPath)
	}
	blocker := requireObject(t, audit, "frontend_backend_acceptance")
	if got := requireString(t, blocker, "acceptance_status"); got != requireString(t, contract, "acceptance_status") {
		t.Fatalf("audit frontend_backend_acceptance.acceptance_status = %s, want %s", got, requireString(t, contract, "acceptance_status"))
	}
	requireBool(t, blocker, "can_count_toward_commercial_ready", false)
	requireStringSliceContains(t, asStringSlice(t, blocker["required_repositories"]), "truzhenos")
	requireStringSliceContains(t, asStringSlice(t, blocker["required_repositories"]), "truzhen-client-web-desktop")
	for _, group := range []string{
		"required_backend_endpoints",
		"required_frontend_surfaces",
		"required_e2e_scenarios",
		"required_forbidden_action_checks",
	} {
		requireStringSliceContains(t, asStringSlice(t, blocker["pending_evidence_groups"]), group)
	}
	requireStringSliceContains(t, asStringSlice(t, audit["non_sufficient_evidence"]), "frontend_backend_contract_without_cross_repo_evidence")
}

func TestShortVideoCommercialReadinessRequiresIndependentAcceptanceSignoffMatrix(t *testing.T) {
	base := filepath.Join("capability-pack-candidates", "short-video-ops-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	goLive := requireObject(t, candidateSet, "commercial_go_live_evidence_package")
	matrixPath := requireString(t, goLive, "independent_acceptance_signoff_matrix")
	if got := requireString(t, candidateSet, "independent_acceptance_signoff_matrix"); got != matrixPath {
		t.Fatalf("candidate_set.independent_acceptance_signoff_matrix = %s, want %s", got, matrixPath)
	}
	requireExistingPath(t, matrixPath, base)

	matrix := readJSON(t, filepath.Join(base, matrixPath))
	readiness := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_readiness_verifier")))
	goNoGo := readJSON(t, filepath.Join(base, requireString(t, goLive, "commercial_go_no_go_gate")))
	machinePackage := readJSON(t, filepath.Join(base, requireString(t, goLive, "machine_evidence_package")))
	goalMap := readJSON(t, filepath.Join(base, requireString(t, candidateSet, "pack_studio_goal_completion_evidence_map")))

	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	requireBool(t, matrix, "can_mark_commercial_ready", false)
	requireBool(t, matrix, "can_request_owner_signoff", false)
	if got := requireString(t, matrix, "signoff_status"); got != "blocked_pending_independent_acceptance_evidence" {
		t.Fatalf("signoff_status = %s, want blocked_pending_independent_acceptance_evidence", got)
	}

	rows := asObjectSlice(t, matrix["review_rows"])
	expectedRows := []string{
		"p12_safe_lifecycle_sample",
		"p13_gui_lifecycle_panel",
		"p15_gui_walkthrough_three_candidates",
		"p16_controlled_code_assistant_run",
		"p17_provider_adapter_candidate",
		"p18_cloud_market_sandbox",
		"gui_api_receipt_traceability",
		"pack_studio_issue_resolution_evidence_schema",
		"forbidden_action_terminal_checks",
		"evidence_writeback_completeness",
		"owner_base_gate_receipts",
		"machine_go_live_evidence_package",
	}
	if len(rows) < len(expectedRows) {
		t.Fatalf("independent acceptance rows = %d, want at least %d", len(rows), len(expectedRows))
	}
	for _, rowID := range expectedRows {
		row := findObjectByString(t, rows, "review_row_id", rowID)
		if got := requireString(t, row, "required_evidence_source"); got == "" {
			t.Fatalf("%s required_evidence_source missing", rowID)
		}
		if got := requireString(t, row, "required_reviewer_role"); got == "" {
			t.Fatalf("%s required_reviewer_role missing", rowID)
		}
		if got := requireString(t, row, "current_status"); got != "pending_evidence" {
			t.Fatalf("%s current_status = %s, want pending_evidence", rowID, got)
		}
		requireBool(t, row, "can_count_toward_commercial_ready", false)
		if refs := asStringSlice(t, row["required_artifact_refs"]); len(refs) == 0 {
			t.Fatalf("%s required_artifact_refs missing", rowID)
		}
		requireStringSliceContains(t, asStringSlice(t, row["blocked_by"]), "independent_acceptance_evidence_missing")
		requireStringSliceContains(t, asStringSlice(t, row["forbidden_substitutions"]), "self_attestation_without_independent_reviewer")
	}

	completionGate := requireObject(t, matrix, "completion_gate")
	requireBool(t, completionGate, "can_pass", false)
	for _, required := range []string{
		"all_independent_review_rows_verified",
		"all_review_artifact_refs_recorded",
		"all_reviewer_identities_recorded",
		"all_blocked_rows_resolved_or_explicitly_not_ready",
	} {
		requireStringSliceContains(t, asStringSlice(t, completionGate["required_before_pass"]), required)
	}
	for _, nonSufficient := range []string{
		"independent_acceptance_signoff_missing",
		"self_attestation_without_independent_reviewer",
		"green_tests_without_independent_review",
		"manual_summary_without_artifact_refs",
	} {
		requireStringSliceContains(t, asStringSlice(t, matrix["non_sufficient_evidence"]), nonSufficient)
	}

	for label, doc := range map[string]map[string]any{
		"commercial readiness verifier": readiness,
		"commercial go/no-go gate":      goNoGo,
		"machine evidence package":      machinePackage,
		"goal completion map":           goalMap,
	} {
		if label == "" {
			t.Fatalf("unreachable")
		}
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "independent_acceptance_signoff_missing")
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "self_attestation_without_independent_reviewer")
		requireStringSliceContains(t, asStringSlice(t, doc["non_sufficient_evidence"]), "green_tests_without_independent_review")
	}

	if got := requireString(t, readiness, "independent_acceptance_signoff_matrix"); got != matrixPath {
		t.Fatalf("readiness independent_acceptance_signoff_matrix = %s, want %s", got, matrixPath)
	}
	if got := requireString(t, goNoGo, "independent_acceptance_signoff_matrix"); got != matrixPath {
		t.Fatalf("go/no-go independent_acceptance_signoff_matrix = %s, want %s", got, matrixPath)
	}
	if got := requireString(t, requireObject(t, machinePackage, "hard_gate_sources"), "independent_acceptance_signoff_matrix"); got != matrixPath {
		t.Fatalf("machine package hard_gate_sources.independent_acceptance_signoff_matrix = %s, want %s", got, matrixPath)
	}
	if got := requireString(t, goalMap, "source_independent_acceptance_signoff_matrix"); got != matrixPath {
		t.Fatalf("goal map source_independent_acceptance_signoff_matrix = %s, want %s", got, matrixPath)
	}

	for label, values := range map[string][]string{
		"readiness required_before_commercial_ready": asStringSlice(t, readiness["required_before_commercial_ready"]),
		"go/no-go required_before_go_live_signoff":   asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_go_live_signoff"]),
		"machine package required_before_ready":      asStringSlice(t, machinePackage["required_before_commercial_ready"]),
		"goal map required_before_goal_complete":     asStringSlice(t, requireObject(t, goalMap, "completion_claim_policy")["required_before_goal_complete"]),
	} {
		if label == "" {
			t.Fatalf("unreachable")
		}
		requireStringSliceContains(t, values, "independent_acceptance_signoff_matrix_passed")
	}

	requirement := findObjectByString(t, asObjectSlice(t, goalMap["active_goal_requirements"]), "requirement_id", "independent_acceptance_signoff_matrix")
	if got := requireString(t, requirement, "current_status"); got != "missing_authoritative_evidence" {
		t.Fatalf("independent_acceptance_signoff_matrix.current_status = %s, want missing_authoritative_evidence", got)
	}
	requireBool(t, requirement, "can_count_toward_goal_completion", false)
	requireStringSliceContains(t, asStringSlice(t, requirement["current_candidate_evidence"]), matrixPath)
	requireStringSliceContains(t, asStringSlice(t, requirement["blocking_reasons"]), "independent_acceptance_signoff_missing")
}

func requireCapabilityPackHasProviderRequirement(t *testing.T, doc map[string]any) {
	t.Helper()
	requirements, ok := doc["provider_requirements"].([]any)
	if !ok || len(requirements) == 0 {
		t.Fatalf("provider_requirements missing")
	}
	for _, raw := range requirements {
		requirement, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("provider requirement = %T", raw)
		}
		if got := requireString(t, requirement, "fallback_policy"); got == "ready" {
			t.Fatalf("fallback_policy must not claim ready by default")
		}
	}
}

func asObjectSlice(t *testing.T, value any) []map[string]any {
	t.Helper()
	rawItems, ok := value.([]any)
	if !ok {
		t.Fatalf("expected object array, got %T", value)
	}
	out := make([]map[string]any, 0, len(rawItems))
	for _, raw := range rawItems {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("expected object array item, got %T", raw)
		}
		out = append(out, item)
	}
	return out
}

func findObjectByString(t *testing.T, items []map[string]any, key string, want string) map[string]any {
	t.Helper()
	for _, item := range items {
		value, ok := item[key].(string)
		if ok && value == want {
			return item
		}
	}
	t.Fatalf("expected object with %s=%s", key, want)
	return nil
}

func requireContractHasEvidenceGroup(t *testing.T, doc map[string]any) {
	t.Helper()
	for _, key := range []string{
		"required_backend_evidence",
		"required_frontend_evidence",
		"required_execution_evidence",
		"required_provider_evidence",
		"required_cloud_evidence",
		"required_forbidden_action_checks",
	} {
		rawItems, ok := doc[key].([]any)
		if !ok {
			continue
		}
		if len(rawItems) == 0 {
			t.Fatalf("%s must not be empty", key)
		}
		for _, raw := range rawItems {
			item, ok := raw.(map[string]any)
			if !ok {
				t.Fatalf("%s item = %T", key, raw)
			}
			if got := requireString(t, item, "evidence_id"); got == "" {
				t.Fatalf("%s evidence_id missing", key)
			}
			if got := requireString(t, item, "current_status"); !strings.Contains(got, "pending") {
				t.Fatalf("%s current_status = %s, want pending before execution", key, got)
			}
			if got := requireString(t, item, "evidence_required"); got == "" {
				t.Fatalf("%s evidence_required missing", key)
			}
		}
		return
	}
	t.Fatalf("evidence contract has no required evidence group")
}

func collectEvidenceIDsFromContract(t *testing.T, contract map[string]any) []string {
	t.Helper()
	var ids []string
	for key, raw := range contract {
		if !strings.HasPrefix(key, "required_") {
			continue
		}
		items, ok := raw.([]any)
		if !ok {
			continue
		}
		for _, rawItem := range items {
			item, ok := rawItem.(map[string]any)
			if !ok {
				continue
			}
			if _, ok := item["evidence_id"]; !ok {
				continue
			}
			ids = append(ids, requireString(t, item, "evidence_id"))
		}
	}
	if len(ids) == 0 {
		t.Fatalf("evidence contract has no evidence_id entries")
	}
	return ids
}

func collectEvidenceIDsFromItems(t *testing.T, value any) []string {
	t.Helper()
	items := asObjectSlice(t, value)
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, requireString(t, item, "evidence_id"))
	}
	if len(ids) == 0 {
		t.Fatalf("evidence items missing")
	}
	return ids
}

func requireEvidenceItems(t *testing.T, doc map[string]any, key string, expectedIDs []string) {
	t.Helper()
	rawItems, ok := doc[key].([]any)
	if !ok {
		t.Fatalf("%s missing", key)
	}
	seen := map[string]bool{}
	for _, raw := range rawItems {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("%s item = %T", key, raw)
		}
		id := requireString(t, item, "evidence_id")
		seen[id] = true
		if got := requireString(t, item, "current_status"); !strings.Contains(got, "pending") {
			t.Fatalf("%s.%s current_status = %s, want pending before execution", key, id, got)
		}
		if got := requireString(t, item, "evidence_required"); got == "" {
			t.Fatalf("%s.%s evidence_required missing", key, id)
		}
	}
	for _, expectedID := range expectedIDs {
		if !seen[expectedID] {
			t.Fatalf("%s missing evidence_id %s", key, expectedID)
		}
	}
}

func requireAcceptanceEvidenceItems(t *testing.T, doc map[string]any, key string, expectedIDs []string) {
	t.Helper()
	rawItems, ok := doc[key].([]any)
	if !ok {
		t.Fatalf("%s missing", key)
	}
	seen := map[string]bool{}
	for _, raw := range rawItems {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("%s item = %T", key, raw)
		}
		id := requireString(t, item, "evidence_id")
		seen[id] = true
		if got := requireString(t, item, "current_status"); !strings.Contains(got, "pending") {
			t.Fatalf("%s.%s current_status = %s, want pending before execution", key, id, got)
		}
		if got := requireString(t, item, "evidence_required"); got == "" {
			t.Fatalf("%s.%s evidence_required missing", key, id)
		}
		if got := requireString(t, item, "target_repository"); got == "" {
			t.Fatalf("%s.%s target_repository missing", key, id)
		}
		if got := requireString(t, item, "source_slice"); got == "" {
			t.Fatalf("%s.%s source_slice missing", key, id)
		}
	}
	for _, expectedID := range expectedIDs {
		if !seen[expectedID] {
			t.Fatalf("%s missing evidence_id %s", key, expectedID)
		}
	}
}

func requireCommercialSignoffMatrix(t *testing.T, candidateSet map[string]any, matrix map[string]any, requiredSlices []string, ownerSignoffGate map[string]any) {
	t.Helper()
	if got := requireString(t, matrix, "matrix_status"); got != "blocked_pending_authorization_and_evidence" {
		t.Fatalf("commercial_signoff_matrix.matrix_status = %s, want blocked_pending_authorization_and_evidence", got)
	}
	if got := requireString(t, matrix, "owner_authorization_truth"); got == "truzhen-packs" {
		t.Fatalf("commercial_signoff_matrix.owner_authorization_truth must not point to truzhen-packs")
	}
	if got := requireString(t, matrix, "required_before_owner_signoff"); got == "" {
		t.Fatalf("commercial_signoff_matrix.required_before_owner_signoff missing")
	}
	if got := requireString(t, matrix, "next_authorization"); got != requireString(t, ownerSignoffGate, "required_next_authorization") {
		t.Fatalf("commercial_signoff_matrix.next_authorization = %s, want %s", got, requireString(t, ownerSignoffGate, "required_next_authorization"))
	}

	entries := asObjectSlice(t, matrix["entries"])
	if len(entries) != len(requiredSlices) {
		t.Fatalf("commercial_signoff_matrix.entries len = %d, want %d", len(entries), len(requiredSlices))
	}
	for i, sliceKey := range requiredSlices {
		entry := entries[i]
		if got := requireString(t, entry, "slice_key"); got != sliceKey {
			t.Fatalf("commercial_signoff_matrix.entries[%d].slice_key = %s, want %s", i, got, sliceKey)
		}
		slice := requireObject(t, candidateSet, sliceKey)
		if got, want := requireString(t, entry, "authorization_card"), requireString(t, slice, "cross_repo_authorization_card"); got != want {
			t.Fatalf("%s authorization_card = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, entry, "authorization_scope_contract"), requireString(t, slice, "authorization_scope_contract"); got != want {
			t.Fatalf("%s authorization_scope_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, entry, "authorization_evidence_intake_contract"), requireString(t, slice, "authorization_evidence_intake_contract"); got != want {
			t.Fatalf("%s authorization_evidence_intake_contract = %s, want %s", sliceKey, got, want)
		}
		intake := readJSON(t, filepath.Join("capability-pack-candidates", "short-video-ops-v0", requireString(t, entry, "authorization_evidence_intake_contract")))
		intakeEvidence := requireObject(t, intake, "current_authorization_evidence")
		if got := requireString(t, intakeEvidence, "status"); got != "missing" {
			t.Fatalf("%s intake current_authorization_evidence.status = %s, want missing before Owner authorization", sliceKey, got)
		}
		intakeGate := requireObject(t, intake, "cross_repo_work_gate")
		requireBool(t, intakeGate, "can_start_cross_repo_work", false)
		scope := readJSON(t, filepath.Join("capability-pack-candidates", "short-video-ops-v0", requireString(t, entry, "authorization_scope_contract")))
		scopeEvidence := requireObject(t, scope, "current_authorization_evidence")
		if got := requireString(t, scopeEvidence, "status"); got != "missing" {
			t.Fatalf("%s scope current_authorization_evidence.status = %s, want missing before Owner authorization", sliceKey, got)
		}
		if got, want := requireString(t, entry, "evidence_contract"), requireString(t, slice, "evidence_contract"); got != want {
			t.Fatalf("%s evidence_contract = %s, want %s", sliceKey, got, want)
		}
		if got, want := requireString(t, entry, "evidence_ledger"), requireString(t, slice, "evidence_ledger"); got != want {
			t.Fatalf("%s evidence_ledger = %s, want %s", sliceKey, got, want)
		}
		if got := requireString(t, entry, "owner_authorization_status"); got != "pending_authorization" {
			t.Fatalf("%s owner_authorization_status = %s, want pending_authorization", sliceKey, got)
		}
		if got := requireString(t, entry, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", sliceKey, got)
		}
		requireBool(t, entry, "can_count_toward_commercial_ready", false)
		requireStringSliceContains(t, asStringSlice(t, ownerSignoffGate["blocking_required_slices"]), sliceKey)
	}
}

func requireBoolValue(t *testing.T, doc map[string]any, key string) bool {
	t.Helper()
	value, ok := doc[key].(bool)
	if !ok {
		t.Fatalf("expected bool %q", key)
	}
	return value
}

func requireNumber(t *testing.T, doc map[string]any, key string) float64 {
	t.Helper()
	value, ok := doc[key].(float64)
	if !ok {
		t.Fatalf("expected number %q", key)
	}
	return value
}

func requireStringSliceContains(t *testing.T, values []string, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("expected %q in %v", want, values)
}

func requireStringSlicesEqual(t *testing.T, got []string, want []string, label string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s[%d] = %s, want %s; full got %v want %v", label, i, got[i], want[i], got, want)
		}
	}
}

func stringSliceOrEmpty(t *testing.T, doc map[string]any, key string) []string {
	t.Helper()
	raw, ok := doc[key]
	if !ok {
		return nil
	}
	return asStringSlice(t, raw)
}

// processWorktreeRefRegistry 登记已知外部过程 worktree 引用前缀 → 历史本机路径（仅作出处记录）。
// 防什么：P12-P18 执行规格/授权卡产自外部 superpowers 过程 worktree，该目录不随本仓分发、
// 也不保证在任意机器存在（2026-07-10 分账登记的预存断链，干净 main 同样 FAIL）。
// 为什么在这里防：Owner 2026-07-11 交付债务集中处理轮裁方案 b——对 process_worktree_ref://
// 引用只断言「引用登记存在性」（前缀在本登记表内），不 os.Stat 外部路径；仓内路径仍必须真实存在。
var processWorktreeRefRegistry = map[string]string{
	"process_worktree_ref://truzhen-packs/gui-capability-pack-test-plan": "/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan",
}

func requireExistingPath(t *testing.T, path string, base string) {
	t.Helper()
	if strings.HasPrefix(path, "process_worktree_ref://") {
		for prefix := range processWorktreeRefRegistry {
			if strings.HasPrefix(path, prefix) {
				return
			}
		}
		t.Fatalf("process worktree reference %s is not registered in processWorktreeRefRegistry", path)
	}
	if base != "" && !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	if _, err := os.Stat(resolveReferencePath(path)); err != nil {
		t.Fatalf("expected path %s to exist: %v", path, err)
	}
}

func resolveReferencePath(path string) string {
	for prefix, actual := range processWorktreeRefRegistry {
		if strings.HasPrefix(path, prefix) {
			return actual + strings.TrimPrefix(path, prefix)
		}
	}
	return path
}
