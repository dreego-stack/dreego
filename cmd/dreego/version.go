package main

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

// version is injected at build time via:
//
//	go build -ldflags "-X main.version=$(cat VERSION)" ./cmd/dreego
var version string

// dreegoVersion returns the CLI version. Resolution order:
//
//  1. version injected via -ldflags (make build, release)
//  2. module version from build info (go install pkg@tag)
//  3. VERSION file in the current working directory or any parent directory
//  4. VERSION file next to this source file (repo-local dev builds from any cwd)
//  5. "dev" fallback for plain local builds
func dreegoVersion() string {
	if version != "" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	if v := versionFromWalkUp("VERSION"); v != "" {
		return v
	}
	if v := versionFromSourceRoot("VERSION"); v != "" {
		return v
	}
	return "dev"
}

func versionFromWalkUp(name string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		v, ok := readVersionFile(filepath.Join(dir, name))
		if ok {
			return v
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func versionFromSourceRoot(name string) string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file)
	for {
		v, ok := readVersionFile(filepath.Join(dir, name))
		if ok {
			return v
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func readVersionFile(path string) (string, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	v := strings.TrimSpace(string(b))
	if v == "" || v == "(devel)" || v == "dev" {
		return "", true
	}
	return v, true
}
