package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GetManifestFromFile(path string) (*ManifestJson, error) {
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

func isLocalFile(s string) bool {
	return (strings.HasPrefix(s, "./") || strings.HasPrefix(s, "/")) && strings.HasSuffix(s, MANIFEST_EXT)
}
