package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DependencyGraph struct{}

func NewDependencyGraph() (*DependencyGraph, error) {
	return &DependencyGraph{
	}, nil
}

func (dg DependencyGraph) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
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

func (dg DependencyGraph) View(app Application) string {
	helptext := "[q] Quit ● [←] Back "
	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	return marginStyle.Render(helptext)
}

func (dg DependencyGraph) Init(app Application) tea.Cmd {
	return nil
}
