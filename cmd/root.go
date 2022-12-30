package cmd

import (
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd/consts"
	"github.com/dkaslovsky/thread-safe/cmd/html"
	"github.com/dkaslovsky/thread-safe/cmd/thread"
)

// Run executes the top level command
func Run(name string, version string, cliArgs []string) error {
	if len(cliArgs) <= 1 {
		printUsage(name)
		return nil
	}

	envArgs := parseEnv()

	cmd, args := cliArgs[1], cliArgs[2:]
	args = append(args, envArgs...)

	switch cmd {
	case "thread":
		return thread.Run(args)
	case "html":
		return html.Run(args)
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

func parseEnv() []string {
	envArgs := []string{}

	path := "."
	if p, ok := os.LookupEnv(consts.EnvVarPath); ok {
		path = p
	}
	envArgs = append(envArgs, "-path", path)

	token := ""
	if t, ok := os.LookupEnv(consts.EnvVarToken); ok {
		token = t
	}
	envArgs = append(envArgs, "-token", token)

	return envArgs
}

func printUsage(name string) {
	fmt.Printf("%s saves a local copy of a Twitter thread\n", name)

	fmt.Print("\nUsage:\n")
	fmt.Printf("  %s [flags]\n", name)
	fmt.Printf("  %s [command]\n", name)

	fmt.Print("\nAvailable Commands:\n")
	fmt.Print("  thread\tsaves thread content and generates a local html file\n")
	fmt.Print("  html\t\tregenerates an html file from a previously saved thread\n")

	fmt.Print("\nFlags:\n")
	fmt.Printf("  -h, -help\thelp for %s\n", name)
	fmt.Printf("  -v, -version\tversion for %s\n", name)

	fmt.Print("\nEnvironment Variables:\n")
	fmt.Printf("  %s\ttop level path for thread files (current directory if unset)\n", consts.EnvVarPath)
	fmt.Printf("  %s\tbearer token for Twitter API\n", consts.EnvVarToken)

	fmt.Printf("\nUse \"%s [command] --help\" for more information about a command\n", name)
}

func printVersion(name string, version string) {
	fmt.Printf("%s v%s\n", name, version)
}
