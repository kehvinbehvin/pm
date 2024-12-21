package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github/pm/internals/ui/application"
	"os"
	"log"
)

func setupLogger() {
	// Open the file for appending logs, create it if it doesn't exist
	file, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// Set the log output to the file
	log.SetOutput(file)

	// Optional: Customize log prefix and flags
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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
	setupLogger()
	app := application.NewApplication()

	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
