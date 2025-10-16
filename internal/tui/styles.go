package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// DAYLOG TUI DESIGN SYSTEM
// Precision with warmth - Clean, colorful, and functional
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

// Color Palette (Vibrant)
var (
	colorPrimary  = lipgloss.Color("#FF79C6") // Magenta for headers
	colorAccent   = lipgloss.Color("#8BE9FD") // Cyan for interactive
	colorSuccess  = lipgloss.Color("#50FA7B") // Green for wins/positive
	colorWarning  = lipgloss.Color("#FFB86C") // Orange for alerts
	colorError    = lipgloss.Color("#FF5555") // Red for errors
	colorMetadata = lipgloss.Color("#BD93F9") // Purple for labels
	colorDim      = lipgloss.Color("#6272A4") // Gray for helper text
	colorBorder   = lipgloss.Color("#44475A") // Subtle border
	colorText     = lipgloss.Color("#D0D0D0") // Softer off-white text (less harsh)
)

// Typography Styles
var (
	// Header for main titles
	HeaderStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	// Subheader for sections
	SubheaderStyle = lipgloss.NewStyle().
			Foreground(colorMetadata).
			Bold(true)

	// Bold text for emphasis
	BoldStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	// Regular body text
	BodyStyle = lipgloss.NewStyle().
			Foreground(colorText)

	// Dim text for helper/meta info
	DimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// Accent for interactive elements
	AccentStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	// Success for wins/positive
	SuccessStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	// Warning for alerts/drift
	WarningStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)

	// Error for problems
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	// Metadata labels (purple, subtle)
	MetadataStyle = lipgloss.NewStyle().
			Foreground(colorMetadata)

	// Prompt style (cyan, friendly)
	PromptStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	// Alert style (alias for warning)
	AlertStyle = WarningStyle

	// Selected item style (cyan, bold for list selections)
	SelectedStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	// Pattern-specific color styles (for weekly review)
	FlowStyle  = SuccessStyle // Green for [FLOW] patterns
	GoldStyle  = AccentStyle  // Cyan for [GOLD] patterns
	StuckStyle = WarningStyle // Orange for [STUCK] patterns
	LeakStyle  = ErrorStyle   // Red for [LEAK] patterns
)

// Box Styles
var (
	// Main container with borders (no fixed width - adapts to terminal)
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(2, 3)

	// Input field style (neutral)
	InputStyle = lipgloss.NewStyle().
			Foreground(colorText)

	// Focused input highlight (cyan glow)
	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(colorAccent)
)

// Helper Functions
// RenderHeaderBar creates a styled header bar with title and date
func RenderHeaderBar(title, date string) string {
	// Format: ┌─ TITLE ──────────────── DATE ─┐
	titlePart := "─ " + title + " "
	datePart := " " + date + " ─"

	// Calculate dashes needed to fill width
	totalLen := len(titlePart) + len(datePart) + 2 // +2 for corners
	dashesNeeded := 80 - totalLen
	if dashesNeeded < 0 {
		dashesNeeded = 0
	}

	dashes := strings.Repeat("─", dashesNeeded)
	return HeaderStyle.Render("┌" + titlePart + dashes + datePart + "┐")
}

// RenderDivider creates a horizontal divider line
func RenderDivider() string {
	return DimStyle.Render(strings.Repeat("─", 76))
}

// RenderDoubleDivider creates a bold horizontal divider
func RenderDoubleDivider() string {
	return DimStyle.Render(strings.Repeat("━", 76))
}

// FormatTime formats time for display (e.g., "3:04pm")
func FormatTime(t interface{}) string {
	// This will be used in the TUI files
	return fmt.Sprintf("%v", t)
}
