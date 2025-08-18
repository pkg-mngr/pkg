package cmd

import (
	"strings"

	"github.com/noclaps/pkg/internal/manifest"
)

func Info(pkg string) string {
	pkgManifest := manifest.GetManifest(pkg)
	output := "\n\033[32;1m=== \033[0;1m" + pkgManifest.Name + ": \033[0m" + pkgManifest.Version + "\n"
	output += pkgManifest.Description + "\n"
	output += "\033[34;4m" + pkgManifest.Homepage + "\033[0m\n"
	output += "From: \033[34;4m" + pkgManifest.ManifestUrl + "\033[0m\n"

	if len(pkgManifest.Dependencies) > 0 {
		output += "Dependencies: " + strings.Join(pkgManifest.Dependencies, ", ") + "\n"
	}

	if pkgManifest.Caveats != "" {
		output += "Caveats: " + pkgManifest.Caveats + "\n"
	}

	return output
}
