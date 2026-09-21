package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	dreefile "github.com/dreego-stack/dreego/internal/dreefile"
)

func cmdFmt(args []string) {
	if err := cmdFmtE(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// cmdFmtE formats .dreego files and returns an error when --check finds files
// that would change. It never writes in check mode; a caller other than main
// can therefore use it in tests without exiting the process.
func cmdFmtE(args []string, stdout io.Writer) error {
	check := false
	write := true
	var targets []string
	for _, a := range args {
		switch a {
		case "--check":
			check = true
			write = false
		case "--stdout":
			write = false
		default:
			targets = append(targets, a)
		}
	}

	if len(targets) == 0 {
		targets = []string{"."}
	}

	var files []string
	for _, t := range targets {
		filepath.Walk(t, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".dreego") {
				return nil
			}
			files = append(files, path)
			return nil
		})
	}

	changed := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fmt: %s: %v\n", f, err)
			continue
		}

		original := string(data)
		formatted := dreefile.Format(original)

		if formatted == original {
			if !check && !write {
				fmt.Fprint(stdout, formatted)
			}
			continue
		}

		changed++
		if check {
			fmt.Fprintf(os.Stderr, "%s: not formatted\n", f)
			continue
		}
		if write {
			os.WriteFile(f, []byte(formatted), 0644)
			fmt.Fprintf(stdout, "%s\n", f)
		} else {
			fmt.Fprint(stdout, formatted)
		}
	}

	if check && changed > 0 {
		return fmt.Errorf("fmt: %d file(s) need formatting", changed)
	}
	return nil
}
