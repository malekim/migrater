package cmd

import (
	"fmt"

	"github.com/malekim/migrater/pkg/migrater"

	"github.com/spf13/cobra"
)

func migrationRoot(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return fmt.Errorf("%s requires a subcommand", cmd.Name())
}

var migrationCmd = &cobra.Command{
	Use:   "migration:generate",
	Short: "Add migration file",
	RunE:  migrationRoot,
}

func addMongoMigrationFile(cmd *cobra.Command, args []string) error {
	return migrater.AddMongoMigrationFile()
}

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Add mongo migration file",
	RunE:  addMongoMigrationFile,
}

func init() {
	rootCmd.AddCommand(migrationCmd)
	migrationCmd.AddCommand(mongoCmd)
}
