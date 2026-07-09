package packs_test

import (
	"os"
	"regexp"
	"testing"
)

// Y4：Pack 生命周期错误码必须按失败阶段细分，不得所有失败都返回单一通用码。
// 契约形状对齐 truzhen-contracts monitoring error_code pattern（TZ-<仓>-<域>-<NNN>）。
// 本 guard 是回归锁：若有人把 install/uninstall 退回单一通用码，或新增格式非法 /
// 重复的码，或删掉某阶段的细分接线，本测试即红。

// pack_diagnostics.py 里的常量定义形如 `NAME = "TZ-PACK-INSTALL-002"`。
var codeConstAssign = regexp.MustCompile(`(?m)^([A-Z][A-Z0-9_]+)\s*=\s*"(TZ-PACK-[A-Z0-9]{2,10}-\d{3})"`)

// 每个 pack 的 install.py 必须实际接线的阶段常量名——覆盖到这些即证明已从单一通用码
// 升级为阶段化归因（连通性 / readiness / 版本冲突 / 角色绑定 / 知识入库 / Base 门）。
var requiredInstallStageConsts = []string{
	"INSTALL_CONNECTIVITY",
	"INSTALL_READINESS",
	"INSTALL_STATE_CONFLICT",
	"INSTALL_ROLE_BINDING",
	"INSTALL_KNOWLEDGE",
	"INSTALL_BASE_GATE",
}

var requiredUninstallStageConsts = []string{
	"UNINSTALL_CONNECTIVITY",
	"UNINSTALL_LIFECYCLE_HTTP",
}

var packRoots = []string{
	"environmental-enforcement-pack-v0",
	"housekeeping-ops-pack-v0",
	"smart-home-owner-pack-v0",
	"shuxuejia-renovation-pack-v0",
}

var installShape = regexp.MustCompile(`^TZ-PACK-INSTALL-\d{3}$`)
var uninstallShape = regexp.MustCompile(`^TZ-PACK-UNINSTALL-\d{3}$`)

// diagnosticsCodes 解析共享诊断模块，返回 常量名->码 映射。
func diagnosticsCodes(t *testing.T) map[string]string {
	t.Helper()
	data, err := os.ReadFile("pack_diagnostics.py")
	if err != nil {
		t.Fatalf("读取 pack_diagnostics.py: %v", err)
	}
	byName := map[string]string{}
	for _, m := range codeConstAssign.FindAllStringSubmatch(string(data), -1) {
		byName[m[1]] = m[2]
	}
	return byName
}

func TestPackErrorCodeTaxonomyIsSubdivided(t *testing.T) {
	byName := diagnosticsCodes(t)

	// 1. 登记簿细分到位 + 码唯一 + 形状合法：install ≥ 8、uninstall ≥ 3。
	seen := map[string]string{}
	var installN, uninstallN int
	for name, code := range byName {
		if prev, dup := seen[code]; dup {
			t.Errorf("码 %q 重复：%s 与 %s", code, prev, name)
		}
		seen[code] = name
		switch {
		case installShape.MatchString(code):
			installN++
		case uninstallShape.MatchString(code):
			uninstallN++
		}
	}
	if installN < 8 {
		t.Errorf("install 阶段码只有 %d 个，未细分（需 ≥8：连通性/生命周期/readiness/版本冲突/角色/知识/Base门）", installN)
	}
	if uninstallN < 3 {
		t.Errorf("uninstall 阶段码只有 %d 个，未细分（需 ≥3）", uninstallN)
	}

	// 2. 全部 required 常量都在登记簿里定义。
	for _, c := range append(append([]string{}, requiredInstallStageConsts...), requiredUninstallStageConsts...) {
		if _, ok := byName[c]; !ok {
			t.Errorf("pack_diagnostics.py 未定义阶段常量 %s", c)
		}
	}

	// 3. 每个 pack 的 install/uninstall 必须实际接线阶段常量（回归锁：防退回单一通用码），
	//    且不得残留裸硬编码的通用码字面串（证明已改走命名常量）。
	for _, root := range packRoots {
		root := root
		t.Run(root, func(t *testing.T) {
			install, err := os.ReadFile(root + "/install.py")
			if err != nil {
				t.Fatalf("读取 %s/install.py: %v", root, err)
			}
			ic := string(install)
			for _, c := range requiredInstallStageConsts {
				if !regexp.MustCompile(`\b` + c + `\b`).MatchString(ic) {
					t.Errorf("%s/install.py 未接线阶段常量 %s（仍在用单一通用码？）", root, c)
				}
			}
			if regexp.MustCompile(`"TZ-PACK-INSTALL-\d{3}"`).MatchString(ic) {
				t.Errorf("%s/install.py 残留裸码字面串，应改走 pack_diagnostics 命名常量", root)
			}

			uninstall, err := os.ReadFile(root + "/uninstall.py")
			if err != nil {
				t.Fatalf("读取 %s/uninstall.py: %v", root, err)
			}
			uc := string(uninstall)
			for _, c := range requiredUninstallStageConsts {
				if !regexp.MustCompile(`\b` + c + `\b`).MatchString(uc) {
					t.Errorf("%s/uninstall.py 未接线阶段常量 %s", root, c)
				}
			}
			if regexp.MustCompile(`"TZ-PACK-UNINSTALL-\d{3}"`).MatchString(uc) {
				t.Errorf("%s/uninstall.py 残留裸码字面串，应改走命名常量", root)
			}
		})
	}
}
