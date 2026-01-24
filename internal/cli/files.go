package cli

import (
	"os"
	"path/filepath"
	"strings"
)

// isDirectory checks if the given path is a directory
func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// findMarkdownFiles recursively finds all markdown files in a directory
func findMarkdownFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has markdown extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".markdown" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// resolvePaths takes a path argument and returns a list of markdown files to process.
// If the path is a file, it returns a slice with just that file.
// If the path is a directory, it recursively finds all markdown files.
func resolvePaths(path string) ([]string, error) {
	isDir, err := isDirectory(path)
	if err != nil {
		return nil, err
	}

	if isDir {
		return findMarkdownFiles(path)
	}

	return []string{path}, nil
}
