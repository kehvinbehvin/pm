package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github/pm/internals/ui/application"
	"os"
)

func main() {
	// Create the initial model
	// m := ui.NewModel()
	//
	// // Run the Bubble Tea program
	// p := tea.NewProgram(m, tea.WithAltScreen())
	// finalModel, err := p.Run()
	//
	// if err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
	//
	// // Access the final state of the model
	// if sm, ok := finalModel.(ui.ApplicationModel); ok {
	// 	epicChoice := sm.Epic
	// 	storyChoice := sm.Story
	// 	taskChoice := sm.Task
	// 	fmt.Println("Selected Epic:", epicChoice)
	// 	fmt.Println("Selected Story:", storyChoice)
	// 	fmt.Println("Selected Task:", taskChoice)
	// } else {
	// 	fmt.Println("Unexpected model type")
	// }
	app := application.NewApplication()

	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
