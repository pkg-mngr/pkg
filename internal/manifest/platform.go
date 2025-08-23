package manifest

import (
	"runtime"
)

type Platform string

var PLATFORM = GetPlatform()

func GetPlatform() Platform {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if os == "darwin" {
		os = "macos"
	}
	if arch == "amd64" {
		arch = "x64"
	}
	return Platform(os + "-" + arch)
}
