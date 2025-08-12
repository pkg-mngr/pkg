package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
)

func main() {
	files := os.Args[1:]
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Println("Reading " + file)
			f, err := os.OpenFile(file, os.O_RDWR, 0o644)
			if err != nil {
				log.Fatalln("Error opening " + file)
			}

			log.Println("Decoding manifest from " + file)
			pkgManifest := new(manifest.Manifest)
			if err := json.NewDecoder(f).Decode(pkgManifest); err != nil {
				log.Fatalln("Error unmarshalling JSON from " + file)
			}

			url := strings.ReplaceAll(pkgManifest.Url, "{{ version }}", pkgManifest.Version)

			log.Println("Fetching file from " + url)
			res, err := http.Get(url)
			if err != nil {
				log.Fatalln("Error fetching from " + url)
			}

			log.Println("Reading file from response from " + url)
			downloadedFile, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("Error reading file from response body")
			}
			res.Body.Close()

			log.Println("Validating checksum for " + path.Base(url))
			checksum := fmt.Sprintf("%x", sha256.Sum256(downloadedFile))
			if checksum != pkgManifest.Sha256 {
				log.Errorln("Calculated checksum did not match checksum in " + file)
				log.Errorln("The calculated checksum was: " + checksum)
				log.Fatalln("The checksum in " + file + " was " + pkgManifest.Sha256)
			}

			log.Println("Everything looks good!")
		}()
	}

	wg.Wait()
}
