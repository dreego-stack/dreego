package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestBuildRejectsTwoHopRedirectCycle(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/a", "/b", http.StatusFound); err != nil {
		t.Fatalf("register /a -> /b: %v", err)
	}
	if err := app.RegisterRedirect("/b", "/a", http.StatusFound); err != nil {
		t.Fatalf("register /b -> /a: %v", err)
	}
	assertBuildCycleError(t, app)
}

func TestBuildRejectsThreeHopRedirectCycle(t *testing.T) {
	t.Parallel()
	app := New()
	for _, c := range []struct{ from, to string }{
		{"/a", "/b"},
		{"/b", "/c"},
		{"/c", "/a"},
	} {
		if err := app.RegisterRedirect(c.from, c.to, http.StatusFound); err != nil {
			t.Fatalf("register %s -> %s: %v", c.from, c.to, err)
		}
	}
	assertBuildCycleError(t, app)
}

func TestBuildRejectsMixedRewriteRedirectCycle(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/a", "/b"); err != nil {
		t.Fatalf("register rewrite /a -> /b: %v", err)
	}
	if err := app.RegisterRedirect("/b", "/a", http.StatusFound); err != nil {
		t.Fatalf("register redirect /b -> /a: %v", err)
	}
	assertBuildCycleError(t, app)
}

func TestBuildRejectsCycleAcrossDuplicateFrom(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRewrite("/a", "/b"); err != nil {
		t.Fatalf("register rewrite /a -> /b: %v", err)
	}
	if err := app.RegisterRedirect("/a", "/c", http.StatusFound); err != nil {
		t.Fatalf("register redirect /a -> /c: %v", err)
	}
	if err := app.RegisterRedirect("/c", "/a", http.StatusFound); err != nil {
		t.Fatalf("register redirect /c -> /a: %v", err)
	}
	assertBuildCycleError(t, app)
}

func TestBuildAcceptsAcyclicRedirectChain(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/a", "/b", http.StatusFound); err != nil {
		t.Fatalf("register /a -> /b: %v", err)
	}
	if err := app.RegisterRewrite("/b", "/c"); err != nil {
		t.Fatalf("register rewrite /b -> /c: %v", err)
	}
	assertBuildSucceeds(t, app)
}

func assertBuildCycleError(t *testing.T, app *App) {
	t.Helper()
	if err := app.Build(); err == nil {
		t.Fatal("expected Build to return an error on redirect/rewrite cycle")
	}
	if app.built {
		t.Fatal("expected app state to be reset after a failed Build")
	}
	if h := app.Handler(); h == nil {
		t.Fatal("expected Handler to return a handler after a failed Build")
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	done := make(chan struct{})
	go func() {
		app.Handler().ServeHTTP(rec, req)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Handler blocked after a failed Build; app state was not reset")
	}
}

func assertBuildSucceeds(t *testing.T, app *App) {
	t.Helper()
	if err := app.Build(); err != nil {
		t.Fatalf("Build failed unexpectedly: %v", err)
	}
}
