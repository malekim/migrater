package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestMigrationRoot(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}
	args := []string{}
	err := migrationRoot(cmd, args)
	if err == nil {
		t.Error("Expected root to return an error")
	}
}

func TestAddMongoMigrationFile(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}
	args := []string{}
	err := addMongoMigrationFile(cmd, args)
	if err != nil {
		t.Error("Error during call addMongoMigrationFile command")
	}
	// clean after test
	dir := filepath.Join("app")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Dir %s should exist", dir)
	}
	// remove testing dir
	err = os.RemoveAll(dir)
	if err != nil {
		t.Errorf("Unsuccessful clear %s", dir)
	}
}
