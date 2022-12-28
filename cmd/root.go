package cmd

import (
	"fmt"

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
	case "all":
		// run all three commands
		return fmt.Errorf("not yet implemented")
	case "thread":
		return thread.Run(args)
	case "-help", "-h":
		// TODO: printUsage(name)
		return nil
	case "-version", "-v":
		// TODO: printVersion(name, version)
		return nil
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}
}
