package env

import "os"

const (
	Path  = "THREAD_SAFE_PATH"
	Token = "THREAD_SAFE_TOKEN" // nolint:gosec
)

type Args struct {
	Path  string
	Token string
}

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
