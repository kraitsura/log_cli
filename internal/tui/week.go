package tui

import (
	"fmt"
	"strings"

	"github.com/aaryareddy/log_cli/internal/analytics"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// WeekModel is the model for displaying weekly pattern analysis
type WeekModel struct {
	summary  *analytics.WeeklyPatternSummary
	viewport viewport.Model
	ready    bool
}

// NewWeekModel creates a new week review model
func NewWeekModel(summary *analytics.WeeklyPatternSummary) WeekModel {
	return WeekModel{
		summary: summary,
	}
}

// Init initializes the model
func (m WeekModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m WeekModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Initialize viewport on first window size message
			m.viewport = viewport.New(msg.Width, msg.Height-2)
			m.viewport.SetContent(m.generateWeekContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
	}

	// Update viewport (handles scrolling)
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// generateWeekContent generates the weekly review content
func (m WeekModel) generateWeekContent() string {
	var b strings.Builder

	// Header (keep this as is - top border preserved)
	b.WriteString(RenderHeaderBar("WEEKLY REVIEW", fmt.Sprintf("%s to %s", m.summary.StartDate, m.summary.EndDate)))
	b.WriteString("\n\n")

	// Summary stats (no emoji, bold white)
	b.WriteString(BoldStyle.Render(fmt.Sprintf("%d entries logged this week", m.summary.TotalEntries)))
	b.WriteString("\n\n")

	if m.summary.TotalEntries == 0 {
		b.WriteString(DimStyle.Render("No entries to analyze yet. Start logging to see patterns!"))
		return b.String()
	}

	// Pattern groups header (purple metadata style with double line)
	b.WriteString(MetadataStyle.Render("PATTERN ANALYSIS"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(strings.Repeat("━", 70)))
	b.WriteString("\n\n")

	// Display each pattern type
	patternOrder := []string{"[FLOW]", "[GOLD]", "[STUCK]", "[LEAK]"}
	for _, flagType := range patternOrder {
		if entries, exists := m.summary.PatternGroups[flagType]; exists && len(entries) > 0 {
			formatted := analytics.FormatPatternGroup(flagType, entries)
			b.WriteString(formatted)
			b.WriteString("\n\n")
		}
	}

	// Momentum distribution
	b.WriteString("\n")
	b.WriteString(analytics.FormatMomentumStats(m.summary.MomentumStats))
	b.WriteString("\n")

	// Waste patterns (if any)
	if len(m.summary.WastePatterns) > 0 {
		b.WriteString(analytics.FormatWastePatterns(m.summary.WastePatterns))
		b.WriteString("\n")
	}

	// Insights section (cyan accent style with double line)
	b.WriteString("\n")
	b.WriteString(AccentStyle.Render("INSIGHTS"))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render(strings.Repeat("━", 70)))
	b.WriteString("\n\n")

	// Generate insights based on patterns
	insights := generateInsights(m.summary)
	for _, insight := range insights {
		b.WriteString(DimStyle.Render("  → "))
		b.WriteString(insight)
		b.WriteString("\n")
	}

	if len(insights) == 0 {
		b.WriteString(DimStyle.Render("  Keep logging to generate personalized insights!"))
		b.WriteString("\n")
	}

	return b.String()
}

// View renders the UI
func (m WeekModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	viewContent := m.viewport.View() + "\n"
	viewContent += DimStyle.Render("↑/↓ or j/k to scroll • q/esc to exit")
	return viewContent
}

// generateInsights creates actionable insights from the week's data
func generateInsights(summary *analytics.WeeklyPatternSummary) []string {
	var insights []string

	// Momentum insights
	if summary.MomentumStats.TotalCount > 0 {
		upPct := float64(summary.MomentumStats.UpCount) / float64(summary.MomentumStats.TotalCount) * 100
		backPct := float64(summary.MomentumStats.BackCount) / float64(summary.MomentumStats.TotalCount) * 100

		if upPct > 60 {
			insights = append(insights, "Strong momentum this week! You logged ↑ on over 60% of marked entries.")
		} else if upPct < 30 {
			insights = append(insights, "Momentum was lower this week. Consider what conditions help you feel more energized.")
		}

		if backPct > 10 {
			insights = append(insights, fmt.Sprintf("%.0f%% of entries were marked as waste (←). Review these patterns to reclaim time.", backPct))
		}
	}

	// Pattern insights
	if len(summary.PatternGroups["[FLOW]"]) > 0 {
		insights = append(insights, fmt.Sprintf("You hit flow %d times this week. What conditions enabled those states?", len(summary.PatternGroups["[FLOW]"])))
	}

	if len(summary.PatternGroups["[LEAK]"]) > 0 {
		insights = append(insights, fmt.Sprintf("Identified %d leak patterns. Common themes in what pulled you off track?", len(summary.PatternGroups["[LEAK]"])))
	}

	if len(summary.PatternGroups["[STUCK]"]) > 0 {
		insights = append(insights, fmt.Sprintf("You got stuck %d times. Consider documenting solutions when you break through.", len(summary.PatternGroups["[STUCK]"])))
	}

	if len(summary.PatternGroups["[GOLD]"]) > 0 {
		insights = append(insights, fmt.Sprintf("Captured %d gold moments! Celebrate these wins.", len(summary.PatternGroups["[GOLD]"])))
	}

	// Entry frequency insight
	avgPerDay := float64(summary.TotalEntries) / 7.0
	if avgPerDay < 3 {
		insights = append(insights, "Log frequency is low. More frequent logs = better awareness and pattern detection.")
	} else if avgPerDay > 15 {
		insights = append(insights, "High logging frequency! You're building strong awareness habits.")
	}

	return insights
}
