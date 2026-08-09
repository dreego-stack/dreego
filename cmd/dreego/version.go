package main

import (
	"runtime/debug"
)

// version is injected at build time via:
//
//	go build -ldflags "-X main.version=$(git describe --tags --abbrev=0)" ./cmd/dreego
var version string

// dreegoVersion returns the CLI version. Resolution order:
//
//  1. version injected via -ldflags (make build, release)
//  2. module version from build info (go install pkg@tag)
//  3. "dev" fallback for plain local builds
func dreegoVersion() string {
	if version != "" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}
