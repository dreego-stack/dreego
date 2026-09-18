package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func collectComponentAliases(gen *Generator, root string) error {
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("error reading %s: %w", path, readErr)
		}
		header, _, headerErr := ParseFileHeaderStrict(string(data))
		if headerErr != nil {
			return fmt.Errorf("%s:%w", path, headerErr)
		}
		for _, imp := range header.Imports {
			aliases := make([]string, 0, len(imp.Aliases))
			for alias := range imp.Aliases {
				aliases = append(aliases, alias)
			}
			sort.Strings(aliases)
			for _, alias := range aliases {
				if err := gen.RegisterCompAlias(alias, imp.Aliases[alias]); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
			}
		}
		return nil
	})
	return err
}

func validateComponentAliases(gen *Generator) error {
	aliases := make([]string, 0, len(gen.CompAliases))
	for alias := range gen.CompAliases {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		target := gen.CompAliases[alias]
		if _, clash := gen.Defs[alias]; clash {
			return fmt.Errorf("component alias %s conflicts with component %s", alias, alias)
		}
		if _, ok := gen.Defs[target]; !ok {
			return fmt.Errorf("component alias %s targets unknown component %s", alias, target)
		}
	}
	return nil
}
