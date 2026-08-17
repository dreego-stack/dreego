package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type genPlan struct {
	files  map[string]string
	genDir string
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

func readDiskFiles(genDir string) (map[string]string, error) {
	disk := map[string]string{}
	info, err := os.Stat(genDir)
	if err != nil {
		if os.IsNotExist(err) {
			return disk, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return disk, nil
	}
	err = filepath.WalkDir(genDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
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
	return disk, err
}

func applyPlan(plan genPlan) error {
	if err := os.MkdirAll(plan.genDir, 0755); err != nil {
		return err
	}
	disk, err := readDiskFiles(plan.genDir)
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
		if isUpToDate(p, content) {
			continue
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			return fmt.Errorf("error writing %s: %w", p, err)
		}
	}
	return nil
}
