#!/bin/sh
# Mirror the repo's documentation into cli/dreego/embedded/ so the CLI can
# embed it via //go:embed all:embedded. Run after any change to _docs/,
# README.md, or CHANGELOG.md, and before committing.
# Usage: _scripts/sync-embedded-docs.sh
set -e

cd "$(dirname "$0")/.."

dest="cli/dreego/embedded"

rm -rf "$dest"
mkdir -p "$dest"

cp -R _docs "$dest/_docs"
cp README.md "$dest/README.md"
cp CHANGELOG.md "$dest/CHANGELOG.md"

echo "synced docs into $dest"
