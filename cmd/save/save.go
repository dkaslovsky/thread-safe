package save

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/dkaslovsky/thread-safe/cmd/env"
	"github.com/dkaslovsky/thread-safe/cmd/errs"
	"github.com/dkaslovsky/thread-safe/pkg/thread"
	"github.com/dkaslovsky/thread-safe/pkg/twitter"
)

// Run executes the package's (sub)command
func Run(appName string, args []string) error {
	cmd := flag.NewFlagSet("save", flag.ExitOnError)
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
	th := thread.New(opts.path, opts.name)

	if th.Dir.Exists() {
		return fmt.Errorf("%s already exists, rename or delete instead of overwriting", th.Dir)
	}

	err := th.Load(twitter.NewClient(opts.token), opts.tweetID)
	if err != nil {
		return fmt.Errorf("failed to parse thread: %w", err)
	}

	dErr := th.Dir.Create()
	if dErr != nil {
		return fmt.Errorf("failed to create thread directory %s: %w", th.Dir, dErr)
	}

	fErr := th.ToJSON()
	if fErr != nil {
		return fmt.Errorf("failed to write thread JSON file: %w", fErr)
	}

	if !opts.noAttachments {
		err := th.DownloadAttachments()
		if err != nil {
			return fmt.Errorf("failed to save thread attachment files: %w", err)
		}
	}

	tErr := th.ToHTML(opts.template, opts.css)
	if tErr != nil {
		return fmt.Errorf("failed to write thread HTML file: %w", tErr)
	}

	return nil
}

type cmdOpts struct {
	// Args
	name    string
	tweetID string
	// Flags
	css           string
	template      string
	noAttachments bool
	// Environment variables
	path  string
	token string
}

func attachOpts(cmd *flag.FlagSet, opts *cmdOpts) {
	cmd.StringVar(&opts.css, "c", "", "path to optional CSS file")
	cmd.StringVar(&opts.css, "css", "", "path to optional CSS file")

	cmd.StringVar(&opts.template, "t", "", "optional path to template file")
	cmd.StringVar(&opts.template, "template", "", "optional path to template file")

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

	tweetID, tErr := parseTweetID(cmd.Arg(1))
	if tErr != nil {
		return tErr
	}
	opts.tweetID = tweetID
	opts.name = cmd.Arg(0)

	envArgs := env.Parse()
	opts.path = envArgs.Path
	opts.token = envArgs.Token

	if opts.path == "" {
		return errs.ErrEmptyPath
	}
	if opts.token == "" {
		return fmt.Errorf("token must be specified in %s or by the environment variable %s", env.TokenFilePath(), env.VarToken)
	}
	if opts.name == "" {
		return errors.New("argument 'name' cannot be empty")
	}
	if opts.tweetID == "" {
		return errors.New("argument 'last-tweet' cannot be empty")
	}
	return nil
}

// parseTweetID extracts a tweet ID from its URL or returns the original input if provided the ID
func parseTweetID(urlOrID string) (string, error) {
	u, err := url.Parse(urlOrID)
	if err != nil {
		// Input is not a URL so return as-is
		return urlOrID, nil
	}

	// Parse ID from URL
	urlParts := strings.Split(u.Path, "/")
	if len(urlParts) == 0 {
		return "", fmt.Errorf("failed to parse tweet ID from URL %s", urlOrID)
	}
	return urlParts[len(urlParts)-1], nil
}

func setUsage(appName string, cmd *flag.FlagSet) {
	cmd.Usage = func() {
		fmt.Printf(usage, cmd.Name(), appName, cmd.Name())
		fmt.Printf("\n\n%s\n", env.Usage())
	}
}

const usage = `'%s' saves thread content and generates a local html file

Usage:
  %s %s [flags] <name> <last-tweet>

Args:
  name           string  name to use for the thread
  last-tweet     string  URL or ID of the last tweet in a single-author thread

Flags:
  -c, --css             string  optional path to CSS file
  -t, --template        string  optional path to template file
      --no-attachments          do not download attachments`
