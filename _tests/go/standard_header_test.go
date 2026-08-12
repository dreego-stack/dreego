package tests

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const usingStandardLine = "# Using standard: _tests/how-to-test-sh.md"

var coreTestsDir = filepath.Join("..", "..", "_tests", "core")

func lineAt(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}
	return ""
}

func checkStandardHeader(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"read " + path + ": " + err.Error()}
	}
	lines := strings.Split(string(data), "\n")
	var violations []string
	if got := lineAt(lines, 0); got != "#!/bin/sh" {
		violations = append(violations, fmt.Sprintf("line 1: expected \"#!/bin/sh\", got %q", got))
	}
	if got := lineAt(lines, 1); got != usingStandardLine {
		violations = append(violations, fmt.Sprintf("line 2: expected %q, got %q", usingStandardLine, got))
	}
	if got := lineAt(lines, 2); !strings.HasPrefix(got, "# What: ") {
		violations = append(violations, fmt.Sprintf("line 3: expected line starting with \"# What: \", got %q", got))
	}
	return violations
}

func TestCheckStandardHeader(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		content string
		want    int
	}{
		{"compliant", "#!/bin/sh\n# Using standard: _tests/how-to-test-sh.md\n# What: Test something\nset -e\n", 0},
		{"missing using standard", "#!/bin/sh\n# What: Test something\nset -e\n", 2},
		{"wrong order", "#!/bin/sh\n# What: Test something\n# Using standard: _tests/how-to-test-sh.md\n", 2},
		{"missing what", "#!/bin/sh\n# Using standard: _tests/how-to-test-sh.md\nset -e\n", 1},
		{"missing shebang", "# Using standard: _tests/how-to-test-sh.md\n# What: Test something\n", 3},
		{"empty file", "", 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.sh")
			if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			got := checkStandardHeader(path)
			if len(got) != tc.want {
				t.Fatalf("expected %d violations, got %d: %v", tc.want, len(got), got)
			}
		})
	}
}

func TestStandardHeaderAllTests(t *testing.T) {
	t.Parallel()
	found := 0
	var violations []string
	err := filepath.WalkDir(coreTestsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "test.sh" {
			return nil
		}
		found++
		violations = append(violations, checkStandardHeader(path)...)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", coreTestsDir, err)
	}
	if found == 0 {
		t.Fatalf("no test.sh found under %s", coreTestsDir)
	}
	if len(violations) > 0 {
		t.Fatalf("standard header violations:\n%s", strings.Join(violations, "\n"))
	}
}
