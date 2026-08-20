package transpiler

import "strings"

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
			end := tagEnd(remaining)
			if end < 0 {
				*pos += len(remaining)
				return Token{Type: TokenText, Value: remaining, Pos: start}
			}
			attrs := strings.TrimSpace(remaining[len(opener):end])
			*pos += end + 1
			selfClose := strings.HasSuffix(attrs, "/")
			attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
			return Token{Type: TokenTagOpen, Tag: tag, Attr: attrs, SelfClose: selfClose, Pos: start}
		}
	}

	if strings.HasPrefix(remaining, "</") {
		end := tagEnd(remaining)
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
		end := tagEnd(remaining)
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
		selfClose := strings.HasSuffix(attrs, "/")
		attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
		*pos += end + 1
		return Token{Type: TokenTagOpen, Tag: strings.TrimSpace(tag), Attr: attrs, SelfClose: selfClose, Pos: start}
	}

	end := tagEnd(remaining)
	if end >= 0 {
		*pos += end + 1
		return Token{Type: TokenText, Value: remaining[:end+1], Pos: start}
	}
	*pos += len(remaining)
	return Token{Type: TokenText, Value: remaining, Pos: start}
}
