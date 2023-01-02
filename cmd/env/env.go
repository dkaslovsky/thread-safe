package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Path is the name of the environment variable indicating path for saving threads
	Path = "THREAD_SAFE_PATH"
	// Token is the name of the environment variable containing the Twitter API bearer token
	Token = "THREAD_SAFE_TOKEN" // nolint:gosec
	// tokenFileName is the name of the file in the user's $HOME directory containing the Twitter API bearer token
	tokenFileName = ".thread-safe"
)

// Args holds environment variable values
type Args struct {
	Path  string
	Token string
}

// Parse parses values from the environment
func Parse() *Args {
	path := "."
	if p, ok := os.LookupEnv(Path); ok {
		path = p
	}

	token := ""
	if t, ok := os.LookupEnv(Token); ok {
		token = t
	} else {
		token = readTokenFile()
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

func readTokenFile() string {
	tokenFilePath := filepath.Clean(filepath.Join(os.ExpandEnv("$HOME"), tokenFileName))

	file, err := os.Open(tokenFilePath)
	if err != nil {
		return ""
	}
	defer func() {
		_ = file.Close()
	}()

	var token string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			continue
		}
		if strings.ToLower(strings.TrimSpace(lineParts[0])) == "token" {
			token = strings.TrimSpace(lineParts[1])
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return token
}

var usage = `Environment Variables:
  %s	top level path for thread files (current directory if unset)
  %s	bearer token for Twitter API`
