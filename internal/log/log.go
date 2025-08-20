package log

import (
	"fmt"
	"os"
)

func Fatalf(format string, a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func Errorf(format string, a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintf(os.Stderr, format, a...)
}

func Printf(format string, a ...any) {
	fmt.Fprint(os.Stderr, "\033[34mINFO:\033[0m ")
	fmt.Fprintf(os.Stderr, format, a...)
}
