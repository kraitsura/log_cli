package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// WinModel is the model for the 10-entry win prompt
type WinModel struct {
	input     textinput.Model
	submitted bool
	win       string
}

// NewWinModel creates a new win model
func NewWinModel() WinModel {
	ti := textinput.New()
	ti.Placeholder = "Any wins to celebrate?"
	ti.Focus()
	ti.Width = 70

	return WinModel{
		input:     ti,
		submitted: false,
	}
}

// Init initializes the model
func (m WinModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m WinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit win (can be empty to skip)
			m.submitted = true
			m.win = strings.TrimSpace(m.input.Value())
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			// Cancel without setting win
			m.submitted = true
			m.win = ""
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the UI
func (m WinModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Celebration
	b.WriteString(HeaderStyle.Render("You've logged 10 entries today!"))
	b.WriteString("\n\n")

	// Prompt
	b.WriteString(BoldStyle.Render("Any wins today?"))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// Helper text
	b.WriteString(DimStyle.Render("Enter to continue (or leave blank to skip)"))

	return BoxStyle.Render(b.String())
}

// GetWin returns the submitted win
func (m WinModel) GetWin() string {
	return m.win
}

// WasSubmitted returns whether the form was submitted
func (m WinModel) WasSubmitted() bool {
	return m.submitted
}
