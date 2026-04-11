package main

import (
	"catendar-wails/models"
	"catendar-wails/repository"
)

type CalendarApp struct{}

func (a *CalendarApp) GetEvents(start, end string) ([]models.Event, error) {
	if start != "" && end != "" {
		return repository.GetEventsByDateRange(start, end)
	}
	return repository.GetAllEvents()
}

func (a *CalendarApp) CreateEvent(input models.EventInput) (*models.Event, error) {
	return repository.CreateEvent(input)
}

func (a *CalendarApp) UpdateEvent(id int64, input models.EventInput) (*models.Event, error) {
	return repository.UpdateEvent(id, input)
}

func (a *CalendarApp) DeleteEvent(id int64) error {
	return repository.DeleteEvent(id)
}
