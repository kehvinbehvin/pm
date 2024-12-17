package application

import (
	tea "github.com/charmbracelet/bubbletea"
)

type WelcomeFrame struct{}

func (wf WelcomeFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "c":
			frame := NewCreateFormFrame()
			app.History.Push(frame)
		}
	}

	return app, nil
}

func (wf WelcomeFrame) View(app Application) string {
	return "WelcomeFrame, Press c to create file"
}

func (wf WelcomeFrame) Init() tea.Cmd {
	return nil
}
