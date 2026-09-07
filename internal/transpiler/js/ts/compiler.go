package ts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	jsoutput "github.com/dreego-stack/dreego/internal/transpiler/js/output"
)

const Version = "7.0.2"

func Process(code, declarations, sourcePath string, lineOffset int) (jsoutput.Artifact, error) {
	compiler, err := compilerPath()
	if err != nil {
		return jsoutput.Artifact{}, err
	}
	dir, err := os.MkdirTemp("", "dreego-typescript-")
	if err != nil {
		return jsoutput.Artifact{}, err
	}
	defer os.RemoveAll(dir)

	inputPath := filepath.Join(dir, "input.ts")
	if err := os.WriteFile(inputPath, []byte(code), 0600); err != nil {
		return jsoutput.Artifact{}, err
	}
	declarationPath := filepath.Join(dir, "dreego.d.ts")
	if err := os.WriteFile(declarationPath, []byte(declarations), 0600); err != nil {
		return jsoutput.Artifact{}, err
	}

	args := []string{"--ignoreConfig", "--pretty", "false", "--strict", "--noEmitOnError", "--target", "ES2022", "--module", "preserve", "--outDir", dir, inputPath, declarationPath}
	args = append(strings.Fields(os.Getenv("DREEGO_TYPESCRIPT_COMPILER_ARGS")), args...)
	cmd := exec.Command(compiler, args...)
	compilerOutput, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return jsoutput.Artifact{}, fmt.Errorf("TypeScript %s: %s", sourceLabel(sourcePath, lineOffset), normalizeDiagnostic(string(compilerOutput), inputPath, sourcePath, lineOffset))
	}
	emitted, err := os.ReadFile(filepath.Join(dir, "input.js"))
	if err != nil {
		return jsoutput.Artifact{}, fmt.Errorf("TypeScript compiler produced no JavaScript: %w", err)
	}
	return jsoutput.Artifact{Code: strings.TrimSpace(string(emitted))}, nil
}

func compilerPath() (string, error) {
	if configured := os.Getenv("DREEGO_TYPESCRIPT_COMPILER"); configured != "" {
		return configured, nil
	}
	if cached, err := CachedCompilerPath(); err == nil {
		if _, statErr := os.Stat(cached); statErr == nil {
			return cached, nil
		}
	}
	if path, err := exec.LookPath("tsc"); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("TypeScript %s compiler is not installed; run 'dreego tools install typescript'", Version)
}

func sourceLabel(path string, line int) string {
	if path == "" {
		path = "input.dreego"
	}
	return fmt.Sprintf("error in %s near line %d", path, line)
}

func normalizeDiagnostic(diagnostic, temporaryPath, sourcePath string, lineOffset int) string {
	if sourcePath == "" {
		sourcePath = "input.dreego"
	}
	diagnostic = strings.ReplaceAll(diagnostic, temporaryPath, sourcePath)
	temporaryReference := regexp.MustCompile(`(?:[^\s(]*[/\\])?` + regexp.QuoteMeta(filepath.Base(temporaryPath)))
	diagnostic = temporaryReference.ReplaceAllString(diagnostic, sourcePath)
	location := regexp.MustCompile(regexp.QuoteMeta(sourcePath) + `\((\d+),(\d+)\)`)
	diagnostic = location.ReplaceAllStringFunc(diagnostic, func(match string) string {
		parts := location.FindStringSubmatch(match)
		line, err := strconv.Atoi(parts[1])
		if err != nil {
			return match
		}
		return fmt.Sprintf("%s(%d,%s)", sourcePath, line+lineOffset-1, parts[2])
	})
	return strings.TrimSpace(diagnostic)
}
