package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// IntentionModel is the model for the morning intention prompt
type IntentionModel struct {
	input     textinput.Model
	submitted bool
	intention string
}

// NewIntentionModel creates a new intention model
func NewIntentionModel() IntentionModel {
	ti := textinput.New()
	ti.Placeholder = "What's your intention for today?"
	ti.Focus()
	ti.Width = 70

	return IntentionModel{
		input:     ti,
		submitted: false,
	}
}

// Init initializes the model
func (m IntentionModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m IntentionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit intention (can be empty to skip)
			m.submitted = true
			m.intention = strings.TrimSpace(m.input.Value())
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			// Cancel without setting intention
			m.submitted = true
			m.intention = ""
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the UI
func (m IntentionModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Greeting
	b.WriteString(HeaderStyle.Render("Good morning! ⋆｡˚☀︎｡⋆˚"))
	b.WriteString("\n\n")

	// Prompt
	b.WriteString(BoldStyle.Render("Today's Intention:"))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// Helper text
	b.WriteString(DimStyle.Render("Enter to continue (or leave blank to skip)"))

	return BoxStyle.Render(b.String())
}

// GetIntention returns the submitted intention
func (m IntentionModel) GetIntention() string {
	return m.intention
}

// WasSubmitted returns whether the form was submitted
func (m IntentionModel) WasSubmitted() bool {
	return m.submitted
}
