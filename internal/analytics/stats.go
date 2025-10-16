package analytics

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aaryareddy/log_cli/internal/database"
)

// FormatWeeklyStats formats weekly statistics for display
// Now includes momentum distribution if momentumStats is provided
func FormatWeeklyStats(stats *database.WeeklyStats, momentumStats *MomentumStats) string {
	var b strings.Builder

	b.WriteString("THIS WEEK\n")
	b.WriteString(strings.Repeat("─", 40))
	b.WriteString("\n\n")

	// Total entries
	avgPerDay := float64(stats.TotalEntries) / 7.0
	b.WriteString(fmt.Sprintf("Total logs:  %d\n", stats.TotalEntries))
	b.WriteString(fmt.Sprintf("Avg per day: %.1f\n\n", avgPerDay))

	// Tag distribution
	if len(stats.TagCounts) > 0 {
		b.WriteString("YOUR TIME THIS WEEK:\n\n")

		// Calculate total tags for percentages
		totalTags := 0
		for _, count := range stats.TagCounts {
			totalTags += count
		}

		// Sort tags by count (descending)
		type tagCount struct {
			tag   string
			count int
		}
		var tags []tagCount
		for tag, count := range stats.TagCounts {
			tags = append(tags, tagCount{tag, count})
		}
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].count > tags[j].count
		})

		// Display each tag with bar chart
		for _, tc := range tags {
			percentage := float64(tc.count) / float64(totalTags) * 100
			bar := generateBar(percentage, 20)
			b.WriteString(fmt.Sprintf("%-10s %s %.0f%%\n", tc.tag, bar, percentage))
		}
	} else {
		b.WriteString("No context tags logged this week\n")
	}

	// Momentum distribution (if provided)
	if momentumStats != nil && (momentumStats.UpCount > 0 || momentumStats.DownCount > 0 || momentumStats.NeutralCount > 0 || momentumStats.BackCount > 0) {
		b.WriteString("\n")
		b.WriteString("MOMENTUM THIS WEEK:\n\n")

		totalWithMomentum := momentumStats.UpCount + momentumStats.DownCount + momentumStats.NeutralCount + momentumStats.BackCount

		if totalWithMomentum > 0 {
			// Calculate percentages
			upPct := float64(momentumStats.UpCount) / float64(totalWithMomentum) * 100
			neutralPct := float64(momentumStats.NeutralCount) / float64(totalWithMomentum) * 100
			downPct := float64(momentumStats.DownCount) / float64(totalWithMomentum) * 100
			backPct := float64(momentumStats.BackCount) / float64(totalWithMomentum) * 100

			b.WriteString(fmt.Sprintf("↑ Productive:  %2d (%.0f%%)\n", momentumStats.UpCount, upPct))
			b.WriteString(fmt.Sprintf("→ Neutral:     %2d (%.0f%%)\n", momentumStats.NeutralCount, neutralPct))
			b.WriteString(fmt.Sprintf("↓ Dragging:    %2d (%.0f%%)\n", momentumStats.DownCount, downPct))
			b.WriteString(fmt.Sprintf("← Waste:       %2d (%.0f%%)\n", momentumStats.BackCount, backPct))
		}
	}

	return b.String()
}

// generateBar creates an ASCII bar chart
func generateBar(percentage float64, maxWidth int) string {
	filled := int(percentage / 100.0 * float64(maxWidth))
	if filled > maxWidth {
		filled = maxWidth
	}
	if filled < 0 {
		filled = 0
	}

	empty := maxWidth - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}
