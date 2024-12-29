package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"errors"

	"github.com/charmbracelet/bubbles/list"
)

type GlobalSelectionFrame struct{
	items list.Model
	selectedItem string
}

func NewGlobalSelectionFrame(app Application) (*GlobalSelectionFrame, error) {
	// Fetch all issues
	filesWithTypes, fsErr := app.Fs.ListAllFilesWithTypes()
	if fsErr != nil {
		return &GlobalSelectionFrame{}, fsErr
	}

	// Create list
	fileItemList := make([]list.Item, 0)
	for fileType, files := range filesWithTypes {
		for _, fileName := range files {
			title := "[" + fileType + "] " + fileName
			fileItemList = append(fileItemList, item(title))
		}
	}
	
	delegate := itemDelegate{}

	const defaultWidth = 50
	l := list.New(fileItemList, delegate, defaultWidth, 14)
	l.Title = "Global listing"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)

	maxHeight := 9 // Maximum height of the list
	l.SetHeight(min(len(fileItemList) + 4, maxHeight))

	return &GlobalSelectionFrame{
		items: l,
	}, nil
}

func (gst GlobalSelectionFrame) getFrame(app Application) (*GlobalSelectionFrame, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &GlobalSelectionFrame{}, errors.New("Cannot get self")
	}

	globalFrame := frame.(*GlobalSelectionFrame)
	return globalFrame, nil
}


func (gsf GlobalSelectionFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	globalFrame, frameErr := gsf.getFrame(app)
	if frameErr != nil {
		return app, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "v":
			// Allow the user to view the markdown for more description
			selectedItem := globalFrame.items.SelectedItem().(item)
			issueId := string(selectedItem)
			content, contentErr := app.Fs.RetrieveFileContents(issueId)
			if contentErr != nil {
				return nil, tea.Quit
			}

			mdFrame, frameErr := NewViewMarkdownFrame(issueId,content, app)
			if frameErr != nil {
				return nil, tea.Quit
			}

			app.History.Push(mdFrame)
		case "enter":
			// Set the item to the struct so that it can be accessed by the previous frame. 
			selectedItem := globalFrame.items.SelectedItem().(item)
			issueId := string(selectedItem)
			globalFrame.selectedItem = issueId
			app.History.Pop()	
		}
	}

	var cmd tea.Cmd
	globalFrame.items, cmd = globalFrame.items.Update(msg)
	return app, cmd
}

func (gsf GlobalSelectionFrame) View(app Application) string {
	helptext := "[v] View ● [enter] select ● [q] Quit"
	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	return gsf.items.View() + marginStyle.Render(helptext)
}

func (gsf GlobalSelectionFrame) Init(app Application) tea.Cmd {
	return nil
}
