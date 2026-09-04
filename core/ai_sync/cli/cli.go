package synccli

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"catendar-wails/core/calendar"
)

const (
	EnvSessionDir   = "CATENDAR_AI_SYNC_DIR"
	EnvSessionToken = "CATENDAR_AI_SYNC_TOKEN"
	maxRequestSize  = 1024 * 1024
)

type Request struct {
	ID      string          `json:"id"`
	Token   string          `json:"token"`
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Response struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

type EmailListRequest struct {
	Cursor int `json:"cursor"`
	Limit  int `json:"limit"`
}

type EmailReadRequest struct {
	IDs []string `json:"ids"`
}

type CalendarListRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type CalendarBatchRequest struct {
	Events []calendar.EmailEventInput `json:"events"`
}

func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 2
	}

	var command string
	var payload any
	var err error
	switch args[0] + " " + args[1] {
	case "email list":
		command, payload, err = parseEmailList(args[2:], stderr)
	case "email read":
		command, payload, err = parseEmailRead(args[2:], stderr)
	case "calendar list":
		command, payload, err = parseCalendarList(args[2:], stderr)
	case "calendar batch-upsert":
		command, payload, err = parseCalendarBatch(args[2:], stdin, stderr)
	default:
		printUsage(stderr)
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	response, err := call(ctx, command, payload)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(response); err != nil {
		fmt.Fprintln(stderr, "write CLI response:", err)
		return 1
	}
	if !response.OK {
		return 1
	}
	return 0
}

func parseEmailList(args []string, stderr io.Writer) (string, EmailListRequest, error) {
	flags := flag.NewFlagSet("email list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	cursor := flags.Int("cursor", 0, "zero-based result cursor")
	limit := flags.Int("limit", 50, "number of email summaries to return")
	if err := flags.Parse(args); err != nil {
		return "", EmailListRequest{}, err
	}
	return "email.list", EmailListRequest{Cursor: *cursor, Limit: *limit}, nil
}

func parseEmailRead(args []string, stderr io.Writer) (string, EmailReadRequest, error) {
	flags := flag.NewFlagSet("email read", flag.ContinueOnError)
	flags.SetOutput(stderr)
	values := flags.String("ids", "", "comma-separated email IDs from email list")
	if err := flags.Parse(args); err != nil {
		return "", EmailReadRequest{}, err
	}
	ids := splitValues(*values)
	if len(ids) == 0 {
		return "", EmailReadRequest{}, errors.New("--ids is required")
	}
	return "email.read", EmailReadRequest{IDs: ids}, nil
}

func parseCalendarList(args []string, stderr io.Writer) (string, CalendarListRequest, error) {
	flags := flag.NewFlagSet("calendar list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	from := flags.String("from", "", "start date in YYYY-MM-DD format")
	to := flags.String("to", "", "end date in YYYY-MM-DD format")
	if err := flags.Parse(args); err != nil {
		return "", CalendarListRequest{}, err
	}
	if *from == "" || *to == "" {
		return "", CalendarListRequest{}, errors.New("--from and --to are required")
	}
	return "calendar.list", CalendarListRequest{From: *from, To: *to}, nil
}

func parseCalendarBatch(args []string, stdin io.Reader, stderr io.Writer) (string, CalendarBatchRequest, error) {
	flags := flag.NewFlagSet("calendar batch-upsert", flag.ContinueOnError)
	flags.SetOutput(stderr)
	inputPath := flags.String("input", "", "JSON file containing an event array; use - for stdin")
	if err := flags.Parse(args); err != nil {
		return "", CalendarBatchRequest{}, err
	}
	if *inputPath == "" {
		return "", CalendarBatchRequest{}, errors.New("--input is required")
	}

	var reader io.Reader = stdin
	if *inputPath != "-" {
		path, err := sessionFile(*inputPath)
		if err != nil {
			return "", CalendarBatchRequest{}, err
		}
		file, err := os.Open(path)
		if err != nil {
			return "", CalendarBatchRequest{}, fmt.Errorf("open event input: %w", err)
		}
		defer file.Close()
		reader = file
	}

	var events []calendar.EmailEventInput
	decoder := json.NewDecoder(io.LimitReader(reader, maxRequestSize))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&events); err != nil {
		return "", CalendarBatchRequest{}, fmt.Errorf("decode event input: %w", err)
	}
	return "calendar.batch-upsert", CalendarBatchRequest{Events: events}, nil
}

func call(ctx context.Context, command string, payload any) (Response, error) {
	directory := os.Getenv(EnvSessionDir)
	token := os.Getenv(EnvSessionToken)
	if directory == "" || token == "" {
		return Response{}, errors.New("CATendar AI Sync session is not available")
	}

	id, err := randomID()
	if err != nil {
		return Response{}, err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("encode CLI request: %w", err)
	}
	request := Request{ID: id, Token: token, Command: command, Payload: payloadJSON}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return Response{}, fmt.Errorf("encode CLI request: %w", err)
	}

	requestPath := filepath.Join(directory, "requests", id+".json")
	temporaryPath := requestPath + ".tmp"
	if err := os.WriteFile(temporaryPath, requestJSON, 0o600); err != nil {
		return Response{}, fmt.Errorf("write CLI request: %w", err)
	}
	if err := os.Rename(temporaryPath, requestPath); err != nil {
		os.Remove(temporaryPath)
		return Response{}, fmt.Errorf("publish CLI request: %w", err)
	}

	responsePath := filepath.Join(directory, "responses", id+".json")
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(90 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		case <-deadline.C:
			return Response{}, errors.New("CATendar Core did not answer the CLI request")
		case <-ticker.C:
			data, err := os.ReadFile(responsePath)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return Response{}, fmt.Errorf("read CLI response: %w", err)
			}
			os.Remove(responsePath)
			var response Response
			if err := json.Unmarshal(data, &response); err != nil {
				return Response{}, fmt.Errorf("decode CLI response: %w", err)
			}
			return response, nil
		}
	}
}

func sessionFile(value string) (string, error) {
	directory := os.Getenv(EnvSessionDir)
	if directory == "" {
		return "", errors.New("CATendar AI Sync session is not available")
	}
	path := value
	if !filepath.IsAbs(path) {
		path = filepath.Join(directory, path)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve event input: %w", err)
	}
	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("resolve AI Sync directory: %w", err)
	}
	relative, err := filepath.Rel(absDirectory, absPath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("event input must be inside the AI Sync working directory")
	}
	return absPath, nil
}

func splitValues(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func randomID() (string, error) {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		return "", fmt.Errorf("create request id: %w", err)
	}
	return hex.EncodeToString(data), nil
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `CATendar AI Sync CLI

Commands:
  CATendar cli email list [--cursor N] [--limit N]
  CATendar cli email read --ids ID,ID
  CATendar cli calendar list --from YYYY-MM-DD --to YYYY-MM-DD
  CATendar cli calendar batch-upsert --input events.json`)
}
