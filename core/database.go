package core

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	for _, pragma := range []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("configure database: %w", err)
		}
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}
	if err := migrateEvents(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	const schema = `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			start TEXT NOT NULL,
			end TEXT,
			color TEXT NOT NULL,
			all_day INTEGER NOT NULL DEFAULT 1,
			description TEXT NOT NULL DEFAULT '',
			bold INTEGER NOT NULL DEFAULT 0,
			source_type TEXT NOT NULL DEFAULT '',
			source_account_id TEXT NOT NULL DEFAULT '',
			source_message_id TEXT NOT NULL DEFAULT '',
			source_url TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_events_start ON events(start);

		CREATE TABLE IF NOT EXISTS email_accounts (
			id TEXT PRIMARY KEY,
			address TEXT NOT NULL,
			username TEXT NOT NULL,
			imap_host TEXT NOT NULL,
			imap_port INTEGER NOT NULL,
			folder TEXT NOT NULL DEFAULT 'INBOX',
			use_tls INTEGER NOT NULL DEFAULT 1,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS ai_settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			provider TEXT NOT NULL DEFAULT 'codex',
			codex_path TEXT NOT NULL DEFAULT '',
			claude_path TEXT NOT NULL DEFAULT ''
		);
		INSERT OR IGNORE INTO ai_settings (id) VALUES (1);
	`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("create database schema: %w", err)
	}
	return nil
}

func migrateEvents(db *sql.DB) error {
	columns := []struct {
		name       string
		definition string
	}{
		{"description", "TEXT NOT NULL DEFAULT ''"},
		{"bold", "INTEGER NOT NULL DEFAULT 0"},
		{"source_type", "TEXT NOT NULL DEFAULT ''"},
		{"source_account_id", "TEXT NOT NULL DEFAULT ''"},
		{"source_message_id", "TEXT NOT NULL DEFAULT ''"},
		{"source_url", "TEXT NOT NULL DEFAULT ''"},
	}

	for _, column := range columns {
		var count int
		if err := db.QueryRow(
			"SELECT COUNT(*) FROM pragma_table_info('events') WHERE name = ?",
			column.name,
		).Scan(&count); err != nil {
			return fmt.Errorf("inspect events.%s: %w", column.name, err)
		}
		if count == 0 {
			query := fmt.Sprintf("ALTER TABLE events ADD COLUMN %s %s", column.name, column.definition)
			if _, err := db.Exec(query); err != nil {
				return fmt.Errorf("add events.%s: %w", column.name, err)
			}
		}
	}

	if _, err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_events_email_source
		ON events(source_account_id, source_message_id, title, start)
		WHERE source_type = 'email'
	`); err != nil {
		return fmt.Errorf("create email event index: %w", err)
	}
	return nil
}
