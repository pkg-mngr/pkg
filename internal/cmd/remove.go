package cmd

import (
	"fmt"
	"os"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func Remove(pkg string) {
	lockfile := config.ReadLockfile()

	if _, ok := lockfile.Packages[pkg]; !ok {
		fmt.Println(pkg + " is not installed.")
		return
	}

	fmt.Println("Removing " + pkg + "...")
	removeFiles(lockfile.Packages[pkg].Files)
	lockfile.RemoveFromLockfile(pkg)
}

func removeFiles(files []string) {
	for _, file := range files {
		fmt.Println("Deleting " + file + "...")
		if err := os.RemoveAll(config.PKG_HOME() + "/" + file); err != nil {
			log.Errorln("Error removing file " + config.PKG_HOME() + "/" + file)
		}
	}
}
