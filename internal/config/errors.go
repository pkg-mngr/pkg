package config

type ErrorAlreadyInitialised struct{}

func (ErrorAlreadyInitialised) Error() string {
	return "pkg is already initialised!"
}
