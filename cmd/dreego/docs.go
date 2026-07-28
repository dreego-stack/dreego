package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"regexp"
	"strings"
)

const docsBaseURL = "https://codeberg.org/dreego/dreego/raw/branch/main"
const docsWebBase = "https://codeberg.org/dreego/dreego/src/branch/main"
const feedbackURL = "https://codeberg.org/dreego/dreego/issues/new"

var headingPattern = regexp.MustCompile(`^#{1,6}\s+(.*)`)
var codeBlockPattern = regexp.MustCompile("`{3}[^`]*`{3}")
var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

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

	body, err := fetchDoc(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docs error: %v\n", err)
		os.Exit(1)
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

func fetchDoc(path string) ([]byte, error) {
	url := docsBaseURL + path
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s not found (%d)", path, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
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
		body, err := fetchDoc(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", p, err)
			continue
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
