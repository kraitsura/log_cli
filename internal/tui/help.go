package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// HelpModel is the model for the help screen
type HelpModel struct {
	viewport viewport.Model
	ready    bool
	content  string
}

// NewHelpModel creates a new help model
func NewHelpModel() HelpModel {
	return HelpModel{
		content: generateHelpContent(),
	}
}

// Init initializes the model
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// generateHelpContent generates the help text content
func generateHelpContent() string {
	var b strings.Builder

	b.WriteString("DAYLOG - HELP\n")
	b.WriteString("═════════════\n\n")

	// Commands section
	b.WriteString("COMMANDS\n")
	b.WriteString("  log              Create a new log entry\n")
	b.WriteString("  log view         Display today's log entries\n")
	b.WriteString("  log stats        Show weekly statistics\n")
	b.WriteString("  log edit [n]     Edit most recent entry (or entry #n)\n")
	b.WriteString("  log delete [n]   Delete most recent entry (or entry #n)\n")
	b.WriteString("  log win          Quickly log a win\n")
	b.WriteString("  log thought      Log a quick thought (no tags/momentum)\n")
	b.WriteString("  log help         Show this help screen\n")
	b.WriteString("\n")

	// Momentum markers section
	b.WriteString("MOMENTUM MARKERS\n")
	b.WriteString("  ↑ or ++   Productive/energized\n")
	b.WriteString("  → or ==   Neutral/coasting\n")
	b.WriteString("  ↓ or --   Dragging/unfocused\n")
	b.WriteString("  ← or <<   Waste/destructive action\n")
	b.WriteString("  (Also: -> for →, <- for ←)\n")
	b.WriteString("\n")

	// Context tags section
	b.WriteString("CONTEXT TAGS\n")
	b.WriteString("  @deep     Deep focused work\n")
	b.WriteString("  @social   Meetings, calls, collaboration\n")
	b.WriteString("  @admin    Email, scheduling, life tasks\n")
	b.WriteString("  @break    Intentional rest\n")
	b.WriteString("  @zone     Creative/flow work\n")
	b.WriteString("  @signoff  End of day marker (triggers reflection)\n")
	b.WriteString("\n")

	// Pattern flags section
	b.WriteString("PATTERN FLAGS\n")
	b.WriteString("  [LEAK]    Time drains (social media, news spirals)\n")
	b.WriteString("  [FLOW]    In the zone, highly productive\n")
	b.WriteString("  [STUCK]   Spinning wheels, unclear what to do\n")
	b.WriteString("  [GOLD]    Unusually productive periods\n")
	b.WriteString("  [DRIFT]   More than 90 minutes without logging\n")
	b.WriteString("  [ANCHOR]  Non-negotiable check-in points\n")
	b.WriteString("\n")

	// Examples section
	b.WriteString("EXAMPLES\n")
	b.WriteString("  log\n")
	b.WriteString("  > Writing blog post ++ @deep\n")
	b.WriteString("\n")
	b.WriteString("  log\n")
	b.WriteString("  > Scrolling Twitter [LEAK] --\n")
	b.WriteString("\n")
	b.WriteString("  log\n")
	b.WriteString("  > Finished proposal draft! [FLOW] ++ @zone\n")

	return b.String()
}

// View renders the help screen
func (m HelpModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	helpText := m.viewport.View() + "\n"
	helpText += DimStyle.Render("↑/↓ or j/k to scroll • q/esc to exit")
	return helpText
}
