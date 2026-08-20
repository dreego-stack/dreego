package transpiler

import (
	"testing"
)

func TestIsInsideDreegoProjectRoot(t *testing.T) {
	cases := map[string]bool{
		"dreego/routes":                      true,
		"dreego/routes/about":                true,
		"dreego/components":                  true,
		"dreego/layouts":                     true,
		"dreego/layouts/default.dreego":      true,
		"vendor/x/dreego/routes":             false,
		"node_modules/y/dreego/components":   false,
		"subapp/dreego/routes":               false,
		"./dreego/routes/get.dreego":         true,
		"./vendor/z/dreego/routes":           false,
		"./sub/dreego/components/Foo.dreego": false,
	}
	for in, want := range cases {
		if got := isInDreegoRoot(in); got != want {
			t.Errorf("isInDreegoRoot(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestIsSkippedOutsideRootDir(t *testing.T) {
	cases := map[string]bool{
		"vendor":        true,
		"node_modules":  true,
		".git":          true,
		".worktrees":    true,
		"dreego":        false,
		"dreego/routes": false,
		"build":         false,
		"tmp":           false,
		".tmp":          true,
	}
	for in, want := range cases {
		if got := isSkippedDir(in); got != want {
			t.Errorf("isSkippedDir(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestIsGeneratedDir(t *testing.T) {
	cases := map[string]bool{
		"dreego/gen":               true,
		"dreego/gen/sub":           true,
		"dreego/routes":            false,
		"dreego/routes/dreego/gen": false,
		"vendor/x/dreego/gen":      false,
		"dreego/components":        false,
		"gen":                      false,
	}
	for in, want := range cases {
		if got := isGeneratedDir(in); got != want {
			t.Errorf("isGeneratedDir(%q) = %v, want %v", in, got, want)
		}
	}
}
