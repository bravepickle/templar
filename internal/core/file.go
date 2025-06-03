package core

import (
	"os"
)

// FileExists checks if file exists
func FileExists(path string) bool {
	if path == `` {
		return false
	}

	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}

// FileContents reads all text file contents from file path as string
func FileContents(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return ``, err
	}

	return string(contents), nil
}
