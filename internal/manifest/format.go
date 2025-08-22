package manifest

import (
	"fmt"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
)

const MANIFEST_EXT = ".json"

func isLocalFile(s string) bool {
	return strings.HasPrefix(s, "./") && strings.HasSuffix(s, MANIFEST_EXT)
}

func getRemoteUrl(pkgName string) string {
	return fmt.Sprintf("%s/%s%s", config.MANIFEST_HOST, pkgName, MANIFEST_EXT)
}
