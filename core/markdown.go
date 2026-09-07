package core

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/markdown"
)

var markdownLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

var trustedWarnOnce sync.Once

func MarkdownToHTML(src string) (string, error) {
	return markdownToHTML(src, markdown.ModeSafe)
}

func MarkdownToHTMLTrusted(src string) (string, error) {
	trustedWarnOnce.Do(func() {
		markdownLogger.Warn("MarkdownToHTMLTrusted: raw HTML passthrough enabled — only use with trusted content")
	})
	return markdownToHTML(src, markdown.ModeTrusted)
}

func markdownToHTML(src string, mode markdown.Mode) (string, error) {
	nodes, err := markdown.ToNodes(src, mode)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, n := range nodes {
		if n.Type != ir.NodeText {
			return "", fmt.Errorf("markdown: control flow not supported in runtime markdown")
		}
		b.WriteString(n.Content)
	}
	return b.String(), nil
}
