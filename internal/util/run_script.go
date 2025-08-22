package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
)

func RunScript(script string, skipConfirmation bool) (string, error) {
	if !skipConfirmation && !getConfirmation(script) {
		return "", nil
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	pkgTmp, err := config.PKG_TMP()
	if err != nil {
		return "", err
	}
	script = fmt.Sprintf("set -euo pipefail\ncd %s\n%s", pkgTmp, script)
	cmd := exec.Command(shell, "-c", script)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("Error getting stderr pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("Error getting stderr pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error while starting command: %v", err)
	}

	stdoutData, err := io.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stderr: %v", err)
	}
	stderrData, err := io.ReadAll(stderr)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stderr: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("Error waiting for command: %v\nstdout:\n%s\nstderr:\n%s", err, stdoutData, stderrData)
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
