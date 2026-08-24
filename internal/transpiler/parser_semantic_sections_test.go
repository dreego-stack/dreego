package transpiler

import "testing"

func TestParseSemanticSectionLanguages(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		src  string
		want string
		read func(*File) string
	}{
		{name: "server default", src: `<server></server>`, want: "go", read: func(f *File) string { return f.Server[0].Language }},
		{name: "server explicit", src: `<server lang="go"></server>`, want: "go", read: func(f *File) string { return f.Server[0].Language }},
		{name: "head default", src: `<head></head>`, want: "html", read: func(f *File) string { return f.Head.Language }},
		{name: "body explicit", src: `<body lang="html"></body>`, want: "html", read: func(f *File) string { return f.Body.Language }},
		{name: "style default", src: `<style></style>`, want: "css", read: func(f *File) string { return f.Style.Language }},
		{name: "client explicit", src: `<client lang="js"></client>`, want: "js", read: func(f *File) string { return f.Client.Language }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := Lex(tc.src)
			if err != nil {
				t.Fatalf("Lex: %v", err)
			}
			file, err := NewParser(tokens).Parse()
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if got := tc.read(file); got != tc.want {
				t.Fatalf("language = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseRejectsUnsupportedBuiltInSectionLanguages(t *testing.T) {
	t.Parallel()
	for _, src := range []string{
		`<server lang="lua"></server>`,
		`<head lang="markdown"></head>`,
		`<body lang="markdown"></body>`,
		`<style lang="scss"></style>`,
		`<client lang="ts"></client>`,
	} {
		tokens, err := Lex(src)
		if err != nil {
			t.Fatalf("Lex(%q): %v", src, err)
		}
		if _, err := NewParser(tokens).Parse(); err == nil {
			t.Fatalf("Parse(%q) succeeded, want unsupported-language error", src)
		}
	}
}
