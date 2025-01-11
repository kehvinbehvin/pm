package application

import (
	tea "github.com/charmbracelet/bubbletea"
)

type EditFileFrame struct{}

func (ed EditFileFrame) Refresh(app Application) error {
	return nil
}

func (ed EditFileFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "left":
			app.History.Pop()
		}
	}

	return app, nil
}

func (ed EditFileFrame) View(app Application) string {
	return "Edit fileFrame\n [<-] Back"
}

func (ed EditFileFrame) Init(app Application) tea.Cmd {
	return nil
}
