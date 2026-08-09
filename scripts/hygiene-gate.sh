#!/usr/bin/env bash
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
mode="${1:-}"
[ -n "$mode" ] || { echo "usage: $0 {doctor|staged|range|outgoing|self-test|canary} ..." >&2; exit 2; }
shift || true
toolchain_dir="${TRUZHEN_G0H_TOOLCHAIN_DIR:-$(git config --get truzhen.g0hToolchainDir || true)}"
[ -n "$toolchain_dir" ] && [ -d "$toolchain_dir" ] || { echo "[hygiene-gate] FAIL: toolchain directory is required" >&2; exit 2; }
lock_file="$repo_root/security/hygiene-tool.lock.json"
allowlist_file="$repo_root/security/hygiene-allowlist.json"
command -v jq >/dev/null 2>&1 && [ -f "$lock_file" ] && [ -f "$allowlist_file" ] || { echo "[hygiene-gate] FAIL: jq, lock and allowlist are mandatory" >&2; exit 2; }
scanner_file="$toolchain_dir/$(jq -r '.scanner.filename // empty' "$lock_file")"
rules_file="$toolchain_dir/$(jq -r '.rules.filename // empty' "$lock_file")"
gitleaks_bin="${TRUZHEN_G0H_GITLEAKS_BIN:-$(command -v gitleaks || true)}"
check_digest() { local label="$1" file="$2" expected="$3" actual; [ -n "$expected" ] && [ -f "$file" ] || { echo "[hygiene-gate] FAIL: missing $label" >&2; exit 2; }; actual="$(shasum -a 256 "$file" | awk '{print $1}')"; [ "$actual" = "$expected" ] || { echo "[hygiene-gate] FAIL: $label digest mismatch" >&2; exit 2; }; }
check_digest scanner "$scanner_file" "$(jq -r '.scanner.sha256 // empty' "$lock_file")"
check_digest rules "$rules_file" "$(jq -r '.rules.sha256 // empty' "$lock_file")"
check_digest allowlist "$allowlist_file" "$(jq -r '.allowlist.sha256 // empty' "$lock_file")"
check_digest gitleaks "$gitleaks_bin" "$(jq -r '.gitleaks.sha256 // empty' "$lock_file")"
[ "$($scanner_file --version)" = "truzhen-hygiene-scan $(jq -r '.scanner.version // empty' "$lock_file")" ] || { echo "[hygiene-gate] FAIL: scanner version mismatch" >&2; exit 2; }
[ "$($gitleaks_bin version)" = "$(jq -r '.gitleaks.version // empty' "$lock_file")" ] || { echo "[hygiene-gate] FAIL: gitleaks version mismatch" >&2; exit 2; }
scanner_args=(--repo "$repo_root" --rules "$rules_file" --allowlist "$allowlist_file" --gitleaks "$gitleaks_bin")
case "$mode" in doctor|staged|self-test) exec "$scanner_file" "${scanner_args[@]}" "$mode" ;; range) exec "$scanner_file" "${scanner_args[@]}" range "$@" ;; outgoing) exec "$scanner_file" "${scanner_args[@]}" outgoing "$@" ;; canary) exec "$scanner_file" "${scanner_args[@]}" canary "$@" ;; *) echo "[hygiene-gate] FAIL: unsupported mode=$mode" >&2; exit 2 ;; esac
