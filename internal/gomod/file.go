package gomod

import (
	"os"

	"golang.org/x/mod/modfile"
)

type File struct {
	Module   string
	Requires map[string]string
}

func Read(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	parsed, err := modfile.Parse(path, data, nil)
	if err != nil {
		return File{}, err
	}
	file := File{Requires: make(map[string]string, len(parsed.Require))}
	if parsed.Module != nil {
		file.Module = parsed.Module.Mod.Path
	}
	for _, requirement := range parsed.Require {
		file.Requires[requirement.Mod.Path] = requirement.Mod.Version
	}
	return file, nil
}
