package transpiler

import (
	"fmt"
	"os/exec"
	"strings"
)

func moduleDirectory(modulePath string) (string, error) {
	command := exec.Command("go", "list", "-m", "-f={{.Dir}}", modulePath)
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("resolve component module %s: %w", modulePath, err)
	}
	directory := strings.TrimSpace(string(output))
	if directory == "" {
		return "", fmt.Errorf("resolve component module %s: empty module directory", modulePath)
	}
	return directory, nil
}
