package cmd

import (
	"fmt"
	"strings"

	"github.com/noclaps/pkg/internal/manifest"
)

func Info(pkg string) string {
	pkgManifest := manifest.GetManifest(pkg)
	output := fmt.Sprintf("\n\033[32;1m=== \033[0;1m%s: \033[0m%s\n", pkgManifest.Name, pkgManifest.Version)
	output += fmt.Sprintln(pkgManifest.Description)
	output += fmt.Sprintf("\033[34;4m%s\033[0m\n", pkgManifest.Homepage)
	output += fmt.Sprintf("From: \033[34;4m%s\033[0m\n", pkgManifest.ManifestUrl)

	if len(pkgManifest.Dependencies) > 0 {
		output += fmt.Sprintf("Dependencies: %s\n", strings.Join(pkgManifest.Dependencies, ", "))
	}

	if pkgManifest.Caveats != "" {
		output += "Caveats: " + pkgManifest.Caveats + "\n"
	}

	return output
}
