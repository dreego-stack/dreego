package core

import "strings"

// ParseFrontmatter splits a leading YAML-like frontmatter block off src.
// The opening "---\n" and a closing "---" on its own line delimit the block
// (only at the very top of src). Each "key: value" line becomes a map entry;
// the first ':' splits key from value, so a ':' inside the value is preserved.
// A list value "tags: [go, web]" is normalized to the string "go, web".
// src without a leading frontmatter block yields nil and the whole src as body.
func ParseFrontmatter(src string) (frontmatter map[string]string, body string) {
	if !strings.HasPrefix(src, "---\n") {
		return nil, src
	}
	rest := src[len("---\n"):]
	end := 0
	if !strings.HasPrefix(rest, "---\n") {
		end = strings.Index(rest, "\n---")
		if end < 0 {
			return nil, src
		}
	}
	block := rest[:end]
	body = rest[end+len("\n---"):]
	body = strings.TrimPrefix(body, "\n")
	return parseBlock(block), body
}

func parseBlock(block string) map[string]string {
	fm := make(map[string]string)
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if key == "" {
			continue
		}
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		fm[key] = normalizeValue(val)
	}
	return fm
}

func normalizeValue(v string) string {
	if !strings.HasPrefix(v, "[") || !strings.HasSuffix(v, "]") {
		return v
	}
	inner := strings.TrimSpace(v[1 : len(v)-1])
	parts := strings.Split(inner, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, ", ")
}
