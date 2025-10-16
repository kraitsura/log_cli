package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SignoffModel is the model for the end-of-day reflection
type SignoffModel struct {
	inputs    []textinput.Model
	current   int
	questions []string
	submitted bool
	answers   []string
	intention string
}

// NewSignoffModel creates a new sign-off reflection model
func NewSignoffModel(intention string) SignoffModel {
	questions := []string{
		"What pulled you off track today?",
		"What kept you on track today?",
		"One thing you'll protect tomorrow?",
	}

	inputs := make([]textinput.Model, len(questions))
	for i := range inputs {
		ti := textinput.New()
		ti.Width = 70
		ti.CharLimit = 200
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	return SignoffModel{
		inputs:    inputs,
		current:   0,
		questions: questions,
		submitted: false,
		answers:   make([]string, len(questions)),
		intention: intention,
	}
}

// Init initializes the model
func (m SignoffModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m SignoffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Save current answer
			m.answers[m.current] = strings.TrimSpace(m.inputs[m.current].Value())

			// Move to next question or submit
			if m.current < len(m.questions)-1 {
				m.inputs[m.current].Blur()
				m.current++
				m.inputs[m.current].Focus()
				return m, textinput.Blink
			} else {
				// All questions answered, submit
				m.submitted = true
				return m, tea.Quit
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			// Cancel sign-off
			m.submitted = true
			// Clear answers to indicate cancellation
			m.answers = make([]string, len(m.questions))
			return m, tea.Quit

		case tea.KeyShiftTab:
			// Go back to previous question
			if m.current > 0 {
				m.inputs[m.current].Blur()
				m.current--
				m.inputs[m.current].Focus()
				return m, textinput.Blink
			}
		}
	}

	var cmd tea.Cmd
	m.inputs[m.current], cmd = m.inputs[m.current].Update(msg)
	return m, cmd
}

// View renders the UI
func (m SignoffModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString(HeaderStyle.Render("DAY COMPLETE 🌙"))
	b.WriteString("\n\n")

	// Show intention if set
	if m.intention != "" {
		b.WriteString(BoldStyle.Render("Today's Intention: "))
		b.WriteString(m.intention)
		b.WriteString("\n\n")
	}

	// Show all questions and inputs
	for i, question := range m.questions {
		if i == m.current {
			// Current question - highlighted
			b.WriteString(PromptStyle.Render(question))
			b.WriteString("\n")
			b.WriteString(m.inputs[i].View())
			b.WriteString("\n")
		} else if i < m.current {
			// Previous question - show answer
			b.WriteString(DimStyle.Render(question))
			b.WriteString("\n")
			if m.answers[i] != "" {
				b.WriteString(DimStyle.Render(m.answers[i]))
			} else {
				b.WriteString(DimStyle.Render("(skipped)"))
			}
			b.WriteString("\n")
		} else {
			// Future question - dimmed
			b.WriteString(DimStyle.Render(question))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Helper text
	b.WriteString(DimStyle.Render("Enter to continue | Shift+Tab to go back | Esc to skip"))

	return BoxStyle.Render(b.String())
}

// GetReflections returns the submitted reflections
func (m SignoffModel) GetReflections() (pulledOffTrack, keptOnTrack, tomorrowProtect string) {
	if len(m.answers) >= 3 {
		return m.answers[0], m.answers[1], m.answers[2]
	}
	return "", "", ""
}

// WasSubmitted returns whether the form was submitted
func (m SignoffModel) WasSubmitted() bool {
	return m.submitted
}
