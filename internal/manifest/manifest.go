package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
	"github.com/noclaps/pkg/internal/platforms"
)

type PlatformScripts struct {
	Install     []string `json:"install,omitempty"`
	Latest      []string `json:"latest,omitempty"`
	Completions []string `json:"completions,omitempty"`
}

type PlatformConfig struct {
	Url     string           `json:"url"`
	Sha256  string           `json:"sha256"`
	Scripts *PlatformScripts `json:"scripts,omitempty"`
}

type Manifest struct {
	Schema       string                 `json:"$schema,omitempty"`
	ManifestUrl  string                 `json:"-"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Homepage     string                 `json:"homepage"`
	Version      string                 `json:"version"`
	Sha256       map[string]string      `json:"sha256"` // Platform-specific SHA256 checksums
	Url          map[string]string      `json:"url"`    // Platform-specific URLs
	Dependencies []string               `json:"dependencies,omitempty"`
	Caveats      string                 `json:"caveats,omitempty"`
	Scripts      map[string]interface{} `json:"scripts"` // Scripts can be global arrays or platform-specific objects
}

// GetURL returns the appropriate URL for the current platform
func (m *Manifest) GetURL() string {
	// Get platform-specific URL
	currentPlatform := config.GetCurrentPlatform()
	if url, exists := m.Url[currentPlatform.String()]; exists {
		return url
	}

	return ""
}

// GetSHA256 returns the appropriate SHA256 for the current platform
func (m *Manifest) GetSHA256() string {
	// Get platform-specific SHA256
	currentPlatform := config.GetCurrentPlatform()
	if sha256, exists := m.Sha256[currentPlatform.String()]; exists {
		return sha256
	}

	return ""
}

// GetInstallScripts returns the appropriate install scripts for the current platform
func (m *Manifest) GetInstallScripts(platform platforms.Platform) []string {
	// Handle platform-specific scripts
	if installScripts, exists := m.Scripts["install"]; exists {
		if platformScripts, ok := installScripts.(map[string]interface{}); ok {
			// Platform-specific install scripts
			if scripts, exists := platformScripts[platform.String()]; exists {
				if scriptsArr, ok := scripts.([]interface{}); ok {
					result := make([]string, len(scriptsArr))
					for i, v := range scriptsArr {
						if str, ok := v.(string); ok {
							result[i] = str
						}
					}
					return result
				}
			}
		} else if globalScripts, ok := installScripts.([]interface{}); ok {
			// Global install scripts
			result := make([]string, len(globalScripts))
			for i, v := range globalScripts {
				if str, ok := v.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}

	return []string{}
}

// GetLatestScripts returns the appropriate latest version scripts for the current platform
func (m *Manifest) GetLatestScripts(platform platforms.Platform) []string {
	// Handle platform-specific scripts
	if latestScripts, exists := m.Scripts["latest"]; exists {
		if platformScripts, ok := latestScripts.(map[string]interface{}); ok {
			// Platform-specific latest scripts
			if scripts, exists := platformScripts[platform.String()]; exists {
				if scriptsArr, ok := scripts.([]interface{}); ok {
					result := make([]string, len(scriptsArr))
					for i, v := range scriptsArr {
						if str, ok := v.(string); ok {
							result[i] = str
						}
					}
					return result
				}
			}
		} else if globalScripts, ok := latestScripts.([]interface{}); ok {
			// Global latest scripts
			result := make([]string, len(globalScripts))
			for i, v := range globalScripts {
				if str, ok := v.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}

	return []string{}
}

// GetCompletionsScripts returns the appropriate completions scripts for the current platform
func (m *Manifest) GetCompletionsScripts(platform platforms.Platform) []string {
	// Handle platform-specific scripts
	if completionsScripts, exists := m.Scripts["completions"]; exists {
		if platformScripts, ok := completionsScripts.(map[string]interface{}); ok {
			// Platform-specific completions scripts
			if scripts, exists := platformScripts[platform.String()]; exists {
				if scriptsArr, ok := scripts.([]interface{}); ok {
					result := make([]string, len(scriptsArr))
					for i, v := range scriptsArr {
						if str, ok := v.(string); ok {
							result[i] = str
						}
					}
					return result
				}
			}
		} else if globalScripts, ok := completionsScripts.([]interface{}); ok {
			// Global completions scripts
			result := make([]string, len(globalScripts))
			for i, v := range globalScripts {
				if str, ok := v.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}

	return []string{}
}

// ValidatePlatformSupport checks if the current platform is supported
func (m *Manifest) ValidatePlatformSupport() error {
	currentPlatform := config.GetCurrentPlatform()

	// Check if we have platform-specific URL/SHA256
	if _, exists := m.Url[currentPlatform.String()]; exists {
		if _, exists := m.Sha256[currentPlatform.String()]; exists {
			return nil
		}
	}

	return fmt.Errorf("package '%s' does not support the current platform '%s'", m.Name, currentPlatform)
}

func GetManifest(pkgName string) Manifest {
	if strings.HasPrefix(pkgName, "./") && strings.HasSuffix(pkgName, ".json") {
		return getManifestFromFile(pkgName)
	}

	manifestUrl, err := url.JoinPath(config.MANIFEST_HOST(), pkgName+".json")
	if err != nil {
		log.Fatalf("Error creating URL to %s manifest: %v\n", pkgName, err)
	}
	res, err := http.Get(manifestUrl)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Printf("Package %s does not exist\n", pkgName)
		os.Exit(0)
	}
	defer res.Body.Close()

	manifest := new(Manifest)
	manifest.ManifestUrl = manifestUrl
	if err := json.NewDecoder(res.Body).Decode(manifest); err != nil {
		log.Fatalf("Error decoding data from manifest: %v\n", err)
	}

	// Format platform-specific URL and SHA256
	for platform, url := range manifest.Url {
		manifest.Url[platform] = formatData(url, *manifest)
	}

	// Format platform-specific scripts
	for scriptType, scriptConfig := range manifest.Scripts {
		if platformScripts, ok := scriptConfig.(map[string]interface{}); ok {
			// Platform-specific scripts
			for platform, scripts := range platformScripts {
				if scriptsArr, ok := scripts.([]interface{}); ok {
					for i, line := range scriptsArr {
						if lineStr, ok := line.(string); ok {
							scriptsArr[i] = formatData(lineStr, *manifest)
						}
					}
					platformScripts[platform] = scriptsArr
				}
			}
			manifest.Scripts[scriptType] = platformScripts
		} else if globalScripts, ok := scriptConfig.([]interface{}); ok {
			// Global scripts
			for i, line := range globalScripts {
				if lineStr, ok := line.(string); ok {
					globalScripts[i] = formatData(lineStr, *manifest)
				}
			}
			manifest.Scripts[scriptType] = globalScripts
		}
	}

	manifest.Caveats = formatData(manifest.Caveats, *manifest)

	return *manifest
}

func getManifestFromFile(path string) Manifest {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading data from %s: %v\n", path, err)
	}

	manifest := new(Manifest)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute filepath to %s: %v\n", path, err)
	}
	manifest.ManifestUrl = absPath
	if err := json.Unmarshal(data, manifest); err != nil {
		log.Fatalf("Error unmarshalling data from %s: %v\n", path, err)
	}

	// Format platform-specific URL and SHA256
	for platform, url := range manifest.Url {
		manifest.Url[platform] = formatData(url, *manifest)
	}

	// Format platform-specific scripts
	for scriptType, scriptConfig := range manifest.Scripts {
		if platformScripts, ok := scriptConfig.(map[string]interface{}); ok {
			// Platform-specific scripts
			for platform, scripts := range platformScripts {
				if scriptsArr, ok := scripts.([]interface{}); ok {
					for i, line := range scriptsArr {
						if lineStr, ok := line.(string); ok {
							scriptsArr[i] = formatData(lineStr, *manifest)
						}
					}
					platformScripts[platform] = scriptsArr
				}
			}
			manifest.Scripts[scriptType] = platformScripts
		} else if globalScripts, ok := scriptConfig.([]interface{}); ok {
			// Global scripts
			for i, line := range globalScripts {
				if lineStr, ok := line.(string); ok {
					globalScripts[i] = formatData(lineStr, *manifest)
				}
			}
			manifest.Scripts[scriptType] = globalScripts
		}
	}

	manifest.Caveats = formatData(manifest.Caveats, *manifest)

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
