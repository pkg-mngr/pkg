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

	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)
		if pkgManifest.Version == lockfile.Packages[pkg].Version {
			fmt.Println(pkg + " is already up to date")
			continue
		}

		fmt.Println("Updating " + pkg + "...")
		Remove(pkg)
		Add(pkg)
	}
}
