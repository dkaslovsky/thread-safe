package thread

import (
	"os"
	"path/filepath"
	"strings"
)

// Directory represents a directory on the file system for Thread files
type Directory struct {
	path string
}

// NewDirectory constructs a new Directory by formatting the name and joining with topLevelDir
func NewDirectory(topLevelDir string, name string) *Directory {
	return &Directory{
		path: filepath.Clean(filepath.Join(topLevelDir, strings.Replace(name, " ", "_", -1))),
	}
}

// Create creates a Directory
func (d *Directory) Create() error {
	return os.MkdirAll(d.path, 0o750)
}

// Exists evaluates if a Directory exists
func (d *Directory) Exists() bool {
	_, err := os.Stat(d.path)
	return !os.IsNotExist(err)
}

// SubDir constructs a file path by joining any provided subpath strings to the Directory path and
// returns the path and a bool indicating if the constructed path exists
func (d *Directory) SubDir(subpaths ...string) (string, bool) {
	path := d.Join(subpaths...)
	_, err := os.Stat(path)
	return path, !os.IsNotExist(err)
}

// Join constructs a string representing a path to a subdirectory or file of a Directory
func (d *Directory) Join(subpaths ...string) string {
	parts := append([]string{d.path}, subpaths...)
	return filepath.Clean(filepath.Join(parts...))
}

func (d *Directory) String() string {
	return d.path
}
