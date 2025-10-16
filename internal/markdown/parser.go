package markdown

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aaryareddy/log_cli/internal/database"
)

// Parser handles parsing markdown files back into database models
type Parser struct {
	// Compiled regex patterns for efficiency
	titlePattern     *regexp.Regexp
	intentionPattern *regexp.Regexp
	entryPattern     *regexp.Regexp
	momentumPattern  *regexp.Regexp
	tagPattern       *regexp.Regexp
	reflectionStart  *regexp.Regexp
}

// NewParser creates a new markdown parser
func NewParser() *Parser {
	return &Parser{
		titlePattern:     regexp.MustCompile(`^# DAYLOG - (.+)$`),
		intentionPattern: regexp.MustCompile(`^\*\*Intention:\*\* (.+)$`),
		entryPattern:     regexp.MustCompile(`^- (\d{1,2}:\d{2}(?:am|pm)) \| (.+)$`),
		momentumPattern:  regexp.MustCompile(`(↑|↓|→|←)$`),
		tagPattern:       regexp.MustCompile(`(@\w+|\[[\w\s]+\])(?:\s|$)`),
		reflectionStart:  regexp.MustCompile(`^\*\*Reflection:\*\*$`),
	}
}

// ParseFile reads a markdown file and returns the Day and Entry structs
func (p *Parser) ParseFile(filePath string) (*database.Day, []*database.Entry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Extract date from filename (YYYY-MM-DD.md format)
	filename := filepath.Base(filePath)
	dateStr := strings.TrimSuffix(filename, ".md")
	dayDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid date in filename: %w", err)
	}

	// Initialize Day struct
	day := &database.Day{
		Date: dayDate,
	}

	var entries []*database.Entry
	scanner := bufio.NewScanner(file)
	inReflectionSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and separators
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "---") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "**Reflection:**") {
			inReflectionSection = true
			day.Completed = true
			continue
		}

		if strings.HasPrefix(line, "**After-Hours:**") {
			inReflectionSection = false
			continue
		}

		// Parse title line (extract date validation)
		if match := p.titlePattern.FindStringSubmatch(line); match != nil {
			// Title line matched - date already extracted from filename
			continue
		}

		// Parse intention
		if match := p.intentionPattern.FindStringSubmatch(line); match != nil {
			intention := match[1]
			day.Intention = &intention
			continue
		}

		// Parse old-style win (separate section)
		if strings.HasPrefix(line, "**Win:**") {
			win := strings.TrimPrefix(line, "**Win:** ")
			day.Win = &win
			continue
		}

		// Parse reflection section
		if inReflectionSection {
			if strings.HasPrefix(line, "- Pulled off track: ") {
				reflection := strings.TrimPrefix(line, "- Pulled off track: ")
				day.PulledOffTrack = &reflection
				continue
			}
			if strings.HasPrefix(line, "- Kept on track: ") {
				reflection := strings.TrimPrefix(line, "- Kept on track: ")
				day.KeptOnTrack = &reflection
				continue
			}
			if strings.HasPrefix(line, "- Tomorrow protect: ") {
				reflection := strings.TrimPrefix(line, "- Tomorrow protect: ")
				day.TomorrowProtect = &reflection
				continue
			}
		}

		// Parse entry lines
		if match := p.entryPattern.FindStringSubmatch(line); match != nil {
			timeStr := match[1]
			entryText := match[2]

			// Parse timestamp
			entryTime, err := p.parseTimeWithDate(timeStr, dayDate)
			if err != nil {
				// Skip invalid entries
				continue
			}

			// Extract momentum
			var momentum *string
			if p.momentumPattern.MatchString(entryText) {
				// Find last character
				lastChar := entryText[len(entryText)-1:]
				entryText = strings.TrimSpace(entryText[:len(entryText)-1])

				var mom string
				switch lastChar {
				case "↑":
					mom = "up"
				case "↓":
					mom = "down"
				case "→":
					mom = "neutral"
				case "←":
					mom = "back"
				}
				momentum = &mom
			}

			// Extract tags
			tags := p.extractTags(entryText)

			// Remove tags from text
			cleanText := p.stripTags(entryText)

			entry := &database.Entry{
				Timestamp: entryTime,
				EntryText: cleanText,
				Momentum:  momentum,
				Tags:      tags,
			}

			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}

	return day, entries, nil
}

// parseTimeWithDate combines a time string (3:04pm) with a date
func (p *Parser) parseTimeWithDate(timeStr string, date time.Time) (time.Time, error) {
	// Parse time in 12-hour format
	t, err := time.Parse("3:04pm", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	// Combine with date
	timestamp := time.Date(
		date.Year(), date.Month(), date.Day(),
		t.Hour(), t.Minute(), 0, 0,
		time.Local,
	)

	return timestamp, nil
}

// extractTags finds all tags in the entry text
func (p *Parser) extractTags(text string) []database.Tag {
	var tags []database.Tag

	matches := p.tagPattern.FindAllString(text, -1)
	for _, match := range matches {
		tagValue := strings.TrimSpace(match)

		var tagType string
		if strings.HasPrefix(tagValue, "@") {
			tagType = "context"
		} else if strings.HasPrefix(tagValue, "[") {
			tagType = "flag"
		} else {
			continue
		}

		tags = append(tags, database.Tag{
			TagType:  tagType,
			TagValue: tagValue,
		})
	}

	return tags
}

// stripTags removes all tags from the entry text
func (p *Parser) stripTags(text string) string {
	return strings.TrimSpace(p.tagPattern.ReplaceAllString(text, ""))
}
