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
	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
	"github.com/noclaps/pkg/internal/util"
)

func Add(pkg string, skipConfirmation bool, lockfile config.Lockfile) {
	pkgManifest := manifest.GetManifest(pkg)
	pkg = pkgManifest.Name

	// Validate platform support
	if err := pkgManifest.ValidatePlatformSupport(); err != nil {
		log.Errorf("Validation failed: %v\n", err)
		return
	}

	if entry, ok := lockfile[pkg]; ok {
		wasDep := false
		for installed, data := range lockfile {
			// if package was installed as a dependency, remove it from the dependency
			// list of all those packages
			if i := slices.Index(data.Dependencies, pkg); i != -1 {
				data.Dependencies = slices.Delete(data.Dependencies, i, i+1)
				lockfile[installed] = data
				wasDep = true
				lockfile.Write()
			}
		}

		if wasDep || entry.Version == pkgManifest.Version {
			fmt.Printf("%s is already installed.\n", pkg)
			return
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

	// Show which platform configuration is being used
	if _, exists := pkgManifest.GetPlatformConfig(); exists {
		currentPlatform := config.GetCurrentPlatform()
		fmt.Printf("Using platform-specific configuration for: %s\n", currentPlatform)
	} else {
		fmt.Println("Using fallback configuration")
	}

	filesBefore := listFiles()
	if err := fetchPackage(pkgManifest); err != nil {
		log.Errorf("%v\n", err)
		return
	}

	fmt.Println("Running install script...")
	installScript := strings.Join(pkgManifest.GetInstallScripts(config.GetCurrentPlatform()), "\n")
	if _, err := util.RunScript(installScript, skipConfirmation); err != nil && err.Error() != "" {
		log.Errorf("%v\n", err)
		return
	}

	if len(pkgManifest.GetCompletionsScripts(config.GetCurrentPlatform())) != 0 {
		fmt.Println("Running completions script...")
		completionsScript := strings.Join(pkgManifest.GetCompletionsScripts(config.GetCurrentPlatform()), "\n")
		if _, err := util.RunScript(completionsScript, skipConfirmation); err != nil && err.Error() != "" {
			log.Errorf("%v\n", err)
			return
		}
	}

	filesAfter := listFiles()

	lockfile.NewEntry(
		pkgManifest.Name,
		pkgManifest.ManifestUrl,
		pkgManifest.Version,
		dependencies,
		diffFiles(filesBefore, filesAfter),
	)

	files, err := os.ReadDir(config.PKG_TMP())
	if err != nil {
		log.Errorf("Error opening %s directory: %v\n", config.PKG_TMP(), err)
	}
	for _, file := range files {
		filename := filepath.Join(config.PKG_TMP(), file.Name())
		if err := os.RemoveAll(filename); err != nil {
			log.Errorf("Error deleting %s: %v\n", filename, err)
		}
	}

	if pkgManifest.Caveats != "" {
		fmt.Printf("\nCaveats:\n %s\n\n", pkgManifest.Caveats)
	}
	fmt.Printf("Finished installing %s.\n", pkg)
}

func listFiles() []string {
	files := []string{}

	entries, err := os.ReadDir(config.PKG_BIN())
	if err != nil {
		log.Fatalf("Error listing %s directory: %v\n", config.PKG_BIN(), err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("bin", entry.Name()))
	}

	entries, err = os.ReadDir(config.PKG_OPT())
	if err != nil {
		log.Fatalf("Error listing %s directory: %v\n", config.PKG_OPT(), err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("opt", entry.Name()))
	}

	entries, err = os.ReadDir(config.PKG_ZSH_COMPLETIONS())
	if err != nil {
		log.Fatalf("Error listing %s directory: %v\n", config.PKG_ZSH_COMPLETIONS(), err)
	}
	for _, entry := range entries {
		files = append(files, filepath.Join("share/zsh/site-functions", entry.Name()))
	}

	return files
}

func fetchPackage(pkgManifest manifest.Manifest) error {
	url := pkgManifest.GetURL()
	expectedSHA256 := pkgManifest.GetSHA256()

	if url == "" {
		return fmt.Errorf("%s: No URL found for current platform", pkgManifest.Name)
	}

	if expectedSHA256 == "" {
		return fmt.Errorf("%s: No SHA256 found for current platform", pkgManifest.Name)
	}

	filename := filepath.Join(config.PKG_TMP(), path.Base(url))

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
	if err := g.Download(url, filename); err != nil {
		return fmt.Errorf("%s: Error while downloading %s: %v", pkgManifest.Name, filename, err)
	}
	fmt.Println("\033[2KDownloaded 100%")

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s: Error reading data from response", pkgManifest.Name)
	}

	fmt.Print("Verifying checksum...")
	checksum := fmt.Sprintf("%x", sha256.Sum256(data))
	if checksum != expectedSHA256 {
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
