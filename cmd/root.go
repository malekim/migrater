package cmd

import (
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
	rootCmd.Execute()
}
