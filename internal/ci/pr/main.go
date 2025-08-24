package main

import (
	"errors"
	"os"

	"github.com/pkg-mngr/pkg/internal/cmd"
	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
	"github.com/pkg-mngr/pkg/internal/manifest"
)

func main() {
	files := os.Args[1:]
	if err := config.Init(); err != nil {
		log.Fatalf("%v\n", err)
	}

	os.Setenv("PATH", config.PKG_BIN+":"+os.Getenv("PATH"))
	lockfile, err := config.ReadLockfile()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	for _, file := range files {
		log.Printf("Checking if installation works...\n")
		if err := cmd.Add("./"+file, true, lockfile); err != nil {
			errPu := manifest.ErrorPackageUnsupported{}
			switch {
			case errors.As(err, &errPu):
				log.Errorf("%v\n", errPu)
			default:
				log.Fatalf("%v\n", err)
			}
			continue
		}
		log.Printf("Everything looks good!\n")
	}
}
