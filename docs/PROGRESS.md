# DAYLOG Development Progress

**⚠️ UPDATE RULE**: This file tracks ONLY upcoming work and planned features. After completing any task, update this file to reflect current status and move completed work to `CHANGELOG.md`.

---

## Current Status

**Phase:** Phase 6 - Intelligence & Awareness Features
**Last Updated:** 2025-10-16
**Recent Completion:** Phase 5 (Historical Viewing & Navigation) ✓
**Next Focus:** Smart prompts and behavioral awareness features

---

## Phase 6: Intelligence & Awareness Features (Next Phase)

**Goal:** Add smart prompts and awareness features to improve user behavior.

### Daily Ritual Settings
- [ ] `log settings` - Configure which daily rituals to enable/disable
  - [ ] Toggle morning intention prompt
  - [ ] Toggle 10-entry win prompt
  - [ ] Toggle drift alerts (and customize threshold)
  - [ ] Toggle sign-off reflection prompts
  - [ ] Toggle after-hours logging messages
  - [ ] Persistent settings storage in database config table
  - [ ] Interactive TUI for settings management

### Anchor Point Suggestions
- [ ] Detect first log after 12:00pm
- [ ] Prepopulate `[ANCHOR - MIDDAY]` tag suggestion
- [ ] Detect first log after 6:00pm
- [ ] Prepopulate `[ANCHOR - EVENING]` tag suggestion
- [ ] Allow user to delete or keep suggestions

### Redirect Prompts
- [ ] Detect distraction keywords (scrolling, browsing, checking, wandering)
- [ ] Display redirect prompt: "What would On-Track You do right now?"
- [ ] Non-judgmental awareness message
- [ ] Press Enter to dismiss

### Energy Pattern Recognition
- [ ] Track ↑ momentum entries over 2+ weeks
- [ ] Identify time ranges with most ↑ entries
- [ ] Display pattern insights: "You log ↑ most often between 9-11am"

### Daily Rituals
- [ ] Display "The Honesty Pact" on first log each day
- [ ] Commitment message: "I commit to honest entries today."

### Time Bookends
- [ ] Detect "Start:" and "Done:" entry pairs
- [ ] Calculate actual time spent on tasks
- [ ] Compare perceived vs actual time
- [ ] Display insights on time estimation

---

## Phase 7: Advanced Analytics (Planned)

**Goal:** Deep pattern analysis and actionable insights.

### Flag Grouping Analysis
- [ ] Identify what activities trigger [LEAK] patterns
- [ ] Identify what activities trigger [FLOW] patterns
- [ ] Identify what activities trigger [STUCK] patterns

### Context Tag Correlation
- [ ] Which tags correlate with high productivity (↑)
- [ ] Which tags correlate with low energy (↓)

### Optimal Work Time Identification
- [ ] Best times for @deep work based on flow states
- [ ] Best times for @admin tasks

### Weekly Insights Generation
- [ ] Productivity trends
- [ ] Distraction patterns
- [ ] Momentum trends

---

## Phase 8: Extended Commands (Planned)

**Goal:** Add utility commands for search, export, and tracking.

### Search & Export
- [ ] `log search [keyword]` - Find past entries containing keyword
  - [ ] Search across all markdown files
  - [ ] Display results with date and context
- [ ] `log export` - Export week/month as formatted document
  - [ ] Generate comprehensive weekly report
  - [ ] Generate monthly summary
  - [ ] Export to PDF or markdown

### Tracking & Insights
- [ ] `log streaks` - Track consecutive days of logging
  - [ ] Display current streak
  - [ ] Display longest streak
  - [ ] Celebrate milestones
- [ ] `log insights` - AI-generated patterns and suggestions
  - [ ] Weekly pattern summary
  - [ ] Personalized recommendations
  - [ ] Productivity optimization tips

### Edge Case Handling
- [ ] Late night logging (after midnight)
  - [ ] Prompt: "Log this as part of today or start of tomorrow?"
  - [ ] Allow user to choose date
- [ ] Empty log entry handling
  - [ ] Display friendly message: "Can't log emptiness! What are you actually doing? 😊"
- [ ] Multiple sign-offs handling
  - [ ] Only last @signoff triggers full ritual
- [ ] No logs today handling
  - [ ] Display helpful prompt for `log view` and `log stats`

---

## Phase 9: Customization & Configuration (Planned)

**Goal:** Allow user customization and improve visual presentation.

### Configuration
- [ ] Config file support (~/.daylog/config.yaml)
  - [ ] Custom markdown output directory
  - [ ] Custom database location
  - [ ] Drift alert threshold (default: 90 minutes)
  - [ ] Enable/disable features
- [ ] Custom user-defined tags
  - [ ] Allow users to add their own @ tags
  - [ ] Allow users to add their own [ ] flags

### Visual Enhancements
- [ ] Dark/light theme toggle
  - [ ] Detect system theme preference
  - [ ] Manual override in config
- [ ] Glamour rendering for markdown display
  - [ ] Render `log view` with Glamour
  - [ ] Render `log week` with Glamour
- [ ] Enhanced styling with Lip Gloss
  - [ ] Consistent color scheme
  - [ ] Improved box borders and spacing
  - [ ] Better typography
- [ ] Optional sound effects toggle
  - [ ] Confetti celebration sound
  - [ ] Entry logged confirmation sound

---

## Phase 10: Polish & Distribution (Planned)

**Goal:** Production-ready distribution and deployment.

### Code Quality
- [ ] Comprehensive error handling
  - [ ] Database connection errors
  - [ ] File system errors
  - [ ] Invalid input handling
- [ ] Input validation and sanitization
- [ ] Performance optimization
  - [ ] Database query optimization
  - [ ] Large file handling
- [ ] Comprehensive test coverage
  - [ ] Unit tests for all packages
  - [ ] Integration tests for flows
  - [ ] TUI component tests

### Distribution
- [ ] Installation script
  - [ ] One-command install
  - [ ] Automatic PATH configuration
  - [ ] Database initialization
- [ ] Homebrew formula
  - [ ] Create tap repository
  - [ ] Submit to Homebrew core
- [ ] Release automation
  - [ ] GitHub Actions CI/CD
  - [ ] Multi-platform builds (macOS, Linux, Windows)
  - [ ] Automated releases

### Optional Integrations
- [ ] Calendar integration
  - [ ] Correlate logs with calendar meetings
  - [ ] Auto-tag @social for meeting times
- [ ] Export to external tools
  - [ ] Notion integration
  - [ ] Obsidian compatibility
  - [ ] Google Calendar sync

---

## Next Steps

### Immediate Priority (Phase 6)
1. Implement `log settings` command with interactive TUI
2. Add anchor point suggestions (midday/evening)
3. Implement redirect prompts for distraction keywords
4. Add energy pattern recognition after 2 weeks of data
5. Display "The Honesty Pact" on first daily log
6. Implement time bookends tracking (Start/Done analysis)

### Medium Priority (Phase 7)
1. Build flag grouping analysis engine
2. Implement context tag correlation
3. Create optimal work time identification
4. Generate weekly insights and trends

### Long-term Priority (Phase 8-10)
1. Search and export functionality
2. Streak tracking and AI insights
3. Full customization and configuration
4. Production release with Homebrew distribution

---

## Blockers & Issues

*None currently*

---

## Notes

- Go version: 1.23.2
- Target platforms: macOS (arm64), Linux (amd64), macOS (amd64)
- Database location: `~/.daylog/daylog.db`
- Markdown output: `~/Documents/daylogs/` (configurable in Phase 9)
- For completed work history, see `CHANGELOG.md`

---

**Remember: This file tracks upcoming work only. Update after completing tasks and move finished work to CHANGELOG.md!**
