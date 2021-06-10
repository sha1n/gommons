package os

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandUserPath expands the specified user home relative path to an absolute path.
// Supports '~', '$HOME' for Unix based environments and '%userprofile%' '%homepath%' for windows.
// If the provided path has no home dir prefix the path returns as is. Relative paths are not expanded.
func ExpandUserPath(path string) (expandedPath string, err error) {
	expandedPath = path
	if hasHomeDirPrefix(path) {
		var homeDir string
		if homeDir, err = os.UserHomeDir(); err == nil {
			expandedPath = filepath.Join(homeDir, path[1:])
		}
	}

	return expandedPath, err
}

func hasHomeDirPrefix(path string) bool {
	return strings.HasPrefix(path, "~") ||
		strings.HasPrefix(path, "$HOME") ||
		strings.HasPrefix(strings.ToUpper(path), "%HOMEPATH%") ||
		strings.HasPrefix(strings.ToUpper(path), "%USERPROFILE%")
}
