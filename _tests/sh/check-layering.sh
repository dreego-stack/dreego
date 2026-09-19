#!/bin/sh
# Enforce the internal layering and dreefile dependency rule over non-test files.
# Rule: _docs/decisions/internal-layering-and-dreefile.md
set -eu

cd "$(dirname "$0")/../.."

module="github.com/dreego-stack/dreego/"

if ! imports="$(go list -e -f '@{{.ImportPath}}{{range .Imports}}{{"\n"}}{{.}}{{end}}' ./internal/... 2>/dev/null)"; then
	echo "FAIL: could not list internal packages"
	exit 1
fi

violations="$(printf '%s\n' "$imports" | awk -v module="$module" '
function in_list(x, list,   n, a, i) {
  n = split(list, a, " ")
  for (i = 1; i <= n; i++) if (a[i] == x) return 1
  return 0
}
function allowed(pkg, imp,   g) {
  g = "internal/dreefile/codegen internal/dreefile/dreecode internal/dreefile/gogen internal/dreefile/ir internal/dreefile/jsoutput"
  if (pkg == "internal/dreefile")
    return in_list(imp, "internal/dreefile/sections/head internal/dreefile/sections/style internal/dreefile/sections/body/html internal/dreefile/sections/body/md internal/dreefile/sections/client internal/dreefile/sections/client/js internal/dreefile/sections/client/ts internal/dreefile/sections/client/lua internal/dreefile/tokens internal/dreefile/lexer internal/dreefile/parser internal/dreefile/i18n internal/gomod " g)
  if (pkg == "internal/dreefile/ir" || pkg == "internal/dreefile/tokens" || pkg == "internal/dreefile/i18n") return 0
  if (pkg == "internal/dreefile/dreecode" || pkg == "internal/dreefile/gogen" || pkg == "internal/dreefile/codegen") return in_list(imp, "internal/dreefile/ir")
  if (pkg == "internal/dreefile/jsoutput") return in_list(imp, "internal/dreefile/gogen")
  if (pkg == "internal/dreefile/lexer" || pkg == "internal/dreefile/parser") return in_list(imp, "internal/dreefile/dreecode internal/dreefile/ir internal/dreefile/tokens")
  if (pkg == "internal/dreefile/sections/client/js" || pkg == "internal/dreefile/sections/client/ts" || pkg == "internal/dreefile/sections/client/lua") return in_list(imp, "internal/dreefile/codegen internal/dreefile/ir internal/dreefile/jsoutput")
  if (pkg == "internal/dreefile/sections/client") return in_list(imp, g " internal/dreefile/sections/client/js internal/dreefile/sections/client/ts internal/dreefile/sections/client/lua")
  if (pkg == "internal/dreefile/sections/body/html") return in_list(imp, g " internal/dreefile/sections/head internal/dreefile/sections/style internal/dreefile/sections/client")
  if (pkg == "internal/dreefile/sections/body/md") return in_list(imp, g " internal/md")
  if (pkg ~ /^internal\/dreefile\/sections\//) return in_list(imp, g)
  if (pkg ~ /^internal\/dreefile\//) return 0
  return -1
}
/^@/ { pkg = substr($0, 2); next }
{
  if ($0 !~ ("^" module)) next
  imp = substr($0, length(module) + 1)
  pkgrel = substr(pkg, length(module) + 1)
  if (pkgrel == "internal/dreefile" || pkgrel ~ /^internal\/dreefile\//) {
    if (imp == "core" || imp == "adapter" || imp ~ /^adapter\//) { print pkgrel "|" imp; next }
    if (pkgrel == "internal/dreefile/sections/body/md" && imp == "internal/md") next
    if (imp == "internal/md" && pkgrel != "internal/dreefile/sections/body/md") { print pkgrel "|" imp; next }
    if (!allowed(pkgrel, imp)) print pkgrel "|" imp
    next
  }
  if (pkgrel ~ /^internal\//) {
    if (imp == "core" || imp == "adapter" || imp ~ /^adapter\// || imp ~ /^internal\/dreefile(\/|$)/) print pkgrel "|" imp
  }
}')"

if [ -z "$violations" ]; then
	echo "PASS: internal layering and dreefile dependency rule holds"
	exit 0
fi

echo "FAIL: internal layering violations"
printf '%s\n' "$violations" | while IFS='|' read -r pkgrel imp; do
	full="$module$imp"
	loc="$(grep -E -n -H "^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*[[:space:]]+)?\"$full\"" "$pkgrel"/*.go 2>/dev/null | grep -v -- '_test.go:' | head -n 1 || true)"
	if [ -n "$loc" ]; then
		echo "  $loc"
	else
		echo "  $pkgrel imports $full"
	fi
done
exit 1
