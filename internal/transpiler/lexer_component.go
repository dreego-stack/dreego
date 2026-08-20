package transpiler

import "strings"

func scanComponentTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	prefix := "<@"
	if strings.HasPrefix(remaining, "</@") {
		prefix = "</@"
	}

	remaining = remaining[len(prefix):]

	end := tagEnd(remaining)
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
