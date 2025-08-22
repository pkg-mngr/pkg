package config

import (
	"os"
	"path/filepath"

	"github.com/pkg-mngr/pkg/internal/log"
)

var (
	PKG_HOME            = getPkgHome()
	PKG_BIN             = filepath.Join(PKG_HOME, "bin")
	PKG_OPT             = filepath.Join(PKG_HOME, "opt")
	PKG_TMP             = filepath.Join(PKG_HOME, "tmp")
	LOCKFILE            = filepath.Join(PKG_HOME, "pkg.lock")
	PKG_ZSH_COMPLETIONS = filepath.Join(PKG_HOME, "share/zsh/site-functions")
	MANIFEST_HOST       = getManifestHost()
)

func getPkgHome() string {
	pkgHome := os.Getenv("PKG_HOME")
	if pkgHome != "" {
		path, err := filepath.Abs(pkgHome)
		if err != nil {
			// crash here because failing to make an absolute path means something has
			// gone very wrong in the system
			log.Fatalf("Error making absolute path to .pkg directory: %v", err)
		}
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// crash here because failing to get the home directory means something has
		// gone very wrong in the system
		log.Fatalf("User home directory not found: %v", err)
	}

	return filepath.Join(home, ".pkg")
}

func getManifestHost() string {
	pkgManifestHost := os.Getenv("PKG_MANIFEST_HOST")
	if pkgManifestHost != "" {
		return pkgManifestHost
	}

	return "https://pkg.zerolimits.dev"
}
