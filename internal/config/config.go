package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func PKG_HOME() (string, error) {
	pkgHome := os.Getenv("PKG_HOME")
	if pkgHome != "" {
		path, err := filepath.Abs(pkgHome)
		if err != nil {
			return "", fmt.Errorf("Error making absolute path to .pkg directory: %v", err)
		}
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("User home directory not found: %v", err)
	}

	pkgHome = filepath.Join(home, ".pkg")
	return pkgHome, nil
}

func PKG_BIN() (string, error) {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return "", err
	}
	return filepath.Join(pkgHome, "bin"), nil
}

func PKG_OPT() (string, error) {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return "", err
	}
	return filepath.Join(pkgHome, "opt"), nil
}

func PKG_ZSH_COMPLETIONS() (string, error) {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return "", err
	}
	return filepath.Join(pkgHome, "share/zsh/site-functions"), nil
}

func PKG_TMP() (string, error) {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return "", err
	}
	return filepath.Join(pkgHome, "tmp"), nil
}

func LOCKFILE() (string, error) {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return "", err
	}
	return filepath.Join(pkgHome, "pkg.lock"), nil
}

func MANIFEST_HOST() string {
	pkgManifestHost := os.Getenv("PKG_MANIFEST_HOST")
	if pkgManifestHost != "" {
		return pkgManifestHost
	}

	return "https://pkg.zerolimits.dev"
}
