package dreefile

import (
	"go/scanner"
	"go/token"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func splitServerSections(sections []ServerSection, emitted map[string]bool) (pkgCode string, inlineCode string) {
	var pkg []string
	var inl []string
	seen := map[string]bool{}
	for _, g := range sections {
		if g.ContentType != "" && g.ContentType != "custom" {
			continue
		}
		code := strings.TrimSpace(ir.TranslateMdtohtml(g.Code))
		if code == "" {
			continue
		}
		leading := true
		for _, chunk := range splitTopLevelChunks(code) {
			clean := strings.TrimSpace(unindent(chunk))
			if clean == "" {
				continue
			}
			if isFuncDeclaration(clean) || (leading && isPackageDeclaration(clean)) {
				if name := declarationName(clean); name != "" {
					if seen[name] || emitted[name] {
						continue
					}
					seen[name] = true
					emitted[name] = true
				}
				pkg = append(pkg, clean)
				continue
			}
			leading = false
			inl = append(inl, clean)
		}
	}
	result := strings.Join(pkg, "\n")
	if result != "" {
		result += "\n"
	}
	return result, strings.Join(inl, "\n")
}

// hoistedDeclarationNames lists the package-level names a route file hoists to
// the shared routes package. Method sections repeat the same leading
// declaration block, so the set is deduplicated.
func hoistedDeclarationNames(file *File) []string {
	seen := map[string]bool{}
	var names []string
	for _, g := range file.Server {
		if g.ContentType != "" && g.ContentType != "custom" {
			continue
		}
		code := strings.TrimSpace(ir.TranslateMdtohtml(g.Code))
		leading := true
		for _, chunk := range splitTopLevelChunks(code) {
			clean := strings.TrimSpace(unindent(chunk))
			if clean == "" {
				continue
			}
			if !isFuncDeclaration(clean) && !(leading && isPackageDeclaration(clean)) {
				leading = false
				continue
			}
			name := declarationName(clean)
			if name != "" && !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// splitTopLevelChunks splits Go source at automatic semicolons (newlines) that
// sit outside any brace, bracket, or parenthesis. Explicit semicolons in for
// and if headers are preserved: the scanner reports them with literal ";",
// while automatic newline semicolons carry the literal "\n".
func splitTopLevelChunks(code string) []string {
	fset := token.NewFileSet()
	file := fset.AddFile("server.go", -1, len(code))
	var s scanner.Scanner
	s.Init(file, []byte(code), func(token.Position, string) {}, scanner.ScanComments)
	depth := 0
	start := 0
	hasCode := false
	var chunks []string
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		off := file.Offset(pos)
		switch tok {
		case token.LPAREN, token.LBRACK, token.LBRACE:
			depth++
		case token.RPAREN, token.RBRACK, token.RBRACE:
			if depth > 0 {
				depth--
			}
		case token.SEMICOLON:
			if depth == 0 && hasCode && lit == "\n" {
				chunks = append(chunks, code[start:off])
				start = off
				hasCode = false
			}
		}
		if tok != token.COMMENT && tok != token.SEMICOLON {
			hasCode = true
		}
	}
	if strings.TrimSpace(code[start:]) != "" {
		chunks = append(chunks, code[start:])
	}
	return chunks
}

func isPackageDeclaration(code string) bool {
	switch firstCodeToken(code) {
	case token.TYPE, token.FUNC, token.VAR, token.CONST:
		return true
	default:
		return false
	}
}

func isFuncDeclaration(code string) bool {
	return firstCodeToken(code) == token.FUNC
}

func firstCodeToken(code string) token.Token {
	fset := token.NewFileSet()
	file := fset.AddFile("server.go", -1, len(code))
	var s scanner.Scanner
	s.Init(file, []byte(code), func(token.Position, string) {}, scanner.ScanComments)
	for {
		_, tok, _ := s.Scan()
		if tok == token.EOF || tok != token.COMMENT {
			return tok
		}
	}
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
