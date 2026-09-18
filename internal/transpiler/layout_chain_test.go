package transpiler

import (
	"path/filepath"
	"strings"
	"testing"
)

func testLayoutEntry(source, declared string) *layoutEntry {
	return &layoutEntry{
		rel:    "",
		source: source,
		file:   &File{Layout: declared, SourcePath: source},
		name:   "Layout",
	}
}

func TestResolveLayoutChainWithoutLayout(t *testing.T) {
	start := testLayoutEntry("/tmp/site/www/layouts/base.dreego", "")
	chain, err := resolveLayoutChain(start, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 1 || chain[0] != start {
		t.Fatalf("chain = %#v, want just the start entry", chain)
	}
}

func TestResolveLayoutChainOneLevel(t *testing.T) {
	parent := testLayoutEntry("/tmp/site/www/layouts/base.dreego", "")
	child := testLayoutEntry("/tmp/site/www/layouts/admin.dreego", "layouts/base.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/base.dreego":  parent,
		"layouts/admin.dreego": child,
	})
	chain, err := resolveLayoutChain(child, index)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 2 || chain[0] != child || chain[1] != parent {
		t.Fatalf("chain = %#v, want [child parent]", chain)
	}
}

func TestResolveLayoutChainTwoLevels(t *testing.T) {
	base := testLayoutEntry("/tmp/site/www/layouts/base.dreego", "")
	admin := testLayoutEntry("/tmp/site/www/layouts/admin.dreego", "layouts/base.dreego")
	dashboard := testLayoutEntry("/tmp/site/www/layouts/dashboard.dreego", "layouts/admin.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/base.dreego":      base,
		"layouts/admin.dreego":     admin,
		"layouts/dashboard.dreego": dashboard,
	})
	chain, err := resolveLayoutChain(dashboard, index)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 3 || chain[0] != dashboard || chain[1] != admin || chain[2] != base {
		t.Fatalf("chain = %#v, want [dashboard admin base]", chain)
	}
}

func TestResolveLayoutChainPrefixedPathResolvesToSameEntry(t *testing.T) {
	parent := testLayoutEntry("/tmp/site/www/layouts/base.dreego", "")
	child := testLayoutEntry("/tmp/site/www/layouts/admin.dreego", "www/layouts/base.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/base.dreego":  parent,
		"layouts/admin.dreego": child,
	})
	chain, err := resolveLayoutChain(child, index)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 2 || chain[1] != parent {
		t.Fatalf("prefixed LAYOUT must resolve to the root-relative entry, got %#v", chain)
	}
}

func TestResolveLayoutChainMissingTarget(t *testing.T) {
	child := testLayoutEntry("/tmp/site/www/layouts/admin.dreego", "layouts/missing.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/admin.dreego": child,
	})
	chain, err := resolveLayoutChain(child, index)
	if err == nil {
		t.Fatalf("expected error, got chain %#v", chain)
	}
	for _, want := range []string{"admin.dreego", "layouts/missing.dreego"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q must mention %q", err, want)
		}
	}
}

func TestResolveLayoutChainCycle(t *testing.T) {
	a := testLayoutEntry("/tmp/site/www/layouts/a.dreego", "layouts/b.dreego")
	b := testLayoutEntry("/tmp/site/www/layouts/b.dreego", "layouts/a.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/a.dreego": a,
		"layouts/b.dreego": b,
	})
	chain, err := resolveLayoutChain(a, index)
	if err == nil {
		t.Fatalf("expected cycle error, got chain %#v", chain)
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error %q must mention cycle", err)
	}
	for _, want := range []string{"a.dreego", "b.dreego"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q must name %q", err, want)
		}
	}
}

func TestResolveLayoutChainSelfCycle(t *testing.T) {
	self := testLayoutEntry("/tmp/site/www/layouts/self.dreego", "layouts/self.dreego")
	index := buildLayoutIndex("/tmp/site/www", map[string]*layoutEntry{
		"layouts/self.dreego": self,
	})
	chain, err := resolveLayoutChain(self, index)
	if err == nil {
		t.Fatalf("expected self cycle error, got chain %#v", chain)
	}
	if !strings.Contains(err.Error(), "cycle") || !strings.Contains(err.Error(), "self.dreego") {
		t.Fatalf("error %q must mention cycle and self.dreego", err)
	}
}

func TestNormaliseLayoutPath(t *testing.T) {
	cases := map[string]string{
		"layouts/base.dreego":     "layouts/base.dreego",
		"/layouts/base.dreego":    "layouts/base.dreego",
		"./layouts/base.dreego":   "layouts/base.dreego",
		"www/layouts/base.dreego": "www/layouts/base.dreego",
		"layouts/../x.dreego":     "x.dreego",
		"":                        "",
	}
	for in, want := range cases {
		if got := normaliseLayoutPath(in); got != want {
			t.Errorf("normaliseLayoutPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDiscoverLayoutsBuildsIndex(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"www/layouts/default.dreego":      "DREEFILE layout\n\n<body><p>base</p></body>",
		"www/admin/layouts/layout.dreego": "DREEFILE layout\n\nLAYOUT \"www/layouts/default.dreego\"\n\n<body><p>admin</p></body>",
	})
	siteRoot := filepath.Join(root, "www")
	entries, index, err := discoverLayouts(siteRoot)
	if err != nil {
		t.Fatalf("discoverLayouts: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}
	admin := index["admin/layouts/layout.dreego"]
	if admin == nil {
		t.Fatalf("index missing root-relative admin entry: %#v", index)
	}
	chain, err := resolveLayoutChain(admin, index)
	if err != nil {
		t.Fatalf("resolveLayoutChain: %v", err)
	}
	if len(chain) != 2 || chain[1] != index["layouts/default.dreego"] {
		t.Fatalf("chain = %#v, want admin then base", chain)
	}
}
