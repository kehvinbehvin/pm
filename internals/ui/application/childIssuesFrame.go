package application

import (
	"errors"
	"github/pm/pkg/fileSystem"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChildIssueFrame struct {
	children list.Model
	fileType string
}

func NewChildIssueFrame(app Application, fileName string, childRelationship string) (ApplicationFrame, error) {
	var issueItems []list.Item
	issues, err := app.Fs.ListRelatedIssues(fileName, childRelationship)
	if err != nil {
		app.History.Pop();
	}

	for _, issue := range issues {
		issueItems = append(issueItems, item(issue))
	}

	var pageTitle string;
	switch(childRelationship) {
	case fileSystem.FILE_RELATIONSHIP_DEPENDENCY:
		pageTitle = "Blocking"
	case fileSystem.FILE_RELATIONSHIPS_HIERARCHY:
		pageTitle = "Child"
	}

	fileType, typeErr := app.Fs.GetFileType(fileName);
	if typeErr != nil {
		return ChildIssueFrame{}, typeErr
	}

	const defaultWidth = 50
	l := list.New(issueItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "[" + fileName + "] " + pageTitle + " issues"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)
	
	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(issueItems) + 4, maxHeight))

	bf := ChildIssueFrame{
		children: l,
		fileType: fileType,
	}

	return &bf, nil

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
		case "c":
			selectedItem := browseFrame.children.SelectedItem().(item) 
			issueId := string(selectedItem)
			childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIPS_HIERARCHY)
			if issueErr != nil {
				return app, tea.Quit
			}

			app.History.Push(childFrame)
		case "d":
			selectedItem := browseFrame.children.SelectedItem().(item) 
			issueId := string(selectedItem)
			childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY)
			if issueErr != nil {
				return app, tea.Quit
			}

			app.History.Push(childFrame)
		case "o":
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
	var cmd tea.Cmd
	browseFrame.children, cmd = browseFrame.children.Update(msg)
	return app, cmd
}

func (cif ChildIssueFrame) View(app Application) string {
	helptext := "[v] View File ● [c] list Children ● [d] List Downstream dependencies ● [u] List Upstream depedencies\n[q] Quit ● [←] Back \n[e] All epics [s] All stories [t] All tasks"
	marginStyle := lipgloss.NewStyle().Margin(1, 2)

	return cif.children.View() + marginStyle.Render(helptext)
}

func (cif ChildIssueFrame) Init(app Application) tea.Cmd {
	return nil
}
