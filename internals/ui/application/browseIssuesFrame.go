package application

import (
	"errors"
	"github/pm/pkg/fileSystem"
	"log"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BrowseFrame struct {
	epics    list.Model
	fileType string
}

func NewBrowseFrame(app Application, fileType string) ApplicationFrame {
	var epicItems []list.Item
	epics, err := app.Fs.ListFileNamesByType(fileType)
	if err != nil {
		return WelcomeFrame{}
	}

	for _, epic := range epics {
		epicItems = append(epicItems, item(epic))
	}

	delegate := itemDelegate{}

	const defaultWidth = 50
	l := list.New(epicItems, delegate, defaultWidth, 14)
	caser := cases.Title(language.English)
	l.Title = caser.String(fileType) + " listing"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)
	l.Styles.NoItems = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("240"))

	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(epicItems)+4, maxHeight))
	bf := BrowseFrame{
		epics:    l,
		fileType: fileType,
	}

	return &bf
}

func (bf BrowseFrame) Refresh(app Application) error {
	browseFrame, frameErr := bf.getFrame(app)
	if frameErr != nil {
		return frameErr
	}

	var epicItems []list.Item
	epics, err := app.Fs.ListFileNamesByType(bf.fileType)
	if err != nil {
		return err
	}

	for _, epic := range epics {
		epicItems = append(epicItems, item(epic))
	}

	browseFrame.epics.SetItems(epicItems)
	return nil
}

func (bg BrowseFrame) getFrame(app Application) (*BrowseFrame, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &BrowseFrame{}, errors.New("Cannot get self")
	}

	browseFrame := frame.(*BrowseFrame)
	return browseFrame, nil
}

func (bf BrowseFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	browseFrame, frameErr := bf.getFrame(app)
	if frameErr != nil {
		return app, tea.Quit
	}

	if len(browseFrame.epics.Items()) == 0 {
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
			case "d":
				selectedItem := browseFrame.epics.SelectedItem().(item)
				issueId := string(selectedItem)
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY, true)
				if issueErr != nil {
					return app, tea.Quit
				}
				app.History.Push(childFrame)

			case "c":
				selectedItem := browseFrame.epics.SelectedItem().(item)
				issueId := string(selectedItem)
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIPS_HIERARCHY, true)
				if issueErr != nil {
					return app, tea.Quit
				}

				app.History.Push(childFrame)
			case "u":
				selectedItem := browseFrame.epics.SelectedItem().(item)
				issueId := string(selectedItem)
				childFrame, issueErr := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY, false)
				if issueErr != nil {
					return app, tea.Quit
				}

				app.History.Push(childFrame)
			case "v":
				selectedItem := browseFrame.epics.SelectedItem().(item)
				issueId := string(selectedItem)
				content, contentErr := app.Fs.RetrieveFileContents(issueId)
				log.Println("Loaded content")
				if contentErr != nil {
					return nil, tea.Quit
				}

				mdFrame, frameErr := NewViewMarkdownFrame(issueId, content, app)
				log.Println("Created MD Frame")
				if frameErr != nil {
					return nil, tea.Quit
				}

				app.History.Push(mdFrame)
			case "e":
				if browseFrame.fileType != "epic" {
					frame := NewBrowseFrame(app, "epic")
					app.History.Push(frame)
				}

			case "s":
				if browseFrame.fileType != "story" {
					frame := NewBrowseFrame(app, "story")
					app.History.Push(frame)
				}

			case "t":
				if browseFrame.fileType != "task" {
					frame := NewBrowseFrame(app, "task")
					app.History.Push(frame)
				}
			case "i":
				frame, frameErr := NewCreateFormFrame(app, "")
				if frameErr != nil {
					return app, tea.Quit
				}

				app.History.Push(frame)
			}
		}

	}

	var cmd tea.Cmd
	browseFrame.epics, cmd = browseFrame.epics.Update(msg)
	return app, cmd
}

func (bf BrowseFrame) View(app Application) string {
	browseFrame, frameErr := bf.getFrame(app)
	if frameErr != nil {
		return ""
	}

	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	if len(browseFrame.epics.Items()) == 0 {
		return bf.epics.View() + marginStyle.Render("[q] Quit ● [←] Back")
	}

	helptext := "[v] View File ● [i] Create Issue [c] list Children ● [d] List Downstream dependencies ● [u] List Upstream depedencies\n[q] Quit ● [←] Back \n"
	epicText := "[e] All epics "
	storyText := "[s] All stories "
	taskText := "[t] All tasks "

	if browseFrame.fileType != "epic" {
		helptext = helptext + epicText
	}

	if browseFrame.fileType != "story" {
		helptext = helptext + storyText
	}

	if browseFrame.fileType != "task" {
		helptext = helptext + taskText
	}

	return bf.epics.View() + marginStyle.Render(helptext)
}

func (bf BrowseFrame) Init(app Application) tea.Cmd {
	return nil
}
