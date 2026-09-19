package gogen

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func GoLiteral(s string) string {
	if strings.Contains(s, "`") {
		return strconv.Quote(s)
	}
	return "`" + s + "`"
}

func ToPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		result.WriteByte(byte(unicode.ToUpper(rune(p[0]))))
		if len(p) > 1 {
			result.WriteString(strings.ToLower(p[1:]))
		}
	}
	return result.String()
}

func SourceRef(src string, pos int) string {
	loc := SourceLocation(src, pos)
	if src == "" {
		return "?:" + loc
	}
	return fmt.Sprintf("%s:%s", src, loc)
}
