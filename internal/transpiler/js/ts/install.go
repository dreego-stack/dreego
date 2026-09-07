package ts

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type releaseAsset struct {
	name   string
	digest string
}

var releaseAssets = map[string]releaseAsset{
	"darwin/amd64":  {"typescript-darwin-x64.tgz", "eba158cb54050f723d5ff781438f33de5640054440bb4f2bd170cfe9bc2eb551"},
	"darwin/arm64":  {"typescript-darwin-arm64.tgz", "902e2fe1cf0799198ef902c6b8c310a450fef629a6baba41d45641ef75c04ebd"},
	"linux/amd64":   {"typescript-linux-x64.tgz", "7ecad6f67377e831856367ab062ef394f21506a611405bf8ac0ff039348637d3"},
	"linux/arm64":   {"typescript-linux-arm64.tgz", "c83d931ac9dd7549cde6e71246aa9d6a9812843023df3e277fe3b5dcf41dd0ea"},
	"windows/amd64": {"typescript-win32-x64.tgz", "61fc4e141d2bc687db580e71bbfa63b9c209f0310645d82ca1b457eb3a24fd19"},
	"windows/arm64": {"typescript-win32-arm64.tgz", "0a73534e6ee50cdbb2a29ac48657ca0ad13cf0f424cf63808e4df7baeb87b8be"},
}

func CachedCompilerPath() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	name := "tsc"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(cache, "dreego", "tools", "typescript", Version, "lib", name), nil
}

func Install() (string, error) {
	asset, ok := releaseAssets[runtime.GOOS+"/"+runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("TypeScript %s has no supported binary for %s/%s", Version, runtime.GOOS, runtime.GOARCH)
	}
	target, err := CachedCompilerPath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	tempDir, err := os.MkdirTemp("", "dreego-typescript-install-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)
	archivePath := filepath.Join(tempDir, asset.name)
	url := "https://github.com/microsoft/typescript-go/releases/download/typescript/v" + Version + "/" + asset.name
	if err := download(url, archivePath); err != nil {
		return "", err
	}
	if err := verifyDigest(archivePath, asset.digest); err != nil {
		return "", err
	}
	installRoot := filepath.Dir(filepath.Dir(target))
	if err := os.MkdirAll(installRoot, 0755); err != nil {
		return "", err
	}
	if err := extractCompiler(archivePath, installRoot); err != nil {
		return "", err
	}
	if err := os.Chmod(target, 0755); err != nil {
		return "", err
	}
	return target, nil
}

func download(url, destination string) error {
	client := &http.Client{Timeout: 2 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download TypeScript: %s", response.Status)
	}
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, io.LimitReader(response.Body, 64<<20))
	return err
}

func verifyDigest(path, expected string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf("TypeScript archive checksum mismatch: got %s", actual)
	}
	return nil
}

func extractCompiler(archivePath, destination string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !strings.HasPrefix(header.Name, "package/lib/") || header.Typeflag != tar.TypeReg {
			continue
		}
		rel := strings.TrimPrefix(header.Name, "package/")
		if strings.Contains(rel, "..") {
			return fmt.Errorf("unsafe TypeScript archive path %q", header.Name)
		}
		path := filepath.Join(destination, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, io.LimitReader(reader, 32<<20))
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
}
