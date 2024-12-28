package application

import (
	"errors"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

	const defaultWidth = 50
	l := list.New(epicItems, itemDelegate{}, defaultWidth, 14)
	caser := cases.Title(language.English)
	l.Title = caser.String(fileType) + " listing"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

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
		case "enter":
			selectedItem := browseFrame.epics.SelectedItem().(item) 
			issueId := string(selectedItem)
			childFrame := NewChildIssueFrame(app, issueId)
			app.History.Push(childFrame)
		case "e":
			selectedItem := browseFrame.epics.SelectedItem().(item)
			issueId := string(selectedItem)
			app.Fs.EditFile(issueId)
		case "v":
			selectedItem := browseFrame.epics.SelectedItem().(item)
			issueId := string(selectedItem)
			content, contentErr := app.Fs.RetrieveFileContents(issueId)
			if contentErr != nil {
				return nil, tea.Quit
			}

			mdFrame, frameErr := NewViewMarkdownFrame(content)
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
	return bf.epics.View()
}

func (bf BrowseFrame) Init() tea.Cmd {
	return nil
}
