package tests

import (
	"sort"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadHelpersEmittedOnceAndLayoutStyleKept(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"www/layouts/default.dreego": `<body>
<html>
<head>{#head}</head>
<body><main>{#slot}</main></body>
</html>
</body>

<style>
.layout-style-marker { color: #123456; }
</style>`,
		"www/routes/get.dreego": "<head><title>Home</title></head>\n<body><h1>Home</h1></body>",
		"www/routes/posts/get.dreego": "<head><title>Post</title><meta name=\"description\" content=\"post\"></head>\n<body><h1>Post</h1></body>",
	}

	generated := dreegotest.Build(t, files)

	keys := make([]string, 0, len(generated))
	for k := range generated {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var all strings.Builder
	for _, k := range keys {
		all.WriteString(generated[k])
		all.WriteString("\n")
	}
	out := all.String()

	if n := strings.Count(out, "func stripTitleTag("); n != 1 {
		t.Fatalf("expected stripTitleTag defined exactly once in generated code, got %d:\n%s", n, out)
	}
	if n := strings.Count(out, "func stripMetaDescriptionTag("); n != 1 {
		t.Fatalf("expected stripMetaDescriptionTag defined exactly once in generated code, got %d:\n%s", n, out)
	}
	layoutStart := strings.Index(out, "func Default(")
	layoutEnd := strings.Index(out, "package routes")
	if layoutStart < 0 || layoutEnd < 0 || layoutEnd <= layoutStart {
		t.Fatalf("could not locate generated layout function in output:\n%s", out)
	}
	layout := out[layoutStart:layoutEnd]
	if !strings.Contains(layout, ".layout-style-marker { color: #123456; }") {
		t.Fatalf("layout <style> content must be emitted raw without a scope prefix:\n%s", layout)
	}
	if strings.Contains(layout, "[data-scope=") {
		t.Fatalf("layout <style> content must not be scoped with a data-scope attribute:\n%s", layout)
	}
}
