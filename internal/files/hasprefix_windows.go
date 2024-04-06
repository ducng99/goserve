//go:build windows
// +build windows

package files

import "strings"

func hasPrefix(prefix, path string) bool {
	// Paths in windows are case-insensitive
	return strings.EqualFold(prefix, path[:len(prefix)])
}
