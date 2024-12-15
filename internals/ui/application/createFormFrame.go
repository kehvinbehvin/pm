package application

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CreateFormFrame struct {}

func (cf CreateFormFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "p":
			app.History.Pop()
		}
	}

	return app, nil
}

func (cf CreateFormFrame) View() (string) {
	return "This is the create form frame. Press q to quit and p to go back"
}

func (cf CreateFormFrame) Init() (tea.Cmd) {
	return nil
}
