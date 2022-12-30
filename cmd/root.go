package cmd

import (
	"fmt"

	"github.com/dkaslovsky/thread-safe/cmd/render"
	"github.com/dkaslovsky/thread-safe/cmd/thread"
)

// Run executes the top level command
func Run(name string, version string, cliArgs []string) error {
	if len(cliArgs) <= 1 {
		printUsage(name)
		return nil
	}

	cmd, args := cliArgs[1], cliArgs[2:]

	switch cmd {
	case "thread":
		return thread.Run(args)
	case "render":
		return render.Run(args)
	case "--help", "-help", "-h":
		printUsage(name)
		return nil
	case "--version", "-version", "-v":
		printVersion(name, version)
		return nil
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}
}

func printUsage(name string) {
	fmt.Printf("%s saves a local copy of a Twitter thread\n", name)

	fmt.Print("\nUsage:\n")
	fmt.Printf("  %s [flags]\n", name)
	fmt.Printf("  %s [command]\n", name)

	fmt.Print("\nAvailable Commands:\n")
	fmt.Print("  thread\tsaves thread content and generates a corresponding html file\n")
	fmt.Print("  render\tregenerates an html file from a previously saved thread\n")

	fmt.Print("\nFlags:\n")
	fmt.Printf("  -h, -help\thelp for %s\n", name)
	fmt.Printf("  -v, -version\tversion for %s\n", name)
}

func printVersion(name string, version string) {
	fmt.Printf("%s v%s\n", name, version)
}
