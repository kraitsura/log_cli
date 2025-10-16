package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/aaryareddy/log_cli/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

// ConfettiModel shows an animated celebration for day completion
type ConfettiModel struct {
	day           *database.Day
	entries       []*database.Entry
	frame         int
	totalFrames   int
	autoCloseTime time.Time
}

// tickMsg is sent on each animation frame
type tickMsg time.Time

// NewConfettiModel creates a new confetti animation model
func NewConfettiModel(day *database.Day, entries []*database.Entry) ConfettiModel {
	return ConfettiModel{
		day:           day,
		entries:       entries,
		frame:         0,
		totalFrames:   20, // ~2 seconds at 10fps
		autoCloseTime: time.Now().Add(2 * time.Second),
	}
}

// Init starts the animation ticker
func (m ConfettiModel) Init() tea.Cmd {
	return tick()
}

// tick sends a message every 100ms for animation
func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (m ConfettiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Any key press closes immediately
		return m, tea.Quit

	case tickMsg:
		m.frame++

		// Auto-close after animation completes
		if m.frame >= m.totalFrames {
			return m, tea.Quit
		}

		// Continue ticking
		return m, tick()
	}

	return m, nil
}

// View renders the animated celebration
func (m ConfettiModel) View() string {
	var b strings.Builder

	// Animated confetti border
	confetti := m.getConfettiFrame()
	b.WriteString(confetti)
	b.WriteString("\n")

	// Header
	b.WriteString(HeaderStyle.Render("DAYLOG COMPLETE!"))
	b.WriteString("\n\n")

	// Date
	b.WriteString(BoldStyle.Render(m.day.Date.Format("Monday, January 2, 2006")))
	b.WriteString("\n\n")

	// Intention
	if m.day.Intention != nil && *m.day.Intention != "" {
		b.WriteString(DimStyle.Render("Intention: "))
		b.WriteString(*m.day.Intention)
		b.WriteString("\n\n")
	}

	// Entry count
	b.WriteString(SuccessStyle.Render(fmt.Sprintf("✓ %d entries logged", len(m.entries))))
	b.WriteString("\n\n")

	// Win if present
	if m.day.Win != nil && *m.day.Win != "" {
		b.WriteString(SuccessStyle.Render("Win: "))
		b.WriteString(*m.day.Win)
		b.WriteString("\n\n")
	}

	// Reflections summary
	if m.day.Completed {
		b.WriteString(DimStyle.Render("─────────────────────────────────────"))
		b.WriteString("\n")

		if m.day.PulledOffTrack != nil && *m.day.PulledOffTrack != "" {
			b.WriteString(DimStyle.Render("✗ "))
			b.WriteString(*m.day.PulledOffTrack)
			b.WriteString("\n")
		}

		if m.day.KeptOnTrack != nil && *m.day.KeptOnTrack != "" {
			b.WriteString(DimStyle.Render("✓ "))
			b.WriteString(*m.day.KeptOnTrack)
			b.WriteString("\n")
		}

		if m.day.TomorrowProtect != nil && *m.day.TomorrowProtect != "" {
			b.WriteString(DimStyle.Render("→ "))
			b.WriteString(*m.day.TomorrowProtect)
			b.WriteString("\n")
		}
	}

	// Bottom confetti
	b.WriteString("\n")
	b.WriteString(confetti)

	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("Great work today! (closing in a moment...)"))

	return BoxStyle.Render(b.String())
}

// getConfettiFrame returns animated confetti pattern for current frame
func (m ConfettiModel) getConfettiFrame() string {
	// Cycle through different confetti patterns
	patterns := []string{
		"    ✨  🎉  ✨  🎊  ✨  🎉  ✨    ",
		"    🎊  ✨  🎉  ✨  🎊  ✨  🎉    ",
		"    🎉  🎊  ✨  🎉  ✨  🎊  ✨    ",
		"    ✨  🎉  🎊  ✨  🎉  ✨  🎊    ",
	}

	index := m.frame % len(patterns)
	return patterns[index]
}
