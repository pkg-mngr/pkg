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

func formatData(val string, manifest ManifestJson) string {
	val = strings.ReplaceAll(val, "{{ version }}", manifest.Version)
	val = strings.ReplaceAll(val, "{{ pkg.opt_dir }}", config.PKG_OPT)
	val = strings.ReplaceAll(val, "{{ pkg.bin_dir }}", config.PKG_BIN)
	val = strings.ReplaceAll(val, "{{ pkg.tmp_dir }}", config.PKG_TMP)
	val = strings.ReplaceAll(val, "{{ pkg.completions.zsh }}", config.PKG_ZSH_COMPLETIONS)

	return val
}
