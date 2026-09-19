package lexer

import "strings"

func isTemplateBrace(input string) bool {
	for _, prefix := range []string{
		"{{", "{#else", "{/if}", "{/each}", "{#if ", "{#each ",
		"{/slot}", "{#slot", "{#verbatim}",
	} {
		if strings.HasPrefix(input, prefix) {
			return true
		}
	}
	return false
}
