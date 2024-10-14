package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
	   	fmt.Println("Bye");
// 	   	switch args[0] {
// 	   	  default:
// 	   	    fmt.Println(args[0])
// 	   	}
      value, _ := cmd.Flags().GetString("epic")
      fmt.Println(value);
	},
}

var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long:  `To load completions run: source <(pm completion bash)`,
		Args:  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				fmt.Println("Unsupported shell")
			}
		},
}

func suggestFromTrie(trie index) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Check if the toComplete string is empty before proceeding
		if len(toComplete) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Trigger suggestions from the Trie based on user input
		suggestions := trie.loadWordsFromPrefix(toComplete)
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}

func init() {
  rootCmd.AddCommand(addCmd)

  // Add flags to the add command for epic, task, and story
  addCmd.Flags().StringP("epic", "e", "", "Add an epic")
  addCmd.Flags().StringP("task", "t", "", "Add a task")
  addCmd.Flags().StringP("story", "s", "", "Add a story")

  rootCmd.AddCommand(completionCmd)

  epicTrie := index.Load("epic")
  taskTrie := index.Load("task")
  storyTrie := index.Load("story")

  // Attach autocompletion to the flags
  addCmd.RegisterFlagCompletionFunc("epic", suggestFromTrie(epicTrie))
  addCmd.RegisterFlagCompletionFunc("task", suggestFromTrie(taskTrie))
  addCmd.RegisterFlagCompletionFunc("story", suggestFromTrie(storyTrie))
}
