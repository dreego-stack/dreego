package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dreego "codeberg.org/dreego/dreego/core"
)

func cmdFmt(args []string) {
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
		formatted := dreego.Format(original)

		if formatted == original {
			if !check && !write {
				fmt.Print(formatted)
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
			fmt.Printf("%s\n", f)
		} else {
			fmt.Print(formatted)
		}
	}

	if check && changed > 0 {
		fmt.Fprintf(os.Stderr, "fmt: %d file(s) need formatting\n", changed)
		os.Exit(1)
	}
}
