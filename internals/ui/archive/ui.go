package archive

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	// "os"
	"fmt"
	"strings"
)

type item struct {
	title string
	desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }


var (
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	menuStyle   = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle             = lipgloss.NewStyle()
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var (
	epic = "epic"
	story = "story"
	task = "task"
)

type menuOption int

const (
	Menu menuOption = iota
	Browse
	Create
	View
	Edit
)

var menuOptions = []string{
	"Menu",
	"Browse",
	"Create",
	"View",
	"Edit",
}

type ApplicationModel struct {
	epicList list.Model
	storyList list.Model
	taskList list.Model
	quitting bool
	Epic string
	Story string
	Task string
	mode menuOption
	cursor  int
	createInput []textinput.Model
	focusIndex int
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (sm ApplicationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			sm.quitting = true
			return sm, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Recalculate the size for all lists
		h, v := docStyle.GetFrameSize()
		sm.epicList.SetSize(msg.Width-h, msg.Height-v)
		sm.storyList.SetSize(msg.Width-h, msg.Height-v)
		sm.taskList.SetSize(msg.Width-h, msg.Height-v)
	}

	switch sm.mode {
	case Menu:
	return sm.menuUpdate(msg);
	case Browse:
	return sm.browseUpdate(msg);
	case Create:
	return sm.createUpdate(msg);
	}
	return sm, tea.Quit
}

func (m ApplicationModel) createUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.createInput) {
				newItem := item{
					title: m.createInput[0].Value(),
					desc:  m.createInput[1].Value(),
				}

				// Add the new item to the epicList
				m.epicList.InsertItem(len(m.epicList.Items()), newItem)

				// Reset inputs
				for i := range m.createInput {
					m.createInput[i].Reset()
				}

				m.mode = Menu
				m.focusIndex = 0;
				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.createInput) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.createInput)
			}

			cmds := make([]tea.Cmd, len(m.createInput))
			for i := 0; i <= len(m.createInput)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.createInput[i].Focus()
					m.createInput[i].PromptStyle = focusedStyle
					m.createInput[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.createInput[i].Blur()
				m.createInput[i].PromptStyle = noStyle
				m.createInput[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m ApplicationModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.createInput))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.createInput {
		m.createInput[i], cmds[i] = m.createInput[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ApplicationModel) menuUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuOptions)-1 {
				m.cursor++
			}
		case "enter":
			m.mode = menuOption(m.cursor)
		}
	}
	return m, nil
}

func (sm ApplicationModel) browseUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sm.Epic == "" {
		return sm.updateEpicList(msg)
	} else if sm.Story == "" {
		return sm.updateStoryList(msg)
	} else if sm.Task == "" {
		return sm.updateTaskList(msg)
	}

	var cmd tea.Cmd
	return sm, cmd
}

func (sm ApplicationModel) updateEpicList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	switch keyPress := msg.String(); keyPress {
		case "enter":
		if sm.epicList.FilterState() == list.Filtering {				// Allow the list to handle the filtering first
			break
		}

		i, ok := sm.epicList.SelectedItem().(item)
			if ok {
				sm.Epic = i.title;
			}
		}
	}

	var cmd tea.Cmd
	sm.epicList, cmd = sm.epicList.Update(msg)
	return sm, cmd
} 

func (sm ApplicationModel) updateStoryList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	switch keyPress := msg.String(); keyPress {
		case "enter":
		if sm.storyList.FilterState() == list.Filtering {				// Allow the list to handle the filtering first
			break
		}

		i, ok := sm.storyList.SelectedItem().(item)
			if ok {
				sm.Story = i.title;
			}
		}
	}

	var cmd tea.Cmd
	sm.storyList, cmd = sm.storyList.Update(msg)
	return sm, cmd
} 

func (sm ApplicationModel) updateTaskList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	switch keyPress := msg.String(); keyPress {
		case "enter":
		if sm.taskList.FilterState() == list.Filtering {				// Allow the list to handle the filtering first
			break
		}
		i, ok := sm.taskList.SelectedItem().(item)
			if ok {
				sm.Task = i.title;
				return sm, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	sm.taskList, cmd = sm.taskList.Update(msg)
	return sm, cmd
}

func (m ApplicationModel) View() (string) {
	switch m.mode {
	case Menu:
	return m.MenuView()
	case Browse:
	return m.BrowseView()
	case Create:
	return m.CreateView()
	}

	return docStyle.Render()
}

func (m ApplicationModel) CreateView() (string) {
	var b strings.Builder

	for i := range m.createInput {
		b.WriteString(m.createInput[i].View())
		if i < len(m.createInput)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.createInput) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func (m ApplicationModel) MenuView() (string) {
	var menu string
	for i, option := range menuOptions {
		cursor := " " // No cursor by default
		style := unselectedStyle
		if m.cursor == i {
			cursor = ">" // Cursor for the selected item
			style = cursorStyle
		}
		if menuOption(i) == m.mode {
			style = selectedStyle
		}
		menu += fmt.Sprintf("%s %s\n", cursor, style.Render(option))
	}

	// Render output
	return fmt.Sprintf("%s\n%s\n%s", "Select an option:", menu)
}

func (m ApplicationModel) BrowseView() (string) {
	if m.Epic == "" {
		return docStyle.Render(m.epicList.View())
	} else if m.Story == "" {
		return docStyle.Render(m.storyList.View())
	} else {
		return docStyle.Render(m.taskList.View())
	}
}

func (sm ApplicationModel) Init() tea.Cmd {
	return textinput.Blink
}

func NewModel() (tea.Model) {
	epics := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
	}

	stories := []list.Item{
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
	}

	tasks := []list.Item{
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	epicList :=  list.New(epics, list.NewDefaultDelegate(), 0, 0);
	storyList := list.New(stories, list.NewDefaultDelegate(), 0, 0)
	taskList := list.New(tasks, list.NewDefaultDelegate(), 0, 0)

	titleInput := textinput.New();
	titleInput.Placeholder = "Title here"
	titleInput.Focus()
	titleInput.CharLimit = 156
	titleInput.Width = 20

	descInput := textinput.New();
	descInput.Placeholder = "Description here"
	descInput.CharLimit = 156
	descInput.Width = 20

	var inputs []textinput.Model
	inputs = append(inputs, titleInput) 
	inputs = append(inputs, descInput) 
	
	m := ApplicationModel{epicList: epicList, storyList: storyList, taskList: taskList , mode: Menu, createInput: inputs}
	m.epicList.Title = "My Fave epics"
	m.storyList.Title = "My Fave stories"
	m.taskList.Title = "My Fave tasks"
	return m
}
