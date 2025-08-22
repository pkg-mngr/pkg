package cmd

import (
	"fmt"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/manifest"
)

func Update(pkgs []string, skipConfirmation bool, lockfile config.Lockfile) error {
	allUpToDate := true

	for _, pkg := range pkgs {
		pkgManifest, err := manifest.GetManifest(pkg)
		if err != nil {
			return err
		}
		if pkgManifest.Version == lockfile[pkg].Version {
			continue
		}

		allUpToDate = false
		fmt.Printf("Updating %s...\n", pkg)
		if err := Remove(pkg, lockfile, true); err != nil {
			return err
		}
		if err := Add(pkg, skipConfirmation, lockfile); err != nil {
			return err
		}
	}

	if allUpToDate {
		fmt.Println("All packages are up to date")
	}
	return nil
}
