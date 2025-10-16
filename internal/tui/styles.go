package tui

import "github.com/charmbracelet/lipgloss"

var (
	// BoxStyle is the main container style
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(80)

	// HeaderStyle is for main headers
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	// DimStyle is for helper text and labels
	DimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// BoldStyle is for emphasis
	BoldStyle = lipgloss.NewStyle().
		Bold(true)

	// AlertStyle is for warnings and drift alerts
	AlertStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")).
		Bold(true)

	// SuccessStyle is for positive messages
	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	// ErrorStyle is for error messages
	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	// PromptStyle is for input prompts
	PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true)
)
