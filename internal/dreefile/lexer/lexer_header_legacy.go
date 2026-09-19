package lexer

import (
	"fmt"
	"strings"
)

func legacyHeaderError(line string, index int, keyword, replacement string) error {
	col := strings.Index(line, keyword) + 1
	return &HeaderError{
		Line: index + 1,
		Col:  col,
		Err:  fmt.Errorf("legacy %s header is no longer supported: %s", keyword, replacement),
	}
}
