package core

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenDatabaseMigratesLegacyEventTableBeforeCreatingSourceIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "calendar.db")
	legacy, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatalf("open legacy database: %v", err)
	}
	if _, err := legacy.Exec(`
		CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			start TEXT NOT NULL,
			end TEXT,
			color TEXT NOT NULL,
			all_day INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		legacy.Close()
		t.Fatalf("create legacy schema: %v", err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatalf("close legacy database: %v", err)
	}

	db, err := openDatabase(path)
	if err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	for _, column := range []string{"description", "bold", "source_type", "source_account_id", "source_message_id", "source_url"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('events') WHERE name = ?", column).Scan(&count); err != nil {
			t.Fatalf("inspect column %s: %v", column, err)
		}
		if count != 1 {
			t.Fatalf("column %s was not migrated", column)
		}
	}

	var indexCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_events_email_source'").Scan(&indexCount); err != nil {
		t.Fatalf("inspect source index: %v", err)
	}
	if indexCount != 1 {
		t.Fatal("email source index was not created")
	}
}
