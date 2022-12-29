package cmd

import (
	"fmt"

	"github.com/dkaslovsky/thread-safe/cmd/render"
	"github.com/dkaslovsky/thread-safe/cmd/thread"
)

// Run executes the top level command
func Run(name string, version string, cliArgs []string) error {
	if len(cliArgs) <= 1 {
		// TODO: printUsage(name)
		return nil
	}

	cmd, args := cliArgs[1], cliArgs[2:]

	switch cmd {
	case "thread":
		return thread.Run(args)
	case "render":
		return render.Run(args)
	case "--help", "-help", "-h":
		// TODO: printUsage(name)
		return nil
	case "--version", "-version", "-v":
		// TODO: printVersion(name, version)
		return nil
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}
}
