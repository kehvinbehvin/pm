package application

import (
	"errors"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/glamour"
)

type ViewMarkdownFrame struct{
	viewport viewport.Model
	content string
}

func NewViewMarkdownFrame(content string) (*ViewMarkdownFrame, error) {
	const width = 78

	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	
	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)

	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &ViewMarkdownFrame{
		viewport: vp,
		content: content,
	}, nil
}

func (vmdf *ViewMarkdownFrame) getFrame(app Application) (*ViewMarkdownFrame, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &ViewMarkdownFrame{}, errors.New("Cannot get self")
	}

	viewMarkDownFrame := frame.(*ViewMarkdownFrame)
	return viewMarkDownFrame, nil
}


func (vmdf *ViewMarkdownFrame) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	viewMarkdownFrame, frameErr := vmdf.getFrame(app)
	if frameErr != nil {
		return app, tea.Quit
	}


	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update the viewport size when the terminal is resized
		viewMarkdownFrame.viewport.Width = msg.Width
		viewMarkdownFrame.viewport.Height = msg.Height
		
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(78),
		)

		// Re-render content to fit the new size
		rendered, err := renderer.Render(viewMarkdownFrame.content)
		if err != nil {
			return app, tea.Quit
		} else {
			viewMarkdownFrame.viewport.SetContent(rendered)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "left":
			app.History.Pop()
		default:
			var cmd tea.Cmd
			viewMarkdownFrame.viewport, cmd = viewMarkdownFrame.viewport.Update(msg)
			return app, cmd
		}
	default:
		return app, nil
	}

	return app, nil
}

func (vmdf *ViewMarkdownFrame) View(app Application) string {
	return vmdf.viewport.View()
}

func (vmdf *ViewMarkdownFrame) Init() tea.Cmd {
	return nil
}
