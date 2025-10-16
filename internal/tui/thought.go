package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ThoughtModel is the model for quick thought entry
type ThoughtModel struct {
	input       textarea.Model
	timestamp   time.Time
	width       int
	height      int
	submitted   bool
	thoughtText string
}

// NewThoughtModel creates a new thought model
func NewThoughtModel() ThoughtModel {
	ta := textarea.New()
	ta.Placeholder = "What's on your mind? Share what you're working on, feeling, or thinking about..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0

	// Start with generous height for expansive feel
	ta.SetHeight(6)
	ta.SetWidth(70) // Will be updated in WindowSizeMsg

	return ThoughtModel{
		input:     ta,
		timestamp: time.Now(),
		submitted: false,
	}
}

// Init initializes the model
func (m ThoughtModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages
func (m ThoughtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD:
			// Ctrl+D to submit thought (Enter creates new lines)
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

		// Set textarea width to 85% of terminal width
		newWidth := int(float64(msg.Width) * 0.85)
		if newWidth < 40 {
			newWidth = 40 // minimum width
		}
		m.input.SetWidth(newWidth)
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

	// Minimal header
	b.WriteString(HeaderStyle.Render("💭 THOUGHT"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(m.timestamp.Format("3:04pm")))
	b.WriteString("\n\n")

	// Expansive text area
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// Simplified helper text
	b.WriteString(DimStyle.Render("No tags or momentum needed - just your thought"))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Ctrl+D to save • Ctrl+C to cancel"))

	// Use minimal, open box style - no visible border, generous padding
	openBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Padding(2, 3)

	return openBoxStyle.Render(b.String())
}

// GetThought returns the thought text
func (m ThoughtModel) GetThought() string {
	return m.thoughtText
}

// WasSubmitted returns whether the thought was submitted
func (m ThoughtModel) WasSubmitted() bool {
	return m.submitted
}
