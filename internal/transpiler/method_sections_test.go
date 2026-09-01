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
	file := parseMethodFixture(t, `<server>getValue := "get"</server>
<server method="post">postValue := "post"</server>
<body><p>{{ getValue }}</p></body>
<body method="post"><p>{{ postValue }}</p></body>`)
	if len(file.Server) != 2 || len(file.Bodies) != 2 {
		t.Fatalf("expected two Go and two template sections, got %d and %d", len(file.Server), len(file.Bodies))
	}
	if file.Server[0].Method != "GET" || file.Server[1].Method != "POST" {
		t.Fatalf("unexpected Go methods: %#v", file.Server)
	}
	if file.Bodies[0].Method != "GET" || file.Bodies[1].Method != "POST" {
		t.Fatalf("unexpected template methods: %#v", file.Bodies)
	}
	if !file.Server[1].MethodExplicit || !file.Bodies[1].MethodExplicit {
		t.Fatal("explicit method attributes were not preserved")
	}
}

func TestParseMethodSectionsSupportAllCommonHTTPMethods(t *testing.T) {
	file := parseMethodFixture(t, `<server>get := 1</server><body>GET</body>
<server method="post">post := 1</server><body method="post">POST</body>
<server method="put">put := 1</server><body method="put">PUT</body>
<server method="delete">del := 1</server><body method="delete">DELETE</body>`)
	for i, want := range []string{"GET", "POST", "PUT", "DELETE"} {
		if file.Server[i].Method != want || file.Bodies[i].Method != want {
			t.Errorf("section %d methods = %q/%q, want %q", i, file.Server[i].Method, file.Bodies[i].Method, want)
		}
	}
}

func TestGenerateMethodHandlerDoesNotMixMethodTemplates(t *testing.T) {
	file := parseMethodFixture(t, `<server>value := "get"</server><server method="post">value := "post"</server>
<body><p>{{ value }}</p></body><body method="post"><p>{{ value }}</p></body>`)
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "main", "account", "/account", "abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !strings.Contains(out, "HandleAccount") || !strings.Contains(out, "HandleAccountPOST") {
		t.Fatalf("generated output missing GET and POST handlers:\n%s", out)
	}
}

func TestGenerateMethodHandlerRegistersEachMethod(t *testing.T) {
	file := parseMethodFixture(t, `<body>get</body><body method="post">post</body><body method="put">put</body>`)
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
	source := `<server method="post">x := "<x>"</server><body method="post"><p>{{ x }}</p></body>`
	first := parseMethodFixture(t, source)
	second := parseMethodFixture(t, source)
	if first.Server[0].Code != second.Server[0].Code || first.Bodies[0].Nodes[0].Content != second.Bodies[0].Nodes[0].Content {
		t.Fatal("method section parsing is not deterministic")
	}
}
