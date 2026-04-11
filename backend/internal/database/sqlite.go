package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = migrate(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func migrate() error {
	// Check if description column exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('events') WHERE name = 'description'").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		log.Println("Migrating events table: adding description column")
		_, err = DB.Exec("ALTER TABLE events ADD COLUMN description TEXT DEFAULT ''")
		if err != nil {
			return err
		}
	}

	// Check if bold column exists
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('events') WHERE name = 'bold'").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		log.Println("Migrating events table: adding bold column")
		_, err = DB.Exec("ALTER TABLE events ADD COLUMN bold INTEGER DEFAULT 0")
		if err != nil {
			return err
		}
	}
	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		start TEXT NOT NULL,
		end TEXT,
		color TEXT NOT NULL,
		all_day INTEGER DEFAULT 1,
		description TEXT DEFAULT '',
		bold INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_events_start ON events(start);
	`
	_, err := DB.Exec(query)
	return err
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
