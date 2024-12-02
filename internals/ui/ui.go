package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	// "os"
	// "fmt"
)

type item struct {
	title string
	desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type SelectorModel struct {
	epicList list.Model
	storyList list.Model
	taskList list.Model
	quitting bool
	Epic string
	Story string
	Task string
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (sm SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	if sm.Epic == "" {
		return sm.updateEpicList(msg)
	} else if sm.Story == "" {
		return sm.updateStoryList(msg)
	} else if sm.Task == "" {
		return sm.updateTaskList(msg)
	}

	return sm, tea.Quit
}

func (sm SelectorModel) updateEpicList(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (sm SelectorModel) updateStoryList(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (sm SelectorModel) updateTaskList(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (sm SelectorModel) View() (string) {
	if sm.Epic == "" {
		return docStyle.Render(sm.epicList.View())
	} else if sm.Story == "" {
		return docStyle.Render(sm.storyList.View())
	} else if sm.Task == "" {
		return docStyle.Render(sm.taskList.View())
	}

	return docStyle.Render(sm.epicList.View())
}

func (sm SelectorModel) Init() tea.Cmd {
	return nil
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

	m := SelectorModel{epicList: list.New(epics, list.NewDefaultDelegate(), 0, 0), storyList: list.New(stories, list.NewDefaultDelegate(), 0, 0), taskList: list.New(tasks, list.NewDefaultDelegate(), 0, 0)}
	m.epicList.Title = "My Fave epics"
	m.storyList.Title = "My Fave stories"
	m.taskList.Title = "My Fave tasks"
	return m
}
