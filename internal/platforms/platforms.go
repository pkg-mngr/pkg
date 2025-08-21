package platforms

import (
	"fmt"
	"runtime"
)

type Platform string

func GetPlatform() Platform {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if arch == "amd64" {
		arch = "x64"
	}
	if os == "darwin" {
		os = "macos"
	}

	return Platform(fmt.Sprintf("%s-%s", os, arch))
}
