package cmd

import (
	"maps"
	"slices"

	"github.com/noclaps/pkg/internal/config"
)

func List() []string {
	lockfile := config.ReadLockfile()
	keys := slices.Collect(maps.Keys(lockfile.Packages))
	output := make([]string, len(lockfile.Packages))

	for i, key := range keys {
		output[i] = key + " " + lockfile.Packages[key].Version
	}

	return output
}
