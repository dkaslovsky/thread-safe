package regen

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/dkaslovsky/thread-safe/cmd/env"
	"github.com/dkaslovsky/thread-safe/cmd/errs"
	"github.com/dkaslovsky/thread-safe/pkg/thread"
)

// Run executes the package's (sub)command
func Run(appName string, args []string) error {
	cmd := flag.NewFlagSet("regen", flag.ExitOnError)
	opts := &cmdOpts{}
	attachOpts(cmd, opts)
	setUsage(appName, cmd)

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
	th, err := thread.FromJSON(opts.path, opts.name)
	if err != nil {
		return fmt.Errorf("failed to load thread from file: %w", err)
	}

	tErr := th.ToHTML(opts.template, opts.css)
	if tErr != nil {
		return fmt.Errorf("failed to write thread HTML file: %w", tErr)
	}

	return nil
}

type cmdOpts struct {
	// Args
	name string
	// Flags
	css      string
	template string
	// Environment variables
	path string
}

func attachOpts(cmd *flag.FlagSet, opts *cmdOpts) {
	cmd.StringVar(&opts.css, "c", "", "optional path to CSS file")
	cmd.StringVar(&opts.css, "css", "", "optional path to CSS file")

	cmd.StringVar(&opts.template, "t", "", "optional path to template file")
	cmd.StringVar(&opts.template, "template", "", "optional path to template file")
}

func parseArgs(cmd *flag.FlagSet, opts *cmdOpts, args []string) error {
	if len(args) == 0 {
		return errs.ErrNoArgs
	}
	err := cmd.Parse(args)
	if err != nil {
		return err
	}
	opts.name = cmd.Arg(0)

	envArgs := env.Parse()
	opts.path = envArgs.Path

	if opts.path == "" {
		return errs.ErrEmptyPath
	}
	if strings.TrimSpace(opts.name) == "" {
		return errors.New("argument 'name' cannot be empty")
	}
	return nil
}

func setUsage(appName string, cmd *flag.FlagSet) {
	cmd.Usage = func() {
		fmt.Printf(usage, cmd.Name(), appName, cmd.Name())
		fmt.Printf("\n\n%s\n", env.Usage())
	}
}

const usage = `'%s' regenerates an html file from a previously saved thread

Usage:
  %s %s [flags] <name>

Args:
  name  string  name given to the thread

Flags:
  -c, --css       string  optional path to CSS file
  -t, --template  string  optional path to template file`
