package main

import (
	"os"

	"github.com/noclaps/pkg/internal/cmd"
	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

func main() {
	files := os.Args[1:]
	config.Init()
	os.Setenv("PATH", config.PKG_BIN()+":"+os.Getenv("PATH"))
	lockfile := config.ReadLockfile()

	for _, file := range files {
		log.Printf("Checking if installation works...\n")
		cmd.Add("./"+file, true, lockfile)
		log.Printf("Everything looks good!\n")
	}
}
