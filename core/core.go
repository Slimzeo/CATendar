package core

import (
	"database/sql"
	"fmt"

	"catendar-wails/core/ai_sync"
	"catendar-wails/core/calendar"
	"catendar-wails/core/email"
)

type Options struct {
	ExecutablePath string
}

type Core struct {
	db       *sql.DB
	Calendar *calendar.Service
	Email    *email.Service
	AISync   *ai_sync.Service
}

func Open(databasePath string, options Options) (*Core, error) {
	db, err := openDatabase(databasePath)
	if err != nil {
		return nil, err
	}

	calendarService := calendar.NewService(db)
	emailService := email.NewService(db, email.NewKeyringCredentials(), email.NewIMAPReader())
	aiSyncService := ai_sync.NewService(db, emailService, calendarService, options.ExecutablePath)

	return &Core{
		db:       db,
		Calendar: calendarService,
		Email:    emailService,
		AISync:   aiSyncService,
	}, nil
}

func (c *Core) Close() error {
	c.AISync.Shutdown()
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}
