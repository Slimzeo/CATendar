package handlers

import (
	"calendar-backend/internal/models"
	"calendar-backend/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetEvents(c echo.Context) error {
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	if start != "" && end != "" {
		events, err := repository.GetEventsByDateRange(start, end)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, events)
	}

	events, err := repository.GetAllEvents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, events)
}

func GetEvent(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	event, err := repository.GetEventByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
	}
	return c.JSON(http.StatusOK, event)
}

func CreateEvent(c echo.Context) error {
	var input models.EventInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if input.Title == "" || input.Start == "" || input.Color == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title, start, and color are required"})
	}

	event, err := repository.CreateEvent(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, event)
}

func UpdateEvent(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var input models.EventInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	event, err := repository.UpdateEvent(id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, event)
}

func DeleteEvent(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := repository.DeleteEvent(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func parseID(c echo.Context) (int64, error) {
	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return 0, err
	}
	return id, nil
}
