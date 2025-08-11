package config

import (
	"os"
	"path/filepath"

	"github.com/noclaps/pkg/internal/log"
)

func PKG_HOME() string {
	pkgHome := os.Getenv("PKG_HOME")
	if pkgHome != "" {
		path, err := filepath.Abs(pkgHome)
		if err != nil {
			log.Fatalln("Error making absolute path to .pkg directory")
		}
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("User home directory not found")
	}

	pkgHome = home + "/.pkg"
	return pkgHome
}

func PKG_BIN() string {
	return PKG_HOME() + "/bin"
}

func PKG_OPT() string {
	return PKG_HOME() + "/opt"
}

func PKG_ZSH_COMPLETIONS() string {
	return PKG_HOME() + "/share/zsh/site-functions"
}

func PKG_TMP() string {
	return PKG_HOME() + "/tmp"
}

func LOCKFILE() string {
	return PKG_HOME() + "/pkg.lock"
}

func MANIFEST_HOST() string {
	pkgManifestHost := os.Getenv("PKG_MANIFEST_HOST")
	if pkgManifestHost != "" {
		return pkgManifestHost
	}

	return "https://pkg.zerolimits.dev"
}
