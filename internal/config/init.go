package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/noclaps/pkg/internal/log"
)

var alreadyInitialised = true

func Init() {
	initDir(PKG_HOME())
	initDir(PKG_BIN())
	initDir(PKG_OPT())
	initDir(PKG_ZSH_COMPLETIONS())
	initDir(PKG_TMP())

	initLockfile()

	if alreadyInitialised {
		fmt.Println("pkg is already initialised!")
		return
	}

	fmt.Printf(`pkg has been installed to %s!
Add the following to your ~/.zshrc to complete installation:

export PKG_HOME="%s"
export PATH="$PKG_HOME:$PATH"
export FPATH="$PKG_HOME/share/zsh/site-functions:$FPATH"

`, PKG_HOME(), PKG_HOME())
}

func initDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		alreadyInitialised = false
		if err := os.MkdirAll(dir, 0o750); err != nil {
			log.Fatalln("Error creating " + dir + " directory")
		}
	}
}

func initLockfile() {
	file := LOCKFILE()
	if _, err := os.Stat(LOCKFILE()); err != nil {
		alreadyInitialised = false
		f, err := os.Create(LOCKFILE())
		if err != nil {
			log.Fatalln("Error creating " + file + " file")
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(Lockfile{
			Version:  LOCKFILE_VERSION,
			Packages: map[string]LockfilePackage{},
		}); err != nil {
			log.Fatalln("Error writing to lockfile")
		}
	}
}
