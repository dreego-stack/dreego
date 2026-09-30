package urlrule

import (
	"net/http"
	"strings"
	"testing"
)

func TestPattern(t *testing.T) {
	cases := []struct {
		name string
		in   string
		ok   bool
	}{
		{"root", "/", true},
		{"single segment", "/a", true},
		{"nested", "/a/b/c", true},
		{"wildcard", "/a/*", true},
		{"wildcard root", "/api/*", true},
		{"dynamic", "/users/{id}", true},
		{"catchall", "/blog/{rest...}", true},
		{"chars", "/a-b_c.d", true},
		{"empty", "", false},
		{"missing slash", "a", false},
		{"double slash", "/a//b", false},
		{"trailing slash", "/a/", false},
		{"bare wildcard", "/*", false},
		{"wildcard mid", "/a/*/b", false},
		{"star no slash", "/a*b", false},
		{"double star", "/a/**", false},
		{"two wildcards", "/a/*/*", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Pattern(tc.in)
			if tc.ok && err != nil {
				t.Fatalf("Pattern(%q) = %v, want nil", tc.in, err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("Pattern(%q) = nil, want error", tc.in)
			}
		})
	}
}

func TestTarget(t *testing.T) {
	cases := []struct {
		name string
		in   string
		ok   bool
	}{
		{"root", "/", true},
		{"single segment", "/a", true},
		{"nested", "/a/b/c", true},
		{"wildcard", "/a/*", true},
		{"chars", "/a-b_c.d", true},
		{"empty", "", false},
		{"missing slash", "a", false},
		{"double slash", "/a//b", false},
		{"trailing slash", "/a/", false},
		{"star no suffix", "/a*b", false},
		{"wildcard mid", "/a/*/b", false},
		{"double star", "/a/**", false},
		{"two wildcards", "/a/*/*", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Target(tc.in)
			if tc.ok && err != nil {
				t.Fatalf("Target(%q) = %v, want nil", tc.in, err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("Target(%q) = nil, want error", tc.in)
			}
		})
	}
}

func TestPair(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
		ok       bool
	}{
		{"exact to exact", "/a", "/b", true},
		{"wildcard to wildcard", "/a/*", "/b/*", true},
		{"exact to root", "/a", "/", true},
		{"wildcard to root", "/a/*", "/", false},
		{"wildcard to wildcard target", "/a/*", "/*", false},
		{"exact to wildcard", "/a", "/b/*", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Pair(tc.from, tc.to)
			if tc.ok && err != nil {
				t.Fatalf("Pair(%q,%q) = %v, want nil", tc.from, tc.to, err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("Pair(%q,%q) = nil, want error", tc.from, tc.to)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
		status   int
		ok       bool
	}{
		{"301", "/old", "/new", http.StatusMovedPermanently, true},
		{"302", "/old", "/new", http.StatusFound, true},
		{"303", "/old", "/new", http.StatusSeeOther, true},
		{"307", "/old", "/new", http.StatusTemporaryRedirect, true},
		{"308", "/old", "/new", http.StatusPermanentRedirect, true},
		{"wildcard", "/api/*", "/v2/*", http.StatusMovedPermanently, true},
		{"status 200", "/old", "/new", http.StatusOK, false},
		{"status 0", "/old", "/new", 0, false},
		{"status 999", "/old", "/new", 999, false},
		{"empty pattern", "", "/new", http.StatusMovedPermanently, false},
		{"bad target", "/old", "new", http.StatusMovedPermanently, false},
		{"self loop", "/a", "/a", http.StatusMovedPermanently, false},
		{"wildcard loop", "/api/*", "/api/v2/*", http.StatusMovedPermanently, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Redirect(tc.from, tc.to, tc.status)
			if tc.ok && err != nil {
				t.Fatalf("Redirect(%q,%q,%d) = %v, want nil", tc.from, tc.to, tc.status, err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("Redirect(%q,%q,%d) = nil, want error", tc.from, tc.to, tc.status)
			}
		})
	}
}

func TestRewrite(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
		ok       bool
	}{
		{"exact", "/a", "/b", true},
		{"wildcard", "/api/*", "/v2/*", true},
		{"empty pattern", "", "/b", false},
		{"bad target", "/a", "b", false},
		{"self loop", "/a", "/a", false},
		{"wildcard loop", "/api/*", "/api/v2/*", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Rewrite(tc.from, tc.to)
			if tc.ok && err != nil {
				t.Fatalf("Rewrite(%q,%q) = %v, want nil", tc.from, tc.to, err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("Rewrite(%q,%q) = nil, want error", tc.from, tc.to)
			}
		})
	}
}

func TestLoops(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
		want     bool
	}{
		{"identical exact", "/a", "/a", true},
		{"different exact", "/a", "/b", false},
		{"identical wildcard", "/api/*", "/api/*", true},
		{"wildcard to child prefix", "/api/*", "/api/v2/*", true},
		{"wildcard to own prefix", "/api/*", "/api", true},
		{"wildcard to unrelated", "/api/*", "/v2/*", false},
		{"near prefix not loop", "/api/*", "/apix", false},
		{"exact to deeper path", "/a", "/a/b", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Loops(tc.from, tc.to); got != tc.want {
				t.Fatalf("Loops(%q,%q) = %v, want %v", tc.from, tc.to, got, tc.want)
			}
		})
	}
}

func TestCycle(t *testing.T) {
	cases := []struct {
		name      string
		redirects [][2]string
		rewrites  [][2]string
		wantErr   bool
	}{
		{"empty", nil, nil, false},
		{"linear redirect", [][2]string{{"/a", "/b"}, {"/b", "/c"}}, nil, false},
		{"self redirect", [][2]string{{"/a", "/a"}}, nil, true},
		{"two hop", [][2]string{{"/a", "/b"}, {"/b", "/a"}}, nil, true},
		{"three hop", [][2]string{{"/a", "/b"}, {"/b", "/c"}, {"/c", "/a"}}, nil, true},
		{"linear rewrite", nil, [][2]string{{"/a", "/b"}, {"/b", "/c"}}, false},
		{"rewrite redirect mix", [][2]string{{"/b", "/a"}}, [][2]string{{"/a", "/b"}}, true},
		{"duplicate from cycle", [][2]string{{"/a", "/b"}, {"/a", "/c"}, {"/c", "/a"}}, nil, true},
		{"duplicate from acyclic", [][2]string{{"/a", "/b"}, {"/a", "/c"}, {"/c", "/d"}}, nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Cycle(tc.redirects, tc.rewrites)
			if tc.wantErr && err == nil {
				t.Fatal("Cycle = nil, want error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Cycle = %v, want nil", err)
			}
		})
	}
}

func TestValidStatus(t *testing.T) {
	valid := []int{http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect}
	for _, s := range valid {
		if !ValidStatus(s) {
			t.Errorf("ValidStatus(%d) = false, want true", s)
		}
	}
	invalid := []int{0, 200, 201, 204, 300, 304, 306, 309, 400, 500, 999}
	for _, s := range invalid {
		if ValidStatus(s) {
			t.Errorf("ValidStatus(%d) = true, want false", s)
		}
	}
}

func TestCycleErrorNamesPath(t *testing.T) {
	err := Cycle([][2]string{{"/loop-a", "/loop-b"}, {"/loop-b", "/loop-a"}}, nil)
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("Cycle error = %v, want a cycle diagnostic", err)
	}
}
