package lexer

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func scanTag(input string, pos *int) tokens.Token {
	start := *pos
	remaining := input[start:]

	if strings.HasPrefix(remaining, "<@") || strings.HasPrefix(remaining, "</@") {
		return scanComponentTag(input, pos)
	}

	knownTags := []string{"server", "head", "body", "style", "client", "go", "div", "script"}

	for _, tag := range knownTags {
		closer := "</" + tag + ">"
		if strings.HasPrefix(remaining, closer) {
			*pos += len(closer)
			return tokens.Token{Type: tokens.TokenTagClose, Tag: tag, Pos: start}
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
			end := tagEnd(remaining)
			if end < 0 {
				*pos += len(remaining)
				return tokens.Token{Type: tokens.TokenText, Value: remaining, Pos: start}
			}
			attrs := strings.TrimSpace(remaining[len(opener):end])
			*pos += end + 1
			selfClose := strings.HasSuffix(attrs, "/")
			attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
			return tokens.Token{Type: tokens.TokenTagOpen, Tag: tag, Attr: attrs, SelfClose: selfClose, Pos: start}
		}
	}

	if strings.HasPrefix(remaining, "</") {
		end := tagEnd(remaining)
		if end < 0 {
			*pos += len(remaining)
			return tokens.Token{Type: tokens.TokenText, Value: remaining, Pos: start}
		}
		tag := remaining[2:end]
		if idx := strings.IndexByte(tag, ' '); idx >= 0 {
			tag = tag[:idx]
		}
		*pos += end + 1
		return tokens.Token{Type: tokens.TokenTagClose, Tag: tag, Pos: start}
	}

	if remaining[0] == '<' {
		end := tagEnd(remaining)
		if end < 0 {
			*pos += len(remaining)
			return tokens.Token{Type: tokens.TokenText, Value: remaining, Pos: start}
		}
		body := remaining[1:end]
		tag := body
		if idx := strings.IndexByte(body, ' '); idx >= 0 {
			tag = body[:idx]
		}
		attrs := strings.TrimSpace(strings.TrimPrefix(body, tag))
		selfClose := strings.HasSuffix(attrs, "/")
		attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
		*pos += end + 1
		return tokens.Token{Type: tokens.TokenTagOpen, Tag: strings.TrimSpace(tag), Attr: attrs, SelfClose: selfClose, Pos: start}
	}

	end := tagEnd(remaining)
	if end >= 0 {
		*pos += end + 1
		return tokens.Token{Type: tokens.TokenText, Value: remaining[:end+1], Pos: start}
	}
	*pos += len(remaining)
	return tokens.Token{Type: tokens.TokenText, Value: remaining, Pos: start}
}
