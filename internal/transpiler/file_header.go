package transpiler

import (
	"fmt"
	"path/filepath"
	"strings"
)

func componentNameFromPath(path string) (string, error) {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ".dreego")
	if name == base {
		return "", nil
	}
	if !IsExportedGoIdentifier(name) {
		return "", fmt.Errorf("invalid component filename %q: the component name comes from the file name and must be an exported Go identifier (start with an uppercase letter, then letters or digits), for example Card.dreego or ProductCard.dreego", base)
	}
	return name, nil
}

// IsExportedGoIdentifier reports whether name is a valid exported Go
// identifier, which is the requirement for a DREEFILE component file name.
func IsExportedGoIdentifier(name string) bool {
	if name == "" || name[0] < 'A' || name[0] > 'Z' {
		return false
	}
	for i := 1; i < len(name); i++ {
		c := name[i]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			continue
		}
		return false
	}
	return true
}
