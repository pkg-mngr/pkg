package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
)

func Search(query string) ([]string, error) {
	resp, err := http.Get(config.MANIFEST_HOST + "/index.json")
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s/index.json not found", config.MANIFEST_HOST)
	}
	defer resp.Body.Close()

	var index map[string]struct{ Version, Description string }
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf(
			"Error decoding %s/index.json, expected format {name: {version: string, description: string}}.",
			config.MANIFEST_HOST)
	}

	packages := make([]string, 0, len(index))
	for name, data := range index {
		line := fmt.Sprintf("\033[1m%s:\033[0m %s - %s", name, data.Version, data.Description)
		searchLine := strings.ToLower(fmt.Sprintf("%s %s", name, data.Description))
		if strings.Contains(searchLine, strings.ToLower(query)) {
			packages = append(packages, line)
		}
	}

	slices.Sort(packages)

	return packages, nil
}
