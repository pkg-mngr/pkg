package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/manifest"
)

func main() {
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
			pkgManifest := new(manifest.Manifest)
			if err := json.NewDecoder(f).Decode(pkgManifest); err != nil {
				log.Errorf("Error unmarshalling JSON from %s: %v\n", file.Name(), err)
				return
			}

			latestScript := strings.Join(pkgManifest.Scripts.Latest, "\n")
			log.Println("Running `latest` script: \n" + latestScript)
			output, err := runScript(latestScript)
			if err != nil && err.Error() != "" {
				log.Printf("stdout: %s\n", output)
				log.Printf("stderr: %v\n", err)
				log.Errorf("Error running latest script in %s\n", file.Name())
				return
			}

			latestVersion := strings.TrimSpace(output)
			log.Printf("Found latest version: %s\n", latestVersion)
			if pkgManifest.Version == latestVersion {
				log.Printf("%s is already up to date\n", pkgManifest.Name)
				return
			}

			pkgManifest.Version = latestVersion
			url := strings.ReplaceAll(pkgManifest.Url, "{{ version }}", pkgManifest.Version)

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
			pkgManifest.Sha256 = checksum

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

func runScript(script string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("Error getting stdout pipe")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("Error getting stderr pipe")
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error while starting command")
	}

	stdoutData, err := io.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stdout")
	}
	stderrData, err := io.ReadAll(stderr)
	if err != nil {
		return "", fmt.Errorf("Error getting data from stderr")
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Error waiting for command")
		return "", err
	}

	return string(stdoutData), fmt.Errorf("%s", stderrData)
}
