#!/bin/sh
# Verify core/ has no imports outside the standard library and the dreego
# module itself. Replaces the compile-time module-boundary enforcement that
# existed when core was a separate Go module.
set -e

cd "$(dirname "$0")/.."

deps=$(go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ./core/ 2>/dev/null | grep -v '^github.com/dreego-stack/dreego' | grep -v '^$' || true)

if [ -n "$deps" ]; then
	echo "FAIL: core/ imports external packages:"
	echo "$deps"
	exit 1
fi

echo "PASS: core/ has no external deps"
