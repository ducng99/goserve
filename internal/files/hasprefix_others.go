//go:build !windows
// +build !windows

package files

func hasPrefix(prefix, path string) bool {
	return prefix == path[:len(prefix)]
}
