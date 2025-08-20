package cmd

import (
	"fmt"
	"maps"
	"slices"

	"github.com/noclaps/pkg/internal/config"
)

func List(lockfile config.Lockfile) []string {
	keys := slices.Collect(maps.Keys(lockfile))
	output := make([]string, len(lockfile))

	for i, key := range keys {
		output[i] = fmt.Sprintf("\033[1m%s:\033[0m %s", key, lockfile[key].Version)
	}

	slices.Sort(output)

	return output
}
