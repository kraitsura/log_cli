package database

import "time"

// Day represents a single day's metadata and reflections
type Day struct {
	ID              int       `db:"id"`
	Date            time.Time `db:"date"`
	Intention       *string   `db:"intention"`
	Win             *string   `db:"win"`
	PulledOffTrack  *string   `db:"pulled_off_track"`
	KeptOnTrack     *string   `db:"kept_on_track"`
	TomorrowProtect *string   `db:"tomorrow_protect"`
	Completed       bool      `db:"completed"`
	CreatedAt       time.Time `db:"created_at"`
}

// Entry represents a single log entry
type Entry struct {
	ID        int       `db:"id"`
	DayID     int       `db:"day_id"`
	Timestamp time.Time `db:"timestamp"`
	EntryText string    `db:"entry_text"`
	Momentum  *string   `db:"momentum"` // "up", "neutral", "down"
	CreatedAt time.Time `db:"created_at"`
	Tags      []Tag     `db:"-"` // Loaded separately
}

// Tag represents a context tag or pattern flag
type Tag struct {
	ID       int    `db:"id"`
	EntryID  int    `db:"entry_id"`
	TagType  string `db:"tag_type"`  // "context" or "flag"
	TagValue string `db:"tag_value"` // e.g., "@deep", "[LEAK]"
}

// Momentum types
type Momentum string

const (
	MomentumUp      Momentum = "up"
	MomentumNeutral Momentum = "neutral"
	MomentumDown    Momentum = "down"
	MomentumBack    Momentum = "back" // Waste/destructive action
)

// Context tags
type ContextTag string

const (
	TagDeep    ContextTag = "@deep"
	TagSocial  ContextTag = "@social"
	TagAdmin   ContextTag = "@admin"
	TagBreak   ContextTag = "@break"
	TagZone    ContextTag = "@zone"
	TagSignoff ContextTag = "@signoff"
)

// Pattern flags
type FlagTag string

const (
	FlagLeak   FlagTag = "[LEAK]"
	FlagFlow   FlagTag = "[FLOW]"
	FlagStuck  FlagTag = "[STUCK]"
	FlagGold   FlagTag = "[GOLD]"
	FlagDrift  FlagTag = "[DRIFT]"
	FlagAnchor FlagTag = "[ANCHOR]"
)

// WeeklyStats holds statistics for a week
type WeeklyStats struct {
	TotalEntries int
	TagCounts    map[string]int
}
