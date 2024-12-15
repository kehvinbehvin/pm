package application

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"fmt"
	"errors"
	"io"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
type CreateFormFrame struct {
	step int
	title     textinput.Model
	list      list.Model
	fileType  string
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

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

func NewCreateFormFrame() (ApplicationFrame) {
	ti := textinput.New()
	ti.Placeholder = "Title"
	ti.Focus()

	items := []list.Item{
		item("PRD"),
		item("Epic"),
		item("Story"),
		item("Task"),
	}

	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, 14)
	l.Title = "What type of file do you want to create"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	m := CreateFormFrame{
		step: 0,
		title: ti,
		list: l,
		fileType: "",
	}

	return &m
}

func (cf CreateFormFrame) getFrame(app Application) (*CreateFormFrame, error) {
	frame, error := app.History.Peek();
	if error != nil {
		return &CreateFormFrame{}, errors.New("Cannot get self")
	}

	createFormFrame := frame.(*CreateFormFrame)
	return createFormFrame, nil
}

func (cf CreateFormFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	createFormFrame, frameErr := cf.getFrame(app);
	if frameErr != nil {
		return app, tea.Quit
	}

	switch (createFormFrame.step) {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return app, tea.Quit
			case "left":
				app.History.Pop()
			case "enter":
				createFormFrame.step++
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
				      return app, tea.Quit
			       }
			    }

		    var cmd tea.Cmd
		    createFormFrame.list, cmd = createFormFrame.list.Update(msg)
		    return app, cmd
	}

	return app, nil 
}

func (cf CreateFormFrame) View(app Application) (string) {
	createFormFrame, frameErr := cf.getFrame(app);
	if frameErr != nil {
		return ""
	}

	if (createFormFrame.step == 0) {
		return fmt.Sprintf(
		"Enter the title of your file\n\n%s\n\n%s",
		createFormFrame.title.View(),
		"[esc] Quit [left arrow] Back",
		)
	} else {
		return createFormFrame.list.View()
	}
}

func (cf CreateFormFrame) Init() (tea.Cmd) {
	return nil
}
