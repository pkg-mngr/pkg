package main

import (
	"os"

	"github.com/noclaps/pkg/internal/cmd"
	"github.com/noclaps/pkg/internal/log"
)

func main() {
	files := os.Args[1:]

	for _, file := range files {
		log.Println("Checking if installation works...")
		cmd.Add("./" + file)
		log.Println("Everything looks good!")
	}
}
