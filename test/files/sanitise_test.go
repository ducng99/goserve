package files_test

import (
	"path/filepath"
	"testing"

	"r.tomng.dev/goserve/internal/files"
)

var (
	RootDir = filepath.Join("..", "testdata", "root")
)

var (
	File1Path = "file1.txt"
	Dir1Path = "dir1"
	Dir1File1Path = filepath.Join("dir1", "file1.txt")
	FileNotExists = "thisisdefinitelynotafilelookawayrightnow_test.go"
	DirNotExists = "thisisdefinitelynotadirwhoevercreatesthisdirisamonster"
	FileInaccessible = filepath.Join("..", "inaccessible.txt")
	File1Symlink = "file1_lnk"
	FileSymlinkInaccessible = "inaccessible_lnk"
	Dir1Symlink = "dir1_lnk"
	DirSymlinkInaccessible = "inaccessible_dir_lnk"
)

func TestSanitiseFile1OK(t *testing.T) {
	rootDir := RootDir
	path := File1Path

	_, err := files.SanitisePath(rootDir, path)

	if err != nil {
		t.Fatalf("Root dir: %s - Path: %s - Error: %v", rootDir, path, err)
	}
}

func TestSanitiseDir1OK(t *testing.T) {
	rootDir := RootDir
	path := Dir1Path

	_, err := files.SanitisePath(rootDir, path)

	if err != nil {
		t.Fatalf("Root dir: %s - Path: %s - Error: %v", rootDir, path, err)
	}
}

func TestSanitiseDir1File1OK(t *testing.T) {
	rootDir := RootDir
	path := Dir1File1Path

	_, err := files.SanitisePath(rootDir, path)

	if err != nil {
		t.Fatalf("Root dir: %s - Path: %s - Error: %v", rootDir, path, err)
	}
}

func TestSanitiseFileNotExists(t *testing.T) {
	rootDir := RootDir
	path := FileNotExists

	_, err := files.SanitisePath(rootDir, path)

	if err == nil {
		t.Fatalf("Root dir: %s - Path: %s. No error", rootDir, path)
	}
}

func TestSanitiseDirNotExists(t *testing.T) {
	rootDir := RootDir
	path := DirNotExists

	_, err := files.SanitisePath(rootDir, path)

	if err == nil {
		t.Fatalf("Root dir: %s - Path: %s. No error", rootDir, path)
	}
}

func TestSanitiseFileInaccessible(t *testing.T) {
	rootDir := RootDir
	path := FileInaccessible

	_, err := files.SanitisePath(rootDir, path)

	if err == nil {
		t.Fatalf("Root dir: %s - Path: %s. No error", rootDir, path)
	}
}

func TestSanitiseFile1SymlinkOK(t *testing.T) {
	rootDir := RootDir
	path := File1Symlink

	_, err := files.SanitisePath(rootDir, path)

	if err != nil {
		t.Fatalf("Root dir: %s - Path: %s - Error: %v", rootDir, path, err)
	}
}

func TestSanitiseFileSymlinkInaccessible(t *testing.T) {
	rootDir := RootDir
	path := FileSymlinkInaccessible

	_, err := files.SanitisePath(rootDir, path)

	if err == nil {
		t.Fatalf("Root dir: %s - Path: %s. No error", rootDir, path)
	}
}

func TestSanitiseDir1SymlinkOK(t *testing.T) {
	rootDir := RootDir
	path := Dir1Symlink

	_, err := files.SanitisePath(rootDir, path)

	if err != nil {
		t.Fatalf("Root dir: %s - Path: %s - Error: %v", rootDir, path, err)
	}
}

func TestSanitiseDirSymlinkInaccessible(t *testing.T) {
	rootDir := RootDir
	path := DirSymlinkInaccessible

	_, err := files.SanitisePath(rootDir, path)

	if err == nil {
		t.Fatalf("Root dir: %s - Path: %s. No error", rootDir, path)
	}
}
