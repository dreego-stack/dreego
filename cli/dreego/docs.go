package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const coreModule = "github.com/dreego-stack/dreego"
const pluginOrgPrefix = "github.com/dreego-stack/"

var modCacheDir = ""

var wdFunc = func() string {
	d, _ := os.Getwd()
	return d
}

type goMod struct {
	Module   string
	Requires map[string]string
}

type sitemapDoc struct {
	Module string        `json:"module"`
	Pages  []sitemapPage `json:"pages"`
}

type sitemapPage struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

func parseGoMod(path string) (*goMod, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gm := &goMod{Requires: map[string]string{}}
	scanner := bufio.NewScanner(f)
	inRequire := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if inRequire {
			if line == ")" {
				inRequire = false
				continue
			}
			if f := strings.Fields(line); len(f) >= 2 {
				gm.Requires[f[0]] = f[1]
			}
			continue
		}
		switch {
		case strings.HasPrefix(line, "module "):
			gm.Module = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		case strings.HasPrefix(line, "require ("):
			inRequire = true
		case strings.HasPrefix(line, "require "):
			if f := strings.Fields(strings.TrimPrefix(line, "require ")); len(f) >= 2 {
				gm.Requires[f[0]] = f[1]
			}
		}
	}
	return gm, scanner.Err()
}

// findModDir locates the on-disk directory for a module path, reading go.mod
// in cwd. Priority: the module itself (self-repo) → vendor/ → module cache.
func findModDir(cwd, modPath string) (string, error) {
	modFile := filepath.Join(cwd, "go.mod")
	gm, err := parseGoMod(modFile)
	if err != nil {
		return "", err
	}
	if gm.Module == modPath {
		return cwd, nil
	}
	version, ok := gm.Requires[modPath]
	if !ok {
		return "", fmt.Errorf("%s not found in %s", modPath, modFile)
	}
	vendorDir := filepath.Join(cwd, "vendor", filepath.FromSlash(modPath))
	if _, err := os.Stat(vendorDir); err == nil {
		return vendorDir, nil
	}
	cacheDir := filepath.Join(goModCache(), filepath.FromSlash(modPath+"@"+version))
	if _, err := os.Stat(cacheDir); err == nil {
		return cacheDir, nil
	}
	return "", fmt.Errorf("module %s (%s) not downloaded; run go mod download", modPath, version)
}

func goModCache() string {
	if modCacheDir != "" {
		return modCacheDir
	}
	out, err := exec.Command("go", "env", "GOMODCACHE").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func readDocFrom(dir, path string) ([]byte, error) {
	rel := strings.TrimPrefix(path, "/")
	full := filepath.Join(dir, rel)
	cleanDir := filepath.Clean(dir)
	cleanFull := filepath.Clean(full)
	if cleanFull != cleanDir && !strings.HasPrefix(cleanFull, cleanDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("path escapes module dir: %s", path)
	}
	return os.ReadFile(full)
}

func readSitemap(dir string) (*sitemapDoc, error) {
	body, err := os.ReadFile(filepath.Join(dir, "_docs", "sitemap.json"))
	if err != nil {
		return nil, err
	}
	var sm sitemapDoc
	if err := json.Unmarshal(body, &sm); err != nil {
		return nil, err
	}
	return &sm, nil
}

func sitemapPaths(dir string) []string {
	sm, err := readSitemap(dir)
	if err != nil {
		return nil
	}
	var paths []string
	for _, p := range sm.Pages {
		paths = append(paths, p.Path)
	}
	return paths
}

func webBases(modPath string) (web, raw string) {
	repo := strings.TrimPrefix(modPath, "github.com/")
	return "https://github.com/" + repo + "/blob/main",
		"https://raw.githubusercontent.com/" + repo + "/main"
}

func cmdDocs(args []string) {
	web, dump, jsonOut, list := false, false, false, false
	plugin := ""
	var remaining []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--web":
			web = true
		case a == "--dump":
			dump = true
		case a == "--json":
			jsonOut = true
		case a == "--list":
			list = true
		case a == "-p" && i+1 < len(args):
			plugin = args[i+1]
			i++
		case strings.HasPrefix(a, "-p="):
			plugin = strings.TrimPrefix(a, "-p=")
		default:
			remaining = append(remaining, a)
		}
	}

	if list {
		cmdList()
		return
	}

	path := "/_docs/index.md"
	if len(remaining) > 0 {
		path = remaining[0]
	}
	dumpAll := path == "all"
	if !strings.HasPrefix(path, "/") && !dumpAll {
		path = "/" + path
	}

	modPath := coreModule
	if plugin != "" {
		modPath = pluginOrgPrefix + plugin
	}
	webBase, rawBase := webBases(modPath)

	dir, err := findModDir(wdFunc(), modPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docs error: %v\n", err)
		os.Exit(1)
	}

	if web {
		openBrowser(webBase + path)
		return
	}

	if dump {
		cmdDump(dir, path, webBase, rawBase)
		return
	}

	body, err := readDocFrom(dir, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docs error: %v\n", err)
		os.Exit(1)
	}

	if jsonOut {
		printJSON(webBase, path, body)
		return
	}

	printDoc(body, webBase, rawBase)
}

func cmdDump(dir, path, webBase, rawBase string) {
	var pages []string
	if path == "/_docs/index.md" || path == "all" {
		pages = sitemapPaths(dir)
	} else {
		pages = strings.Split(path, ",")
	}
	for _, p := range pages {
		p = strings.TrimSpace(p)
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		body, err := readDocFrom(dir, p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", p, err)
			continue
		}
		fmt.Printf("\n--- %s ---\n\n", p)
		printDoc(body, webBase, rawBase)
	}
	fmt.Println()
}

func cmdList() {
	cwd := wdFunc()
	gm, err := parseGoMod(filepath.Join(cwd, "go.mod"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "docs error: no go.mod found (%v)\n", err)
		os.Exit(1)
	}
	seen := map[string]bool{}
	mods := []string{}
	if !seen[gm.Module] {
		seen[gm.Module] = true
		mods = append(mods, gm.Module)
	}
	var plugins []string
	for path := range gm.Requires {
		if strings.HasPrefix(path, pluginOrgPrefix) && !seen[path] {
			seen[path] = true
			plugins = append(plugins, path)
		}
	}
	sort.Strings(plugins)
	mods = append(mods, plugins...)

	for _, m := range mods {
		dir, err := findModDir(cwd, m)
		if err != nil {
			fmt.Printf("\n[%s]\n  (not downloaded)\n", m)
			continue
		}
		sm, err := readSitemap(dir)
		if err != nil {
			fmt.Printf("\n[%s]\n  (no _docs/sitemap.json)\n", m)
			continue
		}
		fmt.Printf("\n[%s]\n", m)
		for _, p := range sm.Pages {
			if p.Title != "" {
				fmt.Printf("  %-35s %s\n", p.Path, p.Title)
			} else {
				fmt.Printf("  %s\n", p.Path)
			}
		}
	}
	fmt.Println()
}
