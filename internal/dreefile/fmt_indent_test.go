package dreefile

import (
	"strings"
	"testing"
)

// HTML sections are indented by nesting depth with four spaces per level, so a
// nested element is no longer flattened to column zero.
func TestFormatIndentsNestedHTMLByDepth(t *testing.T) {
	in := "<body>\n<span>\n<div>\ntext\n</div>\n</span>\n</body>\n"
	out := Format(in)
	want := "<body>\n    <span>\n        <div>\n            text\n        </div>\n    </span>\n</body>"
	if !strings.Contains(out, want) {
		t.Fatalf("html not indented by nesting depth\nwant:\n%s\ngot:\n%s", want, out)
	}
}

func TestFormatIndentsNestedHTMLIdempotent(t *testing.T) {
	in := "<body>\n<span><div>x</div></span>\n</body>\n"
	once := Format(in)
	if twice := Format(once); twice != once {
		t.Fatalf("not idempotent\nonce:  %q\ntwice: %q", once, twice)
	}
}

func TestFormatHTMLClosingTagAlignsWithOpening(t *testing.T) {
	in := "<body>\n<section>\n<p>x</p>\n</section>\n</body>\n"
	out := Format(in)
	want := "<body>\n    <section>\n        <p>x</p>\n    </section>\n</body>"
	if !strings.Contains(out, want) {
		t.Fatalf("closing tag must align with its opening tag\nwant:\n%s\ngot:\n%s", want, out)
	}
}

func TestFormatKeepsInlineBodySingleLine(t *testing.T) {
	in := "<body><p>hi</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "<body><p>hi</p></body>") {
		t.Fatalf("single-line body must stay inline, got:\n%s", out)
	}
}

// A server section is Go source. Its content is shifted to one indent level
// under the root tag while the code's own relative indentation is preserved.
func TestFormatIndentsServerBlock(t *testing.T) {
	in := "<server>\nnum := int1 + int2\n</server>\n\n<body><p>{{ num }}</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "<server>\n    num := int1 + int2\n</server>") {
		t.Fatalf("server block must be indented one level, got:\n%s", out)
	}
}

func TestFormatPreservesServerRelativeIndent(t *testing.T) {
	in := "<server>\nif x {\n    log(x)\n}\n</server>\n\n<body><p>ok</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "    if x {\n        log(x)\n    }") {
		t.Fatalf("server inner indentation must stay relative to the block, got:\n%s", out)
	}
}

func TestFormatServerIndentIdempotent(t *testing.T) {
	in := "<server>\nnum := 1\n</server>\n\n<body><p>ok</p></body>\n"
	once := Format(in)
	if twice := Format(once); twice != once {
		t.Fatalf("not idempotent\nonce:  %q\ntwice: %q", once, twice)
	}
}

func TestFormatServerIndentPreservesRawString(t *testing.T) {
	in := "<server>\nraw := `line1\n\n\nline3`\n</server>\n\n<body><p>ok</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "line1\n\n\nline3`") {
		t.Fatalf("raw string content must be preserved, got:\n%s", out)
	}
	if twice := Format(out); twice != out {
		t.Fatalf("not idempotent with a raw string\nonce:  %q\ntwice: %q", out, twice)
	}
}

func TestFormatIndentsStyleAndClientBlocks(t *testing.T) {
	in := "<body><p>ok</p></body>\n\n<style>\n.a { color: red; }\n</style>\n\n<client>\nconsole.log(\"hi\")\n</client>\n"
	out := Format(in)
	if !strings.Contains(out, "<style>\n    .a { color: red; }\n</style>") {
		t.Fatalf("style block must be indented, got:\n%s", out)
	}
	if !strings.Contains(out, "<client>\n    console.log(\"hi\")\n</client>") {
		t.Fatalf("client block must be indented, got:\n%s", out)
	}
}

// Root section tags are Dreego keywords and stay at column zero; only their
// content is indented.
func TestFormatKeepsRootSectionTagsAtColumnZero(t *testing.T) {
	in := "  <server>\n  x := 1\n  </server>\n\n  <body>\n  <p>ok</p>\n  </body>\n"
	out := Format(in)
	if !strings.Contains(out, "<server>\n    x := 1\n</server>") {
		t.Fatalf("root server tag/content not normalized, got:\n%s", out)
	}
	if !strings.Contains(out, "<body>\n    <p>ok</p>\n</body>") {
		t.Fatalf("root body tag/content not normalized, got:\n%s", out)
	}
}

func TestFormatIndentsControlFlowBlocks(t *testing.T) {
	in := "<body>\n{#if ok}\n<p>yes</p>\n{#else}\n<p>no</p>\n{/if}\n</body>\n"
	out := Format(in)
	want := "<body>\n    {#if ok}\n        <p>yes</p>\n    {#else}\n        <p>no</p>\n    {/if}\n</body>"
	if !strings.Contains(out, want) {
		t.Fatalf("control-flow blocks must indent their children\nwant:\n%s\ngot:\n%s", want, out)
	}
	if twice := Format(out); twice != out {
		t.Fatalf("not idempotent\nonce:  %q\ntwice: %q", out, twice)
	}
}

func TestFormatIndentsEachBlock(t *testing.T) {
	in := "<body>\n<ul>\n{#each items as item}\n<li>{{ item }}</li>\n{/each}\n</ul>\n</body>\n"
	out := Format(in)
	want := "<body>\n    <ul>\n        {#each items as item}\n            <li>{{ item }}</li>\n        {/each}\n    </ul>\n</body>"
	if !strings.Contains(out, want) {
		t.Fatalf("each blocks must indent their children\nwant:\n%s\ngot:\n%s", want, out)
	}
}

// A markdown body is prose: indenting it would turn lines into code blocks and
// change the rendered document, so it is preserved verbatim.
func TestFormatPreservesMarkdownBody(t *testing.T) {
	in := "<body lang=\"md\">\n# Title\n\n- one\n- two\n</body>\n"
	out := Format(in)
	if !strings.Contains(out, "# Title\n\n- one\n- two") {
		t.Fatalf("markdown body must not be re-indented, got:\n%s", out)
	}
	if twice := Format(out); twice != out {
		t.Fatalf("not idempotent\nonce:  %q\ntwice: %q", out, twice)
	}
}

// Content inside {#verbatim} is literal and must not be re-indented.
func TestFormatPreservesVerbatimContent(t *testing.T) {
	in := "<body>\n{#verbatim}\n  keep   this\n    and this\n{/verbatim}\n</body>\n"
	out := Format(in)
	if !strings.Contains(out, "  keep   this\n    and this") {
		t.Fatalf("verbatim content must be preserved, got:\n%s", out)
	}
	if twice := Format(out); twice != out {
		t.Fatalf("not idempotent with verbatim\nonce:  %q\ntwice: %q", out, twice)
	}
}

// Content inside <pre> is literal and must not be re-indented.
func TestFormatPreservesPreContent(t *testing.T) {
	in := "<body>\n<pre>\n  keep   this\n    and this\n</pre>\n</body>\n"
	out := Format(in)
	if !strings.Contains(out, "  keep   this\n    and this") {
		t.Fatalf("pre content must be preserved, got:\n%s", out)
	}
	if twice := Format(out); twice != out {
		t.Fatalf("not idempotent with pre\nonce:  %q\ntwice: %q", out, twice)
	}
}
