package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
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
		var err error
		manifestUrl, err = getRemoteUrl(pkgName)
		if err != nil {
			return Manifest{}, err
		}

		res, err := http.Get(manifestUrl)
		if err != nil || res.StatusCode != http.StatusOK {
			return Manifest{}, ErrorPackageNotFound{Name: pkgName}
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(manifestJson); err != nil {
			return Manifest{}, fmt.Errorf("Error decoding data from manifest: %v", err)
		}
	}

	platform := GetPlatform()

	manifest := Manifest{
		ManifestUrl:  manifestUrl,
		Name:         manifestJson.Name,
		Description:  manifestJson.Description,
		Homepage:     manifestJson.Homepage,
		Version:      manifestJson.Version,
		Dependencies: manifestJson.Dependencies,
	}

	formatted, err := formatData(manifestJson.Caveats, *manifestJson)
	if err != nil {
		return Manifest{}, err
	}
	manifest.Caveats = formatted

	// url
	if _, ok := manifestJson.Url[platform]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: platform}
	}
	formatted, err = formatData(manifestJson.Url[platform], *manifestJson)
	if err != nil {
		return Manifest{}, err
	}
	manifest.Url = formatted

	// sha256
	if _, ok := manifestJson.Sha256[platform]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: platform}
	}
	manifest.Sha256 = manifestJson.Sha256[platform]

	// install script
	if _, ok := manifestJson.Scripts.Install[platform]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: pkgName, Platform: platform}
	}
	installScript := manifestJson.Scripts.Install[platform]
	for i, line := range installScript {
		formatted, err := formatData(line, *manifestJson)
		if err != nil {
			return Manifest{}, err
		}
		installScript[i] = formatted
	}
	manifest.Scripts.Install = installScript

	// latest script
	latestScript := manifestJson.Scripts.Latest
	for i, line := range latestScript {
		formatted, err := formatData(line, *manifestJson)
		if err != nil {
			return Manifest{}, err
		}
		latestScript[i] = formatted
	}
	manifest.Scripts.Latest = latestScript

	// completions script
	if completion, ok := manifestJson.Scripts.Completions[platform]; ok {
		for i, line := range completion {
			formatted, err := formatData(line, *manifestJson)
			if err != nil {
				return Manifest{}, err
			}
			completion[i] = formatted
		}
		manifest.Scripts.Completions = completion
	}

	return manifest, nil
}

func formatData(val string, manifest ManifestJson) (string, error) {
	val = strings.ReplaceAll(val, "{{ version }}", manifest.Version)

	pkgOpt, err := config.PKG_OPT()
	if err != nil {
		return "", err
	}
	val = strings.ReplaceAll(val, "{{ pkg.opt_dir }}", pkgOpt)

	pkgBin, err := config.PKG_BIN()
	if err != nil {
		return "", err
	}
	val = strings.ReplaceAll(val, "{{ pkg.bin_dir }}", pkgBin)

	pkgTmp, err := config.PKG_TMP()
	if err != nil {
		return "", err
	}
	val = strings.ReplaceAll(val, "{{ pkg.tmp_dir }}", pkgTmp)

	pkgZshCompletions, err := config.PKG_ZSH_COMPLETIONS()
	if err != nil {
		return "", err
	}
	val = strings.ReplaceAll(val, "{{ pkg.completions.zsh }}", pkgZshCompletions)

	return val, nil
}
