package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
					return
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

			epicFile, epicErr := os.Create("./.pm/trie/epic")
			if epicErr != nil {
				fmt.Printf("Error creating epic trie: %v\n", epicErr)
			}
			epicTrie := NewTrie("epic")
			epicTrie.Save()

			defer epicFile.Close()

			storyFile, storyErr := os.Create("./.pm/trie/story")
			if storyErr != nil {
				fmt.Printf("Error creating story trie: %v\n", storyErr)
			}
			storyTrie := NewTrie("story")
			storyTrie.Save()

			defer storyFile.Close()

			taskFile, taskErr := os.Create("./.pm/trie/task")
			if taskErr != nil {
				fmt.Printf("Error creating task trie: %v\n", taskErr)
			}
			taskTrie := NewTrie("task")
			taskTrie.Save()

			defer taskFile.Close()

			err = os.Mkdir("./.pm/dag", os.ModePerm)
			if err != nil && !os.IsExist(err) {
				fmt.Printf("Error creating dag directory: %v\n", err)
			}
			pmDag := newDag("pmDag")
			pmDag.SaveDag()

			tmpFile, tmpErr := os.Create("./.pm/tmp")
			if tmpErr != nil {
				fmt.Printf("Error creating tmp file: %v\n", tmpErr)
			}

			defer tmpFile.Close()
		},
	}

	var addCmd = &cobra.Command{
		Use: "add",
		Short: "Add epics, stories or tasks",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			pmDag := LoadDag("pmDag")
			defer pmDag.SaveDag()
			
			var nodesToSave []*Vertex;

			if epics > 0 {
				for _, value := range eValues {
					epic := newVertex(value)
					nodesToSave = append(nodesToSave, epic)
				}

			}

			if stories > 0 {
				for _, value := range sValues {
					story := newVertex(value)
					nodesToSave = append(nodesToSave, story)
				}

			}

			if tasks > 0 {
				for _, value := range tValues {
					task := newVertex(value)
					nodesToSave = append(nodesToSave, task)
				}

			}

			for _, vs := range nodesToSave {
				pmDag.addVertex(vs);
			}

			epicTrie := Load("epic")
			for _, value := range eValues {
				commit(value, "", epicTrie)
			}

			storyTrie := Load("story")
			for _, value := range sValues {
				commit(value, "", storyTrie)
			}

			taskTrie := Load("task")
			for _, value := range tValues {
				commit(value, "", taskTrie)
			}
		},
	}

	var linkCmd = &cobra.Command{
		Use:   "link",
		Short: "Initialize a new .pm project",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			pmDag := LoadDag("pmDag")
			defer pmDag.SaveDag()

			if tasks > 0 {
				if epics > 1 || stories > 1 {
					// Not allowed, invalid relationship.
					fmt.Println("A task can only belong to 1 story and epic")
					return
				}

				if epics == 1 && stories == 1 {
					// 1 Epic and 1 Story per multiple task entry
					epicValue := eValues[0]
					epicVertex := pmDag.retrieveVertex(epicValue)

					storyValue := sValues[0]
					storyVertex := pmDag.retrieveVertex(storyValue)

					pmDag.addEdge(epicVertex, storyVertex)
					for _, value := range tValues {
						taskVertex := pmDag.retrieveVertex(value)
						pmDag.addEdge(storyVertex, taskVertex)
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
					epicValue := eValues[0]
					epicVertex := pmDag.retrieveVertex(epicValue)

					for _, value := range sValues {
						storyVertex := pmDag.retrieveVertex(value)
						pmDag.addEdge(epicVertex, storyVertex)
					}
				}
			} else {
				fmt.Println("Cannot create dependencies between epics")
				return
			}

		},
	}

	var listCmd = &cobra.Command{
		Use: "list",
		Short: "List all nodes of a type",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please choose what to display");
		},

	}

	var listEpicsCmd = &cobra.Command{
		Use: "epic",
		Short: "List all epics",
		Run: func(cmd *cobra.Command, args []string) {
			epicTrie := Load("epic")
			allEpics, err := epicTrie.loadAllWords();
			if err != nil {
				return
			}

			buildList(allEpics, "Epics");
		},
	}

	var listStoriesCmd = &cobra.Command{
		Use: "story",
		Short: "List all story",
		Run: func(cmd *cobra.Command, args []string) {
			storyTrie := Load("story")
			allEpics, err := storyTrie.loadAllWords();
			if err != nil {
				return
			}

			buildList(allEpics, "Story");
		},
	}

	var listTasksCmd = &cobra.Command{
		Use: "task",
		Short: "List all task",
		Run: func(cmd *cobra.Command, args []string) {
			taskTrie := Load("task")
			allTasks, err := taskTrie.loadAllWords();
			if err != nil {
				return
			}

			buildList(allTasks, "Tasks");
		},
	}





	// Create a new command for viewing epics
	var viewCmd = &cobra.Command{
		Use:   "view",
		Short: "View epics, stories or tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)
			pmDag := LoadDag("pmDag")

			total := epics + stories + tasks
			if total > 1 {
				fmt.Println("Only allow to list 1 of a kind at a time")
				return
			}

			var node *Vertex
			var nodeType string

			if epics > 0 {
				// Simulate getting epic, stories, and tasks data
				node = pmDag.retrieveVertex(eValues[0])
				nodeType = "Stories"
			} else if stories > 0 {
				node = pmDag.retrieveVertex(sValues[0])
				nodeType = "Tasks"
			} else if tasks > 0 {
				node = pmDag.retrieveVertex(tValues[0])
				nodeType = "SubTasks"
			} else {
				fmt.Println("Not sure what you want to display")	
				return
			}
			
			if (len(node.Children) == 0) {
				fmt.Println("Nothing to list")
				return
			}

			nodeKeys := make([]string, len(node.Children))
			i := 0
			for k := range node.Children {
				nodeKeys[i] = k
				i++
			}
			
			buildList(nodeKeys, nodeType)

		},
	}

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Open the ./tmp file in the user's preferred editor",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			total := epics + stories + tasks
			if total > 1 {
				fmt.Println("Only allow to 1 file at a time")
				return
			}

			// Define the path to the file
			filePath := "./.pm/tmp"

			// Get the user's preferred editor from the $EDITOR environment variable
			editor := os.Getenv("EDITOR")
			if editor == "" {
				// Fallback to a default editor if $EDITOR is not set
				editor = "vim"
			}

			// Open the file in the editor
			err := openEditor(editor, filePath)
			if err != nil {
				return
			}

			// Read the file contents after the user has saved and exited
			fileContent, err := os.Open(filePath)
			if err != nil {
				return
			}
			defer fileContent.Close()

			// Print the file content (or you can process it as needed)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return
			}

			if epics > 0 {
				epicTrie := Load("epic")
				defer epicTrie.Save()

				updateErr := updateBlobContent(eValues[0], string(content), epicTrie)
				if updateErr != nil {
					fmt.Println("Update file error")
					return
				}
			} else if stories > 0 {
				storyTrie := Load("story")
				defer storyTrie.Save()

				updateErr := updateBlobContent(sValues[0], string(content), storyTrie)
				if updateErr != nil {
					fmt.Println("Update file error")
					return
				}
			} else if tasks > 0 {
				taskTrie := Load("task")
				defer taskTrie.Save()

				updateErr := updateBlobContent(tValues[0], string(content), taskTrie)
				if updateErr != nil {
					fmt.Println("Update file error")
					return
				}
			}

			err = emptyFile(filePath)
			if err != nil {
				fmt.Println("Cannot empty file")
				return
			}
			return
		},
	}

	var epicCmd = &cobra.Command{
		Use:   "epics",
		Short: "Epic suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			epicTrie := Load("epic")
			if epicTrie != nil {
				suggestions := epicTrie.loadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	var storyCmd = &cobra.Command{
		Use:   "story",
		Short: "Story suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			storyTrie := Load("story")
			if storyTrie != nil {
				suggestions := storyTrie.loadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	var taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Task suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			taskTrie := Load("task")
			if taskTrie != nil {
				suggestions := taskTrie.loadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(epicCmd)
	rootCmd.AddCommand(storyCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(listCmd)
	
	listCmd.AddCommand(listEpicsCmd)
	listCmd.AddCommand(listStoriesCmd)
	listCmd.AddCommand(listTasksCmd)

	// Add flags to the add command for epic, task, and story
	linkCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	linkCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add a story")
	linkCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add an task")

	addCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	addCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add a story")
	addCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add an task")

	viewCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	viewCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add an story")
	viewCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add an task")

	editCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	editCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add an story")
	editCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add an task")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func walkDAG(writer *tabwriter.Writer, nodeType string, epic *Vertex) {
	// epicTree := tree.Root(epic.ID);
	//
	// // Iterate over the stories (children of the epic)
	// for _, story := range epic.Children {
	// 	storyTree := tree.Root(story.ID)
	//
	// 	// Iterate over the tasks (children of the story)
	// 	for _, task := range story.Children {
	// 		// Print the task (level 2)
	// 		storyTree.Child(task.ID)
	//
	// 	}
	// 	epicTree.Child(storyTree)
	// }
	//
	// fmt.Println(epicTree);
}

func buildList(rows []string, listHeader string) {
	fmt.Println(rows);
	rowIds := [][]string{}
	for _, value := range rows {
		row := []string{value}
		rowIds = append(rowIds, row)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers(listHeader).
		Rows(rowIds...)
	fmt.Println(t)
}

func openEditor(editor string, filePath string) error {
	// Create an exec command to open the file in the editor
	cmd := exec.Command(editor, filePath)

	// Set the command to use the same standard input, output, and error streams as the Go process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open file in editor: %v", err)
	}

	return nil
}

func emptyFile(filePath string) error {
	// Truncate the file to zero length
	err := os.Truncate(filePath, 0)
	if err != nil {
		return fmt.Errorf("failed to truncate the file: %v", err)
	}
	return nil
}
