package transpiler

import (
	"strings"
	"testing"
)

func TestLayoutHeadDedupeWarningOnNonLiteralPrefix(t *testing.T) {
	file := parseFile(t, "<body>\n<html>\n<head>\n[[ site.head ]]\n{#head}\n<title>Site</title>\n</head>\n<body>{#slot}</body>\n</html>\n</body>\n")
	setNodeSource(file.Body.Nodes, "www/layouts/default.dreego", 0)
	file.SourcePath = "www/layouts/default.dreego"

	warning, ok := layoutHeadDedupeWarning(file)
	if !ok {
		t.Fatal("non-literal layout head must emit a dedupe diagnostic")
	}
	for _, want := range []string{"www/layouts/default.dreego", "Fix:", "dedupe"} {
		if !strings.Contains(warning, want) {
			t.Errorf("diagnostic missing %q, got: %q", want, warning)
		}
	}
}

func TestLayoutHeadDedupeWarningLiteralPrefixSilent(t *testing.T) {
	file := parseFile(t, "<body>\n<html>\n<head>\n<title>Site</title>\n{#head}\n</head>\n<body>{#slot}</body>\n</html>\n</body>\n")
	setNodeSource(file.Body.Nodes, "www/layouts/default.dreego", 0)

	if _, ok := layoutHeadDedupeWarning(file); ok {
		t.Fatal("literal layout head must not emit a dedupe diagnostic")
	}
}

func TestLayoutHeadDedupeWarningWithoutDedupeTagSilent(t *testing.T) {
	file := parseFile(t, "<body>\n<html>\n<head>\n[[ site.head ]]\n{#head}\n</head>\n<body>{#slot}</body>\n</html>\n</body>\n")
	setNodeSource(file.Body.Nodes, "www/layouts/default.dreego", 0)

	if _, ok := layoutHeadDedupeWarning(file); ok {
		t.Fatal("layout without a dedupe tag must not emit a diagnostic")
	}
}

func TestGeneratedDedupeChecksAreCaseInsensitive(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<TITLE>Page</TITLE>`},
		Body: &BodySection{
			Nodes: []TemplateNode{{Type: NodeText, Content: "<p>page</p>"}},
		},
	}
	layout := &layoutEntry{
		file: &File{
			Head: &HeadSection{Content: "<TITLE>Site</TITLE>\n    {#head}"},
			Body: &BodySection{
				Nodes: []TemplateNode{{Type: NodeText, Content: "<body><main>{#slot}</main></body>"}},
			},
		},
		name: "Default",
	}

	out, err := genTempl(NewGenerator(), file, layout, "abc123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `strings.ToLower(pageHead)`) {
		t.Errorf("runtime dedupe must fold the head markup to lower case, got:\n%s", out)
	}

	helpers := headMergeHelpers()
	if !strings.Contains(helpers, "strings.ToLower") {
		t.Errorf("strip helpers must fold markup to lower case, got:\n%s", helpers)
	}
}
