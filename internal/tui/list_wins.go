package tui

import (
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ListWinsModel is the model for viewing all wins for a day
type ListWinsModel struct {
	day      *database.Day
	wins     []*database.Entry
	viewport viewport.Model
	ready    bool
}

// NewListWinsModel creates a new list wins model
func NewListWinsModel(day *database.Day, wins []*database.Entry) ListWinsModel {
	return ListWinsModel{
		day:  day,
		wins: wins,
	}
}

// Init initializes the model
func (m ListWinsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ListWinsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

// generateContent generates the wins list content
func (m ListWinsModel) generateContent() string {
	var b strings.Builder

	// Header with date
	dateStr := m.day.Date.Format("Monday, January 2, 2006")
	b.WriteString(RenderHeaderBar("WINS", dateStr))
	b.WriteString("\n\n")

	if len(m.wins) == 0 {
		b.WriteString(DimStyle.Render("No wins logged today yet!"))
		b.WriteString("\n\n")
		b.WriteString(DimStyle.Render("Type: log win"))
	} else {
		for i, win := range m.wins {
			if i > 0 {
				b.WriteString("\n")
			}

			// Time
			timeStr := win.Timestamp.Format("3:04pm")
			b.WriteString(DimStyle.Render(timeStr))
			b.WriteString(" | ")

			// Win text (remove the 🌟 prefix for display)
			winText := win.EntryText
			winText = strings.TrimPrefix(winText, "🌟 ")
			b.WriteString(SuccessStyle.Render(winText))
			b.WriteString(" 🌟")
		}
	}

	return b.String()
}

// View renders the UI
func (m ListWinsModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	viewContent := m.viewport.View() + "\n"
	viewContent += DimStyle.Render("↑/↓ or j/k to scroll • q/esc to exit")
	return viewContent
}
