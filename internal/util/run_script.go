package util

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func RunScript(script string, skipConfirmation bool) (string, error) {
	if !skipConfirmation && !getConfirmation(script) {
		return "", nil
	}
	script = fmt.Sprintf("set -euo pipefail\ncd %s\n%s", config.PKG_TMP(), script)
	cmd := exec.Command("/bin/sh", "-c", script)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("Error getting stderr pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error while starting command: %v", err)
	}

	stdoutData, err := io.ReadAll(stderr)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stderr: %v", err)
	}
	stderrData, err := io.ReadAll(stderr)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stderr: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Errorf("Error waiting for command: %v\n", err)
	}

	return string(stdoutData), fmt.Errorf("%s", stderrData)
}

func getConfirmation(script string) bool {
	fmt.Printf("Commands to run:\n")
	for line := range strings.Lines(script) {
		fmt.Printf("  %s", SyntaxHighlight(line))
	}
	confirmation := "N"
	fmt.Print("\nProceed? [y/N]: ")
	fmt.Scanln(&confirmation)

	return strings.ToLower(confirmation) == "y"
}
