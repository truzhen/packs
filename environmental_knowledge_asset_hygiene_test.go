package packs_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const environmentalPackDir = "environmental-enforcement-pack-v0"

type environmentalManifestForHygiene struct {
	KnowledgeScopes []string `json:"knowledge_scopes"`
}

type environmentalScopesForHygiene struct {
	Scopes []struct {
		ScopeRef string `json:"scope_ref"`
	} `json:"scopes"`
}

type environmentalIndexForHygiene struct {
	Count   int `json:"count"`
	Entries []struct {
		File               string `json:"file"`
		Kind               string `json:"kind"`
		KnowledgeScopeRef  string `json:"knowledge_scope_ref"`
		Authority          string `json:"authority"`
		VerificationStatus string `json:"verification_status"`
		Checksum           string `json:"checksum"`
	} `json:"entries"`
}

func readEnvironmentalJSONForHygiene(t *testing.T, relative string, target any) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(environmentalPackDir, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatalf("读取 %s: %v", relative, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("解析 %s: %v", relative, err)
	}
}

func TestEnvironmentalKnowledgeAssetHygiene(t *testing.T) {
	var manifest environmentalManifestForHygiene
	var scopes environmentalScopesForHygiene
	var index environmentalIndexForHygiene
	readEnvironmentalJSONForHygiene(t, "manifest.json", &manifest)
	readEnvironmentalJSONForHygiene(t, "knowledge/knowledge-scopes.json", &scopes)
	readEnvironmentalJSONForHygiene(t, "knowledge/knowledge-index.json", &index)

	if got := len(manifest.KnowledgeScopes); got != 15 {
		t.Errorf("manifest knowledge scopes = %d，期望 15", got)
	}
	if got := len(scopes.Scopes); got != 15 {
		t.Errorf("knowledge-scopes scopes = %d，期望 15", got)
	}
	manifestScopes := make(map[string]struct{}, len(manifest.KnowledgeScopes))
	for _, ref := range manifest.KnowledgeScopes {
		manifestScopes[ref] = struct{}{}
	}
	for _, scope := range scopes.Scopes {
		if _, ok := manifestScopes[scope.ScopeRef]; !ok {
			t.Errorf("knowledge-scopes 出现 manifest 未声明 scope: %s", scope.ScopeRef)
		}
		delete(manifestScopes, scope.ScopeRef)
	}
	for ref := range manifestScopes {
		t.Errorf("manifest scope 未出现在 knowledge-scopes: %s", ref)
	}

	if index.Count != 45 || len(index.Entries) != 45 {
		t.Errorf("knowledge index count/entries = %d/%d，期望 45/45", index.Count, len(index.Entries))
	}
	caseEntries := 0
	indexedFiles := make(map[string]struct{}, len(index.Entries))
	for _, entry := range index.Entries {
		if entry.Kind == "case" {
			caseEntries++
		}
		if entry.VerificationStatus != "pending_human_review" {
			t.Errorf("%s verification_status = %q", entry.File, entry.VerificationStatus)
		}
		if entry.Authority != "reference_only" {
			t.Errorf("%s authority = %q", entry.File, entry.Authority)
		}
		indexedFiles[filepath.Clean(filepath.FromSlash(entry.File))] = struct{}{}
		data, err := os.ReadFile(filepath.Join(environmentalPackDir, filepath.FromSlash(entry.File)))
		if err != nil {
			t.Errorf("index 文件不存在 %s: %v", entry.File, err)
			continue
		}
		sum := sha256.Sum256(data)
		want := "sha256:" + hex.EncodeToString(sum[:])
		if entry.Checksum != want {
			t.Errorf("%s checksum 漂移：got %s want %s", entry.File, entry.Checksum, want)
		}
	}
	if caseEntries != 30 {
		t.Errorf("index case entries = %d，期望 30", caseEntries)
	}

	actualCases, err := filepath.Glob(filepath.Join(environmentalPackDir, "knowledge", "cases", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got := len(actualCases); got != 30 {
		t.Errorf("实际 cases Markdown = %d，期望 30", got)
	}
	readme, err := os.ReadFile(filepath.Join(environmentalPackDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readme), "cases/                     # 30 个真实案例") {
		t.Errorf("README 未按实际资产声明 30 个真实案例")
	}
	if strings.Contains(string(readme), "31 个真实案例") {
		t.Errorf("README 仍残留 31 个真实案例漂移")
	}

	linkPattern := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	indexDocuments := []struct {
		Relative string
		Links    int
	}{
		{Relative: "knowledge/index/法条速查索引.md", Links: 66},
		{Relative: "knowledge/index/违法行为分类索引.md", Links: 50},
	}
	totalLinks := 0
	brokenLinks := 0
	oldDirectoryLinks := 0
	for _, document := range indexDocuments {
		path := filepath.Join(environmentalPackDir, filepath.FromSlash(document.Relative))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		matches := linkPattern.FindAllStringSubmatch(string(data), -1)
		if len(matches) != document.Links {
			t.Errorf("%s links = %d，期望 %d", document.Relative, len(matches), document.Links)
		}
		totalLinks += len(matches)
		for _, match := range matches {
			target := match[1]
			if strings.Contains(target, "../02_执法指南/") || strings.Contains(target, "../03_案例库/") {
				oldDirectoryLinks++
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if _, err := os.Stat(resolved); err != nil {
				brokenLinks++
			}
		}
	}
	if totalLinks != 116 {
		t.Errorf("两个索引 links = %d，期望 116", totalLinks)
	}
	if oldDirectoryLinks != 0 {
		t.Errorf("两个索引仍有旧目录链接 %d 个", oldDirectoryLinks)
	}
	if brokenLinks != 0 {
		t.Errorf("两个索引仍有失效链接 %d 个", brokenLinks)
	}

	indexedMarkdown := 0
	err = filepath.WalkDir(filepath.Join(environmentalPackDir, "knowledge"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		relative, err := filepath.Rel(environmentalPackDir, path)
		if err != nil {
			return err
		}
		indexedMarkdown++
		if _, ok := indexedFiles[filepath.Clean(relative)]; !ok {
			t.Errorf("未登记 Markdown: %s", relative)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if indexedMarkdown != 45 {
		t.Errorf("knowledge Markdown = %d，期望 45", indexedMarkdown)
	}
}
