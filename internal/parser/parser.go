package parser

import (
	"regexp"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
)

// ParseEntry parses an entry text and extracts momentum, tags, and clean text
func ParseEntry(text string) (cleanText string, momentum *string, tags []database.Tag) {
	// Extract momentum first
	momentum, remaining := extractMomentum(text)

	// Extract tags from remaining text
	tags, cleanText = extractTags(remaining)

	// Trim whitespace from clean text
	cleanText = strings.TrimSpace(cleanText)

	return cleanText, momentum, tags
}

// extractMomentum finds and extracts momentum markers (↑, ↓, →)
func extractMomentum(text string) (momentum *string, remaining string) {
	momentumRegex := regexp.MustCompile(`[↑↓→]`)

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
		}
		momentum = &m

		// Remove momentum marker from text
		remaining = momentumRegex.ReplaceAllString(text, "")
	} else {
		remaining = text
	}

	return momentum, remaining
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
