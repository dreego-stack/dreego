package dreefile

import (
	"strings"
	"testing"
)

type fmtCase struct {
	name string
	src  string
}

func fmtRoundTripCases() []fmtCase {
	return []fmtCase{
		{
			name: "scaffold web-minimal layout",
			src: `<body>
<!DOCTYPE html>
<html lang="en">
<head>
    <title>App</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {#head}
</head>
<body>
    <a href="#main" class="skip-link">skip to content</a>
    <main id="main">
    {#slot}
    </main>
</body>
</html>
</body>

<style>
.skip-link { position: absolute; left: -9999px; }
</style>
`,
		},
		{
			name: "scaffold web-app layout",
			src: `<body>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {#head}
</head>
<body>
    <a href="#main" class="skip-link">skip to content</a>
    <header class="site-header">
        <span class="brand">App</span>
        <@Nav current={c.Get("current")}/>
    </header>
    <main id="main">
    {#slot}
    </main>
    <footer class="site-footer">
        <p>App — built with Dreego.</p>
    </footer>
</body>
</html>
</body>

<style>
:root { color-scheme: light; }
</style>
`,
		},
		{
			name: "feedback body-level layout",
			src: `DREEFILE layout

<body>
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
		},
		{
			name: "document skeleton without doctype",
			src: `<body>
<html>
<head><title>x</title></head>
<body>{#slot}</body>
</html>
</body>
`,
		},
		{
			name: "body-level head with trailing style",
			src: `<body>
<html>
<head><meta charset="utf-8">{#head}</head>
<body><main>{#slot}</main></body>
</html>
</body>

<style>
body { margin: 0; }
</style>
`,
		},
		{
			name: "route with all sections and prose",
			src: `<server>
msg := "hi"
</server>

<head><title>t</title></head>

<body><p>{{  msg  }}</p>{#if msg}<span>x</span>{/if}</body>

<style>.a { color: red; }</style>
`,
		},
		{
			name: "prose around a single body",
			src: `Intro text

<body><p>x</p></body>

Outro text
`,
		},
		{
			name: "nested body with inner style",
			src: `<body>
<section>
<style>.b { color: blue; }</style>
<p>nested</p>
</section>
</body>
`,
		},
	}
}

func TestFormatRoundTripPreservesSemantics(t *testing.T) {
	for _, tc := range fmtRoundTripCases() {
		t.Run(tc.name, func(t *testing.T) {
			inSig, err := fmtTokenSignature(tc.src)
			if err != nil {
				t.Fatalf("input does not lex: %v", err)
			}
			out := Format(tc.src)
			if again := Format(out); again != out {
				t.Fatalf("Format is not idempotent\nonce:  %q\ntwice: %q", out, again)
			}
			outSig, err := fmtTokenSignature(out)
			if err != nil {
				t.Fatalf("formatted output does not lex: %v\n%s", err, out)
			}
			if outSig != inSig {
				t.Fatalf("Format changed the section structure\nin:  %s\nout: %s\nformatted:\n%s", inSig, outSig, out)
			}
			if got, want := fmtTextSignature(out), fmtTextSignature(tc.src); got != want {
				t.Fatalf("Format changed section text\nwant: %q\ngot:  %q\nformatted:\n%s", want, got, out)
			}
			for _, tag := range []string{"</html>", "</body>", "</head>"} {
				if strings.Count(out, tag) != strings.Count(tc.src, tag) {
					t.Fatalf("Format changed the number of %q\ngot:\n%s", tag, out)
				}
			}
		})
	}
}

func TestFormatBodyLevelLayoutKeepsDocumentSkeleton(t *testing.T) {
	src := `DREEFILE layout

<body>
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
`
	out := Format(src)

	outer := strings.Index(out, "<body>")
	innerHead := strings.Index(out, "<head>")
	closeHTML := strings.Index(out, "</html>")
	if outer < 0 || innerHead < 0 || closeHTML < 0 {
		t.Fatalf("document skeleton lost, got:\n%s", out)
	}
	if innerHead < outer {
		t.Fatalf("the document <head> must stay inside the outer <body>, got:\n%s", out)
	}
	if outer > innerHead || innerHead > closeHTML {
		t.Fatalf("skeleton order broken (body=%d head=%d /html=%d), got:\n%s", outer, innerHead, closeHTML, out)
	}
	if strings.HasPrefix(strings.TrimSpace(strings.TrimPrefix(out, "DREEFILE layout")), "<head>") {
		t.Fatalf("the document <head> must not be hoisted to a root section, got:\n%s", out)
	}
	if !strings.Contains(out, "{#slot}") || !strings.Contains(out, "<main") {
		t.Fatalf("slot/main lost, got:\n%s", out)
	}
}

func TestFormatPreservesServerCodeExact(t *testing.T) {
	src := "<server>\nmsg   :=   compute()\n\nif  msg != \"\"  {\n    log(msg)\n}\n</server>\n\n<body><p>{{ msg }}</p></body>\n"
	out := Format(src)
	want := "msg   :=   compute()\n\nif  msg != \"\"  {\n    log(msg)\n}"
	if !strings.Contains(out, want) {
		t.Fatalf("server section content must not be rewritten\nwant substring: %q\ngot:\n%s", want, out)
	}
}

func fmtTokenSignature(src string) (string, error) {
	toks, err := Lex(src)
	if err != nil {
		return "", err
	}
	var parts []string
	for _, tok := range toks {
		switch tok.Type {
		case TokenTagOpen, TokenTagClose:
			parts = append(parts, tok.Type.String()+":"+tok.Tag)
		case TokenIfOpen, TokenIfClose, TokenEachOpen, TokenEachClose,
			TokenEachElse, TokenElse, TokenElseIf,
			TokenSlot, TokenSlotOpen, TokenSlotClose,
			TokenComponentTagOpen, TokenComponentTagClose, TokenComponentSelfClose,
			TokenVerbatim, TokenExpression, TokenMessage:
			parts = append(parts, tok.Type.String())
		}
	}
	return strings.Join(parts, " "), nil
}

func fmtTextSignature(src string) string {
	toks, err := Lex(src)
	if err != nil {
		return "LEX_ERROR: " + err.Error()
	}
	var text strings.Builder
	for _, tok := range toks {
		if tok.Type == TokenText {
			text.WriteString(tok.Value)
			text.WriteString(" ")
		}
	}
	return strings.Join(strings.Fields(text.String()), " ")
}

func TestFmtTokenSignatureDetectsDroppedCloseTags(t *testing.T) {
	full := "<body>\n<html>\n<head><title>x</title></head>\n<body>x</body>\n</html>\n</body>\n"
	broken := "<head><title>x</title></head>\n\n<body>\n<html>\n<body>x</body>\n</body>\n"
	a, err := fmtTokenSignature(full)
	if err != nil {
		t.Fatalf("lex full: %v", err)
	}
	b, err := fmtTokenSignature(broken)
	if err != nil {
		t.Fatalf("lex broken: %v", err)
	}
	if a == b {
		t.Fatal("signature must distinguish a hoisted head and dropped close tags")
	}
}
