package database

// Schema holds all SQL schema definitions
const (
	// SchemaDays creates the days table
	SchemaDays = `
CREATE TABLE IF NOT EXISTS days (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date DATE UNIQUE NOT NULL,
	intention TEXT,
	win TEXT,
	pulled_off_track TEXT,
	kept_on_track TEXT,
	tomorrow_protect TEXT,
	completed BOOLEAN DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_days_date ON days(date);
`

	// SchemaEntries creates the entries table
	SchemaEntries = `
CREATE TABLE IF NOT EXISTS entries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	day_id INTEGER NOT NULL,
	timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	entry_text TEXT NOT NULL,
	momentum TEXT CHECK(momentum IN ('up', 'neutral', 'down')),
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_entries_day ON entries(day_id);
CREATE INDEX IF NOT EXISTS idx_entries_timestamp ON entries(timestamp);
`

	// SchemaTags creates the tags table
	SchemaTags = `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	entry_id INTEGER NOT NULL,
	tag_type TEXT NOT NULL CHECK(tag_type IN ('context', 'flag')),
	tag_value TEXT NOT NULL,
	FOREIGN KEY (entry_id) REFERENCES entries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tags_entry ON tags(entry_id);
CREATE INDEX IF NOT EXISTS idx_tags_type_value ON tags(tag_type, tag_value);
`

	// SchemaPatternCache creates the pattern_cache table
	SchemaPatternCache = `
CREATE TABLE IF NOT EXISTS pattern_cache (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date DATE NOT NULL,
	pattern_type TEXT NOT NULL,
	pattern_value TEXT NOT NULL,
	count INTEGER DEFAULT 1,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pattern_date ON pattern_cache(date);
`

	// SchemaConfig creates the config table
	SchemaConfig = `
CREATE TABLE IF NOT EXISTS config (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);
`

	// SchemaVersion creates the schema_version table
	SchemaVersion = `
CREATE TABLE IF NOT EXISTS schema_version (
	version INTEGER NOT NULL,
	applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
)

// AllSchemas is an ordered list of all schema creation statements
var AllSchemas = []string{
	SchemaVersion,
	SchemaDays,
	SchemaEntries,
	SchemaTags,
	SchemaPatternCache,
	SchemaConfig,
}
