package tokens

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
	TokenElse
	TokenElseIf
	TokenEachOpen
	TokenEachClose
	TokenEachElse
	TokenSlot
	TokenSlotOpen
	TokenSlotClose
	TokenComponentTagOpen
	TokenComponentTagClose
	TokenComponentSelfClose
	TokenVerbatim
)

type Token struct {
	Type      TokenType
	Value     string
	Tag       string
	Attr      string
	Pos       int
	SelfClose bool
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
	case TokenElse:
		return "Else"
	case TokenElseIf:
		return "ElseIf"
	case TokenEachOpen:
		return "EachOpen"
	case TokenEachClose:
		return "EachClose"
	case TokenEachElse:
		return "EachElse"
	case TokenVerbatim:
		return "Verbatim"
	}
	return "UNKNOWN"
}
