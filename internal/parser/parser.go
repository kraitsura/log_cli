package parser

import (
	"regexp"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
)

// ParseEntry parses an entry text and extracts momentum, tags, and clean text
func ParseEntry(text string) (cleanText string, momentum *string, tags []database.Tag) {
	// Convert shortcuts to arrows first
	text = convertMomentumShortcuts(text)

	// Extract momentum
	momentum, remaining := extractMomentum(text)

	// Extract tags from remaining text
	tags, cleanText = extractTags(remaining)

	// Trim whitespace from clean text
	cleanText = strings.TrimSpace(cleanText)

	return cleanText, momentum, tags
}

// convertMomentumShortcuts converts text shortcuts to arrow symbols
// Only converts double characters to preserve single characters for normal text use
// Supports: ++ → ↑, -- → ↓, == → →, << → ←, also -> and <-
func convertMomentumShortcuts(text string) string {
	// Replace shortcuts with arrows (order matters - replace longer patterns first)
	replacements := map[string]string{
		"++": "↑",
		"--": "↓",
		"==": "→",
		"<<": "←",
		"->": "→",
		"<-": "←",
	}

	// Apply all replacements
	for shortcut, arrow := range replacements {
		text = strings.ReplaceAll(text, shortcut, arrow)
	}

	return text
}

// extractMomentum finds and extracts momentum markers (↑, ↓, →, ←)
func extractMomentum(text string) (momentum *string, remaining string) {
	momentumRegex := regexp.MustCompile(`[↑↓→←]`)

	match := momentumRegex.FindString(text)
	if match != "" {
		var m string
		switch match {
		case "↑":
			m = "up"
		case "↓":
			m = "down"
		case "→":
			m = "neutral"
		case "←":
			m = "back"
		}
		momentum = &m

		// Remove momentum marker from text
		remaining = momentumRegex.ReplaceAllString(text, "")
	} else {
		remaining = text
	}

	return momentum, remaining
}

// ReconstructEntryText reconstructs the original entry text from parsed components
// This is used for editing entries
func ReconstructEntryText(cleanText string, momentum *string, tags []database.Tag) string {
	var parts []string

	parts = append(parts, cleanText)

	// Add momentum marker
	if momentum != nil {
		switch *momentum {
		case "up":
			parts = append(parts, "↑")
		case "down":
			parts = append(parts, "↓")
		case "neutral":
			parts = append(parts, "→")
		case "back":
			parts = append(parts, "←")
		}
	}

	// Add tags
	for _, tag := range tags {
		parts = append(parts, tag.TagValue)
	}

	return strings.Join(parts, " ")
}

// extractTags finds and extracts context tags (@word) and flag tags ([WORD])
func extractTags(text string) (tags []database.Tag, remaining string) {
	remaining = text

	// Extract context tags (@word)
	contextRegex := regexp.MustCompile(`@(deep|social|admin|break|zone|signoff)\b`)
	contextMatches := contextRegex.FindAllStringSubmatch(text, -1)

	for _, match := range contextMatches {
		tags = append(tags, database.Tag{
			TagType:  "context",
			TagValue: "@" + match[1],
		})
	}

	// Remove context tags from text
	remaining = contextRegex.ReplaceAllString(remaining, "")

	// Extract flag tags ([WORD])
	flagRegex := regexp.MustCompile(`\[(LEAK|FLOW|STUCK|GOLD|DRIFT|ANCHOR[^\]]*)\]`)
	flagMatches := flagRegex.FindAllStringSubmatch(remaining, -1)

	for _, match := range flagMatches {
		tags = append(tags, database.Tag{
			TagType:  "flag",
			TagValue: "[" + match[1] + "]",
		})
	}

	// Remove flag tags from text
	remaining = flagRegex.ReplaceAllString(remaining, "")

	// Clean up extra whitespace
	remaining = strings.Join(strings.Fields(remaining), " ")

	return tags, remaining
}
