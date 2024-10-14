package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Mkdir(".pm", 0755)
		if err != nil {
			fmt.Println("pm failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
