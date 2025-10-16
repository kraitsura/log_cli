package tui

import (
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ListThoughtsModel is the model for viewing all thoughts for a day
type ListThoughtsModel struct {
	day      *database.Day
	thoughts []*database.Entry
	viewport viewport.Model
	ready    bool
}

// NewListThoughtsModel creates a new list thoughts model
func NewListThoughtsModel(day *database.Day, thoughts []*database.Entry) ListThoughtsModel {
	return ListThoughtsModel{
		day:      day,
		thoughts: thoughts,
	}
}

// Init initializes the model
func (m ListThoughtsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ListThoughtsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Initialize viewport on first window size message
			m.viewport = viewport.New(msg.Width, msg.Height-2)
			m.viewport.SetContent(m.generateContent())
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

// generateContent generates the thoughts list content
func (m ListThoughtsModel) generateContent() string {
	var b strings.Builder

	// Header with date
	dateStr := m.day.Date.Format("Monday, January 2, 2006")
	b.WriteString(RenderHeaderBar("THOUGHTS", dateStr))
	b.WriteString("\n\n")

	if len(m.thoughts) == 0 {
		b.WriteString(DimStyle.Render("No thoughts logged today yet!"))
		b.WriteString("\n\n")
		b.WriteString(DimStyle.Render("Type: log thought"))
	} else {
		for i, thought := range m.thoughts {
			if i > 0 {
				b.WriteString("\n\n")
			}

			// Time
			timeStr := thought.Timestamp.Format("3:04pm")
			b.WriteString(DimStyle.Render(timeStr))
			b.WriteString("\n")

			// Thought text (keep the 💭 prefix for display)
			thoughtText := thought.EntryText
			b.WriteString(thoughtText)
		}
	}

	return b.String()
}

// View renders the UI
func (m ListThoughtsModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	viewContent := m.viewport.View() + "\n"
	viewContent += DimStyle.Render("↑/↓ or j/k to scroll • q/esc to exit")
	return viewContent
}
