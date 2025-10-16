package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
)

// Writer handles markdown file generation
type Writer struct {
	outputDir string
}

// NewWriter creates a new markdown writer
func NewWriter(outputDir string) (*Writer, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(outputDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		outputDir = filepath.Join(homeDir, outputDir[1:])
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &Writer{outputDir: outputDir}, nil
}

// AppendEntry appends an entry to the day's markdown file
func (w *Writer) AppendEntry(day *database.Day, entry *database.Entry) error {
	filename := filepath.Join(w.outputDir, day.Date.Format("2006-01-02")+".md")

	// Check if file exists
	_, err := os.Stat(filename)
	fileExists := !os.IsNotExist(err)

	// Open file for appending (create if doesn't exist)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open markdown file: %w", err)
	}
	defer f.Close()

	// If file is new, write header
	if !fileExists {
		if err := w.writeHeader(f, day); err != nil {
			return err
		}
	}

	// Format and write entry
	entryLine := w.formatEntry(entry)
	if _, err := f.WriteString(entryLine + "\n"); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return nil
}

// writeHeader writes the markdown file header
func (w *Writer) writeHeader(f *os.File, day *database.Day) error {
	var header strings.Builder

	// Title with formatted date
	header.WriteString(fmt.Sprintf("# DAYLOG - %s\n\n", day.Date.Format("Monday, January 2, 2006")))

	// Intention if present
	if day.Intention != nil && *day.Intention != "" {
		header.WriteString(fmt.Sprintf("**Intention:** %s\n\n", *day.Intention))
	}

	// Separator
	header.WriteString("---\n\n")

	if _, err := f.WriteString(header.String()); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	return nil
}

// formatEntry formats an entry as a markdown list item
func (w *Writer) formatEntry(entry *database.Entry) string {
	var line strings.Builder

	// Time and entry text
	line.WriteString(fmt.Sprintf("- %s | %s",
		entry.Timestamp.Format("3:04pm"),
		entry.EntryText))

	// Add momentum if present
	if entry.Momentum != nil {
		switch *entry.Momentum {
		case "up":
			line.WriteString(" ↑")
		case "down":
			line.WriteString(" ↓")
		case "neutral":
			line.WriteString(" →")
		}
	}

	// Add tags if present
	for _, tag := range entry.Tags {
		line.WriteString(" " + tag.TagValue)
	}

	return line.String()
}

// GenerateCompleteDaylog generates a complete daylog with sign-off reflections
// This will be used for the sign-off ritual in Phase 3
func (w *Writer) GenerateCompleteDaylog(day *database.Day, entries []*database.Entry) error {
	filename := filepath.Join(w.outputDir, day.Date.Format("2006-01-02")+".md")

	var content strings.Builder

	// Header
	content.WriteString(fmt.Sprintf("# DAYLOG - %s\n\n", day.Date.Format("Monday, January 2, 2006")))

	// Intention
	if day.Intention != nil && *day.Intention != "" {
		content.WriteString(fmt.Sprintf("**Intention:** %s\n\n", *day.Intention))
	}

	content.WriteString("---\n\n")

	// All entries
	for _, entry := range entries {
		content.WriteString(w.formatEntry(entry) + "\n")
	}

	// Win if present
	if day.Win != nil && *day.Win != "" {
		content.WriteString(fmt.Sprintf("\n**Win:** %s\n", *day.Win))
	}

	// Reflection section if sign-off completed
	if day.PulledOffTrack != nil || day.KeptOnTrack != nil || day.TomorrowProtect != nil {
		content.WriteString("\n---\n\n**Reflection:**\n")

		if day.PulledOffTrack != nil && *day.PulledOffTrack != "" {
			content.WriteString(fmt.Sprintf("- Pulled off track: %s\n", *day.PulledOffTrack))
		}
		if day.KeptOnTrack != nil && *day.KeptOnTrack != "" {
			content.WriteString(fmt.Sprintf("- Kept on track: %s\n", *day.KeptOnTrack))
		}
		if day.TomorrowProtect != nil && *day.TomorrowProtect != "" {
			content.WriteString(fmt.Sprintf("- Tomorrow protect: %s\n", *day.TomorrowProtect))
		}
	}

	// Write complete file (overwrites existing)
	if err := os.WriteFile(filename, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}
