package thread

import (
	"errors"
	"flag"
	"path/filepath"

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
		return err // TODO: wrap or provide user-friendly message?
	}

	threadDir := filepath.Join(opts.path, opts.name)

	fErr := th.ToJSON(threadDir)
	if fErr != nil {
		return fErr // TODO: wrap or provide user-friendly message?
	}

	if !opts.noAttachments {
		err := th.DownloadAttachments(threadDir)
		if err != nil {
			return err // TODO: wrap or provide user-friendly message?
		}
	}

	tErr := th.ToHTML(threadDir)
	if tErr != nil {
		return tErr // TODO: wrap or provide user-friendly message?
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
	cmd.StringVar(&opts.token, "t", "", "twitter API bearer token")
	cmd.StringVar(&opts.token, "token", "", "twitter API bearer token")

	cmd.StringVar(&opts.tweetID, "i", "", "id of the last tweet in a single-author thread")
	cmd.StringVar(&opts.tweetID, "id", "", "id of the last tweet in a single-author thread")

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

// TODO
func setUsage(cmd *flag.FlagSet) {
	cmd.Usage = func() {}
}
