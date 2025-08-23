package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg-mngr/pkg/internal/log"
)

var alreadyInitialised = true

func Init() error {
	err := initDirs(PKG_HOME, PKG_BIN, PKG_OPT, PKG_TMP, PKG_ZSH_COMPLETIONS)
	if err != nil {
		return err
	}

	if err := initLockfile(); err != nil {
		return err
	}

	if alreadyInitialised {
		log.Printf("pkg is already initialised")
		return nil
	}

	fmt.Printf(`pkg has been installed to %[1]s!
Add the following to your ~/.zshrc to complete installation:

export PKG_HOME="%[1]s"
export PATH="$PKG_HOME/bin:$PATH"
export FPATH="$PKG_HOME/share/zsh/site-functions:$FPATH"

`, PKG_HOME)

	return nil
}

func initDirs(dirs ...string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			continue
		}
		alreadyInitialised = false
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("Error creating %s directory: %v\n", dir, err)
		}

	}
	return nil
}

func initLockfile() error {
	if _, err := os.Stat(LOCKFILE); err == nil {
		return nil
	}
	alreadyInitialised = false
	f, err := os.Create(LOCKFILE)
	if err != nil {
		return fmt.Errorf("Error creating %s file: %v\n", LOCKFILE, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(map[string]LockfilePackage{}); err != nil {
		return fmt.Errorf("Error writing to lockfile: %v\n", err)
	}

	return nil
}
