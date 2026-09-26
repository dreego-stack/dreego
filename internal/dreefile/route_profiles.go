package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func discoverRouteProfiles(root string) (map[string]string, error) {
	profiles := map[string]string{}
	sources := map[string]string{}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if !d.IsDir() || !isRoutesDir(root, path) {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", path, err)
		}
		rel := routeDirRel(root, path)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".dreego") {
				continue
			}
			full := filepath.Join(path, e.Name())
			data, err := os.ReadFile(full)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", full, err)
			}
			header, _, err := ParseFileHeaderStrict(string(data))
			if err != nil {
				return fmt.Errorf("%s:%w", full, err)
			}
			if header.Profile == "" {
				continue
			}
			if prev, ok := profiles[rel]; ok && prev != header.Profile {
				return fmt.Errorf("conflicting PROFILE in route folder %q: %q (%s) and %q (%s)",
					displayRouteFolder(rel), prev, sources[rel], header.Profile, full)
			}
			profiles[rel] = header.Profile
			sources[rel] = full
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func resolveRouteProfile(profiles map[string]string, routeDir string) string {
	for {
		if name, ok := profiles[routeDir]; ok {
			return name
		}
		if routeDir == "" {
			return ""
		}
		if idx := strings.LastIndexByte(routeDir, '/'); idx >= 0 {
			routeDir = routeDir[:idx]
		} else {
			routeDir = ""
		}
	}
}

func displayRouteFolder(rel string) string {
	if rel == "" {
		return "routes"
	}
	return "routes/" + rel
}
