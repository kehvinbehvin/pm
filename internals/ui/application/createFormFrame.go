package application

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().Margin(1, 0, 0, 0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type CreateFormFrame struct {
	step     int
	title    textinput.Model
	list     list.Model
	fileType string
	postActionList list.Model
	postAction string
	parent string
	error bool
	errorMessage string
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{}
}
func (d itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewCreateFormFrame(parent string) ApplicationFrame {
	ti := textinput.New()
	ti.Placeholder = "Title"
	ti.Focus()

	items := []list.Item{
		item("prd"),
		item("epic"),
		item("story"),
		item("task"),
	}

	actionItems := []list.Item{
		item("Edit file"),
		item("Create file"),
		item("Create child issue"),
		item("Menu"),
	}

	const defaultWidth = 100
	delegate := itemDelegate{}
	al := list.New(actionItems, delegate, defaultWidth, 14)
	al.Title = "What do you want to do next"
	al.SetShowStatusBar(false)
	al.SetFilteringEnabled(false)
	al.Styles.Title = titleStyle
	al.Styles.PaginationStyle = paginationStyle
	al.SetShowHelp(false)

	l := list.New(items, itemDelegate{}, defaultWidth, 14)
	l.Title = "What type of file do you want to create"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)
	m := CreateFormFrame{
		step:     0,
		title:    ti,
		list:     l,
		fileType: "",
		postActionList: al,
		parent: parent,
	}

	return &m
}

func (cf CreateFormFrame) getFrame(app Application) (*CreateFormFrame, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &CreateFormFrame{}, errors.New("Cannot get self")
	}

	createFormFrame := frame.(*CreateFormFrame)
	return createFormFrame, nil
}

func (cf CreateFormFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	createFormFrame, frameErr := cf.getFrame(app)
	if frameErr != nil {
		return app, tea.Quit
	}

	switch createFormFrame.step {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return app, tea.Quit
			case "left":
				app.History.Pop()
			case "enter":
				i := createFormFrame.title.Value()
				leftBracket := strings.Index(string(i), "[")
				rightBracket := strings.Index(string(i), "]")
				if leftBracket != -1 || rightBracket != -1 {
					createFormFrame.error = true
					createFormFrame.errorMessage = "Cannot use reserved characters : []"
					return app, nil
				}

				createFormFrame.step++
				createFormFrame.error = false
				createFormFrame.errorMessage = ""
			}
		}

		var cmd tea.Cmd
		createFormFrame.title, cmd = createFormFrame.title.Update(msg)

		return app, cmd
	case 1:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			createFormFrame.list.SetWidth(msg.Width)
			return app, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return app, tea.Quit

			case "enter":
				i, ok := createFormFrame.list.SelectedItem().(item)	
				if ok {
					createFormFrame.fileType = string(i)
				}

				app.Fs.CreateFile(createFormFrame.title.Value(), createFormFrame.fileType)
				app.Fs.LinkHierarchy(createFormFrame.parent, createFormFrame.title.Value())
				// Create file here
				createFormFrame.step = 2
			}
		}

		var cmd tea.Cmd
		createFormFrame.list, cmd = createFormFrame.list.Update(msg)
		return app, cmd
	case 2:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			createFormFrame.list.SetWidth(msg.Width)
			return app, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return app, tea.Quit

			case "enter":
				i, ok := createFormFrame.postActionList.SelectedItem().(item)	
				if ok {
					createFormFrame.postAction = string(i)
					switch (createFormFrame.postAction) {
					case "Edit file":
						app.History.Pop()
						frame := EditFileFrame{}	
						app.History.Push(frame)
					case "Create file":
						app.History.Pop()
						frame := NewCreateFormFrame("")
						app.History.Push(frame)
					case "Create child issue":
						app.History.Pop()
						frame := NewCreateFormFrame(createFormFrame.title.Value())
						app.History.Push(frame)
					case "Menu":
						app.History.Pop()
						frame := WelcomeFrame{}
						app.History.Push(frame);
					}
				}
			}
		}

		var cmd tea.Cmd
		createFormFrame.postActionList, cmd = createFormFrame.postActionList.Update(msg)
		return app, cmd

	}

	return app, nil
}

func (cf CreateFormFrame) View(app Application) string {
	createFormFrame, frameErr := cf.getFrame(app)
	if frameErr != nil {
		return ""
	}

	if createFormFrame.step == 0 {
		title := "Issue Title"
		helptext := "\n[←] Back"
		marginStyle := lipgloss.NewStyle().Margin(1, 2)
		
		if createFormFrame.error && createFormFrame.errorMessage != "" {
			errorText := "\n " + createFormFrame.errorMessage
			return marginStyle.Render(title) + "\n" + createFormFrame.title.View() + marginStyle.Render(errorText) + marginStyle.Render(helptext)
		}

		return marginStyle.Render(title) + "\n" + createFormFrame.title.View() + marginStyle.Render(helptext)
	} else if (createFormFrame.step == 1) {
		helptext := "\n[q] Quit ● [enter] Enter"
		marginStyle := lipgloss.NewStyle().Margin(1, 2)

		return createFormFrame.list.View() + marginStyle.Render(helptext)
	} else {
		helptext := "\n[q] Quit ● [enter] Enter"
		marginStyle := lipgloss.NewStyle().Margin(1, 2)

		return createFormFrame.postActionList.View() + marginStyle.Render(helptext)

	}
}

func (cf CreateFormFrame) Init(app Application) tea.Cmd {
	return nil
}
