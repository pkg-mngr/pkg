package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg-mngr/pkg/internal/manifest"
	"github.com/pkg-mngr/pkg/internal/util"
)

func Info(pkg string) (string, error) {
	pkgManifest, err := manifest.GetManifest(pkg)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("\n\033[32;1m=== \033[0;1m%s: \033[0m%s\n", pkgManifest.Name, pkgManifest.Version)
	output += fmt.Sprintln(pkgManifest.Description)
	output += fmt.Sprintf("\033[34;4m%s\033[0m\n", pkgManifest.Homepage)
	output += fmt.Sprintf("From: \033[34;4m%s\033[0m\n", pkgManifest.ManifestUrl)

	if len(pkgManifest.Dependencies) > 0 {
		output += fmt.Sprintf("Dependencies: %s\n", strings.Join(pkgManifest.Dependencies, ", "))
	}

	if pkgManifest.Caveats != "" {
		output += util.WrapText(fmt.Sprintf("Caveats: %s\n", pkgManifest.Caveats), 90)
	}

	output += "\nInstall:\n"
	for _, line := range pkgManifest.Scripts.Install {
		output += fmt.Sprintf("  %s\n", util.SyntaxHighlight(line))
	}
	if len(pkgManifest.Scripts.Completions) > 0 {
		output += "Completions:\n"
		for _, line := range pkgManifest.Scripts.Completions {
			output += fmt.Sprintf("  %s\n", util.SyntaxHighlight(line))
		}
	}

	return output, nil
}
