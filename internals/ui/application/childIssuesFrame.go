package application

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

type ChildIssueFrame struct {
	children list.Model
}

func NewChildIssueFrame(app Application, fileName string) ApplicationFrame {
	var issueItems []list.Item
	issues, err := app.Fs.ListRelatedHierarchy(fileName)
	if err != nil {
		app.History.Pop();
	}

	for _, issue := range issues {
		issueItems = append(issueItems, item(issue))
	}

	const defaultWidth = 50
	l := list.New(issueItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "[" + fileName + "] child issues"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)
	
	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(issueItems) + 4, maxHeight))

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
	browseFrame, frameErr := cif.getFrame(app)
	if frameErr != nil {
		return app, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "left":
			app.History.Pop()
		case "enter":
			selectedItem := browseFrame.children.SelectedItem().(item) 
			issueId := string(selectedItem)
			childFrame := NewChildIssueFrame(app, issueId)
			app.History.Push(childFrame)
		case "e":
			selectedItem := browseFrame.children.SelectedItem().(item)
			issueId := string(selectedItem)
			app.Fs.EditFile(issueId)
		case "v":
			selectedItem := browseFrame.children.SelectedItem().(item)
			issueId := string(selectedItem)
			content, contentErr := app.Fs.RetrieveFileContents(issueId)
			if contentErr != nil {
				return nil, tea.Quit
			}

			mdFrame, frameErr := NewViewMarkdownFrame(issueId, content, app)
			if frameErr != nil {
				return nil, tea.Quit
			}

			app.History.Push(mdFrame)
		}
	}

	return app, nil
}

func (cif ChildIssueFrame) View(app Application) string {
	helptext := "[e] Edit ● [v] View ● [enter] Children\n[q] Quit ● [←] Back "
	marginStyle := lipgloss.NewStyle().Margin(1, 2)

	return cif.children.View() + marginStyle.Render(helptext)
}

func (cif ChildIssueFrame) Init() tea.Cmd {
	return nil
}
