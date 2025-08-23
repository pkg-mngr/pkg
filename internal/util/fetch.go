package util

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/melbahja/got"
)

func Fetch(url, dest, name string) error {
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
	if err := g.Download(url, dest); err != nil {
		return fmt.Errorf("%s: Error while downloading %s: %v", name, dest, err)
	}
	fmt.Println("\033[2KDownloaded 100%")
	return nil
}

// Returns an error if the sha256 checksum of the file does not match the
// provided shasum, otherwise returns nil
func VerifyChecksum(filename, shasum, name string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s: Error reading data from response", name)
	}

	fmt.Print("Verifying checksum...")
	checksum := fmt.Sprintf("%x", sha256.Sum256(data))
	if checksum != shasum {
		return fmt.Errorf("\n%s: Checksum of data did not match in package manifest", name)
	}
	fmt.Println(" Looks good!")

	return nil
}
