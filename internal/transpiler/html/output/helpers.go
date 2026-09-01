package output

import (
	"fmt"
	"strconv"
	"strings"
)

func ExtractAttrValues(attrs string) string {
	if attrs == "" {
		return ""
	}
	var vals []string
	inQuote := false
	braceDepth := 0
	start := 0
	for i := 0; i < len(attrs); i++ {
		ch := attrs[i]
		if ch == '"' && braceDepth == 0 {
			inQuote = !inQuote
		}
		if ch == '{' && !inQuote {
			braceDepth++
		}
		if ch == '}' && !inQuote {
			braceDepth--
		}
		if ch == ' ' && !inQuote && braceDepth == 0 {
			if start < i {
				vals = append(vals, AttrVal(attrs[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(attrs) {
		vals = append(vals, AttrVal(attrs[start:]))
	}
	return strings.Join(vals, ", ")
}

func AttrVal(part string) string {
	eq := strings.IndexByte(part, '=')
	if eq < 0 {
		return fmt.Sprintf("%q", part)
	}
	val := strings.TrimSpace(part[eq+1:])
	if val == "" {
		return fmt.Sprintf("%q", "")
	}
	if val[0] == '{' && val[len(val)-1] == '}' && strings.Count(val, "{") == 1 && strings.Count(val, "}") == 1 {
		return val[1 : len(val)-1]
	}
	val = strings.Trim(val, "\"")
	if len(val) >= 2 && val[0] == '{' && val[len(val)-1] == '}' && strings.Count(val, "{") == 1 && strings.Count(val, "}") == 1 {
		return val[1 : len(val)-1]
	}
	if strings.Contains(val, "{") {
		return ConcatPlaceholders(val)
	}
	return fmt.Sprintf("%q", val)
}

func ConcatPlaceholders(val string) string {
	var parts []string
	start := 0
	inExpr := false
	for i := 0; i < len(val); i++ {
		switch val[i] {
		case '{':
			if !inExpr {
				if start < i {
					parts = append(parts, strconv.Quote(val[start:i]))
				}
				inExpr = true
				start = i + 1
			}
		case '}':
			if inExpr {
				expr := val[start:i]
				parts = append(parts, fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr))
				inExpr = false
				start = i + 1
			}
		}
	}
	if start < len(val) {
		parts = append(parts, strconv.Quote(val[start:]))
	}
	return strings.Join(parts, " + ")
}
