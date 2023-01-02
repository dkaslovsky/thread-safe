package env

import (
	"fmt"
	"os"
)

const (
	// Path is the name of the environment variable indicating path for saving threads
	Path = "THREAD_SAFE_PATH"
	// Token is the name of the environment variable indicating Twitter API bearer token
	Token = "THREAD_SAFE_TOKEN" // nolint:gosec
)

// Args holds environment variable values
type Args struct {
	Path  string
	Token string
}

// Parse loads environment variables or sets defaults
func Parse() *Args {
	path := "."
	if p, ok := os.LookupEnv(Path); ok {
		path = p
	}

	token := ""
	if t, ok := os.LookupEnv(Token); ok {
		token = t
	}

	return &Args{
		Path:  path,
		Token: token,
	}
}

// Usage returns a string describing the environment variables
func Usage() string {
	return fmt.Sprintf(usage, Path, Token)
}

var usage = `Environment Variables:
  %s	top level path for thread files (current directory if unset)
  %s	bearer token for Twitter API`
