package packs_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var forbiddenPackAssetPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{
		name:    "正式 receipt ref 不得写入 Pack 资产",
		pattern: regexp.MustCompile(`(?i)"[^"]*receipt[^"]*"\s*:\s*"receipt://`),
	},
	{
		name:    "Base/Gate decision ref 不得由 Pack 资产自带",
		pattern: regexp.MustCompile(`(?i)"[^"]*decision[^"]*"\s*:\s*"(decision://|gated_owner_decision|base_decision://)`),
	},
	{
		name:    "正式 pack version ref 不得写入 Pack 资产",
		pattern: regexp.MustCompile(`(?i)"[^"]*pack_version_ref"\s*:\s*"(pack_version://|enabled_pack_version://)`),
	},
	{
		name:    "raw secret 字段不得写入 Pack 资产",
		pattern: regexp.MustCompile(`(?i)"[^"]*(password|api_key|private_key|access_token|refresh_token|cookie|payment_secret)[^"]*"\s*:\s*"[^"]+`),
	},
	{
		name:    "云端业务真相 ID 不得写入 Pack 资产",
		pattern: regexp.MustCompile(`(?i)"(buyer_id|order_id|entitlement_id|payment_id|license_id)"\s*:\s*"[^"]+`),
	},
}

func TestPackAssetsDoNotCarryBusinessDataFormalRefsOrRawSecrets(t *testing.T) {
	assetRoots := []string{
		"environmental-enforcement-pack-v0",
		"housekeeping-ops-pack-v0",
		"smart-home-owner-pack-v0",
		"templates/scene-pack-software-template",
	}

	for _, root := range assetRoots {
		root := root
		t.Run(root, func(t *testing.T) {
			if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if filepath.Ext(path) != ".json" {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				content := string(data)
				for _, forbidden := range forbiddenPackAssetPatterns {
					if loc := forbidden.pattern.FindStringIndex(content); loc != nil {
						start := loc[0]
						line := 1 + strings.Count(content[:start], "\n")
						t.Fatalf("%s: %s at line %d", path, forbidden.name, line)
					}
				}
				return nil
			}); err != nil {
				t.Fatalf("scan %s: %v", root, err)
			}
		})
	}
}
