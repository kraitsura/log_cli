package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// EditModel is the model for the edit entry screen
type EditModel struct {
	input      textinput.Model
	entryIndex int
	timestamp  time.Time
	width      int
	height     int
	submitted  bool
	entryText  string
}

// NewEditModel creates a new edit entry model with pre-filled text
func NewEditModel(entryIndex int, timestamp time.Time, originalText string) EditModel {
	ti := textinput.New()
	ti.Placeholder = "What are you doing right now?"
	ti.Focus()
	ti.Width = 70
	ti.SetValue(originalText) // Pre-fill with existing text

	return EditModel{
		input:      ti,
		entryIndex: entryIndex,
		timestamp:  timestamp,
		submitted:  false,
	}
}

// Init initializes the model
func (m EditModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit entry
			if m.input.Value() != "" {
				m.submitted = true
				m.entryText = m.input.Value()
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
func (m EditModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Header
	header := "EDIT ENTRY"
	b.WriteString(HeaderStyle.Render(header))
	b.WriteString("\n\n")

	// Entry info
	info := fmt.Sprintf("Editing entry #%d from %s", m.entryIndex, m.timestamp.Format("3:04pm"))
	b.WriteString(DimStyle.Render(info))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// Helper text
	helperText := "Momentum: ++ -- == << (or -> <-) for ↑ ↓ → ←"
	b.WriteString(DimStyle.Render(helperText))
	b.WriteString("\n")
	helperTags := "@deep @social @admin @break @zone | [LEAK] [FLOW] [STUCK] [GOLD]"
	b.WriteString(DimStyle.Render(helperTags))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Enter to save • Ctrl+C to cancel"))

	return BoxStyle.Render(b.String())
}

// GetEntryText returns the submitted entry text
func (m EditModel) GetEntryText() string {
	return m.entryText
}

// WasSubmitted returns whether the entry was submitted
func (m EditModel) WasSubmitted() bool {
	return m.submitted
}
