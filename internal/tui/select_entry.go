package tui

import (
	"fmt"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectEntryModel is the model for selecting an entry from a list
type SelectEntryModel struct {
	entries       []*database.Entry
	selectedIndex int
	confirmed     bool
	cancelled     bool
	width         int
	height        int
	title         string // e.g., "SELECT ENTRY TO DELETE"
}

// NewSelectEntryModel creates a new select entry model
func NewSelectEntryModel(entries []*database.Entry, title string) SelectEntryModel {
	return SelectEntryModel{
		entries:       entries,
		selectedIndex: len(entries) - 1, // Start at most recent (bottom)
		confirmed:     false,
		cancelled:     false,
		title:         title,
	}
}

// Init initializes the model
func (m SelectEntryModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m SelectEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			// Move down (toward more recent entries)
			if m.selectedIndex < len(m.entries)-1 {
				m.selectedIndex++
			}
		case "k", "up":
			// Move up (toward older entries)
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "enter":
			// Confirm selection
			m.confirmed = true
			return m, tea.Quit
		case "esc", "ctrl+c":
			// Cancel
			m.cancelled = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the selection UI
func (m SelectEntryModel) View() string {
	if m.confirmed || m.cancelled {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(HeaderStyle.Render(m.title))
	b.WriteString("\n\n")

	// Entry list
	for i, entry := range m.entries {
		entryNum := i + 1
		timeStr := entry.Timestamp.Format("3:04pm")

		// Build entry display text
		entryText := entry.EntryText

		// Add momentum
		if entry.Momentum != nil && *entry.Momentum != "" {
			entryText += " " + formatMomentum(*entry.Momentum)
		}

		// Add tags
		if len(entry.Tags) > 0 {
			tagStrs := []string{}
			for _, tag := range entry.Tags {
				tagStrs = append(tagStrs, tag.TagValue)
			}
			entryText += " " + strings.Join(tagStrs, " ")
		}

		// Truncate if too long
		maxLen := 60
		if len(entryText) > maxLen {
			entryText = entryText[:maxLen-3] + "..."
		}

		// Format line
		line := fmt.Sprintf("%2d. %s | %s", entryNum, timeStr, entryText)

		// Highlight selected
		if i == m.selectedIndex {
			b.WriteString(SelectedStyle.Render("› " + line))
		} else {
			b.WriteString(DimStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(DimStyle.Render("j/k or ↑/↓ to navigate • Enter to select • Esc to cancel"))

	return BoxStyle.Render(b.String())
}

// GetSelectedEntry returns the selected entry
func (m SelectEntryModel) GetSelectedEntry() *database.Entry {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.entries) {
		return m.entries[m.selectedIndex]
	}
	return nil
}

// GetSelectedIndex returns the 1-indexed position of the selected entry
func (m SelectEntryModel) GetSelectedIndex() int {
	return m.selectedIndex + 1
}

// WasConfirmed returns whether an entry was selected
func (m SelectEntryModel) WasConfirmed() bool {
	return m.confirmed
}

// WasCancelled returns whether the selection was cancelled
func (m SelectEntryModel) WasCancelled() bool {
	return m.cancelled
}
