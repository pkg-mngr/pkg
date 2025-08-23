package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/melbahja/got"
	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
	"github.com/pkg-mngr/pkg/internal/manifest"
	"github.com/pkg-mngr/pkg/internal/util"
)

func Add(pkg string, skipConfirmation bool, lockfile config.Lockfile) error {
	manifestJson, err := manifest.GetManifest(pkg)
	if err != nil {
		return err
	}
	pkgManifest, err := manifestJson.Process()
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
			Add(dep, skipConfirmation, lockfile)
		}
	}

	fmt.Printf("Installing %s...\n", pkg)

	filesBefore, err := listFiles()
	if err != nil {
		return err
	}
	if err := fetchPackage(pkgManifest); err != nil {
		return err
	}

	fmt.Println("Running install script...")
	installScript := strings.Join(pkgManifest.Scripts.Install, "\n")
	if _, err := util.RunScript(installScript, skipConfirmation); err != nil && err.Error() != "" {
		return err
	}

	if len(pkgManifest.Scripts.Completions) != 0 {
		fmt.Println("Running completions script...")
		completionsScript := strings.Join(pkgManifest.Scripts.Completions, "\n")
		if _, err := util.RunScript(completionsScript, skipConfirmation); err != nil && err.Error() != "" {
			return err
		}
	}

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

	files, err := os.ReadDir(config.PKG_TMP)
	if err != nil {
		return fmt.Errorf("Error opening %s directory: %v\n", config.PKG_TMP, err)
	}
	for _, file := range files {
		filename := filepath.Join(config.PKG_TMP, file.Name())
		if err := os.RemoveAll(filename); err != nil {
			return fmt.Errorf("Error deleting %s: %v\n", filename, err)
		}
	}

	if pkgManifest.Caveats != "" {
		fmt.Printf("\nCaveats:\n %s\n\n", pkgManifest.Caveats)
	}
	fmt.Printf("Finished installing %s.\n", pkg)

	return nil
}

func listFiles() ([]string, error) {
	files := []string{}

	entries, err := os.ReadDir(config.PKG_BIN)
	if err != nil {
		return nil, fmt.Errorf("Error listing %s directory: %v", config.PKG_BIN, err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("bin", entry.Name()))
	}

	entries, err = os.ReadDir(config.PKG_OPT)
	if err != nil {
		return nil, fmt.Errorf("Error listing %s directory: %v\n", config.PKG_OPT, err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("opt", entry.Name()))
	}

	entries, err = os.ReadDir(config.PKG_ZSH_COMPLETIONS)
	if err != nil {
		return nil, fmt.Errorf("Error listing %s directory: %v\n", config.PKG_ZSH_COMPLETIONS, err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("share/zsh/site-functions", entry.Name()))
	}

	return files, nil
}

func fetchPackage(pkgManifest manifest.Manifest) error {
	filename := filepath.Join(config.PKG_TMP, path.Base(pkgManifest.Url))

	g := got.New()
	g.ProgressFunc = func(d *got.Download) {
		percent := float64(d.Size()) / float64(d.TotalSize()) * 100
		speed := float64(d.AvgSpeed())
		speedStr := fmt.Sprintf("%.2f kB/s", speed/1024)
		if speed/1024/1024 >= 5 {
			speedStr = fmt.Sprintf("%.2f MB/s", speed/1024/1024)
		}
		fmt.Printf("\033[2KDownloaded %.2f%% (%s)\r", percent, speedStr)
	}
	if err := g.Download(pkgManifest.Url, filename); err != nil {
		return fmt.Errorf("%s: Error while downloading %s: %v", pkgManifest.Name, filename, err)
	}
	fmt.Println("\033[2KDownloaded 100%")

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s: Error reading data from response", pkgManifest.Name)
	}

	fmt.Print("Verifying checksum...")
	checksum := fmt.Sprintf("%x", sha256.Sum256(data))
	if checksum != pkgManifest.Sha256 {
		return fmt.Errorf("\n%s: Checksum of data did not match in package manifest", pkgManifest.Name)
	}
	fmt.Println(" Looks good!")

	return nil
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
