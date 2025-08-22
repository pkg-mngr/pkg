package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Lockfile map[string]LockfilePackage

type LockfilePackage struct {
	Manifest     string   `json:"manifest"`
	Version      string   `json:"version"`
	Dependencies []string `json:"dependencies,omitempty"`
	Files        []string `json:"files"`
}

func ReadLockfile() (Lockfile, error) {
	data, err := os.ReadFile(LOCKFILE)
	if err != nil {
		return nil, fmt.Errorf("Error reading lockfile: %v\n", err)
	}

	lf := new(Lockfile)
	if err := json.Unmarshal(data, lf); err != nil {
		return nil, fmt.Errorf("Error unmarshalling lockfile: %v\n", err)
	}

	return *lf, nil
}

func (lf Lockfile) NewEntry(name, manifest, version string, dependencies, files []string) error {
	lf[name] = LockfilePackage{
		Manifest:     manifest,
		Version:      version,
		Dependencies: dependencies,
		Files:        files,
	}

	return lf.Write()
}

func (lf Lockfile) Write() error {
	f, err := os.Create(LOCKFILE)
	if err != nil {
		return fmt.Errorf("Error opening lockfile: %v\n", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(lf); err != nil {
		return fmt.Errorf("Error writing to lockfile: %v\n", err)
	}

	return nil
}

func (lf Lockfile) Remove(name string) error {
	delete(lf, name)
	return lf.Write()
}
