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
		log.Fatalln("Error reading packages/ directory")
	}
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Println("Reading " + file.Name())
			f, err := os.OpenFile(filepath.Join("packages", file.Name()), os.O_RDWR, 0o644)
			if err != nil {
				log.Errorln("Error opening " + file.Name())
				return
			}

			log.Println("Decoding manifest from " + file.Name())
			pkgManifest := new(manifest.Manifest)
			if err := json.NewDecoder(f).Decode(pkgManifest); err != nil {
				log.Errorln("Error unmarshalling JSON from " + file.Name())
				return
			}

			latestScript := strings.Join(pkgManifest.Scripts.Latest, "\n")
			log.Println("Running `latest` script: \n" + latestScript)
			output, err := runScript(latestScript)
			if err != nil && err.Error() != "" {
				log.Println("stdout: " + output)
				log.Println("stderr: " + err.Error())
				log.Errorln("Error running latest script in " + file.Name())
				return
			}

			latestVersion := strings.TrimSpace(output)
			log.Println("Found latest version: " + latestVersion)
			if pkgManifest.Version == latestVersion {
				log.Println(pkgManifest.Name + " is already up to date")
				return
			}

			pkgManifest.Version = latestVersion
			url := strings.ReplaceAll(pkgManifest.Url, "{{ version }}", pkgManifest.Version)

			log.Println("Fetching file from " + url)
			res, err := http.Get(url)
			if err != nil {
				log.Errorln("Error fetching from " + url)
				return
			}

			log.Println("Reading file from response from " + url)
			downloadedFile, err := io.ReadAll(res.Body)
			if err != nil {
				log.Errorln("Error reading file from response body")
				return
			}
			res.Body.Close()

			log.Println("Validating checksum for ", path.Base(url))
			checksum := fmt.Sprintf("%x", sha256.Sum256(downloadedFile))
			pkgManifest.Sha256 = checksum

			log.Println("Updating " + file.Name())
			if err := f.Truncate(0); err != nil {
				log.Errorln("Error truncating " + file.Name())
				return
			}
			writer := io.NewOffsetWriter(f, 0)
			encoder := json.NewEncoder(writer)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(*pkgManifest); err != nil {
				log.Errorln("Error writing manifest to " + file.Name())
			}
			f.Close()

			log.Println("Done updating " + file.Name())
		}()
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
