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
		path: filepath.Join(topLevelDir, strings.Replace(name, " ", "_", -1)),
	}
}

// Exists evaluates if a Directory exists
func (d *Directory) Exists(subpaths ...string) bool {
	_, err := os.Stat(d.Join(subpaths...))
	return !os.IsNotExist(err)
}

// Create creates a Directory
func (d *Directory) Create() error {
	return os.MkdirAll(filepath.Clean(d.path), 0o750)
}

// Join constructs a string representing a path to a subdirectory or file of a Directory
func (d *Directory) Join(subpaths ...string) string {
	parts := append([]string{d.path}, subpaths...)
	return filepath.Join(parts...)
}

func (d *Directory) String() string {
	return d.path
}
