package manifest

import "fmt"

type ErrorPackageNotFound struct {
	Name string
}

func (e ErrorPackageNotFound) Error() string {
	return fmt.Sprintf("Package %s does not exist", e.Name)
}

type ErrorPackageUnsupported struct {
	Name     string
	Platform Platform
}

func (e ErrorPackageUnsupported) Error() string {
	return fmt.Sprintf("Package %s is not supported on this platform (%s)", e.Name, e.Platform)
}
