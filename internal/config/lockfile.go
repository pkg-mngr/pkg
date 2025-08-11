package config

import (
	"encoding/json"
	"os"

	"github.com/noclaps/pkg/internal/log"
)

const LOCKFILE_VERSION = 1

type Lockfile struct {
	Version  int                        `json:"version"`
	Packages map[string]LockfilePackage `json:"packages"`
}

type LockfilePackage struct {
	Manifest string   `json:"manifest"`
	Version  string   `json:"version"`
	Files    []string `json:"files"`
}

func ReadLockfile() *Lockfile {
	data, err := os.ReadFile(LOCKFILE())
	if err != nil {
		log.Fatalln("Error reading lockfile")
	}

	lf := new(Lockfile)
	if err := json.Unmarshal(data, lf); err != nil {
		log.Fatalln("Error unmarshalling lockfile")
	}

	return lf
}

func (lf *Lockfile) WriteToLockfile(name, manifest, version string, files []string) {
	lf.Packages[name] = LockfilePackage{
		Manifest: manifest,
		Version:  version,
		Files:    files,
	}

	f, err := os.Create(LOCKFILE())
	if err != nil {
		log.Fatalln("Error opening lockfile")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(*lf); err != nil {
		log.Fatalln("Error writing to lockfile")
	}
}

func (lf *Lockfile) RemoveFromLockfile(name string) {
	delete(lf.Packages, name)

	f, err := os.Create(LOCKFILE())
	if err != nil {
		log.Fatalln("Error opening lockfile")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(*lf); err != nil {
		log.Fatalln("Error writing to lockfile")
	}
}
