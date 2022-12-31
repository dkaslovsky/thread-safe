package thread

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd/env"
	"github.com/dkaslovsky/thread-safe/cmd/errs"
	"github.com/dkaslovsky/thread-safe/pkg/thread"
	"github.com/dkaslovsky/thread-safe/pkg/twitter"
)

// Run executes the package's (sub)command
func Run(args []string) error {
	cmd := flag.NewFlagSet("thread", flag.ExitOnError)
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
	if _, err := os.Stat(threadDir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists, rename or delete instead of overwriting", threadDir)
	}

	client := twitter.NewTwitterClient(opts.token)
	th, err := thread.NewThread(client, opts.name, opts.tweetID)
	if err != nil {
		return fmt.Errorf("failed to parse thread: %w", err)
	}

	fErr := th.ToJSON(threadDir)
	if fErr != nil {
		return fmt.Errorf("failed to write thread JSON file: %w", err)
	}

	if !opts.noAttachments {
		err := th.DownloadAttachments(threadDir)
		if err != nil {
			return fmt.Errorf("failed to save thread attachment files: %w", err)
		}
	}

	tErr := th.ToHTML(threadDir)
	if tErr != nil {
		return fmt.Errorf("failed to write thread HTML file: %w", err)
	}

	return nil
}

type cmdOpts struct {
	// Args
	name    string
	tweetID string
	// Flags
	noAttachments bool
	// Environment variables
	path  string
	token string
}

func attachOpts(cmd *flag.FlagSet, opts *cmdOpts) {
	cmd.BoolVar(&opts.noAttachments, "no-attachments", false, "do not download media attachments")
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
	opts.tweetID = cmd.Arg(1)

	envArgs := env.Parse()
	opts.path = envArgs.Path
	opts.token = envArgs.Token

	if opts.path == "" {
		return errs.ErrEmptyPath
	}
	if opts.token == "" {
		return fmt.Errorf("token must be specifed by the environment variable %s and must not be empty", env.Token)
	}
	if opts.name == "" {
		return errors.New("argument 'name' cannot be empty")
	}
	if opts.tweetID == "" {
		return errors.New("argument 'id' cannot be empty")
	}
	return nil
}

func setUsage(cmd *flag.FlagSet) {
	cmd.Usage = func() {
		fmt.Printf(fmt.Sprintf("%s\n", usage), cmd.Name(), cmd.Name())
		fmt.Printf("\n%s\n", env.Usage())
	}
}

const usage = `%s saves thread content and generates a local html file

Usage:
  %s [flags] name last-tweet-id

Args:
  name           string  name to use for the thread
  last-tweet-id  string  id of the last tweet in a single-author thread

Flags:
  --no-attachments  do not download attachments`
