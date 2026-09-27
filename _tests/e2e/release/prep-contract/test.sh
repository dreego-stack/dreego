#!/bin/sh
# What: the release-prep script keeps its contract (tag verification, version
# file handling).
set -eu

REPO_DIR="${REPO_DIR:-$(cd "$(dirname "$0")/../../../.." && pwd)}"
cd "$REPO_DIR"

python3 _tests/e2e/release/prep-contract/release-prep-test.py
echo ok
