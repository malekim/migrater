package utils

import (
	"os"
	"path/filepath"
)

// Ensure that dir exists
func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}
