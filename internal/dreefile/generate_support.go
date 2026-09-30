package dreefile

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dreego-stack/dreego/internal/gomod"
	"github.com/dreego-stack/dreego/internal/urlrule"
)

func modulePath() string {
	file, err := gomod.Read("go.mod")
	if err == nil {
		return file.Module
	}
	return ""
}

func moduleRequires() map[string]string {
	file, err := gomod.Read("go.mod")
	if err != nil {
		return nil
	}
	return file.Requires
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

// loadAppSettings merges an app's dreego.config.json on top of the website
// root's config. Both are optional; a missing or invalid file falls back to the
// other layer.
func loadAppSettings(rootDir, appDir string) (*Settings, error) {
	root, err := loadSettings(rootDir)
	if err != nil {
		return nil, err
	}
	if appDir == "" {
		return root, nil
	}
	app, err := loadSettings(appDir)
	if err != nil {
		return nil, err
	}
	return MergeSettings(root, app), nil
}

// validateUrlRules checks every redirect and rewrite rule at generation time so
// an invalid rule fails `dreego generate` instead of panicking when the
// generated app calls dreego.New.
func validateUrlRules(settings *Settings) error {
	if settings == nil {
		return nil
	}
	var redirects, rewrites [][2]string
	for _, rd := range settings.Redirects {
		if err := urlrule.Redirect(rd.From, rd.To, rd.Status); err != nil {
			return err
		}
		redirects = append(redirects, [2]string{rd.From, rd.To})
	}
	for _, rw := range settings.Rewrites {
		if err := urlrule.Rewrite(rw.From, rw.To); err != nil {
			return err
		}
		rewrites = append(rewrites, [2]string{rw.From, rw.To})
	}
	return urlrule.Cycle(redirects, rewrites)
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
