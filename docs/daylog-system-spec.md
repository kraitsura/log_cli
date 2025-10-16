# DAYLOG System Specification
**Live Ongoing Genuine Accountability**

## Overview

DAYLOG is a CLI-based accountability system that helps users manage unstructured time through honest, real-time logging. The system creates awareness at decision points while maintaining zero friction for actual logging.

---

## Core Philosophy

**Awareness precedes change.** The act of writing down what you're doing creates a moment of consciousness that can redirect behavior. DAYLOG is descriptive accountability, not prescriptive scheduling or surveillance.

---

## System Components

### 1. LOG Acronym
**Live Ongoing Genuine** - Emphasizes real-time, honest tracking

### 2. Momentum Markers
- `↑` feeling productive/energized
- `→` neutral/coasting  
- `↓` dragging/unfocused

### 3. Context Tags
- `@deep` - deep focused work
- `@social` - meetings, calls, collaboration  
- `@admin` - email, scheduling, life stuff
- `@break` - intentional rest
- `@zone` - creative/flow work
- `@signoff` - end of day marker

### 4. Pattern Flags
- `[LEAK]` - time drains (social media, news spirals)
- `[FLOW]` - in the zone, highly productive
- `[STUCK]` - spinning wheels, unclear what to do
- `[GOLD]` - unusually productive periods
- `[DRIFT]` - more than 90 minutes without logging
- `[ANCHOR]` - non-negotiable check-in points

---

## CLI Application Flow

### Command: `log`

#### First Log of Day
```
Good morning! ☀️

Today's Intention:
_
```

After entering intention, transitions immediately to first log entry.

#### Standard Log Entry Interface
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DAYLOG - Wednesday, October 15, 2025
Started: 9:15am
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

10:47am | What are you doing right now?

_

↑ ↓ → | @deep @social @admin @break @zone | [LEAK] [FLOW] [STUCK] [GOLD] [DRIFT]
```

**Behavior:**
- Single input box with focus
- Autocomplete for `@` tags and `[` flags
- Press Enter → logs entry, TUI closes immediately
- Can type intention on first log or press Enter to skip

#### 10-Entry Win Prompt
After the 10th log entry of the day:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
You've logged 10 entries today! 💪

Any wins today?
_

(Press Enter to skip)
```

After answering or skipping, proceeds to normal log entry.

#### Drift Alert
If 90+ minutes since last log:
```
⚠️ DRIFT ALERT - Last log was 2h 15m ago

10:47am | What are you doing right now?

_
```

#### Anchor Point Suggestions
System suggests anchor tags at key times:
- First log after 12:00pm → `[ANCHOR - MIDDAY]` prepopulated
- First log after 6:00pm → `[ANCHOR - EVENING]` prepopulated

User can delete or keep suggestions.

---

### Sign-Off Flow

When user includes `@signoff` tag in any log:

**Step 1:** Log the entry normally

**Step 2:** Show sign-off questions
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DAY COMPLETE 🌙
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Today's Intention: Finish proposal draft and clear email backlog

What pulled you off track today?
_

What kept you on track today?
_

One thing you'll protect tomorrow?
_
```

**Step 3:** Display full daylog with confetti animation
```
    ✨  🎉  ✨  🎊  ✨  🎉  ✨
    
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DAYLOG - Wednesday, October 15, 2025
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Intention: Finish proposal draft and clear email backlog

9:15am | Morning coffee, reviewing today ↑
9:30am | Start: Email inbox zero attempt → @admin
10:00am | Got sidetracked reading article [DRIFT]
10:20am | Back to emails →
10:45am | Done: Inbox at 5 messages ↑

12:30pm | [ANCHOR - MIDDAY] Morning was scattered but making progress →

1:15pm | Lunch + walk ↑ @break
2:00pm | Client call → @social
2:45pm | Start: Project proposal ↑ @deep
4:30pm | Still on proposal, in flow state ↑ [FLOW] @zone
5:15pm | Done: Rough draft complete ↑

Win: Finished proposal draft despite rough start 🌟

6:00pm | Wrapping up for the day @signoff

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Pulled off track: Article rabbit hole
Kept on track: Client deadline pressure  
Tomorrow protect: First 2 hours for deep work
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    ✨  🎉  ✨  🎊  ✨  🎉  ✨
```

Confetti animation plays for 2-3 seconds, then TUI closes.

---

## Additional Commands

### `log view`
Shows current day's log without opening input interface.

### `log yesterday`  
Shows previous day's complete log.

### `log week`
Shows weekly review with pattern analysis:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
WEEKLY REVIEW
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[LEAK] patterns found 8 times:
- Reading articles (3x)
- Social media (3x)  
- News checking (2x)

[FLOW] patterns found 5 times:
- Morning deep work sessions (3x)
- Afternoon project work (2x)

[STUCK] patterns found 4 times:
- Email decision paralysis (2x)
- Project planning (2x)

[GOLD] patterns found 2 times:
- Tuesday 9am-11am deep work
- Thursday 2pm-5pm proposal writing
```

### `log stats`
Shows basic metrics:
```
This Week:
- Total logs: 67
- Avg per day: 9.6
- Drift alerts: 5
- Flow states: 8
- Wins recorded: 6

Your time this week:
@deep:   ████████░░ 35%
@admin:  ████░░░░░░ 18%
@social: ███░░░░░░░ 12%
@break:  ███░░░░░░░ 15%
@zone:   ████░░░░░░ 20%
```

---

## Smart Features

### The Honesty Pact
Displayed at top of first log each day:
```
"I commit to honest entries today."
```

### Redirect Prompt
If system detects distraction keywords (scrolling, browsing, checking, wandering):
```
Logged! 📝

Quick question: What would On-Track You do right now?
(Just awareness - no judgment)

Press Enter to dismiss
```

### Energy Pattern Recognition
After 2 weeks of logging:
```
💡 Pattern noticed: You log ↑ most often between 9-11am
```

### Time Bookends
Encourage logging both start and completion:
```
10:00am | Start: Writing blog post ↑ @deep
11:30am | Done: Blog post first draft ↑ @deep
```

Shows actual vs perceived time spent.

---

## Data Storage

### File Structure
```
daylogs/
  2025-10-15.md
  2025-10-14.md
  2025-10-13.md
```

Each day creates a human-readable markdown file that can be opened in any text editor.

### Sample File Format
```markdown
# DAYLOG - Wednesday, October 15, 2025

**Intention:** Finish proposal draft and clear email backlog

---

- 9:15am | Morning coffee, reviewing today ↑
- 9:30am | Start: Email inbox zero attempt → @admin
- 10:00am | Got sidetracked reading article [DRIFT]
- 10:20am | Back to emails →
- 10:45am | Done: Inbox at 5 messages ↑
- 12:30pm | [ANCHOR - MIDDAY] Morning was scattered but making progress →
- 1:15pm | Lunch + walk ↑ @break
- 2:00pm | Client call → @social
- 2:45pm | Start: Project proposal ↑ @deep
- 4:30pm | Still on proposal, in flow state ↑ [FLOW] @zone
- 5:15pm | Done: Rough draft complete ↑

**Win:** Finished proposal draft despite rough start 🌟

- 6:00pm | Wrapping up for the day @signoff

---

**Reflection:**
- Pulled off track: Article rabbit hole
- Kept on track: Client deadline pressure  
- Tomorrow protect: First 2 hours for deep work
```

---

## Edge Cases

### Multiple Sign-offs
Only the last `@signoff` entry triggers the full sign-off ritual.

### Skipping Intention
First log prompt allows skipping:
```
Today's Intention:
(Press Enter to skip)
```

### Empty Log Entry
```
Can't log emptiness! What are you actually doing? 😊
```

### Late Night Logging
After midnight:
```
It's past midnight! Log this as:
[1] Part of today (Oct 15)
[2] Start of tomorrow (Oct 16)
```

### No Logs Today
If user runs `log view` or `log stats` with no entries:
```
No logs yet today. Start logging to build your daylog!
Type: log
```

---

## Why This Works

1. **Zero friction to log** - Just type `log` and one line
2. **Honest defaults** - Assumes you'll tell the truth
3. **Gentle nudges** - Redirect prompts, drift alerts (not harsh)
4. **Celebration** - Confetti at end of day, wins recognition
5. **Awareness over enforcement** - You're tracking yourself, not being tracked
6. **Beautiful conclusion** - Sign-off ritual makes completion satisfying
7. **Portable data** - Markdown files you own forever
8. **Pattern recognition** - System learns and reflects insights back

---

## Design Principles

### Descriptive, Not Prescriptive
DAYLOG doesn't tell you what to do. It asks what you're doing and creates awareness.

### Interruption as Feature
The act of opening the TUI and typing interrupts automatic behavior—that's the point.

### No Judgment, Just Facts
The system never scolds. It observes, reflects, and celebrates.

### Friction at Decision Points, Frictionless for Logging
Hard to drift unconsciously. Easy to log honestly.

### Respect for Human Nature
People forget, get distracted, have bad days. The system accommodates this with drift alerts and gentle reminders, not punishment.

---

## Future Enhancements (Optional)

- `log search [keyword]` - Find past entries
- `log export` - Export week/month as formatted document
- `log insights` - AI-generated patterns and suggestions
- `log streaks` - Track consecutive days of logging
- Custom tags defined by user
- Integration with calendar for meeting correlation
- Dark/light theme toggle
- Sound effects toggle for confetti

---

## Success Metrics

DAYLOG succeeds when:
- Users feel more aware of how they spend time
- Automatic distraction behaviors decrease
- Users can identify their best working conditions
- The daily ritual becomes satisfying, not burdensome
- Users continue logging without external accountability

The goal isn't perfection—it's consciousness.