package manifest

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
)

const MANIFEST_EXT = ".json"

func isLocalFile(s string) bool {
	return strings.HasPrefix(s, "./") && strings.HasSuffix(s, MANIFEST_EXT)
}

func getRemoteUrl(pkgName string) (string, error) {
	url, err := url.JoinPath(config.MANIFEST_HOST(), pkgName+MANIFEST_EXT)
	if err != nil {
		return "", fmt.Errorf("Error creating URL to %s manifest: %v", pkgName, err)
	}

	return url, nil
}
