package application

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/glamour"
	"log"
)

type ViewMarkdownFrame struct{
	content string
	fileName string
}

func NewViewMarkdownFrame(fileName string, content string, app Application) (*ViewMarkdownFrame, error) {
	str, err := app.Renderer.Render(content)
	if err != nil {
		return nil, err
	}

	log.Println("Rendered content");
	app.ViewPort.SetContent(str)

	log.Println("Created viewport");
	return &ViewMarkdownFrame{
		fileName: fileName,
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
		case "e":
			app.Fs.EditFile(viewMarkdownFrame.fileName)

			// Push back to previous page
			// Issue with terminal resizing after editor process hands back control
			app.History.Pop()
		default:
			var cmd tea.Cmd
			*app.ViewPort, cmd = app.ViewPort.Update(msg)
			return app, tea.Batch(tea.ClearScreen, cmd)
		}
	default:
		return app, nil
	}

	return app, nil
}

func (vmdf *ViewMarkdownFrame) View(app Application) string {
	log.Println("Viewing Markdown");
	helptext := "\n[q] Quit ● [←] Back ● [e] Edit"
	marginStyle := lipgloss.NewStyle().Margin(1, 2)

	return app.ViewPort.View() + marginStyle.Render(helptext)
}

func (vmdf *ViewMarkdownFrame) Init() tea.Cmd {
	return nil
}
