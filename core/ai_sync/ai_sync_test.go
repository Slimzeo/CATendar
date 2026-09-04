package ai_sync

import (
	"context"
	"encoding/json"
	"testing"

	synccli "catendar-wails/core/ai_sync/cli"
	"catendar-wails/core/calendar"
	"catendar-wails/core/email"
)

func TestCalendarBatchMayOnlyBeSubmittedOnce(t *testing.T) {
	service := &Service{calendar: calendar.NewService(nil)}

	payload, err := json.Marshal(synccli.CalendarBatchRequest{Events: []calendar.EmailEventInput{}})
	if err != nil {
		t.Fatalf("encode payload: %v", err)
	}
	request := synccli.Request{Command: "calendar.batch-upsert", Payload: payload}

	stateValue := &runState{available: map[string]email.MessageSummary{}}
	if _, err := service.handleCLIRequest(context.Background(), stateValue, request); err != nil {
		t.Fatalf("first batch failed: %v", err)
	}
	if _, err := service.handleCLIRequest(context.Background(), stateValue, request); err == nil {
		t.Fatal("second calendar batch was accepted")
	}
}
