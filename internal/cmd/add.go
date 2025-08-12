package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/melbahja/got"
	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
)

func Add(pkg string) {
	lockfile := config.ReadLockfile()

	pkgManifest := manifest.GetManifest(pkg)

	if entry, ok := lockfile.Packages[pkg]; ok {
		if entry.Version == pkgManifest.Version {
			fmt.Println(pkg + " is already installed.")
			return
		}
	}

	fmt.Println("Installing " + pkg + "...")

	filesBefore := listFiles()
	if err := fetchPackage(pkgManifest); err != nil {
		log.Errorln(err)
		return
	}

	fmt.Println("Running install script...")
	installScript := strings.Join(pkgManifest.Scripts.Install, "\n")
	installScript = "cd " + config.PKG_TMP() + "\n" + installScript
	if err := runScript(installScript); err != nil && err.Error() != "" {
		log.Errorln(err)
		return
	}

	if len(pkgManifest.Scripts.Completions) != 0 {
		fmt.Println("Running completions script...")
		completionsScript := strings.Join(pkgManifest.Scripts.Completions, "\n")
		completionsScript = "cd " + config.PKG_TMP() + "\n" + completionsScript
		if err := runScript(completionsScript); err != nil && err.Error() != "" {
			log.Errorln(err)
			return
		}
	}

	filesAfter := listFiles()

	lockfile.WriteToLockfile(
		pkgManifest.Name,
		pkgManifest.ManifestUrl,
		pkgManifest.Version,
		diffFiles(filesBefore, filesAfter),
	)

	files, err := os.ReadDir(config.PKG_TMP())
	if err != nil {
		log.Errorln("Error opening " + config.PKG_TMP() + " directory")
	}
	for _, file := range files {
		filename := filepath.Join(config.PKG_TMP(), file.Name())
		if err := os.RemoveAll(filename); err != nil {
			log.Errorln("Error deleting " + filename)
		}
	}

	if pkgManifest.Caveats != "" {
		fmt.Println("\nCaveats:\n" + pkgManifest.Caveats + "\n")
	}
	fmt.Println("Finished installing " + pkg + ".")
}

func listFiles() []string {
	files := []string{}

	entries, err := os.ReadDir(config.PKG_BIN())
	if err != nil {
		log.Fatalln("Error listing " + config.PKG_BIN() + " directory")
	}
	for _, entry := range entries {
		files = append(files, "bin/"+entry.Name())
	}

	entries, err = os.ReadDir(config.PKG_OPT())
	if err != nil {
		log.Fatalln("Error listing " + config.PKG_OPT() + " directory")
	}
	for _, entry := range entries {
		files = append(files, "opt/"+entry.Name())
	}

	entries, err = os.ReadDir(config.PKG_ZSH_COMPLETIONS())
	if err != nil {
		log.Fatalln("Error listing " + config.PKG_ZSH_COMPLETIONS() + " directory")
	}
	for _, entry := range entries {
		files = append(files, "share/zsh/site-functions/"+entry.Name())
	}

	return files
}

func fetchPackage(pkgManifest manifest.Manifest) error {
	filename := filepath.Join(config.PKG_TMP(), path.Base(pkgManifest.Url))

	g := got.New()
	g.ProgressFunc = func(d *got.Download) {
		percent := float64(d.Size()) / float64(d.TotalSize()) * 100
		speed := float64(d.AvgSpeed())
		speedStr := ""
		if speed/1024/1024 < 5 {
			speedStr = fmt.Sprintf("%.2f kB/s", speed/1024)
		} else {
			speedStr = fmt.Sprintf("%.2f MB/s", speed/1024/1024)
		}
		fmt.Printf("\033[2KDownloaded %.2f%% (%s)\r", percent, speedStr)
	}
	if err := g.Download(pkgManifest.Url, filename); err != nil {
		return fmt.Errorf("Error while downloading %s: %v", filename, err)
	}
	fmt.Println("\033[2KDownloaded 100%")

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading data from response")
	}

	fmt.Print("Verifying checksum...")
	checksum := fmt.Sprintf("%x", sha256.Sum256(data))
	if checksum != pkgManifest.Sha256 {
		return fmt.Errorf("\nChecksum of data did not match in package manifest")
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

func runScript(script string) error {
	cmd := exec.Command("/bin/sh", "-c", script)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error getting stderr pipe")
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error while starting command")
	}

	stderrData, err := io.ReadAll(stderr)
	if err != nil {
		return fmt.Errorf("Error getting data from stderr")
	}

	if err := cmd.Wait(); err != nil {
		log.Errorln("Error waiting for command")
		return nil
	}

	return fmt.Errorf("%s", stderrData)
}
