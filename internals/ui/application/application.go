package application

import (
	tea "github.com/charmbracelet/bubbletea"
)

func NewApplication() (tea.Model) {
	stack := NewApplicationStack()

	// Initialise the first frame of the application
	welcome := &WelcomeFrame{} 
	stack.Push(welcome)

	return Application{
		History: stack,
	}
}

type Application struct {
	History *ApplicationStack 
}

func (a Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentFrame, error := a.History.Peek();
	if error != nil {
		return a, tea.Quit 
	}

	return currentFrame.Update(msg, a)
}

func (a Application) View() (string) {
	currentFrame, error := a.History.Peek();
	if error != nil {
		// Maybe introduce an error page here
		return ""
	}

	return currentFrame.View(a)
}
func (a Application) Init() (tea.Cmd) {
	currentFrame, error := a.History.Peek();
	if error != nil {
		return nil
	}

	return currentFrame.Init()
}
