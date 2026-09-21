package dreefile

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
)

// declarationName returns the package-level identifier declared by a hoisted
// server declaration (type, func, var, const). It returns "" when the chunk is
// not a simple single-name declaration, for example a grouped const/var block
// or a method declaration, which the caller reports without a name.
func declarationName(code string) string {
	fset := token.NewFileSet()
	file := fset.AddFile("server.go", -1, len(code))
	var s scanner.Scanner
	s.Init(file, []byte(code), func(token.Position, string) {}, scanner.ScanComments)
	seenKeyword := false
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			return ""
		case token.COMMENT, token.SEMICOLON:
			continue
		case token.TYPE, token.FUNC, token.VAR, token.CONST:
			if seenKeyword {
				return ""
			}
			seenKeyword = true
			continue
		case token.IDENT:
			if seenKeyword {
				return lit
			}
			return ""
		default:
			return ""
		}
	}
}

// serverDeclarationConflict builds a generate-time error when two route files
// hoist the same package-level name. All route files share one Go package, so
// the duplicate would otherwise surface as a raw compiler "redeclared in this
// block" error with no dreego source location.
func serverDeclarationConflict(name, first, second string) error {
	d := Diagnostic{
		File:  second,
		Line:  1,
		Cause: fmt.Sprintf("duplicate package-level declaration %q; it is already declared in %s", name, first),
		Fix:   "declare it once and reuse it, or give the second declaration a distinct name",
	}
	return errors.New(d.String())
}
