package thread

import (
	"errors"
	"flag"
	"fmt"

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
	client := twitter.NewTwitterClient(opts.token)
	th, err := thread.NewThread(client, opts.name, opts.tweetID)
	if err != nil {
		return fmt.Errorf("failed to parse thread: %w", err)
	}

	threadDir := thread.Dir(opts.path, opts.name)

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
	token         string // TODO: env
	tweetID       string
	name          string
	path          string // TODO: env
	noAttachments bool
}

func attachOpts(cmd *flag.FlagSet, opts *cmdOpts) {
	cmd.StringVar(&opts.token, "t", "", "Twitter API bearer token")
	cmd.StringVar(&opts.token, "token", "", "Twitter API bearer token")

	cmd.StringVar(&opts.tweetID, "i", "", "ID of the last tweet in a single-author thread")
	cmd.StringVar(&opts.tweetID, "id", "", "ID of the last tweet in a single-author thread")

	cmd.StringVar(&opts.name, "n", "", "name for the thread")
	cmd.StringVar(&opts.name, "name", "", "name for the thread")

	cmd.StringVar(&opts.path, "p", "", "top-level path for thread files")
	cmd.StringVar(&opts.path, "path", "", "top-level path for thread files")

	cmd.BoolVar(&opts.noAttachments, "no-attachments", false, "do not download media attachments")
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

const usage = `%s saves thread content and generates a local html file

Usage:
  %s [flags]

Flags:
  -t, --token string  Twitter API bearer token
  -i, --id    string  id of the last tweet in a single-author thread
  -n, --name  string  name of the thread
  -p, --path  string  top-level path for thread files`
