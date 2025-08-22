package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var alreadyInitialised = true

func Init() error {
	pkgHome, err := PKG_HOME()
	if err != nil {
		return err
	}
	if err := initDir(pkgHome); err != nil {
		return err
	}

	pkgBin, err := PKG_BIN()
	if err != nil {
		return err
	}
	if err := initDir(pkgBin); err != nil {
		return err
	}

	pkgOpt, err := PKG_OPT()
	if err != nil {
		return err
	}
	if err := initDir(pkgOpt); err != nil {
		return err
	}

	pkgZshCompletions, err := PKG_ZSH_COMPLETIONS()
	if err != nil {
		return err
	}
	if err := initDir(pkgZshCompletions); err != nil {
		return err
	}

	pkgTmp, err := PKG_TMP()
	if err != nil {
		return err
	}
	if err := initDir(pkgTmp); err != nil {
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

`, pkgHome)

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
	file, err := LOCKFILE()
	if err != nil {
		return err
	}
	if _, err := os.Stat(file); err != nil {
		alreadyInitialised = false
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("Error creating %s file: %v\n", file, err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(map[string]LockfilePackage{}); err != nil {
			return fmt.Errorf("Error writing to lockfile: %v\n", err)
		}
	}
	return nil
}
