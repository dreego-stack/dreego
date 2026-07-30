package core

import (
	"fmt"
	"strings"
)

func Lex(input string) ([]Token, error) {
	var tokens []Token
	pos := 0
	sectionStack := []string{}
	sectionTags := map[string]bool{"go": true, "head": true, "script": true, "style": true}

	for pos < len(input) {
		inSection := len(sectionStack) > 0

		nextPos := -1
		nextCh := byte(0)

		for i := pos; i < len(input); i++ {
			if input[i] == '<' {
				nextPos = i
				nextCh = '<'
				break
			}
			if !inSection && input[i] == '{' {
				nextPos = i
				nextCh = '{'
				break
			}
		}

		if nextPos < 0 {
			nextPos = len(input)
		}

		if nextPos > pos {
			tokens = append(tokens, Token{
				Type:  TokenText,
				Value: input[pos:nextPos],
				Pos:   pos,
			})
		}

		if nextPos >= len(input) {
			break
		}

		pos = nextPos

		if nextCh == '<' {
			tok := scanTag(input, &pos)
			tokens = append(tokens, tok)

			if tok.Type == TokenTagOpen && sectionTags[tok.Tag] {
				sectionStack = append(sectionStack, tok.Tag)
			}
			if tok.Type == TokenTagClose && sectionTags[tok.Tag] {
				if len(sectionStack) == 0 {
					return nil, fmt.Errorf("unexpected closing tag </%s> at position %d", tok.Tag, tok.Pos)
				}
				if sectionStack[len(sectionStack)-1] != tok.Tag {
					return nil, fmt.Errorf("mismatched closing tag </%s>, expected </%s> at position %d",
						tok.Tag, sectionStack[len(sectionStack)-1], tok.Pos)
				}
				sectionStack = sectionStack[:len(sectionStack)-1]
			}
		} else if nextCh == '{' {
			tok, err := scanBrace(input, &pos)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		}
	}

	if len(sectionStack) > 0 {
		return nil, fmt.Errorf("unclosed tag <%s>", sectionStack[len(sectionStack)-1])
	}

	tokens = append(tokens, Token{Type: TokenEOF, Pos: pos})
	return tokens, nil
}

func ParseHeader(input string) (comp *ComponentDef, imports []Import, body string) {
	lines := strings.Split(input, "\n")
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])

		if strings.HasPrefix(trimmed, "Component ") {
			comp = parseComponentHeader(trimmed)
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "import ") {
			imp := parseImportLine(trimmed)
			if imp != nil {
				imports = append(imports, *imp)
			}
			i++
			continue
		}

		if trimmed == "" {
			i++
			continue
		}

		break
	}

	body = strings.Join(lines[i:], "\n")
	return
}

func parseComponentHeader(line string) *ComponentDef {
	line = strings.TrimPrefix(line, "Component ")
	openParen := strings.IndexByte(line, '(')
	if openParen < 0 {
		return &ComponentDef{Name: strings.TrimSpace(line)}
	}
	name := strings.TrimSpace(line[:openParen])
	rest := line[openParen:]

	closeParen := strings.LastIndexByte(rest, ')')
	if closeParen < 0 {
		return &ComponentDef{Name: name}
	}

	params := strings.TrimSpace(rest[1:closeParen])

	comp := &ComponentDef{Name: name}
	if strings.Contains(params, "(") || !strings.Contains(params, ",") {
		comp.Props = parseProps(params)
	} else {
		parts := strings.SplitN(params, ") (", 2)
		if len(parts) > 0 && parts[0] != "" {
			comp.Props = parseProps(strings.Trim(parts[0], "() "))
		}
		if len(parts) > 1 && parts[1] != "" {
			for _, s := range strings.Split(strings.Trim(parts[1], "() "), ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					comp.Slots = append(comp.Slots, s)
				}
			}
		}
	}

	return comp
}

func parseProps(s string) []Prop {
	var props []Prop
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		p := Prop{Name: fields[0]}
		if len(fields) >= 2 {
			p.Type = fields[1]
		}
		if len(fields) >= 4 && fields[2] == "=" {
			p.Default = fields[3]
		}
		if p.Type == "" {
			p.Type = "string"
		}
		props = append(props, p)
	}
	return props
}

func parseImportLine(line string) *Import {
	line = strings.TrimPrefix(line, "import ")
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	imp := &Import{Path: strings.Trim(fields[len(fields)-1], "\"")}
	if len(fields) >= 3 {
		imp.Alias = fields[0]
	}
	return imp
}

func scanTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	if strings.HasPrefix(remaining, "<@") || strings.HasPrefix(remaining, "</@") {
		return scanComponentTag(input, pos)
	}

	knownTags := []string{"go", "div", "head", "script", "style"}

	for _, tag := range knownTags {
		closer := "</" + tag + ">"
		if strings.HasPrefix(remaining, closer) {
			*pos += len(closer)
			return Token{Type: TokenTagClose, Tag: tag, Pos: start}
		}
		opener := "<" + tag
		if strings.HasPrefix(remaining, opener) {
			next := byte(0)
			if len(remaining) > len(opener) {
				next = remaining[len(opener)]
			}
			if next != ' ' && next != '>' && next != '/' {
				continue
			}
			end := strings.IndexByte(remaining, '>')
			if end < 0 {
				*pos += len(remaining)
				return Token{Type: TokenText, Value: remaining, Pos: start}
			}
			attrs := strings.TrimSpace(remaining[len(opener):end])
			*pos += end + 1
			return Token{Type: TokenTagOpen, Tag: tag, Attr: attrs, Pos: start}
		}
	}

	if strings.HasPrefix(remaining, "</") {
		end := strings.IndexByte(remaining, '>')
		if end < 0 {
			*pos += len(remaining)
			return Token{Type: TokenText, Value: remaining, Pos: start}
		}
		tag := remaining[2:end]
		if idx := strings.IndexByte(tag, ' '); idx >= 0 {
			tag = tag[:idx]
		}
		*pos += end + 1
		return Token{Type: TokenTagClose, Tag: tag, Pos: start}
	}

	if remaining[0] == '<' {
		end := strings.IndexByte(remaining, '>')
		if end < 0 {
			*pos += len(remaining)
			return Token{Type: TokenText, Value: remaining, Pos: start}
		}
		body := remaining[1:end]
		tag := body
		if idx := strings.IndexByte(body, ' '); idx >= 0 {
			tag = body[:idx]
		}
		attrs := strings.TrimSpace(strings.TrimPrefix(body, tag))
		attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
		*pos += end + 1
		return Token{Type: TokenTagOpen, Tag: strings.TrimSpace(tag), Attr: attrs, Pos: start}
	}

	end := strings.IndexByte(remaining, '>')
	if end >= 0 {
		*pos += end + 1
		return Token{Type: TokenText, Value: remaining[:end+1], Pos: start}
	}
	*pos += len(remaining)
	return Token{Type: TokenText, Value: remaining, Pos: start}
}

func scanComponentTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	prefix := "<@"
	if strings.HasPrefix(remaining, "</@") {
		prefix = "</@"
	}

	remaining = remaining[len(prefix):]

	end := strings.IndexByte(remaining, '>')
	if end < 0 {
		*pos += len(input) - start
		return Token{Type: TokenText, Value: input[start:], Pos: start}
	}

	body := strings.TrimSpace(remaining[:end])
	selfClose := strings.HasSuffix(body, "/")
	if selfClose {
		body = strings.TrimSuffix(body, "/")
		body = strings.TrimSpace(body)
	}

	parts := strings.Fields(body)
	tag := ""
	attrs := ""
	if len(parts) > 0 {
		tag = parts[0]
		attrs = strings.TrimSpace(strings.TrimPrefix(body, tag))
	}

	*pos = start + len(prefix) + end + 1

	if prefix == "</@" {
		return Token{Type: TokenComponentTagClose, Tag: tag, Pos: start}
	}
	if selfClose {
		return Token{Type: TokenComponentSelfClose, Tag: tag, Attr: attrs, Pos: start}
	}
	return Token{Type: TokenComponentTagOpen, Tag: tag, Attr: attrs, Pos: start}
}
