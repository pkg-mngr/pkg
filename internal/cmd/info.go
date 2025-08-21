package cmd

import (
	"fmt"
	"strings"

	"github.com/noclaps/pkg/internal/manifest"
	"github.com/noclaps/pkg/internal/util"
)

func Info(pkg string) string {
	pkgManifest := manifest.GetManifest(pkg)
	output := fmt.Sprintf("\n\033[32;1m=== \033[0;1m%s: \033[0m%s\n", pkgManifest.Name, pkgManifest.Version)
	output += fmt.Sprintln(pkgManifest.Description)
	output += fmt.Sprintf("\033[34;4m%s\033[0m\n", pkgManifest.Homepage)
	output += fmt.Sprintf("From: \033[34;4m%s\033[0m\n", pkgManifest.ManifestUrl)

	// Show supported platforms
	output += "\033[33;1mSupported Platforms:\033[0m\n"
	supportedPlatforms := []string{}
	for platform := range pkgManifest.Sha256 {
		supportedPlatforms = append(supportedPlatforms, platform)
	}
	if len(supportedPlatforms) > 0 {
		output += "  " + strings.Join(supportedPlatforms, ", ") + "\n"
	} else {
		output += "  None specified\n"
	}
	output += "\n"

	// Show platform-specific details
	output += "\033[33;1mPlatform Details:\033[0m\n"
	for platform := range pkgManifest.Sha256 {
		output += fmt.Sprintf("  \033[1m%s:\033[0m\n", platform)
		if url, exists := pkgManifest.Url[platform]; exists {
			output += fmt.Sprintf("    URL: %s\n", url)
		}
		if sha256, exists := pkgManifest.Sha256[platform]; exists {
			output += fmt.Sprintf("    SHA256: %s\n", sha256)
		}
		output += "\n"

		if len(pkgManifest.GetInstallScripts()) > 0 {
			output += "    Install:\n"
			for _, line := range pkgManifest.GetInstallScripts() {
				output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
			}
		}

		if len(pkgManifest.GetCompletionsScripts()) > 0 {
			output += "    Completions:\n"
			for _, line := range pkgManifest.GetCompletionsScripts() {
				output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
			}
		}

		if len(pkgManifest.GetLatestScripts()) > 0 {
			output += "    Latest:\n"
			for _, line := range pkgManifest.GetLatestScripts() {
				output += fmt.Sprintf("      %s\n", util.SyntaxHighlight(line))
			}
		}

		output += "\n"
	}

	// Show dependencies
	if len(pkgManifest.Dependencies) > 0 {
		output += fmt.Sprintf("Dependencies: %s\n", strings.Join(pkgManifest.Dependencies, ", "))
	}

	// Show caveats
	if pkgManifest.Caveats != "" {
		output += util.WrapText(fmt.Sprintf("Caveats: %s\n", pkgManifest.Caveats), 90)
	}

	// Show current platform compatibility
	if err := pkgManifest.ValidatePlatformSupport(); err != nil {
		output += "\033[31;1mPlatform Compatibility:\033[0m\n"
		output += "  \033[31m✗ Not compatible with current platform\033[0m\n"
	} else {
		output += "\033[32;1mPlatform Compatibility:\033[0m\n"
		output += "  \033[32m✓ Compatible with current platform\033[0m\n"
	}

	return output
}
