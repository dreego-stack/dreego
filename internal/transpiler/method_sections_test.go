package transpiler

import (
	"strings"
	"testing"
)

func parseMethodFixture(t *testing.T, source string) *File {
	t.Helper()
	tokens, err := Lex(source)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

func TestParseMethodSectionsRemainIsolated(t *testing.T) {
	file := parseMethodFixture(t, `<go>getValue := "get"</go>
<go method="post">postValue := "post"</go>
<div><p>{{ getValue }}</p></div>
<div method="post"><p>{{ postValue }}</p></div>`)
	if len(file.Go) != 2 || len(file.Templates) != 2 {
		t.Fatalf("expected two Go and two template sections, got %d and %d", len(file.Go), len(file.Templates))
	}
	if file.Go[0].Method != "GET" || file.Go[1].Method != "POST" {
		t.Fatalf("unexpected Go methods: %#v", file.Go)
	}
	if file.Templates[0].Method != "GET" || file.Templates[1].Method != "POST" {
		t.Fatalf("unexpected template methods: %#v", file.Templates)
	}
	if !file.Go[1].MethodExplicit || !file.Templates[1].MethodExplicit {
		t.Fatal("explicit method attributes were not preserved")
	}
}

func TestParseDuplicateMethodTemplateReportsMethod(t *testing.T) {
	parseExpectError(t,
		`<div method="post">one</div><div method="POST">two</div>`,
		"duplicate <div> section for method POST")
}

func TestParseMethodSectionsSupportAllCommonHTTPMethods(t *testing.T) {
	file := parseMethodFixture(t, `<go>get := 1</go><div>GET</div>
<go method="post">post := 1</go><div method="post">POST</div>
<go method="put">put := 1</go><div method="put">PUT</div>
<go method="delete">del := 1</go><div method="delete">DELETE</div>`)
	for i, want := range []string{"GET", "POST", "PUT", "DELETE"} {
		if file.Go[i].Method != want || file.Templates[i].Method != want {
			t.Errorf("section %d methods = %q/%q, want %q", i, file.Go[i].Method, file.Templates[i].Method, want)
		}
	}
}

func TestGenerateMethodHandlerDoesNotMixMethodTemplates(t *testing.T) {
	file := parseMethodFixture(t, `<go>value := "get"</go><go method="post">value := "post"</go>
<div><p>{{ value }}</p></div><div method="post"><p>{{ value }}</p></div>`)
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "main", "account", "/account", "abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !strings.Contains(out, "HandleAccount") || !strings.Contains(out, "HandleAccountPOST") {
		t.Fatalf("generated output missing GET and POST handlers:\n%s", out)
	}
}

func TestGenerateMethodHandlerRegistersEachMethod(t *testing.T) {
	file := parseMethodFixture(t, `<div>get</div><div method="post">post</div><div method="put">put</div>`)
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "main", "settings", "/settings", "abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	for _, want := range []string{"HandleSettings(", "HandleSettingsPOST(", "HandleSettingsPUT("} {
		if !strings.Contains(out, want) {
			t.Errorf("generated output missing %q:\n%s", want, out)
		}
	}
}

func TestParseMethodSectionsIsDeterministic(t *testing.T) {
	source := `<go method="post">x := "<x>"</go><div method="post"><p>{{ x }}</p></div>`
	first := parseMethodFixture(t, source)
	second := parseMethodFixture(t, source)
	if first.Go[0].Code != second.Go[0].Code || first.Templates[0].Nodes[0].Content != second.Templates[0].Nodes[0].Content {
		t.Fatal("method section parsing is not deterministic")
	}
}
