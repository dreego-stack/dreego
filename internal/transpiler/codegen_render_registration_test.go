package transpiler

import "testing"

func TestDesktopRoutePathConvertsHTTPExactPatterns(t *testing.T) {
	for _, test := range []struct {
		pattern string
		want    string
	}{
		{pattern: "/{$}", want: "/"},
		{pattern: "/news/{$}", want: "/news/"},
		{pattern: "/about", want: "/about"},
	} {
		got, ok := desktopRoutePath(test.pattern)
		if !ok || got != test.want {
			t.Fatalf("desktopRoutePath(%q) = %q, %t, want %q, true", test.pattern, got, ok, test.want)
		}
	}
}

func TestDesktopRoutePathRejectsDynamicPattern(t *testing.T) {
	if got, ok := desktopRoutePath("/users/{id}"); ok || got != "" {
		t.Fatalf("desktopRoutePath dynamic = %q, %t, want empty, false", got, ok)
	}
}
