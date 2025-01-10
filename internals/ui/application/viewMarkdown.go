package application

import (
	"errors"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/glamour"
	"log"
)

type ViewMarkdownFrame struct{
	content string
	fileName string
	subStack *ApplicationStack
	selectedItem string
	linkChild bool
	linkDownStream bool
	linkUpsteam bool
}
// TODO: Need to handle window sizing.
func NewViewMarkdownFrame(fileName string, content string, app Application) (*ViewMarkdownFrame, error) {
	str, err := app.Renderer.Render(content)
	if err != nil {
		return nil, err
	}

	log.Println("Rendered content");
	app.ViewPort.SetContent(str)

	log.Println("Created viewport");

	stack := NewApplicationStack()
	return &ViewMarkdownFrame{
		fileName: fileName,
		content: content,
		subStack: stack,
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

	// Register search results
	if viewMarkdownFrame.subStack.Size() != 0 {
		frame, frameErr := viewMarkdownFrame.subStack.Peek();
		if frameErr != nil {
			return app, tea.Quit
		}

		globalSelectFrame := frame.(*GlobalSelectionFrame)
		selectedItem := globalSelectFrame.selectedItem
		
		if selectedItem != "" {
			index := strings.Index(selectedItem, "]")
			fileName := strings.TrimSpace(selectedItem[index + 1:])
			viewMarkdownFrame.selectedItem = fileName
		}
		
		viewMarkdownFrame.subStack.Pop()
	}

	if viewMarkdownFrame.linkChild && viewMarkdownFrame.selectedItem != "" {
		app.Fs.LinkHierarchy(viewMarkdownFrame.fileName, viewMarkdownFrame.selectedItem)
		viewMarkdownFrame.linkChild = false
		viewMarkdownFrame.selectedItem = ""
	} else if viewMarkdownFrame.linkDownStream {
		app.Fs.LinkDependency(viewMarkdownFrame.fileName, viewMarkdownFrame.selectedItem)
		viewMarkdownFrame.linkDownStream = false
		viewMarkdownFrame.selectedItem = ""
	} else if viewMarkdownFrame.linkUpsteam {
		app.Fs.LinkDependency(viewMarkdownFrame.selectedItem, viewMarkdownFrame.fileName)
		viewMarkdownFrame.linkUpsteam = false
		viewMarkdownFrame.selectedItem = ""
	}

	// Continue operation
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
		case "o":
			app.Fs.EditFile(viewMarkdownFrame.fileName)

			// Push back to previous page
			// Issue with terminal resizing after editor process hands back control
			app.History.Pop()
			vmdf.subStack.ClearStack()
		case "i":
			createFormFrame, frameErr := NewCreateFormFrame(app, viewMarkdownFrame.fileName)
			if frameErr != nil {
				return app, tea.Quit
			}

			app.History.Push(createFormFrame)
		case "l":
			// Remove the ability to link children outside of creating a new issue
			// to prevent multiple parents
			// Push fileName to a global search of all issues which exlcudes itself
			// relatedIssues, issuesErr := app.Fs.ListRelatedHierarchy(viewMarkdownFrame.fileName)
			// if issuesErr != nil {
			// 	log.Println("Error fetching related children issues")
			// 	return app, tea.Quit
			// }
			// globalSearchFrame, frameErr := NewGlobalSelectionFrame(app, viewMarkdownFrame.fileName, relatedIssues)
			// if frameErr != nil {
			// 	return app, nil 
			// }
			//
			// app.History.Push(globalSearchFrame)
			// viewMarkdownFrame.subStack.Push(globalSearchFrame)
			// viewMarkdownFrame.linkChild = true
		case "d":
			// Push fileName to a global search of all issues which exlcudes itself
			relatedIssues, issuesErr := app.Fs.ListRelatedDependency(viewMarkdownFrame.fileName)
			log.Println(relatedIssues)
			if issuesErr != nil {
				log.Println("Error fetching related children issues")
				return app, tea.Quit
			}

			globalSearchFrame, frameErr := NewGlobalSelectionFrame(app, viewMarkdownFrame.fileName, relatedIssues)
			if frameErr != nil {
				return app, nil 
			}

			app.History.Push(globalSearchFrame)
			viewMarkdownFrame.subStack.Push(globalSearchFrame)
			viewMarkdownFrame.linkDownStream = true
		case "u":
			// Push fileName to a global search of all issues which exlcudes itself
			relatedIssues, issuesErr := app.Fs.ListRelatedDependency(viewMarkdownFrame.fileName)
			if issuesErr != nil {
				log.Println("Error fetching related children issues")
				return app, tea.Quit
			}

			globalSearchFrame, frameErr := NewGlobalSelectionFrame(app, viewMarkdownFrame.fileName, relatedIssues)
			if frameErr != nil {
				return app, nil 
			}

			app.History.Push(globalSearchFrame)
			viewMarkdownFrame.subStack.Push(globalSearchFrame)
			viewMarkdownFrame.linkUpsteam = true
		case "e":
			frame := NewBrowseFrame(app, "epic")
			app.History.Push(frame)
		case "s":
			frame := NewBrowseFrame(app, "story")
			app.History.Push(frame)

		case "t":
			frame := NewBrowseFrame(app, "task")
			app.History.Push(frame)
		}
	default:
		return app, nil
	}

	return app, nil
}

func (vmdf *ViewMarkdownFrame) View(app Application) string {
	helptext := "[o] Open [d] Link Downstream blocker [u] Link Upstream blocker\n[i] Create child issue\n[q] Quit ● [←] Back\n[e] All epics [s] All stories [t] All tasks"
	marginStyle := lipgloss.NewStyle().Margin(1, 2)

	return app.ViewPort.View() + marginStyle.Render(helptext)
}

func (vmdf *ViewMarkdownFrame) Init(app Application) tea.Cmd {
	return nil
}
