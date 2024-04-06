package files

import (
	"io/fs"
	"strconv"
)

type DirEntry struct {
	fs.DirEntry
}

// Gets the name of the entry. If the entry is a directory, a "/" is appended to the name
func (e DirEntry) Name() string {
	fileName := e.DirEntry.Name()
	if e.IsDir() {
		fileName += "/"
	}

	return fileName
}

func (e DirEntry) Permissions() string {
	info, err := e.Info()
	if err != nil {
		return "???"
	}

	return info.Mode().String()
}

func (e DirEntry) Size() string {
	info, err := e.Info()
	if err != nil {
		return "0"
	}

	return strconv.FormatInt(info.Size(), 10)
}
