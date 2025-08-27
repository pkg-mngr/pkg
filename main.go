package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/noclaps/applause"
	"github.com/pkg-mngr/pkg/internal/cmd"
	"github.com/pkg-mngr/pkg/internal/config"
	"github.com/pkg-mngr/pkg/internal/log"
	"github.com/pkg-mngr/pkg/internal/manifest"
)

type Args struct {
	Add struct {
		Packages []string `help:"Packages to install."`
		Yes      bool     `type:"option" short:"y" help:"Skip confirmation to run scripts."`
	} `help:"Install packages."`
	Update *struct {
		Packages []string `help:"Packages to update." completion:"$(jq -r 'keys[]' $PKG_HOME/pkg.lock | tr '\n' ' ')"`
		Yes      bool     `type:"option" short:"y" help:"Skip confirmation to run scripts."`
	} `help:"Update packages."`
	Remove struct {
		Packages []string `help:"Packages to remove." completion:"$(jq -r 'keys[]' $PKG_HOME/pkg.lock | tr '\n' ' ')"`
	} `help:"Remove packages."`
	Info struct {
		Package string `help:"The package to get the info for"`
	} `help:"Get the info for a package."`
	Search struct {
		Name string `help:"The search query"`
	} `help:"Search for packages"`
	List bool `type:"command" help:"List installed packages"`
	Init bool `type:"option" help:"Initialise pkg"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		log.Fatalf("%v\n", err)
	}

	if args.Info.Package != "" {
		info, err := cmd.Info(args.Info.Package)
		if err != nil {
			errPnf := manifest.ErrorPackageNotFound{}
			errPu := manifest.ErrorPackageUnsupported{}
			switch {
			case errors.As(err, &errPnf):
				log.Errorf("%v\n", errPnf)
			case errors.As(err, &errPu):
				log.Errorf("%v\n", errPu)
			default:
				log.Fatalf("%v\n", err)
			}
			return
		}
		fmt.Println(info)
		return
	}

	if args.Init {
		if err := config.Init(); err != nil {
			log.Fatalf("%v\n", err)
		}
		return
	}

	if args.Search.Name != "" {
		results, err := cmd.Search(args.Search.Name)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		if len(results) == 0 {
			fmt.Println("\n\033[31;1m ===\033[0;1m No results found\033[0m")
			return
		}
		fmt.Println("\n\033[32;1m===\033[0;1m Search results\033[0m")
		for _, result := range results {
			fmt.Println(result)
		}
		fmt.Println()
		return
	}

	lockfile, err := config.ReadLockfile()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer lockfile.Write()

	if len(args.Add.Packages) != 0 {
		for _, pkg := range args.Add.Packages {
			if err := cmd.Add(pkg, args.Add.Yes, lockfile); err != nil {
				errPnf := manifest.ErrorPackageNotFound{}
				errPu := manifest.ErrorPackageUnsupported{}
				switch {
				case errors.As(err, &errPnf):
					log.Errorf("%v\n", errPnf)
				case errors.As(err, &errPu):
					log.Errorf("%v\n", errPu)
				default:
					log.Fatalf("%v\n", err)
				}
				continue
			}
		}
		return
	}

	if args.Update != nil {
		pkgs := slices.Collect(maps.Keys(lockfile))
		if len(args.Update.Packages) > 0 {
			pkgs = args.Update.Packages
		}
		if err := cmd.Update(pkgs, args.Update.Yes, lockfile); err != nil {
			errPnf := manifest.ErrorPackageNotFound{}
			errPu := manifest.ErrorPackageUnsupported{}
			switch {
			case errors.As(err, &errPnf):
				log.Errorf("%v\n", errPnf)
			case errors.As(err, &errPu):
				log.Errorf("%v\n", errPu)
			default:
				log.Fatalf("%v\n", err)
			}
			return
		}
		return
	}

	if len(args.Remove.Packages) != 0 {
		for _, pkg := range args.Remove.Packages {
			if err := cmd.Remove(pkg, lockfile, false); err != nil {
				log.Fatalf("%v\n", err)
			}
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
