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
	VarPath = "THREAD_SAFE_PATH"
	// Token is the name of the environment variable containing the Twitter API bearer token
	VarToken = "THREAD_SAFE_TOKEN" // nolint:gosec
	// tokenFileDir is the directory containing the file tokenFileName
	tokenFileDir = "${HOME}"
	// tokenFileName is the name of the file in the user's $HOME directory containing the Twitter API bearer token
	tokenFileName = ".thread-safe" // nolint:gosec
)

// Args holds environment variable values
type Args struct {
	Path  string
	Token string
}

// Parse parses values from the environment
func Parse() *Args {
	path := "."
	if p, ok := os.LookupEnv(VarPath); ok {
		path = p
	}

	token := ""
	if t, ok := os.LookupEnv(VarToken); ok {
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
	return fmt.Sprintf(usage, VarPath, VarToken, TokenFilePath())
}

// TokenFilePath returns the unexpanded path to the file containing the Twitter API bearer token
func TokenFilePath() string {
	return filepath.Clean(filepath.Join(tokenFileDir, tokenFileName))
}

// tokenFilePathExpanded returns the full path to the file containing the Twitter API bearer token
func tokenFilePathExpanded() string {
	return filepath.Clean(filepath.Join(os.ExpandEnv(tokenFileDir), tokenFileName))
}

func readTokenFile() string {
	file, err := os.Open(tokenFilePathExpanded())
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
  %s	bearer token for Twitter API (overrides value read from "%s" if set)`
