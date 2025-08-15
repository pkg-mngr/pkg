package cmd

import (
	"fmt"
	"maps"
	"slices"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/manifest"
)

func Update() {
	lockfile := config.ReadLockfile()
	pkgs := slices.Collect(maps.Keys(lockfile.Packages))
	allUpToDate := true

	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)
		if pkgManifest.Version == lockfile.Packages[pkg].Version {
			continue
		}

		allUpToDate = false
		fmt.Println("Updating " + pkg + "...")
		Remove(pkg)
		Add(pkg)
	}

	if allUpToDate {
		fmt.Println("All packages are up to date")
	}
}
