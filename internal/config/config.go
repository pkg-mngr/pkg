package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/platforms"
)

func PKG_HOME() string {
	pkgHome := os.Getenv("PKG_HOME")
	if pkgHome != "" {
		path, err := filepath.Abs(pkgHome)
		if err != nil {
			log.Fatalf("Error making absolute path to .pkg directory: %v\n", err)
		}
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("User home directory not found: %v\n", err)
	}

	pkgHome = filepath.Join(home, ".pkg")
	return pkgHome
}

func PKG_BIN() string {
	return filepath.Join(PKG_HOME(), "bin")
}

func PKG_OPT() string {
	return filepath.Join(PKG_HOME(), "opt")
}

func PKG_ZSH_COMPLETIONS() string {
	return filepath.Join(PKG_HOME(), "share/zsh/site-functions")
}

func PKG_TMP() string {
	return filepath.Join(PKG_HOME(), "tmp")
}

func LOCKFILE() string {
	return filepath.Join(PKG_HOME(), "pkg.lock")
}

func MANIFEST_HOST() string {
	pkgManifestHost := os.Getenv("PKG_MANIFEST_HOST")
	if pkgManifestHost != "" {
		return pkgManifestHost
	}

	return "https://pkg.zerolimits.dev"
}

// GetCurrentPlatform returns the current platform identifier in the format "os-arch"
func GetCurrentPlatform() platforms.Platform {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Normalize architecture names to match common package naming conventions
	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		// macOS uses "arm64" but packages often use "arm64"
		if os == "darwin" {
			arch = "arm64"
		}
	}

	// Normalize OS names
	switch os {
	case "darwin":
		os = "macos"
	case "windows":
		os = "windows"
	case "linux":
		os = "linux"
	}

	return platforms.Platform{
		Name: os,
		Arch: arch,
	}
}
