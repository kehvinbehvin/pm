package application

import (
	"github/pm/pkg/fileSystem"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/glamour"

)

func NewApplication(fs *fileSystem.FileSystem) (tea.Model, error) {
	stack := NewApplicationStack()
	welcome := &WelcomeFrame{}
	stack.Push(welcome)

	const width = 78
	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, renderErr := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)

	if renderErr != nil {
		return nil, errors.New("Error creating renderer");
	}

	graphRenderer := fileSystem.FileGraphRenderer{
		Fs: *fs,
	}

	return Application{
		History: stack,
		Fs: fs,
		Renderer: renderer,
		ViewPort: &vp,
		GraphRenderer: &graphRenderer,
	}, nil
}

type Application struct {
	History *ApplicationStack
	Fs *fileSystem.FileSystem
	Renderer *glamour.TermRenderer
	ViewPort *viewport.Model
	GraphRenderer *fileSystem.FileGraphRenderer
}

func (a Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentFrame, error := a.History.Peek()
	if error != nil {
		return a, tea.Quit
	}

	return currentFrame.Update(msg, a)
}

func (a Application) View() string {
	currentFrame, error := a.History.Peek()
	if error != nil {
		// Maybe introduce an error page here
		return ""
	}

	currentFrame.Refresh(a)

	return currentFrame.View(a)
}
func (a Application) Init() tea.Cmd {
	currentFrame, error := a.History.Peek()
	if error != nil {
		return nil
	}

	return currentFrame.Init(a)
}
