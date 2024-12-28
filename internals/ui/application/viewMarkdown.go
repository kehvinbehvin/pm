package application

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"log"
)

type ViewMarkdownFrame struct{
	content string
}

func NewViewMarkdownFrame(content string, app Application) (*ViewMarkdownFrame, error) {
	str, err := app.Renderer.Render(content)
	if err != nil {
		return nil, err
	}

	log.Println("Rendered content");
	app.ViewPort.SetContent(str)

	log.Println("Created viewport");
	return &ViewMarkdownFrame{
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
		app.ViewPort.Width = msg.Width
		app.ViewPort.Height = msg.Height
		
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(78),
		)

		// Re-render content to fit the new size
		rendered, err := renderer.Render(viewMarkdownFrame.content)
		if err != nil {
			return app, tea.Quit
		} else {
			app.ViewPort.SetContent(rendered)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return app, tea.Quit
		case "left":
			app.History.Pop()
		default:
			var cmd tea.Cmd
			*app.ViewPort, cmd = app.ViewPort.Update(msg)
			return app, cmd
		}
	default:
		return app, nil
	}

	return app, nil
}

func (vmdf *ViewMarkdownFrame) View(app Application) string {
	log.Println("Viewing Markdown");
	return app.ViewPort.View()
}

func (vmdf *ViewMarkdownFrame) Init() tea.Cmd {
	return nil
}
