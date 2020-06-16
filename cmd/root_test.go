package cmd

import (
	"testing"
)

func TestRoot(t *testing.T) {
	args := []string{}
	root(rootCmd, args)
}

func TestExecute(t *testing.T) {
	Execute()
}
