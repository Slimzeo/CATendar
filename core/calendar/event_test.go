package calendar

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "calendar.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(`
		CREATE TABLE events (
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
		CREATE UNIQUE INDEX idx_events_email_source
		ON events(source_account_id, source_message_id, title, start)
		WHERE source_type = 'email';
	`); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	return NewService(db)
}

func TestGetEventsByDateRangeIncludesEntireEndDate(t *testing.T) {
	service := testService(t)
	ctx := context.Background()
	for _, input := range []EventInput{
		{Title: "Last day", Start: "2026-09-30T23:30:00", End: "2026-09-30T23:45:00", Color: defaultEventColor},
		{Title: "Next month", Start: "2026-10-01T00:00:00", End: "2026-10-01T00:30:00", Color: defaultEventColor},
	} {
		if _, err := service.CreateEvent(ctx, input); err != nil {
			t.Fatalf("create %q: %v", input.Title, err)
		}
	}

	events, err := service.GetEventsByDateRange(ctx, "2026-09-01", "2026-09-30")
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	if len(events) != 1 || events[0].Title != "Last day" {
		t.Fatalf("expected only the end-date event, got %#v", events)
	}
}

func TestBatchCreateFromEmailAddsSourceAndDeduplicates(t *testing.T) {
	service := testService(t)
	ctx := context.Background()
	sources := map[string]EmailSource{
		"42:7": {
			AccountID: "primary",
			MessageID: "message@example.com",
			URL:       "catendar://email/primary?uid=7",
		},
	}
	input := []EmailEventInput{{
		Title:         "Design review",
		Start:         "2026-09-08T14:00:00",
		End:           "2026-09-08T15:00:00",
		Description:   "Review the proposal",
		SourceEmailID: "42:7",
	}}

	first, err := service.BatchCreateFromEmail(ctx, input, sources)
	if err != nil {
		t.Fatalf("first batch: %v", err)
	}
	second, err := service.BatchCreateFromEmail(ctx, input, sources)
	if err != nil {
		t.Fatalf("second batch: %v", err)
	}
	if first.Created != 1 || second.Skipped != 1 || second.Conflicts != 0 {
		t.Fatalf("unexpected batch results: first=%+v second=%+v", first, second)
	}

	events, err := service.GetEventsByDateRange(ctx, "2026-09-08", "2026-09-08")
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	if len(events) != 1 || !strings.Contains(events[0].Description, sources["42:7"].URL) {
		t.Fatalf("source link missing from event: %#v", events)
	}
}
