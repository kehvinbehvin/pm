package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "pm",
  Short: "pm is your best friend",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Root Command")
  },
}

func init() {
  var eValues, sValues, tValues []string

  var initCmd = &cobra.Command{
    	Use:   "init",
    	Short: "Initialize a new .pm project",
    	Run: func(cmd *cobra.Command, args []string) {
        err := os.Mkdir("./.pm", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating pm directory: %v\n", err)
        }

        err = os.Mkdir("./.pm/blobs", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating blobs directory: %v\n", err)
        }

        err = os.Mkdir("./.pm/trie", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating trie directory: %v\n", err)
        }

        err = os.Mkdir("./.pm/trie/epic", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating trie directory: %v\n", err)
        }

        err = os.Mkdir("./.pm/trie/story", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating trie directory: %v\n", err)
        }

        err = os.Mkdir("./.pm/trie/task", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating trie directory: %v\n", err)
        }


        err = os.Mkdir("./.pm/dag", os.ModePerm)
        if err != nil && !os.IsExist(err) {
          fmt.Printf("Error creating dag directory: %v\n", err)
        }
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
  	   	for _, value := range eValues {
  	   	  fmt.Println(value);
  	   	}
  	},
  }

  var epicCmd = &cobra.Command{
  	Use:   "epics",
  	Short: "Epic suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	epicTrie := Load("epic");
  	   	suggestions := epicTrie.loadWordsFromPrefix(toComplete);
  	   	fmt.Println(strings.Join(suggestions, "\n"))
  	},
  }

  var storyCmd = &cobra.Command{
  	Use:   "story",
  	Short: "Story suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	storyTrie := Load("story");
  	   	suggestions := storyTrie.loadWordsFromPrefix(toComplete);
  	   	fmt.Println(strings.Join(suggestions, "\n"))
  	},
  }

  var taskCmd = &cobra.Command{
  	Use:   "task",
  	Short: "Task suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	taskTrie := Load("task");
  	   	suggestions := taskTrie.loadWordsFromPrefix(toComplete);
  	   	fmt.Println(strings.Join(suggestions, "\n"))
  	},
  }

  rootCmd.AddCommand(initCmd)
  rootCmd.AddCommand(addCmd)
  rootCmd.AddCommand(epicCmd)
  rootCmd.AddCommand(storyCmd)
  rootCmd.AddCommand(taskCmd)

  // Add flags to the add command for epic, task, and story
  addCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
  addCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add a story")
  addCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add an task")

  rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
