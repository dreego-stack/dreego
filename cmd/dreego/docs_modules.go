package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
)

type listedModule struct {
	Dir     string
	Version string
	Error   string
}

func findModDir(cwd, modPath string) (string, error) {
	modFile := filepath.Join(cwd, "go.mod")
	gm, err := parseGoMod(modFile)
	if err != nil {
		return "", err
	}
	if gm.Module == modPath {
		return cwd, nil
	}
	version, required := gm.Requires[modPath]
	vendorDir := filepath.Join(cwd, "vendor", filepath.FromSlash(modPath))
	if _, err := os.Stat(vendorDir); err == nil {
		return vendorDir, nil
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
		download := exec.Command("go", "mod", "download", "-json", modPath+"@"+version)
		download.Dir = cwd
		if out, err := download.Output(); err == nil {
			var module listedModule
			if json.Unmarshal(out, &module) == nil && module.Dir != "" && module.Error == "" {
				return module.Dir, nil
			}
		}
	}
	if listErr != nil {
		if !required {
			return "", fmt.Errorf("%s not found in %s or CLI build information", modPath, modFile)
		}
		return "", fmt.Errorf("resolve module %s: %w", modPath, listErr)
	}
	return "", fmt.Errorf("module %s (%s) has no local directory", modPath, version)
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
