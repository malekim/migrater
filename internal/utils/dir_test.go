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

	// clear after test
	os.RemoveAll(dir)
}

func TestEnsureDirError(t *testing.T) {
	dir := filepath.Join("tmptest")
	// force error and create dir
	// with chmod only to read
	err := os.MkdirAll(dir, 0444)
	if err != nil {
		t.Error("Error during creating the directory")
	}
	file := filepath.Join("tmptest", "app", "test.txt")
	err = EnsureDir(file)
	if err == nil {
		t.Error("There should be an error")
	}
	// remove testing dir
	err = os.RemoveAll(dir)
	if err != nil {
		t.Errorf("Unsuccessful clear %s", dir)
	}
}
