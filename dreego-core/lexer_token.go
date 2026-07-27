package core

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF

	TokenTagOpen
	TokenTagClose
	TokenText
	TokenExpression
	TokenIfOpen
	TokenIfClose
	TokenEachOpen
	TokenEachClose
	TokenSlot
	TokenComponentHeader
	TokenImport
	TokenComponentTagOpen
	TokenComponentTagClose
	TokenComponentSelfClose
)

type Token struct {
	Type  TokenType
	Value string
	Tag   string
	Attr  string
	Pos   int
}

func (t TokenType) String() string {
	switch t {
	case TokenError:
		return "ERROR"
	case TokenEOF:
		return "EOF"
	case TokenTagOpen:
		return "TagOpen"
	case TokenTagClose:
		return "TagClose"
	case TokenText:
		return "Text"
	case TokenExpression:
		return "Expression"
	case TokenIfOpen:
		return "IfOpen"
	case TokenIfClose:
		return "IfClose"
	case TokenEachOpen:
		return "EachOpen"
	case TokenEachClose:
		return "EachClose"
	}
	return "UNKNOWN"
}
