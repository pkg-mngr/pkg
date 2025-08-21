package cmd

import (
	"fmt"
	"strings"

	"github.com/noclaps/pkg/internal/manifest"
	"github.com/noclaps/pkg/internal/platforms"
	"github.com/noclaps/pkg/internal/util"
)

func Info(pkg string) string {
	pkgManifest := manifest.GetManifest(pkg)
	currentPlatform := platforms.GetPlatform()

	if pkgManifest.Url[string(currentPlatform)] == "" {
		return fmt.Sprintf("\033[31;1mError:\033[0m Package '%s' is not available for platform '%s'.", pkg, currentPlatform)
	}

	output := fmt.Sprintf("\n\033[32;1m=== \033[0;1m%s: \033[0m%s\n", pkgManifest.Name, pkgManifest.Version)
	output += fmt.Sprintln(pkgManifest.Description)
	output += fmt.Sprintf("\033[34;4m%s\033[0m\n", pkgManifest.Homepage)
	output += fmt.Sprintf("From: \033[34;4m%s\033[0m\n", pkgManifest.ManifestUrl)

	// Show platform-specific details
	output += "\033[33;1mPlatform Details:\033[0m\n"
	output += fmt.Sprintf("  \033[1m%s:\033[0m\n", currentPlatform)
	if url, exists := pkgManifest.Url[string(currentPlatform)]; exists {
		output += fmt.Sprintf("    URL: %s\n", url)
	}
	if sha256, exists := pkgManifest.Sha256[string(currentPlatform)]; exists {
		output += fmt.Sprintf("    SHA256: %s\n", sha256)
	}
	output += "\n"

	if len(pkgManifest.GetInstallScripts(currentPlatform)) > 0 {
		output += "    Install:\n"
		for _, line := range pkgManifest.GetInstallScripts(currentPlatform) {
			output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
		}
	}

	if len(pkgManifest.GetCompletionsScripts(currentPlatform)) > 0 {
		output += "    Completions:\n"
		for _, line := range pkgManifest.GetCompletionsScripts(currentPlatform) {
			output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
		}
	}

	if len(pkgManifest.GetLatestScripts(currentPlatform)) > 0 {
		output += "    Latest:\n"
		for _, line := range pkgManifest.GetLatestScripts(currentPlatform) {
			output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
		}
	}

	// Show dependencies
	if len(pkgManifest.Dependencies) > 0 {
		output += "\n"
		output += fmt.Sprintf("Dependencies: %s\n", strings.Join(pkgManifest.Dependencies, ", "))
	}

	// Show caveats
	if pkgManifest.Caveats != "" {
		output += "\n"
		output += util.WrapText(fmt.Sprintf("Caveats: %s\n", pkgManifest.Caveats), 90)
	}

	return output
}
