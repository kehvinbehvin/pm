package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "pm",
  Short: "pm is a way of working",
  Run: func(cmd *cobra.Command, args []string) {
   fmt.Println("Hello world");
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(os.Stderr, err)
    os.Exit(1)
  }
}
