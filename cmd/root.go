package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

var rootCmd = &cobra.Command{
	Use:   "Migrater",
	Short: "A package to handle migrations written in GO",
	Run:   root,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
