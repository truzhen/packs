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
	{
		name:    "真实手机号不得写入 Pack 资产（业务 PII）",
		pattern: regexp.MustCompile(`:\s*"1[3-9]\d{9}"`),
	},
	{
		name:    "真实身份证号不得写入 Pack 资产（业务 PII）",
		pattern: regexp.MustCompile(`:\s*"\d{17}[\dXx]"`),
	},
}

// forbiddenDatabaseArtifact 判断文件是否为数据库禁品：数据库扩展名，或（防改扩展名绕过）
// 以 SQLite 魔数开头。业务运行态数据库不得进 Pack 资产。
func forbiddenDatabaseArtifact(ext string, head []byte) bool {
	switch strings.ToLower(ext) {
	case ".db", ".sqlite", ".sqlite3":
		return true
	}
	return strings.HasPrefix(string(head), "SQLite format 3")
}

// T5.1 扫描增强：真实手机号 / 身份证号（业务 PII）必须被禁品扫描命中。
func TestForbiddenPatternsCatchBusinessPII(t *testing.T) {
	cases := map[string]string{
		"手机号":  `{"contact":"13800138000"}`,
		"身份证号": `{"id":"11010119900307451X"}`,
	}
	for name, content := range cases {
		matched := false
		for _, f := range forbiddenPackAssetPatterns {
			if f.pattern.FindStringIndex(content) != nil {
				matched = true
				break
			}
		}
		if !matched {
			t.Fatalf("%s PII 必须被禁品扫描命中: %s", name, content)
		}
	}
}

// T5.1 扫描增强：数据库文件（扩展名或 SQLite 魔数）不得进 Pack 资产。
func TestForbiddenDatabaseArtifact(t *testing.T) {
	if !forbiddenDatabaseArtifact(".db", nil) || !forbiddenDatabaseArtifact(".sqlite", nil) || !forbiddenDatabaseArtifact(".SQLite3", nil) {
		t.Fatal("数据库扩展名必须被判为禁品")
	}
	if !forbiddenDatabaseArtifact(".bin", []byte("SQLite format 3\x00rest")) {
		t.Fatal("SQLite 魔数（改扩展名绕过）必须被判为禁品")
	}
	if forbiddenDatabaseArtifact(".json", []byte(`{"ok":true}`)) {
		t.Fatal("普通 json 不应误判为数据库禁品")
	}
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
				// 数据库文件（扩展名或 SQLite 魔数）不得进 Pack 资产（T5.1 扫描增强）。
				var head []byte
				if f, err := os.Open(path); err == nil {
					buf := make([]byte, 16)
					n, _ := f.Read(buf)
					f.Close()
					head = buf[:n]
				}
				if forbiddenDatabaseArtifact(filepath.Ext(path), head) {
					t.Fatalf("%s: 数据库文件不得写入 Pack 资产（业务运行态数据）", path)
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

func TestPackGlueDoesNotMintOwnerActionEvidence(t *testing.T) {
	installers, err := filepath.Glob("*-v0/install.py")
	if err != nil {
		t.Fatal(err)
	}
	uninstallers, err := filepath.Glob("*-v0/uninstall.py")
	if err != nil {
		t.Fatal(err)
	}
	paths := append(installers, uninstallers...)
	if len(paths) == 0 {
		t.Fatal("未发现 Pack glue，静态主权门不能空跑")
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		source := string(data)
		if strings.Contains(source, "owner_action_evidence://") {
			t.Fatalf("%s: Pack glue 禁止自铸 owner_action_evidence_ref", path)
		}
		if strings.HasSuffix(path, "uninstall.py") && !strings.Contains(source, `"action_type": "14.pack-studio.lifecycle.disable"`) {
			t.Fatalf("%s: 卸载缺少 os-14 canonical action_type", path)
		}
	}
}
