#!/bin/sh
# Create release tags for all modules from the single source VERSION file.
# Usage: _scripts/release.sh
# Only tags locally. Push manually after review:
#   git push origin core/<V> cmd/dreego/<V> plugins/sample/<V>
set -e

cd "$(dirname "$0")/.."

V=$(cat VERSION)
# Sanity: VERSION must be a valid v-prefixed semver, e.g. v0.0.22
case "$V" in
v[0-9]*\.[0-9]*\.[0-9]*) ;;
*)
	echo "error: VERSION must look like v0.0.22, got: '$V'"
	exit 1
	;;
esac

for mod in core cmd/dreego plugins/sample; do
	git tag "$mod/$V"
	echo "tagged $mod/$V"
done

echo ""
echo "Next (manual):"
echo "  git push origin core/$V cmd/dreego/$V plugins/sample/$V"
