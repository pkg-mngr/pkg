package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
	"github.com/pkg-mngr/pkg/internal/manifest"
	"github.com/pkg-mngr/pkg/internal/util"
)

func Add(pkg string, skipConfirmation bool, lockfile config.Lockfile) error {
	pkgManifest, err := manifest.Get(pkg)
	if err != nil {
		return err
	}
	pkg = pkgManifest.Name

	if entry, ok := lockfile[pkg]; ok {
		wasDep := false
		for installed, data := range lockfile {
			// if package was installed as a dependency, remove it from the dependency
			// list of all those packages
			if i := slices.Index(data.Dependencies, pkg); i != -1 {
				data.Dependencies = slices.Delete(data.Dependencies, i, i+1)
				lockfile[installed] = data
				wasDep = true
			}
		}

		if wasDep || entry.Version == pkgManifest.Version {
			log.Printf("%s is already installed", pkg)
			return nil
		}
	}

	fmt.Printf("Installing %s...\n", pkg)

	filename := filepath.Join(config.PKG_TMP, path.Base(pkgManifest.Url))

	// download from url
	err = util.Fetch(pkgManifest.Url, filename, pkgManifest.Name)
	if err != nil {
		return err
	}

	// verify checksum
	err = util.VerifyChecksum(filename, pkgManifest.Sha256, pkgManifest.Name)
	if err != nil {
		return err
	}

	if err := install(lockfile, pkgManifest, skipConfirmation); err != nil {
		return err
	}
	fmt.Printf("Finished installing %s\n", pkg)

	return nil
}

func listFiles() ([]string, error) {
	files := []string{}
	pkgDirs := map[string]string{
		config.PKG_BIN:             "bin",
		config.PKG_OPT:             "opt",
		config.PKG_ZSH_COMPLETIONS: "share/zsh/site-functions",
	}
	for dir, dirName := range pkgDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("Error listing %s directory: %v", dir, err)
		}
		for _, entry := range entries {
			files = append(files, filepath.Join(dirName, entry.Name()))
		}
	}

	return files, nil
}

func diffFiles(before, after []string) []string {
	diff := []string{}
	for _, item := range after {
		if !slices.Contains(before, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

func install(lockfile config.Lockfile, pkgManifest manifest.Manifest, skipConfirmation bool) error {
	// skip adding dependencies that are already installed. these were installed
	// either due to some other package or manually by the user, so they're not
	// lockfile dependencies of this package
	dependencies := slices.DeleteFunc(pkgManifest.Dependencies, func(dep string) bool {
		_, ok := lockfile[dep]
		return ok
	})
	if len(dependencies) > 0 {
		fmt.Println("Installing dependencies...")
		for _, dep := range dependencies {
			// add dependencies before in case they're needed for installation of
			// current package
			Add(dep, skipConfirmation, lockfile)
		}
	}

	// list files before installation
	filesBefore, err := listFiles()
	if err != nil {
		return err
	}

	// run install and completions scripts
	fmt.Println("Running install script...")
	installScript := strings.Join(pkgManifest.Scripts.Install, "\n")
	_, err = util.RunScript(installScript, skipConfirmation)
	if err != nil {
		return err
	}
	if len(pkgManifest.Scripts.Completions) != 0 {
		fmt.Println("Running completions script...")
		completionsScript := strings.Join(pkgManifest.Scripts.Completions, "\n")
		_, err := util.RunScript(completionsScript, skipConfirmation)
		if err != nil {
			return err
		}
	}

	// list files after installation
	filesAfter, err := listFiles()
	if err != nil {
		return err
	}

	// add to lockfile
	lockfile[pkgManifest.Name] = config.LockfilePackage{
		Manifest:     pkgManifest.ManifestUrl,
		Version:      pkgManifest.Version,
		Dependencies: dependencies,
		Files:        diffFiles(filesBefore, filesAfter),
	}

	// delete and recreate .pkg/tmp
	if err := os.RemoveAll(config.PKG_TMP); err != nil {
		return fmt.Errorf("Error deleting %s: %v\n", config.PKG_TMP, err)
	}
	if err := os.Mkdir(config.PKG_TMP, 0o755); err != nil {
		return fmt.Errorf("Error creating %s: %v\n", config.PKG_TMP, err)
	}

	if pkgManifest.Caveats != "" {
		fmt.Printf("\nCaveats:\n %s\n\n", pkgManifest.Caveats)
	}

	return nil
}
