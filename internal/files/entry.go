package files

import (
	"fmt"
	"io/fs"
	"slices"
)

type DirEntry struct {
	fs.DirEntry
}

// Gets the name of the entry. If the entry is a directory, a "/" is appended to the name
func (e DirEntry) Name(addSlash bool) string {
	fileName := e.DirEntry.Name()
	if addSlash && e.IsDir() {
		fileName += "/"
	}

	return fileName
}

// Gets string representation of the permissions of the entry (e.g. "drwxr-xr-x").
// If the permissions cannot be determined, "???" is returned.
//
// Uses [io/fs.FileMode.String] internally
func (e DirEntry) Permissions() string {
	info, err := e.Info()
	if err != nil {
		return "???"
	}

	return info.Mode().String()
}

// Gets the size of the entry in SI format.
// E.g. "1.2 kB", "3.4 MB", "5.6 GB", etc.
//
// If the entry is a directory, an empty string is returned.
func (e DirEntry) Size() string {
	info, err := e.Info()
	if err != nil {
		return "0"
	}

	if info.IsDir() {
		return ""
	}

	// From https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
	numBytes := info.Size()

	const unit = 1000
	if numBytes < unit {
		return fmt.Sprintf("%d B", numBytes)
	}

	div, exp := int64(unit), 0
	for n := numBytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(numBytes)/float64(div), "kMGTPE"[exp])
}

// Sorts by directories first, then by name
func Sort(entries []DirEntry) {
	sortDir := make([]DirEntry, 0, len(entries))
	sortFile := make([]DirEntry, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			sortDir = append(sortDir, entry)
		} else {
			sortFile = append(sortFile, entry)
		}
	}

	sortEntry := func(a, b DirEntry) int {
		if a.Name(false) < b.Name(false) {
			return -1
		}

		if a.Name(false) > b.Name(false) {
			return 1
		}

		return 0
	}

	slices.SortFunc(sortDir, sortEntry)
	slices.SortFunc(sortFile, sortEntry)

	copy(entries, append(sortDir, sortFile...))
}
