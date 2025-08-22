package config

import (
	"os"
	"path/filepath"

	"github.com/pkg-mngr/pkg/internal/log"
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
