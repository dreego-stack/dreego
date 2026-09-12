#!/bin/sh
# Verify the root implementation, core, and SSR adapter use only the standard library, the
# dreego module, and modules maintained by the Go project under golang.org/x/.
set -e

cd "$(dirname "$0")/.."

for pkg in ./internal/... ./core/... ./adapter/ssr/... ./adapter/wails/...; do
	deps=$(go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' "$pkg" 2>/dev/null | grep -v '^github.com/dreego-stack/dreego' | grep -v '^golang.org/x/' | grep -v '^$' || true)

	if [ -n "$deps" ]; then
		echo "FAIL: $pkg imports external packages:"
		echo "$deps"
		exit 1
	fi
done

echo "PASS: root internal, core, adapter/ssr, and adapter/wails use only approved dependencies"
