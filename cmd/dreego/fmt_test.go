package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
	dreefile "github.com/dreego-stack/dreego/internal/dreefile"
)

var fmtSkeletonTags = []string{"<head>", "</head>", "<body>", "</body>", "</html>"}

func TestFmtScaffoldLayoutsRoundTrip(t *testing.T) {
	for _, meta := range templates.List() {
		t.Run(meta.Name, func(t *testing.T) {
			dir := t.TempDir()
			if err := templates.Install(dir, "example.com/app", meta.Name); err != nil {
				t.Fatalf("install %q: %v", meta.Name, err)
			}
			layouts, err := filepath.Glob(filepath.Join(dir, "www", "layouts", "*.dreego"))
			if err != nil || len(layouts) == 0 {
				t.Fatalf("no scaffold layouts for %q: %v", meta.Name, err)
			}
			for _, path := range layouts {
				src, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertFmtPreservesSkeleton(t, path, string(src))
			}
		})
	}
}

func assertFmtPreservesSkeleton(t *testing.T, name, src string) {
	t.Helper()
	out := dreefile.Format(src)
	if again := dreefile.Format(out); again != out {
		t.Fatalf("%s: Format is not idempotent\nonce:  %q\ntwice: %q", name, out, again)
	}
	inSig := fmtSignature(t, name, src)
	outSig := fmtSignature(t, name, out)
	if inSig != outSig {
		t.Fatalf("%s: Format changed the section structure\nin:  %s\nout: %s\nformatted:\n%s", name, inSig, outSig, out)
	}
	for _, tag := range fmtSkeletonTags {
		if strings.Count(out, tag) != strings.Count(src, tag) {
			t.Fatalf("%s: Format changed the number of %q\ngot:\n%s", name, tag, out)
		}
	}
}

func fmtSignature(t *testing.T, name, src string) string {
	t.Helper()
	toks, err := dreefile.Lex(src)
	if err != nil {
		t.Fatalf("%s: lex: %v", name, err)
	}
	var parts []string
	for _, tok := range toks {
		switch tok.Type {
		case dreefile.TokenTagOpen, dreefile.TokenTagClose:
			parts = append(parts, tok.Type.String()+":"+tok.Tag)
		case dreefile.TokenIfOpen, dreefile.TokenIfClose, dreefile.TokenEachOpen,
			dreefile.TokenEachClose, dreefile.TokenEachElse, dreefile.TokenElse,
			dreefile.TokenElseIf, dreefile.TokenSlot, dreefile.TokenSlotOpen,
			dreefile.TokenSlotClose, dreefile.TokenComponentTagOpen,
			dreefile.TokenComponentTagClose, dreefile.TokenComponentSelfClose,
			dreefile.TokenVerbatim, dreefile.TokenExpression, dreefile.TokenMessage:
			parts = append(parts, tok.Type.String())
		}
	}
	return strings.Join(parts, " ")
}

func TestFmtBodyLevelLayoutVariantsRoundTrip(t *testing.T) {
	variants := map[string]string{
		"doctype+html": `<body>
<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="utf-8">
    {#head}
</head>
<body>
    <main id="main">{#slot}</main>
</body>
</html>
</body>
`,
		"html only": `<body>
<html>
<head><title>x</title></head>
<body>{#slot}</body>
</html>
</body>
`,
		"with trailing style": `<body>
<html>
<head>{#head}</head>
<body>{#slot}</body>
</html>
</body>

<style>
body { margin: 0; }
</style>
`,
		"prose around body": "Intro\n\n<body><p>x</p></body>\n\nOutro\n",
	}
	for name, src := range variants {
		t.Run(name, func(t *testing.T) {
			assertFmtPreservesSkeleton(t, name, src)
		})
	}
}

func TestFmtCheckNeverWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "messy.dreego")
	src := "<head>\n  <title> t </title>\n</head>\n\n\n<body>\n  <p>{{  msg  }}</p>\n</body>\n"
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := cmdFmtE([]string{"--check", path}, &out); err == nil {
		t.Fatal("--check must report an error for a file that needs formatting")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Fatalf("--check modified the file\nbefore: %q\nafter:  %q", src, after)
	}
}

func TestFmtCheckNeverWritesBodyLevelLayout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "default.dreego")
	src := "DREEFILE layout\n\n\n<body>\n<!DOCTYPE html>\n<html lang=\"de\">\n<head>\n    <meta charset=\"utf-8\">\n    {#head}\n</head>\n<body>\n    <main id=\"main\">{#slot}</main>\n</body>\n</html>\n</body>\n"
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	if err := cmdFmtE([]string{"--check", path}, &bytes.Buffer{}); err == nil {
		t.Fatal("--check must report an error for a layout that needs formatting")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Fatalf("--check destroyed the layout\nbefore: %q\nafter:  %q", src, after)
	}

	formatted := dreefile.Format(src)
	if again := dreefile.Format(formatted); again != formatted {
		t.Fatalf("body-level layout is not idempotent\nonce:  %q\ntwice: %q", formatted, again)
	}
	if _, err := dreefile.Lex(formatted); err != nil {
		t.Fatalf("formatted layout must still lex: %v\n%s", err, formatted)
	}
	if !strings.Contains(formatted, "</html>") || !strings.Contains(formatted, "<head>") {
		t.Fatalf("formatted layout lost its document skeleton:\n%s", formatted)
	}
}
