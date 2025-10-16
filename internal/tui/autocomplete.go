package tui

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// AutocompleteState tracks the state of autocomplete suggestions
type AutocompleteState struct {
	Active         bool
	TriggerChar    string   // "@" or "["
	TriggerPos     int      // Position in text where trigger was typed (in runes)
	TriggerPosBytes int     // Position in bytes for string slicing
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

// runeIndexToByte converts a rune index to a byte index in a string
func runeIndexToByte(s string, runeIndex int) int {
	if runeIndex <= 0 {
		return 0
	}
	byteIndex := 0
	for i := 0; i < runeIndex && byteIndex < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[byteIndex:])
		byteIndex += size
	}
	return byteIndex
}

// byteIndexToRune converts a byte index to a rune index in a string
func byteIndexToRune(s string, byteIndex int) int {
	if byteIndex <= 0 {
		return 0
	}
	return utf8.RuneCountInString(s[:byteIndex])
}

// Update updates autocomplete state based on current text and cursor position (in runes)
func (a *AutocompleteState) Update(text string, cursorPosRunes int) {
	// Convert rune position to byte position for string operations
	cursorPosBytes := runeIndexToByte(text, cursorPosRunes)

	// Find if there's a trigger character before cursor
	triggerPosBytes := -1
	var triggerChar string

	// Look backwards from cursor (in bytes) to find trigger
	i := cursorPosBytes
	for i > 0 {
		// Decode the previous rune
		r, size := utf8.DecodeLastRuneInString(text[:i])
		i -= size

		ch := string(r)
		if ch == "@" || ch == "[" {
			triggerChar = ch
			triggerPosBytes = i
			break
		}
		// Stop if we hit whitespace
		if ch == " " {
			break
		}
	}

	// If no trigger found, deactivate
	if triggerPosBytes == -1 {
		a.Active = false
		return
	}

	// Check if there's text between trigger and cursor that would break autocomplete
	// (e.g., closing bracket for [ or space)
	textBetween := text[triggerPosBytes:cursorPosBytes]
	if triggerChar == "[" && strings.Contains(textBetween, "]") {
		a.Active = false
		return
	}
	if len(textBetween) > 1 && strings.Contains(textBetween[1:], " ") {
		a.Active = false
		return
	}

	// Activate autocomplete
	a.Active = true
	a.TriggerChar = triggerChar
	a.TriggerPosBytes = triggerPosBytes
	a.TriggerPos = byteIndexToRune(text, triggerPosBytes)

	// Extract filter text (skip the trigger character itself)
	triggerSize := len(triggerChar)
	a.FilterText = text[triggerPosBytes+triggerSize : cursorPosBytes]

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

// InsertSuggestion returns the text with the selected suggestion inserted and the new cursor position (in runes)
func (a *AutocompleteState) InsertSuggestion(currentText string) (string, int) {
	if !a.Active {
		return currentText, utf8.RuneCountInString(currentText)
	}

	suggestion := a.GetSelectedSuggestion()
	if suggestion == "" {
		return currentText, utf8.RuneCountInString(currentText)
	}

	// Use byte positions for string slicing
	// Find cursor position in bytes (end of filter text)
	filterTextBytes := len(a.FilterText)
	triggerSizeBytes := len(a.TriggerChar)
	cursorPosBytes := a.TriggerPosBytes + triggerSizeBytes + filterTextBytes

	// Build new text: before trigger + suggestion + space + after cursor
	before := currentText[:a.TriggerPosBytes]
	after := currentText[cursorPosBytes:]

	// Add space after suggestion if not already present
	newText := before + suggestion
	addedSpace := false
	if !strings.HasPrefix(after, " ") && after != "" {
		newText += " "
		addedSpace = true
	}
	newText += after

	// Calculate new cursor position in RUNES (for textinput)
	// Count runes in: before + suggestion + optional space
	beforeRunes := utf8.RuneCountInString(before)
	suggestionRunes := utf8.RuneCountInString(suggestion)
	newCursorPosRunes := beforeRunes + suggestionRunes
	if addedSpace {
		newCursorPosRunes++ // Add 1 for the space
	}

	return newText, newCursorPosRunes
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

	// Minimalistic styles using project colors
	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")) // Dim gray

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")). // Cyan accent
		Bold(true)

	// Build dropdown content - clean, no border
	for i, suggestion := range a.Suggestions {
		if i == a.SelectedIndex {
			b.WriteString(selectedItemStyle.Render("› " + suggestion))
		} else {
			b.WriteString(itemStyle.Render("  " + suggestion))
		}
		if i < len(a.Suggestions)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
