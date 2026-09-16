package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
	"github.com/dreego-stack/dreego/internal/gomod"
)

func cmdInit(args []string) {
	flags, err := parseScaffoldFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if flags.list {
		printTemplates()
		return
	}

	if len(flags.args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dreego init <path> [-t <template>]\n")
		os.Exit(1)
	}
	target := flags.args[0]
	if !validInitPath(target) {
		fmt.Fprintf(os.Stderr, "error: invalid init path %q\n", target)
		fmt.Fprintf(os.Stderr, "  the path must be a non-empty directory name (use '.' for the current directory).\n")
		os.Exit(1)
	}

	meta, err := resolveTemplate(flags.template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !goAvailable() {
		fmt.Fprintf(os.Stderr, "error: 'go' executable not found on PATH.\n")
		fmt.Fprintf(os.Stderr, "  Dreego requires Go 1.27 or newer. Install it from https://go.dev/doc/install and retry.\n")
		os.Exit(1)
	}

	if err := templates.Install(target, moduleName(target), meta.Name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("initialized %s\n", target)
}

// moduleName returns the module path declared in the go.mod of dir (first
// "module " line), or the directory base name as fallback when no go.mod
// exists. Used to qualify the generated package import in the scaffold.
func moduleName(dir string) string {
	file, err := gomod.Read(filepath.Join(dir, "go.mod"))
	if err == nil && file.Module != "" {
		return file.Module
	}
	return filepath.Base(dir)
}

// validInitPath reports whether path is a usable target directory for
// `dreego init`. It accepts '.' and relative/absolute single-segment names
// but rejects empty paths and names that contain shell metacharacters or
// path separators in unusual combinations.
func validInitPath(path string) bool {
	if path == "" {
		return false
	}
	if strings.ContainsAny(path, " \t\"'\\`$;|&<>(){}[]!*?") {
		return false
	}
	return true
}
