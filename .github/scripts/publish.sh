#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:?version missing}"
BRANCH="${2:-}"
COMMITS_LEN="${3:-0}"
NOW_MS="${4:-}"

# Make sure local installs are on PATH (Go & Cargo bins)
export PATH="$HOME/go/bin:$HOME/.cargo/bin:$PATH"

echo "Publishing v${VERSION} on branch ${BRANCH} (${COMMITS_LEN} commits; ts=${NOW_MS})"

# Generate incremental notes for this release and pass to GoReleaser
tmp_notes="$(mktemp)"
git-cliff --current --tag "v${VERSION}" > "${tmp_notes}"

# Build & publish via GoReleaser using those notes
goreleaser release --clean --release-notes "${tmp_notes}"
