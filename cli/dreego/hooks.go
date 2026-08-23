package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type pluginManifest struct {
	Build pluginBuild `json:"build"`
}

type pluginBuild struct {
	Steps []pluginBuildStep `json:"steps"`
}

type pluginBuildStep struct {
	Cmd  string `json:"cmd"`
	When string `json:"when"`
}

func runBuildHooks(cwd string) error {
	gm, err := parseGoMod(filepath.Join(cwd, "go.mod"))
	if err != nil {
		return nil
	}
	var plugins []string
	for path := range gm.Requires {
		if strings.HasPrefix(path, pluginOrgPrefix) && path != coreModule {
			plugins = append(plugins, path)
		}
	}
	sort.Strings(plugins)
	for _, modPath := range plugins {
		dir, err := findModDir(cwd, modPath)
		if err != nil {
			continue
		}
		manifestPath := filepath.Join(dir, "dreego-plugin.json")
		body, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}
		var manifest pluginManifest
		if err := json.Unmarshal(body, &manifest); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %s: invalid dreego-plugin.json: %v\n", modPath, err)
			continue
		}
		for _, step := range manifest.Build.Steps {
			if step.When != "pre-build" {
				continue
			}
			if step.Cmd == "" {
				continue
			}
			fmt.Printf("plugin %s: %s\n", modPath, step.Cmd)
			c := exec.Command("sh", "-c", step.Cmd)
			c.Dir = cwd
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("plugin %s build step failed: %w", modPath, err)
			}
		}
	}
	return nil
}