package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"errors"
	"sort"

	"github.com/charmbracelet/bubbles/list"
)

type GlobalSelectionFrame struct{
	items list.Model
	selectedItem string
}

func indexOf(haystack []string, needle string) (int) {
	for index, value := range haystack {
		if needle == value {
			return index
		}
	}

	return -1
}

func NewGlobalSelectionFrame(app Application, currentFile string, excludeFiles []string) (*GlobalSelectionFrame, error) {
	// Fetch all issues
	filesWithTypes, fsErr := app.Fs.ListAllFilesWithTypes()
	if fsErr != nil {
		return &GlobalSelectionFrame{}, fsErr
	}

	// Create list
	fileItemList := make([]list.Item, 0)
	fileList := make([]string, 0)
	for fileType, files := range filesWithTypes {
		for _, fileName := range files {
			if fileName == currentFile {
				continue;
			}

			if indexOf(excludeFiles, fileName) != -1 {
				continue;
			}

			title := "[" + fileType + "] " + fileName
			fileList = append(fileList, title)
		}
	}

	sort.Strings(fileList)
	if len(fileList) == 0 {
		return &GlobalSelectionFrame{}, errors.New("No items to link to")
	}

	for _, file := range fileList {
		fileItemList = append(fileItemList, item(file))
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
		case "left":
			app.History.Pop()
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
	helptext := "[v] View ● [enter] select\n[q] Quit ● [←] Back "
	marginStyle := lipgloss.NewStyle().Margin(1, 2)
	return gsf.items.View() + marginStyle.Render(helptext)
}

func (gsf GlobalSelectionFrame) Init(app Application) tea.Cmd {
	return nil
}
