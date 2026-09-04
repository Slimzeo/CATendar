package main

import (
	"context"

	"catendar-wails/core/calendar"
)

type CalendarApp struct {
	calendar *calendar.Service
}

func NewCalendarApp(calendarService *calendar.Service) *CalendarApp {
	return &CalendarApp{calendar: calendarService}
}

func (a *CalendarApp) GetEvents(start, end string) ([]calendar.Event, error) {
	return a.calendar.GetEventsByDateRange(context.Background(), start, end)
}

func (a *CalendarApp) CreateEvent(input calendar.EventInput) (*calendar.Event, error) {
	return a.calendar.CreateEvent(context.Background(), input)
}

func (a *CalendarApp) UpdateEvent(id int64, input calendar.EventInput) (*calendar.Event, error) {
	return a.calendar.UpdateEvent(context.Background(), id, input)
}

func (a *CalendarApp) DeleteEvent(id int64) error {
	return a.calendar.DeleteEvent(context.Background(), id)
}
