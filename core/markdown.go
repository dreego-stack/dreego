package core

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html/md"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func MarkdownToHTML(src string) (string, error) {
	nodes, err := md.ToNodes(src)
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
