package cmd

import (
	"maps"
	"slices"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
)

func Update() {
	lockfile := config.ReadLockfile()
	pkgs := slices.Collect(maps.Keys(lockfile.Packages))

	pkgsToUpdate := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)
		if pkgManifest.Version == lockfile.Packages[pkg].Version {
			log.Println(pkg + " is already up to date")
			continue
		}

		pkgsToUpdate = append(pkgsToUpdate, pkg)
	}

	Remove(pkgsToUpdate)
	Add(pkgsToUpdate)
}
