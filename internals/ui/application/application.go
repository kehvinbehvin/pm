package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github/pm/pkg/filesystem"
)

func NewApplication(fs *filesystem.FileSystem) tea.Model {
	stack := NewApplicationStack()
	
	welcome := &WelcomeFrame{}
	stack.Push(welcome)

	return Application{
		History: stack,
		Fs: fs,
	}
}

type Application struct {
	History *ApplicationStack
	Fs *filesystem.FileSystem
}

func (a Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentFrame, error := a.History.Peek()
	if error != nil {
		return a, tea.Quit
	}

	return currentFrame.Update(msg, a)
}

func (a Application) View() string {
	currentFrame, error := a.History.Peek()
	if error != nil {
		// Maybe introduce an error page here
		return ""
	}

	return currentFrame.View(a)
}
func (a Application) Init() tea.Cmd {
	currentFrame, error := a.History.Peek()
	if error != nil {
		return nil
	}

	return currentFrame.Init()
}
