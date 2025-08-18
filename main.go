package main

import (
	"fmt"

	"github.com/noclaps/applause"
	"github.com/noclaps/pkg/internal/cmd"
	"github.com/noclaps/pkg/internal/config"
	"github.com/noclaps/pkg/internal/log"
)

type Args struct {
	Add struct {
		Packages []string `help:"Packages to install."`
	} `help:"Install packages."`
	Update bool `type:"command" help:"Update packages."`
	Remove struct {
		Packages []string `help:"Packages to remove."`
	} `help:"Remove packages."`
	Info struct {
		Package string `help:"The package to get the info for"`
	} `help:"Get the info for a package."`
	List bool `type:"command" help:"List installed packages"`
	Init bool `type:"option" help:"Initialise pkg"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		log.Fatalln(err)
	}

	if args.Info.Package != "" {
		fmt.Println(cmd.Info(args.Info.Package))
		return
	}

	if args.Init {
		config.Init()
		return
	}

	lockfile := config.ReadLockfile()
	if len(args.Add.Packages) != 0 {
		for _, pkg := range args.Add.Packages {
			cmd.Add(pkg, lockfile)
		}
		return
	}

	if args.Update {
		cmd.Update(lockfile)
		return
	}

	if len(args.Remove.Packages) != 0 {
		for _, pkg := range args.Remove.Packages {
			cmd.Remove(pkg, lockfile)
		}
		return
	}

	if args.List {
		pkgs := cmd.List(lockfile)
		if len(pkgs) == 0 {
			fmt.Println("No packages installed!")
			return
		}

		fmt.Println("\n\033[32;1m===\033[0;1m Installed\033[0m")
		for _, pkg := range pkgs {
			fmt.Println(pkg)
		}
		fmt.Println()
		return
	}
}
