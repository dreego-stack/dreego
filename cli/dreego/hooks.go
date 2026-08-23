package main

import (
	"encoding/json"
	"fmt"
	"io"
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

type approvalsFile struct {
	ApprovedHooks map[string]bool `json:"approvedHooks"`
}

var isTerminalFn = isTerminal

func runBuildHooks(cwd string, autoApprove bool, stdin io.Reader) error {
	approvals := readApprovals(cwd)
	changed := false

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
			if step.When != "pre-build" || step.Cmd == "" {
				continue
			}
			key := modPath + ":" + step.Cmd
			if !approvals[key] {
				if autoApprove {
					approvals[key] = true
					changed = true
				} else if !isTerminalFn(stdin) {
					return fmt.Errorf("plugin %s build step not approved: %s\nRun 'dreego build' interactively to approve, or use 'dreego build --yes', or pre-approve in dreego-build.json", modPath, step.Cmd)
				} else {
					approved := promptApproval(modPath, step.Cmd, stdin)
					if !approved {
						return fmt.Errorf("plugin %s build step not approved by user", modPath)
					}
					approvals[key] = true
					changed = true
				}
			}
			fmt.Printf("plugin %s: %s\n", modPath, step.Cmd)
			c := exec.Command("sh", "-c", step.Cmd)
			c.Dir = cwd
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				if changed {
					saveApprovals(cwd, approvals)
				}
				return fmt.Errorf("plugin %s build step failed: %w", modPath, err)
			}
		}
	}
	if changed {
		if err := saveApprovals(cwd, approvals); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not save build approvals: %v\n", err)
		}
	}
	return nil
}

func readApprovals(cwd string) map[string]bool {
	path := filepath.Join(cwd, "dreego-build.json")
	body, err := os.ReadFile(path)
	if err != nil {
		return map[string]bool{}
	}
	var af approvalsFile
	if err := json.Unmarshal(body, &af); err != nil {
		return map[string]bool{}
	}
	if af.ApprovedHooks == nil {
		return map[string]bool{}
	}
	return af.ApprovedHooks
}

func saveApprovals(cwd string, approvals map[string]bool) error {
	af := approvalsFile{ApprovedHooks: approvals}
	data, err := json.MarshalIndent(af, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(cwd, "dreego-build.json"), data, 0644)
}

func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func promptApproval(modPath, cmd string, stdin io.Reader) bool {
	fmt.Printf("plugin %q wants to run %q in this repo. Approve? y/N: ", modPath, cmd)
	var answer string
	fmt.Fscanln(stdin, &answer)
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}