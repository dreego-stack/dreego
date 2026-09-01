package transpiler

import "github.com/dreego-stack/dreego/internal/transpiler/tokens"

type TokenType = tokens.TokenType
type Token = tokens.Token

const (
	TokenError = tokens.TokenError
	TokenEOF   = tokens.TokenEOF

	TokenTagOpen           = tokens.TokenTagOpen
	TokenTagClose          = tokens.TokenTagClose
	TokenText              = tokens.TokenText
	TokenExpression        = tokens.TokenExpression
	TokenIfOpen            = tokens.TokenIfOpen
	TokenIfClose           = tokens.TokenIfClose
	TokenElse              = tokens.TokenElse
	TokenElseIf            = tokens.TokenElseIf
	TokenEachOpen          = tokens.TokenEachOpen
	TokenEachClose         = tokens.TokenEachClose
	TokenEachElse          = tokens.TokenEachElse
	TokenSlot              = tokens.TokenSlot
	TokenSlotOpen          = tokens.TokenSlotOpen
	TokenSlotClose         = tokens.TokenSlotClose
	TokenComponentTagOpen  = tokens.TokenComponentTagOpen
	TokenComponentTagClose = tokens.TokenComponentTagClose
	TokenComponentSelfClose = tokens.TokenComponentSelfClose
	TokenVerbatim          = tokens.TokenVerbatim
)
