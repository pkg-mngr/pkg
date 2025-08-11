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
		output[i] = "\033[1m" + key + ":\033[0m " + lockfile.Packages[key].Version
	}

	return output
}
