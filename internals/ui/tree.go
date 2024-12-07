package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/tree"
)

type GraphModel struct {}

func NewGraph() tea.Model {
	return GraphModel{}
}

func (gm GraphModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return gm, nil
}

func (gm GraphModel) View() string {
	t :=  tree.Root(".").Child("A", "B", "C")
	return t.String(); 
}

func (gm GraphModel) Init() tea.Cmd {
	return nil
}
