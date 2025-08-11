package cmd

import (
	"os"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func Remove(pkgs []string) {
	lockfile := config.ReadLockfile()

	for _, pkg := range pkgs {
		if _, ok := lockfile.Packages[pkg]; !ok {
			log.Errorln(pkg + " is not installed.")
			continue
		}
		removeFiles(lockfile.Packages[pkg].Files)
		lockfile.RemoveFromLockfile(pkg)
	}

}

func removeFiles(files []string) {
	for _, file := range files {
		if err := os.RemoveAll(config.PKG_HOME() + "/" + file); err != nil {
			log.Errorln("Error removing file " + config.PKG_HOME() + "/" + file)
		}
	}
}
