package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// QuickWinModel is the model for quick win entry
type QuickWinModel struct {
	input     textinput.Model
	width     int
	height    int
	submitted bool
	winText   string
}

// NewQuickWinModel creates a new quick win model
func NewQuickWinModel() QuickWinModel {
	ti := textinput.New()
	ti.Placeholder = "What's your win today?"
	ti.Focus()
	ti.Width = 70

	return QuickWinModel{
		input:     ti,
		submitted: false,
	}
}

// Init initializes the model
func (m QuickWinModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m QuickWinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit win
			if m.input.Value() != "" {
				m.submitted = true
				m.winText = m.input.Value()
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
func (m QuickWinModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	b.WriteString(SuccessStyle.Render("🌟 LOG A WIN"))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	b.WriteString(DimStyle.Render("Enter to save • Ctrl+C to cancel"))

	return BoxStyle.Render(b.String())
}

// GetWin returns the win text
func (m QuickWinModel) GetWin() string {
	return m.winText
}

// WasSubmitted returns whether the win was submitted
func (m QuickWinModel) WasSubmitted() bool {
	return m.submitted
}
