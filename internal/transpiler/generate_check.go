package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type genPlan struct {
	files map[string]string
	roots []string
}

type diffKind int

const (
	diffMissing diffKind = iota
	diffExtra
	diffStale
)

type fileDiff struct {
	path string
	kind diffKind
}

func (p genPlan) diff(disk map[string]string) []fileDiff {
	var diffs []fileDiff
	for path, want := range p.files {
		got, ok := disk[path]
		if !ok {
			diffs = append(diffs, fileDiff{path: path, kind: diffMissing})
			continue
		}
		if got != want {
			diffs = append(diffs, fileDiff{path: path, kind: diffStale})
		}
	}
	for path := range disk {
		if _, ok := p.files[path]; !ok {
			diffs = append(diffs, fileDiff{path: path, kind: diffExtra})
		}
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i].path < diffs[j].path })
	return diffs
}

func diffReport(d []fileDiff) string {
	if len(d) == 0 {
		return ""
	}
	var b strings.Builder
	for _, df := range d {
		switch df.kind {
		case diffMissing:
			fmt.Fprintf(&b, "missing: %s\n", df.path)
		case diffExtra:
			fmt.Fprintf(&b, "extra:   %s\n", df.path)
		case diffStale:
			fmt.Fprintf(&b, "stale:   %s\n", df.path)
		}
	}
	return b.String()
}

func readDiskFiles(roots []string) (map[string]string, error) {
	disk := map[string]string{}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("error walking %s: %w", path, walkErr)
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Base(path) != "dree.go" {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			rel, _ := filepath.Rel(".", path)
			disk[filepath.ToSlash(rel)] = string(data)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return disk, nil
}

func applyPlan(plan genPlan, force bool) error {
	for _, root := range plan.roots {
		if err := os.MkdirAll(root, 0755); err != nil {
			return err
		}
	}
	disk, err := readDiskFiles(plan.roots)
	if err != nil {
		return err
	}
	for path := range disk {
		if _, ok := plan.files[path]; !ok {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("error removing %s: %w", path, err)
			}
		}
	}
	var paths []string
	for p := range plan.files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		content := plan.files[p]
		if !force && isUpToDate(p, content) {
			continue
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			return fmt.Errorf("error writing %s: %w", p, err)
		}
	}
	return nil
}
