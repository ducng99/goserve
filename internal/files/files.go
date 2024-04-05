package files

import (
	"os"
	"path/filepath"
)

type Entry struct {
	Name  string
	Path  string
	IsDir bool
}

type PathType uint8

const (
	PathTypeFile PathType = iota
	PathTypeDirectory
)

// Get the type of a path - file or directory
func GetPathType(path string) (PathType, error) {
	f, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if f.IsDir() {
		return PathTypeDirectory, nil
	}

	return PathTypeFile, nil
}

// Get a list of files and directories in a directory.
// Uses [os.ReadDir] internally
//
// Path returned is relative to the root directory
func GetEntries(rootDir, path string) ([]Entry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		fileName := entry.Name()
		if entry.IsDir() {
			fileName += "/"
		}

		filePath, err := filepath.Rel(rootDir, filepath.Join(path, fileName))
		if err != nil {
			return nil, err
		}

		result = append(result, Entry{
			Name:  fileName,
			Path:  filePath,
			IsDir: entry.IsDir(),
		})
	}

	return result, nil
}
