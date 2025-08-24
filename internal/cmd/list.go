package cmd

import (
	"fmt"
	"maps"
	"slices"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/util"
)

func List(lockfile config.Lockfile) []string {
	keys := slices.Collect(maps.Keys(lockfile))

	output := util.Map(keys, func(key string, i int) string {
		return fmt.Sprintf("\033[1m%s:\033[0m %s", key, lockfile[key].Version)
	})

	slices.Sort(output)

	return output
}
