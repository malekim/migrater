package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/malekim/migrater/internal/utils"
	"github.com/malekim/migrater/pkg/migrater"

	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration:generate",
	Short: "Add migration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return fmt.Errorf("%s requires a subcommand", cmd.Name())
	},
}

func addMongoMigrationFile(cmd *cobra.Command, args []string) error {
	tmpl := migrater.MongoStub
	timestamp := time.Now().Unix()
	name := fmt.Sprintf("%d.go", timestamp)
	t := template.Must(template.New("").Parse(tmpl))
	path := filepath.Join("app", "migrations", name)
	if err := utils.EnsureDir(path); err != nil {
		log.Printf("Error creating dir: %s", err.Error())
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return err
	}
	defer f.Close()

	vars := struct {
		Timestamp int64
	}{
		timestamp,
	}

	err = t.Execute(f, vars)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Created %s\n", name)
	return nil
}

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Add mongo migration file",
	RunE: addMongoMigrationFile,
}

func init() {
	rootCmd.AddCommand(migrationCmd)
	migrationCmd.AddCommand(mongoCmd)
}
