package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func Remove(pkg string, lockfile config.Lockfile, isForUpdate bool) {
	if _, ok := lockfile[pkg]; !ok {
		fmt.Println(pkg + " is not installed.")
		return
	}
	if !isForUpdate {
		for installed := range lockfile {
			if slices.Contains(lockfile[installed].Dependencies, pkg) {
				fmt.Println("Cannot uninstall " + pkg + " as it is a dependency of " + installed)
				return
			}
		}
	}

	fmt.Println("Removing " + pkg + "...")
	removeFiles(lockfile[pkg].Files)

	for _, dep := range lockfile[pkg].Dependencies {
		Remove(dep, lockfile, isForUpdate)
	}

	lockfile.Remove(pkg)
}

func removeFiles(files []string) {
	for _, file := range files {
		fmt.Printf("Deleting %s...\n", file)
		if err := os.RemoveAll(filepath.Join(config.PKG_HOME(), file)); err != nil {
			log.Errorf("Error removing file %s: %v\n", filepath.Join(config.PKG_HOME(), file), err)
		}
	}
}
