package parser

import "strings"

// bodySectionAttrs returns the attributes on the <body> section tag that are
// not Dreego directives (lang, method). They are HTML attributes the author
// expects on the rendered <body>; layout rendering cannot apply them.
func bodySectionAttrs(attrs string) string {
	if attrs == "" {
		return ""
	}
	var kept []string
	for part := range strings.FieldsSeq(attrs) {
		if strings.HasPrefix(part, "lang=") || strings.HasPrefix(part, "method=") {
			continue
		}
		kept = append(kept, part)
	}
	return strings.Join(kept, " ")
}
