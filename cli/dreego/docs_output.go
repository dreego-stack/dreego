package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const feedbackURL = "https://github.com/dreego-stack/dreego/issues/new"

var headingPattern = regexp.MustCompile(`^#{1,6}\s+(.*)`)
var codeBlockPattern = regexp.MustCompile("`{3}[^`]*`{3}")
var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func printJSON(webBase, path string, body []byte) {
	text := string(body)
	var headings []string
	for _, line := range strings.Split(text, "\n") {
		if match := headingPattern.FindStringSubmatch(line); match != nil {
			headings = append(headings, strings.TrimSpace(match[1]))
		}
	}
	codeBlocks := codeBlockPattern.FindAllString(text, -1)
	var links [][2]string
	for _, match := range linkPattern.FindAllStringSubmatch(text, -1) {
		links = append(links, [2]string{match[1], match[2]})
	}
	doc := map[string]any{
		"path": path, "source": webBase + path, "headings": headings,
		"code_blocks": codeBlocks, "links": links, "length": len(text),
	}
	out, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(out))
}

func printDoc(body []byte, webBase, rawBase string) {
	out := strings.ReplaceAll(string(body), webBase, "")
	out = strings.ReplaceAll(out, rawBase, "")
	fmt.Println(out)
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
		fmt.Fprintf(os.Stderr, "could not open browser: %v\nvisit: %s\n", err, url)
	}
}
