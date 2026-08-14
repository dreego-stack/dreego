package core

// tagEnd returns the index of the first '>' that is not inside a quoted
// attribute value, or -1 if the tag is unclosed.
func tagEnd(input string) int {
	var quote byte
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '"', '\'':
			if quote == 0 {
				quote = input[i]
			} else if quote == input[i] {
				quote = 0
			}
		case '>':
			if quote == 0 {
				return i
			}
		}
	}
	return -1
}
