package files

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ErrorSanitiseUnauthorized = errors.New("unauthorized path")
	ErrorSanitiseNotExists	= errors.New("path does not exist")
)

// Sanitises a path by resolving symlinks and checking if it is within the root directory.
//
// Returns the absolute path if it is within the root directory, otherwise an error.
func SanitisePath(rootDir, path string) (string, error) {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(rootDir, path))
	if err != nil {
		return "", err
	}

	absPath, err = filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", ErrorSanitiseNotExists
	}

	if !hasPrefix(rootDir, absPath) {
		return "", ErrorSanitiseUnauthorized
	}

	return absPath, nil
}

func hasPrefix(prefix, path string) bool {
	// Paths in windows are case-insensitive
	if runtime.GOOS == "windows" {
		return strings.EqualFold(prefix, path[:len(prefix)])
	}

	return prefix == path[:len(prefix)]
}

// Gets a path starts with '/' and relative to the actual rootDir
// Should be used for display only
func RelativeRoot(rootDir string, path string) string {
	relativePath, err := filepath.Rel(rootDir, path)
	if err != nil {
		return "/"
	}

	if relativePath == "." {
		return "/"
	}

	return "/" + relativePath
}
