package application

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
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
	l.Title = "What type of file do you want to create"
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
