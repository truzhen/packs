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
