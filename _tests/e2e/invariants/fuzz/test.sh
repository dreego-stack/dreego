#!/bin/sh
# What: the committed fuzz seed corpus for the Markdown-to-HTML renderer stays
# safe. Long active fuzzing is a soak job, not a per-PR gate; this runs the
# deterministic seed corpus (including committed crashers).
set -eu

REPO_DIR="${REPO_DIR:-$(cd "$(dirname "$0")/../../../.." && pwd)}"
cd "$REPO_DIR"

go test ./core/ -run FuzzMarkdownToHTMLSafe -count=1
echo ok
