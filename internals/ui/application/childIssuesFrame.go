package application

import (
	"errors"
	"github/pm/pkg/fileSystem"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChildIssueFrame struct {
	children     list.Model
	fileType     string
	relationship string
	fileName     string
	direction    bool
}

func NewChildIssueFrame(app Application, fileName string, childRelationship string, direction bool) (ApplicationFrame, error) {
	issueItems, asyncErr := asycData(app, fileName, childRelationship, direction)
	if asyncErr != nil {
		return ChildIssueFrame{}, asyncErr
	}

	var pageTitle string
	switch childRelationship {
	case fileSystem.FILE_RELATIONSHIP_DEPENDENCY:
		pageTitle = "Blocking"
	case fileSystem.FILE_RELATIONSHIPS_HIERARCHY:
		pageTitle = "Child"
	}

	fileType, typeErr := app.Fs.GetFileType(fileName)
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
	l.Styles.NoItems = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("240"))
	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(issueItems)+4, maxHeight))

	bf := ChildIssueFrame{
		children:     l,
		fileType:     fileType,
		relationship: childRelationship,
		fileName:     fileName,
		direction:    direction,
	}

	return &bf, nil

}

func asycData(app Application, fileName string, childRelationship string, direction bool) ([]list.Item, error) {
	var issueItems []list.Item
	var issues []string
	var err error

	if direction {
		issues, err = app.Fs.ListRelatedIssues(fileName, childRelationship)
	} else {
		issues, err = app.Fs.ListRelatedParentDependency(fileName)
	}

	if err != nil {
		return issueItems, err
	}

	for _, issue := range issues {
		issueItems = append(issueItems, item(issue))
	}

	return issueItems, nil
}

func (cif ChildIssueFrame) Refresh(app Application) error {
	frame, error := app.History.Peek()
	if error != nil {
		return errors.New("Cannot get self")
	}
	childIssueFrame := frame.(*ChildIssueFrame)
	updatedChildren, asyncErr := asycData(app, childIssueFrame.fileName, childIssueFrame.relationship, childIssueFrame.direction)
	if asyncErr != nil {
		return asyncErr
	}

	childIssueFrame.children.SetItems(updatedChildren)

	return nil
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

	if len(browseFrame.children.Items()) == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return app, tea.Quit
			case "left":
				app.History.Pop()
			}
		}
	} else {
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
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIPS_HIERARCHY, true)
				if issueErr != nil {
					return app, tea.Quit
				}

				app.History.Push(childFrame)
			case "d":
				selectedItem := browseFrame.children.SelectedItem().(item)
				issueId := string(selectedItem)
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY, true)
				if issueErr != nil {
					return app, tea.Quit
				}

				app.History.Push(childFrame)
			case "u":
				selectedItem := browseFrame.children.SelectedItem().(item)
				issueId := string(selectedItem)
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY, false)
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
			case "r":
				selectedItem := browseFrame.children.SelectedItem().(item)
				issueId := string(selectedItem)

				switch browseFrame.relationship {
				case fileSystem.FILE_RELATIONSHIPS_HIERARCHY:
					var parent string
					var child string
					if browseFrame.direction {
						parent = browseFrame.fileName
						child = issueId
					} else {
						parent = issueId
						child = browseFrame.fileName
					}

					app.Fs.UnLinkHierarchy(parent, child)
				case fileSystem.FILE_RELATIONSHIP_DEPENDENCY:
					var parent string
					var child string
					if browseFrame.direction {
						parent = browseFrame.fileName
						child = issueId
					} else {
						parent = issueId
						child = browseFrame.fileName
					}

					app.Fs.UnLinkDependency(parent, child)
				}

				app.History.Pop()
				childFrame, issueErr := NewChildIssueFrame(app, browseFrame.fileName, browseFrame.relationship, browseFrame.direction)
				if issueErr != nil {
					return app, tea.Quit
				}

				app.History.Push(childFrame)
			}
		}

	}

	var cmd tea.Cmd
	browseFrame.children, cmd = browseFrame.children.Update(msg)
	return app, cmd
}

func (cif ChildIssueFrame) View(app Application) string {
	browseFrame, frameErr := cif.getFrame(app)
	if frameErr != nil {
		return ""
	}

	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	if len(browseFrame.children.Items()) == 0 {
		return cif.children.View() + marginStyle.Render("[q] Quit ● [←] Back")
	}

	helptext := "[v] View File ● [c] list Children ● [d] List Downstream dependencies ● [u] List Upstream depedencies [r] Unlink issue\n[q] Quit ● [←] Back \n[e] All epics [s] All stories [t] All tasks"
	return cif.children.View() + marginStyle.Render(helptext)
}

func (cif ChildIssueFrame) Init(app Application) tea.Cmd {
	return nil
}
