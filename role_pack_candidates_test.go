package packs_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return value
}

func TestTeamOfficeCommercialReadinessAndGoNoGoConsumeIssueLedgerBeforeCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	issueLedgerPath := "docs/role-studio-issue-ledger.md"
	issueLedgerBytes, err := os.ReadFile(filepath.Join(base, issueLedgerPath))
	if err != nil {
		t.Fatalf("read issue ledger: %v", err)
	}
	issueLedger := string(issueLedgerBytes)
	for _, issueID := range []string{
		"ROLE-STUDIO-GUI-001",
		"ROLE-STUDIO-CLOUD-001",
		"ROLE-STUDIO-GO-NO-GO-001",
		"ROLE-STUDIO-P11-EVIDENCE-PACKAGE-001",
	} {
		if !strings.Contains(issueLedger, issueID) {
			t.Fatalf("issue ledger missing %s", issueID)
		}
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_role_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("readiness source_role_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	readinessIssueSummary := requireObject(t, requireObject(t, readiness, "current_blockers"), "issue_ledger_writeback_summary")
	if got := requireString(t, readinessIssueSummary, "source_role_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("readiness issue ledger summary source = %s, want %s", got, issueLedgerPath)
	}
	requireBool(t, readinessIssueSummary, "can_count_toward_commercial_ready", false)
	readinessTerminal := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "all_issue_ledger_entries_triaged")
	if got := requireString(t, readinessTerminal, "current_result"); got != "pending" {
		t.Fatalf("readiness issue ledger current_result = %s, want pending", got)
	}
	requireBool(t, readinessTerminal, "can_count_toward_commercial_ready", false)
	if evidence := requireString(t, readinessTerminal, "evidence_required"); !strings.Contains(evidence, issueLedgerPath) {
		t.Fatalf("readiness issue ledger evidence_required missing %s", issueLedgerPath)
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "issue_ledger_all_entries_triaged")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "issue_ledger_missing_or_untriaged")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_role_studio_issue_ledger"); got != issueLedgerPath {
		t.Fatalf("go/no-go source_role_studio_issue_ledger = %s, want %s", got, issueLedgerPath)
	}
	goNoGoTerminal := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "all_issue_ledger_entries_triaged")
	if got := requireString(t, goNoGoTerminal, "current_result"); got != "pending" {
		t.Fatalf("go/no-go issue ledger current_result = %s, want pending", got)
	}
	requireBool(t, goNoGoTerminal, "required_final_value", true)
	requireBool(t, goNoGoTerminal, "can_pass_gate", false)
	if evidence := requireString(t, goNoGoTerminal, "evidence_required"); !strings.Contains(evidence, issueLedgerPath) {
		t.Fatalf("go/no-go issue ledger evidence_required missing %s", issueLedgerPath)
	}
	completionRule := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_final_decision"]), "issue_ledger_all_entries_triaged")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "issue_ledger_missing_or_untriaged")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["required_before_completion_claim"]), "issue_ledger_all_entries_triaged")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	goalPolicy := requireObject(t, evidenceMap, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, goalPolicy["required_before_goal_complete"]), "issue_ledger_all_entries_triaged")
}

func TestTeamOfficeCommercialReadinessAndGoNoGoConsumeProductStageFrontendBackendClosureReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	reportPath := "tests/product-stage-frontend-backend-closure-report-candidate.json"
	proof := "product_stage_frontend_backend_closure_report_verified"

	report := readJSON(t, filepath.Join(base, reportPath))
	if got := requireString(t, report, "closure_status"); got != "not_achieved_requires_cross_repo_execution" {
		t.Fatalf("closure_status = %s, want not_achieved_requires_cross_repo_execution", got)
	}
	for _, participant := range []string{"user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent"} {
		requireStringSliceContains(t, asStringSlice(t, report["required_participants"]), participant)
	}
	reportPolicy := requireObject(t, report, "completion_claim_policy")
	for _, required := range []string{"frontend_product_stage_smoke", "backend_receipt_lookup", "frontend_backend_field_match_report"} {
		requireStringSliceContains(t, asStringSlice(t, reportPolicy["required_before_completion_claim"]), required)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_product_stage_frontend_backend_closure_report"); got != reportPath {
		t.Fatalf("readiness source_product_stage_frontend_backend_closure_report = %s, want %s", got, reportPath)
	}
	frontendBackend := requireObject(t, requireObject(t, readiness, "current_blockers"), "frontend_backend_acceptance")
	if got := requireString(t, frontendBackend, "closure_report_ref"); got != reportPath {
		t.Fatalf("readiness frontend_backend_acceptance.closure_report_ref = %s, want %s", got, reportPath)
	}
	readinessTerminal := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", proof)
	if got := requireString(t, readinessTerminal, "current_result"); got != "pending" {
		t.Fatalf("readiness product-stage current_result = %s, want pending", got)
	}
	requireBool(t, readinessTerminal, "can_count_toward_commercial_ready", false)
	if evidence := requireString(t, readinessTerminal, "evidence_required"); !strings.Contains(evidence, reportPath) {
		t.Fatalf("readiness product-stage evidence_required missing %s", reportPath)
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), proof)
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_product_stage_frontend_backend_closure_report")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_product_stage_frontend_backend_closure_report"); got != reportPath {
		t.Fatalf("go/no-go source_product_stage_frontend_backend_closure_report = %s, want %s", got, reportPath)
	}
	goNoGoTerminal := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", proof)
	if got := requireString(t, goNoGoTerminal, "current_result"); got != "pending" {
		t.Fatalf("go/no-go product-stage current_result = %s, want pending", got)
	}
	requireBool(t, goNoGoTerminal, "required_final_value", true)
	requireBool(t, goNoGoTerminal, "can_pass_gate", false)
	if evidence := requireString(t, goNoGoTerminal, "evidence_required"); !strings.Contains(evidence, reportPath) {
		t.Fatalf("go/no-go product-stage evidence_required missing %s", reportPath)
	}
	completionRule := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_final_decision"]), proof)
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_product_stage_frontend_backend_closure_report")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["required_before_completion_claim"]), proof)

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	goalPolicy := requireObject(t, evidenceMap, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, goalPolicy["required_before_goal_complete"]), proof)
}

func TestTeamOfficeCommercialReadinessAndGoNoGoConsumeRuntimeUsageSmokeBeforeCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	runtimeUsagePath := "usage/team-office-runtime-usage-candidate.json"
	proof := "team_office_runtime_usage_smoke_verified"

	usage := readJSON(t, filepath.Join(base, runtimeUsagePath))
	if got := requireString(t, usage, "runtime_truth_source"); got != "truzhenos" {
		t.Fatalf("runtime_truth_source = %s, want truzhenos", got)
	}
	requireBool(t, usage, "candidate_only", true)
	scenarios := asObjectSlice(t, usage["usage_scenarios"])
	for _, scenarioID := range []string{
		"secretary_orchestrates_five_advisors",
		"advisors_return_candidate_outputs",
		"capability_invocation_uses_role_context",
		"owner_gate_before_formalization",
	} {
		scenario := findObjectByString(t, scenarios, "scenario_id", scenarioID)
		if got := requireString(t, scenario, "output_status"); got != "candidate_only" {
			t.Fatalf("%s output_status = %s, want candidate_only", scenarioID, got)
		}
		requireStringSliceContains(t, asStringSlice(t, scenario["required_evidence"]), "runtime_receipt_candidate_ref")
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_runtime_usage_candidate"); got != runtimeUsagePath {
		t.Fatalf("readiness source_runtime_usage_candidate = %s, want %s", got, runtimeUsagePath)
	}
	frontendBackend := requireObject(t, requireObject(t, readiness, "current_blockers"), "frontend_backend_acceptance")
	requireStringSliceContains(t, asStringSlice(t, frontendBackend["required_before_pass"]), "runtime_usage_smoke_receipts_present")
	readinessTerminal := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", proof)
	if got := requireString(t, readinessTerminal, "current_result"); got != "pending" {
		t.Fatalf("readiness runtime usage current_result = %s, want pending", got)
	}
	requireBool(t, readinessTerminal, "can_count_toward_commercial_ready", false)
	if evidence := requireString(t, readinessTerminal, "evidence_required"); !strings.Contains(evidence, runtimeUsagePath) {
		t.Fatalf("readiness runtime usage evidence_required missing %s", runtimeUsagePath)
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), proof)
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_runtime_usage_smoke")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_runtime_usage_candidate"); got != runtimeUsagePath {
		t.Fatalf("go/no-go source_runtime_usage_candidate = %s, want %s", got, runtimeUsagePath)
	}
	goNoGoTerminal := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", proof)
	if got := requireString(t, goNoGoTerminal, "current_result"); got != "pending" {
		t.Fatalf("go/no-go runtime usage current_result = %s, want pending", got)
	}
	requireBool(t, goNoGoTerminal, "required_final_value", true)
	requireBool(t, goNoGoTerminal, "can_pass_gate", false)
	if evidence := requireString(t, goNoGoTerminal, "evidence_required"); !strings.Contains(evidence, runtimeUsagePath) {
		t.Fatalf("go/no-go runtime usage evidence_required missing %s", runtimeUsagePath)
	}
	completionRule := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, completionRule["required_before_final_decision"]), proof)
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_runtime_usage_smoke")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["required_before_completion_claim"]), proof)

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	goalPolicy := requireObject(t, evidenceMap, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, goalPolicy["required_before_goal_complete"]), proof)
}

func asStringSlice(t *testing.T, value any) []string {
	t.Helper()
	raw, ok := value.([]any)
	if !ok {
		t.Fatalf("expected array, got %T", value)
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		text, ok := item.(string)
		if !ok {
			t.Fatalf("expected string array item, got %T", item)
		}
		out = append(out, text)
	}
	return out
}

func requireString(t *testing.T, doc map[string]any, key string) string {
	t.Helper()
	value, ok := doc[key].(string)
	if !ok || value == "" {
		t.Fatalf("expected non-empty string %q", key)
	}
	return value
}

func requireBool(t *testing.T, doc map[string]any, key string, want bool) {
	t.Helper()
	value, ok := doc[key].(bool)
	if !ok {
		t.Fatalf("expected bool %q", key)
	}
	if value != want {
		t.Fatalf("%s = %v, want %v", key, value, want)
	}
}

func requireInt(t *testing.T, doc map[string]any, key string, want int) {
	t.Helper()
	value, ok := doc[key].(float64)
	if !ok {
		t.Fatalf("expected number %q", key)
	}
	if int(value) != want || value != float64(want) {
		t.Fatalf("%s = %v, want %d", key, value, want)
	}
}

func requireObject(t *testing.T, doc map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := doc[key].(map[string]any)
	if !ok {
		t.Fatalf("expected object %q", key)
	}
	return value
}

func requireStringIn(t *testing.T, got string, allowed ...string) {
	t.Helper()
	for _, want := range allowed {
		if got == want {
			return
		}
	}
	t.Fatalf("%q not in allowed values %v", got, allowed)
}

func compareFileSets(t *testing.T, leftName string, left map[string]bool, rightName string, right map[string]bool) {
	t.Helper()
	for rel := range left {
		if !right[rel] {
			t.Fatalf("%s contains %s but %s does not", leftName, rel, rightName)
		}
	}
	for rel := range right {
		if !left[rel] {
			t.Fatalf("%s contains %s but %s does not", rightName, rel, leftName)
		}
	}
}

func TestTeamOfficeRolePackCandidateSet(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))

	if got := requireString(t, candidateSet, "candidate_set_ref"); got != "role-pack-candidate-set://team-office-v0" {
		t.Fatalf("candidate_set_ref = %s", got)
	}
	requireBool(t, candidateSet, "candidate_only", true)
	requireBool(t, candidateSet, "non_formal", true)
	if got := requireString(t, candidateSet, "lifecycle_status"); got != "设计中" {
		t.Fatalf("lifecycle_status = %s", got)
	}

	roleRefs := asStringSlice(t, candidateSet["role_pack_refs"])
	if len(roleRefs) != 6 {
		t.Fatalf("role_pack_refs len = %d, want 6", len(roleRefs))
	}
	expectedRefs := map[string]bool{
		"role_pack://team-office-secretary-general":  false,
		"role_pack://team-office-strategy-advisor":   false,
		"role_pack://team-office-product-advisor":    false,
		"role_pack://team-office-operations-advisor": false,
		"role_pack://team-office-finance-advisor":    false,
		"role_pack://team-office-legal-risk-advisor": false,
	}
	for _, ref := range roleRefs {
		if _, ok := expectedRefs[ref]; !ok {
			t.Fatalf("unexpected role ref %s", ref)
		}
		expectedRefs[ref] = true
	}
	for ref, seen := range expectedRefs {
		if !seen {
			t.Fatalf("missing role ref %s", ref)
		}
	}

	truthSources, ok := candidateSet["truth_sources"].(map[string]any)
	if !ok {
		t.Fatalf("missing truth_sources")
	}
	if truthSources["cloud_commerce_truth"] != "truzhen-cloud" {
		t.Fatalf("cloud_commerce_truth = %v, want truzhen-cloud", truthSources["cloud_commerce_truth"])
	}
}

func TestTeamOfficeRolePacksAreCandidateOnlyProposers(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0", "role-packs")
	files := []string{
		"team-office-secretary-general.rolepack.json",
		"team-office-strategy-advisor.rolepack.json",
		"team-office-product-advisor.rolepack.json",
		"team-office-operations-advisor.rolepack.json",
		"team-office-finance-advisor.rolepack.json",
		"team-office-legal-risk-advisor.rolepack.json",
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			doc := readJSON(t, filepath.Join(base, file))
			requireBool(t, doc, "candidate_only", true)
			requireBool(t, doc, "non_formal", true)
			if got := requireString(t, doc, "formal_authority"); got != "none" {
				t.Fatalf("formal_authority = %s, want none", got)
			}
			if got := requireString(t, doc, "role_mode"); got != "proposer" {
				t.Fatalf("role_mode = %s, want proposer", got)
			}
			if len(asStringSlice(t, doc["compatible_slots"])) == 0 {
				t.Fatalf("compatible_slots must not be empty")
			}
			communication, ok := doc["communication_style"].(map[string]any)
			if !ok {
				t.Fatalf("missing communication_style")
			}
			forbidden := strings.Join(asStringSlice(t, communication["forbidden_phrases"]), "\n")
			for _, phrase := range []string{"我已批准", "我已执行", "已正式生效"} {
				if !strings.Contains(forbidden, phrase) {
					t.Fatalf("forbidden_phrases missing %q", phrase)
				}
			}
		})
	}
}

func TestTeamOfficeRoleSlotsAndCapabilityRequirements(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	slots := readJSON(t, filepath.Join(base, "role-slots", "team-office-role-slots.json"))
	rawSlots, ok := slots["role_slots"].([]any)
	if !ok {
		t.Fatalf("role_slots missing")
	}
	if len(rawSlots) != 6 {
		t.Fatalf("role_slots len = %d, want 6", len(rawSlots))
	}
	seenSlots := map[string]bool{}
	for _, raw := range rawSlots {
		slot, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("slot item = %T", raw)
		}
		slotID := requireString(t, slot, "slot_id")
		seenSlots[slotID] = true
		if requireString(t, slot, "default_role_pack_ref") == "" {
			t.Fatalf("slot %s missing default_role_pack_ref", slotID)
		}
		if got := requireString(t, slot, "formal_authority"); got != "none" {
			t.Fatalf("slot %s formal_authority = %s, want none", slotID, got)
		}
	}
	for _, required := range []string{
		"team_office.secretary_general",
		"team_office.advisor.strategy",
		"team_office.advisor.product",
		"team_office.advisor.operations",
		"team_office.advisor.finance",
		"team_office.advisor.legal_risk",
	} {
		if !seenSlots[required] {
			t.Fatalf("missing slot %s", required)
		}
	}

	requirements := readJSON(t, filepath.Join(base, "capability-role-requirements", "sample-team-research.role-requirements.json"))
	requireBool(t, requirements, "candidate_only", true)
	rawReqs, ok := requirements["role_requirements"].([]any)
	if !ok || len(rawReqs) == 0 {
		t.Fatalf("role_requirements missing")
	}
	for _, raw := range rawReqs {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("role requirement = %T", raw)
		}
		if got := requireString(t, req, "formal_authority"); got != "none" {
			t.Fatalf("role requirement formal_authority = %s, want none", got)
		}
	}
}

func TestTeamOfficeCapabilityRequirementsCoverAllSixProducedRoles(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	requirements := readJSON(t, filepath.Join(base, "capability-role-requirements", "sample-team-research.role-requirements.json"))
	if got := requireString(t, requirements, "capability_pack_ref"); got != "capability-pack://sample-team-research" {
		t.Fatalf("capability_pack_ref = %s", got)
	}

	requiredRoles := map[string]bool{}
	for _, ref := range asStringSlice(t, candidateSet["role_pack_refs"]) {
		requiredRoles[ref] = false
	}
	rawReqs, ok := requirements["role_requirements"].([]any)
	if !ok {
		t.Fatalf("role_requirements missing")
	}
	seenRequirementRefs := map[string]bool{}
	for _, raw := range rawReqs {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("role requirement = %T", raw)
		}
		reqRef := requireString(t, req, "requirement_ref")
		if seenRequirementRefs[reqRef] {
			t.Fatalf("duplicate requirement_ref %s", reqRef)
		}
		seenRequirementRefs[reqRef] = true
		if got := requireString(t, req, "formal_authority"); got != "none" {
			t.Fatalf("%s formal_authority = %s, want none", reqRef, got)
		}
		requireBool(t, req, "post_install_reference_required", true)
		if got := requireString(t, req, "reference_scope"); got != "installed_buyer_team_only" {
			t.Fatalf("%s reference_scope = %s, want installed_buyer_team_only", reqRef, got)
		}
		if len(asStringSlice(t, req["accepted_slot_refs"])) == 0 {
			t.Fatalf("%s accepted_slot_refs missing", reqRef)
		}
		for _, roleRef := range asStringSlice(t, req["accepted_role_pack_refs"]) {
			if _, ok := requiredRoles[roleRef]; !ok {
				t.Fatalf("%s references unexpected role %s", reqRef, roleRef)
			}
			requiredRoles[roleRef] = true
		}
	}
	for roleRef, seen := range requiredRoles {
		if !seen {
			t.Fatalf("produced role %s is not covered by capability role requirements", roleRef)
		}
	}
}

func TestSecretaryAppearanceAndCloudListingAreRefsOnly(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	appearance := readJSON(t, filepath.Join(base, "appearance", "secretary-appearance-preferences.json"))
	requireBool(t, appearance, "candidate_only", true)
	requireBool(t, appearance, "asset_ref_only", true)
	for _, key := range []string{"voice_asset_ref", "vrm_asset_ref", "voice_provider_ref", "vrm_provider_ref"} {
		value := requireString(t, appearance, key)
		if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "file:") {
			t.Fatalf("%s must be asset/provider ref, got %s", key, value)
		}
	}

	listing := readJSON(t, filepath.Join(base, "commerce", "cloud-listing-candidate.json"))
	requireBool(t, listing, "candidate_only", true)
	if got := requireString(t, listing, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, listing, "payment_mode"); got != "sandbox_only_until_owner_authorizes_real_payment" {
		t.Fatalf("payment_mode = %s", got)
	}
	for _, key := range []string{"listing_status", "license_status", "entitlement_status"} {
		if got := requireString(t, listing, key); got != "candidate_not_formal" {
			t.Fatalf("%s = %s, want candidate_not_formal", key, got)
		}
	}
	if got := requireString(t, listing, "download_status"); got != "candidate_requires_entitlement" {
		t.Fatalf("download_status = %s, want candidate_requires_entitlement", got)
	}
	if got := requireString(t, listing, "install_status"); got != "candidate_requires_truzhenos_install" {
		t.Fatalf("install_status = %s, want candidate_requires_truzhenos_install", got)
	}
}

func TestTeamOfficeSecretaryAppearancePreferenceDefinesGuiSelectableClearableDefaults(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	appearance := readJSON(t, filepath.Join(base, "appearance", "secretary-appearance-preferences.json"))
	requireBool(t, appearance, "candidate_only", true)
	requireBool(t, appearance, "asset_ref_only", true)

	policy := requireObject(t, appearance, "gui_control_policy")
	if got := requireString(t, policy, "selection_surface"); got != "team_settings_secretary_appearance_tab" {
		t.Fatalf("selection_surface = %s", got)
	}
	for _, key := range []string{
		"voice_selector_required",
		"vrm_selector_required",
		"clear_to_none_supported",
		"reset_to_default_supported",
		"provider_readiness_visible",
		"default_restore_receipt_required",
	} {
		requireBool(t, policy, key, true)
	}
	if got := requireString(t, policy, "blocked_status_when_raw_asset"); got != "blocked_raw_asset_reference" {
		t.Fatalf("blocked_status_when_raw_asset = %s", got)
	}

	assertOptions := func(kind string, raw any, assetPrefix string) {
		t.Helper()
		options := asObjectSlice(t, raw)
		if len(options) == 0 {
			t.Fatalf("%s options missing", kind)
		}
		hasDefault := false
		for _, option := range options {
			assetRef := requireString(t, option, "asset_ref")
			if !strings.HasPrefix(assetRef, assetPrefix) {
				t.Fatalf("%s asset_ref = %s, want prefix %s", kind, assetRef, assetPrefix)
			}
			if strings.HasPrefix(assetRef, "/") || strings.HasPrefix(assetRef, "file:") {
				t.Fatalf("%s asset_ref must not be raw path: %s", kind, assetRef)
			}
			if providerRef := requireString(t, option, "provider_ref"); !strings.HasPrefix(providerRef, "provider://") {
				t.Fatalf("%s provider_ref = %s, want provider:// prefix", kind, providerRef)
			}
			requireString(t, option, "display_name")
			requireString(t, option, "readiness_status")
			requireBool(t, option, "licensed_for_role_pack_distribution", true)
			requireBool(t, option, "raw_asset_included", false)
			if option["default"] == true {
				hasDefault = true
			}
		}
		if !hasDefault {
			t.Fatalf("%s options must include one default option", kind)
		}
	}
	assertOptions("voice", appearance["available_voice_options"], "voice_asset://")
	assertOptions("vrm", appearance["available_vrm_options"], "vrm_asset://")

	actionIDs := map[string]bool{
		"clear_voice_asset":                     false,
		"clear_vrm_asset":                       false,
		"reset_secretary_appearance_defaults":   false,
		"show_provider_readiness_without_claim": false,
	}
	for _, action := range asObjectSlice(t, appearance["appearance_reset_actions"]) {
		actionID := requireString(t, action, "action_id")
		if _, ok := actionIDs[actionID]; ok {
			actionIDs[actionID] = true
		}
		requireString(t, action, "receipt_required")
		switch actionID {
		case "clear_voice_asset":
			if got := requireString(t, action, "result_voice_asset_ref"); got != "none" {
				t.Fatalf("clear_voice_asset result_voice_asset_ref = %s", got)
			}
		case "clear_vrm_asset":
			if got := requireString(t, action, "result_vrm_asset_ref"); got != "none" {
				t.Fatalf("clear_vrm_asset result_vrm_asset_ref = %s", got)
			}
		case "reset_secretary_appearance_defaults":
			if got := requireString(t, action, "result_voice_asset_ref"); got != "voice_asset://default-secretary-neutral" {
				t.Fatalf("reset result_voice_asset_ref = %s", got)
			}
			if got := requireString(t, action, "result_vrm_asset_ref"); got != "vrm_asset://default-secretary-vrm" {
				t.Fatalf("reset result_vrm_asset_ref = %s", got)
			}
		}
	}
	for actionID, seen := range actionIDs {
		if !seen {
			t.Fatalf("appearance_reset_actions missing %s", actionID)
		}
	}

	negative := requireObject(t, appearance, "negative_cases")
	for caseID, wantStatus := range map[string]string{
		"raw_voice_file_path": "blocked_raw_asset_reference",
		"raw_vrm_file_path":   "blocked_raw_asset_reference",
		"provider_missing":    "degraded_text_role_still_available",
	} {
		item := requireObject(t, negative, caseID)
		if got := requireString(t, item, "expected_status"); got != wantStatus {
			t.Fatalf("%s expected_status = %s, want %s", caseID, got, wantStatus)
		}
		requireString(t, item, "expected_evidence")
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	completionPolicy := requireObject(t, matrix, "completion_claim_policy")
	required := asStringSlice(t, completionPolicy["required_before_completion_claim"])
	requireStringSliceContains(t, required, "secretary_appearance_gui_controls_verified")
}

func TestTeamOfficeCloudListingAndReleaseDeclareSixRoleProductComponents(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	roleRefs := asStringSlice(t, candidateSet["role_pack_refs"])
	if len(roleRefs) != 6 {
		t.Fatalf("role_pack_refs len = %d, want 6", len(roleRefs))
	}
	expectedRoleRefs := make(map[string]bool, len(roleRefs))
	for _, roleRef := range roleRefs {
		expectedRoleRefs[roleRef] = false
	}

	slotsDoc := readJSON(t, filepath.Join(base, "role-slots", "team-office-role-slots.json"))
	slotByRoleRef := map[string]map[string]any{}
	for _, slot := range asObjectSlice(t, slotsDoc["role_slots"]) {
		roleRef := requireString(t, slot, "default_role_pack_ref")
		slotByRoleRef[roleRef] = slot
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	manifestPaths := map[string]bool{}
	for _, item := range asObjectSlice(t, manifest["files"]) {
		manifestPaths[requireString(t, item, "path")] = true
	}

	listing := readJSON(t, filepath.Join(base, "commerce", "cloud-listing-candidate.json"))
	summary := requireObject(t, listing, "role_component_summary")
	requireInt(t, summary, "role_count", len(roleRefs))
	requireBool(t, summary, "buyer_visible", true)
	requireBool(t, summary, "team_settings_replaceable", true)
	if got := requireString(t, summary, "component_truth_source"); got != "truzhen-packs" {
		t.Fatalf("component_truth_source = %s, want truzhen-packs", got)
	}

	appearanceBoundary := requireObject(t, summary, "secretary_appearance_boundary")
	if got := requireString(t, appearanceBoundary, "secretary_role_pack_ref"); got != "role_pack://team-office-secretary-general" {
		t.Fatalf("secretary_role_pack_ref = %s", got)
	}
	for _, key := range []string{"voice_asset_ref_only", "vrm_asset_ref_only", "raw_voice_or_vrm_not_included"} {
		requireBool(t, appearanceBoundary, key, true)
	}
	if got := requireString(t, appearanceBoundary, "appearance_preferences_ref"); got != "appearance/secretary-appearance-preferences.json" {
		t.Fatalf("appearance_preferences_ref = %s", got)
	}
	if got := requireString(t, appearanceBoundary, "asset_rights_ref"); got != "appearance/secretary-appearance-asset-rights-candidate.json" {
		t.Fatalf("asset_rights_ref = %s", got)
	}

	components := asObjectSlice(t, summary["components"])
	if len(components) != len(roleRefs) {
		t.Fatalf("components len = %d, want %d", len(components), len(roleRefs))
	}
	for _, component := range components {
		roleRef := requireString(t, component, "role_pack_ref")
		if _, ok := expectedRoleRefs[roleRef]; !ok {
			t.Fatalf("unexpected role component %s", roleRef)
		}
		if expectedRoleRefs[roleRef] {
			t.Fatalf("duplicate role component %s", roleRef)
		}
		expectedRoleRefs[roleRef] = true

		slot, ok := slotByRoleRef[roleRef]
		if !ok {
			t.Fatalf("role component %s has no role slot", roleRef)
		}
		if got, want := requireString(t, component, "slot_id"), requireString(t, slot, "slot_id"); got != want {
			t.Fatalf("%s slot_id = %s, want %s", roleRef, got, want)
		}
		if got, want := requireString(t, component, "display_name"), requireString(t, slot, "display_name"); got != want {
			t.Fatalf("%s display_name = %s, want %s", roleRef, got, want)
		}
		if got := requireString(t, component, "formal_authority"); got != "none" {
			t.Fatalf("%s formal_authority = %s, want none", roleRef, got)
		}
		requireBool(t, component, "candidate_only", true)
		requireBool(t, component, "buyer_visible", true)
		requireBool(t, component, "team_settings_replaceable", true)

		rolePackFile := requireString(t, component, "role_pack_file")
		if !manifestPaths[rolePackFile] {
			t.Fatalf("%s role_pack_file %s is not in artifact manifest", roleRef, rolePackFile)
		}
		readJSON(t, filepath.Join(base, rolePackFile))
	}
	for roleRef, seen := range expectedRoleRefs {
		if !seen {
			t.Fatalf("role component summary missing %s", roleRef)
		}
	}

	release := readJSON(t, filepath.Join(base, "commerce", "release-candidate-package.json"))
	releaseContents := requireObject(t, release, "release_contents")
	if got := requireString(t, releaseContents, "role_component_summary_source"); got != "commerce/cloud-listing-candidate.json#role_component_summary" {
		t.Fatalf("role_component_summary_source = %s", got)
	}
	requireInt(t, releaseContents, "role_count", len(roleRefs))
	releaseRoleRefs := asStringSlice(t, releaseContents["role_pack_refs"])
	for _, roleRef := range roleRefs {
		requireStringSliceContains(t, releaseRoleRefs, roleRef)
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	listingMetadata := requireObject(t, submission, "listing_metadata")
	if got := requireString(t, listingMetadata, "role_component_summary_ref"); got != "commerce/cloud-listing-candidate.json#role_component_summary" {
		t.Fatalf("listing_metadata.role_component_summary_ref = %s", got)
	}
	checklistItem := false
	for _, item := range asObjectSlice(t, submission["review_checklist"]) {
		if requireString(t, item, "check_id") == "six_role_component_listing" {
			checklistItem = true
			if got := requireString(t, item, "status"); got != "candidate_ready" {
				t.Fatalf("six_role_component_listing status = %s, want candidate_ready", got)
			}
		}
	}
	if !checklistItem {
		t.Fatalf("review_checklist missing six_role_component_listing")
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policy := requireObject(t, matrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, policy["required_before_completion_claim"]), "six_role_component_listing_verified")
}

func TestSecretaryAppearanceAssetRightsCandidateBlocksUnlicensedVoiceAndVrm(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	rights := readJSON(t, filepath.Join(base, "appearance", "secretary-appearance-asset-rights-candidate.json"))
	requireBool(t, rights, "candidate_only", true)
	requireBool(t, rights, "non_formal", true)
	if got := requireString(t, rights, "asset_rights_truth_source"); got == "truzhen-packs" {
		t.Fatalf("asset_rights_truth_source cannot be truzhen-packs")
	}
	if got := requireString(t, rights, "appearance_preferences_ref"); got != "appearance/secretary-appearance-preferences.json" {
		t.Fatalf("appearance_preferences_ref = %s", got)
	}

	rawAssets, ok := rights["assets"].([]any)
	if !ok {
		t.Fatalf("assets missing")
	}
	requiredAssets := map[string]bool{
		"voice_asset://default-secretary-neutral": false,
		"vrm_asset://default-secretary-vrm":       false,
	}
	for _, raw := range rawAssets {
		asset, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("asset item = %T", raw)
		}
		assetRef := requireString(t, asset, "asset_ref")
		if _, ok := requiredAssets[assetRef]; ok {
			requiredAssets[assetRef] = true
		}
		if strings.HasPrefix(assetRef, "/") || strings.HasPrefix(assetRef, "file:") {
			t.Fatalf("asset_ref must not be raw path: %s", assetRef)
		}
		for _, key := range []string{"asset_kind", "license_status", "usage_scope", "provider_readiness", "review_evidence_required"} {
			if got := requireString(t, asset, key); got == "" {
				t.Fatalf("%s %s missing", assetRef, key)
			}
		}
	}
	for assetRef, seen := range requiredAssets {
		if !seen {
			t.Fatalf("missing asset rights entry %s", assetRef)
		}
	}

	review := requireObject(t, rights, "marketplace_review")
	for _, key := range []string{"asset_ref_only", "license_evidence_required", "raw_asset_scan_required", "human_review_required_if_external_persona"} {
		if got := requireString(t, review, key); got != "true" {
			t.Fatalf("marketplace_review.%s = %s, want true", key, got)
		}
	}

	preflight := requireObject(t, rights, "install_preflight")
	for _, key := range []string{"asset_ref_validation", "provider_readiness_check", "license_revocation_check"} {
		if got := requireString(t, preflight, key); got != "required" {
			t.Fatalf("install_preflight.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, rights, "negative_cases")
	for _, key := range []string{
		"raw_voice_file",
		"raw_vrm_file",
		"voice_clone_without_authorization",
		"unlicensed_vrm_avatar",
		"revoked_asset_license",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "appearance/secretary-appearance-asset-rights-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing appearance/secretary-appearance-asset-rights-candidate.json")
}

func TestTeamOfficeArtifactSecretRawAssetScanCandidateProtectsCloudUploadBundle(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	scanPath := "tests/artifact-secret-raw-asset-scan-candidate.json"

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "artifact_secret_raw_asset_scan"); got != scanPath {
		t.Fatalf("artifact_secret_raw_asset_scan = %s, want %s", got, scanPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, scanPath) {
		t.Fatalf("candidate set artifact_files missing %s", scanPath)
	}

	scan := readJSON(t, filepath.Join(base, scanPath))
	requireBool(t, scan, "candidate_only", true)
	requireBool(t, scan, "non_formal", true)
	if got := requireString(t, scan, "scan_status"); got != "current_repo_candidate_scan_required_before_cloud_upload" {
		t.Fatalf("scan_status = %s", got)
	}
	for key, want := range map[string]string{
		"artifact_manifest_ref":                 "commerce/artifact-manifest.json",
		"artifact_bundle_digest_ref":            "commerce/artifact-bundle-digest-candidate.json",
		"product_readiness_evidence_matrix_ref": "tests/product-readiness-evidence-matrix.json",
		"commercial_readiness_verifier_ref":     "tests/commercial-readiness-verifier-candidate.json",
	} {
		if got := requireString(t, scan, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}
	for _, key := range []string{"required_before", "forbidden_file_extensions", "forbidden_payload_patterns", "scan_result_fields", "blocking_if_detected"} {
		if values := asStringSlice(t, scan[key]); len(values) == 0 {
			t.Fatalf("%s must not be empty", key)
		}
	}
	requiredBefore := strings.Join(asStringSlice(t, scan["required_before"]), "\n")
	for _, gate := range []string{"cloud_upload_request", "marketplace_review_submission", "signed_download_enable", "install_preflight", "production_promotion_controls"} {
		if !strings.Contains(requiredBefore, gate) {
			t.Fatalf("required_before missing %s", gate)
		}
	}

	forbiddenExts := map[string]bool{}
	for _, ext := range asStringSlice(t, scan["forbidden_file_extensions"]) {
		forbiddenExts[strings.ToLower(ext)] = true
	}
	for _, ext := range []string{".wav", ".mp3", ".flac", ".ogg", ".vrm", ".glb", ".gltf", ".fbx", ".blend", ".vroid", ".png", ".jpg", ".jpeg"} {
		if !forbiddenExts[ext] {
			t.Fatalf("forbidden_file_extensions missing %s", ext)
		}
	}

	patternIDs := strings.Join(asStringSlice(t, scan["forbidden_payload_patterns"]), "\n")
	secretPatterns := map[string]*regexp.Regexp{
		"pem_private_key":      regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
		"aws_access_key":       regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"github_token":         regexp.MustCompile(`gh[pousr]_[A-Za-z0-9_]{20,}`),
		"openai_api_key":       regexp.MustCompile(`sk-[A-Za-z0-9]{20,}`),
		"slack_token":          regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{20,}`),
		"stripe_secret_key":    regexp.MustCompile(`sk_(?:live|test)_[A-Za-z0-9]{16,}`),
		"signed_url_signature": regexp.MustCompile(`X-Amz-Signature=[0-9a-fA-F]{16,}`),
	}
	for patternID := range secretPatterns {
		if !strings.Contains(patternIDs, patternID) {
			t.Fatalf("forbidden_payload_patterns missing %s", patternID)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	scanInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		rel := requireString(t, item, "path")
		if rel == scanPath {
			scanInManifest = true
			requireStringIn(t, requireString(t, item, "required_for"), "audit")
		}
		if forbiddenExts[strings.ToLower(filepath.Ext(rel))] {
			t.Fatalf("artifact manifest includes forbidden raw asset file %s", rel)
		}
		data, err := os.ReadFile(filepath.Join(base, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		text := string(data)
		for patternID, pattern := range secretPatterns {
			if pattern.MatchString(text) {
				t.Fatalf("%s contains forbidden payload pattern %s", rel, patternID)
			}
		}
	}
	if !scanInManifest {
		t.Fatalf("artifact manifest missing %s", scanPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	rawPayload, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	scanInDigest := false
	for _, raw := range rawPayload {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == scanPath {
			scanInDigest = true
		}
	}
	if !scanInDigest {
		t.Fatalf("artifact bundle digest payload missing %s", scanPath)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requiredProofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	for _, proof := range []string{"no_secret_scan", "artifact_secret_raw_asset_scan_verified"} {
		if !strings.Contains(requiredProofs, proof) {
			t.Fatalf("product readiness matrix missing %s", proof)
		}
	}
	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, verifier, "source_artifact_secret_raw_asset_scan"); got != scanPath {
		t.Fatalf("commercial readiness verifier source_artifact_secret_raw_asset_scan = %s, want %s", got, scanPath)
	}
}

func TestTeamOfficeArtifactManifestSupportsUploadDownloadInstall(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	requireBool(t, manifest, "candidate_only", true)
	if got := requireString(t, manifest, "artifact_ref"); got != "role-pack-artifact://team-office-v0" {
		t.Fatalf("artifact_ref = %s", got)
	}
	if got := requireString(t, manifest, "artifact_type"); got != "role_pack_candidate_bundle" {
		t.Fatalf("artifact_type = %s", got)
	}
	requireStringIn(t, requireString(t, manifest, "upload_status"), "candidate_ready_for_cloud_upload")
	requireStringIn(t, requireString(t, manifest, "download_status"), "candidate_requires_entitlement")
	requireStringIn(t, requireString(t, manifest, "install_status"), "candidate_requires_truzhenos_install")
	requireStringIn(t, requireString(t, manifest, "hash_algorithm"), "sha256")

	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("files missing")
	}
	if len(rawFiles) < 12 {
		t.Fatalf("files len = %d, want at least 12", len(rawFiles))
	}
	seen := map[string]bool{}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		path := requireString(t, item, "path")
		if filepath.IsAbs(path) || strings.Contains(path, "..") {
			t.Fatalf("artifact path must be relative and contained, got %s", path)
		}
		if _, err := os.Stat(filepath.Join(base, path)); err != nil {
			t.Fatalf("artifact file %s missing: %v", path, err)
		}
		requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
		seen[path] = true
	}
	for _, required := range []string{
		"candidate-set.json",
		"role-packs/team-office-secretary-general.rolepack.json",
		"role-slots/team-office-role-slots.json",
		"capability-role-requirements/sample-team-research.role-requirements.json",
		"appearance/secretary-appearance-preferences.json",
		"commerce/cloud-listing-candidate.json",
	} {
		if !seen[required] {
			t.Fatalf("artifact manifest missing %s", required)
		}
	}

	download := requireObject(t, manifest, "download_requirements")
	if got := requireString(t, download, "entitlement_required"); got != "true" {
		t.Fatalf("download entitlement_required = %s, want true", got)
	}
	if got := requireString(t, download, "hash_verification"); got != "required" {
		t.Fatalf("download hash_verification = %s, want required", got)
	}

	install := requireObject(t, manifest, "install_requirements")
	if got := requireString(t, install, "installer_truth_source"); got != "truzhenos" {
		t.Fatalf("installer_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, install, "owner_gate_required"); got != "true" {
		t.Fatalf("owner_gate_required = %s, want true", got)
	}
	if got := requireString(t, install, "receipt_required"); got != "true" {
		t.Fatalf("receipt_required = %s, want true", got)
	}
}

func TestTeamOfficeArtifactManifestFileHashesMatch(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	if got := requireString(t, manifest, "manifest_hash_policy"); got != "self_hash_excluded" {
		t.Fatalf("manifest_hash_policy = %s, want self_hash_excluded", got)
	}
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("files missing")
	}
	hex64 := regexp.MustCompile(`^[a-f0-9]{64}$`)
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		path := requireString(t, item, "path")
		if path == "commerce/artifact-manifest.json" {
			if got := requireString(t, item, "sha256"); got != "self_hash_excluded" {
				t.Fatalf("%s sha256 = %s, want self_hash_excluded", path, got)
			}
			continue
		}
		got := requireString(t, item, "sha256")
		if !hex64.MatchString(got) {
			t.Fatalf("%s sha256 is not 64 lowercase hex chars: %s", path, got)
		}
		data, err := os.ReadFile(filepath.Join(base, path))
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		want := fmt.Sprintf("%x", sha256.Sum256(data))
		if got != want {
			t.Fatalf("%s sha256 = %s, want %s", path, got, want)
		}
	}
}

func TestTeamOfficeArtifactManifestClosureGateMatchesCurrentBundleTree(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))

	gatePath := requireString(t, candidateSet, "artifact_manifest_closure_gate")
	requireStringSliceContains(t, asStringSlice(t, candidateSet["artifact_files"]), gatePath)
	gate := readJSON(t, filepath.Join(base, gatePath))
	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	if got := requireString(t, gate, "closure_status"); got != "candidate_ready_manifest_matches_current_tree" {
		t.Fatalf("closure_status = %s, want candidate_ready_manifest_matches_current_tree", got)
	}
	if got := requireString(t, gate, "artifact_manifest"); got != "commerce/artifact-manifest.json" {
		t.Fatalf("artifact_manifest = %s, want commerce/artifact-manifest.json", got)
	}
	if got := requireString(t, gate, "candidate_set"); got != "candidate-set.json" {
		t.Fatalf("candidate_set = %s, want candidate-set.json", got)
	}
	requireBool(t, gate, "blocks_upload_if_unregistered_file_exists", true)
	requireBool(t, gate, "blocks_download_if_manifest_file_missing", true)
	requireBool(t, gate, "blocks_install_if_candidate_set_manifest_drift", true)
	requireStringSliceContains(t, asStringSlice(t, gate["required_closure_checks"]), "actual_tree_equals_candidate_set_artifact_files")
	requireStringSliceContains(t, asStringSlice(t, gate["required_closure_checks"]), "actual_tree_equals_artifact_manifest_files")
	requireStringSliceContains(t, asStringSlice(t, gate["required_closure_checks"]), "candidate_set_artifact_files_equals_artifact_manifest_files")
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_cloud_upload"]), "artifact_manifest_closure_gate_verified")
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_signed_download"]), "artifact_manifest_closure_gate_verified")
	requireStringSliceContains(t, asStringSlice(t, gate["required_before_local_install"]), "artifact_manifest_closure_gate_verified")

	actualFiles := map[string]bool{}
	if err := filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		actualFiles[filepath.ToSlash(rel)] = true
		return nil
	}); err != nil {
		t.Fatalf("walk %s: %v", base, err)
	}

	candidateFiles := map[string]bool{}
	for _, rel := range asStringSlice(t, candidateSet["artifact_files"]) {
		candidateFiles[rel] = true
	}
	manifestFiles := map[string]bool{}
	for _, raw := range asObjectSlice(t, manifest["files"]) {
		manifestFiles[requireString(t, raw, "path")] = true
	}
	compareFileSets(t, "actual tree", actualFiles, "candidate_set.artifact_files", candidateFiles)
	compareFileSets(t, "actual tree", actualFiles, "artifact_manifest.files", manifestFiles)
	compareFileSets(t, "candidate_set.artifact_files", candidateFiles, "artifact_manifest.files", manifestFiles)

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "artifact_manifest_closure_gate_verified")

	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, verifier, "source_artifact_manifest_closure_gate"); got != gatePath {
		t.Fatalf("source_artifact_manifest_closure_gate = %s, want %s", got, gatePath)
	}
}

func TestTeamOfficeSandboxCommerceFlowCandidate(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	flow := readJSON(t, filepath.Join(base, "commerce", "sandbox-commerce-flow-candidates.json"))
	requireBool(t, flow, "candidate_only", true)
	requireBool(t, flow, "non_formal", true)
	if got := requireString(t, flow, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, flow, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	rawStages, ok := flow["stages"].([]any)
	if !ok {
		t.Fatalf("stages missing")
	}
	wantStages := []string{
		"upload_draft",
		"review_candidate",
		"sandbox_order",
		"sandbox_payment",
		"entitlement_issue",
		"download_signed_artifact",
		"local_install_via_truzhenos",
	}
	if len(rawStages) != len(wantStages) {
		t.Fatalf("stages len = %d, want %d", len(rawStages), len(wantStages))
	}
	for i, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage item = %T", raw)
		}
		if got := requireString(t, stage, "stage"); got != wantStages[i] {
			t.Fatalf("stage[%d] = %s, want %s", i, got, wantStages[i])
		}
		requireStringIn(t, requireString(t, stage, "status"), "candidate_pending", "blocked_without_owner_authorization")
		if requireString(t, stage, "truth_source") == "truzhen-packs" {
			t.Fatalf("stage %s cannot use truzhen-packs as truth source", wantStages[i])
		}
		requireString(t, stage, "evidence_ref")
	}

	negative := requireObject(t, flow, "negative_cases")
	for _, key := range []string{"without_entitlement", "artifact_hash_mismatch", "real_payment_without_owner_authorization"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}
}

func TestTeamOfficeProductHandoffCandidateDefinesCloudAndInstallEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	handoff := readJSON(t, filepath.Join(base, "commerce", "product-handoff-candidate.json"))
	requireBool(t, handoff, "candidate_only", true)
	requireBool(t, handoff, "non_formal", true)
	if got := requireString(t, handoff, "artifact_ref"); got != "role-pack-artifact://team-office-v0" {
		t.Fatalf("artifact_ref = %s", got)
	}
	if got := requireString(t, handoff, "commercial_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commercial_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, handoff, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}

	cloud := requireObject(t, handoff, "cloud_handoff")
	for _, key := range []string{"upload_request", "listing_draft_request", "sandbox_purchase_request", "entitlement_issue_request", "download_request"} {
		req := requireObject(t, cloud, key)
		if got := requireString(t, req, "truth_source"); got != "truzhen-cloud" {
			t.Fatalf("%s truth_source = %s, want truzhen-cloud", key, got)
		}
		if got := requireString(t, req, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
		if got := requireString(t, req, "success_status"); strings.HasPrefix(got, "formal") {
			t.Fatalf("%s success_status must remain candidate/sandbox scoped, got %s", key, got)
		}
	}
	purchase := requireObject(t, cloud, "sandbox_purchase_request")
	if got := requireString(t, purchase, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("sandbox real_payment_policy = %s", got)
	}

	install := requireObject(t, handoff, "install_handoff")
	for _, key := range []string{"artifact_verification", "entitlement_verification", "role_pack_install", "team_role_binding"} {
		req := requireObject(t, install, key)
		if got := requireString(t, req, "truth_source"); got != "truzhenos" {
			t.Fatalf("%s truth_source = %s, want truzhenos", key, got)
		}
		if got := requireString(t, req, "receipt_required"); got != "true" {
			t.Fatalf("%s receipt_required = %s, want true", key, got)
		}
	}
	binding := requireObject(t, install, "team_role_binding")
	if got := requireString(t, binding, "owner_gate_required"); got != "true" {
		t.Fatalf("team_role_binding owner_gate_required = %s, want true", got)
	}

	negative := requireObject(t, handoff, "negative_cases")
	for _, key := range []string{"download_without_entitlement", "install_hash_mismatch", "bind_without_owner_gate", "real_payment_without_owner_authorization"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/product-handoff-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/product-handoff-candidate.json")
}

func TestTeamOfficeGuiUserAgentScenariosCoverFullProductPath(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	scenarios := readJSON(t, filepath.Join(base, "tests", "gui-user-agent-scenarios.json"))
	requireBool(t, scenarios, "candidate_only", true)
	requireBool(t, scenarios, "non_formal", true)
	if got := requireString(t, scenarios, "actor_role"); got != "user_view_gui_agent" {
		t.Fatalf("actor_role = %s, want user_view_gui_agent", got)
	}

	policy := requireObject(t, scenarios, "operation_policy")
	requireBool(t, policy, "gui_only", true)
	forbidden := strings.Join(asStringSlice(t, policy["forbidden_direct_access"]), "\n")
	for _, phrase := range []string{"direct_api_call", "direct_database_write", "manual_json_edit", "filesystem_copy_to_cloud"} {
		if !strings.Contains(forbidden, phrase) {
			t.Fatalf("forbidden_direct_access missing %q", phrase)
		}
	}

	rawScenarios, ok := scenarios["scenarios"].([]any)
	if !ok {
		t.Fatalf("scenarios missing")
	}
	required := map[string]bool{
		"GUI-ROLE-CREATE-SECRETARY":                  false,
		"GUI-ROLE-CREATE-FIVE-ADVISORS":              false,
		"GUI-CAPABILITY-ROLE-REFERENCE":              false,
		"GUI-TEAM-SETTINGS-ROLE-REPLACE":             false,
		"GUI-SECRETARY-VOICE-VRM":                    false,
		"GUI-CLOUD-UPLOAD-PURCHASE-DOWNLOAD-INSTALL": false,
		"GUI-NEGATIVE-SOVEREIGNTY-AND-ASSET-GATES":   false,
	}
	for _, raw := range rawScenarios {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("scenario item = %T", raw)
		}
		id := requireString(t, item, "scenario_id")
		if _, ok := required[id]; ok {
			required[id] = true
		}
		requireBool(t, item, "gui_only", true)
		if got := requireString(t, item, "entry_surface"); got == "" {
			t.Fatalf("%s entry_surface missing", id)
		}
		steps, ok := item["steps"].([]any)
		if !ok || len(steps) == 0 {
			t.Fatalf("%s steps missing", id)
		}
		evidence := strings.Join(asStringSlice(t, item["evidence_required"]), "\n")
		for _, proof := range []string{"screenshot", "page_state", "candidate_or_receipt_ref"} {
			if !strings.Contains(evidence, proof) {
				t.Fatalf("%s evidence_required missing %q", id, proof)
			}
		}
	}
	for id, seen := range required {
		if !seen {
			t.Fatalf("missing GUI scenario %s", id)
		}
	}

	cloudScenario := map[string]any{}
	for _, raw := range rawScenarios {
		item := raw.(map[string]any)
		if requireString(t, item, "scenario_id") == "GUI-CLOUD-UPLOAD-PURCHASE-DOWNLOAD-INSTALL" {
			cloudScenario = item
			break
		}
	}
	targets := strings.Join(asStringSlice(t, cloudScenario["backend_truth_sources"]), "\n")
	for _, source := range []string{"truzhen-cloud", "truzhenos"} {
		if !strings.Contains(targets, source) {
			t.Fatalf("cloud/install scenario backend_truth_sources missing %s", source)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/gui-user-agent-scenarios.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/gui-user-agent-scenarios.json")
}

func TestTeamOfficeGuiUserAgentExecutionScriptDefinesOrderedUserOnlySteps(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	scriptPath := "tests/gui-user-agent-execution-script-candidate.json"
	script := readJSON(t, filepath.Join(base, scriptPath))
	requireBool(t, script, "candidate_only", true)
	requireBool(t, script, "non_formal", true)
	if got := requireString(t, script, "script_status"); got != "not_run_requires_cross_repo_gui" {
		t.Fatalf("script_status = %s, want not_run_requires_cross_repo_gui", got)
	}
	if got := requireString(t, script, "actor_role"); got != "user_view_gui_agent" {
		t.Fatalf("actor_role = %s, want user_view_gui_agent", got)
	}
	if got := requireString(t, script, "operation_mode"); got != "gui_only" {
		t.Fatalf("operation_mode = %s, want gui_only", got)
	}

	sourceRefs := strings.Join(asStringSlice(t, script["source_refs"]), "\n")
	for _, ref := range []string{
		"tests/gui-user-agent-scenarios.json",
		"tests/gui-evidence-capture-protocol.json",
		"integration/cross-repo-execution-readiness-package.json",
	} {
		if !strings.Contains(sourceRefs, ref) {
			t.Fatalf("source_refs missing %s", ref)
		}
	}

	stepContract := requireObject(t, script, "step_contract")
	requiredFields := strings.Join(asStringSlice(t, stepContract["required_fields"]), "\n")
	for _, field := range []string{
		"step_id",
		"step_order",
		"scenario_ref",
		"surface_id",
		"action_type",
		"user_action",
		"expected_page_state",
		"truth_source",
		"required_evidence_slots",
		"failure_issue_id",
	} {
		if !strings.Contains(requiredFields, field) {
			t.Fatalf("step_contract required_fields missing %s", field)
		}
	}

	rawSteps, ok := script["ordered_steps"].([]any)
	if !ok {
		t.Fatalf("ordered_steps missing")
	}
	requiredSteps := map[string]bool{
		"open_role_studio":                    false,
		"create_secretary_candidate":          false,
		"create_five_advisor_candidates":      false,
		"select_secretary_voice_vrm":          false,
		"create_capability_role_reference":    false,
		"export_candidate_bundle":             false,
		"upload_cloud_listing_draft":          false,
		"submit_marketplace_review_candidate": false,
		"sandbox_purchase":                    false,
		"download_purchased_artifact":         false,
		"install_downloaded_role_pack":        false,
		"replace_team_roles_after_install":    false,
		"run_team_office_runtime_use":         false,
		"run_negative_cases":                  false,
	}
	lastOrder := 0.0
	for _, raw := range rawSteps {
		step, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("ordered step = %T", raw)
		}
		stepID := requireString(t, step, "step_id")
		if _, ok := requiredSteps[stepID]; ok {
			requiredSteps[stepID] = true
		}
		order, ok := step["step_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s step_order = %v, want increasing number greater than %v", stepID, step["step_order"], lastOrder)
		}
		lastOrder = order
		requireStringIn(t, requireString(t, step, "action_type"), "open_page", "click", "type", "select", "upload_via_gui", "download_via_gui", "confirm_owner_gate", "observe")
		requireStringIn(t, requireString(t, step, "truth_source"), "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo")
		for _, key := range []string{"scenario_ref", "surface_id", "user_action", "expected_page_state", "failure_issue_id"} {
			if got := requireString(t, step, key); got == "" {
				t.Fatalf("%s %s missing", stepID, key)
			}
		}
		slots := strings.Join(asStringSlice(t, step["required_evidence_slots"]), "\n")
		for _, slot := range []string{"screenshot_path", "page_state_ref", "candidate_or_receipt_ref"} {
			if !strings.Contains(slots, slot) {
				t.Fatalf("%s required_evidence_slots missing %s", stepID, slot)
			}
		}
		if strings.Contains(stepID, "cloud") || strings.Contains(stepID, "purchase") || strings.Contains(stepID, "download") {
			if !strings.Contains(slots, "cloud_receipt_ref") && !strings.Contains(slots, "download_receipt_ref") {
				t.Fatalf("%s must require cloud or download receipt evidence", stepID)
			}
		}
		if strings.Contains(stepID, "install") {
			if !strings.Contains(slots, "install_receipt_ref") {
				t.Fatalf("%s must require install_receipt_ref", stepID)
			}
		}
		if strings.Contains(stepID, "replace_team") {
			if !strings.Contains(slots, "team_binding_receipt_ref") {
				t.Fatalf("%s must require team_binding_receipt_ref", stepID)
			}
		}
	}
	for stepID, seen := range requiredSteps {
		if !seen {
			t.Fatalf("missing ordered step %s", stepID)
		}
	}

	negative := requireObject(t, script, "negative_case_script")
	cases := strings.Join(asStringSlice(t, negative["required_cases"]), "\n")
	for _, caseID := range []string{"download_without_entitlement", "tampered_artifact_hash", "raw_voice_or_vrm_asset", "real_payment_without_owner_authorization"} {
		if !strings.Contains(cases, caseID) {
			t.Fatalf("negative_case_script missing %s", caseID)
		}
	}
	requireBool(t, negative, "must_capture_blocked_receipt", true)

	forbidden := strings.Join(asStringSlice(t, script["forbidden"]), "\n")
	for _, item := range []string{
		"direct_api_call_as_user_action",
		"manual_json_edit_as_gui_success",
		"filesystem_copy_to_cloud",
		"backend_fixture_success_without_gui",
		"completion_claim_from_script_candidate_only",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "gui_user_agent_execution_script"); got != scriptPath {
		t.Fatalf("gui_user_agent_execution_script = %s, want %s", got, scriptPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, scriptPath) {
		t.Fatalf("candidate set artifact_files missing %s", scriptPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == scriptPath {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", scriptPath)
}

func TestTeamOfficeGuiEvidenceCaptureProtocolDefinesUserViewProofChain(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	protocol := readJSON(t, filepath.Join(base, "tests", "gui-evidence-capture-protocol.json"))
	requireBool(t, protocol, "candidate_only", true)
	requireBool(t, protocol, "non_formal", true)
	if got := requireString(t, protocol, "actor_role"); got != "user_view_gui_agent" {
		t.Fatalf("actor_role = %s, want user_view_gui_agent", got)
	}
	if got := requireString(t, protocol, "coordinator_role"); got != "organizer_coordinator_recorder" {
		t.Fatalf("coordinator_role = %s, want organizer_coordinator_recorder", got)
	}

	policy := requireObject(t, protocol, "operation_policy")
	requireBool(t, policy, "gui_only", true)
	forbidden := strings.Join(asStringSlice(t, policy["forbidden_actions"]), "\n")
	for _, action := range []string{"direct_api_call", "manual_json_edit", "filesystem_copy_to_cloud", "backend_fixture_success"} {
		if !strings.Contains(forbidden, action) {
			t.Fatalf("forbidden_actions missing %s", action)
		}
	}

	capture := requireObject(t, protocol, "capture_contract")
	requiredFields := strings.Join(asStringSlice(t, capture["required_fields_for_each_step"]), "\n")
	for _, field := range []string{"screenshot_path", "page_state_ref", "timestamp", "actor_id", "surface_id", "action_summary", "candidate_or_receipt_ref", "network_or_receipt_summary", "blocked_reason_if_any"} {
		if !strings.Contains(requiredFields, field) {
			t.Fatalf("capture required fields missing %s", field)
		}
	}
	if got := requireString(t, capture, "redaction_required"); got != "true" {
		t.Fatalf("redaction_required = %s, want true", got)
	}

	rawStages, ok := protocol["stage_capture_requirements"].([]any)
	if !ok {
		t.Fatalf("stage_capture_requirements missing")
	}
	requiredStages := map[string]bool{
		"role_create":             false,
		"candidate_bundle_export": false,
		"cloud_upload_draft":      false,
		"sandbox_purchase":        false,
		"entitlement_download":    false,
		"local_install":           false,
		"team_settings_replace":   false,
		"runtime_use":             false,
		"negative_case":           false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage requirement = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		requireStringIn(t, requireString(t, stage, "truth_source"), "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo")
		evidence := strings.Join(asStringSlice(t, stage["required_evidence"]), "\n")
		for _, proof := range []string{"screenshot_path", "page_state_ref", "candidate_or_receipt_ref"} {
			if !strings.Contains(evidence, proof) {
				t.Fatalf("%s required_evidence missing %s", stageID, proof)
			}
		}
		if got := requireString(t, stage, "issue_when_missing"); got == "" {
			t.Fatalf("%s issue_when_missing missing", stageID)
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing stage capture requirement %s", stageID)
		}
	}

	chain := requireObject(t, protocol, "proof_chain")
	for _, key := range []string{"correlation_keys_required", "hash_or_receipt_continuity_required", "completion_claim_allowed"} {
		if _, ok := chain[key]; !ok {
			t.Fatalf("proof_chain missing %s", key)
		}
	}
	requireBool(t, chain, "completion_claim_allowed", false)

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/gui-evidence-capture-protocol.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/gui-evidence-capture-protocol.json")
}

func TestTeamOfficeRoleStudioExportProvenanceAttestationRequiresGuiExportEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	attestationPath := requireString(t, candidateSet, "role_studio_export_provenance_attestation")
	requireStringSliceContains(t, asStringSlice(t, candidateSet["artifact_files"]), attestationPath)

	attestation := readJSON(t, filepath.Join(base, attestationPath))
	requireBool(t, attestation, "candidate_only", true)
	requireBool(t, attestation, "non_formal", true)
	if got := requireString(t, attestation, "attestation_status"); got != "not_run_requires_user_view_gui_export_receipt" {
		t.Fatalf("attestation_status = %s, want not_run_requires_user_view_gui_export_receipt", got)
	}
	if got := requireString(t, attestation, "actor_role"); got != "user_view_gui_agent" {
		t.Fatalf("actor_role = %s, want user_view_gui_agent", got)
	}
	if got := requireString(t, attestation, "coordinator_role"); got != "organizer_coordinator_recorder" {
		t.Fatalf("coordinator_role = %s, want organizer_coordinator_recorder", got)
	}
	if got := requireString(t, attestation, "source_gui_script"); got != "tests/gui-user-agent-execution-script-candidate.json" {
		t.Fatalf("source_gui_script = %s, want tests/gui-user-agent-execution-script-candidate.json", got)
	}
	if got := requireString(t, attestation, "source_gui_evidence_protocol"); got != "tests/gui-evidence-capture-protocol.json" {
		t.Fatalf("source_gui_evidence_protocol = %s, want tests/gui-evidence-capture-protocol.json", got)
	}
	if got := requireString(t, attestation, "artifact_manifest"); got != "commerce/artifact-manifest.json" {
		t.Fatalf("artifact_manifest = %s, want commerce/artifact-manifest.json", got)
	}
	if got := requireString(t, attestation, "artifact_bundle_digest"); got != "commerce/artifact-bundle-digest-candidate.json" {
		t.Fatalf("artifact_bundle_digest = %s, want commerce/artifact-bundle-digest-candidate.json", got)
	}
	if got := requireString(t, attestation, "commercial_receipt_schema"); got != "commerce/commercial-distribution-receipt-schema-candidate.json" {
		t.Fatalf("commercial_receipt_schema = %s, want commerce/commercial-distribution-receipt-schema-candidate.json", got)
	}

	requiredEvidence := asStringSlice(t, attestation["required_gui_export_evidence"])
	for _, proof := range []string{
		"role_creation_screenshot_refs",
		"candidate_bundle_export_click_screenshot",
		"candidate_bundle_export_receipt_ref",
		"six_role_candidate_refs",
		"bundle_tree_sha256",
		"artifact_manifest_sha256",
		"no_manual_json_edit_evidence",
	} {
		requireStringSliceContains(t, requiredEvidence, proof)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	manifestFiles := map[string]string{}
	for _, item := range asObjectSlice(t, manifest["files"]) {
		manifestFiles[requireString(t, item, "path")] = requireString(t, item, "required_for")
	}

	roleRefs := map[string]bool{}
	for _, ref := range asStringSlice(t, candidateSet["role_pack_refs"]) {
		roleRefs[ref] = false
	}
	roleCandidates := asObjectSlice(t, attestation["required_role_candidate_refs"])
	if len(roleCandidates) != len(roleRefs) {
		t.Fatalf("required_role_candidate_refs len = %d, want %d", len(roleCandidates), len(roleRefs))
	}
	for _, role := range roleCandidates {
		roleRef := requireString(t, role, "role_pack_ref")
		if _, ok := roleRefs[roleRef]; !ok {
			t.Fatalf("unexpected role_pack_ref %s", roleRef)
		}
		roleRefs[roleRef] = true
		roleFile := requireString(t, role, "role_pack_file")
		if _, ok := manifestFiles[roleFile]; !ok {
			t.Fatalf("role_pack_file %s missing from artifact manifest", roleFile)
		}
		for _, key := range []string{"slot_ref", "candidate_ref", "gui_creation_step_ref"} {
			if got := requireString(t, role, key); got == "" {
				t.Fatalf("%s missing for %s", key, roleRef)
			}
		}
		if got := requireString(t, role, "forbidden_phrase_check_status"); got != "required_pass" {
			t.Fatalf("%s forbidden_phrase_check_status = %s, want required_pass", roleRef, got)
		}
	}
	for ref, seen := range roleRefs {
		if !seen {
			t.Fatalf("required_role_candidate_refs missing %s", ref)
		}
	}

	policy := requireObject(t, attestation, "upload_blocking_policy")
	requireBool(t, policy, "blocks_cloud_upload_without_gui_export_receipt", true)
	requireBool(t, policy, "blocks_cloud_upload_if_manual_json_edit_detected", true)
	requireBool(t, policy, "blocks_cloud_upload_if_bundle_hash_mismatch", true)
	requireBool(t, policy, "blocks_cloud_upload_if_any_role_candidate_missing", true)

	negative := requireObject(t, attestation, "negative_cases")
	for _, caseID := range []string{
		"manual_json_edit_as_export_success",
		"backend_api_call_without_gui",
		"bundle_hash_mismatch",
		"missing_six_role_candidates",
	} {
		caseDoc := requireObject(t, negative, caseID)
		if got := requireString(t, caseDoc, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", caseID, got)
		}
		if got := requireString(t, caseDoc, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", caseID)
		}
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "role_studio_export_provenance_attestation_verified")

	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, verifier, "source_role_studio_export_provenance_attestation"); got != attestationPath {
		t.Fatalf("source_role_studio_export_provenance_attestation = %s, want %s", got, attestationPath)
	}

	if got, ok := manifestFiles[attestationPath]; !ok {
		t.Fatalf("artifact manifest missing %s", attestationPath)
	} else if got != "audit" {
		t.Fatalf("%s required_for = %s, want audit", attestationPath, got)
	}
}

func TestTeamOfficeProductReadinessEvidenceMatrixBlocksPrematureCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	if got := requireString(t, matrix, "completion_status"); got != "not_verified_requires_cross_repo_execution" {
		t.Fatalf("completion_status = %s", got)
	}

	policy := requireObject(t, matrix, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	requiredProofs := strings.Join(asStringSlice(t, policy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{"gui_agent_screenshots", "cloud_receipts", "install_receipts", "team_binding_receipts", "negative_case_block_receipts"} {
		if !strings.Contains(requiredProofs, proof) {
			t.Fatalf("completion policy missing %s", proof)
		}
	}

	rawGates, ok := matrix["readiness_gates"].([]any)
	if !ok {
		t.Fatalf("readiness_gates missing")
	}
	requiredGates := map[string]bool{
		"frontend_role_studio_gui":            false,
		"frontend_team_settings_gui":          false,
		"backend_role_candidate_gate_receipt": false,
		"cloud_upload_purchase_download":      false,
		"local_install_and_team_binding":      false,
		"negative_cases":                      false,
	}
	for _, raw := range rawGates {
		gate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("readiness gate = %T", raw)
		}
		gateID := requireString(t, gate, "gate_id")
		if _, ok := requiredGates[gateID]; ok {
			requiredGates[gateID] = true
		}
		if got := requireString(t, gate, "status"); got != "pending_cross_repo_execution" {
			t.Fatalf("%s status = %s, want pending_cross_repo_execution", gateID, got)
		}
		truthSource := requireString(t, gate, "truth_source")
		requireStringIn(t, truthSource, "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo")
		evidence := strings.Join(asStringSlice(t, gate["required_evidence"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "screenshot") && !strings.Contains(evidence, "blocked") {
			t.Fatalf("%s required_evidence must include receipt, screenshot, or blocked proof", gateID)
		}
	}
	for gateID, seen := range requiredGates {
		if !seen {
			t.Fatalf("missing readiness gate %s", gateID)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/product-readiness-evidence-matrix.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/product-readiness-evidence-matrix.json")
}

func TestTeamOfficeNormalCommercializationCompletionAuditMapsObjectiveToEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	auditPath := "tests/normal-commercialization-completion-audit-candidate.json"
	audit := readJSON(t, filepath.Join(base, auditPath))
	requireBool(t, audit, "candidate_only", true)
	requireBool(t, audit, "non_formal", true)
	if got := requireString(t, audit, "completion_status"); got != "not_achieved_requires_cross_repo_execution" {
		t.Fatalf("completion_status = %s, want not_achieved_requires_cross_repo_execution", got)
	}
	if got := requireString(t, audit, "objective_source"); !strings.Contains(got, "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source = %s", got)
	}

	actors := strings.Join(asStringSlice(t, audit["required_participants"]), "\n")
	for _, actor := range []string{"user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent"} {
		if !strings.Contains(actors, actor) {
			t.Fatalf("required_participants missing %s", actor)
		}
	}

	rawRequirements, ok := audit["goal_requirements"].([]any)
	if !ok {
		t.Fatalf("goal_requirements missing")
	}
	required := map[string]string{
		"role_candidate_creation_gui":             "truzhen-client-web-desktop",
		"capability_pack_role_reference":          "multi_repo",
		"team_settings_secretary_replacement":     "multi_repo",
		"team_settings_five_advisors_replacement": "multi_repo",
		"secretary_voice_vrm_asset_ref":           "multi_repo",
		"team_office_runtime_use":                 "multi_repo",
		"cloud_upload_draft":                      "truzhen-cloud",
		"marketplace_review_listing":              "truzhen-cloud",
		"sandbox_payment_purchase":                "truzhen-cloud",
		"license_entitlement_issue":               "truzhen-cloud",
		"entitled_download_hash":                  "truzhen-cloud",
		"local_install_enabled_version":           "truzhenos",
		"post_install_team_settings_binding":      "multi_repo",
		"negative_cases_blocked":                  "multi_repo",
		"frontend_product_stage":                  "truzhen-client-web-desktop",
		"backend_product_stage":                   "truzhenos",
		"independent_acceptance_signoff":          "independent_acceptance_agent",
	}
	for _, raw := range rawRequirements {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("goal requirement = %T", raw)
		}
		reqID := requireString(t, req, "requirement_id")
		if wantTruth, ok := required[reqID]; ok {
			if got := requireString(t, req, "truth_source"); got != wantTruth {
				t.Fatalf("%s truth_source = %s, want %s", reqID, got, wantTruth)
			}
			required[reqID] = ""
		}
		requireStringIn(t, requireString(t, req, "status"), "candidate_evidence_defined_pending_cross_repo_execution", "pending_owner_decision")
		if len(asStringSlice(t, req["current_repo_evidence"])) == 0 {
			t.Fatalf("%s current_repo_evidence missing", reqID)
		}
		missing := strings.Join(asStringSlice(t, req["missing_authoritative_evidence"]), "\n")
		if !strings.Contains(missing, "receipt") && !strings.Contains(missing, "screenshot") && !strings.Contains(missing, "signoff") {
			t.Fatalf("%s missing_authoritative_evidence must mention receipt, screenshot, or signoff", reqID)
		}
		if got := requireString(t, req, "acceptance_gate_ref"); got == "" {
			t.Fatalf("%s acceptance_gate_ref missing", reqID)
		}
	}
	for reqID, missing := range required {
		if missing != "" {
			t.Fatalf("missing goal requirement %s", reqID)
		}
	}

	barriers := strings.Join(asStringSlice(t, audit["completion_barriers"]), "\n")
	for _, barrier := range []string{
		"cross_repo_authorization_missing",
		"gui_user_agent_evidence_missing",
		"truzhen_cloud_receipts_missing",
		"truzhenos_install_receipts_missing",
		"independent_acceptance_signoff_missing",
	} {
		if !strings.Contains(barriers, barrier) {
			t.Fatalf("completion_barriers missing %s", barrier)
		}
	}

	nonSufficient := strings.Join(asStringSlice(t, audit["non_sufficient_evidence"]), "\n")
	for _, item := range []string{"candidate_json_only", "current_repo_go_test_only", "mock_success_only", "organizer_self_attestation_only"} {
		if !strings.Contains(nonSufficient, item) {
			t.Fatalf("non_sufficient_evidence missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "normal_commercialization_completion_audit"); got != auditPath {
		t.Fatalf("normal_commercialization_completion_audit = %s, want %s", got, auditPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, auditPath) {
		t.Fatalf("candidate set artifact_files missing %s", auditPath)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policy := requireObject(t, matrix, "completion_claim_policy")
	requiredProofs := strings.Join(asStringSlice(t, policy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "normal_commercialization_completion_audit_verified") {
		t.Fatalf("product readiness matrix missing normal_commercialization_completion_audit_verified completion proof")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == auditPath {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", auditPath)
}

func TestTeamOfficeCompletionAuditConsumesRoleStudioExportProvenance(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	attestationPath := "tests/role-studio-export-provenance-attestation-candidate.json"
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "role_studio_export_provenance_attestation_ref"); got != attestationPath {
		t.Fatalf("role_studio_export_provenance_attestation_ref = %s, want %s", got, attestationPath)
	}

	rawRequirements, ok := audit["goal_requirements"].([]any)
	if !ok {
		t.Fatalf("goal_requirements missing")
	}
	for _, raw := range rawRequirements {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("goal requirement = %T", raw)
		}
		if requireString(t, req, "requirement_id") != "cloud_upload_draft" {
			continue
		}
		requireStringSliceContains(t, asStringSlice(t, req["current_repo_evidence"]), attestationPath)
		missing := strings.Join(asStringSlice(t, req["missing_authoritative_evidence"]), "\n")
		if !strings.Contains(missing, "candidate bundle export receipt") {
			t.Fatalf("cloud_upload_draft missing_authoritative_evidence must mention candidate bundle export receipt, got %s", missing)
		}
		return
	}
	t.Fatalf("cloud_upload_draft requirement missing")
}

func TestTeamOfficeGoalCompletionEvidenceMapBlocksGoalCompletionUntilAllAuthoritativeEvidenceExists(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	mapPath := "tests/role-studio-goal-completion-evidence-map-candidate.json"
	evidenceMap := readJSON(t, filepath.Join(base, mapPath))
	requireBool(t, evidenceMap, "candidate_only", true)
	requireBool(t, evidenceMap, "non_formal", true)
	requireBool(t, evidenceMap, "can_mark_goal_complete", false)
	if got := requireString(t, evidenceMap, "completion_status"); got != "not_achieved_requires_cross_repo_execution" {
		t.Fatalf("completion_status = %s, want not_achieved_requires_cross_repo_execution", got)
	}
	if got := requireString(t, evidenceMap, "source_plan"); !strings.Contains(got, "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("source_plan = %s", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref": "role-pack-candidate-set://team-office-v0",
		"source_normal_commercialization_completion_audit": "tests/normal-commercialization-completion-audit-candidate.json",
		"source_role_studio_test_case_coverage_matrix":     "tests/role-studio-test-case-coverage-matrix-candidate.json",
		"source_role_studio_export_provenance_attestation": "tests/role-studio-export-provenance-attestation-candidate.json",
		"source_commercial_readiness_verifier":             "tests/commercial-readiness-verifier-candidate.json",
		"source_commercial_go_no_go_gate":                  "tests/commercial-go-no-go-gate-candidate.json",
		"source_commercial_production_promotion_gate":      "commerce/commercial-production-promotion-gate-candidate.json",
		"source_product_readiness_evidence_matrix":         "tests/product-readiness-evidence-matrix.json",
		"source_p11_commercial_go_live_evidence_package":   "tests/p11-commercial-go-live-evidence-package-template.json",
		"next_required_authorization":                      "owner_authorization_intake",
		"next_authorization_card":                          "integration/owner-authorization-evidence-intake-candidate.json",
	} {
		if got := requireString(t, evidenceMap, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawRequirements, ok := evidenceMap["active_goal_requirements"].([]any)
	if !ok {
		t.Fatalf("active_goal_requirements missing")
	}
	expectedRequirements := map[string]string{
		"role_creation_gui":                       "角色制作",
		"role_runtime_use":                        "使用",
		"cloud_upload":                            "上传",
		"sandbox_or_authorized_purchase":          "支付购买",
		"entitled_download":                       "下载",
		"local_install_enabled":                   "安装",
		"team_settings_replacement_after_install": "团队设置替换",
		"normal_commercialization_go_no_go":       "正常商品化",
		"production_promotion_gate":               "生产发布晋级门",
	}
	if len(rawRequirements) != len(expectedRequirements) {
		t.Fatalf("active_goal_requirements len = %d, want %d", len(rawRequirements), len(expectedRequirements))
	}
	seen := map[string]bool{}
	for _, raw := range rawRequirements {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("active goal requirement = %T", raw)
		}
		reqID := requireString(t, req, "requirement_id")
		goalKeyword, ok := expectedRequirements[reqID]
		if !ok {
			t.Fatalf("unexpected requirement_id %s", reqID)
		}
		seen[reqID] = true
		if got := requireString(t, req, "goal_keyword"); !strings.Contains(got, goalKeyword) {
			t.Fatalf("%s goal_keyword = %s, want contains %s", reqID, got, goalKeyword)
		}
		if got := requireString(t, req, "current_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("%s current_status = %s, want missing_authoritative_evidence", reqID, got)
		}
		requireBool(t, req, "can_count_toward_goal_completion", false)
		for _, key := range []string{
			"truth_sources",
			"required_authoritative_evidence",
			"current_repo_candidate_evidence",
			"missing_authoritative_evidence",
			"blocking_reasons",
		} {
			if values := asStringSlice(t, req[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", reqID, key)
			}
		}
		evidence := strings.Join(asStringSlice(t, req["required_authoritative_evidence"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "screenshot") && !strings.Contains(evidence, "signoff") {
			t.Fatalf("%s required_authoritative_evidence must include receipt, screenshot, or signoff", reqID)
		}
	}
	for reqID := range expectedRequirements {
		if !seen[reqID] {
			t.Fatalf("missing active goal requirement %s", reqID)
		}
	}

	barriers := strings.Join(asStringSlice(t, evidenceMap["goal_completion_barriers"]), "\n")
	for _, barrier := range []string{
		"cross_repo_authorization_missing",
		"user_view_gui_evidence_missing",
		"truzhen_cloud_receipts_missing",
		"truzhenos_install_and_binding_receipts_missing",
		"independent_acceptance_and_owner_go_no_go_missing",
		"commercial_production_promotion_gate_missing",
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
		"role_studio_export_provenance_attestation_verified",
		"role_studio_goal_completion_evidence_map_verified",
		"commercial_readiness_verifier_passed",
		"commercial_go_no_go_gate_passed",
		"commercial_production_promotion_gate_verified",
		"owner_go_no_go_decision_recorded",
	} {
		if !strings.Contains(requiredBeforeCompletion, proof) {
			t.Fatalf("completion_claim_policy required_before_goal_complete missing %s", proof)
		}
	}
	nonSufficient := strings.Join(asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "\n")
	for _, item := range []string{"candidate_assets_only", "single_repo_go_test_only", "manifest_digest_only", "organizer_summary_without_receipts"} {
		if !strings.Contains(nonSufficient, item) {
			t.Fatalf("non_sufficient_evidence missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "role_studio_goal_completion_evidence_map"); got != mapPath {
		t.Fatalf("role_studio_goal_completion_evidence_map = %s, want %s", got, mapPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, mapPath) {
		t.Fatalf("candidate set artifact_files missing %s", mapPath)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_goal_completion_evidence_map"); got != mapPath {
		t.Fatalf("commercial readiness source_goal_completion_evidence_map = %s, want %s", got, mapPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "role_studio_goal_completion_evidence_map_verified") {
		t.Fatalf("product readiness matrix missing role_studio_goal_completion_evidence_map_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	mapInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != mapPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, mapPath))
		if err != nil {
			t.Fatalf("read %s: %v", mapPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", mapPath, got, wantHash)
		}
		mapInManifest = true
	}
	if !mapInManifest {
		t.Fatalf("artifact manifest missing %s", mapPath)
	}
}

func TestTeamOfficeGoalCompletionEvidenceMapRequiresGoLiveApprovalAndP11Acceptance(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))

	approvalPath := "commerce/commercial-go-live-approval-candidate.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	if got := requireString(t, evidenceMap, "source_commercial_go_live_approval"); got != approvalPath {
		t.Fatalf("source_commercial_go_live_approval = %s, want %s", got, approvalPath)
	}
	if got := requireString(t, evidenceMap, "source_p11_evidence_acceptance_checklist"); got != checklistPath {
		t.Fatalf("source_p11_evidence_acceptance_checklist = %s, want %s", got, checklistPath)
	}

	requirement := findObjectByString(t, asObjectSlice(t, evidenceMap["active_goal_requirements"]), "requirement_id", "normal_commercialization_go_no_go")
	for _, evidence := range []string{
		"commercial_go_live_approval_verified",
		"p11_evidence_acceptance_checklist_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, requirement["required_authoritative_evidence"]), evidence)
	}
	for _, candidateEvidence := range []string{approvalPath, checklistPath} {
		requireStringSliceContains(t, asStringSlice(t, requirement["current_repo_candidate_evidence"]), candidateEvidence)
	}
	for _, missing := range []string{
		"commercial_go_live_approval_verified_missing",
		"p11_evidence_acceptance_checklist_verified_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, requirement["missing_authoritative_evidence"]), missing)
	}

	for _, barrier := range []string{
		"commercial_go_live_approval_missing",
		"p11_evidence_acceptance_checklist_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), barrier)
	}
	policy := requireObject(t, evidenceMap, "completion_claim_policy")
	for _, proof := range []string{
		"commercial_go_live_approval_verified",
		"p11_evidence_acceptance_checklist_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, policy["required_before_goal_complete"]), proof)
	}
	for _, item := range []string{
		"go_live_approval_without_p11_acceptance_checklist",
		"commercial_go_no_go_without_go_live_approval",
	} {
		requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), item)
	}
}

func TestTeamOfficeCommercialReadinessGoNoGoAndProductMatrixRequireGoLiveApproval(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	approvalPath := "commerce/commercial-go-live-approval-candidate.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_commercial_go_live_approval"); got != approvalPath {
		t.Fatalf("readiness source_commercial_go_live_approval = %s, want %s", got, approvalPath)
	}
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "commercial_go_live_approval_verified")
	if got := requireString(t, readinessGate, "current_result"); got != "pending" {
		t.Fatalf("readiness commercial_go_live_approval_verified current_result = %s, want pending", got)
	}
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, want := range []string{approvalPath, checklistPath, packagePath, "owner_go_no_go_decision_recorded", "p11_evidence_acceptance_checklist_verified"} {
		if !strings.Contains(readinessEvidence, want) {
			t.Fatalf("readiness commercial_go_live_approval_verified evidence_required missing %s: %s", want, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "commercial_go_live_approval_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_commercial_go_live_approval")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_commercial_go_live_approval"); got != approvalPath {
		t.Fatalf("go/no-go source_commercial_go_live_approval = %s, want %s", got, approvalPath)
	}
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "commercial_go_live_approval_verified")
	if got := requireString(t, goNoGoGate, "current_result"); got != "pending" {
		t.Fatalf("go/no-go commercial_go_live_approval_verified current_result = %s, want pending", got)
	}
	requireBool(t, goNoGoGate, "required_final_value", true)
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	goNoGoEvidence := requireString(t, goNoGoGate, "evidence_required")
	for _, want := range []string{approvalPath, checklistPath, packagePath, "commercial_go_live_approval_verified", "owner_go_no_go_decision_recorded"} {
		if !strings.Contains(goNoGoEvidence, want) {
			t.Fatalf("go/no-go commercial_go_live_approval_verified evidence_required missing %s: %s", want, goNoGoEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "commercial_go_live_approval_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_commercial_go_live_approval")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "commercial_go_live_approval_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "commercial_go_live_approval_missing")
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "commercial_go_live_approval_verified")
	if got := requireString(t, productGate, "truth_source"); got != "multi_repo" {
		t.Fatalf("product matrix commercial_go_live_approval_verified truth_source = %s, want multi_repo", got)
	}
	if got := requireString(t, productGate, "status"); got != "pending_cross_repo_execution" {
		t.Fatalf("product matrix commercial_go_live_approval_verified status = %s, want pending_cross_repo_execution", got)
	}
	for _, evidence := range []string{approvalPath, "commercial_go_no_go_gate_passed", "p11_evidence_acceptance_checklist_verified", "owner_go_no_go_decision_recorded"} {
		requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), evidence)
	}
	if got := requireString(t, productGate, "blocking_if_missing"); got != "cannot_claim_go_live_approval_ready" {
		t.Fatalf("product matrix commercial_go_live_approval_verified blocking_if_missing = %s", got)
	}
}

func TestTeamOfficeCommercialReadinessGoNoGoProductAndGoalRequireCloudUploadListingVerified(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	listingPath := "commerce/cloud-listing-candidate.json"
	manifestPath := "commerce/artifact-manifest.json"
	digestPath := "commerce/artifact-bundle-digest-candidate.json"
	releasePath := "commerce/release-candidate-package.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":   listingPath,
		"source_artifact_manifest":         manifestPath,
		"source_artifact_bundle_digest":    digestPath,
		"source_release_candidate_package": releasePath,
	} {
		if got := requireString(t, readiness, key); got != want {
			t.Fatalf("readiness %s = %s, want %s", key, got, want)
		}
	}
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "cloud_upload_listing_verified")
	if got := requireString(t, readinessGate, "current_result"); got != "pending" {
		t.Fatalf("readiness cloud_upload_listing_verified current_result = %s, want pending", got)
	}
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, want := range []string{packagePath + "#cloud_upload_listing_report", checklistPath, recordPath + "#cloud_upload_listing_result", listingPath, manifestPath, digestPath, releasePath, "cloud_upload_receipt_ref", "listing_draft_ref"} {
		if !strings.Contains(readinessEvidence, want) {
			t.Fatalf("readiness cloud_upload_listing_verified evidence_required missing %s: %s", want, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_cloud_upload_listing_report")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":   listingPath,
		"source_artifact_manifest":         manifestPath,
		"source_artifact_bundle_digest":    digestPath,
		"source_release_candidate_package": releasePath,
	} {
		if got := requireString(t, goNoGo, key); got != want {
			t.Fatalf("go/no-go %s = %s, want %s", key, got, want)
		}
	}
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "cloud_upload_listing_verified")
	requireBool(t, goNoGoGate, "required_final_value", true)
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	goNoGoEvidence := requireString(t, goNoGoGate, "evidence_required")
	for _, want := range []string{packagePath + "#cloud_upload_listing_report", checklistPath, recordPath + "#cloud_upload_listing_result", listingPath, "cloud_upload_receipt_ref", "listing_draft_ref"} {
		if !strings.Contains(goNoGoEvidence, want) {
			t.Fatalf("go/no-go cloud_upload_listing_verified evidence_required missing %s: %s", want, goNoGoEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_cloud_upload_listing_report")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "cloud_upload_listing_missing")
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "cloud_upload_listing_verified")
	if got := requireString(t, productGate, "truth_source"); got != "truzhen-cloud" {
		t.Fatalf("product matrix cloud_upload_listing_verified truth_source = %s, want truzhen-cloud", got)
	}
	for _, evidence := range []string{"cloud_upload_listing_report", "cloud_upload_receipt_ref", "listing_draft_ref", "artifact_manifest_ref", "bundle_tree_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), evidence)
	}

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	uploadReq := findObjectByString(t, asObjectSlice(t, evidenceMap["active_goal_requirements"]), "requirement_id", "cloud_upload")
	requireStringSliceContains(t, asStringSlice(t, uploadReq["required_authoritative_evidence"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, uploadReq["current_repo_candidate_evidence"]), packagePath+"#cloud_upload_listing_report")
	requireStringSliceContains(t, asStringSlice(t, uploadReq["missing_authoritative_evidence"]), "cloud_upload_listing_verified_missing")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "cloud_upload_listing_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_cloud_upload_listing_verified")
}

func TestTeamOfficeCommercialReadinessGoNoGoProductAndGoalRequireDownloadArtifactDeliveryVerified(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	buyerLibraryPath := "commerce/buyer-library-install-state-candidate.json"
	entitlementPolicyPath := "commerce/license-entitlement-policy-candidate.json"
	digestPath := "commerce/artifact-bundle-digest-candidate.json"
	releasePath := "commerce/release-candidate-package.json"
	receiptChainPath := "commerce/commercial-receipt-chain-candidate.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	for key, want := range map[string]string{
		"source_buyer_library_install_state": buyerLibraryPath,
		"source_license_entitlement_policy":  entitlementPolicyPath,
		"source_artifact_bundle_digest":      digestPath,
		"source_release_candidate_package":   releasePath,
	} {
		if got := requireString(t, readiness, key); got != want {
			t.Fatalf("readiness %s = %s, want %s", key, got, want)
		}
	}
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "download_artifact_delivery_verified")
	if got := requireString(t, readinessGate, "current_result"); got != "pending" {
		t.Fatalf("readiness download_artifact_delivery_verified current_result = %s, want pending", got)
	}
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, want := range []string{
		packagePath + "#download_artifact_delivery_report",
		recordPath + "#download_artifact_delivery_result",
		checklistPath,
		buyerLibraryPath,
		entitlementPolicyPath,
		digestPath,
		releasePath,
		receiptChainPath,
		"download_receipt_ref",
		"signed_download_ref",
		"download_artifact_sha256",
	} {
		if !strings.Contains(readinessEvidence, want) {
			t.Fatalf("readiness download_artifact_delivery_verified evidence_required missing %s: %s", want, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_download_artifact_delivery_report")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	for key, want := range map[string]string{
		"source_buyer_library_install_state": buyerLibraryPath,
		"source_license_entitlement_policy":  entitlementPolicyPath,
		"source_artifact_bundle_digest":      digestPath,
		"source_release_candidate_package":   releasePath,
	} {
		if got := requireString(t, goNoGo, key); got != want {
			t.Fatalf("go/no-go %s = %s, want %s", key, got, want)
		}
	}
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "download_artifact_delivery_verified")
	requireBool(t, goNoGoGate, "required_final_value", true)
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	goNoGoEvidence := requireString(t, goNoGoGate, "evidence_required")
	for _, want := range []string{
		packagePath + "#download_artifact_delivery_report",
		recordPath + "#download_artifact_delivery_result",
		checklistPath,
		"download_receipt_ref",
		"signed_download_ref",
		"download_artifact_sha256",
	} {
		if !strings.Contains(goNoGoEvidence, want) {
			t.Fatalf("go/no-go download_artifact_delivery_verified evidence_required missing %s: %s", want, goNoGoEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_download_artifact_delivery_report")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "download_artifact_delivery_missing")
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "download_artifact_delivery_verified")
	if got := requireString(t, productGate, "truth_source"); got != "truzhen-cloud" {
		t.Fatalf("product matrix download_artifact_delivery_verified truth_source = %s, want truzhen-cloud", got)
	}
	for _, evidence := range []string{"download_artifact_delivery_report", "entitlement_ref", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256", "bundle_tree_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), evidence)
	}

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	downloadReq := findObjectByString(t, asObjectSlice(t, evidenceMap["active_goal_requirements"]), "requirement_id", "entitled_download")
	requireStringSliceContains(t, asStringSlice(t, downloadReq["required_authoritative_evidence"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, downloadReq["current_repo_candidate_evidence"]), packagePath+"#download_artifact_delivery_report")
	requireStringSliceContains(t, asStringSlice(t, downloadReq["missing_authoritative_evidence"]), "download_artifact_delivery_verified_missing")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "download_artifact_delivery_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_download_artifact_delivery_verified")
}

func TestTeamOfficeCompletionAuditRequiresProductionPromotionGateForNormalCommercialization(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gatePath := "commerce/commercial-production-promotion-gate-candidate.json"
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "source_commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("normal commercialization audit source_commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}

	rawRequirements, ok := audit["goal_requirements"].([]any)
	if !ok {
		t.Fatalf("goal_requirements missing")
	}
	foundRequirement := false
	for _, raw := range rawRequirements {
		req, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("goal requirement = %T", raw)
		}
		if requireString(t, req, "requirement_id") != "production_promotion_gate" {
			continue
		}
		foundRequirement = true
		if got := requireString(t, req, "truth_source"); got != "multi_repo" {
			t.Fatalf("production_promotion_gate truth_source = %s, want multi_repo", got)
		}
		requireStringIn(t, requireString(t, req, "status"), "pending_owner_decision")
		requireBool(t, req, "can_count_toward_goal_completion", false)
		missing := strings.Join(asStringSlice(t, req["missing_authoritative_evidence"]), "\n")
		for _, proof := range []string{"production promotion gate pass receipt", "Owner go/no-go signoff"} {
			if !strings.Contains(missing, proof) {
				t.Fatalf("production_promotion_gate missing_authoritative_evidence missing %s", proof)
			}
		}
		if got := requireString(t, req, "acceptance_gate_ref"); got != gatePath {
			t.Fatalf("production_promotion_gate acceptance_gate_ref = %s, want %s", got, gatePath)
		}
	}
	if !foundRequirement {
		t.Fatalf("normal commercialization audit missing production_promotion_gate requirement")
	}

	barriers := strings.Join(asStringSlice(t, audit["completion_barriers"]), "\n")
	if !strings.Contains(barriers, "commercial_production_promotion_gate_missing") {
		t.Fatalf("completion_barriers missing commercial_production_promotion_gate_missing")
	}
	nonSufficient := strings.Join(asStringSlice(t, audit["non_sufficient_evidence"]), "\n")
	if !strings.Contains(nonSufficient, "sandbox_pass_without_production_promotion_gate") {
		t.Fatalf("non_sufficient_evidence missing sandbox_pass_without_production_promotion_gate")
	}

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	if got := requireString(t, evidenceMap, "source_commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("goal completion map source_commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}
	policy := requireObject(t, evidenceMap, "completion_claim_policy")
	requiredBeforeGoalComplete := strings.Join(asStringSlice(t, policy["required_before_goal_complete"]), "\n")
	if !strings.Contains(requiredBeforeGoalComplete, "commercial_production_promotion_gate_verified") {
		t.Fatalf("goal completion map required_before_goal_complete missing commercial_production_promotion_gate_verified")
	}
}

func TestTeamOfficeP11EvidenceAcceptanceChecklistDefinesExecutableEvidenceGates(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	checklist := readJSON(t, filepath.Join(base, checklistPath))
	requireBool(t, checklist, "candidate_only", true)
	requireBool(t, checklist, "non_formal", true)
	requireBool(t, checklist, "can_mark_p11_pass", false)
	if got := requireString(t, checklist, "checklist_status"); got != "not_run_requires_authoritative_evidence" {
		t.Fatalf("checklist_status = %s, want not_run_requires_authoritative_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                   "role-pack-candidate-set://team-office-v0",
		"source_runbook":                      "tests/p11-sandbox-execution-runbook-candidate.json",
		"source_verification_record":          "tests/p11-normal-commercialization-verification-record-template.json",
		"source_evidence_ingestion_binder":    "tests/p11-evidence-ingestion-binder-candidate.json",
		"source_go_live_evidence_package":     "tests/p11-commercial-go-live-evidence-package-template.json",
		"source_commercial_receipt_schema":    "commerce/commercial-distribution-receipt-schema-candidate.json",
		"source_commercial_receipt_chain":     "commerce/commercial-receipt-chain-candidate.json",
		"source_product_readiness_matrix":     "tests/product-readiness-evidence-matrix.json",
		"source_goal_completion_evidence_map": "tests/role-studio-goal-completion-evidence-map-candidate.json",
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawStages, ok := checklist["stage_acceptance_checks"].([]any)
	if !ok {
		t.Fatalf("stage_acceptance_checks missing")
	}
	expectedStages := []string{
		"preflight_authorization_and_status",
		"role_candidate_bundle_export",
		"cloud_upload_listing_draft",
		"cloud_upload_listing_report",
		"marketplace_review_candidate",
		"marketplace_listing_review_compliance_report",
		"sandbox_order_payment_entitlement",
		"purchase_entitlement_report",
		"entitled_signed_download",
		"download_artifact_delivery_report",
		"local_install_enabled_version",
		"download_install_access_matrix_report",
		"post_install_team_binding_runtime",
		"role_studio_lineage_report",
		"gui_api_traceability_matrix",
		"negative_cases_and_independent_acceptance",
	}
	if len(rawStages) != len(expectedStages) {
		t.Fatalf("stage_acceptance_checks len = %d, want %d", len(rawStages), len(expectedStages))
	}
	for i, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage item = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if stageID != expectedStages[i] {
			t.Fatalf("stage[%d] = %s, want %s", i, stageID, expectedStages[i])
		}
		if got := requireString(t, stage, "current_status"); got != "not_run" {
			t.Fatalf("%s current_status = %s, want not_run", stageID, got)
		}
		requireBool(t, stage, "can_count_toward_p11_pass", false)
		for _, key := range []string{
			"required_gui_evidence",
			"required_receipt_evidence",
			"required_hash_or_correlation_checks",
			"blocking_if_missing",
			"writeback_targets",
		} {
			if values := asStringSlice(t, stage[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
		evidence := strings.Join(asStringSlice(t, stage["required_receipt_evidence"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "signoff") && !strings.Contains(evidence, "owner_go_no_go") {
			t.Fatalf("%s required_receipt_evidence must include receipt, signoff, or owner go/no-go proof", stageID)
		}
		correlation := strings.Join(asStringSlice(t, stage["required_hash_or_correlation_checks"]), "\n")
		if !strings.Contains(correlation, "bundle_tree_sha256") && !strings.Contains(correlation, "candidate_set_ref") {
			t.Fatalf("%s required_hash_or_correlation_checks must include bundle_tree_sha256 or candidate_set_ref", stageID)
		}
	}

	negativeCases := strings.Join(asStringSlice(t, checklist["required_negative_cases"]), "\n")
	for _, caseID := range []string{
		"download_without_entitlement",
		"artifact_hash_mismatch",
		"install_without_entitlement",
		"real_payment_without_owner_authorization",
		"raw_voice_or_vrm_asset",
		"team_binding_without_owner_gate",
	} {
		if !strings.Contains(negativeCases, caseID) {
			t.Fatalf("required_negative_cases missing %s", caseID)
		}
	}
	forbiddenPayloads := strings.Join(asStringSlice(t, checklist["forbidden_payloads"]), "\n")
	for _, payload := range []string{"raw_payment_token", "signed_download_url_secret", "cloud_access_token", "raw_voice_asset", "raw_vrm_asset"} {
		if !strings.Contains(forbiddenPayloads, payload) {
			t.Fatalf("forbidden_payloads missing %s", payload)
		}
	}

	policy := requireObject(t, checklist, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	requiredBeforeP11Pass := strings.Join(asStringSlice(t, policy["required_before_p11_pass"]), "\n")
	for _, proof := range []string{
		"all_stage_acceptance_checks_verified",
		"all_required_negative_cases_blocked",
		"hash_continuity_verified",
		"independent_acceptance_signed",
		"owner_go_no_go_decision_recorded",
	} {
		if !strings.Contains(requiredBeforeP11Pass, proof) {
			t.Fatalf("completion_claim_policy missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_evidence_acceptance_checklist"); got != checklistPath {
		t.Fatalf("p11_evidence_acceptance_checklist = %s, want %s", got, checklistPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, checklistPath) {
		t.Fatalf("candidate set artifact_files missing %s", checklistPath)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_p11_evidence_acceptance_checklist"); got != checklistPath {
		t.Fatalf("commercial readiness source_p11_evidence_acceptance_checklist = %s, want %s", got, checklistPath)
	}
	evidencePackage := readJSON(t, filepath.Join(base, "tests", "p11-commercial-go-live-evidence-package-template.json"))
	if got := requireString(t, evidencePackage, "evidence_acceptance_checklist_ref"); got != checklistPath {
		t.Fatalf("evidence package evidence_acceptance_checklist_ref = %s, want %s", got, checklistPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "p11_evidence_acceptance_checklist_verified") {
		t.Fatalf("product readiness matrix missing p11_evidence_acceptance_checklist_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != checklistPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, checklistPath))
		if err != nil {
			t.Fatalf("read %s: %v", checklistPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", checklistPath, got, wantHash)
		}
		return
	}
	t.Fatalf("artifact manifest missing %s", checklistPath)
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresCloudUploadListingWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	listingPath := "commerce/cloud-listing-candidate.json"
	manifestPath := "commerce/artifact-manifest.json"
	digestPath := "commerce/artifact-bundle-digest-candidate.json"
	releasePath := "commerce/release-candidate-package.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":   listingPath,
		"source_artifact_manifest":         manifestPath,
		"source_artifact_bundle_digest":    digestPath,
		"source_release_candidate_package": releasePath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "cloud_upload_listing_report")
	for _, guiEvidence := range []string{"cloud_upload_draft_screen", "listing_draft_state_screen", "artifact_hash_match_review_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"candidate_bundle_export_receipt", "cloud_upload_receipt", "listing_draft_receipt", "artifact_hash_match_receipt", "forbidden_artifact_scan_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "artifact_manifest_ref", "listing_draft_ref", "cloud_upload_receipt_ref", "uploaded_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#cloud_upload_listing_result", packagePath + "#cloud_upload_listing_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "cloud_upload_listing_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "cloud_upload_listing_verification_record_target"); got != "cloud_upload_listing_result" {
		t.Fatalf("cloud_upload_listing_verification_record_target = %s, want cloud_upload_listing_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "cloud_upload_listing_records")
	uploadIngestion := requireObject(t, binder, "cloud_upload_listing_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "cloud_upload_listing_result",
		"package_section":     "cloud_upload_listing_report",
	} {
		if got := requireString(t, uploadIngestion, key); got != want {
			t.Fatalf("cloud_upload_listing_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, sourceRef := range []string{listingPath, manifestPath, digestPath, releasePath} {
		requireStringSliceContains(t, asStringSlice(t, uploadIngestion["source_refs"]), sourceRef)
	}
	for _, slot := range []string{"candidate_set_ref", "bundle_tree_sha256", "artifact_manifest_ref", "cloud_upload_receipt_ref", "listing_draft_ref", "listing_draft_receipt_ref", "uploaded_artifact_sha256", "artifact_hash_match_receipt_ref", "forbidden_artifact_scan_receipt_ref"} {
		requireStringSliceContains(t, asStringSlice(t, uploadIngestion["required_evidence_slots"]), slot)
	}
	for _, correlation := range []string{"artifact_manifest_ref", "cloud_upload_receipt_ref", "listing_draft_ref"} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), correlation)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "cloud_upload_listing_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "cloud_upload_listing_result")
	result := requireObject(t, record, "cloud_upload_listing_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("cloud_upload_listing_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"candidate_set_ref", "bundle_tree_sha256", "artifact_manifest_ref", "cloud_upload_receipt_ref", "listing_draft_ref", "listing_draft_receipt", "uploaded_artifact_sha256", "artifact_hash_match_receipt", "forbidden_artifact_scan_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"cloud_upload_receipt_bundle_tree_sha256_matches_manifest", "listing_draft_ref_matches_cloud_upload_receipt", "uploaded_artifact_sha256_matches_artifact_manifest", "forbidden_artifact_scan_passed_before_cloud_upload"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "cloud_upload_listing_verified")

	pkg := readJSON(t, filepath.Join(base, packagePath))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":   listingPath,
		"source_artifact_manifest":         manifestPath,
		"source_artifact_bundle_digest":    digestPath,
		"source_release_candidate_package": releasePath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("package %s = %s, want %s", key, got, want)
		}
	}
	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "cloud_upload_listing_report")
	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "cloud_upload_listing_report")
	for _, sourceRef := range []string{listingPath, manifestPath, digestPath, releasePath, "commerce/commercial-receipt-chain-candidate.json"} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{"cloud_upload_listing_verified", "candidate_set_ref", "bundle_tree_sha256", "artifact_manifest_ref", "cloud_upload_receipt", "cloud_upload_receipt_ref", "listing_draft_ref", "listing_draft_receipt", "uploaded_artifact_sha256", "artifact_hash_match_receipt", "forbidden_artifact_scan_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{"cloud_upload_listing_report_missing", "cloud_upload_receipt_missing", "listing_draft_receipt_missing", "artifact_hash_match_missing", "forbidden_artifact_scan_receipt_missing"} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "cloud_upload_listing_verified")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["non_sufficient_evidence"]), "cloud_upload_listing_report_missing")
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresMarketplaceListingReviewComplianceWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	listingPath := "commerce/cloud-listing-candidate.json"
	reviewPath := "commerce/marketplace-review-submission-candidate.json"
	supportRefundPath := "commerce/support-refund-revocation-policy-candidate.json"
	termsPrivacyPath := "commerce/commercial-terms-privacy-policy-candidate.json"
	publisherPath := "commerce/publisher-account-settlement-policy-candidate.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":             listingPath,
		"source_marketplace_review_submission":       reviewPath,
		"source_support_refund_revocation_policy":    supportRefundPath,
		"source_commercial_terms_privacy_policy":     termsPrivacyPath,
		"source_publisher_account_settlement_policy": publisherPath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "marketplace_listing_review_compliance_report")
	for _, guiEvidence := range []string{"cloud_listing_draft_review_screen", "marketplace_review_submission_screen", "pre_purchase_disclosure_review_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"listing_draft_receipt", "marketplace_review_candidate_receipt", "pre_purchase_disclosure_ack_receipt", "terms_privacy_acceptance_receipt", "publisher_identity_verification_receipt", "pricing_approval_receipt", "production_publish_block_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "listing_draft_ref", "marketplace_review_candidate_ref"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#marketplace_listing_review_result", packagePath + "#marketplace_listing_review_compliance_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "marketplace_listing_review_compliance_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "marketplace_listing_review_verification_record_target"); got != "marketplace_listing_review_result" {
		t.Fatalf("marketplace_listing_review_verification_record_target = %s, want marketplace_listing_review_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "marketplace_listing_review_records")
	marketplaceIngestion := requireObject(t, binder, "marketplace_listing_review_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "marketplace_listing_review_result",
		"package_section":     "marketplace_listing_review_compliance_report",
	} {
		if got := requireString(t, marketplaceIngestion, key); got != want {
			t.Fatalf("marketplace_listing_review_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, sourceRef := range []string{listingPath, reviewPath, supportRefundPath, termsPrivacyPath, publisherPath} {
		requireStringSliceContains(t, asStringSlice(t, marketplaceIngestion["source_refs"]), sourceRef)
	}
	for _, slot := range []string{"listing_draft_ref", "listing_draft_receipt_ref", "marketplace_review_candidate_ref", "marketplace_review_candidate_receipt_ref", "pre_purchase_disclosure_ack_receipt_ref", "terms_privacy_acceptance_receipt_ref", "publisher_identity_verification_receipt_ref", "pricing_approval_receipt_ref", "production_publish_block_receipt_ref", "asset_rights_review_receipt_ref"} {
		requireStringSliceContains(t, asStringSlice(t, marketplaceIngestion["required_evidence_slots"]), slot)
	}
	for _, correlation := range []string{"listing_draft_ref", "marketplace_review_candidate_ref"} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), correlation)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "marketplace_listing_review_compliance_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "marketplace_listing_review_result")
	result := requireObject(t, record, "marketplace_listing_review_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("marketplace_listing_review_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"listing_draft_ref", "listing_draft_receipt", "marketplace_review_candidate_ref", "marketplace_review_candidate_receipt", "pre_purchase_disclosure_ack_receipt", "terms_privacy_acceptance_receipt", "publisher_identity_verification_receipt", "pricing_approval_receipt", "production_publish_block_receipt", "asset_rights_review_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"listing_component_list_matches_candidate_bundle", "marketplace_review_state_allows_sandbox_purchase_only", "pre_purchase_disclosure_ack_present", "terms_privacy_acceptance_present", "publisher_identity_and_pricing_approved", "production_publish_block_receipt_present"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "marketplace_listing_review_compliance_verified")
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresPurchaseEntitlementWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	orderPaymentPath := "commerce/order-payment-state-machine-candidate.json"
	entitlementPolicyPath := "commerce/license-entitlement-policy-candidate.json"
	buyerLibraryPath := "commerce/buyer-library-install-state-candidate.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_order_payment_state_machine": orderPaymentPath,
		"source_license_entitlement_policy":  entitlementPolicyPath,
		"source_buyer_library_install_state": buyerLibraryPath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "purchase_entitlement_report")
	for _, guiEvidence := range []string{"sandbox_purchase_screen", "entitlement_library_screen", "buyer_library_purchased_available_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"sandbox_order_receipt", "sandbox_payment_receipt", "no_real_payment_capture_receipt", "real_payment_block_policy_receipt", "entitlement_receipt", "entitlement_team_scope_verification_receipt_ref", "payment_failed_block_receipt", "refund_or_chargeback_revocation_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "sandbox_order_ref", "license_ref", "entitlement_ref"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#purchase_entitlement_result", packagePath + "#purchase_entitlement_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "purchase_entitlement_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "purchase_entitlement_verification_record_target"); got != "purchase_entitlement_result" {
		t.Fatalf("purchase_entitlement_verification_record_target = %s, want purchase_entitlement_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "purchase_entitlement_records")
	purchaseIngestion := requireObject(t, binder, "purchase_entitlement_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "purchase_entitlement_result",
		"package_section":     "purchase_entitlement_report",
	} {
		if got := requireString(t, purchaseIngestion, key); got != want {
			t.Fatalf("purchase_entitlement_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, slot := range []string{"sandbox_order_ref", "sandbox_order_receipt_ref", "sandbox_payment_receipt_ref", "no_real_payment_capture_receipt", "real_payment_block_policy_receipt", "license_ref", "entitlement_ref", "entitlement_team_scope_verification_receipt_ref", "buyer_library_state_ref", "payment_failed_block_receipt_ref", "refund_or_chargeback_revocation_receipt_ref"} {
		requireStringSliceContains(t, asStringSlice(t, purchaseIngestion["required_evidence_slots"]), slot)
	}
	for _, correlation := range []string{"sandbox_order_ref", "license_ref", "entitlement_ref"} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), correlation)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "purchase_entitlement_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "purchase_entitlement_result")
	result := requireObject(t, record, "purchase_entitlement_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("purchase_entitlement_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"sandbox_order_ref", "sandbox_order_receipt", "sandbox_payment_receipt", "no_real_payment_capture_receipt", "real_payment_block_policy_receipt", "license_ref", "entitlement_ref", "entitlement_team_scope_verification_receipt_ref", "buyer_library_state_ref", "refund_or_chargeback_revocation_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"sandbox_payment_receipt_issues_entitlement", "no_real_payment_capture_receipt_present", "entitlement_team_scope_matches_target_team", "refund_or_chargeback_revokes_entitlement"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "purchase_entitlement_verified")
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresDownloadArtifactDeliveryWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	buyerLibraryPath := "commerce/buyer-library-install-state-candidate.json"
	entitlementPolicyPath := "commerce/license-entitlement-policy-candidate.json"
	digestPath := "commerce/artifact-bundle-digest-candidate.json"
	releasePath := "commerce/release-candidate-package.json"
	receiptChainPath := "commerce/commercial-receipt-chain-candidate.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_buyer_library_install_state": buyerLibraryPath,
		"source_license_entitlement_policy":  entitlementPolicyPath,
		"source_artifact_bundle_digest":      digestPath,
		"source_release_candidate_package":   releasePath,
		"source_commercial_receipt_chain":    receiptChainPath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "download_artifact_delivery_report")
	for _, guiEvidence := range []string{"owned_library_download_screen", "signed_download_result_screen", "download_hash_state_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"entitlement_check_receipt", "download_receipt", "signed_download_receipt", "download_artifact_hash_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "entitlement_ref", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#download_artifact_delivery_result", packagePath + "#download_artifact_delivery_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "download_artifact_delivery_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "download_artifact_delivery_verification_record_target"); got != "download_artifact_delivery_result" {
		t.Fatalf("download_artifact_delivery_verification_record_target = %s, want download_artifact_delivery_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "download_artifact_delivery_records")
	deliveryIngestion := requireObject(t, binder, "download_artifact_delivery_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "download_artifact_delivery_result",
		"package_section":     "download_artifact_delivery_report",
	} {
		if got := requireString(t, deliveryIngestion, key); got != want {
			t.Fatalf("download_artifact_delivery_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, sourceRef := range []string{buyerLibraryPath, entitlementPolicyPath, digestPath, releasePath, receiptChainPath} {
		requireStringSliceContains(t, asStringSlice(t, deliveryIngestion["source_refs"]), sourceRef)
	}
	for _, slot := range []string{"entitlement_ref", "entitlement_check_receipt_ref", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256", "bundle_tree_sha256", "download_artifact_hash_receipt_ref"} {
		requireStringSliceContains(t, asStringSlice(t, deliveryIngestion["required_evidence_slots"]), slot)
	}
	for _, correlation := range []string{"download_receipt_ref", "signed_download_ref", "download_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), correlation)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "download_artifact_delivery_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "download_artifact_delivery_result")
	result := requireObject(t, record, "download_artifact_delivery_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("download_artifact_delivery_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"entitlement_ref", "entitlement_check_receipt", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256", "bundle_tree_sha256", "download_artifact_hash_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"entitlement_allows_signed_download", "download_receipt_references_same_bundle_tree_sha256", "download_artifact_sha256_matches_release_candidate", "signed_download_ref_not_stored_as_secret"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "download_artifact_delivery_verified")

	pkg := readJSON(t, filepath.Join(base, packagePath))
	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "download_artifact_delivery_report")
	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "download_artifact_delivery_report")
	for _, sourceRef := range []string{buyerLibraryPath, entitlementPolicyPath, digestPath, releasePath, receiptChainPath} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{"download_artifact_delivery_verified", "entitlement_ref", "entitlement_check_receipt", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256", "bundle_tree_sha256", "download_artifact_hash_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{"download_artifact_delivery_report_missing", "entitlement_check_receipt_missing", "download_receipt_missing", "signed_download_ref_missing", "download_hash_mismatch"} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "download_artifact_delivery_verified")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["non_sufficient_evidence"]), "download_artifact_delivery_report_missing")
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresDownloadInstallAccessWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	matrixPath := "commerce/download-install-access-matrix.json"
	entitlementPolicyPath := "commerce/license-entitlement-policy-candidate.json"
	buyerLibraryPath := "commerce/buyer-library-install-state-candidate.json"
	installPreflightPath := "install/install-preflight-request-candidate.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	if got := requireString(t, checklist, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "download_install_access_matrix_report")
	for _, guiEvidence := range []string{"download_install_access_matrix_review_screen", "negative_download_install_block_screens", "install_preflight_block_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"unpaid_download_blocked_receipt", "refund_revoked_download_blocked_receipt", "authorization_expired_download_install_blocked_receipt", "version_unpublished_or_revoked_download_install_blocked_receipt", "artifact_hash_mismatch_blocked_receipt", "entitlement_verification_receipt", "artifact_hash_verification_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "entitlement_ref", "release_ref", "download_artifact_sha256", "install_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#download_install_access_result", packagePath + "#download_install_access_matrix_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "download_install_access_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "download_install_access_matrix_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "download_install_access_verification_record_target"); got != "download_install_access_result" {
		t.Fatalf("download_install_access_verification_record_target = %s, want download_install_access_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "download_install_access_records")
	accessIngestion := requireObject(t, binder, "download_install_access_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "download_install_access_result",
		"package_section":     "download_install_access_matrix_report",
	} {
		if got := requireString(t, accessIngestion, key); got != want {
			t.Fatalf("download_install_access_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, sourceRef := range []string{matrixPath, entitlementPolicyPath, buyerLibraryPath, installPreflightPath} {
		requireStringSliceContains(t, asStringSlice(t, accessIngestion["source_refs"]), sourceRef)
	}
	for _, slot := range []string{"unpaid_download_block_receipt_ref", "refund_revoked_download_block_receipt_ref", "authorization_expired_download_install_block_receipt_ref", "version_unpublished_or_revoked_download_install_block_receipt_ref", "artifact_hash_mismatch_install_block_receipt_ref", "entitlement_verification_receipt_ref", "release_state_ref", "download_artifact_sha256", "install_artifact_sha256", "user_visible_reason_refs"} {
		requireStringSliceContains(t, asStringSlice(t, accessIngestion["required_evidence_slots"]), slot)
	}
	for _, correlation := range []string{"release_ref", "download_artifact_sha256", "install_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), correlation)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "download_install_access_matrix_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "download_install_access_result")
	result := requireObject(t, record, "download_install_access_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("download_install_access_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"unpaid_download_blocked_receipt", "refund_revoked_download_blocked_receipt", "authorization_expired_download_install_blocked_receipt", "version_unpublished_or_revoked_download_install_blocked_receipt", "artifact_hash_mismatch_blocked_receipt", "entitlement_verification_receipt", "artifact_hash_verification_receipt", "download_artifact_sha256", "install_artifact_sha256"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"unpaid_user_cannot_download", "refund_revoked_blocks_download_and_install", "authorization_expired_blocks_download_install", "version_unpublished_or_revoked_blocks_download_install", "artifact_hash_mismatch_blocks_install", "download_install_hashes_match_allowed_artifact"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "download_install_access_matrix_verified")
}

func TestTeamOfficeP11EvidenceAcceptanceRequiresRoleStudioLineageWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	recordPath := "tests/p11-normal-commercialization-verification-record-template.json"
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	receiptChainPath := "commerce/commercial-receipt-chain-candidate.json"
	runtimeUsagePath := "usage/team-office-runtime-usage-candidate.json"
	installMapPath := "install/install-runtime-activation-map-candidate.json"
	bindingCatalogPath := "bindings/team-settings-installed-role-catalog-candidate.json"

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_commercial_receipt_chain":             receiptChainPath,
		"source_install_runtime_activation_map":       installMapPath,
		"source_runtime_usage_candidate":              runtimeUsagePath,
		"source_team_settings_installed_role_catalog": bindingCatalogPath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "role_studio_lineage_report")
	for _, guiEvidence := range []string{"role_studio_lineage_review_screen", "runtime_role_set_review_screen"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), guiEvidence)
	}
	for _, receipt := range []string{"candidate_bundle_export_receipt", "cloud_upload_receipt", "download_receipt", "install_receipt", "team_binding_receipt", "runtime_candidate_output_receipt"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), receipt)
	}
	for _, correlation := range []string{"candidate_set_ref", "bundle_tree_sha256", "six_role_pack_refs", "role_pack_ref_set_hash", "enabled_role_pack_version_refs"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), correlation)
	}
	for _, target := range []string{recordPath + "#role_studio_lineage_result", packagePath + "#role_studio_lineage_report"} {
		requireStringSliceContains(t, asStringSlice(t, stage["writeback_targets"]), target)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["required_before_p11_pass"]), "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, checklist, "completion_claim_policy")["non_sufficient_evidence"]), "role_studio_lineage_missing")

	binder := readJSON(t, filepath.Join(base, binderPath))
	if got := requireString(t, binder, "role_studio_lineage_verification_record_target"); got != "role_studio_lineage_result" {
		t.Fatalf("role_studio_lineage_verification_record_target = %s, want role_studio_lineage_result", got)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "ingestion_policy")["required_source_records"]), "role_studio_lineage_records")
	lineageIngestion := requireObject(t, binder, "role_studio_lineage_ingestion")
	for key, want := range map[string]string{
		"output_record_field": "role_studio_lineage_result",
		"package_section":     "role_studio_lineage_report",
	} {
		if got := requireString(t, lineageIngestion, key); got != want {
			t.Fatalf("role_studio_lineage_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, slot := range []string{"candidate_set_ref", "bundle_tree_sha256", "six_role_pack_refs", "role_pack_ref_set_hash", "export_receipt_ref", "cloud_upload_receipt_ref", "download_receipt_ref", "install_receipt_ref", "team_binding_receipt_ref", "runtime_session_ref"} {
		requireStringSliceContains(t, asStringSlice(t, lineageIngestion["required_evidence_slots"]), slot)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "receipt_correlation_rules")["required_correlation_keys"]), "role_pack_ref_set_hash")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, binder, "status_transition_rules")["required_before_passed_verified"]), "role_studio_lineage_verified")

	record := readJSON(t, filepath.Join(base, recordPath))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "record_schema")["required_top_level_fields"]), "role_studio_lineage_result")
	result := requireObject(t, record, "role_studio_lineage_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("role_studio_lineage_result.result_status = %s, want not_run", got)
	}
	for _, ref := range []string{"candidate_set_ref", "bundle_tree_sha256", "six_role_pack_refs", "role_pack_ref_set_hash", "export_receipt_ref", "cloud_upload_receipt_ref", "download_receipt_ref", "install_receipt_ref", "team_binding_receipt_ref", "runtime_session_ref"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	for _, check := range []string{"same_six_role_pack_refs_from_export_to_runtime", "same_role_pack_ref_set_hash_across_export_upload_download_install_binding_runtime"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_match_checks"]), check)
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, record, "final_decision")["required_before_passed_verified"]), "role_studio_lineage_verified")
}

func TestTeamOfficeP11EvidenceAcceptanceChecklistRequiresExecutionQueuePhaseDependencyProof(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	checklist := readJSON(t, filepath.Join(base, "tests", "p11-evidence-acceptance-checklist-candidate.json"))
	request := readJSON(t, filepath.Join(base, requestPath))
	queue := readJSON(t, filepath.Join(base, queuePath))

	contract := requireObject(t, request, "phase_dependency_contract")
	handoffs := asObjectSlice(t, contract["phase_handoffs"])
	links := asObjectSlice(t, queue["p11_phase_dependency_links"])
	if len(links) != len(handoffs) {
		t.Fatalf("p11_phase_dependency_links len = %d, want %d", len(links), len(handoffs))
	}

	if got := requireString(t, checklist, "source_execution_queue"); got != queuePath {
		t.Fatalf("source_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, checklist, "source_p11_phase_dependency_contract"); got != requestPath+"#phase_dependency_contract" {
		t.Fatalf("source_p11_phase_dependency_contract = %s", got)
	}

	gate := requireObject(t, checklist, "phase_dependency_acceptance_gate")
	for key, want := range map[string]string{
		"phase_dependency_links_ref":         queuePath + "#p11_phase_dependency_links",
		"phase_dependency_contract_ref":      requestPath + "#phase_dependency_contract",
		"verification_record_target":         "p11_phase_dependency_result",
		"go_live_evidence_package_section":   "p11_phase_dependency_report",
		"required_completion_proof":          "p11_execution_queue_phase_dependencies_verified",
		"missing_dependency_blocking_status": "blocked_previous_phase_evidence_missing",
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("phase_dependency_acceptance_gate.%s = %s, want %s", key, got, want)
		}
	}
	if got, ok := gate["expected_phase_dependency_count"].(float64); !ok || int(got) != len(handoffs) {
		t.Fatalf("expected_phase_dependency_count = %v, want %d", gate["expected_phase_dependency_count"], len(handoffs))
	}
	for _, evidence := range []string{
		"p11_phase_dependency_links",
		"previous_receipt_ref",
		"bundle_tree_sha256",
		"correlation_id",
		"blocked_previous_phase_evidence_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["required_evidence"]), evidence)
	}
	for _, blocker := range []string{
		"p11_execution_queue_phase_dependencies_verified_missing",
		"blocked_previous_phase_evidence_missing",
		"previous_receipt_ref_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, gate["blocking_if_missing"]), blocker)
	}

	policy := requireObject(t, checklist, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, policy["required_before_p11_pass"]), "p11_execution_queue_phase_dependencies_verified")
	requireStringSliceContains(t, asStringSlice(t, policy["non_sufficient_evidence"]), "final_receipts_without_phase_dependency_chain")
}

func TestTeamOfficeP11AcceptanceRequiresSixInstalledRoleCatalogBeforeBindingAndRuntime(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	activationPath := "install/install-runtime-activation-map-candidate.json"
	checklist := readJSON(t, filepath.Join(base, "tests", "p11-evidence-acceptance-checklist-candidate.json"))
	if got := requireString(t, checklist, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}

	stage := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "post_install_team_binding_runtime")
	for _, proof := range []string{
		"team_settings_installed_role_catalog_screenshot",
		"six_replaceable_roles_visible_screenshot",
		"secretary_and_five_advisors_slot_mapping_screenshot",
	} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_gui_evidence"]), proof)
	}
	for _, proof := range []string{
		"team_settings_catalog_refresh_receipt",
		"enabled_role_pack_version_refs_for_six_roles",
		"team_binding_receipt",
		"runtime_candidate_output_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_receipt_evidence"]), proof)
	}
	for _, proof := range []string{"bundle_tree_sha256", "team_ref", "enabled_role_pack_version_refs"} {
		requireStringSliceContains(t, asStringSlice(t, stage["required_hash_or_correlation_checks"]), proof)
	}
	for _, blocker := range []string{"six_role_catalog_incomplete", "enabled_role_pack_version_refs_missing"} {
		requireStringSliceContains(t, asStringSlice(t, stage["blocking_if_missing"]), blocker)
	}

	binder := readJSON(t, filepath.Join(base, "tests", "p11-evidence-ingestion-binder-candidate.json"))
	if got := requireString(t, binder, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("binder source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}
	replaceStep := findObjectByString(t, asObjectSlice(t, binder["gui_step_to_verification_map"]), "step_id", "replace_team_roles_after_install")
	for _, slot := range []string{
		"team_settings_catalog_refresh_receipt_ref",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		requireStringSliceContains(t, asStringSlice(t, replaceStep["required_evidence_slots"]), slot)
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	roleRefs := map[string]bool{}
	for _, ref := range asStringSlice(t, candidateSet["role_pack_refs"]) {
		roleRefs[ref] = false
	}
	activation := readJSON(t, filepath.Join(base, activationPath))
	activations := asObjectSlice(t, activation["required_role_activations"])
	if len(activations) != len(roleRefs) {
		t.Fatalf("required_role_activations len = %d, want %d", len(activations), len(roleRefs))
	}
	for _, item := range activations {
		roleRef := requireString(t, item, "role_pack_ref")
		if _, ok := roleRefs[roleRef]; !ok {
			t.Fatalf("unexpected required_role_activations role_pack_ref %s", roleRef)
		}
		roleRefs[roleRef] = true
		requireBool(t, item, "replaceable_in_team_settings", true)
		requireBool(t, item, "owner_gate_required_for_binding", true)
		requireBool(t, item, "runtime_usage_required", true)
		for _, key := range []string{"slot_ref", "enabled_role_pack_version_ref", "install_receipt_ref"} {
			if got := requireString(t, item, key); got == "" {
				t.Fatalf("%s missing for %s", key, roleRef)
			}
		}
	}
	for roleRef, seen := range roleRefs {
		if !seen {
			t.Fatalf("required_role_activations missing %s", roleRef)
		}
	}
}

func TestTeamOfficeCommercialEvidenceGateRequiresReceiptsAndIndependentAcceptance(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gate := readJSON(t, filepath.Join(base, "tests", "commercial-evidence-gate-candidate.json"))
	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	if got := requireString(t, gate, "completion_status"); got != "not_verified_requires_cross_repo_execution" {
		t.Fatalf("completion_status = %s", got)
	}

	truthSources := requireObject(t, gate, "truth_sources")
	for _, key := range []string{"frontend_truth", "cloud_truth", "install_truth", "acceptance_truth", "candidate_asset_truth"} {
		if got := requireString(t, truthSources, key); got == "" {
			t.Fatalf("truth_sources.%s missing", key)
		}
	}
	if got := requireString(t, truthSources, "cloud_truth"); got != "truzhen-cloud" {
		t.Fatalf("cloud_truth = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, truthSources, "install_truth"); got != "truzhenos" {
		t.Fatalf("install_truth = %s, want truzhenos", got)
	}

	schema := requireObject(t, gate, "evidence_record_schema")
	requiredFields := strings.Join(asStringSlice(t, schema["required_fields"]), "\n")
	for _, field := range []string{
		"evidence_ref",
		"stage_id",
		"actor",
		"truth_source",
		"gui_screenshot_path",
		"page_state_ref",
		"candidate_or_receipt_ref",
		"receipt_ref",
		"timestamp",
		"correlation_id",
		"artifact_sha256",
		"redaction_status",
		"verifier_ref",
		"result_status",
	} {
		if !strings.Contains(requiredFields, field) {
			t.Fatalf("evidence_record_schema missing field %s", field)
		}
	}
	statuses := strings.Join(asStringSlice(t, schema["allowed_result_statuses"]), "\n")
	for _, status := range []string{"pending", "passed_verified", "blocked_verified", "failed_requires_issue"} {
		if !strings.Contains(statuses, status) {
			t.Fatalf("allowed_result_statuses missing %s", status)
		}
	}

	rawStages, ok := gate["stage_evidence_gates"].([]any)
	if !ok {
		t.Fatalf("stage_evidence_gates missing")
	}
	requiredStages := map[string]bool{
		"publisher_create_role_candidates":      false,
		"publisher_export_candidate_bundle":     false,
		"publisher_upload_cloud_draft":          false,
		"cloud_review_candidate":                false,
		"buyer_sandbox_purchase":                false,
		"buyer_entitlement_and_signed_download": false,
		"buyer_local_install":                   false,
		"buyer_team_settings_replace_and_use":   false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage evidence gate = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		requireStringIn(t, requireString(t, stage, "truth_source"), "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo")
		if got := requireString(t, stage, "status"); got != "pending_cross_repo_execution" {
			t.Fatalf("%s status = %s, want pending_cross_repo_execution", stageID, got)
		}
		requiredEvidence := strings.Join(asStringSlice(t, stage["required_evidence_records"]), "\n")
		for _, proof := range []string{"gui_screenshot_path", "page_state_ref", "candidate_or_receipt_ref", "correlation_id"} {
			if !strings.Contains(requiredEvidence, proof) {
				t.Fatalf("%s required_evidence_records missing %s", stageID, proof)
			}
		}
		if strings.Contains(stageID, "cloud") || strings.Contains(stageID, "purchase") || strings.Contains(stageID, "download") {
			if !strings.Contains(requiredEvidence, "cloud_receipt_ref") && !strings.Contains(requiredEvidence, "download_receipt_ref") {
				t.Fatalf("%s must require cloud/download receipt evidence", stageID)
			}
		}
		if strings.Contains(stageID, "install") {
			if !strings.Contains(requiredEvidence, "install_receipt_ref") {
				t.Fatalf("%s must require install receipt evidence", stageID)
			}
		}
		if strings.Contains(stageID, "team_settings") {
			if !strings.Contains(requiredEvidence, "team_binding_receipt_ref") {
				t.Fatalf("%s must require team binding receipt evidence", stageID)
			}
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing stage evidence gate %s", stageID)
		}
	}

	acceptance := requireObject(t, gate, "independent_acceptance_policy")
	requireBool(t, acceptance, "independent_acceptance_agent_required", true)
	requireBool(t, acceptance, "organizer_self_attestation_allowed_for_completion", false)
	for _, key := range []string{"acceptance_artifact", "required_signoff", "issue_when_missing"} {
		if got := requireString(t, acceptance, key); got == "" {
			t.Fatalf("independent_acceptance_policy.%s missing", key)
		}
	}

	negative := requireObject(t, gate, "negative_evidence_gates")
	for _, key := range []string{"download_without_purchase", "tampered_artifact_hash", "expired_entitlement_install", "real_payment_without_owner_authorization"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "required_block_evidence"); got == "" {
			t.Fatalf("%s required_block_evidence missing", key)
		}
	}

	completion := requireObject(t, gate, "completion_gate")
	requireBool(t, completion, "completion_claim_allowed", false)
	requiredBeforeCompletion := strings.Join(asStringSlice(t, completion["required_before_completion_claim"]), "\n")
	for _, requirement := range []string{"all_required_evidence_records_verified", "cross_repo_authorization_records", "hash_continuity", "negative_cases_blocked_verified", "independent_acceptance_signoff"} {
		if !strings.Contains(requiredBeforeCompletion, requirement) {
			t.Fatalf("completion gate missing %s", requirement)
		}
	}

	forbidden := strings.Join(asStringSlice(t, gate["forbidden"]), "\n")
	for _, item := range []string{"completion_claim_from_candidate_json_only", "organizer_self_attestation_as_final_acceptance", "cloud_success_without_local_install_receipt", "gui_success_without_backend_receipt"} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_evidence_gate_candidate"); got != "tests/commercial-evidence-gate-candidate.json" {
		t.Fatalf("commercial_evidence_gate_candidate = %s", got)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, "tests/commercial-evidence-gate-candidate.json") {
		t.Fatalf("candidate set artifact_files missing commercial evidence gate")
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policy := requireObject(t, matrix, "completion_claim_policy")
	requiredProofs := strings.Join(asStringSlice(t, policy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "commercial_evidence_gate_verified") {
		t.Fatalf("product readiness matrix missing commercial_evidence_gate_verified completion proof")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/commercial-evidence-gate-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/commercial-evidence-gate-candidate.json")
}

func TestTeamOfficeSandboxEnvironmentReadinessDefinesSafeCommercialFixtures(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	readiness := readJSON(t, filepath.Join(base, "tests", "sandbox-environment-readiness-candidate.json"))
	requireBool(t, readiness, "candidate_only", true)
	requireBool(t, readiness, "non_formal", true)
	if got := requireString(t, readiness, "readiness_status"); got != "not_ready_requires_cross_repo_authorization" {
		t.Fatalf("readiness_status = %s, want not_ready_requires_cross_repo_authorization", got)
	}

	truthSources := requireObject(t, readiness, "truth_sources")
	if got := requireString(t, truthSources, "frontend_truth"); got != "truzhen-client-web-desktop" {
		t.Fatalf("frontend_truth = %s, want truzhen-client-web-desktop", got)
	}
	if got := requireString(t, truthSources, "cloud_truth"); got != "truzhen-cloud" {
		t.Fatalf("cloud_truth = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, truthSources, "install_truth"); got != "truzhenos" {
		t.Fatalf("install_truth = %s, want truzhenos", got)
	}

	policy := requireObject(t, readiness, "operation_policy")
	requireBool(t, policy, "sandbox_only", true)
	requireBool(t, policy, "no_real_credentials_in_pack", true)
	requireBool(t, policy, "production_endpoint_blocked", true)
	requireBool(t, policy, "real_payment_capture_blocked", true)
	requireBool(t, policy, "signed_download_secret_not_in_pack", true)

	actors := strings.Join(asStringSlice(t, readiness["sandbox_actors"]), "\n")
	for _, actor := range []string{"publisher_user_view_gui_agent", "buyer_user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent"} {
		if !strings.Contains(actors, actor) {
			t.Fatalf("sandbox_actors missing %s", actor)
		}
	}

	rawBindings, ok := readiness["required_environment_bindings"].([]any)
	if !ok {
		t.Fatalf("required_environment_bindings missing")
	}
	requiredBindings := map[string]bool{
		"cloud_sandbox_base_url_ref":            false,
		"cloud_marketplace_sandbox_project_ref": false,
		"publisher_sandbox_identity_ref":        false,
		"buyer_sandbox_identity_ref":            false,
		"sandbox_payment_method_ref":            false,
		"artifact_storage_sandbox_ref":          false,
		"signed_download_service_ref":           false,
		"local_install_target_ref":              false,
		"team_ref":                              false,
	}
	for _, raw := range rawBindings {
		binding, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("environment binding = %T", raw)
		}
		bindingID := requireString(t, binding, "binding_id")
		if _, ok := requiredBindings[bindingID]; ok {
			requiredBindings[bindingID] = true
		}
		requireStringIn(t, requireString(t, binding, "truth_source"), "truzhen-client-web-desktop", "truzhen-cloud", "truzhenos")
		if got := requireString(t, binding, "value_policy"); got != "ref_only_no_secret" {
			t.Fatalf("%s value_policy = %s, want ref_only_no_secret", bindingID, got)
		}
		if got := requireString(t, binding, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", bindingID)
		}
	}
	for bindingID, seen := range requiredBindings {
		if !seen {
			t.Fatalf("missing environment binding %s", bindingID)
		}
	}

	rawChecks, ok := readiness["preflight_checks"].([]any)
	if !ok {
		t.Fatalf("preflight_checks missing")
	}
	requiredChecks := map[string]bool{
		"cloud_sandbox_reachable":       false,
		"publisher_identity_sandbox":    false,
		"buyer_identity_sandbox":        false,
		"sandbox_payment_stub_ready":    false,
		"artifact_storage_ready":        false,
		"signed_download_ready":         false,
		"local_install_target_ready":    false,
		"team_settings_role_tab_ready":  false,
		"no_production_endpoint_loaded": false,
	}
	for _, raw := range rawChecks {
		check, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("preflight check = %T", raw)
		}
		checkID := requireString(t, check, "check_id")
		if _, ok := requiredChecks[checkID]; ok {
			requiredChecks[checkID] = true
		}
		requireStringIn(t, requireString(t, check, "status"), "pending_cross_repo_execution", "blocked_until_owner_authorizes")
		if got := requireString(t, check, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", checkID)
		}
	}
	for checkID, seen := range requiredChecks {
		if !seen {
			t.Fatalf("missing preflight check %s", checkID)
		}
	}

	negative := requireObject(t, readiness, "negative_cases")
	for _, key := range []string{"production_endpoint_configured", "real_payment_method_configured", "raw_secret_in_environment_binding", "signed_download_secret_in_pack"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, readiness["forbidden"]), "\n")
	for _, item := range []string{
		"store_cloud_access_token_in_pack",
		"store_payment_method_secret_in_pack",
		"use_production_endpoint_for_sandbox",
		"real_payment_capture_without_owner_authorization",
		"store_signed_download_secret_in_pack",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "sandbox_environment_readiness"); got != "tests/sandbox-environment-readiness-candidate.json" {
		t.Fatalf("sandbox_environment_readiness = %s", got)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, "tests/sandbox-environment-readiness-candidate.json") {
		t.Fatalf("candidate set artifact_files missing sandbox environment readiness")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/sandbox-environment-readiness-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/sandbox-environment-readiness-candidate.json")
}

func TestTeamOfficeP11SandboxPreflightGateBlocksRunUntilEnvironmentEvidenceExists(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gatePath := "tests/p11-sandbox-preflight-gate-candidate.json"
	gate := readJSON(t, filepath.Join(base, gatePath))
	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "can_start_p11_sandbox_run", false)
	if got := requireString(t, gate, "gate_status"); got != "not_ready_requires_authorization_and_environment_evidence" {
		t.Fatalf("gate_status = %s, want not_ready_requires_authorization_and_environment_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                        "role-pack-candidate-set://team-office-v0",
		"source_sandbox_environment_readiness":     "tests/sandbox-environment-readiness-candidate.json",
		"source_owner_authorization_intake":        "integration/owner-authorization-evidence-intake-candidate.json",
		"source_execution_queue":                   "integration/commercial-cross-repo-execution-queue-candidate.json",
		"source_runbook":                           "tests/p11-sandbox-execution-runbook-candidate.json",
		"source_cross_repo_readiness_package":      "integration/cross-repo-execution-readiness-package.json",
		"source_p11_evidence_acceptance_checklist": "tests/p11-evidence-acceptance-checklist-candidate.json",
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawChecks, ok := gate["preflight_gate_checks"].([]any)
	if !ok {
		t.Fatalf("preflight_gate_checks missing")
	}
	expectedChecks := map[string]bool{
		"owner_authorization_recorded":  false,
		"cloud_sandbox_reachable":       false,
		"publisher_identity_sandbox":    false,
		"buyer_identity_sandbox":        false,
		"sandbox_payment_stub_ready":    false,
		"artifact_storage_ready":        false,
		"signed_download_ready":         false,
		"local_install_target_ready":    false,
		"team_settings_role_tab_ready":  false,
		"no_production_endpoint_loaded": false,
	}
	if len(rawChecks) != len(expectedChecks) {
		t.Fatalf("preflight_gate_checks len = %d, want %d", len(rawChecks), len(expectedChecks))
	}
	for _, raw := range rawChecks {
		check, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("preflight gate check = %T", raw)
		}
		checkID := requireString(t, check, "check_id")
		if _, ok := expectedChecks[checkID]; !ok {
			t.Fatalf("unexpected preflight gate check %s", checkID)
		}
		expectedChecks[checkID] = true
		if got := requireString(t, check, "current_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("%s current_status = %s, want missing_authoritative_evidence", checkID, got)
		}
		requireBool(t, check, "can_start_p11_sandbox_run", false)
		requireStringIn(t, requireString(t, check, "truth_source"), "Owner", "truzhen-client-web-desktop", "truzhen-cloud", "truzhenos", "multi_repo")
		evidence := strings.Join(asStringSlice(t, check["required_authoritative_evidence"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "screenshot") && !strings.Contains(evidence, "scan") {
			t.Fatalf("%s required_authoritative_evidence must include receipt, screenshot, or scan", checkID)
		}
		if blockers := asStringSlice(t, check["blocking_if_missing"]); len(blockers) == 0 {
			t.Fatalf("%s blocking_if_missing must not be empty", checkID)
		}
	}
	for checkID, seen := range expectedChecks {
		if !seen {
			t.Fatalf("missing preflight gate check %s", checkID)
		}
	}

	policy := requireObject(t, gate, "start_policy")
	requireBool(t, policy, "can_start_now", false)
	requiredBeforeStart := strings.Join(asStringSlice(t, policy["required_before_start"]), "\n")
	for _, proof := range []string{
		"owner_authorization_recorded",
		"all_sandbox_environment_bindings_verified",
		"no_production_endpoint_loaded",
		"no_real_payment_method_loaded",
		"preflight_gate_recorded",
	} {
		if !strings.Contains(requiredBeforeStart, proof) {
			t.Fatalf("start_policy missing %s", proof)
		}
	}
	forbiddenPayloads := strings.Join(asStringSlice(t, gate["forbidden_payloads"]), "\n")
	for _, payload := range []string{"cloud_access_token", "payment_method_secret", "signed_download_secret", "production_endpoint", "raw_identity_document"} {
		if !strings.Contains(forbiddenPayloads, payload) {
			t.Fatalf("forbidden_payloads missing %s", payload)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_sandbox_preflight_gate"); got != gatePath {
		t.Fatalf("p11_sandbox_preflight_gate = %s, want %s", got, gatePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, gatePath) {
		t.Fatalf("candidate set artifact_files missing %s", gatePath)
	}
	runbook := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-execution-runbook-candidate.json"))
	if got := requireString(t, runbook, "sandbox_preflight_gate_ref"); got != gatePath {
		t.Fatalf("runbook sandbox_preflight_gate_ref = %s, want %s", got, gatePath)
	}
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	if got := requireString(t, queue, "p11_sandbox_preflight_gate_ref"); got != gatePath {
		t.Fatalf("execution queue p11_sandbox_preflight_gate_ref = %s, want %s", got, gatePath)
	}
	readinessPackage := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-readiness-package.json"))
	sourceRefs := strings.Join(asStringSlice(t, readinessPackage["source_refs"]), "\n")
	if !strings.Contains(sourceRefs, gatePath) {
		t.Fatalf("cross repo readiness package source_refs missing %s", gatePath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "p11_sandbox_preflight_gate_verified") {
		t.Fatalf("product readiness matrix missing p11_sandbox_preflight_gate_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != gatePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, gatePath))
		if err != nil {
			t.Fatalf("read %s: %v", gatePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", gatePath, got, wantHash)
		}
		return
	}
	t.Fatalf("artifact manifest missing %s", gatePath)
}

func TestTeamOfficeP11SandboxRunRequestDefinesAuthorizedEndToEndExecutionPacket(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	request := readJSON(t, filepath.Join(base, requestPath))
	requireBool(t, request, "candidate_only", true)
	requireBool(t, request, "non_formal", true)
	requireBool(t, request, "can_execute_now", false)
	if got := requireString(t, request, "request_status"); got != "not_ready_requires_owner_authorization_and_preflight_gate" {
		t.Fatalf("request_status = %s, want not_ready_requires_owner_authorization_and_preflight_gate", got)
	}
	if got := requireString(t, request, "candidate_set_ref"); got != "role-pack-candidate-set://team-office-v0" {
		t.Fatalf("candidate_set_ref = %s, want role-pack-candidate-set://team-office-v0", got)
	}

	sourceRefs := requireObject(t, request, "source_refs")
	for key, want := range map[string]string{
		"source_plan_ref":                      "docs/plans/role-pack-studio-team-office-test-plan-20260704.md",
		"source_preflight_gate":                "tests/p11-sandbox-preflight-gate-candidate.json",
		"source_runbook":                       "tests/p11-sandbox-execution-runbook-candidate.json",
		"source_evidence_acceptance_checklist": "tests/p11-evidence-acceptance-checklist-candidate.json",
		"source_gui_script":                    "tests/gui-user-agent-execution-script-candidate.json",
		"source_sandbox_environment_readiness": "tests/sandbox-environment-readiness-candidate.json",
		"source_cross_repo_execution_queue":    "integration/commercial-cross-repo-execution-queue-candidate.json",
		"source_owner_authorization_intake":    "integration/owner-authorization-evidence-intake-candidate.json",
		"source_artifact_bundle_digest":        "commerce/artifact-bundle-digest-candidate.json",
	} {
		if got := requireString(t, sourceRefs, key); got != want {
			t.Fatalf("source_refs.%s = %s, want %s", key, got, want)
		}
	}

	actors := strings.Join(asStringSlice(t, request["required_actors"]), "\n")
	for _, actor := range []string{
		"publisher_user_view_gui_agent",
		"buyer_user_view_gui_agent",
		"organizer_coordinator_recorder",
		"independent_acceptance_agent",
	} {
		if !strings.Contains(actors, actor) {
			t.Fatalf("required_actors missing %s", actor)
		}
	}

	rawPhases, ok := request["execution_phase_requests"].([]any)
	if !ok {
		t.Fatalf("execution_phase_requests missing")
	}
	expectedPhases := map[string]bool{
		"preflight_authorization_and_status":        false,
		"role_candidate_bundle_export":              false,
		"cloud_upload_listing_draft":                false,
		"marketplace_review_candidate":              false,
		"sandbox_order_payment_entitlement":         false,
		"entitled_signed_download":                  false,
		"local_install_enabled_version":             false,
		"post_install_team_binding_runtime":         false,
		"negative_cases_and_independent_acceptance": false,
		"production_promotion_controls":             false,
	}
	if len(rawPhases) != len(expectedPhases) {
		t.Fatalf("execution_phase_requests len = %d, want %d", len(rawPhases), len(expectedPhases))
	}
	for _, raw := range rawPhases {
		phase, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution phase = %T", raw)
		}
		phaseID := requireString(t, phase, "phase_id")
		if _, ok := expectedPhases[phaseID]; !ok {
			t.Fatalf("unexpected phase_id %s", phaseID)
		}
		expectedPhases[phaseID] = true
		if got := requireString(t, phase, "current_status"); got != "not_started" {
			t.Fatalf("%s current_status = %s, want not_started", phaseID, got)
		}
		requireBool(t, phase, "can_execute", false)
		for _, key := range []string{"target_repositories", "required_inputs", "required_evidence_outputs", "blocked_until", "writeback_targets"} {
			if got := asStringSlice(t, phase[key]); len(got) == 0 {
				t.Fatalf("%s %s must not be empty", phaseID, key)
			}
		}
		outputs := strings.Join(asStringSlice(t, phase["required_evidence_outputs"]), "\n")
		if !strings.Contains(outputs, "receipt") && !strings.Contains(outputs, "screenshot") && !strings.Contains(outputs, "signoff") && !strings.Contains(outputs, "hash") {
			t.Fatalf("%s required_evidence_outputs must include receipt, screenshot, signoff, or hash", phaseID)
		}
	}
	for phaseID, seen := range expectedPhases {
		if !seen {
			t.Fatalf("missing phase_id %s", phaseID)
		}
	}

	controls := requireObject(t, request, "execution_controls")
	for _, key := range []string{
		"user_view_gui_only",
		"organizer_may_record_only",
		"sandbox_only",
		"real_payment_blocked",
		"production_publish_blocked",
		"no_cross_repo_write_without_owner_authorization",
	} {
		requireBool(t, controls, key, true)
	}
	forbidden := strings.Join(asStringSlice(t, request["forbidden_actions"]), "\n")
	for _, action := range []string{
		"direct_api_call_as_user_action",
		"manual_json_success",
		"real_payment_capture",
		"production_publish",
		"store_cloud_access_token_in_pack",
		"store_signed_download_secret_in_pack",
	} {
		if !strings.Contains(forbidden, action) {
			t.Fatalf("forbidden_actions missing %s", action)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_sandbox_run_request"); got != requestPath {
		t.Fatalf("p11_sandbox_run_request = %s, want %s", got, requestPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, requestPath) {
		t.Fatalf("candidate set artifact_files missing %s", requestPath)
	}
	runbook := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-execution-runbook-candidate.json"))
	if got := requireString(t, runbook, "sandbox_run_request_ref"); got != requestPath {
		t.Fatalf("runbook sandbox_run_request_ref = %s, want %s", got, requestPath)
	}
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	if got := requireString(t, queue, "p11_sandbox_run_request_ref"); got != requestPath {
		t.Fatalf("execution queue p11_sandbox_run_request_ref = %s, want %s", got, requestPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "p11_sandbox_run_request_verified") {
		t.Fatalf("product readiness matrix missing p11_sandbox_run_request_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != requestPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, requestPath))
		if err != nil {
			t.Fatalf("read %s: %v", requestPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", requestPath, got, wantHash)
		}
		return
	}
	t.Fatalf("artifact manifest missing %s", requestPath)
}

func TestTeamOfficeP11SandboxRunRequestFeedsEvidenceClosureArtifacts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	request := readJSON(t, filepath.Join(base, requestPath))

	for file, fields := range map[string]map[string]string{
		"tests/p11-evidence-acceptance-checklist-candidate.json": {
			"source_p11_sandbox_run_request": requestPath,
		},
		"tests/p11-evidence-ingestion-binder-candidate.json": {
			"p11_sandbox_run_request_ref": requestPath,
		},
		"tests/p11-commercial-go-live-evidence-package-template.json": {
			"sandbox_run_request_ref": requestPath,
		},
		"tests/p11-normal-commercialization-verification-record-template.json": {
			"sandbox_run_request_ref": requestPath,
		},
		"tests/commercial-chain-verifier-candidate.json": {
			"source_p11_sandbox_run_request": requestPath,
		},
		"tests/commercial-readiness-verifier-candidate.json": {
			"source_p11_sandbox_run_request": requestPath,
		},
		"tests/commercial-go-no-go-gate-candidate.json": {
			"source_p11_sandbox_run_request": requestPath,
		},
		"docs/commercial-cross-repo-evidence-ledger.json": {
			"p11_sandbox_run_request_ref": requestPath,
		},
	} {
		doc := readJSON(t, filepath.Join(base, file))
		for key, want := range fields {
			if got := requireString(t, doc, key); got != want {
				t.Fatalf("%s %s = %s, want %s", file, key, got, want)
			}
		}
	}

	rawRequestPhases, ok := request["execution_phase_requests"].([]any)
	if !ok {
		t.Fatalf("execution_phase_requests missing")
	}
	expectedPhases := make(map[string]bool, len(rawRequestPhases))
	for _, raw := range rawRequestPhases {
		phase, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("request phase = %T", raw)
		}
		expectedPhases[requireString(t, phase, "phase_id")] = false
	}

	ledger := readJSON(t, filepath.Join(base, "docs", "commercial-cross-repo-evidence-ledger.json"))
	rawLinks, ok := ledger["p11_run_request_phase_links"].([]any)
	if !ok {
		t.Fatalf("p11_run_request_phase_links missing")
	}
	if len(rawLinks) != len(expectedPhases) {
		t.Fatalf("p11_run_request_phase_links len = %d, want %d", len(rawLinks), len(expectedPhases))
	}
	for _, raw := range rawLinks {
		link, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("phase link = %T", raw)
		}
		phaseID := requireString(t, link, "run_request_phase_id")
		if _, ok := expectedPhases[phaseID]; !ok {
			t.Fatalf("unexpected run_request_phase_id %s", phaseID)
		}
		expectedPhases[phaseID] = true
		if got := requireString(t, link, "current_status"); got != "pending_authorization" {
			t.Fatalf("%s current_status = %s, want pending_authorization", phaseID, got)
		}
		requireBool(t, link, "can_count_toward_goal_completion", false)
		for _, key := range []string{"evidence_ids", "writeback_targets", "blocking_if_missing"} {
			if values := asStringSlice(t, link[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", phaseID, key)
			}
		}
		ids := strings.Join(asStringSlice(t, link["evidence_ids"]), "\n")
		if !strings.Contains(ids, "role_studio_") {
			t.Fatalf("%s evidence_ids must link to role_studio evidence rows", phaseID)
		}
	}
	for phaseID, seen := range expectedPhases {
		if !seen {
			t.Fatalf("missing p11_run_request_phase_links for %s", phaseID)
		}
	}
}

func TestTeamOfficeCommercialEvidenceLedgerMirrorsP11PhaseDependencyContract(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	request := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-run-request-candidate.json"))
	contract := requireObject(t, request, "phase_dependency_contract")
	expected := map[string]map[string]any{}
	for _, handoff := range asObjectSlice(t, contract["phase_handoffs"]) {
		expected[requireString(t, handoff, "phase_id")] = handoff
	}

	ledger := readJSON(t, filepath.Join(base, "docs", "commercial-cross-repo-evidence-ledger.json"))
	links := asObjectSlice(t, ledger["p11_run_request_phase_links"])
	linkByPhase := map[string]map[string]any{}
	for _, link := range links {
		linkByPhase[requireString(t, link, "run_request_phase_id")] = link
	}

	for phaseID, handoff := range expected {
		link, ok := linkByPhase[phaseID]
		if !ok {
			t.Fatalf("p11_run_request_phase_links missing %s", phaseID)
		}
		for _, dep := range asStringSlice(t, handoff["depends_on_phase_ids"]) {
			requireStringSliceContains(t, asStringSlice(t, link["depends_on_phase_ids"]), dep)
		}
		for _, receipt := range asStringSlice(t, handoff["required_previous_receipts"]) {
			requireStringSliceContains(t, asStringSlice(t, link["required_previous_receipts"]), receipt)
		}
		if got := requireString(t, link, "dependency_blocking_status"); got != "blocked_previous_phase_evidence_missing" {
			t.Fatalf("%s dependency_blocking_status = %s", phaseID, got)
		}
		requireBool(t, link, "dependency_can_count_toward_goal_completion", false)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	required := asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"])
	requireStringSliceContains(t, required, "p11_evidence_ledger_phase_dependencies_verified")
}

func TestTeamOfficeCommercialExecutionQueueMirrorsP11PhaseDependencyContract(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	request := readJSON(t, filepath.Join(base, requestPath))
	contract := requireObject(t, request, "phase_dependency_contract")
	handoffs := asObjectSlice(t, contract["phase_handoffs"])

	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	if got := requireString(t, queue, "p11_phase_dependency_contract_ref"); got != requestPath+"#phase_dependency_contract" {
		t.Fatalf("p11_phase_dependency_contract_ref = %s, want %s", got, requestPath+"#phase_dependency_contract")
	}
	enforcement := requireObject(t, queue, "phase_dependency_enforcement")
	for _, key := range []string{
		"enabled",
		"no_stage_can_start_if_previous_receipts_missing",
		"same_correlation_id_required",
		"same_bundle_tree_sha256_required_from_export_onward",
	} {
		requireBool(t, enforcement, key, true)
	}
	if got := requireString(t, enforcement, "missing_dependency_result"); got != "blocked_previous_phase_evidence_missing" {
		t.Fatalf("phase_dependency_enforcement missing_dependency_result = %s", got)
	}

	rawLinks := asObjectSlice(t, queue["p11_phase_dependency_links"])
	if len(rawLinks) != len(handoffs) {
		t.Fatalf("p11_phase_dependency_links len = %d, want %d", len(rawLinks), len(handoffs))
	}
	linkByPhase := map[string]map[string]any{}
	for _, link := range rawLinks {
		linkByPhase[requireString(t, link, "run_request_phase_id")] = link
	}
	for _, handoff := range handoffs {
		phaseID := requireString(t, handoff, "phase_id")
		link, ok := linkByPhase[phaseID]
		if !ok {
			t.Fatalf("p11_phase_dependency_links missing %s", phaseID)
		}
		for _, dep := range asStringSlice(t, handoff["depends_on_phase_ids"]) {
			requireStringSliceContains(t, asStringSlice(t, link["depends_on_phase_ids"]), dep)
		}
		for _, receipt := range asStringSlice(t, handoff["required_previous_receipts"]) {
			requireStringSliceContains(t, asStringSlice(t, link["required_previous_receipts"]), receipt)
		}
		if got := requireString(t, link, "queue_blocking_status"); got != "blocked_previous_phase_evidence_missing" {
			t.Fatalf("%s queue_blocking_status = %s", phaseID, got)
		}
		if got := requireString(t, link, "evidence_ledger_ref"); got != "docs/commercial-cross-repo-evidence-ledger.json" {
			t.Fatalf("%s evidence_ledger_ref = %s", phaseID, got)
		}
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	required := asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"])
	requireStringSliceContains(t, required, "p11_execution_queue_phase_dependencies_verified")
}

func TestTeamOfficeP11SandboxRunRequestBindsCommercePhasesToGuiEvidenceAndReceiptContinuity(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	request := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-run-request-candidate.json"))
	protocol := readJSON(t, filepath.Join(base, "tests", "gui-evidence-capture-protocol.json"))

	protocolStages := map[string]bool{}
	for _, raw := range asObjectSlice(t, protocol["stage_capture_requirements"]) {
		protocolStages[requireString(t, raw, "stage_id")] = true
	}

	expected := map[string]struct {
		stages   []string
		receipts []string
	}{
		"role_candidate_bundle_export": {
			stages:   []string{"candidate_bundle_export"},
			receipts: []string{"candidate_bundle_export_receipt", "artifact_bundle_hash_receipt"},
		},
		"cloud_upload_listing_draft": {
			stages:   []string{"cloud_upload_draft"},
			receipts: []string{"cloud_upload_receipt", "listing_draft_receipt"},
		},
		"marketplace_review_candidate": {
			stages:   []string{"marketplace_review_candidate"},
			receipts: []string{"marketplace_review_candidate_receipt", "production_publish_block_receipt"},
		},
		"sandbox_order_payment_entitlement": {
			stages:   []string{"sandbox_purchase"},
			receipts: []string{"sandbox_order_receipt", "sandbox_payment_receipt", "entitlement_receipt"},
		},
		"entitled_signed_download": {
			stages:   []string{"entitlement_download"},
			receipts: []string{"download_receipt", "download_artifact_hash_receipt"},
		},
		"local_install_enabled_version": {
			stages:   []string{"local_install"},
			receipts: []string{"install_preflight_receipt", "role_pack_install_receipt", "enabled_role_pack_version_receipt"},
		},
		"post_install_team_binding_runtime": {
			stages:   []string{"team_settings_replace", "runtime_use"},
			receipts: []string{"team_binding_receipt", "runtime_candidate_output_receipt"},
		},
		"negative_cases_and_independent_acceptance": {
			stages:   []string{"negative_case"},
			receipts: []string{"negative_case_block_receipts", "independent_acceptance_signoff"},
		},
	}

	seen := map[string]bool{}
	for _, phase := range asObjectSlice(t, request["execution_phase_requests"]) {
		phaseID := requireString(t, phase, "phase_id")
		want, ok := expected[phaseID]
		if !ok {
			continue
		}
		seen[phaseID] = true
		contract := requireObject(t, phase, "gui_evidence_contract")
		if got := requireString(t, contract, "protocol_ref"); got != "tests/gui-evidence-capture-protocol.json" {
			t.Fatalf("%s protocol_ref = %s", phaseID, got)
		}
		requireBool(t, contract, "user_action_must_be_gui", true)
		requireBool(t, contract, "backend_only_success_forbidden", true)
		if got := requireString(t, contract, "cross_repo_writeback_status"); got != "pending_cross_repo_execution" {
			t.Fatalf("%s cross_repo_writeback_status = %s", phaseID, got)
		}
		stageIDs := asStringSlice(t, contract["gui_stage_ids"])
		for _, stageID := range want.stages {
			requireStringSliceContains(t, stageIDs, stageID)
			if !protocolStages[stageID] {
				t.Fatalf("%s references missing GUI protocol stage %s", phaseID, stageID)
			}
		}
		minimumEvidence := asStringSlice(t, contract["minimum_gui_evidence"])
		for _, field := range []string{"screenshot_path", "page_state_ref", "candidate_or_receipt_ref"} {
			requireStringSliceContains(t, minimumEvidence, field)
		}
		authorityReceipts := asStringSlice(t, contract["authority_receipts"])
		for _, receipt := range want.receipts {
			requireStringSliceContains(t, authorityReceipts, receipt)
		}
		correlationKeys := asStringSlice(t, contract["correlation_keys_required"])
		for _, key := range []string{"candidate_set_ref", "artifact_ref", "bundle_tree_sha256", "correlation_id"} {
			requireStringSliceContains(t, correlationKeys, key)
		}
	}
	for phaseID := range expected {
		if !seen[phaseID] {
			t.Fatalf("missing commerce phase %s", phaseID)
		}
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	required := asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"])
	requireStringSliceContains(t, required, "p11_gui_receipt_continuity_verified")
}

func TestTeamOfficeP11SandboxRunRequestDeclaresPhaseDependencyReceiptHandoffs(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	request := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-run-request-candidate.json"))

	rawPhases := asObjectSlice(t, request["execution_phase_requests"])
	phaseOrder := make([]string, 0, len(rawPhases))
	for _, phase := range rawPhases {
		phaseOrder = append(phaseOrder, requireString(t, phase, "phase_id"))
	}

	contract := requireObject(t, request, "phase_dependency_contract")
	requireBool(t, contract, "same_correlation_id_required", true)
	requireBool(t, contract, "same_bundle_tree_sha256_required_from_export_onward", true)
	requireBool(t, contract, "no_phase_can_skip_predecessor_receipt", true)
	requireBool(t, contract, "manual_phase_reorder_forbidden", true)
	requireBool(t, contract, "completion_claim_blocked_if_dependency_missing", true)

	chainOrder := asStringSlice(t, contract["chain_order"])
	if len(chainOrder) != len(phaseOrder) {
		t.Fatalf("chain_order len = %d, want %d", len(chainOrder), len(phaseOrder))
	}
	for i := range phaseOrder {
		if chainOrder[i] != phaseOrder[i] {
			t.Fatalf("chain_order[%d] = %s, want %s", i, chainOrder[i], phaseOrder[i])
		}
	}

	requiredHandoffs := map[string]struct {
		previous string
		receipts []string
	}{
		"role_candidate_bundle_export": {
			previous: "preflight_authorization_and_status",
			receipts: []string{
				"owner_authorization_receipt",
				"preflight_gate_receipt",
			},
		},
		"cloud_upload_listing_draft": {
			previous: "role_candidate_bundle_export",
			receipts: []string{
				"candidate_bundle_export_receipt",
				"artifact_bundle_hash_receipt",
			},
		},
		"marketplace_review_candidate": {
			previous: "cloud_upload_listing_draft",
			receipts: []string{
				"cloud_upload_receipt",
				"listing_draft_receipt",
			},
		},
		"sandbox_order_payment_entitlement": {
			previous: "marketplace_review_candidate",
			receipts: []string{
				"listing_draft_receipt",
				"marketplace_review_candidate_receipt",
			},
		},
		"entitled_signed_download": {
			previous: "sandbox_order_payment_entitlement",
			receipts: []string{
				"entitlement_receipt",
				"sandbox_payment_receipt",
			},
		},
		"local_install_enabled_version": {
			previous: "entitled_signed_download",
			receipts: []string{
				"download_receipt",
				"download_artifact_hash_receipt",
			},
		},
		"post_install_team_binding_runtime": {
			previous: "local_install_enabled_version",
			receipts: []string{
				"role_pack_install_receipt",
				"enabled_role_pack_version_receipt",
			},
		},
		"negative_cases_and_independent_acceptance": {
			previous: "post_install_team_binding_runtime",
			receipts: []string{
				"team_binding_receipt",
				"runtime_candidate_output_receipt",
			},
		},
		"production_promotion_controls": {
			previous: "negative_cases_and_independent_acceptance",
			receipts: []string{
				"independent_acceptance_signoff",
				"owner_go_no_go_signoff",
			},
		},
	}

	seen := map[string]bool{}
	for _, handoff := range asObjectSlice(t, contract["phase_handoffs"]) {
		phaseID := requireString(t, handoff, "phase_id")
		want, ok := requiredHandoffs[phaseID]
		if !ok {
			t.Fatalf("unexpected handoff phase_id %s", phaseID)
		}
		seen[phaseID] = true
		requireStringSliceContains(t, asStringSlice(t, handoff["depends_on_phase_ids"]), want.previous)
		requiredReceipts := asStringSlice(t, handoff["required_previous_receipts"])
		for _, receipt := range want.receipts {
			requireStringSliceContains(t, requiredReceipts, receipt)
		}
		references := asStringSlice(t, handoff["outputs_must_reference"])
		for _, key := range []string{"candidate_set_ref", "artifact_ref", "bundle_tree_sha256", "correlation_id", "previous_receipt_ref"} {
			requireStringSliceContains(t, references, key)
		}
		if got := requireString(t, handoff, "missing_dependency_result"); got != "blocked_previous_phase_evidence_missing" {
			t.Fatalf("%s missing_dependency_result = %s", phaseID, got)
		}
	}
	for phaseID := range requiredHandoffs {
		if !seen[phaseID] {
			t.Fatalf("phase_handoffs missing %s", phaseID)
		}
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	required := asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"])
	requireStringSliceContains(t, required, "p11_phase_dependency_continuity_verified")
}

func TestTeamOfficeCommercialObservabilityDiagnosticsLinksGuiCloudInstallReceipts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	diagnostics := readJSON(t, filepath.Join(base, "tests", "commercial-observability-diagnostics-candidate.json"))
	requireBool(t, diagnostics, "candidate_only", true)
	requireBool(t, diagnostics, "non_formal", true)
	if got := requireString(t, diagnostics, "diagnostics_status"); got != "candidate_not_wired" {
		t.Fatalf("diagnostics_status = %s, want candidate_not_wired", got)
	}

	truthSources := requireObject(t, diagnostics, "truth_sources")
	for _, key := range []string{"frontend_truth", "cloud_truth", "install_truth", "candidate_asset_truth"} {
		if got := requireString(t, truthSources, key); got == "" {
			t.Fatalf("truth_sources.%s missing", key)
		}
	}

	policy := requireObject(t, diagnostics, "data_policy")
	requireBool(t, policy, "no_pii_payload_in_packs", true)
	requireBool(t, policy, "no_secret_payload_in_packs", true)
	requireBool(t, policy, "redaction_required", true)
	requireBool(t, policy, "receipt_ref_not_payload", true)

	rawEvents, ok := diagnostics["required_events"].([]any)
	if !ok {
		t.Fatalf("required_events missing")
	}
	requiredEvents := map[string]bool{
		"role_candidate_created":     false,
		"candidate_bundle_exported":  false,
		"cloud_upload_draft_created": false,
		"sandbox_order_created":      false,
		"sandbox_payment_confirmed":  false,
		"entitlement_issued":         false,
		"signed_download_created":    false,
		"local_install_completed":    false,
		"team_role_binding_replaced": false,
		"runtime_candidate_used":     false,
		"negative_case_blocked":      false,
	}
	for _, raw := range rawEvents {
		event, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("diagnostic event = %T", raw)
		}
		eventID := requireString(t, event, "event_id")
		if _, ok := requiredEvents[eventID]; ok {
			requiredEvents[eventID] = true
		}
		requireStringIn(t, requireString(t, event, "truth_source"), "truzhen-client-web-desktop", "truzhen-cloud", "truzhenos", "multi_repo")
		fields := strings.Join(asStringSlice(t, event["required_fields"]), "\n")
		for _, field := range []string{"correlation_id", "stage_id", "timestamp", "actor_ref", "status"} {
			if !strings.Contains(fields, field) {
				t.Fatalf("%s required_fields missing %s", eventID, field)
			}
		}
		if !strings.Contains(fields, "receipt_ref") && !strings.Contains(fields, "candidate_ref") {
			t.Fatalf("%s required_fields must include receipt_ref or candidate_ref", eventID)
		}
	}
	for eventID, seen := range requiredEvents {
		if !seen {
			t.Fatalf("missing diagnostic event %s", eventID)
		}
	}

	rawSignals, ok := diagnostics["diagnostic_signals"].([]any)
	if !ok {
		t.Fatalf("diagnostic_signals missing")
	}
	requiredSignals := map[string]bool{
		"gui_step_trace":             false,
		"cloud_commerce_trace":       false,
		"install_preflight_trace":    false,
		"artifact_hash_continuity":   false,
		"entitlement_decision_trace": false,
		"owner_gate_decision_trace":  false,
	}
	for _, raw := range rawSignals {
		signal, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("diagnostic signal = %T", raw)
		}
		signalID := requireString(t, signal, "signal_id")
		if _, ok := requiredSignals[signalID]; ok {
			requiredSignals[signalID] = true
		}
		if got := requireString(t, signal, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", signalID)
		}
		requireStringIn(t, requireString(t, signal, "storage_policy"), "receipt_ref_only", "redacted_summary_only", "hash_only")
	}
	for signalID, seen := range requiredSignals {
		if !seen {
			t.Fatalf("missing diagnostic signal %s", signalID)
		}
	}

	rawAlerts, ok := diagnostics["alert_rules"].([]any)
	if !ok {
		t.Fatalf("alert_rules missing")
	}
	requiredAlerts := map[string]bool{
		"missing_receipt_after_gui_success": false,
		"artifact_hash_mismatch":            false,
		"entitlement_install_denied":        false,
		"production_endpoint_detected":      false,
		"secret_or_pii_detected":            false,
	}
	for _, raw := range rawAlerts {
		alert, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("alert rule = %T", raw)
		}
		alertID := requireString(t, alert, "alert_id")
		if _, ok := requiredAlerts[alertID]; ok {
			requiredAlerts[alertID] = true
		}
		requireStringIn(t, requireString(t, alert, "severity"), "yellow", "orange", "red")
		if got := requireString(t, alert, "expected_action"); got == "" {
			t.Fatalf("%s expected_action missing", alertID)
		}
	}
	for alertID, seen := range requiredAlerts {
		if !seen {
			t.Fatalf("missing alert rule %s", alertID)
		}
	}

	forbidden := strings.Join(asStringSlice(t, diagnostics["forbidden"]), "\n")
	for _, item := range []string{
		"upload_raw_log_payload_to_packs",
		"store_personal_data_in_diagnostics",
		"store_secret_in_diagnostics",
		"claim_gui_success_without_receipt_trace",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_observability_diagnostics"); got != "tests/commercial-observability-diagnostics-candidate.json" {
		t.Fatalf("commercial_observability_diagnostics = %s", got)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, "tests/commercial-observability-diagnostics-candidate.json") {
		t.Fatalf("candidate set artifact_files missing commercial observability diagnostics")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/commercial-observability-diagnostics-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/commercial-observability-diagnostics-candidate.json")
}

func TestTeamOfficeLicenseEntitlementPolicyKeepsCommerceTruthInCloud(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "license-entitlement-policy-candidate.json"))
	requireBool(t, policy, "candidate_only", true)
	requireBool(t, policy, "non_formal", true)
	if got := requireString(t, policy, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, policy, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s", got)
	}

	license := requireObject(t, policy, "license_candidate")
	if got := requireString(t, license, "status"); got != "candidate_not_formal" {
		t.Fatalf("license status = %s, want candidate_not_formal", got)
	}
	requireStringIn(t, requireString(t, license, "scope"), "team_scoped", "owner_scoped")

	entitlement := requireObject(t, policy, "entitlement_candidate")
	for _, key := range []string{"issue_trigger", "download_requirement", "install_requirement"} {
		if got := requireString(t, entitlement, key); got == "" {
			t.Fatalf("entitlement %s missing", key)
		}
	}

	rawChecks, ok := policy["validation_checks"].([]any)
	if !ok {
		t.Fatalf("validation_checks missing")
	}
	requiredChecks := map[string]bool{
		"purchase_before_download":           false,
		"entitlement_before_install":         false,
		"artifact_hash_match":                false,
		"refund_blocks_future_download":      false,
		"expired_entitlement_blocks_install": false,
	}
	for _, raw := range rawChecks {
		check, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("validation check = %T", raw)
		}
		id := requireString(t, check, "check_id")
		if _, ok := requiredChecks[id]; ok {
			requiredChecks[id] = true
		}
		requireStringIn(t, requireString(t, check, "truth_source"), "truzhen-cloud", "truzhenos", "truzhen-cloud + truzhenos")
		if got := requireString(t, check, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", id)
		}
	}
	for id, seen := range requiredChecks {
		if !seen {
			t.Fatalf("missing validation check %s", id)
		}
	}

	negative := requireObject(t, policy, "negative_cases")
	for _, key := range []string{"download_without_purchase", "install_without_entitlement", "refund_after_purchase", "expired_entitlement", "real_payment_without_owner_authorization"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}

	forbidden := strings.Join(asStringSlice(t, policy["forbidden"]), "\n")
	for _, phrase := range []string{"store_order_truth_in_truzhen_packs", "store_license_truth_in_truzhen_packs", "store_entitlement_truth_in_truzhen_packs", "real_payment_capture_without_owner_authorization"} {
		if !strings.Contains(forbidden, phrase) {
			t.Fatalf("forbidden missing %s", phrase)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/license-entitlement-policy-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/license-entitlement-policy-candidate.json")
}

func TestTeamOfficeTeamScopedEntitlementMustMatchInstallTargetTeam(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "license-entitlement-policy-candidate.json"))

	license := requireObject(t, policy, "license_candidate")
	if got := requireString(t, license, "scope"); got != "team_scoped" {
		t.Fatalf("license scope = %s, want team_scoped", got)
	}

	entitlement := requireObject(t, policy, "entitlement_candidate")
	teamScope := requireObject(t, entitlement, "team_scope_binding")
	if got := requireString(t, teamScope, "licensed_team_ref_source"); got != "truzhen-cloud_entitlement" {
		t.Fatalf("licensed_team_ref_source = %s, want truzhen-cloud_entitlement", got)
	}
	requireBool(t, teamScope, "target_team_ref_must_match", true)
	if got := requireString(t, teamScope, "cross_team_install_policy"); got != "blocked_entitlement_team_scope_mismatch" {
		t.Fatalf("cross_team_install_policy = %s", got)
	}
	if got := requireString(t, teamScope, "expected_evidence"); got != "entitlement_team_scope_verification_receipt_ref" {
		t.Fatalf("expected_evidence = %s", got)
	}

	scopeCheckFound := false
	for _, check := range asObjectSlice(t, policy["validation_checks"]) {
		if requireString(t, check, "check_id") != "entitlement_team_scope_matches_target_team" {
			continue
		}
		scopeCheckFound = true
		if got := requireString(t, check, "truth_source"); got != "truzhen-cloud + truzhenos" {
			t.Fatalf("scope check truth_source = %s", got)
		}
		if got := requireString(t, check, "blocked_status"); got != "blocked_entitlement_team_scope_mismatch" {
			t.Fatalf("scope check blocked_status = %s", got)
		}
	}
	if !scopeCheckFound {
		t.Fatalf("validation_checks missing entitlement_team_scope_matches_target_team")
	}

	negative := requireObject(t, policy, "negative_cases")
	crossTeam := requireObject(t, negative, "cross_team_install")
	if got := requireString(t, crossTeam, "expected_status"); got != "blocked_entitlement_team_scope_mismatch" {
		t.Fatalf("cross_team_install expected_status = %s", got)
	}

	preflight := readJSON(t, filepath.Join(base, "install", "install-preflight-request-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, preflight["required_inputs"]), "licensed_team_ref")
	preflightScopeCheckFound := false
	for _, check := range asObjectSlice(t, preflight["preflight_checks"]) {
		if requireString(t, check, "check_id") != "entitlement_team_scope_match" {
			continue
		}
		preflightScopeCheckFound = true
		if got := requireString(t, check, "required_input"); got != "entitlement_ref + licensed_team_ref + target_team_ref" {
			t.Fatalf("preflight scope required_input = %s", got)
		}
		if got := requireString(t, check, "success_status"); got != "candidate_entitlement_team_scope_matches_target_team" {
			t.Fatalf("preflight scope success_status = %s", got)
		}
		if got := requireString(t, check, "expected_evidence"); got != "entitlement_team_scope_verification_receipt_ref" {
			t.Fatalf("preflight scope expected_evidence = %s", got)
		}
	}
	if !preflightScopeCheckFound {
		t.Fatalf("preflight_checks missing entitlement_team_scope_match")
	}

	preflightNegative := requireObject(t, preflight, "negative_cases")
	preflightCrossTeam := requireObject(t, preflightNegative, "cross_team_install")
	if got := requireString(t, preflightCrossTeam, "expected_status"); got != "blocked_entitlement_team_scope_mismatch" {
		t.Fatalf("preflight cross_team_install expected_status = %s", got)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policyBlock := requireObject(t, matrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, policyBlock["required_before_completion_claim"]), "entitlement_team_scope_verified")
}

func TestTeamOfficeOrderPaymentStateMachineDefinesSandboxPurchaseAndRefund(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	stateMachine := readJSON(t, filepath.Join(base, "commerce", "order-payment-state-machine-candidate.json"))
	requireBool(t, stateMachine, "candidate_only", true)
	requireBool(t, stateMachine, "non_formal", true)
	if got := requireString(t, stateMachine, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, stateMachine, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	rawStates, ok := stateMachine["states"].([]any)
	if !ok {
		t.Fatalf("states missing")
	}
	requiredStates := map[string]string{
		"listing_draft_ready":     "candidate_ready_for_sandbox_order",
		"sandbox_order_created":   "candidate_order_created",
		"sandbox_payment_pending": "candidate_payment_pending",
		"sandbox_payment_paid":    "candidate_payment_paid_sandbox_only",
		"entitlement_issued":      "candidate_entitlement_issued",
		"payment_failed":          "blocked_payment_failed",
		"refund_succeeded":        "blocked_entitlement_revocation_required",
		"chargeback_received":     "blocked_entitlement_revocation_required",
		"entitlement_revoked":     "blocked_future_download_and_install",
	}
	for _, raw := range rawStates {
		state, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("state = %T", raw)
		}
		stateID := requireString(t, state, "state_id")
		wantStatus, ok := requiredStates[stateID]
		if !ok {
			continue
		}
		delete(requiredStates, stateID)
		if got := requireString(t, state, "status"); got != wantStatus {
			t.Fatalf("%s status = %s, want %s", stateID, got, wantStatus)
		}
		if got := requireString(t, state, "truth_source"); got != "truzhen-cloud" {
			t.Fatalf("%s truth_source = %s, want truzhen-cloud", stateID, got)
		}
		evidence := strings.Join(asStringSlice(t, state["evidence_required"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "state_ref") {
			t.Fatalf("%s evidence_required must include receipt or state_ref", stateID)
		}
	}
	for stateID := range requiredStates {
		t.Fatalf("missing state %s", stateID)
	}

	rawTransitions, ok := stateMachine["transitions"].([]any)
	if !ok {
		t.Fatalf("transitions missing")
	}
	requiredTransitions := map[string]bool{
		"create_sandbox_order":           false,
		"mark_sandbox_payment_paid":      false,
		"issue_entitlement_after_paid":   false,
		"fail_payment":                   false,
		"refund_revokes_entitlement":     false,
		"chargeback_revokes_entitlement": false,
	}
	for _, raw := range rawTransitions {
		transition, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("transition = %T", raw)
		}
		id := requireString(t, transition, "transition_id")
		if _, ok := requiredTransitions[id]; ok {
			requiredTransitions[id] = true
		}
		if got := requireString(t, transition, "truth_source"); got != "truzhen-cloud" {
			t.Fatalf("%s truth_source = %s, want truzhen-cloud", id, got)
		}
		requireString(t, transition, "required_evidence")
	}
	for id, seen := range requiredTransitions {
		if !seen {
			t.Fatalf("missing transition %s", id)
		}
	}

	negative := requireObject(t, stateMachine, "negative_cases")
	for _, key := range []string{"real_payment_without_owner_authorization", "payment_success_without_order", "entitlement_without_paid_receipt", "refund_without_revocation", "store_payment_truth_in_packs"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		requireString(t, item, "expected_evidence")
	}

	forbidden := strings.Join(asStringSlice(t, stateMachine["forbidden"]), "\n")
	for _, phrase := range []string{"store_order_truth_in_truzhen_packs", "store_payment_truth_in_truzhen_packs", "store_payment_method_in_truzhen_packs", "real_payment_capture_without_owner_authorization"} {
		if !strings.Contains(forbidden, phrase) {
			t.Fatalf("forbidden missing %s", phrase)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/order-payment-state-machine-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/order-payment-state-machine-candidate.json")
}

func TestTeamOfficeDownloadInstallAccessMatrixBlocksInvalidCommerceStates(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrix := readJSON(t, filepath.Join(base, "commerce", "download-install-access-matrix.json"))
	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	if got := requireString(t, matrix, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, matrix, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}

	surfaces := strings.Join(asStringSlice(t, matrix["entry_surfaces"]), "\n")
	for _, surface := range []string{"cloud_purchased_page", "signed_download", "local_pack_manager", "team_settings_role_tab"} {
		if !strings.Contains(surfaces, surface) {
			t.Fatalf("entry_surfaces missing %s", surface)
		}
	}

	rawStates, ok := matrix["access_states"].([]any)
	if !ok {
		t.Fatalf("access_states missing")
	}
	requiredStates := map[string][2]string{
		"unpaid_user":                    {"blocked_entitlement_missing", "blocked_entitlement_missing"},
		"sandbox_paid_entitled_user":     {"allowed_candidate_signed_download", "allowed_candidate_install_preflight"},
		"refund_revoked_entitlement":     {"blocked_entitlement_revoked", "blocked_entitlement_revoked"},
		"expired_entitlement":            {"blocked_entitlement_expired", "blocked_entitlement_expired"},
		"unpublished_or_revoked_version": {"blocked_release_not_available", "blocked_release_revoked"},
		"artifact_hash_mismatch":         {"blocked_artifact_hash_mismatch", "blocked_artifact_hash_mismatch"},
	}
	for _, raw := range rawStates {
		state, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("access state = %T", raw)
		}
		stateID := requireString(t, state, "state_id")
		want, ok := requiredStates[stateID]
		if !ok {
			continue
		}
		delete(requiredStates, stateID)
		if got := requireString(t, state, "download_status"); got != want[0] {
			t.Fatalf("%s download_status = %s, want %s", stateID, got, want[0])
		}
		if got := requireString(t, state, "install_status"); got != want[1] {
			t.Fatalf("%s install_status = %s, want %s", stateID, got, want[1])
		}
		if got := requireString(t, state, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s cannot use truzhen-packs as truth source", stateID)
		}
		for _, key := range []string{"download_evidence", "install_evidence", "user_visible_reason"} {
			requireString(t, state, key)
		}
	}
	for stateID := range requiredStates {
		t.Fatalf("missing access state %s", stateID)
	}

	downloadGate := requireObject(t, matrix, "download_gate")
	for _, key := range []string{"entitlement_required", "release_available_required", "signature_required", "artifact_hash_required"} {
		if got := requireString(t, downloadGate, key); got != "true" {
			t.Fatalf("download_gate.%s = %s, want true", key, got)
		}
	}
	installGate := requireObject(t, matrix, "install_gate")
	for _, key := range []string{"entitlement_verification", "artifact_hash_verification", "schema_validation", "forbidden_artifact_scan", "owner_gate_for_team_binding"} {
		if got := requireString(t, installGate, key); got != "required" {
			t.Fatalf("install_gate.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, matrix, "negative_cases")
	for _, key := range []string{"unpaid_download", "refund_after_purchase_download", "expired_entitlement_install", "version_unpublished_download", "artifact_hash_mismatch_install", "team_binding_without_owner_gate"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		requireString(t, item, "expected_evidence")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/download-install-access-matrix.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/download-install-access-matrix.json")
}

func TestTeamOfficeReleaseCandidatePackageDefinesSignedDistributionBoundary(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	release := readJSON(t, filepath.Join(base, "commerce", "release-candidate-package.json"))
	requireBool(t, release, "candidate_only", true)
	requireBool(t, release, "non_formal", true)
	if got := requireString(t, release, "release_status"); got != "candidate_not_published" {
		t.Fatalf("release_status = %s, want candidate_not_published", got)
	}
	if got := requireString(t, release, "release_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("release_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, release, "artifact_ref"); got != "role-pack-artifact://team-office-v0" {
		t.Fatalf("artifact_ref = %s", got)
	}
	requireStringIn(t, requireString(t, release, "package_format"), "role_pack_candidate_bundle")

	channels := strings.Join(asStringSlice(t, release["allowed_channels"]), "\n")
	for _, channel := range []string{"cloud_draft", "sandbox_purchase", "signed_download"} {
		if !strings.Contains(channels, channel) {
			t.Fatalf("allowed_channels missing %s", channel)
		}
	}
	if strings.Contains(channels, "production_publish") {
		t.Fatalf("production_publish cannot be an allowed channel without Owner authorization")
	}

	signing := requireObject(t, release, "signing_policy")
	if got := requireString(t, signing, "signature_required"); got != "true" {
		t.Fatalf("signature_required = %s, want true", got)
	}
	if got := requireString(t, signing, "signing_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("signing_truth_source = %s, want truzhen-cloud", got)
	}
	forbiddenSigning := strings.Join(asStringSlice(t, signing["forbidden"]), "\n")
	if !strings.Contains(forbiddenSigning, "private_key_in_truzhen_packs") {
		t.Fatalf("signing forbidden must reject private_key_in_truzhen_packs")
	}

	upgrade := requireObject(t, release, "upgrade_policy")
	if got := requireString(t, upgrade, "compatibility_check_required"); got != "true" {
		t.Fatalf("compatibility_check_required = %s, want true", got)
	}
	revocation := requireObject(t, release, "revocation_policy")
	if got := requireString(t, revocation, "revocation_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("revocation_truth_source = %s, want truzhen-cloud", got)
	}

	negative := requireObject(t, release, "negative_cases")
	for _, key := range []string{"production_publish_without_owner_authorization", "unsigned_artifact_download", "stale_version_install", "signature_key_in_pack"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/release-candidate-package.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/release-candidate-package.json")
}

func TestTeamOfficeRuntimeUsageCandidateKeepsOutputsCandidateOnly(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	usage := readJSON(t, filepath.Join(base, "usage", "team-office-runtime-usage-candidate.json"))
	requireBool(t, usage, "candidate_only", true)
	requireBool(t, usage, "non_formal", true)
	if got := requireString(t, usage, "runtime_truth_source"); got != "truzhenos" {
		t.Fatalf("runtime_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, usage, "gui_truth_source"); got != "truzhen-client-web-desktop" {
		t.Fatalf("gui_truth_source = %s, want truzhen-client-web-desktop", got)
	}

	rawScenarios, ok := usage["usage_scenarios"].([]any)
	if !ok {
		t.Fatalf("usage_scenarios missing")
	}
	required := map[string]bool{
		"secretary_orchestrates_five_advisors":    false,
		"advisors_return_candidate_outputs":       false,
		"capability_invocation_uses_role_context": false,
		"owner_gate_before_formalization":         false,
	}
	for _, raw := range rawScenarios {
		scenario, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("usage scenario = %T", raw)
		}
		id := requireString(t, scenario, "scenario_id")
		if _, ok := required[id]; ok {
			required[id] = true
		}
		if got := requireString(t, scenario, "output_status"); got != "candidate_only" {
			t.Fatalf("%s output_status = %s, want candidate_only", id, got)
		}
		evidence := strings.Join(asStringSlice(t, scenario["required_evidence"]), "\n")
		for _, proof := range []string{"gui_screenshot", "runtime_receipt_candidate_ref"} {
			if !strings.Contains(evidence, proof) {
				t.Fatalf("%s required_evidence missing %s", id, proof)
			}
		}
	}
	for id, seen := range required {
		if !seen {
			t.Fatalf("missing usage scenario %s", id)
		}
	}

	negative := requireObject(t, usage, "negative_cases")
	for _, key := range []string{"formal_task_without_owner_gate", "direct_execution_from_role", "memory_write_without_gate", "send_message_without_gate"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "usage/team-office-runtime-usage-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing usage/team-office-runtime-usage-candidate.json")
}

func TestTeamOfficeRoleBindingCandidateRequiresOwnerGateAndRollback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	binding := readJSON(t, filepath.Join(base, "bindings", "team-office-role-binding-candidate.json"))
	requireBool(t, binding, "candidate_only", true)
	requireBool(t, binding, "non_formal", true)
	if got := requireString(t, binding, "binding_truth_source"); got != "truzhenos" {
		t.Fatalf("binding_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, binding, "entry_surface"); got != "Team Settings Role Tab" {
		t.Fatalf("entry_surface = %s, want Team Settings Role Tab", got)
	}

	rawBindings, ok := binding["slot_bindings"].([]any)
	if !ok {
		t.Fatalf("slot_bindings missing")
	}
	if len(rawBindings) != 6 {
		t.Fatalf("slot_bindings len = %d, want 6", len(rawBindings))
	}
	requiredSlots := map[string]bool{
		"team_office.secretary_general":  false,
		"team_office.advisor.strategy":   false,
		"team_office.advisor.product":    false,
		"team_office.advisor.operations": false,
		"team_office.advisor.finance":    false,
		"team_office.advisor.legal_risk": false,
	}
	for _, raw := range rawBindings {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("slot binding = %T", raw)
		}
		slotID := requireString(t, item, "slot_ref")
		if _, ok := requiredSlots[slotID]; ok {
			requiredSlots[slotID] = true
		}
		if got := requireString(t, item, "binding_status"); got != "candidate_pending_owner" {
			t.Fatalf("%s binding_status = %s, want candidate_pending_owner", slotID, got)
		}
		if got := requireString(t, item, "owner_gate_required"); got != "true" {
			t.Fatalf("%s owner_gate_required = %s, want true", slotID, got)
		}
		if got := requireString(t, item, "history_policy"); got != "do_not_rewrite_existing_receipts" {
			t.Fatalf("%s history_policy = %s", slotID, got)
		}
	}
	for slotID, seen := range requiredSlots {
		if !seen {
			t.Fatalf("missing slot binding %s", slotID)
		}
	}

	rollback := requireObject(t, binding, "rollback_policy")
	if got := requireString(t, rollback, "rollback_creates_new_receipt"); got != "true" {
		t.Fatalf("rollback_creates_new_receipt = %s, want true", got)
	}
	if got := requireString(t, rollback, "history_policy"); got != "do_not_rewrite_existing_receipts" {
		t.Fatalf("rollback history_policy = %s", got)
	}

	negative := requireObject(t, binding, "negative_cases")
	for _, key := range []string{"replace_without_owner_gate", "incompatible_role_slot", "rewrite_existing_receipt", "capability_pack_changes_team_binding"} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "bindings/team-office-role-binding-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing bindings/team-office-role-binding-candidate.json")
}

func TestTeamOfficeFrontendBackendContractMapRequiresEvidenceForEverySurface(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contractMap := readJSON(t, filepath.Join(base, "integration", "frontend-backend-contract-map.json"))
	requireBool(t, contractMap, "candidate_only", true)
	requireBool(t, contractMap, "non_formal", true)
	if got := requireString(t, contractMap, "frontend_truth_source"); got != "truzhen-client-web-desktop" {
		t.Fatalf("frontend_truth_source = %s, want truzhen-client-web-desktop", got)
	}
	backendTruth := strings.Join(asStringSlice(t, contractMap["backend_truth_sources"]), "\n")
	for _, source := range []string{"truzhenos", "truzhen-cloud"} {
		if !strings.Contains(backendTruth, source) {
			t.Fatalf("backend_truth_sources missing %s", source)
		}
	}

	rawSurfaces, ok := contractMap["surfaces"].([]any)
	if !ok {
		t.Fatalf("surfaces missing")
	}
	requiredSurfaces := map[string]bool{
		"role_studio":                     false,
		"capability_role_reference":       false,
		"team_settings_role_tab":          false,
		"secretary_appearance":            false,
		"cloud_listing_purchase_download": false,
		"local_install_and_runtime_usage": false,
	}
	for _, raw := range rawSurfaces {
		surface, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("surface item = %T", raw)
		}
		surfaceID := requireString(t, surface, "surface_id")
		if _, ok := requiredSurfaces[surfaceID]; ok {
			requiredSurfaces[surfaceID] = true
		}
		requireStringIn(t, requireString(t, surface, "status"), "pending_cross_repo_execution")
		if len(asStringSlice(t, surface["frontend_fields"])) == 0 {
			t.Fatalf("%s frontend_fields must not be empty", surfaceID)
		}
		evidence := strings.Join(asStringSlice(t, surface["backend_evidence_required"]), "\n")
		if evidence == "" {
			t.Fatalf("%s backend_evidence_required must not be empty", surfaceID)
		}
		if surfaceID == "cloud_listing_purchase_download" {
			if !strings.Contains(evidence, "cloud_receipt_ref") {
				t.Fatalf("%s backend_evidence_required missing cloud_receipt_ref", surfaceID)
			}
		} else if !strings.Contains(evidence, "candidate_or_receipt_ref") {
			t.Fatalf("%s backend_evidence_required missing candidate_or_receipt_ref", surfaceID)
		}
	}
	for surfaceID, seen := range requiredSurfaces {
		if !seen {
			t.Fatalf("missing surface %s", surfaceID)
		}
	}

	forbidden := strings.Join(asStringSlice(t, contractMap["forbidden"]), "\n")
	for _, phrase := range []string{
		"frontend_only_success_state",
		"backend_success_without_gui_evidence",
		"cloud_payment_truth_in_truzhen_packs",
		"formal_binding_without_owner_gate",
	} {
		if !strings.Contains(forbidden, phrase) {
			t.Fatalf("forbidden missing %s", phrase)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "integration/frontend-backend-contract-map.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing integration/frontend-backend-contract-map.json")
}

func TestTeamOfficeGuiScriptP11StepsMapToCommercialAPITraceability(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contractMap := readJSON(t, filepath.Join(base, "integration", "frontend-backend-contract-map.json"))
	guiScript := readJSON(t, filepath.Join(base, "tests", "gui-user-agent-execution-script-candidate.json"))
	apiContract := readJSON(t, filepath.Join(base, "integration", "commercial-api-contract-candidate.json"))

	scriptSteps := map[string]map[string]any{}
	for _, step := range asObjectSlice(t, guiScript["ordered_steps"]) {
		scriptSteps[requireString(t, step, "step_id")] = step
	}
	apiEndpoints := map[string]map[string]any{}
	apiFailureStatuses := map[string]bool{}
	for _, endpoint := range asObjectSlice(t, apiContract["endpoints"]) {
		endpointID := requireString(t, endpoint, "endpoint_id")
		apiEndpoints[endpointID] = endpoint
		for _, status := range asStringSlice(t, endpoint["failure_statuses"]) {
			apiFailureStatuses[status] = true
		}
	}

	traceRows := map[string]map[string]any{}
	for _, row := range asObjectSlice(t, contractMap["gui_api_traceability_matrix"]) {
		traceRows[requireString(t, row, "script_step_id")] = row
	}

	required := map[string]struct {
		endpointRefs  []string
		receiptFields []string
		blockStatuses []string
	}{
		"export_candidate_bundle": {
			endpointRefs:  []string{"role_bundle_export"},
			receiptFields: []string{"bundle_export_receipt_ref", "bundle_tree_sha256", "artifact_sha256"},
		},
		"upload_cloud_listing_draft": {
			endpointRefs:  []string{"cloud_upload_draft"},
			receiptFields: []string{"cloud_upload_receipt_ref", "listing_draft_ref", "bundle_tree_sha256"},
		},
		"submit_marketplace_review_candidate": {
			endpointRefs:  []string{"marketplace_review_submit"},
			receiptFields: []string{"marketplace_review_candidate_receipt_ref", "production_publish_block_receipt_ref"},
			blockStatuses: []string{"blocked_production_publish_without_owner_authorization"},
		},
		"sandbox_purchase": {
			endpointRefs:  []string{"sandbox_order_create", "sandbox_payment_confirm", "entitlement_issue"},
			receiptFields: []string{"sandbox_order_ref", "sandbox_payment_receipt_ref", "entitlement_ref"},
			blockStatuses: []string{"blocked_real_payment_without_owner_authorization"},
		},
		"download_purchased_artifact": {
			endpointRefs:  []string{"signed_download_create"},
			receiptFields: []string{"download_receipt_ref", "signed_download_ref", "artifact_sha256"},
			blockStatuses: []string{"blocked_entitlement_missing"},
		},
		"install_downloaded_role_pack": {
			endpointRefs:  []string{"local_install_request"},
			receiptFields: []string{"install_receipt_ref", "entitlement_verification_ref", "artifact_hash_verification_ref", "enabled_role_pack_version"},
			blockStatuses: []string{"blocked_entitlement_missing", "blocked_artifact_hash_mismatch"},
		},
		"replace_team_roles_after_install": {
			endpointRefs:  []string{"team_settings_installed_roles_refresh", "team_role_binding_replace"},
			receiptFields: []string{"team_binding_receipt_ref", "owner_gate_decision_ref"},
			blockStatuses: []string{"blocked_owner_gate_required"},
		},
		"run_team_office_runtime_use": {
			receiptFields: []string{"runtime_receipt_candidate_ref", "secretary_candidate_output_ref", "five_advisor_candidate_output_refs"},
		},
		"run_negative_cases": {
			endpointRefs: []string{
				"marketplace_review_submit",
				"sandbox_payment_confirm",
				"signed_download_create",
				"local_install_request",
				"team_role_binding_replace",
			},
			receiptFields: []string{"blocked_receipt_ref", "blocked_reason", "actual_block_status"},
			blockStatuses: []string{
				"blocked_production_publish_without_owner_authorization",
				"blocked_real_payment_without_owner_authorization",
				"blocked_entitlement_missing",
				"blocked_artifact_hash_mismatch",
				"blocked_owner_gate_required",
			},
		},
	}

	for stepID, want := range required {
		scriptStep, ok := scriptSteps[stepID]
		if !ok {
			t.Fatalf("gui script missing step %s", stepID)
		}
		row, ok := traceRows[stepID]
		if !ok {
			t.Fatalf("gui_api_traceability_matrix missing %s", stepID)
		}
		if got, wantSurface := requireString(t, row, "surface_id"), requireString(t, scriptStep, "surface_id"); got != wantSurface {
			t.Fatalf("%s surface_id = %s, want %s", stepID, got, wantSurface)
		}
		if got := requireString(t, row, "frontend_event_id"); !strings.HasPrefix(got, "gui_event://") {
			t.Fatalf("%s frontend_event_id = %s, want gui_event://*", stepID, got)
		}
		requireBool(t, row, "gui_only", true)
		forbidden := strings.Join(asStringSlice(t, row["forbidden_execution_paths"]), "\n")
		for _, forbiddenPath := range []string{"direct_api_call", "manual_json_edit", "filesystem_copy_to_cloud"} {
			if !strings.Contains(forbidden, forbiddenPath) {
				t.Fatalf("%s forbidden_execution_paths missing %s", stepID, forbiddenPath)
			}
		}

		endpointRefs := asStringSlice(t, row["commercial_api_endpoint_refs"])
		for _, endpointID := range want.endpointRefs {
			requireStringSliceContains(t, endpointRefs, endpointID)
			if _, ok := apiEndpoints[endpointID]; !ok {
				t.Fatalf("%s references unknown commercial API endpoint %s", stepID, endpointID)
			}
		}
		if len(want.endpointRefs) == 0 {
			if actions := asStringSlice(t, row["backend_action_refs"]); len(actions) == 0 {
				t.Fatalf("%s must declare backend_action_refs when no commercial API endpoint applies", stepID)
			}
		}

		receiptFields := asStringSlice(t, row["required_receipt_or_result_fields"])
		for _, field := range want.receiptFields {
			requireStringSliceContains(t, receiptFields, field)
		}
		evidenceSlots := asStringSlice(t, row["required_gui_evidence_slots"])
		for _, slot := range asStringSlice(t, scriptStep["required_evidence_slots"]) {
			requireStringSliceContains(t, evidenceSlots, slot)
		}
		for _, status := range want.blockStatuses {
			requireStringSliceContains(t, asStringSlice(t, row["required_block_statuses"]), status)
			if !apiFailureStatuses[status] {
				t.Fatalf("%s required_block_statuses references status not declared by commercial API contract: %s", stepID, status)
			}
		}
	}
}

func TestTeamOfficeProductStageClosureRequiresGuiAPITraceabilityMatrix(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contractMap := readJSON(t, filepath.Join(base, "integration", "frontend-backend-contract-map.json"))
	report := readJSON(t, filepath.Join(base, "tests", "product-stage-frontend-backend-closure-report-candidate.json"))

	if rows := asObjectSlice(t, contractMap["gui_api_traceability_matrix"]); len(rows) == 0 {
		t.Fatalf("frontend-backend contract map gui_api_traceability_matrix must not be empty")
	}
	policy := requireObject(t, report, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, policy["required_before_completion_claim"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, policy["non_sufficient_evidence"]), "gui_api_traceability_matrix_missing")
	if !strings.Contains(requireString(t, report, "non_completion_clause"), "GUI/API traceability") {
		t.Fatalf("non_completion_clause must mention GUI/API traceability")
	}
}

func TestTeamOfficeCommercialAPIContractMapsProductActionsToReceipts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contract := readJSON(t, filepath.Join(base, "integration", "commercial-api-contract-candidate.json"))
	requireBool(t, contract, "candidate_only", true)
	requireBool(t, contract, "non_formal", true)
	if got := requireString(t, contract, "lifecycle_status"); got != "设计中" {
		t.Fatalf("lifecycle_status = %s, want 设计中", got)
	}

	truthSources := requireObject(t, contract, "truth_sources")
	if got := requireString(t, truthSources, "frontend_truth"); got != "truzhen-client-web-desktop" {
		t.Fatalf("frontend_truth = %s, want truzhen-client-web-desktop", got)
	}
	if got := requireString(t, truthSources, "cloud_truth"); got != "truzhen-cloud" {
		t.Fatalf("cloud_truth = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, truthSources, "install_truth"); got != "truzhenos" {
		t.Fatalf("install_truth = %s, want truzhenos", got)
	}

	policy := requireObject(t, contract, "operation_policy")
	requireBool(t, policy, "user_view_agents_gui_only", true)
	requireBool(t, policy, "receipt_ref_required_for_all_mutations", true)
	requireBool(t, policy, "idempotency_key_required_for_all_mutations", true)
	if got := requireString(t, policy, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	rawEndpoints, ok := contract["endpoints"].([]any)
	if !ok {
		t.Fatalf("endpoints missing")
	}
	requiredEndpoints := map[string]bool{
		"role_bundle_export":                    false,
		"cloud_upload_draft":                    false,
		"marketplace_review_submit":             false,
		"sandbox_order_create":                  false,
		"sandbox_payment_confirm":               false,
		"entitlement_issue":                     false,
		"signed_download_create":                false,
		"local_install_request":                 false,
		"team_settings_installed_roles_refresh": false,
		"team_role_binding_replace":             false,
	}
	for _, raw := range rawEndpoints {
		endpoint, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("endpoint = %T", raw)
		}
		endpointID := requireString(t, endpoint, "endpoint_id")
		if _, ok := requiredEndpoints[endpointID]; ok {
			requiredEndpoints[endpointID] = true
		}
		requireStringIn(t, requireString(t, endpoint, "status"), "pending_cross_repo_execution", "pending_owner_decision")
		requireStringIn(t, requireString(t, endpoint, "truth_source"), "truzhenos", "truzhen-cloud")
		requireStringIn(t, requireString(t, endpoint, "method"), "GET", "POST")
		if got := requireString(t, endpoint, "surface_id"); got == "" {
			t.Fatalf("%s surface_id missing", endpointID)
		}
		if got := requireString(t, endpoint, "path_candidate"); !strings.HasPrefix(got, "/candidate/") {
			t.Fatalf("%s path_candidate = %s, want /candidate/*", endpointID, got)
		}
		requestFields := strings.Join(asStringSlice(t, endpoint["request_fields"]), "\n")
		for _, field := range []string{"actor_ref", "correlation_id"} {
			if !strings.Contains(requestFields, field) {
				t.Fatalf("%s request_fields missing %s", endpointID, field)
			}
		}
		responseFields := strings.Join(asStringSlice(t, endpoint["response_fields"]), "\n")
		for _, field := range []string{"status", "receipt_ref"} {
			if !strings.Contains(responseFields, field) {
				t.Fatalf("%s response_fields missing %s", endpointID, field)
			}
		}
		if requireString(t, endpoint, "method") != "GET" {
			requireBool(t, endpoint, "idempotency_key_required", true)
			requireBool(t, endpoint, "receipt_ref_required", true)
		}
		failureStatuses := strings.Join(asStringSlice(t, endpoint["failure_statuses"]), "\n")
		if !strings.Contains(failureStatuses, "blocked") && !strings.Contains(failureStatuses, "not_ready") {
			t.Fatalf("%s failure_statuses must include blocked or not_ready", endpointID)
		}
	}
	for endpointID, seen := range requiredEndpoints {
		if !seen {
			t.Fatalf("missing endpoint %s", endpointID)
		}
	}

	rawHandoffs, ok := contract["handoff_contracts"].([]any)
	if !ok {
		t.Fatalf("handoff_contracts missing")
	}
	requiredHandoffs := map[string]bool{
		"bundle_export_to_cloud_upload": false,
		"upload_to_review":              false,
		"purchase_to_entitlement":       false,
		"entitlement_to_download":       false,
		"download_to_install":           false,
		"install_to_team_settings":      false,
	}
	for _, raw := range rawHandoffs {
		handoff, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("handoff = %T", raw)
		}
		handoffID := requireString(t, handoff, "handoff_id")
		if _, ok := requiredHandoffs[handoffID]; ok {
			requiredHandoffs[handoffID] = true
		}
		for _, key := range []string{"source_endpoint", "target_endpoint", "required_evidence"} {
			if got := requireString(t, handoff, key); got == "" {
				t.Fatalf("%s %s missing", handoffID, key)
			}
		}
		keys := strings.Join(asStringSlice(t, handoff["required_correlation_keys"]), "\n")
		if !strings.Contains(keys, "correlation_id") {
			t.Fatalf("%s required_correlation_keys missing correlation_id", handoffID)
		}
		if strings.Contains(handoffID, "download") || strings.Contains(handoffID, "install") {
			if !strings.Contains(keys, "artifact_sha256") {
				t.Fatalf("%s required_correlation_keys missing artifact_sha256", handoffID)
			}
		}
	}
	for handoffID, seen := range requiredHandoffs {
		if !seen {
			t.Fatalf("missing handoff %s", handoffID)
		}
	}

	forbidden := strings.Join(asStringSlice(t, contract["forbidden"]), "\n")
	for _, item := range []string{
		"real_payment_capture_without_owner_authorization",
		"production_publish_without_owner_authorization",
		"download_without_entitlement",
		"install_without_hash_verification",
		"frontend_success_without_receipt_ref",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_api_contract_candidate"); got != "integration/commercial-api-contract-candidate.json" {
		t.Fatalf("commercial_api_contract_candidate = %s", got)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, "integration/commercial-api-contract-candidate.json") {
		t.Fatalf("candidate set artifact_files missing commercial API contract")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "integration/commercial-api-contract-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing integration/commercial-api-contract-candidate.json")
}

func TestTeamOfficeCommercialAPIExampleCasesCoverEveryEndpoint(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	examplesPath := "integration/commercial-api-example-cases-candidate.json"
	contractPath := "integration/commercial-api-contract-candidate.json"
	contract := readJSON(t, filepath.Join(base, contractPath))

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_api_example_cases"); got != examplesPath {
		t.Fatalf("commercial_api_example_cases = %s, want %s", got, examplesPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, examplesPath) {
		t.Fatalf("candidate set artifact_files missing %s", examplesPath)
	}

	examples := readJSON(t, filepath.Join(base, examplesPath))
	requireBool(t, examples, "candidate_only", true)
	requireBool(t, examples, "non_formal", true)
	if got := requireString(t, examples, "commercial_api_contract_ref"); got != contractPath {
		t.Fatalf("commercial_api_contract_ref = %s, want %s", got, contractPath)
	}
	policy := requireObject(t, examples, "example_policy")
	requireBool(t, policy, "examples_are_non_executable", true)
	requireBool(t, policy, "refs_only_no_secret_payload", true)
	requireBool(t, policy, "user_view_gui_evidence_required_before_use", true)

	rawEndpoints, ok := contract["endpoints"].([]any)
	if !ok {
		t.Fatalf("contract endpoints missing")
	}
	endpointSpecs := map[string]map[string]any{}
	for _, raw := range rawEndpoints {
		endpoint, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("endpoint = %T", raw)
		}
		endpointSpecs[requireString(t, endpoint, "endpoint_id")] = endpoint
	}

	rawExamples, ok := examples["endpoint_examples"].([]any)
	if !ok {
		t.Fatalf("endpoint_examples missing")
	}
	if len(rawExamples) != len(endpointSpecs) {
		t.Fatalf("endpoint_examples len = %d, want endpoint len %d", len(rawExamples), len(endpointSpecs))
	}
	seen := map[string]bool{}
	for _, raw := range rawExamples {
		example, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("endpoint example = %T", raw)
		}
		endpointID := requireString(t, example, "endpoint_id")
		spec, ok := endpointSpecs[endpointID]
		if !ok {
			t.Fatalf("example for unknown endpoint %s", endpointID)
		}
		if seen[endpointID] {
			t.Fatalf("duplicate endpoint example %s", endpointID)
		}
		seen[endpointID] = true
		if got := requireString(t, example, "method"); got != requireString(t, spec, "method") {
			t.Fatalf("%s method = %s, want %s", endpointID, got, requireString(t, spec, "method"))
		}
		if got := requireString(t, example, "path_candidate"); got != requireString(t, spec, "path_candidate") {
			t.Fatalf("%s path_candidate = %s, want %s", endpointID, got, requireString(t, spec, "path_candidate"))
		}
		requestExample := requireObject(t, example, "request_example")
		for _, field := range asStringSlice(t, spec["request_fields"]) {
			if _, ok := requestExample[field]; !ok {
				t.Fatalf("%s request_example missing %s", endpointID, field)
			}
		}
		responseExample := requireObject(t, example, "success_response_example")
		for _, field := range asStringSlice(t, spec["response_fields"]) {
			if _, ok := responseExample[field]; !ok {
				t.Fatalf("%s success_response_example missing %s", endpointID, field)
			}
		}
		if _, ok := responseExample["receipt_ref"]; !ok {
			t.Fatalf("%s success_response_example missing receipt_ref", endpointID)
		}
		failureExample := requireObject(t, example, "failure_response_example")
		failureStatus := requireString(t, failureExample, "status")
		failureStatuses := strings.Join(asStringSlice(t, spec["failure_statuses"]), "\n")
		if !strings.Contains(failureStatuses, failureStatus) {
			t.Fatalf("%s failure status %s not in contract failure_statuses", endpointID, failureStatus)
		}
		for _, field := range []string{"status", "receipt_ref", "blocked_reason", "candidate_or_receipt_ref"} {
			if _, ok := failureExample[field]; !ok {
				t.Fatalf("%s failure_response_example missing %s", endpointID, field)
			}
		}
		serialized := fmt.Sprintf("%v", example)
		for _, forbidden := range []string{"sk_live_", "sk_test_", "X-Amz-Signature=", "-----BEGIN", "AKIA", "ghp_", "xoxb-"} {
			if strings.Contains(serialized, forbidden) {
				t.Fatalf("%s example contains forbidden secret-like payload %s", endpointID, forbidden)
			}
		}
	}
	for endpointID := range endpointSpecs {
		if !seen[endpointID] {
			t.Fatalf("missing endpoint example for %s", endpointID)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	inManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != examplesPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		inManifest = true
	}
	if !inManifest {
		t.Fatalf("artifact manifest missing %s", examplesPath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requiredProofs := strings.Join(asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "commercial_api_example_cases_verified") {
		t.Fatalf("product readiness matrix missing commercial_api_example_cases_verified")
	}
}

func TestTeamOfficeProductionPromotionAPIExamplesExposeP11PackageAndRuntimeSmokeFailures(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contract := readJSON(t, filepath.Join(base, "integration", "commercial-api-contract-candidate.json"))
	examples := readJSON(t, filepath.Join(base, "integration", "commercial-api-example-cases-candidate.json"))

	productionEndpoints := []string{
		"production_go_live_request",
		"real_payment_enable_request",
		"production_signed_download_enable",
		"production_listing_publish_request",
	}
	requiredFailureStatuses := []string{
		"blocked_p11_evidence_package_incomplete",
		"blocked_team_office_runtime_usage_smoke_missing",
	}

	contractEndpoints := map[string]map[string]any{}
	for _, endpoint := range asObjectSlice(t, contract["endpoints"]) {
		contractEndpoints[requireString(t, endpoint, "endpoint_id")] = endpoint
	}
	exampleEndpoints := map[string]map[string]any{}
	for _, example := range asObjectSlice(t, examples["endpoint_examples"]) {
		exampleEndpoints[requireString(t, example, "endpoint_id")] = example
	}

	for _, endpointID := range productionEndpoints {
		endpoint, ok := contractEndpoints[endpointID]
		if !ok {
			t.Fatalf("contract missing endpoint %s", endpointID)
		}
		failureStatuses := strings.Join(asStringSlice(t, endpoint["failure_statuses"]), "\n")
		for _, status := range requiredFailureStatuses {
			if !strings.Contains(failureStatuses, status) {
				t.Fatalf("%s failure_statuses missing %s", endpointID, status)
			}
		}

		example, ok := exampleEndpoints[endpointID]
		if !ok {
			t.Fatalf("examples missing endpoint %s", endpointID)
		}
		failureExamples := asObjectSlice(t, example["failure_response_examples"])
		seen := map[string]bool{}
		for _, failureExample := range failureExamples {
			status := requireString(t, failureExample, "status")
			seen[status] = true
			if !strings.Contains(failureStatuses, status) {
				t.Fatalf("%s failure example status %s not declared in contract", endpointID, status)
			}
			for _, field := range []string{"status", "receipt_ref", "blocked_reason", "candidate_or_receipt_ref"} {
				if _, ok := failureExample[field]; !ok {
					t.Fatalf("%s failure_response_examples status %s missing %s", endpointID, status, field)
				}
			}
			if status == "blocked_p11_evidence_package_incomplete" {
				reason := requireString(t, failureExample, "blocked_reason")
				if !strings.Contains(reason, "P11") || !strings.Contains(reason, "evidence package") {
					t.Fatalf("%s P11 failure blocked_reason = %s", endpointID, reason)
				}
			}
			if status == "blocked_team_office_runtime_usage_smoke_missing" {
				reason := requireString(t, failureExample, "blocked_reason")
				if !strings.Contains(reason, "runtime usage") || !strings.Contains(reason, "team_office_runtime_usage_smoke_verified") {
					t.Fatalf("%s runtime failure blocked_reason = %s", endpointID, reason)
				}
			}
		}
		for _, status := range requiredFailureStatuses {
			if !seen[status] {
				t.Fatalf("%s failure_response_examples missing %s", endpointID, status)
			}
		}
	}
}

func TestTeamOfficeCommercialAPIContractDefinesProductionPromotionEndpoints(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	contractPath := "integration/commercial-api-contract-candidate.json"
	contract := readJSON(t, filepath.Join(base, contractPath))
	gatePath := "commerce/commercial-production-promotion-gate-candidate.json"
	if got := requireString(t, contract, "source_commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("source_commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}

	rawEndpoints, ok := contract["endpoints"].([]any)
	if !ok {
		t.Fatalf("endpoints missing")
	}
	expected := map[string]struct {
		responseField string
	}{
		"production_go_live_request":         {"production_go_live_request_receipt_ref"},
		"real_payment_enable_request":        {"real_payment_enable_request_receipt_ref"},
		"production_signed_download_enable":  {"production_signed_download_enable_receipt_ref"},
		"production_listing_publish_request": {"production_listing_publish_request_receipt_ref"},
	}
	seen := map[string]bool{}
	for _, raw := range rawEndpoints {
		endpoint, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("endpoint = %T", raw)
		}
		endpointID := requireString(t, endpoint, "endpoint_id")
		want, ok := expected[endpointID]
		if !ok {
			continue
		}
		seen[endpointID] = true
		if got := requireString(t, endpoint, "status"); got != "pending_owner_decision" {
			t.Fatalf("%s status = %s, want pending_owner_decision", endpointID, got)
		}
		if got := requireString(t, endpoint, "truth_source"); got != "truzhen-cloud" {
			t.Fatalf("%s truth_source = %s, want truzhen-cloud", endpointID, got)
		}
		if got := requireString(t, endpoint, "surface_id"); got != "production_promotion" {
			t.Fatalf("%s surface_id = %s, want production_promotion", endpointID, got)
		}
		requireBool(t, endpoint, "idempotency_key_required", true)
		requireBool(t, endpoint, "receipt_ref_required", true)
		requestFields := strings.Join(asStringSlice(t, endpoint["request_fields"]), "\n")
		for _, field := range []string{
			"actor_ref",
			"correlation_id",
			"idempotency_key",
			"production_promotion_gate_ref",
			"owner_go_no_go_decision_ref",
			"p11_go_live_evidence_package_ref",
			"bundle_tree_sha256",
		} {
			if !strings.Contains(requestFields, field) {
				t.Fatalf("%s request_fields missing %s", endpointID, field)
			}
		}
		responseFields := strings.Join(asStringSlice(t, endpoint["response_fields"]), "\n")
		for _, field := range []string{"status", "receipt_ref", want.responseField, "blocked_or_allowed_result"} {
			if !strings.Contains(responseFields, field) {
				t.Fatalf("%s response_fields missing %s", endpointID, field)
			}
		}
		failureStatuses := strings.Join(asStringSlice(t, endpoint["failure_statuses"]), "\n")
		for _, status := range []string{
			"blocked_commercial_production_promotion_gate_missing",
			"blocked_owner_go_no_go_missing",
			"blocked_p11_evidence_package_incomplete",
		} {
			if !strings.Contains(failureStatuses, status) {
				t.Fatalf("%s failure_statuses missing %s", endpointID, status)
			}
		}
	}
	for endpointID := range expected {
		if !seen[endpointID] {
			t.Fatalf("missing production promotion endpoint %s", endpointID)
		}
	}

	rawHandoffs, ok := contract["handoff_contracts"].([]any)
	if !ok {
		t.Fatalf("handoff_contracts missing")
	}
	foundHandoff := false
	for _, raw := range rawHandoffs {
		handoff, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("handoff = %T", raw)
		}
		if requireString(t, handoff, "handoff_id") != "promotion_gate_to_production_controls" {
			continue
		}
		foundHandoff = true
		if got := requireString(t, handoff, "source_endpoint"); got != "production_go_live_request" {
			t.Fatalf("promotion handoff source_endpoint = %s", got)
		}
		if got := requireString(t, handoff, "target_endpoint"); got != "production_listing_publish_request" {
			t.Fatalf("promotion handoff target_endpoint = %s", got)
		}
		keys := strings.Join(asStringSlice(t, handoff["required_correlation_keys"]), "\n")
		for _, key := range []string{"correlation_id", "production_promotion_gate_ref", "owner_go_no_go_decision_ref", "bundle_tree_sha256"} {
			if !strings.Contains(keys, key) {
				t.Fatalf("promotion handoff required_correlation_keys missing %s", key)
			}
		}
	}
	if !foundHandoff {
		t.Fatalf("missing handoff promotion_gate_to_production_controls")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	for _, raw := range manifest["files"].([]any) {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") != contractPath {
			continue
		}
		data, err := os.ReadFile(filepath.Join(base, contractPath))
		if err != nil {
			t.Fatalf("read %s: %v", contractPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", contractPath, got, wantHash)
		}
		return
	}
	t.Fatalf("artifact manifest missing %s", contractPath)
}

func TestTeamOfficeProductionPromotionControlsAreQueuedAndLedgered(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	ledger := readJSON(t, filepath.Join(base, "docs", "commercial-cross-repo-evidence-ledger.json"))

	order := strings.Join(asStringSlice(t, queue["execution_order"]), "\n")
	if !strings.Contains(order, "production_promotion_controls") {
		t.Fatalf("execution_order missing production_promotion_controls")
	}

	expectedOutputs := []string{
		"production_go_live_request_receipt",
		"real_payment_enable_request_receipt",
		"production_signed_download_enable_receipt",
		"production_listing_publish_request_receipt",
		"production_install_observability_receipt",
	}
	rawEntries, ok := queue["execution_entries"].([]any)
	if !ok {
		t.Fatalf("execution_entries missing")
	}
	foundStage := false
	for _, raw := range rawEntries {
		entry, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution entry = %T", raw)
		}
		if requireString(t, entry, "stage_id") != "production_promotion_controls" {
			continue
		}
		foundStage = true
		if got := requireString(t, entry, "status"); got != "blocked_pending_prior_stage_evidence" {
			t.Fatalf("production_promotion_controls status = %s, want blocked_pending_prior_stage_evidence", got)
		}
		targets := strings.Join(asStringSlice(t, entry["target_repositories"]), "\n")
		for _, repo := range []string{"truzhen-cloud", "truzhenos", "truzhen-client-web-desktop"} {
			if !strings.Contains(targets, repo) {
				t.Fatalf("production_promotion_controls target_repositories missing %s", repo)
			}
		}
		inputs := strings.Join(asStringSlice(t, entry["required_input_refs"]), "\n")
		for _, ref := range []string{
			"integration/commercial-api-contract-candidate.json",
			"commerce/commercial-production-promotion-gate-candidate.json",
			"tests/p11-commercial-go-live-evidence-package-template.json",
		} {
			if !strings.Contains(inputs, ref) {
				t.Fatalf("production_promotion_controls required_input_refs missing %s", ref)
			}
		}
		outputs := strings.Join(asStringSlice(t, entry["required_evidence_outputs"]), "\n")
		for _, output := range expectedOutputs {
			if !strings.Contains(outputs, output) {
				t.Fatalf("production_promotion_controls required_evidence_outputs missing %s", output)
			}
		}
		writebacks := strings.Join(asStringSlice(t, entry["evidence_record_targets"]), "\n")
		for _, target := range []string{
			"docs/commercial-cross-repo-evidence-ledger.json",
			"tests/normal-commercialization-completion-audit-candidate.json",
			"tests/product-readiness-evidence-matrix.json",
		} {
			if !strings.Contains(writebacks, target) {
				t.Fatalf("production_promotion_controls evidence_record_targets missing %s", target)
			}
		}
		gate := requireObject(t, entry, "cross_repo_work_gate")
		requireBool(t, gate, "can_start_cross_repo_work", false)
		if got := requireString(t, gate, "required_status_before_cross_repo_work"); got != "owner_authorization_recorded_and_production_promotion_gate_passed" {
			t.Fatalf("production_promotion_controls required_status_before_cross_repo_work = %s", got)
		}
	}
	if !foundStage {
		t.Fatalf("missing production_promotion_controls execution entry")
	}

	phaseLinks, ok := ledger["p11_run_request_phase_links"].([]any)
	if !ok {
		t.Fatalf("p11_run_request_phase_links missing")
	}
	foundPhase := false
	for _, raw := range phaseLinks {
		phase, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("phase link = %T", raw)
		}
		if requireString(t, phase, "run_request_phase_id") != "production_promotion_controls" {
			continue
		}
		foundPhase = true
		if got := requireString(t, phase, "current_status"); got != "pending_authorization" {
			t.Fatalf("production_promotion_controls phase current_status = %s", got)
		}
		requireBool(t, phase, "can_count_toward_goal_completion", false)
		evidenceIDs := strings.Join(asStringSlice(t, phase["evidence_ids"]), "\n")
		if !strings.Contains(evidenceIDs, "role_studio_production_promotion_receipts") {
			t.Fatalf("production_promotion_controls phase missing role_studio_production_promotion_receipts")
		}
	}
	if !foundPhase {
		t.Fatalf("missing production_promotion_controls phase link")
	}

	rawRows, ok := ledger["evidence_rows"].([]any)
	if !ok {
		t.Fatalf("evidence_rows missing")
	}
	foundEvidence := false
	for _, raw := range rawRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("evidence row = %T", raw)
		}
		if requireString(t, row, "evidence_id") != "role_studio_production_promotion_receipts" {
			continue
		}
		foundEvidence = true
		if got := requireString(t, row, "stage_id"); got != "production_promotion_controls" {
			t.Fatalf("production promotion evidence stage_id = %s", got)
		}
		if got := requireString(t, row, "target_repository"); got != "multi_repo" {
			t.Fatalf("production promotion evidence target_repository = %s", got)
		}
		refs := strings.Join(asStringSlice(t, row["required_evidence_refs"]), "\n")
		for _, ref := range expectedOutputs {
			if !strings.Contains(refs, ref) {
				t.Fatalf("production promotion evidence required_evidence_refs missing %s", ref)
			}
		}
		blockers := strings.Join(asStringSlice(t, row["blocking_if_missing"]), "\n")
		for _, blocker := range []string{"production_go_live_request_receipt", "production_listing_publish_request_receipt"} {
			if !strings.Contains(blockers, blocker) {
				t.Fatalf("production promotion evidence blocking_if_missing missing %s", blocker)
			}
		}
		requireBool(t, row, "required_before_stage_complete", true)
		requireBool(t, row, "raw_payload_forbidden", true)
	}
	if !foundEvidence {
		t.Fatalf("missing role_studio_production_promotion_receipts evidence row")
	}
}

func TestTeamOfficeCommercializationExecutionPacketCoversPublisherBuyerAndNegativeEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packet := readJSON(t, filepath.Join(base, "tests", "commercialization-execution-packet.json"))
	requireBool(t, packet, "candidate_only", true)
	requireBool(t, packet, "non_formal", true)
	if got := requireString(t, packet, "completion_status"); got != "not_verified_requires_cross_repo_execution" {
		t.Fatalf("completion_status = %s, want not_verified_requires_cross_repo_execution", got)
	}

	actors := strings.Join(asStringSlice(t, packet["actors"]), "\n")
	for _, actor := range []string{"publisher_user_view_gui_agent", "buyer_user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent"} {
		if !strings.Contains(actors, actor) {
			t.Fatalf("actors missing %s", actor)
		}
	}

	rawStages, ok := packet["execution_stages"].([]any)
	if !ok {
		t.Fatalf("execution_stages missing")
	}
	requiredStages := map[string]bool{
		"publisher_create_role_candidates":      false,
		"publisher_export_candidate_bundle":     false,
		"publisher_upload_cloud_draft":          false,
		"cloud_review_candidate":                false,
		"buyer_sandbox_purchase":                false,
		"buyer_entitlement_and_signed_download": false,
		"buyer_local_install":                   false,
		"buyer_team_settings_replace_and_use":   false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		requireStringIn(t, requireString(t, stage, "status"), "pending_cross_repo_execution")
		truthSource := requireString(t, stage, "truth_source")
		requireStringIn(t, truthSource, "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo")
		evidence := strings.Join(asStringSlice(t, stage["required_evidence"]), "\n")
		for _, proof := range []string{"gui_screenshot", "page_state", "receipt_or_candidate_ref"} {
			if !strings.Contains(evidence, proof) {
				t.Fatalf("%s required_evidence missing %s", stageID, proof)
			}
		}
		if strings.Contains(stageID, "download") && !strings.Contains(evidence, "artifact_hash_match") {
			t.Fatalf("%s required_evidence missing artifact_hash_match", stageID)
		}
		if strings.Contains(stageID, "install") && !strings.Contains(evidence, "install_receipt") {
			t.Fatalf("%s required_evidence missing install_receipt", stageID)
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing execution stage %s", stageID)
		}
	}

	negative := requireObject(t, packet, "negative_cases")
	for _, key := range []string{
		"download_without_purchase",
		"expired_entitlement_install",
		"tampered_artifact_hash",
		"team_binding_without_owner_gate",
		"raw_voice_or_vrm_asset",
		"real_payment_without_owner_authorization",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "required_evidence"); got == "" {
			t.Fatalf("%s required_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/commercialization-execution-packet.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/commercialization-execution-packet.json")
}

func TestTeamOfficeE2EEvidenceRunRecordRequiresEvidenceBeforeCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	record := readJSON(t, filepath.Join(base, "tests", "e2e-evidence-run-record.json"))
	requireBool(t, record, "candidate_only", true)
	requireBool(t, record, "non_formal", true)
	if got := requireString(t, record, "run_status"); got != "not_run_requires_cross_repo_execution" {
		t.Fatalf("run_status = %s, want not_run_requires_cross_repo_execution", got)
	}

	authorization := requireObject(t, record, "cross_repo_authorization")
	for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "truzhen-contracts"} {
		if got := requireString(t, authorization, repo); got != "requires_owner_authorization" {
			t.Fatalf("%s authorization = %s, want requires_owner_authorization", repo, got)
		}
	}

	rawStages, ok := record["stage_records"].([]any)
	if !ok {
		t.Fatalf("stage_records missing")
	}
	requiredStages := map[string]bool{
		"publisher_create_role_candidates":      false,
		"publisher_export_candidate_bundle":     false,
		"publisher_upload_cloud_draft":          false,
		"cloud_review_candidate":                false,
		"buyer_sandbox_purchase":                false,
		"buyer_entitlement_and_signed_download": false,
		"buyer_local_install":                   false,
		"buyer_team_settings_replace_and_use":   false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage record = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		if got := requireString(t, stage, "status"); got != "not_run" {
			t.Fatalf("%s status = %s, want not_run", stageID, got)
		}
		slots := strings.Join(asStringSlice(t, stage["required_evidence_slots"]), "\n")
		for _, slot := range []string{"gui_screenshot_path", "page_state_ref", "receipt_or_candidate_ref", "timestamp", "truth_source"} {
			if !strings.Contains(slots, slot) {
				t.Fatalf("%s required_evidence_slots missing %s", stageID, slot)
			}
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing stage record %s", stageID)
		}
	}

	rawNegative, ok := record["negative_case_records"].([]any)
	if !ok {
		t.Fatalf("negative_case_records missing")
	}
	requiredNegative := map[string]bool{
		"download_without_purchase":                false,
		"expired_entitlement_install":              false,
		"tampered_artifact_hash":                   false,
		"team_binding_without_owner_gate":          false,
		"raw_voice_or_vrm_asset":                   false,
		"real_payment_without_owner_authorization": false,
	}
	for _, raw := range rawNegative {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("negative record = %T", raw)
		}
		caseID := requireString(t, item, "case_id")
		if _, ok := requiredNegative[caseID]; ok {
			requiredNegative[caseID] = true
		}
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", caseID, got)
		}
		slots := strings.Join(asStringSlice(t, item["required_evidence_slots"]), "\n")
		for _, slot := range []string{"gui_screenshot_path", "blocked_reason", "receipt_or_candidate_ref"} {
			if !strings.Contains(slots, slot) {
				t.Fatalf("%s required_evidence_slots missing %s", caseID, slot)
			}
		}
	}
	for caseID, seen := range requiredNegative {
		if !seen {
			t.Fatalf("missing negative case record %s", caseID)
		}
	}

	completion := requireObject(t, record, "completion_gate")
	requireBool(t, completion, "completion_claim_allowed", false)
	requiredBeforeCompletion := strings.Join(asStringSlice(t, completion["required_before_completion"]), "\n")
	for _, proof := range []string{"all_stage_records_passed", "all_negative_case_records_blocked", "artifact_hash_match", "no_secret_scan", "owner_gate_evidence"} {
		if !strings.Contains(requiredBeforeCompletion, proof) {
			t.Fatalf("completion gate missing %s", proof)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "tests/e2e-evidence-run-record.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing tests/e2e-evidence-run-record.json")
}

func TestTeamOfficeCrossRepoExecutionCardsDefineAuthorizedProductWorkstreams(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	cards := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-cards.json"))
	requireBool(t, cards, "candidate_only", true)
	requireBool(t, cards, "non_formal", true)
	if got := requireString(t, cards, "execution_status"); got != "pending_owner_authorization" {
		t.Fatalf("execution_status = %s, want pending_owner_authorization", got)
	}

	rawCards, ok := cards["execution_cards"].([]any)
	if !ok {
		t.Fatalf("execution_cards missing")
	}
	requiredCards := map[string]bool{
		"contracts_schema_impact_review":         false,
		"backend_role_candidate_gate_receipt":    false,
		"frontend_role_studio_and_team_settings": false,
		"cloud_upload_purchase_download":         false,
		"local_install_and_runtime_binding":      false,
		"gui_e2e_acceptance_recording":           false,
	}
	requiredRepos := map[string]bool{
		"truzhen-contracts":          false,
		"truzhenos":                  false,
		"truzhen-client-web-desktop": false,
		"truzhen-cloud":              false,
		"multi_repo":                 false,
	}
	for _, raw := range rawCards {
		card, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution card = %T", raw)
		}
		cardID := requireString(t, card, "card_id")
		if _, ok := requiredCards[cardID]; ok {
			requiredCards[cardID] = true
		}
		repo := requireString(t, card, "target_repo")
		if _, ok := requiredRepos[repo]; ok {
			requiredRepos[repo] = true
		}
		if got := requireString(t, card, "authorization_status"); got != "requires_owner_authorization" {
			t.Fatalf("%s authorization_status = %s, want requires_owner_authorization", cardID, got)
		}
		requireStringIn(t, requireString(t, card, "risk_color"), "黄", "橙", "红")
		if len(asStringSlice(t, card["implementation_targets"])) == 0 {
			t.Fatalf("%s implementation_targets must not be empty", cardID)
		}
		if len(asStringSlice(t, card["verification_commands"])) == 0 {
			t.Fatalf("%s verification_commands must not be empty", cardID)
		}
		evidence := strings.Join(asStringSlice(t, card["required_evidence"]), "\n")
		if !strings.Contains(evidence, "gui_screenshot") && !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "schema_impact") {
			t.Fatalf("%s required_evidence must include GUI, receipt, or schema impact evidence", cardID)
		}
	}
	for cardID, seen := range requiredCards {
		if !seen {
			t.Fatalf("missing execution card %s", cardID)
		}
	}
	for repo, seen := range requiredRepos {
		if !seen {
			t.Fatalf("missing target repo %s", repo)
		}
	}

	forbidden := strings.Join(asStringSlice(t, cards["forbidden"]), "\n")
	for _, phrase := range []string{
		"no_cross_repo_write_without_owner_authorization",
		"no_real_payment_without_owner_authorization",
		"no_production_publish_without_owner_authorization",
		"no_formal_binding_without_owner_gate",
	} {
		if !strings.Contains(forbidden, phrase) {
			t.Fatalf("forbidden missing %s", phrase)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "integration/cross-repo-execution-cards.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing integration/cross-repo-execution-cards.json")
}

func TestTeamOfficeCrossRepoExecutionReadinessPackageDefinesAuthorizationReadyHandoff(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "integration/cross-repo-execution-readiness-package.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))
	requireBool(t, pkg, "candidate_only", true)
	requireBool(t, pkg, "non_formal", true)
	if got := requireString(t, pkg, "readiness_status"); got != "ready_for_owner_authorization_not_executed" {
		t.Fatalf("readiness_status = %s, want ready_for_owner_authorization_not_executed", got)
	}

	sourceRefs := strings.Join(asStringSlice(t, pkg["source_refs"]), "\n")
	for _, ref := range []string{
		"integration/cross-repo-execution-cards.json",
		"tests/normal-commercialization-completion-audit-candidate.json",
		"tests/e2e-evidence-run-record.json",
		"commerce/commercial-go-live-approval-candidate.json",
	} {
		if !strings.Contains(sourceRefs, ref) {
			t.Fatalf("source_refs missing %s", ref)
		}
	}

	rawRepos, ok := pkg["target_repositories"].([]any)
	if !ok {
		t.Fatalf("target_repositories missing")
	}
	requiredRepos := map[string]string{
		"truzhen-contracts":          "/Users/li/Documents/truzhen-contracts",
		"truzhenos":                  "/Users/li/Documents/truzhenos",
		"truzhen-client-web-desktop": "/Users/li/Documents/truzhen-client-web-desktop",
		"truzhen-cloud":              "/Users/li/Documents/truzhen-cloud",
		"truzhen-packs-current":      "/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan",
	}
	for _, raw := range rawRepos {
		repo, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("target repository = %T", raw)
		}
		repoID := requireString(t, repo, "repo_id")
		if wantPath, ok := requiredRepos[repoID]; ok {
			if got := requireString(t, repo, "repo_path"); got != wantPath {
				t.Fatalf("%s repo_path = %s, want %s", repoID, got, wantPath)
			}
			requiredRepos[repoID] = ""
		}
		if got := requireString(t, repo, "authorization_status"); got != "requires_owner_authorization" && repoID != "truzhen-packs-current" {
			t.Fatalf("%s authorization_status = %s, want requires_owner_authorization", repoID, got)
		}
		if got := requireString(t, repo, "required_status_command"); got != "git status --short --branch" {
			t.Fatalf("%s required_status_command = %s", repoID, got)
		}
		if len(asStringSlice(t, repo["allowed_actions_after_authorization"])) == 0 {
			t.Fatalf("%s allowed_actions_after_authorization missing", repoID)
		}
		evidence := strings.Join(asStringSlice(t, repo["required_evidence_outputs"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "screenshot") && !strings.Contains(evidence, "impact_report") {
			t.Fatalf("%s required_evidence_outputs must include receipt, screenshot, or impact_report", repoID)
		}
		forbidden := strings.Join(asStringSlice(t, repo["forbidden_boundaries"]), "\n")
		if forbidden == "" {
			t.Fatalf("%s forbidden_boundaries missing", repoID)
		}
	}
	for repoID, missing := range requiredRepos {
		if missing != "" {
			t.Fatalf("missing target repository %s", repoID)
		}
	}

	rawStages, ok := pkg["execution_stage_order"].([]any)
	if !ok {
		t.Fatalf("execution_stage_order missing")
	}
	requiredStages := map[string]bool{
		"contracts_schema_review":        false,
		"backend_candidate_gate_receipt": false,
		"frontend_user_gui_flow":         false,
		"cloud_sandbox_commerce":         false,
		"local_install_binding":          false,
		"gui_e2e_acceptance":             false,
		"independent_acceptance_review":  false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		requireStringIn(t, requireString(t, stage, "status"), "pending_owner_authorization", "blocked_until_prior_stage_evidence")
		if got := requireString(t, stage, "entry_condition"); got == "" {
			t.Fatalf("%s entry_condition missing", stageID)
		}
		if got := requireString(t, stage, "completion_evidence"); got == "" {
			t.Fatalf("%s completion_evidence missing", stageID)
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing execution stage %s", stageID)
		}
	}

	questions := strings.Join(asStringSlice(t, pkg["owner_authorization_questions"]), "\n")
	for _, item := range []string{
		"truzhen-contracts",
		"truzhenos",
		"truzhen-client-web-desktop",
		"truzhen-cloud",
		"sandbox payment",
		"real payment remains blocked",
		"production publish remains blocked",
	} {
		if !strings.Contains(questions, item) {
			t.Fatalf("owner_authorization_questions missing %s", item)
		}
	}

	completion := requireObject(t, pkg, "completion_update_policy")
	requireBool(t, completion, "may_mark_goal_complete_from_this_package", false)
	if got := requireString(t, completion, "required_before_goal_completion"); got == "" {
		t.Fatalf("completion_update_policy.required_before_goal_completion missing")
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "cross_repo_execution_readiness_package"); got != packagePath {
		t.Fatalf("cross_repo_execution_readiness_package = %s, want %s", got, packagePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, packagePath) {
		t.Fatalf("candidate set artifact_files missing %s", packagePath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == packagePath {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", packagePath)
}

func TestTeamOfficeCrossRepoExecutionHandoffIncludesProductionPromotionControls(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	cards := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-cards.json"))
	readiness := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-readiness-package.json"))

	expectedEvidence := []string{
		"production_go_live_request_receipt",
		"real_payment_enable_request_receipt",
		"production_signed_download_enable_receipt",
		"production_listing_publish_request_receipt",
		"production_install_observability_receipt",
	}

	rawCards, ok := cards["execution_cards"].([]any)
	if !ok {
		t.Fatalf("execution_cards missing")
	}
	foundCard := false
	for _, raw := range rawCards {
		card, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution card = %T", raw)
		}
		if requireString(t, card, "card_id") != "production_promotion_controls" {
			continue
		}
		foundCard = true
		if got := requireString(t, card, "target_repo"); got != "multi_repo" {
			t.Fatalf("production_promotion_controls target_repo = %s, want multi_repo", got)
		}
		if got := requireString(t, card, "authorization_status"); got != "requires_owner_authorization" {
			t.Fatalf("production_promotion_controls authorization_status = %s", got)
		}
		requireStringIn(t, requireString(t, card, "risk_color"), "红")
		targets := strings.Join(asStringSlice(t, card["implementation_targets"]), "\n")
		for _, target := range []string{
			"production_go_live_request",
			"real_payment_enable_request",
			"production_signed_download_enable",
			"production_listing_publish_request",
			"production_install_observability",
		} {
			if !strings.Contains(targets, target) {
				t.Fatalf("production_promotion_controls implementation_targets missing %s", target)
			}
		}
		evidence := strings.Join(asStringSlice(t, card["required_evidence"]), "\n")
		for _, item := range expectedEvidence {
			if !strings.Contains(evidence, item) {
				t.Fatalf("production_promotion_controls required_evidence missing %s", item)
			}
		}
		forbidden := strings.Join(asStringSlice(t, card["forbidden_boundaries"]), "\n")
		for _, item := range []string{
			"no_production_publish_without_owner_authorization",
			"no_real_payment_without_owner_authorization",
			"no_copy_sandbox_receipt_as_production_receipt",
		} {
			if !strings.Contains(forbidden, item) {
				t.Fatalf("production_promotion_controls forbidden_boundaries missing %s", item)
			}
		}
	}
	if !foundCard {
		t.Fatalf("missing execution card production_promotion_controls")
	}

	rawStages, ok := readiness["execution_stage_order"].([]any)
	if !ok {
		t.Fatalf("execution_stage_order missing")
	}
	foundStage := false
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution stage = %T", raw)
		}
		if requireString(t, stage, "stage_id") != "production_promotion_controls" {
			continue
		}
		foundStage = true
		if got := requireString(t, stage, "status"); got != "blocked_until_prior_stage_evidence" {
			t.Fatalf("production_promotion_controls status = %s", got)
		}
		if got := requireString(t, stage, "target_repo"); got != "multi_repo" {
			t.Fatalf("production_promotion_controls target_repo = %s", got)
		}
		entry := requireString(t, stage, "entry_condition")
		for _, item := range []string{
			"commercial-production-promotion-gate",
			"Owner go/no-go",
			"independent acceptance",
		} {
			if !strings.Contains(entry, item) {
				t.Fatalf("production_promotion_controls entry_condition missing %s", item)
			}
		}
		completion := requireString(t, stage, "completion_evidence")
		for _, item := range expectedEvidence {
			if !strings.Contains(completion, item) {
				t.Fatalf("production_promotion_controls completion_evidence missing %s", item)
			}
		}
	}
	if !foundStage {
		t.Fatalf("missing readiness execution stage production_promotion_controls")
	}

	completionPolicy := requireObject(t, readiness, "completion_update_policy")
	requiredBeforeGoal := requireString(t, completionPolicy, "required_before_goal_completion")
	for _, item := range []string{
		"production_go_live_request_receipt",
		"production_listing_publish_request_receipt",
		"real_payment_enable_request_receipt",
	} {
		if !strings.Contains(requiredBeforeGoal, item) {
			t.Fatalf("completion_update_policy.required_before_goal_completion missing %s", item)
		}
	}
}

func TestTeamOfficeInstallCompatibilityMatrixCoversDownloadInstallUpgradeAndRollback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrix := readJSON(t, filepath.Join(base, "commerce", "install-compatibility-matrix.json"))
	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	if got := requireString(t, matrix, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, matrix, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}

	targets := strings.Join(asStringSlice(t, matrix["supported_install_targets"]), "\n")
	for _, target := range []string{"local_pack_manager", "team_settings_role_tab", "cloud_purchased_page_trigger"} {
		if !strings.Contains(targets, target) {
			t.Fatalf("supported_install_targets missing %s", target)
		}
	}

	rawDimensions, ok := matrix["compatibility_dimensions"].([]any)
	if !ok {
		t.Fatalf("compatibility_dimensions missing")
	}
	requiredDimensions := map[string]bool{
		"truzhenos_runtime_version": false,
		"role_pack_schema_version":  false,
		"team_office_role_slots":    false,
		"entitlement_scope":         false,
		"artifact_signature_hash":   false,
		"frontend_route_version":    false,
	}
	for _, raw := range rawDimensions {
		dimension, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("compatibility dimension = %T", raw)
		}
		dimensionID := requireString(t, dimension, "dimension_id")
		if _, ok := requiredDimensions[dimensionID]; ok {
			requiredDimensions[dimensionID] = true
		}
		if got := requireString(t, dimension, "check_required"); got != "true" {
			t.Fatalf("%s check_required = %s, want true", dimensionID, got)
		}
		if got := requireString(t, dimension, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", dimensionID)
		}
	}
	for dimensionID, seen := range requiredDimensions {
		if !seen {
			t.Fatalf("missing compatibility dimension %s", dimensionID)
		}
	}

	rawPaths, ok := matrix["install_paths"].([]any)
	if !ok {
		t.Fatalf("install_paths missing")
	}
	requiredPaths := map[string]bool{
		"fresh_install_after_purchase": false,
		"upgrade_existing_role_pack":   false,
		"rollback_to_previous_binding": false,
		"reinstall_after_reenable":     false,
	}
	for _, raw := range rawPaths {
		path, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("install path = %T", raw)
		}
		pathID := requireString(t, path, "path_id")
		if _, ok := requiredPaths[pathID]; ok {
			requiredPaths[pathID] = true
		}
		evidence := strings.Join(asStringSlice(t, path["required_evidence"]), "\n")
		for _, proof := range []string{"install_receipt", "artifact_hash_match", "owner_gate_evidence"} {
			if !strings.Contains(evidence, proof) {
				t.Fatalf("%s required_evidence missing %s", pathID, proof)
			}
		}
	}
	for pathID, seen := range requiredPaths {
		if !seen {
			t.Fatalf("missing install path %s", pathID)
		}
	}

	negative := requireObject(t, matrix, "negative_cases")
	for _, key := range []string{
		"missing_entitlement",
		"expired_entitlement",
		"incompatible_schema",
		"missing_team_slot",
		"unsupported_truzhenos_version",
		"artifact_hash_mismatch",
		"signature_missing",
		"production_publish_without_owner_authorization",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/install-compatibility-matrix.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/install-compatibility-matrix.json")
}

func TestTeamOfficeMarketplaceReviewSubmissionDefinesCommercialReviewEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	requireBool(t, submission, "candidate_only", true)
	requireBool(t, submission, "non_formal", true)
	if got := requireString(t, submission, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, submission, "review_status"); got != "candidate_not_submitted" {
		t.Fatalf("review_status = %s, want candidate_not_submitted", got)
	}
	if got := requireString(t, submission, "production_publish_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("production_publish_policy = %s, want blocked_until_owner_authorizes", got)
	}

	listing := requireObject(t, submission, "listing_metadata")
	for _, key := range []string{"display_name", "short_description", "category", "artifact_ref", "artifact_manifest", "risk_statement"} {
		if got := requireString(t, listing, key); got == "" {
			t.Fatalf("listing_metadata.%s missing", key)
		}
	}

	pricing := requireObject(t, submission, "pricing_and_license")
	for _, key := range []string{"pricing_status", "currency", "license_scope", "refund_policy", "real_payment_policy"} {
		if got := requireString(t, pricing, key); got == "" {
			t.Fatalf("pricing_and_license.%s missing", key)
		}
	}
	if got := requireString(t, pricing, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	rights := requireObject(t, submission, "asset_rights_and_safety")
	for _, key := range []string{"voice_asset_policy", "vrm_asset_policy", "raw_asset_policy", "privacy_policy", "secret_policy"} {
		if got := requireString(t, rights, key); got == "" {
			t.Fatalf("asset_rights_and_safety.%s missing", key)
		}
	}

	support := requireObject(t, submission, "support_and_operations")
	for _, key := range []string{"support_status", "support_channel_policy", "version_support_policy", "refund_handling_policy", "revocation_notice_policy"} {
		if got := requireString(t, support, key); got == "" {
			t.Fatalf("support_and_operations.%s missing", key)
		}
	}

	rawChecklist, ok := submission["review_checklist"].([]any)
	if !ok {
		t.Fatalf("review_checklist missing")
	}
	requiredChecklist := map[string]bool{
		"artifact_manifest_hashes":    false,
		"role_pack_candidate_only":    false,
		"owner_gate_boundaries":       false,
		"license_entitlement_policy":  false,
		"asset_rights_asset_ref_only": false,
		"support_refund_revocation":   false,
		"negative_case_evidence_plan": false,
	}
	for _, raw := range rawChecklist {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("review checklist item = %T", raw)
		}
		checkID := requireString(t, item, "check_id")
		if _, ok := requiredChecklist[checkID]; ok {
			requiredChecklist[checkID] = true
		}
		requireStringIn(t, requireString(t, item, "status"), "candidate_ready", "pending_owner_decision")
		if got := requireString(t, item, "required_evidence"); got == "" {
			t.Fatalf("%s required_evidence missing", checkID)
		}
	}
	for checkID, seen := range requiredChecklist {
		if !seen {
			t.Fatalf("missing review checklist %s", checkID)
		}
	}

	negative := requireObject(t, submission, "negative_cases")
	for _, key := range []string{
		"submit_review_without_artifact_hash",
		"publish_without_owner_authorization",
		"real_payment_without_owner_authorization",
		"raw_voice_or_vrm_asset",
		"missing_refund_policy",
		"missing_support_policy",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/marketplace-review-submission-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/marketplace-review-submission-candidate.json")
}

func TestTeamOfficePublisherAccountSettlementPolicyBlocksCommercialListingWithoutCloudSellerControls(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "publisher-account-settlement-policy-candidate.json"))
	requireBool(t, policy, "candidate_only", true)
	requireBool(t, policy, "non_formal", true)
	if got := requireString(t, policy, "cloud_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "package_contents_truth_source"); got != "truzhen-packs" {
		t.Fatalf("package_contents_truth_source = %s, want truzhen-packs", got)
	}
	if got := requireString(t, policy, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	identity := requireObject(t, policy, "publisher_identity_policy")
	requireBool(t, identity, "publisher_account_required", true)
	requireBool(t, identity, "sandbox_identity_allowed", true)
	requireBool(t, identity, "no_identity_document_in_packs", true)
	for _, key := range []string{"publisher_identity_truth_source", "kyc_status_policy", "evidence_required"} {
		if got := requireString(t, identity, key); got == "" {
			t.Fatalf("publisher_identity_policy.%s missing", key)
		}
	}

	pricing := requireObject(t, policy, "pricing_approval_policy")
	requireBool(t, pricing, "price_candidate_required", true)
	requireBool(t, pricing, "owner_price_approval_required", true)
	for _, key := range []string{"currency", "pricing_truth_source", "real_price_status", "evidence_required"} {
		if got := requireString(t, pricing, key); got == "" {
			t.Fatalf("pricing_approval_policy.%s missing", key)
		}
	}

	settlement := requireObject(t, policy, "settlement_policy")
	requireBool(t, settlement, "settlement_receipt_required", true)
	requireBool(t, settlement, "no_bank_account_in_packs", true)
	for _, key := range []string{"payout_truth_source", "sandbox_payout_policy", "real_payout_policy", "evidence_required"} {
		if got := requireString(t, settlement, key); got == "" {
			t.Fatalf("settlement_policy.%s missing", key)
		}
	}

	taxInvoice := requireObject(t, policy, "tax_invoice_policy")
	requireBool(t, taxInvoice, "tax_profile_required_before_real_payout", true)
	requireBool(t, taxInvoice, "tax_profile_in_packs", false)
	for _, key := range []string{"invoice_truth_source", "buyer_invoice_policy", "evidence_required"} {
		if got := requireString(t, taxInvoice, key); got == "" {
			t.Fatalf("tax_invoice_policy.%s missing", key)
		}
	}

	rawGates, ok := policy["commercial_go_live_gates"].([]any)
	if !ok {
		t.Fatalf("commercial_go_live_gates missing")
	}
	requiredGates := map[string]bool{
		"publisher_identity_verified":        false,
		"pricing_owner_approved":             false,
		"tax_invoice_policy_recorded":        false,
		"settlement_policy_recorded":         false,
		"support_refund_policy_ready":        false,
		"real_payment_authorized_or_blocked": false,
	}
	for _, raw := range rawGates {
		gate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("go-live gate = %T", raw)
		}
		gateID := requireString(t, gate, "gate_id")
		if _, ok := requiredGates[gateID]; ok {
			requiredGates[gateID] = true
		}
		if got := requireString(t, gate, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs", gateID)
		}
		requireStringIn(t, requireString(t, gate, "status"), "candidate_ready", "pending_owner_decision")
		if got := requireString(t, gate, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", gateID)
		}
	}
	for gateID, seen := range requiredGates {
		if !seen {
			t.Fatalf("missing commercial go-live gate %s", gateID)
		}
	}

	negative := requireObject(t, policy, "negative_cases")
	for _, key := range []string{
		"publish_without_publisher_identity",
		"real_payout_without_tax_profile",
		"bank_account_in_pack",
		"price_live_without_owner_approval",
		"invoice_truth_in_packs",
		"settlement_truth_in_packs",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, policy["forbidden"]), "\n")
	for _, item := range []string{
		"store_bank_account_in_truzhen_packs",
		"store_tax_profile_truth_in_truzhen_packs",
		"store_settlement_truth_in_truzhen_packs",
		"real_payout_without_owner_authorization",
		"publish_without_publisher_identity_verification",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	rawChecklist, ok := submission["review_checklist"].([]any)
	if !ok {
		t.Fatalf("review_checklist missing")
	}
	for _, raw := range rawChecklist {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("review checklist item = %T", raw)
		}
		if requireString(t, item, "check_id") == "publisher_identity_settlement" {
			requireStringIn(t, requireString(t, item, "status"), "candidate_ready", "pending_owner_decision")
			if got := requireString(t, item, "required_evidence"); !strings.Contains(got, "publisher-account-settlement-policy-candidate.json") {
				t.Fatalf("publisher_identity_settlement required_evidence = %s", got)
			}
			goto manifestCheck
		}
	}
	t.Fatalf("marketplace review missing publisher_identity_settlement checklist")

manifestCheck:
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/publisher-account-settlement-policy-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/publisher-account-settlement-policy-candidate.json")
}

func TestTeamOfficeCommercialTermsPrivacyPolicyDisclosesCandidateRoleAndDataBoundaries(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "commercial-terms-privacy-policy-candidate.json"))
	requireBool(t, policy, "candidate_only", true)
	requireBool(t, policy, "non_formal", true)
	if got := requireString(t, policy, "cloud_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "package_contents_truth_source"); got != "truzhen-packs" {
		t.Fatalf("package_contents_truth_source = %s, want truzhen-packs", got)
	}
	if got := requireString(t, policy, "runtime_truth_source"); got != "truzhenos" {
		t.Fatalf("runtime_truth_source = %s, want truzhenos", got)
	}

	terms := requireObject(t, policy, "terms_policy")
	requireBool(t, terms, "terms_surface_required", true)
	requireBool(t, terms, "buyer_acceptance_required_before_purchase", true)
	requireBool(t, terms, "role_output_candidate_disclaimer_required", true)
	for _, key := range []string{"license_terms_surface", "acceptable_use_surface", "evidence_required"} {
		if got := requireString(t, terms, key); got == "" {
			t.Fatalf("terms_policy.%s missing", key)
		}
	}

	privacy := requireObject(t, policy, "privacy_policy")
	requireBool(t, privacy, "privacy_notice_required", true)
	requireBool(t, privacy, "no_customer_runtime_data_in_pack", true)
	requireBool(t, privacy, "no_cloud_order_personal_data_in_pack", true)
	for _, key := range []string{"privacy_surface", "data_controller_truth_source", "evidence_required"} {
		if got := requireString(t, privacy, key); got == "" {
			t.Fatalf("privacy_policy.%s missing", key)
		}
	}

	dataHandling := requireObject(t, policy, "data_handling_policy")
	for _, key := range []string{"pack_content_truth_source", "runtime_data_truth_source", "cloud_order_data_truth_source", "telemetry_policy", "deletion_export_request_route", "evidence_required"} {
		if got := requireString(t, dataHandling, key); got == "" {
			t.Fatalf("data_handling_policy.%s missing", key)
		}
	}
	if got := requireString(t, dataHandling, "pack_content_truth_source"); got != "truzhen-packs" {
		t.Fatalf("pack_content_truth_source = %s, want truzhen-packs", got)
	}
	if got := requireString(t, dataHandling, "runtime_data_truth_source"); got != "truzhenos" {
		t.Fatalf("runtime_data_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, dataHandling, "cloud_order_data_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_order_data_truth_source = %s, want truzhen-cloud", got)
	}

	disclosure := requireObject(t, policy, "buyer_disclosure_policy")
	for _, key := range []string{"candidate_only_notice_surface", "no_professional_advice_notice", "owner_gate_notice", "asset_ref_notice", "evidence_required"} {
		if got := requireString(t, disclosure, key); got == "" {
			t.Fatalf("buyer_disclosure_policy.%s missing", key)
		}
	}

	rawGates, ok := policy["commercial_go_live_gates"].([]any)
	if !ok {
		t.Fatalf("commercial_go_live_gates missing")
	}
	requiredGates := map[string]bool{
		"terms_acceptance_before_purchase": false,
		"privacy_notice_visible":           false,
		"data_boundary_disclosed":          false,
		"candidate_disclaimer_visible":     false,
		"owner_gate_notice_visible":        false,
		"asset_ref_boundary_visible":       false,
	}
	for _, raw := range rawGates {
		gate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("terms gate = %T", raw)
		}
		gateID := requireString(t, gate, "gate_id")
		if _, ok := requiredGates[gateID]; ok {
			requiredGates[gateID] = true
		}
		requireStringIn(t, requireString(t, gate, "status"), "candidate_ready", "pending_owner_decision")
		if got := requireString(t, gate, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs", gateID)
		}
		if got := requireString(t, gate, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", gateID)
		}
	}
	for gateID, seen := range requiredGates {
		if !seen {
			t.Fatalf("missing terms go-live gate %s", gateID)
		}
	}

	negative := requireObject(t, policy, "negative_cases")
	for _, key := range []string{
		"purchase_without_terms_acceptance",
		"listing_without_privacy_notice",
		"pack_contains_customer_runtime_data",
		"role_output_marked_formal_advice",
		"cloud_order_personal_data_in_pack",
		"missing_data_deletion_export_route",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, policy["forbidden"]), "\n")
	for _, item := range []string{
		"store_customer_runtime_data_in_truzhen_packs",
		"store_cloud_order_personal_data_in_truzhen_packs",
		"present_role_output_as_formal_advice",
		"purchase_without_terms_acceptance",
		"hide_data_boundary_from_buyer",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	rawChecklist, ok := submission["review_checklist"].([]any)
	if !ok {
		t.Fatalf("review_checklist missing")
	}
	for _, raw := range rawChecklist {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("review checklist item = %T", raw)
		}
		if requireString(t, item, "check_id") == "terms_privacy_data_policy" {
			requireStringIn(t, requireString(t, item, "status"), "candidate_ready", "pending_owner_decision")
			if got := requireString(t, item, "required_evidence"); !strings.Contains(got, "commercial-terms-privacy-policy-candidate.json") {
				t.Fatalf("terms_privacy_data_policy required_evidence = %s", got)
			}
			goto manifestCheck
		}
	}
	t.Fatalf("marketplace review missing terms_privacy_data_policy checklist")

manifestCheck:
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/commercial-terms-privacy-policy-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/commercial-terms-privacy-policy-candidate.json")
}

func TestTeamOfficeCommercialGoLiveApprovalBlocksProductionUntilEvidenceAndOwnerDecision(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	approvalPath := "commerce/commercial-go-live-approval-candidate.json"
	approval := readJSON(t, filepath.Join(base, approvalPath))
	requireBool(t, approval, "candidate_only", true)
	requireBool(t, approval, "non_formal", true)
	if got := requireString(t, approval, "approval_status"); got != "not_approved_requires_cross_repo_evidence" {
		t.Fatalf("approval_status = %s, want not_approved_requires_cross_repo_evidence", got)
	}

	truthSources := requireObject(t, approval, "truth_sources")
	expectedTruthSources := map[string]string{
		"package_contents_truth": "truzhen-packs",
		"cloud_commerce_truth":   "truzhen-cloud",
		"runtime_install_truth":  "truzhenos",
		"gui_truth":              "truzhen-client-web-desktop",
	}
	for key, want := range expectedTruthSources {
		if got := requireString(t, truthSources, key); got != want {
			t.Fatalf("truth_sources.%s = %s, want %s", key, got, want)
		}
	}

	policies := requireObject(t, approval, "production_policies")
	expectedPolicies := map[string]string{
		"production_publish_policy":     "blocked_until_owner_authorizes_and_all_go_live_gates_pass",
		"real_payment_policy":           "blocked_until_owner_authorizes",
		"production_endpoint_policy":    "blocked_until_cross_repo_sandbox_evidence_passes",
		"signed_download_secret_policy": "not_stored_in_pack",
		"go_live_from_candidate_policy": "blocked_candidate_json_only_is_not_go_live",
	}
	for key, want := range expectedPolicies {
		if got := requireString(t, policies, key); got != want {
			t.Fatalf("production_policies.%s = %s, want %s", key, got, want)
		}
	}

	rawGates, ok := approval["required_go_live_gates"].([]any)
	if !ok {
		t.Fatalf("required_go_live_gates missing")
	}
	requiredGates := map[string]bool{
		"role_candidate_gui_verified":                    false,
		"candidate_bundle_exported":                      false,
		"artifact_manifest_hash_verified":                false,
		"terms_privacy_data_policy_accepted":             false,
		"publisher_identity_settlement_verified":         false,
		"asset_rights_verified":                          false,
		"sandbox_environment_ready":                      false,
		"commercial_api_contract_reviewed":               false,
		"commercial_observability_ready":                 false,
		"sandbox_purchase_entitlement_download_verified": false,
		"local_install_team_binding_verified":            false,
		"negative_cases_blocked_verified":                false,
		"independent_acceptance_signoff":                 false,
	}
	for _, raw := range rawGates {
		gate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("go-live gate = %T", raw)
		}
		gateID := requireString(t, gate, "gate_id")
		if _, ok := requiredGates[gateID]; ok {
			requiredGates[gateID] = true
		}
		requireStringIn(t, requireString(t, gate, "status"), "pending_cross_repo_execution", "pending_owner_decision")
		if got := requireString(t, gate, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs for go-live evidence", gateID)
		}
		if got := requireString(t, gate, "required_evidence"); got == "" {
			t.Fatalf("%s required_evidence missing", gateID)
		}
		requireBool(t, gate, "blocking_if_missing", true)
	}
	for gateID, seen := range requiredGates {
		if !seen {
			t.Fatalf("missing required go-live gate %s", gateID)
		}
	}

	decision := requireObject(t, approval, "owner_decision_record")
	for _, key := range []string{
		"owner_decision_ref",
		"scope",
		"approved_channels",
		"explicit_real_payment_authorization",
		"explicit_production_publish_authorization",
		"rollback_authorization",
		"evidence_pack_ref",
	} {
		if got := requireString(t, decision, key); got == "" {
			t.Fatalf("owner_decision_record.%s missing", key)
		}
	}

	rawRollback, ok := approval["rollback_plan"].([]any)
	if !ok {
		t.Fatalf("rollback_plan missing")
	}
	requiredRollbackStages := map[string]bool{
		"listing_unpublish":                       false,
		"entitlement_revoke_future_download":      false,
		"release_revoke":                          false,
		"installed_role_disable_or_block_upgrade": false,
		"buyer_notice":                            false,
		"receipt_preserve":                        false,
	}
	for _, raw := range rawRollback {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("rollback stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredRollbackStages[stageID]; ok {
			requiredRollbackStages[stageID] = true
		}
		requireStringIn(t, requireString(t, stage, "truth_source"), "truzhen-cloud", "truzhenos")
		if got := requireString(t, stage, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", stageID)
		}
	}
	for stageID, seen := range requiredRollbackStages {
		if !seen {
			t.Fatalf("missing rollback stage %s", stageID)
		}
	}

	negative := requireObject(t, approval, "negative_cases")
	for _, key := range []string{
		"approval_without_sandbox_evidence",
		"production_publish_without_owner_decision",
		"real_payment_without_owner_authorization",
		"go_live_with_missing_observability",
		"go_live_with_unblocked_negative_case",
		"package_hash_changed_after_approval",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, approval["forbidden"]), "\n")
	for _, item := range []string{
		"production_publish_without_owner_authorization",
		"real_payment_capture_without_owner_authorization",
		"go_live_from_candidate_json_only",
		"store_cloud_order_payment_or_entitlement_truth_in_truzhen_packs",
		"delete_receipts_during_rollback",
		"change_artifact_hash_after_approval",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, approvalPath) {
		t.Fatalf("candidate set artifact_files missing %s", approvalPath)
	}
	commercialization := requireObject(t, candidateSet, "commercialization")
	if got := requireString(t, commercialization, "commercial_go_live_approval"); got != approvalPath {
		t.Fatalf("commercial_go_live_approval = %s, want %s", got, approvalPath)
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	rawChecklist, ok := submission["review_checklist"].([]any)
	if !ok {
		t.Fatalf("review_checklist missing")
	}
	foundChecklist := false
	for _, raw := range rawChecklist {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("review checklist item = %T", raw)
		}
		if requireString(t, item, "check_id") != "commercial_go_live_approval" {
			continue
		}
		foundChecklist = true
		requireStringIn(t, requireString(t, item, "status"), "pending_owner_decision")
		if got := requireString(t, item, "required_evidence"); !strings.Contains(got, approvalPath) {
			t.Fatalf("commercial_go_live_approval required_evidence = %s", got)
		}
	}
	if !foundChecklist {
		t.Fatalf("marketplace review missing commercial_go_live_approval checklist")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == approvalPath {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", approvalPath)
}

func TestTeamOfficeCommercialGoLiveApprovalRequiresP11EvidenceAcceptanceChecklist(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	approval := readJSON(t, filepath.Join(base, "commerce", "commercial-go-live-approval-candidate.json"))

	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	contractRef := "tests/p11-sandbox-run-request-candidate.json#phase_dependency_contract"
	if got := requireString(t, approval, "source_p11_evidence_acceptance_checklist"); got != checklistPath {
		t.Fatalf("source_p11_evidence_acceptance_checklist = %s, want %s", got, checklistPath)
	}
	if got := requireString(t, approval, "source_p11_go_live_evidence_package"); got != packagePath {
		t.Fatalf("source_p11_go_live_evidence_package = %s, want %s", got, packagePath)
	}
	if got := requireString(t, approval, "source_execution_queue"); got != queuePath {
		t.Fatalf("source_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, approval, "source_p11_phase_dependency_contract"); got != contractRef {
		t.Fatalf("source_p11_phase_dependency_contract = %s, want %s", got, contractRef)
	}

	acceptance := requireObject(t, approval, "go_live_evidence_acceptance_gate")
	for key, want := range map[string]string{
		"gate_id":                         "p11_evidence_acceptance_checklist_verified",
		"checklist_ref":                   checklistPath,
		"go_live_evidence_package_ref":    packagePath,
		"phase_dependency_gate_ref":       checklistPath + "#phase_dependency_acceptance_gate",
		"required_phase_dependency_proof": "p11_execution_queue_phase_dependencies_verified",
		"expected_verified_status":        "verified_from_authoritative_evidence",
		"current_status":                  "pending_cross_repo_execution",
	} {
		if got := requireString(t, acceptance, key); got != want {
			t.Fatalf("go_live_evidence_acceptance_gate.%s = %s, want %s", key, got, want)
		}
	}
	requireBool(t, acceptance, "required_before_go_live_approval", true)
	requireBool(t, acceptance, "can_approve_go_live_without_this_gate", false)
	for _, evidence := range []string{
		"all_stage_acceptance_checks_verified",
		"p11_execution_queue_phase_dependencies_verified",
		"p11_phase_dependency_report",
		"owner_go_no_go_decision_recorded",
	} {
		requireStringSliceContains(t, asStringSlice(t, acceptance["required_evidence"]), evidence)
	}
	for _, blocker := range []string{
		"p11_evidence_acceptance_checklist_verified_missing",
		"p11_execution_queue_phase_dependencies_verified_missing",
		"owner_go_no_go_decision_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, acceptance["blocking_if_missing"]), blocker)
	}

	gate := findObjectByString(t, asObjectSlice(t, approval["required_go_live_gates"]), "gate_id", "p11_evidence_acceptance_checklist_verified")
	if got := requireString(t, gate, "truth_source"); got != "multi_repo" {
		t.Fatalf("p11_evidence_acceptance_checklist_verified truth_source = %s, want multi_repo", got)
	}
	requireBool(t, gate, "blocking_if_missing", true)
	if got := requireString(t, gate, "required_evidence"); !strings.Contains(got, checklistPath) || !strings.Contains(got, "p11_execution_queue_phase_dependencies_verified") {
		t.Fatalf("p11_evidence_acceptance_checklist_verified required_evidence = %s", got)
	}

	policies := requireObject(t, approval, "production_policies")
	if got := requireString(t, policies, "go_live_evidence_policy"); got != "blocked_until_p11_evidence_acceptance_checklist_and_package_verified" {
		t.Fatalf("production_policies.go_live_evidence_policy = %s", got)
	}
	forbidden := strings.Join(asStringSlice(t, approval["forbidden"]), "\n")
	if !strings.Contains(forbidden, "approve_go_live_without_p11_evidence_acceptance_checklist") {
		t.Fatalf("forbidden missing approve_go_live_without_p11_evidence_acceptance_checklist")
	}
}

func TestTeamOfficeCommercialGoLiveApprovalRequiresDownloadInstallAccessMatrix(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	approval := readJSON(t, filepath.Join(base, "commerce", "commercial-go-live-approval-candidate.json"))

	matrixPath := "commerce/download-install-access-matrix.json"
	if got := requireString(t, approval, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}

	policies := requireObject(t, approval, "production_policies")
	if got := requireString(t, policies, "download_install_access_policy"); got != "blocked_until_download_install_access_matrix_verified" {
		t.Fatalf("production_policies.download_install_access_policy = %s", got)
	}

	gate := findObjectByString(t, asObjectSlice(t, approval["required_go_live_gates"]), "gate_id", "download_install_access_matrix_verified")
	if got := requireString(t, gate, "truth_source"); got != "multi_repo" {
		t.Fatalf("download_install_access_matrix_verified truth_source = %s, want multi_repo", got)
	}
	if got := requireString(t, gate, "status"); got != "pending_cross_repo_execution" {
		t.Fatalf("download_install_access_matrix_verified status = %s", got)
	}
	requireBool(t, gate, "blocking_if_missing", true)
	evidence := requireString(t, gate, "required_evidence")
	for _, want := range []string{
		matrixPath,
		"unpaid_download",
		"refund_revoked_download",
		"version_unpublished_or_revoked_download_install",
		"artifact_hash_mismatch",
	} {
		if !strings.Contains(evidence, want) {
			t.Fatalf("download_install_access_matrix_verified required_evidence missing %s: %s", want, evidence)
		}
	}

	acceptance := requireObject(t, approval, "go_live_evidence_acceptance_gate")
	requireStringSliceContains(t, asStringSlice(t, acceptance["required_evidence"]), "download_install_access_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, acceptance["blocking_if_missing"]), "download_install_access_matrix_verified_missing")

	forbidden := strings.Join(asStringSlice(t, approval["forbidden"]), "\n")
	if !strings.Contains(forbidden, "approve_go_live_without_download_install_access_matrix") {
		t.Fatalf("forbidden missing approve_go_live_without_download_install_access_matrix")
	}
}

func TestTeamOfficeCommercialProductionPromotionGateBlocksProductionUntilP11EvidenceAndOwnerDecision(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gatePath := "commerce/commercial-production-promotion-gate-candidate.json"
	gate := readJSON(t, filepath.Join(base, gatePath))

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	requireBool(t, gate, "formal_write_allowed", false)
	requireBool(t, gate, "can_promote_to_production", false)
	requireBool(t, gate, "can_enable_real_payment", false)
	requireBool(t, gate, "can_publish_production_listing", false)
	requireBool(t, gate, "can_enable_production_download", false)
	if got := requireString(t, gate, "promotion_status"); got != "blocked_requires_p11_verified_evidence_and_owner_decision" {
		t.Fatalf("promotion_status = %s, want blocked_requires_p11_verified_evidence_and_owner_decision", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                        "role-pack-candidate-set://team-office-v0",
		"source_p11_sandbox_run_request":           "tests/p11-sandbox-run-request-candidate.json",
		"source_p11_go_live_evidence_package":      "tests/p11-commercial-go-live-evidence-package-template.json",
		"source_commercial_go_no_go_gate":          "tests/commercial-go-no-go-gate-candidate.json",
		"source_commercial_go_live_approval":       "commerce/commercial-go-live-approval-candidate.json",
		"source_commercial_readiness_verifier":     "tests/commercial-readiness-verifier-candidate.json",
		"source_commercial_chain_verifier":         "tests/commercial-chain-verifier-candidate.json",
		"source_product_readiness_evidence_matrix": "tests/product-readiness-evidence-matrix.json",
		"source_p0_p11_blocker_register":           "tests/p0-p11-commercialization-blocker-register-candidate.json",
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	controls := requireObject(t, gate, "production_controls")
	for _, key := range []string{
		"production_publish_blocked",
		"real_payment_capture_blocked",
		"production_download_distribution_blocked",
		"sandbox_to_production_auto_promotion_forbidden",
		"owner_explicit_go_decision_required",
		"independent_acceptance_required",
		"rollback_plan_required",
	} {
		requireBool(t, controls, key, true)
	}
	for key, want := range map[string]string{
		"cloud_listing_truth_source":          "truzhen-cloud",
		"payment_truth_source":                "truzhen-cloud",
		"license_entitlement_truth_source":    "truzhen-cloud",
		"production_download_truth_source":    "truzhen-cloud",
		"production_install_truth_source":     "truzhenos",
		"gui_evidence_truth_source":           "truzhen-client-web-desktop",
		"candidate_asset_truth_source":        "truzhen-packs",
		"production_promotion_decision_owner": "Owner + truzhen-cloud go-live gate",
	} {
		if got := requireString(t, controls, key); got != want {
			t.Fatalf("production_controls.%s = %s, want %s", key, got, want)
		}
	}

	rawChecks, ok := gate["promotion_gate_checks"].([]any)
	if !ok {
		t.Fatalf("promotion_gate_checks missing")
	}
	expectedChecks := map[string]string{
		"p11_sandbox_run_request_authorized":                  "owner_authorization",
		"role_candidate_bundle_hash_locked":                   "hash",
		"gui_user_view_flow_verified":                         "screenshot",
		"cloud_sandbox_upload_listing_receipts_verified":      "receipt",
		"sandbox_order_payment_entitlement_receipts_verified": "receipt",
		"signed_download_hash_continuity_verified":            "hash",
		"p11_go_live_evidence_package_verified":               "tests/p11-commercial-go-live-evidence-package-template.json",
		"download_install_access_matrix_verified":             "commerce/download-install-access-matrix.json",
		"truzhenos_install_team_binding_receipts_verified":    "receipt",
		"team_office_runtime_usage_smoke_verified":            "usage/team-office-runtime-usage-candidate.json",
		"negative_cases_blocked_verified":                     "blocked receipt",
		"independent_acceptance_signed":                       "signoff",
		"owner_go_no_go_decision_recorded":                    "owner_go_no_go",
	}
	if len(rawChecks) != len(expectedChecks) {
		t.Fatalf("promotion_gate_checks len = %d, want %d", len(rawChecks), len(expectedChecks))
	}
	seenChecks := map[string]bool{}
	for _, raw := range rawChecks {
		check, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("promotion gate check = %T", raw)
		}
		checkID := requireString(t, check, "check_id")
		wantEvidenceMarker, ok := expectedChecks[checkID]
		if !ok {
			t.Fatalf("unexpected promotion gate check %s", checkID)
		}
		seenChecks[checkID] = true
		if got := requireString(t, check, "current_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("%s current_status = %s, want missing_authoritative_evidence", checkID, got)
		}
		requireBool(t, check, "can_pass_gate", false)
		requireBool(t, check, "can_count_toward_production", false)
		if got := requireString(t, check, "truth_source"); got == "" {
			t.Fatalf("%s truth_source missing", checkID)
		}
		requiredEvidence := requireString(t, check, "required_evidence")
		if !strings.Contains(requiredEvidence, wantEvidenceMarker) {
			t.Fatalf("%s required_evidence = %s, want marker %s", checkID, requiredEvidence, wantEvidenceMarker)
		}
		if got := requireString(t, check, "required_before_production"); got == "" {
			t.Fatalf("%s required_before_production missing", checkID)
		}
	}
	for checkID := range expectedChecks {
		if !seenChecks[checkID] {
			t.Fatalf("missing promotion gate check %s", checkID)
		}
	}

	allowedAfterApproval := strings.Join(asStringSlice(t, gate["allowed_after_owner_cloud_approval"]), "\n")
	for _, item := range []string{
		"production_listing_publish_request",
		"real_payment_enable_request",
		"production_release_sign_request",
		"production_download_distribution_request",
		"production_install_observability_required",
	} {
		if !strings.Contains(allowedAfterApproval, item) {
			t.Fatalf("allowed_after_owner_cloud_approval missing %s", item)
		}
	}
	forbidden := strings.Join(asStringSlice(t, gate["forbidden_shortcuts"]), "\n")
	for _, item := range []string{
		"sandbox_success_auto_promotes_to_production",
		"real_payment_without_owner_decision",
		"production_listing_truth_in_pack_repo",
		"copy_sandbox_receipt_as_production_receipt",
		"store_payment_secret_in_pack",
		"store_cloud_token_in_pack",
		"skip_independent_acceptance",
		"skip_hash_continuity",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden_shortcuts missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, gatePath) {
		t.Fatalf("candidate set artifact_files missing %s", gatePath)
	}
	commercialization := requireObject(t, candidateSet, "commercialization")
	if got := requireString(t, commercialization, "commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("commercialization.commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}

	approval := readJSON(t, filepath.Join(base, "commerce", "commercial-go-live-approval-candidate.json"))
	if got := requireString(t, approval, "production_promotion_gate_ref"); got != gatePath {
		t.Fatalf("commercial-go-live approval production_promotion_gate_ref = %s, want %s", got, gatePath)
	}
	policies := requireObject(t, approval, "production_policies")
	if got := requireString(t, policies, "production_promotion_policy"); got != "blocked_until_p11_evidence_and_owner_go_decision" {
		t.Fatalf("production_policies.production_promotion_policy = %s", got)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("commercial-readiness source_commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}
	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_commercial_production_promotion_gate"); got != gatePath {
		t.Fatalf("commercial go/no-go source_commercial_production_promotion_gate = %s, want %s", got, gatePath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "commercial_production_promotion_gate_verified") {
		t.Fatalf("product readiness matrix missing commercial_production_promotion_gate_verified")
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	rawChecklist, ok := submission["review_checklist"].([]any)
	if !ok {
		t.Fatalf("review_checklist missing")
	}
	foundChecklist := false
	for _, raw := range rawChecklist {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("review checklist item = %T", raw)
		}
		if requireString(t, item, "check_id") != "commercial_production_promotion_gate" {
			continue
		}
		foundChecklist = true
		requireStringIn(t, requireString(t, item, "status"), "pending_owner_decision")
		if got := requireString(t, item, "required_evidence"); !strings.Contains(got, gatePath) {
			t.Fatalf("commercial_production_promotion_gate required_evidence = %s", got)
		}
	}
	if !foundChecklist {
		t.Fatalf("marketplace review missing commercial_production_promotion_gate checklist")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != gatePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, gatePath))
		if err != nil {
			t.Fatalf("read %s: %v", gatePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", gatePath, got, wantHash)
		}
		return
	}
	t.Fatalf("artifact manifest missing %s", gatePath)
}

func TestTeamOfficeCommercialProductionPromotionGateRequiresDownloadInstallAccessMatrix(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gate := readJSON(t, filepath.Join(base, "commerce", "commercial-production-promotion-gate-candidate.json"))

	matrixPath := "commerce/download-install-access-matrix.json"
	if got := requireString(t, gate, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}

	controls := requireObject(t, gate, "production_controls")
	requireBool(t, controls, "download_install_access_matrix_required", true)
	requireBool(t, controls, "can_enable_production_download_without_access_matrix", false)

	check := findObjectByString(t, asObjectSlice(t, gate["promotion_gate_checks"]), "check_id", "download_install_access_matrix_verified")
	if got := requireString(t, check, "truth_source"); got != "multi_repo" {
		t.Fatalf("download_install_access_matrix_verified truth_source = %s, want multi_repo", got)
	}
	requireBool(t, check, "can_pass_gate", false)
	requireBool(t, check, "can_count_toward_production", false)
	evidence := requireString(t, check, "required_evidence")
	for _, want := range []string{
		matrixPath,
		"unpaid_download",
		"refund_revoked_download",
		"version_unpublished_or_revoked_download_install",
		"artifact_hash_mismatch",
	} {
		if !strings.Contains(evidence, want) {
			t.Fatalf("download_install_access_matrix_verified required_evidence missing %s: %s", want, evidence)
		}
	}

	requiredBeforeProduction := requireString(t, check, "required_before_production")
	if !strings.Contains(requiredBeforeProduction, "production signed download") || !strings.Contains(requiredBeforeProduction, "matrix") {
		t.Fatalf("download_install_access_matrix_verified required_before_production = %s", requiredBeforeProduction)
	}

	forbidden := strings.Join(asStringSlice(t, gate["forbidden_shortcuts"]), "\n")
	if !strings.Contains(forbidden, "enable_production_download_without_download_install_access_matrix") {
		t.Fatalf("forbidden_shortcuts missing enable_production_download_without_download_install_access_matrix")
	}
}

func TestTeamOfficeCommercialProductionPromotionGateRequiresP11PackageAndRuntimeSmoke(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gate := readJSON(t, filepath.Join(base, "commerce", "commercial-production-promotion-gate-candidate.json"))

	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	runtimePath := "usage/team-office-runtime-usage-candidate.json"
	if got := requireString(t, gate, "source_p11_go_live_evidence_package"); got != packagePath {
		t.Fatalf("source_p11_go_live_evidence_package = %s, want %s", got, packagePath)
	}
	if got := requireString(t, gate, "source_runtime_usage_candidate"); got != runtimePath {
		t.Fatalf("source_runtime_usage_candidate = %s, want %s", got, runtimePath)
	}

	controls := requireObject(t, gate, "production_controls")
	requireBool(t, controls, "p11_go_live_evidence_package_required", true)
	requireBool(t, controls, "runtime_usage_smoke_required", true)
	requireBool(t, controls, "can_enable_real_payment_without_p11_package", false)
	requireBool(t, controls, "can_enable_production_download_without_runtime_usage_smoke", false)

	checks := asObjectSlice(t, gate["promotion_gate_checks"])
	p11Check := findObjectByString(t, checks, "check_id", "p11_go_live_evidence_package_verified")
	if got := requireString(t, p11Check, "truth_source"); got != "multi_repo" {
		t.Fatalf("p11_go_live_evidence_package_verified truth_source = %s, want multi_repo", got)
	}
	requireBool(t, p11Check, "can_pass_gate", false)
	requireBool(t, p11Check, "can_count_toward_production", false)
	p11Evidence := requireString(t, p11Check, "required_evidence")
	for _, want := range []string{packagePath, "final_go_live_decision", "production_promotion_receipts", "owner_go_no_go_decision"} {
		if !strings.Contains(p11Evidence, want) {
			t.Fatalf("p11_go_live_evidence_package_verified required_evidence missing %s: %s", want, p11Evidence)
		}
	}

	runtimeCheck := findObjectByString(t, checks, "check_id", "team_office_runtime_usage_smoke_verified")
	if got := requireString(t, runtimeCheck, "truth_source"); got != "multi_repo" {
		t.Fatalf("team_office_runtime_usage_smoke_verified truth_source = %s, want multi_repo", got)
	}
	requireBool(t, runtimeCheck, "can_pass_gate", false)
	requireBool(t, runtimeCheck, "can_count_toward_production", false)
	runtimeEvidence := requireString(t, runtimeCheck, "required_evidence")
	for _, want := range []string{runtimePath, "secretary and five advisor call trace", "capability invocation candidate", "blocked receipts"} {
		if !strings.Contains(runtimeEvidence, want) {
			t.Fatalf("team_office_runtime_usage_smoke_verified required_evidence missing %s: %s", want, runtimeEvidence)
		}
	}

	forbidden := strings.Join(asStringSlice(t, gate["forbidden_shortcuts"]), "\n")
	for _, item := range []string{
		"enable_real_payment_without_p11_go_live_evidence_package",
		"enable_production_download_without_runtime_usage_smoke",
		"publish_production_listing_without_runtime_usage_smoke",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden_shortcuts missing %s", item)
		}
	}
}

func TestTeamOfficeSupportRefundRevocationPolicyDefinesCommercialAfterSalesBoundary(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "support-refund-revocation-policy-candidate.json"))
	requireBool(t, policy, "candidate_only", true)
	requireBool(t, policy, "non_formal", true)
	if got := requireString(t, policy, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "support_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("support_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "real_payment_policy"); got != "blocked_until_owner_authorizes" {
		t.Fatalf("real_payment_policy = %s, want blocked_until_owner_authorizes", got)
	}

	supportPolicy := requireObject(t, policy, "support_policy")
	requireBool(t, supportPolicy, "support_channel_required", true)
	for _, key := range []string{"support_channel_surface", "response_sla_policy", "evidence_required"} {
		if got := requireString(t, supportPolicy, key); got == "" {
			t.Fatalf("support_policy.%s missing", key)
		}
	}

	refundPolicy := requireObject(t, policy, "refund_policy")
	requireBool(t, refundPolicy, "sandbox_refund_supported", true)
	requireBool(t, refundPolicy, "entitlement_revocation_required", true)
	if got := requireString(t, refundPolicy, "historical_receipts_policy"); got != "do_not_delete_historical_receipts" {
		t.Fatalf("historical_receipts_policy = %s, want do_not_delete_historical_receipts", got)
	}
	for _, key := range []string{"refund_window_policy", "refund_decision_truth_source", "evidence_required"} {
		if got := requireString(t, refundPolicy, key); got == "" {
			t.Fatalf("refund_policy.%s missing", key)
		}
	}

	revocationPolicy := requireObject(t, policy, "revocation_policy")
	for _, key := range []string{"block_future_download", "block_new_install", "preserve_installed_history", "notice_required"} {
		requireBool(t, revocationPolicy, key, true)
	}
	for _, key := range []string{"revocation_truth_source", "revocation_reason_policy", "evidence_required"} {
		if got := requireString(t, revocationPolicy, key); got == "" {
			t.Fatalf("revocation_policy.%s missing", key)
		}
	}

	buyerNotice := requireObject(t, policy, "buyer_notice_policy")
	requireBool(t, buyerNotice, "requires_owner_visible_notice", true)
	for _, key := range []string{"support_contact_surface", "refund_terms_surface", "revocation_reason_surface"} {
		if got := requireString(t, buyerNotice, key); got == "" {
			t.Fatalf("buyer_notice_policy.%s missing", key)
		}
	}

	rawStates, ok := policy["states"].([]any)
	if !ok {
		t.Fatalf("states missing")
	}
	requiredStates := map[string]bool{
		"support_request_opened":         false,
		"refund_requested":               false,
		"refund_approved":                false,
		"refund_rejected":                false,
		"entitlement_revocation_pending": false,
		"entitlement_revoked":            false,
		"release_revoked_notice_sent":    false,
	}
	for _, raw := range rawStates {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("state item = %T", raw)
		}
		stateID := requireString(t, item, "state_id")
		if _, ok := requiredStates[stateID]; ok {
			requiredStates[stateID] = true
		}
		if got := requireString(t, item, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source must not be truzhen-packs", stateID)
		}
		if got := requireString(t, item, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", stateID)
		}
	}
	for stateID, seen := range requiredStates {
		if !seen {
			t.Fatalf("missing after-sales state %s", stateID)
		}
	}

	negative := requireObject(t, policy, "negative_cases")
	for _, key := range []string{
		"missing_support_channel",
		"refund_without_entitlement_revocation",
		"delete_historical_receipts_after_refund",
		"silent_release_revocation",
		"store_support_ticket_truth_in_packs",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, policy["forbidden"]), "\n")
	for _, item := range []string{
		"store_support_ticket_truth_in_truzhen_packs",
		"store_refund_truth_in_truzhen_packs",
		"delete_historical_receipts_after_refund",
		"silent_release_revocation",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/support-refund-revocation-policy-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/support-refund-revocation-policy-candidate.json")
}

func TestTeamOfficePurchaseRequiresSupportRefundRevocationDisclosures(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	policy := readJSON(t, filepath.Join(base, "commerce", "support-refund-revocation-policy-candidate.json"))

	disclosure := requireObject(t, policy, "pre_purchase_disclosure_policy")
	requireBool(t, disclosure, "must_show_before_purchase", true)
	requireBool(t, disclosure, "buyer_acknowledgement_required", true)
	if got := requireString(t, disclosure, "blocking_status_if_missing"); got != "blocked_commercial_purchase_disclosures_missing" {
		t.Fatalf("blocking_status_if_missing = %s", got)
	}
	if got := requireString(t, disclosure, "expected_evidence"); got != "pre_purchase_disclosure_ack_receipt_ref" {
		t.Fatalf("expected_evidence = %s", got)
	}
	for _, surface := range []string{
		"cloud_listing_support_contact_surface",
		"cloud_listing_refund_terms_surface",
		"cloud_listing_revocation_notice_surface",
		"cloud_listing_candidate_role_disclaimer_surface",
	} {
		requireStringSliceContains(t, asStringSlice(t, disclosure["required_surfaces"]), surface)
	}

	negative := requireObject(t, policy, "negative_cases")
	missing := requireObject(t, negative, "purchase_without_support_refund_revocation_disclosure")
	if got := requireString(t, missing, "expected_status"); got != "blocked_commercial_purchase_disclosures_missing" {
		t.Fatalf("purchase_without_support_refund_revocation_disclosure expected_status = %s", got)
	}

	listing := readJSON(t, filepath.Join(base, "commerce", "cloud-listing-candidate.json"))
	for _, surface := range []string{
		"support_contact_surface",
		"refund_terms_surface",
		"revocation_notice_surface",
		"candidate_role_disclaimer_surface",
	} {
		if got := requireString(t, listing, surface); got == "" {
			t.Fatalf("cloud listing missing %s", surface)
		}
	}

	submission := readJSON(t, filepath.Join(base, "commerce", "marketplace-review-submission-candidate.json"))
	foundChecklist := false
	for _, item := range asObjectSlice(t, submission["review_checklist"]) {
		if requireString(t, item, "check_id") != "pre_purchase_support_refund_revocation_disclosures" {
			continue
		}
		foundChecklist = true
		if got := requireString(t, item, "status"); got != "candidate_ready" {
			t.Fatalf("pre_purchase_support_refund_revocation_disclosures status = %s", got)
		}
		if got := requireString(t, item, "required_evidence"); !strings.Contains(got, "pre_purchase_disclosure_policy") {
			t.Fatalf("pre_purchase_support_refund_revocation_disclosures required_evidence = %s", got)
		}
	}
	if !foundChecklist {
		t.Fatalf("review_checklist missing pre_purchase_support_refund_revocation_disclosures")
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policyBlock := requireObject(t, matrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, policyBlock["required_before_completion_claim"]), "pre_purchase_support_refund_revocation_disclosures_verified")
}

func TestTeamOfficeInstallPreflightRequestDefinesInstallerEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	preflight := readJSON(t, filepath.Join(base, "install", "install-preflight-request-candidate.json"))
	requireBool(t, preflight, "candidate_only", true)
	requireBool(t, preflight, "non_formal", true)
	if got := requireString(t, preflight, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, preflight, "commerce_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("commerce_truth_source = %s, want truzhen-cloud", got)
	}

	entrySurfaces := strings.Join(asStringSlice(t, preflight["entry_surfaces"]), "\n")
	for _, surface := range []string{"purchased_product_download", "local_pack_manager", "team_settings_role_tab"} {
		if !strings.Contains(entrySurfaces, surface) {
			t.Fatalf("entry_surfaces missing %s", surface)
		}
	}

	requiredInputs := strings.Join(asStringSlice(t, preflight["required_inputs"]), "\n")
	for _, input := range []string{"downloaded_artifact_ref", "artifact_manifest", "entitlement_ref", "download_receipt_ref", "target_team_ref"} {
		if !strings.Contains(requiredInputs, input) {
			t.Fatalf("required_inputs missing %s", input)
		}
	}

	rawChecks, ok := preflight["preflight_checks"].([]any)
	if !ok {
		t.Fatalf("preflight_checks missing")
	}
	requiredChecks := map[string]bool{
		"entitlement_validation":              false,
		"artifact_manifest_hash_verification": false,
		"artifact_signature_verification":     false,
		"role_pack_schema_validation":         false,
		"forbidden_artifact_scan":             false,
		"team_slot_compatibility":             false,
		"appearance_asset_ref_validation":     false,
		"owner_gate_before_binding":           false,
	}
	for _, raw := range rawChecks {
		check, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("preflight check = %T", raw)
		}
		checkID := requireString(t, check, "check_id")
		if _, ok := requiredChecks[checkID]; ok {
			requiredChecks[checkID] = true
		}
		if got := requireString(t, check, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs", checkID)
		}
		if got := requireString(t, check, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", checkID)
		}
	}
	for checkID, seen := range requiredChecks {
		if !seen {
			t.Fatalf("missing preflight check %s", checkID)
		}
	}

	outputs := requireObject(t, preflight, "output_records")
	for _, key := range []string{"preflight_result_ref", "install_plan_candidate_ref", "install_receipt_ref", "blocked_reason", "evidence_bundle_ref"} {
		if got := requireString(t, outputs, key); got == "" {
			t.Fatalf("output_records.%s missing", key)
		}
	}

	negative := requireObject(t, preflight, "negative_cases")
	for _, key := range []string{
		"missing_entitlement",
		"artifact_hash_mismatch",
		"unsigned_artifact",
		"forbidden_artifact",
		"incompatible_team_slot",
		"raw_asset_reference",
		"owner_gate_missing",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "install/install-preflight-request-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing install/install-preflight-request-candidate.json")
}

func TestTeamOfficeCommercialReceiptChainCorrelatesProductLifecycle(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	chain := readJSON(t, filepath.Join(base, "commerce", "commercial-receipt-chain-candidate.json"))
	requireBool(t, chain, "candidate_only", true)
	requireBool(t, chain, "non_formal", true)
	if got := requireString(t, chain, "completion_status"); got != "not_run_requires_cross_repo_receipts" {
		t.Fatalf("completion_status = %s, want not_run_requires_cross_repo_receipts", got)
	}

	truthSources := requireObject(t, chain, "truth_sources")
	for _, key := range []string{"bundle_export", "cloud_commerce", "local_install", "team_binding", "runtime_usage"} {
		if got := requireString(t, truthSources, key); got == "truzhen-packs" {
			t.Fatalf("truth_sources.%s cannot be truzhen-packs", key)
		}
	}

	correlationKeys := strings.Join(asStringSlice(t, chain["correlation_keys"]), "\n")
	for _, key := range []string{"artifact_ref", "artifact_sha256", "listing_draft_ref", "sandbox_order_ref", "entitlement_ref", "download_receipt_ref", "install_receipt_ref", "team_binding_receipt_ref"} {
		if !strings.Contains(correlationKeys, key) {
			t.Fatalf("correlation_keys missing %s", key)
		}
	}

	rawStages, ok := chain["receipt_sequence"].([]any)
	if !ok {
		t.Fatalf("receipt_sequence missing")
	}
	requiredStages := map[string]bool{
		"role_candidate_bundle_export": false,
		"cloud_upload_draft":           false,
		"marketplace_review_candidate": false,
		"sandbox_order_payment":        false,
		"entitlement_issue":            false,
		"signed_download":              false,
		"local_install_preflight":      false,
		"local_install_enable":         false,
		"team_role_binding":            false,
		"team_office_runtime_use":      false,
	}
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("receipt stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		if got := requireString(t, stage, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs", stageID)
		}
		requiredEvidence := strings.Join(asStringSlice(t, stage["required_evidence"]), "\n")
		if !strings.Contains(requiredEvidence, "receipt") && !strings.Contains(requiredEvidence, "candidate_ref") {
			t.Fatalf("%s required_evidence must include receipt or candidate_ref", stageID)
		}
		stageKeys := strings.Join(asStringSlice(t, stage["required_correlation_keys"]), "\n")
		if !strings.Contains(stageKeys, "artifact_ref") {
			t.Fatalf("%s required_correlation_keys missing artifact_ref", stageID)
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing receipt stage %s", stageID)
		}
	}

	completionGate := requireObject(t, chain, "completion_gate")
	for _, key := range []string{"all_receipts_present", "artifact_hash_consistent", "entitlement_valid_for_team", "owner_gate_recorded_for_binding", "runtime_outputs_candidate_only"} {
		if got := requireString(t, completionGate, key); got != "required" {
			t.Fatalf("completion_gate.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, chain, "negative_cases")
	for _, key := range []string{
		"missing_cloud_upload_receipt",
		"download_hash_differs_from_upload",
		"install_without_entitlement_receipt",
		"binding_without_owner_gate_receipt",
		"runtime_output_without_candidate_flag",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/commercial-receipt-chain-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/commercial-receipt-chain-candidate.json")
}

func TestTeamOfficeCommercialReceiptChainCarriesRoleStudioLineageIntoRuntimeUse(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	chain := readJSON(t, filepath.Join(base, "commerce", "commercial-receipt-chain-candidate.json"))

	correlationKeys := asStringSlice(t, chain["correlation_keys"])
	for _, key := range []string{
		"candidate_set_ref",
		"bundle_tree_sha256",
		"six_role_pack_refs",
		"role_pack_ref_set_hash",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		requireStringSliceContains(t, correlationKeys, key)
	}

	lineage := requireObject(t, chain, "role_studio_lineage_policy")
	requireBool(t, lineage, "same_six_role_pack_refs_required_from_export_to_runtime", true)
	requireBool(t, lineage, "capability_reference_must_use_installed_enabled_role_versions", true)
	requireBool(t, lineage, "team_settings_binding_must_reference_installed_role_pack_versions", true)
	requireBool(t, lineage, "runtime_usage_must_reference_team_binding_receipt", true)
	for key, want := range map[string]string{
		"candidate_set_source":        "candidate-set.json",
		"role_pack_refs_source":       "candidate-set.json#role_pack_refs",
		"enabled_versions_source":     "install/install-runtime-activation-map-candidate.json#required_role_activations",
		"team_binding_source":         "bindings/team-office-role-binding-candidate.json",
		"runtime_usage_source":        "usage/team-office-runtime-usage-candidate.json",
		"capability_reference_source": "capabilities/team-office-capability-requirements-candidate.json",
	} {
		if got := requireString(t, lineage, key); got != want {
			t.Fatalf("role_studio_lineage_policy.%s = %s, want %s", key, got, want)
		}
	}

	stages := asObjectSlice(t, chain["receipt_sequence"])
	stageKeyRequirements := map[string][]string{
		"role_candidate_bundle_export": {
			"candidate_set_ref",
			"bundle_tree_sha256",
			"six_role_pack_refs",
			"role_pack_ref_set_hash",
		},
		"local_install_enable": {
			"candidate_set_ref",
			"bundle_tree_sha256",
			"six_role_pack_refs",
			"role_pack_ref_set_hash",
			"enabled_role_pack_version_refs",
		},
		"team_role_binding": {
			"candidate_set_ref",
			"enabled_role_pack_version_refs",
			"slot_mapping_refs",
			"role_pack_ref_set_hash",
			"team_binding_receipt_ref",
		},
		"team_office_runtime_use": {
			"candidate_set_ref",
			"enabled_role_pack_version_refs",
			"slot_mapping_refs",
			"role_pack_ref_set_hash",
			"team_binding_receipt_ref",
			"runtime_receipt_candidate_ref",
		},
	}
	for stageID, expectedKeys := range stageKeyRequirements {
		stage := findObjectByString(t, stages, "stage_id", stageID)
		stageKeys := asStringSlice(t, stage["required_correlation_keys"])
		for _, key := range expectedKeys {
			requireStringSliceContains(t, stageKeys, key)
		}
	}

	completionGate := requireObject(t, chain, "completion_gate")
	for _, key := range []string{"role_studio_lineage_verified", "same_six_role_pack_refs_from_export_to_runtime"} {
		if got := requireString(t, completionGate, key); got != "required" {
			t.Fatalf("completion_gate.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, chain, "negative_cases")
	mismatch := requireObject(t, negative, "runtime_role_set_differs_from_installed_lineage")
	if got := requireString(t, mismatch, "expected_status"); got != "blocked_role_lineage_mismatch" {
		t.Fatalf("runtime_role_set_differs_from_installed_lineage.expected_status = %s, want blocked_role_lineage_mismatch", got)
	}

	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	verifierKeys := asStringSlice(t, requireObject(t, verifier, "correlation_requirements")["required_correlation_keys"])
	for _, key := range []string{"six_role_pack_refs", "role_pack_ref_set_hash"} {
		requireStringSliceContains(t, verifierKeys, key)
	}
	verifierPass := requireObject(t, verifier, "normal_commercialization_pass_conditions")
	if got := requireString(t, verifierPass, "role_studio_lineage_verified"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.role_studio_lineage_verified = %s, want required", got)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "role_studio_lineage_verified")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "role_studio_lineage_verified")
}

func TestTeamOfficeArtifactBundleLayoutDefinesDeterministicUploadPackage(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	layout := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-layout-candidate.json"))
	requireBool(t, layout, "candidate_only", true)
	requireBool(t, layout, "non_formal", true)
	if got := requireString(t, layout, "build_truth_source"); got != "truzhenos" {
		t.Fatalf("build_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, layout, "distribution_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("distribution_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, layout, "package_contents_truth_source"); got != "truzhen-packs" {
		t.Fatalf("package_contents_truth_source = %s, want truzhen-packs", got)
	}

	if got := requireString(t, layout, "package_format"); got != "role_pack_candidate_bundle_v0" {
		t.Fatalf("package_format = %s, want role_pack_candidate_bundle_v0", got)
	}
	if got := requireString(t, layout, "archive_format"); got == "" {
		t.Fatalf("archive_format missing")
	}

	inputs := strings.Join(asStringSlice(t, layout["required_inputs"]), "\n")
	for _, input := range []string{"candidate-set.json", "commerce/artifact-manifest.json", "commerce/release-candidate-package.json", "commerce/marketplace-review-submission-candidate.json"} {
		if !strings.Contains(inputs, input) {
			t.Fatalf("required_inputs missing %s", input)
		}
	}

	policy := requireObject(t, layout, "deterministic_build_policy")
	for _, key := range []string{"path_normalization", "file_order", "timestamp_policy", "compression_policy", "manifest_hash_policy"} {
		if got := requireString(t, policy, key); got == "" {
			t.Fatalf("deterministic_build_policy.%s missing", key)
		}
	}

	outputs := requireObject(t, layout, "package_outputs")
	for _, key := range []string{"candidate_bundle_ref", "bundle_sha256", "upload_request_ref", "signing_request_ref", "artifact_manifest_ref"} {
		if got := requireString(t, outputs, key); got == "" {
			t.Fatalf("package_outputs.%s missing", key)
		}
	}

	consumers := strings.Join(asStringSlice(t, layout["install_and_distribution_consumers"]), "\n")
	for _, consumer := range []string{"truzhen-cloud upload draft", "truzhen-cloud signed download", "truzhenos install preflight", "team settings role tab"} {
		if !strings.Contains(consumers, consumer) {
			t.Fatalf("install_and_distribution_consumers missing %s", consumer)
		}
	}

	negative := requireObject(t, layout, "negative_cases")
	for _, key := range []string{
		"missing_artifact_manifest",
		"file_not_declared_in_manifest",
		"nondeterministic_timestamp",
		"private_key_in_bundle",
		"bundle_hash_mismatch_before_upload",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/artifact-bundle-layout-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/artifact-bundle-layout-candidate.json")
}

func TestTeamOfficeInstalledRoleCatalogFeedsTeamSettingsReplacement(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	catalog := readJSON(t, filepath.Join(base, "bindings", "team-settings-installed-role-catalog-candidate.json"))
	requireBool(t, catalog, "candidate_only", true)
	requireBool(t, catalog, "non_formal", true)
	if got := requireString(t, catalog, "catalog_truth_source"); got != "truzhenos" {
		t.Fatalf("catalog_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, catalog, "entry_surface"); got != "Team Settings Role Tab" {
		t.Fatalf("entry_surface = %s, want Team Settings Role Tab", got)
	}
	if got := requireString(t, catalog, "install_receipt_required"); got != "true" {
		t.Fatalf("install_receipt_required = %s, want true", got)
	}

	requiredInputs := strings.Join(asStringSlice(t, catalog["required_inputs"]), "\n")
	for _, input := range []string{"install_receipt_ref", "enabled_role_pack_version", "team_ref", "role_slots_ref", "artifact_ref"} {
		if !strings.Contains(requiredInputs, input) {
			t.Fatalf("required_inputs missing %s", input)
		}
	}

	rawItems, ok := catalog["replaceable_role_items"].([]any)
	if !ok {
		t.Fatalf("replaceable_role_items missing")
	}
	if len(rawItems) != 6 {
		t.Fatalf("replaceable_role_items len = %d, want 6", len(rawItems))
	}
	requiredSlots := map[string]bool{
		"team_office.secretary_general":  false,
		"team_office.advisor.strategy":   false,
		"team_office.advisor.product":    false,
		"team_office.advisor.operations": false,
		"team_office.advisor.finance":    false,
		"team_office.advisor.legal_risk": false,
	}
	for _, raw := range rawItems {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("catalog item = %T", raw)
		}
		slotRef := requireString(t, item, "slot_ref")
		if _, ok := requiredSlots[slotRef]; ok {
			requiredSlots[slotRef] = true
		}
		if got := requireString(t, item, "replaceable_in_team_settings"); got != "true" {
			t.Fatalf("%s replaceable_in_team_settings = %s, want true", slotRef, got)
		}
		if got := requireString(t, item, "owner_gate_required_for_binding"); got != "true" {
			t.Fatalf("%s owner_gate_required_for_binding = %s, want true", slotRef, got)
		}
		if got := requireString(t, item, "source_install_receipt_ref"); got == "" {
			t.Fatalf("%s source_install_receipt_ref missing", slotRef)
		}
		if got := requireString(t, item, "role_pack_ref"); !strings.HasPrefix(got, "role_pack://team-office-") {
			t.Fatalf("%s role_pack_ref = %s", slotRef, got)
		}
	}
	for slotRef, seen := range requiredSlots {
		if !seen {
			t.Fatalf("missing replaceable role slot %s", slotRef)
		}
	}

	gui := requireObject(t, catalog, "gui_projection")
	for _, key := range []string{"tab_id", "page_state", "required_badges", "required_actions"} {
		if got := requireString(t, gui, key); got == "" {
			t.Fatalf("gui_projection.%s missing", key)
		}
	}

	negative := requireObject(t, catalog, "negative_cases")
	for _, key := range []string{
		"missing_install_receipt",
		"entitlement_revoked",
		"incompatible_slot",
		"bind_without_owner_gate",
		"catalog_from_uninstalled_bundle",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "bindings/team-settings-installed-role-catalog-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing bindings/team-settings-installed-role-catalog-candidate.json")
}

func TestTeamOfficeBuyerLibraryInstallStateConnectsPurchaseDownloadInstallAndTeamSettings(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	library := readJSON(t, filepath.Join(base, "commerce", "buyer-library-install-state-candidate.json"))
	requireBool(t, library, "candidate_only", true)
	requireBool(t, library, "non_formal", true)
	if got := requireString(t, library, "cloud_library_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_library_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, library, "local_install_truth_source"); got != "truzhenos" {
		t.Fatalf("local_install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, library, "gui_truth_source"); got != "truzhen-client-web-desktop" {
		t.Fatalf("gui_truth_source = %s, want truzhen-client-web-desktop", got)
	}

	surfaces := strings.Join(asStringSlice(t, library["buyer_library_surfaces"]), "\n")
	for _, surface := range []string{"cloud_purchased_library", "buyer_order_detail", "local_pack_manager", "team_settings_role_tab"} {
		if !strings.Contains(surfaces, surface) {
			t.Fatalf("buyer_library_surfaces missing %s", surface)
		}
	}

	inputs := strings.Join(asStringSlice(t, library["required_inputs"]), "\n")
	for _, input := range []string{"sandbox_order_ref", "sandbox_payment_receipt_ref", "entitlement_ref", "download_receipt_ref", "artifact_sha256", "install_receipt_ref", "enabled_role_pack_version_ref", "team_ref"} {
		if !strings.Contains(inputs, input) {
			t.Fatalf("required_inputs missing %s", input)
		}
	}

	rawStates, ok := library["library_states"].([]any)
	if !ok {
		t.Fatalf("library_states missing")
	}
	requiredStates := map[string]bool{
		"purchased_available":         false,
		"download_ready":              false,
		"downloaded_verified":         false,
		"install_available":           false,
		"installed_enabled":           false,
		"team_settings_replaceable":   false,
		"reinstall_available":         false,
		"entitlement_revoked_blocked": false,
	}
	for _, raw := range rawStates {
		state, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("library state = %T", raw)
		}
		stateID := requireString(t, state, "state_id")
		if _, ok := requiredStates[stateID]; ok {
			requiredStates[stateID] = true
		}
		if got := requireString(t, state, "truth_source"); got == "truzhen-packs" {
			t.Fatalf("%s truth_source cannot be truzhen-packs", stateID)
		}
		if got := requireString(t, state, "required_evidence"); got == "" {
			t.Fatalf("%s required_evidence missing", stateID)
		}
		if got := requireString(t, state, "user_visible_status"); got == "" {
			t.Fatalf("%s user_visible_status missing", stateID)
		}
	}
	for stateID, seen := range requiredStates {
		if !seen {
			t.Fatalf("missing buyer library state %s", stateID)
		}
	}

	rawActions, ok := library["gui_action_matrix"].([]any)
	if !ok {
		t.Fatalf("gui_action_matrix missing")
	}
	requiredActions := map[string]bool{
		"download_from_library":        false,
		"install_to_local_team":        false,
		"open_team_settings_replace":   false,
		"reinstall_same_version":       false,
		"rollback_binding_to_previous": false,
		"view_receipts":                false,
	}
	for _, raw := range rawActions {
		action, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("gui action = %T", raw)
		}
		actionID := requireString(t, action, "action_id")
		if _, ok := requiredActions[actionID]; ok {
			requiredActions[actionID] = true
		}
		for _, key := range []string{"entry_surface", "backend_evidence_required", "blocked_reason_if_unavailable"} {
			if got := requireString(t, action, key); got == "" {
				t.Fatalf("%s %s missing", actionID, key)
			}
		}
	}
	for actionID, seen := range requiredActions {
		if !seen {
			t.Fatalf("missing buyer library action %s", actionID)
		}
	}

	completionGate := requireObject(t, library, "completion_gate")
	for _, key := range []string{"receipt_chain_complete", "entitlement_valid", "artifact_hash_verified", "install_receipt_present", "team_settings_catalog_refreshed"} {
		if got := requireString(t, completionGate, key); got != "required" {
			t.Fatalf("completion_gate.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, library, "negative_cases")
	for _, key := range []string{
		"library_visible_without_purchase",
		"download_without_entitlement",
		"install_from_unverified_hash",
		"show_replaceable_without_install_receipt",
		"reinstall_after_entitlement_revoked",
		"hide_refund_revocation_reason",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	forbidden := strings.Join(asStringSlice(t, library["forbidden"]), "\n")
	for _, item := range []string{
		"store_order_truth_in_truzhen_packs",
		"show_unpurchased_as_owned",
		"install_without_entitlement",
		"team_settings_replace_without_owner_gate",
		"claim_gui_download_equals_install",
		"delete_existing_receipts_on_reinstall",
	} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == "commerce/buyer-library-install-state-candidate.json" {
			requireStringIn(t, requireString(t, item, "required_for"), "download", "install", "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing commerce/buyer-library-install-state-candidate.json")
}

func TestTeamOfficeInstallRuntimeActivationMapConnectsInstallTeamSettingsAndUsage(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	mapPath := requireString(t, candidateSet, "install_runtime_activation_map")
	requireStringSliceContains(t, asStringSlice(t, candidateSet["artifact_files"]), mapPath)

	activation := readJSON(t, filepath.Join(base, mapPath))
	requireBool(t, activation, "candidate_only", true)
	requireBool(t, activation, "non_formal", true)
	if got := requireString(t, activation, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, activation, "cloud_entitlement_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_entitlement_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, activation, "gui_truth_source"); got != "truzhen-client-web-desktop" {
		t.Fatalf("gui_truth_source = %s, want truzhen-client-web-desktop", got)
	}
	for key, want := range map[string]string{
		"artifact_manifest":              "commerce/artifact-manifest.json",
		"install_preflight_request":      "install/install-preflight-request-candidate.json",
		"installed_role_catalog":         "bindings/team-settings-installed-role-catalog-candidate.json",
		"buyer_library_install_state":    "commerce/buyer-library-install-state-candidate.json",
		"runtime_usage_candidate":        "usage/team-office-runtime-usage-candidate.json",
		"team_role_binding_candidate":    "bindings/team-office-role-binding-candidate.json",
		"commercial_receipt_chain":       "commerce/commercial-receipt-chain-candidate.json",
		"artifact_manifest_closure_gate": "tests/artifact-manifest-closure-gate-candidate.json",
	} {
		if got := requireString(t, activation, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	steps := asObjectSlice(t, activation["activation_steps"])
	wantSteps := []string{
		"download_receipt_verified",
		"install_preflight_passed",
		"role_pack_versions_enabled",
		"capability_role_reference_ready",
		"team_settings_catalog_refreshed",
		"team_binding_owner_gate_ready",
		"runtime_usage_smoke_ready",
	}
	if len(steps) != len(wantSteps) {
		t.Fatalf("activation_steps len = %d, want %d", len(steps), len(wantSteps))
	}
	for i, step := range steps {
		if got := requireString(t, step, "step_id"); got != wantSteps[i] {
			t.Fatalf("activation_steps[%d] = %s, want %s", i, got, wantSteps[i])
		}
		if got := requireString(t, step, "required_evidence"); got == "" {
			t.Fatalf("%s required_evidence missing", wantSteps[i])
		}
		if got := requireString(t, step, "blocking_if_missing"); got == "" {
			t.Fatalf("%s blocking_if_missing missing", wantSteps[i])
		}
	}

	roleRefs := map[string]bool{}
	for _, ref := range asStringSlice(t, candidateSet["role_pack_refs"]) {
		roleRefs[ref] = false
	}
	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	manifestFiles := map[string]bool{}
	for _, raw := range asObjectSlice(t, manifest["files"]) {
		manifestFiles[requireString(t, raw, "path")] = true
	}
	rawActivations := asObjectSlice(t, activation["required_role_activations"])
	if len(rawActivations) != len(roleRefs) {
		t.Fatalf("required_role_activations len = %d, want %d", len(rawActivations), len(roleRefs))
	}
	for _, item := range rawActivations {
		roleRef := requireString(t, item, "role_pack_ref")
		if _, ok := roleRefs[roleRef]; !ok {
			t.Fatalf("unexpected role_pack_ref %s", roleRef)
		}
		roleRefs[roleRef] = true
		roleFile := requireString(t, item, "role_pack_file")
		if !manifestFiles[roleFile] {
			t.Fatalf("%s role_pack_file %s missing from artifact manifest", roleRef, roleFile)
		}
		if got := requireString(t, item, "enabled_role_pack_version_ref"); !strings.HasPrefix(got, "enabled_role_pack_version://") {
			t.Fatalf("%s enabled_role_pack_version_ref = %s", roleRef, got)
		}
		requireString(t, item, "slot_ref")
		requireString(t, item, "install_receipt_ref")
		requireBool(t, item, "replaceable_in_team_settings", true)
		requireBool(t, item, "owner_gate_required_for_binding", true)
		requireBool(t, item, "runtime_usage_required", true)
	}
	for roleRef, seen := range roleRefs {
		if !seen {
			t.Fatalf("missing activation for %s", roleRef)
		}
	}

	negative := requireObject(t, activation, "negative_cases")
	for _, key := range []string{
		"install_without_entitlement",
		"missing_enabled_role_pack_version",
		"catalog_refresh_without_install_receipt",
		"runtime_use_without_team_binding",
		"role_output_formalization_without_owner_gate",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "install_runtime_activation_map_verified")

	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, verifier, "source_install_runtime_activation_map"); got != mapPath {
		t.Fatalf("source_install_runtime_activation_map = %s, want %s", got, mapPath)
	}
}

func TestTeamOfficeCapabilityRoleReferenceRequiresInstalledEntitledRoleActivation(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requirementsPath := "capability-role-requirements/sample-team-research.role-requirements.json"
	activationPath := "install/install-runtime-activation-map-candidate.json"

	requirements := readJSON(t, filepath.Join(base, requirementsPath))
	policy := requireObject(t, requirements, "post_install_reference_policy")
	requireBool(t, policy, "reference_available_only_after_install", true)
	if got := requireString(t, policy, "cloud_entitlement_truth_source"); got != "truzhen-cloud" {
		t.Fatalf("cloud_entitlement_truth_source = %s, want truzhen-cloud", got)
	}
	if got := requireString(t, policy, "install_truth_source"); got != "truzhenos" {
		t.Fatalf("install_truth_source = %s, want truzhenos", got)
	}
	if got := requireString(t, policy, "installed_role_activation_map"); got != activationPath {
		t.Fatalf("installed_role_activation_map = %s, want %s", got, activationPath)
	}
	for _, proof := range []string{
		"entitlement_verification_receipt",
		"install_receipt_ref",
		"enabled_role_pack_version_ref",
		"role_slot_compatibility_receipt",
		"role_body_not_copied_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, policy["required_before_capability_reference"]), proof)
	}
	for _, blocker := range []string{
		"candidate_json_only_without_install_receipt",
		"entitlement_missing_or_revoked",
		"enabled_role_pack_version_missing",
		"role_slot_incompatible",
		"capability_pack_attempts_to_copy_role_body",
	} {
		requireStringSliceContains(t, asStringSlice(t, policy["blocking_cases"]), blocker)
	}

	activation := readJSON(t, filepath.Join(base, activationPath))
	if got := requireString(t, activation, "capability_role_requirements"); got != requirementsPath {
		t.Fatalf("activation capability_role_requirements = %s, want %s", got, requirementsPath)
	}
	activationStep := findObjectByString(t, asObjectSlice(t, activation["activation_steps"]), "step_id", "capability_role_reference_ready")
	if got := requireString(t, activationStep, "actor"); got != "truzhenos + truzhen-client-web-desktop" {
		t.Fatalf("capability_role_reference_ready actor = %s", got)
	}
	for _, proof := range []string{"CapabilityRoleRequirement validation receipt", "enabled_role_pack_version_ref", "entitlement_verification_receipt"} {
		if !strings.Contains(requireString(t, activationStep, "required_evidence"), proof) {
			t.Fatalf("capability_role_reference_ready required_evidence missing %s", proof)
		}
	}

	activationsByRole := map[string]map[string]any{}
	for _, item := range asObjectSlice(t, activation["required_role_activations"]) {
		activationsByRole[requireString(t, item, "role_pack_ref")] = item
	}
	for _, raw := range asObjectSlice(t, requirements["role_requirements"]) {
		reqRoleRefs := asStringSlice(t, raw["accepted_role_pack_refs"])
		for _, roleRef := range reqRoleRefs {
			activation, ok := activationsByRole[roleRef]
			if !ok {
				t.Fatalf("role requirement references %s without required_role_activation", roleRef)
			}
			if got := requireString(t, activation, "enabled_role_pack_version_ref"); !strings.HasPrefix(got, "enabled_role_pack_version://") {
				t.Fatalf("%s enabled_role_pack_version_ref = %s", roleRef, got)
			}
		}
		requireBool(t, raw, "post_install_reference_required", true)
		requireStringIn(t, requireString(t, raw, "reference_scope"), "installed_buyer_team_only")
		for _, proof := range []string{"entitlement_verification_receipt", "install_receipt_ref", "enabled_role_pack_version_ref"} {
			requireStringSliceContains(t, asStringSlice(t, raw["activation_evidence_required"]), proof)
		}
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "capability_post_install_reference_gate_verified")
	localInstallGate := findObjectByString(t, asObjectSlice(t, matrix["readiness_gates"]), "gate_id", "local_install_and_team_binding")
	requireStringSliceContains(t, asStringSlice(t, localInstallGate["required_evidence"]), "capability_role_reference_after_install_receipt")

	chain := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	pass := requireObject(t, chain, "normal_commercialization_pass_conditions")
	if got := requireString(t, pass, "capability_role_reference_after_install_verified"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.capability_role_reference_after_install_verified = %s, want required", got)
	}
}

func TestTeamOfficeCommercialChainVerifierRequiresAuthoritativeReceiptsBeforeCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	verifierPath := "tests/commercial-chain-verifier-candidate.json"
	verifier := readJSON(t, filepath.Join(base, verifierPath))
	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	if got := requireString(t, verifier, "verifier_status"); got != "not_run_requires_cross_repo_receipts" {
		t.Fatalf("verifier_status = %s, want not_run_requires_cross_repo_receipts", got)
	}
	if got := requireString(t, verifier, "verifier_role"); got != "organizer_coordinator_recorder" {
		t.Fatalf("verifier_role = %s, want organizer_coordinator_recorder", got)
	}
	if got := requireString(t, verifier, "objective_source"); !strings.Contains(got, "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source = %s", got)
	}

	policy := requireObject(t, verifier, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	forbiddenClaims := strings.Join(asStringSlice(t, policy["non_sufficient_completion_evidence"]), "\n")
	for _, item := range []string{"candidate_json_only", "current_repo_tests_only", "sandbox_design_without_receipts", "organizer_self_attestation_only"} {
		if !strings.Contains(forbiddenClaims, item) {
			t.Fatalf("non_sufficient_completion_evidence missing %s", item)
		}
	}

	correlation := requireObject(t, verifier, "correlation_requirements")
	keys := strings.Join(asStringSlice(t, correlation["required_correlation_keys"]), "\n")
	for _, key := range []string{
		"candidate_set_ref",
		"artifact_ref",
		"artifact_sha256",
		"upload_receipt_ref",
		"listing_draft_ref",
		"sandbox_order_ref",
		"sandbox_payment_receipt_ref",
		"entitlement_ref",
		"download_receipt_ref",
		"install_receipt_ref",
		"team_binding_receipt_ref",
	} {
		if !strings.Contains(keys, key) {
			t.Fatalf("required_correlation_keys missing %s", key)
		}
	}
	if got := requireString(t, correlation, "hash_continuity"); got != "upload_download_install_hashes_must_match" {
		t.Fatalf("hash_continuity = %s", got)
	}

	rawStages, ok := verifier["ordered_verification_gates"].([]any)
	if !ok {
		t.Fatalf("ordered_verification_gates missing")
	}
	requiredStages := map[string]bool{
		"role_candidate_export_gui":      false,
		"cloud_upload_listing_draft":     false,
		"marketplace_review_candidate":   false,
		"sandbox_order_payment":          false,
		"license_entitlement_issue":      false,
		"entitled_download_hash":         false,
		"local_install_enabled_version":  false,
		"team_settings_role_replacement": false,
		"runtime_use_candidate_only":     false,
		"negative_cases_blocked":         false,
		"independent_acceptance_signoff": false,
	}
	lastOrder := 0.0
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("verification gate = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := requiredStages[stageID]; ok {
			requiredStages[stageID] = true
		}
		order, ok := stage["stage_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s stage_order = %v, want increasing number greater than %v", stageID, stage["stage_order"], lastOrder)
		}
		lastOrder = order
		requireStringIn(t, requireString(t, stage, "truth_source"), "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo", "independent_acceptance_agent")
		requireStringIn(t, requireString(t, stage, "stage_status"), "pending_cross_repo_execution", "blocked_until_prior_stage_evidence")
		evidence := strings.Join(asStringSlice(t, stage["required_authoritative_evidence"]), "\n")
		if !strings.Contains(evidence, "gui_screenshot") && !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "signoff") {
			t.Fatalf("%s required_authoritative_evidence must include GUI, receipt, or signoff evidence", stageID)
		}
		if strings.Contains(stageID, "cloud") || strings.Contains(stageID, "review") || strings.Contains(stageID, "payment") || strings.Contains(stageID, "entitlement") || strings.Contains(stageID, "download") {
			if !strings.Contains(evidence, "truzhen-cloud") {
				t.Fatalf("%s required_authoritative_evidence missing truzhen-cloud evidence", stageID)
			}
		}
		if strings.Contains(stageID, "install") || strings.Contains(stageID, "team_settings") || strings.Contains(stageID, "runtime") {
			if !strings.Contains(evidence, "truzhenos") {
				t.Fatalf("%s required_authoritative_evidence missing truzhenos evidence", stageID)
			}
		}
		if got := requireString(t, stage, "failure_issue_route"); got == "" {
			t.Fatalf("%s failure_issue_route missing", stageID)
		}
	}
	for stageID, seen := range requiredStages {
		if !seen {
			t.Fatalf("missing verification gate %s", stageID)
		}
	}

	pass := requireObject(t, verifier, "normal_commercialization_pass_conditions")
	for _, key := range []string{
		"all_gui_steps_have_user_view_evidence",
		"cloud_receipt_chain_complete",
		"artifact_hash_consistent_upload_download_install",
		"entitlement_valid_at_download_and_install",
		"truzhenos_install_receipt_present",
		"team_binding_owner_gate_receipt_present",
		"negative_cases_blocked_with_receipts",
		"independent_acceptance_signed",
	} {
		if got := requireString(t, pass, key); got != "required" {
			t.Fatalf("normal_commercialization_pass_conditions.%s = %s, want required", key, got)
		}
	}

	negative := requireObject(t, verifier, "blocking_cases")
	for _, key := range []string{
		"missing_gui_screenshot",
		"missing_cloud_upload_receipt",
		"payment_without_sandbox_receipt",
		"download_without_entitlement",
		"artifact_hash_mismatch",
		"install_without_truzhenos_receipt",
		"team_binding_without_owner_gate",
		"missing_independent_acceptance",
	} {
		item := requireObject(t, negative, key)
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", key, got)
		}
		if got := requireString(t, item, "expected_evidence"); got == "" {
			t.Fatalf("%s expected_evidence missing", key)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_chain_verifier"); got != verifierPath {
		t.Fatalf("commercial_chain_verifier = %s, want %s", got, verifierPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, verifierPath) {
		t.Fatalf("candidate set artifact_files missing %s", verifierPath)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	readinessPolicy := requireObject(t, matrix, "completion_claim_policy")
	requiredProofs := strings.Join(asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "commercial_chain_verifier_passed") {
		t.Fatalf("product readiness matrix missing commercial_chain_verifier_passed completion proof")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("file item = %T", raw)
		}
		if requireString(t, item, "path") == verifierPath {
			requireStringIn(t, requireString(t, item, "required_for"), "audit")
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", verifierPath)
}

func TestTeamOfficeCommercialChainVerifierRequiresP11ExecutionQueuePhaseDependencies(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	request := readJSON(t, filepath.Join(base, requestPath))
	contract := requireObject(t, request, "phase_dependency_contract")
	handoffs := asObjectSlice(t, contract["phase_handoffs"])

	if got := requireString(t, verifier, "source_execution_queue"); got != queuePath {
		t.Fatalf("source_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, verifier, "source_p11_phase_dependency_contract"); got != requestPath+"#phase_dependency_contract" {
		t.Fatalf("source_p11_phase_dependency_contract = %s, want %s", got, requestPath+"#phase_dependency_contract")
	}
	verification := requireObject(t, verifier, "phase_dependency_verification")
	for key, want := range map[string]string{
		"p11_phase_dependency_links_ref": queuePath + "#p11_phase_dependency_links",
		"required_completion_proof":      "p11_execution_queue_phase_dependencies_verified",
		"missing_dependency_status":      "blocked_previous_phase_evidence_missing",
	} {
		if got := requireString(t, verification, key); got != want {
			t.Fatalf("phase_dependency_verification.%s = %s, want %s", key, got, want)
		}
	}
	for _, key := range []string{"all_run_request_handoffs_must_be_verified", "no_stage_can_pass_if_previous_receipts_missing"} {
		requireBool(t, verification, key, true)
	}
	if got, ok := verification["expected_phase_dependency_count"].(float64); !ok || int(got) != len(handoffs) {
		t.Fatalf("expected_phase_dependency_count = %v, want %d", verification["expected_phase_dependency_count"], len(handoffs))
	}

	pass := requireObject(t, verifier, "normal_commercialization_pass_conditions")
	if got := requireString(t, pass, "p11_execution_queue_phase_dependencies_verified"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.p11_execution_queue_phase_dependencies_verified = %s, want required", got)
	}
	related := strings.Join(asStringSlice(t, verifier["related_artifacts"]), "\n")
	for _, artifact := range []string{requestPath, queuePath} {
		if !strings.Contains(related, artifact) {
			t.Fatalf("related_artifacts missing %s", artifact)
		}
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	required := asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"])
	requireStringSliceContains(t, required, "p11_execution_queue_phase_dependencies_verified")
}

func TestTeamOfficeArtifactBundleDigestDefinesUploadDownloadInstallHashContinuity(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	digestPath := "commerce/artifact-bundle-digest-candidate.json"
	digest := readJSON(t, filepath.Join(base, digestPath))
	requireBool(t, digest, "candidate_only", true)
	requireBool(t, digest, "non_formal", true)
	if got := requireString(t, digest, "digest_status"); got != "current_repo_candidate_digest_not_cloud_receipt" {
		t.Fatalf("digest_status = %s, want current_repo_candidate_digest_not_cloud_receipt", got)
	}
	if got := requireString(t, digest, "hash_algorithm"); got != "sha256" {
		t.Fatalf("hash_algorithm = %s, want sha256", got)
	}
	if got := requireString(t, digest, "payload_hash_policy"); got != "manifest_and_digest_self_excluded" {
		t.Fatalf("payload_hash_policy = %s, want manifest_and_digest_self_excluded", got)
	}
	requireBool(t, digest, "cloud_upload_receipt_required", true)
	requireBool(t, digest, "download_install_must_match_bundle_tree_sha256", true)

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	manifestPayload := map[string]string{}
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		path := requireString(t, item, "path")
		if path == "commerce/artifact-manifest.json" || path == digestPath {
			continue
		}
		manifestPayload[path] = requireString(t, item, "sha256")
	}

	rawPayload, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	if len(rawPayload) != len(manifestPayload) {
		t.Fatalf("payload_file_hashes len = %d, want manifest payload len %d", len(rawPayload), len(manifestPayload))
	}
	var treeInput strings.Builder
	for _, raw := range rawPayload {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload item = %T", raw)
		}
		path := requireString(t, item, "path")
		gotHash := requireString(t, item, "sha256")
		wantHash, ok := manifestPayload[path]
		if !ok {
			t.Fatalf("payload file %s not found in manifest payload", path)
		}
		if gotHash != wantHash {
			t.Fatalf("%s payload hash = %s, want manifest hash %s", path, gotHash, wantHash)
		}
		data, err := os.ReadFile(filepath.Join(base, path))
		if err != nil {
			t.Fatalf("read payload %s: %v", path, err)
		}
		actual := fmt.Sprintf("%x", sha256.Sum256(data))
		if gotHash != actual {
			t.Fatalf("%s payload hash = %s, want actual hash %s", path, gotHash, actual)
		}
		treeInput.WriteString(path)
		treeInput.WriteByte('\x00')
		treeInput.WriteString(gotHash)
		treeInput.WriteByte('\n')
		delete(manifestPayload, path)
	}
	if len(manifestPayload) != 0 {
		t.Fatalf("payload_file_hashes missing manifest files: %v", manifestPayload)
	}
	wantTree := fmt.Sprintf("%x", sha256.Sum256([]byte(treeInput.String())))
	if got := requireString(t, digest, "bundle_tree_sha256"); got != wantTree {
		t.Fatalf("bundle_tree_sha256 = %s, want %s", got, wantTree)
	}

	for _, key := range []string{"cloud_upload_request_ref", "download_hash_ref", "install_hash_ref"} {
		if got := requireString(t, digest, key); got == "" {
			t.Fatalf("%s missing", key)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "artifact_bundle_digest"); got != digestPath {
		t.Fatalf("artifact_bundle_digest = %s, want %s", got, digestPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, digestPath) {
		t.Fatalf("candidate set artifact_files missing %s", digestPath)
	}

	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	correlation := requireObject(t, chainVerifier, "correlation_requirements")
	keys := strings.Join(asStringSlice(t, correlation["required_correlation_keys"]), "\n")
	if !strings.Contains(keys, "bundle_tree_sha256") {
		t.Fatalf("commercial chain verifier missing bundle_tree_sha256 correlation key")
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	readinessPolicy := requireObject(t, matrix, "completion_claim_policy")
	requiredProofs := strings.Join(asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "artifact_bundle_digest_verified") {
		t.Fatalf("product readiness matrix missing artifact_bundle_digest_verified completion proof")
	}

	rawManifestFiles, ok = manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") == digestPath {
			requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
			data, err := os.ReadFile(filepath.Join(base, digestPath))
			if err != nil {
				t.Fatalf("read %s: %v", digestPath, err)
			}
			wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
			if got := requireString(t, item, "sha256"); got != wantHash {
				t.Fatalf("%s manifest sha256 = %s, want %s", digestPath, got, wantHash)
			}
			return
		}
	}
	t.Fatalf("artifact manifest missing %s", digestPath)
}

func TestTeamOfficeCommercialDistributionReceiptSchemaDefinesUploadPurchaseDownloadInstallReceipts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	schemaPath := "commerce/commercial-distribution-receipt-schema-candidate.json"
	schema := readJSON(t, filepath.Join(base, schemaPath))
	requireBool(t, schema, "candidate_only", true)
	requireBool(t, schema, "non_formal", true)
	if got := requireString(t, schema, "schema_status"); got != "candidate_not_run_requires_cross_repo_receipts" {
		t.Fatalf("schema_status = %s, want candidate_not_run_requires_cross_repo_receipts", got)
	}
	if got := requireString(t, schema, "bundle_tree_sha256_source"); got != "commerce/artifact-bundle-digest-candidate.json" {
		t.Fatalf("bundle_tree_sha256_source = %s", got)
	}

	truthSources := requireObject(t, schema, "truth_sources")
	for key, want := range map[string]string{
		"cloud":         "truzhen-cloud",
		"install":       "truzhenos",
		"frontend":      "truzhen-client-web-desktop",
		"acceptance":    "independent_acceptance_agent",
		"candidate_set": "truzhen-packs",
	} {
		if got := requireString(t, truthSources, key); got != want {
			t.Fatalf("truth_sources.%s = %s, want %s", key, got, want)
		}
	}

	recordSchema := requireObject(t, schema, "record_schema")
	requiredFields := strings.Join(asStringSlice(t, recordSchema["required_fields"]), "\n")
	for _, field := range []string{
		"receipt_ref",
		"stage_id",
		"actor_ref",
		"truth_source",
		"timestamp",
		"correlation_id",
		"candidate_set_ref",
		"artifact_ref",
		"bundle_tree_sha256",
		"artifact_sha256",
		"previous_receipt_ref",
		"result_status",
		"blocked_reason",
		"redaction_status",
	} {
		if !strings.Contains(requiredFields, field) {
			t.Fatalf("record_schema.required_fields missing %s", field)
		}
	}

	rawTypes, ok := schema["required_receipt_types"].([]any)
	if !ok {
		t.Fatalf("required_receipt_types missing")
	}
	expectedTypes := map[string]string{
		"candidate_bundle_export_receipt": "truzhenos",
		"cloud_upload_receipt":            "truzhen-cloud",
		"listing_review_receipt":          "truzhen-cloud",
		"sandbox_order_receipt":           "truzhen-cloud",
		"sandbox_payment_receipt":         "truzhen-cloud",
		"license_entitlement_receipt":     "truzhen-cloud",
		"signed_download_receipt":         "truzhen-cloud",
		"install_preflight_receipt":       "truzhenos",
		"install_receipt":                 "truzhenos",
		"team_binding_receipt":            "truzhenos",
		"runtime_candidate_receipt":       "truzhenos",
		"negative_block_receipt":          "multi_repo",
		"independent_acceptance_signoff":  "independent_acceptance_agent",
	}
	seenTypes := map[string]bool{}
	for _, raw := range rawTypes {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("required receipt type = %T", raw)
		}
		receiptType := requireString(t, item, "receipt_type")
		wantTruth, ok := expectedTypes[receiptType]
		if !ok {
			t.Fatalf("unexpected receipt_type %s", receiptType)
		}
		seenTypes[receiptType] = true
		if got := requireString(t, item, "truth_source"); got != wantTruth {
			t.Fatalf("%s truth_source = %s, want %s", receiptType, got, wantTruth)
		}
		requireBool(t, item, "must_reference_bundle_tree_sha256", true)
		fields := strings.Join(asStringSlice(t, item["required_fields"]), "\n")
		for _, field := range []string{"bundle_tree_sha256", "correlation_id", "result_status"} {
			if !strings.Contains(fields, field) {
				t.Fatalf("%s required_fields missing %s", receiptType, field)
			}
		}
	}
	for receiptType := range expectedTypes {
		if !seenTypes[receiptType] {
			t.Fatalf("missing receipt_type %s", receiptType)
		}
	}

	rawStages, ok := schema["stage_receipt_chain"].([]any)
	if !ok {
		t.Fatalf("stage_receipt_chain missing")
	}
	expectedStages := map[string]bool{
		"candidate_bundle_export":        false,
		"cloud_upload":                   false,
		"listing_review":                 false,
		"sandbox_order":                  false,
		"sandbox_payment":                false,
		"license_entitlement":            false,
		"signed_download":                false,
		"install_preflight":              false,
		"install":                        false,
		"team_binding":                   false,
		"runtime_candidate_use":          false,
		"negative_block_cases":           false,
		"independent_acceptance_signoff": false,
	}
	lastOrder := 0.0
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage receipt = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := expectedStages[stageID]; ok {
			expectedStages[stageID] = true
		}
		order, ok := stage["stage_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s stage_order = %v, want increasing number greater than %v", stageID, stage["stage_order"], lastOrder)
		}
		lastOrder = order
		if got := requireString(t, stage, "correlation_policy"); got != "same_correlation_id_across_chain" {
			t.Fatalf("%s correlation_policy = %s", stageID, got)
		}
		if got := requireString(t, stage, "required_receipt_type"); got == "" {
			t.Fatalf("%s required_receipt_type missing", stageID)
		}
	}
	for stageID, seen := range expectedStages {
		if !seen {
			t.Fatalf("missing stage_receipt_chain stage %s", stageID)
		}
	}

	rawNegative, ok := schema["negative_receipt_requirements"].([]any)
	if !ok {
		t.Fatalf("negative_receipt_requirements missing")
	}
	expectedNegative := map[string]bool{
		"download_without_entitlement":             false,
		"artifact_hash_mismatch":                   false,
		"install_without_receipt":                  false,
		"team_binding_without_owner_gate":          false,
		"real_payment_without_owner_authorization": false,
	}
	for _, raw := range rawNegative {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("negative receipt requirement = %T", raw)
		}
		caseID := requireString(t, item, "case_id")
		if _, ok := expectedNegative[caseID]; ok {
			expectedNegative[caseID] = true
		}
		if got := requireString(t, item, "expected_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s expected_status = %s, want blocked*", caseID, got)
		}
		if got := requireString(t, item, "required_receipt_type"); got != "negative_block_receipt" {
			t.Fatalf("%s required_receipt_type = %s, want negative_block_receipt", caseID, got)
		}
	}
	for caseID, seen := range expectedNegative {
		if !seen {
			t.Fatalf("missing negative receipt requirement %s", caseID)
		}
	}

	policy := requireObject(t, schema, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, policy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"all_required_receipts_present",
		"bundle_tree_sha256_consistent",
		"negative_block_receipts_present",
		"independent_acceptance_signoff",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_distribution_receipt_schema"); got != schemaPath {
		t.Fatalf("commercial_distribution_receipt_schema = %s, want %s", got, schemaPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, schemaPath) {
		t.Fatalf("candidate set artifact_files missing %s", schemaPath)
	}

	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	correlation := requireObject(t, chainVerifier, "correlation_requirements")
	keys := strings.Join(asStringSlice(t, correlation["required_correlation_keys"]), "\n")
	if !strings.Contains(keys, "distribution_receipt_schema_ref") {
		t.Fatalf("commercial chain verifier missing distribution_receipt_schema_ref correlation key")
	}
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, schemaPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", schemaPath)
	}

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	readinessPolicy := requireObject(t, matrix, "completion_claim_policy")
	requiredProofs := strings.Join(asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredProofs, "commercial_distribution_receipt_schema_verified") {
		t.Fatalf("product readiness matrix missing commercial_distribution_receipt_schema_verified completion proof")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	schemaInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != schemaPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "upload", "download", "install", "audit")
		data, err := os.ReadFile(filepath.Join(base, schemaPath))
		if err != nil {
			t.Fatalf("read %s: %v", schemaPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", schemaPath, got, wantHash)
		}
		schemaInManifest = true
	}
	if !schemaInManifest {
		t.Fatalf("artifact manifest missing %s", schemaPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == schemaPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", schemaPath)
}

func TestTeamOfficeRoleStudioPhaseCoverageMatrixMapsP0ToP11ToEvidenceAndBlockers(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrixPath := "tests/role-studio-phase-coverage-matrix-candidate.json"
	matrix := readJSON(t, filepath.Join(base, matrixPath))
	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	if got := requireString(t, matrix, "coverage_status"); got != "candidate_coverage_defined_pending_cross_repo_execution" {
		t.Fatalf("coverage_status = %s, want candidate_coverage_defined_pending_cross_repo_execution", got)
	}
	if got := requireString(t, matrix, "objective_source"); !strings.Contains(got, "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source = %s", got)
	}

	rawPhases, ok := matrix["phase_coverage"].([]any)
	if !ok {
		t.Fatalf("phase_coverage missing")
	}
	expectedPhases := map[string]bool{
		"P0":  false,
		"P1":  false,
		"P2":  false,
		"P3":  false,
		"P4":  false,
		"P5":  false,
		"P6":  false,
		"P7":  false,
		"P8":  false,
		"P9":  false,
		"P10": false,
		"P11": false,
	}
	requiredTruth := map[string][]string{
		"P0":  {"truzhen-packs"},
		"P1":  {"truzhen-client-web-desktop", "truzhenos"},
		"P2":  {"truzhen-client-web-desktop", "truzhenos"},
		"P3":  {"truzhenos", "truzhen-packs"},
		"P4":  {"truzhen-client-web-desktop", "truzhenos"},
		"P5":  {"truzhen-client-web-desktop", "truzhenos"},
		"P6":  {"truzhen-client-web-desktop", "truzhenos"},
		"P7":  {"truzhen-client-web-desktop", "truzhenos"},
		"P8":  {"multi_repo"},
		"P9":  {"truzhen-packs"},
		"P10": {"truzhen-client-web-desktop", "truzhenos"},
		"P11": {"truzhen-cloud", "truzhenos"},
	}
	lastOrder := 0.0
	for _, raw := range rawPhases {
		phase, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("phase coverage = %T", raw)
		}
		phaseID := requireString(t, phase, "phase_id")
		if _, ok := expectedPhases[phaseID]; !ok {
			t.Fatalf("unexpected phase_id %s", phaseID)
		}
		expectedPhases[phaseID] = true
		order, ok := phase["phase_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s phase_order = %v, want increasing number greater than %v", phaseID, phase["phase_order"], lastOrder)
		}
		lastOrder = order
		if got := requireString(t, phase, "coverage_status"); got != "pending_cross_repo_execution" && got != "pending_owner_authorization" {
			t.Fatalf("%s coverage_status = %s", phaseID, got)
		}
		for _, key := range []string{
			"required_user_view_evidence",
			"required_backend_evidence",
			"current_repo_candidate_artifacts",
			"missing_authoritative_evidence",
			"target_repos",
		} {
			values := asStringSlice(t, phase[key])
			if len(values) == 0 {
				t.Fatalf("%s %s must not be empty", phaseID, key)
			}
		}
		truthSources := strings.Join(asStringSlice(t, phase["truth_sources"]), "\n")
		for _, wantTruth := range requiredTruth[phaseID] {
			if !strings.Contains(truthSources, wantTruth) {
				t.Fatalf("%s truth_sources missing %s", phaseID, wantTruth)
			}
		}
		if got := requireString(t, phase, "blocking_status"); !strings.HasPrefix(got, "blocked") {
			t.Fatalf("%s blocking_status = %s, want blocked*", phaseID, got)
		}
		requireBool(t, phase, "must_not_claim_complete", true)
	}
	for phaseID, seen := range expectedPhases {
		if !seen {
			t.Fatalf("missing phase coverage %s", phaseID)
		}
	}

	claimPolicy := requireObject(t, matrix, "completion_claim_policy")
	requireBool(t, claimPolicy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"all_p0_to_p11_phases_have_authoritative_evidence",
		"user_view_gui_agent_evidence_present",
		"truzhenos_receipts_present",
		"truzhen_cloud_receipts_present",
		"independent_acceptance_signoff",
		"no_direct_json_edit_claim",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "role_studio_phase_coverage_matrix"); got != matrixPath {
		t.Fatalf("role_studio_phase_coverage_matrix = %s, want %s", got, matrixPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, matrixPath) {
		t.Fatalf("candidate set artifact_files missing %s", matrixPath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "role_studio_phase_coverage_matrix_verified") {
		t.Fatalf("product readiness matrix missing role_studio_phase_coverage_matrix_verified")
	}

	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "phase_coverage_matrix_ref"); got != matrixPath {
		t.Fatalf("phase_coverage_matrix_ref = %s, want %s", got, matrixPath)
	}

	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, matrixPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", matrixPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	matrixInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != matrixPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, matrixPath))
		if err != nil {
			t.Fatalf("read %s: %v", matrixPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", matrixPath, got, wantHash)
		}
		matrixInManifest = true
	}
	if !matrixInManifest {
		t.Fatalf("artifact manifest missing %s", matrixPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == matrixPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", matrixPath)
}

func TestTeamOfficeRoleStudioTestCaseCoverageMatrixMapsPlanTCsToEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrixPath := "tests/role-studio-test-case-coverage-matrix-candidate.json"
	matrix := readJSON(t, filepath.Join(base, matrixPath))
	requireBool(t, matrix, "candidate_only", true)
	requireBool(t, matrix, "non_formal", true)
	if got := requireString(t, matrix, "coverage_status"); got != "candidate_coverage_defined_pending_cross_repo_execution" {
		t.Fatalf("coverage_status = %s, want candidate_coverage_defined_pending_cross_repo_execution", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":             "role-pack-candidate-set://team-office-v0",
		"source_plan":                   "/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/role-pack-studio-team-office-test-plan-20260704.md",
		"phase_coverage_matrix_ref":     "tests/role-studio-phase-coverage-matrix-candidate.json",
		"gui_execution_script_ref":      "tests/gui-user-agent-execution-script-candidate.json",
		"evidence_capture_protocol_ref": "tests/gui-evidence-capture-protocol.json",
		"commercial_evidence_gate_ref":  "tests/commercial-evidence-gate-candidate.json",
		"p11_acceptance_gate_ref":       "tests/p11-normal-commercialization-acceptance-gate-candidate.json",
	} {
		if got := requireString(t, matrix, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawRows, ok := matrix["test_case_coverage"].([]any)
	if !ok {
		t.Fatalf("test_case_coverage missing")
	}
	expectedTCs := map[string]string{
		"TC-ROLE-01":        "P2",
		"TC-ROLE-02":        "P2",
		"TC-ROLE-03":        "P8",
		"TC-GUI-USER-01":    "P1",
		"TC-EXPORT-01":      "P3",
		"TC-CAP-ROLE-01":    "P4",
		"TC-CAP-ROLE-02":    "P8",
		"TC-TEAM-01":        "P5",
		"TC-TEAM-02":        "P5",
		"TC-TEAM-03":        "P5",
		"TC-VOICE-01":       "P6",
		"TC-VRM-01":         "P6",
		"TC-NEG-01":         "P8",
		"TC-NEG-02":         "P8",
		"TC-NEG-03":         "P8",
		"TC-PRODUCT-FE-01":  "P10",
		"TC-PRODUCT-BE-01":  "P10",
		"TC-CLOUD-01":       "P11",
		"TC-CLOUD-02":       "P11",
		"TC-PAY-01":         "P11",
		"TC-DOWNLOAD-01":    "P11",
		"TC-INSTALL-01":     "P11",
		"TC-INSTALL-NEG-01": "P11",
		"TC-INSTALL-NEG-02": "P11",
	}
	if len(rawRows) != len(expectedTCs) {
		t.Fatalf("test_case_coverage len = %d, want %d", len(rawRows), len(expectedTCs))
	}
	seen := map[string]bool{}
	for _, raw := range rawRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("test case row = %T", raw)
		}
		tcID := requireString(t, row, "test_case_id")
		wantPhase, ok := expectedTCs[tcID]
		if !ok {
			t.Fatalf("unexpected test_case_id %s", tcID)
		}
		seen[tcID] = true
		if got := requireString(t, row, "plan_phase"); got != wantPhase {
			t.Fatalf("%s plan_phase = %s, want %s", tcID, got, wantPhase)
		}
		if got := requireString(t, row, "execution_mode"); got != "user_view_gui_agent_only" {
			t.Fatalf("%s execution_mode = %s, want user_view_gui_agent_only", tcID, got)
		}
		if got := requireString(t, row, "coverage_status"); got != "pending_cross_repo_execution" && got != "pending_owner_authorization" {
			t.Fatalf("%s coverage_status = %s", tcID, got)
		}
		requireBool(t, row, "can_count_toward_completion_now", false)
		for _, key := range []string{
			"target_repos",
			"required_gui_evidence",
			"required_authoritative_evidence",
			"current_repo_candidate_artifacts",
			"missing_authoritative_evidence",
		} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", tcID, key)
			}
		}
		evidence := strings.Join(asStringSlice(t, row["required_authoritative_evidence"]), "\n")
		if !strings.Contains(evidence, "receipt") && !strings.Contains(evidence, "blocked") && !strings.Contains(evidence, "signoff") {
			t.Fatalf("%s required_authoritative_evidence must include receipt, blocked, or signoff evidence", tcID)
		}
	}
	for tcID := range expectedTCs {
		if !seen[tcID] {
			t.Fatalf("missing test case coverage %s", tcID)
		}
	}

	completion := requireObject(t, matrix, "completion_claim_policy")
	requireBool(t, completion, "completion_claim_allowed", false)
	requiredBeforeCompletion := strings.Join(asStringSlice(t, completion["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"all_plan_test_cases_have_authoritative_evidence",
		"user_view_gui_agent_operation_log_present",
		"cloud_install_and_binding_receipts_present",
		"all_negative_cases_blocked_verified",
		"independent_acceptance_signed",
	} {
		if !strings.Contains(requiredBeforeCompletion, proof) {
			t.Fatalf("completion_claim_policy missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "role_studio_test_case_coverage_matrix"); got != matrixPath {
		t.Fatalf("role_studio_test_case_coverage_matrix = %s, want %s", got, matrixPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, matrixPath) {
		t.Fatalf("candidate set artifact_files missing %s", matrixPath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "role_studio_test_case_coverage_matrix_verified") {
		t.Fatalf("product readiness matrix missing role_studio_test_case_coverage_matrix_verified")
	}

	acceptanceGate := readJSON(t, filepath.Join(base, "tests", "p11-normal-commercialization-acceptance-gate-candidate.json"))
	if got := requireString(t, acceptanceGate, "test_case_coverage_matrix_ref"); got != matrixPath {
		t.Fatalf("p11 acceptance gate test_case_coverage_matrix_ref = %s, want %s", got, matrixPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	matrixInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != matrixPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, matrixPath))
		if err != nil {
			t.Fatalf("read %s: %v", matrixPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", matrixPath, got, wantHash)
		}
		matrixInManifest = true
	}
	if !matrixInManifest {
		t.Fatalf("artifact manifest missing %s", matrixPath)
	}
}

func TestTeamOfficeIndependentAcceptanceSignoffMatrixRequiresExternalReviewBeforeCommercialCompletion(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	signoffPath := "tests/independent-acceptance-signoff-matrix-candidate.json"
	signoff := readJSON(t, filepath.Join(base, signoffPath))
	requireBool(t, signoff, "candidate_only", true)
	requireBool(t, signoff, "non_formal", true)
	if got := requireString(t, signoff, "signoff_status"); got != "candidate_not_run_requires_authoritative_evidence" {
		t.Fatalf("signoff_status = %s, want candidate_not_run_requires_authoritative_evidence", got)
	}
	if got := requireString(t, signoff, "objective_source"); !strings.Contains(got, "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source = %s", got)
	}

	reviewerPolicy := requireObject(t, signoff, "reviewer_policy")
	requireBool(t, reviewerPolicy, "independent_acceptance_agent_required", true)
	requireBool(t, reviewerPolicy, "organizer_self_attestation_allowed_for_completion", false)
	requireBool(t, reviewerPolicy, "user_view_gui_evidence_required", true)
	requireBool(t, reviewerPolicy, "cross_repo_receipts_required", true)

	actors := strings.Join(asStringSlice(t, signoff["required_actor_separation"]), "\n")
	for _, actor := range []string{"user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent", "owner_go_no_go_decider"} {
		if !strings.Contains(actors, actor) {
			t.Fatalf("required_actor_separation missing %s", actor)
		}
	}

	rawRows, ok := signoff["stage_signoff_rows"].([]any)
	if !ok {
		t.Fatalf("stage_signoff_rows missing")
	}
	expectedPhases := map[string]bool{
		"P0":  false,
		"P1":  false,
		"P2":  false,
		"P3":  false,
		"P4":  false,
		"P5":  false,
		"P6":  false,
		"P7":  false,
		"P8":  false,
		"P9":  false,
		"P10": false,
		"P11": false,
	}
	lastOrder := 0.0
	for _, raw := range rawRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage signoff row = %T", raw)
		}
		phaseID := requireString(t, row, "phase_id")
		if _, ok := expectedPhases[phaseID]; !ok {
			t.Fatalf("unexpected phase_id %s", phaseID)
		}
		expectedPhases[phaseID] = true
		order, ok := row["phase_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s phase_order = %v, want increasing number greater than %v", phaseID, row["phase_order"], lastOrder)
		}
		lastOrder = order
		if got := requireString(t, row, "acceptance_status"); got != "not_run" {
			t.Fatalf("%s acceptance_status = %s, want not_run", phaseID, got)
		}
		requireBool(t, row, "must_reference_phase_coverage_matrix", true)
		requireBool(t, row, "must_reference_gui_evidence", true)
		requireBool(t, row, "must_reference_backend_receipt", true)
		for _, key := range []string{"required_review_evidence", "blocking_if_missing", "target_repos"} {
			values := asStringSlice(t, row[key])
			if len(values) == 0 {
				t.Fatalf("%s %s must not be empty", phaseID, key)
			}
		}
		if phaseID == "P11" {
			requireBool(t, row, "must_reference_distribution_receipt_schema", true)
			evidence := strings.Join(asStringSlice(t, row["required_review_evidence"]), "\n")
			for _, proof := range []string{"cloud_upload_receipt", "sandbox_payment_receipt", "download_receipt", "install_receipt"} {
				if !strings.Contains(evidence, proof) {
					t.Fatalf("P11 required_review_evidence missing %s", proof)
				}
			}
		}
	}
	for phaseID, seen := range expectedPhases {
		if !seen {
			t.Fatalf("missing stage signoff row %s", phaseID)
		}
	}

	decisions := strings.Join(asStringSlice(t, signoff["allowed_acceptance_decisions"]), "\n")
	for _, decision := range []string{"passed_verified", "failed_requires_issue", "blocked_missing_evidence", "not_run"} {
		if !strings.Contains(decisions, decision) {
			t.Fatalf("allowed_acceptance_decisions missing %s", decision)
		}
	}

	blockingIssues := strings.Join(asStringSlice(t, signoff["blocking_issue_codes"]), "\n")
	for _, issue := range []string{
		"missing_gui_evidence",
		"missing_cloud_receipt",
		"missing_install_receipt",
		"missing_team_binding_receipt",
		"missing_negative_block_receipt",
		"missing_hash_continuity",
		"organizer_self_attestation_only",
	} {
		if !strings.Contains(blockingIssues, issue) {
			t.Fatalf("blocking_issue_codes missing %s", issue)
		}
	}

	claimPolicy := requireObject(t, signoff, "completion_claim_policy")
	requireBool(t, claimPolicy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"all_stage_signoffs_passed",
		"all_negative_cases_reviewed",
		"independent_acceptance_agent_receipt",
		"owner_go_no_go_decision",
		"no_open_p0_p11_blockers",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "independent_acceptance_signoff_matrix"); got != signoffPath {
		t.Fatalf("independent_acceptance_signoff_matrix = %s, want %s", got, signoffPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, signoffPath) {
		t.Fatalf("candidate set artifact_files missing %s", signoffPath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "independent_acceptance_signoff_matrix_verified") {
		t.Fatalf("product readiness matrix missing independent_acceptance_signoff_matrix_verified")
	}

	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "independent_acceptance_signoff_matrix_ref"); got != signoffPath {
		t.Fatalf("independent_acceptance_signoff_matrix_ref = %s, want %s", got, signoffPath)
	}

	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, signoffPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", signoffPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	signoffInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != signoffPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, signoffPath))
		if err != nil {
			t.Fatalf("read %s: %v", signoffPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", signoffPath, got, wantHash)
		}
		signoffInManifest = true
	}
	if !signoffInManifest {
		t.Fatalf("artifact manifest missing %s", signoffPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == signoffPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", signoffPath)
}

func TestTeamOfficeIndependentAcceptanceSignoffRequiresGuiAPITraceabilityReview(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	signoffPath := "tests/independent-acceptance-signoff-matrix-candidate.json"
	signoff := readJSON(t, filepath.Join(base, signoffPath))

	p11Row := findObjectByString(t, asObjectSlice(t, signoff["stage_signoff_rows"]), "phase_id", "P11")
	for _, evidence := range []string{
		"gui_api_traceability_matrix_verified",
		"commercial_api_endpoint_refs",
		"required_receipt_or_result_fields",
		"required_gui_evidence_slots",
	} {
		requireStringSliceContains(t, asStringSlice(t, p11Row["required_review_evidence"]), evidence)
	}
	for _, blocker := range []string{
		"gui_api_traceability_matrix_missing",
		"commercial_api_endpoint_trace_missing",
		"receipt_field_trace_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, p11Row["blocking_if_missing"]), blocker)
	}
	requireStringSliceContains(t, asStringSlice(t, signoff["blocking_issue_codes"]), "missing_gui_api_traceability_matrix")

	claimPolicy := requireObject(t, signoff, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, claimPolicy["required_before_completion_claim"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, claimPolicy["non_sufficient_evidence"]), "independent_acceptance_without_gui_api_traceability")

	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))
	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "independent_acceptance_signoff")
	requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), signoffPath)
	requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), "tests/p11-evidence-acceptance-checklist-candidate.json")
	requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), "integration/frontend-backend-contract-map.json")
	requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), "gui_api_traceability_matrix_reviewed")
	requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), "gui_api_traceability_matrix_review_missing")
}

func TestTeamOfficeP0ToP11CommercializationBlockerRegisterTurnsMissingEvidenceIntoActionableIssues(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	registerPath := "tests/p0-p11-commercialization-blocker-register-candidate.json"
	register := readJSON(t, filepath.Join(base, registerPath))

	requireBool(t, register, "candidate_only", true)
	requireBool(t, register, "non_formal", true)
	if got := requireString(t, register, "register_status"); got != "candidate_blockers_defined_pending_cross_repo_execution" {
		t.Fatalf("register_status = %s, want candidate_blockers_defined_pending_cross_repo_execution", got)
	}
	if !strings.Contains(requireString(t, register, "objective_source"), "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source must reference the role pack studio test plan")
	}

	rawRows, ok := register["blocker_rows"].([]any)
	if !ok {
		t.Fatalf("blocker_rows missing")
	}
	if len(rawRows) != 12 {
		t.Fatalf("blocker_rows len = %d, want 12", len(rawRows))
	}

	expectedPhases := map[string]bool{}
	for i := 0; i < 12; i++ {
		expectedPhases[fmt.Sprintf("P%d", i)] = false
	}
	for i, raw := range rawRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("blocker row = %T", raw)
		}
		phaseID := requireString(t, row, "phase_id")
		if _, ok := expectedPhases[phaseID]; !ok {
			t.Fatalf("unexpected phase_id %s", phaseID)
		}
		expectedPhases[phaseID] = true
		if order, ok := row["phase_order"].(float64); !ok || int(order) != i+1 {
			t.Fatalf("%s phase_order = %v, want %d", phaseID, row["phase_order"], i+1)
		}
		if !strings.Contains(requireString(t, row, "issue_id"), phaseID) {
			t.Fatalf("%s issue_id must include phase id", phaseID)
		}
		requireStringIn(t, requireString(t, row, "risk_color"), "green", "yellow", "orange", "red")
		requireStringIn(t, requireString(t, row, "contract_impact"), "none", "implementation_only", "schema_or_dto_review", "cross_repo_boundary", "sovereignty_chain")
		requireStringIn(t, requireString(t, row, "lifecycle_status"), "设计中", "契约已定", "已接线")
		requireBool(t, row, "owner_authorization_required", true)
		requireBool(t, row, "acceptance_signoff_required", true)
		requireBool(t, row, "must_not_claim_complete", true)

		for _, key := range []string{"target_repos", "required_authoritative_evidence", "reproduction_or_execution_steps", "blocking_receipt_or_issue_requirements"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", phaseID, key)
			}
		}
		if requireString(t, row, "current_blocking_reason") == "" {
			t.Fatalf("%s current_blocking_reason must not be empty", phaseID)
		}

		targetRepos := strings.Join(asStringSlice(t, row["target_repos"]), "\n")
		evidence := strings.Join(asStringSlice(t, row["required_authoritative_evidence"]), "\n")
		switch phaseID {
		case "P2", "P3", "P4", "P5", "P6", "P7":
			for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos"} {
				if !strings.Contains(targetRepos, repo) {
					t.Fatalf("%s target_repos missing %s", phaseID, repo)
				}
			}
		case "P11":
			for _, repo := range []string{"truzhen-cloud", "truzhenos", "truzhen-client-web-desktop"} {
				if !strings.Contains(targetRepos, repo) {
					t.Fatalf("P11 target_repos missing %s", repo)
				}
			}
			for _, proof := range []string{"cloud_upload_receipt", "sandbox_payment_receipt", "download_receipt", "install_receipt"} {
				if !strings.Contains(evidence, proof) {
					t.Fatalf("P11 required_authoritative_evidence missing %s", proof)
				}
			}
		}
	}
	for phaseID, seen := range expectedPhases {
		if !seen {
			t.Fatalf("missing blocker row %s", phaseID)
		}
	}

	claimPolicy := requireObject(t, register, "completion_claim_policy")
	requireBool(t, claimPolicy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"all_blockers_resolved_or_owner_accepted",
		"all_p0_to_p11_authoritative_evidence_present",
		"independent_acceptance_signoff",
		"owner_go_no_go_decision",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p0_p11_commercialization_blocker_register"); got != registerPath {
		t.Fatalf("p0_p11_commercialization_blocker_register = %s, want %s", got, registerPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, registerPath) {
		t.Fatalf("candidate set artifact_files missing %s", registerPath)
	}

	phaseMatrix := readJSON(t, filepath.Join(base, "tests", "role-studio-phase-coverage-matrix-candidate.json"))
	if got := requireString(t, phaseMatrix, "blocker_register_ref"); got != registerPath {
		t.Fatalf("phase matrix blocker_register_ref = %s, want %s", got, registerPath)
	}
	signoff := readJSON(t, filepath.Join(base, "tests", "independent-acceptance-signoff-matrix-candidate.json"))
	if got := requireString(t, signoff, "blocker_register_ref"); got != registerPath {
		t.Fatalf("signoff blocker_register_ref = %s, want %s", got, registerPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p0_p11_blocker_register_verified") {
		t.Fatalf("product readiness matrix missing p0_p11_blocker_register_verified")
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, registerPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", registerPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	registerInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != registerPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, registerPath))
		if err != nil {
			t.Fatalf("read %s: %v", registerPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", registerPath, got, wantHash)
		}
		registerInManifest = true
	}
	if !registerInManifest {
		t.Fatalf("artifact manifest missing %s", registerPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == registerPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", registerPath)
}

func TestTeamOfficeProductStageFrontendBackendClosureReportMakesP10Auditable(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	reportPath := "tests/product-stage-frontend-backend-closure-report-candidate.json"
	report := readJSON(t, filepath.Join(base, reportPath))

	requireBool(t, report, "candidate_only", true)
	requireBool(t, report, "non_formal", true)
	if got := requireString(t, report, "closure_status"); got != "not_achieved_requires_cross_repo_execution" {
		t.Fatalf("closure_status = %s, want not_achieved_requires_cross_repo_execution", got)
	}
	if !strings.Contains(requireString(t, report, "objective_source"), "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source must reference the role pack studio test plan")
	}

	participants := strings.Join(asStringSlice(t, report["required_participants"]), "\n")
	for _, actor := range []string{"user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent"} {
		if !strings.Contains(participants, actor) {
			t.Fatalf("required_participants missing %s", actor)
		}
	}

	rawFrontendRows, ok := report["frontend_capability_rows"].([]any)
	if !ok {
		t.Fatalf("frontend_capability_rows missing")
	}
	expectedFrontend := map[string]bool{
		"role_studio_create_edit_preview":    false,
		"candidate_export":                   false,
		"capability_role_reference":          false,
		"team_settings_replacement":          false,
		"secretary_voice_vrm":                false,
		"cloud_purchase_download_install":    false,
		"production_go_live_controls":        false,
		"error_recovery_and_receipt_display": false,
	}
	if len(rawFrontendRows) != len(expectedFrontend) {
		t.Fatalf("frontend_capability_rows len = %d, want %d", len(rawFrontendRows), len(expectedFrontend))
	}
	for _, raw := range rawFrontendRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("frontend row = %T", raw)
		}
		rowID := requireString(t, row, "capability_id")
		if _, ok := expectedFrontend[rowID]; !ok {
			t.Fatalf("unexpected frontend capability_id %s", rowID)
		}
		expectedFrontend[rowID] = true
		if got := requireString(t, row, "truth_source"); got != "truzhen-client-web-desktop" {
			t.Fatalf("%s truth_source = %s", rowID, got)
		}
		if got := requireString(t, row, "status"); got != "pending_cross_repo_execution" {
			t.Fatalf("%s status = %s", rowID, got)
		}
		requireBool(t, row, "backend_pairing_required", true)
		for _, key := range []string{"current_repo_evidence", "gui_evidence_required", "backend_pairing_evidence_required", "missing_authoritative_evidence"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", rowID, key)
			}
		}
	}
	for rowID, seen := range expectedFrontend {
		if !seen {
			t.Fatalf("missing frontend capability row %s", rowID)
		}
	}

	rawBackendRows, ok := report["backend_capability_rows"].([]any)
	if !ok {
		t.Fatalf("backend_capability_rows missing")
	}
	expectedBackend := map[string]bool{
		"role_pack_candidate_lifecycle":         false,
		"capability_role_requirement":           false,
		"team_role_slot_binding":                false,
		"secretary_appearance_preference":       false,
		"commercial_upload_purchase_download":   false,
		"install_preflight_and_receipts":        false,
		"observability_and_audit":               false,
		"production_promotion_go_live_controls": false,
	}
	if len(rawBackendRows) != len(expectedBackend) {
		t.Fatalf("backend_capability_rows len = %d, want %d", len(rawBackendRows), len(expectedBackend))
	}
	for _, raw := range rawBackendRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("backend row = %T", raw)
		}
		rowID := requireString(t, row, "capability_id")
		if _, ok := expectedBackend[rowID]; !ok {
			t.Fatalf("unexpected backend capability_id %s", rowID)
		}
		expectedBackend[rowID] = true
		requireStringIn(t, requireString(t, row, "truth_source"), "truzhenos", "truzhen-cloud", "multi_repo")
		if got := requireString(t, row, "status"); got != "pending_cross_repo_execution" {
			t.Fatalf("%s status = %s", rowID, got)
		}
		for _, key := range []string{"current_repo_evidence", "receipt_evidence_required", "missing_authoritative_evidence"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", rowID, key)
			}
		}
	}
	for rowID, seen := range expectedBackend {
		if !seen {
			t.Fatalf("missing backend capability row %s", rowID)
		}
	}

	rawParityRows, ok := report["field_parity_matrix"].([]any)
	if !ok {
		t.Fatalf("field_parity_matrix missing")
	}
	expectedSurfaces := map[string]bool{
		"role_studio":                     false,
		"capability_role_reference":       false,
		"team_settings_role_tab":          false,
		"secretary_appearance":            false,
		"cloud_listing_purchase_download": false,
		"local_install_and_runtime_usage": false,
		"production_promotion_go_live":    false,
	}
	if len(rawParityRows) != len(expectedSurfaces) {
		t.Fatalf("field_parity_matrix len = %d, want %d", len(rawParityRows), len(expectedSurfaces))
	}
	for _, raw := range rawParityRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("parity row = %T", raw)
		}
		surfaceID := requireString(t, row, "surface_id")
		if _, ok := expectedSurfaces[surfaceID]; !ok {
			t.Fatalf("unexpected surface_id %s", surfaceID)
		}
		expectedSurfaces[surfaceID] = true
		if got := requireString(t, row, "parity_status"); got != "not_verified_requires_cross_repo_execution" {
			t.Fatalf("%s parity_status = %s", surfaceID, got)
		}
		for _, key := range []string{"frontend_fields", "backend_receipt_fields", "missing_authoritative_evidence"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", surfaceID, key)
			}
		}
	}
	for surfaceID, seen := range expectedSurfaces {
		if !seen {
			t.Fatalf("missing field parity row %s", surfaceID)
		}
	}

	rawGaps, ok := report["p0_p1_remaining_gaps"].([]any)
	if !ok {
		t.Fatalf("p0_p1_remaining_gaps missing")
	}
	gaps := strings.Join(asStringSlice(t, report["completion_barriers"]), "\n")
	for _, barrier := range []string{"frontend_smoke_missing", "backend_receipt_lookup_missing", "field_parity_report_missing", "cross_repo_authorization_missing"} {
		if !strings.Contains(gaps, barrier) {
			t.Fatalf("completion_barriers missing %s", barrier)
		}
	}
	seenGapPhases := map[string]bool{"P0": false, "P1": false}
	for _, raw := range rawGaps {
		gap, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("gap row = %T", raw)
		}
		phaseID := requireString(t, gap, "phase_id")
		if _, ok := seenGapPhases[phaseID]; !ok {
			t.Fatalf("unexpected p0/p1 gap phase %s", phaseID)
		}
		seenGapPhases[phaseID] = true
		if values := asStringSlice(t, gap["required_resolution_evidence"]); len(values) == 0 {
			t.Fatalf("%s required_resolution_evidence must not be empty", phaseID)
		}
	}
	for phaseID, seen := range seenGapPhases {
		if !seen {
			t.Fatalf("missing p0_p1_remaining_gaps row %s", phaseID)
		}
	}

	rawCards, ok := report["next_execution_cards"].([]any)
	if !ok {
		t.Fatalf("next_execution_cards missing")
	}
	expectedRepos := map[string]bool{
		"truzhen-client-web-desktop": false,
		"truzhenos":                  false,
		"truzhen-cloud":              false,
		"truzhen-contracts":          false,
	}
	for _, raw := range rawCards {
		card, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution card = %T", raw)
		}
		repo := requireString(t, card, "target_repo")
		if _, ok := expectedRepos[repo]; ok {
			expectedRepos[repo] = true
		}
		for _, key := range []string{"execution_card_ref", "required_action", "required_evidence"} {
			if requireString(t, card, key) == "" {
				t.Fatalf("%s missing for execution card %s", key, repo)
			}
		}
	}
	for repo, seen := range expectedRepos {
		if !seen {
			t.Fatalf("next_execution_cards missing repo %s", repo)
		}
	}

	claimPolicy := requireObject(t, report, "completion_claim_policy")
	requireBool(t, claimPolicy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"frontend_product_stage_smoke",
		"backend_receipt_lookup",
		"frontend_backend_field_match_report",
		"independent_acceptance_signoff",
		"owner_go_no_go_decision",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "product_stage_frontend_backend_closure_report"); got != reportPath {
		t.Fatalf("product_stage_frontend_backend_closure_report = %s, want %s", got, reportPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, reportPath) {
		t.Fatalf("candidate set artifact_files missing %s", reportPath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "product_stage_frontend_backend_closure_report_verified") {
		t.Fatalf("product readiness matrix missing product_stage_frontend_backend_closure_report_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "product_stage_frontend_backend_closure_report_ref"); got != reportPath {
		t.Fatalf("normal commercialization audit closure report ref = %s, want %s", got, reportPath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, reportPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", reportPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	reportInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != reportPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, reportPath))
		if err != nil {
			t.Fatalf("read %s: %v", reportPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", reportPath, got, wantHash)
		}
		reportInManifest = true
	}
	if !reportInManifest {
		t.Fatalf("artifact manifest missing %s", reportPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == reportPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", reportPath)
}

func TestTeamOfficeProductStageClosureReportIncludesProductionPromotionControls(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	report := readJSON(t, filepath.Join(base, "tests", "product-stage-frontend-backend-closure-report-candidate.json"))

	rawFrontendRows, ok := report["frontend_capability_rows"].([]any)
	if !ok {
		t.Fatalf("frontend_capability_rows missing")
	}
	foundFrontend := false
	for _, raw := range rawFrontendRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("frontend row = %T", raw)
		}
		if requireString(t, row, "capability_id") != "production_go_live_controls" {
			continue
		}
		foundFrontend = true
		if got := requireString(t, row, "truth_source"); got != "truzhen-client-web-desktop" {
			t.Fatalf("production_go_live_controls truth_source = %s", got)
		}
		guiEvidence := strings.Join(asStringSlice(t, row["gui_evidence_required"]), "\n")
		for _, evidence := range []string{
			"production_promotion_gate_review_screenshot",
			"production_go_live_request_screenshot",
			"real_payment_enable_block_or_confirm_screenshot",
			"production_signed_download_enable_screenshot",
			"production_listing_publish_screenshot",
			"production_install_observability_screenshot",
		} {
			if !strings.Contains(guiEvidence, evidence) {
				t.Fatalf("production_go_live_controls gui_evidence_required missing %s", evidence)
			}
		}
		backendPairing := strings.Join(asStringSlice(t, row["backend_pairing_evidence_required"]), "\n")
		for _, evidence := range []string{
			"production_go_live_request_receipt",
			"real_payment_enable_request_receipt",
			"production_signed_download_enable_receipt",
			"production_listing_publish_request_receipt",
			"production_install_observability_receipt",
		} {
			if !strings.Contains(backendPairing, evidence) {
				t.Fatalf("production_go_live_controls backend_pairing_evidence_required missing %s", evidence)
			}
		}
	}
	if !foundFrontend {
		t.Fatalf("missing frontend capability row production_go_live_controls")
	}

	rawBackendRows, ok := report["backend_capability_rows"].([]any)
	if !ok {
		t.Fatalf("backend_capability_rows missing")
	}
	foundBackend := false
	for _, raw := range rawBackendRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("backend row = %T", raw)
		}
		if requireString(t, row, "capability_id") != "production_promotion_go_live_controls" {
			continue
		}
		foundBackend = true
		if got := requireString(t, row, "truth_source"); got != "multi_repo" {
			t.Fatalf("production_promotion_go_live_controls truth_source = %s", got)
		}
		currentEvidence := strings.Join(asStringSlice(t, row["current_repo_evidence"]), "\n")
		for _, evidence := range []string{
			"commerce/commercial-production-promotion-gate-candidate.json",
			"integration/commercial-api-contract-candidate.json",
			"tests/p11-commercial-go-live-evidence-package-template.json",
		} {
			if !strings.Contains(currentEvidence, evidence) {
				t.Fatalf("production_promotion_go_live_controls current_repo_evidence missing %s", evidence)
			}
		}
		receiptEvidence := strings.Join(asStringSlice(t, row["receipt_evidence_required"]), "\n")
		for _, evidence := range []string{
			"production_go_live_request_receipt",
			"real_payment_enable_request_receipt",
			"production_signed_download_enable_receipt",
			"production_listing_publish_request_receipt",
			"production_install_observability_receipt",
		} {
			if !strings.Contains(receiptEvidence, evidence) {
				t.Fatalf("production_promotion_go_live_controls receipt_evidence_required missing %s", evidence)
			}
		}
	}
	if !foundBackend {
		t.Fatalf("missing backend capability row production_promotion_go_live_controls")
	}

	rawParityRows, ok := report["field_parity_matrix"].([]any)
	if !ok {
		t.Fatalf("field_parity_matrix missing")
	}
	foundParity := false
	for _, raw := range rawParityRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("parity row = %T", raw)
		}
		if requireString(t, row, "surface_id") != "production_promotion_go_live" {
			continue
		}
		foundParity = true
		frontendFields := strings.Join(asStringSlice(t, row["frontend_fields"]), "\n")
		for _, field := range []string{
			"production_promotion_gate_ref",
			"owner_go_no_go_decision_ref",
			"production_go_live_request",
			"real_payment_enable_request",
			"production_signed_download_enable",
			"production_listing_publish",
			"production_install_observability",
		} {
			if !strings.Contains(frontendFields, field) {
				t.Fatalf("production_promotion_go_live frontend_fields missing %s", field)
			}
		}
		backendFields := strings.Join(asStringSlice(t, row["backend_receipt_fields"]), "\n")
		for _, field := range []string{
			"production_go_live_request_receipt_ref",
			"real_payment_enable_request_receipt_ref",
			"production_signed_download_enable_receipt_ref",
			"production_listing_publish_request_receipt_ref",
			"production_install_observability_receipt_ref",
		} {
			if !strings.Contains(backendFields, field) {
				t.Fatalf("production_promotion_go_live backend_receipt_fields missing %s", field)
			}
		}
	}
	if !foundParity {
		t.Fatalf("missing field parity row production_promotion_go_live")
	}

	barriers := strings.Join(asStringSlice(t, report["completion_barriers"]), "\n")
	if !strings.Contains(barriers, "production_promotion_receipts_missing") {
		t.Fatalf("completion_barriers missing production_promotion_receipts_missing")
	}
	claimPolicy := requireObject(t, report, "completion_claim_policy")
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(requiredBeforeClaim, "production_promotion_receipts_verified") {
		t.Fatalf("completion_claim_policy.required_before_completion_claim missing production_promotion_receipts_verified")
	}
}

func TestTeamOfficeP11NormalCommercializationAcceptanceGateRequiresUploadDownloadInstallEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gatePath := "tests/p11-normal-commercialization-acceptance-gate-candidate.json"
	gate := readJSON(t, filepath.Join(base, gatePath))

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	if got := requireString(t, gate, "gate_status"); got != "not_run_requires_cross_repo_authoritative_receipts" {
		t.Fatalf("gate_status = %s, want not_run_requires_cross_repo_authoritative_receipts", got)
	}
	if !strings.Contains(requireString(t, gate, "objective_source"), "role-pack-studio-team-office-test-plan-20260704.md") {
		t.Fatalf("objective_source must reference role pack studio test plan")
	}

	testCases := strings.Join(asStringSlice(t, gate["plan_test_case_refs"]), "\n")
	for _, tc := range []string{
		"TC-CLOUD-01",
		"TC-CLOUD-02",
		"TC-PAY-01",
		"TC-DOWNLOAD-01",
		"TC-INSTALL-01",
		"TC-INSTALL-NEG-01",
		"TC-INSTALL-NEG-02",
	} {
		if !strings.Contains(testCases, tc) {
			t.Fatalf("plan_test_case_refs missing %s", tc)
		}
	}

	rawStages, ok := gate["ordered_decision_stages"].([]any)
	if !ok {
		t.Fatalf("ordered_decision_stages missing")
	}
	expectedStages := map[string]bool{
		"role_candidate_bundle_export":              false,
		"cloud_upload_listing_draft":                false,
		"marketplace_review_candidate":              false,
		"sandbox_order_payment_entitlement":         false,
		"entitled_signed_download":                  false,
		"local_install_enabled_version":             false,
		"post_install_team_binding_runtime":         false,
		"negative_cases_and_independent_acceptance": false,
	}
	if len(rawStages) != len(expectedStages) {
		t.Fatalf("ordered_decision_stages len = %d, want %d", len(rawStages), len(expectedStages))
	}
	for i, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := expectedStages[stageID]; !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		expectedStages[stageID] = true
		if order, ok := stage["stage_order"].(float64); !ok || int(order) != i+1 {
			t.Fatalf("%s stage_order = %v, want %d", stageID, stage["stage_order"], i+1)
		}
		requireStringIn(t, requireString(t, stage, "truth_source"), "truzhen-client-web-desktop", "truzhenos", "truzhen-cloud", "multi_repo", "independent_acceptance_agent")
		requireStringIn(t, requireString(t, stage, "stage_status"), "pending_cross_repo_execution", "blocked_until_prior_stage_evidence")
		if requireString(t, stage, "go_condition") == "" {
			t.Fatalf("%s go_condition missing", stageID)
		}
		if requireString(t, stage, "issue_route") == "" {
			t.Fatalf("%s issue_route missing", stageID)
		}
		for _, key := range []string{"gui_evidence_required", "receipt_evidence_required", "no_go_conditions"} {
			if values := asStringSlice(t, stage[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
		evidence := strings.Join(asStringSlice(t, stage["receipt_evidence_required"]), "\n")
		switch stageID {
		case "entitled_signed_download":
			for _, proof := range []string{"download_receipt", "artifact_hash_matches_upload", "entitlement_ref"} {
				if !strings.Contains(evidence, proof) {
					t.Fatalf("%s receipt_evidence_required missing %s", stageID, proof)
				}
			}
		case "local_install_enabled_version":
			for _, proof := range []string{"install_receipt", "artifact_hash_verification_receipt", "enabled_role_pack_version_ref", "forbidden_artifact_scan"} {
				if !strings.Contains(evidence, proof) {
					t.Fatalf("%s receipt_evidence_required missing %s", stageID, proof)
				}
			}
		}
	}
	for stageID, seen := range expectedStages {
		if !seen {
			t.Fatalf("missing ordered decision stage %s", stageID)
		}
	}

	intake := requireObject(t, gate, "evidence_intake_contract")
	requiredTopLevel := strings.Join(asStringSlice(t, intake["required_top_level_fields"]), "\n")
	for _, field := range []string{
		"gui_evidence_records",
		"cloud_receipt_records",
		"install_receipt_records",
		"negative_case_records",
		"independent_acceptance_ref",
		"owner_go_no_go_decision_ref",
	} {
		if !strings.Contains(requiredTopLevel, field) {
			t.Fatalf("evidence_intake_contract.required_top_level_fields missing %s", field)
		}
	}
	forbiddenIntake := strings.Join(asStringSlice(t, intake["forbidden_payloads"]), "\n")
	for _, forbidden := range []string{"raw_payment_token", "signed_download_url_secret", "cloud_access_token", "raw_voice_asset", "raw_vrm_asset"} {
		if !strings.Contains(forbiddenIntake, forbidden) {
			t.Fatalf("evidence_intake_contract.forbidden_payloads missing %s", forbidden)
		}
	}

	rawNegative, ok := gate["negative_case_decisions"].([]any)
	if !ok {
		t.Fatalf("negative_case_decisions missing")
	}
	expectedNegative := map[string]bool{
		"unpaid_download":                                 false,
		"refund_revoked_download":                         false,
		"expired_entitlement_install":                     false,
		"version_unpublished_or_revoked_download_install": false,
		"artifact_hash_mismatch_install":                  false,
		"real_payment_without_owner_authorization":        false,
		"production_publish_without_owner_authorization":  false,
	}
	if len(rawNegative) != len(expectedNegative) {
		t.Fatalf("negative_case_decisions len = %d, want %d", len(rawNegative), len(expectedNegative))
	}
	for _, raw := range rawNegative {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("negative case = %T", raw)
		}
		caseID := requireString(t, row, "case_id")
		if _, ok := expectedNegative[caseID]; !ok {
			t.Fatalf("unexpected negative case %s", caseID)
		}
		expectedNegative[caseID] = true
		if !strings.HasPrefix(requireString(t, row, "expected_status"), "blocked_") {
			t.Fatalf("%s expected_status must be blocked", caseID)
		}
		for _, key := range []string{"required_gui_evidence", "required_block_receipt", "user_visible_reason"} {
			if requireString(t, row, key) == "" {
				t.Fatalf("%s %s missing", caseID, key)
			}
		}
	}
	for caseID, seen := range expectedNegative {
		if !seen {
			t.Fatalf("missing negative case %s", caseID)
		}
	}

	algorithm := requireObject(t, gate, "pass_fail_algorithm")
	requireBool(t, algorithm, "normal_commercialization_pass_allowed_now", false)
	requiredAll := strings.Join(asStringSlice(t, algorithm["required_all"]), "\n")
	for _, proof := range []string{
		"cloud_upload_receipt",
		"sandbox_payment_receipt",
		"download_receipt",
		"install_receipt",
		"team_binding_receipt",
		"independent_acceptance_signoff",
	} {
		if !strings.Contains(requiredAll, proof) {
			t.Fatalf("pass_fail_algorithm.required_all missing %s", proof)
		}
	}
	decisions := strings.Join(asStringSlice(t, algorithm["allowed_output_decisions"]), "\n")
	for _, decision := range []string{"passed_verified", "blocked_missing_authoritative_evidence", "failed_requires_issue", "not_run"} {
		if !strings.Contains(decisions, decision) {
			t.Fatalf("allowed_output_decisions missing %s", decision)
		}
	}

	claimPolicy := requireObject(t, gate, "completion_claim_policy")
	requireBool(t, claimPolicy, "completion_claim_allowed", false)
	requiredBeforeClaim := strings.Join(asStringSlice(t, claimPolicy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{
		"cloud_upload_receipt",
		"sandbox_payment_receipt",
		"download_receipt",
		"install_receipt",
		"team_binding_receipt",
		"independent_acceptance_signoff",
		"owner_go_no_go_decision",
	} {
		if !strings.Contains(requiredBeforeClaim, proof) {
			t.Fatalf("completion_claim_policy.required_before_completion_claim missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_normal_commercialization_acceptance_gate"); got != gatePath {
		t.Fatalf("p11_normal_commercialization_acceptance_gate = %s, want %s", got, gatePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, gatePath) {
		t.Fatalf("candidate set artifact_files missing %s", gatePath)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p11_normal_commercialization_acceptance_gate_verified") {
		t.Fatalf("product readiness matrix missing p11_normal_commercialization_acceptance_gate_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "p11_normal_commercialization_acceptance_gate_ref"); got != gatePath {
		t.Fatalf("normal commercialization audit P11 acceptance gate ref = %s, want %s", got, gatePath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, gatePath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", gatePath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	gateInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != gatePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, gatePath))
		if err != nil {
			t.Fatalf("read %s: %v", gatePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", gatePath, got, wantHash)
		}
		gateInManifest = true
	}
	if !gateInManifest {
		t.Fatalf("artifact manifest missing %s", gatePath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == gatePath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", gatePath)
}

func TestTeamOfficeP11CommercializationRequiresReleaseRevocationNegativeEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	negativeCaseID := "version_unpublished_or_revoked_download_install"

	accessMatrix := readJSON(t, filepath.Join(base, "commerce", "download-install-access-matrix.json"))
	state := findObjectByString(t, asObjectSlice(t, accessMatrix["access_states"]), "state_id", "unpublished_or_revoked_version")
	if got := requireString(t, state, "download_status"); got != "blocked_release_not_available" {
		t.Fatalf("unpublished_or_revoked_version download_status = %s, want blocked_release_not_available", got)
	}
	if got := requireString(t, state, "install_status"); got != "blocked_release_revoked" {
		t.Fatalf("unpublished_or_revoked_version install_status = %s, want blocked_release_revoked", got)
	}

	checklist := readJSON(t, filepath.Join(base, "tests", "p11-evidence-acceptance-checklist-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, checklist["required_negative_cases"]), negativeCaseID)

	acceptanceGate := readJSON(t, filepath.Join(base, "tests", "p11-normal-commercialization-acceptance-gate-candidate.json"))
	negativeCase := findObjectByString(t, asObjectSlice(t, acceptanceGate["negative_case_decisions"]), "case_id", negativeCaseID)
	if got := requireString(t, negativeCase, "expected_status"); got != "blocked_release_revoked_or_not_available" {
		t.Fatalf("%s expected_status = %s, want blocked_release_revoked_or_not_available", negativeCaseID, got)
	}
	for _, key := range []string{"required_gui_evidence", "required_block_receipt", "user_visible_reason"} {
		if got := requireString(t, negativeCase, key); got == "" {
			t.Fatalf("%s %s missing", negativeCaseID, key)
		}
	}

	chain := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	negativeStage := findObjectByString(t, asObjectSlice(t, chain["ordered_verification_gates"]), "stage_id", "negative_cases_blocked")
	requireStringSliceContains(t, asStringSlice(t, negativeStage["required_authoritative_evidence"]), "blocked_release_revoked_or_not_available_receipt")
	blockingCases := requireObject(t, chain, "blocking_cases")
	releaseCase := requireObject(t, blockingCases, negativeCaseID)
	if got := requireString(t, releaseCase, "expected_status"); got != "blocked_release_revoked_or_not_available" {
		t.Fatalf("chain %s expected_status = %s, want blocked_release_revoked_or_not_available", negativeCaseID, got)
	}

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	stageGate := findObjectByString(t, asObjectSlice(t, goNoGo["required_stage_gates"]), "stage_id", "negative_cases_and_observability")
	requireStringSliceContains(t, asStringSlice(t, stageGate["blocked_by"]), "release_revoked_or_unpublished_block_receipt_missing")
	requireStringSliceContains(t, asStringSlice(t, stageGate["required_before_pass"]), "version unpublished or revoked download/install blocked")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	negativeGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "negative_cases")
	requireStringSliceContains(t, asStringSlice(t, negativeGate["required_evidence"]), "blocked_release_revoked_or_not_available_receipt")
}

func TestTeamOfficeP11CommercializationVerificationRecordTemplateCapturesRealRunDecision(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	templatePath := "tests/p11-normal-commercialization-verification-record-template.json"
	template := readJSON(t, filepath.Join(base, templatePath))

	requireBool(t, template, "candidate_only", true)
	requireBool(t, template, "non_formal", true)
	if got := requireString(t, template, "record_status"); got != "not_run_requires_authoritative_evidence" {
		t.Fatalf("record_status = %s, want not_run_requires_authoritative_evidence", got)
	}
	if got := requireString(t, template, "acceptance_gate_ref"); got != "tests/p11-normal-commercialization-acceptance-gate-candidate.json" {
		t.Fatalf("acceptance_gate_ref = %s", got)
	}
	if got := requireString(t, template, "e2e_run_record_ref"); got != "tests/e2e-evidence-run-record.json" {
		t.Fatalf("e2e_run_record_ref = %s", got)
	}

	recordSchema := requireObject(t, template, "record_schema")
	requiredTopLevel := strings.Join(asStringSlice(t, recordSchema["required_top_level_fields"]), "\n")
	for _, field := range []string{
		"stage_result_rows",
		"hash_continuity_result",
		"negative_case_result_rows",
		"independent_acceptance_result",
		"owner_go_no_go_decision",
		"final_decision",
	} {
		if !strings.Contains(requiredTopLevel, field) {
			t.Fatalf("record_schema.required_top_level_fields missing %s", field)
		}
	}
	forbiddenPayloads := strings.Join(asStringSlice(t, recordSchema["forbidden_payloads"]), "\n")
	for _, forbidden := range []string{"raw_payment_token", "signed_download_url_secret", "cloud_access_token", "raw_voice_asset", "raw_vrm_asset"} {
		if !strings.Contains(forbiddenPayloads, forbidden) {
			t.Fatalf("record_schema.forbidden_payloads missing %s", forbidden)
		}
	}

	rawStages, ok := template["stage_result_rows"].([]any)
	if !ok {
		t.Fatalf("stage_result_rows missing")
	}
	expectedStages := map[string]bool{
		"role_candidate_bundle_export":              false,
		"cloud_upload_listing_draft":                false,
		"marketplace_review_candidate":              false,
		"sandbox_order_payment_entitlement":         false,
		"entitled_signed_download":                  false,
		"local_install_enabled_version":             false,
		"post_install_team_binding_runtime":         false,
		"negative_cases_and_independent_acceptance": false,
	}
	if len(rawStages) != len(expectedStages) {
		t.Fatalf("stage_result_rows len = %d, want %d", len(rawStages), len(expectedStages))
	}
	for i, raw := range rawStages {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage result row = %T", raw)
		}
		stageID := requireString(t, row, "stage_id")
		if _, ok := expectedStages[stageID]; !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		expectedStages[stageID] = true
		if order, ok := row["stage_order"].(float64); !ok || int(order) != i+1 {
			t.Fatalf("%s stage_order = %v, want %d", stageID, row["stage_order"], i+1)
		}
		if got := requireString(t, row, "result_status"); got != "not_run" {
			t.Fatalf("%s result_status = %s, want not_run", stageID, got)
		}
		if got := requireString(t, row, "evidence_status"); got != "missing_authoritative_evidence" {
			t.Fatalf("%s evidence_status = %s, want missing_authoritative_evidence", stageID, got)
		}
		if requireString(t, row, "blocking_issue_id") == "" {
			t.Fatalf("%s blocking_issue_id missing", stageID)
		}
		for _, key := range []string{"required_gui_evidence_refs", "required_receipt_refs", "blocking_if_missing"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
	}
	for stageID, seen := range expectedStages {
		if !seen {
			t.Fatalf("missing stage result row %s", stageID)
		}
	}

	hashResult := requireObject(t, template, "hash_continuity_result")
	if got := requireString(t, hashResult, "result_status"); got != "not_run" {
		t.Fatalf("hash_continuity_result.result_status = %s", got)
	}
	hashRefs := strings.Join(asStringSlice(t, hashResult["required_refs"]), "\n")
	for _, ref := range []string{"bundle_tree_sha256", "cloud_upload_receipt", "download_receipt", "install_receipt"} {
		if !strings.Contains(hashRefs, ref) {
			t.Fatalf("hash_continuity_result.required_refs missing %s", ref)
		}
	}

	rawNegative, ok := template["negative_case_result_rows"].([]any)
	if !ok {
		t.Fatalf("negative_case_result_rows missing")
	}
	expectedNegative := map[string]bool{
		"unpaid_download":                                false,
		"refund_revoked_download":                        false,
		"expired_entitlement_install":                    false,
		"artifact_hash_mismatch_install":                 false,
		"real_payment_without_owner_authorization":       false,
		"production_publish_without_owner_authorization": false,
	}
	if len(rawNegative) != len(expectedNegative) {
		t.Fatalf("negative_case_result_rows len = %d, want %d", len(rawNegative), len(expectedNegative))
	}
	for _, raw := range rawNegative {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("negative result row = %T", raw)
		}
		caseID := requireString(t, row, "case_id")
		if _, ok := expectedNegative[caseID]; !ok {
			t.Fatalf("unexpected negative case %s", caseID)
		}
		expectedNegative[caseID] = true
		if got := requireString(t, row, "result_status"); got != "not_run" {
			t.Fatalf("%s result_status = %s, want not_run", caseID, got)
		}
		if !strings.HasPrefix(requireString(t, row, "expected_status"), "blocked_") {
			t.Fatalf("%s expected_status must be blocked", caseID)
		}
		for _, key := range []string{"required_gui_evidence_ref", "required_block_receipt_ref", "blocking_if_missing"} {
			if requireString(t, row, key) == "" {
				t.Fatalf("%s %s missing", caseID, key)
			}
		}
	}
	for caseID, seen := range expectedNegative {
		if !seen {
			t.Fatalf("missing negative case result row %s", caseID)
		}
	}

	acceptance := requireObject(t, template, "independent_acceptance_result")
	requireBool(t, acceptance, "required", true)
	if got := requireString(t, acceptance, "result_status"); got != "not_run" {
		t.Fatalf("independent_acceptance_result.result_status = %s", got)
	}
	ownerDecision := requireObject(t, template, "owner_go_no_go_decision")
	requireBool(t, ownerDecision, "required", true)
	if got := requireString(t, ownerDecision, "decision_status"); got != "pending_owner_decision" {
		t.Fatalf("owner_go_no_go_decision.decision_status = %s", got)
	}

	finalDecision := requireObject(t, template, "final_decision")
	if got := requireString(t, finalDecision, "current_decision"); got != "not_run" {
		t.Fatalf("final_decision.current_decision = %s", got)
	}
	requireBool(t, finalDecision, "completion_claim_allowed", false)
	allowedDecisions := strings.Join(asStringSlice(t, finalDecision["allowed_decisions"]), "\n")
	for _, decision := range []string{"passed_verified", "blocked_missing_authoritative_evidence", "failed_requires_issue", "not_run"} {
		if !strings.Contains(allowedDecisions, decision) {
			t.Fatalf("final_decision.allowed_decisions missing %s", decision)
		}
	}
	requiredBeforePass := strings.Join(asStringSlice(t, finalDecision["required_before_passed_verified"]), "\n")
	for _, proof := range []string{"all_stage_result_rows_passed", "hash_continuity_verified", "all_negative_cases_blocked_verified", "independent_acceptance_signed", "owner_go_no_go_decision"} {
		if !strings.Contains(requiredBeforePass, proof) {
			t.Fatalf("final_decision.required_before_passed_verified missing %s", proof)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_normal_commercialization_verification_record_template"); got != templatePath {
		t.Fatalf("p11_normal_commercialization_verification_record_template = %s, want %s", got, templatePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, templatePath) {
		t.Fatalf("candidate set artifact_files missing %s", templatePath)
	}

	p11Gate := readJSON(t, filepath.Join(base, "tests", "p11-normal-commercialization-acceptance-gate-candidate.json"))
	if got := requireString(t, p11Gate, "verification_record_template_ref"); got != templatePath {
		t.Fatalf("P11 gate verification_record_template_ref = %s, want %s", got, templatePath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p11_verification_record_template_verified") {
		t.Fatalf("product readiness matrix missing p11_verification_record_template_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "p11_normal_commercialization_verification_record_template_ref"); got != templatePath {
		t.Fatalf("normal commercialization audit verification template ref = %s, want %s", got, templatePath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, templatePath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", templatePath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	templateInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != templatePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, templatePath))
		if err != nil {
			t.Fatalf("read %s: %v", templatePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", templatePath, got, wantHash)
		}
		templateInManifest = true
	}
	if !templateInManifest {
		t.Fatalf("artifact manifest missing %s", templatePath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == templatePath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", templatePath)
}

func TestTeamOfficeP11EvidenceIngestionBinderMapsGuiRunEvidenceIntoVerificationRecord(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	binderPath := "tests/p11-evidence-ingestion-binder-candidate.json"
	binder := readJSON(t, filepath.Join(base, binderPath))

	requireBool(t, binder, "candidate_only", true)
	requireBool(t, binder, "non_formal", true)
	if got := requireString(t, binder, "binder_status"); got != "not_run_requires_real_evidence" {
		t.Fatalf("binder_status = %s, want not_run_requires_real_evidence", got)
	}
	for key, want := range map[string]string{
		"gui_execution_script_ref":           "tests/gui-user-agent-execution-script-candidate.json",
		"e2e_run_record_ref":                 "tests/e2e-evidence-run-record.json",
		"p11_verification_record_template":   "tests/p11-normal-commercialization-verification-record-template.json",
		"commercial_distribution_schema_ref": "commerce/commercial-distribution-receipt-schema-candidate.json",
		"commercial_evidence_gate_ref":       "tests/commercial-evidence-gate-candidate.json",
	} {
		if got := requireString(t, binder, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	policy := requireObject(t, binder, "ingestion_policy")
	for _, key := range []string{
		"user_agents_gui_only",
		"organizer_may_record_only",
		"no_manual_json_success",
		"no_backend_only_success",
		"redaction_required",
		"no_secret_payload_storage",
	} {
		requireBool(t, policy, key, true)
	}
	requiredSources := strings.Join(asStringSlice(t, policy["required_source_records"]), "\n")
	for _, source := range []string{"gui_step_evidence_records", "cloud_receipt_records", "install_receipt_records", "negative_case_records", "independent_acceptance_record", "owner_go_no_go_decision"} {
		if !strings.Contains(requiredSources, source) {
			t.Fatalf("ingestion_policy.required_source_records missing %s", source)
		}
	}

	rawMappings, ok := binder["gui_step_to_verification_map"].([]any)
	if !ok {
		t.Fatalf("gui_step_to_verification_map missing")
	}
	expectedSteps := map[string]string{
		"open_role_studio":                    "pre_p11_role_creation_prerequisite",
		"create_secretary_candidate":          "pre_p11_role_creation_prerequisite",
		"create_five_advisor_candidates":      "pre_p11_role_creation_prerequisite",
		"select_secretary_voice_vrm":          "pre_p11_role_creation_prerequisite",
		"create_capability_role_reference":    "pre_p11_capability_reference_prerequisite",
		"export_candidate_bundle":             "role_candidate_bundle_export",
		"upload_cloud_listing_draft":          "cloud_upload_listing_draft",
		"submit_marketplace_review_candidate": "marketplace_review_candidate",
		"sandbox_purchase":                    "sandbox_order_payment_entitlement",
		"download_purchased_artifact":         "entitled_signed_download",
		"install_downloaded_role_pack":        "local_install_enabled_version",
		"replace_team_roles_after_install":    "post_install_team_binding_runtime",
		"run_team_office_runtime_use":         "post_install_team_binding_runtime",
		"run_negative_cases":                  "negative_cases_and_independent_acceptance",
	}
	if len(rawMappings) != len(expectedSteps) {
		t.Fatalf("gui_step_to_verification_map len = %d, want %d", len(rawMappings), len(expectedSteps))
	}
	lastOrder := 0.0
	for _, raw := range rawMappings {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("mapping row = %T", raw)
		}
		stepID := requireString(t, row, "step_id")
		stageID := requireString(t, row, "verification_stage_id")
		wantStage, ok := expectedSteps[stepID]
		if !ok {
			t.Fatalf("unexpected step_id %s", stepID)
		}
		if stageID != wantStage {
			t.Fatalf("%s verification_stage_id = %s, want %s", stepID, stageID, wantStage)
		}
		order, ok := row["step_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s step_order = %v, want increasing number greater than %v", stepID, row["step_order"], lastOrder)
		}
		lastOrder = order
		requireStringIn(t, requireString(t, row, "result_binding_status"), "not_run", "requires_authoritative_evidence")
		slots := strings.Join(asStringSlice(t, row["required_evidence_slots"]), "\n")
		for _, slot := range []string{"gui_screenshot_path", "page_state_ref", "candidate_or_receipt_ref", "truth_source", "timestamp"} {
			if !strings.Contains(slots, slot) {
				t.Fatalf("%s required_evidence_slots missing %s", stepID, slot)
			}
		}
		switch stepID {
		case "upload_cloud_listing_draft", "submit_marketplace_review_candidate", "sandbox_purchase", "download_purchased_artifact":
			if !strings.Contains(slots, "cloud_receipt_ref") {
				t.Fatalf("%s required_evidence_slots missing cloud_receipt_ref", stepID)
			}
		case "install_downloaded_role_pack":
			for _, slot := range []string{"install_receipt_ref", "entitlement_verification_ref", "artifact_hash_verification_ref"} {
				if !strings.Contains(slots, slot) {
					t.Fatalf("%s required_evidence_slots missing %s", stepID, slot)
				}
			}
		case "replace_team_roles_after_install":
			for _, slot := range []string{"team_binding_receipt_ref", "owner_gate_decision_ref"} {
				if !strings.Contains(slots, slot) {
					t.Fatalf("%s required_evidence_slots missing %s", stepID, slot)
				}
			}
		case "run_negative_cases":
			if !strings.Contains(slots, "blocked_receipt_ref") {
				t.Fatalf("%s required_evidence_slots missing blocked_receipt_ref", stepID)
			}
		}
		if requireString(t, row, "missing_evidence_issue_id") == "" {
			t.Fatalf("%s missing_evidence_issue_id must not be empty", stepID)
		}
	}

	correlation := requireObject(t, binder, "receipt_correlation_rules")
	requireBool(t, correlation, "bundle_hash_must_match_upload_download_install", true)
	keys := strings.Join(asStringSlice(t, correlation["required_correlation_keys"]), "\n")
	for _, key := range []string{"candidate_set_ref", "artifact_ref", "bundle_tree_sha256", "correlation_id", "upload_receipt_ref", "download_receipt_ref", "install_receipt_ref", "team_binding_receipt_ref"} {
		if !strings.Contains(keys, key) {
			t.Fatalf("receipt_correlation_rules.required_correlation_keys missing %s", key)
		}
	}

	transitions := requireObject(t, binder, "status_transition_rules")
	if got := requireString(t, transitions, "initial_status"); got != "not_run" {
		t.Fatalf("initial_status = %s, want not_run", got)
	}
	allowedFinal := strings.Join(asStringSlice(t, transitions["allowed_final_statuses"]), "\n")
	for _, status := range []string{"passed_verified", "blocked_missing_authoritative_evidence", "failed_requires_issue"} {
		if !strings.Contains(allowedFinal, status) {
			t.Fatalf("allowed_final_statuses missing %s", status)
		}
	}
	requiredPass := strings.Join(asStringSlice(t, transitions["required_before_passed_verified"]), "\n")
	for _, proof := range []string{"all_gui_steps_have_evidence", "all_required_receipts_present", "bundle_hash_continuity_verified", "all_negative_cases_blocked_verified", "independent_acceptance_signed", "owner_go_no_go_decision_recorded"} {
		if !strings.Contains(requiredPass, proof) {
			t.Fatalf("required_before_passed_verified missing %s", proof)
		}
	}

	outputs := requireObject(t, binder, "output_record_bindings")
	for key, want := range map[string]string{
		"verification_record_output": "tests/p11-normal-commercialization-verification-record-template.json",
		"blocker_register_output":    "tests/p0-p11-commercialization-blocker-register-candidate.json",
		"chain_verifier_output":      "tests/commercial-chain-verifier-candidate.json",
	} {
		if got := requireString(t, outputs, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_evidence_ingestion_binder"); got != binderPath {
		t.Fatalf("p11_evidence_ingestion_binder = %s, want %s", got, binderPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, binderPath) {
		t.Fatalf("candidate set artifact_files missing %s", binderPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p11_evidence_ingestion_binder_verified") {
		t.Fatalf("product readiness matrix missing p11_evidence_ingestion_binder_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "p11_evidence_ingestion_binder_ref"); got != binderPath {
		t.Fatalf("normal commercialization audit binder ref = %s, want %s", got, binderPath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, binderPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", binderPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	binderInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != binderPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, binderPath))
		if err != nil {
			t.Fatalf("read %s: %v", binderPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", binderPath, got, wantHash)
		}
		binderInManifest = true
	}
	if !binderInManifest {
		t.Fatalf("artifact manifest missing %s", binderPath)
	}

	digest := readJSON(t, filepath.Join(base, "commerce", "artifact-bundle-digest-candidate.json"))
	payloadFiles, ok := digest["payload_file_hashes"].([]any)
	if !ok {
		t.Fatalf("payload_file_hashes missing")
	}
	for _, raw := range payloadFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("payload file item = %T", raw)
		}
		if requireString(t, item, "path") == binderPath {
			return
		}
	}
	t.Fatalf("artifact bundle digest payload missing %s", binderPath)
}

func TestTeamOfficeP11EvidenceIngestionBinderMapsPhaseDependenciesIntoVerificationRecord(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	binder := readJSON(t, filepath.Join(base, "tests", "p11-evidence-ingestion-binder-candidate.json"))
	template := readJSON(t, filepath.Join(base, "tests", "p11-normal-commercialization-verification-record-template.json"))
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))

	if got := requireString(t, binder, "source_execution_queue"); got != "integration/commercial-cross-repo-execution-queue-candidate.json" {
		t.Fatalf("source_execution_queue = %s", got)
	}
	if got := requireString(t, binder, "phase_dependency_verification_record_target"); got != "p11_phase_dependency_result" {
		t.Fatalf("phase_dependency_verification_record_target = %s", got)
	}

	policy := requireObject(t, binder, "ingestion_policy")
	requireStringSliceContains(t, asStringSlice(t, policy["required_source_records"]), "p11_phase_dependency_records")

	dependency := requireObject(t, binder, "phase_dependency_ingestion")
	for key, want := range map[string]string{
		"phase_dependency_links_ref": "integration/commercial-cross-repo-execution-queue-candidate.json#p11_phase_dependency_links",
		"missing_dependency_result":  "blocked_previous_phase_evidence_missing",
		"output_record_field":        "p11_phase_dependency_result",
	} {
		if got := requireString(t, dependency, key); got != want {
			t.Fatalf("phase_dependency_ingestion.%s = %s, want %s", key, got, want)
		}
	}
	for _, key := range []string{"all_phase_handoffs_required", "same_correlation_id_required", "same_bundle_tree_sha256_required"} {
		requireBool(t, dependency, key, true)
	}
	if got, ok := dependency["expected_phase_dependency_count"].(float64); !ok || int(got) != len(asObjectSlice(t, queue["p11_phase_dependency_links"])) {
		t.Fatalf("expected_phase_dependency_count = %v", dependency["expected_phase_dependency_count"])
	}
	for _, field := range []string{
		"run_request_phase_id",
		"depends_on_phase_ids",
		"required_previous_receipts",
		"previous_receipt_ref",
		"bundle_tree_sha256",
		"correlation_id",
		"dependency_status",
	} {
		requireStringSliceContains(t, asStringSlice(t, dependency["required_evidence_slots"]), field)
	}
	transitions := requireObject(t, binder, "status_transition_rules")
	requireStringSliceContains(t, asStringSlice(t, transitions["required_before_passed_verified"]), "p11_execution_queue_phase_dependencies_verified")

	recordSchema := requireObject(t, template, "record_schema")
	requireStringSliceContains(t, asStringSlice(t, recordSchema["required_top_level_fields"]), "p11_phase_dependency_result")
	result := requireObject(t, template, "p11_phase_dependency_result")
	if got := requireString(t, result, "result_status"); got != "not_run" {
		t.Fatalf("p11_phase_dependency_result.result_status = %s", got)
	}
	if got, ok := result["expected_phase_dependency_count"].(float64); !ok || int(got) != len(asObjectSlice(t, queue["p11_phase_dependency_links"])) {
		t.Fatalf("p11_phase_dependency_result.expected_phase_dependency_count = %v", result["expected_phase_dependency_count"])
	}
	for _, ref := range []string{"p11_phase_dependency_links", "previous_receipt_ref", "bundle_tree_sha256", "correlation_id"} {
		requireStringSliceContains(t, asStringSlice(t, result["required_refs"]), ref)
	}
	requireStringSliceContains(t, asStringSlice(t, result["blocking_if_missing"]), "blocked_previous_phase_evidence_missing")
	finalDecision := requireObject(t, template, "final_decision")
	requireStringSliceContains(t, asStringSlice(t, finalDecision["required_before_passed_verified"]), "p11_execution_queue_phase_dependencies_verified")
}

func TestTeamOfficeP11SandboxExecutionRunbookDefinesAuthorizedGuiCloudInstallRun(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	runbookPath := "tests/p11-sandbox-execution-runbook-candidate.json"
	runbook := readJSON(t, filepath.Join(base, runbookPath))

	requireBool(t, runbook, "candidate_only", true)
	requireBool(t, runbook, "non_formal", true)
	if got := requireString(t, runbook, "runbook_status"); got != "pending_owner_authorization" {
		t.Fatalf("runbook_status = %s, want pending_owner_authorization", got)
	}
	for key, want := range map[string]string{
		"source_plan_ref":               "docs/plans/role-pack-studio-team-office-test-plan-20260704.md",
		"gui_execution_script_ref":      "tests/gui-user-agent-execution-script-candidate.json",
		"evidence_ingestion_binder_ref": "tests/p11-evidence-ingestion-binder-candidate.json",
		"sandbox_readiness_ref":         "tests/sandbox-environment-readiness-candidate.json",
		"cross_repo_readiness_ref":      "integration/cross-repo-execution-readiness-package.json",
		"p11_acceptance_gate_ref":       "tests/p11-normal-commercialization-acceptance-gate-candidate.json",
		"p11_verification_record_ref":   "tests/p11-normal-commercialization-verification-record-template.json",
		"commercial_api_contract_ref":   "integration/commercial-api-contract-candidate.json",
		"commercial_receipt_schema_ref": "commerce/commercial-distribution-receipt-schema-candidate.json",
		"commercial_receipt_chain_ref":  "commerce/commercial-receipt-chain-candidate.json",
	} {
		if got := requireString(t, runbook, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	auth := requireObject(t, runbook, "authorization_policy")
	requireBool(t, auth, "explicit_owner_authorization_required_before_cross_repo_run", true)
	requireBool(t, auth, "sandbox_only", true)
	requireBool(t, auth, "real_payment_blocked", true)
	requireBool(t, auth, "production_publish_blocked", true)
	requireBool(t, auth, "no_push_or_merge_without_separate_authorization", true)
	if got := requireString(t, auth, "organizer_role"); got != "coordinate_record_and_route_issues_only" {
		t.Fatalf("organizer_role = %s", got)
	}

	rawRepos, ok := runbook["target_repositories"].([]any)
	if !ok {
		t.Fatalf("target_repositories missing")
	}
	expectedRepos := map[string]bool{
		"truzhen-client-web-desktop": false,
		"truzhenos":                  false,
		"truzhen-cloud":              false,
		"truzhen-contracts":          false,
	}
	for _, raw := range rawRepos {
		repo, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("target repository = %T", raw)
		}
		repoID := requireString(t, repo, "repo_id")
		if _, ok := expectedRepos[repoID]; !ok {
			t.Fatalf("unexpected repo_id %s", repoID)
		}
		expectedRepos[repoID] = true
		if requireString(t, repo, "absolute_path") == "" {
			t.Fatalf("%s absolute_path missing", repoID)
		}
		for _, key := range []string{"allowed_actions", "status_command", "evidence_outputs", "forbidden_actions"} {
			if key == "status_command" {
				requireString(t, repo, key)
				continue
			}
			if values := asStringSlice(t, repo[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", repoID, key)
			}
		}
		if repoID == "truzhen-contracts" {
			actions := strings.Join(asStringSlice(t, repo["allowed_actions"]), "\n")
			if !strings.Contains(actions, "read_only_schema_check") {
				t.Fatalf("truzhen-contracts allowed_actions missing read_only_schema_check")
			}
		}
	}
	for repoID, seen := range expectedRepos {
		if !seen {
			t.Fatalf("target_repositories missing %s", repoID)
		}
	}

	rawStages, ok := runbook["execution_sequence"].([]any)
	if !ok {
		t.Fatalf("execution_sequence missing")
	}
	expectedStages := map[string]bool{
		"preflight_authorization_and_status":        false,
		"role_candidate_bundle_export":              false,
		"cloud_upload_listing_draft":                false,
		"marketplace_review_candidate":              false,
		"sandbox_order_payment_entitlement":         false,
		"entitled_signed_download":                  false,
		"local_install_enabled_version":             false,
		"post_install_team_binding_runtime":         false,
		"negative_cases":                            false,
		"independent_acceptance_and_owner_go_no_go": false,
		"production_promotion_controls":             false,
	}
	if len(rawStages) != len(expectedStages) {
		t.Fatalf("execution_sequence len = %d, want %d", len(rawStages), len(expectedStages))
	}
	lastOrder := 0.0
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution stage = %T", raw)
		}
		stageID := requireString(t, stage, "stage_id")
		if _, ok := expectedStages[stageID]; !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		expectedStages[stageID] = true
		order, ok := stage["stage_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s stage_order = %v, want increasing number greater than %v", stageID, stage["stage_order"], lastOrder)
		}
		lastOrder = order
		requireStringIn(t, requireString(t, stage, "actor"), "publisher_user_view_gui_agent", "buyer_user_view_gui_agent", "organizer_coordinator_recorder", "independent_acceptance_agent")
		for _, key := range []string{"target_repos", "gui_step_refs", "required_authoritative_evidence", "pass_condition", "fail_fast_if_missing", "output_record_refs"} {
			if key == "pass_condition" || key == "fail_fast_if_missing" {
				requireString(t, stage, key)
				continue
			}
			if values := asStringSlice(t, stage[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
		repos := strings.Join(asStringSlice(t, stage["target_repos"]), "\n")
		evidence := strings.Join(asStringSlice(t, stage["required_authoritative_evidence"]), "\n")
		switch stageID {
		case "cloud_upload_listing_draft", "marketplace_review_candidate", "sandbox_order_payment_entitlement", "entitled_signed_download":
			if !strings.Contains(repos, "truzhen-cloud") {
				t.Fatalf("%s target_repos missing truzhen-cloud", stageID)
			}
			if !strings.Contains(evidence, "truzhen-cloud") && !strings.Contains(evidence, "cloud") {
				t.Fatalf("%s required evidence must include cloud receipt evidence", stageID)
			}
		case "local_install_enabled_version", "post_install_team_binding_runtime":
			if !strings.Contains(repos, "truzhenos") {
				t.Fatalf("%s target_repos missing truzhenos", stageID)
			}
			if !strings.Contains(evidence, "install_receipt") && !strings.Contains(evidence, "team_binding_receipt") {
				t.Fatalf("%s required evidence must include install or team binding receipt", stageID)
			}
		case "negative_cases":
			if !strings.Contains(evidence, "blocked_receipt") {
				t.Fatalf("negative_cases required evidence missing blocked_receipt")
			}
		}
	}
	for stageID, seen := range expectedStages {
		if !seen {
			t.Fatalf("execution_sequence missing %s", stageID)
		}
	}

	evidencePackage := requireObject(t, runbook, "evidence_output_package")
	requireBool(t, evidencePackage, "raw_evidence_not_stored_in_truzhen_packs", true)
	requiredOutputs := strings.Join(asStringSlice(t, evidencePackage["required_output_records"]), "\n")
	for _, output := range []string{"p11_verification_record", "e2e_run_record", "receipt_index", "screenshot_index", "hash_continuity_report", "negative_case_report", "independent_acceptance_signoff", "owner_go_no_go_decision"} {
		if !strings.Contains(requiredOutputs, output) {
			t.Fatalf("evidence_output_package.required_output_records missing %s", output)
		}
	}

	hash := requireObject(t, runbook, "hash_continuity_check")
	requireBool(t, hash, "required", true)
	for _, key := range []string{"bundle_tree_sha256_source", "upload_receipt_ref", "download_receipt_ref", "install_receipt_ref", "mismatch_result"} {
		requireString(t, hash, key)
	}

	rollback := requireObject(t, runbook, "rollback_and_cleanup_policy")
	requireBool(t, rollback, "sandbox_cleanup_only", true)
	requireBool(t, rollback, "do_not_delete_receipts", true)
	requireBool(t, rollback, "do_not_rewrite_team_history", true)

	forbidden := strings.Join(asStringSlice(t, runbook["forbidden"]), "\n")
	for _, item := range []string{"direct_api_call_as_user_action", "manual_json_edit_as_success", "real_payment_capture", "production_publish", "raw_voice_or_vrm_upload", "store_secret_or_signed_url"} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_sandbox_execution_runbook"); got != runbookPath {
		t.Fatalf("p11_sandbox_execution_runbook = %s, want %s", got, runbookPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, runbookPath) {
		t.Fatalf("candidate set artifact_files missing %s", runbookPath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p11_sandbox_execution_runbook_verified") {
		t.Fatalf("product readiness matrix missing p11_sandbox_execution_runbook_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "p11_sandbox_execution_runbook_ref"); got != runbookPath {
		t.Fatalf("normal commercialization audit runbook ref = %s, want %s", got, runbookPath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, runbookPath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", runbookPath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	runbookInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != runbookPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, runbookPath))
		if err != nil {
			t.Fatalf("read %s: %v", runbookPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", runbookPath, got, wantHash)
		}
		runbookInManifest = true
	}
	if !runbookInManifest {
		t.Fatalf("artifact manifest missing %s", runbookPath)
	}
}

func TestTeamOfficeP11CommercialGoLiveEvidencePackageTemplateCollectsFinalAcceptanceProof(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))

	requireBool(t, pkg, "candidate_only", true)
	requireBool(t, pkg, "non_formal", true)
	if got := requireString(t, pkg, "package_status"); got != "not_run_requires_cross_repo_evidence" {
		t.Fatalf("package_status = %s, want not_run_requires_cross_repo_evidence", got)
	}
	for key, want := range map[string]string{
		"runbook_ref":                   "tests/p11-sandbox-execution-runbook-candidate.json",
		"evidence_ingestion_binder_ref": "tests/p11-evidence-ingestion-binder-candidate.json",
		"verification_record_ref":       "tests/p11-normal-commercialization-verification-record-template.json",
		"e2e_run_record_ref":            "tests/e2e-evidence-run-record.json",
		"acceptance_gate_ref":           "tests/p11-normal-commercialization-acceptance-gate-candidate.json",
		"receipt_schema_ref":            "commerce/commercial-distribution-receipt-schema-candidate.json",
		"receipt_chain_ref":             "commerce/commercial-receipt-chain-candidate.json",
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requiredRecords := strings.Join(asStringSlice(t, schema["required_top_level_records"]), "\n")
	for _, record := range []string{"package_manifest", "gui_evidence_index", "secretary_appearance_asset_report", "receipt_index", "install_catalog_and_slot_mapping_report", "capability_role_reference_report", "hash_continuity_report", "p11_phase_dependency_report", "cloud_upload_listing_report", "marketplace_listing_review_compliance_report", "purchase_entitlement_report", "download_artifact_delivery_report", "role_studio_lineage_report", "download_install_access_matrix_report", "negative_case_report", "independent_acceptance_signoff", "owner_go_no_go_decision", "production_promotion_receipts", "final_go_live_decision"} {
		if !strings.Contains(requiredRecords, record) {
			t.Fatalf("package_schema.required_top_level_records missing %s", record)
		}
	}
	forbiddenPayloads := strings.Join(asStringSlice(t, schema["forbidden_payloads"]), "\n")
	for _, forbidden := range []string{"raw_payment_token", "signed_download_url_secret", "cloud_access_token", "raw_voice_asset", "raw_vrm_asset", "identity_document"} {
		if !strings.Contains(forbiddenPayloads, forbidden) {
			t.Fatalf("package_schema.forbidden_payloads missing %s", forbidden)
		}
	}
	requireBool(t, schema, "raw_evidence_must_remain_in_execution_repos", true)
	requireBool(t, schema, "redaction_required", true)

	truth := requireObject(t, pkg, "truth_source_map")
	for key, want := range map[string]string{
		"gui_evidence_truth":           "truzhen-client-web-desktop",
		"cloud_receipt_truth":          "truzhen-cloud",
		"install_receipt_truth":        "truzhenos",
		"team_binding_truth":           "truzhenos",
		"candidate_asset_truth":        "truzhen-packs",
		"independent_acceptance_truth": "acceptance_agent",
	} {
		if got := requireString(t, truth, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawSections, ok := pkg["package_sections"].([]any)
	if !ok {
		t.Fatalf("package_sections missing")
	}
	expectedSections := map[string]bool{
		"package_manifest":                             false,
		"gui_evidence_index":                           false,
		"gui_api_traceability_report":                  false,
		"secretary_appearance_asset_report":            false,
		"receipt_index":                                false,
		"install_catalog_and_slot_mapping_report":      false,
		"capability_role_reference_report":             false,
		"hash_continuity_report":                       false,
		"p11_phase_dependency_report":                  false,
		"cloud_upload_listing_report":                  false,
		"marketplace_listing_review_compliance_report": false,
		"purchase_entitlement_report":                  false,
		"download_artifact_delivery_report":            false,
		"role_studio_lineage_report":                   false,
		"download_install_access_matrix_report":        false,
		"negative_case_report":                         false,
		"independent_acceptance_signoff":               false,
		"owner_go_no_go_decision":                      false,
		"production_promotion_receipts":                false,
		"final_go_live_decision":                       false,
	}
	if len(rawSections) != len(expectedSections) {
		t.Fatalf("package_sections len = %d, want %d", len(rawSections), len(expectedSections))
	}
	lastOrder := 0.0
	for _, raw := range rawSections {
		section, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("package section = %T", raw)
		}
		sectionID := requireString(t, section, "section_id")
		if _, ok := expectedSections[sectionID]; !ok {
			t.Fatalf("unexpected section_id %s", sectionID)
		}
		expectedSections[sectionID] = true
		order, ok := section["section_order"].(float64)
		if !ok || order <= lastOrder {
			t.Fatalf("%s section_order = %v, want increasing number greater than %v", sectionID, section["section_order"], lastOrder)
		}
		lastOrder = order
		requireStringIn(t, requireString(t, section, "section_status"), "not_run", "pending_authoritative_evidence")
		for _, key := range []string{"required_source_refs", "required_fields", "blocking_if_missing"} {
			values := asStringSlice(t, section[key])
			if len(values) == 0 {
				t.Fatalf("%s %s must not be empty", sectionID, key)
			}
		}
		fields := strings.Join(asStringSlice(t, section["required_fields"]), "\n")
		switch sectionID {
		case "secretary_appearance_asset_report":
			for _, field := range []string{"secretary_appearance_gui_controls_verified", "voice_asset_ref", "vrm_asset_ref", "voice_vrm_selection_screenshot", "asset_ref_validation_receipt", "raw_asset_absence_scan_receipt"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("secretary_appearance_asset_report required_fields missing %s", field)
				}
			}
		case "receipt_index":
			for _, field := range []string{"cloud_upload_receipt", "sandbox_payment_receipt", "download_receipt", "install_receipt", "team_binding_receipt"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("receipt_index required_fields missing %s", field)
				}
			}
		case "install_catalog_and_slot_mapping_report":
			for _, field := range []string{"team_settings_catalog_refresh_receipt_ref", "six_replaceable_role_refs", "enabled_role_pack_version_refs", "slot_mapping_refs", "secretary_and_five_advisors_slot_mapping_screenshot"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("install_catalog_and_slot_mapping_report required_fields missing %s", field)
				}
			}
		case "hash_continuity_report":
			for _, field := range []string{"bundle_tree_sha256", "upload_hash", "download_hash", "install_hash", "mismatch_result"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("hash_continuity_report required_fields missing %s", field)
				}
			}
		case "capability_role_reference_report":
			for _, field := range []string{"capability_role_reference_after_install_verified", "capability_role_requirement_validation_receipt", "enabled_role_pack_version_ref", "role_body_not_copied_receipt", "capability_invocation_candidate_ref"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("capability_role_reference_report required_fields missing %s", field)
				}
			}
		case "p11_phase_dependency_report":
			for _, field := range []string{"phase_id", "required_previous_receipts", "previous_receipt_ref", "p11_execution_queue_phase_dependencies_verified"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("p11_phase_dependency_report required_fields missing %s", field)
				}
			}
		case "download_install_access_matrix_report":
			for _, field := range []string{"download_install_access_matrix_verified", "unpaid_download_blocked_receipt", "refund_revoked_download_blocked_receipt", "version_unpublished_or_revoked_download_install_blocked_receipt", "artifact_hash_mismatch_blocked_receipt"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("download_install_access_matrix_report required_fields missing %s", field)
				}
			}
		case "purchase_entitlement_report":
			for _, field := range []string{"purchase_entitlement_verified", "sandbox_order_receipt", "sandbox_payment_receipt", "real_payment_block_policy_receipt", "no_real_payment_capture_receipt", "license_ref", "entitlement_receipt", "entitlement_team_scope_verification_receipt_ref"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("purchase_entitlement_report required_fields missing %s", field)
				}
			}
		case "download_artifact_delivery_report":
			for _, field := range []string{"download_artifact_delivery_verified", "entitlement_ref", "entitlement_check_receipt", "download_receipt_ref", "signed_download_ref", "download_artifact_sha256", "bundle_tree_sha256", "download_artifact_hash_receipt"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("download_artifact_delivery_report required_fields missing %s", field)
				}
			}
		case "role_studio_lineage_report":
			for _, field := range []string{"role_studio_lineage_verified", "candidate_set_ref", "bundle_tree_sha256", "six_role_pack_refs", "role_pack_ref_set_hash", "export_receipt_ref", "cloud_upload_receipt_ref", "download_receipt_ref", "install_receipt_ref", "team_binding_receipt_ref", "runtime_session_ref", "same_six_role_pack_refs_from_export_to_runtime"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("role_studio_lineage_report required_fields missing %s", field)
				}
			}
		case "final_go_live_decision":
			for _, field := range []string{"current_decision", "completion_claim_allowed", "owner_decision_ref", "independent_acceptance_ref"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("final_go_live_decision required_fields missing %s", field)
				}
			}
		case "production_promotion_receipts":
			for _, field := range []string{"production_promotion_gate_pass_receipt", "production_go_live_request_receipt", "real_payment_enable_request_receipt", "production_signed_download_enable_receipt", "production_listing_publish_request_receipt", "production_install_observability_receipt"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("production_promotion_receipts required_fields missing %s", field)
				}
			}
		}
	}
	for sectionID, seen := range expectedSections {
		if !seen {
			t.Fatalf("package_sections missing %s", sectionID)
		}
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireBool(t, decisionGate, "completion_claim_allowed_now", false)
	requiredBefore := strings.Join(asStringSlice(t, decisionGate["required_before_go_live_pass"]), "\n")
	for _, proof := range []string{"all_package_sections_present", "all_receipts_authoritative", "secretary_appearance_gui_controls_verified", "secretary_appearance_asset_rights_verified", "hash_continuity_verified", "capability_role_reference_after_install_verified", "cloud_upload_listing_verified", "marketplace_listing_review_compliance_verified", "purchase_entitlement_verified", "download_artifact_delivery_verified", "role_studio_lineage_verified", "download_install_access_matrix_verified", "all_negative_cases_blocked", "independent_acceptance_signed", "owner_go_no_go_decision", "production_promotion_receipts_verified"} {
		if !strings.Contains(requiredBefore, proof) {
			t.Fatalf("decision_gate.required_before_go_live_pass missing %s", proof)
		}
	}
	allowed := strings.Join(asStringSlice(t, decisionGate["allowed_final_decisions"]), "\n")
	for _, decision := range []string{"passed_verified", "blocked_missing_authoritative_evidence", "failed_requires_issue", "not_run"} {
		if !strings.Contains(allowed, decision) {
			t.Fatalf("decision_gate.allowed_final_decisions missing %s", decision)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "p11_commercial_go_live_evidence_package_template"); got != packagePath {
		t.Fatalf("p11_commercial_go_live_evidence_package_template = %s, want %s", got, packagePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, packagePath) {
		t.Fatalf("candidate set artifact_files missing %s", packagePath)
	}
	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productClaimPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	productProofs := strings.Join(asStringSlice(t, productClaimPolicy["required_before_completion_claim"]), "\n")
	if !strings.Contains(productProofs, "p11_go_live_evidence_package_template_verified") {
		t.Fatalf("product readiness matrix missing p11_go_live_evidence_package_template_verified")
	}
	audit := readJSON(t, filepath.Join(base, "tests", "normal-commercialization-completion-audit-candidate.json"))
	if got := requireString(t, audit, "p11_commercial_go_live_evidence_package_template_ref"); got != packagePath {
		t.Fatalf("normal commercialization audit package ref = %s, want %s", got, packagePath)
	}
	chainVerifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	related := strings.Join(asStringSlice(t, chainVerifier["related_artifacts"]), "\n")
	if !strings.Contains(related, packagePath) {
		t.Fatalf("commercial chain verifier related_artifacts missing %s", packagePath)
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawManifestFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	packageInManifest := false
	for _, raw := range rawManifestFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != packagePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, packagePath))
		if err != nil {
			t.Fatalf("read %s: %v", packagePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", packagePath, got, wantHash)
		}
		packageInManifest = true
	}
	if !packageInManifest {
		t.Fatalf("artifact manifest missing %s", packagePath)
	}
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresCompleteUserViewGuiStepCoverage(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	scriptPath := "tests/gui-user-agent-execution-script-candidate.json"

	pkg := readJSON(t, filepath.Join(base, packagePath))
	script := readJSON(t, filepath.Join(base, scriptPath))
	scriptSteps := asObjectSlice(t, script["ordered_steps"])
	wantStepIDs := make([]string, 0, len(scriptSteps))
	for _, step := range scriptSteps {
		wantStepIDs = append(wantStepIDs, requireString(t, step, "step_id"))
	}

	guiSection := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "gui_evidence_index")
	requireInt(t, guiSection, "expected_gui_step_count", len(wantStepIDs))
	gotStepIDs := asStringSlice(t, guiSection["required_gui_step_ids"])
	if len(gotStepIDs) != len(wantStepIDs) {
		t.Fatalf("gui_evidence_index required_gui_step_ids len = %d, want %d", len(gotStepIDs), len(wantStepIDs))
	}
	for i, want := range wantStepIDs {
		if got := gotStepIDs[i]; got != want {
			t.Fatalf("required_gui_step_ids[%d] = %s, want %s", i, got, want)
		}
	}
	for _, field := range []string{
		"expected_gui_step_count",
		"required_gui_step_ids",
		"observed_gui_step_ids",
		"missing_gui_step_ids",
		"user_view_gui_agent_operation_log_ref",
		"gui_step_coverage_result",
	} {
		requireStringSliceContains(t, asStringSlice(t, guiSection["required_fields"]), field)
	}
	for _, blocker := range []string{"all_user_view_gui_steps_covered_missing", "gui_step_coverage_result_missing"} {
		requireStringSliceContains(t, asStringSlice(t, guiSection["blocking_if_missing"]), blocker)
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "all_user_view_gui_steps_covered")

	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	pass := requireObject(t, verifier, "normal_commercialization_pass_conditions")
	if got := requireString(t, pass, "all_user_view_gui_steps_covered"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.all_user_view_gui_steps_covered = %s, want required", got)
	}

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "all_user_view_gui_steps_covered")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "all_user_view_gui_steps_covered")
}

func TestTeamOfficeP11FinalEvidenceChainRequiresGuiAPITraceability(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	contractMapPath := "integration/frontend-backend-contract-map.json"
	productReportPath := "tests/product-stage-frontend-backend-closure-report-candidate.json"
	guiScriptPath := "tests/gui-user-agent-execution-script-candidate.json"
	commercialAPIPath := "integration/commercial-api-contract-candidate.json"

	pkg := readJSON(t, filepath.Join(base, packagePath))
	for key, want := range map[string]string{
		"source_frontend_backend_contract_map":                 contractMapPath,
		"source_product_stage_frontend_backend_closure_report": productReportPath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("package %s = %s, want %s", key, got, want)
		}
	}
	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "gui_api_traceability_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "gui_api_traceability_report")
	for _, sourceRef := range []string{contractMapPath, productReportPath, guiScriptPath, commercialAPIPath} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"gui_api_traceability_matrix_verified",
		"script_step_id",
		"commercial_api_endpoint_refs",
		"required_receipt_or_result_fields",
		"required_block_statuses",
		"required_gui_evidence_slots",
		"candidate_or_receipt_ref",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"gui_api_traceability_matrix_missing",
		"commercial_api_endpoint_trace_missing",
		"receipt_field_trace_missing",
		"gui_evidence_slot_trace_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "gui_api_traceability_matrix_verified")

	checklist := readJSON(t, filepath.Join(base, checklistPath))
	for key, want := range map[string]string{
		"source_frontend_backend_contract_map":                 contractMapPath,
		"source_product_stage_frontend_backend_closure_report": productReportPath,
	} {
		if got := requireString(t, checklist, key); got != want {
			t.Fatalf("checklist %s = %s, want %s", key, got, want)
		}
	}
	traceabilityCheck := findObjectByString(t, asObjectSlice(t, checklist["stage_acceptance_checks"]), "stage_id", "gui_api_traceability_matrix")
	requireStringSliceContains(t, asStringSlice(t, traceabilityCheck["required_gui_evidence"]), "gui_step_to_api_traceability_review_screen")
	requireStringSliceContains(t, asStringSlice(t, traceabilityCheck["required_receipt_evidence"]), "commercial_api_endpoint_receipt_lookup_result")
	requireStringSliceContains(t, asStringSlice(t, traceabilityCheck["required_hash_or_correlation_checks"]), "candidate_or_receipt_ref")
	requireStringSliceContains(t, asStringSlice(t, traceabilityCheck["blocking_if_missing"]), "gui_api_traceability_matrix_missing")
	requireStringSliceContains(t, asStringSlice(t, traceabilityCheck["writeback_targets"]), packagePath+"#gui_api_traceability_report")
	checklistPolicy := requireObject(t, checklist, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, checklistPolicy["required_before_p11_pass"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, checklistPolicy["non_sufficient_evidence"]), "gui_api_traceability_matrix_missing")

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	for key, want := range map[string]string{
		"source_frontend_backend_contract_map":                 contractMapPath,
		"source_product_stage_frontend_backend_closure_report": productReportPath,
	} {
		if got := requireString(t, readiness, key); got != want {
			t.Fatalf("readiness %s = %s, want %s", key, got, want)
		}
	}
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "gui_api_traceability_matrix_verified")
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, ref := range []string{contractMapPath, productReportPath, "commercial_api_endpoint_refs", "required_receipt_or_result_fields"} {
		if !strings.Contains(readinessEvidence, ref) {
			t.Fatalf("readiness gui_api_traceability evidence_required missing %s: %s", ref, readinessEvidence)
		}
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_gui_api_traceability_matrix")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	for key, want := range map[string]string{
		"source_frontend_backend_contract_map":                 contractMapPath,
		"source_product_stage_frontend_backend_closure_report": productReportPath,
	} {
		if got := requireString(t, goNoGo, key); got != want {
			t.Fatalf("go/no-go %s = %s, want %s", key, got, want)
		}
	}
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "gui_api_traceability_matrix_verified")
	requireBool(t, goNoGoGate, "required_final_value", true)
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	goNoGoEvidence := requireString(t, goNoGoGate, "evidence_required")
	for _, ref := range []string{contractMapPath, productReportPath, packagePath, checklistPath} {
		if !strings.Contains(goNoGoEvidence, ref) {
			t.Fatalf("go/no-go gui_api_traceability evidence_required missing %s: %s", ref, goNoGoEvidence)
		}
	}
	goNoGoRule := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, goNoGoRule["required_before_final_decision"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_gui_api_traceability_matrix")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["required_before_completion_claim"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["non_sufficient_evidence"]), "gui_api_traceability_matrix_missing")
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), "frontend_backend_contract_map_gui_api_traceability_matrix")
	requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), "commercial_api_receipt_field_trace")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	goalPolicy := requireObject(t, evidenceMap, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, goalPolicy["required_before_goal_complete"]), "gui_api_traceability_matrix_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "gui_api_traceability_matrix_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_gui_api_traceability_matrix")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresMarketplaceListingReviewCompliance(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	listingPath := "commerce/cloud-listing-candidate.json"
	reviewPath := "commerce/marketplace-review-submission-candidate.json"
	supportRefundPath := "commerce/support-refund-revocation-policy-candidate.json"
	termsPrivacyPath := "commerce/commercial-terms-privacy-policy-candidate.json"
	publisherPath := "commerce/publisher-account-settlement-policy-candidate.json"

	pkg := readJSON(t, filepath.Join(base, packagePath))
	for key, want := range map[string]string{
		"source_cloud_listing_candidate":             listingPath,
		"source_marketplace_review_submission":       reviewPath,
		"source_support_refund_revocation_policy":    supportRefundPath,
		"source_commercial_terms_privacy_policy":     termsPrivacyPath,
		"source_publisher_account_settlement_policy": publisherPath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("package %s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "marketplace_listing_review_compliance_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "marketplace_listing_review_compliance_report")
	for _, sourceRef := range []string{listingPath, reviewPath, supportRefundPath, termsPrivacyPath, publisherPath} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"marketplace_listing_review_compliance_verified",
		"listing_draft_ref",
		"six_role_component_listing_verified",
		"marketplace_review_candidate_receipt",
		"terms_privacy_data_policy_accepted",
		"pre_purchase_support_refund_revocation_disclosures_verified",
		"publisher_identity_settlement_verified",
		"pricing_owner_approval_receipt",
		"candidate_role_disclaimer_visible",
		"production_publish_block_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"marketplace_listing_review_compliance_missing",
		"listing_draft_receipt_missing",
		"marketplace_review_candidate_receipt_missing",
		"pre_purchase_disclosure_receipt_missing",
		"terms_privacy_acceptance_receipt_missing",
		"publisher_identity_settlement_receipt_missing",
		"production_publish_block_receipt_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "marketplace_listing_review_compliance_verified")

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "marketplace_listing_review_compliance_verified")
	evidence := requireString(t, readinessGate, "evidence_required")
	for _, ref := range []string{listingPath, reviewPath, supportRefundPath, termsPrivacyPath, publisherPath, packagePath} {
		if !strings.Contains(evidence, ref) {
			t.Fatalf("readiness marketplace compliance evidence missing %s: %s", ref, evidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_marketplace_listing_review_compliance")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "marketplace_listing_review_compliance_verified")
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_marketplace_listing_review_compliance")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), "marketplace_review_candidate_receipt")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "marketplace_listing_review_compliance_missing")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "marketplace_listing_review_compliance_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "marketplace_listing_review_compliance_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_marketplace_listing_review_compliance")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresPurchaseEntitlementReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	orderPaymentPath := "commerce/order-payment-state-machine-candidate.json"
	entitlementPolicyPath := "commerce/license-entitlement-policy-candidate.json"
	buyerLibraryPath := "commerce/buyer-library-install-state-candidate.json"

	pkg := readJSON(t, filepath.Join(base, packagePath))
	for key, want := range map[string]string{
		"source_order_payment_state_machine": orderPaymentPath,
		"source_license_entitlement_policy":  entitlementPolicyPath,
		"source_buyer_library_install_state": buyerLibraryPath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("package %s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "purchase_entitlement_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "purchase_entitlement_report")
	for _, sourceRef := range []string{orderPaymentPath, entitlementPolicyPath, buyerLibraryPath, "commerce/commercial-distribution-receipt-schema-candidate.json", "commerce/commercial-receipt-chain-candidate.json"} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"purchase_entitlement_verified",
		"sandbox_order_receipt",
		"sandbox_payment_receipt",
		"real_payment_block_policy_receipt",
		"no_real_payment_capture_receipt",
		"license_ref",
		"entitlement_receipt",
		"entitlement_team_scope_verification_receipt_ref",
		"buyer_library_purchased_available_screenshot",
		"owned_library_entitlement_screenshot",
		"payment_failed_block_receipt",
		"entitlement_without_paid_receipt_blocked",
		"refund_or_chargeback_revocation_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"purchase_entitlement_report_missing",
		"sandbox_order_receipt_missing",
		"sandbox_payment_receipt_missing",
		"entitlement_receipt_missing",
		"no_real_payment_capture_receipt_missing",
		"entitlement_team_scope_verification_missing",
		"payment_failure_block_receipt_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["non_sufficient_evidence"]), "purchase_entitlement_report_missing")

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "purchase_entitlement_verified")
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, ref := range []string{orderPaymentPath, entitlementPolicyPath, buyerLibraryPath, packagePath} {
		if !strings.Contains(readinessEvidence, ref) {
			t.Fatalf("readiness purchase entitlement evidence missing %s: %s", ref, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_purchase_entitlement")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "purchase_entitlement_verified")
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_purchase_entitlement")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), "sandbox_payment_receipt")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "purchase_entitlement_missing")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "purchase_entitlement_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "purchase_entitlement_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_purchase_entitlement")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresRoleStudioLineageReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	receiptChainPath := "commerce/commercial-receipt-chain-candidate.json"
	chainVerifierPath := "tests/commercial-chain-verifier-candidate.json"
	runtimeUsagePath := "usage/team-office-runtime-usage-candidate.json"
	installMapPath := "install/install-runtime-activation-map-candidate.json"
	bindingCatalogPath := "bindings/team-settings-installed-role-catalog-candidate.json"

	pkg := readJSON(t, filepath.Join(base, packagePath))
	for key, want := range map[string]string{
		"receipt_chain_ref":                           receiptChainPath,
		"source_runtime_usage_candidate":              runtimeUsagePath,
		"source_install_runtime_activation_map":       installMapPath,
		"source_team_settings_installed_role_catalog": bindingCatalogPath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("package %s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "role_studio_lineage_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "role_studio_lineage_report")
	for _, sourceRef := range []string{receiptChainPath, chainVerifierPath, runtimeUsagePath, installMapPath, bindingCatalogPath} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"role_studio_lineage_verified",
		"candidate_set_ref",
		"bundle_tree_sha256",
		"six_role_pack_refs",
		"role_pack_ref_set_hash",
		"export_receipt_ref",
		"cloud_upload_receipt_ref",
		"download_receipt_ref",
		"install_receipt_ref",
		"team_binding_receipt_ref",
		"runtime_session_ref",
		"same_six_role_pack_refs_from_export_to_runtime",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"role_studio_lineage_report_missing",
		"role_pack_ref_set_hash_mismatch",
		"six_role_pack_refs_missing",
		"runtime_role_set_differs_from_installed_lineage",
		"team_binding_receipt_ref_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "role_studio_lineage_verified")

	chainVerifier := readJSON(t, filepath.Join(base, chainVerifierPath))
	passConditions := requireObject(t, chainVerifier, "normal_commercialization_pass_conditions")
	if got := requireString(t, passConditions, "role_studio_lineage_verified"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.role_studio_lineage_verified = %s, want required", got)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "role_studio_lineage_verified")
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, ref := range []string{packagePath, receiptChainPath, chainVerifierPath, runtimeUsagePath, installMapPath, bindingCatalogPath} {
		if !strings.Contains(readinessEvidence, ref) {
			t.Fatalf("readiness role studio lineage evidence missing %s: %s", ref, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_role_studio_lineage")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "role_studio_lineage_verified")
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_role_studio_lineage")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), "role_pack_ref_set_hash")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "role_studio_lineage_missing")

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "role_studio_lineage_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "role_studio_lineage_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_role_studio_lineage")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresDownloadInstallAccessMatrixReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	matrixPath := "commerce/download-install-access-matrix.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))

	if got := requireString(t, pkg, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "download_install_access_matrix_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "download_install_access_matrix_report")
	for _, sourceRef := range []string{
		matrixPath,
		"tests/commercial-readiness-verifier-candidate.json",
		"tests/commercial-go-no-go-gate-candidate.json",
		"commerce/commercial-go-live-approval-candidate.json",
		"commerce/commercial-production-promotion-gate-candidate.json",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"download_install_access_matrix_verified",
		"unpaid_download_blocked_receipt",
		"refund_revoked_download_blocked_receipt",
		"version_unpublished_or_revoked_download_install_blocked_receipt",
		"artifact_hash_mismatch_blocked_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"download_install_access_matrix_verified_missing",
		"unpaid_download_blocked_receipt_missing",
		"artifact_hash_mismatch_blocked_receipt_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "download_install_access_matrix_verified")
}

func TestTeamOfficeP11GoLiveEvidencePackageCollectsExecutionQueuePhaseDependencyProof(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	requestPath := "tests/p11-sandbox-run-request-candidate.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))
	queue := readJSON(t, filepath.Join(base, queuePath))

	if got := requireString(t, pkg, "source_execution_queue"); got != queuePath {
		t.Fatalf("source_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, pkg, "source_p11_phase_dependency_contract"); got != requestPath+"#phase_dependency_contract" {
		t.Fatalf("source_p11_phase_dependency_contract = %s, want %s", got, requestPath+"#phase_dependency_contract")
	}
	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "p11_phase_dependency_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "p11_phase_dependency_report")
	for _, sourceRef := range []string{queuePath, requestPath, "tests/commercial-chain-verifier-candidate.json"} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"phase_id",
		"depends_on_phase_ids",
		"required_previous_receipts",
		"previous_receipt_ref",
		"bundle_tree_sha256",
		"correlation_id",
		"dependency_status",
		"p11_execution_queue_phase_dependencies_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"blocked_previous_phase_evidence_missing",
		"previous_receipt_ref_missing",
		"p11_execution_queue_phase_dependencies_verified_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}
	if got, ok := section["expected_phase_dependency_count"].(float64); !ok || int(got) != len(asObjectSlice(t, queue["p11_phase_dependency_links"])) {
		t.Fatalf("expected_phase_dependency_count = %v", section["expected_phase_dependency_count"])
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "p11_execution_queue_phase_dependencies_verified")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresSixRoleInstallCatalogEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	activationPath := "install/install-runtime-activation-map-candidate.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))

	if got := requireString(t, pkg, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}
	if got := requireString(t, pkg, "evidence_acceptance_checklist_ref"); got != checklistPath {
		t.Fatalf("evidence_acceptance_checklist_ref = %s, want %s", got, checklistPath)
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "install_catalog_and_slot_mapping_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "install_catalog_and_slot_mapping_report")
	requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), activationPath)
	requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), checklistPath)
	for _, field := range []string{
		"team_settings_catalog_refresh_receipt_ref",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
		"team_settings_installed_role_catalog_screenshot",
		"six_replaceable_roles_visible_screenshot",
		"secretary_and_five_advisors_slot_mapping_screenshot",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"six_role_catalog_incomplete",
		"enabled_role_pack_version_refs_missing",
		"slot_mapping_refs_missing",
		"team_settings_catalog_refresh_receipt_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "six_installed_role_catalog_evidence_verified")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	productPolicy := requireObject(t, productMatrix, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, productPolicy["required_before_completion_claim"]), "six_installed_role_catalog_evidence_verified")
	localInstallGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "local_install_and_team_binding")
	for _, proof := range []string{
		"team_settings_catalog_refresh_receipt_ref",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
		"secretary_and_five_advisors_slot_mapping_screenshot",
	} {
		requireStringSliceContains(t, asStringSlice(t, localInstallGate["required_evidence"]), proof)
	}

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("commercial readiness source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "six_installed_role_catalog_evidence_verified")
	frontendBackend := requireObject(t, requireObject(t, readiness, "current_blockers"), "frontend_backend_acceptance")
	requireStringSliceContains(t, asStringSlice(t, frontendBackend["required_before_pass"]), "six_installed_role_catalog_evidence_verified")
}

func TestTeamOfficeCommercialReadinessGoNoGoProductAndGoalRequireSixInstalledRoleCatalogVerified(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	checklistPath := "tests/p11-evidence-acceptance-checklist-candidate.json"
	activationPath := "install/install-runtime-activation-map-candidate.json"
	bindingCatalogPath := "bindings/team-settings-installed-role-catalog-candidate.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	for key, want := range map[string]string{
		"source_install_runtime_activation_map":       activationPath,
		"source_team_settings_installed_role_catalog": bindingCatalogPath,
		"source_p11_evidence_acceptance_checklist":    checklistPath,
		"source_p11_go_live_evidence_package":         packagePath,
	} {
		if got := requireString(t, readiness, key); got != want {
			t.Fatalf("readiness %s = %s, want %s", key, got, want)
		}
	}
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "six_installed_role_catalog_evidence_verified")
	if got := requireString(t, readinessGate, "current_result"); got != "pending" {
		t.Fatalf("readiness six_installed_role_catalog_evidence_verified current_result = %s, want pending", got)
	}
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	readinessEvidence := requireString(t, readinessGate, "evidence_required")
	for _, want := range []string{
		packagePath + "#install_catalog_and_slot_mapping_report",
		checklistPath,
		activationPath,
		bindingCatalogPath,
		"team_settings_catalog_refresh_receipt_ref",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		if !strings.Contains(readinessEvidence, want) {
			t.Fatalf("readiness six installed role catalog evidence missing %s: %s", want, readinessEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, readiness, "completion_claim_policy")["required_before_completion_claim"]), "six_installed_role_catalog_evidence_verified")
	requireStringSliceContains(t, asStringSlice(t, readiness["non_sufficient_evidence"]), "commercial_ready_without_six_installed_role_catalog")

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	goNoGoGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "six_installed_role_catalog_evidence_verified")
	requireBool(t, goNoGoGate, "required_final_value", true)
	requireBool(t, goNoGoGate, "can_pass_gate", false)
	goNoGoEvidence := requireString(t, goNoGoGate, "evidence_required")
	for _, want := range []string{
		packagePath + "#install_catalog_and_slot_mapping_report",
		checklistPath,
		activationPath,
		bindingCatalogPath,
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		if !strings.Contains(goNoGoEvidence, want) {
			t.Fatalf("go/no-go six installed role catalog evidence missing %s: %s", want, goNoGoEvidence)
		}
	}
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, goNoGo, "completion_rule")["required_before_final_decision"]), "six_installed_role_catalog_evidence_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "go_no_go_without_six_installed_role_catalog")

	productMatrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["required_before_completion_claim"]), "six_installed_role_catalog_evidence_verified")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, productMatrix, "completion_claim_policy")["non_sufficient_evidence"]), "six_installed_role_catalog_missing")
	productGate := findObjectByString(t, asObjectSlice(t, productMatrix["readiness_gates"]), "gate_id", "local_install_and_team_binding")
	for _, evidence := range []string{
		"team_settings_catalog_refresh_receipt_ref",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
		"secretary_and_five_advisors_slot_mapping_screenshot",
	} {
		requireStringSliceContains(t, asStringSlice(t, productGate["required_evidence"]), evidence)
	}

	evidenceMap := readJSON(t, filepath.Join(base, "tests", "role-studio-goal-completion-evidence-map-candidate.json"))
	installReq := findObjectByString(t, asObjectSlice(t, evidenceMap["active_goal_requirements"]), "requirement_id", "local_install_enabled")
	requireStringSliceContains(t, asStringSlice(t, installReq["required_authoritative_evidence"]), "six_installed_role_catalog_evidence_verified")
	requireStringSliceContains(t, asStringSlice(t, installReq["current_repo_candidate_evidence"]), packagePath+"#install_catalog_and_slot_mapping_report")
	requireStringSliceContains(t, asStringSlice(t, installReq["missing_authoritative_evidence"]), "six_installed_role_catalog_evidence_verified_missing")
	requireStringSliceContains(t, asStringSlice(t, requireObject(t, evidenceMap, "completion_claim_policy")["required_before_goal_complete"]), "six_installed_role_catalog_evidence_verified")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["goal_completion_barriers"]), "six_installed_role_catalog_missing")
	requireStringSliceContains(t, asStringSlice(t, evidenceMap["non_sufficient_evidence"]), "goal_complete_without_six_installed_role_catalog_verified")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresCapabilityRoleReferenceReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	activationPath := "install/install-runtime-activation-map-candidate.json"
	requirementsPath := "capability-role-requirements/sample-team-research.role-requirements.json"
	runtimeUsagePath := "usage/team-office-runtime-usage-candidate.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))

	for key, want := range map[string]string{
		"source_install_runtime_activation_map": activationPath,
		"source_capability_role_requirements":   requirementsPath,
		"source_runtime_usage_candidate":        runtimeUsagePath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "capability_role_reference_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "capability_role_reference_report")
	for _, sourceRef := range []string{
		requirementsPath,
		activationPath,
		runtimeUsagePath,
		"tests/p11-evidence-acceptance-checklist-candidate.json",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"capability_role_reference_after_install_verified",
		"capability_role_requirement_validation_receipt",
		"enabled_role_pack_version_ref",
		"entitlement_verification_receipt",
		"role_body_not_copied_receipt",
		"capability_role_picker_screenshot",
		"capability_invocation_candidate_ref",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"capability_role_reference_after_install_missing",
		"capability_role_requirement_validation_receipt_missing",
		"role_body_not_copied_receipt_missing",
		"capability_invocation_candidate_ref_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "capability_role_reference_after_install_verified")
}

func TestTeamOfficeP11GoLiveEvidencePackageRequiresSecretaryAppearanceAssetReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	preferencesPath := "appearance/secretary-appearance-preferences.json"
	rightsPath := "appearance/secretary-appearance-asset-rights-candidate.json"
	secretScanPath := "tests/artifact-secret-raw-asset-scan-candidate.json"
	pkg := readJSON(t, filepath.Join(base, packagePath))

	for key, want := range map[string]string{
		"source_secretary_appearance_preferences":  preferencesPath,
		"source_secretary_appearance_asset_rights": rightsPath,
	} {
		if got := requireString(t, pkg, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	schema := requireObject(t, pkg, "package_schema")
	requireStringSliceContains(t, asStringSlice(t, schema["required_top_level_records"]), "secretary_appearance_asset_report")

	section := findObjectByString(t, asObjectSlice(t, pkg["package_sections"]), "section_id", "secretary_appearance_asset_report")
	for _, sourceRef := range []string{
		preferencesPath,
		rightsPath,
		secretScanPath,
		"tests/p11-evidence-acceptance-checklist-candidate.json",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_source_refs"]), sourceRef)
	}
	for _, field := range []string{
		"secretary_appearance_gui_controls_verified",
		"voice_asset_ref",
		"vrm_asset_ref",
		"voice_provider_readiness",
		"vrm_provider_readiness",
		"voice_vrm_selection_screenshot",
		"asset_ref_validation_receipt",
		"asset_license_evidence_receipt",
		"raw_asset_absence_scan_receipt",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["required_fields"]), field)
	}
	for _, blocker := range []string{
		"secretary_appearance_gui_controls_verified_missing",
		"asset_license_evidence_receipt_missing",
		"raw_asset_absence_scan_receipt_missing",
		"raw_voice_or_vrm_asset_detected",
		"provider_ready_claim_when_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, section["blocking_if_missing"]), blocker)
	}

	decisionGate := requireObject(t, pkg, "decision_gate")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "secretary_appearance_gui_controls_verified")
	requireStringSliceContains(t, asStringSlice(t, decisionGate["required_before_go_live_pass"]), "secretary_appearance_asset_rights_verified")
}

func TestTeamOfficeCommercialGatesRequireSixInstalledRoleCatalogEvidence(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	activationPath := "install/install-runtime-activation-map-candidate.json"

	chain := readJSON(t, filepath.Join(base, "tests", "commercial-chain-verifier-candidate.json"))
	if got := requireString(t, chain, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("commercial chain source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}
	correlation := requireObject(t, chain, "correlation_requirements")
	for _, key := range []string{
		"team_settings_catalog_refresh_receipt_ref",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
	} {
		requireStringSliceContains(t, asStringSlice(t, correlation["required_correlation_keys"]), key)
	}
	teamStage := findObjectByString(t, asObjectSlice(t, chain["ordered_verification_gates"]), "stage_id", "team_settings_role_replacement")
	for _, proof := range []string{
		"truzhenos_team_settings_catalog_refresh_receipt",
		"six_replaceable_role_refs",
		"enabled_role_pack_version_refs",
		"slot_mapping_refs",
		"gui_screenshot_secretary_and_five_advisors_slot_mapping",
	} {
		requireStringSliceContains(t, asStringSlice(t, teamStage["required_authoritative_evidence"]), proof)
	}
	pass := requireObject(t, chain, "normal_commercialization_pass_conditions")
	if got := requireString(t, pass, "six_installed_role_catalog_evidence_verified"); got != "required" {
		t.Fatalf("normal_commercialization_pass_conditions.six_installed_role_catalog_evidence_verified = %s, want required", got)
	}

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_install_runtime_activation_map"); got != activationPath {
		t.Fatalf("go/no-go source_install_runtime_activation_map = %s, want %s", got, activationPath)
	}
	localStage := findObjectByString(t, asObjectSlice(t, goNoGo["required_stage_gates"]), "stage_id", "local_install_team_binding")
	for _, blocker := range []string{
		"six_role_catalog_incomplete",
		"enabled_role_pack_version_refs_missing",
		"slot_mapping_refs_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, localStage["blocked_by"]), blocker)
	}
	for _, proof := range []string{
		"six_installed_role_catalog_evidence_verified",
		"team settings installed role catalog screenshot recorded",
		"secretary and five advisors slot mapping recorded",
	} {
		requireStringSliceContains(t, asStringSlice(t, localStage["required_before_pass"]), proof)
	}
	terminalGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "six_installed_role_catalog_evidence_verified")
	requireBool(t, terminalGate, "can_pass_gate", false)
	if got := requireString(t, terminalGate, "current_result"); got != "pending" {
		t.Fatalf("six_installed_role_catalog_evidence_verified current_result = %s, want pending", got)
	}
	completion := requireObject(t, goNoGo, "completion_rule")
	requireBool(t, completion, "six_installed_role_catalog_evidence_required", true)
}

func TestTeamOfficeP11RunbookAndGoLiveEvidencePackageCollectProductionPromotionReceipts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	runbook := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-execution-runbook-candidate.json"))
	pkg := readJSON(t, filepath.Join(base, "tests", "p11-commercial-go-live-evidence-package-template.json"))

	expectedReceipts := []string{
		"production_promotion_gate_pass_receipt",
		"production_go_live_request_receipt",
		"real_payment_enable_request_receipt",
		"production_signed_download_enable_receipt",
		"production_listing_publish_request_receipt",
		"production_install_observability_receipt",
	}

	rawStages, ok := runbook["execution_sequence"].([]any)
	if !ok {
		t.Fatalf("execution_sequence missing")
	}
	foundStage := false
	for _, raw := range rawStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("runbook stage = %T", raw)
		}
		if requireString(t, stage, "stage_id") != "production_promotion_controls" {
			continue
		}
		foundStage = true
		if got := requireString(t, stage, "actor"); got != "organizer_coordinator_recorder" {
			t.Fatalf("production_promotion_controls actor = %s", got)
		}
		targets := strings.Join(asStringSlice(t, stage["target_repos"]), "\n")
		for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos", "truzhen-cloud"} {
			if !strings.Contains(targets, repo) {
				t.Fatalf("production_promotion_controls target_repos missing %s", repo)
			}
		}
		evidence := strings.Join(asStringSlice(t, stage["required_authoritative_evidence"]), "\n")
		for _, receipt := range expectedReceipts {
			if !strings.Contains(evidence, receipt) {
				t.Fatalf("production_promotion_controls required_authoritative_evidence missing %s", receipt)
			}
		}
		if got := requireString(t, stage, "pass_condition"); !strings.Contains(got, "Owner go/no-go") {
			t.Fatalf("production_promotion_controls pass_condition missing Owner go/no-go: %s", got)
		}
		if got := requireString(t, stage, "fail_fast_if_missing"); !strings.Contains(got, "production") {
			t.Fatalf("production_promotion_controls fail_fast_if_missing missing production marker: %s", got)
		}
		outputs := strings.Join(asStringSlice(t, stage["output_record_refs"]), "\n")
		for _, output := range []string{"p11_go_live_evidence_package.production_promotion_receipts", "commercial_receipt_chain.production_promotion_receipts"} {
			if !strings.Contains(outputs, output) {
				t.Fatalf("production_promotion_controls output_record_refs missing %s", output)
			}
		}
	}
	if !foundStage {
		t.Fatalf("missing runbook stage production_promotion_controls")
	}

	rawSections, ok := pkg["package_sections"].([]any)
	if !ok {
		t.Fatalf("package_sections missing")
	}
	foundSection := false
	for _, raw := range rawSections {
		section, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("package section = %T", raw)
		}
		if requireString(t, section, "section_id") != "production_promotion_receipts" {
			continue
		}
		foundSection = true
		if got := requireString(t, section, "section_status"); got != "not_run" {
			t.Fatalf("production_promotion_receipts section_status = %s", got)
		}
		fields := strings.Join(asStringSlice(t, section["required_fields"]), "\n")
		for _, receipt := range expectedReceipts {
			if !strings.Contains(fields, receipt) {
				t.Fatalf("production_promotion_receipts required_fields missing %s", receipt)
			}
		}
		blockers := strings.Join(asStringSlice(t, section["blocking_if_missing"]), "\n")
		for _, receipt := range []string{"production_go_live_request_receipt", "production_listing_publish_request_receipt"} {
			if !strings.Contains(blockers, receipt) {
				t.Fatalf("production_promotion_receipts blocking_if_missing missing %s", receipt)
			}
		}
	}
	if !foundSection {
		t.Fatalf("missing package section production_promotion_receipts")
	}

	schema := requireObject(t, pkg, "package_schema")
	topRecords := strings.Join(asStringSlice(t, schema["required_top_level_records"]), "\n")
	if !strings.Contains(topRecords, "production_promotion_receipts") {
		t.Fatalf("package_schema.required_top_level_records missing production_promotion_receipts")
	}
	decisionGate := requireObject(t, pkg, "decision_gate")
	requiredBefore := strings.Join(asStringSlice(t, decisionGate["required_before_go_live_pass"]), "\n")
	if !strings.Contains(requiredBefore, "production_promotion_receipts_verified") {
		t.Fatalf("decision_gate.required_before_go_live_pass missing production_promotion_receipts_verified")
	}
}

func TestTeamOfficeOwnerAuthorizationEvidenceIntakeBlocksCrossRepoWorkUntilScopeRecorded(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	intakePath := "integration/owner-authorization-evidence-intake-candidate.json"
	intake := readJSON(t, filepath.Join(base, intakePath))

	requireBool(t, intake, "candidate_only", true)
	requireBool(t, intake, "non_formal", true)
	if got := requireString(t, intake, "status"); got != "missing_owner_authorization" {
		t.Fatalf("status = %s, want missing_owner_authorization", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":              "role-pack-candidate-set://team-office-v0",
		"authorization_scope_contract":   "integration/cross-repo-execution-readiness-package.json",
		"execution_cards_ref":            "integration/cross-repo-execution-cards.json",
		"p11_runbook_ref":                "tests/p11-sandbox-execution-runbook-candidate.json",
		"go_live_evidence_package_ref":   "tests/p11-commercial-go-live-evidence-package-template.json",
		"accepted_authorization_card":    "/Users/li/.config/superpowers/worktrees/truzhen-packs/gui-capability-pack-test-plan/docs/plans/role-pack-studio-cross-repo-execution-authorization-20260704.md",
		"authorization_truth_source":     "Owner explicit authorization in the current thread or signed authorization card; not inferred from plan text, candidate assets, prior runs, or narrow test success.",
		"evidence_output_policy":         "record_refs_only_in_truzhen_packs_raw_evidence_stays_in_execution_repos",
		"current_authorization_boundary": "no_cross_repo_work_without_recorded_matching_owner_authorization",
	} {
		if got := requireString(t, intake, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	requiredQuote := strings.Join(asStringSlice(t, intake["required_owner_quote_contains"]), "\n")
	for _, phrase := range []string{
		"授权按角色制作台跨仓执行授权卡",
		"truzhen-client-web-desktop",
		"truzhenos",
		"truzhen-cloud",
		"truzhen-contracts",
		"truzhen-packs",
		"用户视角 GUI 智能体",
		"sandbox payment",
		"不真实支付",
		"不生产发布",
		"不推送",
	} {
		if !strings.Contains(requiredQuote, phrase) {
			t.Fatalf("required_owner_quote_contains missing %s", phrase)
		}
	}

	rawRepos, ok := intake["authorized_repository_scopes"].([]any)
	if !ok {
		t.Fatalf("authorized_repository_scopes missing")
	}
	requiredRepos := map[string]bool{
		"truzhen-client-web-desktop": false,
		"truzhenos":                  false,
		"truzhen-cloud":              false,
		"truzhen-contracts":          false,
		"truzhen-packs-current":      false,
	}
	for _, raw := range rawRepos {
		repo, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("repository scope = %T", raw)
		}
		repoID := requireString(t, repo, "repo_id")
		if _, ok := requiredRepos[repoID]; !ok {
			t.Fatalf("unexpected repo_id %s", repoID)
		}
		requiredRepos[repoID] = true
		if requireString(t, repo, "repo_path") == "" {
			t.Fatalf("%s repo_path missing", repoID)
		}
		if got := requireString(t, repo, "authorization_status"); got != "missing" {
			t.Fatalf("%s authorization_status = %s, want missing", repoID, got)
		}
		if got := requireString(t, repo, "required_status_command"); got != "git status --short --branch" {
			t.Fatalf("%s required_status_command = %s", repoID, got)
		}
		for _, key := range []string{"allowed_actions_after_authorization", "required_evidence_outputs", "forbidden_boundaries"} {
			if values := asStringSlice(t, repo[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", repoID, key)
			}
		}
	}
	for repoID, seen := range requiredRepos {
		if !seen {
			t.Fatalf("missing repository scope %s", repoID)
		}
	}

	forbidden := strings.Join(asStringSlice(t, intake["forbidden_actions"]), "\n")
	for _, action := range []string{
		"real_payment",
		"production_publish",
		"push_or_merge",
		"store_raw_secret",
		"store_raw_voice_asset",
		"store_raw_vrm_asset",
		"mark_complete_from_candidate_assets",
		"start_cross_repo_work_without_owner_authorization",
	} {
		if !strings.Contains(forbidden, action) {
			t.Fatalf("forbidden_actions missing %s", action)
		}
	}

	evidence := requireObject(t, intake, "current_authorization_evidence")
	if got := requireString(t, evidence, "status"); got != "missing" {
		t.Fatalf("current_authorization_evidence.status = %s, want missing", got)
	}
	requireBool(t, evidence, "owner_thread_quote_required", true)
	for _, key := range []string{"owner_thread_quote", "recorded_by", "recorded_at", "evidence_ref"} {
		if got, ok := evidence[key].(string); !ok || got != "" {
			t.Fatalf("current_authorization_evidence.%s = %v, want empty string", key, evidence[key])
		}
	}

	gate := requireObject(t, intake, "cross_repo_work_gate")
	requireBool(t, gate, "can_start_cross_repo_work", false)
	for _, key := range []string{"required_status_before_cross_repo_work", "blocking_reason", "must_record_before_cross_repo_work"} {
		if key == "must_record_before_cross_repo_work" {
			requireBool(t, gate, key, true)
			continue
		}
		if got := requireString(t, gate, key); got == "" {
			t.Fatalf("cross_repo_work_gate.%s missing", key)
		}
	}

	controls := requireObject(t, intake, "expiration_and_scope_controls")
	for _, key := range []string{"per_repo_status_required", "isolated_worktree_required", "must_reconfirm_if_repos_or_actions_change", "authorization_expires_on_thread_change"} {
		requireBool(t, controls, key, true)
	}
	fields := strings.Join(asStringSlice(t, requireObject(t, intake, "evidence_recording_contract")["required_fields"]), "\n")
	for _, field := range []string{"owner_quote", "repo_status_snapshot", "allowed_actions_by_repo", "forbidden_actions_ack", "evidence_output_refs", "organizer_recorder", "recorded_at"} {
		if !strings.Contains(fields, field) {
			t.Fatalf("evidence_recording_contract.required_fields missing %s", field)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "owner_authorization_evidence_intake"); got != intakePath {
		t.Fatalf("owner_authorization_evidence_intake = %s, want %s", got, intakePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, intakePath) {
		t.Fatalf("candidate set artifact_files missing %s", intakePath)
	}
	readiness := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-readiness-package.json"))
	if got := requireString(t, readiness, "owner_authorization_evidence_intake_ref"); got != intakePath {
		t.Fatalf("readiness package owner_authorization_evidence_intake_ref = %s, want %s", got, intakePath)
	}
	runbook := readJSON(t, filepath.Join(base, "tests", "p11-sandbox-execution-runbook-candidate.json"))
	if got := requireString(t, runbook, "owner_authorization_evidence_intake_ref"); got != intakePath {
		t.Fatalf("runbook owner_authorization_evidence_intake_ref = %s, want %s", got, intakePath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "owner_authorization_evidence_intake_verified") {
		t.Fatalf("product readiness matrix missing owner_authorization_evidence_intake_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	intakeInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != intakePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, intakePath))
		if err != nil {
			t.Fatalf("read %s: %v", intakePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", intakePath, got, wantHash)
		}
		intakeInManifest = true
	}
	if !intakeInManifest {
		t.Fatalf("artifact manifest missing %s", intakePath)
	}
}

func TestTeamOfficeCommercialCrossRepoExecutionQueueOrdersAuthorizedStagesBeforeCommercialRun(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	queue := readJSON(t, filepath.Join(base, queuePath))

	requireBool(t, queue, "candidate_only", true)
	requireBool(t, queue, "non_formal", true)
	requireBool(t, queue, "can_start_cross_repo_execution", false)
	if got := requireString(t, queue, "queue_status"); got != "blocked_pending_owner_authorization_and_stage_evidence" {
		t.Fatalf("queue_status = %s, want blocked_pending_owner_authorization_and_stage_evidence", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                         "role-pack-candidate-set://team-office-v0",
		"owner_authorization_evidence_intake_ref":   "integration/owner-authorization-evidence-intake-candidate.json",
		"execution_readiness_package_ref":           "integration/cross-repo-execution-readiness-package.json",
		"cross_repo_execution_cards_ref":            "integration/cross-repo-execution-cards.json",
		"p11_runbook_ref":                           "tests/p11-sandbox-execution-runbook-candidate.json",
		"go_live_evidence_package_template_ref":     "tests/p11-commercial-go-live-evidence-package-template.json",
		"commercial_chain_verifier_ref":             "tests/commercial-chain-verifier-candidate.json",
		"product_readiness_evidence_matrix_ref":     "tests/product-readiness-evidence-matrix.json",
		"p0_p11_commercialization_blocker_register": "tests/p0-p11-commercialization-blocker-register-candidate.json",
	} {
		if got := requireString(t, queue, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	order := strings.Join(asStringSlice(t, queue["execution_order"]), "\n")
	expectedOrder := []string{
		"owner_authorization_intake",
		"contracts_schema_review",
		"backend_role_candidate_gate_receipt",
		"frontend_user_view_gui_flow",
		"cloud_sandbox_upload_purchase_download",
		"local_install_team_binding",
		"negative_cases_and_observability",
		"independent_acceptance_go_no_go",
		"production_promotion_controls",
	}
	for _, stageID := range expectedOrder {
		if !strings.Contains(order, stageID) {
			t.Fatalf("execution_order missing %s", stageID)
		}
	}

	rawEntries, ok := queue["execution_entries"].([]any)
	if !ok {
		t.Fatalf("execution_entries missing")
	}
	if len(rawEntries) != len(expectedOrder) {
		t.Fatalf("execution_entries len = %d, want %d", len(rawEntries), len(expectedOrder))
	}
	seen := map[string]bool{}
	for i, raw := range rawEntries {
		entry, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("execution entry = %T", raw)
		}
		stageID := requireString(t, entry, "stage_id")
		if stageID != expectedOrder[i] {
			t.Fatalf("execution_entries[%d] stage_id = %s, want %s", i, stageID, expectedOrder[i])
		}
		seen[stageID] = true
		requireStringIn(t, requireString(t, entry, "status"), "blocked_pending_owner_authorization", "blocked_pending_prior_stage_evidence")
		if got := requireString(t, entry, "owner_authorization_status"); got != "missing" {
			t.Fatalf("%s owner_authorization_status = %s, want missing", stageID, got)
		}
		if got := requireString(t, entry, "evidence_status"); got != "pending" {
			t.Fatalf("%s evidence_status = %s, want pending", stageID, got)
		}
		requireBool(t, entry, "can_start_now", false)
		for _, key := range []string{
			"depends_on",
			"target_repositories",
			"required_input_refs",
			"required_evidence_outputs",
			"evidence_record_targets",
			"verification_commands_or_checks",
		} {
			if values := asStringSlice(t, entry[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
		gate := requireObject(t, entry, "cross_repo_work_gate")
		requireBool(t, gate, "can_start_cross_repo_work", false)
		if got := requireString(t, gate, "required_status_before_cross_repo_work"); got == "" {
			t.Fatalf("%s gate required_status_before_cross_repo_work missing", stageID)
		}
		if got := requireString(t, gate, "blocking_reason"); got == "" {
			t.Fatalf("%s gate blocking_reason missing", stageID)
		}
	}
	for _, stageID := range expectedOrder {
		if !seen[stageID] {
			t.Fatalf("missing execution entry %s", stageID)
		}
	}

	stageExpectations := map[string][]string{
		"owner_authorization_intake":             {"truzhen-packs-current", "owner_authorization_evidence_intake"},
		"contracts_schema_review":                {"truzhen-contracts", "schema_impact_report"},
		"backend_role_candidate_gate_receipt":    {"truzhenos", "RolePackCandidate_receipt"},
		"frontend_user_view_gui_flow":            {"truzhen-client-web-desktop", "gui_operation_log"},
		"cloud_sandbox_upload_purchase_download": {"truzhen-cloud", "cloud_upload_receipt", "sandbox_payment_receipt", "download_receipt"},
		"local_install_team_binding":             {"truzhenos", "install_receipt", "team_binding_receipt"},
		"negative_cases_and_observability":       {"multi_repo", "blocked_receipt"},
		"independent_acceptance_go_no_go":        {"acceptance", "independent_acceptance_signoff", "owner_go_no_go_decision"},
	}
	for _, raw := range rawEntries {
		entry := raw.(map[string]any)
		stageID := requireString(t, entry, "stage_id")
		combined := strings.Join(asStringSlice(t, entry["target_repositories"]), "\n") + "\n" +
			strings.Join(asStringSlice(t, entry["required_evidence_outputs"]), "\n") + "\n" +
			strings.Join(asStringSlice(t, entry["evidence_record_targets"]), "\n")
		for _, want := range stageExpectations[stageID] {
			if !strings.Contains(combined, want) {
				t.Fatalf("%s missing expected queue token %s", stageID, want)
			}
		}
	}

	completionGate := requireObject(t, queue, "completion_gate")
	requireBool(t, completionGate, "can_mark_goal_complete_from_queue", false)
	for _, key := range []string{"required_before_first_cross_repo_run", "required_before_commercial_ready", "commercial_status_after_queue_only"} {
		if got := requireString(t, completionGate, key); got == "" {
			t.Fatalf("completion_gate.%s missing", key)
		}
	}
	nonSufficient := strings.Join(asStringSlice(t, queue["non_sufficient_evidence"]), "\n")
	for _, item := range []string{"queue_defined_without_cross_repo_receipts", "owner_authorization_missing", "candidate_json_without_gui_evidence", "single_repo_tests_only"} {
		if !strings.Contains(nonSufficient, item) {
			t.Fatalf("non_sufficient_evidence missing %s", item)
		}
	}
	forbidden := strings.Join(asStringSlice(t, queue["forbidden"]), "\n")
	for _, item := range []string{"start_cross_repo_work_without_owner_authorization", "real_payment", "production_publish", "mark_commercial_ready_from_queue_only", "store_raw_secret_or_asset"} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_cross_repo_execution_queue"); got != queuePath {
		t.Fatalf("commercial_cross_repo_execution_queue = %s, want %s", got, queuePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, queuePath) {
		t.Fatalf("candidate set artifact_files missing %s", queuePath)
	}
	readiness := readJSON(t, filepath.Join(base, "integration", "cross-repo-execution-readiness-package.json"))
	if got := requireString(t, readiness, "commercial_cross_repo_execution_queue_ref"); got != queuePath {
		t.Fatalf("readiness package commercial_cross_repo_execution_queue_ref = %s, want %s", got, queuePath)
	}
	intake := readJSON(t, filepath.Join(base, "integration", "owner-authorization-evidence-intake-candidate.json"))
	if got := requireString(t, intake, "commercial_cross_repo_execution_queue_ref"); got != queuePath {
		t.Fatalf("intake commercial_cross_repo_execution_queue_ref = %s, want %s", got, queuePath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "commercial_cross_repo_execution_queue_verified") {
		t.Fatalf("product readiness matrix missing commercial_cross_repo_execution_queue_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	queueInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != queuePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, queuePath))
		if err != nil {
			t.Fatalf("read %s: %v", queuePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", queuePath, got, wantHash)
		}
		queueInManifest = true
	}
	if !queueInManifest {
		t.Fatalf("artifact manifest missing %s", queuePath)
	}
}

func TestTeamOfficeCommercialCrossRepoEvidenceLedgerDefinesEvidenceIdsForEveryQueueStage(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	ledgerPath := "docs/commercial-cross-repo-evidence-ledger.json"
	ledger := readJSON(t, filepath.Join(base, ledgerPath))

	requireBool(t, ledger, "candidate_only", true)
	requireBool(t, ledger, "non_formal", true)
	if got := requireString(t, ledger, "ledger_status"); got != "pending_owner_authorization" {
		t.Fatalf("ledger_status = %s, want pending_owner_authorization", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                     "role-pack-candidate-set://team-office-v0",
		"commercial_cross_repo_queue_ref":       "integration/commercial-cross-repo-execution-queue-candidate.json",
		"owner_authorization_intake_ref":        "integration/owner-authorization-evidence-intake-candidate.json",
		"p11_verification_record_template_ref":  "tests/p11-normal-commercialization-verification-record-template.json",
		"go_live_evidence_package_template_ref": "tests/p11-commercial-go-live-evidence-package-template.json",
		"commercial_chain_verifier_ref":         "tests/commercial-chain-verifier-candidate.json",
		"evidence_payload_policy":               "refs_hashes_and_redacted_summaries_only",
	} {
		if got := requireString(t, ledger, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	policy := requireObject(t, ledger, "evidence_id_policy")
	requireBool(t, policy, "unique_evidence_id_required", true)
	requireBool(t, policy, "raw_payload_forbidden", true)
	statuses := strings.Join(asStringSlice(t, policy["allowed_statuses"]), "\n")
	for _, status := range []string{"pending_authorization", "evidence_recorded_pending_review", "evidence_complete_verified", "blocked_missing_authoritative_evidence", "failed_requires_issue"} {
		if !strings.Contains(statuses, status) {
			t.Fatalf("allowed_statuses missing %s", status)
		}
	}
	forbiddenPayloads := strings.Join(asStringSlice(t, policy["forbidden_payloads"]), "\n")
	for _, payload := range []string{"raw_payment_token", "cloud_access_token", "signed_download_url_secret", "raw_voice_asset", "raw_vrm_asset", "browser_cookie"} {
		if !strings.Contains(forbiddenPayloads, payload) {
			t.Fatalf("forbidden_payloads missing %s", payload)
		}
	}

	rawRows, ok := ledger["evidence_rows"].([]any)
	if !ok {
		t.Fatalf("evidence_rows missing")
	}
	expectedStages := map[string][]string{
		"owner_authorization_intake":             {"role_studio_auth_intake_recorded", "truzhen-packs-current"},
		"contracts_schema_review":                {"role_studio_contracts_schema_impact", "truzhen-contracts"},
		"backend_role_candidate_gate_receipt":    {"role_studio_backend_candidate_receipts", "truzhenos"},
		"frontend_user_view_gui_flow":            {"role_studio_gui_user_view_walkthrough", "truzhen-client-web-desktop"},
		"cloud_sandbox_upload_purchase_download": {"role_studio_cloud_sandbox_receipts", "truzhen-cloud"},
		"local_install_team_binding":             {"role_studio_install_team_binding_receipts", "truzhenos"},
		"negative_cases_and_observability":       {"role_studio_negative_observability_receipts", "multi_repo"},
		"independent_acceptance_go_no_go":        {"role_studio_independent_acceptance_go_no_go", "acceptance"},
		"production_promotion_controls":          {"role_studio_production_promotion_receipts", "multi_repo"},
	}
	if len(rawRows) != len(expectedStages) {
		t.Fatalf("evidence_rows len = %d, want %d", len(rawRows), len(expectedStages))
	}
	seenEvidenceIDs := map[string]bool{}
	seenStages := map[string]bool{}
	for _, raw := range rawRows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("evidence row = %T", raw)
		}
		stageID := requireString(t, row, "stage_id")
		expected, ok := expectedStages[stageID]
		if !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		seenStages[stageID] = true
		evidenceID := requireString(t, row, "evidence_id")
		if seenEvidenceIDs[evidenceID] {
			t.Fatalf("duplicate evidence_id %s", evidenceID)
		}
		seenEvidenceIDs[evidenceID] = true
		if evidenceID != expected[0] {
			t.Fatalf("%s evidence_id = %s, want %s", stageID, evidenceID, expected[0])
		}
		if got := requireString(t, row, "target_repository"); got != expected[1] {
			t.Fatalf("%s target_repository = %s, want %s", stageID, got, expected[1])
		}
		if got := requireString(t, row, "current_status"); got != "pending_authorization" {
			t.Fatalf("%s current_status = %s, want pending_authorization", stageID, got)
		}
		requireBool(t, row, "required_before_stage_complete", true)
		requireBool(t, row, "raw_payload_forbidden", true)
		for _, key := range []string{"expected_evidence_location", "required_evidence_refs", "writeback_targets", "blocking_if_missing"} {
			if values := asStringSlice(t, row[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
	}
	for stageID := range expectedStages {
		if !seenStages[stageID] {
			t.Fatalf("missing evidence row for stage %s", stageID)
		}
	}

	completion := requireObject(t, ledger, "completion_gate")
	requireBool(t, completion, "can_mark_ledger_complete", false)
	for _, key := range []string{"required_before_ledger_complete", "required_before_goal_completion", "non_sufficient_evidence"} {
		if key == "non_sufficient_evidence" {
			if values := asStringSlice(t, completion[key]); len(values) == 0 {
				t.Fatalf("completion_gate.%s missing", key)
			}
			continue
		}
		if got := requireString(t, completion, key); got == "" {
			t.Fatalf("completion_gate.%s missing", key)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_cross_repo_evidence_ledger"); got != ledgerPath {
		t.Fatalf("commercial_cross_repo_evidence_ledger = %s, want %s", got, ledgerPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, ledgerPath) {
		t.Fatalf("candidate set artifact_files missing %s", ledgerPath)
	}
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	if got := requireString(t, queue, "commercial_cross_repo_evidence_ledger_ref"); got != ledgerPath {
		t.Fatalf("queue commercial_cross_repo_evidence_ledger_ref = %s, want %s", got, ledgerPath)
	}
	pkg := readJSON(t, filepath.Join(base, "tests", "p11-commercial-go-live-evidence-package-template.json"))
	if got := requireString(t, pkg, "commercial_cross_repo_evidence_ledger_ref"); got != ledgerPath {
		t.Fatalf("go-live evidence package commercial_cross_repo_evidence_ledger_ref = %s, want %s", got, ledgerPath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "commercial_cross_repo_evidence_ledger_verified") {
		t.Fatalf("product readiness matrix missing commercial_cross_repo_evidence_ledger_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	ledgerInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != ledgerPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, ledgerPath))
		if err != nil {
			t.Fatalf("read %s: %v", ledgerPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", ledgerPath, got, wantHash)
		}
		ledgerInManifest = true
	}
	if !ledgerInManifest {
		t.Fatalf("artifact manifest missing %s", ledgerPath)
	}
}

func TestTeamOfficeCommercialExecutionEvidenceCollectsGuiAPITraceabilityWriteback(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	queue := readJSON(t, filepath.Join(base, "integration", "commercial-cross-repo-execution-queue-candidate.json"))
	ledger := readJSON(t, filepath.Join(base, "docs", "commercial-cross-repo-evidence-ledger.json"))
	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))

	queueEntry := findObjectByString(t, asObjectSlice(t, queue["execution_entries"]), "stage_id", "frontend_user_view_gui_flow")
	for _, output := range []string{
		"gui_api_traceability_matrix_verified",
		"commercial_api_endpoint_trace",
		"receipt_field_trace",
		"gui_evidence_slot_trace",
	} {
		requireStringSliceContains(t, asStringSlice(t, queueEntry["required_evidence_outputs"]), output)
	}
	for _, target := range []string{
		"tests/p11-commercial-go-live-evidence-package-template.json#gui_api_traceability_report",
		"tests/p11-evidence-acceptance-checklist-candidate.json#gui_api_traceability_matrix",
		"tests/product-stage-frontend-backend-closure-report-candidate.json#gui_api_traceability_matrix_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, queueEntry["evidence_record_targets"]), target)
	}
	checks := strings.Join(asStringSlice(t, queueEntry["verification_commands_or_checks"]), "\n")
	for _, check := range []string{"GUI/API traceability matrix check", "commercial API endpoint and receipt field trace"} {
		if !strings.Contains(checks, check) {
			t.Fatalf("frontend_user_view_gui_flow verification checks missing %s", check)
		}
	}

	ledgerRow := findObjectByString(t, asObjectSlice(t, ledger["evidence_rows"]), "evidence_id", "role_studio_gui_user_view_walkthrough")
	for _, ref := range []string{
		"gui_api_traceability_matrix_verified",
		"commercial_api_endpoint_refs",
		"required_receipt_or_result_fields",
		"required_gui_evidence_slots",
	} {
		requireStringSliceContains(t, asStringSlice(t, ledgerRow["required_evidence_refs"]), ref)
	}
	for _, target := range []string{
		"tests/p11-commercial-go-live-evidence-package-template.json#gui_api_traceability_report",
		"tests/p11-evidence-acceptance-checklist-candidate.json#gui_api_traceability_matrix",
		"tests/product-stage-frontend-backend-closure-report-candidate.json#gui_api_traceability_matrix_verified",
	} {
		requireStringSliceContains(t, asStringSlice(t, ledgerRow["writeback_targets"]), target)
	}
	for _, blocker := range []string{
		"gui_api_traceability_matrix_missing",
		"commercial_api_endpoint_trace_missing",
		"receipt_field_trace_missing",
	} {
		requireStringSliceContains(t, asStringSlice(t, ledgerRow["blocking_if_missing"]), blocker)
	}

	writeback := requireObject(t, requireObject(t, readiness, "current_blockers"), "evidence_writeback_summary")
	requireStringSliceContains(t, asStringSlice(t, writeback["required_before_pass"]), "gui_api_traceability_matrix_verified")
	frontendBackend := requireObject(t, requireObject(t, readiness, "current_blockers"), "frontend_backend_acceptance")
	requireStringSliceContains(t, asStringSlice(t, frontendBackend["required_before_pass"]), "gui_api_traceability_matrix_verified")

	goNoGoStage := findObjectByString(t, asObjectSlice(t, goNoGo["required_stage_gates"]), "stage_id", "frontend_user_view_gui_flow")
	requireStringSliceContains(t, asStringSlice(t, goNoGoStage["required_before_pass"]), "GUI/API traceability matrix verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGoStage["blocked_by"]), "gui_api_traceability_matrix_missing")
}

func TestTeamOfficeCommercialGoNoGoGateBlocksCommercialReadyUntilEveryEvidenceGatePasses(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	gatePath := "tests/commercial-go-no-go-gate-candidate.json"
	gate := readJSON(t, filepath.Join(base, gatePath))

	requireBool(t, gate, "candidate_only", true)
	requireBool(t, gate, "non_formal", true)
	if got := requireString(t, gate, "decision_status"); got != "blocked_not_ready_for_commercial_go_live" {
		t.Fatalf("decision_status = %s, want blocked_not_ready_for_commercial_go_live", got)
	}
	requireBool(t, gate, "can_mark_commercial_ready", false)
	requireBool(t, gate, "can_request_owner_go_live_signoff", false)
	for key, want := range map[string]string{
		"candidate_set_ref":                        "role-pack-candidate-set://team-office-v0",
		"source_execution_queue":                   "integration/commercial-cross-repo-execution-queue-candidate.json",
		"source_cross_repo_evidence_ledger":        "docs/commercial-cross-repo-evidence-ledger.json",
		"source_commercial_chain_verifier":         "tests/commercial-chain-verifier-candidate.json",
		"source_p11_acceptance_gate":               "tests/p11-normal-commercialization-acceptance-gate-candidate.json",
		"source_p11_go_live_evidence_package":      "tests/p11-commercial-go-live-evidence-package-template.json",
		"source_commercial_go_live_approval":       "commerce/commercial-go-live-approval-candidate.json",
		"source_product_readiness_evidence_matrix": "tests/product-readiness-evidence-matrix.json",
		"next_required_authorization":              "owner_authorization_intake",
		"next_authorization_card":                  "integration/owner-authorization-evidence-intake-candidate.json",
	} {
		if got := requireString(t, gate, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	rawStageGates, ok := gate["required_stage_gates"].([]any)
	if !ok {
		t.Fatalf("required_stage_gates missing")
	}
	expectedStages := map[string]string{
		"owner_authorization_intake":             "role_studio_auth_intake_recorded",
		"contracts_schema_review":                "role_studio_contracts_schema_impact",
		"backend_role_candidate_gate_receipt":    "role_studio_backend_candidate_receipts",
		"frontend_user_view_gui_flow":            "role_studio_gui_user_view_walkthrough",
		"cloud_sandbox_upload_purchase_download": "role_studio_cloud_sandbox_receipts",
		"local_install_team_binding":             "role_studio_install_team_binding_receipts",
		"negative_cases_and_observability":       "role_studio_negative_observability_receipts",
		"independent_acceptance_go_no_go":        "role_studio_independent_acceptance_go_no_go",
		"production_promotion_controls":          "role_studio_production_promotion_receipts",
	}
	if len(rawStageGates) != len(expectedStages) {
		t.Fatalf("required_stage_gates len = %d, want %d", len(rawStageGates), len(expectedStages))
	}
	seenStages := map[string]bool{}
	for _, raw := range rawStageGates {
		stageGate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage gate = %T", raw)
		}
		stageID := requireString(t, stageGate, "stage_id")
		wantEvidenceID, ok := expectedStages[stageID]
		if !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		seenStages[stageID] = true
		if got := requireString(t, stageGate, "evidence_id"); got != wantEvidenceID {
			t.Fatalf("%s evidence_id = %s, want %s", stageID, got, wantEvidenceID)
		}
		if got := requireString(t, stageGate, "required_final_status"); got != "evidence_complete_verified" {
			t.Fatalf("%s required_final_status = %s", stageID, got)
		}
		if got := requireString(t, stageGate, "current_status"); got != "pending_authorization" {
			t.Fatalf("%s current_status = %s, want pending_authorization", stageID, got)
		}
		requireBool(t, stageGate, "can_pass_gate", false)
		for _, key := range []string{"blocked_by", "required_before_pass"} {
			if values := asStringSlice(t, stageGate[key]); len(values) == 0 {
				t.Fatalf("%s %s must not be empty", stageID, key)
			}
		}
	}
	for stageID := range expectedStages {
		if !seenStages[stageID] {
			t.Fatalf("missing stage gate %s", stageID)
		}
	}

	rawTerminal, ok := gate["terminal_commercial_gates"].([]any)
	if !ok {
		t.Fatalf("terminal_commercial_gates missing")
	}
	requiredTerminal := map[string]bool{
		"no_raw_asset_payloads":                       false,
		"no_real_payment_without_authorization":       false,
		"no_production_publish_without_authorization": false,
		"candidate_formal_boundary_preserved":         false,
		"bundle_tree_sha256_consistent":               false,
		"independent_acceptance_signed":               false,
		"owner_go_no_go_decision_recorded":            false,
	}
	for _, raw := range rawTerminal {
		terminal, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("terminal gate = %T", raw)
		}
		gateID := requireString(t, terminal, "gate_id")
		if _, ok := requiredTerminal[gateID]; ok {
			requiredTerminal[gateID] = true
		}
		if got := requireString(t, terminal, "current_result"); got != "pending" {
			t.Fatalf("%s current_result = %s, want pending", gateID, got)
		}
		requireBool(t, terminal, "can_pass_gate", false)
		if got := requireString(t, terminal, "evidence_required"); got == "" {
			t.Fatalf("%s evidence_required missing", gateID)
		}
	}
	for gateID, seen := range requiredTerminal {
		if !seen {
			t.Fatalf("missing terminal gate %s", gateID)
		}
	}

	completion := requireObject(t, gate, "completion_rule")
	for _, key := range []string{
		"all_stage_gates_must_pass",
		"all_terminal_gates_must_pass",
		"all_authoritative_receipts_required",
		"p11_evidence_package_reviewed",
		"owner_go_no_go_decision_required",
		"independent_acceptance_required",
	} {
		requireBool(t, completion, key, true)
	}
	if got := requireString(t, completion, "required_final_decision_status"); got != "commercial_ready_candidate" {
		t.Fatalf("required_final_decision_status = %s", got)
	}

	nonSufficient := strings.Join(asStringSlice(t, gate["non_sufficient_evidence"]), "\n")
	for _, item := range []string{"candidate_assets_only", "ledger_template_only", "green_pack_tests_without_gui_cloud_install_receipts", "manual_summary_without_receipts"} {
		if !strings.Contains(nonSufficient, item) {
			t.Fatalf("non_sufficient_evidence missing %s", item)
		}
	}
	forbidden := strings.Join(asStringSlice(t, gate["forbidden"]), "\n")
	for _, item := range []string{"mark_commercial_ready_from_candidate_assets_only", "request_owner_go_live_signoff_before_all_gates_pass", "start_cross_repo_work_without_owner_authorization", "execute_real_payment_without_owner_authorization"} {
		if !strings.Contains(forbidden, item) {
			t.Fatalf("forbidden missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_go_no_go_gate"); got != gatePath {
		t.Fatalf("commercial_go_no_go_gate = %s, want %s", got, gatePath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, gatePath) {
		t.Fatalf("candidate set artifact_files missing %s", gatePath)
	}
	approval := readJSON(t, filepath.Join(base, "commerce", "commercial-go-live-approval-candidate.json"))
	if got := requireString(t, approval, "commercial_go_no_go_gate_ref"); got != gatePath {
		t.Fatalf("commercial-go-live approval commercial_go_no_go_gate_ref = %s, want %s", got, gatePath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "commercial_go_no_go_gate_verified") {
		t.Fatalf("product readiness matrix missing commercial_go_no_go_gate_verified")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	gateInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != gatePath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, gatePath))
		if err != nil {
			t.Fatalf("read %s: %v", gatePath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", gatePath, got, wantHash)
		}
		gateInManifest = true
	}
	if !gateInManifest {
		t.Fatalf("artifact manifest missing %s", gatePath)
	}
}

func TestTeamOfficeCommercialReadinessAndGoNoGoConsumeDownloadInstallAccessMatrix(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	matrixPath := "commerce/download-install-access-matrix.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	if got := requireString(t, readiness, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("readiness source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), "download_install_access_matrix_verified")
	readinessGate := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "download_install_access_matrix_verified")
	requireBool(t, readinessGate, "can_count_toward_commercial_ready", false)
	if evidence := requireString(t, readinessGate, "evidence_required"); !strings.Contains(evidence, matrixPath) || !strings.Contains(evidence, "unpaid_download") || !strings.Contains(evidence, "artifact_hash_mismatch") {
		t.Fatalf("readiness download_install_access_matrix_verified evidence_required = %s", evidence)
	}

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, goNoGo, "source_download_install_access_matrix"); got != matrixPath {
		t.Fatalf("go/no-go source_download_install_access_matrix = %s, want %s", got, matrixPath)
	}
	terminalGate := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "download_install_access_matrix_verified")
	requireBool(t, terminalGate, "can_pass_gate", false)
	if evidence := requireString(t, terminalGate, "evidence_required"); !strings.Contains(evidence, matrixPath) || !strings.Contains(evidence, "refund_revoked_download") || !strings.Contains(evidence, "version_unpublished_or_revoked_download_install") {
		t.Fatalf("go/no-go download_install_access_matrix_verified evidence_required = %s", evidence)
	}
	completion := requireObject(t, goNoGo, "completion_rule")
	requireBool(t, completion, "download_install_access_matrix_required", true)
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p11_negative_cases_without_download_install_access_matrix")
}

func TestTeamOfficeCommercialReadinessAndGoNoGoRequireSecretaryAppearanceAssetReport(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	packagePath := "tests/p11-commercial-go-live-evidence-package-template.json"
	preferencesPath := "appearance/secretary-appearance-preferences.json"
	rightsPath := "appearance/secretary-appearance-asset-rights-candidate.json"
	rawAssetScanPath := "tests/artifact-secret-raw-asset-scan-candidate.json"

	readiness := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	for key, want := range map[string]string{
		"source_p11_go_live_evidence_package":      packagePath,
		"source_secretary_appearance_preferences":  preferencesPath,
		"source_secretary_appearance_asset_rights": rightsPath,
		"source_artifact_secret_raw_asset_scan":    rawAssetScanPath,
	} {
		if got := requireString(t, readiness, key); got != want {
			t.Fatalf("readiness %s = %s, want %s", key, got, want)
		}
	}
	readinessTerminal := findObjectByString(t, asObjectSlice(t, readiness["terminal_checks"]), "gate_id", "secretary_appearance_asset_report_verified")
	if got := requireString(t, readinessTerminal, "current_result"); got != "pending" {
		t.Fatalf("readiness secretary appearance current_result = %s, want pending", got)
	}
	requireBool(t, readinessTerminal, "can_count_toward_commercial_ready", false)
	evidence := requireString(t, readinessTerminal, "evidence_required")
	for _, ref := range []string{packagePath, preferencesPath, rightsPath, rawAssetScanPath, "secretary_appearance_asset_report", "asset_license_evidence_receipt", "raw_asset_absence_scan_receipt"} {
		if !strings.Contains(evidence, ref) {
			t.Fatalf("readiness secretary appearance evidence_required missing %s: %s", ref, evidence)
		}
	}
	readinessPolicy := requireObject(t, readiness, "completion_claim_policy")
	for _, proof := range []string{"secretary_appearance_asset_report_verified", "secretary_appearance_asset_rights_verified", "secretary_appearance_gui_controls_verified"} {
		requireStringSliceContains(t, asStringSlice(t, readinessPolicy["required_before_completion_claim"]), proof)
	}

	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	for key, want := range map[string]string{
		"source_p11_go_live_evidence_package":      packagePath,
		"source_secretary_appearance_preferences":  preferencesPath,
		"source_secretary_appearance_asset_rights": rightsPath,
		"source_artifact_secret_raw_asset_scan":    rawAssetScanPath,
	} {
		if got := requireString(t, goNoGo, key); got != want {
			t.Fatalf("go/no-go %s = %s, want %s", key, got, want)
		}
	}
	goNoGoTerminal := findObjectByString(t, asObjectSlice(t, goNoGo["terminal_commercial_gates"]), "gate_id", "secretary_appearance_asset_report_verified")
	if got := requireString(t, goNoGoTerminal, "current_result"); got != "pending" {
		t.Fatalf("go/no-go secretary appearance current_result = %s, want pending", got)
	}
	requireBool(t, goNoGoTerminal, "required_final_value", true)
	requireBool(t, goNoGoTerminal, "can_pass_gate", false)
	evidence = requireString(t, goNoGoTerminal, "evidence_required")
	for _, ref := range []string{packagePath, preferencesPath, rightsPath, rawAssetScanPath, "secretary_appearance_asset_report", "voice_vrm_selection_screenshot", "asset_license_evidence_receipt", "raw_asset_absence_scan_receipt"} {
		if !strings.Contains(evidence, ref) {
			t.Fatalf("go/no-go secretary appearance evidence_required missing %s: %s", ref, evidence)
		}
	}
	completion := requireObject(t, goNoGo, "completion_rule")
	requireStringSliceContains(t, asStringSlice(t, completion["required_before_final_decision"]), "secretary_appearance_asset_report_verified")
	requireStringSliceContains(t, asStringSlice(t, goNoGo["non_sufficient_evidence"]), "p11_go_live_without_secretary_appearance_asset_report")

	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	policy := requireObject(t, matrix, "completion_claim_policy")
	for _, proof := range []string{"secretary_appearance_asset_report_verified", "secretary_appearance_asset_rights_verified"} {
		requireStringSliceContains(t, asStringSlice(t, policy["required_before_completion_claim"]), proof)
	}
}

func TestTeamOfficeCommercialGoNoGoGateRequiresP11ExecutionQueuePhaseDependencies(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	queuePath := "integration/commercial-cross-repo-execution-queue-candidate.json"
	gate := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	queue := readJSON(t, filepath.Join(base, queuePath))

	if got := requireString(t, gate, "source_execution_queue"); got != queuePath {
		t.Fatalf("source_execution_queue = %s, want %s", got, queuePath)
	}
	if got := requireString(t, gate, "source_p11_phase_dependency_contract"); got != "tests/p11-sandbox-run-request-candidate.json#phase_dependency_contract" {
		t.Fatalf("source_p11_phase_dependency_contract = %s", got)
	}
	dependencyGate := requireObject(t, gate, "p11_phase_dependency_go_no_go_gate")
	for key, want := range map[string]string{
		"phase_dependency_links_ref": queuePath + "#p11_phase_dependency_links",
		"required_completion_proof":  "p11_execution_queue_phase_dependencies_verified",
		"blocked_status_if_missing":  "blocked_previous_phase_evidence_missing",
	} {
		if got := requireString(t, dependencyGate, key); got != want {
			t.Fatalf("p11_phase_dependency_go_no_go_gate.%s = %s, want %s", key, got, want)
		}
	}
	for _, key := range []string{
		"all_phase_handoffs_required",
		"previous_receipts_required_before_owner_go_no_go",
		"payment_download_install_cannot_skip_upload_review_entitlement",
	} {
		requireBool(t, dependencyGate, key, true)
	}
	wantCount := len(asObjectSlice(t, queue["p11_phase_dependency_links"]))
	if got, ok := dependencyGate["expected_phase_dependency_count"].(float64); !ok || int(got) != wantCount {
		t.Fatalf("expected_phase_dependency_count = %v, want %d", dependencyGate["expected_phase_dependency_count"], wantCount)
	}

	rawTerminal, ok := gate["terminal_commercial_gates"].([]any)
	if !ok {
		t.Fatalf("terminal_commercial_gates missing")
	}
	found := false
	for _, raw := range rawTerminal {
		terminal, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("terminal gate = %T", raw)
		}
		if requireString(t, terminal, "gate_id") != "p11_execution_queue_phase_dependencies_verified" {
			continue
		}
		found = true
		if got := requireString(t, terminal, "current_result"); got != "pending" {
			t.Fatalf("p11 dependency terminal current_result = %s, want pending", got)
		}
		requireBool(t, terminal, "required_final_value", true)
		requireBool(t, terminal, "can_pass_gate", false)
		if got := requireString(t, terminal, "evidence_required"); !strings.Contains(got, "p11_phase_dependency_links") {
			t.Fatalf("p11 dependency terminal evidence_required = %s", got)
		}
	}
	if !found {
		t.Fatalf("terminal_commercial_gates missing p11_execution_queue_phase_dependencies_verified")
	}

	completion := requireObject(t, gate, "completion_rule")
	requireBool(t, completion, "p11_execution_queue_phase_dependencies_required", true)
	nonSufficient := strings.Join(asStringSlice(t, gate["non_sufficient_evidence"]), "\n")
	if !strings.Contains(nonSufficient, "final_receipts_without_phase_dependency_chain") {
		t.Fatalf("non_sufficient_evidence missing final_receipts_without_phase_dependency_chain")
	}
}

func TestTeamOfficeCommercialReadinessAndGoNoGoRequireProductionPromotionStageGate(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	verifier := readJSON(t, filepath.Join(base, "tests", "commercial-readiness-verifier-candidate.json"))
	goNoGo := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))

	blockers := requireObject(t, verifier, "current_blockers")
	rawVerifierStages, ok := blockers["required_stage_gates"].([]any)
	if !ok {
		t.Fatalf("verifier current_blockers.required_stage_gates missing")
	}
	foundVerifierStage := false
	for _, raw := range rawVerifierStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("verifier stage gate = %T", raw)
		}
		if requireString(t, stage, "stage_id") != "production_promotion_controls" {
			continue
		}
		foundVerifierStage = true
		if got := requireString(t, stage, "evidence_id"); got != "role_studio_production_promotion_receipts" {
			t.Fatalf("verifier production stage evidence_id = %s", got)
		}
		if got := requireString(t, stage, "authorization_status"); got != "missing" {
			t.Fatalf("verifier production stage authorization_status = %s", got)
		}
		requireBool(t, stage, "can_count_toward_commercial_ready", false)
		blockedBy := strings.Join(asStringSlice(t, stage["blocked_by"]), "\n")
		for _, item := range []string{
			"production_go_live_request_receipt_missing",
			"real_payment_enable_request_receipt_missing",
			"production_signed_download_enable_receipt_missing",
			"production_listing_publish_request_receipt_missing",
			"production_install_observability_receipt_missing",
		} {
			if !strings.Contains(blockedBy, item) {
				t.Fatalf("verifier production stage blocked_by missing %s", item)
			}
		}
	}
	if !foundVerifierStage {
		t.Fatalf("commercial readiness verifier missing production_promotion_controls stage gate")
	}

	writeback := requireObject(t, blockers, "evidence_writeback_summary")
	if got, ok := writeback["total_required_entry_count"].(float64); !ok || got != 9 {
		t.Fatalf("total_required_entry_count = %v, want 9", writeback["total_required_entry_count"])
	}
	if got, ok := writeback["total_pending_entry_count"].(float64); !ok || got != 9 {
		t.Fatalf("total_pending_entry_count = %v, want 9", writeback["total_pending_entry_count"])
	}
	writebackRequired := strings.Join(asStringSlice(t, writeback["required_before_pass"]), "\n")
	if !strings.Contains(writebackRequired, "production_promotion_receipts_verified") {
		t.Fatalf("evidence_writeback_summary.required_before_pass missing production_promotion_receipts_verified")
	}
	frontendBackend := requireObject(t, blockers, "frontend_backend_acceptance")
	frontendRequired := strings.Join(asStringSlice(t, frontendBackend["required_before_pass"]), "\n")
	if !strings.Contains(frontendRequired, "production_promotion_receipts_present") {
		t.Fatalf("frontend_backend_acceptance.required_before_pass missing production_promotion_receipts_present")
	}

	rawGoNoGoStages, ok := goNoGo["required_stage_gates"].([]any)
	if !ok {
		t.Fatalf("go/no-go required_stage_gates missing")
	}
	foundGoNoGoStage := false
	for _, raw := range rawGoNoGoStages {
		stage, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("go/no-go stage gate = %T", raw)
		}
		if requireString(t, stage, "stage_id") != "production_promotion_controls" {
			continue
		}
		foundGoNoGoStage = true
		if got := requireString(t, stage, "evidence_id"); got != "role_studio_production_promotion_receipts" {
			t.Fatalf("go/no-go production stage evidence_id = %s", got)
		}
		if got := requireString(t, stage, "required_final_status"); got != "evidence_complete_verified" {
			t.Fatalf("go/no-go production stage required_final_status = %s", got)
		}
		requireBool(t, stage, "can_pass_gate", false)
		requiredBeforePass := strings.Join(asStringSlice(t, stage["required_before_pass"]), "\n")
		for _, item := range []string{
			"production go-live request receipt recorded",
			"real payment enable decision receipt recorded",
			"production signed download enable receipt recorded",
			"production listing publish receipt recorded",
			"production install observability receipt recorded",
		} {
			if !strings.Contains(requiredBeforePass, item) {
				t.Fatalf("go/no-go production stage required_before_pass missing %s", item)
			}
		}
	}
	if !foundGoNoGoStage {
		t.Fatalf("commercial go/no-go gate missing production_promotion_controls stage gate")
	}
}

func TestTeamOfficeCommercialReadinessVerifierSummarizesCurrentBlockersAndWritebackCounts(t *testing.T) {
	base := filepath.Join("role-pack-candidates", "team-office-v0")
	verifierPath := "tests/commercial-readiness-verifier-candidate.json"
	verifier := readJSON(t, filepath.Join(base, verifierPath))

	requireBool(t, verifier, "candidate_only", true)
	requireBool(t, verifier, "non_formal", true)
	requireBool(t, verifier, "can_mark_commercial_ready", false)
	requireBool(t, verifier, "formal_write_allowed", false)
	if got := requireString(t, verifier, "verification_status"); got != "blocked_not_commercial_ready" {
		t.Fatalf("verification_status = %s, want blocked_not_commercial_ready", got)
	}
	if got := requireString(t, verifier, "verified_current_go_live_status"); got != "not_commercial_ready" {
		t.Fatalf("verified_current_go_live_status = %s, want not_commercial_ready", got)
	}
	for key, want := range map[string]string{
		"candidate_set_ref":                        "role-pack-candidate-set://team-office-v0",
		"source_execution_queue":                   "integration/commercial-cross-repo-execution-queue-candidate.json",
		"source_cross_repo_evidence_ledger":        "docs/commercial-cross-repo-evidence-ledger.json",
		"source_commercial_go_no_go_gate":          "tests/commercial-go-no-go-gate-candidate.json",
		"source_commercial_chain_verifier":         "tests/commercial-chain-verifier-candidate.json",
		"source_p11_go_live_evidence_package":      "tests/p11-commercial-go-live-evidence-package-template.json",
		"source_product_readiness_evidence_matrix": "tests/product-readiness-evidence-matrix.json",
		"next_required_authorization":              "owner_authorization_intake",
		"next_authorization_card":                  "integration/owner-authorization-evidence-intake-candidate.json",
	} {
		if got := requireString(t, verifier, key); got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}

	blockers := requireObject(t, verifier, "current_blockers")
	rawStageGates, ok := blockers["required_stage_gates"].([]any)
	if !ok {
		t.Fatalf("current_blockers.required_stage_gates missing")
	}
	expectedStages := map[string]string{
		"owner_authorization_intake":             "role_studio_auth_intake_recorded",
		"contracts_schema_review":                "role_studio_contracts_schema_impact",
		"backend_role_candidate_gate_receipt":    "role_studio_backend_candidate_receipts",
		"frontend_user_view_gui_flow":            "role_studio_gui_user_view_walkthrough",
		"cloud_sandbox_upload_purchase_download": "role_studio_cloud_sandbox_receipts",
		"local_install_team_binding":             "role_studio_install_team_binding_receipts",
		"negative_cases_and_observability":       "role_studio_negative_observability_receipts",
		"independent_acceptance_go_no_go":        "role_studio_independent_acceptance_go_no_go",
		"production_promotion_controls":          "role_studio_production_promotion_receipts",
	}
	if len(rawStageGates) != len(expectedStages) {
		t.Fatalf("required_stage_gates len = %d, want %d", len(rawStageGates), len(expectedStages))
	}
	seenStages := map[string]bool{}
	for _, raw := range rawStageGates {
		stageGate, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("stage gate = %T", raw)
		}
		stageID := requireString(t, stageGate, "stage_id")
		wantEvidenceID, ok := expectedStages[stageID]
		if !ok {
			t.Fatalf("unexpected stage_id %s", stageID)
		}
		seenStages[stageID] = true
		if got := requireString(t, stageGate, "evidence_id"); got != wantEvidenceID {
			t.Fatalf("%s evidence_id = %s, want %s", stageID, got, wantEvidenceID)
		}
		if got := requireString(t, stageGate, "authorization_status"); got != "missing" {
			t.Fatalf("%s authorization_status = %s, want missing", stageID, got)
		}
		if got := requireString(t, stageGate, "evidence_status"); got != "pending_authorization" {
			t.Fatalf("%s evidence_status = %s, want pending_authorization", stageID, got)
		}
		requireBool(t, stageGate, "can_count_toward_commercial_ready", false)
		if blockedBy := asStringSlice(t, stageGate["blocked_by"]); len(blockedBy) == 0 {
			t.Fatalf("%s blocked_by must not be empty", stageID)
		}
	}
	for stageID := range expectedStages {
		if !seenStages[stageID] {
			t.Fatalf("missing stage gate %s", stageID)
		}
	}

	writeback := requireObject(t, blockers, "evidence_writeback_summary")
	if got := requireString(t, writeback, "source_cross_repo_evidence_ledger"); got != "docs/commercial-cross-repo-evidence-ledger.json" {
		t.Fatalf("source_cross_repo_evidence_ledger = %s", got)
	}
	if got := requireString(t, writeback, "writeback_status"); got != "blocked_pending_authorization_and_evidence" {
		t.Fatalf("writeback_status = %s", got)
	}
	requireBool(t, writeback, "can_count_toward_commercial_ready", false)
	for key, want := range map[string]float64{
		"total_required_entry_count":  9,
		"total_pending_entry_count":   9,
		"total_completed_entry_count": 0,
	} {
		if got, ok := writeback[key].(float64); !ok || got != want {
			t.Fatalf("%s = %v, want %.0f", key, writeback[key], want)
		}
	}
	requiredBeforePass := strings.Join(asStringSlice(t, writeback["required_before_pass"]), "\n")
	for _, proof := range []string{"all_evidence_rows_complete_verified", "all_authoritative_receipts_present", "independent_acceptance_signed", "owner_go_no_go_decision_recorded", "production_promotion_receipts_verified"} {
		if !strings.Contains(requiredBeforePass, proof) {
			t.Fatalf("evidence_writeback_summary required_before_pass missing %s", proof)
		}
	}

	frontendBackend := requireObject(t, blockers, "frontend_backend_acceptance")
	if got := requireString(t, frontendBackend, "acceptance_status"); got != "blocked_pending_cross_repo_authorization_and_evidence" {
		t.Fatalf("frontend_backend_acceptance.acceptance_status = %s", got)
	}
	requireBool(t, frontendBackend, "can_count_toward_commercial_ready", false)
	repos := strings.Join(asStringSlice(t, frontendBackend["required_repositories"]), "\n")
	for _, repo := range []string{"truzhen-client-web-desktop", "truzhenos", "truzhen-cloud"} {
		if !strings.Contains(repos, repo) {
			t.Fatalf("frontend_backend_acceptance required_repositories missing %s", repo)
		}
	}

	rawTerminal, ok := verifier["terminal_checks"].([]any)
	if !ok {
		t.Fatalf("terminal_checks missing")
	}
	requiredTerminal := map[string]bool{
		"no_raw_asset_payloads":                       false,
		"no_real_payment_without_authorization":       false,
		"no_production_publish_without_authorization": false,
		"candidate_formal_boundary_preserved":         false,
		"bundle_tree_sha256_consistent":               false,
		"independent_acceptance_signed":               false,
		"owner_go_no_go_decision_recorded":            false,
	}
	for _, raw := range rawTerminal {
		terminal, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("terminal check = %T", raw)
		}
		gateID := requireString(t, terminal, "gate_id")
		if _, ok := requiredTerminal[gateID]; ok {
			requiredTerminal[gateID] = true
		}
		if got := requireString(t, terminal, "current_result"); got != "pending" {
			t.Fatalf("%s current_result = %s, want pending", gateID, got)
		}
		requireBool(t, terminal, "can_count_toward_commercial_ready", false)
	}
	for gateID, seen := range requiredTerminal {
		if !seen {
			t.Fatalf("missing terminal check %s", gateID)
		}
	}

	policy := requireObject(t, verifier, "completion_claim_policy")
	requireBool(t, policy, "completion_claim_allowed", false)
	requiredBeforeCompletion := strings.Join(asStringSlice(t, policy["required_before_completion_claim"]), "\n")
	for _, proof := range []string{"commercial_readiness_verifier_passed", "commercial_go_no_go_gate_passed", "all_evidence_writebacks_completed", "owner_go_no_go_decision_recorded"} {
		if !strings.Contains(requiredBeforeCompletion, proof) {
			t.Fatalf("completion_claim_policy required_before_completion_claim missing %s", proof)
		}
	}

	nonSufficient := strings.Join(asStringSlice(t, verifier["non_sufficient_evidence"]), "\n")
	for _, item := range []string{"candidate_assets_only", "green_pack_tests_without_gui_cloud_install_receipts", "manual_summary_without_receipts"} {
		if !strings.Contains(nonSufficient, item) {
			t.Fatalf("non_sufficient_evidence missing %s", item)
		}
	}

	candidateSet := readJSON(t, filepath.Join(base, "candidate-set.json"))
	if got := requireString(t, candidateSet, "commercial_readiness_verifier"); got != verifierPath {
		t.Fatalf("commercial_readiness_verifier = %s, want %s", got, verifierPath)
	}
	files := strings.Join(asStringSlice(t, candidateSet["artifact_files"]), "\n")
	if !strings.Contains(files, verifierPath) {
		t.Fatalf("candidate set artifact_files missing %s", verifierPath)
	}
	gate := readJSON(t, filepath.Join(base, "tests", "commercial-go-no-go-gate-candidate.json"))
	if got := requireString(t, gate, "source_commercial_readiness_verifier"); got != verifierPath {
		t.Fatalf("commercial go/no-go source_commercial_readiness_verifier = %s, want %s", got, verifierPath)
	}
	matrix := readJSON(t, filepath.Join(base, "tests", "product-readiness-evidence-matrix.json"))
	proofs := strings.Join(asStringSlice(t, requireObject(t, matrix, "completion_claim_policy")["required_before_completion_claim"]), "\n")
	if !strings.Contains(proofs, "commercial_readiness_verifier_passed") {
		t.Fatalf("product readiness matrix missing commercial_readiness_verifier_passed")
	}

	manifest := readJSON(t, filepath.Join(base, "commerce", "artifact-manifest.json"))
	rawFiles, ok := manifest["files"].([]any)
	if !ok {
		t.Fatalf("artifact manifest files missing")
	}
	verifierInManifest := false
	for _, raw := range rawFiles {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("manifest file item = %T", raw)
		}
		if requireString(t, item, "path") != verifierPath {
			continue
		}
		requireStringIn(t, requireString(t, item, "required_for"), "audit")
		data, err := os.ReadFile(filepath.Join(base, verifierPath))
		if err != nil {
			t.Fatalf("read %s: %v", verifierPath, err)
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(data))
		if got := requireString(t, item, "sha256"); got != wantHash {
			t.Fatalf("%s manifest sha256 = %s, want %s", verifierPath, got, wantHash)
		}
		verifierInManifest = true
	}
	if !verifierInManifest {
		t.Fatalf("artifact manifest missing %s", verifierPath)
	}
}
