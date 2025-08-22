package cmd

import (
	"fmt"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/manifest"
)

func Update(pkgs []string, skipConfirmation bool, lockfile config.Lockfile) {
	allUpToDate := true

	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)
		if pkgManifest.Version == lockfile[pkg].Version {
			continue
		}

		allUpToDate = false
		fmt.Printf("Updating %s...\n", pkg)
		Remove(pkg, lockfile, true)
		Add(pkg, skipConfirmation, lockfile)
	}

	if allUpToDate {
		fmt.Println("All packages are up to date")
	}
}
