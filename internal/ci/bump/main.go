package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
			f, err := os.OpenFile(filepath.Join("packages", file.Name()), os.O_RDWR, 0o644)
			if err != nil {
				log.Errorf("Error opening %s: %v\n", file.Name(), err)
				return
			}

			log.Printf("Decoding manifest from %s\n", file.Name())
			pkgManifest := new(manifest.ManifestJson)
			if err := json.NewDecoder(f).Decode(pkgManifest); err != nil {
				log.Errorf("Error unmarshalling JSON from %s: %v\n", file.Name(), err)
				return
			}

			latestScript := strings.Join(pkgManifest.Scripts.Latest, "\n")
			log.Printf("Running `latest` script: \n%s\n", latestScript)
			output, err := util.RunScript(latestScript, true)
			if err != nil && err.Error() != "" {
				log.Errorf(
					"stdout: %s\nstderr: %v\n, Error running latest script in %s\n",
					output, err, file.Name())
				return
			}

			latestVersion := strings.TrimSpace(output)
			log.Printf("Found latest version: %s\n", latestVersion)
			if pkgManifest.Version == latestVersion {
				log.Printf("%s is already up to date\n", pkgManifest.Name)
				return
			}

			pkgManifest.Version = latestVersion

			for platform := range pkgManifest.Url {
				url := strings.ReplaceAll(pkgManifest.Url[platform], "{{ version }}", pkgManifest.Version)

				log.Printf("Fetching file from %s\n", url)
				res, err := http.Get(url)
				if err != nil {
					log.Errorf("Error fetching from %s: %v\n", url, err)
					return
				}

				log.Printf("Reading file from response from %s\n", url)
				downloadedFile, err := io.ReadAll(res.Body)
				if err != nil {
					log.Errorf("Error reading file from response body: %v\n", err)
					return
				}
				res.Body.Close()

				log.Printf("Validating checksum for %s\n", path.Base(url))
				checksum := fmt.Sprintf("%x", sha256.Sum256(downloadedFile))
				pkgManifest.Sha256[platform] = checksum
			}

			log.Printf("Updating %s\n", file.Name())
			if err := f.Truncate(0); err != nil {
				log.Errorf("Error truncating %s\n", file.Name())
				return
			}
			writer := io.NewOffsetWriter(f, 0)
			encoder := json.NewEncoder(writer)
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
