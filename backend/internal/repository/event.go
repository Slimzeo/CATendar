package repository

import (
	"calendar-backend/internal/database"
	"calendar-backend/internal/models"
	"database/sql"
	"time"
)

func GetAllEvents() ([]models.Event, error) {
	rows, err := database.DB.Query(`
		SELECT id, title, start, end, color, all_day, description, bold, created_at, updated_at
		FROM events
		ORDER BY start ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		var allDay, bold int
		var end, description sql.NullString
		err := rows.Scan(&e.ID, &e.Title, &e.Start, &end, &e.Color, &allDay, &description, &bold, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		e.AllDay = allDay == 1
		e.Bold = bold == 1
		if end.Valid {
			e.End = end.String
		}
		if description.Valid {
			e.Description = description.String
		}
		events = append(events, e)
	}
	return events, nil
}

func GetEventsByDateRange(start, end string) ([]models.Event, error) {
	rows, err := database.DB.Query(`
		SELECT id, title, start, end, color, all_day, description, bold, created_at, updated_at
		FROM events
		WHERE start >= ? AND start <= ?
		ORDER BY start ASC
	`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		var allDay, bold int
		var endNull, description sql.NullString
		err := rows.Scan(&e.ID, &e.Title, &e.Start, &endNull, &e.Color, &allDay, &description, &bold, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		e.AllDay = allDay == 1
		e.Bold = bold == 1
		if endNull.Valid {
			e.End = endNull.String
		}
		if description.Valid {
			e.Description = description.String
		}
		events = append(events, e)
	}
	return events, nil
}

func GetEventByID(id int64) (*models.Event, error) {
	var e models.Event
	var allDay, bold int
	var end, description sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, title, start, end, color, all_day, description, bold, created_at, updated_at
		FROM events WHERE id = ?
	`, id).Scan(&e.ID, &e.Title, &e.Start, &end, &e.Color, &allDay, &description, &bold, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	e.AllDay = allDay == 1
	e.Bold = bold == 1
	if end.Valid {
		e.End = end.String
	}
	if description.Valid {
		e.Description = description.String
	}
	return &e, nil
}

func CreateEvent(input models.EventInput) (*models.Event, error) {
	now := time.Now()
	result, err := database.DB.Exec(`
		INSERT INTO events (title, start, end, color, all_day, description, bold, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, input.Title, input.Start, input.End, input.Color, boolToInt(input.AllDay), input.Description, boolToInt(input.Bold), now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetEventByID(id)
}

func UpdateEvent(id int64, input models.EventInput) (*models.Event, error) {
	now := time.Now()
	_, err := database.DB.Exec(`
		UPDATE events SET title = ?, start = ?, end = ?, color = ?, all_day = ?, description = ?, bold = ?, updated_at = ?
		WHERE id = ?
	`, input.Title, input.Start, input.End, input.Color, boolToInt(input.AllDay), input.Description, boolToInt(input.Bold), now, id)
	if err != nil {
		return nil, err
	}
	return GetEventByID(id)
}

func DeleteEvent(id int64) error {
	_, err := database.DB.Exec(`DELETE FROM events WHERE id = ?`, id)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
