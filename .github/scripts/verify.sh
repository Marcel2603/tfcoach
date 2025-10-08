#!/usr/bin/env bash
set -euo pipefail

# Ensure tools are available
command -v git >/dev/null || { echo "git not found"; exit 1; }
command v git-cliff >/dev/null 2>&1 || command -v git-cliff >/dev/null || { echo "git-cliff not found"; exit 1; }
command -v goreleaser >/dev/null || { echo "goreleaser not found"; exit 1; }

# Ensure token present for GitHub Release
: "${GITHUB_TOKEN:?GITHUB_TOKEN env var is required}"

echo "verify: ok"
