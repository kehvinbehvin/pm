package application

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"fmt"
	"strings"
	"errors"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type CreateFormFrame struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func NewCreateFormFrame() (ApplicationFrame) {
	m := CreateFormFrame{
		inputs: make([]textinput.Model, 3),
		focusIndex: 0,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Others"
			t.CharLimit = 64
}

		m.inputs[i] = t
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "p":
			app.History.Pop()
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && createFormFrame.focusIndex == len(createFormFrame.inputs) {
				return app, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				createFormFrame.focusIndex--
			} else {
				createFormFrame.focusIndex++
			}

			if createFormFrame.focusIndex > len(createFormFrame.inputs) {
				createFormFrame.focusIndex = 0
			} else if createFormFrame.focusIndex < 0 {
				createFormFrame.focusIndex = len(createFormFrame.inputs)
			}

			cmds := make([]tea.Cmd, len(createFormFrame.inputs))
			for i := 0; i <= len(createFormFrame.inputs)-1; i++ {
				if i == createFormFrame.focusIndex {
					// Set focused state
					cmds[i] = createFormFrame.inputs[i].Focus()
					createFormFrame.inputs[i].PromptStyle = focusedStyle
					createFormFrame.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				createFormFrame.inputs[i].Blur()
				createFormFrame.inputs[i].PromptStyle = noStyle
				createFormFrame.inputs[i].TextStyle = noStyle
			}

			return app, tea.Batch(cmds...)
		}
	}

	cmd := createFormFrame.updateInputs(msg, app)

	return app, cmd
}

func (cf CreateFormFrame) updateInputs(msg tea.Msg, app Application) tea.Cmd {
	createFormFrame, frameErr := cf.getFrame(app);
	if frameErr != nil {
		return tea.Quit
	}

	cmds := make([]tea.Cmd, len(createFormFrame.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range createFormFrame.inputs {
		createFormFrame.inputs[i], cmds[i] = createFormFrame.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (cf CreateFormFrame) View(app Application) (string) {
	createFormFrame, frameErr := cf.getFrame(app);
	if frameErr != nil {
		return ""
	}

	var b strings.Builder

	for i := range createFormFrame.inputs {
		b.WriteString(createFormFrame.inputs[i].View())
		if i < len(createFormFrame.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if createFormFrame.focusIndex == len(createFormFrame.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func (cf CreateFormFrame) Init() (tea.Cmd) {
	return textinput.Blink
}
