package archive

import (
	tea "github.com/charmbracelet/bubbletea"
)

type GraphModel struct {
	altscreen bool
}

func NewGraph() tea.Model {
	return GraphModel{}
}

func (gm GraphModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return gm, tea.Quit
		case " ":
			var cmd tea.Cmd
			if gm.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			gm.altscreen = !gm.altscreen
			return gm, cmd
		}
	}
	return gm, nil

}

func (gm GraphModel) View() string {
	return "* End\n| \n* Epic: UAT Testing\n| \\\n|  |\n*  |  Epic: Create test cases\n|  |\n|  * Epic: Create Playwright stories\n| /\n+\n|\n|\n* Epic: Add user account"
}

func (gm GraphModel) Init() tea.Cmd {
	return nil
}
