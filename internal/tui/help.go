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

	// Header
	b.WriteString(HeaderStyle.Render("DAYLOG - HELP"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", 78))
	b.WriteString("\n\n")

	// Commands section
	b.WriteString(SubheaderStyle.Render("COMMANDS"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")
	b.WriteString(MetadataStyle.Render("  log              "))
	b.WriteString("Create a new log entry\n")
	b.WriteString(MetadataStyle.Render("  log view         "))
	b.WriteString("Display today's log entries\n")
	b.WriteString(MetadataStyle.Render("  log stats        "))
	b.WriteString("Show weekly statistics\n")
	b.WriteString(MetadataStyle.Render("  log edit [n]     "))
	b.WriteString("Edit most recent entry (or entry #n)\n")
	b.WriteString(MetadataStyle.Render("  log delete [n]   "))
	b.WriteString("Delete most recent entry (or entry #n)\n")
	b.WriteString(MetadataStyle.Render("  log win          "))
	b.WriteString("Quickly log a win\n")
	b.WriteString(MetadataStyle.Render("  log thought      "))
	b.WriteString("Log a quick thought (no tags/momentum)\n")
	b.WriteString(MetadataStyle.Render("  log help         "))
	b.WriteString("Show this help screen\n")
	b.WriteString("\n")

	// Momentum markers section
	b.WriteString(SubheaderStyle.Render("MOMENTUM MARKERS"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")
	b.WriteString(MetadataStyle.Render("  ↑  or  ++        "))
	b.WriteString("Productive/energized\n")
	b.WriteString(MetadataStyle.Render("  →  or  ==        "))
	b.WriteString("Neutral/coasting\n")
	b.WriteString(MetadataStyle.Render("  ↓  or  --        "))
	b.WriteString("Dragging/unfocused\n")
	b.WriteString(MetadataStyle.Render("  ←  or  <<        "))
	b.WriteString("Waste/destructive action\n")
	b.WriteString(DimStyle.Render("  (Also: -> for →, <- for ←)"))
	b.WriteString("\n\n")

	// Context tags section
	b.WriteString(SubheaderStyle.Render("CONTEXT TAGS"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")
	b.WriteString(MetadataStyle.Render("  @deep            "))
	b.WriteString("Deep focused work\n")
	b.WriteString(MetadataStyle.Render("  @social          "))
	b.WriteString("Meetings, calls, collaboration\n")
	b.WriteString(MetadataStyle.Render("  @admin           "))
	b.WriteString("Email, scheduling, life tasks\n")
	b.WriteString(MetadataStyle.Render("  @break           "))
	b.WriteString("Intentional rest\n")
	b.WriteString(MetadataStyle.Render("  @zone            "))
	b.WriteString("Creative/flow work\n")
	b.WriteString(MetadataStyle.Render("  @signoff         "))
	b.WriteString("End of day marker (triggers reflection)\n")
	b.WriteString("\n")

	// Pattern flags section
	b.WriteString(SubheaderStyle.Render("PATTERN FLAGS"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")
	b.WriteString(MetadataStyle.Render("  [LEAK]           "))
	b.WriteString("Time drains (social media, news spirals)\n")
	b.WriteString(MetadataStyle.Render("  [FLOW]           "))
	b.WriteString("In the zone, highly productive\n")
	b.WriteString(MetadataStyle.Render("  [STUCK]          "))
	b.WriteString("Spinning wheels, unclear what to do\n")
	b.WriteString(MetadataStyle.Render("  [GOLD]           "))
	b.WriteString("Unusually productive periods\n")
	b.WriteString(MetadataStyle.Render("  [DRIFT]          "))
	b.WriteString("More than 90 minutes without logging\n")
	b.WriteString(MetadataStyle.Render("  [ANCHOR]         "))
	b.WriteString("Non-negotiable check-in points\n")
	b.WriteString("\n")

	// Examples section
	b.WriteString(SubheaderStyle.Render("EXAMPLES"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("  $ log"))
	b.WriteString("\n")
	b.WriteString(AccentStyle.Render("  > "))
	b.WriteString("Writing blog post ++ @deep\n")
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("  $ log"))
	b.WriteString("\n")
	b.WriteString(AccentStyle.Render("  > "))
	b.WriteString("Scrolling Twitter [LEAK] --\n")
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("  $ log"))
	b.WriteString("\n")
	b.WriteString(AccentStyle.Render("  > "))
	b.WriteString("Finished proposal draft! [FLOW] ++ @zone\n")

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
