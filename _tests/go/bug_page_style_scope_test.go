package tests

import (
	"regexp"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugPageStyleScopeStaysOnDocumentBody(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<head><title>Scoped page</title></head>
<body><main class="card">Styled</main></body>
<style>.card { color: red; }</style>`,
	})

	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected HTTP 200, got %d", code)
	}
	scope := regexp.MustCompile(`<div data-scope="([a-f0-9]+)">`).FindStringSubmatch(body)
	if len(scope) != 2 {
		t.Fatalf("document page must carry a style scope, got: %s", body)
	}
	selector := `[data-scope="` + scope[1] + `"] .card`
	if !regexp.MustCompile(regexp.QuoteMeta(selector)).MatchString(body) {
		t.Fatalf("style scope value must be quoted for digit-leading hashes, got: %s", body)
	}
}
