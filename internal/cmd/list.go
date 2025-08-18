package cmd

import (
	"maps"
	"slices"

	"github.com/noclaps/pkg/internal/config"
)

func List() []string {
	lockfile := config.ReadLockfile()
	keys := slices.Collect(maps.Keys(lockfile))
	output := make([]string, len(lockfile))

	for i, key := range keys {
		output[i] = "\033[1m" + key + ":\033[0m " + lockfile[key].Version
	}

	slices.Sort(output)

	return output
}
