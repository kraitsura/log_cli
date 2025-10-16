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
	}
}

// Init initializes the model
func (m LogEntryModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m LogEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	b.WriteString("\n\n")

	// Helper text
	helperText := "↑ ↓ → | @deep @social @admin @break @zone | [LEAK] [FLOW] [STUCK] [GOLD]"
	b.WriteString(DimStyle.Render(helperText))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Enter to submit • Ctrl+C to cancel"))

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
