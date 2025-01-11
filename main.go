package main

import (
	"github/pm/internals/ui/application"
	"github/pm/pkg/fileSystem"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	// 	log.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
	//
	// // Access the final state of the model
	// if sm, ok := finalModel.(ui.ApplicationModel); ok {
	// 	epicChoice := sm.Epic
	// 	storyChoice := sm.Story
	// 	taskChoice := sm.Task
	// 	log.Println("Selected Epic:", epicChoice)
	// 	log.Println("Selected Story:", storyChoice)
	// 	log.Println("Selected Task:", taskChoice)
	// } else {
	// 	log.Println("Unexpected model type")
	// }
	setupLogger()
	fs := fileSystem.NewFileSystem()
	fs.Boot()

	app, appErr := application.NewApplication(fs)
	if appErr != nil {
		os.Exit(1)
	}

	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Println("Error running program:", err)
		fs.ShutDown()
		os.Exit(1)
	}

	fs.ShutDown()
}
