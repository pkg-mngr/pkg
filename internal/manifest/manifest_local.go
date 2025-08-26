package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

func FromFile(path string) (*ManifestJson, error) {
	manifestJson := new(ManifestJson)
	manifestJson.ManifestUrl = path

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading data from %s: %v", path, err)
	}
	if err := json.Unmarshal(data, manifestJson); err != nil {
		return nil, fmt.Errorf("Error unmarshalling data from %s: %v", path, err)
	}
	return manifestJson, nil
}

func IsLocalFile(p string) bool {
	return (path.IsAbs(p) || strings.HasPrefix(p, "./")) && strings.HasSuffix(p, MANIFEST_EXT)
}
