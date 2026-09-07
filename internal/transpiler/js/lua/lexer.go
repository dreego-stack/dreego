package lua

import (
	"fmt"
	"strconv"
	"unicode"
)

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenSeparator
	tokenIdentifier
	tokenNumber
	tokenString
	tokenLocal
	tokenIf
	tokenThen
	tokenElseIf
	tokenElse
	tokenEnd
	tokenAnd
	tokenOr
	tokenNot
	tokenTrue
	tokenFalse
	tokenNil
	tokenAssign
	tokenEqual
	tokenNotEqual
	tokenLess
	tokenLessEqual
	tokenGreater
	tokenGreaterEqual
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenPercent
	tokenConcat
	tokenDot
	tokenLeftParen
	tokenRightParen
	tokenComma
	tokenLeftBrace
	tokenRightBrace
)

type token struct {
	kind   tokenKind
	value  string
	line   int
	column int
}

var keywords = map[string]tokenKind{
	"local": tokenLocal, "if": tokenIf, "then": tokenThen,
	"elseif": tokenElseIf, "else": tokenElse, "end": tokenEnd,
	"and": tokenAnd, "or": tokenOr, "not": tokenNot,
	"true": tokenTrue, "false": tokenFalse, "nil": tokenNil,
}

func lex(source string) ([]token, error) {
	lexer := sourceLexer{source: []rune(source), line: 1, column: 1}
	return lexer.scan()
}

type sourceLexer struct {
	source []rune
	index  int
	line   int
	column int
}

func (l *sourceLexer) scan() ([]token, error) {
	var tokens []token
	for l.index < len(l.source) {
		current := l.source[l.index]
		if current == ' ' || current == '\t' || current == '\r' {
			l.advance()
			continue
		}
		if current == '\n' {
			l.advance()
			continue
		}
		if current == ';' {
			tokens = append(tokens, l.makeToken(tokenSeparator, ""))
			l.advance()
			continue
		}
		if current == '-' && l.peek(1) == '-' {
			for l.index < len(l.source) && l.source[l.index] != '\n' {
				l.advance()
			}
			continue
		}
		if unicode.IsLetter(current) || current == '_' {
			tokens = append(tokens, l.identifier())
			continue
		}
		if unicode.IsDigit(current) {
			tokens = append(tokens, l.number())
			continue
		}
		if current == '"' || current == '\'' {
			value, err := l.stringToken()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, value)
			continue
		}
		value, width, ok := l.operator()
		if !ok {
			return nil, fmt.Errorf("Lua %d:%d: unexpected character %q", l.line, l.column, current)
		}
		tokens = append(tokens, value)
		for range width {
			l.advance()
		}
	}
	tokens = append(tokens, token{kind: tokenEOF, line: l.line, column: l.column})
	return tokens, nil
}

func (l *sourceLexer) identifier() token {
	line, column, start := l.line, l.column, l.index
	for l.index < len(l.source) && (unicode.IsLetter(l.source[l.index]) || unicode.IsDigit(l.source[l.index]) || l.source[l.index] == '_') {
		l.advance()
	}
	value := string(l.source[start:l.index])
	kind := tokenIdentifier
	if keyword, ok := keywords[value]; ok {
		kind = keyword
	}
	return token{kind: kind, value: value, line: line, column: column}
}

func (l *sourceLexer) number() token {
	line, column, start := l.line, l.column, l.index
	for l.index < len(l.source) && unicode.IsDigit(l.source[l.index]) {
		l.advance()
	}
	if l.index < len(l.source) && l.source[l.index] == '.' && l.peek(1) != '.' {
		l.advance()
		for l.index < len(l.source) && unicode.IsDigit(l.source[l.index]) {
			l.advance()
		}
	}
	if l.index < len(l.source) && (l.source[l.index] == 'e' || l.source[l.index] == 'E') {
		l.advance()
		if l.index < len(l.source) && (l.source[l.index] == '+' || l.source[l.index] == '-') {
			l.advance()
		}
		for l.index < len(l.source) && unicode.IsDigit(l.source[l.index]) {
			l.advance()
		}
	}
	return token{kind: tokenNumber, value: string(l.source[start:l.index]), line: line, column: column}
}

func (l *sourceLexer) stringToken() (token, error) {
	line, column, quote := l.line, l.column, l.source[l.index]
	l.advance()
	var value []rune
	for l.index < len(l.source) && l.source[l.index] != quote {
		if l.source[l.index] == '\n' {
			return token{}, fmt.Errorf("Lua %d:%d: unclosed string", line, column)
		}
		if l.source[l.index] == '\\' && l.index+1 < len(l.source) {
			value = append(value, l.source[l.index], l.source[l.index+1])
			l.advance()
			l.advance()
			continue
		}
		value = append(value, l.source[l.index])
		l.advance()
	}
	if l.index == len(l.source) {
		return token{}, fmt.Errorf("Lua %d:%d: unclosed string", line, column)
	}
	l.advance()
	decoded, err := decodeString(string(value), byte(quote))
	if err != nil {
		return token{}, fmt.Errorf("Lua %d:%d: invalid string: %w", line, column, err)
	}
	return token{kind: tokenString, value: decoded, line: line, column: column}, nil
}

func decodeString(value string, quote byte) (string, error) {
	var decoded []rune
	for value != "" {
		current, _, tail, err := strconv.UnquoteChar(value, quote)
		if err != nil {
			return "", err
		}
		decoded = append(decoded, current)
		value = tail
	}
	return string(decoded), nil
}

func (l *sourceLexer) operator() (token, int, bool) {
	line, column := l.line, l.column
	pair := string([]rune{l.source[l.index], l.peek(1)})
	pairs := map[string]tokenKind{"==": tokenEqual, "~=": tokenNotEqual, "<=": tokenLessEqual, ">=": tokenGreaterEqual, "..": tokenConcat}
	if kind, ok := pairs[pair]; ok {
		return token{kind: kind, value: pair, line: line, column: column}, 2, true
	}
	singles := map[rune]tokenKind{'=': tokenAssign, '<': tokenLess, '>': tokenGreater, '+': tokenPlus, '-': tokenMinus, '*': tokenStar, '/': tokenSlash, '%': tokenPercent, '.': tokenDot, '(': tokenLeftParen, ')': tokenRightParen, ',': tokenComma, '{': tokenLeftBrace, '}': tokenRightBrace}
	kind, ok := singles[l.source[l.index]]
	return token{kind: kind, value: string(l.source[l.index]), line: line, column: column}, 1, ok
}

func (l *sourceLexer) makeToken(kind tokenKind, value string) token {
	return token{kind: kind, value: value, line: l.line, column: l.column}
}

func (l *sourceLexer) peek(offset int) rune {
	if l.index+offset >= len(l.source) {
		return 0
	}
	return l.source[l.index+offset]
}

func (l *sourceLexer) advance() {
	if l.source[l.index] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.index++
}
