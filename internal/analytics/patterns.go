package analytics

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
	"github.com/charmbracelet/lipgloss"
)

// Color palette (duplicated from tui to avoid import cycle)
var (
	colorSuccess  = lipgloss.Color("#50FA7B") // Green
	colorAccent   = lipgloss.Color("#8BE9FD") // Cyan
	colorWarning  = lipgloss.Color("#FFB86C") // Orange
	colorError    = lipgloss.Color("#FF5555") // Red
	colorMetadata = lipgloss.Color("#BD93F9") // Purple
	colorDim      = lipgloss.Color("#6272A4") // Gray
	colorText     = lipgloss.Color("#D0D0D0") // Off-white
)

// Styles for pattern formatting
var (
	flowStyle     = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	goldStyle     = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	stuckStyle    = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	leakStyle     = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	metadataStyle = lipgloss.NewStyle().Foreground(colorMetadata).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(colorDim)
	bodyStyle     = lipgloss.NewStyle().Foreground(colorText)
	successStyle  = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	warningStyle  = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(colorError).Bold(true)
)

// PatternGroup represents entries grouped by a specific pattern flag
type PatternGroup struct {
	FlagType string
	Entries  []*database.Entry
}

// MomentumStats represents momentum distribution statistics
type MomentumStats struct {
	UpCount      int
	DownCount    int
	NeutralCount int
	BackCount    int
	TotalCount   int
}

// GroupByPatternFlag groups entries by pattern flags ([LEAK], [FLOW], [STUCK], [GOLD])
// Returns a map of flag type to entries containing that flag
func GroupByPatternFlag(entries []*database.Entry) map[string][]*database.Entry {
	groups := make(map[string][]*database.Entry)

	// Initialize pattern types
	patternTypes := []string{"[LEAK]", "[FLOW]", "[STUCK]", "[GOLD]"}
	for _, pt := range patternTypes {
		groups[pt] = []*database.Entry{}
	}

	// Group entries by flags
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			if tag.TagType == "flag" {
				// Add to corresponding group
				if _, exists := groups[tag.TagValue]; exists {
					groups[tag.TagValue] = append(groups[tag.TagValue], entry)
				}
			}
		}
	}

	return groups
}

// CalculateMomentumStats calculates momentum distribution from entries
func CalculateMomentumStats(entries []*database.Entry) *MomentumStats {
	stats := &MomentumStats{}

	for _, entry := range entries {
		if entry.Momentum != nil {
			switch *entry.Momentum {
			case "up":
				stats.UpCount++
			case "down":
				stats.DownCount++
			case "neutral":
				stats.NeutralCount++
			case "back":
				stats.BackCount++
			}
		}
		stats.TotalCount++
	}

	return stats
}

// GetWastePatterns filters entries with back (←) momentum
func GetWastePatterns(entries []*database.Entry) []*database.Entry {
	var wasteEntries []*database.Entry

	for _, entry := range entries {
		if entry.Momentum != nil && *entry.Momentum == "back" {
			wasteEntries = append(wasteEntries, entry)
		}
	}

	return wasteEntries
}

// FormatPatternGroup formats a pattern group for display with colored box borders
func FormatPatternGroup(flagType string, entries []*database.Entry) string {
	if len(entries) == 0 {
		return ""
	}

	var b strings.Builder

	// Header with colored box border and description
	var title, description, styledTitle string

	switch flagType {
	case "[LEAK]":
		title = "LEAK PATTERNS"
		description = "What pulled you off track"
		styledTitle = leakStyle.Render("┌─ " + title)
	case "[FLOW]":
		title = "FLOW PATTERNS"
		description = "When you were in the zone"
		styledTitle = flowStyle.Render("┌─ " + title)
	case "[STUCK]":
		title = "STUCK PATTERNS"
		description = "Where you got blocked"
		styledTitle = stuckStyle.Render("┌─ " + title)
	case "[GOLD]":
		title = "GOLD PATTERNS"
		description = "Wins and breakthroughs"
		styledTitle = goldStyle.Render("┌─ " + title)
	default:
		title = fmt.Sprintf("%s PATTERNS", flagType)
		description = ""
		styledTitle = bodyStyle.Render("┌─ " + title)
	}

	// Fancy box border with colored title
	b.WriteString(styledTitle)
	b.WriteString("\n")
	if description != "" {
		b.WriteString(dimStyle.Render("│ " + description))
		b.WriteString("\n")
	}
	b.WriteString(dimStyle.Render("├" + strings.Repeat("─", 59)))
	b.WriteString("\n")

	// Sort entries by timestamp
	sortedEntries := make([]*database.Entry, len(entries))
	copy(sortedEntries, entries)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Timestamp.Before(sortedEntries[j].Timestamp)
	})

	// Format each entry with box prefix
	for _, entry := range sortedEntries {
		// Date and time
		dateStr := entry.Timestamp.Format("Mon 1/2")
		timeStr := entry.Timestamp.Format("3:04pm")
		b.WriteString(dimStyle.Render("│ "))
		b.WriteString(fmt.Sprintf("%s  %s | %s", dateStr, timeStr, entry.EntryText))

		// Add momentum indicator if present
		if entry.Momentum != nil {
			switch *entry.Momentum {
			case "up":
				b.WriteString(" " + successStyle.Render("↑"))
			case "down":
				b.WriteString(" " + warningStyle.Render("↓"))
			case "neutral":
				b.WriteString(" " + bodyStyle.Render("→"))
			case "back":
				b.WriteString(" " + errorStyle.Render("←"))
			}
		}

		// Add context tags (exclude the current flag)
		for _, tag := range entry.Tags {
			if tag.TagType == "context" {
				b.WriteString(" " + dimStyle.Render(tag.TagValue))
			}
		}

		b.WriteString("\n")
	}

	// Closing border
	b.WriteString(dimStyle.Render("└" + strings.Repeat("─", 59)))

	return b.String()
}

// FormatMomentumStats formats momentum statistics for display with colors and bars
func FormatMomentumStats(stats *MomentumStats) string {
	var b strings.Builder

	// Colored header with double line divider
	b.WriteString(metadataStyle.Render("MOMENTUM DISTRIBUTION"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("━", 70)))
	b.WriteString("\n\n")

	// Calculate percentages
	if stats.TotalCount > 0 {
		upPct := float64(stats.UpCount) / float64(stats.TotalCount) * 100
		downPct := float64(stats.DownCount) / float64(stats.TotalCount) * 100
		neutralPct := float64(stats.NeutralCount) / float64(stats.TotalCount) * 100
		backPct := float64(stats.BackCount) / float64(stats.TotalCount) * 100

		// Color-coded momentum lines with visual bars
		b.WriteString(successStyle.Render("↑") + " Productive   ")
		b.WriteString(fmt.Sprintf("%2d (%.0f%%)  %s\n", stats.UpCount, upPct, generateMiniBar(upPct, 10, successStyle)))

		b.WriteString(bodyStyle.Render("→") + " Neutral      ")
		b.WriteString(fmt.Sprintf("%2d (%.0f%%)  %s\n", stats.NeutralCount, neutralPct, generateMiniBar(neutralPct, 10, bodyStyle)))

		b.WriteString(warningStyle.Render("↓") + " Dragging     ")
		b.WriteString(fmt.Sprintf("%2d (%.0f%%)  %s\n", stats.DownCount, downPct, generateMiniBar(downPct, 10, warningStyle)))

		b.WriteString(errorStyle.Render("←") + " Waste        ")
		b.WriteString(fmt.Sprintf("%2d (%.0f%%)  %s\n", stats.BackCount, backPct, generateMiniBar(backPct, 10, errorStyle)))
	} else {
		b.WriteString(dimStyle.Render("  No entries with momentum markers"))
		b.WriteString("\n")
	}

	return b.String()
}

// generateMiniBar creates a small visual bar chart with color-coded filled portion
func generateMiniBar(percentage float64, maxWidth int, filledStyle lipgloss.Style) string {
	filled := int(percentage / 100.0 * float64(maxWidth))
	if filled > maxWidth {
		filled = maxWidth
	}
	if filled < 0 {
		filled = 0
	}

	empty := maxWidth - filled
	filledBar := filledStyle.Render(strings.Repeat("█", filled))
	emptyBar := dimStyle.Render(strings.Repeat("░", empty))
	return filledBar + emptyBar
}

// FormatWastePatterns formats waste pattern entries for display with warning colors
func FormatWastePatterns(entries []*database.Entry) string {
	if len(entries) == 0 {
		return ""
	}

	var b strings.Builder

	// Warning-styled header with double line divider
	b.WriteString("\n")
	b.WriteString(warningStyle.Render("WASTE PATTERNS"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("━", 70)))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Activities marked with ← (back arrow)"))
	b.WriteString("\n\n")

	// Sort by timestamp
	sortedEntries := make([]*database.Entry, len(entries))
	copy(sortedEntries, entries)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Timestamp.Before(sortedEntries[j].Timestamp)
	})

	for _, entry := range sortedEntries {
		dateStr := entry.Timestamp.Format("Mon 1/2")
		timeStr := entry.Timestamp.Format("3:04pm")
		b.WriteString(fmt.Sprintf("  %s  %s | %s ", dateStr, timeStr, entry.EntryText))
		b.WriteString(errorStyle.Render("←"))
		b.WriteString("\n")
	}

	return b.String()
}

// WeeklyPatternSummary contains all pattern analysis for a week
type WeeklyPatternSummary struct {
	StartDate      string
	EndDate        string
	TotalDays      int
	TotalEntries   int
	PatternGroups  map[string][]*database.Entry
	MomentumStats  *MomentumStats
	WastePatterns  []*database.Entry
}

// AnalyzeWeek performs comprehensive pattern analysis on a week of entries
func AnalyzeWeek(entries []*database.Entry, startDate, endDate string) *WeeklyPatternSummary {
	return &WeeklyPatternSummary{
		StartDate:     startDate,
		EndDate:       endDate,
		TotalEntries:  len(entries),
		PatternGroups: GroupByPatternFlag(entries),
		MomentumStats: CalculateMomentumStats(entries),
		WastePatterns: GetWastePatterns(entries),
	}
}
