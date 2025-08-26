package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg-mngr/pkg/internal/config"
)

func FromRemote(url string) (*ManifestJson, error) {
	manifestJson := new(ManifestJson)
	manifestJson.ManifestUrl = url

	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, ErrorPackageNotFound{Url: url}
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(manifestJson); err != nil {
		return nil, fmt.Errorf("Error decoding data from manifest: %v", err)
	}

	return manifestJson, nil
}

func getRemoteUrl(pkgName string) string {
	return fmt.Sprintf("%s/%s%s", config.MANIFEST_HOST, pkgName, MANIFEST_EXT)
}
