# DAYLOG

Command-line accountability system for managing unstructured time through real-time logging.

---

## Overview

DAYLOG creates consciousness at decision points. Awareness precedes change.

**Core Principle:** Descriptive accountability, not prescriptive scheduling.

### Key Features

- Zero-friction logging with instant TUI
- Morning intention setting
- Real-time momentum tracking (`↑` `→` `↓`)
- Context tagging (`@deep` `@social` `@admin`)
- Pattern recognition (`[FLOW]` `[LEAK]` `[STUCK]` `[GOLD]`)
- Local-first: SQLite + Markdown export
- Drift alerts for extended gaps

---

## Installation

```bash
git clone https://github.com/aaryareddy/log_cli.git
cd log_cli
go install ./cmd/log
```

**Requirements:** Go 1.21+ · Unicode-capable terminal

---

## Usage

### Basic Flow

```bash
log
```

Opens timestamped TUI. Type what you're doing. Press Enter. Done.

### Entry Examples

```
Morning coffee, reviewing today ↑
```

```
Client call about timeline @social
```

```
Got distracted reading news [LEAK]
```

```
Finished proposal section ↑ @deep [FLOW]
```

### Markers & Tags

**Momentum:**
- `↑` Productive, energized
- `→` Neutral, coasting
- `↓` Unfocused, dragging

**Context:**
- `@deep` Focused cognitive work
- `@social` Meetings, collaboration
- `@admin` Email, organization
- `@break` Intentional rest
- `@zone` Creative flow

**Patterns:**
- `[FLOW]` High productivity
- `[LEAK]` Time drains
- `[STUCK]` Blocked, unclear
- `[GOLD]` Peak performance

---

## Data

### Storage

**SQLite** (`~/.daylog/daylog.db`)
Source of truth. Enables analytics.

**Markdown** (`~/Documents/daylogs/YYYY-MM-DD.md`)
Human-readable daily logs. One file per day.

### Sample Output

```markdown
# DAYLOG - Wednesday, October 15, 2025

**Intention:** Finish proposal draft and clear email backlog

---

- 9:15am | Morning coffee, reviewing today ↑
- 9:30am | Email inbox zero attempt @admin
- 10:00am | Got sidetracked reading article [LEAK]
- 10:45am | Inbox at 5 messages ↑
- 1:15pm | Lunch and walk ↑ @break
- 2:45pm | Project proposal ↑ @deep
- 4:30pm | Still on proposal, in flow state ↑ [FLOW] @zone

**Win:** Finished proposal draft despite rough start
```

---

## Intelligence

**Drift Alert:** Notification after 90+ minutes without logging

**Win Prompt:** Celebration trigger at 10 entries

**Morning Intention:** First log of day prompts daily focus

**Pattern Analysis:** (Planned) Weekly insights from tags and flags

---

## Architecture

```
log_cli/
├── cmd/log/              # CLI entry
├── internal/
│   ├── database/         # SQLite operations
│   ├── markdown/         # File generation
│   ├── parser/           # Entry text parsing
│   └── tui/              # Bubble Tea components
└── docs/                 # Specifications
```

**Stack:** Go · SQLite · Bubble Tea · Lip Gloss

---

## Development

```bash
# Run
go run cmd/log/main.go

# Build
go build -o log cmd/log/main.go

# Test
go test ./...

# Format
go fmt ./...
```

---

**Built with Go and Bubble Tea**
