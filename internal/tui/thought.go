package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ThoughtModel is the model for quick thought entry
type ThoughtModel struct {
	input       textinput.Model
	timestamp   time.Time
	width       int
	height      int
	submitted   bool
	thoughtText string
}

// NewThoughtModel creates a new thought model
func NewThoughtModel() ThoughtModel {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.Focus()
	ti.Width = 70

	return ThoughtModel{
		input:     ti,
		timestamp: time.Now(),
		submitted: false,
	}
}

// Init initializes the model
func (m ThoughtModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m ThoughtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit thought
			if m.input.Value() != "" {
				m.submitted = true
				m.thoughtText = m.input.Value()
				return m, tea.Quit
			}
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the UI
func (m ThoughtModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	b.WriteString(HeaderStyle.Render("💭 QUICK THOUGHT"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(m.timestamp.Format("3:04pm")))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	b.WriteString(DimStyle.Render("No tags or momentum needed - just your thought"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("Enter to save • Ctrl+C to cancel"))

	return BoxStyle.Render(b.String())
}

// GetThought returns the thought text
func (m ThoughtModel) GetThought() string {
	return m.thoughtText
}

// WasSubmitted returns whether the thought was submitted
func (m ThoughtModel) WasSubmitted() bool {
	return m.submitted
}
