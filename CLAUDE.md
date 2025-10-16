# DAYLOG - Development Guide for Claude

## Project Overview

DAYLOG is a CLI-based accountability system written in Go that helps users manage unstructured time through honest, real-time logging. The system creates awareness at decision points while maintaining zero friction for actual logging.

**Tagline:** Live Ongoing Genuine Accountability

## Architecture

### Technology Stack
- **Language:** Go 1.23.2
- **Database:** SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- **TUI Framework:** Bubble Tea (`github.com/charmbracelet/bubbletea`)
- **UI Components:** Bubbles (`github.com/charmbracelet/bubbles`)
- **Styling:** Lip Gloss (`github.com/charmbracelet/lipgloss`)
- **Markdown:** Glamour (`github.com/charmbracelet/glamour`)

### Project Structure
```
log_cli/
├── cmd/
│   └── log/
│       └── main.go                 # CLI entry point
├── internal/
│   ├── app/
│   │   └── app.go                  # Application coordinator
│   ├── database/
│   │   ├── db.go                   # Database interface
│   │   ├── schema.go               # Schema definitions
│   │   ├── migrations.go           # Schema migrations
│   │   └── models.go               # Data models
│   ├── markdown/
│   │   ├── writer.go               # Generate .md from entries
│   │   ├── parser.go               # Parse .md to entries (recovery)
│   │   └── formatter.go            # Format helpers
│   ├── tui/
│   │   ├── log_entry.go            # Main log entry screen
│   │   ├── intention.go            # Morning intention screen
│   │   ├── signoff.go              # End-of-day sign-off
│   │   ├── view.go                 # View existing logs
│   │   ├── confetti.go             # Confetti animation
│   │   └── styles.go               # Lip Gloss style definitions
│   ├── analytics/
│   │   ├── patterns.go             # Pattern detection
│   │   ├── stats.go                # Statistics calculations
│   │   └── insights.go             # Smart suggestions
│   └── config/
│       └── config.go               # App configuration
├── docs/
│   ├── daylog-system-spec.md       # System specification
│   ├── daylog_implementation.txt   # Implementation guide
│   ├── PROGRESS.md                 # Upcoming work tracking
│   └── CHANGELOG.md                # Completed work history
├── daylogs/                        # Generated markdown files
│   └── .gitkeep
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md                       # This file
└── .gitignore
```

## Core Concepts

### Data Models
- **Day**: Represents a single day with intention, win, reflections
- **Entry**: Individual log entry with timestamp, text, momentum, tags
- **Tag**: Context tags (@deep, @social) and pattern flags ([LEAK], [FLOW])
- **Momentum**: Up (↑), Neutral (→), Down (↓)

### User Flow
1. `log` → Opens TUI with timestamp
2. User types what they're doing
3. Can add momentum markers and tags
4. Press Enter → Saves to SQLite + Appends to markdown file
5. TUI closes immediately (zero friction)

### Key Features
- **Morning Intention**: First log prompts for daily intention
- **Drift Alerts**: Warns if 90+ minutes since last log
- **Win Prompt**: After 10 entries, asks for wins
- **Sign-off Ritual**: `@signoff` tag triggers reflection questions + confetti
- **Pattern Recognition**: Weekly analysis of LEAK/FLOW/STUCK/GOLD patterns

## Development Workflow

### Running the Application
```bash
# Run directly
go run cmd/log/main.go

# Build and run
go build -o log cmd/log/main.go
./log

# Install locally
go install ./cmd/log
log
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/database/...

# Verbose output
go test -v ./...
```

### Code Formatting
```bash
# Format all code
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

## Go Conventions for This Project

### 1. Error Handling
Always handle errors explicitly. Never ignore errors with `_`.

```go
// Good
if err := store.InsertEntry(entry); err != nil {
    log.Printf("Failed to insert entry: %v", err)
    return err
}

// Bad
store.InsertEntry(entry)
```

### 2. Database Transactions
Use transactions for multi-step database operations.

```go
tx, err := s.db.Begin()
if err != nil {
    return err
}
defer tx.Rollback() // Safe to call even after commit

// ... operations ...

return tx.Commit()
```

### 3. Pointers for Optional Fields
Use pointers for fields that can be NULL in database.

```go
type Day struct {
    ID        int
    Date      time.Time
    Intention *string  // Can be nil
    Win       *string  // Can be nil
}
```

### 4. Struct Tags
Use struct tags for database field mapping.

```go
type Entry struct {
    ID        int       `db:"id"`
    Timestamp time.Time `db:"timestamp"`
    EntryText string    `db:"entry_text"`
}
```

### 5. Package Organization
- `cmd/`: Executable entry points
- `internal/`: Private application code (not importable by other projects)
- Each package has a clear, single responsibility

### 6. Bubble Tea Patterns
Follow the Elm Architecture:
- **Model**: Application state
- **Init**: Initialize commands
- **Update**: Handle messages, update state
- **View**: Render UI

```go
type Model struct {
    // state fields
}

func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (m Model) View() string { ... }
```

### 7. Time Handling
Always use `time.Time` for timestamps. Format for display:
- Database: `time.Now()`
- Display: `time.Now().Format("3:04pm")`
- Date only: `time.Now().Format("2006-01-02")`

### 8. SQL Best Practices
- Use parameterized queries (never string concatenation)
- Create indexes on frequently queried columns
- Use prepared statements for repeated queries

```go
// Good
rows, err := db.Query("SELECT * FROM entries WHERE day_id = ?", dayID)

// Bad
query := fmt.Sprintf("SELECT * FROM entries WHERE day_id = %d", dayID)
```

## Key Development Rules

### 1. Progress Tracking
**MUST UPDATE** `docs/PROGRESS.md` after completing any task:
- Update current status and mark checkboxes for upcoming work
- PROGRESS.md tracks ONLY planned/upcoming features (forward-looking)
- Move completed work to `docs/CHANGELOG.md` with timestamps and details
- Document any blockers or decisions in appropriate file

### 2. Testing Requirements
- Write unit tests for database operations
- Write integration tests for end-to-end flows
- Test TUI components with mock data

### 3. Data Integrity
- SQLite database is source of truth
- Markdown files are human-readable backups
- Implement recovery function to rebuild DB from markdown if needed

### 4. User Experience Principles
- **Zero friction for logging**: TUI should open instantly, close on Enter
- **No judgment**: Never scold or criticize user entries
- **Gentle nudges**: Drift alerts are informative, not punitive
- **Celebration**: Show wins, confetti, positive reinforcement

### 5. Code Quality
- All exported functions must have comments
- Keep functions focused and small
- Use meaningful variable names
- No magic numbers (use constants)

## Common Tasks

### Adding a New Tag Type
1. Add constant to `internal/database/models.go`
2. Update tag parsing in entry parser
3. Update markdown writer to format new tag
4. Add to TUI helper text
5. Update documentation

### Adding a New Command
1. Update `cmd/log/main.go` switch statement
2. Create new function for command logic
3. Add TUI screen if interactive (in `internal/tui/`)
4. Update README with command usage
5. Test end-to-end

### Modifying Database Schema
1. Create new migration in `internal/database/migrations.go`
2. Increment migration version
3. Update models in `internal/database/models.go`
4. Update queries in `internal/database/db.go`
5. Test with fresh database and existing database

## Reference Documentation

- **System Spec**: `docs/daylog-system-spec.md` - User-facing features and flows
- **Implementation Guide**: `docs/daylog_implementation.txt` - Detailed technical guide
- **Progress Tracking**: `docs/PROGRESS.md` - Current status and upcoming planned work
- **Changelog**: `docs/CHANGELOG.md` - Complete history of completed phases and features

## Philosophy

**Awareness precedes change.** DAYLOG is descriptive accountability, not prescriptive scheduling. The goal isn't perfection—it's consciousness.

## Development Phases

### Phase 1: Core Functionality (Current)
- Set up project structure
- Database schema and migrations
- Basic log entry TUI
- Markdown generation
- End-to-end flow working

### Phase 2: Smart Features
- Tag parsing and momentum tracking
- Drift detection
- Morning intention
- Win prompts
- Sign-off ritual

### Phase 3: Viewing & Stats
- `log view` command
- `log stats` command
- Weekly pattern analysis

### Phase 4: Intelligence
- Pattern recognition
- Redirect prompts
- Energy pattern detection
- Anchor point suggestions

### Phase 5: Polish
- Confetti animation
- Better styling
- Error handling
- Configuration file support

## Getting Help

When working on this project, Claude should:
1. Always reference the spec files for feature details
2. Follow Go best practices and idioms
3. Keep code simple and maintainable
4. Update PROGRESS.md after each completed task (keep forward-looking)
5. Move completed work to CHANGELOG.md with full details
6. Test database operations thoroughly
7. Ensure TUI responsiveness and good UX

## Quick Commands Reference

```bash
# Development
go run cmd/log/main.go          # Run app
go test ./...                    # Run tests
go build -o log cmd/log/main.go  # Build binary

# Dependencies
go get <package>                 # Add dependency
go mod tidy                      # Clean up dependencies
go mod download                  # Download dependencies

# Code Quality
go fmt ./...                     # Format code
go vet ./...                     # Vet code
staticcheck ./...                # Static analysis (if installed)
```
