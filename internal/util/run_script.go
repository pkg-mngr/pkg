package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
)

func RunScript(script string, skipConfirmation bool) (string, error) {
	if !skipConfirmation && !getConfirmation(script) {
		return "", nil
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	script = fmt.Sprintf("set -euo pipefail\ncd %s\n%s", config.PKG_TMP, script)
	cmd := exec.Command(shell, "-c", script)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("Error running command: %v", err)
		return stdout.String(), fmt.Errorf("%s", stderr.String())
	}

	return stdout.String(), nil
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
