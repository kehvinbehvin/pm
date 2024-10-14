package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type trieIndex interface {
  loadWordsFromPrefix(string) []string
}

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "pm is a way of working",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world")
	},
	PersistentPreRun: initialize,
}

func initialize(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing common setup for every command")
	// Add initialization logic here
}

func Execute(container []trieIndex) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}
