#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
[ -f main.go ] || { echo "missing main.go"; exit 1; }
[ -f dreego/routes/get.dreego ] || { echo "missing get.dreego"; exit 1; }
echo ok
