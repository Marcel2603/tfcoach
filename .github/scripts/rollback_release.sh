#!/usr/bin/env bash
set -euo pipefail

VERSION="$1"
REPO="$2"
TOKEN="$3"
TAG="v${VERSION}"

api() {
  curl -sS -H "Authorization: Bearer ${TOKEN}" -H "Accept: application/vnd.github+json" "$@"
}

echo "Rolling back release/tag ${TAG} on ${REPO} ..."

RELEASE_JSON="$(api "https://api.github.com/repos/${REPO}/releases/tags/${TAG}" || true)"
RELEASE_ID="$(echo "${RELEASE_JSON}" | jq -r '.id // empty')"

if [[ -n "${RELEASE_ID}" && "${RELEASE_ID}" != "null" ]]; then
  echo "Deleting GitHub release id ${RELEASE_ID}"
  api -X DELETE "https://api.github.com/repos/${REPO}/releases/${RELEASE_ID}" >/dev/null || true
fi

echo "Deleting remote tag ${TAG} (if exists)"
api -X DELETE "https://api.github.com/repos/${REPO}/git/refs/tags/${TAG}" >/dev/null || true

if git rev-parse "${TAG}" >/dev/null 2>&1; then
  git tag -d "${TAG}" || true
fi

echo "Rollback complete."
