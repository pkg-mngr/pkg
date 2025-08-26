package manifest

import "fmt"

type ErrorPackageNotFound struct {
	Url string
}

func (e ErrorPackageNotFound) Error() string {
	return fmt.Sprintf("Package not found at url: %s", e.Url)
}

type ErrorPackageUnsupported struct {
	Name     string
	Platform Platform
}

func (e ErrorPackageUnsupported) Error() string {
	return fmt.Sprintf("Package %s is not supported on this platform (%s)", e.Name, e.Platform)
}
