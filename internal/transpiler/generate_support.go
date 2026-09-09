package transpiler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func modulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

func loadSettings(root string) (*Settings, error) {
	settingsPath := filepath.Join(root, configFileName)
	settings, err := LoadConfig(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		if errors.Is(err, ErrInvalidI18n) {
			return nil, fmt.Errorf("%s is invalid: %w", settingsPath, err)
		}
		slog.Warn("dreego: "+configFileName+" is invalid; using defaults", "path", settingsPath, "error", err)
		return nil, nil
	}
	return settings, nil
}

func isUpToDate(path, content string) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return string(existing) == content
}

func hashOf(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[:12]
}
