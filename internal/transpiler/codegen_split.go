package transpiler

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func splitServerSections(sections []ServerSection, hasFormActions bool) (pkgCode string, inlineCode string) {
	var pkg []string
	var inl []string
	for _, g := range sections {
		if g.ContentType != "" && g.ContentType != "custom" {
			continue
		}
		trimmed := strings.TrimSpace(g.Code)
		if trimmed == "" {
			continue
		}
		firstLine := ""
		for _, line := range strings.Split(trimmed, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			firstLine = line
			break
		}
		isDeclaration := strings.HasPrefix(firstLine, "type ") || strings.HasPrefix(firstLine, "func ")
		if isDeclaration && hasFormActions {
			pkg = append(pkg, ir.TranslateMdtohtml(unindent(g.Code)))
		} else {
			inl = append(inl, ir.TranslateMdtohtml(trimmed))
		}
	}
	result := strings.Join(pkg, "\n")
	if result != "" {
		result += "\n"
	}
	return result, strings.Join(inl, "\n")
}

func unindent(code string) string {
	lines := strings.Split(code, "\n")
	minIndent := -1
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, " \t"))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent <= 0 {
		return code
	}
	var out []string
	for _, l := range lines {
		if len(l) >= minIndent {
			out = append(out, l[minIndent:])
		} else {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}
