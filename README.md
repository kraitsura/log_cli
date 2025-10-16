# DAYLOG

**Live Ongoing Genuine Accountability**

A CLI-based system to manage unstructured time through honest, real-time logging.

## Philosophy

**Awareness precedes change.** DAYLOG creates consciousness at decision points while maintaining zero friction for actual logging. This is descriptive accountability, not prescriptive scheduling. The goal isn't perfection—it's consciousness.

## Features

- **Morning Intention** - Set your daily intention on first log
- **Real-Time Logging** - Quick, frictionless log entries with timestamp
- **Momentum Tracking** - Track your energy with ↑ ↓ → markers
- **Context Tags** - Categorize activities with @tags
- **Pattern Flags** - Mark patterns like [LEAK] or [FLOW]
- **Drift Alerts** - Gentle reminder if 90+ minutes since last log
- **Win Prompts** - Celebrate progress at 10 entries
- **Markdown Export** - Human-readable daily logs in `~/Documents/daylogs/`

## Installation

### From Source

```bash
git clone https://github.com/aaryareddy/log_cli.git
cd log_cli
go install ./cmd/log
```

The binary will be installed to your `$GOPATH/bin` as `log`. Make sure `$GOPATH/bin` is in your `PATH`.

### Requirements

- Go 1.21 or higher
- Terminal with Unicode support (for momentum markers)

## Quick Start

### First Log of the Day

```bash
log
```

You'll see a morning intention prompt:

```
┌──────────────────────────────────────────┐
│ Good morning! ⋆｡˚☀︎｡⋆˚                  │
│                                          │
│ Today's Intention:                       │
│                                          │
│ Finish proposal draft                    │
│                                          │
│ Enter to continue (or leave blank)       │
└──────────────────────────────────────────┘
```

After setting your intention (or skipping), you'll see the log entry screen:

```
┌────────────────────────────────────────────────────┐
│ DAYLOG - Wednesday, October 15, 2025              │
│ Started: 9:15am                                    │
│                                                    │
│ 9:15am | What are you doing right now?            │
│                                                    │
│ Morning coffee, reviewing today ↑                  │
│                                                    │
│ ↑ ↓ → | @deep @social @admin @break @zone        │
│ [LEAK] [FLOW] [STUCK] [GOLD]                      │
│                                                    │
│ Enter to submit • Ctrl+C to cancel                │
└────────────────────────────────────────────────────┘
```

### Logging Throughout the Day

Every time you run `log`, you can:

1. **Describe what you're doing** in your own words
2. **Add a momentum marker** (optional):
   - `↑` - feeling productive/energized
   - `→` - neutral/coasting
   - `↓` - dragging/unfocused

3. **Add context tags** (optional):
   - `@deep` - deep focused work
   - `@social` - meetings, calls, collaboration
   - `@admin` - email, scheduling, life stuff
   - `@break` - intentional rest
   - `@zone` - creative/flow work

4. **Add pattern flags** (optional):
   - `[LEAK]` - time drains (social media, news)
   - `[FLOW]` - in the zone, highly productive
   - `[STUCK]` - spinning wheels, unclear what to do
   - `[GOLD]` - unusually productive periods

### Example Entries

```bash
# Simple entry
log
> Working on proposal draft

# With momentum
log
> Deep work session on proposal ↑

# With context tag
log
> Client call about timeline @social

# With pattern flag
log
> Got distracted reading news [LEAK]

# Full example
log
> Finished proposal section ↑ @deep [FLOW]
```

## How It Works

### Data Storage

DAYLOG uses two storage methods:

1. **SQLite Database** (`~/.daylog/daylog.db`)
   - Source of truth for all entries
   - Enables future analytics and stats features
   - Parsed tags and momentum stored separately

2. **Markdown Files** (`~/Documents/daylogs/YYYY-MM-DD.md`)
   - Human-readable daily logs
   - One file per day
   - Can be viewed/edited in any text editor
   - Great for reflection and journaling

### Sample Markdown Output

```markdown
# DAYLOG - Wednesday, October 15, 2025

**Intention:** Finish proposal draft and clear email backlog

---

- 9:15am | Morning coffee, reviewing today ↑
- 9:30am | Email inbox zero attempt @admin
- 10:00am | Got sidetracked reading article [LEAK]
- 10:20am | Back to emails →
- 10:45am | Inbox at 5 messages ↑
- 12:30pm | Morning was scattered but making progress →
- 1:15pm | Lunch and walk ↑ @break
- 2:00pm | Client call @social
- 2:45pm | Project proposal ↑ @deep
- 4:30pm | Still on proposal, in flow state ↑ [FLOW] @zone

**Win:** Finished proposal draft despite rough start
```

## Smart Features

### Drift Alert

If you haven't logged in 90+ minutes, you'll see:

```
[!] DRIFT ALERT - Last log was 2h 15m ago
```

This gentle reminder helps maintain awareness without being punitive.

### Win Prompt

After your 10th log entry of the day, you'll be prompted:

```
┌────────────────────────────────────┐
│ You've logged 10 entries today!   │
│                                    │
│ Any wins today?                    │
│                                    │
│ _                                  │
│                                    │
│ Enter to continue (or skip)        │
└────────────────────────────────────┘
```

This helps you celebrate progress and maintain positive momentum.

## Tag Reference

### Momentum Markers

| Marker | Meaning | When to Use |
|--------|---------|-------------|
| `↑` | Up | Feeling productive, energized, in flow |
| `→` | Neutral | Coasting, steady state, normal energy |
| `↓` | Down | Dragging, unfocused, low energy |

### Context Tags

| Tag | Purpose | Examples |
|-----|---------|----------|
| `@deep` | Focused, cognitive work | Writing, coding, analysis |
| `@social` | Collaborative work | Meetings, calls, discussions |
| `@admin` | Administrative tasks | Email, scheduling, organization |
| `@break` | Intentional rest | Lunch, walks, meditation |
| `@zone` | Creative flow work | Design, brainstorming, prototyping |

### Pattern Flags

| Flag | Purpose | Examples |
|------|---------|----------|
| `[LEAK]` | Time drains | Social media, news spirals, rabbit holes |
| `[FLOW]` | High productivity | Deep work sessions, creative bursts |
| `[STUCK]` | Blocked or unclear | Analysis paralysis, waiting, confusion |
| `[GOLD]` | Exceptional periods | Best work, breakthroughs, peak performance |

## Building from Source

### Development

```bash
# Clone repository
git clone https://github.com/aaryareddy/log_cli.git
cd log_cli

# Install dependencies
go mod download

# Run directly
go run cmd/log/main.go

# Build binary
go build -o log cmd/log/main.go

# Run tests
go test ./...
```

### Project Structure

```
log_cli/
├── cmd/log/main.go           # CLI entry point
├── internal/
│   ├── database/             # SQLite operations
│   ├── markdown/             # Markdown file generation
│   ├── parser/               # Entry text parsing
│   └── tui/                  # Bubble Tea UI components
├── docs/                     # Documentation
└── daylogs/                  # Sample markdown output
```

## Tips for Success

1. **Be honest** - The system only works if you log truthfully
2. **Keep it simple** - Don't overthink entries, just describe what you're doing
3. **Don't batch log** - Real-time logging creates awareness; retroactive logging doesn't
4. **Tags are optional** - Use them when helpful, skip when not
5. **Embrace drift** - The goal isn't perfection, it's consciousness
6. **Review weekly** - Look back at your markdown files to spot patterns

## Privacy & Security

- **100% local** - All data stays on your machine
- **No network calls** - Completely offline
- **No tracking** - No analytics, no telemetry
- **Your data** - Markdown files are yours to keep, share, or delete

## Roadmap

### Phase 2 - COMPLETE
- [x] Entry parsing with tags and momentum
- [x] Markdown generation
- [x] Morning intention prompt
- [x] Win prompt at 10 entries

### Phase 3 (Planned)
- [ ] `log view` - View today's entries
- [ ] `log stats` - Weekly statistics
- [ ] Sign-off ritual with reflection questions
- [ ] Animated confetti on day completion

### Phase 4 (Future)
- [ ] Pattern analysis and insights
- [ ] Energy pattern recognition
- [ ] Weekly review with trends
- [ ] Search and filter entries

## Contributing

This is currently a personal project, but feedback and suggestions are welcome! Please open an issue on GitHub.

## License

MIT License - See LICENSE file for details

## Philosophy & Design

DAYLOG is built on these principles:

- **Awareness over enforcement** - You're tracking yourself, not being tracked
- **Zero friction to log** - Opens instantly, closes on Enter
- **Honest defaults** - Assumes you'll tell the truth
- **Gentle nudges** - Alerts inform, they don't punish
- **Beautiful conclusion** - End-of-day rituals make completion satisfying
- **Portable data** - Markdown files you own forever

The goal isn't to optimize every minute. The goal is to become more conscious of how you spend your time, so you can make better choices.

---

**Made with Go and Bubble Tea**
