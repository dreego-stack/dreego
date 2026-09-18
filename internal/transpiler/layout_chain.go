package transpiler

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// buildLayoutIndex maps normalised layout paths to their entries. Each file is
// registered both relative to the website root ("layouts/base.dreego") and
// relative to the working directory ("www/layouts/base.dreego") so a LAYOUT
// value written in either form resolves to the same entry.
func buildLayoutIndex(root string, entries map[string]*layoutEntry) map[string]*layoutEntry {
	index := map[string]*layoutEntry{}
	prefix := filepath.Base(root)
	for _, e := range entries {
		if e == nil || e.source == "" {
			continue
		}
		rel := relToRoot(root, e.source)
		for _, raw := range []string{rel, relToRoot(".", e.source), prefix + "/" + rel} {
			if key := normaliseLayoutPath(raw); key != "" {
				index[key] = e
			}
		}
	}
	return index
}

func normaliseLayoutPath(p string) string {
	p = filepath.ToSlash(strings.TrimSpace(p))
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return ""
	}
	p = path.Clean(p)
	if p == "." {
		return ""
	}
	return strings.TrimPrefix(p, "./")
}

func resolveLayoutChain(start *layoutEntry, index map[string]*layoutEntry) ([]*layoutEntry, error) {
	if start == nil {
		return nil, fmt.Errorf("resolve layout chain: nil start")
	}
	chain := []*layoutEntry{start}
	seen := map[*layoutEntry]int{start: 0}
	current := start
	for {
		if current.file == nil || current.file.Layout == "" {
			return chain, nil
		}
		declared := normaliseLayoutPath(current.file.Layout)
		next, ok := index[declared]
		if !ok {
			return nil, layoutChainMissingError(current, current.file.Layout)
		}
		if at, dup := seen[next]; dup {
			return nil, layoutChainCycleError(current, chain[at:], next)
		}
		seen[next] = len(chain)
		chain = append(chain, next)
		current = next
	}
}

func layoutChainMissingError(source *layoutEntry, declared string) error {
	file := source.file
	path := source.source
	if file != nil && file.SourcePath != "" {
		path = file.SourcePath
	}
	return fmt.Errorf("%s:%s: LAYOUT %q could not be resolved", path, layoutDeclPosition(file, declared), declared)
}

func layoutDeclPosition(file *File, declared string) string {
	if file == nil || file.SourceContent == "" || declared == "" {
		return "?:?"
	}
	if idx := strings.Index(file.SourceContent, `"`+declared+`"`); idx >= 0 {
		return sourceLocation(file.SourceContent, idx)
	}
	return "?:?"
}

func layoutChainCycleError(from *layoutEntry, cycle []*layoutEntry, repeated *layoutEntry) error {
	names := make([]string, 0, len(cycle)+1)
	for _, e := range cycle {
		names = append(names, e.source)
	}
	names = append(names, repeated.source)
	loc := "?:?"
	path := ""
	if from != nil {
		loc = layoutDeclPosition(from.file, from.file.Layout)
		path = from.source
	}
	return fmt.Errorf("%s:%s: layout cycle: %s", path, loc, strings.Join(names, " -> "))
}
