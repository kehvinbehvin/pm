package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"errors"
)

type DependencyGraph struct{
	fileName string
}

func NewDependencyGraph(fileName string) (*DependencyGraph, error) {
	return &DependencyGraph{
		fileName: fileName,
	}, nil
}

func (dg DependencyGraph) Update(msg tea.Msg, app Application) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return app, tea.Quit
			case "left":
				app.History.Pop()
		}
	}

	return app, nil
}

func (dg DependencyGraph) getFrame(app Application) (*DependencyGraph, error) {
	frame, error := app.History.Peek()
	if error != nil {
		return &DependencyGraph{}, errors.New("Cannot get self")
	}

	depGraphframe := frame.(*DependencyGraph)
	return depGraphframe, nil
}


func (dg DependencyGraph) View(app Application) string {
	helptext := "\n[q] Quit ● [←] Back "
	marginStyle := lipgloss.NewStyle().Margin(1, 2)

	depGraphFrame, frameErr := dg.getFrame(app);
	if frameErr != nil {
		return ""
	}

	// vertex := app.Fs.GetFileChildMeta(depGraphFrame.fileName)

	graph, renderErr := app.GraphRenderer.Build(depGraphFrame.fileName)
	if renderErr != nil {
		return ""
	}

	return graph + marginStyle.Render(helptext)
}

func (dg DependencyGraph) Init(app Application) tea.Cmd {
	return nil
}
