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
      info, err := os.Stat("./.pm")
      if !os.IsNotExist(err) {
        if info.IsDir() {
          fmt.Printf("Directory already is managed by pm")
          return;
        }
      }

      err = os.Mkdir("./.pm", os.ModePerm)
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

      epicFile, epicErr := os.Create("./.pm/trie/epic");
      if epicErr != nil {
        fmt.Printf("Error creating epic trie: %v\n", epicErr)
      }
      epicTrie := NewTrie("epic");
      epicTrie.Save();

      defer epicFile.Close()

      storyFile, storyErr := os.Create("./.pm/trie/story");
      if storyErr != nil {
        fmt.Printf("Error creating story trie: %v\n", storyErr)
      }
      storyTrie := NewTrie("story");
      storyTrie.Save();

      defer storyFile.Close()

      taskFile, taskErr := os.Create("./.pm/trie/task");
      if taskErr != nil {
        fmt.Printf("Error creating task trie: %v\n", taskErr)
      }
      taskTrie := NewTrie("task");
      taskTrie.Save()

      defer taskFile.Close()

      err = os.Mkdir("./.pm/dag", os.ModePerm)
      if err != nil && !os.IsExist(err) {
        fmt.Printf("Error creating dag directory: %v\n", err)
      }
      pmDag := newDag("pmDag");
      pmDag.SaveDag();
    },
  }

  var addCmd = &cobra.Command{
  	Use:   "add",
  	Short: "Initialize a new .pm project",
  	Run: func(cmd *cobra.Command, args []string) {
  	    epics := len(eValues);
  	    stories := len(sValues);
  	    tasks := len(tValues);

  	    fmt.Println(eValues);
  	    fmt.Println(sValues);
  	    fmt.Println(tValues);

        pmDag := LoadDag("pmDag");
        defer pmDag.SaveDag()

  	    if tasks > 0 {
  	      if epics > 1 || stories > 1 {
  	        // Not allowed, invalid relationship.
  	        fmt.Println("A task can only belong to 1 story and epic")
  	        return;
  	      }

  	      if epics == 1 && stories == 1 {
            // 1 Epic and 1 Story per multiple task entry
            epicValue := eValues[0];
            epicVertex := newVertex(epicValue);
            pmDag.addVertex(epicVertex);
            
            storyValue := sValues[0]
            storyVertex := newVertex(storyValue);
            pmDag.addVertex(storyVertex);

            pmDag.addEdge(epicVertex, storyVertex);
            for _, value := range tValues {
                taskVertex := newVertex(value);
                pmDag.addVertex(taskVertex);
                pmDag.addEdge(storyVertex, taskVertex);
            }
  	      }
  	    } else if stories > 0 {
          if epics > 1 {
            // Not allowed, invalid relationship
  	        fmt.Println("A story can only belong to 1 epic")
            return
          }

          if epics == 1 {
            // 1 Epic per multiple story entry
            epicValue := eValues[0];
            epicVertex := newVertex(epicValue);

            for _, value := range sValues {
                storyVertex := newVertex(value);
                pmDag.addVertex(storyVertex);
                pmDag.addEdge(epicVertex, storyVertex);
            }
          }
  	    } else {
          // No stories or tasks, just epics. No relationships
          for _, value := range eValues {
               vertex := newVertex(value);
               pmDag.addVertex(vertex);
           }
  	    }

  	    epicTrie := Load("epic");
  	    for _, value := range eValues {
  	          commit(value, "", epicTrie);
        }

        storyTrie := Load("story");
        for _, value := range sValues {
              commit(value, "", storyTrie);
        }

        taskTrie := Load("task");
        for _, value := range tValues {
              commit(value, "", taskTrie);
        }
  	},
  }

  var epicCmd = &cobra.Command{
  	Use:   "epics",
  	Short: "Epic suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	epicTrie := Load("epic");
                if epicTrie != nil {
                  suggestions := epicTrie.loadWordsFromPrefix(toComplete);
  	   	  fmt.Println(strings.Join(suggestions, "\n"))
                }
  	},
  }

  var storyCmd = &cobra.Command{
  	Use:   "story",
  	Short: "Story suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	storyTrie := Load("story");
                if storyTrie != nil {
                  suggestions := storyTrie.loadWordsFromPrefix(toComplete);
  	   	  fmt.Println(strings.Join(suggestions, "\n"))
                }
  	},
  }

  var taskCmd = &cobra.Command{
  	Use:   "task",
  	Short: "Task suggestions",
  	Run: func(cmd *cobra.Command, args []string) {
  	   	toComplete := args[0];
  	   	taskTrie := Load("task");
                if taskTrie != nil {
  	   	  suggestions := taskTrie.loadWordsFromPrefix(toComplete);
  	   	  fmt.Println(strings.Join(suggestions, "\n"))
                }
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
