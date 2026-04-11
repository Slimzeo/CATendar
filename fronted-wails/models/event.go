package models

import "time"

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Start       string    `json:"start"`
	End         string    `json:"end,omitempty"`
	Color       string    `json:"color"`
	AllDay      bool      `json:"allDay"`
	Description string    `json:"description"`
	Bold        bool      `json:"bold"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EventInput struct {
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end,omitempty"`
	Color       string `json:"color"`
	AllDay      bool   `json:"allDay"`
	Description string `json:"description"`
	Bold        bool   `json:"bold"`
}