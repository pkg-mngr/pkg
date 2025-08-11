package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
)

func Add(pkgs []string) {
	lockfile := config.ReadLockfile()

	for _, pkg := range pkgs {
		pkgManifest := manifest.GetManifest(pkg)

		if entry, ok := lockfile.Packages[pkg]; ok {
			if entry.Version == pkgManifest.Version {
				log.Errorln(pkg + " is already installed.")
				continue
			}
		}

		fmt.Println("Installing " + pkg)

		filesBefore := listFiles()
		if err := fetchPackage(pkgManifest); err != nil {
			log.Println(err)
			continue
		}

		installScript := strings.Join(pkgManifest.Scripts.Install, "\n")
		installScript = "cd " + config.PKG_TMP() + "\n" + installScript
		if err := runScript(installScript); err != nil && err.Error() != "" {
			log.Println(err)
			continue
		}

		if len(pkgManifest.Scripts.Completions) != 0 {
			completionsScript := strings.Join(pkgManifest.Scripts.Completions, "\n")
			if err := runScript(completionsScript); err != nil && err.Error() != "" {
				log.Println(err)
				continue
			}
		}

		filesAfter := listFiles()

		lockfile.WriteToLockfile(
			pkgManifest.Name,
			pkgManifest.ManifestUrl,
			pkgManifest.Version,
			diffFiles(filesBefore, filesAfter),
		)

		urlParts := strings.Split(pkgManifest.Url, "/")
		filename := config.PKG_TMP() + "/" + urlParts[len(urlParts)-1]
		os.RemoveAll(filename)
	}
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
	filename := config.PKG_TMP() + "/" + path.Base(pkgManifest.Url)

	if _, err := os.Stat(filename); err != nil {
		res, err := http.Get(pkgManifest.Url)
		if err != nil || res.StatusCode != http.StatusOK {
			return fmt.Errorf("Error fetching data from %s", pkgManifest.Url)
		}

		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("Error creating temporary file")
		}

		if _, err := io.Copy(f, res.Body); err != nil {
			return fmt.Errorf("Error writing data to file")
		}
		res.Body.Close()
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading data from response")
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))
	if checksum != pkgManifest.Sha256 {
		return fmt.Errorf("Checksum of data did not match in package manifest")
	}

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
		log.Println("Error waiting for command")
		return nil
	}

	return fmt.Errorf("%s", stderrData)
}
