package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugBodyLevelLayoutHeadDedupe(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html lang="en">
<head>
    <title>Site</title>
    <meta name="description" content="site desc">
    <meta charset="utf-8">
    {#head}
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><title>Page</title><meta name="description" content="route desc"></head>
<body><h1>Page</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Page</title>") {
		t.Fatalf("route title missing in body: %s", body)
	}
	if strings.Contains(body, "<title>Site</title>") {
		t.Fatalf("layout title still present in body: %s", body)
	}
	if n := strings.Count(body, `name="description"`); n != 1 {
		t.Fatalf("expected exactly 1 meta description, got %d: %s", n, body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
	if strings.Contains(body, `content="site desc"`) {
		t.Fatalf("layout meta description still present in body: %s", body)
	}
	for _, want := range []string{
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`</head>`,
		`<main>`,
		`<h1>Page</h1>`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("layout skeleton lost %q in body: %s", want, body)
		}
	}
}

func TestBugBodyLevelLayoutHeadDedupePlaceholderFirst(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html lang="en">
<head>
    {#head}
    <title>Site</title>
    <meta name="description" content="site desc">
    <meta charset="utf-8">
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><title>Page</title><meta name="description" content="route desc"></head>
<body><h1>Page</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Page</title>") {
		t.Fatalf("route title missing in body: %s", body)
	}
	if strings.Contains(body, "<title>Site</title>") {
		t.Fatalf("layout title still present in body: %s", body)
	}
	if n := strings.Count(body, `name="description"`); n != 1 {
		t.Fatalf("expected exactly 1 meta description, got %d: %s", n, body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
	if strings.Contains(body, `content="site desc"`) {
		t.Fatalf("layout meta description still present in body: %s", body)
	}
	for _, want := range []string{
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`</head>`,
		`<main>`,
		`<h1>Page</h1>`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("layout skeleton lost %q in body: %s", want, body)
		}
	}
}

func TestBugBodyLevelLayoutKeepsTailTitleWithoutRouteTitle(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html lang="en">
<head>
    {#head}
    <title>Site</title>
    <meta charset="utf-8">
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><meta name="description" content="route desc"></head>
<body><h1>Page</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Site</title>") {
		t.Fatalf("layout tail title must be kept when the route defines none: %s", body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
}

func TestBugBodyLevelLayoutHeadDedupeBothSides(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html lang="en">
<head>
    <title>Before</title>
    <meta name="description" content="before desc">
    {#head}
    <title>After</title>
    <meta name="description" content="after desc">
    <meta charset="utf-8">
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><title>Page</title><meta name="description" content="route desc"></head>
<body><h1>Page</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Page</title>") {
		t.Fatalf("route title missing in body: %s", body)
	}
	for _, unwanted := range []string{"<title>Before</title>", "<title>After</title>"} {
		if strings.Contains(body, unwanted) {
			t.Fatalf("layout title %q still present in body: %s", unwanted, body)
		}
	}
	if n := strings.Count(body, `name="description"`); n != 1 {
		t.Fatalf("expected exactly 1 meta description, got %d: %s", n, body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
	if !strings.Contains(body, `<meta charset="utf-8">`) {
		t.Fatalf("layout skeleton lost <meta charset> in body: %s", body)
	}
}

func TestBugBodyLevelLayoutKeepsLayoutTitleWithoutRouteTitle(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html lang="en">
<head>
    <title>Site</title>
    <meta charset="utf-8">
    {#head}
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><meta name="description" content="route desc"></head>
<body><h1>Page</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Site</title>") {
		t.Fatalf("layout title must be kept when the route defines none: %s", body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
}
