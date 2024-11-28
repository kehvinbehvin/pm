package cobra

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"

	"github/pm/pkg/common"
	"github/pm/pkg/dag"
	"github/pm/pkg/fileManager"
	"github/pm/pkg/trie"
	"github/pm/resolver"
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
			epicTrie := trie.NewReconcilableTrie("epic")
			epicTrie.SaveReconcilable("./.pm/trie/epic")

			defer epicFile.Close()

			storyFile, storyErr := os.Create("./.pm/trie/story")
			if storyErr != nil {
				fmt.Printf("Error creating story trie: %v\n", storyErr)
			}
			storyTrie := trie.NewReconcilableTrie("story")
			storyTrie.SaveReconcilable("./pm/trie/story")

			defer storyFile.Close()

			taskFile, taskErr := os.Create("./.pm/trie/task")
			if taskErr != nil {
				fmt.Printf("Error creating task trie: %v\n", taskErr)
			}
			taskTrie := trie.NewReconcilableTrie("task")
			taskTrie.SaveReconcilable("./.pm/trie/task")

			defer taskFile.Close()

			err = os.Mkdir("./.pm/dag", os.ModePerm)
			if err != nil && !os.IsExist(err) {
				fmt.Printf("Error creating dag directory: %v\n", err)
			}
			pmDag := dag.NewReconcilableDag("./pm/dag/pmDag")
			pmDag.SaveReconcilable("./.pm/dag/pmDag")

			tmpFile, tmpErr := os.Create("./.pm/tmp")
			if tmpErr != nil {
				fmt.Printf("Error creating tmp file: %v\n", tmpErr)
			}

			defer tmpFile.Close()

			err = os.Mkdir("./.pm/remote", os.ModePerm)
			if err != nil && !os.IsExist(err) {
				fmt.Printf("Error creating remote directory: %v\n", err)
			}

			remoteDeltaFile, remoteDeltaErr := os.Create("./.pm/remote/dag")
			if remoteDeltaErr != nil {
				fmt.Printf("Error creating delta file: %v\n", remoteDeltaErr)
			}
			remoteDeltaTree := dag.NewReconcilableDag("./.pm/remote/dag")
			remoteDeltaTree.SaveReconcilable("./.pm/remote/dag")

			defer remoteDeltaFile.Close()
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete epics, stories or tasks",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			pmDag := dag.LoadDag("pmDag")
			defer pmDag.SaveDag()

			if epics > 0 {
				for _, value := range eValues {
					// Create a Delete Vertex Alpha
					// Update datastructure
				}

			}

			if stories > 0 {
				for _, value := range sValues {
					// Create a Delete Vertex Alpha
					// Update datastructure
				}

			}

			if tasks > 0 {
				for _, value := range tValues {
					// Create a Delete Vertex Alpha
					// Update datastructure
				}

			}

			epicTrie := trie.Load("epic")
			defer epicTrie.SaveReconcilable("./.pm/trie/epic")
			for _, value := range eValues {
				// Create a Delete Word Alpha
				// Update dataStructure
			}

			storyTrie := trie.Load("story")
			defer storyTrie.SaveReconcilable("./.pm/trie/story")
			for _, value := range sValues {
				// Create a Delete Word Alpha
				// Update datastructure
			}

			taskTrie := trie.Load("task")
			defer taskTrie.SaveReconcilable("./.pm/trie/task")
			for _, value := range tValues {
				// Create a Delete Word Alpha
				// Update datastructure
			}

		},
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add epics, stories or tasks",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			pmDag := dag.LoadDag("pmDag")
			defer pmDag.SaveDag()

			var nodesToSave []*dag.Vertex

			if epics > 0 {
				for _, value := range eValues {
					epic := dag.NewVertex(value)
					// Create Add vertex event
					nodesToSave = append(nodesToSave, epic)
				}

			}

			if stories > 0 {
				for _, value := range sValues {
					story := dag.NewVertex(value)
					// Create Add vertex event
					nodesToSave = append(nodesToSave, story)
				}

			}

			if tasks > 0 {
				for _, value := range tValues {
					task := dag.NewVertex(value)
					// Create Add vertex event
					nodesToSave = append(nodesToSave, task)
				}

			}

			for _, vs := range nodesToSave {
				addErr := pmDag.AddVertex(vs)
				if addErr != nil {
					fmt.Println(addErr.Error())
				} else {
					deltaTree.AddVertexEvent(vs)
				}
			}

			epicTrie := trie.Load("epic")
			for _, value := range eValues {
				// Create File
			}

			storyTrie := trie.Load("story")
			for _, value := range sValues {
				// Create File
			}

			taskTrie := trie.Load("task")
			for _, value := range tValues {
				// Create File
			}
		},
	}

	var linkCmd = &cobra.Command{
		Use:   "link",
		Short: "Link epics, stories and tasks",
		Run: func(cmd *cobra.Command, args []string) {
			epics := len(eValues)
			stories := len(sValues)
			tasks := len(tValues)

			pmDag := dag.LoadDag("pmDag")
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
					epicVertex := pmDag.RetrieveVertex(epicValue)

					storyValue := sValues[0]
					storyVertex := pmDag.RetrieveVertex(storyValue)

					edgeErr := pmDag.AddEdge(epicVertex, storyVertex)
					if edgeErr != nil {
						fmt.Println(edgeErr.Error())
					} else {
						deltaTree.AddEdgeEvent(epicVertex, storyVertex)
					}

					for _, value := range tValues {
						taskVertex := pmDag.RetrieveVertex(value)
						edgeErr := pmDag.AddEdge(storyVertex, taskVertex)

						if edgeErr != nil {
							fmt.Println(edgeErr.Error())
						} else {
							deltaTree.AddEdgeEvent(storyVertex, taskVertex)
						}
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
					epicVertex := pmDag.RetrieveVertex(epicValue)

					for _, value := range sValues {
						storyVertex := pmDag.RetrieveVertex(value)
						edgeErr := pmDag.AddEdge(epicVertex, storyVertex)

						if edgeErr != nil {
							fmt.Println(edgeErr.Error())
						} else {
							deltaTree.AddEdgeEvent(epicVertex, storyVertex)
						}
					}
				}
			} else {
				fmt.Println("Cannot create dependencies between epics")
				return
			}

		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all nodes of a type",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please choose what to display")
		},
	}

	var listEpicsCmd = &cobra.Command{
		Use:   "epic",
		Short: "List all epics",
		Run: func(cmd *cobra.Command, args []string) {
			epicTrie := trie.Load("epic").DataStructure.(*trie.Trie)
			allEpics, err := epicTrie.LoadAllWords()
			if err != nil {
				return
			}

			buildList(allEpics, "Epics")
		},
	}

	var listStoriesCmd = &cobra.Command{
		Use:   "story",
		Short: "List all story",
		Run: func(cmd *cobra.Command, args []string) {
			storyTrie := trie.Load("story").DataStructure.(*trie.Trie)
			allEpics, err := storyTrie.LoadAllWords()
			if err != nil {
				return
			}

			buildList(allEpics, "Story")
		},
	}

	var listTasksCmd = &cobra.Command{
		Use:   "task",
		Short: "List all task",
		Run: func(cmd *cobra.Command, args []string) {
			taskTrie := trie.Load("task").DataStructure.(*trie.Trie)
			allTasks, err := taskTrie.LoadAllWords()
			if err != nil {
				return
			}

			buildList(allTasks, "Tasks")
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
			pmDag := dag.LoadDag("./pm/pmDag")

			total := epics + stories + tasks
			if total > 1 {
				fmt.Println("Only allow to list 1 of a kind at a time")
				return
			}

			var node *dag.Vertex
			var nodeType string
			var header string
			var description []byte
			var err error
			var childNodeType string

			if epics > 0 {
				// Simulate getting epic, stories, and tasks data
				node = pmDag.RetrieveVertex(eValues[0])
				header = eValues[0]
				nodeType = "Epic"
				childNodeType = "Stories"

				epicTrie := trie.Load("epic").DataStructure.(*trie.Trie)
				description, err = blob.RetrieveContent(header, epicTrie)
				if err != nil {
					return
				}

			} else if stories > 0 {
				node = pmDag.RetrieveVertex(sValues[0])
				header = sValues[0]
				nodeType = "Story"
				childNodeType = "Tasks"

				storyTrie := trie.Load("story").DataStructure.(*trie.Trie)
				description, err = blob.RetrieveContent(header, storyTrie)
				if err != nil {
					return
				}
			} else if tasks > 0 {
				node = pmDag.RetrieveVertex(tValues[0])
				header = tValues[0]
				nodeType = "Task"

				taskTrie := trie.Load("task").DataStructure.(*trie.Trie)
				description, err = blob.RetrieveContent(header, taskTrie)
				if err != nil {
					return
				}
			} else {
				fmt.Println("Not sure what you want to display")
				return
			}

			nodeKeys := make([]string, len(node.Children))
			i := 0
			for k := range node.Children {
				nodeKeys[i] = k
				i++
			}

			buildSection(description, nodeType+": "+header)
			buildList(nodeKeys, "Related "+childNodeType)

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
				epicTrie := trie.Load("epic")
				defer epicTrie.SaveReconcilable("./.pm/trie/epic")

				updateErr := fileManager.UpdateBlobContent(eValues[0], string(content), epicTrie)
				if updateErr != nil {
					fmt.Println("Update file error")
					return
				}
			} else if stories > 0 {
				storyTrie := trie.Load("story")
				defer storyTrie.SaveReconcilable("./.pm/trie/story")

				updateErr := fileManager.UpdateBlobContent(sValues[0], string(content), storyTrie)
				if updateErr != nil {
					fmt.Println("Update file error")
					return
				}
			} else if tasks > 0 {
				taskTrie := trie.Load("task")
				defer taskTrie.SaveReconcilable("./.pm/trie/task")

				updateErr := fileManager.UpdateBlobContent(tValues[0], string(content), taskTrie)
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

	var attachCmd = &cobra.Command{
		Use:   "attach",
		Short: "Create symlink to master directory",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				info, err := os.Stat("./.pm")
				if !os.IsNotExist(err) {
					if info.IsDir() {
						fmt.Printf("Directory already is managed by pm")
						return
					}
				}
				sourceDir := args[0]
				fmt.Println(sourceDir)
				err = os.Symlink(sourceDir, "./.pm")
				if err != nil {
					fmt.Println("Error creating symlink:", err)
					return
				}
			}
		},
	}
	var detachCmd = &cobra.Command{
		Use:   "detach",
		Short: "Remove symlink to master directory",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat("./.pm")
			if os.IsNotExist(err) {
				fmt.Println("pm does not exist")
				return
			}
			err = os.Remove("./.pm")
			if err != nil {
				fmt.Println("Error removing symlink:", err)
				return
			}
		},
	}

	var epicCmd = &cobra.Command{
		Use:   "epics",
		Short: "Epic suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			epicTrie := trie.Load("epic").DataStructure.(*trie.Trie)
			if epicTrie != nil {
				suggestions := epicTrie.LoadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	var storyCmd = &cobra.Command{
		Use:   "story",
		Short: "Story suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			storyTrie := trie.Load("story").DataStructure.(*trie.Trie)
			if storyTrie != nil {
				suggestions := storyTrie.LoadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	var taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Task suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			toComplete := args[0]
			taskTrie := trie.Load("task").DataStructure.(*trie.Trie)
			if taskTrie != nil {
				suggestions := taskTrie.LoadWordsFromPrefix(toComplete)
				fmt.Println(strings.Join(suggestions, "\n"))
			}
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull delta",
		Run: func(cmd *cobra.Command, args []string) {
			localTree := dag.LoadDelta("./.pm/delta")
			defer localTree.SaveDelta("./.pm/delta")
			if localTree == nil {
				fmt.Printf("Local tree is empty")
				return
			}

			remoteTree := dag.LoadRemoteDelta()
			defer remoteTree.SaveDelta("./.pm/remote/delta")
			if remoteTree == nil {
				fmt.Printf("Remote tree is empty")
				return
			}

			retrieveFile(localTree, remoteTree)
		},
	}

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push delta",
		Run: func(cmd *cobra.Command, args []string) {
			localTree := dag.LoadDelta("./.pm/delta")
			defer localTree.SaveDelta("./.pm/delta")
			if localTree == nil {
				fmt.Printf("Local tree is empty")
				return
			}

			remoteTree := dag.LoadRemoteDelta()
			defer remoteTree.SaveDelta("./.pm/remote/delta")
			if remoteTree == nil {
				fmt.Printf("Remote tree is empty")
				return
			}

			pushDeltas(localTree, remoteTree)
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Test",
		Run: func(cmd *cobra.Command, args []string) {
			deltaTree := dag.LoadDelta("./.pm/delta")
			defer deltaTree.SaveDelta("./.pm/delta")
			if deltaTree == nil || len(deltaTree.Seq) == 0 {
				fmt.Printf("Local tree is empty")
			} else {
				for _, value := range deltaTree.Seq {
					delta := *value
					fmt.Println(delta)
				}

			}

			fmt.Println("DONE")

			remoteTree := dag.LoadRemoteDelta()
			defer remoteTree.SaveDelta("./.pm/remote/delta")
			if remoteTree == nil || len(remoteTree.Seq) == 0 {
				fmt.Printf("Remote tree is empty")
			} else {
				for _, value := range remoteTree.Seq {
					delta := *value
					fmt.Println(delta)
				}

			}

			fmt.Println("DONE")
			pmDag := dag.LoadDag("pmDag")
			defer pmDag.SaveDag()
			mergingErr := resolver.MergeTrees(deltaTree, remoteTree, pmDag)
			if mergingErr != nil {
				fmt.Println("Error mergingTrees", mergingErr.Error())
			}
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(attachCmd)
	rootCmd.AddCommand(detachCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(epicCmd)
	rootCmd.AddCommand(storyCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)

	rootCmd.AddCommand(testCmd)

	listCmd.AddCommand(listEpicsCmd)
	listCmd.AddCommand(listStoriesCmd)
	listCmd.AddCommand(listTasksCmd)

	// Add flags to the add command for epic, task, and story
	linkCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	linkCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add a story")
	linkCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add a task")

	addCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	addCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add a story")
	addCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add a task")

	deleteCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Delete an epic")
	deleteCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Delete a story")
	deleteCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Delete a task")

	viewCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	viewCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add an story")
	viewCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add a task")

	editCmd.Flags().StringSliceVarP(&eValues, "epic", "e", []string{}, "Add an epic")
	editCmd.Flags().StringSliceVarP(&sValues, "story", "s", []string{}, "Add an story")
	editCmd.Flags().StringSliceVarP(&tValues, "task", "t", []string{}, "Add a task")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildSection(description []byte, listHeader string) {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers(listHeader)

	t.Row(string(description))
	fmt.Println(t)
}

func buildList(rows []string, listHeader string) {
	rowIds := [][]string{}
	for _, value := range rows {
		row := []string{value}
		rowIds = append(rowIds, row)
	}

	if len(rowIds) == 0 {
		return
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
