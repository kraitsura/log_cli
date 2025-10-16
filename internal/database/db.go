package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Store handles all database operations
type Store struct {
	db *sql.DB
}

// NewStore creates a new database store and runs migrations
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &Store{db: db}

	// Run migrations
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// GetOrCreateToday gets today's day record or creates it if it doesn't exist
func (s *Store) GetOrCreateToday() (*Day, error) {
	today := time.Now().Format("2006-01-02")

	var day Day
	err := s.db.QueryRow(`
		SELECT id, date, intention, win, pulled_off_track,
		       kept_on_track, tomorrow_protect, completed, created_at
		FROM days WHERE date = ?
	`, today).Scan(
		&day.ID, &day.Date, &day.Intention, &day.Win,
		&day.PulledOffTrack, &day.KeptOnTrack, &day.TomorrowProtect,
		&day.Completed, &day.CreatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new day
		result, err := s.db.Exec(`
			INSERT INTO days (date) VALUES (?)
		`, today)
		if err != nil {
			return nil, fmt.Errorf("failed to create day: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert id: %w", err)
		}

		day.ID = int(id)
		day.Date, _ = time.Parse("2006-01-02", today)
		day.CreatedAt = time.Now()
		day.Completed = false
		return &day, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query day: %w", err)
	}

	return &day, nil
}

// InsertEntry creates a new log entry with tags
func (s *Store) InsertEntry(entry *Entry) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert entry
	result, err := tx.Exec(`
		INSERT INTO entries (day_id, timestamp, entry_text, momentum)
		VALUES (?, ?, ?, ?)
	`, entry.DayID, entry.Timestamp, entry.EntryText, entry.Momentum)
	if err != nil {
		return fmt.Errorf("failed to insert entry: %w", err)
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get entry id: %w", err)
	}
	entry.ID = int(entryID)

	// Insert tags
	for _, tag := range entry.Tags {
		_, err := tx.Exec(`
			INSERT INTO tags (entry_id, tag_type, tag_value)
			VALUES (?, ?, ?)
		`, entryID, tag.TagType, tag.TagValue)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTodayEntries retrieves all entries for today
func (s *Store) GetTodayEntries(dayID int) ([]*Entry, error) {
	rows, err := s.db.Query(`
		SELECT id, day_id, timestamp, entry_text, momentum, created_at
		FROM entries
		WHERE day_id = ?
		ORDER BY timestamp ASC
	`, dayID)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		var e Entry
		err := rows.Scan(&e.ID, &e.DayID, &e.Timestamp,
			&e.EntryText, &e.Momentum, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}

		// Load tags for this entry
		tags, err := s.GetEntryTags(e.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
		e.Tags = tags

		entries = append(entries, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return entries, nil
}

// GetEntryTags retrieves all tags for an entry
func (s *Store) GetEntryTags(entryID int) ([]Tag, error) {
	rows, err := s.db.Query(`
		SELECT id, entry_id, tag_type, tag_value
		FROM tags
		WHERE entry_id = ?
	`, entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.EntryID, &t.TagType, &t.TagValue); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tags, nil
}

// UpdateDayIntention updates the intention for a day
func (s *Store) UpdateDayIntention(dayID int, intention string) error {
	_, err := s.db.Exec(`
		UPDATE days SET intention = ? WHERE id = ?
	`, intention, dayID)
	if err != nil {
		return fmt.Errorf("failed to update intention: %w", err)
	}
	return nil
}

// UpdateDayWin updates the win for a day
func (s *Store) UpdateDayWin(dayID int, win string) error {
	_, err := s.db.Exec(`
		UPDATE days SET win = ? WHERE id = ?
	`, win, dayID)
	if err != nil {
		return fmt.Errorf("failed to update win: %w", err)
	}
	return nil
}

// CompleteDaySignoff marks a day as completed and saves sign-off reflections
func (s *Store) CompleteDaySignoff(dayID int, pulledOff, keptOn, protect string) error {
	_, err := s.db.Exec(`
		UPDATE days
		SET pulled_off_track = ?,
		    kept_on_track = ?,
		    tomorrow_protect = ?,
		    completed = 1
		WHERE id = ?
	`, pulledOff, keptOn, protect, dayID)
	if err != nil {
		return fmt.Errorf("failed to complete sign-off: %w", err)
	}
	return nil
}

// GetWeeklyStats calculates statistics for the past 7 days
func (s *Store) GetWeeklyStats() (*WeeklyStats, error) {
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	// Count entries
	var totalEntries int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM entries e
		JOIN days d ON e.day_id = d.id
		WHERE d.date >= ?
	`, weekAgo).Scan(&totalEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to count entries: %w", err)
	}

	// Tag distribution
	rows, err := s.db.Query(`
		SELECT t.tag_value, COUNT(*) as count
		FROM tags t
		JOIN entries e ON t.entry_id = e.id
		JOIN days d ON e.day_id = d.id
		WHERE d.date >= ? AND t.tag_type = 'context'
		GROUP BY t.tag_value
	`, weekAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag distribution: %w", err)
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tag string
		var count int
		if err := rows.Scan(&tag, &count); err != nil {
			return nil, fmt.Errorf("failed to scan tag count: %w", err)
		}
		tagCounts[tag] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &WeeklyStats{
		TotalEntries: totalEntries,
		TagCounts:    tagCounts,
	}, nil
}
