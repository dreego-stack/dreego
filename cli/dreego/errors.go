package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	transpiler "github.com/dreego-stack/dreego/internal/transpiler"
)

var posRe = regexp.MustCompile(`at position (\d+)`)

func formatGenerateError(err error) string {
	msg := err.Error()
	if strings.Contains(msg, "generated code is out of date") {
		return msg + "Fix: run `dreego generate` to transpile .dreego files into Go code\n"
	}
	file, rest := splitErrorFile(msg)
	if file == "" {
		return msg
	}
	loc := "?:?"
	fix := fixForCause(rest)
	if m := posRe.FindStringSubmatch(rest); m != nil {
		var pos int
		fmt.Sscanf(m[1], "%d", &pos)
		if data, e := os.ReadFile(file); e == nil {
			raw := string(data)
			_, _, body := transpiler.ParseHeader(raw)
			loc = lineCol(raw, pos+len(raw)-len(body))
		}
	}
	cause := strings.TrimSpace(posRe.ReplaceAllString(rest, ""))
	cause = strings.TrimRight(cause, " :")
	return fmt.Sprintf("%s:%s: %s Fix: %s", file, loc, cause, fix)
}

func splitErrorFile(msg string) (file, rest string) {
	for _, prefix := range []string{"error parsing ", "error lexing ", "error generating ", "error generating error page "} {
		if strings.HasPrefix(msg, prefix) {
			after := strings.TrimPrefix(msg, prefix)
			if idx := strings.Index(after, ": "); idx >= 0 {
				return after[:idx], after[idx+2:]
			}
			return after, ""
		}
	}
	return "", msg
}

func fixForCause(cause string) string {
	switch {
	case strings.Contains(cause, "unclosed {#if"):
		return "close the {#if} block with {/if}"
	case strings.Contains(cause, "unclosed {#each"):
		return "close the {#each} block with {/each}"
	case strings.Contains(cause, "unclosed {#slot"):
		return "close the {#slot} block with {/slot}"
	case strings.Contains(cause, "unclosed {#verbatim"):
		return "close the {#verbatim} block with {/verbatim}"
	case strings.Contains(cause, "unclosed <div>"):
		return "close the <div> with </div>"
	case strings.Contains(cause, "unclosed <"):
		return "close the corresponding opening tag"
	default:
		return "fix the reported issue and run dreego generate again"
	}
}

func lineCol(src string, pos int) string {
	line, col := 1, 1
	for i := 0; i < pos && i < len(src); i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return fmt.Sprintf("%d:%d", line, col)
}
