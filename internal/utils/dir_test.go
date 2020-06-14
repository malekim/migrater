package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	file := filepath.Join("tmptest", "test.txt")
	if err := EnsureDir(file); err != nil {
		t.Errorf("EnsureDir was unable to create path for %s", file)
	}

	dir := filepath.Dir(file)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("EnsureDir did not create path for %s", dir)
	}
	// force error
	if err := EnsureDir("~"); err == nil {
		t.Errorf("EnsureDir should return error")
	}

	// clear after test
	os.RemoveAll(dir)
}
