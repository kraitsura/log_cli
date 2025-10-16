package tui

import (
	"strings"
	"time"

	"github.com/aaryareddy/log_cli/internal/database"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewModel is the model for viewing today's log entries
type ViewModel struct {
	day              *database.Day
	entries          []*database.Entry
	viewport         viewport.Model
	ready            bool
	textExpanded     bool // Toggle for rolling/unrolling long text
	maxCollapsedLen  int  // Maximum characters before text is collapsed
}

// NewViewModel creates a new view model
func NewViewModel(day *database.Day, entries []*database.Entry) ViewModel {
	return ViewModel{
		day:             day,
		entries:         entries,
		textExpanded:    false, // Start with text collapsed
		maxCollapsedLen: 100,   // Collapse text longer than 100 chars
	}
}

// Init initializes the model
func (m ViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "shift+r", "R":
			// Toggle text expansion
			m.textExpanded = !m.textExpanded
			// Regenerate content with new expansion state
			m.viewport.SetContent(m.generateViewContent())
			return m, nil
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Initialize viewport on first window size message
			m.viewport = viewport.New(msg.Width, msg.Height-2)
			m.viewport.SetContent(m.generateViewContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
	}

	// Update viewport (handles scrolling)
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// generateViewContent generates the log view content
func (m ViewModel) generateViewContent() string {
	var b strings.Builder

	// Header with date
	dateStr := m.day.Date.Format("Monday, January 2, 2006")
	b.WriteString(RenderHeaderBar("DAYLOG", dateStr))
	b.WriteString("\n\n")

	// Show intention if set
	if m.day.Intention != nil && *m.day.Intention != "" {
		b.WriteString(BoldStyle.Render("Intention: "))
		b.WriteString(*m.day.Intention)
		b.WriteString("\n\n")
	}

	// Show entries
	if len(m.entries) == 0 {
		b.WriteString(DimStyle.Render("No logs yet today. Start logging to build your daylog!"))
		b.WriteString("\n\n")
		b.WriteString(DimStyle.Render("Type: log"))
	} else {
		// Split entries if day is completed
		regularEntries, afterHoursEntries := splitEntries(m.entries, m.day.Completed)

		// Display regular entries (wins now appear inline with timestamps)
		b.WriteString(m.formatEntries(regularEntries))

		// Show reflections if day is completed
		if m.day.Completed {
			b.WriteString("\n\n")
			b.WriteString(DimStyle.Render("─────────────────────────────────────"))
			b.WriteString("\n\n")

			if m.day.PulledOffTrack != nil && *m.day.PulledOffTrack != "" {
				b.WriteString(DimStyle.Render("Pulled off track: "))
				b.WriteString(*m.day.PulledOffTrack)
				b.WriteString("\n")
			}

			if m.day.KeptOnTrack != nil && *m.day.KeptOnTrack != "" {
				b.WriteString(DimStyle.Render("Kept on track: "))
				b.WriteString(*m.day.KeptOnTrack)
				b.WriteString("\n")
			}

			if m.day.TomorrowProtect != nil && *m.day.TomorrowProtect != "" {
				b.WriteString(DimStyle.Render("Tomorrow protect: "))
				b.WriteString(*m.day.TomorrowProtect)
			}
		}

		// Show after-hours section if any entries
		if len(afterHoursEntries) > 0 {
			b.WriteString("\n\n")
			b.WriteString(DimStyle.Render("═════════════════════════════════════"))
			b.WriteString("\n")
			b.WriteString(BoldStyle.Render("After-Hours"))
			b.WriteString("\n")
			b.WriteString(DimStyle.Render("═════════════════════════════════════"))
			b.WriteString("\n\n")
			b.WriteString(m.formatEntries(afterHoursEntries))
		}
	}

	return b.String()
}

// View renders the UI
func (m ViewModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	viewContent := m.viewport.View() + "\n"
	viewContent += DimStyle.Render("↑/↓ or j/k to scroll • Shift+R to toggle text • q/esc to exit")
	return viewContent
}

// splitEntries separates regular entries from after-hours entries
// Returns regular entries and after-hours entries
func splitEntries(entries []*database.Entry, dayCompleted bool) ([]*database.Entry, []*database.Entry) {
	if !dayCompleted {
		// If day not completed, all entries are regular
		return entries, nil
	}

	// Find last @signoff entry timestamp
	var signoffTime time.Time
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			if tag.TagValue == "@signoff" {
				signoffTime = entry.Timestamp
			}
		}
	}

	if signoffTime.IsZero() {
		// No signoff found (shouldn't happen if day.Completed is true, but handle it)
		return entries, nil
	}

	// Split entries
	var regular []*database.Entry
	var afterHours []*database.Entry

	for _, entry := range entries {
		if entry.Timestamp.After(signoffTime) {
			afterHours = append(afterHours, entry)
		} else {
			regular = append(regular, entry)
		}
	}

	return regular, afterHours
}

// formatEntries formats the list of entries for display
func (m ViewModel) formatEntries(entries []*database.Entry) string {
	var b strings.Builder

	for i, entry := range entries {
		if i > 0 {
			b.WriteString("\n")
		}

		// Time
		timeStr := entry.Timestamp.Format("3:04pm")
		b.WriteString(DimStyle.Render(timeStr))
		b.WriteString(" | ")

		// Entry text - truncate if collapsed and text is long
		entryText := entry.EntryText
		if !m.textExpanded && len(entryText) > m.maxCollapsedLen {
			// Truncate and add indicator
			entryText = entryText[:m.maxCollapsedLen] + DimStyle.Render("... [Shift+R to expand]")
		}
		b.WriteString(entryText)

		// Momentum
		if entry.Momentum != nil && *entry.Momentum != "" {
			b.WriteString(" ")
			b.WriteString(formatMomentum(*entry.Momentum))
		}

		// Tags
		if len(entry.Tags) > 0 {
			b.WriteString(" ")
			b.WriteString(formatTags(entry.Tags))
		}
	}

	return b.String()
}

// formatMomentum returns the visual representation of momentum
func formatMomentum(momentum string) string {
	switch momentum {
	case "up":
		return "↑"
	case "down":
		return "↓"
	case "neutral":
		return "→"
	case "back":
		return "←"
	default:
		return ""
	}
}

// formatTags formats tags for display
func formatTags(tags []database.Tag) string {
	var parts []string
	for _, tag := range tags {
		parts = append(parts, tag.TagValue)
	}
	return DimStyle.Render(strings.Join(parts, " "))
}
