package errs

import "errors"

// ErrorNoArgs is returned when no arguments are passed to a command
var ErrNoArgs = errors.New("missing required argument(s)")
