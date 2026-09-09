package transpiler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dreego-stack/dreego/internal/gomod"
)

const maxPluginClientModuleSize = 256 << 10
const maxPluginClientBundleSize = 1 << 20

type pluginClientManifest struct {
	Client pluginClientBundle `json:"client"`
}

type pluginClientBundle struct {
	Format  string               `json:"format"`
	Path    string               `json:"path"`
	Modules []pluginClientModule `json:"modules"`
}

type pluginClientModule struct {
	ID        string   `json:"id"`
	Path      string   `json:"path"`
	DependsOn []string `json:"dependsOn"`
	Required  bool     `json:"required"`
}

func generatePluginClientAssets(project string, settings *Settings, routePatterns map[string]bool) (string, int, error) {
	if settings == nil || len(settings.Plugins) == 0 {
		return "", 0, nil
	}
	modules, err := requiredModules(project)
	if err != nil {
		return "", 0, err
	}
	paths := make([]string, 0, len(settings.Plugins))
	for modulePath := range settings.Plugins {
		paths = append(paths, modulePath)
	}
	sort.Strings(paths)

	var src strings.Builder
	for _, modulePath := range paths {
		if !strings.HasPrefix(modulePath, "github.com/dreego-stack/plugin-") {
			return "", 0, fmt.Errorf("plugin client: unsupported plugin module %q", modulePath)
		}
		version, ok := modules[modulePath]
		if !ok {
			return "", 0, fmt.Errorf("plugin client: %s is not required by go.mod", modulePath)
		}
		dir, err := pluginModuleDir(project, modulePath, version)
		if err != nil {
			return "", 0, fmt.Errorf("plugin client %s: %w", modulePath, err)
		}
		bundle, err := loadPluginClientBundle(dir, settings.Plugins[modulePath].Client)
		if err != nil {
			return "", 0, fmt.Errorf("plugin client %s: %w", modulePath, err)
		}
		if routePatterns["GET "+bundle.path] {
			return "", 0, fmt.Errorf("plugin client %s path %q conflicts with existing route", modulePath, bundle.path)
		}
		routePatterns["GET "+bundle.path] = true
		src.WriteString(registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, []byte(%q))", bundle.path, "application/javascript; charset=utf-8", bundle.content)))
	}
	return src.String(), len(paths), nil
}

type generatedPluginBundle struct {
	path    string
	content string
}

func loadPluginClientBundle(dir string, selected []string) (generatedPluginBundle, error) {
	body, err := os.ReadFile(filepath.Join(dir, "dreego-plugin.json"))
	if err != nil {
		return generatedPluginBundle{}, fmt.Errorf("read dreego-plugin.json: %w", err)
	}
	var manifest pluginClientManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return generatedPluginBundle{}, fmt.Errorf("parse dreego-plugin.json: %w", err)
	}
	if manifest.Client.Format != "modules-v1" {
		return generatedPluginBundle{}, fmt.Errorf("unsupported client format %q", manifest.Client.Format)
	}
	if err := validatePluginClientURL(manifest.Client.Path); err != nil {
		return generatedPluginBundle{}, err
	}
	byID := make(map[string]pluginClientModule, len(manifest.Client.Modules))
	for _, module := range manifest.Client.Modules {
		if !validPluginClientID(module.ID) {
			return generatedPluginBundle{}, fmt.Errorf("invalid client module id %q", module.ID)
		}
		if _, exists := byID[module.ID]; exists {
			return generatedPluginBundle{}, fmt.Errorf("duplicate client module id %q", module.ID)
		}
		if err := validatePluginClientFile(module.Path); err != nil {
			return generatedPluginBundle{}, fmt.Errorf("module %q: %w", module.ID, err)
		}
		byID[module.ID] = module
	}
	allIDs := make([]string, 0, len(byID))
	for id := range byID {
		allIDs = append(allIDs, id)
	}
	sort.Strings(allIDs)
	if _, err := orderPluginClientModules(byID, allIDs); err != nil {
		return generatedPluginBundle{}, err
	}

	roots := append([]string(nil), selected...)
	for _, module := range manifest.Client.Modules {
		if module.Required {
			roots = append(roots, module.ID)
		}
	}
	sort.Strings(roots)
	ordered, err := orderPluginClientModules(byID, roots)
	if err != nil {
		return generatedPluginBundle{}, err
	}

	var content strings.Builder
	for _, id := range ordered {
		module := byID[id]
		data, err := readPluginClientFile(dir, module.Path)
		if err != nil {
			return generatedPluginBundle{}, fmt.Errorf("module %q: %w", id, err)
		}
		if len(data) > maxPluginClientModuleSize {
			return generatedPluginBundle{}, fmt.Errorf("module %q exceeds %d bytes", id, maxPluginClientModuleSize)
		}
		fmt.Fprintf(&content, "// dreego-plugin-module:%s\n%s\n", id, data)
		if content.Len() > maxPluginClientBundleSize {
			return generatedPluginBundle{}, fmt.Errorf("client bundle exceeds %d bytes", maxPluginClientBundleSize)
		}
	}
	digest := sha256.Sum256([]byte(content.String()))
	result := "// dreego-plugin-sha256:" + hex.EncodeToString(digest[:]) + "\n" + content.String()
	return generatedPluginBundle{path: manifest.Client.Path, content: result}, nil
}

func orderPluginClientModules(modules map[string]pluginClientModule, roots []string) ([]string, error) {
	state := map[string]uint8{}
	var ordered []string
	var visit func(string) error
	visit = func(id string) error {
		module, ok := modules[id]
		if !ok {
			return fmt.Errorf("unknown client module %q", id)
		}
		switch state[id] {
		case 1:
			return fmt.Errorf("client module dependency cycle at %q", id)
		case 2:
			return nil
		}
		state[id] = 1
		deps := append([]string(nil), module.DependsOn...)
		sort.Strings(deps)
		for _, dependency := range deps {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		state[id] = 2
		ordered = append(ordered, id)
		return nil
	}
	for _, root := range roots {
		if err := visit(root); err != nil {
			return nil, err
		}
	}
	return ordered, nil
}

func validatePluginClientURL(value string) error {
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.ContainsAny(value, "\\?#") || path.Clean(value) != value || !strings.HasSuffix(value, ".js") {
		return fmt.Errorf("invalid client bundle path %q", value)
	}
	return nil
}

func validatePluginClientFile(value string) error {
	clean := filepath.Clean(filepath.FromSlash(value))
	if value == "" || strings.Contains(value, "\\") || filepath.IsAbs(value) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("unsafe client module path %q", value)
	}
	if filepath.Ext(clean) != ".js" {
		return fmt.Errorf("client module path %q must be a .js file", value)
	}
	return nil
}

func validPluginClientID(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if i > 0 && (r >= '0' && r <= '9' || r == '-' || r == '_') {
			continue
		}
		return false
	}
	return true
}

func readPluginClientFile(dir, relative string) ([]byte, error) {
	root, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return nil, err
	}
	full, err := filepath.EvalSymlinks(filepath.Join(root, filepath.FromSlash(relative)))
	if err != nil {
		return nil, err
	}
	if full != root && !strings.HasPrefix(full, root+string(os.PathSeparator)) {
		return nil, fmt.Errorf("client module path %q escapes plugin directory", relative)
	}
	return os.ReadFile(full)
}

func requiredModules(project string) (map[string]string, error) {
	file, err := gomod.Read(filepath.Join(project, "go.mod"))
	return file.Requires, err
}

func pluginModuleDir(project, modulePath, version string) (string, error) {
	vendor := filepath.Join(project, "vendor", filepath.FromSlash(modulePath))
	if _, err := os.Stat(vendor); err == nil {
		return vendor, nil
	}
	command := exec.Command("go", "list", "-m", "-f={{.Dir}}", modulePath)
	command.Dir = project
	out, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("resolve module %s@%s; run go mod download: %w", modulePath, version, err)
	}
	dir := strings.TrimSpace(string(out))
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("module is not downloaded: %w", err)
	}
	return dir, nil
}
