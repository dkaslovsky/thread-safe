package errs

import (
	"errors"
	"fmt"

	"github.com/dkaslovsky/thread-safe/cmd/env"
)

var (
	// ErrNoArgs is returned when no arguments are passed to a command
	ErrNoArgs = errors.New("missing required argument(s)")

	// ErrEmptyPath is a fatal returned when an empty path is received by a command
	ErrEmptyPath = fmt.Errorf("fatal: path could not be determined from %s or the current directory", env.Path)
)
