package main

import (
	"os"

	"github.com/pkg-mngr/pkg/internal/cmd"
	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
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
