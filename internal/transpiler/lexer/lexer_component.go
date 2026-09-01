package lexer

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func scanComponentTag(input string, pos *int) tokens.Token {
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
		return tokens.Token{Type: tokens.TokenText, Value: input[start:], Pos: start}
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
		return tokens.Token{Type: tokens.TokenComponentTagClose, Tag: tag, Pos: start}
	}
	if selfClose {
		return tokens.Token{Type: tokens.TokenComponentSelfClose, Tag: tag, Attr: attrs, Pos: start}
	}
	return tokens.Token{Type: tokens.TokenComponentTagOpen, Tag: tag, Attr: attrs, Pos: start}
}
