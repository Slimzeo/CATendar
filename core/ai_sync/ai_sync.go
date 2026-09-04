package ai_sync

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"catendar-wails/core/ai_sync/agent"
	synccli "catendar-wails/core/ai_sync/cli"
	"catendar-wails/core/calendar"
	"catendar-wails/core/email"
)

const (
	StateIdle      = "idle"
	StateRunning   = "running"
	StateSucceeded = "succeeded"
	StateFailed    = "failed"
	StateCancelled = "cancelled"

	maxRecentEmail = 200
	syncTimeout    = 20 * time.Minute
)

type Status struct {
	RunID      string         `json:"runId,omitempty"`
	State      string         `json:"state"`
	Stage      string         `json:"stage,omitempty"`
	Message    string         `json:"message,omitempty"`
	Provider   agent.Provider `json:"provider,omitempty"`
	EmailCount int            `json:"emailCount"`
	Created    int            `json:"created"`
	Skipped    int            `json:"skipped"`
	Rejected   int            `json:"rejected"`
	Conflicts  int            `json:"conflicts"`
	StartedAt  time.Time      `json:"startedAt,omitempty"`
	FinishedAt time.Time      `json:"finishedAt,omitempty"`
}

type Service struct {
	db             *sql.DB
	email          *email.Service
	calendar       *calendar.Service
	executablePath string
	runner         agent.Runner

	mu     sync.RWMutex
	status Status
	cancel context.CancelFunc
	done   chan struct{}
}

type runState struct {
	runID      string
	account    email.Account
	summaries  []email.MessageSummary
	available  map[string]email.MessageSummary
	batchCalls int
}

type emailListResult struct {
	Items      []email.MessageSummary `json:"items"`
	NextCursor *int                   `json:"nextCursor,omitempty"`
	Total      int                    `json:"total"`
}

func NewService(
	db *sql.DB,
	emailService *email.Service,
	calendarService *calendar.Service,
	executablePath string,
) *Service {
	return &Service{
		db:             db,
		email:          emailService,
		calendar:       calendarService,
		executablePath: executablePath,
		status:         Status{State: StateIdle},
	}
}

func (s *Service) Start(parent context.Context) (Status, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status.State == StateRunning {
		return s.status, errors.New("AI Sync is already running")
	}
	if s.executablePath == "" {
		return s.status, errors.New("CATendar executable path is unavailable")
	}

	runID, err := newRunID()
	if err != nil {
		return s.status, err
	}
	ctx, cancel := context.WithTimeout(parent, syncTimeout)
	s.cancel = cancel
	s.done = make(chan struct{})
	s.status = Status{
		RunID:     runID,
		State:     StateRunning,
		Stage:     "preparing",
		Message:   "Preparing AI Sync…",
		StartedAt: time.Now(),
	}
	go s.run(ctx, cancel, runID, s.done)
	return s.status, nil
}

func (s *Service) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *Service) Cancel() bool {
	s.mu.RLock()
	cancel := s.cancel
	running := s.status.State == StateRunning
	s.mu.RUnlock()
	if running && cancel != nil {
		cancel()
		return true
	}
	return false
}

func (s *Service) Shutdown() {
	s.mu.RLock()
	cancel := s.cancel
	done := s.done
	s.mu.RUnlock()
	if cancel != nil {
		cancel()
	}
	if done != nil {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
		}
	}
}

func (s *Service) run(ctx context.Context, cancel context.CancelFunc, runID string, done chan struct{}) {
	defer cancel()
	defer close(done)
	settings, err := s.getSettings(ctx)
	if err != nil {
		s.finish(runID, StateFailed, "Could not read AI settings", err)
		return
	}
	provider, err := agent.ValidateProvider(string(settings.Provider))
	if err != nil {
		s.finish(runID, StateFailed, "AI provider is not configured", err)
		return
	}
	adapter, _ := agent.ForProvider(provider)
	availability := agent.Resolve(adapter, configuredPath(settings, provider))
	if !availability.Installed {
		s.finish(runID, StateFailed, "Selected AI CLI is not available", errors.New(availability.Error))
		return
	}
	s.update(runID, func(status *Status) {
		status.Provider = provider
		status.Stage = "email"
		status.Message = "Reading recent email…"
	})

	account, summaries, err := s.email.ListRecent(ctx, time.Now().AddDate(0, 0, -30), maxRecentEmail)
	if err != nil {
		s.finish(runID, StateFailed, "Could not read recent email", err)
		return
	}
	if len(summaries) == 0 {
		s.finish(runID, StateSucceeded, "No email found in the last 30 days", nil)
		return
	}
	s.update(runID, func(status *Status) {
		status.EmailCount = len(summaries)
		status.Stage = "agent"
		status.Message = fmt.Sprintf("%s is planning from %d emails…", providerLabel(provider), len(summaries))
	})

	session, err := newLocalSession(s.executablePath)
	if err != nil {
		s.finish(runID, StateFailed, "Could not prepare the local CLI", err)
		return
	}
	defer session.Close()
	state := &runState{
		runID:     runID,
		account:   account,
		summaries: summaries,
		available: make(map[string]email.MessageSummary, len(summaries)),
	}
	for _, summary := range summaries {
		state.available[summary.ID] = summary
	}

	serverContext, stopServer := context.WithCancel(ctx)
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- session.Serve(serverContext, func(requestContext context.Context, request synccli.Request) (any, error) {
			return s.handleCLIRequest(requestContext, state, request)
		})
	}()

	prompt := agent.Prompt(session.CommandPath(), time.Now())
	err = s.runner.Run(ctx, adapter, agent.RunRequest{
		BinaryPath:  availability.Path,
		WorkDir:     session.WorkDir(),
		Prompt:      prompt,
		CLICommand:  session.CommandPath(),
		Environment: session.Environment(),
	}, func() {
		s.update(runID, func(status *Status) {
			if status.Stage == "agent" {
				status.Message = providerLabel(provider) + " is analyzing email…"
			}
		})
	})
	stopServer()
	serverErr := <-serverDone
	if serverErr != nil && !errors.Is(serverErr, context.Canceled) && err == nil {
		err = serverErr
	}
	if err == nil && state.batchCalls != 1 {
		err = fmt.Errorf("agent submitted the calendar batch %d times; expected exactly once", state.batchCalls)
	}
	if err != nil {
		if agent.IsCancellation(err) {
			s.finish(runID, StateCancelled, "AI Sync cancelled", nil)
		} else {
			s.finish(runID, StateFailed, providerLabel(provider)+" could not finish AI Sync", err)
		}
		return
	}
	s.finish(runID, StateSucceeded, "AI Sync complete", nil)
}

func (s *Service) handleCLIRequest(ctx context.Context, state *runState, request synccli.Request) (any, error) {
	switch request.Command {
	case "email.list":
		var input synccli.EmailListRequest
		if err := json.Unmarshal(request.Payload, &input); err != nil {
			return nil, errors.New("invalid email list request")
		}
		if input.Cursor < 0 || input.Cursor > len(state.summaries) {
			return nil, errors.New("email cursor is out of range")
		}
		if input.Limit < 1 || input.Limit > 100 {
			input.Limit = 50
		}
		end := input.Cursor + input.Limit
		if end > len(state.summaries) {
			end = len(state.summaries)
		}
		result := emailListResult{
			Items: state.summaries[input.Cursor:end],
			Total: len(state.summaries),
		}
		if end < len(state.summaries) {
			result.NextCursor = &end
		}
		return result, nil

	case "email.read":
		var input synccli.EmailReadRequest
		if err := json.Unmarshal(request.Payload, &input); err != nil || len(input.IDs) == 0 {
			return nil, errors.New("invalid email read request")
		}
		if len(input.IDs) > 20 {
			return nil, errors.New("read at most 20 emails per request")
		}
		for _, id := range input.IDs {
			if _, ok := state.available[id]; !ok {
				return nil, fmt.Errorf("email %q is outside this AI Sync run", id)
			}
		}
		s.update(state.runID, func(status *Status) {
			status.Message = fmt.Sprintf("Reading %d selected emails…", len(input.IDs))
		})
		messages, err := s.email.ReadMessagesForAccount(ctx, state.account, input.IDs)
		if err != nil {
			return nil, err
		}
		return messages, nil

	case "calendar.list":
		var input synccli.CalendarListRequest
		if err := json.Unmarshal(request.Payload, &input); err != nil {
			return nil, errors.New("invalid calendar list request")
		}
		if err := validateDateRange(input.From, input.To); err != nil {
			return nil, err
		}
		return s.calendar.GetEventsByDateRange(ctx, input.From, input.To)

	case "calendar.batch-upsert":
		state.batchCalls++
		if state.batchCalls > 1 {
			return nil, errors.New("calendar batch-upsert may only be called once per AI Sync run")
		}
		var input synccli.CalendarBatchRequest
		if err := json.Unmarshal(request.Payload, &input); err != nil {
			return nil, errors.New("invalid calendar batch request")
		}
		sources := make(map[string]calendar.EmailSource, len(state.available))
		for id, summary := range state.available {
			sources[id] = calendar.EmailSource{
				AccountID: state.account.ID,
				MessageID: summary.MessageID,
				Subject:   summary.Subject,
				URL:       summary.SourceURL,
			}
		}
		s.update(state.runID, func(status *Status) {
			status.Stage = "calendar"
			status.Message = "Writing calendar plans…"
		})
		result, err := s.calendar.BatchCreateFromEmail(ctx, input.Events, sources)
		if err != nil {
			return nil, err
		}
		s.update(state.runID, func(status *Status) {
			status.Created += result.Created
			status.Skipped += result.Skipped
			status.Rejected += result.Rejected
			status.Conflicts += result.Conflicts
		})
		return result, nil
	default:
		return nil, errors.New("CLI command is not allowed in this AI Sync session")
	}
}

func (s *Service) update(runID string, update func(*Status)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status.RunID == runID && s.status.State == StateRunning {
		update(&s.status)
	}
}

func (s *Service) finish(runID, state, message string, runErr error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status.RunID != runID {
		return
	}
	s.status.State = state
	s.status.Stage = "done"
	s.status.Message = message
	if runErr != nil {
		detail := strings.TrimSpace(runErr.Error())
		if len(detail) > 300 {
			detail = detail[:300]
		}
		if detail != "" {
			s.status.Message += ": " + detail
		}
	}
	s.status.FinishedAt = time.Now()
	s.cancel = nil
}

func validateDateRange(from, to string) error {
	start, err := time.Parse("2006-01-02", from)
	if err != nil {
		return errors.New("calendar from must use YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", to)
	if err != nil {
		return errors.New("calendar to must use YYYY-MM-DD")
	}
	if end.Before(start) || end.Sub(start) > 730*24*time.Hour {
		return errors.New("calendar range is invalid or exceeds two years")
	}
	return nil
}

func newRunID() (string, error) {
	value := make([]byte, 8)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("create AI Sync id: %w", err)
	}
	return hex.EncodeToString(value), nil
}

func providerLabel(provider agent.Provider) string {
	if provider == agent.ProviderClaude {
		return "Claude Code"
	}
	return "Codex"
}
