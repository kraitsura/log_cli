package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// AutocompleteState tracks the state of autocomplete suggestions
type AutocompleteState struct {
	Active         bool
	TriggerChar    string   // "@" or "["
	TriggerPos     int      // Position in text where trigger was typed
	FilterText     string   // Text after trigger for filtering
	Suggestions    []string // Filtered list of suggestions
	SelectedIndex  int      // Currently selected suggestion
	AllSuggestions map[string][]string
}

// NewAutocompleteState creates a new autocomplete state
func NewAutocompleteState() AutocompleteState {
	return AutocompleteState{
		Active:        false,
		SelectedIndex: 0,
		AllSuggestions: map[string][]string{
			"@": {"@deep", "@social", "@admin", "@break", "@zone", "@signoff"},
			"[": {"[LEAK]", "[FLOW]", "[STUCK]", "[GOLD]", "[DRIFT]", "[ANCHOR]"},
		},
	}
}

// Update updates autocomplete state based on current text and cursor position
func (a *AutocompleteState) Update(text string, cursorPos int) {
	// Find if there's a trigger character before cursor
	triggerPos := -1
	var triggerChar string

	// Look backwards from cursor to find trigger
	for i := cursorPos - 1; i >= 0; i-- {
		ch := string(text[i])
		if ch == "@" || ch == "[" {
			triggerChar = ch
			triggerPos = i
			break
		}
		// Stop if we hit whitespace or another trigger
		if ch == " " || ch == "@" || ch == "[" {
			break
		}
	}

	// If no trigger found or trigger is not immediately accessible, deactivate
	if triggerPos == -1 {
		a.Active = false
		return
	}

	// Check if there's text between trigger and cursor that would break autocomplete
	// (e.g., closing bracket for [ or space)
	textBetween := text[triggerPos:cursorPos]
	if triggerChar == "[" && strings.Contains(textBetween, "]") {
		a.Active = false
		return
	}
	if strings.Contains(textBetween[1:], " ") {
		a.Active = false
		return
	}

	// Activate autocomplete
	a.Active = true
	a.TriggerChar = triggerChar
	a.TriggerPos = triggerPos
	a.FilterText = text[triggerPos+1 : cursorPos]

	// Filter suggestions
	a.Suggestions = a.filterSuggestions()

	// Reset selected index if suggestions changed
	if a.SelectedIndex >= len(a.Suggestions) {
		a.SelectedIndex = 0
	}
}

// filterSuggestions returns suggestions matching the current filter text
func (a *AutocompleteState) filterSuggestions() []string {
	allSuggestions := a.AllSuggestions[a.TriggerChar]
	if a.FilterText == "" {
		return allSuggestions
	}

	var filtered []string
	lowerFilter := strings.ToLower(a.FilterText)

	for _, suggestion := range allSuggestions {
		// For context tags, check if it starts with @filterText
		if a.TriggerChar == "@" {
			if strings.HasPrefix(strings.ToLower(suggestion), "@"+lowerFilter) {
				filtered = append(filtered, suggestion)
			}
		} else if a.TriggerChar == "[" {
			// For flag tags, check if it starts with [FILTERTEXT
			if strings.HasPrefix(strings.ToLower(suggestion), "["+lowerFilter) {
				filtered = append(filtered, suggestion)
			}
		}
	}

	return filtered
}

// MoveUp moves selection up in the list
func (a *AutocompleteState) MoveUp() {
	if !a.Active || len(a.Suggestions) == 0 {
		return
	}
	a.SelectedIndex--
	if a.SelectedIndex < 0 {
		a.SelectedIndex = len(a.Suggestions) - 1
	}
}

// MoveDown moves selection down in the list
func (a *AutocompleteState) MoveDown() {
	if !a.Active || len(a.Suggestions) == 0 {
		return
	}
	a.SelectedIndex++
	if a.SelectedIndex >= len(a.Suggestions) {
		a.SelectedIndex = 0
	}
}

// GetSelectedSuggestion returns the currently selected suggestion
func (a *AutocompleteState) GetSelectedSuggestion() string {
	if !a.Active || len(a.Suggestions) == 0 {
		return ""
	}
	if a.SelectedIndex < 0 || a.SelectedIndex >= len(a.Suggestions) {
		return ""
	}
	return a.Suggestions[a.SelectedIndex]
}

// InsertSuggestion returns the text with the selected suggestion inserted
func (a *AutocompleteState) InsertSuggestion(currentText string) string {
	if !a.Active {
		return currentText
	}

	suggestion := a.GetSelectedSuggestion()
	if suggestion == "" {
		return currentText
	}

	// Find cursor position (end of filter text)
	cursorPos := a.TriggerPos + 1 + len(a.FilterText)

	// Build new text: before trigger + suggestion + space + after cursor
	before := currentText[:a.TriggerPos]
	after := currentText[cursorPos:]

	// Add space after suggestion if not already present
	if !strings.HasPrefix(after, " ") && after != "" {
		suggestion += " "
	}

	return before + suggestion + after
}

// Deactivate turns off autocomplete
func (a *AutocompleteState) Deactivate() {
	a.Active = false
}

// View renders the autocomplete dropdown
func (a *AutocompleteState) View() string {
	if !a.Active || len(a.Suggestions) == 0 {
		return ""
	}

	var b strings.Builder

	// Style for dropdown
	dropdownStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		MaxWidth(30)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("238"))

	// Build dropdown content
	for i, suggestion := range a.Suggestions {
		if i == a.SelectedIndex {
			b.WriteString(selectedItemStyle.Render("▶ " + suggestion))
		} else {
			b.WriteString(itemStyle.Render("  " + suggestion))
		}
		if i < len(a.Suggestions)-1 {
			b.WriteString("\n")
		}
	}

	return dropdownStyle.Render(b.String())
}
