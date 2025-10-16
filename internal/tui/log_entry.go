package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// LogEntryModel is the model for the log entry screen
type LogEntryModel struct {
	input        textinput.Model
	timestamp    time.Time
	lastLogTime  time.Time
	entryCount   int
	isDriftAlert bool
	width        int
	height       int
	submitted    bool
	entryText    string
	autocomplete AutocompleteState
}

// NewLogEntryModel creates a new log entry model
func NewLogEntryModel(lastLog time.Time, count int) LogEntryModel {
	ti := textinput.New()
	ti.Placeholder = "What are you doing right now?"
	ti.Focus()
	ti.Width = 70

	now := time.Now()
	isDrift := false
	if !lastLog.IsZero() && now.Sub(lastLog).Minutes() >= 90 {
		isDrift = true
	}

	return LogEntryModel{
		input:        ti,
		timestamp:    now,
		lastLogTime:  lastLog,
		entryCount:   count,
		isDriftAlert: isDrift,
		submitted:    false,
		autocomplete: NewAutocompleteState(),
	}
}

// Init initializes the model
func (m LogEntryModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m LogEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	previousValue := m.input.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle autocomplete navigation when active
		if m.autocomplete.Active {
			switch msg.Type {
			case tea.KeyUp:
				m.autocomplete.MoveUp()
				return m, nil
			case tea.KeyDown:
				m.autocomplete.MoveDown()
				return m, nil
			case tea.KeyTab, tea.KeyEnter:
				// If Enter and autocomplete has suggestions, insert suggestion
				if msg.Type == tea.KeyEnter && m.autocomplete.GetSelectedSuggestion() != "" {
					newText := m.autocomplete.InsertSuggestion(m.input.Value())
					m.input.SetValue(newText)
					m.autocomplete.Deactivate()
					return m, nil
				}
				// If Tab, always insert suggestion if available
				if msg.Type == tea.KeyTab {
					suggestion := m.autocomplete.GetSelectedSuggestion()
					if suggestion != "" {
						newText := m.autocomplete.InsertSuggestion(m.input.Value())
						m.input.SetValue(newText)
						m.autocomplete.Deactivate()
						return m, nil
					}
				}
				// Fall through if Enter with no suggestion (submit entry)
			case tea.KeyEsc:
				m.autocomplete.Deactivate()
				return m, nil
			}
		}

		// Handle normal key presses
		switch msg.Type {
		case tea.KeyEnter:
			// Submit entry (only if not handled by autocomplete above)
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

	// Update the input
	m.input, cmd = m.input.Update(msg)

	// If text changed, apply momentum conversion and update autocomplete
	if m.input.Value() != previousValue {
		newValue := m.input.Value()

		// Apply real-time momentum marker conversion
		newValue = convertMomentumMarkers(newValue)

		// Update input if changed
		if newValue != m.input.Value() {
			cursorPos := m.input.Position()
			m.input.SetValue(newValue)
			// Adjust cursor position if needed (momentum markers change text length)
			if cursorPos > len(newValue) {
				m.input.SetCursor(len(newValue))
			} else {
				m.input.SetCursor(cursorPos)
			}
		}

		// Update autocomplete state
		m.autocomplete.Update(m.input.Value(), m.input.Position())
	}

	return m, cmd
}

// View renders the UI
func (m LogEntryModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Header with date
	header := fmt.Sprintf("DAYLOG - %s", time.Now().Format("Monday, January 2, 2006"))
	b.WriteString(HeaderStyle.Render(header))
	b.WriteString("\n")

	// Timestamp when logging started
	startTime := DimStyle.Render(fmt.Sprintf("Started: %s", m.timestamp.Format("3:04pm")))
	b.WriteString(startTime)
	b.WriteString("\n\n")

	// Drift alert if applicable
	if m.isDriftAlert {
		duration := time.Since(m.lastLogTime)
		alert := fmt.Sprintf("[!] DRIFT ALERT - Last log was %s ago", formatDuration(duration))
		b.WriteString(AlertStyle.Render(alert))
		b.WriteString("\n\n")
	}

	// Main prompt
	prompt := fmt.Sprintf("%s | What are you doing right now?", m.timestamp.Format("3:04pm"))
	b.WriteString(BoldStyle.Render(prompt))
	b.WriteString("\n\n")

	// Text input
	b.WriteString(m.input.View())
	b.WriteString("\n")

	// Autocomplete dropdown (if active)
	if m.autocomplete.Active && len(m.autocomplete.Suggestions) > 0 {
		b.WriteString("\n")
		b.WriteString(m.autocomplete.View())
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	// Helper text - show shortcuts for momentum
	helperText := "Momentum: ++ -- == << (or -> <-) for ↑ ↓ → ←"
	b.WriteString(DimStyle.Render(helperText))
	b.WriteString("\n")
	helperTags := "@deep @social @admin @break @zone | [LEAK] [FLOW] [STUCK] [GOLD]"
	b.WriteString(DimStyle.Render(helperTags))
	b.WriteString("\n\n")

	// Control hints - update to show Tab option when autocomplete is active
	if m.autocomplete.Active {
		b.WriteString(DimStyle.Render("↑↓ to navigate • Tab/Enter to select • Esc to close • Ctrl+C to cancel"))
	} else {
		b.WriteString(DimStyle.Render("Enter to submit • Ctrl+C to cancel"))
	}

	return BoxStyle.Render(b.String())
}

// GetEntryText returns the submitted entry text
func (m LogEntryModel) GetEntryText() string {
	return m.entryText
}

// WasSubmitted returns whether the entry was submitted
func (m LogEntryModel) WasSubmitted() bool {
	return m.submitted
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

// convertMomentumMarkers converts momentum shortcuts to arrow symbols in real-time
func convertMomentumMarkers(text string) string {
	// Replace shortcuts with arrows
	// Order matters - replace longer patterns first to avoid partial replacements
	replacements := []struct{ old, new string }{
		{"++", "↑"},
		{"--", "↓"},
		{"==", "→"},
		{"<<", "←"},
		{"->", "→"},
		{"<-", "←"},
	}

	for _, r := range replacements {
		text = strings.ReplaceAll(text, r.old, r.new)
	}

	return text
}
