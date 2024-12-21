package application

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type ChildIssueFrame struct {
	children list.Model
}

func NewChildIssueFrame(app Application, fileName string) ApplicationFrame {
	var issueItems []list.Item
	issues, err := app.Fs.ListChildIssues(fileName)
	if err != nil {
		app.History.Pop();
	}

	for _, issue := range issues {
		issueItems = append(issueItems, item(issue))
	}

	const defaultWidth = 50
	l := list.New(issueItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "What type of file do you want to create"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	bf := ChildIssueFrame{
		children: l,
	}

	return &bf

}

func (cif ChildIssueFrame) getFrame(app Application) (*ChildIssueFrame, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &ChildIssueFrame{}, errors.New("Cannot get self")
	}

	childIssueFrame := frame.(*ChildIssueFrame)
	return childIssueFrame, nil
}

func (cif ChildIssueFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
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

func (cif ChildIssueFrame) View(app Application) string {
	return cif.children.View()
}

func (cif ChildIssueFrame) Init() tea.Cmd {
	return nil
}
