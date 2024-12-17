package application

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ApplicationFrame interface {
	Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd)
	View(app Application) string
	Init() tea.Cmd
}
