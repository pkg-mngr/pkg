package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var alreadyInitialised = true

func Init() error {
	if err := initDir(PKG_HOME); err != nil {
		return err
	}
	if err := initDir(PKG_BIN); err != nil {
		return err
	}
	if err := initDir(PKG_OPT); err != nil {
		return err
	}
	if err := initDir(PKG_ZSH_COMPLETIONS); err != nil {
		return err
	}
	if err := initDir(PKG_TMP); err != nil {
		return err
	}

	if err := initLockfile(); err != nil {
		return err
	}

	if alreadyInitialised {
		return ErrorAlreadyInitialised{}
	}

	fmt.Printf(`pkg has been installed to %[1]s!
Add the following to your ~/.zshrc to complete installation:

export PKG_HOME="%[1]s"
export PATH="$PKG_HOME/bin:$PATH"
export FPATH="$PKG_HOME/share/zsh/site-functions:$FPATH"

`, PKG_HOME)

	return nil
}

func initDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		alreadyInitialised = false
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("Error creating %s directory: %v\n", dir, err)
		}
	}
	return nil
}

func initLockfile() error {
	if _, err := os.Stat(LOCKFILE); err != nil {
		alreadyInitialised = false
		f, err := os.Create(LOCKFILE)
		if err != nil {
			return fmt.Errorf("Error creating %s file: %v\n", LOCKFILE, err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(map[string]LockfilePackage{}); err != nil {
			return fmt.Errorf("Error writing to lockfile: %v\n", err)
		}
	}
	return nil
}
