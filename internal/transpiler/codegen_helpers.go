package transpiler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// attrNameAt returns the name of the attribute whose value contains the byte
// at index i of a tag string, or "" when the position is not inside an
// attribute value. Quoted and unquoted values are recognized, and whitespace
// around "=" is tolerated. It is used to pick the safe-value rule for an
// expression placeholder inside an attribute.
func attrNameAt(tag string, i int) string {
	if i < 0 || i >= len(tag) {
		return ""
	}
	pos := 0
	for pos < len(tag) {
		pos = skipAttrSpace(tag, pos)
		if pos >= len(tag) {
			return ""
		}
		nameStart := pos
		for pos < len(tag) && tag[pos] != '=' && !isAttrSpace(tag[pos]) {
			pos++
		}
		name := tag[nameStart:pos]
		pos = skipAttrSpace(tag, pos)
		if pos >= len(tag) || tag[pos] != '=' {
			continue
		}
		pos = skipAttrSpace(tag, pos+1)
		if pos >= len(tag) {
			return ""
		}
		if tag[pos] == '"' || tag[pos] == '\'' {
			quote := tag[pos]
			pos++
			valStart := pos
			for pos < len(tag) && tag[pos] != quote {
				pos++
			}
			if pos >= len(tag) {
				return ""
			}
			if i >= valStart && i < pos {
				return name
			}
			pos++
			continue
		}
		valStart := pos
		for pos < len(tag) && !isAttrSpace(tag[pos]) && tag[pos] != '>' {
			pos++
		}
		if i >= valStart && i < pos {
			return name
		}
	}
	return ""
}

// isURLAttr reports whether an attribute name carries a URL value that must be
// scheme-validated (href, src, action, and the srcset family).
func isURLAttr(name string) bool {
	switch strings.ToLower(name) {
	case "href", "src", "action", "srcset", "poster", "cite", "formaction", "data", "background", "longdesc", "usemap", "xlink:href":
		return true
	}
	return false
}

// isScriptAttr reports whether an attribute name carries a script value (event
// handler attributes such as onclick, onload, onerror, the Alpine/HTMX event
// directives x-on:* and @*, the HTMX response-trigger directives hx-on:* and
// hx-on::*, and the Alpine evaluator directives x-data, x-init, x-effect,
// x-html, x-show, x-model, x-text, and x-transition). The plain "on" prefix
// must be followed by a letter so attributes like "once" or "only" are
// conservatively JSON-encoded rather than mistaken for event handlers that do
// not exist.
func isScriptAttr(name string) bool {
	n := strings.ToLower(name)
	if strings.HasPrefix(n, "on") && len(n) > 2 && n[2] >= 'a' && n[2] <= 'z' {
		return true
	}
	switch {
	case strings.HasPrefix(n, "x-on:"), strings.HasPrefix(n, "x-on."), strings.HasPrefix(n, "@"):
		return true
	case strings.HasPrefix(n, "hx-on:"):
		return true
	}
	switch n {
	case "x-data", "x-init", "x-effect", "x-html", "x-show", "x-model", "x-text", "x-transition":
		return true
	}
	return false
}

// attrContext picks the safe-value rule name for an attribute by name. It
// returns "SafeScript", "SafeURL", "SafeStyle", or "SafeAttr". Alpine style
// bindings (x-bind:style, :style) are treated as style context; other directive
// attributes fall back to SafeAttr, which still prevents attribute breakout.
func attrContext(name string) string {
	n := strings.ToLower(name)
	if n == "srcdoc" {
		return "SafeSrcdoc"
	}
	if isScriptAttr(name) {
		return "SafeScript"
	}
	if isURLAttr(name) {
		return "SafeURL"
	}
	if n == "style" || n == "x-bind:style" || n == ":style" {
		return "SafeStyle"
	}
	return "SafeAttr"
}

func goLiteral(s string) string {
	if strings.Contains(s, "`") {
		return strconv.Quote(s)
	}
	return "`" + s + "`"
}

func toPascalCase(s string) string {
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

func extractAttrValues(attrs string) string {
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
				vals = append(vals, attrVal(attrs[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(attrs) {
		vals = append(vals, attrVal(attrs[start:]))
	}
	return strings.Join(vals, ", ")
}

func attrVal(part string) string {
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
		return concatPlaceholders(val)
	}
	return fmt.Sprintf("%q", val)
}

// concatPlaceholders splits a quoted attribute value that mixes literal text and
// one or more {…} placeholders into a Go concatenation, quoting literal segments
// and emitting each expression as fmt.Sprintf. Escaping is deferred to the prop
// injection point (component bodies escape their own placeholders), matching the
// single-placeholder raw path so multi-placeholder calls are not double-escaped.
func concatPlaceholders(val string) string {
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
