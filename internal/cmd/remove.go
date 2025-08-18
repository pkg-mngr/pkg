package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func Remove(pkg string, lockfile config.Lockfile) {
	if _, ok := lockfile[pkg]; !ok {
		fmt.Println(pkg + " is not installed.")
		return
	}

	fmt.Println("Removing " + pkg + "...")
	removeFiles(lockfile[pkg].Files)

	for _, dep := range lockfile[pkg].Dependencies {
		Remove(dep, lockfile)
	}

	lockfile.Remove(pkg)
}

func removeFiles(files []string) {
	for _, file := range files {
		fmt.Println("Deleting " + file + "...")
		if err := os.RemoveAll(filepath.Join(config.PKG_HOME(), file)); err != nil {
			log.Errorln("Error removing file " + filepath.Join(config.PKG_HOME(), file))
		}
	}
}
