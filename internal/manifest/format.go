package manifest

import (
	"net/url"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

const MANIFEST_EXT = ".json"

func isLocalFile(s string) bool {
	return strings.HasPrefix(s, "./") && strings.HasSuffix(s, MANIFEST_EXT)
}

func getRemoteUrl(pkgName string) string {
	url, err := url.JoinPath(config.MANIFEST_HOST(), pkgName+MANIFEST_EXT)
	if err != nil {
		log.Fatalf("Error creating URL to %s manifest: %v\n", pkgName, err)
	}

	return url
}
