#!/bin/sh
set -eu

repo_dir="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
temporary_dir=""
dreego_bin="${DREEGO_BIN:-}"

if [ -z "$dreego_bin" ]; then
	temporary_dir="$(mktemp -d)"
	trap 'rm -rf "$temporary_dir"' EXIT
	dreego_bin="$temporary_dir/dreego"
	(cd "$repo_dir" && go build -o "$dreego_bin" ./cmd/dreego)
fi

(cd "$repo_dir/demo/demo-wailsv3" && "$dreego_bin" generate && go test ./...)
