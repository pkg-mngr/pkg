package cmd

import (
	"fmt"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/manifest"
)

func Update(pkgs []string, skipConfirmation bool, lockfile config.Lockfile) {
	allUpToDate := true

	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)
		if pkgManifest.Version == lockfile[pkg].Version {
			continue
		}

		allUpToDate = false
		fmt.Println("Updating " + pkg + "...")
		Remove(pkg, lockfile, true)
		Add(pkg, skipConfirmation, lockfile)
	}

	if allUpToDate {
		fmt.Println("All packages are up to date")
	}
}
