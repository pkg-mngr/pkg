package config

import (
	"encoding/json"
	"os"

	"github.com/noclaps/pkg/internal/log"
)

type Lockfile map[string]LockfilePackage

type LockfilePackage struct {
	Manifest     string   `json:"manifest"`
	Version      string   `json:"version"`
	Dependencies []string `json:"dependencies,omitempty"`
	Files        []string `json:"files"`
}

func ReadLockfile() Lockfile {
	data, err := os.ReadFile(LOCKFILE())
	if err != nil {
		log.Fatalf("Error reading lockfile: %v\n", err)
	}

	lf := new(Lockfile)
	if err := json.Unmarshal(data, lf); err != nil {
		log.Fatalf("Error unmarshalling lockfile: %v\n", err)
	}

	return *lf
}

func (lf Lockfile) NewEntry(name, manifest, version string, dependencies, files []string) {
	lf[name] = LockfilePackage{
		Manifest:     manifest,
		Version:      version,
		Dependencies: dependencies,
		Files:        files,
	}

	lf.Write()
}

func (lf Lockfile) Write() {
	f, err := os.Create(LOCKFILE())
	if err != nil {
		log.Fatalf("Error opening lockfile: %v\n", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(lf); err != nil {
		log.Fatalf("Error writing to lockfile: %v\n", err)
	}
}

func (lf Lockfile) Remove(name string) {
	delete(lf, name)
	lf.Write()
}
