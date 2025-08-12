package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

type Manifest struct {
	Schema      string `json:"$schema,omitempty"`
	ManifestUrl string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Version     string `json:"version"`
	Sha256      string `json:"sha256"`
	Url         string `json:"url"`
	Caveats     string `json:"caveats,omitempty"`
	Scripts     struct {
		Install     []string `json:"install"`
		Latest      []string `json:"latest"`
		Completions []string `json:"completions,omitempty"`
	} `json:"scripts"`
}

func GetManifest(pkgName string) Manifest {
	manifestUrl, err := url.JoinPath(config.MANIFEST_HOST(), pkgName+".json")
	if err != nil {
		log.Fatalln("Error creating URL to " + pkgName + " manifest")
	}
	res, err := http.Get(manifestUrl)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Println("Package " + pkgName + " does not exist")
		os.Exit(0)
	}
	defer res.Body.Close()

	manifest := new(Manifest)
	manifest.ManifestUrl = manifestUrl
	if err := json.NewDecoder(res.Body).Decode(manifest); err != nil {
		log.Fatalln("Error decoding data from manifest")
	}

	manifest.Url = formatData(manifest.Url, *manifest)
	manifest.Caveats = formatData(manifest.Caveats, *manifest)

	for i, line := range manifest.Scripts.Install {
		manifest.Scripts.Install[i] = formatData(line, *manifest)
	}
	for i, line := range manifest.Scripts.Latest {
		manifest.Scripts.Latest[i] = formatData(line, *manifest)
	}
	for i, line := range manifest.Scripts.Completions {
		manifest.Scripts.Completions[i] = formatData(line, *manifest)
	}

	return *manifest
}

func GetManifestFromFile(path string) Manifest {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading data from " + path)
	}

	manifest := new(Manifest)
	manifest.ManifestUrl = path
	if err := json.Unmarshal(data, manifest); err != nil {
		log.Fatalln("Error unmarshalling data from " + path)
	}

	manifest.Url = formatData(manifest.Url, *manifest)
	manifest.Caveats = formatData(manifest.Caveats, *manifest)

	for i, line := range manifest.Scripts.Install {
		manifest.Scripts.Install[i] = formatData(line, *manifest)
	}
	for i, line := range manifest.Scripts.Latest {
		manifest.Scripts.Latest[i] = formatData(line, *manifest)
	}
	for i, line := range manifest.Scripts.Completions {
		manifest.Scripts.Completions[i] = formatData(line, *manifest)
	}

	return *manifest
}

func formatData(val string, manifest Manifest) string {
	val = strings.ReplaceAll(val, "{{ version }}", manifest.Version)
	val = strings.ReplaceAll(val, "{{ pkg.opt_dir }}", config.PKG_OPT())
	val = strings.ReplaceAll(val, "{{ pkg.bin_dir }}", config.PKG_BIN())
	val = strings.ReplaceAll(val, "{{ pkg.tmp_dir }}", config.PKG_TMP())
	val = strings.ReplaceAll(val, "{{ pkg.completions.zsh }}", config.PKG_ZSH_COMPLETIONS())
	return val
}
