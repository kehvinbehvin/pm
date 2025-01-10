package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeFrame struct{}

func (wf WelcomeFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "i":
			frame, frameErr := NewCreateFormFrame(app, "")
			if frameErr != nil {
				return app, tea.Quit
			}

			app.History.Push(frame)
		case "e":
			frame := NewBrowseFrame(app, "epic")
			app.History.Push(frame)
		case "s":
			frame := NewBrowseFrame(app, "story")
			app.History.Push(frame)
		case "t":
			frame := NewBrowseFrame(app, "task")
			app.History.Push(frame)
		}
	}

	return app, nil
}

func (wf WelcomeFrame) View(app Application) string {
	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	return marginStyle.Render("Browser\n\n[i] Create issue\n[e] List epics\n[s] List stories\n[t] List tasks\n[q] Quit")
}

func (wf WelcomeFrame) Init(app Application) tea.Cmd {
	return nil
}
