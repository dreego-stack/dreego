#!/bin/sh
# Verify core/ and internal/transpiler/ have no imports outside the standard
# library and the dreego module itself. Replaces the compile-time
# module-boundary enforcement that existed when core was a separate Go module.
set -e

cd "$(dirname "$0")/.."

for pkg in ./core/ ./internal/transpiler/; do
	deps=$(go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' "$pkg" 2>/dev/null | grep -v '^github.com/dreego-stack/dreego' | grep -v '^$' || true)

	if [ -n "$deps" ]; then
		echo "FAIL: $pkg imports external packages:"
		echo "$deps"
		exit 1
	fi
done

echo "PASS: core/ and internal/transpiler/ have no external deps"
