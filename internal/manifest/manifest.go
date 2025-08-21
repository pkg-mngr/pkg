package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
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

func GetManifest(pkgName string) Manifest {
	manifestJson := new(ManifestJson)
	var manifestUrl string

	if isLocalFile(pkgName) {
		path := pkgName

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Error reading data from %s: %v\n", path, err)
		}

		manifestUrl, err = filepath.Abs(path)
		if err != nil {
			log.Fatalf("Error getting absolute filepath to %s: %v\n", path, err)
		}
		if err := json.Unmarshal(data, manifestJson); err != nil {
			log.Fatalf("Error unmarshalling data from %s: %v\n", path, err)
		}
	} else {
		manifestUrl = getRemoteUrl(pkgName)

		res, err := http.Get(manifestUrl)
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Printf("Package %s does not exist\n", pkgName)
			os.Exit(0)
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(manifestJson); err != nil {
			log.Fatalf("Error decoding data from manifest: %v\n", err)
		}
	}

	platform := GetPlatform()

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
	if _, ok := manifestJson.Url[platform]; !ok {
		fmt.Printf("Package %s is not supported on this platform (%s)\n", pkgName, platform)
		os.Exit(0)
	}
	manifest.Url = formatData(manifestJson.Url[platform], *manifestJson)

	// sha256
	if _, ok := manifestJson.Sha256[platform]; !ok {
		fmt.Printf("Package %s is not supported on this platform (%s)\n", pkgName, platform)
		os.Exit(0)
	}
	manifest.Sha256 = manifestJson.Sha256[platform]

	// install script
	if _, ok := manifestJson.Scripts.Install[platform]; !ok {
		fmt.Printf("Package %s is not supported on this platform (%s)\n", pkgName, platform)
		os.Exit(0)
	}
	installScript := manifestJson.Scripts.Install[platform]
	for i, line := range installScript {
		installScript[i] = formatData(line, *manifestJson)
	}
	manifest.Scripts.Install = installScript

	// latest script
	latestScript := manifestJson.Scripts.Latest
	for i, line := range latestScript {
		latestScript[i] = formatData(line, *manifestJson)
	}
	manifest.Scripts.Latest = latestScript

	// completions script
	if completion, ok := manifestJson.Scripts.Completions[platform]; ok {
		for i, line := range completion {
			completion[i] = formatData(line, *manifestJson)
		}
		manifest.Scripts.Completions = completion
	}

	return manifest
}

func formatData(val string, manifest ManifestJson) string {
	val = strings.ReplaceAll(val, "{{ version }}", manifest.Version)
	val = strings.ReplaceAll(val, "{{ pkg.opt_dir }}", config.PKG_OPT())
	val = strings.ReplaceAll(val, "{{ pkg.bin_dir }}", config.PKG_BIN())
	val = strings.ReplaceAll(val, "{{ pkg.tmp_dir }}", config.PKG_TMP())
	val = strings.ReplaceAll(val, "{{ pkg.completions.zsh }}", config.PKG_ZSH_COMPLETIONS())
	return val
}
