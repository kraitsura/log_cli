package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmDeleteModel is the model for the delete confirmation screen
type ConfirmDeleteModel struct {
	entryIndex int
	timestamp  time.Time
	entryText  string
	confirmed  bool
	cancelled  bool
	width      int
	height     int
}

// NewConfirmDeleteModel creates a new delete confirmation model
func NewConfirmDeleteModel(entryIndex int, timestamp time.Time, entryText string) ConfirmDeleteModel {
	return ConfirmDeleteModel{
		entryIndex: entryIndex,
		timestamp:  timestamp,
		entryText:  entryText,
		confirmed:  false,
		cancelled:  false,
	}
}

// Init initializes the model
func (m ConfirmDeleteModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ConfirmDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			return m, tea.Quit
		case "n", "N", "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the confirmation prompt
func (m ConfirmDeleteModel) View() string {
	if m.confirmed || m.cancelled {
		return ""
	}

	var b strings.Builder

	b.WriteString(AlertStyle.Render("⚠️  DELETE ENTRY"))
	b.WriteString("\n\n")

	info := fmt.Sprintf("Entry #%d from %s:", m.entryIndex, m.timestamp.Format("3:04pm"))
	b.WriteString(info)
	b.WriteString("\n\n")

	// Calculate responsive width for entry preview (85% of terminal width, or 60 chars default)
	previewWidth := 60
	if m.width > 0 {
		previewWidth = int(float64(m.width) * 0.85)
		if previewWidth < 40 {
			previewWidth = 40 // minimum width
		}
		if previewWidth > 80 {
			previewWidth = 80 // maximum width for readability
		}
	}

	// Apply width constraint and word wrap to entry text
	entryStyle := lipgloss.NewStyle().
		Width(previewWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#44475A")).
		Padding(1, 2)

	entryPreview := entryStyle.Render(m.entryText)
	b.WriteString(entryPreview)
	b.WriteString("\n\n")

	b.WriteString(BoldStyle.Render("Are you sure you want to delete this entry?"))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("[Y]es to delete • [N]o to cancel"))

	return BoxStyle.Render(b.String())
}

// WasConfirmed returns whether the deletion was confirmed
func (m ConfirmDeleteModel) WasConfirmed() bool {
	return m.confirmed
}

// WasCancelled returns whether the deletion was cancelled
func (m ConfirmDeleteModel) WasCancelled() bool {
	return m.cancelled
}
