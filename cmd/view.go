package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("pm created")
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
