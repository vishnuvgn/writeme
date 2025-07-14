package helpers

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PreviewModel struct {
	linesAbove []string
	linesBelow []string
	input      textinput.Model
	confirmed  bool
	finalNote  string
}

func NewPreviewModel(linesAbove, linesBelow []string, initialNote string) PreviewModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = "+ - " // show + and bullet
	ti.CharLimit = 500
	ti.SetValue(initialNote)
	ti.Focus()
	ti.CursorEnd()

	return PreviewModel{
		linesAbove: linesAbove,
		linesBelow: linesBelow,
		input:      ti,
	}
}

func (m PreviewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.finalNote = m.input.Value()
			m.confirmed = true
			return m, tea.Quit
		case "esc", "q":
			m.confirmed = false
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m PreviewModel) View() string {
	var b strings.Builder

	title := lipgloss.NewStyle().Bold(true).Render("--- Proposed Change Preview ---")
	b.WriteString("\n" + title + "\n\n")

	for _, line := range m.linesAbove {
		b.WriteString("  " + line + "\n")
	}
	greenInput := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(m.input.View())
	b.WriteString(greenInput + "\n")

	// b.WriteString(m.input.View() + "\n")

	for _, line := range m.linesBelow {
		b.WriteString("  " + line + "\n")
	}

	b.WriteString("\n[Enter = confirm, Esc = cancel]\n")
	return b.String()
}

func RunPreviewWithEdit(linesAbove, linesBelow []string, initialNote string) (string, bool, error) {
	m := NewPreviewModel(linesAbove, linesBelow, initialNote)

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", false, err
	}

	mFinal := finalModel.(PreviewModel)
	return mFinal.finalNote, mFinal.confirmed, nil
}

// package helpers

// import (
// 	"strings"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// // previewModel holds the snippet and whether the user confirmed.
// type previewModel struct {
// 	snippet   []string
// 	confirmed bool
// }

// // Init is a no-op.
// func (m previewModel) Init() tea.Cmd {
// 	return nil
// }

// // Update handles keypresses.
// func (m previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "enter":
// 			m.confirmed = true
// 			return m, tea.Quit
// 		case "esc", "q":
// 			m.confirmed = false
// 			return m, tea.Quit
// 		}
// 	}
// 	return m, nil
// }

// // View renders the snippet.
// func (m previewModel) View() string {
// 	var b strings.Builder

// 	b.WriteString("\n--- Proposed Change Preview ---\n\n")

// 	for _, line := range m.snippet {
// 		b.WriteString(line + "\n")
// 	}

// 	b.WriteString("\n[Enter = confirm, Esc = cancel]\n")

// 	return b.String()
// }

// // RunPreview runs the Bubble Tea program and returns confirmed status.
// func RunPreview(snippet []string) (bool, error) {
// 	m := previewModel{snippet: snippet}

// 	prog := tea.NewProgram(m)
// 	finalModel, err := prog.Run()
// 	if err != nil {
// 		return false, err
// 	}

// 	return finalModel.(previewModel).confirmed, nil
// }

// func BuildPreviewSnippet(lines []string, insertAt int) []string {
// 	start := insertAt - 2
// 	if start < 0 {
// 		start = 0
// 	}

// 	end := insertAt + 2
// 	if end >= len(lines) {
// 		end = len(lines) - 1
// 	}

// 	snippet := []string{}

// 	for i := start; i <= end; i++ {
// 		line := lines[i]
// 		if i == insertAt {
// 			// Highlight inserted line with "+" and green color
// 			line = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("+ " + line)
// 		} else {
// 			line = "  " + line
// 		}
// 		snippet = append(snippet, line)
// 	}

// 	return snippet
// }
