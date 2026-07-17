package packs_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const contentOperationsPackRoot = "content-operations-workbench-v0"

func readContentOpsJSON(t *testing.T, rel string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode %s: %v", rel, err)
	}
	return value
}

func stringSlice(t *testing.T, value any, field string) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("%s must be array, got %T", field, value)
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok {
			t.Fatalf("%s item must be string, got %T", field, item)
		}
		result = append(result, text)
	}
	return result
}

func contentOpsContainsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestContentOperationsPackDeclaresSeparatedBoundaries(t *testing.T) {
	manifest := readContentOpsJSON(t, "manifest.json")
	if got := manifest["pack_ref"]; got != "scene_pack://content-operations-workbench" {
		t.Fatalf("pack_ref = %v", got)
	}
	if got := manifest["pack_type"]; got != "domain_work_pack" {
		t.Fatalf("pack_type = %v", got)
	}
	security, ok := manifest["security_profile"].(map[string]any)
	if !ok || security["candidate_only"] != true || security["no_real_publish"] != true || security["no_platform_login"] != true {
		t.Fatalf("security_profile must be candidate-only/no-publish/no-login: %#v", manifest["security_profile"])
	}

	workModes, ok := manifest["work_modes"].([]any)
	if !ok || len(workModes) != 3 {
		t.Fatalf("work_modes = %#v", manifest["work_modes"])
	}
	wants := map[string]bool{"direction_radar": false, "content_production": false, "weekly_review": false}
	for _, raw := range workModes {
		mode, _ := raw.(map[string]any)
		if _, exists := wants[fmt.Sprint(mode["mode_id"])]; exists {
			wants[fmt.Sprint(mode["mode_id"])] = true
		}
	}
	for mode, present := range wants {
		if !present {
			t.Fatalf("missing work mode %s", mode)
		}
	}

	providerRequirements, ok := manifest["provider_requirements"].([]any)
	if !ok || len(providerRequirements) != 1 {
		t.Fatalf("provider_requirements = %#v", manifest["provider_requirements"])
	}
	requirement, _ := providerRequirements[0].(map[string]any)
	if requirement["provider_family"] != "codex-cli-family" || requirement["gateway_class"] != "execution" {
		t.Fatalf("provider must remain external Codex Hands via execution gateway: %#v", requirement)
	}
	moat, ok := manifest["moat_justification"].(map[string]any)
	if !ok {
		t.Fatalf("moat_justification = %#v", manifest["moat_justification"])
	}
	allowedMoatReasons := map[string]bool{
		"real_system_integration":      true,
		"domain_approval_compliance":   true,
		"auditable_evidence_receipt":   true,
		"multi_role_compare_sovereign": true,
		"transaction_lifecycle_replay": true,
	}
	for _, reason := range stringSlice(t, moat["valid_reasons"], "moat_justification.valid_reasons") {
		if !allowedMoatReasons[reason] {
			t.Fatalf("unsupported Pack Studio moat reason %q", reason)
		}
	}
}

func TestContentOperationsFlowCannotPublish(t *testing.T) {
	flow := readContentOpsJSON(t, "flows/content-operations.flow.json")
	allowedTypes := map[string]bool{
		"start": true, "business_object": true, "advice": true, "task": true,
		"challenge": true, "compare_gate": true, "human_approval": true,
		"receipt": true, "end": true,
	}
	nodes, ok := flow["nodes"].([]any)
	if !ok {
		t.Fatalf("nodes = %#v", flow["nodes"])
	}
	for _, raw := range nodes {
		node, _ := raw.(map[string]any)
		nodeType := fmt.Sprint(node["type"])
		if !allowedTypes[nodeType] {
			t.Fatalf("flow contains forbidden or unknown node type %q: %#v", nodeType, node)
		}
		joined := strings.ToLower(fmt.Sprint(node))
		for _, forbidden := range []string{"platform_login", "social_platform_upload", "social_platform_publish", "contact_scraping"} {
			if strings.Contains(joined, forbidden) {
				t.Fatalf("flow node contains forbidden action %q: %#v", forbidden, node)
			}
		}
	}
	edges, ok := flow["edges"].([]any)
	if !ok {
		t.Fatalf("edges = %#v", flow["edges"])
	}
	for _, raw := range edges {
		edge, _ := raw.(map[string]any)
		condition := fmt.Sprint(edge["condition"])
		if strings.Contains(condition, "mode=") {
			t.Fatalf("edge %v uses unsupported condition token %q", edge["id"], condition)
		}
	}
	forbiddenNodes := stringSlice(t, flow["forbidden_nodes"], "forbidden_nodes")
	if !contentOpsContainsString(forbiddenNodes, "social_platform_publish") || !contentOpsContainsString(forbiddenNodes, "social_platform_login") {
		t.Fatalf("flow must explicitly reject publish/login: %v", forbiddenNodes)
	}
}

func TestContentOperationsSkillBundleIntegrityAndCandidateContract(t *testing.T) {
	bundle := readContentOpsJSON(t, "skills/truzhen-content-ops/bundle.json")
	if bundle["provider_family"] != "codex-cli-family" || bundle["runtime_implementation_in_pack"] != false {
		t.Fatalf("bundle runtime boundary = %#v", bundle)
	}
	integrity, ok := bundle["content_integrity"].(map[string]any)
	if !ok || integrity["algorithm"] != "sha256-framed-files-v1" {
		t.Fatalf("content_integrity = %#v", bundle["content_integrity"])
	}
	hasher := sha256.New()
	files := stringSlice(t, integrity["files"], "content_integrity.files")
	if !contentOpsContainsString(files, "model-output.schema.json") || bundle["model_output_schema"] != "model-output.schema.json" {
		t.Fatalf("Pack must own and integrity-cover its model output schema: files=%v schema=%v", files, bundle["model_output_schema"])
	}
	for _, name := range files {
		data, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, "skills/truzhen-content-ops", name))
		if err != nil {
			t.Fatalf("read integrity file %s: %v", name, err)
		}
		_, _ = hasher.Write([]byte(name))
		_, _ = hasher.Write([]byte{0})
		_, _ = hasher.Write([]byte(fmt.Sprint(len(data))))
		_, _ = hasher.Write([]byte{0})
		_, _ = hasher.Write(data)
	}
	if got, want := hex.EncodeToString(hasher.Sum(nil)), fmt.Sprint(integrity["sha256"]); got != want {
		t.Fatalf("skill bundle hash = %s, want %s", got, want)
	}
	skillBytes, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, "skills/truzhen-content-ops/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	skillText := string(skillBytes)
	if !strings.HasPrefix(skillText, "---\nname: truzhen-content-ops\ndescription: ") {
		t.Fatal("SKILL.md must start with name/description-only YAML frontmatter")
	}
	frontmatterEnd := strings.Index(strings.TrimPrefix(skillText, "---\n"), "\n---\n")
	if frontmatterEnd < 0 {
		t.Fatal("SKILL.md frontmatter terminator missing")
	}
	frontmatter := strings.TrimPrefix(skillText, "---\n")[:frontmatterEnd]
	for _, forbidden := range []string{"metadata:", "version:", "tools:"} {
		if strings.Contains(frontmatter, forbidden) {
			t.Fatalf("SKILL.md frontmatter contains unsupported field %q", forbidden)
		}
	}

	contract := readContentOpsJSON(t, "skills/truzhen-content-ops/output-contract.json")
	if contract["candidate_only"] != true {
		t.Fatalf("candidate_only = %v", contract["candidate_only"])
	}
	if outputs := contract["formal_outputs"].([]any); len(outputs) != 0 {
		t.Fatalf("formal_outputs must be empty: %#v", outputs)
	}
	if effects := contract["external_side_effects"].([]any); len(effects) != 0 {
		t.Fatalf("external_side_effects must be empty: %#v", effects)
	}
	modelSchema := readContentOpsJSON(t, "skills/truzhen-content-ops/model-output.schema.json")
	if modelSchema["type"] != "object" {
		t.Fatalf("model output schema must have an object root: %#v", modelSchema)
	}
	publication, _ := contract["publication"].(map[string]any)
	if publication["supported"] != false || publication["platform_login"] == true || publication["publish"] == true {
		t.Fatalf("publication boundary = %#v", publication)
	}
}

func TestContentOperationsSchedulesStayLowRiskAndDoNotAutoProduce(t *testing.T) {
	doc := readContentOpsJSON(t, "schedules/content-operations.schedules.json")
	if doc["scheduler_timezone"] != "UTC" || doc["timezone_display"] != "Asia/Taipei" {
		t.Fatalf("schedule timezone mapping must be explicit: %#v", doc)
	}
	rawSchedules, ok := doc["schedules"].([]any)
	if !ok || len(rawSchedules) != 2 {
		t.Fatalf("schedules = %#v", doc["schedules"])
	}
	wants := map[string]string{
		"content_ops.direction_radar": "30 0 * * 1-5",
		"content_ops.weekly_review":   "30 12 * * 0",
	}
	for _, raw := range rawSchedules {
		schedule, _ := raw.(map[string]any)
		skillID := fmt.Sprint(schedule["skill_id"])
		if wants[skillID] != schedule["cron"] {
			t.Fatalf("unexpected schedule %s: %#v", skillID, schedule)
		}
		if schedule["template_risk_level"] != "low" || schedule["risk_ceiling"] != "low" || schedule["candidate_only"] != true {
			t.Fatalf("schedule must stay low-risk candidate-only: %#v", schedule)
		}
		if skillID == "content_ops.content_production" {
			t.Fatal("content production must wait for Owner direction and cannot have a standing schedule")
		}
	}
	notScheduled, ok := doc["not_scheduled"].([]any)
	if !ok || len(notScheduled) != 1 || notScheduled[0].(map[string]any)["skill_id"] != "content_ops.content_production" {
		t.Fatalf("content production exclusion missing: %#v", doc["not_scheduled"])
	}
}

func TestContentOperationsPackHasNoProviderRuntimeOrClientMintedProof(t *testing.T) {
	for _, forbiddenDir := range []string{"provider", "sidecar", "runtime", "scripts", "bin"} {
		path := filepath.Join(contentOperationsPackRoot, forbiddenDir)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			t.Fatalf("provider/runtime implementation must not exist in Pack: %s", path)
		}
	}
	for _, rel := range []string{"install.py", "uninstall.py"} {
		data, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, rel))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		for _, forbidden := range []string{
			"owner_action_evidence://",
			"decision://content-operations",
			"receipt://content-operations",
			"subprocess.",
			"os.system(",
		} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s contains client-minted proof or runtime execution marker %q", rel, forbidden)
			}
		}
	}
	uninstall, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, "uninstall.py"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(uninstall), "issue[\"owner_action_evidence_ref\"]") {
		t.Fatal("uninstall must consume Base-issued owner_action_evidence_ref")
	}
	install, err := os.ReadFile(filepath.Join(contentOperationsPackRoot, "install.py"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"agent-slots/readmodel",
		"enabled_bindings",
		"binding_key in enabled_bindings",
		"lifecycle_occ",
		"content-ops-pack-reactivate:%s:%d",
		"disabled 版本 reactivate 失败",
	} {
		if !strings.Contains(string(install), required) {
			t.Fatalf("install must prevent duplicate enabled slot bindings; missing %q", required)
		}
	}
	for _, required := range []string{
		"/v3/task-governance/candidates/intake",
		"/v3/task-governance/candidates/submit-review",
		"/v3/task-governance/schedules/approve",
		"/v3/task-governance/schedules/resume",
	} {
		if !strings.Contains(string(install), required) {
			t.Fatalf("install must use governed 07 schedule lifecycle; missing %q", required)
		}
	}
	if !strings.Contains(string(uninstall), "/v3/task-governance/schedules/pause") {
		t.Fatal("uninstall must pause Pack-owned schedules before disabling the Pack")
	}
}
