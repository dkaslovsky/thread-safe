package html

import (
	"errors"
	"flag"
	"fmt"

	"github.com/dkaslovsky/thread-safe/cmd/errs"
	"github.com/dkaslovsky/thread-safe/pkg/thread"
)

// Run executes the package's (sub)command
func Run(args []string) error {
	cmd := flag.NewFlagSet("html", flag.ExitOnError)
	opts := &cmdOpts{}
	attachOpts(cmd, opts)
	setUsage(cmd)

	err := parseArgs(cmd, opts, args)
	if err != nil {
		if errors.Is(err, errs.ErrNoArgs) {
			cmd.Usage()
			return nil
		}
		return err
	}

	return run(opts)
}

func run(opts *cmdOpts) error {
	threadDir := thread.Dir(opts.path, opts.name)

	th, err := thread.NewThreadFromFile(threadDir)
	if err != nil {
		return fmt.Errorf("failed to load thread from file: %w", err)
	}

	tErr := th.ToHTML(threadDir)
	if tErr != nil {
		return fmt.Errorf("failed to write thread HTML file: %w", err)
	}

	return nil
}

type cmdOpts struct {
	name string
	path string // TODO: env
}

func attachOpts(cmd *flag.FlagSet, opts *cmdOpts) {
	cmd.StringVar(&opts.name, "n", "", "name of the thread")
	cmd.StringVar(&opts.name, "name", "", "name of the thread")

	cmd.StringVar(&opts.path, "p", "", "top-level path for thread files")
	cmd.StringVar(&opts.path, "path", "", "top-level path for thread files")
}

// TODO: args vs flags
func parseArgs(cmd *flag.FlagSet, opts *cmdOpts, args []string) error {
	if len(args) == 0 {
		return errs.ErrNoArgs
	}
	err := cmd.Parse(args)
	if err != nil {
		return err
	}

	if opts.name == "" {
		return errors.New("argument 'name' cannot be empty")
	}
	if opts.path == "" {
		return errors.New("argument 'path' cannot be empty")
	}
	return nil
}

func setUsage(cmd *flag.FlagSet) {
	cmd.Usage = func() {
		fmt.Printf(usage, cmd.Name(), cmd.Name())
	}
}

const usage = `%s regenerates an html file from a previously saved thread

Usage:
  %s [flags]

Flags:
  -n, --name string  name of the thread
  -p, --path string  top-level path for thread files`