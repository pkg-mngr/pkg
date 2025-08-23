package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg-mngr/pkg/internal/util"
)

type Manifest struct {
	ManifestUrl  string
	Name         string
	Description  string
	Homepage     string
	Version      string
	Sha256       string
	Url          string
	Dependencies []string
	Caveats      string
	Scripts      struct {
		Install     []string
		Latest      []string
		Completions []string
	}
}

type ManifestJson struct {
	Schema       string              `json:"$schema,omitempty"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Homepage     string              `json:"homepage"`
	Version      string              `json:"version"`
	Sha256       map[Platform]string `json:"sha256"`
	Url          map[Platform]string `json:"url"`
	Dependencies []string            `json:"dependencies,omitempty"`
	Caveats      string              `json:"caveats,omitempty"`
	Scripts      struct {
		Install     map[Platform][]string `json:"install"`
		Latest      []string              `json:"latest"`
		Completions map[Platform][]string `json:"completions,omitempty"`
	} `json:"scripts"`
}

func GetManifest(pkgName string) (Manifest, error) {
	manifestJson := new(ManifestJson)
	var manifestUrl string

	if isLocalFile(pkgName) {
		path := pkgName

		data, err := os.ReadFile(path)
		if err != nil {
			return Manifest{}, fmt.Errorf("Error reading data from %s: %v", path, err)
		}

		manifestUrl, err = filepath.Abs(path)
		if err != nil {
			return Manifest{}, fmt.Errorf("Error getting absolute filepath to %s: %v", path, err)
		}
		if err := json.Unmarshal(data, manifestJson); err != nil {
			return Manifest{}, fmt.Errorf("Error unmarshalling data from %s: %v", path, err)
		}
	} else {
		manifestUrl = getRemoteUrl(pkgName)

		res, err := http.Get(manifestUrl)
		if err != nil || res.StatusCode != http.StatusOK {
			return Manifest{}, ErrorPackageNotFound{Name: pkgName}
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(manifestJson); err != nil {
			return Manifest{}, fmt.Errorf("Error decoding data from manifest: %v", err)
		}
	}

	manifest := Manifest{
		ManifestUrl:  manifestUrl,
		Name:         manifestJson.Name,
		Description:  manifestJson.Description,
		Homepage:     manifestJson.Homepage,
		Version:      manifestJson.Version,
		Caveats:      formatData(manifestJson.Caveats, *manifestJson),
		Dependencies: manifestJson.Dependencies,
	}

	// url
	if _, ok := manifestJson.Url[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: PLATFORM}
	}
	manifest.Url = formatData(manifestJson.Url[PLATFORM], *manifestJson)

	// sha256
	if _, ok := manifestJson.Sha256[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: PLATFORM}
	}
	manifest.Sha256 = manifestJson.Sha256[PLATFORM]

	// install script
	if _, ok := manifestJson.Scripts.Install[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: PLATFORM}
	}
	installScript := manifestJson.Scripts.Install[PLATFORM]
	manifest.Scripts.Install = util.Map(installScript, func(line string, i int) string {
		return formatData(line, *manifestJson)
	})

	// latest script
	latestScript := manifestJson.Scripts.Latest
	manifest.Scripts.Latest = util.Map(latestScript, func(line string, i int) string {
		return formatData(line, *manifestJson)
	})

	// completions script
	if completion, ok := manifestJson.Scripts.Completions[PLATFORM]; ok {
		manifest.Scripts.Completions = util.Map(completion, func(line string, i int) string {
			return formatData(line, *manifestJson)
		})
	}

	return manifest, nil
}
