package manifest

import (
	"strings"

	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/util"
)

const MANIFEST_EXT = ".json"

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
	ManifestUrl  string              `json:"-"`
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

func GetManifest(pkgName string) (*ManifestJson, error) {
	if isLocalFile(pkgName) {
		return GetManifestFromFile(pkgName)
	}
	return GetManifestFromRemote(pkgName)
}

func (manifestJson *ManifestJson) Process() (Manifest, error) {
	manifest := Manifest{
		ManifestUrl:  manifestJson.ManifestUrl,
		Name:         manifestJson.Name,
		Description:  manifestJson.Description,
		Homepage:     manifestJson.Homepage,
		Version:      manifestJson.Version,
		Caveats:      formatData(manifestJson.Caveats, *manifestJson),
		Dependencies: manifestJson.Dependencies,
	}

	// url
	if _, ok := manifestJson.Url[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: manifestJson.Name, Platform: PLATFORM}
	}
	manifest.Url = formatData(manifestJson.Url[PLATFORM], *manifestJson)

	// sha256
	if _, ok := manifestJson.Sha256[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: manifestJson.Name, Platform: PLATFORM}
	}
	manifest.Sha256 = manifestJson.Sha256[PLATFORM]

	// install script
	if _, ok := manifestJson.Scripts.Install[PLATFORM]; !ok {
		return Manifest{}, ErrorPackageUnsupported{Name: manifestJson.Name, Platform: PLATFORM}
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

func formatData(val string, manifest ManifestJson) string {
	val = strings.ReplaceAll(val, "{{ version }}", manifest.Version)
	val = strings.ReplaceAll(val, "{{ pkg.opt_dir }}", config.PKG_OPT)
	val = strings.ReplaceAll(val, "{{ pkg.bin_dir }}", config.PKG_BIN)
	val = strings.ReplaceAll(val, "{{ pkg.tmp_dir }}", config.PKG_TMP)
	val = strings.ReplaceAll(val, "{{ pkg.completions.zsh }}", config.PKG_ZSH_COMPLETIONS)

	return val
}
