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
		case "p":
			frame := NewBrowseFrame(app, "prd")
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
	return "[q] Quit; [c] Create Issue;\n[p] prd; [e] epic; [s] story; [t] task"
}

func (wf WelcomeFrame) Init() tea.Cmd {
	return nil
}
