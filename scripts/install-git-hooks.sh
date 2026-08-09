#!/usr/bin/env bash
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
test -x scripts/hooks/pre-commit && test -x scripts/hooks/pre-push && test -x scripts/hygiene-gate.sh || { echo "[install-git-hooks] mandatory hygiene scripts are missing" >&2; exit 1; }
toolchain_dir="${TRUZHEN_G0H_TOOLCHAIN_DIR:-}"
test -n "$toolchain_dir" && test -d "$toolchain_dir" || { echo "[install-git-hooks] TRUZHEN_G0H_TOOLCHAIN_DIR is required" >&2; exit 1; }
git config core.hooksPath scripts/hooks
git config truzhen.g0hToolchainDir "$toolchain_dir"
test "$(git config --get core.hooksPath)" = "scripts/hooks"
test "$(git config --get truzhen.g0hToolchainDir)" = "$toolchain_dir"
bash scripts/hygiene-gate.sh doctor
echo "[install-git-hooks] PASS" >&2
