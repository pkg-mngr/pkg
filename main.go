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

	if len(args.Add.Packages) != 0 {
		cmd.Add(args.Add.Packages)
		return
	}

	if args.Update {
		cmd.Update()
		return
	}

	if len(args.Remove.Packages) != 0 {
		cmd.Remove(args.Remove.Packages)
		return
	}

	if args.Info.Package != "" {
		fmt.Println(cmd.Info(args.Info.Package))
		return
	}

	if args.List {
		pkgs := cmd.List()
		if len(pkgs) == 0 {
			fmt.Println("No packages installed!")
		}
		for _, pkg := range pkgs {
			fmt.Println(pkg)
		}
		return
	}

	if args.Init {
		config.Init()
		return
	}
}
