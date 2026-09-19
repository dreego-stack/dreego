package lexer

import "fmt"

type HeaderError struct {
	Line int
	Col  int
	Err  error
}

func (e *HeaderError) Error() string {
	return fmt.Sprintf("%d:%d: %v", e.Line, e.Col, e.Err)
}

func (e *HeaderError) Unwrap() error {
	return e.Err
}
