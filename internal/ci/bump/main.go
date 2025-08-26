package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
	"github.com/pkg-mngr/pkg/internal/manifest"
	"github.com/pkg-mngr/pkg/internal/util"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%v\n", err)
	}

	files, err := os.ReadDir("packages")
	if err != nil {
		log.Fatalf("Error reading packages/ directory: %v\n", err)
	}
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Go(func() {
			log.Printf("Reading %s\n", file.Name())

			pkgManifest, stderr := manifest.FromFile("./packages/" + file.Name())
			if stderr != nil {
				log.Errorf("Error getting manifest from %s: %v\n", file.Name(), stderr)
				return
			}

			latestScript := strings.Join(pkgManifest.Scripts.Latest, "\n")
			log.Printf("Running `latest` script: \n%s\n", latestScript)
			stdout, stderr := util.RunScript(latestScript, true)
			if stderr != nil {
				log.Errorf(
					"stdout: %s\nstderr: %v\n, Error running latest script in %s\n",
					stdout, stderr, file.Name())
				return
			}

			latestVersion := strings.TrimSpace(stdout)
			log.Printf("Found latest version: %s\n", latestVersion)
			if pkgManifest.Version == latestVersion {
				log.Printf("%s is already up to date\n", pkgManifest.Name)
				return
			}

			pkgManifest.Version = latestVersion

			for platform := range pkgManifest.Url {
				url := strings.ReplaceAll(pkgManifest.Url[platform], "{{ version }}", pkgManifest.Version)

				log.Printf("Fetching file from %s\n", url)
				filename := filepath.Join(config.PKG_TMP, path.Base(url))
				if err := util.Fetch(url, filename, pkgManifest.Name); err != nil {
					log.Errorf("Error fetching from %s: %v\n", url, err)
					return
				}

				data, err := os.ReadFile(filename)
				if err != nil {
					log.Errorf("Error reading data from %s: %v\n", filename, err)
				}

				log.Printf("Validating checksum for %s\n", path.Base(url))
				checksum := fmt.Sprintf("%x", sha256.Sum256(data))
				pkgManifest.Sha256[platform] = checksum
			}

			log.Printf("Updating %s\n", file.Name())
			f, stderr := os.Create(filepath.Join("packages", file.Name()))
			if stderr != nil {
				log.Errorf("Error creating file %s: %v\n", file.Name(), stderr)
				return
			}
			encoder := json.NewEncoder(f)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(*pkgManifest); err != nil {
				log.Errorf("Error writing manifest to %s\n", file.Name())
			}
			f.Close()

			log.Printf("Done updating %s\n", file.Name())
		})
	}

	wg.Wait()
}
