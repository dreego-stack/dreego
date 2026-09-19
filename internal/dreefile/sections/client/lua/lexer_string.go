package lua

import (
	"fmt"
	"unicode/utf8"
)

func (l *sourceLexer) stringToken() (token, error) {
	line, column, quote := l.line, l.column, l.source[l.index]
	l.advance()
	var value []rune
	for l.index < len(l.source) && l.source[l.index] != quote {
		if l.source[l.index] == '\n' {
			return token{}, fmt.Errorf("Lua %d:%d: unclosed string", line, column)
		}
		if l.source[l.index] == '\\' {
			escaped, err := l.escape(line, column)
			if err != nil {
				return token{}, err
			}
			value = append(value, escaped...)
			continue
		}
		value = append(value, l.source[l.index])
		l.advance()
	}
	if l.index == len(l.source) {
		return token{}, fmt.Errorf("Lua %d:%d: unclosed string", line, column)
	}
	l.advance()
	return token{kind: tokenString, value: string(value), line: line, column: column}, nil
}

func (l *sourceLexer) escape(line, column int) ([]rune, error) {
	l.advance()
	if l.index >= len(l.source) {
		return nil, fmt.Errorf("Lua %d:%d: unfinished escape sequence", line, column)
	}
	current := l.source[l.index]
	simple := map[rune]rune{'a': '\a', 'b': '\b', 'f': '\f', 'n': '\n', 'r': '\r', 't': '\t', 'v': '\v', '\\': '\\', '"': '"', '\'': '\''}
	if decoded, ok := simple[current]; ok {
		l.advance()
		return []rune{decoded}, nil
	}
	if current == '\n' || current == '\r' {
		first := current
		l.advance()
		if l.index < len(l.source) && (first == '\r' && l.source[l.index] == '\n' || first == '\n' && l.source[l.index] == '\r') {
			l.advance()
		}
		return []rune{'\n'}, nil
	}
	if current == 'z' {
		l.advance()
		for l.index < len(l.source) && isLuaSpace(l.source[l.index]) {
			l.advance()
		}
		return nil, nil
	}
	if current == 'x' {
		l.advance()
		return l.fixedEscape(line, column, 16, 2, "hexadecimal")
	}
	if current == 'u' {
		return l.unicodeEscape(line, column)
	}
	if current >= '0' && current <= '9' {
		return l.decimalEscape(line, column)
	}
	return nil, fmt.Errorf("Lua %d:%d: invalid escape sequence \\%c", line, column, current)
}

func isLuaSpace(value rune) bool {
	return value == ' ' || value == '\f' || value == '\n' || value == '\r' || value == '\t' || value == '\v'
}

func (l *sourceLexer) fixedEscape(line, column, base, width int, name string) ([]rune, error) {
	value := 0
	for range width {
		if l.index >= len(l.source) {
			return nil, fmt.Errorf("Lua %d:%d: incomplete %s escape", line, column, name)
		}
		digit := digitValue(l.source[l.index])
		if digit < 0 || digit >= base {
			return nil, fmt.Errorf("Lua %d:%d: invalid %s escape", line, column, name)
		}
		value = value*base + digit
		l.advance()
	}
	return []rune{rune(value)}, nil
}

func (l *sourceLexer) decimalEscape(line, column int) ([]rune, error) {
	value := 0
	for count := 0; count < 3 && l.index < len(l.source) && l.source[l.index] >= '0' && l.source[l.index] <= '9'; count++ {
		value = value*10 + int(l.source[l.index]-'0')
		l.advance()
	}
	if value > 255 {
		return nil, fmt.Errorf("Lua %d:%d: decimal escape exceeds 255", line, column)
	}
	return []rune{rune(value)}, nil
}

func (l *sourceLexer) unicodeEscape(line, column int) ([]rune, error) {
	l.advance()
	if l.index >= len(l.source) || l.source[l.index] != '{' {
		return nil, fmt.Errorf("Lua %d:%d: expected { after \\u", line, column)
	}
	l.advance()
	value, digits := 0, 0
	for l.index < len(l.source) && l.source[l.index] != '}' {
		if digits == 6 {
			return nil, fmt.Errorf("Lua %d:%d: Unicode escape exceeds the browser range", line, column)
		}
		digit := digitValue(l.source[l.index])
		if digit < 0 || digit >= 16 {
			return nil, fmt.Errorf("Lua %d:%d: invalid Unicode escape", line, column)
		}
		value, digits = value*16+digit, digits+1
		l.advance()
	}
	if digits == 0 || l.index >= len(l.source) {
		return nil, fmt.Errorf("Lua %d:%d: invalid Unicode escape", line, column)
	}
	l.advance()
	if value > utf8.MaxRune || value >= 0xD800 && value <= 0xDFFF {
		return nil, fmt.Errorf("Lua %d:%d: invalid Unicode code point", line, column)
	}
	return []rune{rune(value)}, nil
}

func digitValue(value rune) int {
	switch {
	case value >= '0' && value <= '9':
		return int(value - '0')
	case value >= 'a' && value <= 'f':
		return int(value-'a') + 10
	case value >= 'A' && value <= 'F':
		return int(value-'A') + 10
	default:
		return -1
	}
}
