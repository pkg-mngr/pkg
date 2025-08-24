package cmd

import "fmt"

type ErrorPackageNotInstalled struct {
	Name string
}

func (e ErrorPackageNotInstalled) Error() string {
	return fmt.Sprintf("%s is not installed", e.Name)
}

type ErrorPackageDependencyOf struct {
	Name, Dependent string
}

func (e ErrorPackageDependencyOf) Error() string {
	return fmt.Sprintf("Cannot uninstall %s as it is a dependency of %s", e.Name, e.Dependent)
}
