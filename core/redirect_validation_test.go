package core

import (
	"testing"
)

func TestRegisterRedirectRejectsInvalidStatus(t *testing.T) {
	t.Parallel()
	cases := []int{200, 404, 0, 399, 400, 500, 309}
	for _, status := range cases {
		app := New()
		err := app.RegisterRedirect("/old", "/new", status)
		if err == nil {
			t.Errorf("status %d should be rejected", status)
		}
	}
}

func TestRegisterRedirectAcceptsValidStatus(t *testing.T) {
	t.Parallel()
	for _, status := range []int{301, 302, 303, 307, 308} {
		app := New()
		if err := app.RegisterRedirect("/old", "/new", status); err != nil {
			t.Errorf("status %d should be accepted, got %v", status, err)
		}
	}
}

func TestRegisterRedirectRejectsBadPattern(t *testing.T) {
	t.Parallel()
	cases := []struct{ from, to string }{
		{"", "/new"},
		{"old", "/new"},
		{"/old/", "/new"},
		{"//old", "/new"},
		{"/*", "/new"},
		{"/old", ""},
		{"/old", "new"},
		{"/old", "//new"},
		{"/old/*/*", "/new"},
		{"/ol*", "/new"},
	}
	for _, c := range cases {
		app := New()
		if err := app.RegisterRedirect(c.from, c.to, 301); err == nil {
			t.Errorf("from=%q to=%q should be rejected", c.from, c.to)
		}
	}
}

func TestRegisterRedirectRejectsSelfLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/old", "/old", 301); err == nil {
		t.Fatal("self-redirect /old -> /old should be rejected")
	}
	if err := app.RegisterRedirect("/old/*", "/old/*", 301); err == nil {
		t.Fatal("self-redirect /old/* -> /old/* should be rejected")
	}
}

func TestRegisterRedirectRejectsWildcardLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api/*", "/api/v2/*", 301); err == nil {
		t.Fatal("/api/* -> /api/v2/* loops back to /api/*, should be rejected")
	}
}

func TestRegisterRedirectRejectsWildcardExactLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api/*", "/api", 301); err == nil {
		t.Fatal("/api/* -> /api redirects /api/x to itself, should be rejected")
	}
	if err := app.RegisterRedirect("/api/*", "/api/v2", 301); err == nil {
		t.Fatal("/api/* -> /api/v2 loops back into /api/*, should be rejected")
	}
}

func TestRegisterRedirectRejectsWildcardRootTarget(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api/*", "/", 301); err == nil {
		t.Fatal("/api/* -> / emits //-prefixed targets, should be rejected")
	}
}

func TestRegisterRedirectAcceptsSegmentSafeWildcardTarget(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api/*", "/apiary/*", 301); err != nil {
		t.Fatalf("/api/* -> /apiary/* is segment-safe, should be accepted: %v", err)
	}
	if err := app.RegisterRedirect("/api/*", "/apix", 301); err != nil {
		t.Fatalf("/api/* -> /apix is segment-safe, should be accepted: %v", err)
	}
}

func TestRegisterRedirectRejectsExactFromWildcardTo(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api", "/v2/*", 301); err == nil {
		t.Fatal("exact /api -> wildcard /v2/* emits literal * in Location, should be rejected")
	}
}

func TestRegisterRedirectRejectsWhenBuilt(t *testing.T) {
	t.Parallel()
	app := New()
	app.Build()
	if err := app.RegisterRedirect("/old", "/new", 301); err == nil {
		t.Fatal("registering redirect after Build should fail")
	}
}

func TestRegisterRewriteRejectsBadPattern(t *testing.T) {
	t.Parallel()
	cases := []struct{ from, to string }{
		{"", "/new"},
		{"/old", ""},
		{"old", "/new"},
		{"/old", "new"},
		{"//old", "/new"},
		{"/old", "//new"},
		{"/*", "/new"},
		{"/old/*/*", "/new/*"},
		{"/ol*", "/new"},
	}
	for _, c := range cases {
		app := New()
		if err := app.RegisterRewrite(c.from, c.to); err == nil {
			t.Errorf("from=%q to=%q should be rejected", c.from, c.to)
		}
	}
}

func TestRegisterRewriteRejectsSelfLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/old", "/old"); err == nil {
		t.Fatal("self-rewrite /old -> /old should be rejected")
	}
}

func TestRegisterRewriteRejectsWildcardLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/api/*", "/api/v2/*"); err == nil {
		t.Fatal("/api/* -> /api/v2/* loops back, should be rejected")
	}
}

func TestRegisterRewriteRejectsWildcardExactLoop(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/api/*", "/api"); err == nil {
		t.Fatal("/api/* -> /api rewrites /api/x to itself, should be rejected")
	}
	if err := app.RegisterRewrite("/api/*", "/api/v2"); err == nil {
		t.Fatal("/api/* -> /api/v2 loops back into /api/*, should be rejected")
	}
}

func TestRegisterRewriteRejectsWildcardRootTarget(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/api/*", "/"); err == nil {
		t.Fatal("/api/* -> / emits //-prefixed paths, should be rejected")
	}
}

func TestRegisterRewriteAcceptsSegmentSafeWildcardTarget(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/api/*", "/apiary/*"); err != nil {
		t.Fatalf("/api/* -> /apiary/* is segment-safe, should be accepted: %v", err)
	}
}

func TestRegisterRewriteRejectsExactFromWildcardTo(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/api", "/v2/*"); err == nil {
		t.Fatal("exact /api -> wildcard /v2/* emits literal * in path, should be rejected")
	}
}

func TestRegisterRewriteRejectsWhenBuilt(t *testing.T) {
	t.Parallel()
	app := New()
	app.Build()
	if err := app.RegisterRewrite("/old/*", "/new/*"); err == nil {
		t.Fatal("registering rewrite after Build should fail")
	}
}
