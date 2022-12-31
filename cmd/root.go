package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd/env"
	"github.com/dkaslovsky/thread-safe/cmd/regen"
	"github.com/dkaslovsky/thread-safe/cmd/save"
)

// Run executes the top level command
func Run(name string, version string) error {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, fmt.Sprintf("version for %s", name))
	flag.BoolVar(&versionFlag, "version", false, fmt.Sprintf("version for %s", name))

	setUsage(name)
	flag.Parse()

	if versionFlag {
		printVersion(name, version)
		return nil
	}

	if flag.NArg() == 0 {
		flag.Usage()
		return nil
	}

	subCmd, args := flag.Arg(0), os.Args[2:]

	switch subCmd {
	case "save":
		return save.Run(name, args)
	case "regen":
		return regen.Run(name, args)
	case "version":
		printVersion(name, version)
	case "help":
		flag.Usage()
	default:
		return fmt.Errorf("unknown command \"%s\"", subCmd)
	}

	return nil
}

func setUsage(name string) {
	flag.Usage = func() {
		fmt.Printf(usage, name, name, name, name, name, env.Usage(), name)
	}
}

func printVersion(name string, version string) {
	fmt.Printf("%s v%s\n", name, version)
}

const usage = `'%s' saves a local copy of a Twitter thread

Usage:
  %s [flags]
  %s [command]

Available Commands:
  save    saves thread content and generates a local html file
  html    regenerates an html file from a previously saved thread

Flags:
  -h, --help	 help for %s
  -v, --version	 version for %s

%s

Use "%s [command] --help" for more information about a command`
