package database

import (
	"database/sql"
	"fmt"
)

const (
	// CurrentSchemaVersion is the current database schema version
	CurrentSchemaVersion = 1
)

// migrate runs database migrations
func (s *Store) migrate() error {
	// Enable WAL mode for better concurrency
	if _, err := s.db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign key constraints
	if _, err := s.db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Get current schema version
	currentVersion := s.getSchemaVersion()

	// If already at current version, no migration needed
	if currentVersion >= CurrentSchemaVersion {
		return nil
	}

	// Run migrations
	if err := s.runMigrations(currentVersion); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

// getSchemaVersion returns the current schema version
func (s *Store) getSchemaVersion() int {
	var version int
	err := s.db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return 0 // No schema version yet
	}
	if err != nil {
		return 0 // Assume version 0 on any error
	}
	return version
}

// runMigrations runs all necessary migrations
func (s *Store) runMigrations(fromVersion int) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Migration 0 → 1: Initial schema
	if fromVersion < 1 {
		if err := s.migrateV1(tx); err != nil {
			return fmt.Errorf("failed to migrate to v1: %w", err)
		}
	}

	// Future migrations go here
	// if fromVersion < 2 {
	//     if err := s.migrateV2(tx); err != nil {
	//         return fmt.Errorf("failed to migrate to v2: %w", err)
	//     }
	// }

	return tx.Commit()
}

// migrateV1 creates the initial schema
func (s *Store) migrateV1(tx *sql.Tx) error {
	// Create all tables
	for _, schema := range AllSchemas {
		if _, err := tx.Exec(schema); err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	// Record schema version
	_, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", 1)
	if err != nil {
		return fmt.Errorf("failed to record schema version: %w", err)
	}

	return nil
}
