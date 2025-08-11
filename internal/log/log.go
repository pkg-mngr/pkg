package log

import (
	"fmt"
	"os"
)

func Fatalln(a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func Errorln(a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintln(os.Stderr, a...)
}

func Println(a ...any) {
	fmt.Fprint(os.Stderr, "\033[34mINFO:\033[0m ")
	fmt.Fprintln(os.Stderr, a...)
}
