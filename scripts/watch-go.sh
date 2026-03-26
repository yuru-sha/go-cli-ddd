#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

INTERVAL="${WATCH_INTERVAL:-2}"
TMP_FILE="${TMPDIR:-/tmp}/go-cli-ddd-watch-go.state"

snapshot_files() {
  find . \
    -path './.git' -prune -o \
    -path './bin' -prune -o \
    -path './vendor' -prune -o \
    -name '*.go' -type f -print0 |
    sort -z |
    xargs -0 stat -f '%m %N' 2>/dev/null
}

run_checks() {
  echo "[watch-go] formatting Go files"
  gofmt -w $(find . -path './.git' -prune -o -path './bin' -prune -o -path './vendor' -prune -o -name '*.go' -type f -print)
  goimports -w $(find . -path './.git' -prune -o -path './bin' -prune -o -path './vendor' -prune -o -name '*.go' -type f -print)

  echo "[watch-go] running golangci-lint"
  golangci-lint run
  echo "[watch-go] checks passed"
}

echo "[watch-go] watching .go files under $ROOT_DIR"
snapshot_files > "$TMP_FILE"

while true; do
  sleep "$INTERVAL"

  NEXT_FILE="${TMP_FILE}.next"
  snapshot_files > "$NEXT_FILE"

  if ! cmp -s "$TMP_FILE" "$NEXT_FILE"; then
    mv "$NEXT_FILE" "$TMP_FILE"
    run_checks || true
    snapshot_files > "$TMP_FILE"
    continue
  fi

  rm -f "$NEXT_FILE"
done
