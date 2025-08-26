package cmd

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/manifest"
	"github.com/pkg-mngr/pkg/internal/util"
)

func Update(pkgs []string, skipConfirmation bool, lockfile config.Lockfile) error {
	allUpToDate := true

	for _, pkg := range pkgs {
		if _, ok := lockfile[pkg]; !ok {
			return ErrorPackageNotInstalled{Name: pkg}
		}

		manifestUrl := lockfile[pkg].Manifest
		var pkgManifest manifest.Manifest
		if manifest.IsLocalFile(manifestUrl) {
			manifestJson, err := manifest.FromFile(manifestUrl)
			if err != nil {
				return err
			}
			processed, err := manifestJson.Process()
			if err != nil {
				return err
			}
			pkgManifest = processed
		} else {
			manifestJson, err := manifest.FromRemote(manifestUrl)
			if err != nil {
				return err
			}
			processed, err := manifestJson.Process()
			if err != nil {
				return err
			}
			pkgManifest = processed
		}

		if pkgManifest.Version == lockfile[pkg].Version {
			continue
		}

		allUpToDate = false
		fmt.Printf("Updating %s...\n", pkg)

		filename := filepath.Join(config.PKG_TMP, path.Base(pkgManifest.Url))

		// download file and verify checksum
		if err := util.Fetch(pkgManifest.Url, filename, pkgManifest.Name); err != nil {
			return err
		}
		if err := util.VerifyChecksum(filename, pkgManifest.Sha256, pkgManifest.Name); err != nil {
			return err
		}

		// if previous steps were successful, remove existing package
		if err := Remove(pkg, lockfile, true); err != nil {
			return err
		}

		// resume installation after removing old version
		if err := install(lockfile, pkgManifest, skipConfirmation); err != nil {
			return err
		}
		fmt.Printf("Finished updating %s\n", pkg)
	}

	if allUpToDate {
		fmt.Println("All packages are up to date")
	}
	return nil
}
