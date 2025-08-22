package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/pkg-mngr/pkg/internal/config"
)

func Remove(pkg string, lockfile config.Lockfile, isForUpdate bool) error {
	if _, ok := lockfile[pkg]; !ok {
		return ErrorPackageNotInstalled{Name: pkg}
	}
	if !isForUpdate {
		for installed := range lockfile {
			if slices.Contains(lockfile[installed].Dependencies, pkg) {
				return ErrorPackageDependencyOf{Name: pkg, Dependent: installed}
			}
		}
	}

	fmt.Printf("Removing %s...\n", pkg)
	if err := removeFiles(lockfile[pkg].Files); err != nil {
		return err
	}

	for _, dep := range lockfile[pkg].Dependencies {
		if err := Remove(dep, lockfile, isForUpdate); err != nil {
			return err
		}
	}

	return lockfile.Remove(pkg)
}

func removeFiles(files []string) error {
	for _, file := range files {
		fmt.Printf("Deleting %s...\n", file)
		pkgHome, err := config.PKG_HOME()
		if err != nil {
			return err
		}
		if err := os.RemoveAll(filepath.Join(pkgHome, file)); err != nil {
			return fmt.Errorf("Error removing file %s: %v\n", filepath.Join(pkgHome, file), err)
		}
	}

	return nil
}
