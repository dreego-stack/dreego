package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
)

type listedModule struct {
	Dir     string
	Version string
	Error   string
}

func findModDir(cwd, modPath string) (string, error) {
	modFile := filepath.Join(cwd, "go.mod")
	gm, gmErr := parseGoMod(modFile)
	if gmErr == nil && gm.Module == modPath {
		return cwd, nil
	}
	if gmErr == nil {
		vendorDir := filepath.Join(cwd, "vendor", filepath.FromSlash(modPath))
		if _, err := os.Stat(vendorDir); err == nil {
			return vendorDir, nil
		}
	}
	cmd := exec.Command("go", "list", "-m", "-json", modPath)
	cmd.Dir = cwd
	out, listErr := cmd.Output()
	if listErr == nil {
		var module listedModule
		if json.Unmarshal(out, &module) == nil && module.Dir != "" {
			return module.Dir, nil
		}
	}
	if version := builtModuleVersion(modPath); version != "" {
		if dir, ok := moduleCacheDir(modPath, version); ok {
			return dir, nil
		}
		download := exec.Command("go", "mod", "download", "-json", modPath+"@"+version)
		download.Dir = cwd
		if out, err := download.Output(); err == nil {
			var module listedModule
			if json.Unmarshal(out, &module) == nil && module.Dir != "" && module.Error == "" {
				return module.Dir, nil
			}
		}
	} else if dir, ok := moduleCacheDir(modPath, ""); ok {
		return dir, nil
	}
	if gmErr != nil {
		return "", fmt.Errorf("%s not found: %w", modPath, gmErr)
	}
	if listErr != nil {
		if _, required := gm.Requires[modPath]; !required {
			return "", fmt.Errorf("%s not found in %s or CLI build information", modPath, modFile)
		}
		return "", fmt.Errorf("resolve module %s: %w", modPath, listErr)
	}
	return "", fmt.Errorf("module %s has no local directory", modPath)
}

// moduleCacheDir resolves a module from GOMODCACHE. It lets `dreego docs`
// work without a go.mod, where `go mod download` is unavailable. When version
// is empty (for example a repo-local build) it picks the cached version.
func moduleCacheDir(modPath, version string) (string, bool) {
	cache := exec.Command("go", "env", "GOMODCACHE")
	out, err := cache.Output()
	if err != nil {
		return "", false
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", false
	}
	escaped := filepath.FromSlash(escapeModulePath(modPath))
	if version != "" {
		dir := filepath.Join(root, escaped+"@"+version)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, true
		}
		return "", false
	}
	matches, err := filepath.Glob(filepath.Join(root, escaped+"@*"))
	if err != nil || len(matches) == 0 {
		return "", false
	}
	sort.Strings(matches)
	return matches[len(matches)-1], true
}

func escapeModulePath(path string) string {
	var b strings.Builder
	for _, r := range path {
		if r >= 'A' && r <= 'Z' {
			b.WriteByte('!')
			b.WriteRune(r + ('a' - 'A'))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func builtModuleVersion(path string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return moduleVersion(info, path)
}

func moduleVersion(info *debug.BuildInfo, path string) string {
	if info.Main.Path == path && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	for _, dep := range info.Deps {
		module := dep
		if dep.Replace != nil {
			module = dep.Replace
		}
		if dep.Path == path && module.Version != "" && module.Version != "(devel)" {
			return module.Version
		}
	}
	return ""
}
