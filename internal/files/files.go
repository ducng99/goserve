package files

import (
	"os"
)

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
func GetEntries(absPath string) ([]DirEntry, error) {
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	result := make([]DirEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, DirEntry{entry})
	}

	return result, nil
}
