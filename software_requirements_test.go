package packs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSoftwareRequirementsDeclareBaserowConflictAndSharedOCR(t *testing.T) {
	smartHome := readJSON(t, filepath.Join("smart-home-owner-pack-v0", "manifest.json"))
	housekeeping := readJSON(t, filepath.Join("housekeeping-ops-pack-v0", "manifest.json"))

	requireSoftwareRequirement(t, smartHome, "baserow-runtime-a", "baserow-family", ">=2.2.0,<2.3.0", "reuse_preferred")
	requireSoftwareRequirement(t, housekeeping, "baserow-runtime-b", "baserow-family", ">=2.3.0,<3.0.0", "coexist_multi_version")
	requireSoftwareRequirement(t, smartHome, "shared-document-ocr", "pdf-ocr-family", ">=1.0.0,<2.0.0", "reuse_preferred")
	requireSoftwareRequirement(t, housekeeping, "shared-document-ocr", "pdf-ocr-family", ">=1.0.0,<2.0.0", "reuse_preferred")
}

func TestSmartHomeDeclaresFullProjectLifecycleAndOptionalHomeAssistantBoundary(t *testing.T) {
	manifest := readJSON(t, filepath.Join("smart-home-owner-pack-v0", "manifest.json"))
	if got := requireString(t, manifest, "version"); got != "1.1.0" {
		t.Fatalf("smart-home version = %s, want 1.1.0", got)
	}
	for _, capability := range []string{
		"project_opportunity_write_candidate",
		"project_initiation_write_candidate",
		"project_progress_write_candidate",
		"project_material_write_candidate",
		"project_delivery_write_candidate",
	} {
		if !containsString(manifest["capabilities_required"], capability) {
			t.Fatalf("smart-home full lifecycle capability missing: %s", capability)
		}
	}
	requireSoftwareRequirement(t, manifest, "home-assistant-runtime", "home-assistant-core", ">=2025.0.0,<2027.0.0", "reuse_preferred")

	providerRequirements, ok := manifest["provider_requirements"].([]any)
	if !ok {
		t.Fatal("provider_requirements missing")
	}
	foundOptionalHomeAssistant := false
	for _, raw := range providerRequirements {
		requirement, _ := raw.(map[string]any)
		if requirement["requirement_id"] != "home_assistant_device_control_candidate" {
			continue
		}
		foundOptionalHomeAssistant = requirement["optional"] == true &&
			requirement["fallback_policy"] == "not_ready" &&
			requirement["gateway_class"] == "execution"
	}
	if !foundOptionalHomeAssistant {
		t.Fatal("Home Assistant must remain an optional L2 ProviderRequirement with explicit not_ready fallback")
	}

	flow := readJSON(t, filepath.Join("smart-home-owner-pack-v0", "flows", "smart-home-owner-project-ops-flow.flow.json"))
	for _, nodeID := range []string{"opportunity_candidate", "project_initiation_candidate", "progress_candidate", "material_candidate", "delivery_candidate", "history_query"} {
		if !flowContainsNode(flow, nodeID) {
			t.Fatalf("smart-home flow missing project lifecycle node %s", nodeID)
		}
	}
}

func TestPackTreesDoNotContainSoftwareRuntimeArtifacts(t *testing.T) {
	for _, root := range []string{"smart-home-owner-pack-v0", "housekeeping-ops-pack-v0"} {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel := filepath.ToSlash(path)
			lower := strings.ToLower(rel)
			if strings.Contains(lower, "/runtime/") ||
				strings.Contains(lower, "/software/") ||
				strings.Contains(lower, "/models/") ||
				strings.Contains(lower, "/secrets/") ||
				strings.HasSuffix(lower, ".db") ||
				strings.HasSuffix(lower, ".sqlite") ||
				strings.HasSuffix(lower, ".tokens") ||
				strings.HasSuffix(lower, ".pem") {
				t.Fatalf("%s contains forbidden software/runtime artifact: %s", root, rel)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
}

func requireSoftwareRequirement(t *testing.T, manifest map[string]any, id, family, versionRange, isolationPolicy string) {
	t.Helper()
	raw, ok := manifest["software_requirements"].([]any)
	if !ok {
		t.Fatalf("software_requirements missing")
	}
	for _, item := range raw {
		req, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if requireString(t, req, "requirement_id") != id {
			continue
		}
		if got := requireString(t, req, "software_family"); got != family {
			t.Fatalf("%s software_family = %s, want %s", id, got, family)
		}
		if got := requireString(t, req, "version_range"); got != versionRange {
			t.Fatalf("%s version_range = %s, want %s", id, got, versionRange)
		}
		if got := requireString(t, req, "isolation_policy"); got != isolationPolicy {
			t.Fatalf("%s isolation_policy = %s, want %s", id, got, isolationPolicy)
		}
		return
	}
	t.Fatalf("software requirement %s missing", id)
}

func containsString(raw any, want string) bool {
	items, _ := raw.([]any)
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func flowContainsNode(flow map[string]any, want string) bool {
	nodes, _ := flow["nodes"].([]any)
	for _, raw := range nodes {
		node, _ := raw.(map[string]any)
		if node["id"] == want {
			return true
		}
	}
	return false
}
