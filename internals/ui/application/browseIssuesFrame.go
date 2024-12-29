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

type BrowseFrame struct{
	epics list.Model
}

func NewBrowseFrame(app Application, fileType string) (ApplicationFrame) {
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

	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(epicItems) + 4, maxHeight))
	bf := BrowseFrame{
		epics: l,
	}

	return &bf
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
			childFrame := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIP_DEPENDENCY)
			app.History.Push(childFrame)

		case "c":
			selectedItem := browseFrame.epics.SelectedItem().(item) 
			issueId := string(selectedItem)
			childFrame := NewChildIssueFrame(app, issueId, fileSystem.FILE_RELATIONSHIPS_HIERARCHY)
			app.History.Push(childFrame)
		case "v":
			selectedItem := browseFrame.epics.SelectedItem().(item)
			issueId := string(selectedItem)
			content, contentErr := app.Fs.RetrieveFileContents(issueId)
			log.Println("Loaded content");
			if contentErr != nil {
				return nil, tea.Quit
			}

			mdFrame, frameErr := NewViewMarkdownFrame(issueId,content, app)
			log.Println("Created MD Frame");
			if frameErr != nil {
				return nil, tea.Quit
			}

			app.History.Push(mdFrame)
		}
	}

	var cmd tea.Cmd
	browseFrame.epics, cmd = browseFrame.epics.Update(msg)
	return app, cmd
}

func (bf BrowseFrame) View(app Application) string {	
	helptext := "[v] View ● [c] Child issues ● [d] Dependenciees\n[q] Quit ● [←] Back "
	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	return bf.epics.View() + marginStyle.Render(helptext)
}

func (bf BrowseFrame) Init(app Application) tea.Cmd {
	return nil
}
