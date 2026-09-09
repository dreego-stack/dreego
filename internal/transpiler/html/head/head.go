package head

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/parser"
)

func Gen(html string, bufName string) (string, error) {
	return generate(nil, html, bufName, "c")
}

func GenWithMessages(gen *codegen.State, html, bufName, contextName string) (string, error) {
	return generate(gen, html, bufName, contextName)
}

func generate(gen *codegen.State, html, bufName, contextName string) (string, error) {
	var out strings.Builder
	rest := html
	pos := 0
	for rest != "" {
		expressionOpen := strings.Index(rest, "{{")
		messageOpen := strings.Index(rest, "[[")
		open := firstOpen(expressionOpen, messageOpen)
		if open < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest)))
			break
		}
		if open > 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest[:open])))
			pos += open
			rest = rest[open:]
		}
		if strings.HasPrefix(rest, "[[") {
			if insideRawElement(html, pos) {
				out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral("[[")))
				pos += 2
				rest = rest[2:]
				continue
			}
			closeIndex := ir.FindMessageEnd(rest[2:])
			if closeIndex < 0 {
				return "", fmt.Errorf("unclosed message expression at position %d", pos)
			}
			closeIndex += 2
			key, arguments, err := parser.ParseMessageExpression(rest[2:closeIndex], pos)
			if err != nil {
				return "", err
			}
			if gen == nil {
				return "", fmt.Errorf("message expression in head requires generator state")
			}
			names := make([]string, 0, len(arguments))
			parts := make([]string, 0, len(arguments))
			for _, argument := range arguments {
				names = append(names, argument.Name)
				constructor := "StringMessageArg"
				switch gen.MessageArgumentKind(key, argument.Name) {
				case "number":
					constructor = "NumberMessageArg"
				case "time":
					constructor = "TimeMessageArg"
				}
				parts = append(parts, fmt.Sprintf("dreego.%s(%q, %s)", constructor, argument.Name, argument.Expression))
			}
			gen.RegisterMessageUse(key, names)
			call := fmt.Sprintf("dreego.Message(%s, %q", contextName, key)
			if len(parts) > 0 {
				call += ", " + strings.Join(parts, ", ")
			}
			call += ")"
			out.WriteString(fmt.Sprintf("%s.WriteString(dreego.%s(%s))\n", bufName, HeadSafeFunc(html, pos), call))
			pos += closeIndex + 2
			rest = rest[closeIndex+2:]
			continue
		}
		closeIdx := ir.FindExprEnd(rest[2:])
		if closeIdx < 0 {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, ir.GoLiteral(rest)))
			break
		}
		closeIdx += 2
		expr, filters := ir.ParseExpression(strings.TrimSpace(rest[2:closeIdx]))
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
		raw := false
		for _, f := range filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
			default:
				return "", fmt.Errorf("unknown filter '%s' at position %d", f, pos)
			}
		}
		if raw {
			out.WriteString(fmt.Sprintf("%s.WriteString(%s)\n", bufName, code))
		} else {
			out.WriteString(fmt.Sprintf("%s.WriteString(dreego.%s(%s))\n", bufName, HeadSafeFunc(html, pos), code))
		}
		pos += closeIdx + 2
		rest = rest[closeIdx+2:]
	}
	return out.String(), nil
}

func insideRawElement(source string, position int) bool {
	prefix := strings.ToLower(source[:position])
	for _, tag := range []string{"script", "style"} {
		open := strings.LastIndex(prefix, "<"+tag)
		close := strings.LastIndex(prefix, "</"+tag+">")
		tagEnd := strings.LastIndex(prefix, ">")
		if open > close && tagEnd > open {
			return true
		}
	}
	return false
}

func firstOpen(left, right int) int {
	if left < 0 {
		return right
	}
	if right < 0 || left < right {
		return left
	}
	return right
}

func HeadSafeFunc(html string, i int) string {
	if i < 0 || i >= len(html) {
		return "SafeText"
	}
	tagStart := strings.LastIndex(html[:i], "<")
	if tagStart < 0 || html[tagStart+1] == '/' || html[tagStart+1] == '!' {
		return "SafeText"
	}
	tagEnd := ir.TagEnd(html[tagStart:])
	if tagEnd < 0 {
		return "SafeText"
	}
	tagEnd += tagStart
	if isMetaRefresh(html[tagStart:tagEnd]) {
		return "SafeRefresh"
	}
	name := ir.AttrNameAt(html[tagStart:tagEnd], i-tagStart)
	if name == "" {
		return "SafeText"
	}
	return ir.AttrContext(name)
}

func isMetaRefresh(tag string) bool {
	name := strings.TrimSpace(ir.AttrValue(tag, "http-equiv"))
	return strings.EqualFold(name, "refresh")
}
