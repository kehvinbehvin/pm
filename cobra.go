package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "pm",
  Short: "pm is a very fast static site generator",
  Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Root Command")
  },
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
	   	epicFlag, _ := cmd.Flags().GetString("epic")
	   	storyFlag, _ := cmd.Flags().GetString("story")
	   	taskFlag, _ := cmd.Flags().GetString("task")
	   	fmt.Println(epicFlag);
	   	fmt.Println(storyFlag);
	   	fmt.Println(taskFlag);
	},
}

var epicCmd = &cobra.Command{
	Use:   "epics",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
	   	toComplete := args[0];
	   	epicTrie := Load("epic");
	   	suggestions := epicTrie.loadWordsFromPrefix(toComplete);
	   	fmt.Println(strings.Join(suggestions, "\n"))
	},
}

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
	   	toComplete := args[0];
	   	storyTrie := Load("story");
	   	suggestions := storyTrie.loadWordsFromPrefix(toComplete);
	   	fmt.Println(strings.Join(suggestions, "\n"))
	},
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
	   	toComplete := args[0];
	   	taskTrie := Load("task");
	   	suggestions := taskTrie.loadWordsFromPrefix(toComplete);
	   	fmt.Println(strings.Join(suggestions, "\n"))
	},
}

func init() {
  rootCmd.AddCommand(addCmd)
  rootCmd.AddCommand(epicCmd)
  rootCmd.AddCommand(storyCmd)
  rootCmd.AddCommand(taskCmd)

  // Add flags to the add command for epic, task, and story
  addCmd.Flags().StringP("epic", "e", "", "Add an epic")
  addCmd.Flags().StringP("story", "s", "", "Add a story")
  addCmd.Flags().StringP("task", "t", "", "Add an task")
  rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(os.Stderr, err)
    os.Exit(1)
  }
}
