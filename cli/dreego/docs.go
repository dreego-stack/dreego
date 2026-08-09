package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const docsBaseURL = "https://raw.githubusercontent.com/dreego-stack/dreego/main"
const docsWebBase = "https://github.com/dreego-stack/dreego/blob/main"
const feedbackURL = "https://github.com/dreego-stack/dreego/issues/new"

var headingPattern = regexp.MustCompile(`^#{1,6}\s+(.*)`)
var codeBlockPattern = regexp.MustCompile("`{3}[^`]*`{3}")
var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

var pluginDocsRoot = "plugins"
var fetchDocFallback = fetchDocEmbedded

func fetchDocLocal(path string) ([]byte, bool, error) {
	rel := strings.TrimPrefix(path, "/")
	parts := strings.SplitN(rel, "/", 2)
	if len(parts) != 2 || parts[0] != "plugins" {
		// Not a plugin path: no local doc. The caller (cmdDocs/cmdDump)
		// decides whether to invoke the fallback. Do NOT call it here.
		return nil, false, nil
	}
	full := filepath.Join(pluginDocsRoot, parts[1])
	body, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return body, true, nil
}

func cmdDocs(args []string) {
	web := false
	dump := false
	jsonOut := false
	var remaining []string
	for _, a := range args {
		switch a {
		case "--web":
			web = true
		case "--dump":
			dump = true
		case "--json":
			jsonOut = true
		default:
			remaining = append(remaining, a)
		}
	}

	path := "/_docs/index.md"
	if len(remaining) > 0 {
		path = remaining[0]
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if web {
		openBrowser(docsWebBase + path)
		return
	}

	if dump {
		cmdDump(path)
		return
	}

	body, fromLocal, err := fetchDocLocal(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docs error: %v\n", err)
		os.Exit(1)
	}
	if !fromLocal {
		body, err = fetchDocFallback(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "docs error: %v\n", err)
			os.Exit(1)
		}
	}

	if jsonOut {
		printJSON(path, body)
		return
	}

	out := string(body)
	out = strings.ReplaceAll(out, docsWebBase, "")
	out = strings.ReplaceAll(out, docsBaseURL, "")
	fmt.Print(out)
	fmt.Println()
}

func cmdDump(path string) {
	var pages []string
	if path == "/_docs/index.md" || path == "all" {
		pages = []string{
			"/_docs/index.md",
			"/_docs/getting-started.md",
			"/_docs/cli.md",
			"/_docs/routing.md",
			"/_docs/components.md",
			"/_docs/middleware.md",
			"/_docs/config.md",
			"/_docs/runtime.md",
			"/_docs/testing.md",
			"/README.md",
			"/CHANGELOG.md",
		}
	} else {
		pages = strings.Split(path, ",")
	}

	for _, p := range pages {
		p = strings.TrimSpace(p)
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		body, fromLocal, err := fetchDocLocal(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", p, err)
			continue
		}
		if !fromLocal {
			body, err = fetchDocFallback(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip %s: %v\n", p, err)
				continue
			}
		}
		fmt.Printf("\n--- %s ---\n\n", p)
		out := string(body)
		out = strings.ReplaceAll(out, docsWebBase, "")
		out = strings.ReplaceAll(out, docsBaseURL, "")
		fmt.Print(out)
	}
	fmt.Println()
}

func printJSON(path string, body []byte) {
	text := string(body)

	var headings []string
	for _, line := range strings.Split(text, "\n") {
		if m := headingPattern.FindStringSubmatch(line); m != nil {
			headings = append(headings, strings.TrimSpace(m[1]))
		}
	}

	var codeBlocks []string
	for _, m := range codeBlockPattern.FindAllString(text, -1) {
		codeBlocks = append(codeBlocks, m)
	}

	var links [][2]string
	for _, m := range linkPattern.FindAllStringSubmatch(text, -1) {
		if len(m) >= 3 {
			links = append(links, [2]string{m[1], m[2]})
		}
	}

	doc := map[string]any{
		"path":        path,
		"source":      docsWebBase + path,
		"headings":    headings,
		"code_blocks": codeBlocks,
		"links":       links,
		"length":      len(text),
	}

	out, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(out))
}

func cmdFeedback() {
	fmt.Printf("Opening %s\n", feedbackURL)
	openBrowser(feedbackURL)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = exec.Command("cmd", "/c", "start", url).Start()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open browser: %v\n", err)
		fmt.Fprintf(os.Stderr, "visit: %s\n", url)
	}
}
